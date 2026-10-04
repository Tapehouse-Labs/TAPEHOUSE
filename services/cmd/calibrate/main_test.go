// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"bytes"
	"context"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tapehouse/tapehouse/services/internal/backtest"
	"github.com/tapehouse/tapehouse/services/internal/calibrate"
)

func TestARatioIsTheShareOfTheFallTheReopenMade(t *testing.T) {
	if r, ok := ratioOf(100, 90, 95); !ok || math.Abs(r-0.5) > 1e-12 {
		t.Fatalf("%v, %v", r, ok)
	}
	if r, ok := ratioOf(100, 90, 85); !ok || math.Abs(r-1.5) > 1e-12 {
		t.Fatalf("%v, %v", r, ok)
	}
	if _, ok := ratioOf(100, 110, 105); ok {
		t.Fatal("a rise has a ratio")
	}
}

func TestTheMarketIsCutAtItsLastSession(t *testing.T) {
	day := func(d int) time.Time { return time.Date(2026, 9, d, 13, 30, 0, 0, time.UTC) }
	m := backtest.Market{At: []time.Time{day(1), day(2), day(3)}, Open: [][]float64{{1}, {2}, {3}}, Close: [][]float64{{1}, {2}, {3}}, Adj: [][]float64{{1}, {2}, {3}}}
	got := until(m, day(2).Add(10*time.Hour))
	if len(got.At) != 2 || len(got.Adj) != 2 || got.Adj[1][0] != 2 {
		t.Fatalf("%+v", got)
	}
}

func TestSnapshotsAreReadFromTheDirectory(t *testing.T) {
	dir := t.TempDir()
	if got, err := loadSnapshots(filepath.Join(dir, "none")); err != nil || len(got) != 0 {
		t.Fatalf("%v, %v", got, err)
	}
	if err := writeJSON(filepath.Join(dir, "depth-7.json"), calibrate.Snapshot{Block: 7, Time: 9}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "proposal.json"), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := loadSnapshots(dir)
	if err != nil || len(got) != 1 || got[0].Block != 7 || got[0].Time != 9 {
		t.Fatalf("%v, %v", got, err)
	}
}

func TestTheCommandsAreChecked(t *testing.T) {
	none := func(string) string { return "" }
	for _, c := range []struct {
		args []string
		want string
	}{
		{nil, "usage"},
		{[]string{"launch"}, "unknown command"},
		{[]string{"apply"}, "usage"},
		{[]string{"snapshot"}, "TAPEHOUSE_RPC_URL"},
		{[]string{"propose"}, "TAPEHOUSE_RPC_URL"},
	} {
		err := run(context.Background(), c.args, none, &bytes.Buffer{})
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%v: %v, want %q", c.args, err, c.want)
		}
	}
}
