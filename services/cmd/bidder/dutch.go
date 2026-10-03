// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/core/types"
)

// closedAuctionLifetime is the longest an auction lives, in seconds: the liquidator's CLOSED_AUCTION_LIFETIME.
const closedAuctionLifetime = 4 * 60 * 60

// refreshSeconds is how long the block that bounds the auctions' start is kept before it is searched again.
const refreshSeconds = 10 * 60

// Dutch runs one tick of the Dutch loop at head: for each live liquidation auction it buys each Stock Token the
// position holds once the ask has fallen to the bot's own price, with the ask's cost, rounded up, as its limit.
func (b *Bidder) Dutch(ctx context.Context, head *types.Header) error {
	from, err := b.fromBlock(ctx, head)
	if err != nil {
		return err
	}
	auctions, err := liveAuctions(ctx, b.IndexerURL, b.IndexerKey, from)
	if err != nil || len(auctions) == 0 {
		return err
	}
	opts := b.callOpts(ctx, head)
	stocks, err := b.Client.Accounts().Stocks(opts)
	if err != nil {
		return err
	}
	liquidator := b.Client.Liquidator()
	spender, err := liquidator.Address()
	if err != nil {
		return err
	}
	quotes := map[string]*big.Int{}
	var errs []error
	for _, auction := range auctions {
		live, err := liquidator.Auction(opts, auction.Account, auction.Position)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if live.StartedAt == 0 {
			continue
		}
		for i, asset := range stocks.Assets {
			if len(b.Assets) != 0 && !slices.Contains(b.Assets, asset) {
				continue
			}
			token := stocks.Tokens[i]
			held, err := b.Client.Accounts().Collateral(opts, auction.Account, auction.Position, token)
			if err != nil || held.Sign() == 0 {
				errs = append(errs, err)
				continue
			}
			ask, err := liquidator.Price(opts, auction.Account, auction.Position, token)
			if err != nil || ask.Sign() == 0 {
				errs = append(errs, err)
				continue
			}
			own, known := quotes[asset]
			if !known {
				quote, err := b.Client.Band().Quote(opts, asset)
				if err != nil {
					errs = append(errs, err)
					continue
				}
				own, _ = OwnPrice(nil, quote, b.DiscountBps)
				quotes[asset] = own
			}
			if own == nil || ask.Cmp(own) > 0 {
				continue
			}
			amount := new(big.Int).Mul(b.Budget, tokenToUSDG)
			amount.Div(amount, ask)
			if amount.Sign() == 0 {
				continue
			}
			maxCost := Escrow(amount, ask)
			name := fmt.Sprintf("buy %s from %s at %s", asset, auction.Account.Hex(), ask)
			if err := b.approve(ctx, "approve the liquidator for "+name, spender, maxCost); err != nil {
				errs = append(errs, err)
				continue
			}
			tx, err := liquidator.Buy(auction.Account, auction.Position, token, amount, maxCost, b.Opts.From)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			if _, _, err := b.send(ctx, name, tx); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

// fromBlock is the first block of the last closedAuctionLifetime seconds, found by bisection and kept for a while.
func (b *Bidder) fromBlock(ctx context.Context, head *types.Header) (uint64, error) {
	if b.since.block != 0 && head.Time-b.since.at < refreshSeconds {
		return b.since.block, nil
	}
	target := head.Time - closedAuctionLifetime
	low, high := uint64(0), head.Number.Uint64()
	for low < high {
		middle := low + (high-low)/2
		header, err := b.Chain.HeaderByNumber(ctx, new(big.Int).SetUint64(middle))
		if err != nil {
			return 0, err
		}
		if header.Time < target {
			low = middle + 1
		} else {
			high = middle
		}
	}
	b.since.block, b.since.at = low, head.Time
	return low, nil
}
