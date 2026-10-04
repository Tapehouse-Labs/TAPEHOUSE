// SPDX-License-Identifier: MIT OR Apache-2.0

package calibrate

import (
	"cmp"
	"math"
	"math/rand/v2"
	"slices"

	"github.com/tapehouse/tapehouse/services/internal/backtest"
)

// The tail's fit.
const (
	tailQuantile    = 0.95
	bootstraps      = 20_000
	scaleLow        = 1e-4
	scaleHigh       = 10.0
	scaleIterations = 100
	redeployShare   = 0.10
	lowQuantile     = 0.05
)

// Tail is the gap cover's tail refitted: the threshold u at the 95th percentile of the weekend falls in gaps, the
// share of falls above it, the scale of a generalised Pareto of shape ½ fitted to the excesses, the 90% interval of
// the reopen ratio bootstrapped over weeks, and the three constants the cover takes, in millionths of a gap, beside
// the deployed ones.
type Tail struct {
	Threshold   float64
	Scale       float64
	Probability float64
	Excesses    int
	ReopenLow   float64
	ReopenHigh  float64
	Scaled      [3]uint32
	Deployed    [3]uint32
	Redeploy    bool
}

// Falls are the weekend falls of the market in the Years up to its last session, each divided by its asset's gap in
// millionths: the fall from a close to the open that follows a closure of three days less an hour or more, pooled over
// the assets; rises are left out.
func Falls(m backtest.Market, gaps []uint32) []float64 {
	var out []float64
	from := m.At[len(m.At)-1].AddDate(-Years, 0, 0)
	for _, c := range m.Closures() {
		if c.Open.Sub(c.Close) < gapSpan || c.Close.Before(from) {
			continue
		}
		for k, move := range c.Moves {
			if move < 0 && gaps[k] > 0 {
				out = append(out, -move/(float64(gaps[k])/1e6))
			}
		}
	}
	return out
}

// quantile is the nearest-rank q-quantile of sorted.
func quantile(sorted []float64, q float64) float64 {
	return sorted[max(0, int(math.Ceil(q*float64(len(sorted))))-1)]
}

// paretoScale is the scale that maximises the log-likelihood of excesses under a generalised Pareto of shape ½:
// -n·ln β - 3·Σ ln(1 + x/(2β)).
func paretoScale(excesses []float64) float64 {
	return golden(func(beta float64) float64 {
		l := -float64(len(excesses)) * math.Log(beta)
		for _, x := range excesses {
			l -= 3 * math.Log1p(x/(2*beta))
		}
		return l
	}, scaleLow, scaleHigh, scaleIterations)
}

// FitTail fits falls, in gaps, and bootstraps the reopen ratios, by week, with the seed.
func FitTail(falls []float64, ratios [][]float64, seed uint64) Tail {
	var t Tail
	if len(falls) > 0 {
		sorted := slices.Clone(falls)
		slices.SortFunc(sorted, cmp.Compare[float64])
		t.Threshold = quantile(sorted, tailQuantile)
		var excesses []float64
		for _, f := range sorted {
			if f > t.Threshold {
				excesses = append(excesses, f-t.Threshold)
			}
		}
		t.Excesses = len(excesses)
		t.Probability = float64(len(excesses)) / float64(len(sorted))
		if len(excesses) > 0 {
			t.Scale = paretoScale(excesses)
		}
	}
	if len(ratios) > 0 {
		rng := rand.New(rand.NewPCG(seed, seed))
		means := make([]float64, bootstraps)
		for i := range means {
			var sum float64
			var count int
			for range ratios {
				week := ratios[rng.IntN(len(ratios))]
				for _, r := range week {
					sum += r
				}
				count += len(week)
			}
			if count > 0 {
				means[i] = sum / float64(count)
			}
		}
		slices.Sort(means)
		t.ReopenLow, t.ReopenHigh = quantile(means, lowQuantile), quantile(means, tailQuantile)
		t.Scaled = t.Constants()
	}
	return t
}

// Constants are the cover's TAIL_THRESHOLD_PPM, TAIL_SCALE_PPM and TAIL_PROBABILITY_PPM the fit gives: the threshold
// and the scale, in millionths of a gap, scaled by the upper end of the reopen ratio's interval, and the probability.
func (t Tail) Constants() [3]uint32 {
	return [3]uint32{millionths(t.Threshold * t.ReopenHigh), millionths(t.Scale * t.ReopenHigh), millionths(t.Probability)}
}

// Against sets the deployed constants and, once a reopen ratio is known, flags a redeployment where any of the fitted ones differs from its by more
// than a tenth.
func (t Tail) Against(deployed [3]uint32) Tail {
	t.Deployed = deployed
	t.Redeploy = false
	if t.ReopenHigh == 0 {
		return t
	}
	for i, d := range deployed {
		if math.Abs(float64(t.Scaled[i])-float64(d)) > redeployShare*float64(d) {
			t.Redeploy = true
		}
	}
	return t
}
