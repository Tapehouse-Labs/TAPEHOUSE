// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/band"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/bandfeed"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/liquidator"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/marginaccounts"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/reopeningauction"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/stocktoken"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/usdg"
)

var (
	botKey, _ = crypto.HexToECDSA("59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d")
	me        = crypto.PubkeyToAddress(botKey.PublicKey)
	other     = common.HexToAddress("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC")
	borrower  = common.HexToAddress("0x90F79bf6EB2c4f870365E785982E1f101E93b906")
	cross     = [32]byte(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"))
)

const openMs = uint64(1_790_000_000_000)

// call is one transaction the chain mined: its contract, its method and its decoded arguments.
type call struct {
	contract, method string
	args             []any
}

func (c call) String() string { return c.contract + "." + c.method + sprint(c.args...) }

// double is one contract: its ABI, its reads and its writes, each of which a call simulates and a mined transaction
// applies.
type double struct {
	name   string
	abi    abi.ABI
	reads  map[string]func(args []any) []any
	writes map[string]func(args []any, mined bool) error
}

// chain is a test double of a chain whose contracts answer through the ABIs their SDK bindings were generated from, so
// the bidder runs through the real SDK, signer and receipts.
type chain struct {
	t            *testing.T
	mu           sync.Mutex
	deployments  *sdk.Deployments
	contracts    map[common.Address]*double
	names        map[string]common.Address
	time, block  uint64
	blockSeconds uint64
	nonce        uint64
	calls        []call
	receipts     map[common.Hash]*types.Receipt
	logs         []types.Log
	onSend       func(call)
	key          *ecdsa.PrivateKey
}

func newChain(t *testing.T) *chain {
	t.Helper()
	c := &chain{t: t, contracts: map[common.Address]*double{}, names: map[string]common.Address{}, time: 1_789_000_000,
		block: 1000, receipts: map[common.Hash]*types.Receipt{}, key: botKey}
	c.deployments = &sdk.Deployments{ChainID: 412346, Tokens: map[string]common.Address{},
		BandFeeds: map[string]common.Address{}, Tapehouse: map[string]common.Address{}}
	for name, metadata := range map[string]*bind.MetaData{
		"Band": &band.BandMetaData, "MarginAccounts": &marginaccounts.MarginAccountsMetaData,
		"Liquidator": &liquidator.LiquidatorMetaData, "ReopeningAuction": &reopeningauction.ReopeningAuctionMetaData,
	} {
		c.deployments.Tapehouse[name] = c.add(name, metadata)
	}
	c.deployments.Tokens["USDG"] = c.add("USDG", &usdg.UsdgMetaData)
	for _, token := range []string{"SPY", "NVDA", "WETH"} {
		c.deployments.Tokens[token] = c.add("Token:"+token, &stocktoken.StockTokenMetaData)
	}
	for _, asset := range []string{"SPY", "NVDA"} {
		c.deployments.BandFeeds[asset] = c.add("BandFeed:"+asset, &bandfeed.BandFeedMetaData)
	}
	return c
}

func (c *chain) add(name string, metadata *bind.MetaData) common.Address {
	parsed, err := metadata.ParseABI()
	if err != nil {
		c.t.Fatal(err)
	}
	address := common.BytesToAddress(crypto.Keccak256([]byte(name))[:20])
	c.contracts[address] = &double{name: name, abi: *parsed, reads: map[string]func([]any) []any{},
		writes: map[string]func([]any, bool) error{}}
	c.names[name] = address
	return address
}

func (c *chain) read(contract, method string, answer func(args []any) []any) {
	c.contracts[c.names[contract]].reads[method] = answer
}

func (c *chain) write(contract, method string, apply func(args []any, mined bool) error) {
	c.contracts[c.names[contract]].writes[method] = apply
}

func (c *chain) nowMs() uint64 { return c.time * 1000 }

// at sets the chain's clock to ms.
func (c *chain) at(ms uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.time = ms / 1000
}

func (c *chain) head() *types.Header {
	c.mu.Lock()
	defer c.mu.Unlock()
	return &types.Header{Number: new(big.Int).SetUint64(c.block), Time: c.time}
}

// mined returns the transactions mined, as "Contract.method[args]", and forgets them.
func (c *chain) mined() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	var out []string
	for _, call := range c.calls {
		out = append(out, call.String())
	}
	c.calls = nil
	return out
}

func (c *chain) took() []call {
	c.mu.Lock()
	defer c.mu.Unlock()
	calls := c.calls
	c.calls = nil
	return calls
}

