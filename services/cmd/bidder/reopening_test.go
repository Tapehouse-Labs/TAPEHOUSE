// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"context"
	"math/big"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/reopeningauction"
)

func tick(t *testing.T, b *Bidder, m *market) {
	t.Helper()
	if err := b.Reopening(context.Background(), m.head()); err != nil {
		t.Fatal(err)
	}
}

func statePath(t *testing.T) string { return filepath.Join(t.TempDir(), "bidder.json") }

func TestNothingIsCommittedBeforeTheCommitLead(t *testing.T) {
	m := newMarket(t)
	b := m.bidder(t, statePath(t))
	m.at(openMs - 1_800_000 - 16*60_000)
	tick(t, b, m)
	if sent := m.mined(); len(sent) != 0 {
		t.Fatalf("%v", sent)
	}
	m.at(openMs - 1_800_000 - 15*60_000)
	tick(t, b, m)
	if sent := m.mined(); len(sent) != 2 {
		t.Fatalf("at the lead itself: %v", sent)
	}
}

func TestOneBidIsCommittedPerRoundWithinTheLead(t *testing.T) {
	m := newMarket(t)
	b := m.bidder(t, statePath(t))
	tick(t, b, m)
	calls := m.took()
	if len(calls) != 2 || calls[0].String() != "USDG.approve"+sprint(m.names["ReopeningAuction"], big.NewInt(1000e6)) {
		t.Fatalf("%v", calls)
	}
	held := b.State.Bids()
	if len(held) != 1 || held[0].Symbol != "SPY" || held[0].OpenMs != openMs {
		t.Fatalf("%+v", held)
	}
	bid := held[0]
	if bid.Quantity.Cmp(big.NewInt(5e18)) != 0 || bid.Price.Cmp(big.NewInt(9_506_000_000)) != 0 || bid.Deposit.Cmp(big.NewInt(1000e6)) != 0 {
		t.Fatalf("%+v", bid)
	}
	commitment, err := sdk.Commitment(me, "SPY", openMs, bid.Quantity, bid.Price, bid.Salt)
	if err != nil || bid.Commitment != commitment {
		t.Fatal(err)
	}
	if calls[1].String() != "ReopeningAuction.commit"+sprint(symbolOf(t, "SPY"), commitment, big.NewInt(1000e6)) {
		t.Fatalf("%v", calls[1])
	}
	tick(t, b, m)
	if sent := m.mined(); len(sent) != 0 {
		t.Fatalf("a second tick: %v", sent)
	}
}

func TestTheStateIsWrittenBeforeTheCommitIsSent(t *testing.T) {
	m := newMarket(t)
	path := statePath(t)
	b := m.bidder(t, path)
	seen := false
	m.onSend = func(c call) {
		if c.method != "commit" {
			return
		}
		fresh, err := LoadState(path)
		if err != nil || len(fresh.Bids()) != 1 || fresh.Bids()[0].Commitment != c.args[1].([32]byte) {
			t.Errorf("at the commit: %v, %v", fresh.Bids(), err)
		}
		seen = true
	}
	tick(t, b, m)
	if !seen {
		t.Fatal("no commit was sent")
	}
}

func TestABidRevealWouldRefuseIsNotCommitted(t *testing.T) {
	m := newMarket(t)
	m.rounds["SPY"].Floor = big.NewInt(96e8)
	b := m.bidder(t, statePath(t))
	tick(t, b, m)
	if sent := m.mined(); len(sent) != 0 || len(b.State.Bids()) != 0 {
		t.Fatalf("%v, %v", sent, b.State.Bids())
	}
	if !strings.Contains(m.logs.String(), ErrBelowFloor.Error()) {
		t.Fatalf("logs: %s", m.logs.String())
	}
}

func TestEveryCommittedBidIsRevealedInItsWindow(t *testing.T) {
	m := newMarket(t)
	b := m.bidder(t, statePath(t))
	tick(t, b, m)
	m.took()
	bid := b.State.Bids()[0]
	m.at(openMs - 1_800_000 - 1000)
	tick(t, b, m)
	if sent := m.mined(); len(sent) != 0 {
		t.Fatalf("before the window: %v", sent)
	}
	m.at(openMs - 1_800_000 + 1000)
	tick(t, b, m)
	want := "ReopeningAuction.reveal" + sprint(symbolOf(t, "SPY"), openMs, bid.Quantity, bid.Price, bid.Salt)
	if sent := m.mined(); len(sent) != 1 || sent[0] != want {
		t.Fatalf("%v, want %s", sent, want)
	}
	tick(t, b, m)
	if sent := m.mined(); len(sent) != 0 {
		t.Fatalf("a revealed bid again: %v", sent)
	}
}

func TestARestartRevealsFromTheStateFile(t *testing.T) {
	m := newMarket(t)
	path := statePath(t)
	tick(t, m.bidder(t, path), m)
	m.took()
	m.at(openMs - 1_800_000 + 1000)
	restarted := m.bidder(t, path)
	tick(t, restarted, m)
	if sent := m.mined(); len(sent) != 1 || !strings.HasPrefix(sent[0], "ReopeningAuction.reveal") {
		t.Fatalf("%v", sent)
	}
}

func TestACommitThatNeverLandedIsNotRevealed(t *testing.T) {
	m := newMarket(t)
	path := statePath(t)
	state, _ := LoadState(path)
	if err := state.Add(committed("SPY", openMs)); err != nil {
		t.Fatal(err)
	}
	m.at(openMs - 1_800_000 + 1000)
	tick(t, m.bidder(t, path), m)
	if sent := m.mined(); len(sent) != 0 {
		t.Fatalf("%v", sent)
	}
}

