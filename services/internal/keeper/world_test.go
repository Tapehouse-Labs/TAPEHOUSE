// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
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
	"github.com/tapehouse/tapehouse/services/sdk/bindings/aggregator"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/band"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/bandfeed"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/basket"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/gapbackstop"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/gapcover"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/liquidator"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/marginaccounts"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/quoterv2"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/reopeningauction"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/shortpositions"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/stocklendingvault"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/stocktoken"
)

// world is a test double of one chain and its index: each contract's reads and writes are handlers over the ABI its
// SDK binding was generated from, so the keepers run through the real SDK, signer and sender.
type world struct {
	t           *testing.T
	mu          sync.Mutex
	deployments *sdk.Deployments
	contracts   map[common.Address]*double
	names       map[string]common.Address
	time        uint64
	block       uint64
	nonce       uint64
	sent        []string
	receipts    map[common.Hash]*types.Receipt
	debts       []map[string]any
	shorts      []map[string]any
	events      []indexed
	index       *httptest.Server
	logs        bytes.Buffer
	key         *ecdsa.PrivateKey
}

type indexed struct {
	contract, event string
	args            map[string]any
}

// double is one contract: its ABI, its reads, and its writes, each of which a call simulates and a mined
// transaction applies.
type double struct {
	name   string
	abi    abi.ABI
	reads  map[string]func(args []any) []any
	writes map[string]func(args []any, mined bool) error
}

var (
	keeperKey, _ = crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcaf784d7bf4f2ff80")
	alice        = common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	bob          = common.HexToAddress("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC")
)

func newWorld(t *testing.T) *world {
	t.Helper()
	w := &world{t: t, contracts: map[common.Address]*double{}, names: map[string]common.Address{},
		time: 1790775000, block: 100, receipts: map[common.Hash]*types.Receipt{}, key: keeperKey}
	w.deployments = &sdk.Deployments{ChainID: 412346, Tokens: map[string]common.Address{},
		Chainlink: map[string]common.Address{}, BandFeeds: map[string]common.Address{},
		Tapehouse: map[string]common.Address{}, StockLending: map[string]common.Address{},
		Baskets: map[string]common.Address{}, UniswapV3: map[string]common.Address{}}
	for name, metadata := range map[string]*bind.MetaData{
		"Band": &band.BandMetaData, "MarginAccounts": &marginaccounts.MarginAccountsMetaData,
		"Liquidator": &liquidator.LiquidatorMetaData, "GapBackstop": &gapbackstop.GapBackstopMetaData,
		"ReopeningAuction": &reopeningauction.ReopeningAuctionMetaData,
		"ShortPositions":   &shortpositions.ShortPositionsMetaData, "GapCover": &gapcover.GapCoverMetaData,
	} {
		w.deployments.Tapehouse[name] = w.add(name, metadata)
	}
	w.deployments.UniswapV3["QuoterV2"] = w.add("QuoterV2", &quoterv2.QuoterV2MetaData)
	for _, token := range []string{"USDG", "WETH"} {
		w.deployments.Tokens[token] = w.add("Token:"+token, &stocktoken.StockTokenMetaData)
	}
	w.index = httptest.NewServer(http.HandlerFunc(w.serve))
	t.Cleanup(w.index.Close)
	return w
}

// add registers a contract double at a fresh address.
func (w *world) add(name string, metadata *bind.MetaData) common.Address {
	parsed, err := metadata.ParseABI()
	if err != nil {
		w.t.Fatal(err)
	}
	address := common.BytesToAddress(crypto.Keccak256([]byte(name))[:20])
	w.contracts[address] = &double{name: name, abi: *parsed, reads: map[string]func([]any) []any{},
		writes: map[string]func([]any, bool) error{}}
	w.names[name] = address
	return address
}

// asset adds a Stock Token, its Chainlink feed and its BandFeed.
func (w *world) asset(name string) {
	w.deployments.Tokens[name] = w.add("Token:"+name, &stocktoken.StockTokenMetaData)
	w.deployments.Chainlink[name+"_USD"] = w.add("Feed:"+name, &aggregator.AggregatorMetaData)
	w.deployments.BandFeeds[name] = w.add("BandFeed:"+name, &bandfeed.BandFeedMetaData)
}

