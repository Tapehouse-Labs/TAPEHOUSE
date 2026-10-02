// SPDX-License-Identifier: MIT OR Apache-2.0

package indexer_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"path/filepath"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/tapehouse/tapehouse/services/internal/catalog"
	"github.com/tapehouse/tapehouse/services/internal/indexer"
	"github.com/tapehouse/tapehouse/services/internal/store"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/band"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/stocktoken"
)

var (
	bandAddress  = common.HexToAddress("0x00000000000000000000000000000000000000Ba")
	tokenAddress = common.HexToAddress("0x00000000000000000000000000000000000000Aa")
	nvda         = [32]byte{'N', 'V', 'D', 'A'}
	quiet        = slog.New(slog.NewTextHandler(io.Discard, nil))
)

// chain is a double of an Arbitrum chain's RPC: blocks linked by their parent hashes, logs, contracts deployed at a
// block, a finalized block or none, and an eth_getLogs that refuses ranges over a limit.
type chain struct {
	mu        sync.Mutex
	id        uint64
	headers   []*types.Header
	logs      map[uint64][]types.Log
	deployed  map[common.Address]uint64
	finalized uint64
	maxRange  uint64
	stale     bool
	queries   int
}

func newChain() *chain {
	c := &chain{id: 412346, logs: map[uint64][]types.Log{}, deployed: map[common.Address]uint64{}, maxRange: 1 << 62}
	c.headers = []*types.Header{{Number: big.NewInt(0), Time: 1_790_000_000}}
	return c
}

// mine adds a block holding logs, made by fork so that a block mined again after a reorganisation hashes apart.
func (c *chain) mine(fork byte, logs ...types.Log) {
	c.mu.Lock()
	defer c.mu.Unlock()
	parent := c.headers[len(c.headers)-1]
	header := &types.Header{Number: new(big.Int).Add(parent.Number, common.Big1), ParentHash: parent.Hash(),
		Time: parent.Time + 1, Extra: []byte{fork}}
	number := header.Number.Uint64()
	for i := range logs {
		logs[i].BlockNumber, logs[i].BlockHash, logs[i].Index = number, header.Hash(), uint(i)
		logs[i].BlockTimestamp, logs[i].TxHash = header.Time, common.BigToHash(header.Number)
	}
	c.headers = append(c.headers, header)
	c.logs[number] = logs
}

// reorg drops every block from number on.
func (c *chain) reorg(number uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for n := number; n < uint64(len(c.headers)); n++ {
		delete(c.logs, n)
	}
	c.headers = c.headers[:number]
}

func (c *chain) ChainID(context.Context) (*big.Int, error) {
	return new(big.Int).SetUint64(c.id), nil
}

func (c *chain) HeaderByNumber(_ context.Context, number *big.Int) (*types.Header, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch {
	case number == nil:
		return c.headers[len(c.headers)-1], nil
	case number.Int64() == int64(rpc.FinalizedBlockNumber) && c.finalized == 0:
		return nil, errors.New("finalized block not found")
	case number.Int64() == int64(rpc.FinalizedBlockNumber):
		return c.headers[c.finalized], nil
	case number.Uint64() >= uint64(len(c.headers)):
		return nil, ethereum.NotFound
	}
	return c.headers[number.Uint64()], nil
}

func (c *chain) CodeAt(_ context.Context, account common.Address, block *big.Int) ([]byte, error) {
	if deployed, ok := c.deployed[account]; ok && block.Uint64() >= deployed {
		return []byte{0xef}, nil
	}
	return nil, nil
}

func (c *chain) FilterLogs(_ context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.queries++
	from, to := q.FromBlock.Uint64(), q.ToBlock.Uint64()
	if to-from+1 > c.maxRange {
		return nil, fmt.Errorf("eth_getLogs is limited to a %d range", c.maxRange)
	}
	var out []types.Log
	for n := from; n <= to; n++ {
		for _, log := range c.logs[n] {
			if !slices.Contains(q.Addresses, log.Address) {
				continue
			}
			if len(q.Topics) > 0 && !slices.Contains(q.Topics[0], log.Topics[0]) {
				continue
			}
			if c.stale {
				log.BlockHash = common.HexToHash("0xdead")
			}
			out = append(out, log)
		}
	}
	return out, nil
}

