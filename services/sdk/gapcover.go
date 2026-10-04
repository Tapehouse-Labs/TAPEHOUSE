// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/gapcover"
)

// GapCover reads the gap cover and packs its transactions.
type GapCover struct {
	target
	c     *Client
	cover *gapcover.GapCover
}

// SeriesStatus is where a series stands: still to settle, settled at its reopening price, or void and refunding its
// premiums, because it could not settle on a price the band vouched for or no one settled it within a week.
type SeriesStatus uint8

// The statuses of a series.
const (
	SeriesOpen SeriesStatus = iota
	SeriesSettled
	SeriesVoid
)

// Layer is what a cover pays: the fall beyond DeductibleBps, up to LimitBps, on Notional USDG.
type Layer struct {
	Notional      *big.Int
	DeductibleBps *big.Int
	LimitBps      *big.Int
}

// CoverQuote is a cover's premium at the weekend gap of PricingGap and the band's centre now, and the writers' USDG it
// reserves.
type CoverQuote struct {
	Premium *big.Int
	Reserve *big.Int
}

// Series is an asset's series over one closure: the notional covered, where it stands, its reference and reopening
// prices in USD with 8 decimals once settled, whether it settled on the band's median centre, its reopen in
// milliseconds, zero while not recorded, and the slots its window moved because the band could not vouch for the asset.
type Series struct {
	Notional       *big.Int
	Status         SeriesStatus
	ReferencePrice uint64
	Price          uint64
	Flagged        bool
	ReopenMs       *big.Int
	Shift          uint16
}

// PricingGap is the weekend gap the cover prices an asset at, in millionths, and the realised move of the trading week
// before the close on sale that it follows: during the sales twice the move, at least 55% of the engine's gap and at
// most eight times it; outside them the engine's gap and a move of zero. A move of 2^64 - 1 is a week the feed cannot
// show, which prices at the most.
type PricingGap struct {
	Gap      *big.Int
	WeekMove *big.Int
}

// Sales is the closure on sale: the close that keys its series and when its sales end, in milliseconds.
type Sales struct {
	ClosesMs uint64
	EndsMs   uint64
}

// GapCover returns the gap cover of the registry's .tapehouse.GapCover.
func (c *Client) GapCover() *GapCover {
	return &GapCover{
		target: lookup(c.deployments.Tapehouse, "GapCover", ".tapehouse"),
		c:      c,
		cover:  gapcover.NewGapCover(),
	}
}

// Buy buys cover on asset paying layer for holder over the coming closure, for at most maxPremium, from Quote.
func (g *GapCover) Buy(asset string, layer Layer, maxPremium *big.Int, holder common.Address) (Tx, error) {
	to, symbol := g.of(asset)
	if err := present(layer.Notional, layer.DeductibleBps, layer.LimitBps, maxPremium); err != nil {
		return Tx{}, err
	}
	return to.tx(g.cover.TryPackBuy(symbol, layer.Notional, layer.DeductibleBps, layer.LimitBps, maxPremium, holder))
}

// Deposit deposits assets of the sender's USDG with the writers, for shares to receiver.
func (g *GapCover) Deposit(assets *big.Int, receiver common.Address) (Tx, error) {
	if err := present(assets); err != nil {
		return Tx{}, err
	}
	return g.tx(g.cover.TryPackDeposit(assets, receiver))
}

// Redeem redeems owner's shares of the writers' USDG to receiver, once no cover is outstanding.
func (g *GapCover) Redeem(shares *big.Int, receiver, owner common.Address) (Tx, error) {
	if err := present(shares); err != nil {
		return Tx{}, err
	}
	return g.tx(g.cover.TryPackRedeem(shares, receiver, owner))
}

// Record records the session's next close, or the reopen while it is closed. Anyone may.
func (g *GapCover) Record() (Tx, error) {
	return g.tx(g.cover.TryPackRecord())
}

// Measure keeps the realised move of the week before the close on sale for asset's series, so that its first buyer
// does not pay for reading it from the feed: during the sales, once the week's last day, a day before the close, has
// passed. Anyone may.
func (g *GapCover) Measure(asset string) (Tx, error) {
	to, symbol := g.of(asset)
	return to.tx(g.cover.TryPackMeasure(symbol))
}

// Observe records asset's band in this minute's slot of the window after the closure from closesMs. Anyone may.
func (g *GapCover) Observe(asset string, closesMs uint64) (Tx, error) {
	to, symbol := g.of(asset)
	return to.tx(g.cover.TryPackObserve(symbol, closesMs))
}

