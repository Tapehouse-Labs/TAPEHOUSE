// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/gapbackstop"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/liquidator"
)

// Liquidator reads the liquidator and packs its transactions.
type Liquidator struct {
	target
	c          *Client
	liquidator *liquidator.Liquidator
}

// Liquidator returns the liquidator of the registry's .tapehouse.Liquidator.
func (c *Client) Liquidator() *Liquidator {
	return &Liquidator{
		target:     lookup(c.deployments.Tapehouse, "Liquidator", ".tapehouse"),
		c:          c,
		liquidator: liquidator.NewLiquidator(),
	}
}

// Shortfall reads account's position as the liquidator judges it now: equity and requirement in USD with 18
// decimals, whether it falls short, and whether the market is closed. A position that cannot be judged reads as not
// short.
func (l *Liquidator) Shortfall(opts *bind.CallOpts, account common.Address, position [32]byte) (liquidator.ShortfallOutput, error) {
	return read(l.c, opts, l.target, l.liquidator.UnpackShortfall)(l.liquidator.TryPackShortfall(account, position))
}

// Running reports whether an auction of account's position runs in the market's current state, all at opts' block,
// or at the latest block where opts names none.
func (l *Liquidator) Running(opts *bind.CallOpts, account common.Address, position [32]byte) (bool, error) {
	if l.err != nil {
		return false, l.err
	}
	pinned, err := l.c.pin(opts)
	if err != nil {
		return false, err
	}
	auction, err := read(l.c, pinned, l.target, l.liquidator.UnpackAuctions)(l.liquidator.TryPackAuctions(account, position))
	if err != nil || auction.StartedAt == 0 {
		return false, err
	}
	judged, err := l.Shortfall(pinned, account, position)
	if err != nil || auction.Closed != judged.Closed {
		return false, err
	}
	pack, unpack := l.liquidator.TryPackOPENAUCTIONLIFETIME, l.liquidator.UnpackOPENAUCTIONLIFETIME
	if judged.Closed {
		pack, unpack = l.liquidator.TryPackCLOSEDAUCTIONLIFETIME, l.liquidator.UnpackCLOSEDAUCTIONLIFETIME
	}
	lifetime, err := read(l.c, pinned, l.target, unpack)(pack())
	if err != nil {
		return false, err
	}
	header, err := l.c.backend.HeaderByNumber(pinned.Context, pinned.BlockNumber)
	if err != nil {
		return false, err
	}
	end := new(big.Int).Add(new(big.Int).SetUint64(auction.StartedAt), lifetime)
	return new(big.Int).SetUint64(header.Time).Cmp(end) < 0, nil
}

// Price reads what the auction of account's position asks now for token, in USD with 8 decimals per whole token;
// zero while none runs.
func (l *Liquidator) Price(opts *bind.CallOpts, account common.Address, position [32]byte, token common.Address) (*big.Int, error) {
	return read(l.c, opts, l.target, l.liquidator.UnpackPrice)(l.liquidator.TryPackPrice(account, position, token))
}

// HourlyAllowance reads how much of token a purchase from account's position may take now while the market is closed.
func (l *Liquidator) HourlyAllowance(opts *bind.CallOpts, account common.Address, position [32]byte, token common.Address) (*big.Int, error) {
	return read(l.c, opts, l.target, l.liquidator.UnpackHourlyAllowance)(
		l.liquidator.TryPackHourlyAllowance(account, position, token))
}

// HeldUntil reads until when, in seconds, the reopening auction holds account's position out of the Dutch auction.
func (l *Liquidator) HeldUntil(opts *bind.CallOpts, account common.Address, position [32]byte) (uint64, error) {
	return read(l.c, opts, l.target, l.liquidator.UnpackHeldUntil)(l.liquidator.TryPackHeldUntil(account, position))
}

// Start starts an auction of account's position once it falls short. Anyone may.
func (l *Liquidator) Start(account common.Address, position [32]byte) (Tx, error) {
	return l.tx(l.liquidator.TryPackStart(account, position))
}

// Stop stops the auction of account's position once it no longer falls short. Anyone may.
func (l *Liquidator) Stop(account common.Address, position [32]byte) (Tx, error) {
	return l.tx(l.liquidator.TryPackStop(account, position))
}

// Buy buys up to amount of token from account's position at the auction's price, for at most maxCost of the sender's
// USDG, to receiver.
func (l *Liquidator) Buy(account common.Address, position [32]byte, token common.Address, amount, maxCost *big.Int, receiver common.Address) (Tx, error) {
	if err := present(amount, maxCost); err != nil {
		return Tx{}, err
	}
	return l.tx(l.liquidator.TryPackBuy(account, position, token, amount, maxCost, receiver))
}

// SettleCash repays what account's position owes with the USDG it holds, once it falls short. Anyone may.
func (l *Liquidator) SettleCash(account common.Address, position [32]byte) (Tx, error) {
	return l.tx(l.liquidator.TryPackSettleCash(account, position))
}

// WriteOff writes off what account's emptied position owes. Anyone may while the accounts have no backstop.
func (l *Liquidator) WriteOff(account common.Address, position [32]byte) (Tx, error) {
	return l.tx(l.liquidator.TryPackWriteOff(account, position))
}

// Recall recalls all account's position lent of token and has not recalled, once it falls short. Anyone may.
func (l *Liquidator) Recall(account common.Address, position [32]byte, token common.Address) (Tx, error) {
	return l.tx(l.liquidator.TryPackRecall(account, position, token))
}

// Backstop reads the gap backstop and packs its transactions.
type Backstop struct {
	target
	backstop *gapbackstop.GapBackstop
}

// Backstop returns the gap backstop of the registry's .tapehouse.GapBackstop.
func (c *Client) Backstop() *Backstop {
	return &Backstop{
		target:   lookup(c.deployments.Tapehouse, "GapBackstop", ".tapehouse"),
		backstop: gapbackstop.NewGapBackstop(),
	}
}

// Claim claims the premium the accounts set aside for the backstop. Anyone may.
func (b *Backstop) Claim() (Tx, error) {
	return b.tx(b.backstop.TryPackClaim())
}

// Cover repays what account's emptied position owes within the backstop's limits and has the rest written off.
// Anyone may.
func (b *Backstop) Cover(account common.Address, position [32]byte) (Tx, error) {
	return b.tx(b.backstop.TryPackCover(account, position))
}
