// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"errors"
	"math/big"
	"math/rand"
	"testing"

	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/band"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/bandfeed"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/reopeningauction"
)

var constants = sdk.AuctionConstants{Bond: big.NewInt(100e6), MinBid: big.NewInt(100e6), RevealMs: 1_800_000, ClearMs: 3_600_000}

func seal(state uint8, low uint64) *bandfeed.SealsOutput {
	return &bandfeed.SealsOutput{State: state, Live: 1, Low: low, SealedAt: 1}
}

func quote(state, live uint8, low uint64) band.QuoteOutput {
	return band.QuoteOutput{State: state, Live: live, Low: low}
}

func round(floor int64, supply *big.Int) reopeningauction.ReopeningAuctionRound {
	return reopeningauction.ReopeningAuctionRound{Floor: big.NewInt(floor), Supply: supply}
}

func TestOwnPriceIsTheLowerLowLessTheDiscount(t *testing.T) {
	got, ok := OwnPrice(seal(3, 100e8), quote(3, 1, 98e8), 300)
	if !ok || got.Cmp(big.NewInt(95_06e6)) != 0 {
		t.Fatalf("%v, %v", got, ok)
	}
	got, ok = OwnPrice(seal(1, 97e8), quote(3, 1, 98e8), 300)
	if !ok || got.Cmp(big.NewInt(94_09e6)) != 0 {
		t.Fatalf("%v, %v", got, ok)
	}
}

func TestASealOrABandThatCannotVouchIsIgnored(t *testing.T) {
	price := func(sealed *bandfeed.SealsOutput, now band.QuoteOutput) (int64, bool) {
		got, ok := OwnPrice(sealed, now, 0)
		if !ok {
			return 0, false
		}
		return got.Int64(), true
	}
	for _, c := range []struct {
		name   string
		sealed *bandfeed.SealsOutput
		now    band.QuoteOutput
		want   int64
		ok     bool
	}{
		{"a closed seal", seal(2, 90e8), quote(3, 1, 98e8), 98e8, true},
		{"no seal", seal(0, 0), quote(3, 1, 98e8), 98e8, true},
		{"a nil seal", nil, quote(1, 1, 98e8), 98e8, true},
		{"a band that is not live", seal(3, 97e8), quote(3, 0, 90e8), 97e8, true},
		{"a band in state 0", seal(1, 97e8), quote(0, 1, 90e8), 97e8, true},
		{"neither", seal(2, 90e8), quote(0, 0, 90e8), 0, false},
		{"a zero low", seal(3, 0), quote(3, 0, 0), 0, false},
	} {
		if got, ok := price(c.sealed, c.now); got != c.want || ok != c.ok {
			t.Errorf("%s: %d, %v", c.name, got, ok)
		}
	}
	if _, ok := OwnPrice(seal(3, 100e8), quote(3, 1, 98e8), 10_000); ok {
		t.Error("a discount of the whole price priced")
	}
}

func TestABidIsSizedByTheBudgetAndTheSupply(t *testing.T) {
	price := big.NewInt(95e8)
	bid, err := PlanBid(price, round(85e8, big.NewInt(5e18)), big.NewInt(1000e6), constants)
	want := new(big.Int).Div(new(big.Int).Mul(big.NewInt(1e9), tokenToUSDG), price)
	if err != nil || want.Cmp(big.NewInt(5e18)) < 0 {
		t.Fatalf("%v, %v", err, want)
	}
	if bid.Quantity.Cmp(big.NewInt(5e18)) != 0 {
		t.Fatalf("the supply caps the bid: %v", bid.Quantity)
	}
	bid, err = PlanBid(price, round(85e8, big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(20))), big.NewInt(1000e6), constants)
	if err != nil || bid.Quantity.Cmp(want) != 0 || bid.Price.Cmp(price) != 0 {
		t.Fatalf("the budget sizes the bid: %+v, %v, want %v", bid, err, want)
	}
}

