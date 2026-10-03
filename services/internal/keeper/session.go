// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"time"
)

// The session's states and NYSE's, as the band's session() reports them.
const (
	sessionClosed = 1
	sessionOpen   = 2
	nyseRegular   = 1
)

const (
	// SealWindow is how long before the 24/5 reopen a BandFeed may be sealed.
	SealWindow = 10 * time.Minute
	// SealFirst is how long before the reopen the keeper seals a feed not yet sealed for it, so a keeper that then
	// fails leaves a seal.
	SealFirst = 5 * time.Minute
	// SealLast is how long before the reopen the keeper seals every feed again, as late as it safely can: the latest
	// seal before the reopen stands.
	SealLast = 90 * time.Second
)

// seal seals every BandFeed in the window before the 24/5 reopen: once SealFirst before it where none is sealed for
// it, and again SealLast before it.
func (k *Keeper) seal(ctx context.Context) error {
	band := k.client.Band()
	session, err := band.Session(k.opts(ctx))
	if err != nil {
		return err
	}
	now := k.nowMs()
	if session.State != sessionClosed || session.BoundaryMs <= now ||
		session.BoundaryMs-now > uint64(SealWindow.Milliseconds()) {
		return nil
	}
	left := session.BoundaryMs - now
	for _, asset := range slices.Sorted(maps.Keys(k.client.Deployments().BandFeeds)) {
		sealed, err := band.Sealed(k.opts(ctx), asset, session.BoundaryMs)
		if err != nil {
			return err
		}
		first := sealed.SealedAt == 0 && left <= uint64(SealFirst.Milliseconds())
		last := left <= uint64(SealLast.Milliseconds()) &&
			sealed.SealedAt*1000 < session.BoundaryMs-uint64(SealLast.Milliseconds())
		if !first && !last {
			continue
		}
		tx, err := band.Seal(asset)
		if _, err := k.act(ctx, "seal "+asset, tx, err, "NotSealWindow"); err != nil {
			return err
		}
	}
	return nil
}

// accruePremium has the margin accounts record the closure the session shows: its close while it is ahead, its
// reopening once closed, and once more in the seal window before the reopening, so a changed reopening is caught.
func (k *Keeper) accruePremium(ctx context.Context) error {
	session, err := k.client.Band().Session(k.opts(ctx))
	if err != nil {
		return err
	}
	closure, err := k.client.Accounts().Closure(k.opts(ctx))
	if err != nil {
		return err
	}
	now := k.nowMs()
	ahead := session.State == sessionOpen && session.BoundaryMs > now && closure.ClosesMs != session.BoundaryMs
	closed := session.State == sessionClosed && session.BoundaryMs != 0
	reopening := closed && closure.ReopensMs != session.BoundaryMs
	sealing := closed && session.BoundaryMs > now && session.BoundaryMs-now <= uint64(SealWindow.Milliseconds()) &&
		closure.AccruedMs < session.BoundaryMs-uint64(SealWindow.Milliseconds())
	if !ahead && !reopening && !sealing {
		return nil
	}
	tx, err := k.client.Accounts().AccruePremium()
	_, err = k.act(ctx, "accruePremium", tx, err)
	return err
}

// mark has the shorts record each shortable asset's band centre at the regular close, the close the restriction
// after a 10% fall measures from. The keeper learns the close while NYSE is in regular hours and marks once it has
// passed.
func (k *Keeper) mark(ctx context.Context) error {
	session, err := k.client.Band().Session(k.opts(ctx))
	if err != nil {
		return err
	}
	if session.Nyse == nyseRegular {
		k.remember("close", strconv.FormatUint(session.ChangeMs, 10))
		return nil
	}
	k.mu.Lock()
	closeMs := k.memory["close"]
	k.mu.Unlock()
	at, err := strconv.ParseUint(closeMs, 10, 64)
	if err != nil || at > k.nowMs() || k.remember("marked", closeMs) {
		return nil
	}
	assets, err := k.bandAssets(ctx)
	if err != nil {
		return err
	}
	shorts := k.client.Shorts()
	for _, a := range assets {
		fee, err := shorts.Fee(k.opts(ctx), a.name)
		if err != nil {
			return err
		}
		if fee.Sign() == 0 {
			continue
		}
		tx, err := shorts.Mark(a.name)
		if _, err := k.act(ctx, fmt.Sprintf("mark %s", a.name), tx, err); err != nil {
			k.remember("marked", "")
			return err
		}
	}
	return nil
}
