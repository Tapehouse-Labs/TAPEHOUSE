// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk"
)

// recalls moves each stock lending vault's recall queue: it settles the ticket next in turn once the vault holds
// something for it, and has the borrower buy in the latest ticket whose notice ran out. After a buy-in it liquidates
// the shorts in deficit first.
func (k *Keeper) recalls(ctx context.Context) error {
	stocks, err := k.client.Accounts().Stocks(k.opts(ctx))
	if err != nil {
		return err
	}
	for _, asset := range slices.Sorted(maps.Keys(k.client.Deployments().StockLending)) {
		i := slices.Index(stocks.Assets, asset)
		if i < 0 {
			continue
		}
		if err := k.settleHead(ctx, asset, stocks.Tokens[i]); err != nil {
			return err
		}
		bought, err := k.buyIn(ctx, asset)
		if err != nil {
			return err
		}
		if bought {
			if err := k.shorts(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

// settleHead settles the vault's ticket next in turn for the position that recalled it, found by the accounts'
// Recall. NothingToSettle is nothing to do; any other revert, as while the accounts' holding of the asset is written
// off or the issuer blocklists them, stops the queue, and is alerted rather than retried blindly.
func (k *Keeper) settleHead(ctx context.Context, asset string, token common.Address) error {
	vault := k.client.LendingVault(asset)
	queue, err := vault.Queue(k.opts(ctx))
	if err != nil || queue.Head >= queue.Tickets {
		return err
	}
	claimable, err := vault.Claimable(k.opts(ctx), queue.Head)
	if err != nil || claimable.Sign() == 0 {
		return err
	}
	symbol, err := sdk.ToBytes32(asset)
	if err != nil {
		return err
	}
	recalls, err := k.index.Events(ctx, "tapehouse.MarginAccounts", "Recall", map[string]string{
		"symbol": fmt.Sprintf("%#x", symbol), "ticket": strconv.FormatUint(queue.Head, 10),
	})
	if err != nil {
		return err
	}
	if len(recalls) == 0 {
		return fmt.Errorf("the index has no recall of %s's ticket %d yet", asset, queue.Head)
	}
	account := common.HexToAddress(recalls[0].Arg("account"))
	position, err := word(recalls[0].Arg("position"))
	if err != nil {
		return err
	}
	tx, err := k.client.Accounts().Settle(account, position, token)
	_, err = k.act(ctx, fmt.Sprintf("settle %s's ticket %d", asset, queue.Head), tx, err, "NothingToSettle")
	if _, ok := sdk.DecodeRevert(err); ok {
		k.log.Error("the recall queue is stopped; settling waits", "asset", asset, "ticket", queue.Head, "error", err)
		return nil
	}
	return err
}

// buyIn has the borrower buy in the latest ticket whose notice has run out while the vault has not been returned its
// end, which covers every ticket before it. NotDue is nothing owed; Unavailable a token paused, a vault blocklisted or
// a burn not yet synced, which the next pass retries; BuyInFailed the borrower's failure, alerted.
func (k *Keeper) buyIn(ctx context.Context, asset string) (bool, error) {
	vault := k.client.LendingVault(asset)
	queue, err := vault.Queue(k.opts(ctx))
	if err != nil {
		return false, err
	}
	now := uint64(k.now().Unix())
	for id := queue.Tickets; id > queue.Head; id-- {
		ticket, err := vault.Ticket(k.opts(ctx), id-1)
		if err != nil {
			return false, err
		}
		if ticket.DueAt == 0 || ticket.DueAt > now {
			continue
		}
		if queue.Assigned.Cmp(ticket.End) >= 0 {
			return false, nil
		}
		tx, err := vault.BuyIn(id - 1)
		bought, err := k.act(ctx, fmt.Sprintf("buyIn %s's ticket %d", asset, id-1), tx, err, "NotDue", "Unavailable")
		if reverted(err, "BuyInFailed") {
			k.log.Error("the borrower failed a buy-in", "asset", asset, "ticket", id-1, "error", err)
			return false, nil
		}
		return bought, err
	}
	return false, nil
}
