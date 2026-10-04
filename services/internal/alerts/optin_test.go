// SPDX-License-Identifier: MIT OR Apache-2.0

package alerts

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type delivery struct {
	destination string
	message     Message
}

// outbox is a channel that keeps what it was asked to send.
type outbox struct {
	mu   sync.Mutex
	sent []delivery
	bot  string
	fail error
}

func (o *outbox) Send(_ context.Context, destination string, m Message) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.fail != nil {
		return o.fail
	}
	o.sent = append(o.sent, delivery{destination, m})
	return nil
}

func (o *outbox) Username(context.Context) (string, error) { return o.bot, nil }

func (o *outbox) take() []delivery {
	o.mu.Lock()
	defer o.mu.Unlock()
	out := o.sent
	o.sent = nil
	return out
}

type optin struct {
	t       *testing.T
	store   *Store
	server  *Server
	http    *httptest.Server
	email   *outbox
	bot     *outbox
	now     time.Time
	key     *ecdsa.PrivateKey
	account common.Address
}

func newOptin(t *testing.T) *optin {
	t.Helper()
	store, err := OpenStore(filepath.Join(t.TempDir(), "alerts.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	o := &optin{t: t, store: store, email: &outbox{}, bot: &outbox{bot: "tapehouse_bot"}, now: time.Unix(1_790_775_000, 0)}
	o.key, _ = crypto.GenerateKey()
	o.account = crypto.PubkeyToAddress(o.key.PublicKey)
	o.server = NewServer(store, 412346, Links{BaseURL: "https://alerts.example", Secret: bytes.Repeat([]byte{7}, 32)}, o.email, o.bot,
		slog.New(slog.NewTextHandler(io.Discard, nil)))
	o.server.Now = func() time.Time { return o.now }
	o.http = httptest.NewServer(o.server.Handler())
	t.Cleanup(o.http.Close)
	return o
}

func sign(t *testing.T, key *ecdsa.PrivateKey, text string) hexutil.Bytes {
	t.Helper()
	signature, err := crypto.Sign(accounts.TextHash([]byte(text)), key)
	if err != nil {
		t.Fatal(err)
	}
	signature[64] += 27
	return signature
}

func (o *optin) post(body any) (int, map[string]any) {
	o.t.Helper()
	data, _ := json.Marshal(body)
	response, err := http.Post(o.http.URL+"/v1/subscriptions", "application/json", bytes.NewReader(data))
	if err != nil {
		o.t.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()
	var out map[string]any
	_ = json.NewDecoder(response.Body).Decode(&out)
	return response.StatusCode, out
}

func (o *optin) subscribe(channel, destination string, extra map[string]any) (int, map[string]any) {
	o.t.Helper()
	issued := o.now.Unix()
	body := map[string]any{"account": o.account, "channel": channel, "destination": destination, "issuedAt": issued,
		"signature": sign(o.t, o.key, SubscribeText(412346, o.account, channel, issued))}
	for k, v := range extra {
		body[k] = v
	}
	return o.post(body)
}

func (o *optin) get(path string) int {
	o.t.Helper()
	response, err := http.Get(o.http.URL + path)
	if err != nil {
		o.t.Fatal(err)
	}
	_ = response.Body.Close()
	return response.StatusCode
}

func (o *optin) confirmed() []Subscription {
	o.t.Helper()
	subs, err := o.store.Confirmed(context.Background())
	if err != nil {
		o.t.Fatal(err)
	}
	return subs
}

var link = regexp.MustCompile(`https://alerts\.example/v1/confirm/([0-9a-f]{64})`)

func TestOnlyTheAccountSubscribesItAndNothingIsSentUntilConfirmed(t *testing.T) {
	o := newOptin(t)
	other, _ := crypto.GenerateKey()
	issued := o.now.Unix()
	for name, body := range map[string]map[string]any{
		"a signature of another address": {"account": o.account, "channel": Email, "destination": "ada@example.org", "issuedAt": issued,
			"signature": sign(t, other, SubscribeText(412346, o.account, Email, issued))},
		"a signature of another channel": {"account": o.account, "channel": Email, "destination": "ada@example.org", "issuedAt": issued,
			"signature": sign(t, o.key, SubscribeText(412346, o.account, Telegram, issued))},
		"a signature of another chain": {"account": o.account, "channel": Email, "destination": "ada@example.org", "issuedAt": issued,
			"signature": sign(t, o.key, SubscribeText(1, o.account, Email, issued))},
		"a stale issuedAt": {"account": o.account, "channel": Email, "destination": "ada@example.org", "issuedAt": issued - 601,
			"signature": sign(t, o.key, SubscribeText(412346, o.account, Email, issued-601))},
		"an issuedAt in the future": {"account": o.account, "channel": Email, "destination": "ada@example.org", "issuedAt": issued + 120,
			"signature": sign(t, o.key, SubscribeText(412346, o.account, Email, issued+120))},
		"no signature": {"account": o.account, "channel": Email, "destination": "ada@example.org", "issuedAt": issued},
	} {
		if status, _ := o.post(body); status != http.StatusUnauthorized {
			t.Fatalf("%s: status %d", name, status)
		}
	}
	if sent := o.email.take(); len(sent) != 0 {
		t.Fatalf("an email went to someone who did not ask: %v", sent)
	}
	if subs := o.confirmed(); len(subs) != 0 {
		t.Fatalf("subscriptions %v", subs)
	}
	status, out := o.subscribe(Email, "ada@example.org", nil)
	if status != http.StatusAccepted || out["status"] != "pending" {
		t.Fatalf("a signed subscription: %d %v", status, out)
	}
	sent := o.email.take()
	if len(sent) != 1 || sent[0].destination != "ada@example.org" || !strings.Contains(sent[0].message.Body, "/v1/unsubscribe/") {
		t.Fatalf("the confirmation: %v", sent)
	}
	if subs := o.confirmed(); len(subs) != 0 {
		t.Fatalf("an unconfirmed email is subscribed: %v", subs)
	}
	token := link.FindStringSubmatch(sent[0].message.Body)
	if token == nil {
		t.Fatalf("no confirmation link in %q", sent[0].message.Body)
	}
	if status := o.get("/v1/confirm/" + token[1]); status != http.StatusOK {
		t.Fatalf("confirm: %d", status)
	}
	subs := o.confirmed()
	if len(subs) != 1 || subs[0].Account != o.account || subs[0].Destination != "ada@example.org" || subs[0].Threshold != DefaultThreshold {
		t.Fatalf("subscriptions %+v", subs)
	}
	if status := o.get("/v1/confirm/" + token[1]); status != http.StatusNotFound {
		t.Fatalf("a link used twice: %d", status)
	}
}

func TestAnEmailLinkExpiresAfterADayAndTakesItsRowWithIt(t *testing.T) {
	o := newOptin(t)
	o.subscribe(Email, "ada@example.org", nil)
	token := link.FindStringSubmatch(o.email.take()[0].message.Body)[1]
	o.now = o.now.Add(EmailTokenTTL + time.Second)
	if status := o.get("/v1/confirm/" + token); status != http.StatusNotFound {
		t.Fatalf("an expired link: %d", status)
	}
	if err := o.store.Sweep(context.Background(), o.now); err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"subscriptions", "pending"} {
		if rows := count(t, o.store, table); rows != 0 {
			t.Fatalf("%d rows of %s after the link expired", rows, table)
		}
	}
}

func count(t *testing.T, s *Store, table string) int {
	t.Helper()
	var rows int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM ` + table).Scan(&rows); err != nil {
		t.Fatal(err)
	}
	return rows
}

func TestATelegramChatIsLinkedByAOneTimeCode(t *testing.T) {
	o := newOptin(t)
	status, out := o.subscribe(Telegram, "", map[string]any{"threshold": 1.5})
	if status != http.StatusAccepted {
		t.Fatalf("status %d", status)
	}
	code, _ := out["code"].(string)
	if !regexp.MustCompile(`^[A-Za-z0-9_-]{24}$`).MatchString(code) || out["link"] != "https://t.me/tapehouse_bot?start="+code {
		t.Fatalf("code %q link %v", code, out["link"])
	}
	if subs := o.confirmed(); len(subs) != 0 {
		t.Fatalf("an unlinked chat is subscribed: %v", subs)
	}
	o.server.Chat(context.Background(), "42", "/start wrong-code")
	if subs := o.confirmed(); len(subs) != 0 {
		t.Fatalf("a wrong code linked a chat: %v", subs)
	}
	o.server.Chat(context.Background(), "42", "/start "+code)
	subs := o.confirmed()
	if len(subs) != 1 || subs[0].Destination != "42" || subs[0].Channel != Telegram || subs[0].Threshold != 1500 {
		t.Fatalf("subscriptions %+v", subs)
	}
	if sent := o.bot.take(); len(sent) != 2 || sent[1].destination != "42" || !strings.Contains(sent[1].message.Body, "Linked") {
		t.Fatalf("the bot's replies: %v", sent)
	}
	o.server.Chat(context.Background(), "43", "/start "+code)
	if subs := o.confirmed(); len(subs) != 1 || subs[0].Destination != "42" {
		t.Fatalf("a code used twice: %+v", subs)
	}
	o.server.Chat(context.Background(), "42", "/stop")
	if subs := o.confirmed(); len(subs) != 0 || count(t, o.store, "subscriptions") != 0 {
		t.Fatal("/stop left a row")
	}
	o.subscribe(Telegram, "", nil)
	o.now = o.now.Add(TelegramCodeTTL + time.Second)
	if err := o.store.Sweep(context.Background(), o.now); err != nil || count(t, o.store, "subscriptions") != 0 || count(t, o.store, "pending") != 0 {
		t.Fatalf("an expired code left rows: %v", err)
	}
}

func TestUnsubscribingLeavesNoRow(t *testing.T) {
	o := newOptin(t)
	confirm := func() string {
		o.subscribe(Email, "ada@example.org", nil)
		token := link.FindStringSubmatch(o.email.take()[0].message.Body)[1]
		o.get("/v1/confirm/" + token)
		return o.confirmed()[0].ID
	}
	id := confirm()
	unsubscribe := o.server.Links.Unsubscribe(id)
	if status := o.get(strings.TrimPrefix(unsubscribe, "https://alerts.example")); status != http.StatusOK {
		t.Fatalf("the link at the foot of an email: %d", status)
	}
	if count(t, o.store, "subscriptions") != 0 || count(t, o.store, "pending") != 0 {
		t.Fatal("a row is left after the link")
	}
	if status := o.get("/v1/unsubscribe/" + strings.Repeat("0", 64)); status != http.StatusNotFound {
		t.Fatalf("an unknown token: %d", status)
	}
	id = confirm()
	stranger, _ := crypto.GenerateKey()
	issued := o.now.Unix()
	remove := func(key *ecdsa.PrivateKey, at int64) int {
		body, _ := json.Marshal(map[string]any{"issuedAt": at, "signature": sign(t, key, UnsubscribeText(412346, o.account, id, at))})
		request, _ := http.NewRequest(http.MethodDelete, o.http.URL+"/v1/subscriptions/"+id, bytes.NewReader(body))
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		_ = response.Body.Close()
		return response.StatusCode
	}
	if status := remove(stranger, issued); status != http.StatusUnauthorized {
		t.Fatalf("another address deleted a subscription: %d", status)
	}
	if status := remove(o.key, issued-601); status != http.StatusUnauthorized {
		t.Fatalf("a stale signature deleted a subscription: %d", status)
	}
	if status := remove(o.key, issued); status != http.StatusNoContent {
		t.Fatalf("the account's fresh signature: %d", status)
	}
	if count(t, o.store, "subscriptions") != 0 {
		t.Fatal("a row is left after the account deleted its subscription")
	}
}

func TestOnlyWhatAnAlertNeedsIsKept(t *testing.T) {
	o := newOptin(t)
	columns := func(table string) []string {
		rows, err := o.store.db.Query(`SELECT name FROM pragma_table_info(?) ORDER BY cid`, table)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = rows.Close() }()
		var out []string
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				t.Fatal(err)
			}
			out = append(out, name)
		}
		return out
	}
	if got, want := columns("subscriptions"), strings.Fields("id account channel destination threshold confirmed state unsubscribe_digest"); !slices.Equal(got, want) {
		t.Fatalf("subscriptions: %v", got)
	}
	if got, want := columns("pending"), strings.Fields("digest subscription_id kind expires"); !slices.Equal(got, want) {
		t.Fatalf("pending: %v", got)
	}
	var tables int
	if err := o.store.db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table'`).Scan(&tables); err != nil || tables != 2 {
		t.Fatalf("%d tables, %v", tables, err)
	}
	o.subscribe(Email, "ada@example.org", nil)
	var digest []byte
	if err := o.store.db.QueryRow(`SELECT digest FROM pending`).Scan(&digest); err != nil || len(digest) != 32 {
		t.Fatalf("a pending token is kept as its SHA-256: %d bytes, %v", len(digest), err)
	}
	var destination string
	var unsubscribe []byte
	if err := o.store.db.QueryRow(`SELECT destination, unsubscribe_digest FROM subscriptions`).Scan(&destination, &unsubscribe); err != nil || len(unsubscribe) != 32 {
		t.Fatalf("%v", err)
	}
	if token := o.server.Links.UnsubscribeToken(o.confirmedOrPendingID()); bytes.Contains(unsubscribe, []byte(token)) {
		t.Fatal("the unsubscribe token itself is kept")
	}
}

func (o *optin) confirmedOrPendingID() string {
	var id string
	if err := o.store.db.QueryRow(`SELECT id FROM subscriptions`).Scan(&id); err != nil && err != sql.ErrNoRows {
		o.t.Fatal(err)
	}
	return id
}

func TestARequestIsRefusedIfItsFieldsAreWrong(t *testing.T) {
	o := newOptin(t)
	for name, test := range map[string]struct {
		channel, destination string
		extra                map[string]any
	}{
		"a threshold below 1.05":        {Email, "ada@example.org", map[string]any{"threshold": 1.0}},
		"a threshold above 3":           {Email, "ada@example.org", map[string]any{"threshold": 3.5}},
		"no address":                    {Email, "", nil},
		"an address with a name":        {Email, "Ada <ada@example.org>", nil},
		"a header injection":            {Email, "ada@example.org\r\nBcc: eve@example.org", nil},
		"a channel that does not exist": {"sms", "123", nil},
	} {
		if status, _ := o.subscribe(test.channel, test.destination, test.extra); status != http.StatusBadRequest && status != http.StatusUnauthorized {
			t.Errorf("%s: status %d", name, status)
		}
	}
	if count(t, o.store, "subscriptions") != 0 || len(o.email.take()) != 0 {
		t.Fatal("a refused request left a row or an email")
	}
	if status, _ := o.subscribe(Email, "ada@example.org", map[string]any{"threshold": 3.0}); status != http.StatusAccepted {
		t.Fatalf("a threshold of 3: %d", status)
	}
	if status, _ := o.subscribe(Email, "ada@example.org", map[string]any{"threshold": 1.05}); status != http.StatusAccepted {
		t.Fatalf("a threshold of 1.05: %d", status)
	}
	if rows := count(t, o.store, "subscriptions"); rows != 1 {
		t.Fatalf("a second request for the same account and channel must replace the first: %d rows", rows)
	}
}
