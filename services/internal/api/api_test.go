// SPDX-License-Identifier: MIT OR Apache-2.0

package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"net/url"
	"path/filepath"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pb33f/libopenapi"
	validator "github.com/pb33f/libopenapi-validator"
	"github.com/tapehouse/tapehouse/services/internal/api"
	"github.com/tapehouse/tapehouse/services/internal/catalog"
	"github.com/tapehouse/tapehouse/services/internal/codec"
	"github.com/tapehouse/tapehouse/services/internal/indexer"
	"github.com/tapehouse/tapehouse/services/internal/redstone"
	"github.com/tapehouse/tapehouse/services/internal/store"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/band"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/morphobandoracle"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/morphoblue"
)

var (
	alice  = common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	bob    = common.HexToAddress("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC")
	nvda   = [32]byte{'N', 'V', 'D', 'A'}
	spy    = [32]byte{'S', 'P', 'Y'}
	market = common.HexToHash("0x0b4a3f1e9de6e6d6b4f53e31dcfd4bb0d7ab8b2e64e2a12a8e8b5b0d9b0a2f11")
	quiet  = slog.New(slog.NewTextHandler(io.Discard, nil))
)

const closesMs = 1_790_200_800_000

func address(n int64) common.Address {
	return common.BigToAddress(big.NewInt(0x1000 + n))
}

// node is a double of the chain's RPC that answers each view of the catalog's contracts by name, or fails every call.
type node struct {
	bind.ContractBackend
	catalog *catalog.Catalog
	views   map[string]func(args []any) ([]any, error)
	calls   atomic.Int64
	fail    error
}

type rpcError struct{ data string }

func (e rpcError) Error() string          { return "execution reverted" }
func (e rpcError) ErrorCode() int         { return 3 }
func (e rpcError) ErrorData() interface{} { return e.data }

func reverts(metadata *bind.MetaData, name string, args ...any) error {
	parsed, _ := metadata.ParseABI()
	e := parsed.Errors[name]
	data, err := e.Inputs.Pack(args...)
	if err != nil {
		panic(err)
	}
	return rpcError{hexutil.Encode(append(e.ID[:4:4], data...))}
}

func (n *node) CallContract(_ context.Context, call ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	n.calls.Add(1)
	if n.fail != nil {
		return nil, n.fail
	}
	for _, contract := range n.catalog.Contracts {
		if contract.Address != *call.To {
			continue
		}
		method, err := contract.ABI.MethodById(call.Data[:4])
		if err != nil {
			return nil, err
		}
		args, err := method.Inputs.Unpack(call.Data[4:])
		if err != nil {
			return nil, err
		}
		answer, ok := n.views[contract.Name+"."+method.Name]
		if !ok {
			return nil, fmt.Errorf("no answer for %s.%s", contract.Name, method.Name)
		}
		outputs, err := answer(args)
		if err != nil {
			return nil, err
		}
		return method.Outputs.Pack(outputs...)
	}
	return nil, fmt.Errorf("no contract at %s", call.To)
}

type index struct{ finalized uint64 }

func (i index) Start() uint64     { return 7 }
func (i index) Finalized() uint64 { return i.finalized }

type relay struct{ fail bool }

func (r relay) Signed(_ context.Context, ids []string) (*redstone.Signed, error) {
	if r.fail {
		return nil, errors.New(`Get "https://oracle-gateway.a.redstone.finance/secret": EOF`)
	}
	return &redstone.Signed{TimestampMs: 1_790_949_660_000, Payload: []byte{1, 2},
		Packages: []redstone.Package{{DataPackageID: ids[0], Signer: codec.Address(redstone.Signers[0]), TimestampMs: 1_790_949_660_000,
			DataPoints: []redstone.DataPoint{{DataFeedID: ids[0], Value: "23729678801"}}}}}, nil
}

type fixture struct {
	t       *testing.T
	server  *httptest.Server
	node    *node
	store   *store.Store
	hub     *api.Hub
	tiers   *api.Tiers
	catalog *catalog.Catalog
	docs    validator.Validator
	key     string
}

