// SPDX-License-Identifier: MIT OR Apache-2.0

// Package store keeps the indexed events in SQLite, a file the indexer rebuilds from the chain alone: deleting it
// loses nothing.
package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	_ "modernc.org/sqlite" // registers the "sqlite" driver
)

const schema = `
CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT NOT NULL) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS blocks (number INTEGER PRIMARY KEY, hash BLOB NOT NULL) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS events (
	block INTEGER NOT NULL,
	log_index INTEGER NOT NULL,
	block_hash BLOB NOT NULL,
	time INTEGER NOT NULL,
	tx BLOB NOT NULL,
	contract TEXT NOT NULL,
	event TEXT NOT NULL,
	args TEXT NOT NULL,
	PRIMARY KEY (block, log_index)
) WITHOUT ROWID;
CREATE INDEX IF NOT EXISTS events_by_contract ON events (contract, event, block, log_index);
`

// Event is an indexed log, decoded by the ABI of the contract that emitted it.
type Event struct {
	Block     uint64         `json:"block"`
	LogIndex  uint           `json:"logIndex"`
	BlockHash common.Hash    `json:"blockHash"`
	Time      uint64         `json:"time"`
	Tx        common.Hash    `json:"tx"`
	Contract  string         `json:"contract"`
	Event     string         `json:"event"`
	Args      map[string]any `json:"args"`
}

// Head is the last block the index covers.
type Head struct {
	Number uint64      `json:"number"`
	Hash   common.Hash `json:"hash"`
}

// Store is the index.
type Store struct {
	db *sql.DB
}

// Open opens the index at path, creating it where there is none. It is written ahead of a log, so reads never wait
// for a write.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", "file:"+path+"?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{db}, nil
}

// Close closes the index.
func (s *Store) Close() error {
	return s.db.Close()
}

// Meta returns the value of key, empty where none is set.
func (s *Store) Meta(ctx context.Context, key string) (string, error) {
	var value string
	err := s.db.QueryRowContext(ctx, `SELECT value FROM meta WHERE key = ?`, key).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return value, err
}

// SetMeta sets key to value.
func (s *Store) SetMeta(ctx context.Context, key, value string) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO meta (key, value) VALUES (?, ?)
		ON CONFLICT (key) DO UPDATE SET value = excluded.value`, key, value)
	return err
}

// Reset deletes everything the index holds.
func (s *Store) Reset(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM events; DELETE FROM blocks; DELETE FROM meta;`)
	return err
}

// Head returns the last block the index covers; false before the first.
func (s *Store) Head(ctx context.Context) (Head, bool, error) {
	var head Head
	var hash []byte
	err := s.db.QueryRowContext(ctx, `SELECT number, hash FROM blocks ORDER BY number DESC LIMIT 1`).Scan(&head.Number, &hash)
	if errors.Is(err, sql.ErrNoRows) {
		return Head{}, false, nil
	}
	head.Hash = common.BytesToHash(hash)
	return head, err == nil, err
}

// Blocks returns the blocks the index holds a hash of above from, newest first: the blocks with events and every head.
func (s *Store) Blocks(ctx context.Context, from uint64) ([]Head, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT number, hash FROM blocks WHERE number > ? ORDER BY number DESC`, from)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []Head
	for rows.Next() {
		var head Head
		var hash []byte
		if err := rows.Scan(&head.Number, &hash); err != nil {
			return nil, err
		}
		head.Hash = common.BytesToHash(hash)
		out = append(out, head)
	}
	return out, rows.Err()
}

// Append adds events, and head as the last block the index covers, in one transaction.
func (s *Store) Append(ctx context.Context, events []Event, head Head) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	block, err := tx.PrepareContext(ctx, `INSERT OR REPLACE INTO blocks (number, hash) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	insert, err := tx.PrepareContext(ctx, `INSERT INTO events (block, log_index, block_hash, time, tx, contract, event, args)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	for _, event := range events {
		args, err := json.Marshal(event.Args)
		if err != nil {
			return err
		}
		if _, err := insert.ExecContext(ctx, event.Block, event.LogIndex, event.BlockHash.Bytes(), event.Time,
			event.Tx.Bytes(), event.Contract, event.Event, string(args)); err != nil {
			return err
		}
		if _, err := block.ExecContext(ctx, event.Block, event.BlockHash.Bytes()); err != nil {
			return err
		}
	}
	if _, err := block.ExecContext(ctx, head.Number, head.Hash.Bytes()); err != nil {
		return err
	}
	return tx.Commit()
}

// Rewind deletes every event and block hash above block, as a reorganisation replaced them.
func (s *Store) Rewind(ctx context.Context, block uint64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `DELETE FROM events WHERE block > ?`, block); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM blocks WHERE number > ?`, block); err != nil {
		return err
	}
	return tx.Commit()
}

