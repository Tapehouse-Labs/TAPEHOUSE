// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/quoterv2"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/shortpositions"
)

// Shorts reads the short positions and packs their transactions.
type Shorts struct {
	target
	c      *Client
	shorts *shortpositions.ShortPositions
}

// Sale is what a sale pays through the asset's pool now, by QuoterV2's quote, and the least to accept.
type Sale struct {
	Proceeds    *big.Int
	MinProceeds *big.Int
}

// BuyBack is what a buy-back costs through the asset's pool now, by QuoterV2's quote, and the most to pay.
type BuyBack struct {
	Cost    *big.Int
	MaxCost *big.Int
}

// Shorts returns the short positions of the registry's .tapehouse.ShortPositions.
func (c *Client) Shorts() *Shorts {
	return &Shorts{
		target: lookup(c.deployments.Tapehouse, "ShortPositions", ".tapehouse"),
		c:      c,
		shorts: shortpositions.NewShortPositions(),
	}
}

// Deposit adds amount of the sender's USDG to account's short of asset.
func (s *Shorts) Deposit(asset string, amount *big.Int, account common.Address) (Tx, error) {
	to, symbol := s.of(asset)
	if err := present(amount); err != nil {
		return Tx{}, err
	}
	return to.tx(s.shorts.TryPackDeposit(symbol, amount, account))
}

// Withdraw sends amount of USDG from account's short of asset to receiver.
func (s *Shorts) Withdraw(asset string, amount *big.Int, account, receiver common.Address) (Tx, error) {
	to, symbol := s.of(asset)
	if err := present(amount); err != nil {
		return Tx{}, err
	}
	return to.tx(s.shorts.TryPackWithdraw(symbol, amount, account, receiver))
}

// Sell borrows amount of asset and sells it for at least minProceeds of USDG, from QuoteSale.
func (s *Shorts) Sell(asset string, amount, minProceeds *big.Int, account common.Address) (Tx, error) {
	to, symbol := s.of(asset)
	if err := present(amount, minProceeds); err != nil {
		return Tx{}, err
	}
	return to.tx(s.shorts.TryPackSell(symbol, amount, minProceeds, account))
}

// Cover buys back amount of asset for at most maxCost of USDG, from QuoteCover, and repays it.
func (s *Shorts) Cover(asset string, amount, maxCost *big.Int, account common.Address) (Tx, error) {
	to, symbol := s.of(asset)
	if err := present(amount, maxCost); err != nil {
		return Tx{}, err
	}
	return to.tx(s.shorts.TryPackCover(symbol, amount, maxCost, account))
}

// Liquidate buys back all account's short of asset owes once it falls short. Anyone may.
func (s *Shorts) Liquidate(account common.Address, asset string) (Tx, error) {
	to, symbol := s.of(asset)
	return to.tx(s.shorts.TryPackLiquidate(account, symbol))
}

// Mark records asset's band centre for the restriction after a 10% fall. Anyone may.
func (s *Shorts) Mark(asset string) (Tx, error) {
	to, symbol := s.of(asset)
	return to.tx(s.shorts.TryPackMark(symbol))
}

// Position reads account's short of asset: its USDG, negative in deficit, what it owes in the token, its borrow
// shares and its book.
func (s *Shorts) Position(opts *bind.CallOpts, account common.Address, asset string) (shortpositions.PositionOutput, error) {
	to, symbol := s.of(asset)
	return read(s.c, opts, to, s.shorts.UnpackPosition)(s.shorts.TryPackPosition(account, symbol))
}

// Health reads account's short of asset at the band's high edge: equity and requirement in USD with 18 decimals, the
// engine's missing bits and the regime.
func (s *Shorts) Health(opts *bind.CallOpts, account common.Address, asset string) (shortpositions.HealthOutput, error) {
	to, symbol := s.of(asset)
	return read(s.c, opts, to, s.shorts.UnpackHealth)(s.shorts.TryPackHealth(account, symbol))
}

// Epoch reads the epoch of asset's current book.
func (s *Shorts) Epoch(opts *bind.CallOpts, asset string) (*big.Int, error) {
	to, symbol := s.of(asset)
	return read(s.c, opts, to, s.shorts.UnpackEpoch)(s.shorts.TryPackEpoch(symbol))
}

