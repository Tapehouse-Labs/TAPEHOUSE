// SPDX-License-Identifier: MIT OR Apache-2.0

// Command calibrate proposes new parameters for the margin engine and applies them where its key owns the engine.
//
//	calibrate snapshot
//	calibrate propose [-out FILE] [-as-of DATE] [-deployments FILE]
//	calibrate apply FILE
//
// snapshot records each of the engine's pools at the latest block, to TAPEHOUSE_CALIBRATOR_DIR (calibrate/). propose
// estimates every parameter setParameters takes from the latest ten years of daily bars, steps them within what the
// engine accepts of its current values, takes each depth as the lowest of the last four weeks' snapshots where one was
// taken on a weekend, backtests the current and the proposed parameters, and simulates the call from the engine's
// owner; it writes the proposal to FILE (proposal.json) and prints it. -as-of fixes the method's last session.
// apply sends a proposal, only where the engine has not been updated since it was made and PRIVATE_KEY is the owner.
//
// It reads TAPEHOUSE_RPC_URL and TAPEHOUSE_DEPLOYMENTS from the environment.
package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapehouse/tapehouse/services/internal/backtest"
	"github.com/tapehouse/tapehouse/services/internal/band"
	"github.com/tapehouse/tapehouse/services/internal/calibrate"
	"github.com/tapehouse/tapehouse/services/internal/history"
	"github.com/tapehouse/tapehouse/services/internal/keeper"
	"github.com/tapehouse/tapehouse/services/internal/redact"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/margin"
)

const (
	since     = "2010-01-01"
	indexFeed = "USA500.Y---24_7"
	weeks     = 30
	seed      = 1
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err := run(ctx, os.Args[1:], os.Getenv, os.Stdout)
	stop()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type config struct {
	rpcURL, deployments, dir string
}

func loadConfig(getenv func(string) string) (config, error) {
	cfg := config{rpcURL: getenv("TAPEHOUSE_RPC_URL"), deployments: getenv("TAPEHOUSE_DEPLOYMENTS"), dir: getenv("TAPEHOUSE_CALIBRATOR_DIR")}
	if cfg.dir == "" {
		cfg.dir = "calibrate"
	}
	return cfg, nil
}

func run(ctx context.Context, args []string, getenv func(string) string, out io.Writer) error {
	if len(args) == 0 {
		return errors.New("usage: calibrate snapshot | propose [-out FILE] [-as-of DATE] [-deployments FILE] | apply FILE")
	}
	cfg, err := loadConfig(getenv)
	if err != nil {
		return err
	}
	switch args[0] {
	case "snapshot":
		return snapshot(ctx, cfg, out)
	case "propose":
		return propose(ctx, cfg, args[1:], out)
	case "apply":
		if len(args) != 2 {
			return errors.New("usage: calibrate apply FILE")
		}
		return apply(ctx, cfg, getenv, args[1], out)
	}
	return fmt.Errorf("unknown command %q", args[0])
}

// connect dials the RPC and reads the registry, refusing a registry of another chain. Errors never carry the URL.
func connect(ctx context.Context, cfg config) (*ethclient.Client, *sdk.Client, error) {
	if cfg.rpcURL == "" || cfg.deployments == "" {
		return nil, nil, errors.New("set TAPEHOUSE_RPC_URL and TAPEHOUSE_DEPLOYMENTS")
	}
	d, err := sdk.LoadDeployments(cfg.deployments)
	if err != nil {
		return nil, nil, err
	}
	chain, err := ethclient.DialContext(ctx, cfg.rpcURL)
	if err != nil {
		return nil, nil, errors.New("TAPEHOUSE_RPC_URL does not parse")
	}
	id, err := chain.ChainID(ctx)
	if err != nil {
		chain.Close()
		return nil, nil, errors.New(redact.String(err.Error(), cfg.rpcURL))
	}
	if id.Uint64() != d.ChainID {
		chain.Close()
		return nil, nil, fmt.Errorf("the registry is of chain %d and the RPC of chain %s", d.ChainID, id)
	}
	return chain, sdk.NewClient(chain, d), nil
}

func engineAssets(ctx context.Context, client *sdk.Client) ([]string, error) {
	symbols, err := client.Margin().Assets(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, err
	}
	names := make([]string, len(symbols))
	for i, s := range symbols {
		names[i] = strings.TrimRight(string(s[:]), "\x00")
	}
	return names, nil
}

func snapshot(ctx context.Context, cfg config, out io.Writer) error {
	chain, client, err := connect(ctx, cfg)
	if err != nil {
		return err
	}
	defer chain.Close()
	assets, err := engineAssets(ctx, client)
	if err != nil {
		return err
	}
	s, err := calibrate.TakeSnapshot(ctx, client, assets)
	if err != nil {
		return errors.New(redact.String(err.Error(), cfg.rpcURL))
	}
	path := filepath.Join(cfg.dir, fmt.Sprintf("depth-%d.json", s.Block))
	if err := writeJSON(path, s); err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, "snapshot of %d pools at block %d, %s: %s\n", len(s.Pools), s.Block, time.Unix(int64(s.Time), 0).UTC().Format(time.RFC3339), path)
	return err
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

// loadSnapshots reads the snapshots of dir; a directory that does not exist holds none.
func loadSnapshots(dir string) ([]calibrate.Snapshot, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "depth-*.json"))
	if err != nil {
		return nil, err
	}
	var out []calibrate.Snapshot
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var s calibrate.Snapshot
		if err := json.Unmarshal(data, &s); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		out = append(out, s)
	}
	return out, nil
}

