// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk_test

import (
	"context"
	"math/big"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/gapcover"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/margin"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/uniswapv3pool"
)

var poolAddress = common.HexToAddress("0x1111111111111111111111111111111111111111")

type engine struct {
	bind.ContractBackend
	t      *testing.T
	blocks []*big.Int
}

func (e *engine) HeaderByNumber(context.Context, *big.Int) (*types.Header, error) {
	return &types.Header{Number: big.NewInt(7)}, nil
}

func (e *engine) CallContract(_ context.Context, call ethereum.CallMsg, block *big.Int) ([]byte, error) {
	e.blocks = append(e.blocks, block)
	parsed, _ := margin.MarginMetaData.ParseABI()
	method, err := parsed.MethodById(call.Data)
	if err != nil {
		e.t.Fatal(err)
	}
	nvda, spy := bytes32(e.t, "NVDA"), bytes32(e.t, "SPY")
	args, _ := method.Inputs.Unpack(call.Data[4:])
	switch method.Name {
	case "assets":
		return method.Outputs.Pack([][32]byte{nvda, spy})
	case "volatility":
		return method.Outputs.Pack(uint32(31_352), uint32(25_000))
	case "correlation":
		return method.Outputs.Pack(uint16(7_100), uint16(5_000))
	case "weekendGap":
		return method.Outputs.Pack(uint32(118_601), uint32(100_000))
	case "depth":
		return method.Outputs.Pack(uint32(3_561_773), uint32(1_377_156), uint32(4_000_000), uint32(2_000_000))
	case "pool":
		if args[0] != nvda {
			e.t.Fatalf("pool of %v", args)
		}
		return method.Outputs.Pack(poolAddress)
	case "market":
		return method.Outputs.Pack(spy)
	case "lastUpdate":
		return method.Outputs.Pack(uint64(1_790_000_000))
	case "owner":
		return method.Outputs.Pack(alice)
	}
	e.t.Fatalf("unexpected %s", method.Name)
	return nil, nil
}

func TestTheEnginesParametersAreReadAtOneBlock(t *testing.T) {
	backend := &engine{t: t}
	m := sdk.NewClient(backend, registry(t, "46630")).Margin()
	opts := &bind.CallOpts{BlockNumber: big.NewInt(7)}
	nvda, spy := bytes32(t, "NVDA"), bytes32(t, "SPY")
	assets, err := m.Assets(opts)
	if err != nil || !slices.Equal(assets, [][32]byte{nvda, spy}) {
		t.Fatalf("assets %x, %v", assets, err)
	}
	if v, f, err := m.Volatility(opts, nvda); err != nil || v != 31_352 || f != 25_000 {
		t.Fatalf("volatility %d %d, %v", v, f, err)
	}
	if v, f, err := m.Correlation(opts, nvda, spy); err != nil || v != 7_100 || f != 5_000 {
		t.Fatalf("correlation %d %d, %v", v, f, err)
	}
	if v, f, err := m.WeekendGap(opts, nvda); err != nil || v != 118_601 || f != 100_000 {
		t.Fatalf("gap %d %d, %v", v, f, err)
	}
	if d, err := m.Depth(opts, nvda); err != nil || d.Selling != 3_561_773 || d.Buying != 1_377_156 || d.SellingCeiling != 4_000_000 || d.BuyingCeiling != 2_000_000 {
		t.Fatalf("depth %+v, %v", d, err)
	}
	if p, err := m.Pool(opts, nvda); err != nil || p != poolAddress {
		t.Fatalf("pool %v, %v", p, err)
	}
	if market, err := m.Market(opts); err != nil || market != spy {
		t.Fatalf("market %x, %v", market, err)
	}
	if last, err := m.LastUpdate(opts); err != nil || last != 1_790_000_000 {
		t.Fatalf("lastUpdate %d, %v", last, err)
	}
	if owner, err := m.Owner(opts); err != nil || owner != alice {
		t.Fatalf("owner %v, %v", owner, err)
	}
	if len(backend.blocks) != 9 {
		t.Fatalf("%d reads", len(backend.blocks))
	}
	for _, b := range backend.blocks {
		if b == nil || b.Int64() != 7 {
			t.Fatalf("read at %v", backend.blocks)
		}
	}
}

func TestSetParametersPacksTheFourArrays(t *testing.T) {
	m := sdk.NewClient(nil, registry(t, "46630")).Margin()
	tx, err := m.SetParameters([]uint32{31_352}, []uint16{}, []uint32{118_601}, []uint32{3_561_773, 1_377_156})
	if err != nil {
		t.Fatal(err)
	}
	want, _ := margin.NewMargin().TryPackSetParameters([]uint32{31_352}, []uint16{}, []uint32{118_601}, []uint32{3_561_773, 1_377_156})
	address, _ := m.Address()
	if tx.To != common.HexToAddress("0x9B2DB8135222d7B05aEA29B54aE0317E8640D6B0") || tx.To != address || !slices.Equal(tx.Data, want) {
		t.Fatalf("%x to %v", tx.Data, tx.To)
	}
}

type pool struct {
	bind.ContractBackend
	t      *testing.T
	blocks []*big.Int
}

