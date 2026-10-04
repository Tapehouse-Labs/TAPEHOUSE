// SPDX-License-Identifier: MIT OR Apache-2.0

package alerts

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/internal/band"
	"github.com/tapehouse/tapehouse/services/internal/margin"
	"github.com/tapehouse/tapehouse/services/sdk"
	bandbinding "github.com/tapehouse/tapehouse/services/sdk/bindings/band"
)

const (
	scenarioSize  = 256
	feedMaxAgeS   = 86_400 + 60
	stockDecimals = 18
	wethDecimals  = 18
	usdgDecimals  = 6
)

// Reader reads markets and positions from the chain through the Go SDK, every call at one block.
type Reader struct {
	client *sdk.Client
}

// NewReader returns a Reader over client.
func NewReader(client *sdk.Client) *Reader {
	return &Reader{client}
}

func name(symbol [32]byte) string {
	return string(bytes.TrimRight(symbol[:], "\x00"))
}

// closedAt is the liquidator's _closed: the session is not open, its boundary has passed, or NYSE is neither in
// regular hours nor between two trading days.
func closedAt(s bandbinding.SessionOutput, nowMs uint64) bool {
	if s.State != 2 || (s.BoundaryMs != 0 && s.BoundaryMs <= nowMs) {
		return true
	}
	return s.Nyse != 1 && (s.Nyse != 2 || s.NyseNext != 1)
}

// Market reads the band, the margin engine and WETH's price at block number, or at the latest block where number is
// nil.
func (r *Reader) Market(ctx context.Context, number *big.Int) (*Market, error) {
	header, err := r.client.Header(ctx, number)
	if err != nil {
		return nil, err
	}
	opts := &bind.CallOpts{Context: ctx, BlockNumber: header.Number}
	m := &Market{Block: header.Number.Uint64(), TimeMs: header.Time * 1000}
	bnd, engine := r.client.Band(), r.client.Margin()
	if m.Settled, err = bnd.SequencerSettled(opts); err != nil {
		return nil, err
	}
	session, err := bnd.Session(opts)
	if err != nil {
		return nil, err
	}
	m.Closed, m.BoundaryMs = closedAt(session, m.TimeMs), session.BoundaryMs
	m.Regime = margin.RegimeAt(session.State, session.BoundaryMs, m.TimeMs)
	symbols, err := engine.Assets(opts)
	if err != nil {
		return nil, err
	}
	parameters := margin.Parameters{Symbols: symbols, Market: -1}
	if market, err := engine.Market(opts); err != nil {
		return nil, err
	} else if market != ([32]byte{}) {
		for i, s := range symbols {
			if s == market {
				parameters.Market = i
			}
		}
	}
	for i, symbol := range symbols {
		volatility, _, err := engine.Volatility(opts, symbol)
		if err != nil {
			return nil, err
		}
		gap, _, err := engine.WeekendGap(opts, symbol)
		if err != nil {
			return nil, err
		}
		depth, err := engine.Depth(opts, symbol)
		if err != nil {
			return nil, err
		}
		parameters.Volatilities = append(parameters.Volatilities, volatility)
		parameters.Gaps = append(parameters.Gaps, gap)
		m.Depths = append(m.Depths, depth.Selling, depth.Buying)
		for _, other := range symbols[i+1:] {
			correlation, _, err := engine.Correlation(opts, symbol, other)
			if err != nil {
				return nil, err
			}
			parameters.Correlations = append(parameters.Correlations, correlation)
		}
	}
	m.Set = margin.NewSet(parameters, scenarioSize, margin.Horizon)
	engineEth, err := engine.EthUsdFeed(opts)
	if err != nil {
		return nil, err
	}
	engineUSD, err := r.answer(opts, engineEth, header.Time)
	if err != nil {
		return nil, err
	}
	for _, symbol := range symbols {
		asset, token, err := r.asset(opts, name(symbol))
		if err != nil {
			return nil, err
		}
		m.Assets = append(m.Assets, asset)
		pool, err := r.pool(opts, symbol, token, engineUSD)
		if err != nil {
			return nil, err
		}
		m.Pools = append(m.Pools, pool)
	}
	liquidatorEth, err := r.client.Liquidator().EthUsd(opts)
	if err != nil {
		return nil, err
	}
	if m.EthPrice, err = r.answer(opts, liquidatorEth, header.Time); err != nil {
		return nil, err
	}
	m.RecallHaircut, err = r.client.Liquidator().RecallHaircut(opts)
	return m, err
}