func deployments(t *testing.T) *sdk.Deployments {
	t.Helper()
	d, err := sdk.ParseDeployments(fmt.Appendf(nil, `{"chainId":412346,
		"tokens":{"NVDA":"%s","USDG":"%s"},
		"tapehouse":{"HaltSigner":"%s","Band":"%s","MarginAccounts":"%s","ShortPositions":"%s","GapCover":"%s",
			"StockLending":{"NVDA":"%s"}},
		"bandFeeds":{"NVDA":"%s","SPY":"%s"},"morphoOracles":{"NVDA":"%s","SPY":"%s"},
		"morpho":{"Blue":"%s","Markets":{"NVDA_USDG":"%s"}}}`,
		address(1).Hex(), address(2).Hex(), alice.Hex(), address(3).Hex(), address(4).Hex(), address(5).Hex(),
		address(6).Hex(), address(7).Hex(), address(8).Hex(), address(9).Hex(), address(10).Hex(), address(11).Hex(),
		address(12).Hex(), market.Hex()))
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func newFixture(t *testing.T, proxies ...netip.Prefix) *fixture {
	t.Helper()
	d := deployments(t)
	c, err := catalog.New(d)
	if err != nil {
		t.Fatal(err)
	}
	f := &fixture{t: t, catalog: c, hub: api.NewHub(), key: "an integrator's key"}
	f.node = &node{catalog: c, views: answers(d)}
	if f.store, err = store.Open(filepath.Join(t.TempDir(), "index.db")); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.store.Close() })
	if f.tiers, err = api.NewTiers([]string{api.Digest(f.key)}, proxies); err != nil {
		t.Fatal(err)
	}
	server := api.New(api.Config{Deployments: d, Catalog: c, Store: f.store, Index: index{finalized: 20},
		Backend: f.node, Relay: relay{}, Hub: f.hub, Tiers: f.tiers, Log: quiet})
	f.server = httptest.NewServer(server.Handler())
	t.Cleanup(f.server.Close)
	document, err := libopenapi.NewDocument(api.OpenAPI)
	if err != nil {
		t.Fatal(err)
	}
	var errs []error
	if f.docs, errs = validator.NewValidator(document); len(errs) > 0 {
		t.Fatal(errs)
	}
	return f
}

func answers(d *sdk.Deployments) map[string]func([]any) ([]any, error) {
	constant := func(outputs ...any) func([]any) ([]any, error) {
		return func([]any) ([]any, error) { return outputs, nil }
	}
	return map[string]func([]any) ([]any, error){
		"tapehouse.Band.haltSigner": constant(alice),
		"tapehouse.Band.halt": func(args []any) ([]any, error) {
			return []any{args[0] == nvda, uint64(1_790_000_900), uint64(1_790_000_000), false}, nil
		},
		"tapehouse.Band.corporateAction": constant(uint8(0), uint64(0), big.NewInt(0), big.NewInt(0)),
		"tapehouse.Band.price": func([]any) ([]any, error) {
			return nil, reverts(&band.BandMetaData, "PackageNotNewer", nvda, uint64(2), uint64(1))
		},
		"tapehouse.MarginAccounts.debt": func(args []any) ([]any, error) {
			if args[0] == alice {
				return []any{big.NewInt(5_000_000)}, nil
			}
			return []any{big.NewInt(0)}, nil
		},
		"tapehouse.MarginAccounts.health": constant(big.NewInt(1_000), big.NewInt(500), uint8(0), uint8(2)),
		"tapehouse.ShortPositions.position": func(args []any) ([]any, error) {
			if args[0] == alice {
				return []any{big.NewInt(-3), big.NewInt(7), big.NewInt(9), big.NewInt(1)}, nil
			}
			return []any{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(1)}, nil
		},
		"tapehouse.ShortPositions.health": constant(big.NewInt(-1), big.NewInt(2), uint8(0), uint8(1)),
		"tapehouse.GapCover.series": constant(big.NewInt(1_000_000_000), uint64(22_000_000_000),
			uint64(21_000_000_000), uint16(0), uint8(1), true),
		"tapehouse.GapCover.reopenOf":             constant(big.NewInt(1_790_222_400_000)),
		"tapehouse.StockLending.NVDA.utilization": constant(big.NewInt(7_500)),
		"morphoOracles.NVDA.price":                constant(big.NewInt(123)),
		"morphoOracles.SPY.price": func([]any) ([]any, error) {
			return nil, reverts(&morphobandoracle.MorphoBandOracleMetaData, "NoAnswer", spy)
		},
		"morphoOracles.SPY.halt":             constant(true, uint64(1_790_000_900), uint64(1_790_000_000), false),
		"morphoOracles.NVDA.band":            constant(d.Tapehouse["Band"]),
		"morphoOracles.SPY.band":             constant(d.Tapehouse["Band"]),
		"morphoOracles.NVDA.symbol":          constant(nvda),
		"morphoOracles.SPY.symbol":           constant(spy),
		"morphoOracles.NVDA.collateralToken": constant(d.Tokens["NVDA"]),
		"morphoOracles.SPY.collateralToken":  constant(bob),
		"morphoOracles.NVDA.loanToken":       constant(d.Tokens["USDG"]),
		"morphoOracles.SPY.loanToken":        constant(d.Tokens["USDG"]),
		"morphoOracles.NVDA.scaleFactor":     constant(big.NewInt(1e16)),
		"morphoOracles.SPY.scaleFactor":      constant(big.NewInt(1e16)),
		"morphoOracles.NVDA.owner":           constant(alice),
		"morphoOracles.SPY.owner":            constant(alice),
		"morpho.Blue.idToMarketParams": constant(morphoblue.IMorphoMarketParams{LoanToken: d.Tokens["USDG"],
			CollateralToken: d.Tokens["NVDA"], Oracle: d.MorphoOracles["NVDA"], Irm: bob, Lltv: big.NewInt(625e15)}),
		"morpho.Blue.market": constant(big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4), big.NewInt(5),
			big.NewInt(0)),
	}
}

