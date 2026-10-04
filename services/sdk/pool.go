// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/uniswapv3pool"
)

// Pool reads a Uniswap v3 pool's price, liquidity and initialized ticks.
type Pool struct {
	target
	c    *Client
	pool *uniswapv3pool.UniswapV3Pool
}

// Pool returns the Uniswap v3 pool at address.
func (c *Client) Pool(address common.Address) *Pool {
	return &Pool{target: target{address: address}, c: c, pool: uniswapv3pool.NewUniswapV3Pool()}
}

// Slot0 reads the pool's current price and tick.
func (p *Pool) Slot0(opts *bind.CallOpts) (uniswapv3pool.Slot0Output, error) {
	return read(p.c, opts, p.target, p.pool.UnpackSlot0)(p.pool.TryPackSlot0())
}

// Liquidity reads the liquidity in range at the current tick.
func (p *Pool) Liquidity(opts *bind.CallOpts) (*big.Int, error) {
	return read(p.c, opts, p.target, p.pool.UnpackLiquidity)(p.pool.TryPackLiquidity())
}

// TickSpacing reads the distance between the pool's initializable ticks.
func (p *Pool) TickSpacing(opts *bind.CallOpts) (int32, error) {
	spacing, err := read(p.c, opts, p.target, p.pool.UnpackTickSpacing)(p.pool.TryPackTickSpacing())
	if err != nil {
		return 0, err
	}
	return int32(spacing.Int64()), nil
}

// Fee reads the pool's fee in hundredths of a basis point.
func (p *Pool) Fee(opts *bind.CallOpts) (uint32, error) {
	fee, err := read(p.c, opts, p.target, p.pool.UnpackFee)(p.pool.TryPackFee())
	if err != nil {
		return 0, err
	}
	return uint32(fee.Uint64()), nil
}

// Tokens reads the pool's two tokens, the lower address first.
func (p *Pool) Tokens(opts *bind.CallOpts) (token0, token1 common.Address, err error) {
	if token0, err = read(p.c, opts, p.target, p.pool.UnpackToken0)(p.pool.TryPackToken0()); err != nil {
		return token0, token1, err
	}
	token1, err = read(p.c, opts, p.target, p.pool.UnpackToken1)(p.pool.TryPackToken1())
	return token0, token1, err
}

// TickBitmap reads the word of the pool's initialized-tick bitmap at word.
func (p *Pool) TickBitmap(opts *bind.CallOpts, word int16) (*big.Int, error) {
	return read(p.c, opts, p.target, p.pool.UnpackTickBitmap)(p.pool.TryPackTickBitmap(word))
}

// LiquidityNet reads the liquidity added to the pool as its price crosses tick upward.
func (p *Pool) LiquidityNet(opts *bind.CallOpts, tick int32) (*big.Int, error) {
	out, err := read(p.c, opts, p.target, p.pool.UnpackTicks)(p.pool.TryPackTicks(big.NewInt(int64(tick))))
	return out.LiquidityNet, err
}

// Observe reads the pool's cumulative tick and seconds-per-liquidity at each of secondsAgos seconds before the block.
func (p *Pool) Observe(opts *bind.CallOpts, secondsAgos []uint32) (uniswapv3pool.ObserveOutput, error) {
	return read(p.c, opts, p.target, p.pool.UnpackObserve)(p.pool.TryPackObserve(secondsAgos))
}