// emit adds a log to the receipt of the transaction being mined.
func (c *chain) emit(contract, event string) {
	e := c.contracts[c.names[contract]].abi.Events[event]
	c.logs = append(c.logs, types.Log{Address: c.names[contract], Topics: []common.Hash{e.ID}})
}

func (c *chain) dispatch(to *common.Address, data []byte) (*double, *abi.Method, []any) {
	if to == nil {
		c.t.Fatal("a call with no recipient")
	}
	d, ok := c.contracts[*to]
	if !ok {
		c.t.Fatalf("a call to %s, which the chain does not hold", to.Hex())
	}
	method, err := d.abi.MethodById(data)
	if err != nil {
		c.t.Fatalf("%s: %v", d.name, err)
	}
	args, err := method.Inputs.Unpack(data[4:])
	if err != nil {
		c.t.Fatal(err)
	}
	return d, method, args
}

type revertWith struct{ data []byte }

func (r revertWith) Error() string          { return "execution reverted" }
func (r revertWith) ErrorCode() int         { return 3 }
func (r revertWith) ErrorData() interface{} { return hexutil.Encode(r.data) }

// fail is the revert of the error name, as a node answers it, with its data.
func (c *chain) fail(contract, name string, args ...any) error {
	e, ok := c.contracts[c.names[contract]].abi.Errors[name]
	if !ok {
		c.t.Fatalf("%s has no error %s", contract, name)
	}
	packed, err := e.Inputs.Pack(args...)
	if err != nil {
		c.t.Fatal(err)
	}
	return revertWith{append(e.ID[:4:4], packed...)}
}

func (c *chain) CallContract(_ context.Context, msg ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	d, method, args := c.dispatch(msg.To, msg.Data)
	if apply, ok := d.writes[method.Name]; ok {
		if err := apply(args, false); err != nil {
			return nil, err
		}
		return make([]byte, 32*len(method.Outputs)), nil
	}
	answer, ok := d.reads[method.Name]
	if !ok {
		c.t.Fatalf("an unexpected read of %s.%s", d.name, method.Name)
	}
	out, err := method.Outputs.Pack(answer(args)...)
	if err != nil {
		c.t.Fatalf("%s.%s: %v", d.name, method.Name, err)
	}
	return out, nil
}

func (c *chain) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	if _, err := c.CallContract(ctx, msg, nil); err != nil {
		return 0, err
	}
	return 100_000, nil
}

func (c *chain) SendTransaction(_ context.Context, tx *types.Transaction) error {
	c.mu.Lock()
	hook := c.onSend
	d, method, args := c.dispatch(tx.To(), tx.Data())
	sent := call{strings.SplitN(d.name, ":", 2)[0], method.Name, args}
	c.mu.Unlock()
	if hook != nil {
		hook(sent)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil || from != me || tx.Nonce() != c.nonce {
		c.t.Fatalf("a transaction from %s at nonce %d: %v", from.Hex(), tx.Nonce(), err)
	}
	c.nonce++
	c.logs = nil
	status := types.ReceiptStatusSuccessful
	if err := d.writes[method.Name](args, true); err != nil {
		status = types.ReceiptStatusFailed
	}
	c.block++
	c.calls = append(c.calls, sent)
	c.receipts[tx.Hash()] = &types.Receipt{Status: status, BlockNumber: new(big.Int).SetUint64(c.block),
		GasUsed: 100_000, TxHash: tx.Hash(), Logs: toPointers(c.logs)}
	return nil
}

func toPointers(logs []types.Log) []*types.Log {
	out := make([]*types.Log, len(logs))
	for i := range logs {
		out[i] = &logs[i]
	}
	return out
}

func (c *chain) TransactionReceipt(_ context.Context, hash common.Hash) (*types.Receipt, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if receipt, ok := c.receipts[hash]; ok {
		return receipt, nil
	}
	return nil, ethereum.NotFound
}

func (c *chain) HeaderByNumber(_ context.Context, number *big.Int) (*types.Header, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if number == nil {
		number = new(big.Int).SetUint64(c.block)
	}
	at := c.time
	if n := number.Uint64(); n < c.block {
		at -= (c.block - n) * c.blockSeconds
	}
	return &types.Header{Number: number, Time: at, BaseFee: big.NewInt(1e7)}, nil
}

func (c *chain) CodeAt(context.Context, common.Address, *big.Int) ([]byte, error) {
	return []byte{1}, nil
}

func (c *chain) PendingCodeAt(context.Context, common.Address) ([]byte, error) { return []byte{1}, nil }

func (c *chain) PendingNonceAt(context.Context, common.Address) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.nonce, nil
}