func symbolHex(symbol [32]byte) string {
	return hexutil.Encode(symbol[:])
}

// seed indexes blocks 10 to 30 with the events the composite reads follow.
func (f *fixture) seed() []store.Event {
	f.t.Helper()
	cross := symbolHex([32]byte{})
	event := func(block uint64, index uint, contract, name string, args map[string]any) store.Event {
		return store.Event{Block: block, LogIndex: index, BlockHash: common.BigToHash(big.NewInt(int64(block))),
			Time: 1_790_000_000 + block, Tx: common.BigToHash(big.NewInt(int64(block*10 + uint64(index)))),
			Contract: contract, Event: name, Args: args}
	}
	key := func(slot int) map[string]any {
		return map[string]any{"symbol": symbolHex(nvda), "closesMs": fmt.Sprint(closesMs), "slot": fmt.Sprint(slot),
			"mid": "22000000000", "low": "21900000000", "high": "22100000000"}
	}
	events := []store.Event{
		event(10, 0, "tapehouse.MarginAccounts", "Borrow", map[string]any{"account": alice.Hex(), "position": cross,
			"assets": "5000000", "shares": "5000000", "receiver": alice.Hex()}),
		event(11, 0, "tapehouse.MarginAccounts", "Borrow", map[string]any{"account": bob.Hex(), "position": cross,
			"assets": "1", "shares": "1", "receiver": bob.Hex()}),
		event(12, 0, "tapehouse.ShortPositions", "Sell", map[string]any{"account": alice.Hex(), "symbol": symbolHex(nvda),
			"amount": "1", "proceeds": "2", "shares": "9"}),
		event(12, 1, "tapehouse.ShortPositions", "Sell", map[string]any{"account": bob.Hex(), "symbol": symbolHex(nvda),
			"amount": "1", "proceeds": "2", "shares": "9"}),
		event(14, 0, "tapehouse.Band", "HaltWritten", map[string]any{"symbol": symbolHex(nvda), "halted": true,
			"issuedAt": "1790000000", "expiresAt": "1790000900"}),
		event(13, 0, "tapehouse.GapCover", "Measured", map[string]any{"symbol": symbolHex(nvda),
			"closesMs": fmt.Sprint(closesMs), "weekMove": "41250"}),
		event(15, 0, "tapehouse.GapCover", "Bought", map[string]any{"id": "1", "holder": alice.Hex(),
			"symbol": symbolHex(nvda), "closesMs": fmt.Sprint(closesMs), "notional": "1000000000",
			"deductibleBps": "500", "limitBps": "2000", "premium": "7"}),
		event(15, 1, "tapehouse.GapCover", "Bought", map[string]any{"id": "2", "holder": alice.Hex(),
			"symbol": symbolHex(nvda), "closesMs": "1", "notional": "1", "deductibleBps": "500", "limitBps": "2000",
			"premium": "7"}),
		event(21, 0, "tapehouse.GapCover", "Observed", key(0)),
		event(22, 0, "tapehouse.GapCover", "Observed", key(1)),
		event(25, 0, "tapehouse.GapCover", "Settled", map[string]any{"symbol": symbolHex(nvda),
			"closesMs": fmt.Sprint(closesMs), "referenceRound": "100", "referencePrice": "22000000000",
			"firstRound": "101", "price": "21000000000", "flagged": true}),
	}
	if err := f.store.Append(context.Background(), events, store.Head{Number: 30, Hash: common.HexToHash("0x30")}); err != nil {
		f.t.Fatal(err)
	}
	return events
}

