// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"errors"
	"math/big"

	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/band"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/bandfeed"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/reopeningauction"
)

// ErrBelowFloor, ErrBelowMinBid and ErrNoSupply are the reasons a bid is not planned: the first two are two of the
// reveal's three refusals, and a bid refused at reveal forfeits its deposit.
var (
	ErrBelowFloor  = errors.New("the price is below the round's floor")
	ErrBelowMinBid = errors.New("the bid is worth less than the minimum bid at the round's floor")
	ErrNoSupply    = errors.New("the round has no supply")
	errEscrow      = errors.New("the bid's escrow is above its deposit")
)

// tokenToUSDG is 10^20: a quantity with 18 decimals times a price with 8 decimals, in USDG with 6.
var tokenToUSDG = new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil)

// Bid is a planned bid: its quantity with 18 decimals, its price in USD with 8 decimals per token and its deposit in
// USDG with 6 decimals.
type Bid struct {
	Quantity, Price, Deposit *big.Int
}

// OwnPrice is the bot's price for a Stock Token: the lower of the low edge of the band sealed before the reopen and
// of the band now, less discountBps. The seal counts only in state 1 or 3, the band now only in state 1, 2 or 3 and
// live; with neither there is no price.
func OwnPrice(sealed *bandfeed.SealsOutput, now band.QuoteOutput, discountBps uint64) (*big.Int, bool) {
	if discountBps >= 10_000 {
		return nil, false
	}
	var low uint64
	found := false
	if sealed != nil && (sealed.State == 1 || sealed.State == 3) {
		low, found = sealed.Low, true
	}
	if now.State >= 1 && now.State <= 3 && now.Live != 0 && (!found || now.Low < low) {
		low, found = now.Low, true
	}
	if !found || low == 0 {
		return nil, false
	}
	price := new(big.Int).SetUint64(low)
	price.Mul(price, new(big.Int).SetUint64(10_000-discountBps))
	return price.Div(price, big.NewInt(10_000)), true
}

// PlanBid sizes a bid at price by budget and the round's supply, with a deposit of budget, at least the bond, and
// refuses any bid the auction's reveal would refuse: a price below the floor, a bid worth less than the minimum at the
// floor, or an escrow above the deposit.
func PlanBid(price *big.Int, round reopeningauction.ReopeningAuctionRound, budget *big.Int, c sdk.AuctionConstants) (Bid, error) {
	switch {
	case round.Supply.Sign() == 0:
		return Bid{}, ErrNoSupply
	case price.Sign() <= 0 || price.Cmp(round.Floor) < 0:
		return Bid{}, ErrBelowFloor
	}
	quantity := new(big.Int).Mul(budget, tokenToUSDG)
	quantity.Div(quantity, price)
	if quantity.Cmp(round.Supply) > 0 {
		quantity.Set(round.Supply)
	}
	deposit := new(big.Int).Set(budget)
	if deposit.Cmp(c.Bond) < 0 {
		deposit.Set(c.Bond)
	}
	atFloor := new(big.Int).Mul(quantity, round.Floor)
	if atFloor.Div(atFloor, tokenToUSDG).Cmp(c.MinBid) < 0 {
		return Bid{}, ErrBelowMinBid
	}
	if Escrow(quantity, price).Cmp(deposit) > 0 {
		return Bid{}, errEscrow
	}
	return Bid{quantity, new(big.Int).Set(price), deposit}, nil
}

// Escrow is what a bid keeps of its deposit: its worth at its price, rounded up.
func Escrow(quantity, price *big.Int) *big.Int {
	worth := new(big.Int).Mul(quantity, price)
	worth.Add(worth, new(big.Int).Sub(tokenToUSDG, big.NewInt(1)))
	return worth.Div(worth, tokenToUSDG)
}
