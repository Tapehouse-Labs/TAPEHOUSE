// SPDX-License-Identifier: MIT OR Apache-2.0

// Package backtest runs the band and the margin engine over every weekend and holiday closure of the launch assets'
// daily history and over the closures since Robinhood Chain launched, and measures what the protocol's parameters
// and constants would have done there.
package backtest

import (
	"math"
	"math/big"
	"math/rand/v2"
	"time"

	"github.com/tapehouse/tapehouse/services/internal/history"
	"github.com/tapehouse/tapehouse/services/internal/margin"
)

// Unit is the USD one unit of a portfolio's weight holds.
const Unit = 100_000

// Market is the launch assets' daily bars on the sessions every one of them traded, oldest first, indexed
// [session][asset].
type Market struct {
	Names []string
	At    []time.Time
	Open  [][]float64
	Close [][]float64
	Adj   [][]float64
}

// Align keeps the sessions every series traded.
func Align(series []history.Series) Market {
	count := map[string]int{}
	for _, s := range series {
		for _, b := range s.Bars {
			count[b.At.Format(time.DateOnly)]++
		}
	}
	m := Market{}
	for _, s := range series {
		m.Names = append(m.Names, s.Symbol)
	}
	index := map[string]int{}
	for _, b := range series[0].Bars {
		key := b.At.Format(time.DateOnly)
		if count[key] == len(series) {
			index[key] = len(m.At)
			m.At = append(m.At, b.At)
		}
	}
	for _, rows := range []*[][]float64{&m.Open, &m.Close, &m.Adj} {
		*rows = make([][]float64, len(m.At))
		for i := range *rows {
			(*rows)[i] = make([]float64, len(series))
		}
	}
	for k, s := range series {
		for _, b := range s.Bars {
			if i, ok := index[b.At.Format(time.DateOnly)]; ok {
				m.Open[i][k], m.Close[i][k], m.Adj[i][k] = b.Open, b.Close, b.AdjClose
			}
		}
	}
	return m
}

// Closure is a weekend or holiday closure: the last session before it and the first after, how many hours the
// 24/5 session stays closed, each asset's move from the last close to the next open, adjusted for splits and
// dividends, and each asset's move over the eight sessions that end the session before the close.
type Closure struct {
	Close time.Time
	Open  time.Time
	Hours int
	Moves []float64
	Prior []float64
}

func days(a, b time.Time) int {
	return int(math.Round(b.Sub(a).Hours() / 24))
}

// Closures are the market's closures: consecutive sessions two calendar days or more apart.
func (m Market) Closures() []Closure {
	var out []Closure
	for i := 1; i < len(m.At); i++ {
		d := days(m.At[i-1], m.At[i])
		if d < 2 {
			continue
		}
		c := Closure{Close: m.At[i-1], Open: m.At[i], Hours: 24 * (d - 1)}
		for k := range m.Names {
			c.Moves = append(c.Moves, m.Open[i][k]*m.Adj[i][k]/m.Close[i][k]/m.Adj[i-1][k]-1)
			prior := 0.0
			if i >= 10 {
				prior = math.Abs(m.Adj[i-2][k]/m.Adj[i-10][k] - 1)
			}
			c.Prior = append(c.Prior, prior)
		}
		out = append(out, c)
	}
	return out
}

// Returns are each asset's close-to-close returns, adjusted, over sessions that follow another by one day.
func (m Market) Returns() [][]float64 {
	out := make([][]float64, 0, len(m.At))
	for i := 1; i < len(m.At); i++ {
		if days(m.At[i-1], m.At[i]) != 1 {
			continue
		}
		r := make([]float64, len(m.Names))
		for k := range r {
			r[k] = m.Adj[i][k]/m.Adj[i-1][k] - 1
		}
		out = append(out, r)
	}
	return out
}

// Portfolios draws count portfolios of whole weights from low to high over assets, from a PCG seeded with seed,
// skipping a portfolio that holds nothing.
func Portfolios(seed uint64, count, assets, low, high int) [][]int {
	rng := rand.New(rand.NewPCG(seed, seed))
	out := make([][]int, 0, count)
	for len(out) < count {
		w := make([]int, assets)
		empty := true
		for k := range w {
			w[k] = low + rng.IntN(high-low+1)
			empty = empty && w[k] == 0
		}
		if !empty {
			out = append(out, w)
		}
	}
	return out
}

// Book is a portfolio at Unit USD a weight, with the engine's exact requirements over the open market and across a
// closure and its leverage floor, as fractions of its gross exposure, before the liquidity add-on.
type Book struct {
	Weights []int
	Gross   float64
	Open    float64
	Closed  float64
	Floor   float64
	weekday float64
	across  float64
	capped  float64
}

// Books margins each portfolio with the engine's scenario set.
func Books(set *margin.Set, portfolios [][]int) []Book {
	out := make([]Book, len(portfolios))
	unit := new(big.Int).Mul(big.NewInt(Unit), big.NewInt(1e18))
	for i, w := range portfolios {
		exposures := make([]*big.Int, len(w))
		for k, x := range w {
			exposures[k] = new(big.Int).Mul(unit, big.NewInt(int64(x)))
		}
		open, closed := margin.ScenarioRequirements(set, exposures)
		floor := margin.LeverageFloor(exposures)
		gross := new(big.Int)
		for _, e := range exposures {
			gross.Add(gross, new(big.Int).Abs(e))
		}
		g := new(big.Float).SetInt(gross)
		frac := func(v *big.Int) float64 {
			f, _ := new(big.Float).Quo(new(big.Float).SetInt(v), g).Float64()
			return f
		}
		zero := new(big.Int)
		out[i] = Book{
			Weights: w,
			Gross:   float64(Unit) * float64(sumAbs(w)),
			Open:    frac(open),
			Closed:  frac(closed),
			Floor:   frac(floor),
			weekday: frac(margin.Current(open, closed, floor, margin.Regime{Kind: margin.Open})),
			across:  frac(margin.Current(open, closed, zero, margin.Regime{Kind: margin.Closed})),
			capped:  frac(margin.Current(open, closed, floor, margin.Regime{Kind: margin.Closed})),
		}
	}
	return out
}

func sumAbs(w []int) int {
	s := 0
	for _, x := range w {
		s += max(x, -x)
	}
	return s
}

// Weekday is the current requirement over the open market: the buffered open requirement.
func (b Book) Weekday() float64 {
	return b.weekday
}

// Across is the current requirement across a closure, with the weekend leverage cap's floor or without it.
func (b Book) Across(capped bool) float64 {
	if capped {
		return b.capped
	}
	return b.across
}

// Loss is what the book loses in a closure of moves, as a fraction of its gross exposure.
func (b Book) Loss(moves []float64) float64 {
	l := 0.0
	for k, w := range b.Weights {
		l -= float64(w) * moves[k]
	}
	return l / float64(sumAbs(b.Weights))
}
