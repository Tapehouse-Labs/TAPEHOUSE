// SPDX-License-Identifier: MIT OR Apache-2.0

// Package alerts warns the holders of margin positions before a weekend and before a liquidation. It values a
// position as the liquidator does, with the margin engine and the band evaluated in Go, searches the prices at which
// a position falls short, and sends each alert by the channel its holder chose.
package alerts

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/internal/band"
	"github.com/tapehouse/tapehouse/services/internal/margin"
)

var (
	priceUnit     = big.NewInt(100_000_000)
	basisPoints   = big.NewInt(10_000)
	usdgToUSD     = big.NewInt(1_000_000_000_000)
	wethLiquidate = big.NewInt(8_400)
)

// Stock is a Stock Token a position holds: its quantity as the liquidator counts it, collateral plus lent plus what
// its baskets hold, and the part of that on loan.
type Stock struct {
	Asset    string
	Quantity *big.Int
	Lent     *big.Int
}

// Position is what the liquidator reads of one position at one block: what it holds, in token units, and what it owes
// in USDG.
type Position struct {
	Account common.Address
	ID      [32]byte
	Stocks  []Stock
	USDG    *big.Int
	WETH    *big.Int
	Debt    *big.Int
	Premium *big.Int
}

func (p *Position) holds(asset string) bool {
	for _, s := range p.Stocks {
		if s.Asset == asset && s.Quantity.Sign() > 0 {
			return true
		}
	}
	return false
}

// Asset is what the band says of one asset at one block, and what rebuilds its band at another price: the variance
// and basis of its 24/7 leg.
type Asset struct {
	Name        string
	State       uint8
	Unconfirmed bool
	Mid         uint64
	Low         uint64
	High        *big.Int
	Variance    *big.Int
	Basis       uint64
}

// Market is everything a valuation reads of the chain at one block: the band's session and assets, WETH's price, and
// the margin engine's parameters and pools.
type Market struct {
	Block         uint64
	TimeMs        uint64
	Settled       bool
	Closed        bool
	BoundaryMs    uint64
	Regime        margin.Regime
	Assets        []Asset
	EthPrice      *big.Int
	RecallHaircut *big.Int
	Set           *margin.Set
	Depths        []uint32
	Pools         []margin.Pool
}

// Edges are each asset's low and high edge, in the engine's order of assets.
type Edges struct {
	Low, High []*big.Int
}

// Valuation is a position as the liquidator values it: whether it can be judged, whether it holds a Stock Token, and
// its equity and requirement in USD with 18 decimals.
type Valuation struct {
	Judged      bool
	Held        bool
	Equity      *big.Int
	Requirement *big.Int
}

// Short reports whether the position falls short.
func (v Valuation) Short() bool {
	return v.Judged && v.Equity.Cmp(v.Requirement) < 0
}

// Surplus is the equity over the requirement.
func (v Valuation) Surplus() *big.Int {
	return new(big.Int).Sub(v.Equity, v.Requirement)
}

// Current is every asset's band at the block.
func (m *Market) Current() Edges {
	e := Edges{Low: make([]*big.Int, len(m.Assets)), High: make([]*big.Int, len(m.Assets))}
	for i, a := range m.Assets {
		e.Low[i], e.High[i] = new(big.Int).SetUint64(a.Low), new(big.Int).Set(a.High)
	}
	return e
}

func (m *Market) index(asset string) int {
	for i, a := range m.Assets {
		if a.Name == asset {
			return i
		}
	}
	return -1
}

// Value is the liquidator's valuation of p with each Stock Token at its low edge, or at its high edge where high.
// Nothing is judged while the sequencer is unsettled, an asset held is halted or has an unconfirmed multiplier step,
// or WETH is held without a price.
func (m *Market) Value(p *Position, e Edges, regime margin.Regime, high bool) Valuation {
	var v Valuation
	if !m.Settled {
		return v
	}
	quantities, prices := make([]*big.Int, len(m.Assets)), make([]*big.Int, len(m.Assets))
	for i := range quantities {
		quantities[i], prices[i] = new(big.Int), new(big.Int)
	}
	gross, haircut := new(big.Int), new(big.Int)
	for _, s := range p.Stocks {
		if s.Quantity.Sign() == 0 {
			continue
		}
		i := m.index(s.Asset)
		if i < 0 || m.Assets[i].State == uint8(band.Halted) || m.Assets[i].Unconfirmed {
			return v
		}
		unit := e.Low[i]
		if high {
			unit = e.High[i]
		}
		quantities[i], prices[i], v.Held = s.Quantity, unit, true
		gross.Add(gross, new(big.Int).Div(new(big.Int).Mul(s.Quantity, unit), priceUnit))
		cut := new(big.Int).Mul(new(big.Int).Mul(s.Lent, unit), m.RecallHaircut)
		haircut.Add(haircut, cut.Div(cut, new(big.Int).Mul(priceUnit, basisPoints)))
	}
	gross.Sub(gross, haircut)
	cash := new(big.Int).Mul(p.USDG, usdgToUSD)
	eth := new(big.Int)
	if p.WETH.Sign() != 0 {
		if m.EthPrice.Sign() == 0 {
			return v
		}
		eth.Mul(p.WETH, m.EthPrice).Div(eth, priceUnit)
	}
	owed := new(big.Int).Mul(new(big.Int).Add(p.Debt, p.Premium), usdgToUSD)
	v.Equity = gross.Add(gross, cash).Add(gross, eth.Mul(eth, wethLiquidate).Div(eth, basisPoints))
	v.Equity.Sub(v.Equity, owed)
	v.Requirement = new(big.Int)
	if v.Held {
		if regime.Kind == margin.Unknown {
			regime = margin.Regime{Kind: margin.Open}
		}
		requirement, ok := m.Requirement(quantities, prices, regime)
		if !ok {
			return Valuation{}
		}
		v.Requirement = requirement
	}
	v.Judged = true
	return v
}

