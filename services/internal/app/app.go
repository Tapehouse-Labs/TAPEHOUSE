// SPDX-License-Identifier: MIT OR Apache-2.0

// Package app runs the indexer and its API together over one chain's registry.
package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapehouse/tapehouse/services/internal/api"
	"github.com/tapehouse/tapehouse/services/internal/catalog"
	"github.com/tapehouse/tapehouse/services/internal/indexer"
	"github.com/tapehouse/tapehouse/services/internal/payments"
	"github.com/tapehouse/tapehouse/services/internal/redact"
	"github.com/tapehouse/tapehouse/services/internal/redstone"
	"github.com/tapehouse/tapehouse/services/internal/store"
	"github.com/tapehouse/tapehouse/services/sdk"
	x402 "github.com/x402-foundation/x402/go/v2"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

// Config is what the indexer runs on, read from the environment by ConfigFromEnv.
type Config struct {
	// Deployments is the path of the chain's registry, deployments/<chainId>.json: TAPEHOUSE_DEPLOYMENTS.
	Deployments string
	// RPCURL is the chain's RPC endpoint, which may carry a key and is never logged or served: TAPEHOUSE_RPC_URL.
	RPCURL string
	// DB is the index's file, tapehouse-<chainId>.db unless TAPEHOUSE_DB names one.
	DB string
	// Listen is the API's address, 127.0.0.1:8080 unless TAPEHOUSE_LISTEN names one.
	Listen string
	// KeyDigests are the SHA-256 digests of the keyed tier's API keys, hex, separated by commas: TAPEHOUSE_API_KEYS.
	KeyDigests []string
	// Proxies are the networks of the proxies whose X-Forwarded-For is trusted, separated by commas:
	// TAPEHOUSE_TRUSTED_PROXIES.
	Proxies []netip.Prefix
	// ClientIPHeader is the header in which a trusted proxy puts the client's address, read in place of
	// X-Forwarded-For: TAPEHOUSE_CLIENT_IP_HEADER.
	ClientIPHeader string
	// Gateways are RedStone's gateways, the keyed ones with REDSTONE_API_KEY and REDSTONE_BACKUP_API_KEY where set.
	Gateways []redstone.Gateway
	// Interval is how often the indexer polls the chain once it has caught up.
	Interval time.Duration
	// RPCRate is the most requests a second the indexer and the API send the RPC together, 20 unless
	// TAPEHOUSE_RPC_RATE names another: providers refuse a client over its plan's rate.
	RPCRate float64
	// Payments, where set, sells the requests beyond the public allowance one at a time, in USDG, instead of a key.
	Payments *Payments
}

// Payments is how the API is paid per request, over x402.
type Payments struct {
	// Facilitator is the base URL of the x402 facilitator that verifies and settles the payments:
	// TAPEHOUSE_X402_FACILITATOR.
	Facilitator string
	// Token is the facilitator's bearer token: TAPEHOUSE_X402_FACILITATOR_TOKEN.
	Token string
	// PayTo receives the payments: TAPEHOUSE_X402_PAY_TO.
	PayTo common.Address
	// Price is what a request costs, in USDG's atomic units, 1000 (0.001 USDG) unless TAPEHOUSE_X402_PRICE names
	// another.
	Price *big.Int
}

