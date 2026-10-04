// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func encode(text string) string { return base64.StdEncoding.EncodeToString([]byte(text)) }

func TestPrepareWritesPasswordsAndKeystores(t *testing.T) {
	dir := t.TempDir()
	environ := []string{
		"PATH=/usr/bin",
		"TAPEHOUSE_KEEPER_PASSWORD_B64=" + encode("hunter2\n"),
		"TAPEHOUSE_KEEPER_KEYSTORE_B64=" + encode(`{"version":3}`),
		"TAPEHOUSE_RPC_URL=https://rpc.example",
	}
	out, err := prepare(environ, dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"PATH=/usr/bin",
		"TAPEHOUSE_RPC_URL=https://rpc.example",
		"TAPEHOUSE_KEEPER_PASSWORD_FILE=" + filepath.Join(dir, "TAPEHOUSE_KEEPER_PASSWORD"),
		"TAPEHOUSE_KEEPER_KEYSTORE=" + filepath.Join(dir, "TAPEHOUSE_KEEPER_KEYSTORE"),
	} {
		if !slices.Contains(out, want) {
			t.Errorf("environment lacks %q: %v", want, out)
		}
	}
	for _, entry := range out {
		if strings.Contains(entry, "_B64=") {
			t.Errorf("the encoded variable survives: %q", entry)
		}
	}
	for name, content := range map[string]string{"TAPEHOUSE_KEEPER_PASSWORD": "hunter2\n", "TAPEHOUSE_KEEPER_KEYSTORE": `{"version":3}`} {
		path := filepath.Join(dir, name)
		got, err := os.ReadFile(path)
		if err != nil || string(got) != content {
			t.Errorf("%s holds %q, %v", name, got, err)
		}
		info, _ := os.Stat(path)
		if info.Mode().Perm() != 0o600 {
			t.Errorf("%s has mode %v", name, info.Mode().Perm())
		}
	}
}

func TestPrepareWithoutSecretsChangesNothing(t *testing.T) {
	environ := []string{"A=1", "B=2"}
	out, err := prepare(environ, t.TempDir())
	if err != nil || !slices.Equal(out, environ) {
		t.Fatalf("got %v, %v", out, err)
	}
}

func TestPrepareRejectsBadInput(t *testing.T) {
	secret := "c2VjcmV0"
	for name, environ := range map[string][]string{
		"invalid base64": {"TAPEHOUSE_SPONSOR_PASSWORD_B64=%%%not-base64%%%"},
		"empty":          {"TAPEHOUSE_SPONSOR_PASSWORD_B64="},
		"target set":     {"TAPEHOUSE_SPONSOR_PASSWORD_B64=" + secret, "TAPEHOUSE_SPONSOR_PASSWORD_FILE=/elsewhere"},
	} {
		_, err := prepare(environ, t.TempDir())
		if err == nil {
			t.Errorf("%s: no error", name)
			continue
		}
		if strings.Contains(err.Error(), "not-base64") || strings.Contains(err.Error(), secret) {
			t.Errorf("%s: the error leaks the value: %v", name, err)
		}
	}
}
