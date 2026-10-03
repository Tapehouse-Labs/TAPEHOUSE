// SPDX-License-Identifier: MIT OR Apache-2.0

// Command backtest runs the band and the margin engine over every weekend and holiday closure of the launch assets
// since 2010 and over the closures since Robinhood Chain launched, from public data alone, and writes what it
// measures as JSON. Its record command rebuilds a deployed band from its PriceWritten and Anchored events and checks
// the rebuild against the band's own variance, session and quotes.
//
// Usage:
//
//	backtest history [-out FILE] [-deployments DIR] [-parameters FILE] [-debt-cap USD] [-launch DATE]
//	backtest record RPC_URL DEPLOYMENTS_JSON
//
// history reads Robinhood Chain at ROBINHOOD_RPC_URL and Arbitrum One at ARBITRUM_RPC_URL.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapehouse/tapehouse/services/internal/backtest"
	"github.com/tapehouse/tapehouse/services/internal/band"
	"github.com/tapehouse/tapehouse/services/internal/history"
	"github.com/tapehouse/tapehouse/services/internal/margin"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/quoterv2"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/stocktoken"
	"golang.org/x/time/rate"
)

var (
	names  = []string{"NVDA", "TSLA", "AAPL", "MSFT", "GOOGL", "SPY"}
	status = []string{band.CurrentStatus, band.NextStatus, band.NextChangeTime}
	one    = big.NewInt(1_000_000_000_000_000_000)
)

const (
	since       = "2010-01-01"
	indexFeed   = "USA500.Y---24_7"
	blockSample = 1_000_000
	// requestsPerSecond keeps the reads within a public RPC's rate limit.
	requestsPerSecond = 20
	// regularSession is how long NYSE's regular session lasts, from its open to its close.
	regularSession = 6*time.Hour + 30*time.Minute
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: backtest history [flags] | backtest record RPC_URL DEPLOYMENTS_JSON")
	}
	var err error
	switch os.Args[1] {
	case "history":
		err = runHistory(os.Args[2:])
	case "record":
		if len(os.Args) != 4 {
			log.Fatal("usage: backtest record RPC_URL DEPLOYMENTS_JSON")
		}
		err = runRecord(context.Background(), os.Args[2], os.Args[3])
	default:
		log.Fatalf("unknown command %q", os.Args[1])
	}
	if err != nil {
		log.Fatal(err)
	}
}

func asset(name string, regularHours bool) band.Asset {
	if name == "SPY" {
		return band.Asset{Symbol: name, Index: indexFeed, RegularHours: regularHours}
	}
	return band.Asset{Symbol: name, Feed: name + "---24_7", RegularHours: regularHours}
}

type chain struct {
	client      *ethclient.Client
	deployments *sdk.Deployments
	limit       *rate.Limiter
}

func dial(ctx context.Context, url, registry string) (chain, error) {
	if url == "" {
		return chain{}, fmt.Errorf("no RPC URL for %s", registry)
	}
	client, err := ethclient.DialContext(ctx, url)
	if err != nil {
		return chain{}, err
	}
	d, err := sdk.LoadDeployments(registry)
	if err != nil {
		return chain{}, err
	}
	id, err := client.ChainID(ctx)
	if err != nil {
		return chain{}, err
	}
	if id.Uint64() != d.ChainID {
		return chain{}, fmt.Errorf("%s is for chain %d, the RPC serves %v", registry, d.ChainID, id)
	}
	return chain{client, d, rate.NewLimiter(rate.Limit(requestsPerSecond), 1)}, nil
}

