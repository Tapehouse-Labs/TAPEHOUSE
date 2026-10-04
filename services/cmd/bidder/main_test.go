// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

const testKey = "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"

func env(values map[string]string) func(string) string {
	return func(name string) string { return values[name] }
}

func TestTheKeyComesFromTheEnvironmentOnly(t *testing.T) {
	var stderr bytes.Buffer
	if code := run(context.Background(), env(nil), &stderr); code != 2 || !strings.Contains(stderr.String(), "TAPEHOUSE_BIDDER_KEYSTORE or PRIVATE_KEY") {
		t.Fatalf("%d: %s", code, stderr.String())
	}
	stderr.Reset()
	code := run(context.Background(), env(map[string]string{"PRIVATE_KEY": "0x" + testKey}), &stderr)
	if code != 2 || strings.Contains(stderr.String(), testKey) || !strings.Contains(stderr.String(), "TAPEHOUSE_RPC_URL") {
		t.Fatalf("%d: %s", code, stderr.String())
	}
	stderr.Reset()
	if code := run(context.Background(), env(map[string]string{"PRIVATE_KEY": "not a key"}), &stderr); code != 2 ||
		strings.Contains(stderr.String(), "not a key") {
		t.Fatalf("%d: %s", code, stderr.String())
	}
}

func TestTheKeyIsAKeystoreOrHex(t *testing.T) {
	want, _ := crypto.HexToECDSA(testKey)
	dir := t.TempDir()
	store := keystore.NewKeyStore(dir, keystore.LightScryptN, keystore.LightScryptP)
	account, err := store.ImportECDSA(want, "secret")
	if err != nil {
		t.Fatal(err)
	}
	write := func(name, content string) string {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	right, wrong := write("right", "secret\n"), write("wrong", "other\n")
	for name, test := range map[string]struct {
		values map[string]string
		err    string
	}{
		"keystore and password file": {map[string]string{"TAPEHOUSE_BIDDER_KEYSTORE": account.URL.Path,
			"TAPEHOUSE_BIDDER_PASSWORD_FILE": right}, ""},
		"hex fallback": {map[string]string{"PRIVATE_KEY": "0x" + testKey}, ""},
		"missing password file": {map[string]string{"TAPEHOUSE_BIDDER_KEYSTORE": account.URL.Path,
			"TAPEHOUSE_BIDDER_PASSWORD_FILE": filepath.Join(dir, "missing")}, "TAPEHOUSE_BIDDER_PASSWORD_FILE could not be read"},
		"password not given": {map[string]string{"TAPEHOUSE_BIDDER_KEYSTORE": account.URL.Path},
			"TAPEHOUSE_BIDDER_PASSWORD_FILE is required with the keystore"},
		"wrong password": {map[string]string{"TAPEHOUSE_BIDDER_KEYSTORE": account.URL.Path,
			"TAPEHOUSE_BIDDER_PASSWORD_FILE": wrong}, "TAPEHOUSE_BIDDER_KEYSTORE: could not decrypt key with given password"},
		"bad hex": {map[string]string{"PRIVATE_KEY": "nope"}, "PRIVATE_KEY is not a private key"},
		"neither": {map[string]string{}, "the signer's key is required: TAPEHOUSE_BIDDER_KEYSTORE or PRIVATE_KEY"},
	} {
		key, err := loadKey(env(test.values))
		switch {
		case test.err == "" && (err != nil || !key.Equal(want)):
			t.Errorf("%s: %v", name, err)
		case test.err != "" && (err == nil || err.Error() != test.err):
			t.Errorf("%s: got %v, want %q", name, err, test.err)
		}
	}
}

func TestTheConfigurationHasItsDefaults(t *testing.T) {
	cfg, err := loadConfig(env(map[string]string{"TAPEHOUSE_RPC_URL": "http://rpc", "TAPEHOUSE_DEPLOYMENTS": "d.json"}))
	if err != nil || cfg.budget.Int64() != 1000e6 || cfg.discountBps != 300 || cfg.commitLead != 15*time.Minute ||
		len(cfg.assets) != 0 || cfg.dryRun || cfg.indexerURL != "" {
		t.Fatalf("%+v, %v", cfg, err)
	}
	cfg, err = loadConfig(env(map[string]string{"TAPEHOUSE_RPC_URL": "http://rpc", "TAPEHOUSE_DEPLOYMENTS": "d.json",
		"TAPEHOUSE_BIDDER_BUDGET": "50000000", "TAPEHOUSE_BIDDER_DISCOUNT_BPS": "100", "TAPEHOUSE_BIDDER_COMMIT_LEAD": "5m",
		"TAPEHOUSE_BIDDER_ASSETS": "SPY, NVDA", "TAPEHOUSE_BIDDER_DRY_RUN": "true", "TAPEHOUSE_INDEXER_URL": "http://i/"}))
	if err != nil || cfg.budget.Int64() != 50e6 || cfg.discountBps != 100 || cfg.commitLead != 5*time.Minute ||
		strings.Join(cfg.assets, ",") != "SPY,NVDA" || !cfg.dryRun || cfg.indexerURL != "http://i" {
		t.Fatalf("%+v, %v", cfg, err)
	}
	for _, bad := range []map[string]string{
		{"TAPEHOUSE_BIDDER_BUDGET": "-1"}, {"TAPEHOUSE_BIDDER_BUDGET": "many"}, {"TAPEHOUSE_BIDDER_DISCOUNT_BPS": "10000"},
		{"TAPEHOUSE_BIDDER_COMMIT_LEAD": "soon"},
	} {
		bad["TAPEHOUSE_RPC_URL"], bad["TAPEHOUSE_DEPLOYMENTS"] = "http://rpc", "d.json"
		if _, err := loadConfig(env(bad)); err == nil {
			t.Errorf("%v loaded", bad)
		}
	}
}

func TestTheLogsNeverShowTheRPCURLOrTheIndexerKey(t *testing.T) {
	var logs bytes.Buffer
	secret, key := "https://rpc.example/secret-key", "indexer-secret"
	log := newLogger(&logs, secret, key)
	log.Warn("failed at "+secret, "error", errors.New(`Post "`+secret+`": EOF`), "key", key, "block", 7)
	if strings.Contains(logs.String(), "secret") || strings.Count(logs.String(), "redacted") != 3 || !strings.Contains(logs.String(), `"block":7`) {
		t.Fatalf("%s", logs.String())
	}
}
