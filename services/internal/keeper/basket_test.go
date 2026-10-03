// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestABasketIsRebalancedOnlyWhenThePoolsPayOrTheTeamDoes(t *testing.T) {
	w := pricedWorld(t)
	w.basket("PAIR")
	nvda, spy := w.deployments.Tokens["NVDA"], w.deployments.Tokens["SPY"]
	basket := w.deployments.Baskets["PAIR"]
	w.read("Basket:PAIR", "components", func([]any) []any {
		return []any{[][32]byte{symbol(t, "NVDA"), symbol(t, "SPY")}, []common.Address{nvda, spy}}
	})
	w.read("Basket:PAIR", "target", func([]any) []any { return []any{[]*big.Int{n(3e18), n(1e18)}} })
	w.read("Basket:PAIR", "totalSupply", func([]any) []any { return []any{n(10)} })
	balances := map[common.Address]map[common.Address]*big.Int{nvda: {basket: n(20)}, spy: {basket: n(20)}}
	keeperAddress := common.Address{}
	for _, asset := range []string{"NVDA", "SPY"} {
		token := w.deployments.Tokens[asset]
		w.read("Token:"+asset, "balanceOf", func(args []any) []any {
			if b, ok := balances[token][args[0].(common.Address)]; ok {
				return []any{b}
			}
			return []any{n(0)}
		})
		w.write("Token:"+asset, "approve", func([]any, bool) error { return nil })
	}
	w.read("Band", "quote", func(args []any) []any {
		if text(args[0].([32]byte)) == "NVDA" {
			return []any{uint8(3), uint8(2), uint64(100), uint64(10), uint64(99), n(101)}
		}
		return []any{uint8(3), uint8(2), uint64(300), uint64(10), uint64(297), n(303)}
	})
	w.write("Basket:PAIR", "rebalance", func(args []any, _ bool) error {
		in, out := args[0].([]*big.Int), args[1].([]*big.Int)
		if in[0].Int64() != 10 || in[1].Sign() != 0 || out[0].Sign() != 0 || out[1].Int64() != 3 || args[2] != keeperAddress {
			t.Fatalf("rebalanced %v", args)
		}
		return nil
	})
	w.read("ShortPositions", "fee", func([]any) []any { return []any{n(500)} })
	proceeds := int64(800)
	w.read("QuoterV2", "quoteExactInputSingle", func([]any) []any { return []any{n(proceeds), n(0), uint32(0), n(0)} })
	w.read("QuoterV2", "quoteExactOutputSingle", func([]any) []any { return []any{n(1000), n(0), uint32(0), n(0)} })
	k := w.keeper(Config{})
	keeperAddress = k.From()
	balances[nvda][keeperAddress] = n(50)
	if got := w.pass(t, k, "basket"); len(got) != 0 || !strings.Contains(w.logs.String(), "the pools do not pay") {
		t.Fatalf("a rebalance the pools do not pay: %v", got)
	}
	proceeds = 1100
	if got := w.pass(t, k, "basket"); len(got) != 2 || !strings.HasPrefix(got[1], "Basket.rebalance") {
		t.Fatalf("a rebalance the pools pay: %v", got)
	}
	proceeds = 0
	subsidised := w.keeper(Config{Rebalance: true})
	if got := w.pass(t, subsidised, "basket"); len(got) != 2 {
		t.Fatalf("a rebalance the team pays: %v", got)
	}
}
