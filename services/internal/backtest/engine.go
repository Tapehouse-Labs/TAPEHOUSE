// SPDX-License-Identifier: MIT OR Apache-2.0

package backtest

import (
	"math"
	"slices"
	"time"

	"github.com/tapehouse/tapehouse/services/internal/margin"
)

// Coverage is how the requirement in force covered the realised closures: the buffered open requirement against
// the closed one and against each portfolio's worst closure, and the ramp's level at its start, half-way and its
// end against every closure's loss.
type Coverage struct {
	Portfolios         int        `json:"portfolios"`
	BufferCoversClosed int        `json:"bufferCoversClosed"`
	CoveredAtWorst     int        `json:"coveredAtWorst"`
	LargestShortfall   float64    `json:"largestShortfall"`
	Ramp               [3]float64 `json:"ramp"`
}

// Cover measures the buffer and the ramp over closures.
func Cover(books []Book, closures []Closure) Coverage {
	c := Coverage{Portfolios: len(books)}
	var ramp [3]int
	for _, b := range books {
		if b.Weekday() >= b.Closed {
			c.BufferCoversClosed++
		}
		worst := math.Inf(-1)
		for _, cl := range closures {
			loss := b.Loss(cl.Moves)
			worst = max(worst, loss)
			for i, f := range []float64{0, 0.5, 1} {
				if loss <= b.Weekday()+(b.Across(true)-b.Weekday())*f {
					ramp[i]++
				}
			}
		}
		if worst <= b.Weekday() {
			c.CoveredAtWorst++
		}
		c.LargestShortfall = max(c.LargestShortfall, worst-b.Weekday())
	}
	for i := range ramp {
		c.Ramp[i] = float64(ramp[i]) / float64(len(books)*len(closures))
	}
	return c
}

// Horizon is how a closure's tail and variance compare with the engine's diffusion: the closure's ES99, the larger
// of its falls' and its rises', squared in days of a normal diffusion at the asset's daily volatility, and its
// variance in days.
type Horizon struct {
	Asset    string  `json:"asset"`
	Days     float64 `json:"days"`
	Variance float64 `json:"variance"`
}

func shortfall(xs []float64) float64 {
	a := slices.Clone(xs)
	slices.Sort(a)
	n := int(math.Ceil(float64(len(a)) / 100))
	down, up := 0.0, 0.0
	for i := range n {
		down, up = down-a[i], up+a[len(a)-1-i]
	}
	return max(down, up) / float64(n)
}

func variance(xs []float64) float64 {
	mean := 0.0
	for _, x := range xs {
		mean += x
	}
	mean /= float64(len(xs))
	v := 0.0
	for _, x := range xs {
		v += (x - mean) * (x - mean)
	}
	return v / float64(len(xs)-1)
}

func column(rows [][]float64, k int) []float64 {
	out := make([]float64, len(rows))
	for i, r := range rows {
		out[i] = r[k]
	}
	return out
}

func moves(closures []Closure) [][]float64 {
	out := make([][]float64, len(closures))
	for i, c := range closures {
		out[i] = c.Moves
	}
	return out
}

// Horizons measures each asset's closures against a diffusion at the volatility of its sessions that follow
// another by one day.
func Horizons(names []string, closures []Closure, returns [][]float64) []Horizon {
	out := make([]Horizon, len(names))
	for k, n := range names {
		c := column(moves(closures), k)
		v := variance(column(returns, k))
		ratio := shortfall(c) / (margin.ZMarket / 1e6 * math.Sqrt(v))
		out[k] = Horizon{Asset: n, Days: ratio * ratio, Variance: variance(c) / v}
	}
	return out
}

// Worst is the largest closure move of an asset in either direction, and when it reopened.
type Worst struct {
	Asset string  `json:"asset"`
	Move  float64 `json:"move"`
	Date  string  `json:"date"`
}

// Breaches are the portfolios that lose more than their margin in some closure, and those closures.
type Breaches struct {
	Portfolios int      `json:"portfolios"`
	Count      int      `json:"count"`
	Closures   []string `json:"closures"`
}

