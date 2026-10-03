// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"maps"
	"math/big"
	"slices"
)

// RebalanceSlippageBps is the slippage the basket keeper allows its pool quotes when it judges whether a rebalance
// pays for itself.
const RebalanceSlippageBps = 50

// rebalance moves each basket toward the target in effect from the keeper's own Stock Tokens: it gives what the
// basket holds too few of, up to what the keeper holds, and takes what it holds too many of, no more than what comes
// in is worth at the bands' low edges against their high edges. It rebalances only when the pools pay for the trade,
// what goes out selling for at least what comes in costs, or when the team pays the bands' width
// (TAPEHOUSE_KEEPER_REBALANCE=subsidise). A halted band, a pending multiplier change, an unsettled sequencer or a
// paused token waits.
func (k *Keeper) rebalance(ctx context.Context) error {
	for _, key := range slices.Sorted(maps.Keys(k.client.Deployments().Baskets)) {
		if err := k.rebalanceBasket(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

func (k *Keeper) rebalanceBasket(ctx context.Context, key string) error {
	basket := k.client.Basket(key)
	address := k.client.Deployments().Baskets[key]
	components, err := basket.Components(k.opts(ctx))
	if err != nil {
		return err
	}
	units, err := basket.Target(k.opts(ctx))
	if err != nil {
		return err
	}
	supply, err := basket.TotalSupply(k.opts(ctx))
	if err != nil {
		return err
	}
	n := len(components.Tokens)
	in, out := make([]*big.Int, n), make([]*big.Int, n)
	valueIn, valueOut := new(big.Int), new(big.Int)
	lows, highs := make([]*big.Int, n), make([]*big.Int, n)
	moves := false
	for i, token := range components.Tokens {
		in[i], out[i] = new(big.Int), new(big.Int)
		held, err := k.client.BalanceOf(k.opts(ctx), token, address)
		if err != nil {
			return err
		}
		goal := new(big.Int).Mul(units[i], supply)
		have := new(big.Int).Mul(held, wad)
		switch have.Cmp(goal) {
		case -1:
			in[i].Div(goal.Sub(goal, have), wad)
			own, err := k.client.BalanceOf(k.opts(ctx), token, k.sender.From())
			if err != nil {
				return err
			}
			if in[i].Cmp(own) > 0 {
				in[i].Set(own)
			}
		case 1:
			out[i].Div(have.Sub(have, goal), wad)
		}
		if in[i].Sign() == 0 && out[i].Sign() == 0 {
			continue
		}
		quote, err := k.client.Band().Quote(k.opts(ctx), components.Assets[i])
		if err != nil {
			return err
		}
		lows[i], highs[i] = new(big.Int).SetUint64(quote.Low), quote.High
		valueIn.Add(valueIn, new(big.Int).Mul(in[i], lows[i]))
		valueOut.Add(valueOut, new(big.Int).Mul(out[i], highs[i]))
		moves = true
	}
	if !moves || valueIn.Sign() == 0 {
		return nil
	}
	if valueOut.Cmp(valueIn) > 0 {
		for i := range out {
			out[i].Div(out[i].Mul(out[i], valueIn), valueOut)
		}
	}
	if !k.cfg.Rebalance && !k.pays(ctx, components.Assets, in, out) {
		k.log.Info("the basket is off its target; the pools do not pay the bands' width", "basket", key)
		return nil
	}
	for i, token := range components.Tokens {
		if in[i].Sign() == 0 {
			continue
		}
		approve, err := k.client.Approve(token, address, in[i])
		if _, err := k.act(ctx, "approve "+components.Assets[i], approve, err); err != nil {
			return err
		}
	}
	tx, err := basket.Rebalance(in, out, k.sender.From())
	_, err = k.act(ctx, "rebalance "+key, tx, err, "AssetHalted", "CorporateActionPending", "SequencerNotSettled",
		"ComponentPaused", "ValueLost", "PastTarget")
	return err
}

// pays reports whether what goes out sells through the pools for at least what comes in costs through them, by
// QuoterV2's quotes less RebalanceSlippageBps. A token without a pool to quote does not pay.
func (k *Keeper) pays(ctx context.Context, assets []string, in, out []*big.Int) bool {
	proceeds, cost := new(big.Int), new(big.Int)
	shorts := k.client.Shorts()
	for i, asset := range assets {
		if out[i].Sign() > 0 {
			sale, err := shorts.QuoteSale(k.opts(ctx), asset, out[i], RebalanceSlippageBps)
			if err != nil {
				return false
			}
			proceeds.Add(proceeds, sale.MinProceeds)
		}
		if in[i].Sign() > 0 {
			buy, err := shorts.QuoteCover(k.opts(ctx), asset, in[i], RebalanceSlippageBps)
			if err != nil {
				return false
			}
			cost.Add(cost, buy.MaxCost)
		}
	}
	return proceeds.Cmp(cost) >= 0
}
