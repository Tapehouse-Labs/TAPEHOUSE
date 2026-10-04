// SPDX-License-Identifier: MIT OR Apache-2.0

package alerts

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tapehouse/tapehouse/services/internal/margin"
)

const reopeningNote = "From the 24/5 reopen, a position that falls short at the low edges of its bands may be enrolled in the " +
	"reopening auction and sold at the regular open at the clearing price, at least 90% of the sealed low edge."

// ratio shows equity over requirement with three decimals, from integers.
func ratio(equity, requirement *big.Int) string {
	scaled := new(big.Int).Mul(equity, big.NewInt(1000))
	scaled.Quo(scaled, requirement)
	sign := ""
	if scaled.Sign() < 0 {
		sign, scaled = "-", scaled.Neg(scaled)
	}
	whole, part := new(big.Int).QuoRem(scaled, big.NewInt(1000), new(big.Int))
	return fmt.Sprintf("%s%s.%03d", sign, whole, part.Int64())
}

// usd shows an amount of USD with decimals decimals, to the cent and rounded toward zero.
func usd(amount *big.Int, decimals int) string {
	cents := new(big.Int).Mul(amount, big.NewInt(100))
	cents.Quo(cents, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
	sign := ""
	if cents.Sign() < 0 {
		sign, cents = "-", cents.Neg(cents)
	}
	whole, part := new(big.Int).QuoRem(cents, big.NewInt(100), new(big.Int))
	return fmt.Sprintf("%s%s.%02d USD", sign, whole, part.Int64())
}

func price(value uint64) string {
	return usd(new(big.Int).SetUint64(value), 8)
}

func label(key string) string {
	if strings.HasPrefix(key, "short:") {
		return "your short of " + strings.TrimPrefix(key, "short:")
	}
	return "your position " + key[:min(len(key), 10)]
}

// thresholdText is the alert of a ratio under its subscriber's threshold.
func thresholdText(key string, threshold int, equity, requirement *big.Int) Message {
	return Message{
		Subject: "Tapehouse: " + label(key) + " is close to its requirement",
		Body: fmt.Sprintf("The equity of %s is %s against a requirement of %s: a ratio of %s, under your threshold of %d.%03d.\n"+
			"Below 1.000 the position falls short and can be liquidated.\n", label(key), usd(equity, 18), usd(requirement, 18),
			ratio(equity, requirement), threshold/1000, threshold%1000),
	}
}

// liquidationText is the alert of a position the liquidator would start now, or has started.
func liquidationText(key string, equity, requirement *big.Int, auction uint64) Message {
	body := fmt.Sprintf("%s falls short: equity %s against a requirement of %s. Its collateral can be sold in a liquidation auction.\n",
		label(key), usd(equity, 18), usd(requirement, 18))
	if auction != 0 {
		body += "An auction of it started at " + time.Unix(int64(auction), 0).UTC().Format("15:04:05 UTC on 2 Jan") + ".\n"
	}
	return Message{Subject: "Tapehouse: " + label(key) + " falls short", Body: body}
}

// unjudgedText is the notice of a position nobody can judge.
func unjudgedText(key string) Message {
	return Message{
		Subject: "Tapehouse: " + label(key) + " cannot be judged",
		Body: "An asset it holds is halted or has a multiplier change not yet confirmed, or its price feed is not live. " +
			"No price or alert can be given for it until that clears.\n",
	}
}

// assetPrices are what a weekend alert says of one Stock Token a position holds.
type assetPrices struct {
	Asset     string
	Centre    uint64
	Weekend   Price
	Reopening Price
	Open      *big.Int
}

// weekendText is the alert before a closure: when it closes, the closed requirement against the equity now, and each
// asset's weekend and reopening prices.
func weekendText(key string, closure uint64, equity, closed *big.Int, assets []assetPrices) Message {
	var body strings.Builder
	fmt.Fprintf(&body, "The market closes at %s.\nThe equity of %s is %s now; its requirement across the closure is %s.\n\n",
		time.UnixMilli(int64(closure)).UTC().Format("15:04 UTC on Mon 2 Jan"), label(key), usd(equity, 18), usd(closed, 18))
	for _, a := range assets {
		fmt.Fprintf(&body, "%s, now %s:\n", a.Asset, price(a.Centre))
		fmt.Fprintf(&body, "  Weekend liquidation price: %s\n", searched(a.Weekend))
		fmt.Fprintf(&body, "  Reopening price: %s\n", searched(a.Reopening))
		if a.Open != nil && a.Open.Sign() > 0 {
			fmt.Fprintf(&body, "  Open-market liquidation price: %s\n", usd(a.Open, 8))
		}
	}
	body.WriteString("\nThe weekend price is the price of the asset, with the others where they are, at which the position falls short at both " +
		"edges of its band while the market is closed. The reopening price is the one at which it falls short at its low edge against the open " +
		"requirement.\n" + reopeningNote + "\n")
	return Message{Subject: "Tapehouse: the market is about to close; " + label(key), Body: body.String()}
}

func searched(p Price) string {
	switch p.Status {
	case Priced:
		return price(p.Centre)
	case NoPrice:
		return "none: no price makes the position fall short"
	}
	return "none: the position cannot be judged"
}

func closedRegime() margin.Regime {
	return margin.Regime{Kind: margin.Closed}
}

func symbolName(hex string) (string, error) {
	raw, err := hexutil.Decode(hex)
	if err != nil || len(raw) != 32 {
		return "", fmt.Errorf("%q is not a 32-byte word", hex)
	}
	return strings.TrimRight(string(raw), "\x00"), nil
}
