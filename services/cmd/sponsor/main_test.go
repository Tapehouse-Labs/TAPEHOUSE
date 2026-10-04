// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

const testKey = "4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"

func env(values map[string]string) func(string) string {
	return func(name string) string { return values[name] }
}

func TestTheConfigurationIsReadFromTheEnvironment(t *testing.T) {
	base := map[string]string{"TAPEHOUSE_DEPLOYMENTS": "d.json", "TAPEHOUSE_RPC_URL": "http://127.0.0.1:8547"}
	cfg, err := loadConfig(env(base))
	if err != nil || cfg.addr != "127.0.0.1:4338" || cfg.limits.Daily != 30 || len(cfg.proxies) != 0 {
		t.Fatalf("defaults: %+v, %v", cfg, err)
	}
	base["TAPEHOUSE_SPONSOR_ADDR"], base["TAPEHOUSE_SPONSOR_DAILY"] = ":9000", "5"
	base["TAPEHOUSE_TRUSTED_PROXIES"] = "10.0.0.0/8, 192.168.0.0/16"
	if cfg, err = loadConfig(env(base)); err != nil || cfg.addr != ":9000" || cfg.limits.Daily != 5 || len(cfg.proxies) != 2 {
		t.Fatalf("configured: %+v, %v", cfg, err)
	}
	for name, value := range map[string]string{"TAPEHOUSE_SPONSOR_DAILY": "0", "TAPEHOUSE_TRUSTED_PROXIES": "nowhere",
		"TAPEHOUSE_RPC_URL": ""} {
		bad := map[string]string{}
		for k, v := range base {
			bad[k] = v
		}
		bad[name] = value
		if _, err := loadConfig(env(bad)); err == nil {
			t.Fatalf("%s=%q was taken", name, value)
		}
	}
}

func TestTheKeyIsAKeystoreOrHex(t *testing.T) {
	want, _ := crypto.HexToECDSA(testKey)
	key, err := loadKey(env(map[string]string{"TAPEHOUSE_SPONSOR_KEY": "0x" + testKey}))
	if err != nil || !key.Equal(want) {
		t.Fatalf("hex: %v", err)
	}
	dir := t.TempDir()
	store := keystore.NewKeyStore(dir, keystore.LightScryptN, keystore.LightScryptP)
	account, err := store.ImportECDSA(want, "secret")
	if err != nil {
		t.Fatal(err)
	}
	password := filepath.Join(dir, "password")
	if err := os.WriteFile(password, []byte("secret\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	key, err = loadKey(env(map[string]string{"TAPEHOUSE_SPONSOR_KEYSTORE": account.URL.Path,
		"TAPEHOUSE_SPONSOR_PASSWORD_FILE": password}))
	if err != nil || !key.Equal(want) {
		t.Fatalf("keystore: %v", err)
	}
	for want, values := range map[string]map[string]string{
		"the signer's key is required: TAPEHOUSE_SPONSOR_KEYSTORE or TAPEHOUSE_SPONSOR_KEY": {},
		"TAPEHOUSE_SPONSOR_KEY is not a private key":                                        {"TAPEHOUSE_SPONSOR_KEY": "nope"},
		"TAPEHOUSE_SPONSOR_PASSWORD_FILE is required with the keystore":                     {"TAPEHOUSE_SPONSOR_KEYSTORE": account.URL.Path},
		"TAPEHOUSE_SPONSOR_PASSWORD_FILE could not be read": {"TAPEHOUSE_SPONSOR_KEYSTORE": account.URL.Path,
			"TAPEHOUSE_SPONSOR_PASSWORD_FILE": filepath.Join(dir, "missing")},
	} {
		if _, err := loadKey(env(values)); err == nil || err.Error() != want {
			t.Fatalf("got %v, want %q", err, want)
		}
	}
}
