// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"fmt"
	"math/big"
	"slices"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk"
)

// account is a position's state in the margin world: what it holds by token, what it lent, whether it falls short,
// and its auction.
type account struct {
	held       map[common.Address]*big.Int
	lent       map[common.Address]*big.Int
	sellable   map[common.Address]*big.Int
	short      bool
	auction    uint64
	heldUntil  uint64
	inBaskets  *big.Int
	writtenOff bool
}

type pos struct {
	account  common.Address
	position [32]byte
}

// margin is a world with the margin accounts taking NVDA and SPY, and the liquidator over them.
func margin(t *testing.T) (*world, map[pos]*account) {
	w := pricedWorld(t)
	positions := map[pos]*account{}
	of := func(args []any) *account {
		p, ok := positions[pos{args[0].(common.Address), args[1].([32]byte)}]
		if !ok {
			t.Fatalf("an unknown position %v", args[:2])
		}
		return p
	}
	amount := func(m map[common.Address]*big.Int, token common.Address) *big.Int {
		if v, ok := m[token]; ok {
			return v
		}
		return n(0)
	}
	nvda, spy := w.deployments.Tokens["NVDA"], w.deployments.Tokens["SPY"]
	w.read("MarginAccounts", "stocks", func([]any) []any {
		return []any{[][32]byte{symbol(t, "NVDA"), symbol(t, "SPY")}, []common.Address{nvda, spy}}
	})
	w.read("MarginAccounts", "collateral", func(args []any) []any { return []any{amount(of(args).held, args[2].(common.Address))} })
	w.read("MarginAccounts", "lent", func(args []any) []any { return []any{amount(of(args).lent, args[2].(common.Address))} })
	w.read("MarginAccounts", "sellable", func(args []any) []any {
		return []any{amount(of(args).sellable, args[2].(common.Address))}
	})
	w.read("MarginAccounts", "backstop", func([]any) []any { return []any{common.Address{}} })
	w.read("MarginAccounts", "inBaskets", func(args []any) []any {
		p := of(args)
		if p.inBaskets == nil {
			return []any{[]*big.Int{n(0), n(0)}}
		}
		return []any{[]*big.Int{p.inBaskets, n(0)}}
	})
	w.read("Liquidator", "shortfall", func(args []any) []any { return []any{n(1), n(2), of(args).short, false} })
	w.read("Liquidator", "auctions", func(args []any) []any { return []any{of(args).auction, false} })
	w.read("Liquidator", "OPEN_AUCTION_LIFETIME", func([]any) []any { return []any{n(3600)} })
	w.read("Liquidator", "heldUntil", func(args []any) []any { return []any{of(args).heldUntil} })
	w.write("Liquidator", "start", func(args []any, mined bool) error {
		if mined {
			of(args).auction = w.time
		}
		return nil
	})
	w.write("Liquidator", "stop", func(args []any, mined bool) error {
		if mined {
			of(args).auction = 0
		}
		return nil
	})
	w.write("Liquidator", "recall", func(args []any, mined bool) error {
		if mined {
			of(args).lent[args[2].(common.Address)] = n(0)
		}
		return nil
	})
	w.write("Liquidator", "writeOff", func(args []any, mined bool) error {
		p := of(args)
		for _, v := range p.held {
			if v.Sign() > 0 {
				return w.fail("MarginAccounts", "PositionNotEmpty")
			}
		}
		if mined {
			p.writtenOff = true
		}
		return nil
	})
	for _, token := range []string{"NVDA", "SPY"} {
		w.read("Token:"+token, "paused", func([]any) []any { return []any{token == "SPY"} })
	}
	return w, positions
}

