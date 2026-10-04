// SPDX-License-Identifier: MIT OR Apache-2.0

package calibrate

import (
	"context"
	"errors"
	"math/big"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tapehouse/tapehouse/services/internal/backtest"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/margin"
)

var (
	owner       = common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	engineAddr  = common.HexToAddress("0x9B2DB8135222d7B05aEA29B54aE0317E8640D6B0")
	engineNames = []string{"AAA", "BBB", "MKT"}
)

// engine is a margin engine as the chain answers for it: its assets, values, floors, ceilings, last update and
// owner, and a setParameters that succeeds or reverts with revert.
type engine struct {
	bind.ContractBackend
	t          *testing.T
	volatility []uint32
	lastUpdate uint64
	revert     error
	simulated  []ethereum.CallMsg
}

func (e *engine) HeaderByNumber(context.Context, *big.Int) (*types.Header, error) {
	return &types.Header{Number: big.NewInt(7), Time: 1_790_000_000}, nil
}

func (e *engine) CodeAt(context.Context, common.Address, *big.Int) ([]byte, error) {
	return []byte{1}, nil
}

func (e *engine) CallContract(_ context.Context, call ethereum.CallMsg, block *big.Int) ([]byte, error) {
	if block == nil || block.Int64() != 7 {
		e.t.Fatalf("read at %v", block)
	}
	parsed, _ := margin.MarginMetaData.ParseABI()
	method, err := parsed.MethodById(call.Data)
	if err != nil {
		e.t.Fatal(err)
	}
	args, _ := method.Inputs.Unpack(call.Data[4:])
	index := func(symbol [32]byte) int {
		return slices.Index(engineNames, strings.TrimRight(string(symbol[:]), "\x00"))
	}
	switch method.Name {
	case "assets":
		var out [][32]byte
		for _, n := range engineNames {
			var s [32]byte
			copy(s[:], n)
			out = append(out, s)
		}
		return method.Outputs.Pack(out)
	case "volatility":
		return method.Outputs.Pack(e.volatility[index(args[0].([32]byte))], uint32(20_000))
	case "weekendGap":
		return method.Outputs.Pack(uint32(45_000), uint32(40_000))
	case "correlation":
		i, j := index(args[0].([32]byte)), index(args[1].([32]byte))
		return method.Outputs.Pack(uint16(3_000+500*(i+j)), uint16(1_000))
	case "depth":
		return method.Outputs.Pack(uint32(100_000), uint32(100_000), uint32(1_000_000), uint32(1_000_000))
	case "market":
		var s [32]byte
		copy(s[:], "MKT")
		return method.Outputs.Pack(s)
	case "lastUpdate":
		return method.Outputs.Pack(e.lastUpdate)
	case "owner":
		return method.Outputs.Pack(owner)
	case "setParameters":
		e.simulated = append(e.simulated, call)
		return nil, e.revert
	}
	e.t.Fatalf("unexpected %s", method.Name)
	return nil, nil
}

func syntheticEngine(t *testing.T) (*engine, *sdk.Client) {
	t.Helper()
	d, err := sdk.ParseDeployments([]byte(`{"chainId":412346,"tapehouse":{"Margin":"` + engineAddr.Hex() + `"}}`))
	if err != nil {
		t.Fatal(err)
	}
	e := &engine{t: t, volatility: []uint32{25_000, 22_000, 21_000}}
	return e, sdk.NewClient(e, d)
}

func syntheticMarketOf(t *testing.T) backtest.Market {
	t.Helper()
	return backtest.Align(syntheticMarket())
}

func TestNoProposalWithinADayOfTheLastUpdate(t *testing.T) {
	e, client := syntheticEngine(t)
	now := time.Unix(1_790_086_399, 0)
	e.lastUpdate = uint64(now.Unix()) - 86_399
	_, err := Propose(context.Background(), client, syntheticMarketOf(t), nil, now)
	if !errors.Is(err, ErrTooSoon) || !strings.Contains(err.Error(), time.Unix(int64(e.lastUpdate)+86_400, 0).UTC().Format(time.RFC3339)) {
		t.Fatalf("%v", err)
	}
	e.lastUpdate = uint64(now.Unix()) - 86_400
	if _, err := Propose(context.Background(), client, syntheticMarketOf(t), nil, now); err != nil {
		t.Fatal(err)
	}
}

