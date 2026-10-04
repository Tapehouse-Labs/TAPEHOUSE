// SPDX-License-Identifier: MIT OR Apache-2.0

package calibrate

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"slices"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/internal/margin"
	"github.com/tapehouse/tapehouse/services/sdk"
)

// The tokens' decimals: Stock Tokens and WETH have 18, USDG 6.
const (
	stockDecimals = 18
	usdgDecimals  = 6
	wethDecimals  = 18
	// DepthMove is the move, in percent, of an asset's price a depth is measured to.
	DepthMove = 10
	// SnapshotDays is how many days of snapshots a depth is the lowest of.
	SnapshotDays   = 28
	feeDenominator = 1_000_000
	// reach is the ticks a 10% move crosses either way, rounded up.
	reach = 1_100
)

var (
	q96       = new(big.Int).Lsh(big.NewInt(1), 96)
	chainlink = big.NewInt(100_000_000)
)

// PoolState is a Uniswap v3 pool as a walk needs it: its price, liquidity, fee and tick spacing, which of its tokens
// is the stock, the quote token's price in USD, and the liquidity net of every initialized tick near the price.
type PoolState struct {
	Asset         string             `json:"asset"`
	Pool          common.Address     `json:"pool"`
	SqrtPriceX96  *big.Int           `json:"sqrtPriceX96"`
	Liquidity     *big.Int           `json:"liquidity"`
	Tick          int32              `json:"tick"`
	Spacing       int32              `json:"spacing"`
	Fee           uint32             `json:"fee"`
	StockIsToken0 bool               `json:"stockIsToken0"`
	StockDecimals uint8              `json:"stockDecimals"`
	QuoteDecimals uint8              `json:"quoteDecimals"`
	USD           string             `json:"usd"`
	Nets          map[int32]*big.Int `json:"nets"`
}

// Snapshot is the engine's pools at one block.
type Snapshot struct {
	Block uint64      `json:"block"`
	Time  uint64      `json:"time"`
	Pools []PoolState `json:"pools"`
}

// PoolRef names the pool an asset is liquidated in.
type PoolRef struct {
	Asset   string
	Address common.Address
}

// TakeSnapshot reads the pool of each of assets, as the engine names it, at the latest block, with Chainlink's
// ETH/USD for a WETH quote and the liquidity net of every initialized tick within 10% of the price either way.
func TakeSnapshot(ctx context.Context, client *sdk.Client, assets []string) (Snapshot, error) {
	header, err := client.Header(ctx, nil)
	if err != nil {
		return Snapshot{}, err
	}
	opts := &bind.CallOpts{Context: ctx, BlockNumber: header.Number}
	refs := make([]PoolRef, len(assets))
	for i, name := range assets {
		symbol, err := sdk.ToBytes32(name)
		if err != nil {
			return Snapshot{}, err
		}
		address, err := client.Margin().Pool(opts, symbol)
		if err != nil {
			return Snapshot{}, fmt.Errorf("pool of %s: %w", name, err)
		}
		refs[i] = PoolRef{name, address}
	}
	return SnapshotPools(ctx, client, header.Number, header.Time, refs)
}