func (r *Reader) asset(opts *bind.CallOpts, asset string) (Asset, common.Address, error) {
	bnd := r.client.Band()
	quote, err := bnd.Quote(opts, asset)
	if err != nil {
		return Asset{}, common.Address{}, err
	}
	action, err := bnd.CorporateAction(opts, asset)
	if err != nil {
		return Asset{}, common.Address{}, err
	}
	feeds, err := bnd.Asset(opts, asset)
	if err != nil {
		return Asset{}, common.Address{}, err
	}
	feed, basis := feeds.RedstoneFeedId, uint64(0)
	if feed == ([32]byte{}) {
		feed, basis = feeds.IndexFeedId, band.IndexBasisBps
	}
	variance, err := bnd.Variance(opts, feed)
	if err != nil {
		return Asset{}, common.Address{}, err
	}
	return Asset{Name: asset, State: quote.State, Unconfirmed: action.Status == 2, Mid: quote.Mid, Low: quote.Low,
		High: quote.High, Variance: variance, Basis: basis}, feeds.Token, nil
}

// answer is a Chainlink feed's latest answer, zero where there is no feed or its round is not positive, newer than
// the block or more than a day and a minute old.
func (r *Reader) answer(opts *bind.CallOpts, feed common.Address, nowS uint64) (*big.Int, error) {
	if feed == (common.Address{}) {
		return new(big.Int), nil
	}
	round, err := r.client.LatestRound(opts, feed)
	if err != nil {
		return nil, err
	}
	updated := round.UpdatedAt.Uint64()
	if round.Answer.Sign() <= 0 || !round.UpdatedAt.IsUint64() || updated > nowS || nowS-updated > feedMaxAgeS {
		return new(big.Int), nil
	}
	return round.Answer, nil
}

// pool is asset's pool as the engine's liquidity add-on reads it from its mean tick and liquidity over the last
// 1,800 seconds: absent without a pool, unread where the pool cannot say or its quote has no price.
func (r *Reader) pool(opts *bind.CallOpts, symbol [32]byte, stock common.Address, engineUSD *big.Int) (margin.Pool, error) {
	address, err := r.client.Margin().Pool(opts, symbol)
	if err != nil || address == (common.Address{}) {
		return margin.Pool{}, err
	}
	pool := r.client.Pool(address)
	token0, token1, err := pool.Tokens(opts)
	if err != nil {
		return margin.Pool{}, err
	}
	fee, err := pool.Fee(opts)
	if err != nil {
		return margin.Pool{}, err
	}
	unread := margin.Pool{Kind: margin.Unread, Terms: margin.PoolTerms{Fee: fee}}
	tokens := r.client.Deployments().Tokens
	stockIsToken0, quoteIsWeth, ok := margin.Classify(token0, token1, stock, tokens["USDG"], tokens["WETH"])
	if !ok {
		return unread, nil
	}
	observed, err := pool.Observe(opts, []uint32{margin.Window, 0})
	if err != nil || len(observed.TickCumulatives) != 2 || len(observed.SecondsPerLiquidityCumulativeX128s) != 2 {
		return unread, nil
	}
	tick, liquidity, ok := margin.Mean(observed.TickCumulatives[0].Int64(), observed.TickCumulatives[1].Int64(),
		observed.SecondsPerLiquidityCumulativeX128s[0], observed.SecondsPerLiquidityCumulativeX128s[1])
	usd, quoteDecimals := big.NewInt(100_000_000), uint8(usdgDecimals)
	if quoteIsWeth {
		usd, quoteDecimals = engineUSD, wethDecimals
	}
	if !ok || usd.Sign() == 0 {
		return unread, nil
	}
	terms, ok := margin.Terms(stockIsToken0, stockDecimals, quoteDecimals, tick, liquidity, usd)
	if !ok {
		return unread, nil
	}
	terms.Fee = fee
	return margin.Pool{Kind: margin.Read, Terms: terms}, nil
}

