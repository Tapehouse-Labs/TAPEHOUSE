// SPDX-License-Identifier: MIT OR Apache-2.0

// Package keeper runs Tapehouse's keepers on one chain: every permissionless action that keeps the band fresh, the
// session's records current and the positions, vaults, auctions, baskets and covers settled on time. Each pass reads
// the chain, and the index for what the chain cannot enumerate, decides from that alone, simulates each action and
// sends it only if it would succeed, so a keeper that crashes and restarts, or one that runs beside another, repeats
// nothing.
package keeper

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"math/big"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/internal/redstone"
	"github.com/tapehouse/tapehouse/services/sdk"
)

// HaltsURL is Robinhood's quotes of its Stock Tokens, each with whether its trading is halted.
const HaltsURL = "https://api.robinhood.com/rhj/prices"

// Config is what the keepers run on, read from the environment by ConfigFromEnv. Keys are read by Keys.
type Config struct {
	// Deployments is the path of the chain's registry, deployments/<chainId>.json: TAPEHOUSE_DEPLOYMENTS.
	Deployments string
	// RPCURLs are the chain's RPC endpoints, tried in order: TAPEHOUSE_RPC_URL, then TAPEHOUSE_RPC_FALLBACK_URLS,
	// separated by commas. They may carry keys and are never logged.
	RPCURLs []string
	// RPCRate is the most requests a second the keepers send the RPC, 20 unless TAPEHOUSE_RPC_RATE names another.
	RPCRate float64
	// IndexerURL is the indexer's API, which enumerates positions, shorts, recalls, commitments and covers:
	// TAPEHOUSE_INDEXER_URL. Only the keepers that need it require it.
	IndexerURL string
	// Keepers are the keepers to run, separated by commas: TAPEHOUSE_KEEPERS, unless set every one whose contracts the
	// registry names.
	Keepers []string
	// Gateways are RedStone's gateways, the keyed ones with REDSTONE_API_KEY and REDSTONE_BACKUP_API_KEY where set.
	Gateways []redstone.Gateway
	// HaltsURL is where the halt keeper reads trading halts, HaltsURL unless TAPEHOUSE_HALTS_URL names another.
	HaltsURL string
	// Buy has the liquidation keeper buy at an auction as the buyer of last resort, with the keeper's own USDG, once the
	// auction asks no more than the asset's reference price: TAPEHOUSE_KEEPER_BUY=true.
	Buy bool
	// Rebalance has the basket keeper rebalance from its own Stock Tokens even when the pools do not pay for the
	// bands' width, the team paying it: TAPEHOUSE_KEEPER_REBALANCE=subsidise.
	Rebalance bool
}

// ConfigFromEnv reads the configuration from the environment.
func ConfigFromEnv(getenv func(string) string) (Config, error) {
	cfg := Config{
		Deployments: getenv("TAPEHOUSE_DEPLOYMENTS"),
		IndexerURL:  strings.TrimRight(getenv("TAPEHOUSE_INDEXER_URL"), "/"),
		Gateways:    redstone.Gateways(getenv),
		HaltsURL:    getenv("TAPEHOUSE_HALTS_URL"),
		RPCRate:     20,
		Buy:         getenv("TAPEHOUSE_KEEPER_BUY") == "true",
		Rebalance:   getenv("TAPEHOUSE_KEEPER_REBALANCE") == "subsidise",
	}
	if primary := getenv("TAPEHOUSE_RPC_URL"); primary != "" {
		cfg.RPCURLs = append(cfg.RPCURLs, primary)
	}
	cfg.RPCURLs = append(cfg.RPCURLs, list(getenv("TAPEHOUSE_RPC_FALLBACK_URLS"))...)
	if cfg.Deployments == "" || len(cfg.RPCURLs) == 0 {
		return cfg, errors.New("TAPEHOUSE_DEPLOYMENTS and TAPEHOUSE_RPC_URL are required")
	}
	if text := getenv("TAPEHOUSE_RPC_RATE"); text != "" {
		rate, err := strconv.ParseFloat(text, 64)
		if err != nil || rate <= 0 {
			return cfg, errors.New("TAPEHOUSE_RPC_RATE must be a positive number of requests a second")
		}
		cfg.RPCRate = rate
	}
	if cfg.HaltsURL == "" {
		cfg.HaltsURL = HaltsURL
	}
	cfg.Keepers = list(getenv("TAPEHOUSE_KEEPERS"))
	for _, name := range cfg.Keepers {
		if _, ok := jobByName(name); !ok {
			return cfg, fmt.Errorf("TAPEHOUSE_KEEPERS: no keeper %q; the keepers are %s", name, strings.Join(Names(), ", "))
		}
	}
	return cfg, nil
}