func TestTheDepositIsTheBudgetAtLeastTheBond(t *testing.T) {
	loose := sdk.AuctionConstants{Bond: big.NewInt(100e6), MinBid: big.NewInt(1), RevealMs: 1, ClearMs: 1}
	for budget, deposit := range map[int64]int64{1000e6: 1000e6, 50e6: 100e6, 100e6: 100e6} {
		bid, err := PlanBid(big.NewInt(95e8), round(85e8, big.NewInt(5e18)), big.NewInt(budget), loose)
		if err != nil || bid.Deposit.Int64() != deposit {
			t.Errorf("budget %d: %+v, %v", budget, bid, err)
		}
	}
}

func TestABidBelowTheFloorIsNotPlanned(t *testing.T) {
	if _, err := PlanBid(big.NewInt(85e8-1), round(85e8, big.NewInt(5e18)), big.NewInt(1000e6), constants); !errors.Is(err, ErrBelowFloor) {
		t.Fatalf("floor - 1: %v", err)
	}
	if _, err := PlanBid(big.NewInt(85e8), round(85e8, big.NewInt(5e18)), big.NewInt(1000e6), constants); err != nil {
		t.Fatalf("the floor: %v", err)
	}
}

func TestABidWorthLessThanMinBidAtTheFloorIsNotPlanned(t *testing.T) {
	floor := int64(100e8)
	// floor 100e8: quantity q is worth q * 100e8 / 1e20 = q / 1e10 at the floor.
	for supply, want := range map[int64]error{99_999_999e10 + 9_999_999_999: ErrBelowMinBid, 100e6 * 1e10: nil} {
		_, err := PlanBid(big.NewInt(100e8), round(floor, big.NewInt(supply)), big.NewInt(1000e6), constants)
		if !errors.Is(err, want) && (err != nil || want != nil) {
			t.Errorf("supply %d: %v, want %v", supply, err, want)
		}
	}
}

func TestEveryPlannedBidIsOneRevealAccepts(t *testing.T) {
	rng := rand.New(rand.NewSource(8))
	planned := 0
	for range 10_000 {
		floor := big.NewInt(1e8 + rng.Int63n(1000e8))
		price := new(big.Int).Add(floor, big.NewInt(rng.Int63n(200e8)-20e8))
		if price.Sign() <= 0 {
			continue
		}
		supply := new(big.Int).Mul(big.NewInt(1+rng.Int63n(1e6)), big.NewInt(1e12+rng.Int63n(1e12)))
		budget := big.NewInt(1 + rng.Int63n(5000e6))
		r := reopeningauction.ReopeningAuctionRound{Floor: floor, Supply: supply}
		bid, err := PlanBid(price, r, budget, constants)
		if err != nil {
			continue
		}
		planned++
		escrow := new(big.Int).Div(new(big.Int).Add(new(big.Int).Mul(bid.Quantity, bid.Price), new(big.Int).Sub(tokenToUSDG, big.NewInt(1))), tokenToUSDG)
		atFloor := new(big.Int).Div(new(big.Int).Mul(bid.Quantity, floor), tokenToUSDG)
		if bid.Price.Cmp(floor) < 0 || atFloor.Cmp(constants.MinBid) < 0 || escrow.Cmp(bid.Deposit) > 0 ||
			bid.Deposit.Cmp(constants.Bond) < 0 || bid.Quantity.Cmp(supply) > 0 || escrow.Cmp(Escrow(bid.Quantity, bid.Price)) != 0 {
			t.Fatalf("a planned bid reveal refuses: %+v for floor %v, supply %v, budget %v", bid, floor, supply, budget)
		}
	}
	if planned < 1000 {
		t.Fatalf("only %d of 10000 rounds planned a bid", planned)
	}
}

func TestARoundWithNoSupplyGetsNoBid(t *testing.T) {
	if _, err := PlanBid(big.NewInt(95e8), round(85e8, big.NewInt(0)), big.NewInt(1000e6), constants); !errors.Is(err, ErrNoSupply) {
		t.Fatal(err)
	}
}

func TestEscrowRoundsUp(t *testing.T) {
	if got := Escrow(big.NewInt(3), big.NewInt(1e8+1)); got.Int64() != 1 {
		t.Fatalf("%v", got)
	}
	if got := Escrow(big.NewInt(1e12), big.NewInt(1e8)); got.Int64() != 1 {
		t.Fatalf("%v", got)
	}
	if got := Escrow(big.NewInt(0), big.NewInt(1e8)); got.Sign() != 0 {
		t.Fatalf("%v", got)
	}
}
