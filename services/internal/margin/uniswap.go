// SPDX-License-Identifier: MIT OR Apache-2.0

package margin

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Pool arithmetic of the liquidity add-on.
const (
	Window     = 1_800
	SellFactor = 51_316_701
	BuyFactor  = 48_808_848
	MinTick    = -887_272
	MaxTick    = 887_272
)

var (
	two160  = new(big.Int).Lsh(big.NewInt(1), 160)
	two256  = new(big.Int).Lsh(big.NewInt(1), 256)
	maxUint = new(big.Int).Sub(two256, big.NewInt(1))
)

// Classify reports whether the asset's stock token is the pool's token0 and whether the other token is WETH, and
// whether the pool trades stock against usdg or weth at all.
func Classify(token0, token1, stock, usdg, weth common.Address) (stockIsToken0, quoteIsWeth, ok bool) {
	quote := func(t common.Address) bool { return t != (common.Address{}) && (t == usdg || t == weth) }
	switch {
	case token0 == stock && token1 != stock && quote(token1):
		return true, token1 == weth, true
	case token1 == stock && token0 != stock && quote(token0):
		return false, token0 == weth, true
	}
	return false, false, false
}

// Mean is Uniswap's OracleLibrary.consult over two observations Window seconds apart, with the cumulatives wrapping
// as the pool's int56 and uint160 do: the arithmetic-mean tick and the harmonic-mean liquidity.
func Mean(tickThen, tickNow int64, secondsThen, secondsNow *big.Int) (int32, *big.Int, bool) {
	delta := (tickNow - tickThen) << 8 >> 8
	tick := delta / Window
	if delta < 0 && delta%Window != 0 {
		tick--
	}
	if tick < MinTick || tick > MaxTick {
		return 0, nil, false
	}
	perLiquidity := new(big.Int).Sub(secondsNow, secondsThen)
	perLiquidity.Mod(perLiquidity, two160).Lsh(perLiquidity, 32)
	if perLiquidity.Sign() == 0 {
		return 0, nil, false
	}
	liquidity := new(big.Int).Sub(two160, big.NewInt(1))
	liquidity.Mul(liquidity, big.NewInt(Window)).Quo(liquidity, perLiquidity)
	if liquidity.BitLen() > 128 {
		return 0, nil, false
	}
	return int32(tick), liquidity, true
}

var tickFactors = []string{
	"fff97272373d413259a46990580e213a", "fff2e50f5f656932ef12357cf3c7fdcc", "ffe5caca7e10e4e61c3624eaa0941cd0",
	"ffcb9843d60f6159c9db58835c926644", "ff973b41fa98c081472e6896dfb254c0", "ff2ea16466c96a3843ec78b326b52861",
	"fe5dee046a99a2a811c461f1969c3053", "fcbe86c7900a88aedcffc83b479aa3a4", "f987a7253ac413176f2b074cf7815e54",
	"f3392b0822b70005940c7a398e4b70f3", "e7159475a2c29b7443b29c7fa6e889d9", "d097f3bdfd2022b8845ad8f792aa5825",
	"a9f746462d870fdf8a65dc1f90e061e5", "70d869a156d2a1b890bb3df62baf32f7", "31be135f97d08fd981231505542fcfa6",
	"9aa508b5b7a84e1c677de54f3e99bc9", "5d6af8dedb81196699c329225ee604", "2216e584f5fa1ea926041bedfe98",
	"48a170391f7dc42444e8fa2",
}

func hexInt(s string) *big.Int {
	v, _ := new(big.Int).SetString(s, 16)
	return v
}

// SPDX-SnippetBegin
// SPDX-SnippetCopyrightText: 2023 Universal Navigation Inc.
// SPDX-License-Identifier: MIT

// SqrtRatioAtTick is √(1.0001^tick) as a Q64.96, exactly as Uniswap v4-core's TickMath.getSqrtPriceAtTick, whose
// factors this uses.
func SqrtRatioAtTick(tick int32) *big.Int {
	abs := int64(tick)
	if abs < 0 {
		abs = -abs
	}
	ratio := new(big.Int).Lsh(big.NewInt(1), 128)
	if abs&1 != 0 {
		ratio = hexInt("fffcb933bd6fad37aa2d162d1a594001")
	}
	for bit, f := range tickFactors {
		if abs&(2<<bit) != 0 {
			ratio.Mul(ratio, hexInt(f)).Rsh(ratio, 128)
		}
	}
	if tick > 0 {
		ratio.Quo(maxUint, ratio)
	}
	rounded := new(big.Int).And(ratio, big.NewInt(1<<32-1)).Sign() != 0
	ratio.Rsh(ratio, 32)
	if rounded {
		ratio.Add(ratio, big.NewInt(1))
	}
	return ratio
}

// SPDX-SnippetEnd

// Terms are the asset's price in USD with 8 decimals from a pool's mean tick and liquidity, and the USD with 18
// decimals the pool pays out for a 10% fall and takes in for a 10% rise; usd is the quote token's price with 8
// decimals. The fee is left to the caller.
func Terms(stockIsToken0 bool, stockDecimals, quoteDecimals uint8, tick int32, liquidity, usd *big.Int) (PoolTerms, bool) {
	if !stockIsToken0 {
		tick = -tick
	}
	sqrt := SqrtRatioAtTick(tick)
	stockUnit := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(stockDecimals)), nil)
	quoteUnit := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(quoteDecimals)), nil)
	perStock := mulShr(mulShr(sqrt, sqrt, 64), stockUnit, 128)
	price := new(big.Int).Mul(perStock, usd)
	if price.Cmp(maxUint) > 0 {
		return PoolTerms{}, false
	}
	price.Quo(price, quoteUnit)
	moved := mulShr(liquidity, sqrt, 96)
	scale := new(big.Int).Quo(bigWad, quoteUnit)
	value := func(factor int64) (*big.Int, bool) {
		raw := new(big.Int).Mul(moved, big.NewInt(factor))
		if raw.Cmp(maxUint) > 0 {
			return nil, false
		}
		raw.Quo(raw, big.NewInt(1_000_000_000)).Mul(raw, usd)
		if raw.Cmp(maxUint) > 0 {
			return nil, false
		}
		if raw.Mul(raw, scale); raw.Cmp(maxUint) > 0 {
			return nil, false
		}
		return raw.Quo(raw, big.NewInt(100_000_000)), true
	}
	selling, ok := value(SellFactor)
	if !ok {
		return PoolTerms{}, false
	}
	buying, ok := value(BuyFactor)
	return PoolTerms{Price: price, Selling: selling, Buying: buying}, ok
}

func mulShr(a, b *big.Int, shift uint) *big.Int {
	out := new(big.Int).Mul(a, b)
	return out.Rsh(out, shift).Mod(out, two256)
}
