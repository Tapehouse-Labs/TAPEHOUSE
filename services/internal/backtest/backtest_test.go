// SPDX-License-Identifier: MIT OR Apache-2.0

package backtest

import (
	"math"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/tapehouse/tapehouse/services/internal/history"
	"github.com/tapehouse/tapehouse/services/internal/margin"
)

var names = []string{"NVDA", "TSLA", "AAPL", "MSFT", "GOOGL", "SPY"}

func day(s string) time.Time {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		panic(err)
	}
	return t.Add(13*time.Hour + 30*time.Minute)
}

func series(symbol string, days []string, opens, closes []float64) history.Series {
	s := history.Series{Symbol: symbol}
	for i, d := range days {
		s.Bars = append(s.Bars, history.Bar{At: day(d), Open: opens[i], Close: closes[i], AdjClose: closes[i]})
	}
	return s
}

func TestClosuresAreTheSessionsTwoOrMoreDaysApartThatEveryAssetTraded(t *testing.T) {
	days := []string{"2020-03-12", "2020-03-13", "2020-03-16", "2020-03-17", "2020-03-20", "2020-03-23"}
	a := series("A", days, []float64{100, 100, 90, 95, 100, 100}, []float64{100, 100, 95, 100, 100, 110})
	b := series("B", append([]string{"2020-03-11"}, days[:5]...), []float64{10, 10, 10, 12, 10, 10}, []float64{10, 10, 10, 10, 10, 10})
	m := Align([]history.Series{a, b})
	if len(m.At) != 5 || m.At[0] != day("2020-03-12") {
		t.Fatalf("aligned sessions: %v", m.At)
	}
	c := m.Closures()
	if len(c) != 2 || c[0].Hours != 48 || c[1].Hours != 48 {
		t.Fatalf("closures: %+v", c)
	}
	if math.Abs(c[0].Moves[0]+0.10) > 1e-12 || math.Abs(c[0].Moves[1]-0.20) > 1e-12 {
		t.Errorf("Friday close to Monday open: %v", c[0].Moves)
	}
	if c[0].Open != day("2020-03-16") || c[1].Close != day("2020-03-17") {
		t.Errorf("closure days: %v %v", c[0].Open, c[1].Close)
	}
	r := m.Returns()
	if len(r) != 2 || math.Abs(r[0][0]) > 1e-12 || math.Abs(r[1][0]-100.0/95+1) > 1e-12 {
		t.Errorf("returns: %v", r)
	}
}

func TestPortfoliosAreReproducibleAndNeverEmpty(t *testing.T) {
	p := Portfolios(7, 200, 6, 0, 3)
	q := Portfolios(7, 200, 6, 0, 3)
	if len(p) != 200 {
		t.Fatalf("%d portfolios", len(p))
	}
	for i := range p {
		empty := true
		for k, w := range p[i] {
			if w != q[i][k] || w < 0 || w > 3 {
				t.Fatalf("portfolio %d: %v", i, p[i])
			}
			empty = empty && w == 0
		}
		if empty {
			t.Fatalf("portfolio %d is empty", i)
		}
	}
}

func launch(t *testing.T) *margin.Set {
	t.Helper()
	data, err := os.ReadFile("../../../stylus/contracts/margin/parameters.json")
	if err != nil {
		t.Fatal(err)
	}
	p, _, err := margin.Launch(data, names)
	if err != nil {
		t.Fatal(err)
	}
	return margin.NewSet(p, 256, margin.Horizon)
}

func TestABookHoldsTheEngineRequirementsAsFractionsOfItsGross(t *testing.T) {
	set := launch(t)
	books := Books(set, [][]int{{1, 0, 0, 0, 0, 0}, {1, 1, 1, 1, 1, 1}})
	one := books[0]
	exposures := []*big.Int{new(big.Int).Mul(big.NewInt(Unit), big.NewInt(1e18)), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int)}
	open, closed := margin.ScenarioRequirements(set, exposures)
	gross := new(big.Float).SetInt(exposures[0])
	want, _ := new(big.Float).Quo(new(big.Float).SetInt(open), gross).Float64()
	if one.Open != want || one.Weekday() < 0.1875 || one.Floor != 0.2 {
		t.Errorf("one NVDA: %+v", one)
	}
	if c, _ := new(big.Float).Quo(new(big.Float).SetInt(closed), gross).Float64(); one.Closed != c {
		t.Errorf("closed: %v, want %v", one.Closed, c)
	}
	if books[1].Open >= one.Open {
		t.Error("six assets are not margined below one")
	}
	if l := one.Loss([]float64{-0.2, 0.5, 0, 0, 0, 0}); math.Abs(l-0.2) > 1e-12 {
		t.Errorf("loss: %v", l)
	}
	if one.Across(false) < one.Weekday() || one.Across(true) < 0.2 {
		t.Errorf("across a closure: %v, %v", one.Across(false), one.Across(true))
	}
}

func TestLossVersusFairMatchesThePapersExample(t *testing.T) {
	lvf := LossVersusFair(0.05, 1e-4, 12)
	if math.Abs(lvf-0.00133) > 0.00001 {
		t.Errorf("5%% daily volatility, 1 bp a second, 12 s blocks: %v", lvf)
	}
}

func TestTheWeekendLeverageCapIsTheLargestWholeLeverageThatCoversTheWorstClosure(t *testing.T) {
	if got := Cap(0.155); got != 5 {
		t.Errorf("15.5%%: %d", got)
	}
	if got := Cap(0.16); got != 5 {
		t.Errorf("16%%: %d", got)
	}
	if got := Cap(0.161); got != 4 {
		t.Errorf("16.1%%: %d", got)
	}
}

func TestTheIndexOverflowsSoonerAtAHigherRate(t *testing.T) {
	if y := OverflowYears(100_000); math.Abs(y-2.655) > 0.01 {
		t.Errorf("1,000%% a year: %v years", y)
	}
	if y := OverflowYears(5_125); math.Abs(y-51.8) > 0.1 {
		t.Errorf("51.25%% a year: %v years", y)
	}
}

func TestTheEnginesImpactIsLinearToItsDepthAndWholeBeyond(t *testing.T) {
	if c := EngineCost(100_000, 1_000_000, 500); math.Abs(c-0.0055) > 1e-12 {
		t.Errorf("a tenth of the depth: %v", c)
	}
	if c := EngineCost(2_000_000, 1_000_000, 3_000); math.Abs(c-(0.003+(50_000+1_000_000)/2_000_000.0)) > 1e-12 {
		t.Errorf("twice the depth: %v", c)
	}
}
