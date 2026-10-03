// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk"
)

const (
	// SlotMs is one slot of a series' settlement window, in milliseconds.
	SlotMs = 60_000
	// WindowSlots is how many slots a window holds.
	WindowSlots = 15
	// MaxDelaySlots is how many slots the band's failures may move a window by before the series is void.
	MaxDelaySlots = 15
	// SettlementDeadline is how long after its close a series still open may be voided.
	SettlementDeadline = 7 * 24 * time.Hour
	// PostMarketMs is the post-market between a regular close and the 24/5 close, in milliseconds.
	PostMarketMs = 14_400_000
	// ReopenBeforeOpenMs is the time from the 24/5 reopen to the regular open, in milliseconds.
	ReopenBeforeOpenMs = 48_600_000
)

// series is a gap cover series with covers: its asset, the close that keys it, and its covers.
type series struct {
	asset    string
	closesMs uint64
	covers   []*big.Int
}

// gapCover keeps the gap cover's series settling: it records each closure's reopen while the session is closed;
// from the reopen it observes each series with covers once a minute until its window ends, then settles it on the
// feed's rounds; it releases every cover of a series settled or void; and it voids a series still open a week after
// its close.
func (k *Keeper) gapCover(ctx context.Context) error {
	if err := k.recordReopen(ctx); err != nil {
		return err
	}
	all, err := k.allSeries(ctx)
	if err != nil {
		return err
	}
	for _, s := range all {
		if err := k.keepSeries(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

// recordReopen records the reopen of the last recorded closure once the session is closed, unless it holds it already.
func (k *Keeper) recordReopen(ctx context.Context) error {
	cover := k.client.GapCover()
	session, err := k.client.Band().Session(k.opts(ctx))
	if err != nil || session.State != sessionClosed || session.BoundaryMs == 0 {
		return err
	}
	last, err := cover.LastCloseMs(k.opts(ctx))
	if err != nil || last == 0 || last > k.nowMs() {
		return err
	}
	reopen, err := cover.ReopenOf(k.opts(ctx), last)
	if err != nil {
		return err
	}
	regular, err := cover.RegularHours(k.opts(ctx))
	if err != nil {
		return err
	}
	want := session.BoundaryMs
	if regular {
		want += ReopenBeforeOpenMs
	}
	if reopen.Uint64() == want {
		return nil
	}
	tx, err := cover.Record()
	_, err = k.act(ctx, "record the reopen", tx, err)
	return err
}

// allSeries reads every series with covers from the gap cover's Bought, with its covers.
func (k *Keeper) allSeries(ctx context.Context) ([]*series, error) {
	bought, err := k.index.Events(ctx, "tapehouse.GapCover", "Bought", nil)
	if err != nil {
		return nil, err
	}
	var all []*series
	for _, event := range bought {
		symbol, err := word(event.Arg("symbol"))
		if err != nil {
			return nil, err
		}
		closesMs, err := strconv.ParseUint(event.Arg("closesMs"), 10, 64)
		if err != nil {
			return nil, err
		}
		asset := text(symbol)
		i := slices.IndexFunc(all, func(s *series) bool { return s.asset == asset && s.closesMs == closesMs })
		if i < 0 {
			all = append(all, &series{asset: asset, closesMs: closesMs})
			i = len(all) - 1
		}
		all[i].covers = append(all[i].covers, integer(event.Arg("id")))
	}
	return all, nil
}

func (k *Keeper) keepSeries(ctx context.Context, s *series) error {
	cover := k.client.GapCover()
	read, err := cover.Series(k.opts(ctx), s.asset, s.closesMs)
	if err != nil {
		return err
	}
	name := fmt.Sprintf("%s's series from %d", s.asset, s.closesMs)
	if read.Status != sdk.SeriesOpen {
		return k.release(ctx, s)
	}
	now := k.nowMs()
	if now >= s.closesMs+uint64(SettlementDeadline.Milliseconds()) {
		tx, err := cover.Void(s.asset, s.closesMs)
		_, err = k.act(ctx, "void "+name, tx, err)
		return err
	}
	reopen := read.ReopenMs.Uint64()
	if reopen == 0 || now < reopen {
		return nil
	}
	end := reopen + uint64(read.Shift)*SlotMs
	if read.Shift <= MaxDelaySlots {
		end += WindowSlots * SlotMs
	}
	if now < end {
		minute := strconv.FormatUint((now-reopen)/SlotMs, 10)
		if k.remember("observe "+name, minute) {
			return nil
		}
		tx, err := cover.Observe(s.asset, s.closesMs)
		_, err = k.act(ctx, "observe "+name, tx, err, "WindowClosed", "NotReopened")
		return err
	}
	referenceRound, lastRound, err := k.settlementRounds(ctx, s, reopen)
	if err != nil {
		return err
	}
	tx, err := cover.Settle(s.asset, s.closesMs, referenceRound, lastRound)
	_, err = k.act(ctx, "settle "+name, tx, err, "WindowOpen")
	return err
}

// release releases every cover of a series that has settled or is void and is not yet released.
func (k *Keeper) release(ctx context.Context, s *series) error {
	cover := k.client.GapCover()
	for _, id := range s.covers {
		_, open, err := cover.Cover(k.opts(ctx), id)
		if err != nil || !open {
			if err != nil {
				return err
			}
			continue
		}
		tx, err := cover.Release(id)
		if _, err := k.act(ctx, fmt.Sprintf("release cover %s", id), tx, err, "NoCover"); err != nil {
			return err
		}
	}
	return nil
}

// settlementRounds finds a series' rounds: the feed's last round started before its sales ended, at its close as the
// last CloseRecorded of its key revised it, less the post-market where the feeds follow regular hours; and its last
// round started before the reopen.
func (k *Keeper) settlementRounds(ctx context.Context, s *series, reopen uint64) (*big.Int, *big.Int, error) {
	cover := k.client.GapCover()
	feed, err := cover.Feed(k.opts(ctx), s.asset)
	if err != nil {
		return nil, nil, err
	}
	closes, err := k.index.Events(ctx, "tapehouse.GapCover", "CloseRecorded",
		map[string]string{"closesMs": strconv.FormatUint(s.closesMs, 10)})
	if err != nil {
		return nil, nil, err
	}
	salesEnd := s.closesMs
	if len(closes) > 0 {
		if salesEnd, err = strconv.ParseUint(closes[len(closes)-1].Arg("atMs"), 10, 64); err != nil {
			return nil, nil, err
		}
	}
	regular, err := cover.RegularHours(k.opts(ctx))
	if err != nil {
		return nil, nil, err
	}
	if regular {
		salesEnd -= min(salesEnd, PostMarketMs)
	}
	referenceRound, err := k.lastRoundBefore(ctx, feed, salesEnd)
	if err != nil {
		return nil, nil, err
	}
	lastRound, err := k.lastRoundBefore(ctx, feed, reopen)
	return referenceRound, lastRound, err
}

var aggregatorMask = new(big.Int).SetUint64(^uint64(0))

// lastRoundBefore finds a Chainlink feed's last round started before atMs, searching back from its latest round by
// bisection within each phase, across a phase change where the current phase began later.
func (k *Keeper) lastRoundBefore(ctx context.Context, feed common.Address, atMs uint64) (*big.Int, error) {
	latest, err := k.client.LatestRound(k.opts(ctx), feed)
	if err != nil {
		return nil, err
	}
	current := new(big.Int).Rsh(latest.RoundID, 64).Uint64()
	last := new(big.Int).And(latest.RoundID, aggregatorMask).Uint64()
	for phase := current; ; phase-- {
		if phase != current {
			if last, err = k.lastOfPhase(ctx, feed, phase); err != nil {
				return nil, err
			}
		}
		before := func(id uint64) (bool, error) {
			round, ok, err := k.feedRound(ctx, feed, phase, id)
			return ok && round.StartedAt.Uint64()*1000 < atMs, err
		}
		found := false
		if last != 0 {
			if found, err = before(1); err != nil {
				return nil, err
			}
		}
		if found {
			lo, hi := uint64(1), last
			for lo < hi {
				mid := hi - (hi-lo)/2
				ok, err := before(mid)
				if err != nil {
					return nil, err
				}
				if ok {
					lo = mid
				} else {
					hi = mid - 1
				}
			}
			return roundID(phase, lo), nil
		}
		if phase == 0 {
			break
		}
	}
	return nil, fmt.Errorf("the feed %s has no round before %d", feed.Hex(), atMs)
}

// lastOfPhase finds the last round of a phase by doubling, then bisection; zero for a phase with none.
func (k *Keeper) lastOfPhase(ctx context.Context, feed common.Address, phase uint64) (uint64, error) {
	exists := func(id uint64) (bool, error) {
		_, ok, err := k.feedRound(ctx, feed, phase, id)
		return ok, err
	}
	if ok, err := exists(1); err != nil || !ok {
		return 0, err
	}
	lo, hi := uint64(1), uint64(2)
	for {
		ok, err := exists(hi)
		if err != nil {
			return 0, err
		}
		if !ok {
			break
		}
		lo, hi = hi, hi*2
	}
	for hi-lo > 1 {
		mid := lo + (hi-lo)/2
		ok, err := exists(mid)
		if err != nil {
			return 0, err
		}
		if ok {
			lo = mid
		} else {
			hi = mid
		}
	}
	return lo, nil
}

// feedRound reads round id of phase; false for a round the feed does not have, which it reverts on or answers empty.
func (k *Keeper) feedRound(ctx context.Context, feed common.Address, phase, id uint64) (sdk.Round, bool, error) {
	round, err := k.client.RoundData(k.opts(ctx), feed, roundID(phase, id))
	if _, reverted := sdk.DecodeRevert(err); reverted {
		return sdk.Round{}, false, nil
	}
	if err != nil {
		return sdk.Round{}, false, err
	}
	return round, round.UpdatedAt.Sign() != 0, nil
}

func roundID(phase, id uint64) *big.Int {
	out := new(big.Int).Lsh(new(big.Int).SetUint64(phase), 64)
	return out.Or(out, new(big.Int).SetUint64(id))
}
