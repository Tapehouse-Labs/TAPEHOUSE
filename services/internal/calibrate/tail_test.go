// SPDX-License-Identifier: MIT OR Apache-2.0

package calibrate

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
	"time"

	"github.com/tapehouse/tapehouse/services/internal/backtest"
)

func TestTheTailIsAParetoOfShapeOneHalf(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	const u, beta, n = 0.25, 0.13, 400_000
	falls := make([]float64, 0, n)
	for range n * 95 / 100 {
		falls = append(falls, u*rng.Float64())
	}
	for range n * 5 / 100 {
		falls = append(falls, u+beta/0.5*(math.Pow(1-rng.Float64(), -0.5)-1))
	}
	tail := FitTail(falls, nil, 1)
	if math.Abs(tail.Threshold/u-1) > 0.01 || math.Abs(tail.Scale/beta-1) > 0.05 || math.Abs(tail.Probability/0.05-1) > 0.02 {
		t.Fatalf("%+v", tail)
	}
	if tail.Excesses != n*5/100 {
		t.Fatalf("%d excesses", tail.Excesses)
	}
}

func bootstrap(weeks [][]float64, seed uint64) (low, high float64) {
	rng := rand.New(rand.NewPCG(seed, seed))
	means := make([]float64, 20_000)
	for i := range means {
		var sum float64
		var count int
		for range weeks {
			week := weeks[rng.IntN(len(weeks))]
			for _, r := range week {
				sum += r
			}
			count += len(week)
		}
		means[i] = sum / float64(count)
	}
	slices.Sort(means)
	return means[999], means[18_999]
}

func TestTheReopenRatioIsBootstrappedOverWeeks(t *testing.T) {
	rng := rand.New(rand.NewPCG(4, 4))
	weeks := make([][]float64, 13)
	for i := range weeks {
		for range 5 {
			weeks[i] = append(weeks[i], 0.5+0.5*rng.Float64())
		}
	}
	falls := []float64{0.1, 0.2, 0.3, 0.4, 0.5}
	tail := FitTail(falls, weeks, 1)
	low, high := bootstrap(weeks, 1)
	if tail.ReopenLow != low || tail.ReopenHigh != high || low >= high {
		t.Fatalf("[%g, %g], want [%g, %g]", tail.ReopenLow, tail.ReopenHigh, low, high)
	}
	if again := FitTail(falls, weeks, 1); again.ReopenLow != tail.ReopenLow || again.ReopenHigh != tail.ReopenHigh {
		t.Fatal("the same seed gave another interval")
	}
}

func TestTheScaledConstantsAreTheFitTimesTheUpperRatio(t *testing.T) {
	tail := Tail{Threshold: 0.247687, Scale: 0.129509, Probability: 0.050224, ReopenHigh: 0.74}
	if got := tail.Constants(); got != [3]uint32{183288, 95837, 50224} {
		t.Fatalf("constants %v", got)
	}
}

func TestATailMovingMoreThanATenthIsFlagged(t *testing.T) {
	deployed := [3]uint32{183288, 95837, 50224}
	if tail := (Tail{ReopenHigh: 0.74, Scaled: [3]uint32{183288, 105500, 50224}}).Against(deployed); !tail.Redeploy || tail.Deployed != deployed {
		t.Fatalf("%+v", tail)
	}
	if tail := (Tail{ReopenHigh: 0.74, Scaled: [3]uint32{183288, 105000, 50224}}).Against(deployed); tail.Redeploy {
		t.Fatalf("%+v", tail)
	}
}

func TestATailWithoutAReopenRatioIsNotFlagged(t *testing.T) {
	if tail := (Tail{Threshold: 0.2, Scale: 0.1, Probability: 0.05}).Against([3]uint32{183288, 95837, 50224}); tail.Redeploy {
		t.Fatalf("%+v", tail)
	}
}

func TestTheFallsAreWeekendFallsOverTheCurrentGap(t *testing.T) {
	at := func(day int) time.Time { return time.Date(2026, 9, day, 13, 30, 0, 0, time.UTC) }
	m := backtest.Market{
		Names: []string{"AAA"},
		At:    []time.Time{at(1), at(2), at(7), at(8), at(11)},
		Open:  [][]float64{{100}, {100}, {90}, {100}, {105}},
		Close: [][]float64{{100}, {100}, {100}, {100}, {100}},
		Adj:   [][]float64{{100}, {100}, {100}, {100}, {100}},
	}
	// Sessions 2 to 7 are five days apart, 8 to 11 three days apart, 1 to 2 one day: one fall (-10%), one rise.
	got := Falls(m, []uint32{50_000})
	if len(got) != 1 || math.Abs(got[0]-2) > 1e-9 {
		t.Fatalf("falls %v", got)
	}
}