// Settle settles asset's series over the closure from closesMs, naming its feed's last rounds started before the
// sales ended and before the reopen. Anyone may.
func (g *GapCover) Settle(asset string, closesMs uint64, referenceRound, lastRound *big.Int) (Tx, error) {
	to, symbol := g.of(asset)
	if err := present(referenceRound, lastRound); err != nil {
		return Tx{}, err
	}
	return to.tx(g.cover.TryPackSettle(symbol, closesMs, referenceRound, lastRound))
}

// Void voids asset's series over the closure from closesMs once no one has settled it within a week of its close.
// Anyone may.
func (g *GapCover) Void(asset string, closesMs uint64) (Tx, error) {
	to, symbol := g.of(asset)
	return to.tx(g.cover.TryPackVoid(symbol, closesMs))
}

// Release credits cover id's payout or refund to its holder once its series has settled or is void. Anyone may.
func (g *GapCover) Release(id *big.Int) (Tx, error) {
	if err := present(id); err != nil {
		return Tx{}, err
	}
	return g.tx(g.cover.TryPackRelease(id))
}

// Claim sends receiver the USDG credited to the sender.
func (g *GapCover) Claim(receiver common.Address) (Tx, error) {
	return g.tx(g.cover.TryPackClaim(receiver))
}

// Quote reads the premium of cover on asset paying layer at the weekend gap of PricingGap and the band's centre now,
// and gives the writers' USDG it reserves.
func (g *GapCover) Quote(opts *bind.CallOpts, asset string, layer Layer) (CoverQuote, error) {
	to, symbol := g.of(asset)
	if err := present(layer.Notional, layer.DeductibleBps, layer.LimitBps); err != nil {
		return CoverQuote{}, err
	}
	premium, err := read(g.c, opts, to, g.cover.UnpackQuote)(
		g.cover.TryPackQuote(symbol, layer.Notional, layer.DeductibleBps, layer.LimitBps))
	if err != nil {
		return CoverQuote{}, err
	}
	reserve := new(big.Int).Sub(layer.LimitBps, layer.DeductibleBps)
	reserve.Mul(reserve, layer.Notional).Add(reserve, big.NewInt(9_999))
	return CoverQuote{premium, reserve.Div(reserve, bps)}, nil
}

// MinDeductible reads the smallest deductible of cover on asset now, in basis points: the start of the fitted tail at
// the weekend gap of PricingGap, plus how far the feed's last answer stands above the band's centre.
func (g *GapCover) MinDeductible(opts *bind.CallOpts, asset string) (*big.Int, error) {
	to, symbol := g.of(asset)
	return read(g.c, opts, to, g.cover.UnpackMinDeductible)(g.cover.TryPackMinDeductible(symbol))
}

// Tail reads the constants of the cover's fitted tail, in millionths: its threshold and scale, in gaps, and the share of
// falls beyond the threshold.
func (g *GapCover) Tail(opts *bind.CallOpts) (threshold, scale, probability uint32, err error) {
	for _, c := range []struct {
		into   *uint32
		pack   func() ([]byte, error)
		unpack func([]byte) (*big.Int, error)
	}{
		{&threshold, g.cover.TryPackTAILTHRESHOLDPPM, g.cover.UnpackTAILTHRESHOLDPPM},
		{&scale, g.cover.TryPackTAILSCALEPPM, g.cover.UnpackTAILSCALEPPM},
		{&probability, g.cover.TryPackTAILPROBABILITYPPM, g.cover.UnpackTAILPROBABILITYPPM},
	} {
		value, err := read(g.c, opts, g.target, c.unpack)(c.pack())
		if err != nil {
			return 0, 0, 0, err
		}
		*c.into = uint32(value.Uint64())
	}
	return threshold, scale, probability, nil
}

// PricingGap reads the weekend gap the cover prices asset at now and the week's realised move it follows.
func (g *GapCover) PricingGap(opts *bind.CallOpts, asset string) (PricingGap, error) {
	to, symbol := g.of(asset)
	out, err := read(g.c, opts, to, g.cover.UnpackPricingGap)(g.cover.TryPackPricingGap(symbol))
	return PricingGap{out.Gap, out.WeekMove}, err
}

// Capacity reads the writers' USDG not reserved for a cover: the most a new cover may reserve.
func (g *GapCover) Capacity(opts *bind.CallOpts) (*big.Int, error) {
	return read(g.c, opts, g.target, g.cover.UnpackCapacity)(g.cover.TryPackCapacity())
}

// LastCloseMs reads the close that keys the series of the coming or the last recorded closure, in milliseconds.
func (g *GapCover) LastCloseMs(opts *bind.CallOpts) (uint64, error) {
	return read(g.c, opts, g.target, g.cover.UnpackLastCloseMs)(g.cover.TryPackLastCloseMs())
}