// SnapshotPools reads pools at block, whose timestamp is at.
func SnapshotPools(ctx context.Context, client *sdk.Client, block *big.Int, at uint64, pools []PoolRef) (Snapshot, error) {
	opts := &bind.CallOpts{Context: ctx, BlockNumber: block}
	d := client.Deployments()
	snapshot := Snapshot{Block: block.Uint64(), Time: at}
	var eth string
	for _, ref := range pools {
		p := client.Pool(ref.Address)
		token0, token1, err := p.Tokens(opts)
		if err != nil {
			return Snapshot{}, err
		}
		stockIsToken0, quoteIsWeth, ok := margin.Classify(token0, token1, d.Tokens[ref.Asset], d.Tokens["USDG"], d.Tokens["WETH"])
		if !ok {
			return Snapshot{}, fmt.Errorf("the pool of %s does not trade it against USDG or WETH", ref.Asset)
		}
		state := PoolState{Asset: ref.Asset, Pool: ref.Address, StockIsToken0: stockIsToken0, StockDecimals: stockDecimals, QuoteDecimals: usdgDecimals, USD: "1"}
		if quoteIsWeth {
			if eth == "" {
				if eth, err = ethUSD(opts, client); err != nil {
					return Snapshot{}, err
				}
			}
			state.QuoteDecimals, state.USD = wethDecimals, eth
		}
		slot0, err := p.Slot0(opts)
		if err != nil {
			return Snapshot{}, err
		}
		state.SqrtPriceX96, state.Tick = slot0.SqrtPriceX96, int32(slot0.Tick.Int64())
		if state.Liquidity, err = p.Liquidity(opts); err != nil {
			return Snapshot{}, err
		}
		if state.Spacing, err = p.TickSpacing(opts); err != nil {
			return Snapshot{}, err
		}
		if state.Fee, err = p.Fee(opts); err != nil {
			return Snapshot{}, err
		}
		if state.Nets, err = nets(opts, p, state); err != nil {
			return Snapshot{}, err
		}
		snapshot.Pools = append(snapshot.Pools, state)
	}
	return snapshot, nil
}

func ethUSD(opts *bind.CallOpts, client *sdk.Client) (string, error) {
	feed, err := client.Deployments().TokenPriceFeed("ETH_USD")
	if err != nil {
		return "", err
	}
	round, err := client.TokenPrice(opts, feed)
	if err != nil {
		return "", err
	}
	return new(big.Rat).SetFrac(round.Answer, chainlink).FloatString(8), nil
}

// nets reads the liquidity net of the initialized ticks within reach of the pool's tick, through the words of its
// tick bitmap that cover them.
func nets(opts *bind.CallOpts, p *sdk.Pool, s PoolState) (map[int32]*big.Int, error) {
	low, high := floorDiv(s.Tick-reach-s.Spacing, s.Spacing), floorDiv(s.Tick+reach+s.Spacing, s.Spacing)
	out := map[int32]*big.Int{}
	for word := floorDiv(low, 256); word <= floorDiv(high, 256); word++ {
		bits, err := p.TickBitmap(opts, int16(word))
		if err != nil {
			return nil, err
		}
		for bit := range 256 {
			if bits.Bit(bit) == 0 {
				continue
			}
			tick := (word*256 + int32(bit)) * s.Spacing
			net, err := p.LiquidityNet(opts, tick)
			if err != nil {
				return nil, err
			}
			out[tick] = net
		}
	}
	return out, nil
}

func floorDiv(a, b int32) int32 {
	q := a / b
	if a%b != 0 && (a < 0) != (b < 0) {
		q--
	}
	return q
}

func mulDivUp(a, b, d *big.Int) *big.Int {
	v := new(big.Int).Mul(a, b)
	v, m := v.QuoRem(v, d, new(big.Int))
	if m.Sign() != 0 {
		v.Add(v, big.NewInt(1))
	}
	return v
}

// amount0 and amount1 are Uniswap v3's SqrtPriceMath.getAmount0Delta and getAmount1Delta.
func amount0(a, b, liquidity *big.Int, up bool) *big.Int {
	if a.Cmp(b) > 0 {
		a, b = b, a
	}
	numerator := new(big.Int).Lsh(liquidity, 96)
	gap := new(big.Int).Sub(b, a)
	if up {
		v := mulDivUp(numerator, gap, b)
		return v.Add(v, new(big.Int).Sub(a, big.NewInt(1))).Quo(v, a)
	}
	v := new(big.Int).Mul(numerator, gap)
	return v.Quo(v, b).Quo(v, a)
}

func amount1(a, b, liquidity *big.Int, up bool) *big.Int {
	if a.Cmp(b) > 0 {
		a, b = b, a
	}
	gap := new(big.Int).Sub(b, a)
	if up {
		return mulDivUp(liquidity, gap, q96)
	}
	v := new(big.Int).Mul(liquidity, gap)
	return v.Quo(v, q96)
}

