// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/band"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/bandfeed"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/stocktoken"
)

// PackageSource supplies signed RedStone data packages: the integrator's own gateway client, cache or relay. Payload
// returns the packages for feedIDs serialised as RedStone's EVM connector appends them, the payload the band's
// writePrices verifies. The SDK holds no API key and calls no gateway.
type PackageSource interface {
	Payload(ctx context.Context, feedIDs []string) ([]byte, error)
}

// Round is a feed's latest round, as Chainlink's latestRoundData returns it.
type Round struct {
	RoundID         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

// Band reads the price band and each asset's BandFeed, and packs their transactions.
type Band struct {
	target
	c    *Client
	band *band.Band
	feed *bandfeed.BandFeed
}

// Band returns the band of the registry's .tapehouse.Band, with the feeds of its .bandFeeds.
func (c *Client) Band() *Band {
	return &Band{
		target: lookup(c.deployments.Tapehouse, "Band", ".tapehouse"),
		c:      c,
		band:   band.NewBand(),
		feed:   bandfeed.NewBandFeed(),
	}
}

// Quote reads the band of asset: its state (0 halted, 1 degraded, 2 closed, 3 open), live legs, centre, half-width
// and edges.
func (b *Band) Quote(opts *bind.CallOpts, asset string) (band.QuoteOutput, error) {
	symbol, err := ToBytes32(asset)
	return read(b.c, opts, b.with(err), b.band.UnpackQuote)(b.band.TryPackQuote(symbol))
}

// Session reads Chainlink's 24/5 session (0 not known, 1 closed, 2 open), NYSE's state and next state, and their
// boundaries in milliseconds.
func (b *Band) Session(opts *bind.CallOpts) (band.SessionOutput, error) {
	return read(b.c, opts, b.target, b.band.UnpackSession)(b.band.TryPackSession())
}

// Halt reads the trading halt of asset: whether a signed halt holds, until when, when it was issued, and the issuer's
// pause of the token's oracle.
func (b *Band) Halt(opts *bind.CallOpts, asset string) (band.HaltOutput, error) {
	symbol, err := ToBytes32(asset)
	return read(b.c, opts, b.with(err), b.band.UnpackHalt)(b.band.TryPackHalt(symbol))
}

// Price reads the RedStone value the band stores for feedID, its package's timestamp in milliseconds and when it was
// written.
func (b *Band) Price(opts *bind.CallOpts, feedID string) (band.PriceOutput, error) {
	id, err := ToBytes32(feedID)
	return read(b.c, opts, b.with(err), b.band.UnpackPrice)(b.band.TryPackPrice(id))
}

// WritePrices packs the band's writePrices of feedIDs, with the payload source signs for them.
func (b *Band) WritePrices(ctx context.Context, source PackageSource, feedIDs []string) (Tx, error) {
	if b.err != nil {
		return Tx{}, b.err
	}
	ids := make([][32]byte, len(feedIDs))
	for i, feedID := range feedIDs {
		id, err := ToBytes32(feedID)
		if err != nil {
			return Tx{}, err
		}
		ids[i] = id
	}
	payload, err := source.Payload(ctx, feedIDs)
	if err != nil {
		return Tx{}, err
	}
	return b.tx(b.band.TryPackWritePrices(ids, payload))
}

// Asset reads asset's Chainlink feed, RedStone feed ID, index feed ID and Stock Token; all zero for an asset the band
// does not price.
func (b *Band) Asset(opts *bind.CallOpts, asset string) (band.AssetOutput, error) {
	symbol, err := ToBytes32(asset)
	return read(b.c, opts, b.with(err), b.band.UnpackAsset)(b.band.TryPackAsset(symbol))
}

// CorporateAction reads the multiplier change of asset's Stock Token as it affects the band now: its status (0 none,
// 1 scheduled, 2 not yet confirmed by Chainlink), when it takes effect, and the multipliers before and after.
func (b *Band) CorporateAction(opts *bind.CallOpts, asset string) (band.CorporateActionOutput, error) {
	symbol, err := ToBytes32(asset)
	return read(b.c, opts, b.with(err), b.band.UnpackCorporateAction)(b.band.TryPackCorporateAction(symbol))
}

// Terms are a Stock Token's multiplier terms as ERC-8056 reports them: its multiplier now, the new one and when it
// takes effect, in seconds, and whether its issuer has paused its oracle.
type Terms struct {
	Multiplier    *big.Int
	NewMultiplier *big.Int
	EffectiveAt   *big.Int
	OraclePaused  bool
}

// Terms reads the multiplier terms of the Stock Token the band names for asset, at opts' block, or at the latest block
// where opts names none.
func (b *Band) Terms(opts *bind.CallOpts, asset string) (Terms, error) {
	pinned, err := b.c.pin(opts)
	if err != nil {
		return Terms{}, err
	}
	named, err := b.Asset(pinned, asset)
	if err != nil {
		return Terms{}, err
	}
	if named.Token == (common.Address{}) {
		return Terms{}, fmt.Errorf("the band names no Stock Token for %s", asset)
	}
	token, binding := target{address: named.Token}, stocktoken.NewStockToken()
	var terms Terms
	if terms.Multiplier, err = read(b.c, pinned, token, binding.UnpackUiMultiplier)(binding.TryPackUiMultiplier()); err != nil {
		return Terms{}, err
	}
	if terms.NewMultiplier, err = read(b.c, pinned, token, binding.UnpackNewUIMultiplier)(binding.TryPackNewUIMultiplier()); err != nil {
		return Terms{}, err
	}
	if terms.EffectiveAt, err = read(b.c, pinned, token, binding.UnpackEffectiveAt)(binding.TryPackEffectiveAt()); err != nil {
		return Terms{}, err
	}
	terms.OraclePaused, err = read(b.c, pinned, token, binding.UnpackOraclePaused)(binding.TryPackOraclePaused())
	return terms, err
}

// HaltSigner reads the address whose signed trading halts the band accepts.
func (b *Band) HaltSigner(opts *bind.CallOpts) (common.Address, error) {
	return read(b.c, opts, b.target, b.band.UnpackHaltSigner)(b.band.TryPackHaltSigner())
}

// WriteHalt packs the band's writeHalt of asset: a halt the halt signer signed, issued and expiring at the times given
// in seconds, or its lift. Anyone may send it.
func (b *Band) WriteHalt(asset string, halted bool, issuedAt, expiresAt uint64, signature []byte) (Tx, error) {
	symbol, err := ToBytes32(asset)
	return b.with(err).tx(b.band.TryPackWriteHalt(symbol, halted, issuedAt, expiresAt, signature))
}

// SyncMultiplier packs the band's syncMultiplier of asset, which records its Stock Token's multiplier change and
// confirms a material one once Chainlink prices the new terms. Anyone may send it.
func (b *Band) SyncMultiplier(asset string) (Tx, error) {
	symbol, err := ToBytes32(asset)
	return b.with(err).tx(b.band.TryPackSyncMultiplier(symbol))
}

// LatestRound reads the low side of asset's band from its BandFeed, as Chainlink's latestRoundData: its round and
// updatedAt are the block time.
func (b *Band) LatestRound(opts *bind.CallOpts, asset string) (Round, error) {
	feed := lookup(b.c.deployments.BandFeeds, asset, ".bandFeeds")
	out, err := read(b.c, opts, feed, b.feed.UnpackLatestRoundData)(b.feed.TryPackLatestRoundData())
	return Round{out.Arg0, out.Arg1, out.Arg2, out.Arg3, out.Arg4}, err
}

// LatestBand reads the whole band of asset from its BandFeed, the token's market beside it, and both halt sources.
func (b *Band) LatestBand(opts *bind.CallOpts, asset string) (bandfeed.BandFeedBand, error) {
	feed := lookup(b.c.deployments.BandFeeds, asset, ".bandFeeds")
	return read(b.c, opts, feed, b.feed.UnpackLatestBand)(b.feed.TryPackLatestBand())
}

// Sealed reads the band asset's BandFeed sealed before the reopen at reopenMs; all zero where none was sealed.
func (b *Band) Sealed(opts *bind.CallOpts, asset string, reopenMs uint64) (bandfeed.SealsOutput, error) {
	feed := lookup(b.c.deployments.BandFeeds, asset, ".bandFeeds")
	return read(b.c, opts, feed, b.feed.UnpackSeals)(b.feed.TryPackSeals(reopenMs))
}

// Seal packs asset's BandFeed.seal(), which anyone may send in the ten minutes before the session reopens.
func (b *Band) Seal(asset string) (Tx, error) {
	return lookup(b.c.deployments.BandFeeds, asset, ".bandFeeds").tx(b.feed.TryPackSeal())
}

// SequencerSettled reports whether the L2 sequencer is up and has been for over an hour; true without a feed.
func (b *Band) SequencerSettled(opts *bind.CallOpts) (bool, error) {
	return read(b.c, opts, b.target, b.band.UnpackSequencerSettled)(b.band.TryPackSequencerSettled())
}

// Variance reads the EWMA variance of the 24/7 feed feedID, in centi-basis-points squared per minute.
func (b *Band) Variance(opts *bind.CallOpts, feedID [32]byte) (*big.Int, error) {
	return read(b.c, opts, b.target, b.band.UnpackVariance)(b.band.TryPackVariance(feedID))
}