// until keeps the sessions up to last.
func until(m backtest.Market, last time.Time) backtest.Market {
	n := len(m.At)
	for n > 0 && m.At[n-1].After(last) {
		n--
	}
	m.At, m.Open, m.Close, m.Adj = m.At[:n], m.Open[:n], m.Close[:n], m.Adj[:n]
	return m
}

// market is the daily bars of names, and of the stress window's market where it is not one of them, since 2010 and up
// to the last session of asOf.
func market(ctx context.Context, web *http.Client, names []string, asOf time.Time) (backtest.Market, error) {
	from, _ := time.Parse(time.DateOnly, since)
	if !slices.Contains(names, calibrate.StressMarket) {
		names = append(slices.Clone(names), calibrate.StressMarket)
	}
	series := make([]history.Series, len(names))
	for i, name := range names {
		s, err := history.Daily(ctx, web, history.YahooBase, name, from, asOf.Add(48*time.Hour))
		if err != nil {
			return backtest.Market{}, err
		}
		series[i] = s
	}
	return until(backtest.Align(series), asOf), nil
}

func propose(ctx context.Context, cfg config, args []string, out io.Writer) error {
	fs := flag.NewFlagSet("propose", flag.ContinueOnError)
	file := fs.String("out", "proposal.json", "where to write the proposal")
	asOf := fs.String("as-of", "", "the method's last session, YYYY-MM-DD (the latest)")
	registry := fs.String("deployments", "", "the registry, over TAPEHOUSE_DEPLOYMENTS")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *registry != "" {
		cfg.deployments = *registry
	}
	chain, client, err := connect(ctx, cfg)
	if err != nil {
		return err
	}
	defer chain.Close()
	assets, err := engineAssets(ctx, client)
	if err != nil {
		return err
	}
	last := time.Now().UTC()
	if *asOf != "" {
		if last, err = time.Parse(time.DateOnly, *asOf); err != nil {
			return fmt.Errorf("-as-of: %w", err)
		}
		last = last.Add(23 * time.Hour)
	}
	web := &http.Client{Timeout: 2 * time.Minute}
	m, err := market(ctx, web, assets, last)
	if err != nil {
		return err
	}
	snapshots, err := loadSnapshots(cfg.dir)
	if err != nil {
		return err
	}
	p, err := calibrate.Propose(ctx, client, m, snapshots, time.Now())
	if err != nil && !errors.Is(err, calibrate.ErrWorse) {
		return errors.New(redact.String(err.Error(), cfg.rpcURL))
	}
	worse := err
	p.Reference = reference(ctx, web, client, assets)
	if p.Tail, err = tail(ctx, chain, client, m, p); err != nil {
		return errors.New(redact.String(err.Error(), cfg.rpcURL))
	}
	if err := writeJSON(*file, p); err != nil {
		return err
	}
	for _, line := range p.Lines() {
		if _, err := fmt.Fprintln(out, line); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(out, "wrote %s\n", *file); err != nil {
		return err
	}
	return worse
}

func feedOf(name string) band.Asset {
	if name == "SPY" {
		return band.Asset{Symbol: name, Index: indexFeed}
	}
	return band.Asset{Symbol: name, Feed: name + "---24_7"}
}

// reference measures the band against the filter over the last weeks' RedStone values of each asset, leaving out any
// asset whose history could not be read.
func reference(ctx context.Context, web *http.Client, client *sdk.Client, assets []string) []calibrate.Distance {
	var out []calibrate.Distance
	for _, name := range assets {
		a := feedOf(name)
		if a.Feed == "" {
			continue
		}
		var legs [][]history.Point
		for _, days := range []int{1, 7, weeks} {
			points, err := history.RedStone(ctx, web, history.RedStoneBase, a.Feed, days)
			if err != nil {
				legs = nil
				break
			}
			legs = append(legs, points)
		}
		if legs == nil {
			continue
		}
		out = append(out, calibrate.Measure(a, history.Merge(legs...), multiplier(ctx, client, name)))
	}
	return out
}

func apply(ctx context.Context, cfg config, getenv func(string) string, file string, out io.Writer) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	var p calibrate.Proposal
	if err := json.Unmarshal(data, &p); err != nil {
		return fmt.Errorf("%s: %w", file, err)
	}
	hex := getenv("PRIVATE_KEY")
	if hex == "" {
		return errors.New("set PRIVATE_KEY")
	}
	key, err := crypto.HexToECDSA(strings.TrimPrefix(hex, "0x"))
	if err != nil {
		return errors.New("PRIVATE_KEY is not a private key")
	}
	chain, client, err := connect(ctx, cfg)
	if err != nil {
		return err
	}
	defer chain.Close()
	if client.Deployments().ChainID != p.Chain {
		return fmt.Errorf("the proposal is for chain %d, the registry of chain %d", p.Chain, client.Deployments().ChainID)
	}
	if err := calibrate.CheckBase(ctx, client, p); err != nil {
		return err
	}
	owner, err := client.Margin().Owner(&bind.CallOpts{Context: ctx})
	if err != nil {
		return err
	}
	if from := crypto.PubkeyToAddress(key.PublicKey); from != owner {
		return fmt.Errorf("%s is not the engine's owner, %s: hand the proposal's calldata to the owner", from, owner)
	}
	return send(ctx, chain, client, key, p, out)
}