func runHistory(args []string) error {
	fs := flag.NewFlagSet("history", flag.ExitOnError)
	out := fs.String("out", "backtest.json", "where to write the report")
	dir := fs.String("deployments", "../deployments", "the directory of the chains' registries")
	parameters := fs.String("parameters", "../stylus/contracts/margin/parameters.json", "the engine's parameters")
	debtCap := fs.Float64("debt-cap", 100_000, "the margin accounts' debt cap, in USD")
	launchDate := fs.String("launch", "2026-07-01", "when Robinhood Chain launched")
	if err := fs.Parse(args); err != nil {
		return err
	}
	ctx := context.Background()
	web := &http.Client{Timeout: 2 * time.Minute}
	from, _ := time.Parse(time.DateOnly, since)
	launch, err := time.Parse(time.DateOnly, *launchDate)
	if err != nil {
		return err
	}
	today := time.Now().UTC().Truncate(24 * time.Hour)
	series := make([]history.Series, len(names))
	if err := each(len(names), func(i int) (err error) {
		series[i], err = history.Daily(ctx, web, history.YahooBase, names[i], from, today)
		return err
	}); err != nil {
		return err
	}
	data, err := os.ReadFile(*parameters)
	if err != nil {
		return err
	}
	p, depths, err := margin.Launch(data, names)
	if err != nil {
		return err
	}
	robinhood, err := dial(ctx, os.Getenv("ROBINHOOD_RPC_URL"), filepath.Join(*dir, "4663.json"))
	if err != nil {
		return err
	}
	arbitrum, err := dial(ctx, os.Getenv("ARBITRUM_RPC_URL"), filepath.Join(*dir, "42161.json"))
	if err != nil {
		return err
	}
	interval, err := blockInterval(ctx, robinhood.client)
	if err != nil {
		return err
	}
	m := backtest.Align(series)
	report := backtest.Run(m, margin.NewSet(p, 256, margin.Horizon), depths, *debtCap, interval, launch)
	var after []backtest.Closure
	for _, c := range m.Closures() {
		if !c.Open.Before(launch) {
			after = append(after, c)
		}
	}
	report.Reopens, err = reopens(ctx, web, []chain{robinhood, arbitrum}, after)
	if err != nil {
		return err
	}
	report.Accuracy = backtest.Summarize(report.Reopens)
	pools, err := margin.Pools(data)
	if err != nil {
		return err
	}
	if report.Impacts, err = impacts(ctx, robinhood, pools, depths); err != nil {
		return err
	}
	encoded, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(*out, append(encoded, '\n'), 0o644); err != nil {
		return err
	}
	summarize(report, *out)
	return nil
}

