// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/tapehouse/tapehouse/services/sdk"
)

func TestTheHeadIsSettledForItsRecallerAndTheLatestDueTicketBoughtIn(t *testing.T) {
	w, _ := margin(t)
	w.vault("SPY")
	spy := w.deployments.Tokens["SPY"]
	type ticket struct {
		end   int64
		dueAt uint64
	}
	tickets := []ticket{{10, w.time - 7200}, {25, w.time - 60}, {40, w.time + 3600}}
	head, assigned, claimable := uint64(0), int64(10), int64(4)
	w.read("Vault:SPY", "head", func([]any) []any { return []any{head} })
	w.read("Vault:SPY", "tickets", func([]any) []any { return []any{uint64(len(tickets))} })
	w.read("Vault:SPY", "assigned", func([]any) []any { return []any{n(assigned)} })
	w.read("Vault:SPY", "claimable", func(args []any) []any {
		if args[0].(*big.Int).Uint64() != head {
			return []any{n(0)}
		}
		return []any{n(claimable)}
	})
	w.read("Vault:SPY", "ticket", func(args []any) []any {
		tk := tickets[args[0].(*big.Int).Int64()]
		return []any{n(0), n(tk.end), n(0), tk.dueAt}
	})
	w.event("tapehouse.MarginAccounts", "Recall", map[string]any{"account": bob.Hex(), "position": fmt.Sprintf("%#x", sdk.Cross),
		"symbol": fmt.Sprintf("%#x", symbol(t, "SPY")), "ticket": "0", "amount": "10"})
	w.event("tapehouse.MarginAccounts", "Recall", map[string]any{"account": alice.Hex(), "position": fmt.Sprintf("%#x", symbol(t, "SPY")),
		"symbol": fmt.Sprintf("%#x", symbol(t, "SPY")), "ticket": "1", "amount": "15"})
	w.write("MarginAccounts", "settle", func(_ []any, mined bool) error {
		if mined {
			head++
		}
		return nil
	})
	failed := false
	w.write("Vault:SPY", "buyIn", func(args []any, mined bool) error {
		if failed {
			return w.fail("Vault:SPY", "BuyInFailed", n(15))
		}
		if mined {
			assigned = tickets[args[0].(*big.Int).Int64()].end
		}
		return nil
	})
	k := w.keeper(Config{})
	got := w.pass(t, k, "recall")
	want := []string{
		fmt.Sprintf("MarginAccounts.settle[%s %v %s]", bob.Hex(), sdk.Cross, spy.Hex()),
		"Vault.buyIn[1]",
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("mined %v, want %v", got, want)
	}
	claimable = 0
	if got := w.pass(t, k, "recall"); len(got) != 0 {
		t.Fatalf("a head with nothing to take, or a buy-in met: %v", got)
	}
	claimable, assigned = 3, 10
	w.write("MarginAccounts", "settle", func([]any, bool) error {
		return w.fail("MarginAccounts", "AssetWrittenOff", symbol(t, "SPY"))
	})
	failed = true
	if got := w.pass(t, k, "recall"); len(got) != 0 || !strings.Contains(w.logs.String(), "the recall queue is stopped") ||
		!strings.Contains(w.logs.String(), "the borrower failed a buy-in") {
		t.Fatalf("a stopped queue and a failed buy-in were not alerted: %v %s", got, w.logs.String())
	}
}
