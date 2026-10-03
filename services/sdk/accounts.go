// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/marginaccounts"
)

// Accounts reads the margin accounts and packs their transactions.
type Accounts struct {
	target
	c        *Client
	accounts *marginaccounts.MarginAccounts
}

// Repayment is what repays all a position owes: its debt, its premium, and their sum, in USDG.
type Repayment struct {
	Debt    *big.Int
	Premium *big.Int
	Assets  *big.Int
}

// Accounts returns the margin accounts of the registry's .tapehouse.MarginAccounts.
func (c *Client) Accounts() *Accounts {
	return &Accounts{
		target:   lookup(c.deployments.Tapehouse, "MarginAccounts", ".tapehouse"),
		c:        c,
		accounts: marginaccounts.NewMarginAccounts(),
	}
}

// SetAuthorization lets authorized borrow and withdraw for the sender's account, deposit Stock Tokens and USDG into
// it, and act for its shorts, or stops it. Until then the margin accounts and the shorts revert with
// Unauthorized(caller, account).
func (a *Accounts) SetAuthorization(authorized common.Address, allowed bool) (Tx, error) {
	return a.tx(a.accounts.TryPackSetAuthorization(authorized, allowed))
}

// Deposit deposits amount of token from the sender into account's position: Cross or an asset's symbol. A basket's
// address deposits its shares, into Cross alone.
func (a *Accounts) Deposit(position [32]byte, token common.Address, amount *big.Int, account common.Address) (Tx, error) {
	if err := present(amount); err != nil {
		return Tx{}, err
	}
	return a.tx(a.accounts.TryPackDeposit(position, token, amount, account))
}

// Withdraw withdraws amount of token, a basket's shares included, from account's position to receiver.
func (a *Accounts) Withdraw(position [32]byte, token common.Address, amount *big.Int, account, receiver common.Address) (Tx, error) {
	if err := present(amount); err != nil {
		return Tx{}, err
	}
	return a.tx(a.accounts.TryPackWithdraw(position, token, amount, account, receiver))
}

// Borrow borrows assets of USDG against account's position, to receiver.
func (a *Accounts) Borrow(position [32]byte, assets *big.Int, account, receiver common.Address) (Tx, error) {
	if err := present(assets); err != nil {
		return Tx{}, err
	}
	return a.tx(a.accounts.TryPackBorrow(position, assets, account, receiver))
}

// Repay repays assets of USDG for account's position with the sender's USDG: its debt first, then its premium with
// what is left. Anyone may. Repayment reads what a full repayment takes.
func (a *Accounts) Repay(position [32]byte, assets *big.Int, account common.Address) (Tx, error) {
	if err := present(assets); err != nil {
		return Tx{}, err
	}
	return a.tx(a.accounts.TryPackRepay(position, assets, account))
}

// Unwrap redeems shares of the registry's basket key in account's cross position for its Stock Tokens, which the
// position then holds. The account or an address it authorized may, and anyone once the position falls short.
func (a *Accounts) Unwrap(account common.Address, key string, shares *big.Int) (Tx, error) {
	basket, err := entry(a.c.deployments.Baskets, key, ".tapehouse.Baskets")
	if err == nil {
		err = present(shares)
	}
	if err != nil {
		return Tx{}, err
	}
	return a.tx(a.accounts.TryPackUnwrap(account, basket, shares))
}

// IsAuthorized reports whether authorized may act for account.
func (a *Accounts) IsAuthorized(opts *bind.CallOpts, account, authorized common.Address) (bool, error) {
	return read(a.c, opts, a.target, a.accounts.UnpackIsAuthorized)(a.accounts.TryPackIsAuthorized(account, authorized))
}

// Repayment reads what repays all account's position owes: its debt and its premium, both at opts' block, or at the
// latest block where opts names none.
func (a *Accounts) Repayment(opts *bind.CallOpts, account common.Address, position [32]byte) (Repayment, error) {
	if a.err != nil {
		return Repayment{}, a.err
	}
	pinned, err := a.c.pin(opts)
	if err != nil {
		return Repayment{}, err
	}
	debt, err := read(a.c, pinned, a.target, a.accounts.UnpackDebt)(a.accounts.TryPackDebt(account, position))
	if err != nil {
		return Repayment{}, err
	}
	premium, err := read(a.c, pinned, a.target, a.accounts.UnpackPremium)(a.accounts.TryPackPremium(account, position))
	if err != nil {
		return Repayment{}, err
	}
	return Repayment{debt, premium, new(big.Int).Add(debt, premium)}, nil
}

// InBaskets reads what account's position holds of each Stock Token through its baskets, by asset, as its shares
// would redeem now: the engine margins them as those Stock Tokens. Both reads are at opts' block, or at the latest
// block where opts names none.
func (a *Accounts) InBaskets(opts *bind.CallOpts, account common.Address, position [32]byte) (map[string]*big.Int, error) {
	if a.err != nil {
		return nil, a.err
	}
	pinned, err := a.c.pin(opts)
	if err != nil {
		return nil, err
	}
	stocks, err := read(a.c, pinned, a.target, a.accounts.UnpackStocks)(a.accounts.TryPackStocks())
	if err != nil {
		return nil, err
	}
	amounts, err := read(a.c, pinned, a.target, a.accounts.UnpackInBaskets)(a.accounts.TryPackInBaskets(account, position))
	if err != nil {
		return nil, err
	}
	if len(amounts) != len(stocks.Symbols) {
		return nil, fmt.Errorf("the accounts name %d assets and %d amounts", len(stocks.Symbols), len(amounts))
	}
	held := make(map[string]*big.Int, len(amounts))
	for i, symbol := range stocks.Symbols {
		held[assetName(symbol)] = amounts[i]
	}
	return held, nil
}

