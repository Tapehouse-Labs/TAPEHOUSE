// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"context"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapehouse/tapehouse/services/internal/backtest"
	"github.com/tapehouse/tapehouse/services/internal/calibrate"
	"github.com/tapehouse/tapehouse/services/internal/history"
	"github.com/tapehouse/tapehouse/services/sdk"
	"golang.org/x/time/rate"
)

const (
	launchDate = "2026-07-01"
	// requestsPerSecond keeps the reads within a public RPC's rate limit.
	requestsPerSecond = 20
	// weekendSpan is the shortest closure the gap rule counts, three days less an hour.
	weekendSpan = 3*24*time.Hour - time.Hour
)

var one = big.NewInt(1_000_000_000_000_000_000)

// multiplier is the multiplier of name's Stock Token as the band names it, or one where the registry names none.
func multiplier(ctx context.Context, client *sdk.Client, name string) *big.Int {
	terms, err := client.Band().Terms(&bind.CallOpts{Context: ctx}, name)
	if err != nil || terms.Multiplier == nil {
		return one
	}
	return terms.Multiplier
}

// ratioOf is how much of a weekend's fall from the close to the Monday open the 24/5 reopening price had made: its
// move from the close over the open's. A rise has no ratio.
func ratioOf(closePrice, open, reopen float64) (float64, bool) {
	move := open/closePrice - 1
	if move >= 0 || closePrice <= 0 {
		return 0, false
	}
	return (reopen/closePrice - 1) / move, true
}

// tail refits the gap cover's tail where the registry names one, and reports it beside the deployed constants.
func tail(ctx context.Context, chain *ethclient.Client, client *sdk.Client, m backtest.Market, p calibrate.Proposal) (*calibrate.Tail, error) {
	if _, ok := client.Deployments().Tapehouse["GapCover"]; !ok {
		return nil, nil
	}
	restricted, err := calibrate.Restrict(m, p.Assets)
	if err != nil {
		return nil, err
	}
	opts := &bind.CallOpts{Context: ctx}
	threshold, scale, probability, err := client.GapCover().Tail(opts)
	if err != nil {
		return nil, err
	}
	ratios, err := reopenRatios(ctx, chain, client, restricted)
	if err != nil {
		return nil, err
	}
	t := calibrate.FitTail(calibrate.Falls(restricted, p.Base.Gaps), ratios, seed).Against([3]uint32{threshold, scale, probability})
	return &t, nil
}

// reopenRatios are, for each weekend since the chain launched, the ratio of each asset's fall that the first 24/5
// round after the Friday close had made, where the Stock Token's Chainlink feed has such a round and the Monday open
// was a fall. A chain whose feeds price the share has none.
func reopenRatios(ctx context.Context, chain *ethclient.Client, client *sdk.Client, m backtest.Market) ([][]float64, error) {
	d := client.Deployments()
	if d.SharePrices() {
		return nil, nil
	}
	launch, _ := time.Parse(time.DateOnly, launchDate)
	limit := rate.NewLimiter(rate.Limit(requestsPerSecond), 1)
	opts := &bind.CallOpts{Context: ctx}
	type leg struct {
		feed   *history.Feed
		latest history.Round
		mult   float64
	}
	legs := make([]leg, len(m.Names))
	for k, name := range m.Names {
		feed, err := d.TokenPriceFeed(name + "_USD")
		if err != nil {
			return nil, err
		}
		f := history.NewFeed(chain, feed.Address)
		latest, err := f.Latest(opts)
		if err != nil {
			return nil, err
		}
		mult, _ := new(big.Float).Quo(new(big.Float).SetInt(multiplier(ctx, client, name)), new(big.Float).SetInt(one)).Float64()
		legs[k] = leg{f, latest, mult}
	}
	var out [][]float64
	for i := 1; i < len(m.At); i++ {
		if m.At[i].Before(launch) || m.At[i].Sub(m.At[i-1]) < weekendSpan {
			continue
		}
		at := uint64(m.At[i-1].Truncate(24 * time.Hour).Add(36 * time.Hour).Unix())
		var week []float64
		for k, l := range legs {
			_, after, err := history.Around(func(id *big.Int) (history.Round, bool, error) {
				if err := limit.Wait(ctx); err != nil {
					return history.Round{}, false, err
				}
				return l.feed.Get(opts, id)
			}, l.latest, at)
			if err != nil {
				return nil, err
			}
			if after.ID == nil || !after.Answer.IsInt64() {
				continue
			}
			price := float64(after.Answer.Int64()) / 1e8 / l.mult
			if r, ok := ratioOf(m.Close[i-1][k], m.Open[i][k], price); ok && !math.IsNaN(r) {
				week = append(week, r)
			}
		}
		if len(week) > 0 {
			out = append(out, week)
		}
	}
	return out, nil
}
