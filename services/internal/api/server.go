// SPDX-License-Identifier: MIT OR Apache-2.0

// Package api serves the index and the chain's views over REST and a WebSocket stream, as openapi.yaml documents.
// Everything it serves is an event the chain emitted, a view the chain answers at a named block, or a package RedStone
// signed: anyone can reproduce it.
package api

import (
	"context"
	_ "embed" // embeds the OpenAPI document
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/tapehouse/tapehouse/services/internal/catalog"
	"github.com/tapehouse/tapehouse/services/internal/redstone"
	"github.com/tapehouse/tapehouse/services/internal/store"
	"github.com/tapehouse/tapehouse/services/sdk"
)

// OpenAPI is the API's OpenAPI 3.1 document.
//
//go:embed openapi.yaml
var OpenAPI []byte

// Backend is what the API reads of the chain: the views, at a block. *ethclient.Client implements it.
type Backend interface {
	bind.ContractBackend
}

// Index is what the API reads of the indexer.
type Index interface {
	Start() uint64
	Finalized() uint64
}

// Relay supplies RedStone's latest signed packages.
type Relay interface {
	Signed(ctx context.Context, ids []string) (*redstone.Signed, error)
}

// Config is what a server serves.
type Config struct {
	Deployments *sdk.Deployments
	Catalog     *catalog.Catalog
	Store       *store.Store
	Index       Index
	Backend     Backend
	Relay       Relay
	Hub         *Hub
	Tiers       *Tiers
	Log         *slog.Logger
}

// Server serves the API.
type Server struct {
	Config
	sdk   *sdk.Client
	cache *lru
}

// New returns a server of cfg.
func New(cfg Config) *Server {
	s := &Server{Config: cfg, sdk: sdk.NewClient(cfg.Backend, cfg.Deployments), cache: newLRU(4096)}
	cfg.Hub.mu.Lock()
	cfg.Hub.served = s.served
	cfg.Hub.mu.Unlock()
	return s
}

// Route is an operation of the API: its method and path pattern as http.ServeMux takes them, and its OpenAPI
// operationId.
type Route struct {
	Method      string
	Path        string
	OperationID string
	handler     func(*Server, http.ResponseWriter, *http.Request)
}

// Routes are the API's operations.
var Routes = []Route{
	{"GET", "/v1/status", "getStatus", (*Server).status},
	{"GET", "/v1/contracts", "listContracts", (*Server).contracts},
	{"GET", "/v1/events", "listEvents", (*Server).events},
	{"GET", "/v1/contracts/{contract}/views/{function}", "callView", (*Server).view},
	{"GET", "/v1/accounts/debts", "listDebts", (*Server).debts},
	{"GET", "/v1/shorts", "listShorts", (*Server).shorts},
	{"GET", "/v1/halts", "listHalts", (*Server).halts},
	{"GET", "/v1/morpho/oracles", "listMorphoOracles", (*Server).morphoOracles},
	{"GET", "/v1/morpho/markets", "listMorphoMarkets", (*Server).morphoMarkets},
	{"GET", "/v1/gap-cover/series/{symbol}/{closesMs}", "getGapCoverSeries", (*Server).gapCoverSeries},
	{"GET", "/v1/redstone/payload", "getRedStonePayload", (*Server).redstonePayload},
	{"GET", "/v1/stream", "stream", (*Server).stream},
	{"GET", "/v1/openapi.yaml", "getOpenAPI", (*Server).openapi},
}

// Handler returns the API's handler, behind the tiers' keys and rate limits.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	for _, route := range Routes {
		mux.HandleFunc(route.Method+" "+route.Path, func(w http.ResponseWriter, r *http.Request) {
			route.handler(s, w, r)
		})
	}
	return s.Tiers.Middleware(mux)
}

func (s *Server) openapi(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	_, _ = w.Write(OpenAPI)
}

// head is the block reads are made at: the last block the index covers.
func (s *Server) head(ctx context.Context) (store.Head, error) {
	head, ok, err := s.Store.Head(ctx)
	if err == nil && !ok {
		err = errNotIndexed
	}
	return head, err
}

var errNotIndexed = errors.New("the index has not reached its first block yet")

// at returns call options at block.
func at(ctx context.Context, block uint64) *bind.CallOpts {
	return &bind.CallOpts{Context: ctx, BlockNumber: new(big.Int).SetUint64(block)}
}

// badRequest is an error in the request, answered with 400.
type badRequest struct{ error }

func invalid(format string, args ...any) error {
	return badRequest{fmt.Errorf(format, args...)}
}

// reply writes v as JSON, or the error: 400 for a bad request, 404 for something the registry does not name, 503
// before the index's first block, and 502 for a failed RPC, whose details, which may carry the RPC's URL and key,
// are logged and never served.
func (s *Server) reply(w http.ResponseWriter, v any, err error) {
	var bad badRequest
	var missing notFound
	switch {
	case err == nil:
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	case errors.As(err, &bad):
		s.fail(w, http.StatusBadRequest, bad.Error())
	case errors.As(err, &missing):
		s.fail(w, http.StatusNotFound, missing.Error())
	case errors.Is(err, errNotIndexed):
		s.fail(w, http.StatusServiceUnavailable, err.Error())
	case errors.Is(err, context.Canceled):
	default:
		s.Log.Warn("a read failed", "error", err)
		s.fail(w, http.StatusBadGateway, "the chain's RPC failed")
	}
}

type notFound struct{ error }

func missing(format string, args ...any) error {
	return notFound{fmt.Errorf(format, args...)}
}

func (s *Server) fail(w http.ResponseWriter, status int, message string) {
	writeError(w, status, message)
}

// uintParam parses the query parameter name, absent as zero.
func uintParam(r *http.Request, name string) (uint64, error) {
	text := r.URL.Query().Get(name)
	if text == "" {
		return 0, nil
	}
	n, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		return 0, invalid("%s must be a block number", name)
	}
	return n, nil
}

// call runs a view of contract at block. A revert is an answer, decoded where an SDK binding knows it.
func (s *Server) call(ctx context.Context, contract *catalog.Contract, function string, block uint64, args ...any) (map[string]any, *Revert, error) {
	method, ok := contract.ABI.Methods[function]
	if !ok || !method.IsConstant() {
		return nil, nil, missing("%s has no view %s", contract.Name, function)
	}
	input, err := method.Inputs.Pack(args...)
	if err != nil {
		return nil, nil, invalid("%s", err)
	}
	address := contract.Address
	data, err := s.Backend.CallContract(ctx, ethereum.CallMsg{To: &address, Data: append(method.ID, input...)},
		new(big.Int).SetUint64(block))
	if err != nil {
		if revert, ok := revertOf(err); ok {
			return nil, revert, nil
		}
		return nil, nil, err
	}
	values, err := method.Outputs.UnpackValues(data)
	if err != nil {
		return nil, nil, err
	}
	return valuesOf(method.Outputs, values), nil, nil
}