func list(text string) []string {
	var out []string
	for _, item := range strings.Split(text, ",") {
		if item = strings.TrimSpace(item); item != "" {
			out = append(out, item)
		}
	}
	return out
}

// Keeper holds what the keepers share: the SDK client, the sender of their transactions, the index, RedStone's
// packages, the halt signer, and what each remembers between passes, which only spares a repeated call.
type Keeper struct {
	client  *sdk.Client
	chain   Chain
	sender  *Sender
	index   *Index
	prices  sdk.PackageSource
	halts   *halts
	cfg     Config
	log     *slog.Logger
	now     func() time.Time
	sleep   func(context.Context, time.Duration) error
	mu      sync.Mutex
	assets  []asset
	memory  map[string]string
	started map[string]bool
}

// New returns the keepers of cfg over chain, sending from key and signing halts with haltKey, which may be nil.
func New(cfg Config, deployments *sdk.Deployments, chain Chain, key, haltKey *ecdsa.PrivateKey, log *slog.Logger) (*Keeper, error) {
	if len(cfg.Keepers) == 0 {
		for _, j := range jobs {
			if _, ok := deployments.Tapehouse[j.contract]; ok || (j.contract == "Baskets" && len(deployments.Baskets) > 0) {
				cfg.Keepers = append(cfg.Keepers, j.name)
			}
		}
	}
	for _, name := range cfg.Keepers {
		if j, _ := jobByName(name); j.index && cfg.IndexerURL == "" {
			return nil, fmt.Errorf("the %s keeper needs TAPEHOUSE_INDEXER_URL", name)
		}
	}
	client := sdk.NewClient(chain, deployments)
	k := &Keeper{
		client:  client,
		chain:   chain,
		sender:  NewSender(client, chain, key, new(big.Int).SetUint64(deployments.ChainID), log),
		prices:  statusGroup{redstone.New(cfg.Gateways, &http.Client{Timeout: 10 * time.Second})},
		cfg:     cfg,
		log:     log,
		now:     time.Now,
		sleep:   sleep,
		memory:  map[string]string{},
		started: map[string]bool{},
	}
	if cfg.IndexerURL != "" {
		k.index = &Index{base: cfg.IndexerURL, client: &http.Client{Timeout: 30 * time.Second}}
	}
	if haltKey != nil {
		k.halts = &halts{key: haltKey, url: cfg.HaltsURL, client: &http.Client{Timeout: 10 * time.Second}}
	}
	return k, nil
}

// Keepers returns the keepers k runs.
func (k *Keeper) Keepers() []string {
	return k.cfg.Keepers
}

// From returns the address the keepers send from.
func (k *Keeper) From() common.Address {
	return k.sender.From()
}

// job is one keeper: its name, how often it runs a pass, the registry's contract it keeps, and whether it reads the
// index.
type job struct {
	name     string
	every    time.Duration
	contract string
	index    bool
	pass     func(*Keeper, context.Context) error
}

var jobs = []job{
	{"prices", 60 * time.Second, "Band", false, (*Keeper).writePrices},
	{"multiplier", 30 * time.Second, "Band", false, (*Keeper).syncMultipliers},
	{"halt", 30 * time.Second, "Band", false, (*Keeper).writeHalts},
	{"seal", 15 * time.Second, "Band", false, (*Keeper).seal},
	{"premium", 60 * time.Second, "MarginAccounts", false, (*Keeper).accruePremium},
	{"mark", 30 * time.Second, "ShortPositions", false, (*Keeper).mark},
	{"sync", 60 * time.Second, "MarginAccounts", true, (*Keeper).syncHoldings},
	{"liquidation", 15 * time.Second, "Liquidator", true, (*Keeper).liquidate},
	{"backstop", 60 * time.Second, "GapBackstop", true, (*Keeper).backstop},
	{"reopening", 30 * time.Second, "ReopeningAuction", true, (*Keeper).reopening},
	{"recall", 30 * time.Second, "MarginAccounts", true, (*Keeper).recalls},
	{"shorts", 15 * time.Second, "ShortPositions", true, (*Keeper).shorts},
	{"basket", 5 * time.Minute, "Baskets", false, (*Keeper).rebalance},
	{"gapcover", 20 * time.Second, "GapCover", true, (*Keeper).gapCover},
}

