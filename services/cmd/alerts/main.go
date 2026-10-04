// SPDX-License-Identifier: MIT OR Apache-2.0

// Command alerts warns the holders of Tapehouse's margin positions before a weekend and before a liquidation, by email
// or Telegram as each chose. An account subscribes with a signature over POST /v1/subscriptions; the service values its
// positions as the liquidator does and sends the alerts that are due once a minute. It reads its configuration from the
// environment: TAPEHOUSE_DEPLOYMENTS, TAPEHOUSE_RPC_URL, TAPEHOUSE_INDEXER_URL, TAPEHOUSE_ALERTS_URL and
// TAPEHOUSE_ALERTS_SECRET, and optionally TAPEHOUSE_RPC_FALLBACK_URLS, TAPEHOUSE_RPC_RATE, TAPEHOUSE_ALERTS_DB,
// TAPEHOUSE_ALERTS_LISTEN, TAPEHOUSE_ALERTS_INTERVAL, TAPEHOUSE_SMTP_URL with TAPEHOUSE_ALERTS_FROM, and
// TELEGRAM_BOT_TOKEN with TELEGRAM_API_URL. A password or a token is read only from the environment and never logged.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tapehouse/tapehouse/services/internal/alerts"
	"github.com/tapehouse/tapehouse/services/internal/keeper"
	"github.com/tapehouse/tapehouse/services/internal/redact"
	"github.com/tapehouse/tapehouse/services/sdk"
	"golang.org/x/sync/errgroup"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	cfg, err := alerts.ConfigFromEnv(os.Getenv)
	if err != nil {
		log.Error("configuration", "error", err)
		os.Exit(2)
	}
	secrets := append([]string{cfg.BotToken, string(cfg.Secret)}, cfg.RPCURLs...)
	if password := smtpPassword(cfg.SMTPURL); password != "" {
		secrets = append(secrets, password)
	}
	log = slog.New(redact.Handler(log.Handler(), secrets...))
	if err := run(cfg, log); err != nil && !errors.Is(err, context.Canceled) {
		log.Error("the alerts stopped", "error", redact.String(err.Error(), secrets...))
		os.Exit(1)
	}
}

func smtpPassword(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.User == nil {
		return ""
	}
	password, _ := u.User.Password()
	return password
}

func run(cfg alerts.Config, log *slog.Logger) error {
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
	store, err := alerts.OpenStore(cfg.Database)
	if err != nil {
		return err
	}
	defer func() { _ = store.Close() }()
	links := alerts.Links{BaseURL: cfg.PublicURL, Secret: cfg.Secret}
	channels := map[string]alerts.Channel{}
	var email alerts.Channel
	var bot *alerts.Bot
	if cfg.SMTPURL != "" {
		smtp, err := alerts.NewEmail(cfg.SMTPURL, cfg.From)
		if err != nil {
			return err
		}
		email, channels[alerts.Email] = smtp, smtp
	}
	var chatbot alerts.Chatbot
	if cfg.BotToken != "" {
		bot = alerts.NewTelegram(cfg.BotToken, cfg.BotURL)
		chatbot, channels[alerts.Telegram] = bot, bot
	}
	server := alerts.NewServer(store, deployments.ChainID, links, email, chatbot, log)
	service := alerts.NewService(store, alerts.NewReader(sdk.NewClient(chain, deployments)),
		keeper.NewIndex(cfg.IndexerURL, &http.Client{Timeout: 30 * time.Second}), channels, links, server, log)
	listener := &http.Server{Addr: cfg.Listen, Handler: server.Handler(), ReadHeaderTimeout: 10 * time.Second}
	log.Info("running the alerts", "chainId", deployments.ChainID, "listen", cfg.Listen, "email", email != nil, "telegram", bot != nil)
	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error { return service.Run(ctx, cfg.Interval) })
	if bot != nil {
		group.Go(func() error { return service.Poll(ctx, bot) })
	}
	group.Go(func() error {
		if err := listener.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return listener.Shutdown(shutdown)
	})
	return group.Wait()
}