// Absorbed is the USD, to the whole dollar, the pool pays out when the stock is sold into it until its price has
// fallen movePct percent, or takes in, with the fee, when the stock is bought from it until the price has risen
// movePct percent: Uniswap v3's swap, tick by tick, as computeSwapStep does it for an input that reaches its limit.
func Absorbed(p PoolState, movePct int, selling bool) *big.Int {
	factor := big.NewInt(int64(100 + movePct))
	against := big.NewInt(100)
	if selling {
		factor = big.NewInt(int64(100 - movePct))
	}
	if !p.StockIsToken0 {
		factor, against = against, factor
	}
	square := new(big.Int).Mul(p.SqrtPriceX96, p.SqrtPriceX96)
	target := new(big.Int).Sqrt(square.Mul(square, factor).Quo(square, against))
	zeroForOne := selling == p.StockIsToken0
	ticks := make([]int32, 0, len(p.Nets))
	for tick := range p.Nets {
		ticks = append(ticks, tick)
	}
	slices.Sort(ticks)
	sqrt, liquidity, tick := new(big.Int).Set(p.SqrtPriceX96), new(big.Int).Set(p.Liquidity), p.Tick
	fee, rest := big.NewInt(int64(p.Fee)), big.NewInt(feeDenominator-int64(p.Fee))
	total := new(big.Int)
	for sqrt.Cmp(target) != 0 {
		next, found := int32(margin.MinTick), false
		if zeroForOne {
			if i, _ := slices.BinarySearch(ticks, tick+1); i > 0 {
				next, found = ticks[i-1], true
			}
		} else {
			next = margin.MaxTick
			if i, _ := slices.BinarySearch(ticks, tick+1); i < len(ticks) {
				next, found = ticks[i], true
			}
		}
		edge := margin.SqrtRatioAtTick(next)
		stop := edge
		if (zeroForOne && edge.Cmp(target) < 0) || (!zeroForOne && edge.Cmp(target) > 0) {
			stop = target
		}
		var in, out *big.Int
		if zeroForOne {
			in, out = amount0(stop, sqrt, liquidity, true), amount1(stop, sqrt, liquidity, false)
		} else {
			in, out = amount1(sqrt, stop, liquidity, true), amount0(sqrt, stop, liquidity, false)
		}
		if selling {
			total.Add(total, out)
		} else {
			total.Add(total, in).Add(total, mulDivUp(in, fee, rest))
		}
		sqrt = stop
		if sqrt.Cmp(edge) == 0 && found {
			net := p.Nets[next]
			if zeroForOne {
				liquidity.Sub(liquidity, net)
				tick = next - 1
			} else {
				liquidity.Add(liquidity, net)
				tick = next
			}
		}
	}
	usd, _ := new(big.Rat).SetString(p.USD)
	usd.Mul(usd, new(big.Rat).SetFrac(total, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(p.QuoteDecimals)), nil)))
	return new(big.Int).Quo(usd.Num(), usd.Denom())
}

func weekend(unix uint64) bool {
	day := time.Unix(int64(unix), 0).UTC().Weekday()
	return day == time.Saturday || day == time.Sunday
}

// Depths are each asset's selling then buying depth: the lowest Absorbed over a DepthMove% move among the snapshots
// of the last SnapshotDays days before now. It reports false, and no depths, where none of them was taken on a
// weekend, or one holds no pool for an asset.
func Depths(snapshots []Snapshot, now uint64, assets []string) ([]uint32, bool) {
	from := now - min(now, SnapshotDays*86_400)
	depths := make([]uint32, 2*len(assets))
	for i := range depths {
		depths[i] = math.MaxUint32
	}
	var weekends int
	for _, s := range snapshots {
		if s.Time < from || s.Time > now {
			continue
		}
		if weekend(s.Time) {
			weekends++
		}
		for i, name := range assets {
			j := slices.IndexFunc(s.Pools, func(p PoolState) bool { return p.Asset == name })
			if j < 0 {
				return nil, false
			}
			for k, selling := range []bool{true, false} {
				v := Absorbed(s.Pools[j], DepthMove, selling)
				depths[2*i+k] = min(depths[2*i+k], uint32(min(v.Uint64(), math.MaxUint32)))
			}
		}
	}
	if weekends == 0 {
		return nil, false
	}
	return depths, true
}
