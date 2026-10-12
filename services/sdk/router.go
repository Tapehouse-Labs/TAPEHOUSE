// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/swaprouter"
)

// SwapRouter packs the swaps of the registry's Uniswap v3 router.
type SwapRouter struct {
	target
	router *swaprouter.SwapRouter
}

// SwapRouter returns the router of the registry's .uniswapV3.SwapRouter02.
func (c *Client) SwapRouter() *SwapRouter {
	return &SwapRouter{target: lookup(c.deployments.UniswapV3, "SwapRouter02", ".uniswapV3"), router: swaprouter.NewSwapRouter()}
}

// ExactInputSingle swaps amountIn of tokenIn for tokenOut through the pool at fee, in hundredths of a basis point,
// paying recipient at least minOut, and stopping where the pool's price reaches sqrtPriceLimitX96 (zero for no limit).
func (r *SwapRouter) ExactInputSingle(tokenIn, tokenOut common.Address, fee uint32, recipient common.Address, amountIn, minOut, sqrtPriceLimitX96 *big.Int) (Tx, error) {
	if err := present(amountIn); err != nil {
		return Tx{}, err
	}
	return r.tx(r.router.TryPackExactInputSingle(swaprouter.IV3SwapRouterExactInputSingleParams{
		TokenIn: tokenIn, TokenOut: tokenOut, Fee: new(big.Int).SetUint64(uint64(fee)), Recipient: recipient,
		AmountIn: amountIn, AmountOutMinimum: minOut, SqrtPriceLimitX96: sqrtPriceLimitX96,
	}))
}
