// SPDX-License-Identifier: MIT OR Apache-2.0

package app

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"
)

func TestTheLogsNeverShowTheRPCURL(t *testing.T) {
	var logs bytes.Buffer
	secret := "https://rpc.example/secret-key"
	log := slog.New(redacting{slog.NewJSONHandler(&logs, nil), secret}).With("rpc", secret).WithGroup("step")
	log.Warn("indexing step failed at "+secret, "error", errors.New(`Post "`+secret+`": EOF`), "block", 7)
	if strings.Contains(logs.String(), "secret") || strings.Count(logs.String(), "the RPC") != 3 ||
		!strings.Contains(logs.String(), `"block":"7"`) {
		t.Fatalf("logs: %s", logs.String())
	}
}
