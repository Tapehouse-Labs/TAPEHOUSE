// SPDX-License-Identifier: MIT OR Apache-2.0

// Command entrypoint starts a service in a container. Each variable NAME_B64 holds a base64 secret: it is written to
// a 0600 file under /run/secrets and replaced by the variable the services read, NAME_FILE for a NAME ending in
// _PASSWORD and NAME otherwise, so a service still reads its secrets only from files. It then runs its arguments.
package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

const secrets = "/run/secrets"

func prepare(environ []string, dir string) ([]string, error) {
	taken := map[string]bool{}
	for _, entry := range environ {
		name, _, _ := strings.Cut(entry, "=")
		taken[name] = true
	}
	out := make([]string, 0, len(environ)+2)
	for _, entry := range environ {
		name, value, _ := strings.Cut(entry, "=")
		base, ok := strings.CutSuffix(name, "_B64")
		if !ok || base == "" {
			out = append(out, entry)
			continue
		}
		target := base
		if strings.HasSuffix(base, "_PASSWORD") {
			target += "_FILE"
		}
		if taken[target] {
			return nil, fmt.Errorf("%s and %s are both set", name, target)
		}
		data, err := base64.StdEncoding.DecodeString(strings.TrimSpace(value))
		if err != nil || len(data) == 0 {
			return nil, fmt.Errorf("%s is not a non-empty base64 value", name)
		}
		path := filepath.Join(dir, base)
		if err := os.WriteFile(path, data, 0o600); err != nil {
			return nil, fmt.Errorf("%s: %w", name, errors.Unwrap(err))
		}
		out = append(out, target+"="+path)
	}
	return out, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: entrypoint command [argument ...]")
		os.Exit(2)
	}
	environ, err := prepare(os.Environ(), secrets)
	if err != nil {
		fmt.Fprintln(os.Stderr, "entrypoint:", err)
		os.Exit(2)
	}
	binary, err := exec.LookPath(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "entrypoint:", err)
		os.Exit(2)
	}
	fmt.Fprintln(os.Stderr, "entrypoint:", syscall.Exec(binary, os.Args[1:], environ))
	os.Exit(1)
}