func (p *pool) CallContract(_ context.Context, call ethereum.CallMsg, block *big.Int) ([]byte, error) {
	p.blocks = append(p.blocks, block)
	if call.To == nil || *call.To != poolAddress {
		p.t.Fatalf("call to %v", call.To)
	}
	parsed, _ := uniswapv3pool.UniswapV3PoolMetaData.ParseABI()
	method, err := parsed.MethodById(call.Data)
	if err != nil {
		p.t.Fatal(err)
	}
	args, _ := method.Inputs.Unpack(call.Data[4:])
	switch method.Name {
	case "slot0":
		sqrt, _ := new(big.Int).SetString("5237137002579766924522963008881898", 10)
		return method.Outputs.Pack(sqrt, big.NewInt(221_990), uint16(1), uint16(2), uint16(3), uint8(0), true)
	case "liquidity":
		l, _ := new(big.Int).SetString("11245526858841909681", 10)
		return method.Outputs.Pack(l)
	case "tickSpacing":
		return method.Outputs.Pack(big.NewInt(10))
	case "fee":
		return method.Outputs.Pack(big.NewInt(500))
	case "token0":
		return method.Outputs.Pack(alice)
	case "token1":
		return method.Outputs.Pack(bob)
	case "tickBitmap":
		if args[0].(int16) != 866 {
			p.t.Fatalf("tickBitmap of %v", args)
		}
		return method.Outputs.Pack(big.NewInt(0b1010))
	case "ticks":
		if args[0].(*big.Int).Int64() != 221_990 {
			p.t.Fatalf("ticks of %v", args)
		}
		return method.Outputs.Pack(big.NewInt(5), big.NewInt(-5), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), uint32(0), true)
	}
	p.t.Fatalf("unexpected %s", method.Name)
	return nil, nil
}

func TestAPoolsStateAndTicksAreRead(t *testing.T) {
	backend := &pool{t: t}
	p := sdk.NewClient(backend, registry(t, "46630")).Pool(poolAddress)
	opts := &bind.CallOpts{BlockNumber: big.NewInt(9)}
	slot0, err := p.Slot0(opts)
	if err != nil || slot0.SqrtPriceX96.String() != "5237137002579766924522963008881898" || slot0.Tick.Int64() != 221_990 {
		t.Fatalf("slot0 %+v, %v", slot0, err)
	}
	if l, err := p.Liquidity(opts); err != nil || l.String() != "11245526858841909681" {
		t.Fatalf("liquidity %v, %v", l, err)
	}
	if s, err := p.TickSpacing(opts); err != nil || s != 10 {
		t.Fatalf("spacing %d, %v", s, err)
	}
	if f, err := p.Fee(opts); err != nil || f != 500 {
		t.Fatalf("fee %d, %v", f, err)
	}
	if a, b, err := p.Tokens(opts); err != nil || a != alice || b != bob {
		t.Fatalf("tokens %v %v, %v", a, b, err)
	}
	if w, err := p.TickBitmap(opts, 866); err != nil || w.Int64() != 0b1010 {
		t.Fatalf("bitmap %v, %v", w, err)
	}
	if n, err := p.LiquidityNet(opts, 221_990); err != nil || n.Int64() != -5 {
		t.Fatalf("net %v, %v", n, err)
	}
	for _, b := range backend.blocks {
		if b == nil || b.Int64() != 9 {
			t.Fatalf("read at %v", backend.blocks)
		}
	}
}

func TestTheEnginesRefusalsAreDecoded(t *testing.T) {
	nvda := bytes32(t, "NVDA")
	for _, c := range []struct {
		name string
		args []any
		want string
	}{
		{"UpdateTooSoon", []any{uint64(1_790_000_000)}, "UpdateTooSoon(1790000000)"},
		{"VolatilityStepTooLarge", []any{nvda, uint32(31_352), uint32(47_029)}, "VolatilityStepTooLarge(0x4e56444100000000000000000000000000000000000000000000000000000000, 31352, 47029)"},
		{"NotPositiveDefinite", nil, "NotPositiveDefinite()"},
		{"InvalidDepth", []any{nvda, uint32(0), uint32(5)}, "InvalidDepth(0x4e56444100000000000000000000000000000000000000000000000000000000, 0, 5)"},
		{"OwnableUnauthorizedAccount", []any{bob}, "OwnableUnauthorizedAccount(" + bob.Hex() + ")"},
	} {
		revert, ok := sdk.DecodeRevertData(encodeError(t, &margin.MarginMetaData, c.name, c.args...))
		if !ok || revert.Error() != c.want {
			t.Errorf("%s: %v, want %s", c.name, revert, c.want)
		}
	}
}

func TestAHeaderIsReadAtABlock(t *testing.T) {
	header, err := sdk.NewClient(&engine{t: t}, registry(t, "46630")).Header(context.Background(), nil)
	if err != nil || header.Number.Int64() != 7 {
		t.Fatalf("header %v, %v", header, err)
	}
}

type tailBackend struct {
	bind.ContractBackend
}

func (tailBackend) CallContract(_ context.Context, call ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	parsed, _ := gapcover.GapCoverMetaData.ParseABI()
	method, err := parsed.MethodById(call.Data)
	if err != nil {
		return nil, err
	}
	values := map[string]int64{"TAIL_THRESHOLD_PPM": 183_288, "TAIL_SCALE_PPM": 95_837, "TAIL_PROBABILITY_PPM": 50_224}
	return method.Outputs.Pack(big.NewInt(values[method.Name]))
}

func TestTheGapCoversTailConstantsAreRead(t *testing.T) {
	d, err := sdk.ParseDeployments([]byte(`{"chainId":1,"tapehouse":{"GapCover":"` + alice.Hex() + `"}}`))
	if err != nil {
		t.Fatal(err)
	}
	threshold, scale, probability, err := sdk.NewClient(tailBackend{}, d).GapCover().Tail(&bind.CallOpts{})
	if err != nil || threshold != 183_288 || scale != 95_837 || probability != 50_224 {
		t.Fatalf("%d %d %d, %v", threshold, scale, probability, err)
	}
}
