// SPDX-License-Identifier: MIT OR Apache-2.0

// Package calibrate estimates the parameters the margin engine's setParameters takes, by the method that set its
// launch values, and proposes them within what the engine accepts.
//
// Over the latest ten years of adjusted daily closes, an asset's volatility is the sample standard deviation of its
// log returns, floored by the decade's figure and set by the larger of that and the last 252 sessions'. A pair's
// correlation is Pearson's, floored by 0.75 of the decade's and 0.25 of the stress window's, the 60 sessions on which
// the market fell furthest. A weekend gap is the mean of the largest 1% of absolute moves from a close to the open
// that follows a closure of three days less an hour or more, capped at 100%. A depth is what a pool pays out or takes
// in walking its initialized ticks to a 10% move, the lowest of recent snapshots.
package calibrate

import (
	"cmp"
	"fmt"
	"math"
	"slices"
	"sort"
	"time"

	"github.com/tapehouse/tapehouse/services/internal/backtest"
	"github.com/tapehouse/tapehouse/services/internal/history"
)

// The windows of the launch method.
const (
	// Window is the sessions of the recent figure.
	Window = 252
	// StressSessions is the sessions of the stress window.
	StressSessions = 60
	// Years is the years of history behind the floors, counted back from the third day after the last session: the
	// launch figures came from bars pulled that day, the Monday after the last Friday.
	Years = 10
	// StressMarket is the asset whose worst sessions are the stress window, when the engine's market is not an asset.
	StressMarket = "SPY"

	pullDays      = 3
	maxMillionths = 1_000_000
	maxBps        = 10_000
	gapSpan       = 3*24*time.Hour - time.Hour
)

// Targets are the parameters the method estimates, in the engine's formats and order.
type Targets struct {
	Volatilities []uint32
	Gaps         []uint32
	Correlations []uint16
}

// millionths is x in millionths, rounded half to even.
func millionths(x float64) uint32 {
	return uint32(math.RoundToEven(x * 1e6))
}

func basisPoints(x float64) uint16 {
	return uint16(max(0, min(maxBps, math.RoundToEven(x*1e4))))
}

// Returns are the log returns of consecutive adjusted closes of the sessions every series traded in the Years up to
// last, and the dates they end on, indexed [asset][session].
func Returns(series []history.Series, last time.Time) (dates []time.Time, returns [][]float64) {
	m := backtest.Align(series)
	from := last.AddDate(-Years, 0, pullDays)
	first := sort.Search(len(m.At), func(i int) bool { return m.At[i].After(from) })
	if first > 0 {
		first--
	}
	returns = make([][]float64, len(series))
	for i := first + 1; i < len(m.At) && !m.At[i].After(last); i++ {
		dates = append(dates, m.At[i])
		for k := range series {
			returns[k] = append(returns[k], math.Log(m.Adj[i][k]/m.Adj[i-1][k]))
		}
	}
	return dates, returns
}

func deviation(xs []float64) float64 {
	var sum float64
	for _, x := range xs {
		sum += x
	}
	mean := sum / float64(len(xs))
	var squares float64
	for _, x := range xs {
		squares += (x - mean) * (x - mean)
	}
	return math.Sqrt(squares / float64(len(xs)-1))
}

func recent(xs []float64) []float64 {
	return xs[max(0, len(xs)-Window):]
}

// Volatilities are each asset's daily volatility in centi-basis-points: the floor over the whole window and the
// figure over the last Window sessions.
func Volatilities(returns [][]float64) (floors, years []uint32) {
	for _, r := range returns {
		floors = append(floors, millionths(deviation(r)))
		years = append(years, millionths(deviation(recent(r))))
	}
	return floors, years
}

// StressStart is the first of the StressSessions consecutive returns of market with the lowest sum.
func StressStart(market []float64) int {
	best, lowest := 0, math.Inf(1)
	var sum float64
	for i, r := range market {
		sum += r
		if i >= StressSessions {
			sum -= market[i-StressSessions]
		}
		if i >= StressSessions-1 && sum < lowest {
			best, lowest = i-StressSessions+1, sum
		}
	}
	return best
}

