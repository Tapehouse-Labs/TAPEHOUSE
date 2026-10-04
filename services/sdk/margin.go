// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/margin"
)

// Margin reads the margin engine's risk parameters and packs their update.
type Margin struct {
	target
	c      *Client
	engine *margin.Margin
}

// Margin returns the margin engine of the registry's .tapehouse.Margin.
func (c *Client) Margin() *Margin {
	return &Margin{
		target: lookup(c.deployments.Tapehouse, "Margin", ".tapehouse"),
		c:      c,
		engine: margin.NewMargin(),
	}
}

// Address returns the engine's address.
func (m *Margin) Address() (common.Address, error) {
	return m.address, m.err
}

// Assets reads the engine's assets, in the order its parameters are stored.
func (m *Margin) Assets(opts *bind.CallOpts) ([][32]byte, error) {
	return read(m.c, opts, m.target, m.engine.UnpackAssets)(m.engine.TryPackAssets())
}

// Volatility reads the daily volatility of asset and its floor, in centi-basis-points.
func (m *Margin) Volatility(opts *bind.CallOpts, asset [32]byte) (value, floor uint32, err error) {
	out, err := read(m.c, opts, m.target, m.engine.UnpackVolatility)(m.engine.TryPackVolatility(asset))
	return out.Value, out.Floor, err
}

// Correlation reads the correlation of asset and other and its floor, in basis points.
func (m *Margin) Correlation(opts *bind.CallOpts, asset, other [32]byte) (value, floor uint16, err error) {
	out, err := read(m.c, opts, m.target, m.engine.UnpackCorrelation)(m.engine.TryPackCorrelation(asset, other))
	return out.Value, out.Floor, err
}

// WeekendGap reads the weekend gap of asset and its floor, in millionths.
func (m *Margin) WeekendGap(opts *bind.CallOpts, asset [32]byte) (value, floor uint32, err error) {
	out, err := read(m.c, opts, m.target, m.engine.UnpackWeekendGap)(m.engine.TryPackWeekendGap(asset))
	return out.Value, out.Floor, err
}

// Depth reads the USD a liquidation of asset can sell, then buy, within a 10% move of its pool, and their ceilings.
func (m *Margin) Depth(opts *bind.CallOpts, asset [32]byte) (margin.DepthOutput, error) {
	return read(m.c, opts, m.target, m.engine.UnpackDepth)(m.engine.TryPackDepth(asset))
}

// Pool reads the Uniswap v3 pool asset is liquidated in; zero for none.
func (m *Margin) Pool(opts *bind.CallOpts, asset [32]byte) (common.Address, error) {
	return read(m.c, opts, m.target, m.engine.UnpackPool)(m.engine.TryPackPool(asset))
}

// Market reads the asset that stands for the market; zero for the equal-weighted portfolio of the assets.
func (m *Margin) Market(opts *bind.CallOpts) ([32]byte, error) {
	return read(m.c, opts, m.target, m.engine.UnpackMarket)(m.engine.TryPackMarket())
}

// LastUpdate reads when the parameters were last updated, in seconds; zero before the first update.
func (m *Margin) LastUpdate(opts *bind.CallOpts) (uint64, error) {
	return read(m.c, opts, m.target, m.engine.UnpackLastUpdate)(m.engine.TryPackLastUpdate())
}

// Owner reads the engine's owner, who may update the parameters.
func (m *Margin) Owner(opts *bind.CallOpts) (common.Address, error) {
	return read(m.c, opts, m.target, m.engine.UnpackOwner)(m.engine.TryPackOwner())
}

// SetParameters packs the replacement of every parameter, in the order of Assets: the volatilities and gaps, the
// correlations of the upper triangle row by row, and each asset's selling then buying depth.
func (m *Margin) SetParameters(volatilities []uint32, correlations []uint16, gaps []uint32, depths []uint32) (Tx, error) {
	return m.tx(m.engine.TryPackSetParameters(volatilities, correlations, gaps, depths))
}

// EthUsdFeed reads the Chainlink feed that prices WETH for the engine.
func (m *Margin) EthUsdFeed(opts *bind.CallOpts) (common.Address, error) {
	return read(m.c, opts, m.target, m.engine.UnpackEthUsdFeed)(m.engine.TryPackEthUsdFeed())
}

// CurrentRequirement reads the margin a portfolio needs now, in USD with 18 decimals, its missing bits and the regime
// the band's session puts it in. quantities are signed token amounts with 18 decimals and prices USD with 8 decimals,
// in the order of Assets.
func (m *Margin) CurrentRequirement(opts *bind.CallOpts, quantities, prices []*big.Int) (margin.CurrentRequirementOutput, error) {
	return read(m.c, opts, m.target, m.engine.UnpackCurrentRequirement)(m.engine.TryPackCurrentRequirement(quantities, prices))
}