func registry(t *testing.T) *catalog.Catalog {
	t.Helper()
	d, err := sdk.ParseDeployments([]byte(`{"chainId":412346,"tapehouse":{"Band":"` + bandAddress.Hex() +
		`","HaltSigner":"` + tokenAddress.Hex() + `"},"tokens":{"NVDA":"` + tokenAddress.Hex() + `"}}`))
	if err != nil {
		t.Fatal(err)
	}
	c, err := catalog.New(d)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func eventOf(t *testing.T, address common.Address, event abi.Event, args ...any) types.Log {
	t.Helper()
	var indexed [][]any
	var data []any
	for i, input := range event.Inputs {
		if input.Indexed {
			indexed = append(indexed, []any{args[i]})
		} else {
			data = append(data, args[i])
		}
	}
	topics, err := abi.MakeTopics(indexed...)
	if err != nil {
		t.Fatal(err)
	}
	packed, err := event.Inputs.NonIndexed().Pack(data...)
	if err != nil {
		t.Fatal(err)
	}
	log := types.Log{Address: address, Topics: []common.Hash{event.ID}, Data: packed}
	for _, topic := range topics {
		log.Topics = append(log.Topics, topic[0])
	}
	return log
}

type fixture struct {
	chain   *chain
	store   *store.Store
	index   *indexer.Indexer
	updates []indexer.Update
	bandABI *abi.ABI
	token   *abi.ABI
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	f := &fixture{chain: newChain()}
	var err error
	if f.store, err = store.Open(filepath.Join(t.TempDir(), "index.db")); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.store.Close() })
	f.bandABI, _ = band.BandMetaData.ParseABI()
	f.token, _ = stocktoken.StockTokenMetaData.ParseABI()
	f.index = indexer.New(f.chain, f.store, registry(t), func(u indexer.Update) { f.updates = append(f.updates, u) }, quiet)
	return f
}

func (f *fixture) halt(t *testing.T, halted bool) types.Log {
	return eventOf(t, bandAddress, f.bandABI.Events["HaltWritten"], nvda, halted, uint64(1), uint64(2))
}

// sync steps until the index reaches the chain's latest block.
func (f *fixture) sync(t *testing.T) {
	t.Helper()
	for range 64 {
		caughtUp, err := f.index.Step(context.Background())
		if caughtUp {
			return
		}
		_ = err
	}
	t.Fatal("the index did not catch up")
}

func (f *fixture) events(t *testing.T) []store.Event {
	t.Helper()
	events, err := f.store.Events(context.Background(), store.Query{Limit: 100})
	if err != nil {
		t.Fatal(err)
	}
	return events
}

func names(events []store.Event) []string {
	out := make([]string, len(events))
	for i, event := range events {
		out[i] = fmt.Sprintf("%d %s %s", event.Block, event.Contract, event.Event)
	}
	return out
}

