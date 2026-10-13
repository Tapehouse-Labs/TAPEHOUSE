// SPDX-License-Identifier: MIT OR Apache-2.0

package app_test

import (
	"context"
	"log/slog"
	"net/netip"
	"strings"
	"testing"

	"github.com/tapehouse/tapehouse/services/internal/app"
)

func env(values map[string]string) func(string) string {
	return func(name string) string { return values[name] }
}

func TestTheConfigurationComesFromTheEnvironment(t *testing.T) {
	cfg, err := app.ConfigFromEnv(env(map[string]string{
		"TAPEHOUSE_DEPLOYMENTS": "deployments/4663.json", "TAPEHOUSE_RPC_URL": "https://rpc.example/key",
		"TAPEHOUSE_API_KEYS": " aa , bb,", "TAPEHOUSE_TRUSTED_PROXIES": "10.0.0.0/8, 192.168.0.0/16",
		"REDSTONE_API_KEY": "redstone",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Listen != "127.0.0.1:8080" || cfg.DB != "" || len(cfg.KeyDigests) != 2 || cfg.KeyDigests[1] != "bb" ||
		cfg.Proxies[1] != netip.MustParsePrefix("192.168.0.0/16") || cfg.Gateways[0].Key != "redstone" ||
		len(cfg.Gateways) != 3 {
		t.Fatalf("configuration %+v", cfg)
	}
	if _, err := app.ConfigFromEnv(env(map[string]string{"TAPEHOUSE_DEPLOYMENTS": "deployments/4663.json"})); err == nil {
		t.Fatal("no RPC URL was taken")
	}
	if _, err := app.ConfigFromEnv(env(map[string]string{"TAPEHOUSE_DEPLOYMENTS": "d", "TAPEHOUSE_RPC_URL": "r",
		"TAPEHOUSE_TRUSTED_PROXIES": "10.0.0.1"})); err == nil {
		t.Fatal("a proxy that is not a network was taken")
	}
}

func TestTheErrorDoesNotShowTheRPCURL(t *testing.T) {
	log := slog.New(slog.DiscardHandler)
	secret := "http://127.0.0.1:9/secret-key"
	err := app.Run(context.Background(), app.Config{Deployments: "../../../deployments/46630.json", RPCURL: secret,
		DB: t.TempDir() + "/index.db", Listen: "127.0.0.1:0"}, log, nil)
	if err == nil || strings.Contains(err.Error(), "secret") || !strings.Contains(err.Error(), "the RPC") {
		t.Fatalf("an unreachable RPC: %v", err)
	}
	if err := app.Run(context.Background(), app.Config{Deployments: "missing.json", RPCURL: secret}, log, nil); err == nil {
		t.Fatal("a missing registry ran")
	}
}

func TestTheRequestsBeyondTheAllowanceAreSoldOnlyWhereAFacilitatorIsNamed(t *testing.T) {
	base := map[string]string{"TAPEHOUSE_DEPLOYMENTS": "deployments/46630.json", "TAPEHOUSE_RPC_URL": "https://rpc.example/key"}
	with := func(extra map[string]string) map[string]string {
		merged := map[string]string{}
		for name, value := range base {
			merged[name] = value
		}
		for name, value := range extra {
			merged[name] = value
		}
		return merged
	}
	cfg, err := app.ConfigFromEnv(env(base))
	if err != nil || cfg.Payments != nil {
		t.Fatalf("without a facilitator: %+v %v", cfg.Payments, err)
	}
	cfg, err = app.ConfigFromEnv(env(with(map[string]string{"TAPEHOUSE_X402_FACILITATOR": "http://facilitator.internal:4402",
		"TAPEHOUSE_X402_FACILITATOR_TOKEN": "a token", "TAPEHOUSE_X402_PAY_TO": "0x8A2631c8226D7EA5612798e6E8C34F4DC703d0ac"})))
	if err != nil || cfg.Payments == nil || cfg.Payments.Price.Int64() != 1000 || cfg.Payments.Token != "a token" ||
		cfg.Payments.PayTo.Hex() != "0x8A2631c8226D7EA5612798e6E8C34F4DC703d0ac" {
		t.Fatalf("a facilitator: %+v %v", cfg.Payments, err)
	}
	cfg, err = app.ConfigFromEnv(env(with(map[string]string{"TAPEHOUSE_X402_FACILITATOR": "http://f", "TAPEHOUSE_X402_PRICE": "250",
		"TAPEHOUSE_X402_FACILITATOR_TOKEN": "a token", "TAPEHOUSE_X402_PAY_TO": "0x8A2631c8226D7EA5612798e6E8C34F4DC703d0ac"})))
	if err != nil || cfg.Payments.Price.Int64() != 250 {
		t.Fatalf("a price: %+v %v", cfg.Payments, err)
	}
	for name, extra := range map[string]map[string]string{
		"no payee":  {"TAPEHOUSE_X402_FACILITATOR": "http://f", "TAPEHOUSE_X402_FACILITATOR_TOKEN": "a token"},
		"no token":  {"TAPEHOUSE_X402_FACILITATOR": "http://f", "TAPEHOUSE_X402_PAY_TO": "0x8A2631c8226D7EA5612798e6E8C34F4DC703d0ac"},
		"bad price": {"TAPEHOUSE_X402_FACILITATOR": "http://f", "TAPEHOUSE_X402_FACILITATOR_TOKEN": "a token", "TAPEHOUSE_X402_PRICE": "free", "TAPEHOUSE_X402_PAY_TO": "0x8A2631c8226D7EA5612798e6E8C34F4DC703d0ac"},
		"no cost":   {"TAPEHOUSE_X402_FACILITATOR": "http://f", "TAPEHOUSE_X402_FACILITATOR_TOKEN": "a token", "TAPEHOUSE_X402_PRICE": "0", "TAPEHOUSE_X402_PAY_TO": "0x8A2631c8226D7EA5612798e6E8C34F4DC703d0ac"},
	} {
		if _, err := app.ConfigFromEnv(env(with(extra))); err == nil {
			t.Fatalf("%s was taken", name)
		}
	}
}
