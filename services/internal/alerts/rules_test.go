// SPDX-License-Identifier: MIT OR Apache-2.0

package alerts

import (
	"context"
	"path/filepath"
	"reflect"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// at is an observation of a position whose equity is ratio thousandths of its requirement.
func at(ratio int64) Observation {
	return Observation{Judged: true, Held: true, Equity: n(ratio * 1_000_000), Requirement: n(1_000 * 1_000_000)}
}

func kinds(t *testing.T, threshold int, state *PositionState, o Observation) []Kind {
	t.Helper()
	got, next := Decide(threshold, *state, o)
	*state = next
	return got
}

func TestARatioOscillatingAroundTheThresholdAlertsOnceUntilItRearms(t *testing.T) {
	var state PositionState
	for step, test := range []struct {
		ratio int64
		want  []Kind
	}{
		{1300, nil}, {1249, []Kind{Threshold}}, {1100, nil}, {1251, nil}, {1300, nil}, {1249, nil}, {1312, nil},
		{1313, nil}, {1249, []Kind{Threshold}}, {1249, nil},
	} {
		if got := kinds(t, 1250, &state, at(test.ratio)); !slices.Equal(got, test.want) {
			t.Fatalf("step %d, ratio %d: %v, want %v", step, test.ratio, got, test.want)
		}
	}
}

func TestAPositionWithoutRequirementHasNoRatio(t *testing.T) {
	var state PositionState
	o := Observation{Judged: true, Equity: n(5), Requirement: n(0)}
	if got := kinds(t, 1250, &state, o); got != nil {
		t.Fatalf("%v", got)
	}
}

func TestOneWeekendAlertIsSentForEachClosure(t *testing.T) {
	var state PositionState
	o := at(2000)
	o.Closing, o.Closure = true, 1_000
	if got := kinds(t, 1250, &state, o); !slices.Equal(got, []Kind{Weekend}) {
		t.Fatalf("the ramp's first tick: %v", got)
	}
	for range 3 {
		if got := kinds(t, 1250, &state, o); got != nil {
			t.Fatalf("a second weekend alert for one closure: %v", got)
		}
	}
	o.Closing = false
	kinds(t, 1250, &state, o)
	o.Closing, o.Closure = true, 2_000
	if got := kinds(t, 1250, &state, o); !slices.Equal(got, []Kind{Weekend}) {
		t.Fatalf("the next closure: %v", got)
	}
	var nothing PositionState
	empty := Observation{Judged: true, Equity: n(0), Requirement: n(0), Closing: true, Closure: 3_000}
	if got := kinds(t, 1250, &nothing, empty); got != nil {
		t.Fatalf("a position that holds no Stock Token has no weekend: %v", got)
	}
}

func TestOneLiquidationAlertIsSentForEachAuction(t *testing.T) {
	var state PositionState
	short := at(900)
	short.Short = true
	if got := kinds(t, 1000, &state, short); !slices.Equal(got, []Kind{Liquidation}) {
		t.Fatalf("short now: %v", got)
	}
	short.Auction = 5_000
	if got := kinds(t, 1000, &state, short); got != nil {
		t.Fatalf("the keeper's start of the auction already warned of: %v", got)
	}
	if got := kinds(t, 1000, &state, short); got != nil {
		t.Fatalf("%v", got)
	}
	short.Auction = 9_000
	if got := kinds(t, 1000, &state, short); !slices.Equal(got, []Kind{Liquidation}) {
		t.Fatalf("a new auction: %v", got)
	}
	recovered := at(1500)
	kinds(t, 1000, &state, recovered)
	if got := kinds(t, 1000, &state, short); !slices.Equal(got, []Kind{Liquidation}) {
		t.Fatalf("short again after recovering: %v", got)
	}
}

func TestAPositionNobodyCanJudgeGetsOneNoticeAndNothingElse(t *testing.T) {
	var state PositionState
	o := Observation{Closing: true, Closure: 7}
	if got := kinds(t, 1250, &state, o); !slices.Equal(got, []Kind{Unjudged}) {
		t.Fatalf("%v", got)
	}
	if got := kinds(t, 1250, &state, o); got != nil {
		t.Fatalf("a second notice for one halt: %v", got)
	}
	if got := kinds(t, 1250, &state, at(2000)); got != nil {
		t.Fatalf("judged again: %v", got)
	}
	if got := kinds(t, 1250, &state, o); !slices.Equal(got, []Kind{Unjudged}) {
		t.Fatalf("halted again: %v", got)
	}
}

func TestTheAlertStateSurvivesTheStore(t *testing.T) {
	store, err := OpenStore(filepath.Join(t.TempDir(), "alerts.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = store.Close() }()
	ctx := context.Background()
	id, err := store.Subscribe(ctx, common.Address{1}, Email, "ada@example.org", 1250, Links{Secret: []byte("x")}.digestOf)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Confirm(ctx, id, ""); err != nil {
		t.Fatal(err)
	}
	want := map[string]PositionState{"0x01": {Alerted: true, Closure: 9, Short: true, Auction: 4, Unjudged: true}}
	if err := store.SaveState(ctx, id, want); err != nil {
		t.Fatal(err)
	}
	got, ok, err := store.Get(ctx, id)
	if err != nil || !ok || !reflect.DeepEqual(got.State, want) {
		t.Fatalf("%+v, %v", got.State, err)
	}
}
