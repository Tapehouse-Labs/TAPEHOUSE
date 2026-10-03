// SPDX-License-Identifier: MIT OR Apache-2.0

// Command keeper runs Tapehouse's keepers on the chain of a registry: RedStone's prices, multiplier syncs, signed
// halts, seals, the premium's closures, the shorts' marks, burns and dust, liquidations, the backstop, the reopening
// auction, recalls and buy-ins, short liquidations, basket rebalances and the gap cover's settlement. It reads its
// configuration from the environment: TAPEHOUSE_DEPLOYMENTS, TAPEHOUSE_RPC_URL, its key as TAPEHOUSE_KEEPER_KEYSTORE
// and TAPEHOUSE_KEEPER_PASSWORD_FILE or TAPEHOUSE_KEEPER_KEY, and optionally TAPEHOUSE_RPC_FALLBACK_URLS,
// TAPEHOUSE_RPC_RATE, TAPEHOUSE_INDEXER_URL, TAPEHOUSE_KEEPERS, the halt signer's TAPEHOUSE_HALT_KEYSTORE and
// TAPEHOUSE_HALT_PASSWORD_FILE or TAPEHOUSE_HALT_KEY, TAPEHOUSE_HALTS_URL, TAPEHOUSE_KEEPER_BUY,
// TAPEHOUSE_KEEPER_REBALANCE, REDSTONE_API_KEY and REDSTONE_BACKUP_API_KEY. A keystore's password is read only from
// the file its *_PASSWORD_FILE names, never from the environment itself; the hex keys are for a dev node.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/tapehouse/tapehouse/services/internal/keeper"
	"github.com/tapehouse/tapehouse/services/internal/redact"
	"github.com/tapehouse/tapehouse/services/sdk"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	cfg, err := keeper.ConfigFromEnv(os.Getenv)
	if err != nil {
		log.Error("configuration", "error", err)
		os.Exit(2)
	}
	log = slog.New(redact.Handler(log.Handler(), cfg.RPCURLs...))
	if err := run(cfg, log); err != nil {
		log.Error("the keepers stopped", "error", redact.String(err.Error(), cfg.RPCURLs...))
		os.Exit(1)
	}
}

func run(cfg keeper.Config, log *slog.Logger) error {
	key, haltKey, err := keeper.Keys(os.Getenv)
	if err != nil {
		return err
	}
	deployments, err := sdk.LoadDeployments(cfg.Deployments)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	chain, err := keeper.Dial(ctx, cfg.RPCURLs, cfg.RPCRate, log)
	if err != nil {
		return err
	}
	defer chain.Close()
	keepers, err := keeper.New(cfg, deployments, chain, key, haltKey, log)
	if err != nil {
		return err
	}
	log.Info("running the keepers", "chainId", deployments.ChainID, "keepers", keepers.Keepers(),
		"from", keepers.From().Hex(), "halts", haltKey != nil)
	return keepers.Run(ctx)
}