func TestIndexingStartsAtTheEarliestDeploymentAndKeepsEveryEventDecoded(t *testing.T) {
	f := newFixture(t)
	f.chain.deployed[tokenAddress], f.chain.deployed[bandAddress] = 1, 3
	f.chain.mine(0)
	f.chain.mine(0, eventOf(t, tokenAddress, f.token.Events["OraclePaused"]))
	f.chain.mine(0)
	f.chain.mine(0, f.halt(t, true), eventOf(t, tokenAddress, f.token.Events["Transfer"], bandAddress, tokenAddress,
		big.NewInt(1)), eventOf(t, tokenAddress, f.token.Events["OracleUnpaused"]))
	f.chain.mine(0, types.Log{Address: bandAddress, Topics: []common.Hash{common.HexToHash("0x01")}, Data: []byte{7}})
	f.chain.mine(0)
	if err := f.index.Prepare(context.Background()); err != nil {
		t.Fatal(err)
	}
	if f.index.Start() != 3 {
		t.Fatalf("start %d", f.index.Start())
	}
	f.sync(t)
	events := f.events(t)
	want := []string{"4 tapehouse.Band HaltWritten", "4 tokens.NVDA OracleUnpaused", "5 tapehouse.Band "}
	if !slices.Equal(names(events), want) {
		t.Fatalf("events %v, want %v", names(events), want)
	}
	if events[0].Args["symbol"] != "0x4e56444100000000000000000000000000000000000000000000000000000000" ||
		events[0].Args["halted"] != true || events[0].Time != 1_790_000_004 || events[1].LogIndex != 2 {
		t.Fatalf("HaltWritten %+v", events[0])
	}
	if events[2].Args["data"] != "0x07" || events[2].Args["topics"].([]any)[0] != common.HexToHash("0x01").Hex() {
		t.Fatalf("a log the ABI does not know: %+v", events[2])
	}
	head, _, _ := f.store.Head(context.Background())
	if head.Number != 6 || head.Hash != f.chain.headers[6].Hash() || f.index.Finalized() != 0 {
		t.Fatalf("head %v, finalized %d", head, f.index.Finalized())
	}
	if last := f.updates[len(f.updates)-1]; last.Head != head || len(f.updates) != 1 || len(last.Events) != 3 {
		t.Fatalf("updates %+v", f.updates)
	}
}

func TestARangeTheRPCRefusesIsHalvedUntilItIsServed(t *testing.T) {
	f := newFixture(t)
	f.chain.deployed[bandAddress] = 1
	for i := range 9 {
		f.chain.mine(0, f.halt(t, i%2 == 0))
	}
	f.chain.maxRange = 3
	if err := f.index.Prepare(context.Background()); err != nil {
		t.Fatal(err)
	}
	f.sync(t)
	if events := f.events(t); len(events) != 9 || events[8].Block != 9 {
		t.Fatalf("events %v", names(events))
	}
}

func TestAReorganisationRewindsToTheNewestBlockTheChainStillHas(t *testing.T) {
	f := newFixture(t)
	f.chain.deployed[bandAddress] = 1
	f.chain.mine(0, f.halt(t, true))
	f.chain.mine(0)
	f.chain.mine(0, f.halt(t, false))
	f.chain.mine(0, f.halt(t, true))
	if err := f.index.Prepare(context.Background()); err != nil {
		t.Fatal(err)
	}
	f.sync(t)
	f.chain.reorg(3)
	f.chain.mine(1)
	f.chain.mine(1, f.halt(t, false), f.halt(t, true))
	f.chain.mine(1)
	f.sync(t)
	want := []string{"1 tapehouse.Band HaltWritten", "4 tapehouse.Band HaltWritten", "4 tapehouse.Band HaltWritten"}
	if events := f.events(t); !slices.Equal(names(events), want) || events[1].BlockHash != f.chain.headers[4].Hash() {
		t.Fatalf("events after the reorganisation %v, want %v", names(events), want)
	}
	var rewound []uint64
	for _, update := range f.updates {
		if update.Rewound != nil {
			rewound = append(rewound, *update.Rewound)
		}
	}
	if !slices.Equal(rewound, []uint64{1}) {
		t.Fatalf("rewound to %v", rewound)
	}
}

func TestNothingAtOrBelowTheFinalizedBlockIsRewound(t *testing.T) {
	f := newFixture(t)
	f.chain.deployed[bandAddress] = 1
	for range 4 {
		f.chain.mine(0, f.halt(t, true))
	}
	f.chain.finalized = 3
	if err := f.index.Prepare(context.Background()); err != nil {
		t.Fatal(err)
	}
	f.sync(t)
	if f.index.Finalized() != 3 {
		t.Fatalf("finalized %d", f.index.Finalized())
	}
	f.chain.reorg(4)
	f.chain.mine(1)
	f.sync(t)
	if events := f.events(t); len(events) != 3 || events[2].Block != 3 {
		t.Fatalf("events %v", names(events))
	}
}

