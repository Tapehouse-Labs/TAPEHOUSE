// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"fmt"
	"testing"
)

type session struct {
	state, nyse   uint8
	change, bound uint64
}

func (w *world) session(s *session) {
	w.read("Band", "session", func([]any) []any {
		return []any{s.state, s.nyse, uint8(0), s.change, s.bound}
	})
}

func (w *world) pass(t *testing.T, k *Keeper, name string) []string {
	t.Helper()
	if err := k.Pass(context.Background(), name); err != nil {
		t.Fatal(err)
	}
	return w.mined()
}

func TestTheSealGoesInTheWindowOnceEarlyAndAgainAsLateAsItCan(t *testing.T) {
	w := pricedWorld(t)
	reopen := (w.time + 1200) * 1000
	s := &session{state: sessionClosed, nyse: 3, change: reopen, bound: reopen}
	w.session(s)
	sealedAt := map[string]uint64{}
	for _, name := range []string{"NVDA", "TSLA", "SPY"} {
		w.read("BandFeed:"+name, "seals", func(args []any) []any {
			if args[0].(uint64) != reopen {
				t.Fatalf("seals of %d", args[0])
			}
			return []any{uint8(3), uint8(2), uint64(1), uint64(1), uint64(1), n(1), sealedAt[name]}
		})
		w.write("BandFeed:"+name, "seal", func(_ []any, mined bool) error {
			if mined {
				sealedAt[name] = w.time
			}
			return nil
		})
	}
	k := w.keeper(Config{})
	for _, step := range []struct {
		at    uint64
		seals int
	}{{600, 0}, {601, 0}, {899, 0}, {900, 3}, {901, 0}, {1109, 0}, {1110, 3}, {1111, 0}, {1200, 0}} {
		w.time = reopen/1000 - 1200 + step.at
		if got := w.pass(t, k, "seal"); len(got) != step.seals {
			t.Fatalf("%d s before the reopen: %v", 1200-step.at, got)
		}
	}
	s.state = sessionOpen
	w.time = reopen/1000 - 30
	if got := w.pass(t, k, "seal"); len(got) != 0 {
		t.Fatalf("sealed while open: %v", got)
	}
}

func TestThePremiumRecordsTheCloseAheadTheReopeningAndTheSealWindow(t *testing.T) {
	w := newWorld(t)
	closesMs := (w.time + 3600) * 1000
	s := &session{state: sessionOpen, nyse: 1, change: closesMs - 14_400_000, bound: closesMs}
	w.session(s)
	var closure struct{ closes, reopens, accrued uint64 }
	w.read("MarginAccounts", "closure", func([]any) []any { return []any{closure.closes, closure.reopens, closure.accrued} })
	w.write("MarginAccounts", "accruePremium", func(_ []any, mined bool) error {
		if mined {
			if s.state == sessionOpen {
				closure.closes, closure.reopens = s.bound, 0
			} else {
				closure.reopens = s.bound
			}
			closure.accrued = w.time * 1000
		}
		return nil
	})
	k := w.keeper(Config{})
	if got := w.pass(t, k, "premium"); len(got) != 1 || closure.closes != closesMs {
		t.Fatalf("the close ahead: %v", got)
	}
	if got := w.pass(t, k, "premium"); len(got) != 0 {
		t.Fatalf("the close recorded twice: %v", got)
	}
	reopensMs := closesMs + 48*3600*1000
	s.state, s.nyse, s.bound = sessionClosed, 3, reopensMs
	w.time = closesMs/1000 + 60
	if got := w.pass(t, k, "premium"); len(got) != 1 || closure.reopens != reopensMs {
		t.Fatalf("the reopening: %v", got)
	}
	w.time = reopensMs/1000 - 601
	if got := w.pass(t, k, "premium"); len(got) != 0 {
		t.Fatalf("accrued before the seal window: %v", got)
	}
	w.time++
	if got := w.pass(t, k, "premium"); len(got) != 1 {
		t.Fatalf("no accrual in the seal window: %v", got)
	}
	if got := w.pass(t, k, "premium"); len(got) != 0 {
		t.Fatalf("accrued twice in the seal window: %v", got)
	}
	s.bound += 3600 * 1000
	if got := w.pass(t, k, "premium"); len(got) != 1 {
		t.Fatalf("a changed reopening not caught: %v", got)
	}
}

func TestEachShortableAssetIsMarkedOnceAfterTheRegularClose(t *testing.T) {
	w := pricedWorld(t)
	closeMs := (w.time + 60) * 1000
	s := &session{state: sessionOpen, nyse: 1, change: closeMs, bound: closeMs + 14_400_000}
	w.session(s)
	w.read("ShortPositions", "fee", func(args []any) []any {
		if text(args[0].([32]byte)) == "TSLA" {
			return []any{n(0)}
		}
		return []any{n(500)}
	})
	w.write("ShortPositions", "mark", func([]any, bool) error { return nil })
	k := w.keeper(Config{})
	if got := w.pass(t, k, "mark"); len(got) != 0 {
		t.Fatalf("marked in regular hours: %v", got)
	}
	s.nyse, s.change = 3, closeMs+60*3600*1000
	w.time += 61
	got := w.pass(t, k, "mark")
	want := []string{fmt.Sprintf("ShortPositions.mark[%v]", symbol(t, "NVDA")), fmt.Sprintf("ShortPositions.mark[%v]", symbol(t, "SPY"))}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("marked %v, want %v", got, want)
	}
	if got := w.pass(t, k, "mark"); len(got) != 0 {
		t.Fatalf("marked twice: %v", got)
	}
	restarted := w.keeper(Config{})
	if got := w.pass(t, restarted, "mark"); len(got) != 0 {
		t.Fatalf("a keeper that never saw the close marked: %v", got)
	}
}
