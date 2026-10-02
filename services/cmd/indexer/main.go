// SPDX-License-Identifier: MIT OR Apache-2.0

// Command indexer indexes Tapehouse's contracts on the chain of a registry and serves them over REST and a WebSocket
// stream, with RedStone's signed packages relayed beside them. It reads its configuration from the environment:
// TAPEHOUSE_DEPLOYMENTS, TAPEHOUSE_RPC_URL, and optionally TAPEHOUSE_DB, TAPEHOUSE_LISTEN, TAPEHOUSE_API_KEYS,
// TAPEHOUSE_TRUSTED_PROXIES, REDSTONE_API_KEY and REDSTONE_BACKUP_API_KEY.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/tapehouse/tapehouse/services/internal/app"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	cfg, err := app.ConfigFromEnv(os.Getenv)
	if err != nil {
		log.Error("configuration", "error", err)
		os.Exit(2)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err = app.Run(ctx, cfg, log, nil)
	stop()
	if err != nil {
		log.Error("the indexer stopped", "error", err)
		os.Exit(1)
	}
}