func (c *chain) SuggestGasPrice(context.Context) (*big.Int, error)  { return big.NewInt(1e7), nil }
func (c *chain) SuggestGasTipCap(context.Context) (*big.Int, error) { return big.NewInt(0), nil }

func (c *chain) TransactionByHash(context.Context, common.Hash) (*types.Transaction, bool, error) {
	return nil, false, ethereum.NotFound
}

func (c *chain) FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error) {
	return nil, errors.New("the bidder reads events from the indexer")
}

func (c *chain) SubscribeFilterLogs(context.Context, ethereum.FilterQuery, chan<- types.Log) (ethereum.Subscription, error) {
	return nil, errors.New("the bidder reads events from the indexer")
}

func symbolOf(t *testing.T, name string) [32]byte {
	t.Helper()
	word, err := sdk.ToBytes32(name)
	if err != nil {
		t.Fatal(err)
	}
	return word
}

// market is the reopening auction and the band as the tests stage them.
type market struct {
	*chain
	phaseOpenMs uint64
	rounds      map[string]*reopeningauction.ReopeningAuctionRound
	commitments map[[32]byte]reopeningauction.CommitmentsOutput
	bids        map[string][]reopeningauction.ReopeningAuctionBid
	sealed      map[string]bandfeed.SealsOutput
	quotes      map[string]band.QuoteOutput
	outbid      bool
	logs        bytes.Buffer
}

func zeroRound() *reopeningauction.ReopeningAuctionRound {
	zero := func() *big.Int { return new(big.Int) }
	return &reopeningauction.ReopeningAuctionRound{Floor: zero(), Supply: zero(), Deposits: zero(), Pool: zero(),
		Price: zero(), Above: zero(), AtPrice: zero(), Sold: zero(), Paid: zero(), Taken: zero()}
}