func (w *world) vault(asset string) {
	w.deployments.StockLending[asset] = w.add("Vault:"+asset, &stocklendingvault.StockLendingVaultMetaData)
}

func (w *world) basket(key string) {
	w.deployments.Baskets[key] = w.add("Basket:"+key, &basket.BasketMetaData)
}

func (w *world) read(contract, method string, answer func(args []any) []any) {
	w.contracts[w.names[contract]].reads[method] = answer
}

func (w *world) write(contract, method string, apply func(args []any, mined bool) error) {
	w.contracts[w.names[contract]].writes[method] = apply
}

// keeper returns a keeper of every keeper over the world, its clock the world's.
func (w *world) keeper(cfg Config) *Keeper {
	w.t.Helper()
	cfg.IndexerURL = w.index.URL
	k, err := New(cfg, w.deployments, w, w.key, nil, slog.New(slog.NewTextHandler(&w.logs, nil)))
	if err != nil {
		w.t.Fatal(err)
	}
	k.now = func() time.Time {
		w.mu.Lock()
		defer w.mu.Unlock()
		return time.Unix(int64(w.time), 0)
	}
	k.sleep = func(context.Context, time.Duration) error { return nil }
	return k
}

// Sent returns the writes mined, as "Contract.method(args)", and forgets them.
func (w *world) mined() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	sent := w.sent
	w.sent = nil
	return sent
}

func (w *world) dispatch(to *common.Address, data []byte) (*double, *abi.Method, []any) {
	if to == nil {
		w.t.Fatal("a call with no recipient")
	}
	c, ok := w.contracts[*to]
	if !ok {
		w.t.Fatalf("a call to %s, which the world does not hold", to.Hex())
	}
	method, err := c.abi.MethodById(data)
	if err != nil {
		w.t.Fatalf("%s: %v", c.name, err)
	}
	args, err := method.Inputs.Unpack(data[4:])
	if err != nil {
		w.t.Fatal(err)
	}
	return c, method, args
}

// revertWith is a revert the world answers, as an RPC node does, with its data.
type revertWith struct{ data []byte }

func (r revertWith) Error() string          { return "execution reverted" }
func (r revertWith) ErrorCode() int         { return 3 }
func (r revertWith) ErrorData() interface{} { return hexutil.Encode(r.data) }

// fail returns the revert of error name of contract's ABI.
func (w *world) fail(contract, name string, args ...any) error {
	e, ok := w.contracts[w.names[contract]].abi.Errors[name]
	if !ok {
		w.t.Fatalf("%s has no error %s", contract, name)
	}
	packed, err := e.Inputs.Pack(args...)
	if err != nil {
		w.t.Fatal(err)
	}
	return revertWith{append(e.ID[:4:4], packed...)}
}

func (w *world) CallContract(_ context.Context, call ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	c, method, args := w.dispatch(call.To, call.Data)
	if apply, ok := c.writes[method.Name]; ok {
		if err := apply(args, false); err != nil {
			return nil, err
		}
		return make([]byte, 32*len(method.Outputs)), nil
	}
	answer, ok := c.reads[method.Name]
	if !ok {
		w.t.Fatalf("an unexpected read of %s.%s", c.name, method.Name)
	}
	out, err := method.Outputs.Pack(answer(args)...)
	if err != nil {
		w.t.Fatalf("%s.%s: %v", c.name, method.Name, err)
	}
	return out, nil
}

func (w *world) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	if _, err := w.CallContract(ctx, call, nil); err != nil {
		return 0, err
	}
	return 100_000, nil
}