// Query selects events in the order they were emitted.
type Query struct {
	Contract string
	Event    string
	// Args are arguments the event must carry, by name: a string in its JSON form, or 1 or 0 for a bool.
	Args map[string]any
	From uint64
	// To is the last block, none where zero.
	To uint64
	// After is the position to resume after, the Cursor of the last event of a previous page.
	After *Cursor
	Limit int
}

// Cursor is an event's position: its block and log index.
type Cursor struct {
	Block    uint64
	LogIndex uint
}

// String returns the cursor as block:logIndex.
func (c Cursor) String() string {
	return fmt.Sprintf("%d:%d", c.Block, c.LogIndex)
}

// ParseCursor parses block:logIndex.
func ParseCursor(text string) (Cursor, error) {
	block, index, ok := strings.Cut(text, ":")
	b, err1 := strconv.ParseUint(block, 10, 64)
	i, err2 := strconv.ParseUint(index, 10, 32)
	if !ok || err1 != nil || err2 != nil {
		return Cursor{}, fmt.Errorf("%q is not block:logIndex", text)
	}
	return Cursor{b, uint(i)}, nil
}

// Events returns the events q selects.
func (s *Store) Events(ctx context.Context, q Query) ([]Event, error) {
	var where []string
	var args []any
	if q.Contract != "" {
		where, args = append(where, "contract = ?"), append(args, q.Contract)
	}
	if q.Event != "" {
		where, args = append(where, "event = ?"), append(args, q.Event)
	}
	for _, name := range sortedKeys(q.Args) {
		where = append(where, "json_extract(args, ?) = ?")
		args = append(args, "$."+name, q.Args[name])
	}
	if q.From > 0 {
		where, args = append(where, "block >= ?"), append(args, q.From)
	}
	if q.To > 0 {
		where, args = append(where, "block <= ?"), append(args, q.To)
	}
	if q.After != nil {
		where, args = append(where, "(block, log_index) > (?, ?)"), append(args, q.After.Block, q.After.LogIndex)
	}
	query := `SELECT block, log_index, block_hash, time, tx, contract, event, args FROM events`
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY block, log_index LIMIT ?"
	rows, err := s.db.QueryContext(ctx, query, append(args, q.Limit)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	out := []Event{}
	for rows.Next() {
		var event Event
		var blockHash, tx []byte
		var raw string
		if err := rows.Scan(&event.Block, &event.LogIndex, &blockHash, &event.Time, &tx, &event.Contract, &event.Event,
			&raw); err != nil {
			return nil, err
		}
		event.BlockHash, event.Tx = common.BytesToHash(blockHash), common.BytesToHash(tx)
		if err := json.Unmarshal([]byte(raw), &event.Args); err != nil {
			return nil, err
		}
		out = append(out, event)
	}
	return out, rows.Err()
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

// Pairs returns each distinct pair of the arguments first and second that contract's event carried, sorted: the
// accounts' positions that borrowed, or the shorts that sold.
func (s *Store) Pairs(ctx context.Context, contract, event, first, second string) ([][2]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT json_extract(args, ?), json_extract(args, ?) FROM events
		WHERE contract = ? AND event = ? GROUP BY 1, 2 ORDER BY 1, 2`,
		"$."+first, "$."+second, contract, event)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out [][2]string
	for rows.Next() {
		var pair [2]string
		if err := rows.Scan(&pair[0], &pair[1]); err != nil {
			return nil, err
		}
		out = append(out, pair)
	}
	return out, rows.Err()
}
