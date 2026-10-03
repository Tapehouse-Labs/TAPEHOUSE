// SPDX-License-Identifier: MIT OR Apache-2.0

// Package band evaluates the band program's arithmetic bit for bit: the quote, the variance of the 24/7 leg, SPY's
// index leg, the token price through the ERC-8056 multiplier, and Chainlink's 24/5 session from the signed New York
// market status. Prices carry 8 decimals, widths are in basis points, and the variance is in centi-basis-points
// squared per minute.
package band

import (
	"math"
	"math/big"
)

// Constants of the band program.
const (
	FloorBps          = 30
	ZX10              = 30
	LatencyMin        = 3
	SingleSourceBps   = 25
	LiveMaxAgeS       = 120
	MaxHalfBps        = 1500
	EwmaNum           = 94
	EwmaDen           = 100
	SampleMinGapMs    = 50_000
	ClDevBps          = 50
	ClMaxAgeS         = 86_400 + 60
	SequencerGraceS   = 3_600
	IndexBasisBps     = 25
	AnchorMaxLagMs    = 120_000
	sqrtLatencyX100   = 173
	multiplierScale   = 1_000_000_000_000_000_000
	returnCapExponent = 57
)

// State is what the band reports about its legs, from the most to the least restrictive.
type State uint8

// The band's states.
const (
	Halted State = iota
	Degraded
	Closed
	Open
)

// Inputs are both legs of one asset, the Chainlink session and the variance of the 24/7 leg.
type Inputs struct {
	Live247Px     uint64
	Live247AgeS   uint64
	ClPx          uint64
	ClAgeS        uint64
	ClSessionOpen bool
	VarCpb2       *big.Int
	BasisBps      uint64
}

// Quote is the band: its state, how many legs are live, its centre, half-width and bounds.
type Quote struct {
	State   State
	Live    uint8
	Mid     uint64
	HalfBps uint64
	Low     uint64
	High    *big.Int
}

// Compute is the band program's quote. The centre is the live 24/7 leg, else Chainlink.
func Compute(in Inputs) Quote {
	live247 := in.Live247Px > 0 && in.Live247AgeS <= LiveMaxAgeS
	clLive := in.ClPx > 0 && in.ClSessionOpen && in.ClAgeS <= ClMaxAgeS
	if !live247 && !clLive {
		return Quote{High: new(big.Int)}
	}
	mid := in.ClPx
	if live247 {
		mid = in.Live247Px
	}
	vol := new(big.Int).Sqrt(in.VarCpb2)
	vol.Mul(vol, big.NewInt(ZX10*sqrtLatencyX100)).Quo(vol, big.NewInt(10*100*100))
	half := max(vol.Uint64(), FloorBps)
	if !vol.IsUint64() {
		half = MaxHalfBps
	}
	var live uint8
	if live247 {
		live++
	}
	if clLive {
		live++
	}
	if live == 2 {
		diff := in.ClPx - in.Live247Px
		if in.Live247Px > in.ClPx {
			diff = in.Live247Px - in.ClPx
		}
		d := new(big.Int).SetUint64(diff)
		d.Mul(d, big.NewInt(10_000)).Quo(d, new(big.Int).SetUint64(in.Live247Px))
		if d.Cmp(big.NewInt(ClDevBps)) > 0 {
			half = saturating(half, d.Uint64()-ClDevBps, d.IsUint64())
		}
	} else {
		half += SingleSourceBps
	}
	if live247 {
		half = saturating(half, in.BasisBps, true)
	} else {
		half += ClDevBps
	}
	half = min(half, MaxHalfBps)
	expected := uint8(1)
	if in.ClSessionOpen {
		expected = 2
	}
	state := Closed
	switch {
	case live < expected:
		state = Degraded
	case in.ClSessionOpen:
		state = Open
	}
	centre := new(big.Int).SetUint64(mid)
	low := new(big.Int).Mul(centre, big.NewInt(int64(10_000-half)))
	high := new(big.Int).Mul(centre, big.NewInt(int64(10_000+half)))
	return Quote{
		State:   state,
		Live:    live,
		Mid:     mid,
		HalfBps: half,
		Low:     low.Quo(low, big.NewInt(10_000)).Uint64(),
		High:    high.Quo(high, big.NewInt(10_000)),
	}
}

func saturating(a, b uint64, fits bool) uint64 {
	if !fits || a+b < a {
		return MaxHalfBps
	}
	return a + b
}