func TestALogOfABlockReplacedWhileItWasReadIsNotKept(t *testing.T) {
	f := newFixture(t)
	f.chain.deployed[bandAddress] = 1
	f.chain.mine(0, f.halt(t, true))
	f.chain.stale = true
	if err := f.index.Prepare(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := f.index.Step(context.Background()); err == nil || err.Error() != "block 1 changed while it was read" {
		t.Fatalf("a stale log: %v", err)
	}
	if _, ok, _ := f.store.Head(context.Background()); ok || len(f.updates) != 0 {
		t.Fatal("a stale log was kept")
	}
	f.chain.stale = false
	f.sync(t)
	if len(f.events(t)) != 1 {
		t.Fatal("the log was not kept once read again")
	}
}

func TestTheIndexRefusesAnotherChainAndRebuildsForAnotherRegistry(t *testing.T) {
	f := newFixture(t)
	f.chain.deployed[bandAddress] = 1
	f.chain.mine(0, f.halt(t, true))
	ctx := context.Background()
	if err := f.index.Prepare(ctx); err != nil {
		t.Fatal(err)
	}
	f.sync(t)
	again := indexer.New(f.chain, f.store, registry(t), func(indexer.Update) {}, quiet)
	if err := again.Prepare(ctx); err != nil || len(f.events(t)) != 1 {
		t.Fatalf("the same registry: %v", err)
	}
	d, _ := sdk.ParseDeployments([]byte(`{"chainId":412346,"tapehouse":{"Band":"` + bandAddress.Hex() + `"}}`))
	other, _ := catalog.New(d)
	if err := indexer.New(f.chain, f.store, other, func(indexer.Update) {}, quiet).Prepare(ctx); err != nil {
		t.Fatal(err)
	}
	if _, ok, _ := f.store.Head(ctx); ok || len(f.events(t)) != 0 {
		t.Fatal("another registry kept the index")
	}
	f.chain.id = 4663
	if err := again.Prepare(ctx); err == nil || err.Error() != "the RPC serves chain 4663, the registry chain 412346" {
		t.Fatalf("another chain: %v", err)
	}
	f.chain.id = 412346
	d, _ = sdk.ParseDeployments([]byte(`{"chainId":412346,"tokens":{"NVDA":"` + tokenAddress.Hex() + `"}}`))
	none, _ := catalog.New(d)
	if err := indexer.New(f.chain, f.store, none, nil, quiet).Prepare(ctx); err == nil {
		t.Fatal("a registry without Tapehouse's contracts was indexed")
	}
	d, _ = sdk.ParseDeployments([]byte(`{"chainId":412346,"tapehouse":{"Band":"` + tokenAddress.Hex() + `"}}`))
	undeployed, _ := catalog.New(d)
	if err := indexer.New(f.chain, f.store, undeployed, nil, quiet).Prepare(ctx); err == nil {
		t.Fatal("a contract without code was indexed")
	}
}

func TestRunFollowsTheChainUntilItsContextEnds(t *testing.T) {
	f := newFixture(t)
	f.chain.deployed[bandAddress] = 1
	f.chain.mine(0, f.halt(t, true))
	if err := f.index.Prepare(context.Background()); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error)
	go func() { done <- f.index.Run(ctx, time.Millisecond) }()
	deadline := time.Now().Add(5 * time.Second)
	for {
		f.chain.mine(0, f.halt(t, false))
		if head, ok, _ := f.store.Head(context.Background()); ok && head.Number >= 5 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("Run did not follow the chain")
		}
		time.Sleep(5 * time.Millisecond)
	}
	cancel()
	if err := <-done; !errors.Is(err, context.Canceled) {
		t.Fatalf("Run returned %v", err)
	}
}
