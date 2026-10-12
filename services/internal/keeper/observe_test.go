// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// observedWorld is a world whose margin accounts take NVDA and SPY, with each asset's band state, multiplier change
// status and last observation by name.
func observedWorld(t *testing.T) (*world, map[string]uint8, map[string]uint8, map[string]uint64) {
	t.Helper()
	w := pricedWorld(t)
	state := map[string]uint8{"NVDA": 3, "SPY": 3}
	status := map[string]uint8{}
	unusable := map[string]uint64{}
	nvda, spy := w.deployments.Tokens["NVDA"], w.deployments.Tokens["SPY"]
	w.read("MarginAccounts", "stocks", func([]any) []any {
		return []any{[][32]byte{symbol(t, "NVDA"), symbol(t, "SPY")}, []common.Address{nvda, spy}}
	})
	w.read("Band", "quote", func(args []any) []any {
		return []any{state[text(args[0].([32]byte))], uint8(2), uint64(1), uint64(55), uint64(1), n(1)}
	})
	w.read("Band", "corporateAction", func(args []any) []any {
		return []any{status[text(args[0].([32]byte))], uint64(0), n(0), n(0)}
	})
	w.read("MarginAccounts", "unusableAt", func(args []any) []any { return []any{unusable[text(args[0].([32]byte))]} })
	w.read("MarginAccounts", "RECOVERY", func([]any) []any { return []any{n(1800)} })
	w.write("MarginAccounts", "observe", func(args []any, mined bool) error {
		if mined {
			unusable[text(args[0].([32]byte))] = w.time
		}
		return nil
	})
	return w, state, status, unusable
}

func TestTheObserveKeeperSendsNothingWhileEveryAssetIsHealthy(t *testing.T) {
	w, _, _, _ := observedWorld(t)
	k := w.keeper(Config{})
	if got := w.pass(t, k, "observe"); len(got) != 0 {
		t.Fatalf("healthy assets: %v", got)
	}
}

func TestTheObserveKeeperObservesEachUnusableAssetOnceAMinute(t *testing.T) {
	w, state, status, _ := observedWorld(t)
	state["NVDA"] = 0
	status["SPY"] = 2
	k := w.keeper(Config{})
	want := []string{
		fmt.Sprintf("MarginAccounts.observe[%v]", symbol(t, "NVDA")),
		fmt.Sprintf("MarginAccounts.observe[%v]", symbol(t, "SPY")),
	}
	if got := w.pass(t, k, "observe"); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("a halt and a pending step: %v, want %v", got, want)
	}
	w.time += 59
	if got := w.pass(t, k, "observe"); len(got) != 0 {
		t.Fatalf("a fresh observation: %v", got)
	}
	w.time++
	if got := w.pass(t, k, "observe"); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("an observation a minute old: %v, want %v", got, want)
	}
	state["NVDA"] = 3
	w.time += 60
	if got := w.pass(t, k, "observe"); fmt.Sprint(got) != fmt.Sprint(want[1:]) {
		t.Fatalf("a recovered asset: %v, want %v", got, want[1:])
	}
}

func TestTheObserveKeeperLogsAnAssetItCannotObserveAndGoesOn(t *testing.T) {
	w, state, _, _ := observedWorld(t)
	state["NVDA"], state["SPY"] = 0, 0
	w.write("MarginAccounts", "observe", func(args []any, _ bool) error {
		if text(args[0].([32]byte)) == "NVDA" {
			return w.fail("MarginAccounts", "UnknownAsset", args[0])
		}
		return nil
	})
	k := w.keeper(Config{})
	want := []string{fmt.Sprintf("MarginAccounts.observe[%v]", symbol(t, "SPY"))}
	if got := w.pass(t, k, "observe"); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("a refused observation: %v, want %v", got, want)
	}
	if !strings.Contains(w.logs.String(), "observe failed") {
		t.Fatalf("the failure was not logged: %s", w.logs.String())
	}
}

func TestTheObserveKeeperIsListed(t *testing.T) {
	w, _, _, _ := observedWorld(t)
	if !slices.Contains(w.keeper(Config{}).Keepers(), "observe") {
		t.Fatalf("the keepers are %v", w.keeper(Config{}).Keepers())
	}
	if !slices.Contains(Names(), "observe") {
		t.Fatal("observe is not a keeper")
	}
}
