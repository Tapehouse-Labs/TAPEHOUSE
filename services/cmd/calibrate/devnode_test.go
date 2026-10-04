// SPDX-License-Identifier: MIT OR Apache-2.0

//go:build devnode

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapehouse/tapehouse/services/internal/calibrate"
	"github.com/tapehouse/tapehouse/services/sdk"
)

const root = "../../.."

// shell runs a command in dir, relative to the repository's root, and returns its output.
func shell(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = filepath.Join(root, dir)
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("%s: %v\n%s", name, err, out.String())
	}
	return out.String()
}

// raised returns the margin program's assets argument with every volatility a fifth higher, which its floors allow.
func raised(t *testing.T, assets string) string {
	t.Helper()
	tuple := regexp.MustCompile(`\(([^()]*)\)`)
	return tuple.ReplaceAllStringFunc(assets, func(entry string) string {
		fields := strings.Split(strings.Trim(entry, "()"), ",")
		volatility, err := strconv.ParseUint(fields[2], 10, 32)
		if err != nil {
			t.Fatal(err)
		}
		fields[2] = strconv.FormatUint(volatility*6/5, 10)
		return "(" + strings.Join(fields, ",") + ")"
	})
}

// estimate is the L2 gas of the proposal's call from its owner, as Arbitrum's NodeInterface splits its estimate.
func estimate(t *testing.T, rpcURL string, p calibrate.Proposal) uint64 {
	t.Helper()
	parsed, err := abi.JSON(strings.NewReader(`[{"type":"function","name":"gasEstimateComponents","stateMutability":"view",
		"inputs":[{"name":"to","type":"address"},{"name":"contractCreation","type":"bool"},{"name":"data","type":"bytes"}],
		"outputs":[{"name":"gasEstimate","type":"uint64"},{"name":"gasEstimateForL1","type":"uint64"},{"name":"baseFee","type":"uint256"},{"name":"l1BaseFeeEstimate","type":"uint256"}]}]`))
	if err != nil {
		t.Fatal(err)
	}
	input, err := parsed.Pack("gasEstimateComponents", p.Margin, false, []byte(p.Calldata))
	if err != nil {
		t.Fatal(err)
	}
	chain, err := ethclient.Dial(rpcURL)
	if err != nil {
		t.Fatal(err)
	}
	defer chain.Close()
	nodeInterface := common.HexToAddress("0x00000000000000000000000000000000000000C8")
	raw, err := chain.CallContract(context.Background(), ethereum.CallMsg{From: p.Owner, To: &nodeInterface, Data: input}, nil)
	if err != nil {
		t.Fatal(err)
	}
	out, err := parsed.Unpack("gasEstimateComponents", raw)
	if err != nil {
		t.Fatal(err)
	}
	return out[0].(uint64) - out[1].(uint64)
}