// Names are the keepers, in the order they first run.
func Names() []string {
	names := make([]string, len(jobs))
	for i, j := range jobs {
		names[i] = j.name
	}
	return names
}

func jobByName(name string) (job, bool) {
	i := slices.IndexFunc(jobs, func(j job) bool { return j.name == name })
	if i < 0 {
		return job{}, false
	}
	return jobs[i], true
}

// Pass runs one pass of the keeper name.
func (k *Keeper) Pass(ctx context.Context, name string) error {
	j, ok := jobByName(name)
	if !ok {
		return fmt.Errorf("no keeper %q", name)
	}
	if j.index && k.index == nil {
		return fmt.Errorf("the %s keeper needs the index", name)
	}
	return j.pass(k, ctx)
}

// Run runs the configured keepers until ctx ends: one pass of each in order, then each on its own schedule. A failed
// pass is logged and the next one tries again from the chain.
func (k *Keeper) Run(ctx context.Context) error {
	var group sync.WaitGroup
	for _, name := range k.cfg.Keepers {
		j, _ := jobByName(name)
		k.runPass(ctx, j)
		group.Go(func() {
			ticker := time.NewTicker(j.every)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					k.runPass(ctx, j)
				}
			}
		})
	}
	group.Wait()
	return nil
}

func (k *Keeper) runPass(ctx context.Context, j job) {
	if err := j.pass(k, ctx); err != nil && ctx.Err() == nil {
		k.log.Warn("the pass failed; the next one tries again", "keeper", j.name, "error", err)
	}
}

// remember reports whether key last held value, and records it: what spares a keeper a repeated call.
func (k *Keeper) remember(key, value string) bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	same := k.memory[key] == value
	k.memory[key] = value
	return same
}

func sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// nowMs is the keeper's clock in milliseconds, as the contracts compare it with the session's boundaries.
func (k *Keeper) nowMs() uint64 {
	return uint64(k.now().UnixMilli())
}

// asset is an asset the band prices: its name, its feeds as the band configures them, and its Stock Token.
type asset struct {
	name    string
	feedID  string
	indexID string
	token   common.Address
}

// bandAssets reads, once, the assets the band prices among those the registry names: its Stock Tokens, its Chainlink
// feeds and its BandFeeds. The band's assets are set at its deployment.
func (k *Keeper) bandAssets(ctx context.Context) ([]asset, error) {
	k.mu.Lock()
	cached := k.assets
	k.mu.Unlock()
	if cached != nil {
		return cached, nil
	}
	d := k.client.Deployments()
	names := map[string]bool{}
	for name := range d.Tokens {
		names[name] = true
	}
	for name := range d.Chainlink {
		if asset, ok := strings.CutSuffix(name, "_USD"); ok {
			names[asset] = true
		}
	}
	for name := range d.BandFeeds {
		names[name] = true
	}
	var out []asset
	for _, name := range slices.Sorted(maps.Keys(names)) {
		read, err := k.client.Band().Asset(k.opts(ctx), name)
		if err != nil {
			return nil, err
		}
		feedID, indexID := text(read.RedstoneFeedId), text(read.IndexFeedId)
		if read.ChainlinkFeed == (common.Address{}) && feedID == "" && indexID == "" {
			continue
		}
		out = append(out, asset{name, feedID, indexID, read.Token})
	}
	k.mu.Lock()
	k.assets = out
	k.mu.Unlock()
	return out, nil
}

func text(word [32]byte) string {
	return strings.TrimRight(string(word[:]), "\x00")
}
