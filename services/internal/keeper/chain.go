// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"errors"
	"log/slog"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/time/rate"
)

// Chain is what the keepers read and send through: go-ethereum's contract and deploy backends.
type Chain interface {
	bind.ContractBackend
	bind.DeployBackend
}

// limitExceeded is the JSON-RPC error code providers answer a client over its plan's rate with.
const limitExceeded = -32005

// Failover is a Chain over several RPC endpoints of one chain, under one rate. Each request goes to the endpoint that
// last answered; a failed connection, an HTTP error or a provider's rate limit moves it on to the next, while any
// other JSON-RPC error, such as a revert, is the chain's answer.
type Failover struct {
	clients []*ethclient.Client
	limiter *rate.Limiter
	log     *slog.Logger
	mu      sync.Mutex
	current int
}

// Dial returns a Failover over urls, at most perSecond requests a second in all.
func Dial(ctx context.Context, urls []string, perSecond float64, log *slog.Logger) (*Failover, error) {
	f := &Failover{limiter: rate.NewLimiter(rate.Limit(perSecond), 1), log: log}
	for _, url := range urls {
		client, err := ethclient.DialContext(ctx, url)
		if err != nil {
			f.Close()
			return nil, errors.New("an RPC URL does not parse")
		}
		f.clients = append(f.clients, client)
	}
	return f, nil
}

// Close closes every endpoint's connection.
func (f *Failover) Close() {
	for _, client := range f.clients {
		client.Close()
	}
}

// failed reports whether err is the endpoint's failure rather than the chain's answer.
func failed(err error) bool {
	if err == nil || errors.Is(err, ethereum.NotFound) || errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	var answer rpc.Error
	if errors.As(err, &answer) {
		return answer.ErrorCode() == limitExceeded
	}
	return true
}

// do runs call on the current endpoint and, while it fails, on each other in turn, once.
func do[T any](ctx context.Context, f *Failover, call func(*ethclient.Client) (T, error)) (T, error) {
	f.mu.Lock()
	start := f.current
	f.mu.Unlock()
	var out T
	var err error
	for i := range f.clients {
		at := (start + i) % len(f.clients)
		if err = f.limiter.Wait(ctx); err != nil {
			return out, err
		}
		out, err = call(f.clients[at])
		if !failed(err) {
			f.mu.Lock()
			if f.current != at {
				f.log.Warn("moved to another RPC endpoint", "endpoint", at)
			}
			f.current = at
			f.mu.Unlock()
			return out, err
		}
	}
	return out, err
}

// CodeAt returns the code of contract at blockNumber.
func (f *Failover) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return do(ctx, f, func(c *ethclient.Client) ([]byte, error) { return c.CodeAt(ctx, contract, blockNumber) })
}

// CallContract runs call with eth_call at blockNumber.
func (f *Failover) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return do(ctx, f, func(c *ethclient.Client) ([]byte, error) { return c.CallContract(ctx, call, blockNumber) })
}

// HeaderByNumber returns the header of block number, the latest where number is nil.
func (f *Failover) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return do(ctx, f, func(c *ethclient.Client) (*types.Header, error) { return c.HeaderByNumber(ctx, number) })
}

// PendingCodeAt returns the code of account in the pending state.
func (f *Failover) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return do(ctx, f, func(c *ethclient.Client) ([]byte, error) { return c.PendingCodeAt(ctx, account) })
}

// PendingNonceAt returns the next nonce of account.
func (f *Failover) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return do(ctx, f, func(c *ethclient.Client) (uint64, error) { return c.PendingNonceAt(ctx, account) })
}

// SuggestGasPrice returns the gas price to pay.
func (f *Failover) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return do(ctx, f, func(c *ethclient.Client) (*big.Int, error) { return c.SuggestGasPrice(ctx) })
}

// SuggestGasTipCap returns the priority fee to pay.
func (f *Failover) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return do(ctx, f, func(c *ethclient.Client) (*big.Int, error) { return c.SuggestGasTipCap(ctx) })
}

// EstimateGas estimates the gas call needs.
func (f *Failover) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	return do(ctx, f, func(c *ethclient.Client) (uint64, error) { return c.EstimateGas(ctx, call) })
}

// SendTransaction sends tx. An endpoint that already knows it, because another relayed it before failing, has it.
func (f *Failover) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	_, err := do(ctx, f, func(c *ethclient.Client) (struct{}, error) { return struct{}{}, c.SendTransaction(ctx, tx) })
	if err != nil && strings.Contains(err.Error(), "already known") {
		return nil
	}
	return err
}

// TransactionByHash returns the transaction hash.
func (f *Failover) TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	type found struct {
		tx      *types.Transaction
		pending bool
	}
	out, err := do(ctx, f, func(c *ethclient.Client) (found, error) {
		tx, pending, err := c.TransactionByHash(ctx, hash)
		return found{tx, pending}, err
	})
	return out.tx, out.pending, err
}

// TransactionReceipt returns the receipt of transaction txHash.
func (f *Failover) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return do(ctx, f, func(c *ethclient.Client) (*types.Receipt, error) { return c.TransactionReceipt(ctx, txHash) })
}

// FilterLogs returns the logs q matches.
func (f *Failover) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return do(ctx, f, func(c *ethclient.Client) ([]types.Log, error) { return c.FilterLogs(ctx, q) })
}

// SubscribeFilterLogs subscribes to the logs q matches, on the current endpoint.
func (f *Failover) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return do(ctx, f, func(c *ethclient.Client) (ethereum.Subscription, error) { return c.SubscribeFilterLogs(ctx, q, ch) })
}
