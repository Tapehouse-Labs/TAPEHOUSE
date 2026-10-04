// SPDX-License-Identifier: MIT OR Apache-2.0

package calibrate

import (
	"errors"

	"github.com/tapehouse/tapehouse/services/internal/margin"
)

// What the engine accepts of a parameter, as stylus/contracts/margin checks it.
const (
	maxVolatility  = 1_000_000
	maxGap         = 1_000_000
	maxCorrelation = 10_000
	// blends is how many halvings of a step toward the current matrix are tried before it is kept.
	blends = 6
)

// Current is what the engine holds, in its order, with the floors and ceilings its values stay within.
type Current struct {
	Volatilities      []uint32
	Gaps              []uint32
	Correlations      []uint16
	Depths            []uint32
	VolatilityFloors  []uint32
	GapFloors         []uint32
	CorrelationFloors []uint16
	DepthCeilings     []uint32
}

// Proposed is what setParameters takes. Shrunk is k of the first 2^-k blend of the correlations toward the current
// matrix that is positive definite, 0 where the clamped matrix was, and blends+1 where none was and the current
// matrix is kept.
type Proposed struct {
	Volatilities []uint32
	Gaps         []uint32
	Correlations []uint16
	Depths       []uint32
	Shrunk       int
}

// Pairs are the pairs of names the engine stores correlations for: the upper triangle, row by row.
func Pairs(names []string) [][2]string {
	var out [][2]string
	for i, a := range names {
		for _, b := range names[i+1:] {
			out = append(out, [2]string{a, b})
		}
	}
	return out
}

// box is the values the engine's within_step accepts after previous: 2·value ≤ 3·previous and 3·value ≥ 2·previous.
func box(previous uint64) (low, high uint64) {
	return (2*previous + 2) / 3, 3 * previous / 2
}

func clamp(target, low, high uint64) uint64 {
	return min(max(target, low), high)
}

// step is target clamped to [floor, cap], then to the box around previous.
func step(previous, target, floor, ceiling uint64) uint64 {
	low, high := box(previous)
	return clamp(clamp(target, floor, ceiling), low, high)
}

// Step proposes each target from the engine's current value, within its floor or ceiling and the engine's ×1.5
// bound, in integers. A depth falls to its target at once and rises at most ×1.5; depths nil keeps the current
// ones. The correlations are blended toward the current matrix until positive definite.
func Step(current Current, targets Targets, depths []uint32) (Proposed, error) {
	n := len(current.Volatilities)
	if len(targets.Volatilities) != n || len(targets.Gaps) != n || len(current.Gaps) != n ||
		len(targets.Correlations) != len(current.Correlations) || len(current.Depths) != 2*n ||
		(depths != nil && len(depths) != 2*n) || len(current.VolatilityFloors) != n || len(current.GapFloors) != n ||
		len(current.CorrelationFloors) != len(current.Correlations) || len(current.DepthCeilings) != 2*n {
		return Proposed{}, errors.New("the current values and the targets do not fit one set of assets")
	}
	var p Proposed
	for i := range n {
		p.Volatilities = append(p.Volatilities, uint32(step(uint64(current.Volatilities[i]), uint64(targets.Volatilities[i]), uint64(current.VolatilityFloors[i]), maxVolatility)))
		p.Gaps = append(p.Gaps, uint32(step(uint64(current.Gaps[i]), uint64(targets.Gaps[i]), uint64(current.GapFloors[i]), maxGap)))
	}
	clamped := make([]uint16, len(current.Correlations))
	for k := range clamped {
		clamped[k] = uint16(step(uint64(current.Correlations[k]), uint64(targets.Correlations[k]), uint64(current.CorrelationFloors[k]), maxCorrelation))
	}
	p.Correlations, p.Shrunk = positiveDefinite(n, current.Correlations, clamped)
	for k, previous := range current.Depths {
		target := uint64(previous)
		if depths != nil {
			target = uint64(depths[k])
		}
		_, high := box(uint64(previous))
		p.Depths = append(p.Depths, uint32(min(clamp(target, 1, uint64(current.DepthCeilings[k])), high)))
	}
	return p, nil
}

// positiveDefinite returns clamped where it is positive definite, else the first of previous + (clamped - previous)
// / 2^k, each value rounded toward previous, that is, else previous.
func positiveDefinite(n int, previous, clamped []uint16) ([]uint16, int) {
	if margin.PositiveDefinite(n, clamped) {
		return clamped, 0
	}
	for k := 1; k <= blends; k++ {
		blend := make([]uint16, len(clamped))
		for i := range blend {
			blend[i] = uint16(int(previous[i]) + (int(clamped[i])-int(previous[i]))/(1<<k))
		}
		if margin.PositiveDefinite(n, blend) {
			return blend, k
		}
	}
	return previous, blends + 1
}
