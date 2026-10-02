// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/band"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/bandfeed"
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