// Requirement is what the engine's currentRequirement answers for a portfolio of quantities with 18 decimals valued at
// prices with 8, in the regime the band's session puts it in; false where an exposure is beyond what the engine takes.
func (m *Market) Requirement(quantities, prices []*big.Int, regime margin.Regime) (*big.Int, bool) {
	exposures, large := margin.Exposures(quantities, prices)
	if large >= 0 {
		return nil, false
	}
	open, closed, floor, _ := margin.Requirements(m.Set, exposures, prices, m.Depths, m.Pools)
	return margin.Current(open, closed, floor, regime), true
}

// Portfolio is p's Stock Tokens as the engine takes them, quantities and prices in the engine's order of assets, each
// at its low edge in e; false where an asset is one the engine does not list.
func (m *Market) Portfolio(p *Position, e Edges) (quantities, prices []*big.Int, ok bool) {
	quantities, prices = make([]*big.Int, len(m.Assets)), make([]*big.Int, len(m.Assets))
	for i := range quantities {
		quantities[i], prices[i] = new(big.Int), new(big.Int)
	}
	for _, s := range p.Stocks {
		if s.Quantity.Sign() == 0 {
			continue
		}
		i := m.index(s.Asset)
		if i < 0 {
			return nil, nil, false
		}
		quantities[i], prices[i] = s.Quantity, e.Low[i]
	}
	return quantities, prices, true
}

// Assess is the liquidator's judgement of p: while the market is open the low edge decides, and while it is closed
// the position is short only if it is short at both edges, the better surplus of the two being its equity.
func (m *Market) Assess(p *Position, e Edges, closed bool, regime margin.Regime) Valuation {
	v := m.Value(p, e, regime, false)
	if !closed || !v.Short() {
		return v
	}
	if high := m.Value(p, e, regime, true); high.Surplus().Cmp(v.Surplus()) > 0 {
		return high
	}
	return v
}

type bandKind func(a Asset, centre uint64) band.Quote

// closedBand is an asset's band while the 24/5 session is closed: its 24/7 leg alone at centre.
func closedBand(a Asset, centre uint64) band.Quote {
	return band.Compute(band.Inputs{Live247Px: centre, VarCpb2: a.Variance, BasisBps: a.Basis})
}

// openBand is an asset's band while the session is open, both legs at centre.
func openBand(a Asset, centre uint64) band.Quote {
	return band.Compute(band.Inputs{Live247Px: centre, ClPx: centre, ClSessionOpen: true, VarCpb2: a.Variance, BasisBps: a.Basis})
}

func (m *Market) edge(i int, kind bandKind, centre uint64) (low, high *big.Int) {
	if centre == 0 {
		return new(big.Int), new(big.Int)
	}
	q := kind(m.Assets[i], centre)
	return new(big.Int).SetUint64(q.Low), q.High
}

func (m *Market) centres(kind bandKind) Edges {
	e := Edges{Low: make([]*big.Int, len(m.Assets)), High: make([]*big.Int, len(m.Assets))}
	for i, a := range m.Assets {
		e.Low[i], e.High[i] = m.edge(i, kind, a.Mid)
	}
	return e
}

// Status is what a price search found.
type Status uint8

// The outcomes of a price search.
const (
	Priced Status = iota
	NoPrice
	NotJudged
)

// Price is a searched price: the highest centre of the asset's band, in USD with 8 decimals, at which the position
// falls short.
type Price struct {
	Centre uint64
	Status Status
}

// WeekendPrice is the highest centre price of asset at which p falls short at both edges of its band while the
// market is closed, against the closed requirement: the band is the 24/7 leg alone and every other asset stays at its
// centre. A position short at the current centre reports it.
func (m *Market) WeekendPrice(p *Position, asset string) Price {
	return m.search(p, asset, closedBand, margin.Regime{Kind: margin.Closed}, true)
}

// ReopeningPrice is the highest centre price of asset at which p falls short at its low edge against the open
// requirement, with both legs of the band at that centre: the price at which the reopening auction may enrol it.
func (m *Market) ReopeningPrice(p *Position, asset string) Price {
	return m.search(p, asset, openBand, margin.Regime{Kind: margin.Open}, false)
}

func (m *Market) search(p *Position, asset string, kind bandKind, regime margin.Regime, closed bool) Price {
	i := m.index(asset)
	if !m.Value(p, m.Current(), regime, false).Judged {
		return Price{Status: NotJudged}
	}
	if i < 0 || !p.holds(asset) {
		return Price{Status: NoPrice}
	}
	edges := m.centres(kind)
	short := func(centre uint64) bool {
		edges.Low[i], edges.High[i] = m.edge(i, kind, centre)
		return m.Assess(p, edges, closed, regime).Short()
	}
	high := m.Assets[i].Mid
	if short(high) {
		return Price{Centre: high}
	}
	if !short(0) {
		return Price{Status: NoPrice}
	}
	low := uint64(0)
	for high-low > 1 {
		if mid := low + (high-low)/2; short(mid) {
			low = mid
		} else {
			high = mid
		}
	}
	return Price{Centre: low}
}