func TestAProposalCarriesItsCalldataAndBothReports(t *testing.T) {
	e, client := syntheticEngine(t)
	m := syntheticMarketOf(t)
	p, err := Propose(context.Background(), client, m, nil, time.Unix(1_790_000_000, 0))
	if err != nil {
		t.Fatal(err)
	}
	targets, err := Estimate(syntheticMarket(), engineNames, 2, m.At[len(m.At)-1])
	if err != nil {
		t.Fatal(err)
	}
	current := Current{
		Volatilities: []uint32{25_000, 22_000, 21_000}, VolatilityFloors: []uint32{20_000, 20_000, 20_000},
		Gaps: []uint32{45_000, 45_000, 45_000}, GapFloors: []uint32{40_000, 40_000, 40_000},
		Correlations: []uint16{3_500, 4_000, 4_500}, CorrelationFloors: []uint16{1_000, 1_000, 1_000},
		Depths: []uint32{100_000, 100_000, 100_000, 100_000, 100_000, 100_000}, DepthCeilings: slices.Repeat([]uint32{1_000_000}, 6),
	}
	want, err := Step(current, targets, nil)
	if err != nil {
		t.Fatal(err)
	}
	if p.Block != 7 || p.Margin != engineAddr || p.Chain != 412346 || !slices.Equal(p.Assets, engineNames) {
		t.Fatalf("%+v", p)
	}
	if !slices.Equal(p.Base.Volatilities, current.Volatilities) || !slices.Equal(p.Base.Correlations, current.Correlations) || !slices.Equal(p.Base.DepthCeilings, current.DepthCeilings) {
		t.Fatalf("base %+v", p.Base)
	}
	if !slices.Equal(p.Proposed.Volatilities, want.Volatilities) || !slices.Equal(p.Proposed.Correlations, want.Correlations) ||
		!slices.Equal(p.Proposed.Gaps, want.Gaps) || !slices.Equal(p.Proposed.Depths, want.Depths) {
		t.Fatalf("proposed %+v, want %+v", p.Proposed, want)
	}
	tx, _ := client.Margin().SetParameters(want.Volatilities, want.Correlations, want.Gaps, want.Depths)
	if !slices.Equal(p.Calldata, tx.Data) {
		t.Fatalf("calldata %x", p.Calldata)
	}
	if p.Simulation != "ok" || len(e.simulated) != 1 || e.simulated[0].From != owner || e.simulated[0].To == nil || *e.simulated[0].To != engineAddr {
		t.Fatalf("simulation %q from %+v", p.Simulation, e.simulated)
	}
	if p.Current.History.Sessions == 0 || p.Next.History.Sessions != p.Current.History.Sessions || len(p.ReferenceVolatilities) != 3 {
		t.Fatalf("reports %+v / %+v, reference %v", p.Current.History, p.Next.History, p.ReferenceVolatilities)
	}
}

func TestASimulationThatRevertsOtherwiseIsAnError(t *testing.T) {
	e, client := syntheticEngine(t)
	e.revert = rpcRevert(t, "NotPositiveDefinite")
	if _, err := Propose(context.Background(), client, syntheticMarketOf(t), nil, time.Unix(1_790_000_000, 0)); err == nil {
		t.Fatal("a refused simulation was reported as clean")
	}
	e.revert = rpcRevert(t, "UpdateTooSoon", uint64(1_790_086_400))
	p, err := Propose(context.Background(), client, syntheticMarketOf(t), nil, time.Unix(1_790_000_000, 0))
	if err != nil || p.Simulation != "UpdateTooSoon(1790086400)" {
		t.Fatalf("%q, %v", p.Simulation, err)
	}
}