func TestACommitThatFailedIsSentAgain(t *testing.T) {
	m := newMarket(t)
	path := statePath(t)
	state, _ := LoadState(path)
	bid := Committed{Symbol: "SPY", OpenMs: openMs, Quantity: big.NewInt(5e18), Price: big.NewInt(9_506_000_000),
		Deposit: big.NewInt(1000e6), Salt: [32]byte{4}}
	bid.Commitment, _ = sdk.Commitment(me, "SPY", openMs, bid.Quantity, bid.Price, bid.Salt)
	if err := state.Add(bid); err != nil {
		t.Fatal(err)
	}
	tick(t, m.bidder(t, path), m)
	if sent := m.mined(); len(sent) != 2 || !strings.HasPrefix(sent[1], "ReopeningAuction.commit") {
		t.Fatalf("%v", sent)
	}
}

func TestAnOutbidRevealIsNotAnError(t *testing.T) {
	m := newMarket(t)
	b := m.bidder(t, statePath(t))
	tick(t, b, m)
	m.took()
	m.outbid = true
	m.at(openMs - 1_800_000 + 1000)
	tick(t, b, m)
	if sent := m.mined(); len(sent) != 1 {
		t.Fatalf("%v", sent)
	}
	if !strings.Contains(m.logs.String(), "outbid") {
		t.Fatalf("logs: %s", m.logs.String())
	}
	m.phaseOpenMs = 0
	m.rounds["SPY"].Cleared = true
	m.at(openMs + 1000)
	tick(t, b, m)
	if sent := m.mined(); len(sent) != 0 || len(b.State.Bids()) != 0 {
		t.Fatalf("nothing to claim: %v, %v", sent, b.State.Bids())
	}
}

func bid(bidder common.Address, claimed bool) reopeningauction.ReopeningAuctionBid {
	return reopeningauction.ReopeningAuctionBid{Bidder: bidder, Quantity: big.NewInt(1), Price: big.NewInt(1),
		Escrow: big.NewInt(1), Claimed: claimed}
}

func staged(t *testing.T, m *market) *Bidder {
	t.Helper()
	path := statePath(t)
	state, _ := LoadState(path)
	if err := state.Add(committed("SPY", openMs)); err != nil {
		t.Fatal(err)
	}
	m.phaseOpenMs = 0
	return m.bidder(t, path)
}

func TestEachOfItsBidsIsClaimedAfterTheClearing(t *testing.T) {
	m := newMarket(t)
	b := staged(t, m)
	m.rounds["SPY"].Cleared = true
	m.bids["SPY"] = []reopeningauction.ReopeningAuctionBid{bid(other, false), bid(me, false), bid(me, true)}
	m.at(openMs + 1000)
	tick(t, b, m)
	want := "ReopeningAuction.claim" + sprint(symbolOf(t, "SPY"), openMs, big.NewInt(1))
	if sent := m.mined(); len(sent) != 1 || sent[0] != want || len(b.State.Bids()) != 0 {
		t.Fatalf("%v, %v", sent, b.State.Bids())
	}
}

func TestALapsedRoundIsClaimed(t *testing.T) {
	m := newMarket(t)
	b := staged(t, m)
	m.bids["SPY"] = []reopeningauction.ReopeningAuctionBid{bid(me, false), bid(other, false), bid(me, false)}
	m.at(openMs + 3_600_000)
	tick(t, b, m)
	spy := symbolOf(t, "SPY")
	sent := m.mined()
	if len(sent) != 2 || sent[0] != "ReopeningAuction.claim"+sprint(spy, openMs, big.NewInt(0)) ||
		sent[1] != "ReopeningAuction.claim"+sprint(spy, openMs, big.NewInt(2)) || len(b.State.Bids()) != 0 {
		t.Fatalf("%v, %v", sent, b.State.Bids())
	}
}

func TestNothingIsClaimedBeforeTheClearing(t *testing.T) {
	m := newMarket(t)
	b := staged(t, m)
	m.bids["SPY"] = []reopeningauction.ReopeningAuctionBid{bid(me, false)}
	m.at(openMs + 3_600_000 - 1000)
	tick(t, b, m)
	if sent := m.mined(); len(sent) != 0 || len(b.State.Bids()) != 1 {
		t.Fatalf("%v, %v", sent, b.State.Bids())
	}
}

func TestADryRunSendsNothing(t *testing.T) {
	m := newMarket(t)
	path := statePath(t)
	b := m.bidder(t, path)
	b.DryRun = true
	tick(t, b, m)
	if sent := m.mined(); len(sent) != 0 || len(b.State.Bids()) != 0 {
		t.Fatalf("%v, %v", sent, b.State.Bids())
	}
	if !strings.Contains(m.logs.String(), "dry run") {
		t.Fatalf("logs: %s", m.logs.String())
	}
}

func TestACommitRefusedAsNotNowIsLoggedAndRetried(t *testing.T) {
	m := newMarket(t)
	b := m.bidder(t, statePath(t))
	m.write("ReopeningAuction", "commit", func([]any, bool) error { return m.fail("ReopeningAuction", "WrongPhase") })
	tick(t, b, m)
	if sent := m.mined(); len(sent) != 1 || !strings.HasPrefix(sent[0], "USDG.approve") || len(b.State.Bids()) != 1 {
		t.Fatalf("%v, %v", sent, b.State.Bids())
	}
	m.write("ReopeningAuction", "commit", func([]any, bool) error { return m.fail("ReopeningAuction", "InvalidBid") })
	if err := b.Reopening(context.Background(), m.head()); err == nil || !strings.Contains(err.Error(), "InvalidBid") {
		t.Fatalf("an unexpected revert is an error: %v", err)
	}
}
