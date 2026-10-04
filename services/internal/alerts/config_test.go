// SPDX-License-Identifier: MIT OR Apache-2.0

package alerts

import (
	"strings"
	"testing"
	"time"
)

func environment(extra map[string]string) func(string) string {
	base := map[string]string{
		"TAPEHOUSE_DEPLOYMENTS": "deployments/46630.json", "TAPEHOUSE_RPC_URL": "https://rpc.example/key", "TAPEHOUSE_INDEXER_URL": "https://index.example/",
		"TAPEHOUSE_ALERTS_URL": "https://alerts.example/", "TAPEHOUSE_ALERTS_SECRET": strings.Repeat("s", 32),
		"TAPEHOUSE_SMTP_URL": "smtp://user:password@smtp.example:587", "TAPEHOUSE_ALERTS_FROM": "alerts@tapehouse.example",
	}
	for k, v := range extra {
		base[k] = v
	}
	return func(name string) string { return base[name] }
}

func TestTheConfigurationComesFromTheEnvironment(t *testing.T) {
	cfg, err := ConfigFromEnv(environment(map[string]string{"TAPEHOUSE_RPC_FALLBACK_URLS": "https://b.example, https://c.example", "TELEGRAM_BOT_TOKEN": "token"}))
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.RPCURLs) != 3 || cfg.IndexerURL != "https://index.example" || cfg.PublicURL != "https://alerts.example" || cfg.Database != "alerts.db" ||
		cfg.Listen != "127.0.0.1:8090" || cfg.Interval != time.Minute || cfg.BotURL != "https://api.telegram.org" || cfg.RPCRate != 20 {
		t.Fatalf("%+v", cfg)
	}
	for name, extra := range map[string]map[string]string{
		"no registry":                   {"TAPEHOUSE_DEPLOYMENTS": ""},
		"no indexer":                    {"TAPEHOUSE_INDEXER_URL": ""},
		"no public address":             {"TAPEHOUSE_ALERTS_URL": "alerts.example"},
		"a short secret":                {"TAPEHOUSE_ALERTS_SECRET": "short"},
		"an SMTP server without sender": {"TAPEHOUSE_ALERTS_FROM": ""},
		"no channel":                    {"TAPEHOUSE_SMTP_URL": "", "TAPEHOUSE_ALERTS_FROM": ""},
		"a rate of zero":                {"TAPEHOUSE_RPC_RATE": "0"},
		"an interval of no time":        {"TAPEHOUSE_ALERTS_INTERVAL": "0s"},
	} {
		if _, err := ConfigFromEnv(environment(extra)); err == nil {
			t.Errorf("%s: no error", name)
		} else if strings.Contains(err.Error(), "password") || strings.Contains(err.Error(), "rpc.example") {
			t.Errorf("%s: the error names a secret: %v", name, err)
		}
	}
}
