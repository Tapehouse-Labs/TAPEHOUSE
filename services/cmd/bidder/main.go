// SPDX-License-Identifier: MIT OR Apache-2.0

// Command bidder is a reference bidder for Tapehouse's reopening auction and its liquidation auctions, built on the Go
// SDK, go-ethereum and the indexer's public API alone. Its strategy:
//
//   - A Stock Token is worth to it the lower of the low edge of the band sealed before the reopen and of the band now,
//     less TAPEHOUSE_BIDDER_DISCOUNT_BPS: never above what the band's low edge says the token is worth.
//   - It commits one bid for each round, within TAPEHOUSE_BIDDER_COMMIT_LEAD of the commit phase's end, sized by
//     TAPEHOUSE_BIDDER_BUDGET and the round's supply, with a deposit that does not reveal the bid's size.
//   - It commits a bid only if its reveal cannot be refused: a price at or above the round's floor, a bid worth at
//     least the minimum at the floor, an escrow within the deposit. An unrevealed deposit forfeits at least the bond.
//   - It writes each bid to TAPEHOUSE_BIDDER_STATE before it commits, so a restart still reveals it, and claims once the
//     round is cleared or has lapsed.
//   - In a liquidation auction it buys a Stock Token once the ask has fallen to its price, with the ask's cost as its
//     limit. It approves exactly the deposit or the limit each time, never more.
//
// It reads its configuration from the environment: PRIVATE_KEY, TAPEHOUSE_RPC_URL and TAPEHOUSE_DEPLOYMENTS, and
// optionally TAPEHOUSE_INDEXER_URL (without it the liquidation auctions are not watched), TAPEHOUSE_INDEXER_KEY,
// TAPEHOUSE_BIDDER_STATE (bidder-<chainId>.json), TAPEHOUSE_BIDDER_BUDGET (USDG in raw units for each round and each
// purchase, 1000000000), TAPEHOUSE_BIDDER_DISCOUNT_BPS (300), TAPEHOUSE_BIDDER_COMMIT_LEAD (15m),
// TAPEHOUSE_BIDDER_ASSETS (every asset of the accounts) and TAPEHOUSE_BIDDER_DRY_RUN.
package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapehouse/tapehouse/services/sdk"
)

const (
	pollEvery  = 2 * time.Second
	dutchEvery = 15 // seconds of chain time between Dutch ticks
)

// outcomes is what each revert a bidder meets comes to; any other revert is an error.
var outcomes = map[string]outcome{
	"WrongPhase":        retry,
	"NoAuction":         retry,
	"NotLiquidatable":   retry,
	"PositionHeld":      retry,
	"CostAboveLimit":    retry,
	"NothingToBuy":      retry,
	"WrongState":        retry,
	"UnknownCommitment": done,
}

type config struct {
	rpcURL, deployments, indexerURL, indexerKey, statePath string
	budget                                                 *big.Int
	discountBps                                            uint64
	commitLead                                             time.Duration
	assets                                                 []string
	dryRun                                                 bool
}

func loadConfig(getenv func(string) string) (config, error) {
	cfg := config{rpcURL: getenv("TAPEHOUSE_RPC_URL"), deployments: getenv("TAPEHOUSE_DEPLOYMENTS"),
		indexerURL: strings.TrimRight(getenv("TAPEHOUSE_INDEXER_URL"), "/"), indexerKey: getenv("TAPEHOUSE_INDEXER_KEY"),
		statePath: getenv("TAPEHOUSE_BIDDER_STATE"), budget: big.NewInt(1000e6), discountBps: 300,
		commitLead: 15 * time.Minute}
	for _, required := range []string{"TAPEHOUSE_RPC_URL", "TAPEHOUSE_DEPLOYMENTS"} {
		if getenv(required) == "" {
			return cfg, fmt.Errorf("set %s", required)
		}
	}
	if text := getenv("TAPEHOUSE_BIDDER_BUDGET"); text != "" {
		budget, ok := new(big.Int).SetString(text, 10)
		if !ok || budget.Sign() <= 0 {
			return cfg, errors.New("TAPEHOUSE_BIDDER_BUDGET is not a positive integer")
		}
		cfg.budget = budget
	}
	if text := getenv("TAPEHOUSE_BIDDER_DISCOUNT_BPS"); text != "" {
		bps, err := strconv.ParseUint(text, 10, 64)
		if err != nil || bps >= 10_000 {
			return cfg, errors.New("TAPEHOUSE_BIDDER_DISCOUNT_BPS is not a number below 10000")
		}
		cfg.discountBps = bps
	}
	if text := getenv("TAPEHOUSE_BIDDER_COMMIT_LEAD"); text != "" {
		lead, err := time.ParseDuration(text)
		if err != nil || lead < 0 {
			return cfg, errors.New("TAPEHOUSE_BIDDER_COMMIT_LEAD is not a duration")
		}
		cfg.commitLead = lead
	}
	for _, asset := range strings.Split(getenv("TAPEHOUSE_BIDDER_ASSETS"), ",") {
		if asset = strings.TrimSpace(asset); asset != "" {
			cfg.assets = append(cfg.assets, asset)
		}
	}
	dryRun, _ := strconv.ParseBool(getenv("TAPEHOUSE_BIDDER_DRY_RUN"))
	cfg.dryRun = dryRun
	return cfg, nil
}

