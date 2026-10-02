// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
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

// Deposit deposits amount of token from the sender into account's position: Cross or an asset's symbol.
func (a *Accounts) Deposit(position [32]byte, token common.Address, amount *big.Int, account common.Address) (Tx, error) {
	if err := present(amount); err != nil {
		return Tx{}, err
	}
	return a.tx(a.accounts.TryPackDeposit(position, token, amount, account))
}

// Withdraw withdraws amount of token from account's position to receiver.
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
