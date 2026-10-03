// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/basket"
)

// Basket reads one of the registry's baskets and packs its transactions.
type Basket struct {
	target
	c      *Client
	basket *basket.Basket
}

// Components are the assets a basket holds and their Stock Tokens, in the order of every list of amounts.
type Components struct {
	Assets []string
	Tokens []common.Address
}

// Basket returns the registry's basket .tapehouse.Baskets.<key>.
func (c *Client) Basket(key string) *Basket {
	return &Basket{
		target: lookup(c.deployments.Baskets, key, ".tapehouse.Baskets"),
		c:      c,
		basket: basket.NewBasket(),
	}
}

// Mint mints shares to receiver for at most maxAssets of each Stock Token, in the basket's order: PreviewMint reads
// what it takes, its part of every token rounded up.
func (b *Basket) Mint(shares *big.Int, receiver common.Address, maxAssets []*big.Int) (Tx, error) {
	if err := present(append([]*big.Int{shares}, maxAssets...)...); err != nil {
		return Tx{}, err
	}
	return b.tx(b.basket.TryPackMint(shares, receiver, maxAssets))
}

// Redeem redeems shares of owner for the basket's Stock Tokens, to receiver.
func (b *Basket) Redeem(shares *big.Int, receiver, owner common.Address) (Tx, error) {
	if err := present(shares); err != nil {
		return Tx{}, err
	}
	return b.tx(b.basket.TryPackRedeem(shares, receiver, owner))
}

// Rebalance takes assetsIn of the Stock Tokens from the sender and sends assetsOut to receiver, each in the basket's
// order, moving each toward the target in effect. Anyone may.
func (b *Basket) Rebalance(assetsIn, assetsOut []*big.Int, receiver common.Address) (Tx, error) {
	if err := present(append(append([]*big.Int{}, assetsIn...), assetsOut...)...); err != nil {
		return Tx{}, err
	}
	return b.tx(b.basket.TryPackRebalance(assetsIn, assetsOut, receiver))
}

// Components reads the assets the basket holds and their Stock Tokens.
func (b *Basket) Components(opts *bind.CallOpts) (Components, error) {
	out, err := read(b.c, opts, b.target, b.basket.UnpackComponents)(b.basket.TryPackComponents())
	assets := make([]string, len(out.Symbols))
	for i, symbol := range out.Symbols {
		assets[i] = assetName(symbol)
	}
	return Components{assets, out.Tokens}, err
}

// PreviewMint reads what minting shares takes of each Stock Token, in raw units, rounded up.
func (b *Basket) PreviewMint(opts *bind.CallOpts, shares *big.Int) ([]*big.Int, error) {
	if err := present(shares); err != nil {
		return nil, err
	}
	return read(b.c, opts, b.target, b.basket.UnpackPreviewMint)(b.basket.TryPackPreviewMint(shares))
}

// PreviewRedeem reads what redeeming shares gives of each Stock Token, in raw units, rounded down.
func (b *Basket) PreviewRedeem(opts *bind.CallOpts, shares *big.Int) ([]*big.Int, error) {
	if err := present(shares); err != nil {
		return nil, err
	}
	return read(b.c, opts, b.target, b.basket.UnpackPreviewRedeem)(b.basket.TryPackPreviewRedeem(shares))
}

// Target reads the target in effect: raw units of each Stock Token per share.
func (b *Basket) Target(opts *bind.CallOpts) ([]*big.Int, error) {
	return read(b.c, opts, b.target, b.basket.UnpackTarget)(b.basket.TryPackTarget())
}

// TotalSupply reads the basket's shares outstanding.
func (b *Basket) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	return read(b.c, opts, b.target, b.basket.UnpackTotalSupply)(b.basket.TryPackTotalSupply())
}

// PendingTarget reads the target proposed and not yet in effect, and when it takes effect, in seconds; empty and zero
// if none.
func (b *Basket) PendingTarget(opts *bind.CallOpts) (basket.PendingTargetOutput, error) {
	return read(b.c, opts, b.target, b.basket.UnpackPendingTarget)(b.basket.TryPackPendingTarget())
}

// assetName returns an asset's name from its symbol as the contracts hold it, right-padded with zero bytes.
func assetName(symbol [32]byte) string {
	return strings.TrimRight(string(symbol[:]), "\x00")
}
