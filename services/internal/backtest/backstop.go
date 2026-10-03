// SPDX-License-Identifier: MIT OR Apache-2.0

package backtest

import (
	"math"
	"math/big"

	"github.com/tapehouse/tapehouse/services/internal/margin"
)

// Constants the backstop's calibration reads from the liquidator and the margin accounts.
const (
	// MaxDiscountBps is the liquidator's deepest discount, what a closure's liquidation is assumed to sell at.
	MaxDiscountBps = 10_00
	// ReserveShareBps is the share of each premium the margin accounts keep as a reserve.
	ReserveShareBps = 10_00
	// MaxPremiumRateBps is the most the accounts' weekend premium may be, a year while closed.
	MaxPremiumRateBps = 100_00
	yearHours         = 365 * 24
)

// Limit is the most the backstop would have covered for one kind of position in one closure, in USD, and which.
type Limit struct {
	Position string  `json:"position"`
	Limit    float64 `json:"limit"`
	Date     string  `json:"date"`
}

// Backstop is the backstop's exposure limits and the accounts' premium rate calibrated over closures: each Stock
// Token's isolated positions, and the cross positions holding every asset in proportion to its cap, levered to the
// requirement across a closure on the debt cap, or on the caps where those bind first, and sold after the closure at
// the liquidator's deepest discount. The premium rate is the one whose share for the backstop pays what it would have absorbed from the cross
// positions at the debt cap over every closure.
type Backstop struct {
	DebtCap        float64 `json:"debtCap"`
	Caps           []Limit `json:"caps"`
	Limits         []Limit `json:"limits"`
	Absorbed       float64 `json:"absorbed"`
	ClosedYears    float64 `json:"closedYears"`
	PremiumRateBps uint32  `json:"premiumRateBps"`
}

type held struct {
	gross []float64
	debt  float64
}

func hold(set *margin.Set, caps []float64, debtCap float64) held {
	exposures := make([]*big.Int, len(caps))
	gross := 0.0
	for k, c := range caps {
		exposures[k] = new(big.Int).Mul(big.NewInt(int64(c)), big.NewInt(1e18))
		gross += c
	}
	open, closed := margin.ScenarioRequirements(set, exposures)
	need := margin.Current(open, closed, margin.LeverageFloor(exposures), margin.Regime{Kind: margin.Closed})
	required, _ := new(big.Float).Quo(new(big.Float).SetInt(need), big.NewFloat(1e18)).Float64()
	lent := 1 - required/gross
	scale := min(1, debtCap/(gross*lent))
	h := held{debt: gross * lent * scale}
	for _, c := range caps {
		h.gross = append(h.gross, c*scale)
	}
	return h
}

func (h held) absorbed(moves []float64) float64 {
	proceeds := 0.0
	for k, g := range h.gross {
		proceeds += g * (1 + moves[k]) * (1 - MaxDiscountBps/10_000.0)
	}
	return max(0, h.debt-proceeds)
}

// Calibrate sizes the exposure limits and the premium rate. caps holds each asset's cap in USD, its selling depth.
func Calibrate(set *margin.Set, names []string, caps []float64, debtCap float64, closures []Closure) Backstop {
	b := Backstop{DebtCap: debtCap}
	positions := map[string][]float64{"cross": caps}
	order := []string{"cross"}
	for k, n := range names {
		only := make([]float64, len(caps))
		only[k] = caps[k]
		positions[n] = only
		order = append(order, n)
		b.Caps = append(b.Caps, Limit{Position: n, Limit: caps[k]})
	}
	hours := 0
	for _, c := range closures {
		hours += c.Hours
	}
	b.ClosedYears = float64(hours) / yearHours
	for _, p := range order {
		h := hold(set, positions[p], debtCap)
		l := Limit{Position: p}
		for _, c := range closures {
			a := h.absorbed(c.Moves)
			if a > l.Limit {
				l.Limit, l.Date = a, date(c.Open)
			}
			if p == "cross" {
				b.Absorbed += a
			}
		}
		l.Limit = math.Ceil(l.Limit)
		b.Limits = append(b.Limits, l)
		if p == "cross" && h.debt > 0 {
			rate := b.Absorbed / (h.debt * b.ClosedYears * (1 - ReserveShareBps/10_000.0))
			b.PremiumRateBps = uint32(min(math.Ceil(rate*10_000), MaxPremiumRateBps))
		}
	}
	return b
}
