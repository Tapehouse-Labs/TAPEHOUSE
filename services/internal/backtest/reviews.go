// SPDX-License-Identifier: MIT OR Apache-2.0

package backtest

import (
	"math"
	"time"
)

// The liquidator's, the lending vaults', the shorts' and the baskets' constants the reviews weigh.
const (
	OpenDecayBpsPerMinute   = 60
	ClosedDecayBpsPerMinute = 15
	ClosedHourlyBps         = 10_00
	RecallHaircutBps        = 5_00
	MaxPremiumBps           = 1_00
	RestrictionDropBps      = 10_00
	MaxRateBps              = 1000_00
	CurveMaxRateBps         = 51_25
	BasketNoticeDays        = 7
	bandFloorBps            = 30
)

// LossVersusFair is what a Dutch auction whose log price falls deltaPerSecond loses to arbitrageurs relative to the
// fair price, with blocks blockSeconds apart and a daily volatility sigmaDaily (Moallemi and Robinson,
// arXiv:2406.00113).
func LossVersusFair(sigmaDaily, deltaPerSecond, blockSeconds float64) float64 {
	s2 := sigmaDaily * sigmaDaily / 86_400
	return 1 / (1 + deltaPerSecond/s2*(math.Sqrt(1+2*s2/(deltaPerSecond*deltaPerSecond*blockSeconds))-1))
}

// FillSeconds is the expected time such an auction takes to fill from z0 above the fair price.
func FillSeconds(sigmaDaily, deltaPerSecond, blockSeconds, z0 float64) float64 {
	s2 := sigmaDaily * sigmaDaily / 86_400
	return z0/deltaPerSecond + blockSeconds/2*(1+math.Sqrt(1+2*s2/(deltaPerSecond*deltaPerSecond*blockSeconds)))
}

// Decay is one asset's liquidation auctions at the chain's block interval: loss versus fair and the time to fill
// from 60 bps above, while the session is open and while it is closed.
type Decay struct {
	Asset      string  `json:"asset"`
	Volatility float64 `json:"volatility"`
	OpenLvf    float64 `json:"openLvf"`
	OpenFill   float64 `json:"openFill"`
	ClosedLvf  float64 `json:"closedLvf"`
	ClosedFill float64 `json:"closedFill"`
}

// Liquidator is the review of the liquidator's constants: each asset's decay, the minutes each auction takes to
// reach its deepest discount, and the price move a closed hour's share of a position at its cap makes at the
// engine's linear impact.
type Liquidator struct {
	BlockSeconds         float64 `json:"blockSeconds"`
	Decays               []Decay `json:"decays"`
	OpenMinutesToFloor   float64 `json:"openMinutesToFloor"`
	ClosedMinutesToFloor float64 `json:"closedMinutesToFloor"`
	HourlyMove           float64 `json:"hourlyMove"`
}

// ReviewLiquidator weighs the liquidator's decay and hourly share at the measured block interval.
func ReviewLiquidator(names []string, returns [][]float64, blockSeconds float64) Liquidator {
	l := Liquidator{
		BlockSeconds:         blockSeconds,
		OpenMinutesToFloor:   float64(MaxDiscountBps) / OpenDecayBpsPerMinute,
		ClosedMinutesToFloor: float64(MaxDiscountBps) / ClosedDecayBpsPerMinute,
		HourlyMove:           ClosedHourlyBps / 10_000.0 * 0.1,
	}
	open, closed := OpenDecayBpsPerMinute/60.0/10_000, ClosedDecayBpsPerMinute/60.0/10_000
	for k, n := range names {
		sigma := math.Sqrt(variance(column(returns, k)))
		l.Decays = append(l.Decays, Decay{
			Asset:      n,
			Volatility: sigma,
			OpenLvf:    LossVersusFair(sigma, open, blockSeconds),
			OpenFill:   FillSeconds(sigma, open, blockSeconds, 0.006),
			ClosedLvf:  LossVersusFair(sigma, closed, blockSeconds),
			ClosedFill: FillSeconds(sigma, closed, blockSeconds, 0.006),
		})
	}
	return l
}

// OverflowYears is how long a borrow index with 27 decimals in 128 bits takes to overflow at rateBps a year,
// accrued often enough to compound continuously.
func OverflowYears(rateBps float64) float64 {
	return (128*math.Ln2 - 27*math.Ln10) / (rateBps / 10_000)
}

// Asset is one asset's closures and sessions as the lending, recall and short constants see them: the share of
// closures that moved it up by more than the recall haircut, its largest move up over a closure, the sessions that
// closed at least the short-sale restriction's drop below the one before, and the most a buy-in can buy within the
// shorts' premium limit at the engine's linear impact, its buying depth's tenth.
type Asset struct {
	Asset          string  `json:"asset"`
	AboveHaircut   float64 `json:"aboveHaircut"`
	LargestUp      float64 `json:"largestUp"`
	RestrictedDays int     `json:"restrictedDays"`
	BuyIn          float64 `json:"buyIn"`
}

