// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/reopeningauction"
)

// Chain is what the bidder reads and sends through: go-ethereum's contract and deploy backends, as ethclient has.
type Chain interface {
	bind.ContractBackend
	bind.DeployBackend
}

// Bidder bids in the reopening auction and buys in the liquidator's Dutch auctions, from the address of Opts.
type Bidder struct {
	Client      *sdk.Client
	Chain       Chain
	Opts        *bind.TransactOpts
	State       *State
	Budget      *big.Int
	DiscountBps uint64
	CommitLead  time.Duration
	Assets      []string
	DryRun      bool
	Log         *slog.Logger

	IndexerURL, IndexerKey string

	constants *sdk.AuctionConstants
	since     struct{ block, at uint64 }
}

// outcome is what a transaction came to.
type outcome int

const (
	sent   outcome = iota // mined
	retry                 // refused as not now: tried again next tick
	done                  // refused because there is nothing left to do
	dryRun                // simulated only
)

var outbid = func() [32]byte {
	parsed, err := reopeningauction.ReopeningAuctionMetaData.ParseABI()
	if err != nil {
		panic(err)
	}
	return [32]byte(parsed.Events["Outbid"].ID)
}()

// Reopening runs one tick of the reopening-auction loop at head: it commits one bid for each asset's round late in the
// commit phase, reveals each committed bid in its window, and claims once a round is cleared or has lapsed.
func (b *Bidder) Reopening(ctx context.Context, head *types.Header) error {
	opts := b.callOpts(ctx, head)
	auction := b.Client.ReopeningAuction()
	if b.constants == nil {
		c, err := auction.Constants(opts)
		if err != nil {
			return err
		}
		b.constants = &c
	}
	phase, err := auction.Phase(opts)
	if err != nil {
		return err
	}
	nowMs := head.Time * 1000
	var errs []error
	if phase.OpenMs != 0 && !phase.Revealing {
		errs = append(errs, b.commitAll(ctx, opts, phase.OpenMs, nowMs))
	}
	for _, bid := range b.State.Bids() {
		errs = append(errs, b.reveal(ctx, opts, bid, nowMs), b.claim(ctx, opts, bid, nowMs))
	}
	return errors.Join(errs...)
}

func (b *Bidder) callOpts(ctx context.Context, head *types.Header) *bind.CallOpts {
	return &bind.CallOpts{Context: ctx, From: b.Opts.From, BlockNumber: head.Number}
}

func (b *Bidder) assets(opts *bind.CallOpts) ([]string, error) {
	if len(b.Assets) != 0 {
		return b.Assets, nil
	}
	stocks, err := b.Client.Accounts().Stocks(opts)
	return stocks.Assets, err
}

// commitAll commits a bid in each asset's round for the regular open openMs, once within the commit lead of the
// phase's end, and sends again a commit that was written to the state and never landed.
func (b *Bidder) commitAll(ctx context.Context, opts *bind.CallOpts, openMs, nowMs uint64) error {
	end := openMs - b.constants.RevealMs
	if nowMs+uint64(b.CommitLead.Milliseconds()) < end {
		return nil
	}
	assets, err := b.assets(opts)
	if err != nil {
		return err
	}
	var errs []error
	for _, asset := range assets {
		errs = append(errs, b.commitIn(ctx, opts, asset, openMs))
	}
	return errors.Join(errs...)
}

func (b *Bidder) commitIn(ctx context.Context, opts *bind.CallOpts, asset string, openMs uint64) error {
	for _, held := range b.State.Bids() {
		if held.Symbol != asset || held.OpenMs != openMs {
			continue
		}
		landed, err := b.Client.ReopeningAuction().Commitment(opts, asset, openMs, held.Commitment)
		if err != nil || landed.Bidder != (common.Address{}) {
			return err
		}
		return b.commit(ctx, held)
	}
	bid, ok, err := b.plan(opts, asset, openMs)
	if err != nil || !ok {
		return err
	}
	if !b.DryRun {
		if err := b.State.Add(bid); err != nil {
			return err
		}
	}
	return b.commit(ctx, bid)
}

// plan prices and sizes a bid in asset's round; false where there is none to make, with the reason logged.
func (b *Bidder) plan(opts *bind.CallOpts, asset string, openMs uint64) (Committed, bool, error) {
	round, err := b.Client.ReopeningAuction().Round(opts, asset, openMs)
	if err != nil || round.Supply.Sign() == 0 {
		return Committed{}, false, err
	}
	sealed, err := b.Client.Band().Sealed(opts, asset, round.SealMs)
	if err != nil {
		return Committed{}, false, err
	}
	now, err := b.Client.Band().Quote(opts, asset)
	if err != nil {
		return Committed{}, false, err
	}
	price, ok := OwnPrice(&sealed, now, b.DiscountBps)
	if !ok {
		b.Log.Info("no price, no bid", "asset", asset, "openMs", openMs)
		return Committed{}, false, nil
	}
	bid, err := PlanBid(price, round, b.Budget, *b.constants)
	if err != nil {
		b.Log.Warn("no bid: "+err.Error(), "asset", asset, "openMs", openMs, "price", price, "floor", round.Floor)
		return Committed{}, false, nil
	}
	committed := Committed{Symbol: asset, OpenMs: openMs, Quantity: bid.Quantity, Price: bid.Price, Deposit: bid.Deposit}
	if _, err := rand.Read(committed.Salt[:]); err != nil {
		return Committed{}, false, err
	}
	committed.Commitment, err = sdk.Commitment(b.Opts.From, asset, openMs, bid.Quantity, bid.Price, committed.Salt)
	return committed, err == nil, err
}