// get requests path, checks the request and the response against the OpenAPI document, and decodes the response.
func (f *fixture) get(path string, header ...string) (int, map[string]any, http.Header) {
	f.t.Helper()
	request, err := http.NewRequest(http.MethodGet, f.server.URL+path, nil)
	if err != nil {
		f.t.Fatal(err)
	}
	for i := 0; i+1 < len(header); i += 2 {
		request.Header.Set(header[i], header[i+1])
	}
	response, err := f.server.Client().Do(request)
	if err != nil {
		f.t.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		f.t.Fatal(err)
	}
	response.Body = io.NopCloser(strings.NewReader(string(body)))
	request.URL.Host, request.URL.Scheme = "127.0.0.1:8080", "http"
	validate := f.docs.ValidateHttpRequestResponse
	if response.StatusCode >= http.StatusBadRequest {
		validate = f.docs.ValidateHttpResponse
	}
	if ok, errs := validate(request, response); !ok {
		for _, e := range errs {
			f.t.Errorf("%s: %s: %s %v", path, e.Message, e.Reason, e.SchemaValidationErrors)
		}
	}
	var out map[string]any
	if strings.HasPrefix(response.Header.Get("Content-Type"), "application/json") {
		if err := json.Unmarshal(body, &out); err != nil {
			f.t.Fatalf("%s: %v", path, err)
		}
	}
	return response.StatusCode, out, response.Header
}

func (f *fixture) keyed(path string) (int, map[string]any) {
	f.t.Helper()
	status, body, _ := f.get(path, "X-API-Key", f.key)
	return status, body
}

func TestEveryRouteIsDocumentedAndEveryDocumentedOperationIsServed(t *testing.T) {
	document, err := libopenapi.NewDocument(api.OpenAPI)
	if err != nil {
		t.Fatal(err)
	}
	model, err := document.BuildV3Model()
	if err != nil {
		t.Fatal(err)
	}
	var documented []string
	for path, item := range model.Model.Paths.PathItems.FromOldest() {
		for method, operation := range item.GetOperations().FromOldest() {
			documented = append(documented, strings.ToUpper(method)+" "+path+" "+operation.OperationId)
		}
	}
	var served []string
	for _, route := range api.Routes {
		served = append(served, route.Method+" "+route.Path+" "+route.OperationID)
	}
	slices.Sort(documented)
	slices.Sort(served)
	if !slices.Equal(documented, served) {
		t.Fatalf("documented:\n%s\nserved:\n%s", strings.Join(documented, "\n"), strings.Join(served, "\n"))
	}
	if model.Model.Version != "3.1.0" {
		t.Fatalf("OpenAPI %s", model.Model.Version)
	}
	v, errs := validator.NewValidator(document)
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	if ok, problems := v.ValidateDocument(); !ok {
		for _, problem := range problems {
			t.Error(problem.Message, problem.Reason)
		}
	}
}

func TestBeforeItsFirstBlockTheIndexAnswersUnavailable(t *testing.T) {
	f := newFixture(t)
	status, body, _ := f.get("/v1/status")
	if status != http.StatusServiceUnavailable || body["error"] != "the index has not reached its first block yet" {
		t.Fatalf("status %d %v", status, body)
	}
	f.seed()
	status, body, _ = f.get("/v1/status")
	head := body["head"].(map[string]any)
	if status != http.StatusOK || body["chainId"] != 412346.0 || body["start"] != 7.0 || body["finalized"] != 20.0 ||
		head["number"] != 30.0 {
		t.Fatalf("status %d %v", status, body)
	}
	status, body, _ = f.get("/v1/contracts")
	contracts := body["contracts"].([]any)
	if status != http.StatusOK || len(contracts) != len(f.catalog.Contracts) {
		t.Fatalf("contracts %d %v", status, body)
	}
	if status, body, _ := f.get("/v1/openapi.yaml"); status != http.StatusOK || body != nil {
		t.Fatalf("openapi %d", status)
	}
}

