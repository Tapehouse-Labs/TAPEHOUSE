// SPDX-License-Identifier: MIT OR Apache-2.0

// Package sdk reads Tapehouse's price band, its feeds, the margin accounts, the liquidator, the gap backstop, the
// reopening auction, the stock lending vaults, the baskets, the short positions, the Morpho oracles and the gap cover,
// and packs their transactions, over go-ethereum bindings generated from the contracts' ABIs. Addresses come from a
// chain's registry, deployments/<chainId>.json, read at runtime.
package sdk

import (
	"context"
	"errors"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/aggregator"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/stocktoken"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/usdg"
)

// Tx is a transaction to sign and send: its recipient and calldata.
type Tx struct {
	To   common.Address
	Data []byte
}

// Client reads Tapehouse's contracts on one chain and packs their transactions.
type Client struct {
	backend     bind.ContractBackend
	deployments *Deployments
}

// NewClient returns a client of the contracts deployments names, over backend.
func NewClient(backend bind.ContractBackend, deployments *Deployments) *Client {
	return &Client{backend: backend, deployments: deployments}
}

// Deployments returns the registry the client reads its addresses from.
func (c *Client) Deployments() *Deployments {
	return c.deployments
}

// Simulate runs tx from opts.From with eth_call. A revert comes back as a *Revert where its data decodes.
func (c *Client) Simulate(opts *bind.CallOpts, tx Tx) error {
	_, err := c.contract(tx.To).CallRaw(opts, tx.Data)
	return decoded(err)
}

// Send signs tx with opts and sends it. A revert in its gas estimate comes back as a *Revert where its data decodes.
func (c *Client) Send(opts *bind.TransactOpts, tx Tx) (*types.Transaction, error) {
	sent, err := c.contract(tx.To).RawTransact(opts, tx.Data)
	return sent, decoded(err)
}

// TokenPrice reads the latest round of a Chainlink feed that prices the Stock Token.
func (c *Client) TokenPrice(opts *bind.CallOpts, feed TokenPriceFeed) (Round, error) {
	return c.latestRound(opts, feed.Address)
}

// SharePrice reads the latest round of a Chainlink feed that prices the share.
func (c *Client) SharePrice(opts *bind.CallOpts, feed SharePriceFeed) (Round, error) {
	return c.latestRound(opts, feed.Address)
}

// RoundData reads round roundID of a Chainlink feed.
func (c *Client) RoundData(opts *bind.CallOpts, feed common.Address, roundID *big.Int) (Round, error) {
	if err := present(roundID); err != nil {
		return Round{}, err
	}
	binding := aggregator.NewAggregator()
	out, err := read(c, opts, target{address: feed}, binding.UnpackGetRoundData)(binding.TryPackGetRoundData(roundID))
	return Round{out.RoundId, out.Answer, out.StartedAt, out.UpdatedAt, out.AnsweredInRound}, err
}

// LatestRound reads the latest round of a Chainlink feed.
func (c *Client) LatestRound(opts *bind.CallOpts, feed common.Address) (Round, error) {
	return c.latestRound(opts, feed)
}

// BalanceOf reads holder's balance of an ERC-20 token: USDG, WETH, a Stock Token or a basket's shares.
func (c *Client) BalanceOf(opts *bind.CallOpts, token, holder common.Address) (*big.Int, error) {
	binding := usdg.NewUsdg()
	return read(c, opts, target{address: token}, binding.UnpackBalanceOf)(binding.TryPackBalanceOf(holder))
}

// Paused reads whether a Stock Token's transfers are paused, by the issuer's own pause or its registry's: frozen.
func (c *Client) Paused(opts *bind.CallOpts, token common.Address) (bool, error) {
	binding := stocktoken.NewStockToken()
	return read(c, opts, target{address: token}, binding.UnpackPaused)(binding.TryPackPaused())
}

// Approve lets spender take amount of the sender's token.
func (c *Client) Approve(token, spender common.Address, amount *big.Int) (Tx, error) {
	if err := present(amount); err != nil {
		return Tx{}, err
	}
	return target{address: token}.tx(usdg.NewUsdg().TryPackApprove(spender, amount))
}

func (c *Client) latestRound(opts *bind.CallOpts, feed common.Address) (Round, error) {
	binding := aggregator.NewAggregator()
	out, err := read(c, opts, target{address: feed}, binding.UnpackLatestRoundData)(binding.TryPackLatestRoundData())
	return Round{out.RoundId, out.Answer, out.StartedAt, out.UpdatedAt, out.AnsweredInRound}, err
}

func (c *Client) contract(address common.Address) *bind.BoundContract {
	return bind.NewBoundContract(address, abi.ABI{}, c.backend, c.backend, c.backend)
}

// pin returns opts at its block, or at the latest block where opts names none, so that several reads see one state.
func (c *Client) pin(opts *bind.CallOpts) (*bind.CallOpts, error) {
	var pinned bind.CallOpts
	if opts != nil {
		pinned = *opts
	}
	if pinned.BlockNumber == nil {
		ctx := pinned.Context
		if ctx == nil {
			ctx = context.Background()
		}
		header, err := c.backend.HeaderByNumber(ctx, nil)
		if err != nil {
			return nil, err
		}
		pinned.BlockNumber = header.Number
	}
	return &pinned, nil
}

// target is a contract the registry names, or why it names none.
type target struct {
	address common.Address
	err     error
}

func lookup(group map[string]common.Address, name, path string) target {
	address, err := entry(group, name, path)
	return target{address, err}
}

// with returns the target, failing with err where the registry did not.
func (t target) with(err error) target {
	if t.err == nil {
		t.err = err
	}
	return t
}

// tx returns the packed transaction to the target, or why there is none.
func (t target) tx(data []byte, err error) (Tx, error) {
	if t.err != nil {
		return Tx{}, t.err
	}
	if err != nil {
		return Tx{}, err
	}
	return Tx{t.address, data}, nil
}

// read returns a function that runs the packed call on the target with eth_call and unpacks its result.
func read[T any](c *Client, opts *bind.CallOpts, to target, unpack func([]byte) (T, error)) func([]byte, error) (T, error) {
	return func(data []byte, err error) (T, error) {
		var zero T
		if to.err != nil {
			return zero, to.err
		}
		if err != nil {
			return zero, err
		}
		out, err := bind.Call(c.contract(to.address), opts, data, unpack)
		if err != nil {
			return zero, decoded(err)
		}
		return out, nil
	}
}

var bps = big.NewInt(10_000)

// present refuses a nil amount, which go-ethereum's packer cannot encode and panics on.
func present(amounts ...*big.Int) error {
	if slices.Contains(amounts, nil) {
		return errors.New("an amount is nil")
	}
	return nil
}
