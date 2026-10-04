// SPDX-License-Identifier: MIT OR Apache-2.0

package calibrate

import (
	"math"
	"math/big"
	"slices"

	"github.com/tapehouse/tapehouse/services/internal/band"
	"github.com/tapehouse/tapehouse/services/internal/history"
)

// The search of Fit: the bounds of log q and log r, and the iterations of each golden-section search.
const (
	logQLow, logQHigh = -20.0, 0.0
	logRLow, logRHigh = -30.0, -5.0
	searches          = 60
	// bandCoverage is the share of next values the band must hold, and widthRatio how much wider than the reference
	// its median half-width may be, before a redeployment is flagged.
	bandCoverage = 0.99
	widthRatio   = 2.0
	warmup       = 50
)

// Filter is a local-level state-space model of a log price, x(t) = x(t-1) + w with w ~ N(0, Q·dt) in days, observed
// as y = x + v with v ~ N(0, R). It computes in float64: nothing it reports decides money.
type Filter struct {
	Q, R float64
}

// Run filters ys, observed at times in days, started diffuse at the first observation. It returns the mean and the
// variance of each observation's prediction before it is seen, the first being the observation itself with no
// variance, and the log-likelihood of the innovations after the first.
func (f Filter) Run(times, ys []float64) (means, variances []float64, logLikelihood float64) {
	means, variances = make([]float64, len(ys)), make([]float64, len(ys))
	if len(ys) == 0 {
		return means, variances, 0
	}
	x, p := ys[0], f.R
	means[0] = ys[0]
	for i := 1; i < len(ys); i++ {
		p += f.Q * (times[i] - times[i-1])
		means[i], variances[i] = x, p
		s := p + f.R
		e := ys[i] - x
		logLikelihood -= 0.5 * (math.Log(2*math.Pi*s) + e*e/s)
		k := p / s
		x += k * e
		p *= 1 - k
	}
	return means, variances, logLikelihood
}

// golden maximises f over [low, high] by iterations golden-section steps and returns the maximiser.
func golden(f func(float64) float64, low, high float64, iterations int) float64 {
	phi := (math.Sqrt(5) - 1) / 2
	a, b := low, high
	c, d := b-phi*(b-a), a+phi*(b-a)
	fc, fd := f(c), f(d)
	for range iterations {
		if fc > fd {
			b, d, fd = d, c, fc
			c = b - phi*(b-a)
			fc = f(c)
		} else {
			a, c, fc = c, d, fd
			d = a + phi*(b-a)
			fd = f(d)
		}
	}
	return (a + b) / 2
}

// Fit is the filter that maximises the likelihood of ys: a golden-section search on log r, the outer one, around a
// golden-section search on log q, each over bounded intervals.
func Fit(times, ys []float64) Filter {
	best := func(logR float64) (float64, float64) {
		r := math.Exp(logR)
		logQ := golden(func(logQ float64) float64 {
			_, _, l := Filter{math.Exp(logQ), r}.Run(times, ys)
			return l
		}, logQLow, logQHigh, searches)
		_, _, l := Filter{math.Exp(logQ), r}.Run(times, ys)
		return logQ, l
	}
	logR := golden(func(logR float64) float64 {
		_, l := best(logR)
		return l
	}, logRLow, logRHigh, searches)
	logQ, _ := best(logR)
	return Filter{math.Exp(logQ), math.Exp(logR)}
}

// ReferenceVolatility is the fitted filter's daily volatility, sqrt(q) in centi-basis-points, over the last Window
// sessions of the adjusted closes.
func ReferenceVolatility(adjCloses []float64) uint32 {
	closes := adjCloses[max(0, len(adjCloses)-Window-1):]
	times, ys := make([]float64, len(closes)), make([]float64, len(closes))
	for i, c := range closes {
		times[i], ys[i] = float64(i), math.Log(c)
	}
	return millionths(math.Sqrt(Fit(times, ys).Q))
}

// Distance is how far the band replayed from a 24/7 history is from the filter's reference band: at every value after
// the first warmup, the distance of the two centres in basis points, the ratio of the half-widths, and whether each
// band held the next value. Redeploy is set where the band's coverage falls below 99% or its median half-width
// exceeds twice the reference's.
type Distance struct {
	Asset             string
	Points            int
	MedianMidBps      float64
	P95MidBps         float64
	MedianHalfRatio   float64
	BandCoverage      float64
	ReferenceCoverage float64
	Redeploy          bool
}

func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	return sorted[min(len(sorted)-1, int(p*float64(len(sorted))))]
}

// Measure replays leg through the band's record, as the band stores the values of its 24/7 feed, and fits the filter
// to its logarithms, with time in days. multiplier is the Stock Token's multiplier with 18 decimals. An asset the band
// prices through an index leg has no 24/7 feed to replay, and measures nothing.
func Measure(asset band.Asset, leg []history.Point, multiplier *big.Int) Distance {
	d := Distance{Asset: asset.Symbol, Points: len(leg)}
	if asset.Feed == "" || len(leg) <= warmup+1 {
		return d
	}
	times, ys := make([]float64, len(leg)), make([]float64, len(leg))
	for i, p := range leg {
		times[i], ys[i] = float64(p.Ms-leg[0].Ms)/86_400_000, math.Log(p.Value)
	}
	filter := Fit(times, ys)
	means, variances, _ := filter.Run(times, ys)
	scale, _ := new(big.Float).Quo(new(big.Float).SetInt(multiplier), big.NewFloat(1e18)).Float64()
	z := float64(band.ZX10) / 10
	record := band.NewRecord(asset.Feed)
	var mids, ratios []float64
	var inside, referenceInside, next int
	for i, p := range leg {
		record.Write(band.Written{Feed: asset.Feed, Value: p.Scaled(), Ms: p.Ms})
		if i < warmup || i+1 >= len(leg) {
			continue
		}
		q := record.Quote(asset, band.Chainlink{}, multiplier, true, (p.Ms+999)/1000)
		if q.State == band.Halted {
			continue
		}
		next++
		value := leg[i+1].Value * 1e8 * scale
		if float64(q.Low) <= value && value <= float64(q.Mid)*(1+float64(q.HalfBps)/1e4) {
			inside++
		}
		centre := math.Exp(means[i+1])
		half := z * math.Sqrt(variances[i+1]+filter.Q*float64(band.LatencyMin)/1440)
		if centre*math.Exp(-half) <= leg[i+1].Value && leg[i+1].Value <= centre*math.Exp(half) {
			referenceInside++
		}
		mids = append(mids, math.Abs(float64(q.Mid)/(centre*1e8*scale)-1)*1e4)
		ratios = append(ratios, float64(q.HalfBps)/(half*1e4))
	}
	if next == 0 {
		return d
	}
	slices.Sort(mids)
	slices.Sort(ratios)
	d.MedianMidBps, d.P95MidBps = percentile(mids, 0.5), percentile(mids, 0.95)
	d.MedianHalfRatio = percentile(ratios, 0.5)
	d.BandCoverage, d.ReferenceCoverage = float64(inside)/float64(next), float64(referenceInside)/float64(next)
	d.Redeploy = d.BandCoverage < bandCoverage || d.MedianHalfRatio > widthRatio
	return d
}