// Book reads the borrow shares, USDG and cost index of asset's shorts in book epoch.
func (s *Shorts) Book(opts *bind.CallOpts, asset string, epoch *big.Int) (shortpositions.BookOutput, error) {
	to, symbol := s.of(asset)
	if err := present(epoch); err != nil {
		return shortpositions.BookOutput{}, err
	}
	return read(s.c, opts, to, s.shorts.UnpackBook)(s.shorts.TryPackBook(symbol, epoch))
}

// Restriction reports whether a short of asset sells at no less than the band's centre now, and the close it
// measures from.
func (s *Shorts) Restriction(opts *bind.CallOpts, asset string) (shortpositions.RestrictionOutput, error) {
	to, symbol := s.of(asset)
	return read(s.c, opts, to, s.shorts.UnpackRestriction)(s.shorts.TryPackRestriction(symbol))
}

// Fee reads the fee tier of asset's pool with USDG; zero where it may not be shorted.
func (s *Shorts) Fee(opts *bind.CallOpts, asset string) (*big.Int, error) {
	to, symbol := s.of(asset)
	return read(s.c, opts, to, s.shorts.UnpackFee)(s.shorts.TryPackFee(symbol))
}

// QuoteSale quotes what selling amount of asset pays through its pool now, from Uniswap's QuoterV2, and that less
// slippageBps, from 0 to 10000.
func (s *Shorts) QuoteSale(opts *bind.CallOpts, asset string, amount *big.Int, slippageBps int64) (Sale, error) {
	if err := checkSlippage(slippageBps); err != nil {
		return Sale{}, err
	}
	quote, err := s.quote(opts, asset, amount, true)
	if err != nil {
		return Sale{}, err
	}
	least := new(big.Int).Mul(quote, big.NewInt(10_000-slippageBps))
	return Sale{quote, least.Div(least, bps)}, nil
}

// QuoteCover quotes what buying back amount of asset costs through its pool now, from Uniswap's QuoterV2, and that
// plus slippageBps, from 0 to 10000, rounded up.
func (s *Shorts) QuoteCover(opts *bind.CallOpts, asset string, amount *big.Int, slippageBps int64) (BuyBack, error) {
	if err := checkSlippage(slippageBps); err != nil {
		return BuyBack{}, err
	}
	quote, err := s.quote(opts, asset, amount, false)
	if err != nil {
		return BuyBack{}, err
	}
	most := new(big.Int).Mul(quote, big.NewInt(10_000+slippageBps))
	most.Add(most, big.NewInt(9_999))
	return BuyBack{quote, most.Div(most, bps)}, nil
}

func (s *Shorts) of(asset string) (target, [32]byte) {
	symbol, err := ToBytes32(asset)
	return s.with(err), symbol
}

func checkSlippage(slippageBps int64) error {
	if slippageBps < 0 || slippageBps > 10_000 {
		return fmt.Errorf("slippageBps %d is outside 0 to 10000", slippageBps)
	}
	return nil
}

func (s *Shorts) quote(opts *bind.CallOpts, asset string, amount *big.Int, sale bool) (*big.Int, error) {
	if err := present(amount); err != nil {
		return nil, err
	}
	d := s.c.deployments
	quoter := lookup(d.UniswapV3, "QuoterV2", ".uniswapV3")
	token, err := entry(d.Tokens, asset, ".tokens")
	if err != nil {
		return nil, err
	}
	usdg, err := entry(d.Tokens, "USDG", ".tokens")
	if err != nil {
		return nil, err
	}
	fee, err := s.Fee(opts, asset)
	if err != nil {
		return nil, err
	}
	binding := quoterv2.NewQuoterV2()
	if sale {
		params := quoterv2.IQuoterV2QuoteExactInputSingleParams{
			TokenIn: token, TokenOut: usdg, AmountIn: amount, Fee: fee, SqrtPriceLimitX96: new(big.Int),
		}
		out, err := read(s.c, opts, quoter, binding.UnpackQuoteExactInputSingle)(binding.TryPackQuoteExactInputSingle(params))
		return out.AmountOut, err
	}
	params := quoterv2.IQuoterV2QuoteExactOutputSingleParams{
		TokenIn: usdg, TokenOut: token, Amount: amount, Fee: fee, SqrtPriceLimitX96: new(big.Int),
	}
	out, err := read(s.c, opts, quoter, binding.UnpackQuoteExactOutputSingle)(binding.TryPackQuoteExactOutputSingle(params))
	return out.AmountIn, err
}
