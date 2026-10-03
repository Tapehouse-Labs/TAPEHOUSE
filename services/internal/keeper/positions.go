// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	wad         = big.NewInt(1e18)
	tokenToUSDG = new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil)
)

// holding is a token a position holds, its symbol where it is a Stock Token, and how much.
type holding struct {
	asset  string
	token  common.Address
	amount *big.Int
}

// holdings reads what account's position holds of USDG, WETH and each Stock Token the accounts take.
func (k *Keeper) holdings(ctx context.Context, account common.Address, position [32]byte) ([]holding, error) {
	accounts := k.client.Accounts()
	stocks, err := accounts.Stocks(k.opts(ctx))
	if err != nil {
		return nil, err
	}
	tokens := k.client.Deployments().Tokens
	all := []holding{{"USDG", tokens["USDG"], nil}, {"WETH", tokens["WETH"], nil}}
	for i, asset := range stocks.Assets {
		all = append(all, holding{asset, stocks.Tokens[i], nil})
	}
	var held []holding
	for _, h := range all {
		if h.token == (common.Address{}) {
			continue
		}
		if h.amount, err = accounts.Collateral(k.opts(ctx), account, position, h.token); err != nil {
			return nil, err
		}
		if h.amount.Sign() > 0 {
			held = append(held, h)
		}
	}
	return held, nil
}

// liquidate works every position with debt, deepest first: it starts the auction of each that falls short and has
// none running in the market's current state, unless the reopening auction holds it while the market is open; recalls
// what it lent that its vault cannot return now; settles it in cash where a Stock Token it holds is frozen; writes off
// one left with nothing while the accounts have no backstop; buys at its auction as the buyer of last resort where
// configured; and stops the auction of each that recovered.
func (k *Keeper) liquidate(ctx context.Context) error {
	positions, err := k.index.Debts(ctx)
	if err != nil {
		return err
	}
	deepestFirst(positions, func(p Position) (*big.Int, *big.Int) { return p.Equity, p.Requirement })
	var failures []error
	for _, p := range positions {
		if err := k.liquidatePosition(ctx, p.Account, p.Position); err != nil {
			failures = append(failures, err)
		}
	}
	return errors.Join(failures...)
}

func (k *Keeper) liquidatePosition(ctx context.Context, account common.Address, position [32]byte) error {
	liquidator := k.client.Liquidator()
	judged, err := liquidator.Shortfall(k.opts(ctx), account, position)
	if err != nil {
		return err
	}
	running, err := liquidator.Running(k.opts(ctx), account, position)
	if err != nil {
		return err
	}
	name := fmt.Sprintf("%s %s", account.Hex(), hexutil.Encode(position[:]))
	if !judged.Short {
		if !running {
			return nil
		}
		tx, err := liquidator.Stop(account, position)
		_, err = k.act(ctx, "stop "+name, tx, err, "StillShort", "CannotJudge")
		return err
	}
	held, err := liquidator.HeldUntil(k.opts(ctx), account, position)
	if err != nil {
		return err
	}
	if !judged.Closed && held*1000 > k.nowMs() {
		return nil
	}
	if err := k.recall(ctx, account, position); err != nil {
		return err
	}
	if written, err := k.writeOff(ctx, account, position, name); err != nil || written {
		return err
	}
	holdings, err := k.holdings(ctx, account, position)
	if err != nil {
		return err
	}
	if !running {
		tx, err := liquidator.Start(account, position)
		if _, err := k.act(ctx, "start "+name, tx, err, "NotLiquidatable", "PositionHeld"); err != nil {
			return err
		}
	}
	if err := k.settleCash(ctx, account, position, holdings, name); err != nil {
		return err
	}
	if k.cfg.Buy {
		return k.buy(ctx, account, position, holdings, name)
	}
	return nil
}

// recall recalls, for a position that falls short, what it lent of each Stock Token beyond what its vault can return
// now, so the auction can sell it once it comes back.
func (k *Keeper) recall(ctx context.Context, account common.Address, position [32]byte) error {
	accounts := k.client.Accounts()
	stocks, err := accounts.Stocks(k.opts(ctx))
	if err != nil {
		return err
	}
	for i, asset := range stocks.Assets {
		token := stocks.Tokens[i]
		if token == (common.Address{}) {
			continue
		}
		lent, err := accounts.Lent(k.opts(ctx), account, position, token)
		if err != nil || lent.Sign() == 0 {
			if err != nil {
				return err
			}
			continue
		}
		sellable, err := accounts.Sellable(k.opts(ctx), account, position, token)
		if err != nil {
			return err
		}
		held, err := accounts.Collateral(k.opts(ctx), account, position, token)
		if err != nil {
			return err
		}
		if lent.Cmp(new(big.Int).Sub(sellable, held)) <= 0 {
			continue
		}
		tx, err := k.client.Liquidator().Recall(account, position, token)
		if _, err := k.act(ctx, "recall "+asset, tx, err, "NothingToRecall", "NotLiquidatable"); err != nil {
			return err
		}
	}
	return nil
}

