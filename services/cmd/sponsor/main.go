// SPDX-License-Identifier: MIT OR Apache-2.0

// Command sponsor runs Tapehouse's sponsor service: it serves ERC-7677's pm_getPaymasterStubData and
// pm_getPaymasterData over JSON-RPC, deciding which user operations the registry's sponsor paymaster pays for and
// signing them with the paymaster's signer. It reads its configuration from the environment: TAPEHOUSE_DEPLOYMENTS,
// TAPEHOUSE_RPC_URL, its key as TAPEHOUSE_SPONSOR_KEYSTORE and TAPEHOUSE_SPONSOR_PASSWORD_FILE or TAPEHOUSE_SPONSOR_KEY,
// and optionally TAPEHOUSE_SPONSOR_ADDR (127.0.0.1:4338), TAPEHOUSE_SPONSOR_DAILY (the signatures each client gets a
// day, 30) and TAPEHOUSE_TRUSTED_PROXIES (the networks of the proxies whose X-Forwarded-For is trusted, separated by
// commas). A keystore's password is read only from the file TAPEHOUSE_SPONSOR_PASSWORD_FILE names; the hex key is for
// a dev node.
package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/netip"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapehouse/tapehouse/services/internal/redact"
	"github.com/tapehouse/tapehouse/services/internal/sponsor"
	"github.com/tapehouse/tapehouse/services/sdk"
)

type config struct {
	deployments, rpcURL, addr string
	limits                    sponsor.Limits
	proxies                   []netip.Prefix
}

func loadConfig(getenv func(string) string) (config, error) {
	cfg := config{deployments: getenv("TAPEHOUSE_DEPLOYMENTS"), rpcURL: getenv("TAPEHOUSE_RPC_URL"),
		addr: getenv("TAPEHOUSE_SPONSOR_ADDR"), limits: sponsor.DefaultLimits}
	for _, required := range []string{"TAPEHOUSE_DEPLOYMENTS", "TAPEHOUSE_RPC_URL"} {
		if getenv(required) == "" {
			return cfg, fmt.Errorf("set %s", required)
		}
	}
	if cfg.addr == "" {
		cfg.addr = "127.0.0.1:4338"
	}
	if text := getenv("TAPEHOUSE_SPONSOR_DAILY"); text != "" {
		daily, err := strconv.Atoi(text)
		if err != nil || daily <= 0 {
			return cfg, errors.New("TAPEHOUSE_SPONSOR_DAILY is not a positive integer")
		}
		cfg.limits.Daily = daily
	}
	for _, text := range strings.Split(getenv("TAPEHOUSE_TRUSTED_PROXIES"), ",") {
		if text = strings.TrimSpace(text); text == "" {
			continue
		}
		prefix, err := netip.ParsePrefix(text)
		if err != nil {
			return cfg, fmt.Errorf("TAPEHOUSE_TRUSTED_PROXIES: %q is not a network", text)
		}
		cfg.proxies = append(cfg.proxies, prefix)
	}
	return cfg, nil
}

// loadKey reads the signer's key: a Web3 Secret Storage keystore with its password in a file, or a hex key.
func loadKey(getenv func(string) string) (*ecdsa.PrivateKey, error) {
	if path := getenv("TAPEHOUSE_SPONSOR_KEYSTORE"); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("TAPEHOUSE_SPONSOR_KEYSTORE: %w", err)
		}
		passwordPath := getenv("TAPEHOUSE_SPONSOR_PASSWORD_FILE")
		if passwordPath == "" {
			return nil, errors.New("TAPEHOUSE_SPONSOR_PASSWORD_FILE is required with the keystore")
		}
		password, err := os.ReadFile(passwordPath)
		if err != nil {
			return nil, errors.New("TAPEHOUSE_SPONSOR_PASSWORD_FILE could not be read")
		}
		key, err := keystore.DecryptKey(data, strings.TrimSuffix(strings.TrimSuffix(string(password), "\n"), "\r"))
		if err != nil {
			return nil, fmt.Errorf("TAPEHOUSE_SPONSOR_KEYSTORE: %w", err)
		}
		return key.PrivateKey, nil
	}
	if hex := getenv("TAPEHOUSE_SPONSOR_KEY"); hex != "" {
		key, err := crypto.HexToECDSA(strings.TrimPrefix(hex, "0x"))
		if err != nil {
			return nil, errors.New("TAPEHOUSE_SPONSOR_KEY is not a private key")
		}
		return key, nil
	}
	return nil, errors.New("the signer's key is required: TAPEHOUSE_SPONSOR_KEYSTORE or TAPEHOUSE_SPONSOR_KEY")
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	cfg, err := loadConfig(os.Getenv)
	if err != nil {
		log.Error("configuration", "error", err)
		os.Exit(2)
	}
	key, err := loadKey(os.Getenv)
	if err != nil {
		log.Error("configuration", "error", err)
		os.Exit(2)
	}
	log = slog.New(redact.Handler(log.Handler(), cfg.rpcURL))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err = run(ctx, cfg, key, log)
	stop()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("the sponsor service stopped", "error", redact.String(err.Error(), cfg.rpcURL))
		os.Exit(1)
	}
}

func run(ctx context.Context, cfg config, key *ecdsa.PrivateKey, log *slog.Logger) error {
	deployments, err := sdk.LoadDeployments(cfg.deployments)
	if err != nil {
		return err
	}
	chain, err := ethclient.DialContext(ctx, cfg.rpcURL)
	if err != nil {
		return errors.New("TAPEHOUSE_RPC_URL does not parse")
	}
	defer chain.Close()
	chainID, err := chain.ChainID(ctx)
	if err != nil {
		return err
	}
	if chainID.Uint64() != deployments.ChainID {
		return fmt.Errorf("the registry is of chain %d and the RPC of chain %s", deployments.ChainID, chainID)
	}
	service, err := sponsor.New(ctx, chain, deployments, key, cfg.limits, cfg.proxies, log)
	if err != nil {
		return err
	}
	server := &http.Server{Addr: cfg.addr, Handler: service, ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout: 10 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: time.Minute}
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdown)
	}()
	log.Info("serving sponsorships", "chainId", deployments.ChainID, "addr", cfg.addr, "signer", service.Signer().Hex(),
		"daily", cfg.limits.Daily)
	return server.ListenAndServe()
}
