// SPDX-License-Identifier: MIT OR Apache-2.0

package band

import "math/big"

// Written is a stored price: its value and its package timestamp in milliseconds.
type Written struct {
	Feed  string
	Value *big.Int
	Ms    uint64
}

// Record is the band's state rebuilt from its events: each feed's last price from PriceWritten, the variance of the
// tracked 24/7 feeds folded from the same events with the 50 s rule, SPY's anchor from Anchored, and the last
// regular close from the written market status.
type Record struct {
	Prices  map[string]Written
	Samples map[string]Sample
	Anchors map[string]Anchor
	CloseMs uint64
}

// NewRecord starts an empty record whose variance follows the feeds tracked.
func NewRecord(tracked ...string) *Record {
	r := &Record{Prices: map[string]Written{}, Samples: map[string]Sample{}, Anchors: map[string]Anchor{}}
	for _, f := range tracked {
		r.Samples[f] = Sample{Var: new(big.Int)}
	}
	return r
}

// Write applies the PriceWritten events of one writePrices call, in order.
func (r *Record) Write(prices ...Written) {
	status := false
	for _, p := range prices {
		r.Prices[p.Feed] = p
		if last, ok := r.Samples[p.Feed]; ok && p.Value.IsUint64() {
			if next, ok := NextSample(last, p.Value.Uint64(), p.Ms); ok {
				r.Samples[p.Feed] = next
			}
		}
		status = status || p.Feed == CurrentStatus
	}
	if status {
		s, known := r.Status()
		r.CloseMs = NextClose(s, known, r.CloseMs)
	}
}

// Anchored applies an Anchored event.
func (r *Record) Anchored(symbol string, anchor Anchor) {
	r.Anchors[symbol] = anchor
}

// Reanchor applies what writing asset's index price does to its anchor: the anchor moves to Chainlink's latest
// round when NextAnchor allows it, and only in NYSE regular hours where the asset's Chainlink leg follows them.
func (r *Record) Reanchor(asset Asset, cl Chainlink, nowS uint64) {
	p, ok := r.Prices[asset.Index]
	if asset.Index == "" || !ok || !p.Value.IsUint64() {
		return
	}
	if asset.RegularHours && !InRegularHours(r.Session(nowS)) {
		return
	}
	if next, ok := NextAnchor(r.Anchors[asset.Symbol], cl.Answer, cl.UpdatedAt, p.Value.Uint64(), p.Ms); ok {
		r.Anchors[asset.Symbol] = next
	}
}

// Status is the stored market status, and whether it decodes.
func (r *Record) Status() (Status, bool) {
	var values [3]*big.Int
	var ms [3]uint64
	for i, f := range []string{CurrentStatus, NextStatus, NextChangeTime} {
		p, ok := r.Prices[f]
		if !ok {
			return Status{}, false
		}
		values[i], ms[i] = p.Value, p.Ms
	}
	return DecodeStatus(values, ms)
}

// Session is the session at nowS seconds, and whether it is known.
func (r *Record) Session(nowS uint64) (Session, bool) {
	s, known := r.Status()
	return Derive(s, known, r.CloseMs, saturatingMul1000(nowS))
}

// Asset is how the band prices one asset: its 24/7 feed, or the index feed its leg is made from, and whether its
// Chainlink leg counts only in NYSE regular hours.
type Asset struct {
	Symbol       string
	Feed         string
	Index        string
	RegularHours bool
}

// Chainlink is an asset's latest Chainlink round as the band reads it: its answer, zero when absent, and when it was
// updated, in seconds.
type Chainlink struct {
	Answer    uint64
	UpdatedAt uint64
}

// Quote is the band of asset at nowS from the record, Chainlink's latest round, the Stock Token's multiplier with 18
// decimals and whether the L2 sequencer is up and settled. A corporate action or a signed halt is not replayed.
func (r *Record) Quote(asset Asset, cl Chainlink, multiplier *big.Int, sequencerUp bool, nowS uint64) Quote {
	var px, ms uint64
	feed := asset.Feed
	basis := uint64(0)
	if feed != "" {
		if p, ok := r.Prices[feed]; ok && p.Value.IsUint64() {
			px, ms = TokenPx(p.Value.Uint64(), multiplier), p.Ms
		}
	} else if asset.Index != "" {
		feed, basis = asset.Index, IndexBasisBps
		if p, ok := r.Prices[feed]; ok && p.Value.IsUint64() {
			px, ms = LegPx(r.Anchors[asset.Symbol], p.Value.Uint64()), p.Ms
		}
	}
	if px == 0 {
		ms = 0
	}
	clAt := cl.UpdatedAt
	if cl.Answer == 0 {
		clAt = 0
	}
	session, known := r.Session(nowS)
	clPx := cl.Answer
	if asset.RegularHours && !InRegularHours(session, known) {
		clPx = 0
	}
	variance := new(big.Int)
	if s, ok := r.Samples[feed]; ok {
		variance = s.Var
	}
	q := Compute(Inputs{
		Live247Px:     px,
		Live247AgeS:   PackageAgeS(nowS, ms),
		ClPx:          clPx,
		ClAgeS:        AgeS(nowS, clAt),
		ClSessionOpen: known && session.Open,
		VarCpb2:       variance,
		BasisBps:      basis,
	})
	if q.State != Halted && (!known || !sequencerUp) {
		q.State = Degraded
	}
	return q
}

func saturatingMul1000(s uint64) uint64 {
	if s > (1<<64-1)/1000 {
		return 1<<64 - 1
	}
	return s * 1000
}