func TestTheCalibratorUpdatesAFreshEngine(t *testing.T) {
	ctx := context.Background()
	rpcURL, registry, key := os.Getenv("TAPEHOUSE_RPC_URL"), os.Getenv("TAPEHOUSE_DEPLOYMENTS"), os.Getenv("PRIVATE_KEY")
	if rpcURL == "" || registry == "" || key == "" {
		t.Fatal("set TAPEHOUSE_RPC_URL, TAPEHOUSE_DEPLOYMENTS and PRIVATE_KEY for the dev node")
	}
	dir := t.TempDir()
	launch := filepath.Join(dir, "launch.json")
	data, err := os.ReadFile(registry)
	if err != nil {
		t.Fatal(err)
	}
	var fresh map[string]any
	if err := json.Unmarshal(data, &fresh); err != nil {
		t.Fatal(err)
	}
	feeds := fresh["chainlink"].(map[string]any)
	for _, asset := range []string{"AAPL", "MSFT", "GOOGL"} {
		feeds[asset+"_USD"] = feeds["NVDA_USD"]
	}
	if data, err = json.Marshal(fresh); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(launch, data, 0o644); err != nil {
		t.Fatal(err)
	}
	args := strings.Fields(shell(t, ".", "stylus/scripts/margin-args.sh", launch))
	args[0] = raised(t, args[0])
	deployed := shell(t, "stylus/contracts/margin", "cargo", append([]string{"stylus", "deploy", "--no-verify", "-e", rpcURL, "--private-key", key, "--constructor-args"}, args...)...)
	address := regexp.MustCompile(`deployed code at address.*?(0x[0-9a-fA-F]{40})`).FindStringSubmatch(deployed)
	if address == nil {
		t.Fatalf("no deployed address in %s", deployed)
	}
	fresh["tapehouse"].(map[string]any)["Margin"] = address[1]
	data, _ = json.Marshal(fresh)
	copied := filepath.Join(dir, "registry.json")
	if err := os.WriteFile(copied, data, 0o644); err != nil {
		t.Fatal(err)
	}
	env := func(name string) string {
		switch name {
		case "TAPEHOUSE_DEPLOYMENTS":
			return copied
		case "TAPEHOUSE_CALIBRATOR_DIR":
			return filepath.Join(dir, "calibrate")
		}
		return os.Getenv(name)
	}

	proposalFile := filepath.Join(dir, "proposal.json")
	var out bytes.Buffer
	if err := run(ctx, []string{"propose", "-as-of", "2026-09-25", "-out", proposalFile}, env, &out); err != nil {
		t.Fatalf("%v\n%s", err, out.String())
	}
	t.Log(out.String())
	var p calibrate.Proposal
	data, err = os.ReadFile(proposalFile)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &p); err != nil {
		t.Fatal(err)
	}
	if p.Simulation != "ok" || len(p.Assets) != 6 {
		t.Fatalf("simulation %q over %v", p.Simulation, p.Assets)
	}
	want := map[string]int{}
	for i := range p.Assets {
		if p.Proposed.Volatilities[i] != p.Base.Volatilities[i] {
			want["VolatilitySet"]++
		}
		if p.Proposed.Gaps[i] != p.Base.Gaps[i] {
			want["GapSet"]++
		}
		if p.Proposed.Depths[2*i] != p.Base.Depths[2*i] || p.Proposed.Depths[2*i+1] != p.Base.Depths[2*i+1] {
			want["DepthSet"]++
		}
	}
	for k := range p.Proposed.Correlations {
		if p.Proposed.Correlations[k] != p.Base.Correlations[k] {
			want["CorrelationSet"]++
		}
	}
	if want["VolatilitySet"] != 6 {
		t.Fatalf("%d volatilities change, want all six", want["VolatilitySet"])
	}

	gas := estimate(t, rpcURL, p)
	t.Logf("setParameters of six assets: %d L2 gas", gas)
	if gas < 193_983*9/10 || gas > 193_983*11/10 {
		t.Errorf("L2 gas %d, want 193,983 within 10%%", gas)
	}

	out.Reset()
	if err := run(ctx, []string{"apply", proposalFile}, env, &out); err != nil {
		t.Fatalf("%v\n%s", err, out.String())
	}
	t.Log(out.String())
	for _, event := range []string{"VolatilitySet", "GapSet", "DepthSet", "CorrelationSet"} {
		if got := strings.Count(out.String(), "\n"+event+" "); got != want[event] {
			t.Errorf("%d %s events, want %d\n%s", got, event, want[event], out.String())
		}
	}

	chain, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		t.Fatal(err)
	}
	defer chain.Close()
	deployments, err := sdk.LoadDeployments(copied)
	if err != nil {
		t.Fatal(err)
	}
	engine := sdk.NewClient(chain, deployments).Margin()
	opts := &bind.CallOpts{Context: ctx}
	assets, err := engine.Assets(opts)
	if err != nil {
		t.Fatal(err)
	}
	for i, symbol := range assets {
		if v, _, err := engine.Volatility(opts, symbol); err != nil || v != p.Proposed.Volatilities[i] {
			t.Errorf("%s volatility %d, %v, proposed %d", p.Assets[i], v, err, p.Proposed.Volatilities[i])
		}
		if v, _, err := engine.WeekendGap(opts, symbol); err != nil || v != p.Proposed.Gaps[i] {
			t.Errorf("%s gap %d, %v, proposed %d", p.Assets[i], v, err, p.Proposed.Gaps[i])
		}
		if d, err := engine.Depth(opts, symbol); err != nil || d.Selling != p.Proposed.Depths[2*i] || d.Buying != p.Proposed.Depths[2*i+1] {
			t.Errorf("%s depth %+v, %v", p.Assets[i], d, err)
		}
		for j, other := range assets[i+1:] {
			k := i*len(assets) - i*(i+1)/2 + j
			if v, _, err := engine.Correlation(opts, symbol, other); err != nil || v != p.Proposed.Correlations[k] {
				t.Errorf("%s/%s correlation %d, %v, proposed %d", p.Assets[i], p.Assets[i+1+j], v, err, p.Proposed.Correlations[k])
			}
		}
	}
	last, err := engine.LastUpdate(opts)
	if err != nil {
		t.Fatal(err)
	}
	out.Reset()
	err = run(ctx, []string{"propose", "-as-of", "2026-09-25", "-out", proposalFile}, env, &out)
	if want := time.Unix(int64(last)+86_400, 0).UTC().Format(time.RFC3339); err == nil || !strings.Contains(err.Error(), "less than a day") || !strings.Contains(err.Error(), want) {
		t.Fatalf("a second proposal: %v, want a refusal naming %s", err, want)
	}
}
