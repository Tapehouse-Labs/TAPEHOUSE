// SPDX-License-Identifier: MIT OR Apache-2.0

package store_test

import (
	"context"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/internal/store"
)

func open(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.Open(filepath.Join(t.TempDir(), "index.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func event(block uint64, index uint, contract, name string, args map[string]any) store.Event {
	return store.Event{Block: block, LogIndex: index, BlockHash: common.BigToHash(common.Big1), Time: 1_790_000_000 + block,
		Tx: common.BigToHash(common.Big2), Contract: contract, Event: name, Args: args}
}

func seed(t *testing.T, s *store.Store) []store.Event {
	t.Helper()
	events := []store.Event{
		event(10, 0, "tapehouse.MarginAccounts", "Borrow", map[string]any{"account": "0xA", "position": "0x01", "assets": "5"}),
		event(10, 3, "tapehouse.Band", "HaltWritten", map[string]any{"symbol": "0x4e", "halted": true}),
		event(12, 1, "tapehouse.MarginAccounts", "Borrow", map[string]any{"account": "0xB", "position": "0x01", "assets": "7"}),
		event(15, 0, "tapehouse.MarginAccounts", "Borrow", map[string]any{"account": "0xA", "position": "0x01", "assets": "9"}),
		event(15, 2, "tapehouse.Band", "HaltWritten", map[string]any{"symbol": "0x4e", "halted": false}),
	}
	if err := s.Append(context.Background(), events[:2], store.Head{Number: 11, Hash: common.HexToHash("0x11")}); err != nil {
		t.Fatal(err)
	}
	if err := s.Append(context.Background(), events[2:], store.Head{Number: 16, Hash: common.HexToHash("0x16")}); err != nil {
		t.Fatal(err)
	}
	return events
}

func TestEventsComeBackInTheOrderTheyWereEmittedAndPageByCursor(t *testing.T) {
	s := open(t)
	ctx := context.Background()
	events := seed(t, s)
	all, err := s.Events(ctx, store.Query{Limit: 100})
	if err != nil || !reflect.DeepEqual(all, events) {
		t.Fatalf("all: %v, %v", all, err)
	}
	page, _ := s.Events(ctx, store.Query{Limit: 2})
	cursor := store.Cursor{Block: page[1].Block, LogIndex: page[1].LogIndex}
	parsed, err := store.ParseCursor(cursor.String())
	if err != nil || parsed != cursor || cursor.String() != "10:3" {
		t.Fatalf("cursor %v, %v", parsed, err)
	}
	rest, _ := s.Events(ctx, store.Query{After: &parsed, Limit: 100})
	if !reflect.DeepEqual(append(page, rest...), events) {
		t.Fatalf("pages: %v then %v", page, rest)
	}
	for _, bad := range []string{"10", "a:1", "1:b", ""} {
		if _, err := store.ParseCursor(bad); err == nil {
			t.Errorf("%q parsed as a cursor", bad)
		}
	}
}

func TestEventsFilterByContractEventBlocksAndArguments(t *testing.T) {
	s := open(t)
	ctx := context.Background()
	events := seed(t, s)
	for _, c := range []struct {
		query store.Query
		want  []store.Event
	}{
		{store.Query{Contract: "tapehouse.Band"}, []store.Event{events[1], events[4]}},
		{store.Query{Contract: "tapehouse.MarginAccounts", Event: "Borrow", Args: map[string]any{"account": "0xA"}},
			[]store.Event{events[0], events[3]}},
		{store.Query{Contract: "tapehouse.Band", Event: "HaltWritten", Args: map[string]any{"halted": 1}},
			[]store.Event{events[1]}},
		{store.Query{From: 11, To: 14}, []store.Event{events[2]}},
		{store.Query{Event: "Repay"}, []store.Event{}},
	} {
		c.query.Limit = 100
		got, err := s.Events(ctx, c.query)
		if err != nil || !reflect.DeepEqual(got, c.want) {
			t.Errorf("%+v: %v, %v; want %v", c.query, got, err, c.want)
		}
	}
	pairs, err := s.Pairs(ctx, "tapehouse.MarginAccounts", "Borrow", "account", "position")
	if err != nil || !reflect.DeepEqual(pairs, [][2]string{{"0xA", "0x01"}, {"0xB", "0x01"}}) {
		t.Fatalf("pairs %v, %v", pairs, err)
	}
}

func TestTheHeadIsTheLastBlockAppendedAndARewindDropsWhatFollows(t *testing.T) {
	s := open(t)
	ctx := context.Background()
	if _, ok, err := s.Head(ctx); ok || err != nil {
		t.Fatalf("an empty index has a head: %v", err)
	}
	events := seed(t, s)
	head, ok, err := s.Head(ctx)
	if !ok || err != nil || head != (store.Head{Number: 16, Hash: common.HexToHash("0x16")}) {
		t.Fatalf("head %v, %v, %v", head, ok, err)
	}
	blocks, err := s.Blocks(ctx, 10)
	if err != nil || len(blocks) != 4 || blocks[0].Number != 16 || blocks[3].Number != 11 {
		t.Fatalf("blocks above 10: %v, %v", blocks, err)
	}
	if err := s.Rewind(ctx, 12); err != nil {
		t.Fatal(err)
	}
	head, _, _ = s.Head(ctx)
	all, _ := s.Events(ctx, store.Query{Limit: 100})
	if head.Number != 12 || !reflect.DeepEqual(all, events[:3]) {
		t.Fatalf("after the rewind: head %v, events %v", head, all)
	}
}

func TestMetaHoldsKeysUntilAReset(t *testing.T) {
	s := open(t)
	ctx := context.Background()
	if value, err := s.Meta(ctx, "fingerprint"); value != "" || err != nil {
		t.Fatalf("unset: %q, %v", value, err)
	}
	if err := s.SetMeta(ctx, "fingerprint", "a"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetMeta(ctx, "fingerprint", "b"); err != nil {
		t.Fatal(err)
	}
	if value, _ := s.Meta(ctx, "fingerprint"); value != "b" {
		t.Fatalf("fingerprint %q", value)
	}
	seed(t, s)
	if err := s.Reset(ctx); err != nil {
		t.Fatal(err)
	}
	value, _ := s.Meta(ctx, "fingerprint")
	all, _ := s.Events(ctx, store.Query{Limit: 100})
	if _, ok, _ := s.Head(ctx); ok || value != "" || len(all) != 0 {
		t.Fatalf("after a reset: %q, %v", value, all)
	}
}

func TestAnIndexReopensWithWhatItHeld(t *testing.T) {
	path := filepath.Join(t.TempDir(), "index.db")
	s, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	events := seed(t, s)
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()
	all, err := s.Events(context.Background(), store.Query{Limit: 100})
	if err != nil || !reflect.DeepEqual(all, events) {
		t.Fatalf("reopened: %v, %v", all, err)
	}
}
