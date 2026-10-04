// SPDX-License-Identifier: MIT OR Apache-2.0

package alerts

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tapehouse/tapehouse/services/internal/band"
	"github.com/tapehouse/tapehouse/services/internal/margin"
	bandbinding "github.com/tapehouse/tapehouse/services/sdk/bindings/band"
	marginbinding "github.com/tapehouse/tapehouse/services/sdk/bindings/margin"
)

func TestAMarketAndAPositionAreReadAtOneBlock(t *testing.T) {
	w := newWorld(t)
	w.holdings = map[string]*big.Int{"collateral NVDA": units(8, 18), "lent NVDA": units(2, 18), "basket NVDA": units(1, 18),
		"collateral SPY": units(1, 18), "collateral USDG": units(500, 6), "collateral WETH": units(1, 18),
		"debt": units(700, 6), "premium": units(3, 6)}
	r := NewReader(w.client)
	m, err := r.Market(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	account := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	j, err := r.Judge(context.Background(), m, account, [32]byte{7})
	if err != nil {
		t.Fatal(err)
	}
	if m.Block != 100 {
		t.Fatalf("the market is of block %d", m.Block)
	}
	for i, block := range w.calls {
		if block == nil || block.Uint64() != 100 {
			t.Fatalf("call %d read block %v, not 100, while the head moved to %d", i, block, w.head)
		}
	}
	p := j.Position
	if len(p.Stocks) != 2 || p.Stocks[0].Asset != "NVDA" || p.Stocks[0].Quantity.Cmp(units(11, 18)) != 0 ||
		p.Stocks[0].Lent.Cmp(units(2, 18)) != 0 || p.Stocks[1].Asset != "SPY" {
		t.Fatalf("stocks %+v", p.Stocks)
	}
	if p.USDG.Cmp(units(500, 6)) != 0 || p.WETH.Cmp(units(1, 18)) != 0 || p.Debt.Cmp(units(700, 6)) != 0 || p.Premium.Cmp(units(3, 6)) != 0 {
		t.Fatalf("position %+v", p)
	}
	if !j.Evaluated || !j.Valuation.Judged || j.Valuation.Requirement.Sign() == 0 {
		t.Fatalf("judgement %+v", j.Valuation)
	}
	if m.Assets[2].Basis != band.IndexBasisBps || m.Assets[0].Basis != 0 || m.Assets[0].Mid != units(100, 8).Uint64() {
		t.Fatalf("assets %+v", m.Assets)
	}
	if m.Regime.Kind != margin.Open || m.Closed || m.EthPrice.Cmp(units(2000, 8)) != 0 || m.RecallHaircut.Int64() != 500 {
		t.Fatalf("market %+v", m)
	}
}

func TestTheGoRequirementIsCheckedAgainstTheEngines(t *testing.T) {
	account := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	for name, test := range map[string]struct {
		engine func(requirement *big.Int, regime uint8) (*big.Int, uint8)
		want   bool
	}{
		"the engine agrees":          {func(r *big.Int, g uint8) (*big.Int, uint8) { return r, g }, true},
		"a requirement one unit off": {func(r *big.Int, g uint8) (*big.Int, uint8) { return new(big.Int).Add(r, n(1)), g }, false},
		"another regime":             {func(r *big.Int, _ uint8) (*big.Int, uint8) { return r, 1 }, false},
	} {
		t.Run(name, func(t *testing.T) {
			w := newWorld(t)
			w.holdings = map[string]*big.Int{"collateral NVDA": units(10, 18), "debt": units(100, 6)}
			good := w.requirement
			w.requirement = func(q, p []*big.Int) (*big.Int, uint8) {
				r, g := good(q, p)
				return test.engine(r, g)
			}
			r := NewReader(w.client)
			m, err := r.Market(context.Background(), nil)
			if err != nil {
				t.Fatal(err)
			}
			j, err := r.Judge(context.Background(), m, account, [32]byte{})
			if err != nil {
				t.Fatal(err)
			}
			if j.Evaluated != test.want {
				t.Fatalf("evaluated %v, want %v", j.Evaluated, test.want)
			}
		})
	}
}

func TestAnEngineThatRevertsLeavesThePositionNotEvaluated(t *testing.T) {
	w := newWorld(t)
	w.holdings = map[string]*big.Int{"collateral NVDA": units(10, 18), "debt": units(100, 6)}
	parsed, err := marginbinding.MarginMetaData.ParseABI()
	if err != nil {
		t.Fatal(err)
	}
	failure := parsed.Errors["LengthMismatch"]
	w.engine.fail("currentRequirement", revertWith{failure.ID[:4]})
	r := NewReader(w.client)
	m, err := r.Market(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	j, err := r.Judge(context.Background(), m, common.Address{1}, [32]byte{})
	if err != nil || j.Evaluated {
		t.Fatalf("a revert of the engine is its answer: evaluated %v, error %v", j != nil && j.Evaluated, err)
	}
}

type revertWith struct{ data []byte }

func (revertWith) Error() string            { return "execution reverted" }
func (revertWith) ErrorCode() int           { return 3 }
func (r revertWith) ErrorData() interface{} { return hexutil.Encode(r.data) }

func TestPoolsAreReadAsTheEngineReadsThem(t *testing.T) {
	w := newWorld(t)
	observed := func(then, now, secondsThen, secondsNow int64) func() ([]any, error) {
		return func() ([]any, error) {
			return []any{[]*big.Int{n(then), n(now)}, []*big.Int{n(secondsThen), n(secondsNow)}}, nil
		}
	}
	// NVDA against USDG, SPY against WETH, TSLA's pool without enough history.
	w.pool("NVDA", addrNVDAPool, addrNVDAToken, addrUSDG, 500, observed(1_000_000, 1_000_000-1_800*230, 5_000, 5_000+1_800*100))
	w.pool("SPY", addrSPYPool, addrSPYToken, addrWETH, 3000, observed(2_000, 2_000+1_800*(-5), 7_000, 7_000+1_800*50))
	w.pool("TSLA", addrTSLAPool, addrTSLAToken, addrUSDG, 3000, func() ([]any, error) { return nil, errors.New("OLD") })
	r := NewReader(w.client)
	m, err := r.Market(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	expect := func(stockIsToken0 bool, quoteDecimals uint8, then, now int64, secondsThen, secondsNow int64, usd *big.Int, fee uint32) margin.Pool {
		tick, liquidity, ok := margin.Mean(then, now, n(secondsThen), n(secondsNow))
		if !ok {
			t.Fatal("no mean")
		}
		terms, ok := margin.Terms(stockIsToken0, 18, quoteDecimals, tick, liquidity, usd)
		if !ok {
			t.Fatal("no terms")
		}
		terms.Fee = fee
		return margin.Pool{Kind: margin.Read, Terms: terms}
	}
	nvdaIs0 := addrNVDAToken.Cmp(addrUSDG) < 0
	spyIs0 := addrSPYToken.Cmp(addrWETH) < 0
	nvda := expect(nvdaIs0, 6, 1_000_000, 1_000_000-1_800*230, 5_000, 5_000+1_800*100, n(1e8), 500)
	spy := expect(spyIs0, 18, 2_000, 2_000+1_800*(-5), 7_000, 7_000+1_800*50, units(2000, 8), 3000)
	for i, want := range []margin.Pool{nvda, {Kind: margin.Unread, Terms: margin.PoolTerms{Fee: 3000}}, spy} {
		got := m.Pools[i]
		if got.Kind != want.Kind || got.Terms.Fee != want.Terms.Fee {
			t.Fatalf("pool %d: %+v, want %+v", i, got, want)
		}
		if want.Kind == margin.Read && (got.Terms.Price.Cmp(want.Terms.Price) != 0 || got.Terms.Selling.Cmp(want.Terms.Selling) != 0 ||
			got.Terms.Buying.Cmp(want.Terms.Buying) != 0) {
			t.Fatalf("pool %d: %+v, want %+v", i, got.Terms, want.Terms)
		}
	}
	w.ethUpdated = n(1_790_775_000 - 86_461)
	if m, err = r.Market(context.Background(), nil); err != nil || m.Pools[2].Kind != margin.Unread || m.EthPrice.Sign() != 0 {
		t.Fatalf("a stale ETH/USD must leave a WETH pool unread and WETH unpriced: %+v, %v", m.Pools[2], err)
	}
}

func TestTheLiquidatorsClosedMarketIsRebuilt(t *testing.T) {
	for name, test := range map[string]struct {
		session bandbinding.SessionOutput
		nowMs   uint64
		closed  bool
		regime  margin.RegimeKind
	}{
		"open regular hours":        {bandbinding.SessionOutput{State: 2, Nyse: 1, NyseNext: 1}, 10, false, margin.Open},
		"session not known":         {bandbinding.SessionOutput{State: 0, Nyse: 1, NyseNext: 1}, 10, true, margin.Unknown},
		"session closed":            {bandbinding.SessionOutput{State: 1, Nyse: 3, NyseNext: 3}, 10, true, margin.Closed},
		"boundary passed":           {bandbinding.SessionOutput{State: 2, Nyse: 1, NyseNext: 1, BoundaryMs: 10}, 10, true, margin.Closed},
		"post market, trading day":  {bandbinding.SessionOutput{State: 2, Nyse: 2, NyseNext: 1}, 10, false, margin.Open},
		"post market, weekend next": {bandbinding.SessionOutput{State: 2, Nyse: 2, NyseNext: 3}, 10, true, margin.Open},
		"in the ramp":               {bandbinding.SessionOutput{State: 2, Nyse: 1, NyseNext: 2, BoundaryMs: 10 + margin.RampMs/2}, 10, false, margin.Closing},
	} {
		t.Run(name, func(t *testing.T) {
			if got := closedAt(test.session, test.nowMs); got != test.closed {
				t.Fatalf("closed %v, want %v", got, test.closed)
			}
			if got := margin.RegimeAt(test.session.State, test.session.BoundaryMs, test.nowMs).Kind; got != test.regime {
				t.Fatalf("regime %d, want %d", got, test.regime)
			}
		})
	}
}
