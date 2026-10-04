// SPDX-License-Identifier: MIT OR Apache-2.0

package alerts

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	_ "modernc.org/sqlite" // registers the "sqlite" driver
)

const schema = `
CREATE TABLE IF NOT EXISTS subscriptions (
	id TEXT PRIMARY KEY,
	account TEXT NOT NULL,
	channel TEXT NOT NULL,
	destination TEXT NOT NULL,
	threshold INTEGER NOT NULL,
	confirmed INTEGER NOT NULL,
	state TEXT NOT NULL,
	unsubscribe_digest BLOB NOT NULL UNIQUE,
	UNIQUE (account, channel)
) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS pending (
	digest BLOB PRIMARY KEY,
	subscription_id TEXT NOT NULL,
	kind TEXT NOT NULL,
	expires INTEGER NOT NULL
) WITHOUT ROWID;
`

// The channels a subscriber may choose.
const (
	Email    = "email"
	Telegram = "telegram"
)

// Subscription is what an alert needs to know of its subscriber: the account it watches, the channel and destination
// it goes to, the ratio below which it warns, and the state of each position's alerts. A threshold is in thousandths.
type Subscription struct {
	ID          string
	Account     common.Address
	Channel     string
	Destination string
	Threshold   int
	Confirmed   bool
	State       map[string]PositionState
}

// Store keeps the subscriptions in SQLite, and the codes and tokens that confirm them until they are used or expire.
type Store struct {
	db *sql.DB
}

// OpenStore opens the subscriptions at path, creating them where there are none.
func OpenStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite", "file:"+path+"?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(schema); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{db}, nil
}

// Close closes the store.
func (s *Store) Close() error {
	return s.db.Close()
}

// Subscribe records an unconfirmed subscription of account to channel and returns its ID. A subscription of the same
// account and channel is replaced: its destination, its threshold and its confirmation are those of the new request,
// and its alert state and codes are cleared. digest is the SHA-256 of the token that unsubscribes it, for a new row.
func (s *Store) Subscribe(ctx context.Context, account common.Address, channel, destination string, threshold int,
	digest func(id string) [32]byte) (string, error) {
	var raw [8]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	candidate := hex.EncodeToString(raw[:])
	unsubscribe := digest(candidate)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer func() { _ = tx.Rollback() }()
	var id string
	err = tx.QueryRowContext(ctx, `INSERT INTO subscriptions (id, account, channel, destination, threshold, confirmed, state, unsubscribe_digest)
		VALUES (?, ?, ?, ?, ?, 0, '{}', ?)
		ON CONFLICT (account, channel) DO UPDATE SET destination = excluded.destination, threshold = excluded.threshold,
			confirmed = 0, state = '{}'
		RETURNING id`, candidate, account.Hex(), channel, destination, threshold, unsubscribe[:]).Scan(&id)
	if err != nil {
		return "", err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM pending WHERE subscription_id = ?`, id); err != nil {
		return "", err
	}
	return id, tx.Commit()
}

// AddPending keeps the code or token whose SHA-256 is digest until expires: it confirms subscription id, as kind
// Email or Telegram.
func (s *Store) AddPending(ctx context.Context, id, kind string, digest [32]byte, expires time.Time) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO pending (digest, subscription_id, kind, expires) VALUES (?, ?, ?, ?)`,
		digest[:], id, kind, expires.Unix())
	return err
}

// Redeem uses up the code or token of kind whose SHA-256 is digest, if it has not expired, and returns the
// subscription it confirms.
func (s *Store) Redeem(ctx context.Context, kind string, digest [32]byte, now time.Time) (string, bool, error) {
	var id string
	err := s.db.QueryRowContext(ctx, `DELETE FROM pending WHERE digest = ? AND kind = ? AND expires > ? RETURNING subscription_id`,
		digest[:], kind, now.Unix()).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	return id, err == nil, err
}

// Confirm marks subscription id confirmed, with destination where it is not empty. It reports whether there is one.
func (s *Store) Confirm(ctx context.Context, id, destination string) (bool, error) {
	result, err := s.db.ExecContext(ctx, `UPDATE subscriptions SET confirmed = 1, destination = CASE WHEN ? = '' THEN destination ELSE ? END
		WHERE id = ?`, destination, destination, id)
	if err != nil {
		return false, err
	}
	changed, err := result.RowsAffected()
	return changed > 0, err
}

// Delete removes subscription id and what is pending for it.
func (s *Store) Delete(ctx context.Context, id string) error {
	return s.remove(ctx, `id = ?`, id)
}

// DeleteByDigest removes the subscription whose unsubscribe token has SHA-256 digest, and reports whether there was one.
func (s *Store) DeleteByDigest(ctx context.Context, digest [32]byte) (bool, error) {
	var id string
	err := s.db.QueryRowContext(ctx, `SELECT id FROM subscriptions WHERE unsubscribe_digest = ?`, digest[:]).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, s.Delete(ctx, id)
}

// DeleteChat removes every Telegram subscription linked to chat.
func (s *Store) DeleteChat(ctx context.Context, chat string) error {
	return s.remove(ctx, `channel = 'telegram' AND destination = ?`, chat)
}

func (s *Store) remove(ctx context.Context, where string, arg any) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `DELETE FROM pending WHERE subscription_id IN (SELECT id FROM subscriptions WHERE `+where+`)`, arg); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM subscriptions WHERE `+where, arg); err != nil {
		return err
	}
	return tx.Commit()
}

// Sweep deletes the codes and tokens that expired by now, and the unconfirmed subscriptions nothing is pending for.
func (s *Store) Sweep(ctx context.Context, now time.Time) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `DELETE FROM pending WHERE expires <= ?`, now.Unix()); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM subscriptions WHERE confirmed = 0 AND id NOT IN (SELECT subscription_id FROM pending)`); err != nil {
		return err
	}
	return tx.Commit()
}

// Get reads subscription id.
func (s *Store) Get(ctx context.Context, id string) (Subscription, bool, error) {
	subs, err := s.query(ctx, `SELECT id, account, channel, destination, threshold, confirmed, state FROM subscriptions WHERE id = ?`, id)
	if err != nil || len(subs) == 0 {
		return Subscription{}, false, err
	}
	return subs[0], true, nil
}

// Confirmed reads every confirmed subscription.
func (s *Store) Confirmed(ctx context.Context) ([]Subscription, error) {
	return s.query(ctx, `SELECT id, account, channel, destination, threshold, confirmed, state FROM subscriptions WHERE confirmed = 1 ORDER BY id`)
}

func (s *Store) query(ctx context.Context, query string, args ...any) ([]Subscription, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []Subscription
	for rows.Next() {
		var sub Subscription
		var account, state string
		if err := rows.Scan(&sub.ID, &account, &sub.Channel, &sub.Destination, &sub.Threshold, &sub.Confirmed, &state); err != nil {
			return nil, err
		}
		sub.Account = common.HexToAddress(account)
		if err := json.Unmarshal([]byte(state), &sub.State); err != nil {
			return nil, err
		}
		out = append(out, sub)
	}
	return out, rows.Err()
}

// SaveState keeps the alert state of subscription id.
func (s *Store) SaveState(ctx context.Context, id string, state map[string]PositionState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `UPDATE subscriptions SET state = ? WHERE id = ?`, string(data), id)
	return err
}