// settleCash repays a short position from its USDG where a Stock Token it holds is frozen and cannot be auctioned.
func (k *Keeper) settleCash(ctx context.Context, account common.Address, position [32]byte, holdings []holding, name string) error {
	cash, frozen := false, false
	for _, h := range holdings {
		if h.asset == "USDG" {
			cash = true
			continue
		}
		if h.asset == "WETH" {
			continue
		}
		paused, err := k.client.Paused(k.opts(ctx), h.token)
		if err != nil {
			return err
		}
		frozen = frozen || paused
	}
	if !cash || !frozen {
		return nil
	}
	tx, err := k.client.Liquidator().SettleCash(account, position)
	_, err = k.act(ctx, "settleCash "+name, tx, err, "NotLiquidatable", "NothingToBuy")
	return err
}

// writeOff writes off a position left with nothing, or with holdings worth less than the liquidator's dust, while the
// accounts have no backstop; once they have one, the backstop keeper covers it. PositionNotEmpty is a position that
// still holds something to sell.
func (k *Keeper) writeOff(ctx context.Context, account common.Address, position [32]byte, name string) (bool, error) {
	backstop, err := k.client.Accounts().HasBackstop(k.opts(ctx))
	if err != nil || backstop {
		return false, err
	}
	tx, err := k.client.Liquidator().WriteOff(account, position)
	return k.act(ctx, "writeOff "+name, tx, err, "PositionNotEmpty")
}

// buy buys each token a short position holds at its auction's ask once the ask is no more than the token's reference,
// its band's low edge or ETH's price, as the buyer of last resort, within the keeper's USDG.
func (k *Keeper) buy(ctx context.Context, account common.Address, position [32]byte, holdings []holding, name string) error {
	usdg := k.client.Deployments().Tokens["USDG"]
	for _, h := range holdings {
		if h.asset == "USDG" {
			continue
		}
		ask, err := k.client.Liquidator().Price(k.opts(ctx), account, position, h.token)
		if err != nil || ask.Sign() == 0 {
			return err
		}
		reference, err := k.reference(ctx, h.asset)
		if err != nil || reference.Sign() == 0 || ask.Cmp(reference) > 0 {
			return err
		}
		budget, err := k.client.BalanceOf(k.opts(ctx), usdg, k.sender.From())
		if err != nil {
			return err
		}
		amount := new(big.Int).Div(new(big.Int).Mul(budget, tokenToUSDG), ask)
		if amount.Cmp(h.amount) > 0 {
			amount = h.amount
		}
		if amount.Sign() == 0 {
			return nil
		}
		approve, err := k.client.Approve(usdg, k.client.Deployments().Tapehouse["Liquidator"], budget)
		if _, err := k.act(ctx, "approve USDG", approve, err); err != nil {
			return err
		}
		tx, err := k.client.Liquidator().Buy(account, position, h.token, amount, budget, k.sender.From())
		if _, err := k.act(ctx, "buy "+h.asset+" "+name, tx, err, "NothingToBuy", "NotLiquidatable", "NoAuction",
			"PositionHeld"); err != nil {
			return err
		}
	}
	return nil
}

// reference is what the keeper pays at most for a whole token, in USD with 8 decimals: a Stock Token's band low edge,
// WETH's Chainlink price.
func (k *Keeper) reference(ctx context.Context, asset string) (*big.Int, error) {
	if asset == "WETH" {
		feed, ok := k.client.Deployments().Chainlink["ETH_USD"]
		if !ok {
			return new(big.Int), nil
		}
		round, err := k.client.LatestRound(k.opts(ctx), feed)
		return round.Answer, err
	}
	quote, err := k.client.Band().Quote(k.opts(ctx), asset)
	return new(big.Int).SetUint64(quote.Low), err
}

// backstop claims the backstop's premium whenever the accounts hold some for it, and covers each emptied position
// with debt in place of the liquidator's write-off. A cover that reverts otherwise than PositionNotEmpty waits for a
// halted asset or a stale ETH price, or for the recall and sale of what the position lent.
func (k *Keeper) backstop(ctx context.Context) error {
	premium, err := k.client.Accounts().BackstopPremium(k.opts(ctx))
	if err != nil {
		return err
	}
	if premium.Sign() > 0 {
		tx, err := k.client.Backstop().Claim()
		if _, err := k.act(ctx, "claim the backstop's premium", tx, err); err != nil {
			return err
		}
	}
	positions, err := k.index.Debts(ctx)
	if err != nil {
		return err
	}
	for _, p := range positions {
		tx, err := k.client.Backstop().Cover(p.Account, p.Position)
		if _, err := k.act(ctx, fmt.Sprintf("cover %s %s", p.Account.Hex(), hexutil.Encode(p.Position[:])), tx, err,
			"PositionNotEmpty"); err != nil {
			k.log.Warn("the cover waits", "error", err)
		}
	}
	return nil
}

