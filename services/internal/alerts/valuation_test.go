// SPDX-License-Identifier: MIT OR Apache-2.0

package alerts

import (
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/internal/band"
	"github.com/tapehouse/tapehouse/services/internal/margin"
)

func n(v int64) *big.Int { return big.NewInt(v) }

func units(whole int64, decimals int) *big.Int {
	return new(big.Int).Mul(n(whole), new(big.Int).Exp(n(10), n(int64(decimals)), nil))
}

// fixture is the launch parameters of three assets over a dev-node-like market: the session open, no pools.
func fixture(t testing.TB) *Market {
	t.Helper()
	data, err := os.ReadFile("../../../stylus/contracts/margin/parameters.json")
	if err != nil {
		t.Fatal(err)
	}
	names := []string{"NVDA", "TSLA", "SPY"}
	parameters, depths, err := margin.Launch(data, names)
	if err != nil {
		t.Fatal(err)
	}
	m := &Market{
		Block: 100, TimeMs: 1_790_775_000_000, Settled: true, Regime: margin.Regime{Kind: margin.Open},
		Set: margin.NewSet(parameters, 256, margin.Horizon), Depths: depths,
		Pools:         make([]margin.Pool, len(names)),
		EthPrice:      units(2000, 8),
		RecallHaircut: n(500),
	}
	for i, name := range names {
		basis := uint64(0)
		if name == "SPY" {
			basis = band.IndexBasisBps
		}
		a := Asset{Name: name, State: uint8(band.Open), Mid: units(int64(100*(i+1)), 8).Uint64(), Variance: n(int64(40_000 + 10_000*i)), Basis: basis}
		q := openBand(a, a.Mid)
		a.Low, a.High = q.Low, q.High
		m.Assets = append(m.Assets, a)
	}
	return m
}

func holder(debt int64) *Position {
	return &Position{
		Account: common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"), ID: [32]byte{1},
		Stocks: []Stock{{Asset: "NVDA", Quantity: units(10, 18), Lent: n(0)}},
		USDG:   n(0), WETH: n(0), Debt: units(debt, 6), Premium: n(0),
	}
}

func TestTheValuationIsTheLiquidatorsOwn(t *testing.T) {
	m := fixture(t)
	p := &Position{
		Stocks: []Stock{{Asset: "NVDA", Quantity: units(12, 18), Lent: units(2, 18)}, {Asset: "SPY", Quantity: units(1, 18), Lent: n(0)}},
		USDG:   units(500, 6), WETH: units(1, 18), Debt: units(800, 6), Premium: units(5, 6),
	}
	edges := m.Current()
	v := m.Assess(p, edges, false, m.Regime)
	if !v.Judged || !v.Held {
		t.Fatalf("not judged: %+v", v)
	}
	lowNVDA, lowSPY := edges.Low[0], edges.Low[2]
	gross := new(big.Int).Div(new(big.Int).Mul(units(12, 18), lowNVDA), n(1e8))
	gross.Add(gross, new(big.Int).Div(new(big.Int).Mul(units(1, 18), lowSPY), n(1e8)))
	haircut := new(big.Int).Div(new(big.Int).Mul(new(big.Int).Mul(units(2, 18), lowNVDA), n(500)), n(1e12))
	gross.Sub(gross, haircut)
	eth := new(big.Int).Div(new(big.Int).Mul(units(1, 18), units(2000, 8)), n(1e8))
	want := new(big.Int).Add(gross, new(big.Int).Mul(units(500, 6), n(1e12)))
	want.Add(want, new(big.Int).Div(new(big.Int).Mul(eth, n(8400)), n(10_000)))
	want.Sub(want, new(big.Int).Mul(units(805, 6), n(1e12)))
	if v.Equity.Cmp(want) != 0 {
		t.Fatalf("equity %s, want %s", v.Equity, want)
	}
	quantities := []*big.Int{units(12, 18), n(0), units(1, 18)}
	prices := []*big.Int{lowNVDA, n(0), lowSPY}
	exposures, _ := margin.Exposures(quantities, prices)
	open, closed, floor, _ := margin.Requirements(m.Set, exposures, prices, m.Depths, m.Pools)
	if wantReq := margin.Current(open, closed, floor, m.Regime); v.Requirement.Cmp(wantReq) != 0 {
		t.Fatalf("requirement %s, want %s", v.Requirement, wantReq)
	}
}

func TestAClosedMarketJudgesAtBothEdges(t *testing.T) {
	m := fixture(t)
	edges := m.Current()
	var debt int64
	for d := int64(100); ; d++ {
		low := m.Value(holder(d), edges, m.Regime, false)
		if low.Short() {
			debt = d
			break
		}
		if d > 10_000 {
			t.Fatal("no debt makes the position short")
		}
	}
	p := holder(debt)
	if !m.Value(p, edges, m.Regime, false).Short() || m.Value(p, edges, m.Regime, true).Short() {
		t.Fatal("the position must be short at the low edge only")
	}
	if !m.Assess(p, edges, false, m.Regime).Short() {
		t.Fatal("an open market judges at the low edge")
	}
	closed := m.Assess(p, edges, true, m.Regime)
	if closed.Short() {
		t.Fatal("a closed market must not call a position short that is not short at its high edge")
	}
	if high := m.Value(p, edges, m.Regime, true); closed.Equity.Cmp(high.Equity) != 0 {
		t.Fatalf("the closed judgement takes the better surplus: %s, want %s", closed.Equity, high.Equity)
	}
}

