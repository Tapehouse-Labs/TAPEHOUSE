// SPDX-License-Identifier: MIT OR Apache-2.0

package backtest

import (
	"math"
	"math/big"
	"math/rand/v2"
	"os"
	"testing"
	"time"

	"github.com/tapehouse/tapehouse/services/internal/band"
	"github.com/tapehouse/tapehouse/services/internal/history"
	"github.com/tapehouse/tapehouse/services/internal/margin"
)

func synthetic(sessions int) Market {
	rng := rand.New(rand.NewPCG(1, 2))
	m := Market{Names: names}
	at := time.Date(2019, 1, 2, 14, 30, 0, 0, time.UTC)
	price := []float64{100, 100, 100, 100, 100, 100}
	for len(m.At) < sessions {
		if at.Weekday() != time.Saturday && at.Weekday() != time.Sunday {
			opens, closes := make([]float64, 6), make([]float64, 6)
			gap := 0.0
			if len(m.At) > 0 && days(m.At[len(m.At)-1], at) >= 2 {
				gap = rng.NormFloat64() * 0.03
			}
			for k := range price {
				opens[k] = price[k] * (1 + gap + rng.NormFloat64()*0.005)
				closes[k] = opens[k] * (1 + rng.NormFloat64()*0.02)
				price[k] = closes[k]
			}
			m.At, m.Open, m.Close, m.Adj = append(m.At, at), append(m.Open, opens), append(m.Close, closes), append(m.Adj, closes)
		}
		at = at.Add(24 * time.Hour)
	}
	return m
}

func TestTheBacktestMeasuresEverySection(t *testing.T) {
	data, err := os.ReadFile("../../../stylus/contracts/margin/parameters.json")
	if err != nil {
		t.Fatal(err)
	}
	p, depths, err := margin.Launch(data, names)
	if err != nil {
		t.Fatal(err)
	}
	set := margin.NewSet(p, 256, margin.Horizon)
	m := synthetic(3_000)
	r := Run(m, set, depths, 100_000, 0.1, m.At[2_900])
	closures := m.Closures()
	if r.History.Closures != len(closures) || r.History.Since == 0 || r.History.Since >= len(closures) {
		t.Fatalf("history: %+v", r.History)
	}
	worst := 0.0
	for _, c := range closures {
		for _, x := range c.Moves {
			worst = max(worst, math.Abs(x))
		}
	}
	if r.Leverage.Cap != Cap(worst) || r.Leverage.DeployedCap != 5 {
		t.Errorf("cap %d for a worst move of %v", r.Leverage.Cap, worst)
	}
	if r.Leverage.AtCap["long"].Count != 0 || r.Leverage.AtCap["mixed"].Count != 0 {
		t.Errorf("a portfolio lost more than its margin at the recalibrated cap: %+v", r.Leverage.AtCap)
	}
	if r.Capacity.Model <= 0 || r.Capacity.Flat != 0.5 || r.Capacity.BadDebtFlat > r.Capacity.BadDebtModel+1 {
		t.Errorf("capacity: %+v", r.Capacity)
	}
	for _, c := range r.Coverage {
		if c.Portfolios != Count || c.Ramp[0] > c.Ramp[1] || c.Ramp[1] > c.Ramp[2] {
			t.Errorf("coverage: %+v", c)
		}
	}
	if len(r.Horizons) != 6 || r.Horizons[0].Days <= 0 {
		t.Errorf("horizons: %+v", r.Horizons)
	}
	if r.Backstop.PremiumRateBps > MaxPremiumRateBps || len(r.Backstop.Limits) != 7 || r.Backstop.Limits[0].Position != "cross" {
		t.Errorf("backstop: %+v", r.Backstop)
	}
	for _, l := range r.Backstop.Limits {
		if l.Limit < 0 || l.Limit > 100_000 {
			t.Errorf("a limit beyond the debt cap: %+v", l)
		}
	}
	if r.Basket.Rebalances < 30 || r.Basket.MedianCost <= 0 || r.Basket.MedianCost > 2*bandFloorBps/10_000.0 {
		t.Errorf("basket: %+v", r.Basket)
	}
	if d := r.Liquidator.Decays[0]; d.ClosedLvf <= 0 || d.ClosedLvf >= d.OpenLvf || d.ClosedFill <= d.OpenFill {
		t.Errorf("the closed decay is not the slower one: %+v", d)
	}
	if r.Deleveraging.Long < 0 || r.Deleveraging.Long > 1 || r.Liquidator.OpenMinutesToFloor != 1000.0/60 {
		t.Errorf("deleveraging %+v, liquidator %+v", r.Deleveraging, r.Liquidator)
	}
	if g := r.GapCover["history"]; g.Weekends != len(closures) || g.Raised < 0 || g.Raised > 1 {
		t.Errorf("gap cover: %+v", g)
	}
	if r.Reviews.LongestClosureHours != 48 || r.Reviews.ClosedShare < 0.25 || r.Reviews.ClosedShare > 0.3 {
		t.Errorf("reviews: %+v", r.Reviews)
	}
	if s := r.March2020; s == nil || s.Date != "2020-03-16" || len(s.Moves) != 6 || s.Binds["long"] > Count || s.CapacityScenario > r.Leverage.CapacityModel {
		t.Errorf("16 March 2020: %+v", s)
	}
}

func TestAReopenIsReplayedFromTheHistoriesAsTheBandWouldHaveWrittenThem(t *testing.T) {
	const friday = 1_789_761_600
	sunday := uint64(friday + 2*86_400)
	status := [3][]history.Point{
		{{Ms: (sunday - 1_800) * 1000, Value: 1.02}},
		{{Ms: (sunday - 1_800) * 1000, Value: 1}},
		{{Ms: (sunday - 1_800) * 1000, Value: float64(sunday*1000 + 48_600_000)}},
	}
	leg := []history.Point{{Ms: (friday - 600) * 1000, Value: 219}, {Ms: (sunday - 120) * 1000, Value: 220}, {Ms: sunday*1000 - 59_500, Value: 221}, {Ms: (sunday + 60) * 1000, Value: 230}}
	r := Reopening{
		Chain:       4663,
		Asset:       band.Asset{Symbol: "NVDA", Feed: "NVDA---24_7"},
		Date:        "2026-09-21",
		SinceMs:     friday * 1000,
		Leg:         leg,
		Status:      status,
		Before:      history.Round{Answer: big.NewInt(21_900_000_000), UpdatedAt: friday + 14_400},
		After:       history.Round{Answer: big.NewInt(22_150_000_000), UpdatedAt: sunday + 20},
		Multiplier:  big.NewInt(1e18),
		SequencerUp: true,
	}
	got := r.Replay()
	if !got.Evaluated || got.State != uint8(band.Closed) || got.Live != 1 || got.LeadS != 79 || !got.Inside || math.Abs(got.ErrorBps-22.62) > 0.01 {
		t.Errorf("%+v", got)
	}
	r.SinceMs = (sunday - 30) * 1000
	if late := r.Replay(); late.Evaluated || late.Inside {
		t.Errorf("a closure with no 24/7 value inside it was evaluated: %+v", late)
	}
	s := Summarize([]Reopen{got, {Chain: 42161, Evaluated: true, State: 1, ErrorBps: 5, LeadS: 30}, {Chain: 4663}})
	if len(s) != 2 || s[0].Reopens != 2 || s[0].Evaluated != 1 || s[0].Inside != 1 || s[0].MedianError != got.ErrorBps || s[0].MedianLeadS != 79 || s[1].MedianError != 5 {
		t.Errorf("%+v", s)
	}
}
