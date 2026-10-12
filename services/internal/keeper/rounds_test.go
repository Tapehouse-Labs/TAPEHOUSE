// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"errors"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/roundfeed"
)

// roundWorld holds TSLA and NVDA, a gap cover whose feed for TSLA is a RedStone round feed, last updated at `updated`,
// and for NVDA a Chainlink feed, which has no data feed ID.
func roundWorld(t *testing.T, updated func() int64) (*world, *[]string) {
	t.Helper()
	w := newWorld(t)
	w.asset("TSLA")
	w.asset("NVDA")
	round := w.add("RoundFeed:TSLA", &roundfeed.RoundFeedMetaData)
	chainlink := w.add("Chainlink:NVDA", &roundfeed.RoundFeedMetaData)
	w.write("Chainlink:NVDA", "dataFeedId", func([]any, bool) error { return errors.New("execution reverted") })
	w.read("Band", "asset", func(args []any) []any {
		name := text(args[0].([32]byte))
		return []any{common.Address{}, symbol(t, name+"---24_7"), [32]byte{}, w.deployments.Tokens[name]}
	})
	w.read("GapCover", "feed", func(args []any) []any {
		if text(args[0].([32]byte)) == "TSLA" {
			return []any{round}
		}
		return []any{chainlink}
	})
	w.read("RoundFeed:TSLA", "dataFeedId", func([]any) []any { return []any{symbol(t, "TSLA")} })
	w.read("RoundFeed:TSLA", "latestRoundData", func([]any) []any {
		return []any{n(1), n(38_000_000_000), n(updated()), n(updated()), n(1)}
	})
	var updates []string
	w.write("RoundFeed:TSLA", "update", func(_ []any, mined bool) error {
		if mined {
			updates = append(updates, "TSLA")
		}
		return nil
	})
	return w, &updates
}

func TestARoundFeedIsRecordedWhileNYSEIsInRegularHoursAndItsLastRoundIsOld(t *testing.T) {
	var w *world
	w, updates := roundWorld(t, func() int64 { return int64(w.time) - 120 })
	w.session(&session{state: sessionOpen, nyse: 1})
	k := w.keeper(Config{})
	feeds := &source{}
	k.prices = feeds
	if got := w.pass(t, k, "rounds"); len(got) != 1 || len(feeds.asked) != 1 || feeds.asked[0][0] != "TSLA" || !strings.HasPrefix(got[0], "RoundFeed.update") || len(*updates) != 1 {
		t.Fatalf("an old round in regular hours: %v", got)
	}
}

func TestARoundFeedWaitsForAFreshRoundAndForTheRegularHours(t *testing.T) {
	var w *world
	updated := func() int64 { return int64(w.time) - 10 }
	w, _ = roundWorld(t, updated)
	w.session(&session{state: sessionOpen, nyse: 1})
	k := w.keeper(Config{})
	if got := w.pass(t, k, "rounds"); len(got) != 0 {
		t.Fatalf("a round recorded ten seconds ago: %v", got)
	}

	w, _ = roundWorld(t, func() int64 { return int64(w.time) - 3_600 })
	w.session(&session{state: sessionClosed, nyse: 3})
	k = w.keeper(Config{})
	if got := w.pass(t, k, "rounds"); len(got) != 0 {
		t.Fatalf("a closed market: %v", got)
	}
}

func TestAChainlinkFeedIsLeftToItsOwnOracleAndAStaleRoundRecordedTwiceIsNotAnError(t *testing.T) {
	var w *world
	w, _ = roundWorld(t, func() int64 { return int64(w.time) - 200 })
	w.session(&session{state: sessionOpen, nyse: 1})
	w.write("RoundFeed:TSLA", "update", func([]any, bool) error { return w.fail("RoundFeed:TSLA", "NotNewer", n(1), n(1)) })
	k := w.keeper(Config{})
	k.prices = &source{}
	if got := w.pass(t, k, "rounds"); len(got) != 0 {
		t.Fatalf("a round another keeper recorded first: %v", got)
	}
	if strings.Contains(w.logs.String(), "NVDA") {
		t.Fatalf("a Chainlink feed was touched: %s", w.logs.String())
	}
}