// ConfigFromEnv reads the configuration from the environment.
func ConfigFromEnv(getenv func(string) string) (Config, error) {
	cfg := Config{
		Deployments:    getenv("TAPEHOUSE_DEPLOYMENTS"),
		RPCURL:         getenv("TAPEHOUSE_RPC_URL"),
		DB:             getenv("TAPEHOUSE_DB"),
		ClientIPHeader: strings.TrimSpace(getenv("TAPEHOUSE_CLIENT_IP_HEADER")),
		Listen:         getenv("TAPEHOUSE_LISTEN"),
		Gateways:       redstone.Gateways(getenv),
		Interval:       time.Second,
		RPCRate:        20,
	}
	if cfg.Deployments == "" || cfg.RPCURL == "" {
		return cfg, errors.New("TAPEHOUSE_DEPLOYMENTS and TAPEHOUSE_RPC_URL are required")
	}
	if text := getenv("TAPEHOUSE_RPC_RATE"); text != "" {
		rate, err := strconv.ParseFloat(text, 64)
		if err != nil || rate <= 0 {
			return cfg, errors.New("TAPEHOUSE_RPC_RATE must be a positive number of requests a second")
		}
		cfg.RPCRate = rate
	}
	if cfg.Listen == "" {
		cfg.Listen = "127.0.0.1:8080"
	}
	for _, digest := range strings.Split(getenv("TAPEHOUSE_API_KEYS"), ",") {
		if digest = strings.TrimSpace(digest); digest != "" {
			cfg.KeyDigests = append(cfg.KeyDigests, digest)
		}
	}
	for _, network := range strings.Split(getenv("TAPEHOUSE_TRUSTED_PROXIES"), ",") {
		if network = strings.TrimSpace(network); network == "" {
			continue
		}
		prefix, err := netip.ParsePrefix(network)
		if err != nil {
			return cfg, fmt.Errorf("TAPEHOUSE_TRUSTED_PROXIES: %w", err)
		}
		cfg.Proxies = append(cfg.Proxies, prefix)
	}
	if facilitator := getenv("TAPEHOUSE_X402_FACILITATOR"); facilitator != "" {
		payTo := getenv("TAPEHOUSE_X402_PAY_TO")
		if getenv("TAPEHOUSE_X402_FACILITATOR_TOKEN") == "" {
			return cfg, errors.New("TAPEHOUSE_X402_FACILITATOR_TOKEN is required with the facilitator")
		}
		if !common.IsHexAddress(payTo) {
			return cfg, errors.New("TAPEHOUSE_X402_PAY_TO must be the address that receives the payments")
		}
		cfg.Payments = &Payments{Facilitator: facilitator, Token: getenv("TAPEHOUSE_X402_FACILITATOR_TOKEN"), PayTo: common.HexToAddress(payTo), Price: big.NewInt(1000)}
		if text := getenv("TAPEHOUSE_X402_PRICE"); text != "" {
			price, ok := new(big.Int).SetString(text, 10)
			if !ok || price.Sign() <= 0 {
				return cfg, errors.New("TAPEHOUSE_X402_PRICE must be a positive number of USDG's atomic units")
			}
			cfg.Payments.Price = price
		}
	}
	return cfg, nil
}

// sell is the tiers' payment: the paywall at the price in USDG, which the facilitator settles.
func sell(ctx context.Context, chain *throttled, deployments *sdk.Deployments, pay *Payments) (*api.Payment, error) {
	usdg, ok := deployments.Tokens["USDG"]
	if !ok {
		return nil, errors.New("the registry has no USDG to be paid in")
	}
	name, version, err := payments.TokenDomain(ctx, chain, usdg)
	if err != nil {
		return nil, err
	}
	wrap, err := payments.Paywall(payments.Terms{Network: x402.Network(fmt.Sprintf("eip155:%d", deployments.ChainID)),
		Asset: usdg, Name: name, Version: version, PayTo: pay.PayTo, Price: pay.Price, Facilitator: pay.Facilitator, Token: pay.Token,
		Patterns: api.Priced()})
	if err != nil {
		return nil, err
	}
	return &api.Payment{Wrap: wrap, Priced: api.IsPriced, Attempts: api.DefaultAttempts}, nil
}

// Run indexes the registry's chain and serves the API until ctx ends. ready, where set, is called with the API's
// address once it listens. Neither its logs nor its error show the RPC URL, which may carry a key.
func Run(ctx context.Context, cfg Config, log *slog.Logger, ready func(net.Addr)) error {
	log = slog.New(redact.Handler(log.Handler(), cfg.RPCURL))
	if err := run(ctx, cfg, log, ready); err != nil {
		return errors.New(redact.String(err.Error(), cfg.RPCURL))
	}
	return nil
}

