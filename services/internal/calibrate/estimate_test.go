// SPDX-License-Identifier: MIT OR Apache-2.0

package calibrate

import (
	"slices"
	"testing"
	"time"

	"github.com/tapehouse/tapehouse/services/internal/history"
)

type lcg uint64

func (x *lcg) draw() uint64 {
	*x = *x*6364136223846793005 + 1442695040888963407
	return uint64(*x) >> 33
}

var syntheticNames = []string{"AAA", "BBB", "MKT"}

// syntheticMarket is three assets over 700 sessions from Monday 6 January 2020, a three-day gap before every fifth
// session, drawn from 64-bit LCGs; the third asset is the market.
func syntheticMarket() []history.Series {
	start := time.Date(2020, 1, 6, 13, 30, 0, 0, time.UTC)
	common, asset := lcg(7), []lcg{1, 2, 3}
	series := make([]history.Series, 3)
	closes := []float64{100, 50, 200}
	for k := range series {
		series[k].Symbol = syntheticNames[k]
	}
	at := start
	for i := range 700 {
		if i > 0 {
			at = at.AddDate(0, 0, 1)
			if i%5 == 0 {
				at = at.AddDate(0, 0, 2)
			}
		}
		if i == 0 {
			for k := range series {
				series[k].Bars = append(series[k].Bars, history.Bar{At: at, Open: closes[k], Close: closes[k], AdjClose: closes[k]})
			}
			continue
		}
		m := (float64(common.draw()%2001) - 1000) / 40000
		for k := range series {
			e := (float64(asset[k].draw()%2001) - 1000) / float64(30000+10000*k)
			g := 0.0
			if i%5 == 0 {
				g = (float64(asset[k].draw()%2001) - 1000) / 20000
			}
			open := closes[k] * (1 + g)
			closes[k] = open * (1 + m + e)
			series[k].Bars = append(series[k].Bars, history.Bar{At: at, Open: open, Close: closes[k], AdjClose: closes[k]})
		}
	}
	return series
}

func syntheticReturns(t *testing.T) [][]float64 {
	t.Helper()
	series := syntheticMarket()
	_, returns := Returns(series, series[0].Bars[699].At)
	if len(returns) != 3 || len(returns[0]) != 699 {
		t.Fatalf("%d series of %d returns", len(returns), len(returns[0]))
	}
	return returns
}

func TestTheVolatilitiesAreTheLaunchMethods(t *testing.T) {
	floors, years := Volatilities(syntheticReturns(t))
	if !slices.Equal(floors, []uint32{26777, 24129, 23195}) || !slices.Equal(years, []uint32{27165, 23694, 22453}) {
		t.Fatalf("floors %v, last 252 %v", floors, years)
	}
}

func TestTheStressWindowIsTheMarketsWorstSixtySessions(t *testing.T) {
	if got := StressStart(syntheticReturns(t)[2]); got != 421 {
		t.Fatalf("stress window starts at return %d", got)
	}
}

func TestTheCorrelationsAreTheLaunchMethods(t *testing.T) {
	floors, years := Correlations(syntheticReturns(t), 2)
	if !slices.Equal(floors, []uint16{3259, 3260, 4320}) || !slices.Equal(years, []uint16{2779, 3992, 4532}) {
		t.Fatalf("floors %v, last 252 %v", floors, years)
	}
}

func TestTheGapsAreTheMeanOfTheLargestOnePercent(t *testing.T) {
	series := syntheticMarket()
	got := Gaps(series, series[0].Bars[699].At)
	if !slices.Equal(got, []uint32{49825, 49500, 49925}) {
		t.Fatalf("gaps %v", got)
	}
}

func TestATargetIsTheLargerOfFloorAndYear(t *testing.T) {
	series := syntheticMarket()
	got, err := Estimate(series, syntheticNames, 2, series[0].Bars[699].At)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got.Volatilities, []uint32{27165, 24129, 23195}) || !slices.Equal(got.Correlations, []uint16{3259, 3992, 4532}) ||
		!slices.Equal(got.Gaps, []uint32{49825, 49500, 49925}) {
		t.Fatalf("targets %+v", got)
	}
	if _, err := Estimate(series, []string{"AAA", "ZZZ"}, 0, series[0].Bars[699].At); err == nil {
		t.Fatal("an asset without a series was estimated")
	}
}

func TestRoundingIsHalfToEven(t *testing.T) {
	for _, c := range []struct {
		in   float64
		want uint32
	}{{2.5e-6, 2}, {3.5e-6, 4}, {1.5e-6, 2}} {
		if got := millionths(c.in); got != c.want {
			t.Errorf("round(%g · 1e6) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestAGapIsCappedAtOneHundredPercent(t *testing.T) {
	series := syntheticMarket()[:1]
	for i := range series[0].Bars {
		if i > 0 && i%5 == 0 {
			prev := series[0].Bars[i-1].Close
			series[0].Bars[i].Open = prev * 3.5
		}
	}
	if got := Gaps(series, series[0].Bars[699].At); got[0] != 1_000_000 {
		t.Fatalf("gap %d", got[0])
	}
}

func TestTheStressWindowIsSPYsWhereTheEngineHasNoMarketAsset(t *testing.T) {
	series := syntheticMarket()
	series[2].Symbol = StressMarket
	last := series[0].Bars[699].At
	spy, err := Estimate(series, []string{"AAA", "BBB"}, -1, last)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(spy.Correlations, []uint16{3259}) {
		t.Fatalf("correlations %v", spy.Correlations)
	}
	equal, err := Estimate(series[:2], []string{"AAA", "BBB"}, -1, last)
	if err != nil {
		t.Fatal(err)
	}
	if slices.Equal(equal.Correlations, spy.Correlations) {
		t.Fatalf("the equal-weighted stress window gave the same %v", equal.Correlations)
	}
}