// EwmaUpdate folds one sample into the variance: the return from prevPx to px, normalised to one minute over dtS
// seconds. A fall rounds toward minus infinity.
func EwmaUpdate(varCpb2 *big.Int, prevPx, px, dtS uint64) *big.Int {
	if prevPx == 0 || dtS == 0 {
		return new(big.Int).Set(varCpb2)
	}
	diff := px - prevPx
	if prevPx > px {
		diff = prevPx - px
	}
	r := new(big.Int).SetUint64(diff)
	r.Mul(r, big.NewInt(1_000_000))
	prev := new(big.Int).SetUint64(prevPx)
	if px < prevPx {
		r.Add(r, prev).Sub(r, big.NewInt(1))
	}
	r.Quo(r, prev)
	if limit := new(big.Int).Lsh(big.NewInt(1), returnCapExponent); r.Cmp(limit) > 0 {
		r = limit
	}
	r.Mul(r, r).Mul(r, big.NewInt(60)).Quo(r, new(big.Int).SetUint64(dtS))
	out := new(big.Int).Mul(varCpb2, big.NewInt(EwmaNum))
	out.Add(out, r.Mul(r, big.NewInt(EwmaDen-EwmaNum)))
	return out.Quo(out, big.NewInt(EwmaDen))
}

// Sample is the variance, the sampled price and its package timestamp in milliseconds, zero before the first sample.
type Sample struct {
	Var *big.Int
	Px  uint64
	Ms  uint64
}

// NextSample is the sample once a price px with package timestamp ms is written, and whether one is taken: always
// for the first price, then once SampleMinGapMs has passed since the last. A gap never resets the variance.
func NextSample(last Sample, px, ms uint64) (Sample, bool) {
	if px == 0 || (last.Ms != 0 && ms < saturatingAdd(last.Ms, SampleMinGapMs)) {
		return last, false
	}
	return Sample{Var: EwmaUpdate(last.Var, last.Px, px, (ms-last.Ms)/1000), Px: px, Ms: ms}, true
}

// AgeS is the seconds from atS to nowS, zero when atS is ahead.
func AgeS(nowS, atS uint64) uint64 {
	if atS > nowS {
		return 0
	}
	return nowS - atS
}

// PackageAgeS is the whole seconds from a package timestamp in milliseconds to nowS, zero when the package is ahead.
func PackageAgeS(nowS, ms uint64) uint64 {
	now := uint64(math.MaxUint64)
	if nowS <= math.MaxUint64/1000 {
		now = nowS * 1000
	}
	if ms > now {
		return 0
	}
	return (now - ms) / 1000
}

func saturatingAdd(a, b uint64) uint64 {
	if a+b < a {
		return math.MaxUint64
	}
	return a + b
}

// Anchor is a Chainlink print and the index price at that print, all zero before the first.
type Anchor struct {
	ClPx    uint64
	IndexPx uint64
	AtS     uint64
}

// NextAnchor is the anchor after an index price is written at indexMs while Chainlink's latest print is clPx at
// clAtS, and whether it moved: only to a newer print of a different price, within AnchorMaxLagMs of the package.
func NextAnchor(anchor Anchor, clPx, clAtS, indexPx, indexMs uint64) (Anchor, bool) {
	at := new(big.Int).Mul(new(big.Int).SetUint64(clAtS), big.NewInt(1000))
	lag := at.Sub(at, new(big.Int).SetUint64(indexMs)).Abs(at)
	near := lag.Cmp(big.NewInt(AnchorMaxLagMs)) <= 0
	if clPx == 0 || indexPx == 0 || clAtS <= anchor.AtS || clPx == anchor.ClPx || !near {
		return anchor, false
	}
	return Anchor{ClPx: clPx, IndexPx: indexPx, AtS: clAtS}, true
}

// LegPx is the index leg's price at indexPx, rounded down; zero without an anchor or beyond u64.
func LegPx(anchor Anchor, indexPx uint64) uint64 {
	if anchor.IndexPx == 0 {
		return 0
	}
	return mulDiv(anchor.ClPx, new(big.Int).SetUint64(indexPx), new(big.Int).SetUint64(anchor.IndexPx))
}

// TokenPx is the Stock Token's price for a share price through its multiplier with 18 decimals, rounded down; zero
// beyond u64.
func TokenPx(sharePx uint64, multiplier *big.Int) uint64 {
	return mulDiv(sharePx, multiplier, big.NewInt(multiplierScale))
}

func mulDiv(a uint64, b, c *big.Int) uint64 {
	out := new(big.Int).SetUint64(a)
	out.Mul(out, b).Quo(out, c)
	if !out.IsUint64() {
		return 0
	}
	return out.Uint64()
}

// SequencerSettled reports whether an L2 sequencer-uptime round says the sequencer is up and has been for more than
// SequencerGraceS at now.
func SequencerSettled(answer, startedAt, now uint64) bool {
	return answer == 0 && startedAt != 0 && startedAt <= now && now-startedAt > SequencerGraceS
}
