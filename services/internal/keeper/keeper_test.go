// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tapehouse/tapehouse/services/sdk"
)

func env(values map[string]string) func(string) string {
	return func(name string) string { return values[name] }
}

func TestTheConfigurationComesFromTheEnvironment(t *testing.T) {
	cfg, err := ConfigFromEnv(env(map[string]string{
		"TAPEHOUSE_DEPLOYMENTS": "deployments/4663.json", "TAPEHOUSE_RPC_URL": "https://a.example/key",
		"TAPEHOUSE_RPC_FALLBACK_URLS": " https://b.example/key , https://c.example,", "TAPEHOUSE_KEEPERS": "prices,halt",
		"TAPEHOUSE_KEEPER_BUY": "true", "REDSTONE_API_KEY": "redstone",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(cfg.RPCURLs, []string{"https://a.example/key", "https://b.example/key", "https://c.example"}) ||
		!slices.Equal(cfg.Keepers, []string{"prices", "halt"}) || cfg.RPCRate != 20 || cfg.HaltsURL != HaltsURL ||
		!cfg.Buy || cfg.Rebalance || cfg.Gateways[0].Key != "redstone" {
		t.Fatalf("configuration %+v", cfg)
	}
	for values, want := range map[*map[string]string]string{
		{"TAPEHOUSE_DEPLOYMENTS": "d"}: "TAPEHOUSE_DEPLOYMENTS and TAPEHOUSE_RPC_URL are required",
		{"TAPEHOUSE_DEPLOYMENTS": "d", "TAPEHOUSE_RPC_URL": "r", "TAPEHOUSE_KEEPERS": "prices,sweep"}: `no keeper "sweep"`,
		{"TAPEHOUSE_DEPLOYMENTS": "d", "TAPEHOUSE_RPC_URL": "r", "TAPEHOUSE_RPC_RATE": "0"}:           "TAPEHOUSE_RPC_RATE",
	} {
		if _, err := ConfigFromEnv(env(*values)); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%v: %v, want %q", *values, err, want)
		}
	}
}

func TestTheKeysComeFromAKeystoreWithItsPasswordFileOrAHexKeyNeverAnArgument(t *testing.T) {
	signer, err := keystore.StoreKey(t.TempDir(), "secret", keystore.LightScryptN, keystore.LightScryptP)
	if err != nil {
		t.Fatal(err)
	}
	path := signer.URL.Path
	passwordFile := filepath.Join(t.TempDir(), "password")
	if err := os.WriteFile(passwordFile, []byte("secret\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	keeper, halt, err := Keys(env(map[string]string{
		"TAPEHOUSE_KEEPER_KEY":    "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcaf784d7bf4f2ff80",
		"TAPEHOUSE_HALT_KEYSTORE": path, "TAPEHOUSE_HALT_PASSWORD_FILE": passwordFile,
	}))
	if err != nil || !keeper.Equal(keeperKey) || crypto.PubkeyToAddress(halt.PublicKey) != signer.Address {
		t.Fatalf("keys: %v", err)
	}
	if _, halt, err := Keys(env(map[string]string{"TAPEHOUSE_KEEPER_KEY": "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcaf784d7bf4f2ff80"})); err != nil || halt != nil {
		t.Fatalf("a keeper without a halt signer: %v", err)
	}
	for values, want := range map[*map[string]string]string{
		{}:                                  "the keepers' key is required: TAPEHOUSE_KEEPER_KEYSTORE or TAPEHOUSE_KEEPER_KEY",
		{"TAPEHOUSE_KEEPER_KEY": "0x01"}:    "TAPEHOUSE_KEEPER_KEY is not a private key",
		{"TAPEHOUSE_KEEPER_KEYSTORE": path}: "TAPEHOUSE_KEEPER_PASSWORD_FILE is required with the keystore",
	} {
		if _, _, err := Keys(env(*values)); err == nil || err.Error() != want {
			t.Errorf("%v: %v, want %q", *values, err, want)
		}
	}
}

func TestAKeystorePasswordComesFromAFileAndNeverFromTheEnvironment(t *testing.T) {
	signer, err := keystore.StoreKey(t.TempDir(), "secret", keystore.LightScryptN, keystore.LightScryptP)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	write := func(name, content string) string {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	right, wrong, empty := write("right", "secret\n"), write("wrong", "hunter2-wrong\n"), write("empty", "\n")
	for _, role := range []string{"KEEPER", "HALT"} {
		prefix := "TAPEHOUSE_" + role + "_"
		values := map[string]string{prefix + "KEYSTORE": signer.URL.Path, prefix + "PASSWORD_FILE": right}
		if role == "HALT" {
			values["TAPEHOUSE_KEEPER_KEY"] = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcaf784d7bf4f2ff80"
		}
		keeper, halt, err := Keys(env(values))
		got := keeper
		if role == "HALT" {
			got = halt
		}
		if err != nil || got == nil || crypto.PubkeyToAddress(got.PublicKey) != signer.Address {
			t.Fatalf("%s: a password file with a trailing newline: %v", role, err)
		}
		for values, want := range map[*map[string]string]string{
			{prefix + "KEYSTORE": signer.URL.Path}:                                                    prefix + "PASSWORD_FILE is required",
			{prefix + "KEYSTORE": signer.URL.Path, prefix + "PASSWORD": "secret"}:                     prefix + "PASSWORD_FILE is required",
			{prefix + "KEYSTORE": signer.URL.Path, prefix + "PASSWORD_FILE": filepath.Join(dir, "x")}: prefix + "PASSWORD_FILE could not be read",
			{prefix + "KEYSTORE": signer.URL.Path, prefix + "PASSWORD_FILE": empty}:                   prefix + "PASSWORD_FILE is empty",
			{prefix + "KEYSTORE": signer.URL.Path, prefix + "PASSWORD_FILE": wrong}:                   prefix + "KEYSTORE: could not decrypt key with given password",
		} {
			(*values)["TAPEHOUSE_KEEPER_KEY"] = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcaf784d7bf4f2ff80"
			_, _, err := Keys(env(*values))
			if err == nil || !strings.Contains(err.Error(), want) || strings.Contains(err.Error(), "hunter2") || strings.Contains(err.Error(), "secret") {
				t.Errorf("%v: %v, want %q", *values, err, want)
			}
		}
	}
}

func TestTheKeepersRunAreThoseWhoseContractsTheRegistryNames(t *testing.T) {
	arbitrum, err := sdk.LoadDeployments("../../../deployments/42161.json")
	if err != nil {
		t.Fatal(err)
	}
	k, err := New(Config{}, arbitrum, newWorld(t), keeperKey, nil, slog.New(slog.DiscardHandler))
	if err != nil || !slices.Equal(k.Keepers(), []string{"prices", "multiplier", "halt", "seal"}) {
		t.Fatalf("on 42161: %v, %v", k, err)
	}
	w := newWorld(t)
	k, err = New(Config{}, w.deployments, w, keeperKey, nil, slog.New(slog.DiscardHandler))
	if err == nil || err.Error() != "the sync keeper needs TAPEHOUSE_INDEXER_URL" {
		t.Fatalf("index keepers without the index: %v, %v", k, err)
	}
	if k = w.keeper(Config{}); !slices.Equal(k.Keepers(), []string{"prices", "multiplier", "halt", "seal", "premium",
		"mark", "sync", "liquidation", "backstop", "reopening", "recall", "shorts", "gapcover"}) {
		t.Fatalf("every keeper: %v", k.Keepers())
	}
	if err := k.Pass(context.Background(), "basket"); err != nil {
		t.Fatalf("a basket pass with no basket: %v", err)
	}
}
