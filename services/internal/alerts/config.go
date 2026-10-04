// SPDX-License-Identifier: MIT OR Apache-2.0

package alerts

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Config is the service's configuration from the environment.
type Config struct {
	Deployments string
	RPCURLs     []string
	RPCRate     float64
	IndexerURL  string
	Database    string
	Listen      string
	PublicURL   string
	Secret      []byte
	SMTPURL     string
	From        string
	BotToken    string
	BotURL      string
	Interval    time.Duration
}

// ConfigFromEnv reads the configuration: TAPEHOUSE_DEPLOYMENTS, TAPEHOUSE_RPC_URL, TAPEHOUSE_INDEXER_URL,
// TAPEHOUSE_ALERTS_URL (the address links in emails point to) and TAPEHOUSE_ALERTS_SECRET (at least 32 characters) are
// required; TAPEHOUSE_RPC_FALLBACK_URLS, TAPEHOUSE_RPC_RATE, TAPEHOUSE_ALERTS_DB (alerts.db), TAPEHOUSE_ALERTS_LISTEN
// (127.0.0.1:8090), TAPEHOUSE_ALERTS_INTERVAL (60s), TAPEHOUSE_SMTP_URL with TAPEHOUSE_ALERTS_FROM, and
// TELEGRAM_BOT_TOKEN with TELEGRAM_API_URL (https://api.telegram.org) open the email and Telegram channels.
func ConfigFromEnv(getenv func(string) string) (Config, error) {
	cfg := Config{
		Deployments: getenv("TAPEHOUSE_DEPLOYMENTS"),
		IndexerURL:  strings.TrimRight(getenv("TAPEHOUSE_INDEXER_URL"), "/"),
		PublicURL:   strings.TrimRight(getenv("TAPEHOUSE_ALERTS_URL"), "/"),
		Secret:      []byte(getenv("TAPEHOUSE_ALERTS_SECRET")),
		Database:    or(getenv("TAPEHOUSE_ALERTS_DB"), "alerts.db"),
		Listen:      or(getenv("TAPEHOUSE_ALERTS_LISTEN"), "127.0.0.1:8090"),
		SMTPURL:     getenv("TAPEHOUSE_SMTP_URL"),
		From:        getenv("TAPEHOUSE_ALERTS_FROM"),
		BotToken:    getenv("TELEGRAM_BOT_TOKEN"),
		BotURL:      or(getenv("TELEGRAM_API_URL"), "https://api.telegram.org"),
		RPCRate:     20,
		Interval:    time.Minute,
	}
	if primary := getenv("TAPEHOUSE_RPC_URL"); primary != "" {
		cfg.RPCURLs = append(cfg.RPCURLs, primary)
	}
	for _, u := range strings.Split(getenv("TAPEHOUSE_RPC_FALLBACK_URLS"), ",") {
		if u = strings.TrimSpace(u); u != "" {
			cfg.RPCURLs = append(cfg.RPCURLs, u)
		}
	}
	if cfg.Deployments == "" || len(cfg.RPCURLs) == 0 || cfg.IndexerURL == "" {
		return cfg, errors.New("TAPEHOUSE_DEPLOYMENTS, TAPEHOUSE_RPC_URL and TAPEHOUSE_INDEXER_URL are required")
	}
	if u, err := url.Parse(cfg.PublicURL); err != nil || u.Host == "" || (u.Scheme != "https" && u.Scheme != "http") {
		return cfg, errors.New("TAPEHOUSE_ALERTS_URL must be the http(s) address the links in emails point to")
	}
	if len(cfg.Secret) < 32 {
		return cfg, errors.New("TAPEHOUSE_ALERTS_SECRET must be at least 32 characters")
	}
	if (cfg.SMTPURL == "") != (cfg.From == "") {
		return cfg, errors.New("TAPEHOUSE_SMTP_URL and TAPEHOUSE_ALERTS_FROM go together")
	}
	if cfg.SMTPURL == "" && cfg.BotToken == "" {
		return cfg, errors.New("set TAPEHOUSE_SMTP_URL with TAPEHOUSE_ALERTS_FROM, or TELEGRAM_BOT_TOKEN: a channel to send by")
	}
	if text := getenv("TAPEHOUSE_RPC_RATE"); text != "" {
		rate, err := strconv.ParseFloat(text, 64)
		if err != nil || rate <= 0 {
			return cfg, errors.New("TAPEHOUSE_RPC_RATE must be a positive number of requests a second")
		}
		cfg.RPCRate = rate
	}
	if text := getenv("TAPEHOUSE_ALERTS_INTERVAL"); text != "" {
		interval, err := time.ParseDuration(text)
		if err != nil || interval < time.Second {
			return cfg, errors.New("TAPEHOUSE_ALERTS_INTERVAL must be a duration of at least a second")
		}
		cfg.Interval = interval
	}
	return cfg, nil
}

func or(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