// Leverage is the weekend leverage cap recalibrated: each asset's worst closure, the cap the rule gives, the
// portfolios that would lose more than their margin at the model's own weekend limit, at the rule's cap and at the
// deployed cap, and the capacity the cap costs long-only portfolios.
type Leverage struct {
	Worst          []Worst             `json:"worst"`
	Cap            int                 `json:"cap"`
	DeployedCap    int                 `json:"deployedCap"`
	Model          map[string]Breaches `json:"model"`
	AtCap          map[string]Breaches `json:"atCap"`
	AtDeployedCap  map[string]Breaches `json:"atDeployedCap"`
	CapacityModel  float64             `json:"capacityModel"`
	CapacityCapped float64             `json:"capacityCapped"`
	CapacityLost   float64             `json:"capacityLost"`
	CapacityAtCap  float64             `json:"capacityAtCap"`
}

// Cap is the weekend leverage cap's rule: the largest whole L with 1/L at least 1.25 times the worst move.
func Cap(worst float64) int {
	return int(math.Floor(1/(1.25*worst) + 1e-9))
}

func date(t time.Time) string {
	return t.Format(time.DateOnly)
}

func breaches(books []Book, closures []Closure, need func(Book) float64) Breaches {
	b := Breaches{Portfolios: len(books)}
	seen := map[string]bool{}
	for _, book := range books {
		hit := false
		for _, c := range closures {
			if book.Loss(c.Moves) > need(book) {
				hit = true
				if !seen[date(c.Open)] {
					seen[date(c.Open)] = true
					b.Closures = append(b.Closures, date(c.Open))
				}
			}
		}
		if hit {
			b.Count++
		}
	}
	slices.Sort(b.Closures)
	return b
}

func median(xs []float64) float64 {
	s := slices.Clone(xs)
	slices.Sort(s)
	n := len(s)
	if n%2 == 1 {
		return s[n/2]
	}
	return (s[n/2-1] + s[n/2]) / 2
}

// Recalibrate repeats the weekend leverage cap's calibration over closures for long-only and mixed books.
func Recalibrate(names []string, closures []Closure, books map[string][]Book) Leverage {
	l := Leverage{DeployedCap: margin.WeekendLeverage / 10_000}
	all := Worst{}
	for k, n := range names {
		w := Worst{Asset: n}
		for _, c := range closures {
			if math.Abs(c.Moves[k]) > math.Abs(w.Move) {
				w.Move, w.Date = c.Moves[k], date(c.Open)
			}
		}
		l.Worst = append(l.Worst, w)
		if math.Abs(w.Move) > math.Abs(all.Move) {
			all = w
		}
	}
	l.Cap = Cap(math.Abs(all.Move))
	l.Model, l.AtCap, l.AtDeployedCap = map[string]Breaches{}, map[string]Breaches{}, map[string]Breaches{}
	for kind, bs := range books {
		l.Model[kind] = breaches(bs, closures, func(b Book) float64 { return b.Across(false) })
		l.AtCap[kind] = breaches(bs, closures, func(b Book) float64 { return max(b.Across(false), 1/float64(l.Cap)) })
		l.AtDeployedCap[kind] = breaches(bs, closures, func(b Book) float64 { return b.Across(true) })
	}
	var model, capped, lost, atCap []float64
	for _, b := range books["long"] {
		model = append(model, 1-b.Across(false))
		capped = append(capped, 1-b.Across(true))
		lost = append(lost, 1-(1-b.Across(true))/(1-b.Across(false)))
		atCap = append(atCap, 1-max(b.Across(false), 1/float64(l.Cap)))
	}
	l.CapacityModel, l.CapacityCapped, l.CapacityLost = median(model), median(capped), median(lost)
	l.CapacityAtCap = median(atCap)
	return l
}

// Deleveraging is what an account at its weekday limit cuts before the close to stay within the weekend cap: the
// median share of its gross exposure, long-only and mixed, and, with every account at an asset's cap of its selling
// depth, the price move the long-only cut makes in the pool at the engine's linear impact.
type Deleveraging struct {
	Long  float64 `json:"long"`
	Mixed float64 `json:"mixed"`
	Move  float64 `json:"move"`
}