// newMarket stages the round of SPY for the open openMs: supply 5e18, a floor of 88 USD, a band sealed at a low edge
// of 100 USD and a band now at 98 USD, so that the bot's price is 95.06 USD.
func newMarket(t *testing.T) *market {
	t.Helper()
	m := &market{chain: newChain(t), phaseOpenMs: openMs, rounds: map[string]*reopeningauction.ReopeningAuctionRound{},
		commitments: map[[32]byte]reopeningauction.CommitmentsOutput{},
		bids:        map[string][]reopeningauction.ReopeningAuctionBid{},
		sealed:      map[string]bandfeed.SealsOutput{}, quotes: map[string]band.QuoteOutput{}}
	spy := zeroRound()
	spy.Floor, spy.Supply, spy.SealMs = big.NewInt(88e8), big.NewInt(5e18), openMs-48_600_000
	m.rounds["SPY"] = spy
	m.rounds["NVDA"] = zeroRound()
	m.sealed["SPY"] = bandfeed.SealsOutput{State: 3, Live: 1, Low: 100e8, High: big.NewInt(110e8), SealedAt: 1}
	m.quotes["SPY"] = band.QuoteOutput{State: 3, Live: 1, Low: 98e8, High: big.NewInt(108e8)}
	m.at(openMs - 1_800_000 - 10*60_000)

	m.read("ReopeningAuction", "phase", func([]any) []any {
		return []any{m.phaseOpenMs, m.phaseOpenMs != 0 && m.nowMs()+1_800_000 >= m.phaseOpenMs}
	})
	for name, value := range map[string]any{"BOND": big.NewInt(100e6), "MIN_BID": big.NewInt(100e6),
		"REVEAL_MS": uint64(1_800_000), "CLEAR_MS": uint64(3_600_000)} {
		m.read("ReopeningAuction", name, func([]any) []any { return []any{value} })
	}
	asset := func(args []any) string {
		word := args[0].([32]byte)
		return strings.TrimRight(string(word[:]), "\x00")
	}
	m.read("ReopeningAuction", "round", func(args []any) []any { return []any{*m.rounds[asset(args)]} })
	m.read("ReopeningAuction", "bids", func(args []any) []any { return []any{m.bids[asset(args)]} })
	m.read("ReopeningAuction", "commitments", func(args []any) []any {
		c, ok := m.commitments[args[2].([32]byte)]
		if !ok {
			return []any{common.Address{}, new(big.Int)}
		}
		return []any{c.Bidder, c.Deposit}
	})
	m.write("ReopeningAuction", "commit", func(args []any, mined bool) error {
		if m.phaseOpenMs == 0 || m.nowMs()+1_800_000 >= m.phaseOpenMs {
			return m.fail("ReopeningAuction", "WrongPhase")
		}
		commitment, deposit := args[1].([32]byte), args[2].(*big.Int)
		if _, held := m.commitments[commitment]; held || deposit.Cmp(big.NewInt(100e6)) < 0 {
			return m.fail("ReopeningAuction", "InvalidCommitment", commitment)
		}
		if mined {
			m.commitments[commitment] = reopeningauction.CommitmentsOutput{Bidder: me, Deposit: deposit}
		}
		return nil
	})
	m.write("ReopeningAuction", "reveal", func(args []any, mined bool) error {
		now := m.nowMs()
		if open := args[1].(uint64); now+1_800_000 < open || now >= open {
			return m.fail("ReopeningAuction", "WrongPhase")
		}
		name := asset(args)
		commitment, err := sdk.Commitment(me, name, args[1].(uint64), args[2].(*big.Int), args[3].(*big.Int), args[4].([32]byte))
		if err != nil {
			m.t.Fatal(err)
		}
		held, ok := m.commitments[commitment]
		if !ok {
			return m.fail("ReopeningAuction", "UnknownCommitment", commitment)
		}
		if mined {
			delete(m.commitments, commitment)
			if m.outbid {
				m.emit("ReopeningAuction", "Outbid")
				return nil
			}
			m.bids[name] = append(m.bids[name], reopeningauction.ReopeningAuctionBid{Bidder: me, Quantity: args[2].(*big.Int),
				Price: args[3].(*big.Int), Escrow: Escrow(args[2].(*big.Int), args[3].(*big.Int))})
		}
		_ = held
		return nil
	})
	m.write("ReopeningAuction", "claim", func(args []any, mined bool) error {
		name, index := asset(args), args[2].(*big.Int).Int64()
		round := m.rounds[name]
		lapsed := !round.Cleared && m.nowMs() >= args[1].(uint64)+3_600_000
		if m.bids[name][index].Claimed || (!round.Cleared && !lapsed) {
			return m.fail("ReopeningAuction", "WrongState")
		}
		if mined {
			m.bids[name][index].Claimed = true
		}
		return nil
	})
	m.read("Band", "quote", func(args []any) []any {
		q := m.quotes[asset(args)]
		return []any{q.State, q.Live, q.Mid, q.HalfBps, q.Low, bigOr(q.High)}
	})
	for _, name := range []string{"SPY", "NVDA"} {
		m.read("BandFeed:"+name, "seals", func([]any) []any {
			s := m.sealed[name]
			return []any{s.State, s.Live, s.Mid, s.HalfBps, s.Low, bigOr(s.High), s.SealedAt}
		})
	}
	m.read("MarginAccounts", "stocks", func([]any) []any {
		return []any{[][32]byte{symbolOf(t, "SPY"), symbolOf(t, "NVDA")},
			[]common.Address{m.deployments.Tokens["SPY"], m.deployments.Tokens["NVDA"]}}
	})
	m.write("USDG", "approve", func([]any, bool) error { return nil })
	return m
}

func bigOr(v *big.Int) *big.Int {
	if v == nil {
		return new(big.Int)
	}
	return v
}

// bidder returns a bidder from the bot's key over the market, with its state at path.
func (m *market) bidder(t *testing.T, path string) *Bidder {
	t.Helper()
	state, err := LoadState(path)
	if err != nil {
		t.Fatal(err)
	}
	return &Bidder{Client: sdk.NewClient(m.chain, m.deployments), Chain: m.chain,
		Opts: bind.NewKeyedTransactor(botKey, big.NewInt(412346)), State: state, Budget: big.NewInt(1000e6),
		DiscountBps: 300, CommitLead: 15 * time.Minute, Assets: []string{"SPY"},
		Log: slog.New(slog.NewTextHandler(&m.logs, nil))}
}

func (m *market) indexer(t *testing.T, serve http.HandlerFunc) string {
	t.Helper()
	server := httptest.NewServer(serve)
	t.Cleanup(server.Close)
	return server.URL
}

func sprint(args ...any) string {
	parts := make([]string, len(args))
	for i, a := range args {
		if word, ok := a.([32]byte); ok {
			parts[i] = hexutil.Encode(word[:])
		} else {
			parts[i] = fmt.Sprint(a)
		}
	}
	return "[" + strings.Join(parts, " ") + "]"
}
