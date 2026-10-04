// SPDX-License-Identifier: MIT OR Apache-2.0

package calibrate

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tapehouse/tapehouse/services/internal/backtest"
	"github.com/tapehouse/tapehouse/services/internal/history"
	"github.com/tapehouse/tapehouse/services/internal/margin"
	"github.com/tapehouse/tapehouse/services/sdk"
)

// What the backtest that gates a proposal is run with.
const (
	// UpdateInterval is the seconds the engine requires between two updates.
	UpdateInterval = 86_400
	// DebtCap is the margin accounts' debt cap in USD.
	DebtCap = 100_000
	// BlockSeconds is the chain's block interval.
	BlockSeconds = 0.25
	// LaunchDate is when Robinhood Chain launched.
	LaunchDate = "2026-07-01"
)

// The ways a proposal is refused.
var (
	ErrTooSoon   = errors.New("the engine was updated less than a day ago")
	ErrStaleBase = errors.New("the engine no longer holds the values the proposal stepped from")
	ErrWorse     = errors.New("the backtest finds the proposal worse than the current parameters")
)

// Proposal is an update of the engine's parameters: what it stepped from at Block, what the method estimated, what it
// proposes within the engine's rules, the call that makes it and how the call simulated from the owner, and the
// backtest of the current and the proposed parameters. Reference and Tail are what the filter and the cover's refit
// report beside it.
type Proposal struct {
	Chain                 uint64          `json:"chain"`
	Margin                common.Address  `json:"margin"`
	Block                 uint64          `json:"block"`
	Assets                []string        `json:"assets"`
	LastUpdate            uint64          `json:"lastUpdate"`
	Owner                 common.Address  `json:"owner"`
	Base                  Current         `json:"base"`
	Targets               Targets         `json:"targets"`
	Proposed              Proposed        `json:"proposed"`
	Calldata              hexutil.Bytes   `json:"calldata"`
	Simulation            string          `json:"simulation"`
	Current               backtest.Report `json:"current"`
	Next                  backtest.Report `json:"next"`
	Reference             []Distance      `json:"reference,omitempty"`
	ReferenceVolatilities []uint32        `json:"referenceVolatilities"`
	Tail                  *Tail           `json:"tail,omitempty"`
}

// state is what the engine holds at one block.
type state struct {
	assets     []string
	symbols    [][32]byte
	current    Current
	market     int
	lastUpdate uint64
	owner      common.Address
}

func name(symbol [32]byte) string {
	return strings.TrimRight(string(symbol[:]), "\x00")
}

func read(ctx context.Context, client *sdk.Client, block *big.Int) (state, error) {
	opts := &bind.CallOpts{Context: ctx, BlockNumber: block}
	m := client.Margin()
	var s state
	var err error
	if s.symbols, err = m.Assets(opts); err != nil {
		return s, err
	}
	for _, symbol := range s.symbols {
		s.assets = append(s.assets, name(symbol))
	}
	for _, symbol := range s.symbols {
		v, vf, err := m.Volatility(opts, symbol)
		if err != nil {
			return s, err
		}
		g, gf, err := m.WeekendGap(opts, symbol)
		if err != nil {
			return s, err
		}
		d, err := m.Depth(opts, symbol)
		if err != nil {
			return s, err
		}
		c := &s.current
		c.Volatilities, c.VolatilityFloors = append(c.Volatilities, v), append(c.VolatilityFloors, vf)
		c.Gaps, c.GapFloors = append(c.Gaps, g), append(c.GapFloors, gf)
		c.Depths, c.DepthCeilings = append(c.Depths, d.Selling, d.Buying), append(c.DepthCeilings, d.SellingCeiling, d.BuyingCeiling)
	}
	for i, a := range s.symbols {
		for _, b := range s.symbols[i+1:] {
			v, f, err := m.Correlation(opts, a, b)
			if err != nil {
				return s, err
			}
			s.current.Correlations, s.current.CorrelationFloors = append(s.current.Correlations, v), append(s.current.CorrelationFloors, f)
		}
	}
	market, err := m.Market(opts)
	if err != nil {
		return s, err
	}
	s.market = -1
	if market != ([32]byte{}) {
		s.market = slices.Index(s.assets, name(market))
	}
	if s.lastUpdate, err = m.LastUpdate(opts); err != nil {
		return s, err
	}
	s.owner, err = m.Owner(opts)
	return s, err
}

