// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/aggregator"
)

// feed returns a Chainlink feed double whose rounds started at started, by phase and aggregator round.
func (w *world) feed(name string, latest *big.Int, started map[uint64][]uint64) {
	w.read(name, "latestRoundData", func([]any) []any { return []any{latest, n(1), n(1), n(1), latest} })
	w.read(name, "getRoundData", func(args []any) []any {
		id := args[0].(*big.Int)
		phase, round := new(big.Int).Rsh(id, 64).Uint64(), new(big.Int).And(id, aggregatorMask).Uint64()
		rounds := started[phase]
		if round == 0 || round > uint64(len(rounds)) {
			return []any{id, n(0), n(0), n(0), id}
		}
		at := new(big.Int).SetUint64(rounds[round-1])
		return []any{id, n(100), at, at, id}
	})
}

func TestTheLastRoundBeforeATimeIsFoundAcrossAPhaseChange(t *testing.T) {
	w := newWorld(t)
	address := w.add("Feed:X", &aggregator.AggregatorMetaData)
	started := map[uint64][]uint64{1: {100, 200, 300, 400, 500, 600, 700}, 2: {1000, 1100, 1200}}
	w.feed("Feed:X", roundID(2, 3), started)
	k := w.keeper(Config{})
	for _, c := range []struct {
		atMs uint64
		want *big.Int
	}{
		{1_150_000, roundID(2, 2)}, {1_100_001, roundID(2, 2)}, {1_100_000, roundID(2, 1)}, {5_000_000, roundID(2, 3)},
		{999_000, roundID(1, 7)}, {350_000, roundID(1, 3)}, {100_001, roundID(1, 1)},
	} {
		got, err := k.lastRoundBefore(context.Background(), address, c.atMs)
		if err != nil || got.Cmp(c.want) != 0 {
			t.Errorf("before %d: %v, %v; want %v", c.atMs, got, err, c.want)
		}
	}
	if _, err := k.lastRoundBefore(context.Background(), address, 100_000); err == nil {
		t.Fatal("a round before the feed's first")
	}
}

func TestTheReopenIsRecordedAndEachSeriesObservedSettledReleasedOrVoided(t *testing.T) {
	w := pricedWorld(t)
	closesMs := (w.time - 2*86400) * 1000
	reopen := closesMs + 48*3600*1000
	s := &session{state: sessionClosed, nyse: 3, change: reopen, bound: reopen}
	w.session(s)
	recorded := uint64(0)
	w.read("GapCover", "lastCloseMs", func([]any) []any { return []any{closesMs} })
	w.read("GapCover", "reopenOf", func([]any) []any { return []any{new(big.Int).SetUint64(recorded)} })
	w.read("GapCover", "regularHours", func([]any) []any { return []any{false} })
	w.write("GapCover", "record", func(_ []any, mined bool) error {
		if mined {
			recorded = s.bound
		}
		return nil
	})
	type series struct {
		status uint8
		shift  uint16
	}
	all := map[string]*series{"NVDA": {}, "SPY": {}}
	released := map[int64]bool{}
	w.read("GapCover", "series", func(args []any) []any {
		x := all[text(args[0].([32]byte))]
		return []any{n(1000), uint64(0), uint64(0), x.shift, x.status, false}
	})
	w.read("GapCover", "covers", func(args []any) []any {
		holder := alice
		if released[args[0].(*big.Int).Int64()] {
			holder = common.Address{}
		}
		return []any{holder, closesMs, uint16(0), uint16(0), [32]byte{}, n(0), n(0)}
	})
	w.write("GapCover", "observe", func([]any, bool) error { return nil })
	var settled []any
	w.write("GapCover", "settle", func(args []any, mined bool) error {
		if mined {
			settled = args
			all[text(args[0].([32]byte))].status = 1
		}
		return nil
	})
	w.write("GapCover", "release", func(args []any, mined bool) error {
		if mined {
			released[args[0].(*big.Int).Int64()] = true
		}
		return nil
	})
	w.write("GapCover", "void", func(args []any, mined bool) error {
		if mined {
			all[text(args[0].([32]byte))].status = 2
		}
		return nil
	})
	w.read("GapCover", "feed", func([]any) []any { return []any{w.deployments.Chainlink["NVDA_USD"]} })
	w.feed("Feed:NVDA", roundID(1, 4), map[uint64][]uint64{1: {closesMs/1000 - 7200, closesMs/1000 - 60, reopen/1000 - 600, reopen/1000 + 120}})
	for id, asset := range map[int64]string{1: "NVDA", 2: "NVDA", 3: "SPY"} {
		w.event("tapehouse.GapCover", "Bought", map[string]any{"id": fmt.Sprint(id), "holder": alice.Hex(),
			"symbol": fmt.Sprintf("%#x", symbol(t, asset)), "closesMs": fmt.Sprint(closesMs)})
	}
	w.event("tapehouse.GapCover", "CloseRecorded", map[string]any{"closesMs": fmt.Sprint(closesMs), "atMs": fmt.Sprint(closesMs + 1000)})
	k := w.keeper(Config{})
	w.time = reopen/1000 - 3600
	if got := w.pass(t, k, "gapcover"); len(got) != 1 || got[0] != "GapCover.record[]" {
		t.Fatalf("the reopen: %v", got)
	}
	if got := w.pass(t, k, "gapcover"); len(got) != 0 {
		t.Fatalf("before the reopen: %v", got)
	}
	s.state, s.bound = sessionOpen, 0
	w.time = reopen/1000 + 30
	if got := w.pass(t, k, "gapcover"); len(got) != 2 || !strings.HasPrefix(got[0], "GapCover.observe") {
		t.Fatalf("the first slot: %v", got)
	}
	w.time += 20
	if got := w.pass(t, k, "gapcover"); len(got) != 0 {
		t.Fatalf("one slot observed twice: %v", got)
	}
	w.time += 20
	if got := w.pass(t, k, "gapcover"); len(got) != 2 {
		t.Fatalf("the second slot: %v", got)
	}
	all["SPY"].shift = 16
	w.time = reopen/1000 + 15*60
	got := w.pass(t, k, "gapcover")
	if len(got) != 2 || settled == nil || text(settled[0].([32]byte)) != "NVDA" || settled[2].(*big.Int).Cmp(roundID(1, 2)) != 0 || settled[3].(*big.Int).Cmp(roundID(1, 3)) != 0 {
		t.Fatalf("settled %v with %v", got, settled)
	}
	if got := w.pass(t, k, "gapcover"); len(got) != 2 || !released[1] || !released[2] {
		t.Fatalf("NVDA's covers released: %v", got)
	}
	w.time = closesMs/1000 + 7*86400
	if got := w.pass(t, k, "gapcover"); len(got) != 1 || !strings.HasPrefix(got[0], "GapCover.void") {
		t.Fatalf("SPY's series voided: %v", got)
	}
	if got := w.pass(t, k, "gapcover"); len(got) != 1 || !released[3] {
		t.Fatalf("SPY's cover released: %v", got)
	}
}