// commit approves the auction for exactly the bid's deposit and commits it.
func (b *Bidder) commit(ctx context.Context, bid Committed) error {
	auction := b.Client.ReopeningAuction()
	spender, err := auction.Address()
	if err != nil {
		return err
	}
	if err := b.approve(ctx, "approve the auction for "+bid.Symbol+"'s deposit", spender, bid.Deposit); err != nil {
		return err
	}
	tx, err := auction.Commit(bid.Symbol, bid.Commitment, bid.Deposit)
	if err != nil {
		return err
	}
	_, _, err = b.send(ctx, fmt.Sprintf("commit to %s's round at %d", bid.Symbol, bid.OpenMs), tx)
	return err
}

// reveal reveals bid in its window, if its commitment landed.
func (b *Bidder) reveal(ctx context.Context, opts *bind.CallOpts, bid Committed, nowMs uint64) error {
	if nowMs+b.constants.RevealMs < bid.OpenMs || nowMs >= bid.OpenMs {
		return nil
	}
	auction := b.Client.ReopeningAuction()
	pending, err := auction.Commitment(opts, bid.Symbol, bid.OpenMs, bid.Commitment)
	if err != nil || pending.Bidder == (common.Address{}) {
		return err
	}
	tx, err := auction.Reveal(bid.Symbol, bid.OpenMs, bid.Quantity, bid.Price, bid.Salt)
	if err != nil {
		return err
	}
	name := fmt.Sprintf("reveal %s's bid at %d", bid.Symbol, bid.OpenMs)
	result, receipt, err := b.send(ctx, name, tx)
	if err != nil || result != sent {
		return err
	}
	for _, log := range receipt.Logs {
		if len(log.Topics) != 0 && [32]byte(log.Topics[0]) == outbid {
			b.Log.Info("outbid: the deposit is returned", "asset", bid.Symbol, "openMs", bid.OpenMs)
		}
	}
	return nil
}

// claim claims each of the bot's unclaimed bids of bid's round once it is cleared or has lapsed, then forgets the round.
func (b *Bidder) claim(ctx context.Context, opts *bind.CallOpts, bid Committed, nowMs uint64) error {
	if nowMs < bid.OpenMs {
		return nil
	}
	auction := b.Client.ReopeningAuction()
	round, err := auction.Round(opts, bid.Symbol, bid.OpenMs)
	if err != nil {
		return err
	}
	if !round.Cleared && nowMs < bid.OpenMs+b.constants.ClearMs {
		return nil
	}
	bids, err := auction.Bids(opts, bid.Symbol, bid.OpenMs)
	if err != nil {
		return err
	}
	finished := true
	for i, revealed := range bids {
		if revealed.Bidder != b.Opts.From || revealed.Claimed {
			continue
		}
		tx, err := auction.Claim(bid.Symbol, bid.OpenMs, big.NewInt(int64(i)))
		if err != nil {
			return err
		}
		result, _, err := b.send(ctx, fmt.Sprintf("claim %s's bid %d at %d", bid.Symbol, i, bid.OpenMs), tx)
		if err != nil {
			return err
		}
		finished = finished && result == sent
	}
	if !finished || b.DryRun {
		return nil
	}
	return b.State.Remove(bid.Symbol, bid.OpenMs)
}

// approve lets spender take exactly amount of the bot's USDG.
func (b *Bidder) approve(ctx context.Context, name string, spender common.Address, amount *big.Int) error {
	token, ok := b.Client.Deployments().Tokens["USDG"]
	if !ok {
		return errors.New("the registry has no .tokens.USDG")
	}
	tx, err := b.Client.Approve(token, spender, amount)
	if err != nil {
		return err
	}
	result, _, err := b.send(ctx, name, tx)
	if err == nil && result != sent && result != dryRun {
		err = fmt.Errorf("%s: refused", name)
	}
	return err
}

// send simulates tx from the bot's address and, if it would succeed and this is not a dry run, sends it and waits
// for its receipt. A revert named in the outcomes table is not an error.
func (b *Bidder) send(ctx context.Context, name string, tx sdk.Tx) (outcome, *types.Receipt, error) {
	err := b.Client.Simulate(&bind.CallOpts{Context: ctx, From: b.Opts.From}, tx)
	if err == nil && b.DryRun {
		b.Log.Info("dry run: "+name+" would succeed", "to", tx.To)
		return dryRun, nil, nil
	}
	var receipt *types.Receipt
	if err == nil {
		opts := *b.Opts
		opts.Context = ctx
		var pending *types.Transaction
		if pending, err = b.Client.Send(&opts, tx); err == nil {
			if receipt, err = bind.WaitMined(ctx, b.Chain, pending.Hash()); err == nil {
				if receipt.Status != types.ReceiptStatusSuccessful {
					return sent, receipt, fmt.Errorf("%s reverted in block %d, transaction %s", name, receipt.BlockNumber, pending.Hash())
				}
				b.Log.Info(name, "tx", pending.Hash().Hex(), "block", receipt.BlockNumber.Uint64(), "gasUsed", receipt.GasUsed)
				return sent, receipt, nil
			}
		}
	}
	if revert, ok := sdk.DecodeRevert(err); ok {
		if b.DryRun {
			b.Log.Warn("dry run: "+name+" would revert", "reason", revert.Error())
			return dryRun, nil, nil
		}
		switch outcomes[revert.Name] {
		case retry:
			b.Log.Info(name+": not now, tried again next tick", "reason", revert.Error())
			return retry, nil, nil
		case done:
			b.Log.Info(name+": nothing left to do", "reason", revert.Error())
			return done, nil, nil
		}
	}
	return sent, nil, fmt.Errorf("%s: %w", name, err)
}