// Restrict is the market's columns for names, in that order.
func Restrict(m backtest.Market, names []string) (backtest.Market, error) {
	index := make([]int, len(names))
	for i, n := range names {
		if index[i] = slices.Index(m.Names, n); index[i] < 0 {
			return backtest.Market{}, fmt.Errorf("no history for %s", n)
		}
	}
	out := backtest.Market{Names: names, At: m.At}
	for _, rows := range [][2]any{{m.Open, &out.Open}, {m.Close, &out.Close}, {m.Adj, &out.Adj}} {
		from, into := rows[0].([][]float64), rows[1].(*[][]float64)
		for _, row := range from {
			picked := make([]float64, len(names))
			for i, k := range index {
				picked[i] = row[k]
			}
			*into = append(*into, picked)
		}
	}
	return out, nil
}

// Series are the market's assets' daily bars.
func Series(m backtest.Market) []history.Series {
	out := make([]history.Series, len(m.Names))
	for k, n := range m.Names {
		out[k].Symbol = n
		for i, at := range m.At {
			out[k].Bars = append(out[k].Bars, history.Bar{At: at, Open: m.Open[i][k], Close: m.Close[i][k], AdjClose: m.Adj[i][k]})
		}
	}
	return out
}

func symbolsOf(s state) margin.Parameters {
	return margin.Parameters{Symbols: s.symbols, Market: s.market}
}

func report(m backtest.Market, s state, volatilities []uint32, correlations []uint16, gaps, depths []uint32) backtest.Report {
	p := symbolsOf(s)
	p.Volatilities, p.Correlations, p.Gaps = volatilities, correlations, gaps
	launch, _ := time.Parse(time.DateOnly, LaunchDate)
	return backtest.Run(m, margin.NewSet(p, 256, margin.Horizon), depths, DebtCap, BlockSeconds, launch)
}

// Gate refuses a proposal that, at the deployed weekend cap, leaves more portfolios losing more than their margin than
// the current parameters did, or raises the bad-debt rate.
func Gate(current, next backtest.Report) error {
	for _, kind := range []string{"long", "mixed"} {
		was, is := current.Leverage.AtDeployedCap[kind].Portfolios, next.Leverage.AtDeployedCap[kind].Portfolios
		if is > was {
			return fmt.Errorf("%w: %d %s portfolios lose more than their margin at the deployed cap, against %d", ErrWorse, is, kind, was)
		}
	}
	if next.Capacity.BadDebtModel > current.Capacity.BadDebtModel {
		return fmt.Errorf("%w: a bad-debt rate of %.4f%%, against %.4f%%", ErrWorse, 100*next.Capacity.BadDebtModel, 100*current.Capacity.BadDebtModel)
	}
	return nil
}

