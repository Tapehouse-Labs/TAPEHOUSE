// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import "context"

// ObserveEvery is the least time between two observations of an asset, in seconds.
const ObserveEvery = 60

// observe has the margin accounts record each asset the band reports halted or with a multiplier change pending, so
// that its recovery begins, unless an observation of it is less than a minute old. An asset that fails is logged and
// the others go on.
func (k *Keeper) observe(ctx context.Context) error {
	accounts := k.client.Accounts()
	stocks, err := accounts.Stocks(k.opts(ctx))
	if err != nil {
		return err
	}
	now := uint64(k.now().Unix())
	for _, asset := range stocks.Assets {
		if err := k.observeAsset(ctx, asset, now); err != nil {
			k.log.Warn("observe failed; the next pass tries again", "asset", asset, "error", err)
		}
	}
	return nil
}

func (k *Keeper) observeAsset(ctx context.Context, asset string, now uint64) error {
	band := k.client.Band()
	quote, err := band.Quote(k.opts(ctx), asset)
	if err != nil {
		return err
	}
	if quote.State != 0 {
		action, err := band.CorporateAction(k.opts(ctx), asset)
		if err != nil || action.Status == 0 {
			return err
		}
	}
	accounts := k.client.Accounts()
	recovery, err := accounts.Recovery(k.opts(ctx), asset)
	if err != nil {
		return err
	}
	if now < recovery.UnusableAt+ObserveEvery {
		return nil
	}
	tx, err := accounts.Observe(asset)
	_, err = k.act(ctx, "observe "+asset, tx, err)
	return err
}
