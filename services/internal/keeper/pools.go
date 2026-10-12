// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"fmt"
	"maps"
	"math"
	"math/big"
	"slices"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// AlignBps is how far, in basis points, a pool's price may lie from the centre of its asset's band before the pools
// keeper trades it back.
const AlignBps = 25

var (
	maxUint256 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
	q192       = new(big.Int).Lsh(big.NewInt(1), 192)
)

// alignPools trades each pool of the registry's Uniswap v3 against USDG back to the centre of its asset's band from
// the keeper's own tokens, as an arbitrageur would, where the pool lies more than AlignBps from it and nothing else
// has. A pool the margin engine reads and a short sells into must lie near the band, and a chain without arbitrageurs,
// such as a testnet, leaves it where the last swap did. It trades all the keeper holds of the token it pays in, and
// stops where the pool's price reaches the band's centre. A halted band or a keeper holding none waits.
func (k *Keeper) alignPools(ctx context.Context) error {
	d := k.client.Deployments()
	usdg, ok := d.Tokens["USDG"]
	if !ok {
		return nil
	}
	for _, name := range slices.Sorted(maps.Keys(d.UniswapV3)) {
		asset, fee, ok := parsePool(name)
		token, listed := d.Tokens[asset]
		if !ok || !listed {
			continue
		}
		if err := k.alignPool(ctx, asset, token, usdg, d.UniswapV3[name], fee); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	return nil
}

func (k *Keeper) alignPool(ctx context.Context, asset string, token, usdg, pool common.Address, fee uint32) error {
	quote, err := k.client.Band().Quote(k.opts(ctx), asset)
	if err != nil || quote.State == 0 || quote.Mid == 0 {
		return err
	}
	mid := new(big.Int).SetUint64(quote.Mid)
	p := k.client.Pool(pool)
	slot0, err := p.Slot0(k.opts(ctx))
	if err != nil {
		return err
	}
	token0, _, err := p.Tokens(k.opts(ctx))
	if err != nil {
		return err
	}
	stockDecimals, err := k.client.Decimals(k.opts(ctx), token)
	if err != nil {
		return err
	}
	quoteDecimals, err := k.client.Decimals(k.opts(ctx), usdg)
	if err != nil {
		return err
	}
	stockIs0 := token0 == token
	price := poolPrice(slot0.SqrtPriceX96, stockIs0, stockDecimals, quoteDecimals)
	buy, far := deviation(price, mid)
	if !far {
		return nil
	}
	in, out := token, usdg
	if buy {
		in, out = usdg, token
	}
	held, err := k.client.BalanceOf(k.opts(ctx), in, k.From())
	if err != nil {
		return err
	}
	if held.Sign() == 0 {
		if !k.remember("align/"+asset, "empty") {
			k.log.Warn("the pool lies off the band, and the keeper holds none of the token that would trade it back",
				"asset", asset, "pool", pool.Hex(), "token", in.Hex())
		}
		return nil
	}
	router := k.client.SwapRouter()
	allowed, err := k.client.Allowance(k.opts(ctx), in, k.From(), k.client.Deployments().UniswapV3["SwapRouter02"])
	if err != nil {
		return err
	}
	if allowed.Cmp(held) < 0 {
		tx, err := k.client.Approve(in, k.client.Deployments().UniswapV3["SwapRouter02"], maxUint256)
		if _, err := k.act(ctx, "approve the router for "+in.Hex(), tx, err); err != nil {
			return err
		}
	}
	tx, err := router.ExactInputSingle(in, out, fee, k.From(), held, new(big.Int), targetSqrtPrice(mid, stockIs0, stockDecimals, quoteDecimals))
	moved, err := k.act(ctx, "trade "+asset+"'s pool back to the band", tx, err)
	if moved {
		k.log.Info("traded a pool back to the band", "asset", asset, "pool", pool.Hex(), "poolPrice", price.String(), "bandMid", mid.String())
	}
	return err
}

// parsePool reads a registry's pool name, <ASSET>_USDG_<fee>.
func parsePool(name string) (asset string, fee uint32, ok bool) {
	parts := strings.Split(name, "_")
	if len(parts) != 3 || parts[1] != "USDG" {
		return "", 0, false
	}
	n, err := strconv.ParseUint(parts[2], 10, 32)
	if err != nil || n == 0 || n > math.MaxUint32 {
		return "", 0, false
	}
	return parts[0], uint32(n), true
}

// poolPrice is the price of the stock in USD with 8 decimals, from the pool's √price as a Q64.96.
func poolPrice(sqrtPriceX96 *big.Int, stockIs0 bool, stockDecimals, quoteDecimals uint8) *big.Int {
	square := new(big.Int).Mul(sqrtPriceX96, sqrtPriceX96)
	if stockIs0 {
		num := new(big.Int).Mul(square, pow10(int(stockDecimals)+8))
		return num.Div(num, new(big.Int).Mul(q192, pow10(int(quoteDecimals))))
	}
	num := new(big.Int).Mul(q192, pow10(int(stockDecimals)+8))
	return num.Div(num, new(big.Int).Mul(square, pow10(int(quoteDecimals))))
}

// targetSqrtPrice is the pool's √price, as a Q64.96, at which the stock costs mid, in USD with 8 decimals.
func targetSqrtPrice(mid *big.Int, stockIs0 bool, stockDecimals, quoteDecimals uint8) *big.Int {
	num, den := new(big.Int).Mul(mid, pow10(int(quoteDecimals))), pow10(int(stockDecimals)+8)
	if !stockIs0 {
		num, den = den, num
	}
	return new(big.Int).Sqrt(num.Mul(num, q192).Div(num, den))
}

// deviation reports whether price lies more than AlignBps from mid, and whether the stock must be bought, rather than
// sold, to bring it back.
func deviation(price, mid *big.Int) (buy, far bool) {
	diff := new(big.Int).Sub(mid, price)
	limit := new(big.Int).Mul(mid, big.NewInt(AlignBps))
	return diff.Sign() > 0, new(big.Int).Mul(new(big.Int).Abs(diff), big.NewInt(10_000)).Cmp(limit) > 0
}

func pow10(n int) *big.Int {
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n)), nil)
}
