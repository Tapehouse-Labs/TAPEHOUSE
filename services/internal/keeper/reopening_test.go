// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/reopeningauction"
)

func TestPositionsAreEnrolledAndEachRoundClearedForfeitedAndClaimed(t *testing.T) {
	w, positions := margin(t)
	nvda, spy := w.deployments.Tokens["NVDA"], w.deployments.Tokens["SPY"]
	openMs := (w.time + 3*3600) * 1000
	positions[pos{alice, sdk.Cross}] = &account{held: map[common.Address]*big.Int{nvda: n(5)}, inBaskets: n(0)}
	positions[pos{bob, sdk.Cross}] = &account{held: map[common.Address]*big.Int{nvda: n(5), spy: n(2)}}
	positions[pos{bob, symbol(t, "SPY")}] = &account{held: map[common.Address]*big.Int{}, inBaskets: nil}
	w.debt(alice, sdk.Cross, 80, 100)
	w.debt(bob, sdk.Cross, 20, 100)
	w.debt(bob, symbol(t, "SPY"), 150, 100)
	phase := struct {
		open      uint64
		revealing bool
	}{openMs, false}
	w.read("ReopeningAuction", "phase", func([]any) []any { return []any{phase.open, phase.revealing} })
	w.read("ReopeningAuction", "lastOpenMs", func([]any) []any { return []any{openMs} })
	enrolled := map[string]bool{}
	w.read("ReopeningAuction", "enrolled", func(args []any) []any {
		return []any{enrolled[fmt.Sprint(args[1:])]}
	})
	w.write("ReopeningAuction", "enroll", func(args []any, mined bool) error {
		if mined {
			enrolled[fmt.Sprint(args)] = true
		}
		return nil
	})
	k := w.keeper(Config{})
	got := w.pass(t, k, "reopening")
	enroll := func(account common.Address, position [32]byte, asset string) string {
		return fmt.Sprintf("ReopeningAuction.enroll[%s %v %v]", account.Hex(), position, symbol(t, asset))
	}
	want := []string{enroll(bob, sdk.Cross, "NVDA"), enroll(bob, sdk.Cross, "SPY"), enroll(alice, sdk.Cross, "NVDA")}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("mined %v, want %v", got, want)
	}
	if got := w.pass(t, k, "reopening"); len(got) != 0 {
		t.Fatalf("enrolled twice: %v", got)
	}

	phase.open = 0
	w.time = openMs/1000 + 60
	round := reopeningauction.ReopeningAuctionRound{Floor: n(20_000_000_000), Supply: n(10), Deposits: n(0), Pool: n(0),
		Price: n(0), Above: n(0), AtPrice: n(0), Sold: n(0), Paid: n(0), Taken: n(0)}
	bid := func(quantity, price int64) reopeningauction.ReopeningAuctionBid {
		return reopeningauction.ReopeningAuctionBid{Bidder: alice, Quantity: n(quantity), Price: n(price), Escrow: n(1)}
	}
	bids := []reopeningauction.ReopeningAuctionBid{bid(6, 23_000_000_000), bid(6, 22_500_000_000), bid(6, 22_000_000_000)}
	w.read("ReopeningAuction", "round", func(args []any) []any {
		if text(args[0].([32]byte)) != "NVDA" {
			return []any{reopeningauction.ReopeningAuctionRound{Floor: n(0), Supply: n(0), Deposits: n(0), Pool: n(0),
				Price: n(0), Above: n(0), AtPrice: n(0), Sold: n(0), Paid: n(0), Taken: n(0)}}
		}
		return []any{round}
	})
	w.read("ReopeningAuction", "bids", func([]any) []any { return []any{bids} })
	var cleared *big.Int
	w.write("ReopeningAuction", "clear", func(args []any, mined bool) error {
		if mined {
			cleared, round.Cleared = args[2].(*big.Int), true
		}
		return nil
	})
	pending := map[[32]byte]common.Address{{1}: bob, {2}: alice}
	for commitment := range pending {
		w.event("tapehouse.ReopeningAuction", "Committed", map[string]any{"symbol": fmt.Sprintf("%#x", symbol(t, "NVDA")),
			"openMs": fmt.Sprint(openMs), "bidder": alice.Hex(), "commitment": fmt.Sprintf("%#x", commitment), "deposit": "100000000"})
	}
	delete(pending, [32]byte{2})
	w.read("ReopeningAuction", "commitments", func(args []any) []any { return []any{pending[args[2].([32]byte)], n(1)} })
	w.write("ReopeningAuction", "forfeit", func(args []any, mined bool) error {
		if mined {
			delete(pending, args[2].([32]byte))
		}
		return nil
	})
	got = w.pass(t, k, "reopening")
	if len(got) != 2 || cleared.Int64() != 22_500_000_000 || got[1] != fmt.Sprintf("ReopeningAuction.forfeit[%v %d %v]", symbol(t, "NVDA"), openMs, [32]byte{1}) {
		t.Fatalf("mined %v, cleared at %v", got, cleared)
	}
	if got := w.pass(t, k, "reopening"); len(got) != 0 {
		t.Fatalf("a cleared round worked again: %v", got)
	}
	w.time = openMs/1000 + 3600
	claimed := map[int64]bool{}
	w.write("ReopeningAuction", "claim", func(args []any, mined bool) error {
		if mined {
			claimed[args[2].(*big.Int).Int64()] = true
			bids[args[2].(*big.Int).Int64()].Claimed = true
		}
		return nil
	})
	bids[1].Claimed = true
	if got := w.pass(t, k, "reopening"); len(got) != 2 || !claimed[0] || !claimed[2] {
		t.Fatalf("claims once the window closed: %v", got)
	}
}