func pearson(a, b []float64) float64 {
	var sa, sb float64
	for i := range a {
		sa += a[i]
		sb += b[i]
	}
	ma, mb := sa/float64(len(a)), sb/float64(len(b))
	var ab, aa, bb float64
	for i := range a {
		ab += (a[i] - ma) * (b[i] - mb)
		aa += (a[i] - ma) * (a[i] - ma)
		bb += (b[i] - mb) * (b[i] - mb)
	}
	return ab / math.Sqrt(aa*bb)
}

// Correlations are each pair's correlation in basis points, the upper triangle row by row: the floor, 0.75 of the
// whole window's and 0.25 of the stress window's, and the figure over the last Window sessions. The stress window is
// the market asset's worst; market is its position, or -1 for the equal-weighted portfolio.
func Correlations(returns [][]float64, market int) (floors, years []uint16) {
	if market >= 0 {
		return correlations(returns, returns[market])
	}
	mean := make([]float64, len(returns[0]))
	for i := range mean {
		for k := range returns {
			mean[i] += returns[k][i] / float64(len(returns))
		}
	}
	return correlations(returns, mean)
}

func correlations(returns [][]float64, stressMarket []float64) (floors, years []uint16) {
	n := len(returns)
	start := StressStart(stressMarket)
	for i := range n {
		for j := i + 1; j < n; j++ {
			full := pearson(returns[i], returns[j])
			stress := pearson(returns[i][start:start+StressSessions], returns[j][start:start+StressSessions])
			floors = append(floors, basisPoints(0.75*full+0.25*stress))
			years = append(years, basisPoints(pearson(recent(returns[i]), recent(returns[j]))))
		}
	}
	return floors, years
}

// Gaps are each asset's weekend gap in millionths: the mean of the largest 1% of the absolute moves from a close to
// the next open over closures of at least three days less an hour, in the Years up to last, capped at 100%.
func Gaps(series []history.Series, last time.Time) []uint32 {
	from := last.AddDate(-Years, 0, pullDays)
	out := make([]uint32, len(series))
	for k, s := range series {
		var moves []float64
		for i := 1; i < len(s.Bars); i++ {
			b, p := s.Bars[i], s.Bars[i-1]
			if b.At.After(last) || p.At.Before(from) || b.At.Sub(p.At) < gapSpan {
				continue
			}
			moves = append(moves, math.Abs(b.Open*b.AdjClose/b.Close/p.AdjClose-1))
		}
		slices.SortFunc(moves, func(a, b float64) int { return cmp.Compare(b, a) })
		count := int(math.Ceil(float64(len(moves)) / 100))
		var sum float64
		for _, m := range moves[:count] {
			sum += m
		}
		if count > 0 {
			out[k] = min(maxMillionths, millionths(sum/float64(count)))
		}
	}
	return out
}

// Estimate is the method's target for each of names, the larger of the floor's method and the last Window sessions'
// figure for a volatility and a correlation, and the gap's. market is the market asset's position in names; where it
// is -1 the stress window is StressMarket's, if series has it, else the equal-weighted portfolio's.
func Estimate(series []history.Series, names []string, market int, last time.Time) (Targets, error) {
	ordered := make([]history.Series, len(names), len(names)+1)
	for i, name := range names {
		j := slices.IndexFunc(series, func(s history.Series) bool { return s.Symbol == name })
		if j < 0 {
			return Targets{}, fmt.Errorf("no history for %s", name)
		}
		ordered[i] = series[j]
	}
	stress := -1
	if market < 0 {
		if j := slices.IndexFunc(series, func(s history.Series) bool { return s.Symbol == StressMarket }); j >= 0 {
			stress = len(ordered)
			ordered = append(ordered, series[j])
		}
	} else {
		stress = market
	}
	_, all := Returns(ordered, last)
	returns := all[:len(names)]
	floors, years := Volatilities(returns)
	var correlationFloors, correlationYears []uint16
	switch {
	case stress >= 0:
		correlationFloors, correlationYears = correlations(returns, all[stress])
	default:
		correlationFloors, correlationYears = Correlations(returns, -1)
	}
	t := Targets{Gaps: Gaps(ordered[:len(names)], last)}
	for i := range floors {
		t.Volatilities = append(t.Volatilities, max(floors[i], years[i]))
	}
	for i := range correlationFloors {
		t.Correlations = append(t.Correlations, max(correlationFloors[i], correlationYears[i]))
	}
	return t, nil
}
