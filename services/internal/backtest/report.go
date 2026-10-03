// SPDX-License-Identifier: MIT OR Apache-2.0

package backtest

import (
	"time"

	"github.com/tapehouse/tapehouse/services/internal/margin"
)

// Mega5 is the basket MEGA5's components.
var Mega5 = []string{"NVDA", "TSLA", "AAPL", "MSFT", "GOOGL"}

// Seeds and weights of the portfolios: long-only from 0 to 3 a unit, and mixed from −3 to 3.
const (
	LongSeed  = 7
	MixedSeed = 11
	Count     = 200
)

// History is the record the backtest ran on.
type History struct {
	Assets   []string `json:"assets"`
	From     string   `json:"from"`
	To       string   `json:"to"`
	Sessions int      `json:"sessions"`
	Closures int      `json:"closures"`
	Launch   string   `json:"launch"`
	Since    int      `json:"closuresSinceLaunch"`
}

// Accuracy sums the band's closures on one chain: how many there were, how many a public 24/7 value reaches, how
// many reopening rounds fell inside the band's last value, the median error and how long before the reopening
// round that value was.
type Accuracy struct {
	Chain       uint64  `json:"chain"`
	Reopens     int     `json:"reopens"`
	Evaluated   int     `json:"evaluated"`
	Inside      int     `json:"inside"`
	MedianError float64 `json:"medianErrorBps"`
	MedianLeadS float64 `json:"medianLeadS"`
}

// Report is everything the backtest measures.
type Report struct {
	History      History             `json:"history"`
	Coverage     map[string]Coverage `json:"coverage"`
	Horizons     []Horizon           `json:"horizons"`
	Leverage     Leverage            `json:"leverage"`
	Deleveraging Deleveraging        `json:"deleveraging"`
	March2020    *Scenario           `json:"march2020,omitempty"`
	Capacity     Capacity            `json:"capacity"`
	Backstop     Backstop            `json:"backstop"`
	Liquidator   Liquidator          `json:"liquidator"`
	Reviews      Reviews             `json:"reviews"`
	Basket       Basket              `json:"basket"`
	GapCover     map[string]GapCover `json:"gapCover"`
	Reopens      []Reopen            `json:"reopens"`
	Impacts      []Impact            `json:"impacts"`
	Accuracy     []Accuracy          `json:"accuracy"`
}

// Run measures the engine and the constants over the market. depths holds each asset's selling then buying depth in
// USD, debtCap is the accounts' debt cap in USD, blockSeconds the chain's block interval, and launch when the chain
// launched.
func Run(m Market, set *margin.Set, depths []uint32, debtCap, blockSeconds float64, launch time.Time) Report {
	closures := m.Closures()
	books := map[string][]Book{
		"long":  Books(set, Portfolios(LongSeed, Count, len(m.Names), 0, 3)),
		"mixed": Books(set, Portfolios(MixedSeed, Count, len(m.Names), -3, 3)),
	}
	r := Report{
		History: History{
			Assets: m.Names, From: date(m.At[0]), To: date(m.At[len(m.At)-1]), Sessions: len(m.At),
			Closures: len(closures), Launch: date(launch),
		},
		Coverage:     map[string]Coverage{"long": Cover(books["long"], closures), "mixed": Cover(books["mixed"], closures)},
		Horizons:     Horizons(m.Names, closures, m.Returns()),
		Leverage:     Recalibrate(m.Names, closures, books),
		Deleveraging: Delever(books),
		Capacity:     Compare(closures, books),
		Liquidator:   ReviewLiquidator(m.Names, m.Returns(), blockSeconds),
		Reviews:      Review(m, closures, depths),
		Basket:       Rebalance(m, Mega5),
		GapCover: map[string]GapCover{
			"history": ReviewGapCover(closures, set.Parameters.Gaps, m.At[0]),
			"chain":   ReviewGapCover(closures, set.Parameters.Gaps, launch),
		},
	}
	for _, c := range closures {
		if !c.Open.Before(launch) {
			r.History.Since++
		}
	}
	if s, ok := Replay(m.Names, closures, books, "2020-03-16"); ok {
		r.March2020 = &s
	}
	caps := make([]float64, len(m.Names))
	for k := range caps {
		caps[k] = float64(depths[2*k])
	}
	r.Backstop = Calibrate(set, m.Names, caps, debtCap, closures)
	return r
}

// Summarize sums the reopens by chain.
func Summarize(reopens []Reopen) []Accuracy {
	var out []Accuracy
	errors, leads := map[uint64][]float64{}, map[uint64][]float64{}
	index := map[uint64]int{}
	for _, r := range reopens {
		i, ok := index[r.Chain]
		if !ok {
			i = len(out)
			index[r.Chain] = i
			out = append(out, Accuracy{Chain: r.Chain})
		}
		out[i].Reopens++
		if !r.Evaluated {
			continue
		}
		out[i].Evaluated++
		leads[r.Chain] = append(leads[r.Chain], float64(r.LeadS))
		if r.Inside {
			out[i].Inside++
		}
		if r.State != 0 {
			errors[r.Chain] = append(errors[r.Chain], r.ErrorBps)
		}
	}
	for i := range out {
		if e := errors[out[i].Chain]; len(e) > 0 {
			out[i].MedianError = median(e)
		}
		if l := leads[out[i].Chain]; len(l) > 0 {
			out[i].MedianLeadS = median(l)
		}
	}
	return out
}
