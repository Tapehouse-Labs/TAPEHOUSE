// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/roundfeed"
)

// RoundFeed reads a RedStone round feed and packs the update of its rounds.
type RoundFeed struct {
	target
	c    *Client
	feed *roundfeed.RoundFeed
}

// RoundFeed returns the RedStone round feed at address, one of the registry's .roundFeeds.
func (c *Client) RoundFeed(address common.Address) *RoundFeed {
	return &RoundFeed{target: target{address: address}, c: c, feed: roundfeed.NewRoundFeed()}
}

// DataFeedID reads the RedStone data feed the feed follows.
func (f *RoundFeed) DataFeedID(opts *bind.CallOpts) ([32]byte, error) {
	return read(f.c, opts, f.target, f.feed.UnpackDataFeedId)(f.feed.TryPackDataFeedId())
}

// Update packs the feed's update with the payload source signs for the data feed it follows, feedID.
func (f *RoundFeed) Update(ctx context.Context, source PackageSource, feedID string) (Tx, error) {
	if f.err != nil {
		return Tx{}, f.err
	}
	payload, err := source.Payload(ctx, []string{feedID})
	if err != nil {
		return Tx{}, err
	}
	data, err := f.feed.TryPackUpdate()
	if err != nil {
		return Tx{}, err
	}
	return f.tx(append(data, payload...), nil)
}
