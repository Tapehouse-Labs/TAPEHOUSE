// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// dutch stages a live auction of borrower's cross position, holding NVDA and WETH, that the indexer lists.
type dutch struct {
	*market
	ask      *big.Int
	started  uint64
	asked    []common.Address
	queries  []map[string]string
	headers  []string
	pages    [][]map[string]any
	mu       sync.Mutex
	position [32]byte
}

func newDutch(t *testing.T, ask int64) (*dutch, *Bidder) {
	t.Helper()
	m := newMarket(t)
	m.quotes["NVDA"] = m.quotes["SPY"]
	d := &dutch{market: m, ask: big.NewInt(ask), started: m.time - 60, position: cross}
	m.read("Liquidator", "auctions", func([]any) []any { return []any{d.started, false} })
	m.read("Liquidator", "price", func(args []any) []any {
		d.asked = append(d.asked, args[2].(common.Address))
		return []any{d.ask}
	})
	m.read("MarginAccounts", "collateral", func(args []any) []any {
		if args[2].(common.Address) == m.deployments.Tokens["NVDA"] {
			return []any{big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(100))}
		}
		return []any{new(big.Int)}
	})
	m.write("Liquidator", "buy", func([]any, bool) error { return nil })
	url := m.indexer(t, d.serve)
	b := m.bidder(t, statePath(t))
	b.Assets = nil
	b.IndexerURL = url
	d.pages = [][]map[string]any{{{"account": borrower.Hex(), "position": hexutil.Encode(cross[:]), "closed": false}}}
	return d, b
}

func (d *dutch) serve(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()
	query := map[string]string{"path": r.URL.Path}
	for key, values := range r.URL.Query() {
		query[key] = values[0]
	}
	d.queries = append(d.queries, query)
	d.headers = append(d.headers, r.Header.Get("X-API-Key"))
	page := 0
	if after := r.URL.Query().Get("after"); after != "" {
		page = int(after[0] - '0')
	}
	var next *string
	if page+1 < len(d.pages) {
		cursor := string(rune('0' + page + 1))
		next = &cursor
	}
	events := []map[string]any{}
	for i, args := range d.pages[page] {
		events = append(events, map[string]any{"block": 1, "logIndex": i, "args": args})
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"events": events, "next": next})
}

func (d *dutch) tick(t *testing.T, b *Bidder) {
	t.Helper()
	if err := b.Dutch(context.Background(), d.head()); err != nil {
		t.Fatal(err)
	}
}

func TestNothingIsBoughtAboveTheBotsPrice(t *testing.T) {
	d, b := newDutch(t, 99e8)
	d.tick(t, b)
	if sent := d.mined(); len(sent) != 0 {
		t.Fatalf("%v", sent)
	}
}

func TestATokenIsBoughtOnceTheAskReachesTheBotsPrice(t *testing.T) {
	d, b := newDutch(t, 95e8)
	d.tick(t, b)
	amount := new(big.Int).Div(new(big.Int).Mul(big.NewInt(1000e6), tokenToUSDG), big.NewInt(95e8))
	maxCost := Escrow(amount, big.NewInt(95e8))
	if maxCost.Cmp(big.NewInt(1000e6)) > 0 {
		t.Fatalf("the limit %v exceeds the budget", maxCost)
	}
	nvda := d.deployments.Tokens["NVDA"]
	want := []string{
		"USDG.approve" + sprint(d.names["Liquidator"], maxCost),
		"Liquidator.buy" + sprint(borrower, cross, nvda, amount, maxCost, me),
	}
	if sent := d.mined(); len(sent) != 2 || sent[0] != want[0] || sent[1] != want[1] {
		t.Fatalf("%v, want %v", sent, want)
	}
}

func TestMaxCostIsTheAsksCostRoundedUp(t *testing.T) {
	d, b := newDutch(t, 95e8)
	b.Budget = big.NewInt(1)
	d.tick(t, b)
	calls := d.took()
	if len(calls) != 2 || calls[1].args[3].(*big.Int).Int64() != 10_526_315_789 || calls[1].args[4].(*big.Int).Int64() != 1 {
		t.Fatalf("%v", calls)
	}
	if got := Escrow(big.NewInt(3), big.NewInt(1e8+1)); got.Int64() != 1 {
		t.Fatalf("amount 3 at an ask of 1e8+1: %v", got)
	}
}

func TestWETHIsNeverBought(t *testing.T) {
	d, b := newDutch(t, 95e8)
	d.read("MarginAccounts", "collateral", func(args []any) []any {
		if args[2].(common.Address) == d.deployments.Tokens["WETH"] {
			return []any{big.NewInt(1e18)}
		}
		return []any{new(big.Int)}
	})
	d.tick(t, b)
	if sent := d.mined(); len(sent) != 0 {
		t.Fatalf("%v", sent)
	}
	for _, token := range d.asked {
		if token == d.deployments.Tokens["WETH"] {
			t.Fatal("WETH was priced")
		}
	}
}

func TestAnAuctionNoLongerLiveIsSkipped(t *testing.T) {
	d, b := newDutch(t, 95e8)
	d.started = 0
	d.tick(t, b)
	d.started = d.time - 60
	d.ask = new(big.Int)
	d.tick(t, b)
	if sent := d.mined(); len(sent) != 0 {
		t.Fatalf("%v", sent)
	}
}

func TestOnlyTheLastFourHoursOfStartsAreRead(t *testing.T) {
	d, b := newDutch(t, 99e8)
	d.blockSeconds, d.block = 12, 100_000
	b.IndexerKey = "key"
	second := map[string]any{"account": other.Hex(), "position": hexutil.Encode(cross[:]), "closed": true}
	d.pages = append(d.pages, []map[string]any{second})
	d.tick(t, b)
	if len(d.queries) != 2 {
		t.Fatalf("%v", d.queries)
	}
	for i, query := range d.queries {
		if query["from"] != "98800" || query["contract"] != "tapehouse.Liquidator" || query["event"] != "AuctionStarted" ||
			query["path"] != "/v1/events" || d.headers[i] != "key" {
			t.Errorf("%v %q", query, d.headers[i])
		}
	}
	if d.queries[1]["after"] != "1" {
		t.Errorf("the next page was not followed: %v", d.queries[1])
	}
	auctions, err := liveAuctions(context.Background(), b.IndexerURL, "", 98800)
	if err != nil || len(auctions) != 2 || auctions[0].Account != borrower || auctions[1].Account != other {
		t.Fatalf("%v, %v", auctions, err)
	}
}

func TestAPurchaseTheContractRefusesIsRetriedNextTick(t *testing.T) {
	for _, name := range []string{"CostAboveLimit", "NoAuction", "NotLiquidatable"} {
		d, b := newDutch(t, 95e8)
		d.write("Liquidator", "buy", func([]any, bool) error {
			if name == "CostAboveLimit" {
				return d.fail("Liquidator", name, big.NewInt(2), big.NewInt(1))
			}
			return d.fail("Liquidator", name, borrower, cross)
		})
		d.tick(t, b)
		d.tick(t, b)
		approvals := 0
		for _, sent := range d.mined() {
			if strings.HasPrefix(sent, "USDG.approve") {
				approvals++
			}
			if strings.HasPrefix(sent, "Liquidator.buy") {
				t.Errorf("%s: a refused purchase was mined", name)
			}
		}
		if approvals != 2 || !strings.Contains(d.logs.String(), name) {
			t.Errorf("%s: %d approvals, logs %s", name, approvals, d.logs.String())
		}
	}
}