type rpcError struct{ data string }

func (e rpcError) Error() string          { return "execution reverted" }
func (e rpcError) ErrorCode() int         { return 3 }
func (e rpcError) ErrorData() interface{} { return e.data }

func rpcRevert(t *testing.T, name string, args ...any) error {
	t.Helper()
	parsed, _ := margin.MarginMetaData.ParseABI()
	e := parsed.Errors[name]
	packed, err := e.Inputs.Pack(args...)
	if err != nil {
		t.Fatal(err)
	}
	return rpcError{hexutil.Encode(append(e.ID[:4:4], packed...))}
}

func TestAProposalWhoseBaseChangedIsRefused(t *testing.T) {
	e, client := syntheticEngine(t)
	p, err := Propose(context.Background(), client, syntheticMarketOf(t), nil, time.Unix(1_790_000_000, 0))
	if err != nil {
		t.Fatal(err)
	}
	if err := CheckBase(context.Background(), client, p); err != nil {
		t.Fatal(err)
	}
	e.volatility[1]++
	err = CheckBase(context.Background(), client, p)
	if !errors.Is(err, ErrStaleBase) || !strings.Contains(err.Error(), "BBB volatility") {
		t.Fatalf("%v", err)
	}
	e.volatility[1]--
	e.lastUpdate = 1_790_000_001
	if err := CheckBase(context.Background(), client, p); !errors.Is(err, ErrStaleBase) || !strings.Contains(err.Error(), "last update") {
		t.Fatalf("%v", err)
	}
}

func TestAProposalTheBacktestFindsWorseIsRefused(t *testing.T) {
	var current, next backtest.Report
	current.Leverage.AtDeployedCap = map[string]backtest.Breaches{"long": {Portfolios: 0}, "mixed": {Portfolios: 3}}
	next.Leverage.AtDeployedCap = map[string]backtest.Breaches{"long": {Portfolios: 0}, "mixed": {Portfolios: 4}}
	if err := Gate(current, next); !errors.Is(err, ErrWorse) {
		t.Fatalf("%v", err)
	}
	next.Leverage.AtDeployedCap["mixed"] = backtest.Breaches{Portfolios: 3}
	current.Capacity.BadDebtModel, next.Capacity.BadDebtModel = 0.001, 0.0011
	if err := Gate(current, next); !errors.Is(err, ErrWorse) {
		t.Fatalf("%v", err)
	}
	next.Capacity.BadDebtModel = 0.001
	if err := Gate(current, next); err != nil {
		t.Fatal(err)
	}
}

func TestAProposalIsDescribedLineByLine(t *testing.T) {
	p := Proposal{
		Assets:                []string{"NVDA", "TSLA"},
		Base:                  Current{Volatilities: []uint32{31352, 37436}, Gaps: []uint32{1, 1}, Depths: []uint32{9, 9, 9, 9}, Correlations: []uint16{4637}},
		Proposed:              Proposed{Volatilities: []uint32{33000, 37436}, Gaps: []uint32{1, 1}, Depths: []uint32{9, 9, 8, 9}, Correlations: []uint16{4637}},
		Targets:               Targets{Volatilities: []uint32{33000, 37436}, Gaps: []uint32{1, 1}, Correlations: []uint16{4637}},
		ReferenceVolatilities: []uint32{50000, 37000},
		Simulation:            "ok",
	}
	lines := p.Lines()
	for _, want := range []string{
		"NVDA volatility 31352 → 33000 (target 33000)",
		"TSLA selling depth 9 → 8 (target 8)",
		"reference volatility NVDA 50000 against the method's 33000: diverges beyond ×1.5",
		"simulation: ok",
	} {
		if !slices.Contains(lines, want) {
			t.Errorf("no line %q in %q", want, lines)
		}
	}
	for _, l := range lines {
		if strings.Contains(l, "TSLA volatility") || strings.Contains(l, "reference volatility TSLA") && strings.Contains(l, "diverges") {
			t.Errorf("line %q", l)
		}
	}
}