func TestEventsAreFilteredByTheirArgumentsAndPagedInOrder(t *testing.T) {
	f := newFixture(t)
	events := f.seed()
	status, body := f.keyed("/v1/events?contract=tapehouse.GapCover&event=Bought&arg[symbol]=NVDA&arg[closesMs]=" +
		fmt.Sprint(closesMs))
	page := body["events"].([]any)
	first := page[0].(map[string]any)
	if status != http.StatusOK || len(page) != 1 || first["address"] != address(6).Hex() || first["finalized"] != true ||
		first["args"].(map[string]any)["id"] != "1" || body["next"] != nil {
		t.Fatalf("Bought on NVDA's closure: %d %v", status, body)
	}
	_, body = f.keyed("/v1/events?contract=tapehouse.MarginAccounts&event=Borrow&arg[account]=" +
		strings.ToLower(bob.Hex()))
	if page := body["events"].([]any); len(page) != 1 || page[0].(map[string]any)["block"] != 11.0 {
		t.Fatalf("Borrow by bob: %v", body)
	}
	_, body = f.keyed("/v1/events?contract=tapehouse.GapCover&event=Settled&arg[flagged]=true")
	if page := body["events"].([]any); len(page) != 1 || page[0].(map[string]any)["finalized"] != false {
		t.Fatalf("flagged settlements: %v", body)
	}
	var seen []any
	next := ""
	for range len(events) {
		path := "/v1/events?limit=3&from=10&to=30"
		if next != "" {
			path += "&after=" + next
		}
		_, body = f.keyed(path)
		seen = append(seen, body["events"].([]any)...)
		if body["next"] == nil {
			break
		}
		next = body["next"].(string)
	}
	if len(seen) != len(events) || seen[10].(map[string]any)["event"] != "Settled" {
		t.Fatalf("pages: %d events", len(seen))
	}
	for path, want := range map[string]int{
		"/v1/events?limit=1001":                  http.StatusBadRequest,
		"/v1/events?after=12":                    http.StatusBadRequest,
		"/v1/events?from=-1":                     http.StatusBadRequest,
		"/v1/events?arg[account]=" + alice.Hex(): http.StatusBadRequest,
		"/v1/events?contract=tapehouse.Band&event=HaltWritten&arg[x]=1":                                                                http.StatusBadRequest,
		"/v1/events?contract=tapehouse.Band&event=HaltWritten&arg[symbol]=" + url.QueryEscape("a symbol longer than thirty-two bytes"): http.StatusBadRequest,
		"/v1/events?contract=tapehouse.Options":                                                                                        http.StatusNotFound,
		"/v1/events?contract=tapehouse.Band&event=Sealed&arg[x]=1":                                                                     http.StatusNotFound,
		"/v1/events?contracts=tapehouse.Band":                                                                                          http.StatusBadRequest,
	} {
		if status, body := f.keyed(path); status != want || body["error"] == nil {
			t.Errorf("%s: %d %v, want %d", path, status, body, want)
		}
	}
}

func TestAViewAnswersAtTheHeadOrANamedBlockAndARevertByName(t *testing.T) {
	f := newFixture(t)
	f.seed()
	status, body := f.keyed("/v1/contracts/tapehouse.StockLending.NVDA/views/utilization")
	if status != http.StatusOK || body["block"] != 30.0 || body["outputs"].(map[string]any)["arg0"] != "7500" ||
		body["revert"] != nil {
		t.Fatalf("utilization: %d %v", status, body)
	}
	calls := f.node.calls.Load()
	f.keyed("/v1/contracts/tapehouse.StockLending.NVDA/views/utilization")
	_, body = f.keyed("/v1/contracts/tapehouse.StockLending.NVDA/views/utilization?block=12")
	if f.node.calls.Load() != calls+1 || body["block"] != 12.0 {
		t.Fatalf("calls %d after %d: a read at a block is not cached", f.node.calls.Load(), calls)
	}
	_, body = f.keyed("/v1/contracts/tapehouse.Band/views/halt?symbol=NVDA")
	if halt := body["outputs"].(map[string]any); halt["signedHalt"] != true || halt["until"] != "1790000900" {
		t.Fatalf("halt: %v", body)
	}
	_, body = f.keyed("/v1/contracts/tapehouse.Band/views/price?feedId=NVDA---24_7")
	revert := body["revert"].(map[string]any)
	if body["outputs"] != nil || revert["name"] != "PackageNotNewer" || revert["args"].([]any)[1] != "2" {
		t.Fatalf("a revert: %v", body)
	}
	for path, want := range map[string]int{
		"/v1/contracts/tapehouse.Band/views/halt":                        http.StatusBadRequest,
		"/v1/contracts/tapehouse.Band/views/halt?symbol=NVDA&block=31":   http.StatusBadRequest,
		"/v1/contracts/tapehouse.Band/views/halt?symbol=NVDA&size=1":     http.StatusBadRequest,
		"/v1/contracts/tapehouse.Band/views/halt?symbol=NVDA&block=late": http.StatusBadRequest,
		"/v1/contracts/tapehouse.Band/views/writePrices":                 http.StatusNotFound,
		"/v1/contracts/tapehouse.Band/views/quote2":                      http.StatusNotFound,
		"/v1/contracts/tapehouse.Options/views/halt":                     http.StatusNotFound,
	} {
		if status, body := f.keyed(path); status != want || body["error"] == nil {
			t.Errorf("%s: %d %v, want %d", path, status, body, want)
		}
	}
}

