// SPDX-License-Identifier: MIT OR Apache-2.0

package calibrate

import (
	"encoding/json"
	"math/big"
	"os"
	"slices"
	"testing"
)

var launchAssets = []string{"NVDA", "TSLA", "AAPL", "MSFT", "GOOGL", "SPY"}

func recorded(t *testing.T, block string) Snapshot {
	t.Helper()
	data, err := os.ReadFile("testdata/pools-" + block + ".json")
	if err != nil {
		t.Fatal(err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	return s
}

func TestTheTickWalkGivesTheLaunchDepths(t *testing.T) {
	sunday, monday := recorded(t, "74136756"), recorded(t, "75093578")
	got, ok := Depths([]Snapshot{sunday, monday}, monday.Time, launchAssets)
	want := []uint32{3561773, 1377156, 164067, 166576, 169957, 136784, 373481, 143815, 386624, 354428, 320861, 877772}
	if !ok || !slices.Equal(got, want) {
		t.Fatalf("depths %v, %v", got, ok)
	}
}

func TestEachSnapshotsDepthsAreItsOwn(t *testing.T) {
	monday := recorded(t, "75093578")
	nvda := Absorbed(monday.Pools[0], DepthMove, true)
	tsla := Absorbed(monday.Pools[1], DepthMove, true)
	if nvda.Int64() != 3561773 || tsla.Int64() != 164067 {
		t.Fatalf("NVDA selling %v, TSLA selling %v", nvda, tsla)
	}
	if got, ok := Depths([]Snapshot{monday}, monday.Time, launchAssets); ok || got != nil {
		t.Fatalf("a weekday's snapshot alone gave %v", got)
	}
}

func TestASmallerMoveAbsorbsAtLeastItsShare(t *testing.T) {
	for _, block := range []string{"74136756", "75093578"} {
		for _, p := range recorded(t, block).Pools {
			for _, selling := range []bool{true, false} {
				full := Absorbed(p, 10, selling)
				for _, pct := range []int64{1, 2, 5} {
					least := new(big.Int).Mul(full, big.NewInt(pct))
					least.Quo(least, big.NewInt(10))
					if got := Absorbed(p, int(pct), selling); got.Cmp(least) < 0 {
						t.Errorf("%s %s selling %v: %d%% absorbs %v, less than %v", block, p.Asset, selling, pct, got, least)
					}
				}
			}
		}
	}
}

func TestWithoutAWeekendSnapshotTheDepthsAreKept(t *testing.T) {
	monday := recorded(t, "75093578")
	if got, ok := Depths([]Snapshot{monday, monday}, monday.Time+86_400, launchAssets); ok || got != nil {
		t.Fatalf("depths %v, %v", got, ok)
	}
	sunday := recorded(t, "74136756")
	if _, ok := Depths([]Snapshot{sunday, monday}, sunday.Time+29*86_400, launchAssets); ok {
		t.Fatal("a weekend snapshot older than 28 days counted")
	}
}
