// SPDX-License-Identifier: MIT OR Apache-2.0

package backtest

import (
	"cmp"
	"math"
	"math/big"
	"slices"

	"github.com/tapehouse/tapehouse/services/internal/band"
	"github.com/tapehouse/tapehouse/services/internal/history"
)

// Reopening is what the band is replayed from at one closure on one chain: the asset as the chain's band prices
// it, when the closure began in milliseconds, its 24/7 or index history, the signed market status' history,
// Chainlink's last round before the reopen and the first after, the Stock Token's multiplier and whether the
// sequencer was up and settled.
type Reopening struct {
	Chain       uint64
	Asset       band.Asset
	Date        string
	SinceMs     uint64
	Leg         []history.Point
	Status      [3][]history.Point
	Before      history.Round
	After       history.Round
	Multiplier  *big.Int
	SequencerUp bool
}

// Reopen is the band's last value in a closure, quoted when the last public 24/7 value before the reopen was
// written, and whether the first Chainlink round after the reopen fell inside it: the closure's value compared with
// the opening price. Evaluated is false when no public 24/7 value lies inside the closure; LeadS is how long before
// the reopening round the band's value was.
type Reopen struct {
	Chain     uint64  `json:"chain"`
	Asset     string  `json:"asset"`
	Date      string  `json:"date"`
	Evaluated bool    `json:"evaluated"`
	LeadS     uint64  `json:"leadS"`
	State     uint8   `json:"state"`
	Live      uint8   `json:"live"`
	HalfBps   uint64  `json:"halfBps"`
	ErrorBps  float64 `json:"errorBps"`
	Inside    bool    `json:"inside"`
}

type write struct {
	ms     uint64
	prices []band.Written
}

// Replay rebuilds the band from the histories as if each point had been written when signed, then quotes it as the
// closure's last public 24/7 value is written.
func (r Reopening) Replay() Reopen {
	out := Reopen{Chain: r.Chain, Asset: r.Asset.Symbol, Date: r.Date}
	feed := r.Asset.Feed
	if feed == "" {
		feed = r.Asset.Index
	}
	i, _ := slices.BinarySearchFunc(r.Leg, r.After.UpdatedAt*1000, func(p history.Point, ms uint64) int { return cmp.Compare(p.Ms, ms) })
	if i == 0 || r.Leg[i-1].Ms < r.SinceMs {
		return out
	}
	limit := r.Leg[i-1].Ms
	nowS := (limit + 999) / 1000
	out.Evaluated, out.LeadS = true, r.After.UpdatedAt-nowS
	var writes []write
	for _, p := range r.Leg {
		if p.Ms <= limit {
			writes = append(writes, write{p.Ms, []band.Written{{Feed: feed, Value: p.Scaled(), Ms: p.Ms}}})
		}
	}
	at := map[uint64]int{}
	for _, series := range r.Status {
		for _, p := range series {
			at[p.Ms]++
		}
	}
	for _, p := range r.Status[0] {
		if at[p.Ms] != 3 || p.Ms > limit {
			continue
		}
		w := write{ms: p.Ms}
		for i, f := range []string{band.CurrentStatus, band.NextStatus, band.NextChangeTime} {
			j, _ := slices.BinarySearchFunc(r.Status[i], p.Ms, func(q history.Point, ms uint64) int { return cmp.Compare(q.Ms, ms) })
			w.prices = append(w.prices, band.Written{Feed: f, Value: r.Status[i][j].Scaled(), Ms: p.Ms})
		}
		writes = append(writes, w)
	}
	slices.SortStableFunc(writes, func(a, b write) int { return cmp.Compare(a.ms, b.ms) })
	record := band.NewRecord(feed)
	cl := band.Chainlink{}
	if r.Before.Answer != nil && r.Before.Answer.IsUint64() && r.Before.UpdatedAt <= nowS {
		cl = band.Chainlink{Answer: r.Before.Answer.Uint64(), UpdatedAt: r.Before.UpdatedAt}
	}
	for _, w := range writes {
		record.Write(w.prices...)
		if w.prices[0].Feed == feed && cl.UpdatedAt*1000 <= w.ms {
			record.Reanchor(r.Asset, cl, w.ms/1000)
		}
	}
	q := record.Quote(r.Asset, cl, r.Multiplier, r.SequencerUp, nowS)
	out.State, out.Live, out.HalfBps = uint8(q.State), q.Live, q.HalfBps
	if q.State != band.Halted && r.After.Answer != nil {
		after := new(big.Int).Set(r.After.Answer)
		out.Inside = after.Cmp(new(big.Int).SetUint64(q.Low)) >= 0 && after.Cmp(q.High) <= 0
		a, _ := new(big.Float).SetInt(after).Float64()
		out.ErrorBps = math.Abs(a/float64(q.Mid)-1) * 10_000
	}
	return out
}
