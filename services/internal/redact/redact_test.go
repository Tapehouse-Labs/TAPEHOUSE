// SPDX-License-Identifier: MIT OR Apache-2.0

package redact_test

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/tapehouse/tapehouse/services/internal/redact"
)

func TestTheLogsNeverShowTheRPCURL(t *testing.T) {
	var logs bytes.Buffer
	secret := "https://rpc.example/secret-key"
	log := slog.New(redact.Handler(slog.NewJSONHandler(&logs, nil), secret)).With("rpc", secret).WithGroup("step")
	log.Warn("indexing step failed at "+secret, "error", errors.New(`Post "`+secret+`": EOF`), "block", 7)
	if strings.Contains(logs.String(), "secret") || strings.Count(logs.String(), "the RPC") != 3 ||
		!strings.Contains(logs.String(), `"block":"7"`) {
		t.Fatalf("logs: %s", logs.String())
	}
}

func TestEverySecretIsRedacted(t *testing.T) {
	text := redact.String("https://a.example/k1 then https://b.example/k2", "https://a.example/k1", "", "https://b.example/k2")
	if text != "the RPC then the RPC" {
		t.Fatalf("%q", text)
	}
}
