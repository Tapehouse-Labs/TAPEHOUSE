// SPDX-License-Identifier: MIT OR Apache-2.0

package margin

import "math/big"

// Reference floor moves, in millionths: FINRA 4210(g)'s ±15% for a stock, −8% and +6% for the market asset.
const (
	StockMove = 150_000
	IndexDown = 80_000
	IndexUp   = 60_000
	tail      = 3
)

var (
	bigPPM      = big.NewInt(ppm)
	bigWad      = big.NewInt(wad)
	maxExposure = new(big.Int).Mul(big.NewInt(10_000_000_000), bigWad)
)

// PositiveDefinite reports whether the symmetric matrix with a unit diagonal and the upper triangle upper, in basis
// points, is positive definite: Sylvester's criterion through Bareiss's fraction-free elimination.
func PositiveDefinite(n int, upper []uint16) bool {
	a := make([]*big.Int, n*n)
	for i := range n {
		a[i*n+i] = big.NewInt(One)
		for j := i + 1; j < n; j++ {
			a[i*n+j] = big.NewInt(int64(upper[Pair(n, i, j)]))
			a[j*n+i] = a[i*n+j]
		}
	}
	previous := big.NewInt(1)
	for k := range n {
		pivot := a[k*n+k]
		if pivot.Sign() <= 0 {
			return false
		}
		for i := k + 1; i < n; i++ {
			for j := k + 1; j < n; j++ {
				v := new(big.Int).Mul(a[i*n+j], pivot)
				a[i*n+j] = v.Sub(v, new(big.Int).Mul(a[i*n+k], a[k*n+j])).Quo(v, previous)
			}
		}
		previous = pivot
	}
	return true
}

type worst [tail]*big.Int

func newWorst() *worst {
	return &worst{}
}

func (w *worst) push(loss *big.Int) {
	for i := range w {
		if loss == nil {
			return
		}
		if w[i] == nil || loss.Cmp(w[i]) > 0 {
			w[i], loss = loss, w[i]
		}
	}
}

func (w *worst) shortfall(size int) *big.Int {
	whole, part := size/100, int64(size%100)
	sum := new(big.Int)
	for _, l := range w[:whole] {
		sum.Add(sum, l)
	}
	sum.Mul(sum, big.NewInt(100)).Add(sum, new(big.Int).Mul(big.NewInt(part), w[whole]))
	return sum.Quo(sum, big.NewInt(int64(size)))
}

// Exposures are each position's value: its signed quantity with 18 decimals at its price with 8, rounded toward
// zero. The second result is the first asset whose value exceeds 10 billion USD, or -1.
func Exposures(quantities, prices []*big.Int) ([]*big.Int, int) {
	out := make([]*big.Int, len(quantities))
	for i, q := range quantities {
		v := new(big.Int).Mul(q, prices[i])
		out[i] = v.Quo(v, big.NewInt(100_000_000))
		if new(big.Int).Abs(v).Cmp(maxExposure) > 0 {
			return nil, i
		}
	}
	return out, -1
}

// Loss is the loss of exposures in a scenario of returns.
func Loss(exposures []*big.Int, returns []int32) *big.Int {
	sum, term := new(big.Int), new(big.Int)
	for i, e := range exposures {
		term.Mul(e, big.NewInt(int64(returns[i])))
		sum.Sub(sum, term.Quo(term, bigPPM))
	}
	return sum
}

// ScenarioRequirements is the requirement from the scenarios of set over the open market and across a closure.
func ScenarioRequirements(set *Set, exposures []*big.Int) (open, closed *big.Int) {
	size, n := set.Size, len(exposures)
	joint, apart := newWorst(), newWorst()
	alone := make([]*worst, n)
	for i := range alone {
		alone[i] = newWorst()
	}
	stress, gaps := new(big.Int), new(big.Int)
	for index := range Count(size, n) {
		returns := set.Row(index)
		switch {
		case index < size:
			joint.push(Loss(exposures, returns))
		case index < 2*size:
			for i, w := range alone {
				l := new(big.Int).Mul(exposures[i], big.NewInt(int64(returns[i])))
				w.push(l.Neg(l).Quo(l, bigPPM))
			}
		case index < 3*size:
			apart.push(Loss(exposures, returns))
		case index < 3*size+2:
			if l := Loss(exposures, returns); l.Cmp(stress) > 0 {
				stress = l
			}
		default:
			if l := Loss(exposures, returns); l.Cmp(gaps) > 0 {
				gaps = l
			}
		}
	}
	dependent := joint.shortfall(size)
	if a := apart.shortfall(size); a.Cmp(dependent) > 0 {
		dependent = a
	}
	standalone := new(big.Int)
	for _, w := range alone {
		standalone.Add(standalone, w.shortfall(size))
	}
	excess := standalone.Sub(standalone, dependent)
	if excess.Sign() < 0 {
		excess.SetInt64(0)
	}
	open = excess.Quo(excess, big.NewInt(5)).Add(excess, dependent)
	for i, e := range exposures {
		adverse := int64(StockMove)
		if i == set.Parameters.Market {
			adverse = IndexUp
			if e.Sign() > 0 {
				adverse = IndexDown
			}
		}
		f := new(big.Int).Abs(e)
		if f.Mul(f, big.NewInt(adverse)).Quo(f, bigPPM); f.Cmp(open) > 0 {
			open = f
		}
	}
	for _, v := range []*big.Int{stress, new(big.Int)} {
		if v.Cmp(open) > 0 {
			open = v
		}
	}
	closed = open
	if gaps.Cmp(closed) > 0 {
		closed = gaps
	}
	return open, closed
}