func run(ctx context.Context, cfg Config, log *slog.Logger, ready func(net.Addr)) error {
	deployments, err := sdk.LoadDeployments(cfg.Deployments)
	if err != nil {
		return err
	}
	contracts, err := catalog.New(deployments)
	if err != nil {
		return err
	}
	tiers, err := api.NewTiers(cfg.KeyDigests, cfg.Proxies)
	if err != nil {
		return err
	}
	tiers.ClientIPHeader = cfg.ClientIPHeader
	dialed, err := ethclient.DialContext(ctx, cfg.RPCURL)
	if err != nil {
		return errors.New("the RPC URL does not parse")
	}
	defer dialed.Close()
	client := &throttled{dialed, rate.NewLimiter(rate.Limit(cfg.RPCRate), 1)}
	if cfg.Payments != nil {
		if tiers.Payment, err = sell(ctx, client, deployments, cfg.Payments); err != nil {
			return err
		}
	}
	if cfg.DB == "" {
		cfg.DB = fmt.Sprintf("tapehouse-%d.db", deployments.ChainID)
	}
	index, err := store.Open(cfg.DB)
	if err != nil {
		return err
	}
	defer func() { _ = index.Close() }()
	hub := api.NewHub()
	follower := indexer.New(client, index, contracts, hub.Publish, log)
	if err := follower.Prepare(ctx); err != nil {
		return err
	}
	server := api.New(api.Config{Deployments: deployments, Catalog: contracts, Store: index, Index: follower,
		Backend: client, Relay: redstone.New(cfg.Gateways, &http.Client{Timeout: 10 * time.Second}), Hub: hub,
		Tiers: tiers, Log: log})
	if cfg.Payments != nil {
		log.Info("selling the requests beyond the public allowance", "price", cfg.Payments.Price.String(),
			"payTo", cfg.Payments.PayTo.Hex(), "chainId", deployments.ChainID)
	}
	listener, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		return err
	}
	httpServer := &http.Server{Handler: server.Handler(), ReadHeaderTimeout: 10 * time.Second}
	log.Info("serving the API", "address", listener.Addr().String(), "chainId", deployments.ChainID,
		"contracts", len(contracts.Contracts), "start", follower.Start())
	if ready != nil {
		ready(listener.Addr())
	}
	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error { return follower.Run(ctx, cfg.Interval) })
	group.Go(func() error {
		if err := httpServer.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpServer.Shutdown(shutdown)
	})
	if err := group.Wait(); !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

// throttled is the RPC client, holding the indexer's and the API's requests together under the RPC's rate.
type throttled struct {
	*ethclient.Client
	limiter *rate.Limiter
}

func (t *throttled) ChainID(ctx context.Context) (*big.Int, error) {
	if err := t.limiter.Wait(ctx); err != nil {
		return nil, err
	}
	return t.Client.ChainID(ctx)
}

func (t *throttled) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	if err := t.limiter.Wait(ctx); err != nil {
		return nil, err
	}
	return t.Client.HeaderByNumber(ctx, number)
}

func (t *throttled) CodeAt(ctx context.Context, account common.Address, block *big.Int) ([]byte, error) {
	if err := t.limiter.Wait(ctx); err != nil {
		return nil, err
	}
	return t.Client.CodeAt(ctx, account, block)
}

func (t *throttled) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	if err := t.limiter.Wait(ctx); err != nil {
		return nil, err
	}
	return t.Client.FilterLogs(ctx, q)
}

func (t *throttled) CallContract(ctx context.Context, call ethereum.CallMsg, block *big.Int) ([]byte, error) {
	if err := t.limiter.Wait(ctx); err != nil {
		return nil, err
	}
	return t.Client.CallContract(ctx, call, block)
}