func TestTheKeepersReadsListThePositionsWithDebtAndTheOpenShorts(t *testing.T) {
	f := newFixture(t)
	f.seed()
	status, body := f.keyed("/v1/accounts/debts")
	positions := body["positions"].([]any)
	if status != http.StatusOK || len(positions) != 1 {
		t.Fatalf("debts: %d %v", status, body)
	}
	if debt := positions[0].(map[string]any); debt["account"] != alice.Hex() || debt["debt"] != "5000000" ||
		debt["health"].(map[string]any)["requirement"] != "500" {
		t.Fatalf("alice's debt: %v", debt)
	}
	_, body = f.keyed("/v1/shorts")
	shorts := body["shorts"].([]any)
	if len(shorts) != 1 || shorts[0].(map[string]any)["symbol"] != "NVDA" ||
		shorts[0].(map[string]any)["position"].(map[string]any)["usdgHeld"] != "-3" {
		t.Fatalf("shorts: %v", body)
	}
}

func TestTheHaltIsLabelledTapehousesOneSignedInput(t *testing.T) {
	f := newFixture(t)
	f.seed()
	status, body := f.keyed("/v1/halts")
	halts := body["halts"].([]any)
	if status != http.StatusOK || body["haltSigner"] != alice.Hex() || body["signedInput"] != api.SignedInput ||
		len(halts) != 2 {
		t.Fatalf("halts: %d %v", status, body)
	}
	if nvdaHalt := halts[0].(map[string]any); nvdaHalt["symbol"] != "NVDA" || nvdaHalt["signedHalt"] != true {
		t.Fatalf("NVDA: %v", nvdaHalt)
	}
	if spyHalt := halts[1].(map[string]any); spyHalt["symbol"] != "SPY" || spyHalt["signedHalt"] != false {
		t.Fatalf("SPY: %v", spyHalt)
	}
}

func TestAMorphoOracleWithoutAPriceSaysWhyAndTheMarketsAreServedById(t *testing.T) {
	f := newFixture(t)
	f.seed()
	status, body := f.keyed("/v1/morpho/oracles")
	oracles := body["oracles"].([]any)
	if status != http.StatusOK || len(oracles) != 2 {
		t.Fatalf("oracles: %d %v", status, body)
	}
	nvdaOracle, spyOracle := oracles[0].(map[string]any), oracles[1].(map[string]any)
	if nvdaOracle["price"] != "123" || nvdaOracle["noPrice"] != nil || nvdaOracle["symbol"] != "NVDA" ||
		nvdaOracle["scaleFactor"] != "10000000000000000" {
		t.Fatalf("NVDA: %v", nvdaOracle)
	}
	if spyOracle["price"] != nil || spyOracle["noPrice"] != "halted" {
		t.Fatalf("SPY: %v", spyOracle)
	}
	_, body = f.keyed("/v1/morpho/markets")
	markets := body["markets"].([]any)
	nvdaMarket := markets[0].(map[string]any)
	if len(markets) != 1 || nvdaMarket["id"] != market.Hex() ||
		nvdaMarket["params"].(map[string]any)["lltv"] != "625000000000000000" ||
		nvdaMarket["market"].(map[string]any)["totalBorrowAssets"] != "3" {
		t.Fatalf("markets: %v", body)
	}
}

