// SPDX-License-Identifier: MIT OR Apache-2.0

package backtest

import "time"

// The weekend gap cover's pricing, as GapCover's constants: the gap that follows the prior week, its floor and cap as
// shares of the engine's gap, and the share of the move from Friday's close to Monday's open the 24/5 reopen sees.
const (
	MoveToGap      = 2.0
	MinGap         = 0.55
	MaxGapMultiple = 8
	ReopenShare    = 0.74
)

// Weighting is how often a weekend's realised move to the 24/5 reopen exceeded the gap the cover priced with, under
// demand spread uniformly, weighted by the prior week's move and by its square, and how often it exceeded the
// engine's gap alone.
type Weighting struct {
	Uniform   float64 `json:"uniform"`
	Volatile  float64 `json:"volatile"`
	Variance  float64 `json:"variance"`
	EngineGap float64 `json:"engineGap"`
}

// GapCover is the cover's pricing gap reviewed on closures it was not fitted on: the share of weekends whose gap
// the prior week raised above the engine's, the exceedances, and on the weekends left at the engine's gap, the
// realised move over the priced gap and the exceedance at a multiplier below one.
type GapCover struct {
	From       string             `json:"from"`
	Weekends   int                `json:"weekends"`
	Raised     float64            `json:"raised"`
	Exceeded   Weighting          `json:"exceeded"`
	CalmRatio  float64            `json:"calmRatio"`
	CalmExceed map[string]float64 `json:"calmExceed"`
}

// ReviewGapCover measures the cover's pricing gap over the closures from from on; gaps holds the engine's weekend
// gap of each asset, in millionths.
func ReviewGapCover(closures []Closure, gaps []uint32, from time.Time) GapCover {
	g := GapCover{From: date(from), CalmExceed: map[string]float64{}}
	var raised, n, calm float64
	var exceeded, volW, varW, volX, varX, engine, calmMove, calmGap float64
	multipliers := map[string]float64{"1": 1, "0.75": 0.75, "0.5": 0.5}
	calmHits := map[string]float64{}
	for _, c := range closures {
		if c.Open.Before(from) {
			continue
		}
		g.Weekends++
		for k, m := range c.Moves {
			engineGap := float64(gaps[k]) / 1e6
			floor := MinGap * engineGap
			priced := min(max(MoveToGap*c.Prior[k], floor), MaxGapMultiple*engineGap)
			realised := ReopenShare * max(m, -m)
			n++
			hit := 0.0
			if realised > priced {
				hit = 1
			}
			exceeded += hit
			volW, volX = volW+c.Prior[k], volX+hit*c.Prior[k]
			varW, varX = varW+c.Prior[k]*c.Prior[k], varX+hit*c.Prior[k]*c.Prior[k]
			if realised > engineGap {
				engine++
			}
			if MoveToGap*c.Prior[k] > floor {
				raised++
				continue
			}
			calm++
			calmMove, calmGap = calmMove+realised, calmGap+floor
			for label, x := range multipliers {
				if realised > x*floor {
					calmHits[label]++
				}
			}
		}
	}
	if n == 0 {
		return g
	}
	g.Raised = raised / n
	g.Exceeded = Weighting{Uniform: exceeded / n, Volatile: volX / volW, Variance: varX / varW, EngineGap: engine / n}
	if calm > 0 {
		g.CalmRatio = calmMove / calmGap
		for label := range multipliers {
			g.CalmExceed[label] = calmHits[label] / calm
		}
	}
	return g
}