// shorts liquidates every short whose equity is below its requirement, those in deficit first, since until they are
// closed other shorts' closes revert with DeficitOpen, then the deepest. NotShortfall is a short that recovered;
// AboveLimit a pool asking too much above the band's centre, which the next pass tries again.
func (k *Keeper) shorts(ctx context.Context) error {
	open, err := k.index.Shorts(ctx)
	if err != nil {
		return err
	}
	deepestFirst(open, func(s Short) (*big.Int, *big.Int) { return s.Equity, s.Requirement })
	var deficits, short []Short
	for _, s := range open {
		switch {
		case s.USDG.Sign() < 0:
			deficits = append(deficits, s)
		case s.Equity.Cmp(s.Requirement) < 0:
			short = append(short, s)
		}
	}
	for _, s := range append(deficits, short...) {
		tx, err := k.client.Shorts().Liquidate(s.Account, s.Asset)
		if _, err := k.act(ctx, fmt.Sprintf("liquidate %s's short of %s", s.Account.Hex(), s.Asset), tx, err,
			"NotShortfall", "AboveLimit", "DeficitOpen"); err != nil {
			return err
		}
	}
	return nil
}

// syncHoldings syncs each Stock Token the accounts or its lending vault count more of than they hold, as after the
// issuer burns from them, and clears each position's holding a burn left worth less than a raw unit.
func (k *Keeper) syncHoldings(ctx context.Context) error {
	accounts := k.client.Accounts()
	stocks, err := accounts.Stocks(k.opts(ctx))
	if err != nil {
		return err
	}
	address := k.client.Deployments().Tapehouse["MarginAccounts"]
	for i, asset := range stocks.Assets {
		token := stocks.Tokens[i]
		if token == (common.Address{}) {
			continue
		}
		book, err := accounts.Holding(k.opts(ctx), asset)
		if err != nil {
			return err
		}
		balance, err := k.client.BalanceOf(k.opts(ctx), token, address)
		if err != nil {
			return err
		}
		counted := new(big.Int).Div(new(big.Int).Mul(book.Units, book.Scale), wad)
		burnt := balance.Cmp(counted) < 0
		if vault, ok := k.client.Deployments().StockLending[asset]; ok && !burnt {
			if burnt, err = k.vaultBurnt(ctx, asset, token, vault); err != nil {
				return err
			}
		}
		if burnt {
			tx, err := accounts.Sync(asset)
			if _, err := k.act(ctx, "sync "+asset, tx, err); err != nil {
				return err
			}
			if book, err = accounts.Holding(k.opts(ctx), asset); err != nil {
				return err
			}
		}
		if book.Scale.Cmp(wad) < 0 && book.Units.Sign() > 0 {
			if err := k.clearDust(ctx, asset, token, book.Scale); err != nil {
				return err
			}
		}
	}
	return nil
}

func (k *Keeper) vaultBurnt(ctx context.Context, asset string, token, vault common.Address) (bool, error) {
	counted, err := k.client.LendingVault(asset).Counted(k.opts(ctx))
	if err != nil {
		return false, err
	}
	balance, err := k.client.BalanceOf(k.opts(ctx), token, vault)
	return balance.Cmp(counted) < 0, err
}

// clearDust clears, once at each scale, the holding of asset of every position that deposited it and now holds none
// of it as a whole raw unit.
func (k *Keeper) clearDust(ctx context.Context, asset string, token common.Address, scale *big.Int) error {
	deposits, err := k.index.Events(ctx, "tapehouse.MarginAccounts", "Deposit", map[string]string{"token": token.Hex()})
	if err != nil {
		return err
	}
	for _, deposit := range deposits {
		account := common.HexToAddress(deposit.Arg("account"))
		position, err := word(deposit.Arg("position"))
		if err != nil {
			return err
		}
		key := fmt.Sprintf("clear %s %x %s", account.Hex(), position, asset)
		held, err := k.client.Accounts().Collateral(k.opts(ctx), account, position, token)
		if err != nil {
			return err
		}
		if held.Sign() != 0 || k.remember(key, scale.String()) {
			continue
		}
		tx, err := k.client.Accounts().Clear(account, position, asset)
		if _, err := k.act(ctx, "clear "+asset, tx, err, "HoldingNotEmpty"); err != nil {
			k.remember(key, "")
			return err
		}
	}
	return nil
}
