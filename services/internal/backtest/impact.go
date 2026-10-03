// SPDX-License-Identifier: MIT OR Apache-2.0

package backtest

import (
	"math/big"

	"github.com/tapehouse/tapehouse/services/internal/margin"
)

// ImpactShares are the sales the engine's linear impact is checked on, as shares of the asset's selling depth.
var ImpactShares = []float64{0.1, 0.5, 1}

// Impact is a sale of an asset worth a share of its selling depth: what the engine's liquidity add-on charges for it
// and what its pool's quote at a later block costs, each as a fraction of the value sold.
type Impact struct {
	Asset  string  `json:"asset"`
	Block  uint64  `json:"block"`
	Share  float64 `json:"share"`
	Engine float64 `json:"engine"`
	Pool   float64 `json:"pool"`
}

func wad(usd float64) *big.Int {
	v, _ := new(big.Float).Mul(big.NewFloat(usd), big.NewFloat(1e18)).Int(nil)
	return v
}

// EngineCost is the engine's liquidity add-on for selling value USD of an asset whose pool trades at its valuation
// price, with a selling depth of depth USD and a fee in millionths, as a fraction of value.
func EngineCost(value, depth float64, fee uint32) float64 {
	price := big.NewInt(100_000_000)
	addon := margin.LiquidityAddon(wad(value), true, price, price, wad(depth), fee, 0)
	f, _ := new(big.Float).Quo(new(big.Float).SetInt(addon), new(big.Float).SetInt(wad(value))).Float64()
	return f
}