// Sales reads the closure on sale now; false while no cover is sold.
func (g *GapCover) Sales(opts *bind.CallOpts) (Sales, bool, error) {
	out, err := read(g.c, opts, g.target, g.cover.UnpackSales)(g.cover.TryPackSales())
	return Sales{out.ClosesMs, out.EndsMs}, err == nil && out.EndsMs != 0, err
}

// Series reads asset's series over the closure from closesMs, and its reopen, at one block.
func (g *GapCover) Series(opts *bind.CallOpts, asset string, closesMs uint64) (Series, error) {
	to, symbol := g.of(asset)
	if to.err != nil {
		return Series{}, to.err
	}
	pinned, err := g.c.pin(opts)
	if err != nil {
		return Series{}, err
	}
	out, err := read(g.c, pinned, to, g.cover.UnpackSeries)(g.cover.TryPackSeries(symbol, closesMs))
	if err != nil {
		return Series{}, err
	}
	if out.Status > uint8(SeriesVoid) {
		return Series{}, fmt.Errorf("unknown series status %d", out.Status)
	}
	reopen, err := read(g.c, pinned, to, g.cover.UnpackReopenOf)(g.cover.TryPackReopenOf(closesMs))
	if err != nil {
		return Series{}, err
	}
	return Series{out.Notional, SeriesStatus(out.Status), out.ReferencePrice, out.Price, out.Flagged, reopen, out.Shift}, nil
}

// Feed reads the Chainlink feed that settles asset's series; zero where it may not be covered.
func (g *GapCover) Feed(opts *bind.CallOpts, asset string) (common.Address, error) {
	to, symbol := g.of(asset)
	return read(g.c, opts, to, g.cover.UnpackFeed)(g.cover.TryPackFeed(symbol))
}

// RegularHours reads whether the cover's Chainlink feeds follow NYSE regular hours, so that its sales end at the
// regular close and its series reopen at the regular open.
func (g *GapCover) RegularHours(opts *bind.CallOpts) (bool, error) {
	return read(g.c, opts, g.target, g.cover.UnpackRegularHours)(g.cover.TryPackRegularHours())
}

// ReopenOf reads when the closure keyed by closesMs reopens, in milliseconds: the recorded 24/5 reopen, or the regular
// open after it where the feeds follow regular hours; zero while not recorded.
func (g *GapCover) ReopenOf(opts *bind.CallOpts, closesMs uint64) (*big.Int, error) {
	return read(g.c, opts, g.target, g.cover.UnpackReopenOf)(g.cover.TryPackReopenOf(closesMs))
}

// Cover reads cover id: its holder, closure, asset's symbol, layer and premium; false once released.
func (g *GapCover) Cover(opts *bind.CallOpts, id *big.Int) (gapcover.CoversOutput, bool, error) {
	if err := present(id); err != nil {
		return gapcover.CoversOutput{}, false, err
	}
	out, err := read(g.c, opts, g.target, g.cover.UnpackCovers)(g.cover.TryPackCovers(id))
	return out, err == nil && out.Holder != (common.Address{}), err
}

// Payouts reads the USDG credited to holder and not yet claimed.
func (g *GapCover) Payouts(opts *bind.CallOpts, holder common.Address) (*big.Int, error) {
	return read(g.c, opts, g.target, g.cover.UnpackPayouts)(g.cover.TryPackPayouts(holder))
}

// Payout is what a cover paying layer pays once its series settled at price against referencePrice: its notional
// times the fall beyond its deductible, up to its limit, before any shortfall a wipe leaves.
func Payout(layer Layer, referencePrice, price uint64) *big.Int {
	if price >= referencePrice {
		return new(big.Int)
	}
	reference := new(big.Int).SetUint64(referencePrice)
	fall := new(big.Int).Mul(new(big.Int).SetUint64(referencePrice-price), bps)
	if limit := new(big.Int).Mul(reference, layer.LimitBps); fall.Cmp(limit) > 0 {
		fall = limit
	}
	deductible := new(big.Int).Mul(reference, layer.DeductibleBps)
	if fall.Cmp(deductible) <= 0 {
		return new(big.Int)
	}
	paid := new(big.Int).Mul(layer.Notional, fall.Sub(fall, deductible))
	return paid.Div(paid, reference.Mul(reference, bps))
}

func (g *GapCover) of(asset string) (target, [32]byte) {
	symbol, err := ToBytes32(asset)
	return g.with(err), symbol
}