func TestPositionsAreLiquidatedDeepestFirstAndAHeldOneWaits(t *testing.T) {
	w, positions := margin(t)
	nvda, spy, usdg := w.deployments.Tokens["NVDA"], w.deployments.Tokens["SPY"], w.deployments.Tokens["USDG"]
	shallow, deep, held, recovered, empty, lending, frozen := pos{alice, sdk.Cross}, pos{bob, sdk.Cross},
		pos{alice, symbol(t, "NVDA")}, pos{bob, symbol(t, "NVDA")}, pos{alice, symbol(t, "SPY")},
		pos{bob, symbol(t, "SPY")}, pos{alice, symbol(t, "TSLA")}
	positions[shallow] = &account{held: map[common.Address]*big.Int{nvda: n(5)}, short: true}
	positions[deep] = &account{held: map[common.Address]*big.Int{nvda: n(5)}, short: true}
	positions[held] = &account{held: map[common.Address]*big.Int{nvda: n(5)}, short: true, heldUntil: w.time + 60}
	positions[recovered] = &account{held: map[common.Address]*big.Int{nvda: n(5)}, auction: w.time - 60}
	positions[empty] = &account{held: map[common.Address]*big.Int{}, short: true}
	positions[lending] = &account{held: map[common.Address]*big.Int{spy: n(1)}, short: true,
		lent: map[common.Address]*big.Int{spy: n(10)}, sellable: map[common.Address]*big.Int{spy: n(4)}}
	positions[frozen] = &account{held: map[common.Address]*big.Int{spy: n(1), usdg: n(7)}, short: true}
	for p, depth := range map[pos][2]int64{shallow: {90, 100}, deep: {10, 100}, held: {0, 100}, recovered: {200, 100},
		empty: {-5, 100}, lending: {50, 100}, frozen: {80, 100}} {
		w.debt(p.account, p.position, depth[0], depth[1])
	}
	w.write("Liquidator", "settleCash", func(args []any, mined bool) error {
		if mined {
			positions[pos{args[0].(common.Address), args[1].([32]byte)}].held[usdg] = n(0)
		}
		return nil
	})
	k := w.keeper(Config{})
	got := w.pass(t, k, "liquidation")
	start := func(p pos) string { return fmt.Sprintf("Liquidator.start[%s %v]", p.account.Hex(), p.position) }
	want := []string{
		fmt.Sprintf("Liquidator.writeOff[%s %v]", empty.account.Hex(), empty.position),
		start(deep),
		fmt.Sprintf("Liquidator.recall[%s %v %s]", lending.account.Hex(), lending.position, spy.Hex()),
		start(lending),
		start(frozen),
		fmt.Sprintf("Liquidator.settleCash[%s %v]", frozen.account.Hex(), frozen.position),
		start(shallow),
		fmt.Sprintf("Liquidator.stop[%s %v]", recovered.account.Hex(), recovered.position),
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("mined\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	w.debts = slices.DeleteFunc(w.debts, func(d map[string]any) bool {
		return d["account"] == empty.account.Hex() && d["position"] == fmt.Sprintf("%#x", empty.position)
	})
	restarted := w.keeper(Config{})
	if got := w.pass(t, restarted, "liquidation"); len(got) != 0 {
		t.Fatalf("a restarted keeper repeated %v", got)
	}
	w.time += 3600
	if got := w.pass(t, restarted, "liquidation"); len(got) != 5 || got[0] != start(held) || got[1] != start(deep) {
		t.Fatalf("expired auctions and the released hold not started: %v", got)
	}
}

