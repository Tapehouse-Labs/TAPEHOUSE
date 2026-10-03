// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestTheBidderImportsNothingInternal(t *testing.T) {
	out, err := exec.Command("go", "list", "-deps", ".").Output()
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range strings.Fields(string(out)) {
		if strings.HasPrefix(path, "github.com/tapehouse/tapehouse/services/internal") {
			t.Errorf("the bidder imports %s", path)
		}
	}
	if !strings.Contains(string(out), "github.com/tapehouse/tapehouse/services/sdk\n") {
		t.Error("the bidder does not build on the SDK")
	}
}