func TestAGapCoverSeriesCarriesWhatRebuildsItsSettlement(t *testing.T) {
	f := newFixture(t)
	f.seed()
	status, body := f.keyed(fmt.Sprintf("/v1/gap-cover/series/NVDA/%d", closesMs))
	series := body["series"].(map[string]any)
	if status != http.StatusOK || series["status"] != "settled" || series["flagged"] != true ||
		series["reopenMs"] != "1790222400000" || series["price"] != "21000000000" {
		t.Fatalf("series: %d %v", status, body)
	}
	observed, settled, measured := body["observed"].([]any), body["settled"].([]any), body["measured"].([]any)
	if len(measured) != 1 || measured[0].(map[string]any)["args"].(map[string]any)["weekMove"] != "41250" {
		t.Fatalf("the week measured: %v", body["measured"])
	}
	if len(body["bought"].([]any)) != 1 || len(observed) != 2 || len(settled) != 1 || len(body["voided"].([]any)) != 0 {
		t.Fatalf("events: %v", body)
	}
	if observed[1].(map[string]any)["args"].(map[string]any)["slot"] != "1" ||
		settled[0].(map[string]any)["args"].(map[string]any)["referenceRound"] != "100" {
		t.Fatalf("the window and the settlement: %v", body)
	}
	if status, _ := f.keyed("/v1/gap-cover/series/NVDA/soon"); status != http.StatusBadRequest {
		t.Fatalf("a closure that is not a time: %d", status)
	}
}

func TestAFailedRPCIsABadGatewayThatNeverShowsTheRPC(t *testing.T) {
	f := newFixture(t)
	f.seed()
	f.node.fail = errors.New(`Post "https://rpc.example/secret-key": dial tcp: connection refused`)
	for _, path := range []string{"/v1/halts", "/v1/accounts/debts", "/v1/contracts/tapehouse.Band/views/haltSigner"} {
		request, _ := http.NewRequest(http.MethodGet, f.server.URL+path, nil)
		request.Header.Set("X-API-Key", f.key)
		response, err := f.server.Client().Do(request)
		if err != nil {
			t.Fatal(err)
		}
		body, _ := io.ReadAll(response.Body)
		_ = response.Body.Close()
		if response.StatusCode != http.StatusBadGateway || strings.Contains(string(body), "secret") ||
			strings.Contains(string(body), "rpc.example") {
			t.Fatalf("%s: %d %s", path, response.StatusCode, body)
		}
	}
}

func TestTheRelayServesSignedPackagesAndNeverTheGatewaysFailure(t *testing.T) {
	f := newFixture(t)
	status, body := f.keyed("/v1/redstone/payload?feeds=NVDA---24_7")
	if status != http.StatusOK || body["payload"] != "0x0102" ||
		body["packages"].([]any)[0].(map[string]any)["signer"] != redstone.Signers[0].Hex() {
		t.Fatalf("relay: %d %v", status, body)
	}
	if status, _ := f.keyed("/v1/redstone/payload"); status != http.StatusBadRequest {
		t.Fatalf("no feeds: %d", status)
	}
	d := deployments(t)
	failing := api.New(api.Config{Deployments: d, Catalog: f.catalog, Store: f.store, Index: index{}, Backend: f.node,
		Relay: relay{fail: true}, Hub: api.NewHub(), Tiers: f.tiers, Log: quiet})
	recorder := httptest.NewRecorder()
	failing.Handler().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/v1/redstone/payload?feeds=NVDA---24_7", nil))
	if recorder.Code != http.StatusServiceUnavailable || strings.Contains(recorder.Body.String(), "secret") {
		t.Fatalf("a failing gateway: %d %s", recorder.Code, recorder.Body)
	}
}