// Position reads account's position at the block of m: what it holds of each Stock Token (collateral, lent and in
// its baskets), USDG and WETH, and what it owes.
func (r *Reader) Position(ctx context.Context, m *Market, account common.Address, id [32]byte) (*Position, error) {
	opts := &bind.CallOpts{Context: ctx, BlockNumber: new(big.Int).SetUint64(m.Block)}
	accounts := r.client.Accounts()
	stocks, err := accounts.Stocks(opts)
	if err != nil {
		return nil, err
	}
	basketed, err := accounts.InBaskets(opts, account, id)
	if err != nil {
		return nil, err
	}
	p := &Position{Account: account, ID: id}
	for i, asset := range stocks.Assets {
		token := stocks.Tokens[i]
		if token == (common.Address{}) {
			continue
		}
		collateral, err := accounts.Collateral(opts, account, id, token)
		if err != nil {
			return nil, err
		}
		lent, err := accounts.Lent(opts, account, id, token)
		if err != nil {
			return nil, err
		}
		quantity := new(big.Int).Add(collateral, lent)
		if held := basketed[asset]; held != nil {
			quantity.Add(quantity, held)
		}
		if quantity.Sign() > 0 {
			p.Stocks = append(p.Stocks, Stock{Asset: asset, Quantity: quantity, Lent: lent})
		}
	}
	tokens := r.client.Deployments().Tokens
	if p.USDG, err = accounts.Collateral(opts, account, id, tokens["USDG"]); err != nil {
		return nil, err
	}
	if p.WETH, err = accounts.Collateral(opts, account, id, tokens["WETH"]); err != nil {
		return nil, err
	}
	owed, err := accounts.Repayment(opts, account, id)
	if err != nil {
		return nil, err
	}
	p.Debt, p.Premium = owed.Debt, owed.Premium
	return p, nil
}

// Judgement is a position valued as the liquidator values it now. Evaluated is false where the Go requirement is not
// the engine's at the block's own prices, or where an asset the position holds is not one the engine lists: nothing
// is then said of the position.
type Judgement struct {
	Position  *Position
	Valuation Valuation
	Evaluated bool
}

// Judge reads account's position at the block of m and values it at the band's current edges, in the liquidator's
// market state. Its requirement is checked against the engine's currentRequirement at the same block.
func (r *Reader) Judge(ctx context.Context, m *Market, account common.Address, id [32]byte) (*Judgement, error) {
	p, err := r.Position(ctx, m, account, id)
	if err != nil {
		return nil, err
	}
	edges := m.Current()
	j := &Judgement{Position: p, Valuation: m.Assess(p, edges, m.Closed, m.Regime)}
	quantities, prices, listed := m.Portfolio(p, edges)
	if !listed {
		return j, nil
	}
	j.Evaluated = true
	if !j.Valuation.Judged || !j.Valuation.Held {
		return j, nil
	}
	opts := &bind.CallOpts{Context: ctx, BlockNumber: new(big.Int).SetUint64(m.Block)}
	engine, err := r.client.Margin().CurrentRequirement(opts, quantities, prices)
	if _, reverted := sdk.DecodeRevert(err); reverted {
		j.Evaluated = false
		return j, nil
	}
	if err != nil {
		return nil, fmt.Errorf("the engine's current requirement: %w", err)
	}
	want, ok := m.Requirement(quantities, prices, m.Regime)
	j.Evaluated = ok && engine.Regime == m.Regime.Code() && engine.Margin.Cmp(want) == 0
	return j, nil
}

// Auction is when the auction of account's position that runs in the market's current state started, in seconds, zero
// where none runs, at the block of m.
func (r *Reader) Auction(ctx context.Context, m *Market, account common.Address, id [32]byte) (uint64, error) {
	opts := &bind.CallOpts{Context: ctx, BlockNumber: new(big.Int).SetUint64(m.Block)}
	liquidator := r.client.Liquidator()
	running, err := liquidator.Running(opts, account, id)
	if err != nil || !running {
		return 0, err
	}
	auction, err := liquidator.Auction(opts, account, id)
	return auction.StartedAt, err
}

// LiquidationPrice is the price of asset at which account's position stops falling short in the open market, in USD
// with 8 decimals, as the margin accounts read it at the block of m: zero where none does.
func (r *Reader) LiquidationPrice(ctx context.Context, m *Market, p *Position, asset string) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx, BlockNumber: new(big.Int).SetUint64(m.Block)}
	return r.client.Accounts().LiquidationPrice(opts, p.Account, p.ID, asset, new(big.Int))
}

// Short reads account's short of asset at the block of m: its equity and requirement at the band's high edge, and
// whether it is open at all.
func (r *Reader) Short(ctx context.Context, m *Market, account common.Address, asset string) (equity, requirement *big.Int, open bool, err error) {
	opts := &bind.CallOpts{Context: ctx, BlockNumber: new(big.Int).SetUint64(m.Block)}
	shorts := r.client.Shorts()
	position, err := shorts.Position(opts, account, asset)
	if err != nil || position.Shares.Sign() == 0 {
		return nil, nil, false, err
	}
	health, err := shorts.Health(opts, account, asset)
	return health.Equity, health.Requirement, true, err
}
