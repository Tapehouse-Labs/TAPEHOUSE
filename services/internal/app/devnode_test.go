// SPDX-License-Identifier: MIT OR Apache-2.0

//go:build devnode

package app_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapehouse/tapehouse/services/internal/api"
	"github.com/tapehouse/tapehouse/services/internal/app"
	"github.com/tapehouse/tapehouse/services/internal/catalog"
	"github.com/tapehouse/tapehouse/services/internal/codec"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/band"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/bandfeed"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/morphobandoracle"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/stocklendingvault"
	"golang.org/x/time/rate"
)

// P99Budget is the slowest the 99th percentile of the API's reads may be under the load test.
const P99Budget = 25 * time.Millisecond

// served is an indexer and its API, running as cmd/indexer runs them.
type served struct {
	base   string
	cancel func()
	done   chan error
}

// loadKeys are the keyed tier's keys the load test spreads its clients over.
var loadKeys = make([]string, 8)

func serve(t *testing.T, db string, digests string) *served {
	t.Helper()
	cfg, err := app.ConfigFromEnv(func(name string) string {
		switch name {
		case "TAPEHOUSE_DB":
			return db
		case "TAPEHOUSE_LISTEN":
			return "127.0.0.1:0"
		case "TAPEHOUSE_API_KEYS":
			return digests
		case "TAPEHOUSE_RPC_RATE":
			return "1000"
		}
		return os.Getenv(name)
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	ready := make(chan string, 1)
	s := &served{cancel: cancel, done: make(chan error, 1)}
	go func() {
		s.done <- app.Run(ctx, cfg, slog.New(slog.NewTextHandler(io.Discard, nil)), func(address net.Addr) {
			ready <- "http://" + address.String()
		})
	}()
	select {
	case s.base = <-ready:
	case err := <-s.done:
		t.Fatalf("the indexer did not start: %v", err)
	}
	t.Cleanup(func() {
		cancel()
		if err := <-s.done; err != nil {
			t.Errorf("the indexer stopped with %v", err)
		}
	})
	return s
}

func (s *served) get(t *testing.T, path, key string, out any) int {
	t.Helper()
	request, _ := http.NewRequest(http.MethodGet, s.base+path, nil)
	if key != "" {
		request.Header.Set("X-API-Key", key)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()
	if out != nil {
		if err := json.NewDecoder(response.Body).Decode(out); err != nil {
			t.Fatalf("%s: %v", path, err)
		}
	}
	return response.StatusCode
}

type status struct {
	Start     uint64 `json:"start"`
	Finalized uint64 `json:"finalized"`
	Head      struct {
		Number uint64 `json:"number"`
	} `json:"head"`
}

// reach waits until the index covers block.
func (s *served) reach(t *testing.T, key string, block uint64) status {
	t.Helper()
	deadline := time.Now().Add(2 * time.Minute)
	for {
		var st status
		if s.get(t, "/v1/status", key, &st) == http.StatusOK && st.Head.Number >= block {
			return st
		}
		if time.Now().After(deadline) {
			t.Fatalf("the index did not reach block %d", block)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// events reads every event up to block, page by page.
func (s *served) events(t *testing.T, key, query string, block uint64) []api.Event {
	t.Helper()
	var out []api.Event
	after := ""
	for {
		var page struct {
			Events []api.Event `json:"events"`
			Next   *string     `json:"next"`
		}
		path := fmt.Sprintf("/v1/events?limit=1000&to=%d%s", block, query)
		if after != "" {
			path += "&after=" + after
		}
		if s.get(t, path, key, &page) != http.StatusOK {
			t.Fatalf("%s failed", path)
		}
		out = append(out, page.Events...)
		if page.Next == nil {
			return out
		}
		after = *page.Next
	}
}

// relaySource is an integrator's PackageSource over the API's RedStone relay.
type relaySource struct {
	base, key string
}

func (r relaySource) Payload(ctx context.Context, feedIDs []string) ([]byte, error) {
	request, _ := http.NewRequestWithContext(ctx, http.MethodGet,
		r.base+"/v1/redstone/payload?feeds="+url.QueryEscape(strings.Join(feedIDs, ",")), nil)
	request.Header.Set("X-API-Key", r.key)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the relay answered %s", response.Status)
	}
	var signed struct {
		Payload hexutil.Bytes `json:"payload"`
	}
	err = json.NewDecoder(response.Body).Decode(&signed)
	return signed.Payload, err
}

func TestTheIndexerServesTheDevNodesEventsAndViews(t *testing.T) {
	ctx := context.Background()
	client, err := ethclient.DialContext(ctx, os.Getenv("TAPEHOUSE_RPC_URL"))
	if err != nil {
		t.Fatal(err)
	}
	deployments, err := sdk.LoadDeployments(os.Getenv("TAPEHOUSE_DEPLOYMENTS"))
	if err != nil {
		t.Fatal(err)
	}
	contracts, err := catalog.New(deployments)
	if err != nil {
		t.Fatal(err)
	}
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(os.Getenv("PRIVATE_KEY"), "0x"))
	if err != nil {
		t.Fatal("set PRIVATE_KEY to the dev node's funded key")
	}
	key := "the e2e's key " + time.Now().String()
	digests := []string{api.Digest(key)}
	for i := range loadKeys {
		loadKeys[i] = fmt.Sprintf("load key %d %s", i, key)
		digests = append(digests, api.Digest(loadKeys[i]))
	}
	dir := t.TempDir()
	first := serve(t, filepath.Join(dir, "first.db"), strings.Join(digests, ","))
	latest, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	st := first.reach(t, key, latest.Number.Uint64())
	head := st.Head.Number
	indexed := first.events(t, key, "", head)
	fmt.Printf("indexer: %d events of %d contracts from block %d to %d, finalized %d\n", len(indexed),
		len(contracts.Contracts), st.Start, head, st.Finalized)

	// Every log of the catalog's contracts that eth_getLogs returns is indexed, decoded by name, and nothing else is.
	all, some, topics := contracts.Filters()
	var raw []types.Log
	for from := st.Start; from <= head; from += 10_000 {
		to := min(from+9_999, head)
		query := ethereum.FilterQuery{FromBlock: new(big.Int).SetUint64(from), ToBlock: new(big.Int).SetUint64(to),
			Addresses: all}
		logs, err := client.FilterLogs(ctx, query)
		if err != nil {
			t.Fatal(err)
		}
		raw = append(raw, logs...)
		query.Addresses, query.Topics = some, [][]common.Hash{topics}
		if logs, err = client.FilterLogs(ctx, query); err != nil {
			t.Fatal(err)
		}
		raw = append(raw, logs...)
	}
	byPosition := map[[2]uint64]types.Log{}
	for _, log := range raw {
		byPosition[[2]uint64{log.BlockNumber, uint64(log.Index)}] = log
	}
	counts := map[string]int{}
	for _, event := range indexed {
		log, ok := byPosition[[2]uint64{event.Block, uint64(event.LogIndex)}]
		if !ok || log.TxHash != event.Tx || log.BlockHash != event.BlockHash || event.Time != log.BlockTimestamp {
			t.Fatalf("%s %s at %d:%d is not the chain's log", event.Contract, event.Event.Event, event.Block, event.LogIndex)
		}
		if event.Event.Event == "" {
			t.Errorf("%s emitted a log its ABI does not know: %v", event.Contract, event.Args)
		}
		counts[event.Contract+" "+event.Event.Event]++
	}
	if len(raw) != len(indexed) {
		t.Fatalf("eth_getLogs returned %d logs, the index holds %d", len(raw), len(indexed))
	}
	names := make([]string, 0, len(counts))
	for name := range counts {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("  %s %d\n", name, counts[name])
	}
	for _, name := range []string{
		"tapehouse.Band PriceWritten", "tapehouse.Band Anchored", "tapehouse.Band HaltWritten",
		"tapehouse.Margin VolatilitySet", "tapehouse.Margin CorrelationSet", "tapehouse.MarginAccounts Borrow",
		"tapehouse.MarginAccounts ClosureSet", "tapehouse.Liquidator AuctionStarted", "tapehouse.GapBackstop Covered",
		"tapehouse.Band MultiplierRecorded", "tapehouse.ShortPositions Sell", "tapehouse.StockLending.SPY Borrow",
		"tapehouse.Baskets.PAIR Deposit", "tapehouse.GapCover Bought", "morphoOracles.NVDA OwnershipTransferred",
	} {
		if counts[name] == 0 {
			t.Errorf("no %s indexed", name)
		}
	}

	// The generic decoding agrees with the bindings' typed unpackers.
	unpacked := 0
	for _, event := range indexed {
		log := byPosition[[2]uint64{event.Block, uint64(event.LogIndex)}]
		var typed any
		switch {
		case event.Event.Event == "HaltWritten" && event.Contract == "tapehouse.Band":
			typed, err = band.NewBand().UnpackHaltWrittenEvent(&log)
		case event.Event.Event == "Sealed":
			typed, err = bandfeed.NewBandFeed().UnpackSealedEvent(&log)
		case event.Event.Event == "OwnershipTransferred" && strings.HasPrefix(event.Contract, "morphoOracles."):
			typed, err = morphobandoracle.NewMorphoBandOracle().UnpackOwnershipTransferredEvent(&log)
		default:
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		value := reflect.ValueOf(typed).Elem()
		for i := range value.NumField() {
			field := value.Type().Field(i).Name
			if field == "Raw" {
				continue
			}
			name := strings.ToLower(field[:1]) + field[1:]
			if fmt.Sprint(plainOf(value.Field(i).Interface())) != fmt.Sprint(event.Args[name]) {
				t.Errorf("%s %s: %s is %v, the binding unpacks %v", event.Contract, event.Event.Event, name,
					event.Args[name], value.Field(i).Interface())
			}
		}
		unpacked++
	}
	fmt.Printf("indexer: %d events decoded as the bindings' unpackers decode them\n", unpacked)

	// A view answers what the binding reads at the same block.
	var view struct {
		Block   uint64         `json:"block"`
		Outputs map[string]any `json:"outputs"`
	}
	first.get(t, "/v1/contracts/tapehouse.StockLending.SPY/views/utilization", key, &view)
	vault := stocklendingvault.NewStockLendingVault()
	data, err := client.CallContract(ctx, ethereum.CallMsg{To: ptr(deployments.StockLending["SPY"]),
		Data: vault.PackUtilization()}, new(big.Int).SetUint64(view.Block))
	if err != nil {
		t.Fatal(err)
	}
	utilization, _ := vault.UnpackUtilization(data)
	if view.Outputs["arg0"] != utilization.String() {
		t.Fatalf("utilization %v, the binding reads %v", view.Outputs, utilization)
	}
	sdkClient := sdk.NewClient(client, deployments)

	// The halts, labelled, and the Morpho oracles as the SDK reads them.
	var halts struct {
		Block       uint64           `json:"block"`
		HaltSigner  string           `json:"haltSigner"`
		SignedInput string           `json:"signedInput"`
		Halts       []map[string]any `json:"halts"`
	}
	first.get(t, "/v1/halts", key, &halts)
	signer := deployments.Tapehouse["HaltSigner"].Hex()
	for _, event := range indexed {
		if event.Contract == "tapehouse.Band" && event.Event.Event == "HaltSignerUpdated" {
			signer = event.Args["newSigner"].(string)
		}
	}
	if halts.HaltSigner != signer || halts.SignedInput != api.SignedInput ||
		len(halts.Halts) != len(deployments.BandFeeds) {
		t.Fatalf("halts %+v", halts)
	}
	var oracles struct {
		Block   uint64           `json:"block"`
		Oracles []map[string]any `json:"oracles"`
	}
	first.get(t, "/v1/morpho/oracles", key, &oracles)
	for _, oracle := range oracles.Oracles {
		price, err := sdkClient.MorphoOracles().Price(&bind.CallOpts{Context: ctx,
			BlockNumber: new(big.Int).SetUint64(oracles.Block)}, oracle["asset"].(string))
		if err != nil {
			t.Fatal(err)
		}
		if (price.Price == nil) != (oracle["price"] == nil) || (price.Price != nil && price.Price.String() != oracle["price"]) ||
			(price.Price == nil && string(price.NoPrice) != oracle["noPrice"]) {
			t.Fatalf("%s: %v, the SDK reads %+v", oracle["asset"], oracle, price)
		}
		fmt.Printf("indexer: morpho oracle %s price %v noPrice %v at block %d\n", oracle["asset"], oracle["price"],
			oracle["noPrice"], oracles.Block)
	}
	var debts struct {
		Block     uint64           `json:"block"`
		Positions []map[string]any `json:"positions"`
	}
	first.get(t, "/v1/accounts/debts", key, &debts)
	for _, position := range debts.Positions {
		repayment, err := sdkClient.Accounts().Repayment(&bind.CallOpts{Context: ctx,
			BlockNumber: new(big.Int).SetUint64(debts.Block)}, common.HexToAddress(position["account"].(string)),
			[32]byte(common.HexToHash(position["position"].(string))))
		if err != nil || repayment.Debt.String() != position["debt"] || repayment.Debt.Sign() == 0 {
			t.Fatalf("debt %v, the SDK reads %+v, %v", position, repayment, err)
		}
	}
	var shorts struct {
		Shorts []map[string]any `json:"shorts"`
	}
	first.get(t, "/v1/shorts", key, &shorts)
	fmt.Printf("indexer: %d positions with debt and %d open shorts at block %d\n", len(debts.Positions),
		len(shorts.Shorts), debts.Block)

	// The stream sends an authorization the moment it is indexed, and the relay's packages are written by the band.
	stream, _, err := websocket.Dial(ctx, "ws"+strings.TrimPrefix(first.base, "http")+
		"/v1/stream?contract=tapehouse.Band&event=PriceWritten", &websocket.DialOptions{HTTPHeader: http.Header{"X-Api-Key": {key}}})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = stream.CloseNow() }()
	chainID, _ := client.ChainID(ctx)
	owner := bind.NewKeyedTransactor(privateKey, chainID)
	source := relaySource{first.base, key}
	var receipt *types.Receipt
	for attempt := 0; receipt == nil; attempt++ {
		write, err := sdkClient.Band().WritePrices(ctx, source, []string{"NVDA---24_7"})
		if err != nil {
			t.Fatal(err)
		}
		if err := sdkClient.Simulate(&bind.CallOpts{Context: ctx, From: owner.From}, write); err != nil {
			if revert, ok := sdk.DecodeRevert(err); ok && revert.Name == "PackageNotNewer" && attempt < 6 {
				time.Sleep(5 * time.Second)
				continue
			}
			t.Fatalf("the band refused the relay's packages: %v", err)
		}
		sent, err := sdkClient.Send(owner, write)
		if err != nil {
			t.Fatal(err)
		}
		if receipt, err = bind.WaitMined(ctx, client, sent.Hash()); err != nil || receipt.Status != types.ReceiptStatusSuccessful {
			t.Fatalf("writePrices: %v", err)
		}
	}
	sentAt := time.Now()
	for {
		read, cancel := context.WithTimeout(ctx, 30*time.Second)
		_, data, err := stream.Read(read)
		cancel()
		if err != nil {
			t.Fatal(err)
		}
		var message api.Message
		if err := json.Unmarshal(data, &message); err != nil {
			t.Fatal(err)
		}
		if message.Type == "event" && message.Event.Tx == receipt.TxHash {
			fmt.Printf("indexer: the relay's NVDA---24_7 package (%v) written in block %d and streamed %v after its receipt\n",
				message.Event.Args["value"], receipt.BlockNumber, time.Since(sentAt).Round(time.Millisecond))
			break
		}
	}

	// A second indexer rebuilds the same index from the chain alone.
	head = first.reach(t, key, receipt.BlockNumber.Uint64()).Head.Number
	second := serve(t, filepath.Join(dir, "second.db"), strings.Join(digests, ","))
	second.reach(t, key, head)
	rebuilt := second.events(t, key, "", head)
	again := first.events(t, key, "", head)
	if len(rebuilt) != len(again) {
		t.Fatalf("the rebuild holds %d events, the first index %d", len(rebuilt), len(again))
	}
	left, _ := json.Marshal(rebuilt)
	right, _ := json.Marshal(again)
	if !bytes.Equal(left, right) {
		t.Fatal("the rebuilt index differs from the first")
	}
	fmt.Printf("indexer: a second index rebuilt from the chain alone holds the same %d events to block %d\n",
		len(rebuilt), head)

	// The API under load: 64 clients on 8 keys, each key at 95% of the keyed tier's rate; the public tier is limited.
	paths := []string{"/v1/status", "/v1/events?limit=100&contract=tapehouse.Band",
		"/v1/events?contract=tapehouse.MarginAccounts&event=Borrow&arg[position]=0x0000000000000000000000000000000000000000000000000000000000000000",
		"/v1/contracts/tapehouse.StockLending.SPY/views/utilization", "/v1/halts", "/v1/accounts/debts",
		"/v1/morpho/oracles", "/v1/contracts"}
	const requests, clients = 40_000, 64
	latencies := make([]time.Duration, requests)
	var wg sync.WaitGroup
	httpClient := &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: clients}}
	pace := make([]*rate.Limiter, len(loadKeys))
	for i := range pace {
		pace[i] = rate.NewLimiter(api.KeyedLimit.Rate*95/100, 1)
	}
	began := time.Now()
	for c := range clients {
		wg.Add(1)
		go func() {
			defer wg.Done()
			key := c % len(loadKeys)
			for i := c; i < requests; i += clients {
				if err := pace[key].Wait(ctx); err != nil {
					t.Error(err)
					return
				}
				request, _ := http.NewRequest(http.MethodGet, first.base+paths[i%len(paths)], nil)
				request.Header.Set("X-API-Key", loadKeys[key])
				start := time.Now()
				response, err := httpClient.Do(request)
				if err != nil {
					t.Error(err)
					return
				}
				_, _ = io.Copy(io.Discard, response.Body)
				_ = response.Body.Close()
				latencies[i] = time.Since(start)
				if response.StatusCode != http.StatusOK {
					t.Errorf("%s answered %d under load", paths[i%len(paths)], response.StatusCode)
					return
				}
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(began)
	slices.Sort(latencies)
	p50, p99 := latencies[requests/2], latencies[requests*99/100]
	fmt.Printf("indexer: %d requests from %d clients on %d keys in %v (%.0f a second): p50 %v, p99 %v, max %v\n",
		requests, clients, len(loadKeys), elapsed.Round(time.Millisecond), float64(requests)/elapsed.Seconds(),
		p50.Round(time.Microsecond),
		p99.Round(time.Microsecond), latencies[requests-1].Round(time.Microsecond))
	if p99 > P99Budget {
		t.Errorf("p99 %v is over the budget, %v", p99, P99Budget)
	}
	limited := 0
	for range api.PublicLimit.Burst + 10 {
		if first.get(t, "/v1/status", "", nil) == http.StatusTooManyRequests {
			limited++
		}
	}
	if limited == 0 {
		t.Fatal("the public tier was not limited")
	}
	fmt.Printf("indexer: the public tier refused %d of %d requests in a burst\n", limited, api.PublicLimit.Burst+10)
	if !t.Failed() {
		fmt.Println("PASS")
	}
}

func ptr[T any](v T) *T {
	return &v
}

// plainOf is a value a typed unpacker returns, in the API's JSON conventions.
func plainOf(v any) any {
	switch v := v.(type) {
	case common.Address:
		return v.Hex()
	case [32]byte:
		return hexutil.Encode(v[:])
	case *big.Int:
		return v.String()
	case bool:
		return v
	case codec.Address:
		return common.Address(v).Hex()
	}
	return fmt.Sprint(v)
}