func TestThePublicTierIsLimitedPerAddressAndTheKeyedTierPerKey(t *testing.T) {
	f := newFixture(t)
	f.seed()
	limited := 0
	for range api.PublicLimit.Burst + 5 {
		status, body, header := f.get("/v1/status")
		if status == http.StatusTooManyRequests {
			limited++
			if header.Get("Retry-After") == "" || body["error"] == nil {
				t.Fatalf("429 without Retry-After: %v", header)
			}
		}
		if header.Get("Access-Control-Allow-Origin") != "*" {
			t.Fatal("a public read without CORS")
		}
	}
	if limited == 0 {
		t.Fatal("the public tier was not limited")
	}
	for range api.PublicLimit.Burst + 5 {
		if status, _ := f.keyed("/v1/status"); status != http.StatusOK {
			t.Fatalf("the keyed tier was limited at the public tier's rate: %d", status)
		}
	}
	if status, body, _ := f.get("/v1/status", "X-API-Key", "a guess"); status != http.StatusUnauthorized ||
		body["error"] != "unknown API key" {
		t.Fatalf("an unknown key: %d %v", status, body)
	}
	if _, err := api.NewTiers([]string{"not hex"}, nil); err == nil {
		t.Fatal("a digest that is not hex was taken")
	}
}

func TestBehindATrustedProxyTheClientIsTheForwardedAddress(t *testing.T) {
	f := newFixture(t, netip.MustParsePrefix("127.0.0.0/8"))
	f.seed()
	for i := range api.PublicLimit.Burst + 5 {
		status, _, _ := f.get("/v1/status", "X-Forwarded-For", fmt.Sprintf("203.0.113.%d, 127.0.0.1", i))
		if status != http.StatusOK {
			t.Fatalf("client %d was limited with another's allowance: %d", i, status)
		}
	}
	for range api.PublicLimit.Burst + 1 {
		f.get("/v1/status", "X-Forwarded-For", "198.51.100.1")
	}
	if status, _, _ := f.get("/v1/status", "X-Forwarded-For", "198.51.100.1"); status != http.StatusTooManyRequests {
		t.Fatalf("one forwarded client was not limited: %d", status)
	}
}

func dial(t *testing.T, f *fixture, query string, header http.Header) (*websocket.Conn, *http.Response, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return websocket.Dial(ctx, "ws"+strings.TrimPrefix(f.server.URL, "http")+"/v1/stream"+query,
		&websocket.DialOptions{HTTPHeader: header})
}

func read(t *testing.T, conn *websocket.Conn) api.Message {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var message api.Message
	if err := json.Unmarshal(data, &message); err != nil {
		t.Fatal(err)
	}
	return message
}

func TestTheStreamSendsTheEventsItSelectsTheHeadAndRewinds(t *testing.T) {
	f := newFixture(t)
	events := f.seed()
	conn, _, err := dial(t, f, "?contract=tapehouse.GapCover&event=Observed", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.CloseNow() }()
	deadline := time.Now().Add(5 * time.Second)
	for f.hub.Subscribers() == 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	head := store.Head{Number: 30, Hash: common.HexToHash("0x30")}
	f.hub.Publish(indexer.Update{Head: head, Finalized: 20, Events: events})
	if message := read(t, conn); message.Type != "event" || message.Event.Event.Event != "Observed" ||
		message.Event.Block != 21 || common.Address(message.Event.Address) != address(6) || message.Event.Finalized {
		t.Fatalf("first message %+v", message)
	}
	if message := read(t, conn); message.Type != "event" || message.Event.Block != 22 {
		t.Fatalf("second message %+v", message)
	}
	if message := read(t, conn); message.Type != "head" || *message.Head != head || message.Finalized != 20 {
		t.Fatalf("third message %+v", message)
	}
	to := uint64(21)
	f.hub.Publish(indexer.Update{Head: store.Head{Number: 21}, Finalized: 20, Rewound: &to})
	if message := read(t, conn); message.Type != "rewind" || *message.Block != 21 {
		t.Fatalf("rewind %+v", message)
	}
	if message := read(t, conn); message.Type != "head" || message.Head.Number != 21 {
		t.Fatalf("head after the rewind %+v", message)
	}
}

func TestStreamsAreCountedPerTier(t *testing.T) {
	f := newFixture(t)
	f.seed()
	first, _, err := dial(t, f, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = first.CloseNow() }()
	second, _, err := dial(t, f, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = second.CloseNow() }()
	if _, response, err := dial(t, f, "", nil); err == nil || response.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("a third public stream: %v", err)
	}
	keyed, _, err := dial(t, f, "", http.Header{"X-Api-Key": {f.key}})
	if err != nil {
		t.Fatalf("a keyed stream: %v", err)
	}
	defer func() { _ = keyed.CloseNow() }()
	if _, response, err := dial(t, f, "?contract=tapehouse.Options", nil); err == nil ||
		response.StatusCode != http.StatusNotFound {
		t.Fatalf("a stream of an unknown contract: %v", err)
	}
}
