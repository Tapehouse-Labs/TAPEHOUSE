// SPDX-License-Identifier: MIT OR Apache-2.0

// Package indexer follows a chain's logs from the earliest deployment of Tapehouse's contracts, decodes each by the
// ABI of the contract that emitted it, and keeps them in the store. Blocks the chain has not finalized are checked
// against their hashes at every step, and the index rewinds past any block a reorganisation replaced.
package indexer

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/tapehouse/tapehouse/services/internal/catalog"
	"github.com/tapehouse/tapehouse/services/internal/codec"
	"github.com/tapehouse/tapehouse/services/internal/store"
	"golang.org/x/sync/errgroup"
)

// Chain is what the indexer reads of a chain's RPC. *ethclient.Client implements it.
type Chain interface {
	ChainID(ctx context.Context) (*big.Int, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
}

// Update is what one step changed: the events it added, or the block the index rewound to, and the head and the
// finalized block after it.
type Update struct {
	Head      store.Head
	Finalized uint64
	Events    []store.Event
	// Rewound is set when a reorganisation replaced every event above it.
	Rewound *uint64
}

// MaxRange is the most blocks one eth_getLogs query covers, the limit of the RPC providers the registries use.
const MaxRange = 10_000

// Indexer keeps a store in step with a chain.
type Indexer struct {
	chain     Chain
	store     *store.Store
	catalog   *catalog.Catalog
	publish   func(Update)
	log       *slog.Logger
	start     uint64
	span      uint64
	finalized atomic.Uint64
	mu        sync.Mutex
}

// New returns an indexer of the catalog's contracts on chain into s. Each update is passed to publish.
func New(chain Chain, s *store.Store, c *catalog.Catalog, publish func(Update), log *slog.Logger) *Indexer {
	return &Indexer{chain: chain, store: s, catalog: c, publish: publish, log: log, span: MaxRange}
}

// Prepare checks that the chain is the registry's, rebuilds the index where it was built from another catalog, and
// finds the block indexing starts at: the earliest deployment of Tapehouse's contracts.
func (ix *Indexer) Prepare(ctx context.Context) error {
	id, err := ix.chain.ChainID(ctx)
	if err != nil {
		return err
	}
	if id.Uint64() != ix.catalog.ChainID {
		return fmt.Errorf("the RPC serves chain %d, the registry chain %d", id, ix.catalog.ChainID)
	}
	if len(ix.catalog.Tapehouse()) == 0 {
		return fmt.Errorf("the registry of chain %d names none of Tapehouse's contracts", ix.catalog.ChainID)
	}
	fingerprint, err := ix.store.Meta(ctx, "fingerprint")
	if err != nil {
		return err
	}
	if fingerprint != ix.catalog.Fingerprint() {
		if fingerprint != "" {
			ix.log.Info("the registry changed: rebuilding the index from the chain")
		}
		if err := ix.store.Reset(ctx); err != nil {
			return err
		}
		start, err := ix.deployment(ctx)
		if err != nil {
			return err
		}
		if err := ix.store.SetMeta(ctx, "start", strconv.FormatUint(start, 10)); err != nil {
			return err
		}
		if err := ix.store.SetMeta(ctx, "fingerprint", ix.catalog.Fingerprint()); err != nil {
			return err
		}
	}
	start, err := ix.store.Meta(ctx, "start")
	if err != nil {
		return err
	}
	ix.start, err = strconv.ParseUint(start, 10, 64)
	return err
}

// Start is the block indexing starts at.
func (ix *Indexer) Start() uint64 {
	return ix.start
}

// Finalized is the latest block the chain reported finalized; zero where it reports none, as a dev node does.
func (ix *Indexer) Finalized() uint64 {
	return ix.finalized.Load()
}

// deployment finds the earliest block at which one of Tapehouse's contracts has code, by bisection over eth_getCode.
func (ix *Indexer) deployment(ctx context.Context) (uint64, error) {
	latest, err := ix.chain.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	addresses := ix.catalog.Tapehouse()
	blocks := make([]uint64, len(addresses))
	group, ctx := errgroup.WithContext(ctx)
	group.SetLimit(8)
	for i, address := range addresses {
		group.Go(func() error {
			code, err := ix.chain.CodeAt(ctx, address, latest.Number)
			if err != nil {
				return err
			}
			if len(code) == 0 {
				return fmt.Errorf("%s has no code", address)
			}
			low, high := uint64(0), latest.Number.Uint64()
			for low < high {
				middle := low + (high-low)/2
				code, err := ix.chain.CodeAt(ctx, address, new(big.Int).SetUint64(middle))
				if err != nil {
					return err
				}
				if len(code) > 0 {
					high = middle
				} else {
					low = middle + 1
				}
			}
			blocks[i] = low
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return 0, err
	}
	return slices.Min(blocks), nil
}

// Run steps until ctx ends: at once while the index is behind, every interval once it has caught up, and after a
// failure with a backoff of up to a minute.
func (ix *Indexer) Run(ctx context.Context, interval time.Duration) error {
	backoff := interval
	for {
		caughtUp, err := ix.Step(ctx)
		wait := time.Duration(0)
		switch {
		case ctx.Err() != nil:
			return ctx.Err()
		case err != nil:
			ix.log.Warn("indexing step failed", "error", err)
			wait, backoff = backoff, min(2*backoff, time.Minute)
		case caughtUp:
			wait, backoff = interval, interval
		default:
			backoff = interval
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
}

// Step indexes the next range of blocks, up to the latest, and reports whether it reached it. A head whose hash the
// chain no longer has rewinds the index first.
func (ix *Indexer) Step(ctx context.Context) (bool, error) {
	ix.mu.Lock()
	defer ix.mu.Unlock()
	latest, err := ix.chain.HeaderByNumber(ctx, nil)
	if err != nil {
		return false, err
	}
	if err := ix.updateFinalized(ctx); err != nil {
		return false, err
	}
	head, indexed, err := ix.store.Head(ctx)
	if err != nil {
		return false, err
	}
	from := ix.start
	if indexed {
		if head.Number > ix.Finalized() {
			canonical, err := ix.chain.HeaderByNumber(ctx, new(big.Int).SetUint64(head.Number))
			if err != nil {
				return false, err
			}
			if canonical.Hash() != head.Hash {
				return false, ix.rewind(ctx)
			}
		}
		from = head.Number + 1
	}
	if from > latest.Number.Uint64() {
		return true, nil
	}
	to := min(from+ix.span-1, latest.Number.Uint64())
	logs, err := ix.logs(ctx, from, to)
	if err != nil {
		if ix.span > 1 {
			ix.span /= 2
		}
		return false, err
	}
	ix.span = min(2*ix.span, MaxRange)
	hashes, err := ix.hashes(ctx, logs, to, latest)
	if err != nil {
		return false, err
	}
	events := make([]store.Event, 0, len(logs))
	for _, log := range logs {
		if hash, ok := hashes[log.BlockNumber]; ok && hash != log.BlockHash {
			return false, fmt.Errorf("block %d changed while it was read", log.BlockNumber)
		}
		events = append(events, ix.decode(log))
	}
	next := store.Head{Number: to, Hash: hashes[to]}
	if err := ix.store.Append(ctx, events, next); err != nil {
		return false, err
	}
	ix.publish(Update{Head: next, Finalized: ix.Finalized(), Events: events})
	return to == latest.Number.Uint64(), nil
}

func (ix *Indexer) updateFinalized(ctx context.Context) error {
	finalized, err := ix.chain.HeaderByNumber(ctx, big.NewInt(int64(rpc.FinalizedBlockNumber)))
	switch {
	case err != nil && strings.Contains(err.Error(), "finalized block not found"):
		return nil
	case err != nil:
		return err
	}
	ix.finalized.Store(finalized.Number.Uint64())
	return nil
}

// logs reads the logs of the catalog's contracts from from to to, in the order they were emitted: those of the
// contracts whose every event is kept, and those of the events the others keep.
func (ix *Indexer) logs(ctx context.Context, from, to uint64) ([]types.Log, error) {
	all, some, topics := ix.catalog.Filters()
	query := ethereum.FilterQuery{FromBlock: new(big.Int).SetUint64(from), ToBlock: new(big.Int).SetUint64(to)}
	var logs []types.Log
	if len(all) > 0 {
		query.Addresses = all
		found, err := ix.chain.FilterLogs(ctx, query)
		if err != nil {
			return nil, err
		}
		logs = append(logs, found...)
	}
	if len(some) > 0 {
		query.Addresses, query.Topics = some, [][]common.Hash{topics}
		found, err := ix.chain.FilterLogs(ctx, query)
		if err != nil {
			return nil, err
		}
		logs = append(logs, found...)
	}
	slices.SortFunc(logs, func(a, b types.Log) int {
		return cmp.Or(cmp.Compare(a.BlockNumber, b.BlockNumber), cmp.Compare(a.Index, b.Index))
	})
	return logs, nil
}

// hashes reads the hash of block to, and of every block with logs the chain has not finalized, so that a log of a
// block a reorganisation replaced while it was read is never kept.
func (ix *Indexer) hashes(ctx context.Context, logs []types.Log, to uint64, latest *types.Header) (map[uint64]common.Hash, error) {
	blocks := []uint64{to}
	for _, log := range logs {
		if log.BlockNumber > ix.Finalized() && !slices.Contains(blocks, log.BlockNumber) {
			blocks = append(blocks, log.BlockNumber)
		}
	}
	hashes := make([]common.Hash, len(blocks))
	group, ctx := errgroup.WithContext(ctx)
	group.SetLimit(8)
	for i, number := range blocks {
		if number == latest.Number.Uint64() {
			hashes[i] = latest.Hash()
			continue
		}
		group.Go(func() error {
			header, err := ix.chain.HeaderByNumber(ctx, new(big.Int).SetUint64(number))
			if err != nil {
				return err
			}
			hashes[i] = header.Hash()
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return nil, err
	}
	out := make(map[uint64]common.Hash, len(blocks))
	for i, number := range blocks {
		out[number] = hashes[i]
	}
	return out, nil
}

// decode decodes log by the ABI of the contract that emitted it. A log the ABI does not know keeps its topics and data
// under an empty event name, so that nothing the chain emitted is lost.
func (ix *Indexer) decode(log types.Log) store.Event {
	event := store.Event{Block: log.BlockNumber, LogIndex: log.Index, BlockHash: log.BlockHash, Time: log.BlockTimestamp,
		Tx: log.TxHash}
	if len(log.Topics) > 0 {
		if contract, abiEvent, ok := ix.catalog.At(log.Address, log.Topics[0]); ok {
			event.Contract = contract.Name
			if args, err := codec.Event(abiEvent, log); err == nil {
				event.Event, event.Args = abiEvent.RawName, args
				return event
			}
		}
	}
	if event.Contract == "" {
		event.Contract = ix.name(log.Address)
	}
	topics := make([]any, len(log.Topics))
	for i, topic := range log.Topics {
		topics[i] = topic.Hex()
	}
	event.Args = map[string]any{"topics": topics, "data": hexutil.Encode(log.Data)}
	return event
}

func (ix *Indexer) name(address common.Address) string {
	for _, contract := range ix.catalog.Contracts {
		if contract.Address == address {
			return contract.Name
		}
	}
	return address.Hex()
}

// rewind deletes the events of every block the chain replaced: back to the newest block whose hash it still has, or to
// the finalized block.
func (ix *Indexer) rewind(ctx context.Context) error {
	blocks, err := ix.store.Blocks(ctx, ix.Finalized())
	if err != nil {
		return err
	}
	to := ix.Finalized()
	if to < ix.start && ix.start > 0 {
		to = ix.start - 1
	}
	for _, block := range blocks {
		canonical, err := ix.chain.HeaderByNumber(ctx, new(big.Int).SetUint64(block.Number))
		if err != nil && !errors.Is(err, ethereum.NotFound) {
			return err
		}
		if err == nil && canonical.Hash() == block.Hash {
			to = block.Number
			break
		}
	}
	if err := ix.store.Rewind(ctx, to); err != nil {
		return err
	}
	head, _, err := ix.store.Head(ctx)
	if err != nil {
		return err
	}
	ix.log.Warn("a reorganisation replaced indexed blocks", "rewoundTo", to)
	ix.publish(Update{Head: head, Finalized: ix.Finalized(), Rewound: &to})
	return nil
}