// Reviews is the measured side of the constant reviews: per asset, the closures' longest and the time closed, the
// borrow index's life at the vaults' most and the deployed curve's most, and each asset's closures.
type Reviews struct {
	LongestClosureHours int     `json:"longestClosureHours"`
	LongestClosure      string  `json:"longestClosure"`
	ClosedShare         float64 `json:"closedShare"`
	OverflowAtMax       float64 `json:"overflowAtMax"`
	OverflowAtCurve     float64 `json:"overflowAtCurve"`
	Assets              []Asset `json:"assets"`
}

// Review measures the constants' side of the market. depths holds each asset's selling then buying depth in USD.
func Review(m Market, closures []Closure, depths []uint32) Reviews {
	r := Reviews{OverflowAtMax: OverflowYears(MaxRateBps), OverflowAtCurve: OverflowYears(CurveMaxRateBps)}
	hours := 0
	for _, c := range closures {
		hours += c.Hours
		if c.Hours > r.LongestClosureHours {
			r.LongestClosureHours, r.LongestClosure = c.Hours, date(c.Open)
		}
	}
	r.ClosedShare = float64(hours) / m.At[len(m.At)-1].Sub(m.At[0]).Hours()
	returns := m.Returns()
	for k, n := range m.Names {
		a := Asset{Asset: n, BuyIn: float64(depths[2*k+1]) / 10}
		for _, c := range closures {
			if c.Moves[k] > RecallHaircutBps/10_000.0 {
				a.AboveHaircut++
			}
			a.LargestUp = max(a.LargestUp, c.Moves[k])
		}
		a.AboveHaircut /= float64(len(closures))
		for _, ret := range returns {
			if ret[k] <= -RestrictionDropBps/10_000.0 {
				a.RestrictedDays++
			}
		}
		r.Assets = append(r.Assets, a)
	}
	return r
}

// Basket is what keeping an equal-weighted basket costs: its components reset to equal value at the start of each
// quarter, what comes in valued at the low edges and what goes out at the high edges of a band at its floor, per
// rebalance and a year, as shares of the basket's value; and how far its weights drift over the seven days' notice
// a new target waits, the median and the largest.
type Basket struct {
	Components   []string `json:"components"`
	From         string   `json:"from"`
	Rebalances   int      `json:"rebalances"`
	MedianCost   float64  `json:"medianCost"`
	CostPerYear  float64  `json:"costPerYear"`
	MedianDrift  float64  `json:"medianDrift"`
	LargestDrift float64  `json:"largestDrift"`
}

// Rebalance measures the basket of components over the market's last ten years.
func Rebalance(m Market, components []string) Basket {
	b := Basket{Components: components}
	index := make([]int, len(components))
	for i, c := range components {
		for k, n := range m.Names {
			if n == c {
				index[i] = k
			}
		}
	}
	start := len(m.At) - 1
	for start > 0 && m.At[len(m.At)-1].Sub(m.At[start-1]) <= 10*365*24*time.Hour {
		start--
	}
	b.From = date(m.At[start])
	units := make([]float64, len(components))
	reset := func(i int) float64 {
		value := 0.0
		for c, k := range index {
			value += units[c] * m.Adj[i][k]
		}
		if value == 0 {
			value = 100 * float64(len(components))
		}
		turnover := 0.0
		for c, k := range index {
			target := value / float64(len(components)) / m.Adj[i][k]
			turnover += math.Abs(target-units[c]) * m.Adj[i][k]
			units[c] = target
		}
		return turnover / value
	}
	reset(start)
	var costs, drifts []float64
	quarter := func(t time.Time) int { return t.Year()*4 + int(t.Month()-1)/3 }
	for i := start + 1; i < len(m.At); i++ {
		if quarter(m.At[i]) != quarter(m.At[i-1]) {
			costs = append(costs, reset(i)*bandFloorBps/10_000)
		}
		j := i
		for j < len(m.At) && m.At[j].Sub(m.At[i]) < BasketNoticeDays*24*time.Hour {
			j++
		}
		if j < len(m.At) {
			drifts = append(drifts, drift(m, index, i, j))
		}
	}
	b.Rebalances = len(costs)
	b.MedianCost = median(costs)
	total := 0.0
	for _, c := range costs {
		total += c
	}
	b.CostPerYear = total / (m.At[len(m.At)-1].Sub(m.At[start]).Hours() / yearHours)
	b.MedianDrift = median(drifts)
	for _, d := range drifts {
		b.LargestDrift = max(b.LargestDrift, d)
	}
	return b
}

func drift(m Market, index []int, i, j int) float64 {
	weights := func(at int) []float64 {
		w, sum := make([]float64, len(index)), 0.0
		for c, k := range index {
			w[c] = m.Adj[at][k] / m.Adj[i][k]
			sum += w[c]
		}
		for c := range w {
			w[c] /= sum
		}
		return w
	}
	a, b := weights(i), weights(j)
	d := 0.0
	for c := range a {
		d += math.Abs(a[c] - b[c])
	}
	return d / 2
}