// PoolKind is what an asset's pool gives the liquidity add-on.
type PoolKind uint8

// The pool kinds: the asset has none, its oracle cannot be read, or what it says.
const (
	Absent PoolKind = iota
	Unread
	Read
)

// PoolTerms are what a pool says about liquidating an asset: its time-weighted price in USD with 8 decimals, the USD
// with 18 decimals it pays out for a 10% fall and takes in for a 10% rise, and its fee in millionths.
type PoolTerms struct {
	Price   *big.Int
	Selling *big.Int
	Buying  *big.Int
	Fee     uint32
}

// Pool is an asset's pool as the liquidity add-on sees it; an unread pool keeps only its fee.
type Pool struct {
	Kind  PoolKind
	Terms PoolTerms
}

// Liquidity is the add-on of every position on the side a liquidation takes, and a bit for every asset held without
// a read pool. depths holds each asset's selling then buying depth in whole USD.
func Liquidity(exposures, prices []*big.Int, gaps, depths []uint32, pools []Pool) (*big.Int, uint8) {
	total := new(big.Int)
	var missing uint8
	for i, e := range exposures {
		if e.Sign() == 0 {
			continue
		}
		long := e.Sign() > 0
		side := 2*i + 1
		if long {
			side = 2 * i
		}
		governed := new(big.Int).Mul(big.NewInt(int64(depths[side])), bigWad)
		price, depth, fee := prices[i], governed, pools[i].Terms.Fee
		switch pools[i].Kind {
		case Absent:
			missing |= 1 << i
			continue
		case Unread:
			missing |= 1 << i
		case Read:
			t := pools[i].Terms
			price = t.Price
			own := t.Buying
			if long {
				own = t.Selling
			}
			if half := new(big.Int).Rsh(governed, 1); own.Cmp(half) < 0 {
				own = half
			}
			if own.Cmp(governed) < 0 {
				depth = own
			}
		}
		total.Add(total, LiquidityAddon(new(big.Int).Abs(e), long, prices[i], price, depth, fee, gaps[i]))
	}
	return total, missing
}

// LiquidityAddon is what liquidating one position costs beyond its value: the pool's fee, its discount to the
// valuation price when it works against the position, at most the asset's weekend gap, and the price impact, linear
// to a 10% move at depth and in full beyond it.
func LiquidityAddon(value *big.Int, long bool, price, poolPrice, depth *big.Int, fee, gap uint32) *big.Int {
	adverse := new(big.Int).Sub(price, poolPrice)
	if !long {
		adverse.Neg(adverse)
	}
	if adverse.Sign() < 0 {
		adverse.SetInt64(0)
	}
	if bound := new(big.Int).Mul(price, big.NewInt(int64(gap))); bound.Quo(bound, bigPPM).Cmp(adverse) < 0 {
		adverse = bound
	}
	out := new(big.Int)
	if price.Sign() != 0 {
		out.Mul(value, adverse).Quo(out, price)
	}
	f := new(big.Int).Mul(value, big.NewInt(int64(fee)))
	out.Add(out, f.Quo(f, bigPPM))
	impact := new(big.Int)
	switch {
	case value.Cmp(depth) > 0:
		impact.Quo(depth, big.NewInt(20)).Add(impact, new(big.Int).Sub(value, depth))
	case depth.Sign() != 0:
		impact.Mul(value, value).Quo(impact, new(big.Int).Mul(big.NewInt(20), depth))
	}
	return out.Add(out, impact)
}

// Requirements are the requirement over the open market and across a closure, the leverage floor with the same
// liquidity add-on, and the missing bits.
func Requirements(set *Set, exposures, prices []*big.Int, depths []uint32, pools []Pool) (open, closed, floor *big.Int, missing uint8) {
	open, closed = ScenarioRequirements(set, exposures)
	addon, missing := Liquidity(exposures, prices, set.Parameters.Gaps, depths, pools)
	floor = LeverageFloor(exposures)
	return new(big.Int).Add(open, addon), new(big.Int).Add(closed, addon), floor.Add(floor, addon), missing
}

// Requirement is the requirement of exposures valued at prices, over the open market or across a closure.
func Requirement(set *Set, exposures, prices []*big.Int, spansClosure bool, depths []uint32, pools []Pool) (*big.Int, uint8) {
	open, closed, _, missing := Requirements(set, exposures, prices, depths, pools)
	if spansClosure {
		return closed, missing
	}
	return open, missing
}
