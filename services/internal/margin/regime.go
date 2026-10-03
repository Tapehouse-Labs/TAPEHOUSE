// SPDX-License-Identifier: MIT OR Apache-2.0

package margin

import "math/big"

// Constants of the margin program's regimes.
const (
	// Horizon is the margin period of risk in every regime, in seconds.
	Horizon = 172_800
	// RampMs is how long before a weekend or holiday close the requirement rises to the closed one.
	RampMs = 25_200_000
	// WeekendLeverage is the most gross exposure a requirement allows per unit of margin across a closure, in basis
	// points.
	WeekendLeverage = 50_000
)

// RegimeKind is where the current requirement stands.
type RegimeKind uint8

// The regimes, numbered as currentRequirement reports them.
const (
	Unknown RegimeKind = iota
	Closed
	Open
	Closing
)

// Regime is a regime and, while closing, how far into the ramp.
type Regime struct {
	Kind      RegimeKind
	ElapsedMs uint64
}

// Code is the regime's number: 0 unknown, 1 closed, 2 open, 3 closing.
func (r Regime) Code() uint8 {
	return uint8(r.Kind)
}

// RegimeAt is the regime of the band's session, its state open (0 unknown, 1 closed, 2 open) and its boundary, at
// nowMs. An open session at or past its boundary reads as closed.
func RegimeAt(open uint8, boundaryMs, nowMs uint64) Regime {
	switch {
	case open == 1:
		return Regime{Kind: Closed}
	case open != 2:
		return Regime{Kind: Unknown}
	case boundaryMs != 0 && boundaryMs <= nowMs:
		return Regime{Kind: Closed}
	case boundaryMs != 0 && boundaryMs-nowMs < RampMs:
		return Regime{Kind: Closing, ElapsedMs: RampMs - (boundaryMs - nowMs)}
	}
	return Regime{Kind: Open}
}

// LeverageFloor is the margin that holds exposures' gross value to WeekendLeverage, rounded up.
func LeverageFloor(exposures []*big.Int) *big.Int {
	gross := new(big.Int)
	for _, e := range exposures {
		gross.Add(gross, new(big.Int).Abs(e))
	}
	return ceilDiv(gross.Mul(gross, big.NewInt(10_000)), big.NewInt(WeekendLeverage))
}

// Current is the current requirement from the requirement over the open market, across a closure, and the leverage
// floor: the buffered open requirement, rising in a straight line through the ramp to the largest of the three.
func Current(open, closed, floor *big.Int, r Regime) *big.Int {
	buffered := ceilDiv(new(big.Int).Mul(open, big.NewInt(5)), big.NewInt(4))
	across := new(big.Int).Set(buffered)
	for _, v := range []*big.Int{closed, floor} {
		if v.Cmp(across) > 0 {
			across.Set(v)
		}
	}
	switch r.Kind {
	case Open:
		return buffered
	case Closing:
		rise := new(big.Int).Sub(across, buffered)
		rise.Mul(rise, new(big.Int).SetUint64(r.ElapsedMs)).Quo(rise, big.NewInt(RampMs))
		return rise.Add(rise, buffered)
	}
	return across
}

func ceilDiv(a, b *big.Int) *big.Int {
	q, m := new(big.Int).QuoRem(a, b, new(big.Int))
	if m.Sign() != 0 {
		q.Add(q, big.NewInt(1))
	}
	return q
}