// Propose reads the engine at the latest block and proposes its next parameters: the method's estimates over the
// market, with the depths of the snapshots where a weekend's is among them, stepped within the engine's rules from its
// current values. It refuses within a day of the last update, backtests the current and the proposed parameters over
// the market and refuses a worse proposal, with the proposal returned for its reports, and simulates the call from the
// engine's owner. The market's last session is the method's last day.
func Propose(ctx context.Context, client *sdk.Client, market backtest.Market, snapshots []Snapshot, now time.Time) (Proposal, error) {
	header, err := client.Header(ctx, nil)
	if err != nil {
		return Proposal{}, err
	}
	address, err := client.Margin().Address()
	if err != nil {
		return Proposal{}, err
	}
	s, err := read(ctx, client, header.Number)
	if err != nil {
		return Proposal{}, err
	}
	p := Proposal{
		Chain: client.Deployments().ChainID, Margin: address, Block: header.Number.Uint64(), Assets: s.assets,
		LastUpdate: s.lastUpdate, Owner: s.owner, Base: s.current,
	}
	if next := s.lastUpdate + UpdateInterval; s.lastUpdate != 0 && uint64(now.Unix()) < next {
		return p, fmt.Errorf("%w: the next may come at %s", ErrTooSoon, time.Unix(int64(next), 0).UTC().Format(time.RFC3339))
	}
	m, err := Restrict(market, s.assets)
	if err != nil {
		return p, err
	}
	last := m.At[len(m.At)-1]
	if p.Targets, err = Estimate(Series(market), s.assets, s.market, last); err != nil {
		return p, err
	}
	depths, ok := Depths(snapshots, uint64(now.Unix()), s.assets)
	if !ok {
		depths = nil
	}
	if p.Proposed, err = Step(s.current, p.Targets, depths); err != nil {
		return p, err
	}
	for k := range s.assets {
		closes := make([]float64, len(m.At))
		for i := range closes {
			closes[i] = m.Adj[i][k]
		}
		p.ReferenceVolatilities = append(p.ReferenceVolatilities, ReferenceVolatility(closes))
	}
	tx, err := client.Margin().SetParameters(p.Proposed.Volatilities, p.Proposed.Correlations, p.Proposed.Gaps, p.Proposed.Depths)
	if err != nil {
		return p, err
	}
	p.Calldata = tx.Data
	p.Current = report(m, s, s.current.Volatilities, s.current.Correlations, s.current.Gaps, s.current.Depths)
	p.Next = report(m, s, p.Proposed.Volatilities, p.Proposed.Correlations, p.Proposed.Gaps, p.Proposed.Depths)
	if err := Gate(p.Current, p.Next); err != nil {
		return p, err
	}
	if err := client.Simulate(&bind.CallOpts{Context: ctx, From: s.owner, BlockNumber: header.Number}, tx); err != nil {
		revert, ok := sdk.DecodeRevert(err)
		if !ok || revert.Name != "UpdateTooSoon" {
			return p, fmt.Errorf("the engine refuses the proposal: %w", err)
		}
		p.Simulation = revert.Error()
	} else {
		p.Simulation = "ok"
	}
	return p, nil
}

// CheckBase refuses a proposal whose base the engine no longer holds at the latest block: another update landed.
func CheckBase(ctx context.Context, client *sdk.Client, p Proposal) error {
	header, err := client.Header(ctx, nil)
	if err != nil {
		return err
	}
	s, err := read(ctx, client, header.Number)
	if err != nil {
		return err
	}
	if s.lastUpdate != p.LastUpdate {
		return fmt.Errorf("%w: its last update is %d, not %d", ErrStaleBase, s.lastUpdate, p.LastUpdate)
	}
	if !slices.Equal(s.assets, p.Assets) {
		return fmt.Errorf("%w: its assets are %v", ErrStaleBase, s.assets)
	}
	for i, a := range s.assets {
		for _, c := range []struct {
			what    string
			was, is uint64
		}{
			{"volatility", uint64(p.Base.Volatilities[i]), uint64(s.current.Volatilities[i])},
			{"gap", uint64(p.Base.Gaps[i]), uint64(s.current.Gaps[i])},
			{"selling depth", uint64(p.Base.Depths[2*i]), uint64(s.current.Depths[2*i])},
			{"buying depth", uint64(p.Base.Depths[2*i+1]), uint64(s.current.Depths[2*i+1])},
		} {
			if c.was != c.is {
				return fmt.Errorf("%w: %s %s is %d, not %d", ErrStaleBase, a, c.what, c.is, c.was)
			}
		}
	}
	for k, pair := range Pairs(s.assets) {
		if p.Base.Correlations[k] != s.current.Correlations[k] {
			return fmt.Errorf("%w: %s/%s correlation is %d, not %d", ErrStaleBase, pair[0], pair[1], s.current.Correlations[k], p.Base.Correlations[k])
		}
	}
	return nil
}
