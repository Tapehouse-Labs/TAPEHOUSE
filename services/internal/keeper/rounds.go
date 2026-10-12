// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"maps"
	"slices"

	"github.com/ethereum/go-ethereum/common"
)

// RoundInterval is how old, in seconds, a round feed's last round must be for the rounds keeper to record another.
const RoundInterval = 50

// writeRounds records the next round of each of the gap cover's feeds that is a RedStone round feed, from RedStone's
// latest package of the feed it follows, while NYSE is in regular hours: a regular-hours feed has no rounds while the
// market is closed, as a Chainlink feed of those hours has none. A feed that is not a round feed, such as Chainlink's,
// is left to its own oracle.
func (k *Keeper) writeRounds(ctx context.Context) error {
	session, err := k.client.Band().Session(k.opts(ctx))
	if err != nil || session.Nyse != 1 {
		return err
	}
	assets, err := k.bandAssets(ctx)
	if err != nil {
		return err
	}
	done := map[common.Address]bool{}
	for _, name := range slices.Sorted(maps.Keys(k.assetNames(assets))) {
		feed, err := k.client.GapCover().Feed(k.opts(ctx), name)
		if err != nil {
			return err
		}
		if feed == (common.Address{}) || done[feed] {
			continue
		}
		done[feed] = true
		round := k.client.RoundFeed(feed)
		id, err := round.DataFeedID(k.opts(ctx))
		if err != nil {
			continue
		}
		last, err := k.client.LatestRound(k.opts(ctx), feed)
		if err == nil && last.UpdatedAt.IsUint64() && k.nowMs()/1000 < last.UpdatedAt.Uint64()+RoundInterval {
			continue
		}
		tx, err := round.Update(ctx, k.prices, text(id))
		if _, err := k.act(ctx, "record "+name+"'s round", tx, err, "NotNewer"); err != nil {
			return err
		}
	}
	return nil
}

func (k *Keeper) assetNames(assets []asset) map[string]bool {
	names := map[string]bool{}
	for _, a := range assets {
		names[a.name] = true
	}
	return names
}
