// SPDX-License-Identifier: MIT OR Apache-2.0

package calibrate

import (
	"slices"
	"testing"

	"github.com/tapehouse/tapehouse/services/internal/margin"
)

func withinStep(previous, value uint64) bool {
	return 2*value <= 3*previous && 3*value >= 2*previous
}

func single(previous, floor uint32) Current {
	return Current{
		Volatilities: []uint32{previous}, VolatilityFloors: []uint32{floor},
		Gaps: []uint32{previous}, GapFloors: []uint32{floor},
		Depths: []uint32{previous, previous}, DepthCeilings: []uint32{10_000_000, 10_000_000},
	}
}

func TestAStepStaysInItsBox(t *testing.T) {
	for _, c := range []struct{ previous, target, want uint32 }{
		{31352, 60000, 47028},
		{31352, 10000, 20902},
		{30, 45, 45},
		{30, 46, 45},
		{30, 20, 20},
		{30, 19, 20},
	} {
		got, err := Step(single(c.previous, 1), Targets{Volatilities: []uint32{c.target}, Gaps: []uint32{c.target}}, nil)
		if err != nil {
			t.Fatal(err)
		}
		if got.Volatilities[0] != c.want || got.Gaps[0] != c.want || !withinStep(uint64(c.previous), uint64(got.Volatilities[0])) {
			t.Errorf("%d to %d: %d, want %d", c.previous, c.target, got.Volatilities[0], c.want)
		}
	}
}

func TestAValueStaysAboveItsFloorAndCap(t *testing.T) {
	got, _ := Step(single(31352, 31352), Targets{Volatilities: []uint32{100}, Gaps: []uint32{100}}, nil)
	if got.Volatilities[0] != 31352 || got.Gaps[0] != 31352 {
		t.Fatalf("below the floor: %d, %d", got.Volatilities[0], got.Gaps[0])
	}
	got, _ = Step(single(900_000, 1), Targets{Volatilities: []uint32{5_000_000}, Gaps: []uint32{5_000_000}}, nil)
	if got.Volatilities[0] != 1_000_000 || got.Gaps[0] != 1_000_000 {
		t.Fatalf("above the cap: %d, %d", got.Volatilities[0], got.Gaps[0])
	}
	two := Current{
		Volatilities: []uint32{1, 1}, VolatilityFloors: []uint32{1, 1}, Gaps: []uint32{1, 1}, GapFloors: []uint32{1, 1},
		Correlations: []uint16{9_000}, CorrelationFloors: []uint16{5_000}, Depths: []uint32{1, 1, 1, 1}, DepthCeilings: []uint32{9, 9, 9, 9},
	}
	if got := step(9_000, 12_000, 5_000, maxCorrelation); got != 10_000 {
		t.Fatalf("correlation %d", got)
	}
	got, _ = Step(two, Targets{Volatilities: []uint32{1, 1}, Gaps: []uint32{1, 1}, Correlations: []uint16{100}}, nil)
	if got.Correlations[0] != 6_000 {
		t.Fatalf("correlation %d", got.Correlations[0])
	}
	two.CorrelationFloors[0] = 8_000
	got, _ = Step(two, Targets{Volatilities: []uint32{1, 1}, Gaps: []uint32{1, 1}, Correlations: []uint16{100}}, nil)
	if got.Correlations[0] != 8_000 {
		t.Fatalf("correlation %d", got.Correlations[0])
	}
}

func TestADepthFallsFreelyAndRisesByHalfAtMost(t *testing.T) {
	for _, c := range []struct{ previous, ceiling, target, want uint32 }{
		{3561773, 4_000_000, 1000, 1000},
		{100000, 1_000_000, 400000, 150000},
		{100000, 120000, 400000, 120000},
		{100000, 1_000_000, 0, 1},
	} {
		current := single(1, 1)
		current.Depths = []uint32{c.previous, c.previous}
		current.DepthCeilings = []uint32{c.ceiling, c.ceiling}
		got, err := Step(current, Targets{Volatilities: []uint32{1}, Gaps: []uint32{1}}, []uint32{c.target, c.target})
		if err != nil || got.Depths[0] != c.want || got.Depths[1] != c.want {
			t.Errorf("%d to %d under %d: %v, %v, want %d", c.previous, c.target, c.ceiling, got.Depths, err, c.want)
		}
	}
	got, _ := Step(single(1, 1), Targets{Volatilities: []uint32{1}, Gaps: []uint32{1}}, nil)
	if !slices.Equal(got.Depths, []uint32{1, 1}) {
		t.Fatalf("without depths %v", got.Depths)
	}
}