// newLogger returns a JSON logger on w that replaces each secret, wherever it appears, with "redacted".
func newLogger(w io.Writer, secrets ...string) *slog.Logger {
	scrub := func(text string) string {
		for _, secret := range secrets {
			if secret != "" {
				text = strings.ReplaceAll(text, secret, "redacted")
			}
		}
		return text
	}
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		switch value := a.Value.Any().(type) {
		case string:
			a.Value = slog.StringValue(scrub(value))
		case error:
			a.Value = slog.StringValue(scrub(value.Error()))
		}
		return a
	}}))
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := run(ctx, os.Getenv, os.Stderr)
	stop()
	os.Exit(code)
}

func run(ctx context.Context, getenv func(string) string, stderr io.Writer) int {
	log := newLogger(stderr)
	hex := getenv("PRIVATE_KEY")
	if hex == "" {
		log.Error("set PRIVATE_KEY")
		return 2
	}
	key, err := crypto.HexToECDSA(strings.TrimPrefix(hex, "0x"))
	if err != nil {
		log.Error("PRIVATE_KEY is not a private key")
		return 2
	}
	cfg, err := loadConfig(getenv)
	if err != nil {
		log.Error(err.Error())
		return 2
	}
	log = newLogger(stderr, cfg.rpcURL, cfg.indexerKey)
	if err := serve(ctx, cfg, key, log); err != nil && !errors.Is(err, context.Canceled) {
		log.Error("the bidder stopped", "error", err)
		return 1
	}
	return 0
}

func serve(ctx context.Context, cfg config, key *ecdsa.PrivateKey, log *slog.Logger) error {
	deployments, err := sdk.LoadDeployments(cfg.deployments)
	if err != nil {
		return err
	}
	chain, err := ethclient.DialContext(ctx, cfg.rpcURL)
	if err != nil {
		return errors.New("TAPEHOUSE_RPC_URL does not parse")
	}
	defer chain.Close()
	chainID, err := chain.ChainID(ctx)
	if err != nil {
		return err
	}
	if chainID.Uint64() != deployments.ChainID {
		return fmt.Errorf("the registry is of chain %d and the RPC of chain %s", deployments.ChainID, chainID)
	}
	if cfg.statePath == "" {
		cfg.statePath = fmt.Sprintf("bidder-%s.json", chainID)
	}
	state, err := LoadState(cfg.statePath)
	if err != nil {
		return err
	}
	bidder := &Bidder{Client: sdk.NewClient(chain, deployments), Chain: chain, Opts: bind.NewKeyedTransactor(key, chainID),
		State: state, Budget: cfg.budget, DiscountBps: cfg.discountBps, CommitLead: cfg.commitLead, Assets: cfg.assets,
		DryRun: cfg.dryRun, Log: log, IndexerURL: cfg.indexerURL, IndexerKey: cfg.indexerKey}
	log.Info("running the bidder", "chainId", chainID.Uint64(), "address", bidder.Opts.From.Hex(), "dryRun", cfg.dryRun,
		"dutch", cfg.indexerURL != "", "bids", len(state.Bids()))
	if cfg.indexerURL == "" {
		log.Warn("the liquidation auctions are not watched: set TAPEHOUSE_INDEXER_URL")
	}
	var last, dutchAt uint64
	for {
		head, err := chain.HeaderByNumber(ctx, nil)
		switch {
		case ctx.Err() != nil:
			return ctx.Err()
		case err != nil:
			log.Warn("the head could not be read", "error", err)
		case head.Number.Uint64() != last:
			last = head.Number.Uint64()
			if err := bidder.Reopening(ctx, head); err != nil {
				log.Error("the reopening auction", "error", err)
			}
			if cfg.indexerURL != "" && head.Time >= dutchAt+dutchEvery {
				dutchAt = head.Time
				if err := bidder.Dutch(ctx, head); err != nil {
					log.Error("the liquidation auctions", "error", err)
				}
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(pollEvery):
		}
	}
}
