// SPDX-License-Identifier: MIT OR Apache-2.0

package calibrate

import (
	"math"
	"math/big"
	"math/rand/v2"
	"testing"

	"github.com/tapehouse/tapehouse/services/internal/band"
	"github.com/tapehouse/tapehouse/services/internal/history"
)

// walk is n observations of a local-level model: log prices that move by N(0, q) a day, seen through N(0, r).
func walk(n int, q, r float64, a, b uint64) (times, ys []float64) {
	rng := rand.New(rand.NewPCG(a, b))
	x := math.Log(100)
	for i := range n {
		x += math.Sqrt(q) * rng.NormFloat64()
		times = append(times, float64(i))
		ys = append(ys, x+math.Sqrt(r)*rng.NormFloat64())
	}
	return times, ys
}

func TestTheFilterRecoversARandomWalksVariance(t *testing.T) {
	times, ys := walk(10_000, 4e-4, 1e-8, 3, 5)
	if f := Fit(times, ys); math.Abs(f.Q/4e-4-1) > 0.05 {
		t.Fatalf("q %g", f.Q)
	}
}

func TestTheFilterRecoversItsNoise(t *testing.T) {
	times, ys := walk(10_000, 1e-4, 1e-4, 3, 5)
	if f := Fit(times, ys); math.Abs(f.R/1e-4-1) > 0.20 {
		t.Fatalf("r %g", f.R)
	}
}

func TestPredictionsWidenWithTime(t *testing.T) {
	f := Filter{Q: 4e-4, R: 1e-6}
	_, near, _ := f.Run([]float64{0, 1}, []float64{4.6, 4.61})
	_, far, _ := f.Run([]float64{0, 4}, []float64{4.6, 4.61})
	if got := far[1] - near[1]; math.Abs(got-3*f.Q) > 1e-15 {
		t.Fatalf("variance widened by %g, want %g", got, 3*f.Q)
	}
}

func TestTheReferenceVolatilityOfDailyClosesNearsTheSampleOne(t *testing.T) {
	series := syntheticMarket()
	closes := make([]float64, len(series[0].Bars))
	for i, b := range series[0].Bars {
		closes[i] = b.AdjClose
	}
	if got := ReferenceVolatility(closes); math.Abs(float64(got)/27165-1) > 0.05 {
		t.Fatalf("reference volatility %d", got)
	}
}

// leg is points 60 s apart from a random walk whose per-step standard deviation is sigma, with a shock of shock
// from every period-th step on where period is not zero.
func leg(n int, sigma, shock float64, period int) []history.Point {
	rng := rand.New(rand.NewPCG(9, 9))
	price := 100.0
	out := make([]history.Point, n)
	for i := range out {
		price *= 1 + sigma*rng.NormFloat64()
		if period > 0 && i%period == period-1 {
			price *= 1 + shock
		}
		out[i] = history.Point{Ms: 1_790_000_000_000 + uint64(i)*60_000, Value: price}
	}
	return out
}

var distanceAsset = band.Asset{Symbol: "NVDA", Feed: "NVDA---24_7"}

func TestADistanceIsMeasuredOnTheSameHistory(t *testing.T) {
	d := Measure(distanceAsset, leg(2000, 0.002, 0, 0), big.NewInt(1e18))
	if d.Points != 2000 || d.MedianMidBps >= 5 || d.BandCoverage < 0.99 || d.ReferenceCoverage < 0.99 || d.Redeploy {
		t.Fatalf("%+v", d)
	}
}

func TestABandFarNarrowerThanTheReferenceIsFlagged(t *testing.T) {
	d := Measure(distanceAsset, leg(2000, 0.0005, 0.01, 25), big.NewInt(1e18))
	if d.BandCoverage >= 0.99 || d.ReferenceCoverage < 0.99 || !d.Redeploy {
		t.Fatalf("%+v", d)
	}
}

func TestABandFarWiderThanTheReferenceIsFlagged(t *testing.T) {
	d := Measure(distanceAsset, leg(2000, 0.0001, 0, 0), big.NewInt(1e18))
	if d.MedianHalfRatio <= 2 || !d.Redeploy {
		t.Fatalf("%+v", d)
	}
}

func TestAnIndexLegHasNothingToMeasure(t *testing.T) {
	d := Measure(band.Asset{Symbol: "SPY", Index: "USA500.Y---24_7"}, leg(2000, 0.002, 0, 0), big.NewInt(1e18))
	if d.Points != 2000 || d.BandCoverage != 0 || d.Redeploy {
		t.Fatalf("%+v", d)
	}
}