func each(n int, run func(int) error) error {
	errs := make([]error, n)
	var wg sync.WaitGroup
	for i := range n {
		wg.Go(func() { errs[i] = run(i) })
	}
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func blockInterval(ctx context.Context, client *ethclient.Client) (float64, error) {
	latest, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	earlier, err := client.HeaderByNumber(ctx, new(big.Int).Sub(latest.Number, big.NewInt(blockSample)))
	if err != nil {
		return 0, err
	}
	return float64(latest.Time-earlier.Time) / blockSample, nil
}

func redstone(ctx context.Context, web *http.Client, feed string) ([]history.Point, error) {
	var ranges [][]history.Point
	for _, days := range []int{1, 7, 30} {
		points, err := history.RedStone(ctx, web, history.RedStoneBase, feed, days)
		if err != nil {
			return nil, err
		}
		ranges = append(ranges, points)
	}
	return history.Merge(ranges...), nil
}

// rounds returns a reader of feed's rounds at opts, paced by the chain's limit.
func (c chain) rounds(opts *bind.CallOpts, feed *history.Feed) func(*big.Int) (history.Round, bool, error) {
	return func(id *big.Int) (history.Round, bool, error) {
		if err := c.limit.Wait(opts.Context); err != nil {
			return history.Round{}, false, err
		}
		return feed.Get(opts, id)
	}
}

func reopens(ctx context.Context, web *http.Client, chains []chain, closures []backtest.Closure) ([]backtest.Reopen, error) {
	feeds := append(slicesOf(names, func(n string) string { return asset(n, false).Feed }), indexFeed)
	feeds = append(feeds, status...)
	legs := map[string][]history.Point{}
	var mu sync.Mutex
	if err := each(len(feeds), func(i int) error {
		if feeds[i] == "" {
			return nil
		}
		points, err := redstone(ctx, web, feeds[i])
		mu.Lock()
		legs[feeds[i]] = points
		mu.Unlock()
		return err
	}); err != nil {
		return nil, err
	}
	signed := [3][]history.Point{legs[status[0]], legs[status[1]], legs[status[2]]}
	results := make([][]backtest.Reopen, len(chains)*len(names))
	err := each(len(results), func(i int) error {
		c, name := chains[i/len(names)], names[i%len(names)]
		a := asset(name, c.deployments.SharePrices())
		var address common.Address
		if c.deployments.SharePrices() {
			feed, err := c.deployments.SharePriceFeed(name + "_USD")
			if err != nil {
				return err
			}
			address = feed.Address
		} else {
			feed, err := c.deployments.TokenPriceFeed(name + "_USD")
			if err != nil {
				return err
			}
			address = feed.Address
		}
		opts := &bind.CallOpts{Context: ctx}
		feed := history.NewFeed(c.client, address)
		latest, err := feed.Latest(opts)
		if err != nil {
			return err
		}
		multiplier, err := multiplierOf(opts, c, name)
		if err != nil {
			return err
		}
		leg := legs[a.Feed]
		if a.Feed == "" {
			leg = legs[a.Index]
		}
		for _, cl := range closures {
			at := uint64(cl.Close.Truncate(24 * time.Hour).Add(36 * time.Hour).Unix())
			if a.RegularHours {
				at = uint64(cl.Open.Unix()) - 1
			}
			before, after, err := history.Around(c.rounds(opts, feed), latest, at)
			if err != nil {
				return err
			}
			if after.ID == nil {
				continue
			}
			up, err := sequencerUp(opts, c, after.UpdatedAt-1)
			if err != nil {
				return err
			}
			results[i] = append(results[i], backtest.Reopening{
				Chain: c.deployments.ChainID, Asset: a, Date: cl.Open.Format(time.DateOnly),
				SinceMs: uint64(cl.Close.Add(regularSession).UnixMilli()), Leg: leg, Status: signed,
				Before: before, After: after, Multiplier: multiplier, SequencerUp: up,
			}.Replay())
		}
		return nil
	})
	var out []backtest.Reopen
	for _, r := range results {
		out = append(out, r...)
	}
	return out, err
}

// impacts quotes, through Uniswap's QuoterV2 at the latest block, a sale of each asset worth a share of its selling
// depth in the pool the engine liquidates it in, valued at Chainlink's latest round, and its quote token's proceeds in
// USD at Chainlink's ETH/USD where it is WETH.
func impacts(ctx context.Context, c chain, pools map[string]string, depths []uint32) ([]backtest.Impact, error) {
	head, err := c.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	opts := &bind.CallOpts{Context: ctx, BlockNumber: head.Number}
	quoter := quoterv2.NewQuoterV2()
	contract := bind.NewBoundContract(c.deployments.UniswapV3["QuoterV2"], abi.ABI{}, c.client, c.client, c.client)
	usd := func(feed string) (*big.Float, error) {
		f, err := c.deployments.TokenPriceFeed(feed)
		if err != nil {
			return nil, err
		}
		r, err := history.NewFeed(c.client, f.Address).Latest(opts)
		if err != nil {
			return nil, err
		}
		return new(big.Float).Quo(new(big.Float).SetInt(r.Answer), big.NewFloat(1e8)), nil
	}
	var out []backtest.Impact
	for k, name := range names {
		parts := strings.Split(pools[name], "_")
		if len(parts) != 3 {
			return nil, fmt.Errorf("pool %q of %s", pools[name], name)
		}
		fee, err := strconv.ParseUint(parts[2], 10, 32)
		if err != nil {
			return nil, err
		}
		price, err := usd(name + "_USD")
		if err != nil {
			return nil, err
		}
		quote := big.NewFloat(1)
		if parts[1] == "WETH" {
			if quote, err = usd("ETH_USD"); err != nil {
				return nil, err
			}
		}
		depth := float64(depths[2*k])
		for _, share := range backtest.ImpactShares {
			value := share * depth
			amount, _ := new(big.Float).Mul(new(big.Float).Quo(big.NewFloat(value), price), big.NewFloat(1e18)).Int(nil)
			if err := c.limit.Wait(ctx); err != nil {
				return nil, err
			}
			q, err := bind.Call(contract, opts, quoter.PackQuoteExactInputSingle(quoterv2.IQuoterV2QuoteExactInputSingleParams{
				TokenIn: c.deployments.Tokens[name], TokenOut: c.deployments.Tokens[parts[1]], AmountIn: amount,
				Fee: new(big.Int).SetUint64(fee), SqrtPriceLimitX96: new(big.Int),
			}), quoter.UnpackQuoteExactInputSingle)
			if err != nil {
				return nil, fmt.Errorf("quoting %s: %w", name, err)
			}
			decimals := 1e18
			if parts[1] == "USDG" {
				decimals = 1e6
			}
			proceeds, _ := new(big.Float).Mul(new(big.Float).Quo(new(big.Float).SetInt(q.AmountOut), big.NewFloat(decimals)), quote).Float64()
			out = append(out, backtest.Impact{
				Asset: name, Block: head.Number.Uint64(), Share: share,
				Engine: backtest.EngineCost(value, depth, uint32(fee)), Pool: 1 - proceeds/value,
			})
		}
	}
	return out, nil
}

func slicesOf(xs []string, f func(string) string) []string {
	out := make([]string, len(xs))
	for i, x := range xs {
		out[i] = f(x)
	}
	return out
}

func multiplierOf(opts *bind.CallOpts, c chain, name string) (*big.Int, error) {
	token, ok := c.deployments.Tokens[name]
	if !ok {
		return one, nil
	}
	binding := stocktoken.NewStockToken()
	contract := bind.NewBoundContract(token, abi.ABI{}, c.client, c.client, c.client)
	m, err := bind.Call(contract, opts, binding.PackUiMultiplier(), binding.UnpackUiMultiplier)
	if err != nil {
		return one, nil
	}
	return m, nil
}

func sequencerUp(opts *bind.CallOpts, c chain, nowS uint64) (bool, error) {
	address, ok := c.deployments.ChainlinkSequencer["Uptime"]
	if !ok {
		return true, nil
	}
	feed := history.NewFeed(c.client, address)
	latest, err := feed.Latest(opts)
	if err != nil {
		return false, err
	}
	before, _, err := history.Around(c.rounds(opts, feed), latest, nowS)
	if err != nil || before.ID == nil {
		return false, err
	}
	return band.SequencerSettled(before.Answer.Uint64(), before.StartedAt, nowS), nil
}

func summarize(r backtest.Report, out string) {
	h := r.History
	fmt.Printf("history: %d sessions from %s to %s, %d closures, %d since %s\n", h.Sessions, h.From, h.To, h.Closures, h.Since, h.Launch)
	for _, kind := range []string{"long", "mixed"} {
		c := r.Coverage[kind]
		fmt.Printf("buffer (%s): covers the closed requirement for %d of %d, the worst closure for %d, largest shortfall %.1f%%; ramp covers %.1f%% / %.1f%% / %.1f%% of closures\n",
			kind, c.BufferCoversClosed, c.Portfolios, c.CoveredAtWorst, 100*c.LargestShortfall, 100*c.Ramp[0], 100*c.Ramp[1], 100*c.Ramp[2])
	}
	for _, x := range r.Horizons {
		fmt.Printf("horizon %s: ES99² %.2f days, variance %.2f days\n", x.Asset, x.Days, x.Variance)
	}
	l := r.Leverage
	fmt.Printf("leverage: cap %d× (deployed %d×); breaches long %d at the model, %d at the cap, %d at the deployed cap; capacity %.1f%% to %.1f%% at %d×, %.1f%% at %d×\n",
		l.Cap, l.DeployedCap, l.Model["long"].Count, l.AtCap["long"].Count, l.AtDeployedCap["long"].Count, 100*l.CapacityModel,
		100*l.CapacityCapped, l.DeployedCap, 100*l.CapacityAtCap, l.Cap)
	for _, w := range l.Worst {
		fmt.Printf("worst closure %s: %+.1f%% on %s\n", w.Asset, 100*w.Move, w.Date)
	}
	fmt.Printf("deleveraging: median cut %.0f%% long, %.0f%% mixed\n", 100*r.Deleveraging.Long, 100*r.Deleveraging.Mixed)
	if s := r.March2020; s != nil {
		fmt.Printf("16 March 2020: binds %d long and %d mixed; capacity %.1f%% with the scenario, %.1f%% with the cap; mixed pay %.2f%% more with it, %.2f%% with the cap\n",
			s.Binds["long"], s.Binds["mixed"], 100*s.CapacityScenario, 100*s.CapacityCap, 100*s.HedgeTaxScenario, 100*s.HedgeTaxCap)
	}
	c := r.Capacity
	fmt.Printf("capacity: %.1f%% against %.1f%% flat; bad debt %.4f%% against %.4f%%\n", 100*c.Model, 100*c.Flat, 100*c.BadDebtModel, 100*c.BadDebtFlat)
	b := r.Backstop
	limits := make([]string, len(b.Limits))
	for i, x := range b.Limits {
		limits[i] = fmt.Sprintf("%s %.0f", x.Position, x.Limit)
	}
	fmt.Printf("backstop at %.0f USD of debt: limits %s; premium %d bps a year\n", b.DebtCap, strings.Join(limits, ", "), b.PremiumRateBps)
	fmt.Printf("liquidator: blocks %.3f s apart\n", r.Liquidator.BlockSeconds)
	for i := 0; i < len(r.Impacts); i += len(backtest.ImpactShares) {
		x := r.Impacts[i : i+len(backtest.ImpactShares)]
		fmt.Printf("impact %s at block %d, a tenth, half and all of its depth: engine %.2f%% / %.2f%% / %.2f%%, pool %.2f%% / %.2f%% / %.2f%%\n",
			x[0].Asset, x[0].Block, 100*x[0].Engine, 100*x[1].Engine, 100*x[2].Engine, 100*x[0].Pool, 100*x[1].Pool, 100*x[2].Pool)
	}
	for _, a := range r.Accuracy {
		fmt.Printf("band on %d: %d closures, %d with a public 24/7 value, %d reopening rounds inside the band's last value, median error %.1f bps, %.0f s ahead\n",
			a.Chain, a.Reopens, a.Evaluated, a.Inside, a.MedianError, a.MedianLeadS)
	}
	fmt.Printf("wrote %s\n", out)
}
