// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/swaprouter"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/uniswapv3pool"
)

type swap struct {
	TokenIn, TokenOut common.Address
	Fee, AmountIn     *big.Int
	AmountOutMinimum  *big.Int
	SqrtPriceLimitX96 *big.Int
	Recipient         common.Address
}

// poolWorld holds TSLA, a band that quotes it at mid, and its pool against USDG at 0.3%, priced at price, with the
// stock the pool's token0.
func poolWorld(t *testing.T, mid, price int64) (*world, *swap, map[common.Address]map[common.Address]*big.Int) {
	t.Helper()
	w := newWorld(t)
	w.asset("TSLA")
	usdg, tsla := w.deployments.Tokens["USDG"], w.deployments.Tokens["TSLA"]
	w.deployments.UniswapV3["SwapRouter02"] = w.add("SwapRouter02", &swaprouter.SwapRouterMetaData)
	w.deployments.UniswapV3["TSLA_USDG_3000"] = w.add("Pool", &uniswapv3pool.UniswapV3PoolMetaData)
	w.deployments.UniswapV3["Factory"] = common.Address{1}
	w.read("Band", "quote", func([]any) []any { return []any{uint8(3), uint8(2), uint64(mid), uint64(55), uint64(mid), n(mid)} })
	sqrt := targetSqrtPrice(n(price), true, 18, 6)
	w.read("Pool", "slot0", func([]any) []any { return []any{sqrt, n(0), uint16(0), uint16(1801), uint16(1801), uint8(0), true} })
	w.read("Pool", "token0", func([]any) []any { return []any{tsla} })
	w.read("Pool", "token1", func([]any) []any { return []any{usdg} })
	balances := map[common.Address]map[common.Address]*big.Int{usdg: {}, tsla: {}}
	allowances := map[common.Address]*big.Int{}
	for name, decimals := range map[string]uint8{"USDG": 6, "TSLA": 18} {
		token := w.deployments.Tokens[name]
		w.read("Token:"+name, "decimals", func([]any) []any { return []any{decimals} })
		w.read("Token:"+name, "balanceOf", func(args []any) []any {
			if b, ok := balances[token][args[0].(common.Address)]; ok {
				return []any{b}
			}
			return []any{n(0)}
		})
		w.read("Token:"+name, "allowance", func([]any) []any {
			if a, ok := allowances[token]; ok {
				return []any{a}
			}
			return []any{n(0)}
		})
		w.write("Token:"+name, "approve", func(args []any, mined bool) error {
			if mined {
				allowances[token] = args[1].(*big.Int)
			}
			return nil
		})
	}
	var seen swap
	w.write("SwapRouter02", "exactInputSingle", func(args []any, mined bool) error {
		if mined {
			raw, err := json.Marshal(args[0])
			if err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(raw, &seen); err != nil {
				t.Fatal(err)
			}
		}
		return nil
	})
	return w, &seen, balances
}

func TestAPoolOffTheBandIsTradedBackToItsCentreFromTheKeepersOwnTokens(t *testing.T) {
	const mid = 38_000_000_000
	w, seen, balances := poolWorld(t, mid, 37_700_000_000)
	usdg, tsla := w.deployments.Tokens["USDG"], w.deployments.Tokens["TSLA"]
	k := w.keeper(Config{Align: true})
	if !strings.Contains(strings.Join(k.Keepers(), ","), "pools") {
		t.Fatalf("the keepers are %v", k.Keepers())
	}
	balances[usdg][k.From()] = n(7_000_000)

	if got := w.pass(t, k, "pools"); len(got) != 2 || !strings.HasPrefix(got[0], "Token.approve") ||
		!strings.HasPrefix(got[1], "SwapRouter02.exactInputSingle") {
		t.Fatalf("a pool below the band: %v", got)
	}
	if seen.TokenIn != usdg || seen.TokenOut != tsla || seen.Fee.Int64() != 3000 || seen.AmountIn.Int64() != 7_000_000 ||
		seen.Recipient != k.From() || seen.AmountOutMinimum.Sign() != 0 ||
		seen.SqrtPriceLimitX96.Cmp(targetSqrtPrice(n(mid), true, 18, 6)) != 0 {
		t.Fatalf("a buy of the stock: %+v", seen)
	}
	if got := w.pass(t, k, "pools"); len(got) != 1 {
		t.Fatalf("the router is approved once: %v", got)
	}
}