func TestAMatrixThatIsNotPositiveDefiniteIsPulledBack(t *testing.T) {
	current := Current{
		Volatilities: []uint32{1, 1, 1}, VolatilityFloors: []uint32{1, 1, 1}, Gaps: []uint32{1, 1, 1}, GapFloors: []uint32{1, 1, 1},
		Correlations: []uint16{6000, 6000, 6000}, CorrelationFloors: []uint16{0, 0, 0},
		Depths: []uint32{1, 1, 1, 1, 1, 1}, DepthCeilings: []uint32{1, 1, 1, 1, 1, 1},
	}
	targets := Targets{Volatilities: []uint32{1, 1, 1}, Gaps: []uint32{1, 1, 1}, Correlations: []uint16{10000, 10000, 0}}
	if !margin.PositiveDefinite(3, current.Correlations) || margin.PositiveDefinite(3, []uint16{9000, 9000, 4000}) {
		t.Fatal("the fixture's matrices are not what the test says")
	}
	got, err := Step(current, targets, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.Shrunk == 0 || !margin.PositiveDefinite(3, got.Correlations) {
		t.Fatalf("shrunk %d: %v", got.Shrunk, got.Correlations)
	}
	for i, v := range got.Correlations {
		if !withinStep(6000, uint64(v)) {
			t.Errorf("correlation %d, %d, is out of its box", i, v)
		}
	}
	for k := 1; k < got.Shrunk; k++ {
		blend := make([]uint16, 3)
		for i, c := range []int{9000, 9000, 4000} {
			blend[i] = uint16(6000 + (c-6000)/(1<<k))
		}
		if margin.PositiveDefinite(3, blend) {
			t.Fatalf("2^-%d was positive definite but the proposal took 2^-%d", k, got.Shrunk)
		}
	}
}

func TestTheArraysAreInTheEnginesOrder(t *testing.T) {
	want := [][2]string{{"SPY", "NVDA"}, {"SPY", "TSLA"}, {"NVDA", "TSLA"}}
	if got := Pairs([]string{"SPY", "NVDA", "TSLA"}); !slices.Equal(got, want) {
		t.Fatalf("pairs %v", got)
	}
	current := Current{
		Volatilities: []uint32{100, 200, 300}, VolatilityFloors: []uint32{1, 1, 1}, Gaps: []uint32{110, 210, 310}, GapFloors: []uint32{1, 1, 1},
		Correlations: []uint16{5000, 6000, 7000}, CorrelationFloors: []uint16{1, 1, 1},
		Depths: []uint32{10, 11, 20, 21, 30, 31}, DepthCeilings: []uint32{99, 99, 99, 99, 99, 99},
	}
	got, err := Step(current, Targets{Volatilities: []uint32{101, 202, 303}, Gaps: []uint32{111, 212, 313}, Correlations: []uint16{5001, 6002, 7003}}, []uint32{12, 13, 22, 23, 32, 33})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got.Volatilities, []uint32{101, 202, 303}) || !slices.Equal(got.Gaps, []uint32{111, 212, 313}) ||
		!slices.Equal(got.Correlations, []uint16{5001, 6002, 7003}) || !slices.Equal(got.Depths, []uint32{12, 13, 22, 23, 32, 33}) {
		t.Fatalf("%+v", got)
	}
	if _, err := Step(current, Targets{Volatilities: []uint32{1}}, nil); err == nil {
		t.Fatal("targets of the wrong length were stepped")
	}
}
