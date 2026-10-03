// SPDX-License-Identifier: MIT OR Apache-2.0

// Package redact keeps secrets, such as an RPC URL that carries a provider's key, out of logs and errors.
package redact

import (
	"context"
	"log/slog"
	"strings"
)

// String returns text with each non-empty secret written as "the RPC".
func String(text string, secrets ...string) string {
	for _, secret := range secrets {
		if secret != "" {
			text = strings.ReplaceAll(text, secret, "the RPC")
		}
	}
	return text
}

// Handler returns a log handler that writes each secret as "the RPC" wherever a message or an attribute carries it.
func Handler(h slog.Handler, secrets ...string) slog.Handler {
	return redacting{h, secrets}
}

type redacting struct {
	slog.Handler
	secrets []string
}

func (h redacting) Handle(ctx context.Context, record slog.Record) error {
	out := slog.NewRecord(record.Time, record.Level, String(record.Message, h.secrets...), record.PC)
	record.Attrs(func(attr slog.Attr) bool {
		out.AddAttrs(h.attr(attr))
		return true
	})
	return h.Handler.Handle(ctx, out)
}

func (h redacting) WithAttrs(attrs []slog.Attr) slog.Handler {
	redacted := make([]slog.Attr, len(attrs))
	for i, attr := range attrs {
		redacted[i] = h.attr(attr)
	}
	return redacting{h.Handler.WithAttrs(redacted), h.secrets}
}

func (h redacting) attr(attr slog.Attr) slog.Attr {
	return slog.String(attr.Key, String(attr.Value.Resolve().String(), h.secrets...))
}

func (h redacting) WithGroup(name string) slog.Handler {
	return redacting{h.Handler.WithGroup(name), h.secrets}
}