func TestAPositionTheLiquidatorCannotJudgeIsNotPriced(t *testing.T) {
	for name, change := range map[string]func(*Market){
		"a halted band":                func(m *Market) { m.Assets[0].State = uint8(band.Halted) },
		"an unconfirmed multiplier":    func(m *Market) { m.Assets[0].Unconfirmed = true },
		"a sequencer not yet settled":  func(m *Market) { m.Settled = false },
		"a WETH holding without price": func(m *Market) { m.EthPrice = n(0) },
	} {
		t.Run(name, func(t *testing.T) {
			m := fixture(t)
			change(m)
			p := holder(100)
			p.WETH = units(1, 18)
			if v := m.Assess(p, m.Current(), false, m.Regime); v.Judged {
				t.Fatalf("judged: %+v", v)
			}
			if got := m.WeekendPrice(p, "NVDA"); got.Status != NotJudged {
				t.Fatalf("a weekend price for an unjudged position: %+v", got)
			}
			if got := m.ReopeningPrice(p, "NVDA"); got.Status != NotJudged {
				t.Fatalf("a reopening price for an unjudged position: %+v", got)
			}
		})
	}
	m := fixture(t)
	m.Assets[1].State = uint8(band.Halted)
	if v := m.Assess(holder(100), m.Current(), false, m.Regime); !v.Judged {
		t.Fatal("a halted asset the position does not hold must not stop its judgement")
	}
}

// boundary is the largest debt, in USDG, that leaves holder(debt) not short at both edges of a closed market at its
// current centre.
func boundary(m *Market, closed bool) int64 {
	edges := m.Current()
	regime := margin.Regime{Kind: margin.Open}
	if closed {
		regime = margin.Regime{Kind: margin.Closed}
		edges = m.centres(closedBand)
	}
	low, high := int64(0), int64(1_000_000)
	for high-low > 1 {
		mid := (low + high) / 2
		if m.Assess(holder(mid), edges, closed, regime).Short() {
			high = mid
		} else {
			low = mid
		}
	}
	return low
}

func TestTheWeekendPriceIsWhereTheLiquidatorWouldCallItShort(t *testing.T) {
	m := fixture(t)
	p := holder(boundary(m, true) * 97 / 100)
	got := m.WeekendPrice(p, "NVDA")
	if got.Status != Priced || got.Centre == 0 || got.Centre >= m.Assets[0].Mid {
		t.Fatalf("weekend price %+v at a centre of %d", got, m.Assets[0].Mid)
	}
	short := func(centre uint64) bool {
		edges := m.centres(closedBand)
		edges.Low[0], edges.High[0] = m.edge(0, closedBand, centre)
		return m.Assess(p, edges, true, margin.Regime{Kind: margin.Closed}).Short()
	}
	if !short(got.Centre) || short(got.Centre+1) {
		t.Fatalf("the position is short at %d: %v, at one unit above: %v", got.Centre, short(got.Centre), short(got.Centre+1))
	}
	if !short(got.Centre - 1) {
		t.Fatal("a lower price must be short too")
	}
}

func TestAPositionAlreadyShortReportsTheCurrentCentreAndOneNoPriceShortsReportsNone(t *testing.T) {
	m := fixture(t)
	if got := m.WeekendPrice(holder(boundary(m, true)*2), "NVDA"); got.Status != Priced || got.Centre != m.Assets[0].Mid {
		t.Fatalf("a position short now: %+v, want the centre %d", got, m.Assets[0].Mid)
	}
	if got := m.WeekendPrice(holder(0), "NVDA"); got.Status != NoPrice {
		t.Fatalf("a position without debt: %+v", got)
	}
	if got := m.WeekendPrice(holder(1), "TSLA"); got.Status != NoPrice {
		t.Fatalf("an asset the position does not hold: %+v", got)
	}
}

func TestTheReopeningPriceIsWhereTheLowEdgeIsShortAgainstTheOpenRequirement(t *testing.T) {
	m := fixture(t)
	p := holder(boundary(m, false) * 95 / 100)
	got := m.ReopeningPrice(p, "NVDA")
	if got.Status != Priced || got.Centre == 0 || got.Centre >= m.Assets[0].Mid {
		t.Fatalf("reopening price %+v at a centre of %d", got, m.Assets[0].Mid)
	}
	short := func(centre uint64) bool {
		edges := m.centres(openBand)
		edges.Low[0], edges.High[0] = m.edge(0, openBand, centre)
		return m.Assess(p, edges, false, margin.Regime{Kind: margin.Open}).Short()
	}
	if !short(got.Centre) || short(got.Centre+1) {
		t.Fatalf("short at %d: %v, one unit above: %v", got.Centre, short(got.Centre), short(got.Centre+1))
	}
	if weekend := m.WeekendPrice(p, "NVDA"); weekend.Status == Priced && weekend.Centre < got.Centre {
		t.Fatalf("the weekend price %d is below the reopening price %d", weekend.Centre, got.Centre)
	}
}