func send(ctx context.Context, chain *ethclient.Client, client *sdk.Client, key *ecdsa.PrivateKey, p calibrate.Proposal, out io.Writer) error {
	id, err := chain.ChainID(ctx)
	if err != nil {
		return err
	}
	tx, err := client.Margin().SetParameters(p.Proposed.Volatilities, p.Proposed.Correlations, p.Proposed.Gaps, p.Proposed.Depths)
	if err != nil {
		return err
	}
	if !slices.Equal(tx.Data, p.Calldata) {
		return errors.New("the proposal's calldata is not its values' setParameters")
	}
	receipt, err := keeper.NewSender(client, chain, key, id, slog.New(slog.DiscardHandler)).Act(ctx, "setParameters", tx)
	if err != nil {
		return err
	}
	lines := []string{fmt.Sprintf("setParameters in block %d, transaction %s, %d gas", receipt.BlockNumber.Uint64(), receipt.TxHash, receipt.GasUsed)}
	engine := margin.NewMargin()
	for _, l := range receipt.Logs {
		if e, err := engine.UnpackVolatilitySetEvent(l); err == nil {
			lines = append(lines, fmt.Sprintf("VolatilitySet %s %d", name(e.Symbol), e.Value))
		} else if e, err := engine.UnpackCorrelationSetEvent(l); err == nil {
			lines = append(lines, fmt.Sprintf("CorrelationSet %s/%s %d", name(e.Symbol), name(e.Other), e.Value))
		} else if e, err := engine.UnpackGapSetEvent(l); err == nil {
			lines = append(lines, fmt.Sprintf("GapSet %s %d", name(e.Symbol), e.Value))
		} else if e, err := engine.UnpackDepthSetEvent(l); err == nil {
			lines = append(lines, fmt.Sprintf("DepthSet %s %d %d", name(e.Symbol), e.Selling, e.Buying))
		}
	}
	_, err = io.WriteString(out, strings.Join(lines, "\n")+"\n")
	return err
}

func name(symbol [32]byte) string {
	return strings.TrimRight(string(symbol[:]), "\x00")
}
