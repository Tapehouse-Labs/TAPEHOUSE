// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"testing"
)

func committed(symbol string, openMs uint64) Committed {
	return Committed{Symbol: symbol, OpenMs: openMs, Quantity: big.NewInt(5e18), Price: big.NewInt(95e8),
		Deposit: big.NewInt(1000e6), Salt: [32]byte{1, 2, 3}, Commitment: [32]byte{9, 8, 7}}
}

func TestTheStateIsOnDiskBeforeAddReturns(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bidder.json")
	state, err := LoadState(path)
	if err != nil || len(state.Bids()) != 0 {
		t.Fatalf("a missing file is an empty state: %v", err)
	}
	want := committed("SPY", 1_790_000_000_000)
	if err := state.Add(want); err != nil {
		t.Fatal(err)
	}
	fresh, err := LoadState(path)
	if err != nil || len(fresh.Bids()) != 1 {
		t.Fatalf("%v, %v", fresh, err)
	}
	got := fresh.Bids()[0]
	if got.Symbol != "SPY" || got.OpenMs != want.OpenMs || got.Quantity.Cmp(want.Quantity) != 0 ||
		got.Price.Cmp(want.Price) != 0 || got.Deposit.Cmp(want.Deposit) != 0 || got.Salt != want.Salt || got.Commitment != want.Commitment {
		t.Fatalf("%+v, want %+v", got, want)
	}
	if err := state.Add(want); err == nil {
		t.Fatal("a second bid in one round was added")
	}
}

func TestAWriteNeverLeavesAHalfFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bidder.json")
	state, err := LoadState(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := state.Add(committed("SPY", 1)); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(path)
	state.afterTemporary = func() error { return errors.New("the process died") }
	if err := state.Add(committed("NVDA", 2)); err == nil {
		t.Fatal("an interrupted write succeeded")
	}
	after, _ := os.ReadFile(path)
	if string(before) != string(after) {
		t.Fatalf("the old file changed:\n%s\n%s", before, after)
	}
	if len(state.Bids()) != 1 {
		t.Fatalf("a bid that never reached the disk is held: %v", state.Bids())
	}
	state.afterTemporary = nil
	fresh, err := LoadState(path)
	if err != nil || len(fresh.Bids()) != 1 {
		t.Fatalf("%v, %v", fresh.Bids(), err)
	}
}

func TestARemovedRoundIsGone(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bidder.json")
	state, _ := LoadState(path)
	_ = state.Add(committed("SPY", 1))
	_ = state.Add(committed("NVDA", 1))
	if err := state.Remove("SPY", 1); err != nil {
		t.Fatal(err)
	}
	fresh, err := LoadState(path)
	if err != nil || len(state.Bids()) != 1 || len(fresh.Bids()) != 1 || fresh.Bids()[0].Symbol != "NVDA" {
		t.Fatalf("%v, %v, %v", state.Bids(), fresh.Bids(), err)
	}
}

func TestAHalfWrittenStateIsRefused(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bidder.json")
	if err := os.WriteFile(path, []byte(`{"bids":[`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadState(path); err == nil {
		t.Fatal("a corrupt state loaded")
	}
}