// Delever measures the cut before the close.
func Delever(books map[string][]Book) Deleveraging {
	cut := func(bs []Book) float64 {
		var c []float64
		for _, b := range bs {
			c = append(c, max(0, 1-b.Weekday()/b.Floor))
		}
		return median(c)
	}
	d := Deleveraging{Long: cut(books["long"]), Mixed: cut(books["mixed"])}
	d.Move = d.Long * 0.1
	return d
}

// Scenario is a historical joint gap weighed against the weekend cap: the closure's moves, the portfolios the
// scenario would bind beyond the model's own requirement, the long-only capacity under the scenario and under the
// cap, and the median requirement each adds to mixed portfolios, hedges included.
type Scenario struct {
	Date             string             `json:"date"`
	Moves            map[string]float64 `json:"moves"`
	Binds            map[string]int     `json:"binds"`
	CapacityScenario float64            `json:"capacityScenario"`
	CapacityCap      float64            `json:"capacityCap"`
	HedgeTaxScenario float64            `json:"hedgeTaxScenario"`
	HedgeTaxCap      float64            `json:"hedgeTaxCap"`
}

// Replay weighs the closure that reopened on day as a scenario row against the cap; ok is false without it.
func Replay(names []string, closures []Closure, books map[string][]Book, day string) (Scenario, bool) {
	i := slices.IndexFunc(closures, func(c Closure) bool { return date(c.Open) == day })
	if i < 0 {
		return Scenario{}, false
	}
	c := closures[i]
	s := Scenario{Date: day, Moves: map[string]float64{}, Binds: map[string]int{}}
	for k, n := range names {
		s.Moves[n] = c.Moves[k]
	}
	for kind, bs := range books {
		var scenario, capped, taxS, taxC []float64
		for _, b := range bs {
			need := max(b.Across(false), b.Loss(c.Moves))
			if b.Loss(c.Moves) > b.Across(false) {
				s.Binds[kind]++
			}
			scenario, capped = append(scenario, 1-need), append(capped, 1-b.Across(true))
			taxS, taxC = append(taxS, need-b.Across(false)), append(taxC, b.Across(true)-b.Across(false))
		}
		if kind == "long" {
			s.CapacityScenario, s.CapacityCap = median(scenario), median(capped)
		} else {
			s.HedgeTaxScenario, s.HedgeTaxCap = median(taxS), median(taxC)
		}
	}
	return s, true
}

// FlatRequirement is the flat rule the capacity is compared with: Regulation T's 50% of gross exposure.
const FlatRequirement = 0.5

// Capacity is the capacity gain over the flat rule and the bad-debt rate of each: the median long-only share of
// gross exposure a requirement leaves to borrow against across a closure, and the share of portfolio-closures,
// long-only and mixed, whose realised loss exceeded the requirement.
type Capacity struct {
	Method       string  `json:"method"`
	Model        float64 `json:"model"`
	Flat         float64 `json:"flat"`
	BadDebtModel float64 `json:"badDebtModel"`
	BadDebtFlat  float64 `json:"badDebtFlat"`
}

// Compare measures the engine with its weekend cap against the flat rule.
func Compare(closures []Closure, books map[string][]Book) Capacity {
	c := Capacity{
		Method: "Requirement across a closure as the margin accounts hold it (the larger of the buffered open and the " +
			"closed requirement, at least gross over the weekend cap), before the liquidity add-on, for portfolios at " +
			"100,000 USD a weight, against Regulation T's 50%; a bad debt is a closure loss, from the last close to the " +
			"next open, larger than the requirement.",
		Flat: 1 - FlatRequirement,
	}
	var capacity []float64
	pairs, model, flat := 0, 0, 0
	for _, kind := range []string{"long", "mixed"} {
		for _, b := range books[kind] {
			if kind == "long" {
				capacity = append(capacity, 1-b.Across(true))
			}
			for _, cl := range closures {
				pairs++
				loss := b.Loss(cl.Moves)
				if loss > b.Across(true) {
					model++
				}
				if loss > FlatRequirement {
					flat++
				}
			}
		}
	}
	c.Model = median(capacity)
	c.BadDebtModel, c.BadDebtFlat = float64(model)/float64(pairs), float64(flat)/float64(pairs)
	return c
}
