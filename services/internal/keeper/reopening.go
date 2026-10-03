// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tapehouse/tapehouse/services/sdk"
)

// ClearWindow is how long after the regular open a reopening round may be cleared before it lapses.
const ClearWindow = time.Hour

// reopening runs the reopening auction's rounds: while bids are committed, it enrolls every position short at its
// low edges, deepest first, in the round of each Stock Token it holds; in the first hour of regular trading it clears
// each round at the price the revealed bids set; once the bids are revealed it forfeits the unrevealed commitments,
// and once the round's clearing window has closed it claims for the bidders who have not.
func (k *Keeper) reopening(ctx context.Context) error {
	auction := k.client.ReopeningAuction()
	phase, err := auction.Phase(k.opts(ctx))
	if err != nil {
		return err
	}
	if phase.OpenMs != 0 && !phase.Revealing {
		return k.enroll(ctx, phase.OpenMs)
	}
	openMs, err := auction.LastOpenMs(k.opts(ctx))
	if err != nil || openMs == 0 || k.nowMs() < openMs {
		return err
	}
	stocks, err := k.client.Accounts().Stocks(k.opts(ctx))
	if err != nil {
		return err
	}
	for _, asset := range stocks.Assets {
		if err := k.round(ctx, asset, openMs); err != nil {
			return err
		}
	}
	return nil
}

// enroll enrolls each position with debt, deepest first, in the round of each Stock Token it holds. NotEnrollable is a
// position not short at its low edges or with nothing of the asset to sell; TooManyLots a full round of larger lots.
func (k *Keeper) enroll(ctx context.Context, openMs uint64) error {
	positions, err := k.index.Debts(ctx)
	if err != nil {
		return err
	}
	deepestFirst(positions, func(p Position) (*big.Int, *big.Int) { return p.Equity, p.Requirement })
	auction := k.client.ReopeningAuction()
	for _, p := range positions {
		held, err := k.heldAssets(ctx, p.Account, p.Position)
		if err != nil {
			return err
		}
		for _, asset := range held {
			enrolled, err := auction.Enrolled(k.opts(ctx), openMs, p.Account, p.Position, asset)
			if err != nil || enrolled {
				if err != nil {
					return err
				}
				continue
			}
			tx, err := auction.Enroll(p.Account, p.Position, asset)
			name := fmt.Sprintf("enroll %s %s in %s", p.Account.Hex(), hexutil.Encode(p.Position[:]), asset)
			if _, err := k.act(ctx, name, tx, err, "NotEnrollable", "TooManyLots"); err != nil {
				return err
			}
		}
	}
	return nil
}

// heldAssets are the Stock Tokens a position holds, directly or through its baskets.
func (k *Keeper) heldAssets(ctx context.Context, account common.Address, position [32]byte) ([]string, error) {
	accounts := k.client.Accounts()
	stocks, err := accounts.Stocks(k.opts(ctx))
	if err != nil {
		return nil, err
	}
	inBaskets, err := accounts.InBaskets(k.opts(ctx), account, position)
	if err != nil {
		return nil, err
	}
	var held []string
	for i, asset := range stocks.Assets {
		if stocks.Tokens[i] == (common.Address{}) {
			continue
		}
		amount, err := accounts.Collateral(k.opts(ctx), account, position, stocks.Tokens[i])
		if err != nil {
			return nil, err
		}
		if amount.Sign() > 0 || (inBaskets[asset] != nil && inBaskets[asset].Sign() > 0) {
			held = append(held, asset)
		}
	}
	return held, nil
}

// round clears asset's round for the regular open openMs in its clearing window, then forfeits its unrevealed
// commitments and, once the window has closed, claims for each bidder who has not.
func (k *Keeper) round(ctx context.Context, asset string, openMs uint64) error {
	auction := k.client.ReopeningAuction()
	round, err := auction.Round(k.opts(ctx), asset, openMs)
	if err != nil || round.Floor.Sign() == 0 {
		return err
	}
	closed := k.nowMs() >= openMs+uint64(ClearWindow.Milliseconds())
	if !round.Cleared && !closed {
		bids, err := auction.Bids(k.opts(ctx), asset, openMs)
		if err != nil {
			return err
		}
		price, ok := sdk.ClearingPrice(bids, round.Supply, round.Floor)
		if !ok {
			return fmt.Errorf("no price clears %s's round at %d", asset, openMs)
		}
		tx, err := auction.Clear(asset, openMs, price)
		if _, err := k.act(ctx, fmt.Sprintf("clear %s at %s", asset, price), tx, err, "WrongPhase"); err != nil {
			return err
		}
	}
	if err := k.forfeit(ctx, asset, openMs); err != nil {
		return err
	}
	if !closed {
		return nil
	}
	bids, err := auction.Bids(k.opts(ctx), asset, openMs)
	if err != nil {
		return err
	}
	for i, bid := range bids {
		if bid.Claimed {
			continue
		}
		tx, err := auction.Claim(asset, openMs, big.NewInt(int64(i)))
		if _, err := k.act(ctx, fmt.Sprintf("claim %s's bid %d", asset, i), tx, err, "WrongState"); err != nil {
			return err
		}
	}
	return nil
}

// forfeit forfeits every commitment of asset's round still unrevealed.
func (k *Keeper) forfeit(ctx context.Context, asset string, openMs uint64) error {
	symbol, err := sdk.ToBytes32(asset)
	if err != nil {
		return err
	}
	committed, err := k.index.Events(ctx, "tapehouse.ReopeningAuction", "Committed", map[string]string{
		"symbol": hexutil.Encode(symbol[:]), "openMs": strconv.FormatUint(openMs, 10),
	})
	if err != nil {
		return err
	}
	auction := k.client.ReopeningAuction()
	for _, event := range committed {
		commitment, err := word(event.Arg("commitment"))
		if err != nil {
			return err
		}
		pending, err := auction.Commitment(k.opts(ctx), asset, openMs, commitment)
		if err != nil {
			return err
		}
		if pending.Bidder == (common.Address{}) {
			continue
		}
		tx, err := auction.Forfeit(asset, openMs, commitment)
		if _, err := k.act(ctx, "forfeit a commitment to "+asset, tx, err, "WrongPhase", "UnknownCommitment"); err != nil {
			return err
		}
	}
	return nil
}
