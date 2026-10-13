// SPDX-License-Identifier: MIT OR Apache-2.0

// Command facilitator runs an x402 facilitator for one chain: it verifies the payments that clients sign in USDG and
// settles them on the chain, sending each as a transferWithAuthorization from its own account, which pays the gas. It
// serves, to the client that holds its token, GET /supported, POST /verify and POST /settle. It reads its
// configuration from the environment: TAPEHOUSE_DEPLOYMENTS, TAPEHOUSE_RPC_URL, TAPEHOUSE_FACILITATOR_TOKEN (the bearer
// token of the one client it serves), its key as TAPEHOUSE_FACILITATOR_KEYSTORE and TAPEHOUSE_FACILITATOR_PASSWORD_FILE
// or TAPEHOUSE_FACILITATOR_KEY, and optionally TAPEHOUSE_FACILITATOR_ADDR (127.0.0.1:4402). A keystore's password is read
// only from the file TAPEHOUSE_FACILITATOR_PASSWORD_FILE names; the hex key is for a dev node.
package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	x402 "github.com/x402-foundation/x402/go/v2"

	"github.com/tapehouse/tapehouse/services/internal/payments"
	"github.com/tapehouse/tapehouse/services/internal/redact"
	"github.com/tapehouse/tapehouse/services/sdk"
)

type config struct{ deployments, rpcURL, addr, token string }

func loadConfig(getenv func(string) string) (config, error) {
	cfg := config{deployments: getenv("TAPEHOUSE_DEPLOYMENTS"), rpcURL: getenv("TAPEHOUSE_RPC_URL"),
		addr: getenv("TAPEHOUSE_FACILITATOR_ADDR"), token: getenv("TAPEHOUSE_FACILITATOR_TOKEN")}
	for _, required := range []string{"TAPEHOUSE_DEPLOYMENTS", "TAPEHOUSE_RPC_URL", "TAPEHOUSE_FACILITATOR_TOKEN"} {
		if getenv(required) == "" {
			return cfg, fmt.Errorf("set %s", required)
		}
	}
	if cfg.addr == "" {
		cfg.addr = "127.0.0.1:4402"
	}
	return cfg, nil
}

// loadKey reads the facilitator's key: a Web3 Secret Storage keystore with its password in a file, or a hex key.
func loadKey(getenv func(string) string) (*ecdsa.PrivateKey, error) {
	if path := getenv("TAPEHOUSE_FACILITATOR_KEYSTORE"); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("TAPEHOUSE_FACILITATOR_KEYSTORE: %w", err)
		}
		passwordPath := getenv("TAPEHOUSE_FACILITATOR_PASSWORD_FILE")
		if passwordPath == "" {
			return nil, errors.New("TAPEHOUSE_FACILITATOR_PASSWORD_FILE is required with the keystore")
		}
		password, err := os.ReadFile(passwordPath)
		if err != nil {
			return nil, errors.New("TAPEHOUSE_FACILITATOR_PASSWORD_FILE could not be read")
		}
		key, err := keystore.DecryptKey(data, strings.TrimSuffix(strings.TrimSuffix(string(password), "\n"), "\r"))
		if err != nil {
			return nil, fmt.Errorf("TAPEHOUSE_FACILITATOR_KEYSTORE: %w", err)
		}
		return key.PrivateKey, nil
	}
	if hex := getenv("TAPEHOUSE_FACILITATOR_KEY"); hex != "" {
		key, err := crypto.HexToECDSA(strings.TrimPrefix(hex, "0x"))
		if err != nil {
			return nil, errors.New("TAPEHOUSE_FACILITATOR_KEY is not a private key")
		}
		return key, nil
	}
	return nil, errors.New("the facilitator's key is required: TAPEHOUSE_FACILITATOR_KEYSTORE or TAPEHOUSE_FACILITATOR_KEY")
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
		log.Error("the facilitator stopped", "error", redact.String(err.Error(), cfg.rpcURL))
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
	signer, err := payments.NewSigner(ctx, chain, key)
	if err != nil {
		return err
	}
	if chainID, _ := signer.GetChainID(ctx); chainID.Uint64() != deployments.ChainID {
		return fmt.Errorf("the registry is of chain %d and the RPC of chain %s", deployments.ChainID, chainID)
	}
	usdg, ok := deployments.Tokens["USDG"]
	if !ok {
		return errors.New("the registry has no USDG to settle")
	}
	network := x402.Network(fmt.Sprintf("eip155:%d", deployments.ChainID))
	server := &http.Server{Addr: cfg.addr, Handler: payments.NewFacilitator(network, signer, usdg, cfg.token),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second, WriteTimeout: 3 * time.Minute, IdleTimeout: time.Minute}
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdown)
	}()
	log.Info("serving payments", "network", network, "addr", cfg.addr, "signer", signer.Address().Hex())
	return server.ListenAndServe()
}