func TestTheKeeperBuysOnlyOnceTheAskIsAtItsReference(t *testing.T) {
	w, positions := margin(t)
	nvda := w.deployments.Tokens["NVDA"]
	p := pos{alice, sdk.Cross}
	positions[p] = &account{held: map[common.Address]*big.Int{nvda: n(3e18)}, short: true, auction: w.time}
	w.debt(alice, sdk.Cross, 10, 100)
	ask := int64(23_000_000_000)
	w.read("Liquidator", "price", func([]any) []any { return []any{n(ask)} })
	w.read("Band", "quote", func([]any) []any {
		return []any{uint8(3), uint8(2), uint64(22_950_000_000), uint64(10), uint64(22_900_000_000), n(23_000_000_000)}
	})
	w.read("Token:USDG", "balanceOf", func([]any) []any { return []any{n(458_000_000)} })
	w.write("Token:USDG", "approve", func([]any, bool) error { return nil })
	var bought []any
	w.write("Liquidator", "buy", func(args []any, mined bool) error {
		if mined {
			bought = args
		}
		return nil
	})
	k := w.keeper(Config{Buy: true})
	if got := w.pass(t, k, "liquidation"); len(got) != 0 {
		t.Fatalf("bought above the low edge: %v", got)
	}
	ask = 22_900_000_000
	if got := w.pass(t, k, "liquidation"); len(got) != 2 || !strings.HasPrefix(got[0], "Token.approve") {
		t.Fatalf("no purchase at the low edge: %v", got)
	}
	if bought[2] != nvda || bought[3].(*big.Int).Cmp(n(2e18)) != 0 || bought[4].(*big.Int).Int64() != 458_000_000 ||
		bought[5] != k.From() {
		t.Fatalf("bought %v", bought)
	}
}

func TestTheBackstopClaimsItsPremiumAndCoversEmptiedPositions(t *testing.T) {
	w, positions := margin(t)
	nvda := w.deployments.Tokens["NVDA"]
	positions[pos{alice, sdk.Cross}] = &account{held: map[common.Address]*big.Int{nvda: n(2)}, short: true}
	positions[pos{bob, sdk.Cross}] = &account{held: map[common.Address]*big.Int{}, short: true}
	w.debt(alice, sdk.Cross, 10, 100)
	w.debt(bob, sdk.Cross, -10, 100)
	premium := int64(5)
	w.read("MarginAccounts", "backstopPremium", func([]any) []any { return []any{n(premium)} })
	w.write("GapBackstop", "claim", func(_ []any, mined bool) error {
		if mined {
			premium = 0
		}
		return nil
	})
	w.write("GapBackstop", "cover", func(args []any, _ bool) error {
		for _, v := range positions[pos{args[0].(common.Address), args[1].([32]byte)}].held {
			if v.Sign() > 0 {
				return w.fail("MarginAccounts", "PositionNotEmpty")
			}
		}
		return nil
	})
	k := w.keeper(Config{})
	got := w.pass(t, k, "backstop")
	if len(got) != 2 || got[0] != "GapBackstop.claim[]" || got[1] != fmt.Sprintf("GapBackstop.cover[%s %v]", bob.Hex(), sdk.Cross) {
		t.Fatalf("mined %v", got)
	}
}

func TestShortsInDeficitAreLiquidatedFirstThenTheDeepest(t *testing.T) {
	w := pricedWorld(t)
	add := func(account common.Address, asset string, usdg, equity, requirement int64) {
		w.shorts = append(w.shorts, map[string]any{"account": account.Hex(), "symbol": asset,
			"position": map[string]any{"usdgHeld": fmt.Sprint(usdg)},
			"health":   map[string]any{"equity": fmt.Sprint(equity), "requirement": fmt.Sprint(requirement)}})
	}
	add(alice, "NVDA", 100, 95, 100)
	add(bob, "NVDA", 100, 50, 100)
	add(alice, "SPY", -3, 150, 100)
	add(bob, "SPY", 100, 150, 100)
	recovered := map[string]bool{}
	w.write("ShortPositions", "liquidate", func(args []any, _ bool) error {
		if recovered[args[0].(common.Address).Hex()] {
			return w.fail("ShortPositions", "NotShortfall", n(95), n(100))
		}
		return nil
	})
	k := w.keeper(Config{})
	got := w.pass(t, k, "shorts")
	liquidated := func(account common.Address, asset string) string {
		return fmt.Sprintf("ShortPositions.liquidate[%s %v]", account.Hex(), symbol(t, asset))
	}
	if fmt.Sprint(got) != fmt.Sprint([]string{liquidated(alice, "SPY"), liquidated(bob, "NVDA"), liquidated(alice, "NVDA")}) {
		t.Fatalf("mined %v", got)
	}
	recovered[alice.Hex()] = true
	if got := w.pass(t, k, "shorts"); len(got) != 1 || got[0] != liquidated(bob, "NVDA") {
		t.Fatalf("a recovered short: %v", got)
	}
}