// Health reads account's position as the engine sees it: equity and requirement in USD with 18 decimals, the
// engine's missing bits and the regime.
func (a *Accounts) Health(opts *bind.CallOpts, account common.Address, position [32]byte) (marginaccounts.HealthOutput, error) {
	return read(a.c, opts, a.target, a.accounts.UnpackHealth)(a.accounts.TryPackHealth(account, position))
}

// Collateral reads what account's position holds of token: USDG, WETH or a Stock Token.
func (a *Accounts) Collateral(opts *bind.CallOpts, account common.Address, position [32]byte, token common.Address) (*big.Int, error) {
	return read(a.c, opts, a.target, a.accounts.UnpackCollateral)(a.accounts.TryPackCollateral(account, position, token))
}

// Leverage reads the gross exposure of account's position over its equity, in basis points; the largest uint256 for a
// position without equity.
func (a *Accounts) Leverage(opts *bind.CallOpts, account common.Address, position [32]byte) (*big.Int, error) {
	return read(a.c, opts, a.target, a.accounts.UnpackLeverage)(a.accounts.TryPackLeverage(account, position))
}

// LiquidationPrice reads the price of asset, in USD with 8 decimals, at or above which account's position meets its
// requirement once it owes borrowing more USDG: its band's low edge where it falls short already, zero where no price
// leaves it short.
func (a *Accounts) LiquidationPrice(opts *bind.CallOpts, account common.Address, position [32]byte, asset string, borrowing *big.Int) (*big.Int, error) {
	if err := present(borrowing); err != nil {
		return nil, err
	}
	symbol, err := ToBytes32(asset)
	return read(a.c, opts, a.with(err), a.accounts.UnpackLiquidationPrice)(
		a.accounts.TryPackLiquidationPrice(account, position, symbol, borrowing))
}

// Stocks reads the Stock Tokens the accounts take: each asset's name and token, in the accounts' order.
func (a *Accounts) Stocks(opts *bind.CallOpts) (Components, error) {
	out, err := read(a.c, opts, a.target, a.accounts.UnpackStocks)(a.accounts.TryPackStocks())
	assets := make([]string, len(out.Symbols))
	for i, symbol := range out.Symbols {
		assets[i] = assetName(symbol)
	}
	return Components{assets, out.Tokens}, err
}

// Holding reads the accounts' holding of asset: the units the positions hold, each worth scale / 10^18 of a token, and
// the most they may hold.
func (a *Accounts) Holding(opts *bind.CallOpts, asset string) (marginaccounts.HoldingOutput, error) {
	symbol, err := ToBytes32(asset)
	return read(a.c, opts, a.with(err), a.accounts.UnpackHolding)(a.accounts.TryPackHolding(symbol))
}

// Closure reads the closure the premium accrues over and when it last accrued, in milliseconds.
func (a *Accounts) Closure(opts *bind.CallOpts) (marginaccounts.ClosureOutput, error) {
	return read(a.c, opts, a.target, a.accounts.UnpackClosure)(a.accounts.TryPackClosure())
}

// BackstopPremium reads the premium set aside for the backstop and not yet claimed, in USDG.
func (a *Accounts) BackstopPremium(opts *bind.CallOpts) (*big.Int, error) {
	return read(a.c, opts, a.target, a.accounts.UnpackBackstopPremium)(a.accounts.TryPackBackstopPremium())
}

// HasBackstop reports whether the accounts have a backstop, which alone may then write off a position.
func (a *Accounts) HasBackstop(opts *bind.CallOpts) (bool, error) {
	backstop, err := read(a.c, opts, a.target, a.accounts.UnpackBackstop)(a.accounts.TryPackBackstop())
	return backstop != (common.Address{}), err
}

// Lent reads what account's position has lent of token through its lending vault, its fee included.
func (a *Accounts) Lent(opts *bind.CallOpts, account common.Address, position [32]byte, token common.Address) (*big.Int, error) {
	return read(a.c, opts, a.target, a.accounts.UnpackLent)(a.accounts.TryPackLent(account, position, token))
}

// Sellable reads what the liquidator may take of token from account's position now: what it holds and what its
// lending vault can return of what it lent.
func (a *Accounts) Sellable(opts *bind.CallOpts, account common.Address, position [32]byte, token common.Address) (*big.Int, error) {
	return read(a.c, opts, a.target, a.accounts.UnpackSellable)(a.accounts.TryPackSellable(account, position, token))
}

// AccruePremium accrues the premium and records the closure the band's session shows. Anyone may.
func (a *Accounts) AccruePremium() (Tx, error) {
	return a.tx(a.accounts.TryPackAccruePremium())
}

// Sync brings every position's holding of asset down to what the accounts hold after a burn, and syncs its lending
// vault. Anyone may.
func (a *Accounts) Sync(asset string) (Tx, error) {
	symbol, err := ToBytes32(asset)
	return a.with(err).tx(a.accounts.TryPackSync(symbol))
}

// Clear clears account's holding of asset in position once a burn has left it worth less than a raw unit. Anyone may.
func (a *Accounts) Clear(account common.Address, position [32]byte, asset string) (Tx, error) {
	symbol, err := ToBytes32(asset)
	return a.with(err).tx(a.accounts.TryPackClear(account, position, symbol))
}

// Settle takes what the lending vault holds for account's position's recall of token, if it is next in turn, or gives
// up its recalls once nothing it lent is worth anything. Anyone may.
func (a *Accounts) Settle(account common.Address, position [32]byte, token common.Address) (Tx, error) {
	return a.tx(a.accounts.TryPackSettle(account, position, token))
}