func TestAPoolAboveTheBandSellsTheStock(t *testing.T) {
	w, seen, balances := poolWorld(t, 38_000_000_000, 38_300_000_000)
	usdg, tsla := w.deployments.Tokens["USDG"], w.deployments.Tokens["TSLA"]
	k := w.keeper(Config{Align: true})
	balances[tsla][k.From()] = n(2e16)
	balances[usdg][k.From()] = n(9_000_000)
	if got := w.pass(t, k, "pools"); len(got) != 2 {
		t.Fatalf("a pool above the band: %v", got)
	}
	if seen.TokenIn != tsla || seen.TokenOut != usdg || seen.AmountIn.Cmp(n(2e16)) != 0 {
		t.Fatalf("a sale of the stock: %+v", seen)
	}
}

func TestAPoolNearTheBandAHaltedBandAndAnEmptyKeeperAreLeftAlone(t *testing.T) {
	w, _, balances := poolWorld(t, 38_000_000_000, 37_950_000_000)
	k := w.keeper(Config{Align: true})
	balances[w.deployments.Tokens["USDG"]][k.From()] = n(5_000_000)
	if got := w.pass(t, k, "pools"); len(got) != 0 {
		t.Fatalf("a pool within %d bps: %v", AlignBps, got)
	}

	w, _, _ = poolWorld(t, 38_000_000_000, 37_000_000_000)
	k = w.keeper(Config{Align: true})
	if got := w.pass(t, k, "pools"); len(got) != 0 || !strings.Contains(w.logs.String(), "holds none") {
		t.Fatalf("a keeper holding none: %v", got)
	}

	w, _, balances = poolWorld(t, 38_000_000_000, 37_000_000_000)
	w.read("Band", "quote", func([]any) []any { return []any{uint8(0), uint8(0), uint64(0), uint64(0), uint64(0), n(0)} })
	k = w.keeper(Config{Align: true})
	balances[w.deployments.Tokens["USDG"]][k.From()] = n(5_000_000)
	if got := w.pass(t, k, "pools"); len(got) != 0 {
		t.Fatalf("a halted band: %v", got)
	}
}

func TestThePoolsKeeperRunsOnlyWhereItIsAskedFor(t *testing.T) {
	w, _, _ := poolWorld(t, 38_000_000_000, 37_700_000_000)
	if strings.Contains(strings.Join(w.keeper(Config{}).Keepers(), ","), "pools") {
		t.Fatal("the pools keeper runs without being asked for")
	}
	env := map[string]string{"TAPEHOUSE_DEPLOYMENTS": "d.json", "TAPEHOUSE_RPC_URL": "http://rpc", "TAPEHOUSE_KEEPERS": "pools"}
	if _, err := ConfigFromEnv(func(key string) string { return env[key] }); err == nil || !strings.Contains(err.Error(), "TAPEHOUSE_KEEPER_ALIGN") {
		t.Fatalf("the pools keeper named without its consent: %v", err)
	}
	env["TAPEHOUSE_KEEPER_ALIGN"] = "true"
	if cfg, err := ConfigFromEnv(func(key string) string { return env[key] }); err != nil || !cfg.Align {
		t.Fatalf("the pools keeper with its consent: %v", err)
	}
}

func TestAPoolsPriceAndItsTargetAreInverses(t *testing.T) {
	for _, stockIs0 := range []bool{true, false} {
		for _, price := range []int64{26_000_000_000, 37_950_000_000, 1_234_567_890_000} {
			sqrt := targetSqrtPrice(n(price), stockIs0, 18, 6)
			got := poolPrice(sqrt, stockIs0, 18, 6)
			diff := new(big.Int).Abs(new(big.Int).Sub(got, n(price)))
			if diff.Cmp(n(price/1_000_000_000+2)) > 0 {
				t.Fatalf("stock is token0 %v: %d became %v", stockIs0, price, got)
			}
		}
	}
	buy, far := deviation(n(37_700_000_000), n(38_000_000_000))
	sell, farUp := deviation(n(38_300_000_000), n(38_000_000_000))
	near, farNear := deviation(n(37_950_000_000), n(38_000_000_000))
	if !buy || !far || sell || !farUp || !buy == near || farNear {
		t.Fatalf("deviations: %v %v, %v %v, %v %v", buy, far, sell, farUp, near, farNear)
	}
	for name, want := range map[string]bool{"TSLA_USDG_3000": true, "SPY_WETH_500": false, "Factory": false, "TSLA_USDG_0": false} {
		if _, _, ok := parsePool(name); ok != want {
			t.Fatalf("%s: %v", name, ok)
		}
	}
}