func (w *world) SendTransaction(_ context.Context, tx *types.Transaction) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil || from != crypto.PubkeyToAddress(w.key.PublicKey) || tx.Nonce() != w.nonce {
		w.t.Fatalf("a transaction from %s at nonce %d: %v", from.Hex(), tx.Nonce(), err)
	}
	w.nonce++
	c, method, args := w.dispatch(tx.To(), tx.Data())
	status := types.ReceiptStatusSuccessful
	if err := c.writes[method.Name](args, true); err != nil {
		status = types.ReceiptStatusFailed
	}
	w.block++
	w.sent = append(w.sent, fmt.Sprintf("%s.%s%v", strings.SplitN(c.name, ":", 2)[0], method.Name, args))
	w.receipts[tx.Hash()] = &types.Receipt{Status: status, BlockNumber: new(big.Int).SetUint64(w.block),
		GasUsed: 100_000, TxHash: tx.Hash()}
	return nil
}

func (w *world) TransactionReceipt(_ context.Context, hash common.Hash) (*types.Receipt, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if receipt, ok := w.receipts[hash]; ok {
		return receipt, nil
	}
	return nil, ethereum.NotFound
}

func (w *world) HeaderByNumber(_ context.Context, number *big.Int) (*types.Header, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if number == nil {
		number = new(big.Int).SetUint64(w.block)
	}
	return &types.Header{Number: number, Time: w.time, BaseFee: big.NewInt(1e7)}, nil
}

func (w *world) CodeAt(context.Context, common.Address, *big.Int) ([]byte, error) {
	return []byte{1}, nil
}

func (w *world) PendingCodeAt(context.Context, common.Address) ([]byte, error) { return []byte{1}, nil }

func (w *world) PendingNonceAt(context.Context, common.Address) (uint64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.nonce, nil
}

func (w *world) SuggestGasPrice(context.Context) (*big.Int, error) { return big.NewInt(1e7), nil }

func (w *world) SuggestGasTipCap(context.Context) (*big.Int, error) { return big.NewInt(0), nil }

func (w *world) TransactionByHash(context.Context, common.Hash) (*types.Transaction, bool, error) {
	return nil, false, ethereum.NotFound
}

func (w *world) FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error) {
	return nil, errors.New("the keepers read events from the index")
}

func (w *world) SubscribeFilterLogs(context.Context, ethereum.FilterQuery, chan<- types.Log) (ethereum.Subscription, error) {
	return nil, errors.New("the keepers read events from the index")
}

// serve answers the indexer's routes the keepers read.
func (w *world) serve(rw http.ResponseWriter, r *http.Request) {
	w.mu.Lock()
	defer w.mu.Unlock()
	var out any
	switch r.URL.Path {
	case "/v1/accounts/debts":
		out = map[string]any{"positions": w.debts}
	case "/v1/shorts":
		out = map[string]any{"shorts": w.shorts}
	case "/v1/events":
		query := r.URL.Query()
		events := []map[string]any{}
		for _, e := range w.events {
			if e.contract != query.Get("contract") || e.event != query.Get("event") {
				continue
			}
			match := true
			for key, values := range query {
				if name, ok := strings.CutPrefix(key, "arg["); ok && fmt.Sprint(e.args[strings.TrimSuffix(name, "]")]) != values[0] {
					match = false
				}
			}
			if match {
				events = append(events, map[string]any{"block": 1, "logIndex": len(events), "args": e.args})
			}
		}
		out = map[string]any{"events": events, "next": nil}
	default:
		http.NotFound(rw, r)
		return
	}
	_ = json.NewEncoder(rw).Encode(out)
}

func (w *world) debt(account common.Address, position [32]byte, equity, requirement int64) {
	w.debts = append(w.debts, map[string]any{"account": account.Hex(), "position": hexutil.Encode(position[:]),
		"debt": "1", "health": map[string]any{"equity": fmt.Sprint(equity), "requirement": fmt.Sprint(requirement)}})
}

func (w *world) event(contract, event string, args map[string]any) {
	w.events = append(w.events, indexed{contract, event, args})
}

func symbol(t *testing.T, name string) [32]byte {
	t.Helper()
	word, err := sdk.ToBytes32(name)
	if err != nil {
		t.Fatal(err)
	}
	return word
}

func n(value int64) *big.Int { return big.NewInt(value) }