func TestABurnIsSyncedAndTheDustItLeavesCleared(t *testing.T) {
	w, positions := margin(t)
	w.vault("SPY")
	nvda, spy := w.deployments.Tokens["NVDA"], w.deployments.Tokens["SPY"]
	accounts, vault := w.deployments.Tapehouse["MarginAccounts"], w.deployments.StockLending["SPY"]
	balance := map[common.Address]map[common.Address]*big.Int{nvda: {accounts: n(100)}, spy: {accounts: n(50), vault: n(20)}}
	scale := map[string]*big.Int{"NVDA": n(1e18), "SPY": n(1e18)}
	units := map[string]*big.Int{"NVDA": n(100), "SPY": n(50)}
	counted := n(20)
	for _, asset := range []string{"NVDA", "SPY"} {
		token := w.deployments.Tokens[asset]
		w.read("Token:"+asset, "balanceOf", func(args []any) []any { return []any{balance[token][args[0].(common.Address)]} })
	}
	w.read("MarginAccounts", "holding", func(args []any) []any {
		asset := text(args[0].([32]byte))
		return []any{units[asset], scale[asset], n(1e18)}
	})
	w.read("Vault:SPY", "idle", func([]any) []any { return []any{counted} })
	w.read("Vault:SPY", "locked", func([]any) []any { return []any{n(0)} })
	w.write("MarginAccounts", "sync", func(args []any, mined bool) error {
		if mined {
			asset := text(args[0].([32]byte))
			token := w.deployments.Tokens[asset]
			scale[asset] = new(big.Int).Div(new(big.Int).Mul(balance[token][accounts], n(1e18)), units[asset])
			if asset == "SPY" {
				counted = balance[spy][vault]
			}
		}
		return nil
	})
	k := w.keeper(Config{})
	if got := w.pass(t, k, "sync"); len(got) != 0 {
		t.Fatalf("synced with nothing burnt: %v", got)
	}
	balance[spy][vault] = n(15)
	if got := w.pass(t, k, "sync"); len(got) != 1 || got[0] != fmt.Sprintf("MarginAccounts.sync[%v]", symbol(t, "SPY")) {
		t.Fatalf("a burn from the vault: %v", got)
	}
	balance[nvda][accounts] = n(0)
	positions[pos{alice, sdk.Cross}] = &account{held: map[common.Address]*big.Int{}}
	positions[pos{bob, symbol(t, "NVDA")}] = &account{held: map[common.Address]*big.Int{}}
	w.event("tapehouse.MarginAccounts", "Deposit", map[string]any{"account": alice.Hex(), "position": fmt.Sprintf("%#x", sdk.Cross), "token": nvda.Hex()})
	w.event("tapehouse.MarginAccounts", "Deposit", map[string]any{"account": bob.Hex(), "position": fmt.Sprintf("%#x", symbol(t, "NVDA")), "token": nvda.Hex()})
	w.event("tapehouse.MarginAccounts", "Deposit", map[string]any{"account": bob.Hex(), "position": fmt.Sprintf("%#x", symbol(t, "SPY")), "token": spy.Hex()})
	w.write("MarginAccounts", "clear", func([]any, bool) error { return nil })
	got := w.pass(t, k, "sync")
	want := []string{
		fmt.Sprintf("MarginAccounts.sync[%v]", symbol(t, "NVDA")),
		fmt.Sprintf("MarginAccounts.clear[%s %v %v]", alice.Hex(), sdk.Cross, symbol(t, "NVDA")),
		fmt.Sprintf("MarginAccounts.clear[%s %v %v]", bob.Hex(), symbol(t, "NVDA"), symbol(t, "NVDA")),
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("a burn of every NVDA: %v, want %v", got, want)
	}
	if got := w.pass(t, k, "sync"); len(got) != 0 {
		t.Fatalf("dust cleared twice: %v", got)
	}
}
