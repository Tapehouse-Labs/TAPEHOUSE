// SPDX-License-Identifier: MIT OR Apache-2.0

package alerts

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"math/big"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tapehouse/tapehouse/services/internal/band"
	"github.com/tapehouse/tapehouse/services/internal/keeper"
)

var holderAccount = common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")

// service is the whole service over doubles of the chain, the index and both channels.
type service struct {
	*world
	t      *testing.T
	store  *Store
	svc    *Service
	email  *outbox
	chat   *outbox
	logs   *bytes.Buffer
	events []map[string]any
	index  *httptest.Server
	mu     sync.Mutex
}

func newService(t *testing.T) *service {
	t.Helper()
	s := &service{world: newWorld(t), t: t, email: &outbox{}, chat: &outbox{bot: "tapehouse_bot"}, logs: &bytes.Buffer{}}
	s.index = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		query := r.URL.Query()
		var out []map[string]any
		for _, e := range s.events {
			if e["contract"] == query.Get("contract") && e["event"] == query.Get("event") && e["args"].(map[string]any)["account"] == query.Get("arg[account]") {
				out = append(out, map[string]any{"block": 1, "logIndex": len(out), "args": e["args"]})
			}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"events": out, "next": nil})
	}))
	t.Cleanup(s.index.Close)
	var err error
	if s.store, err = OpenStore(filepath.Join(t.TempDir(), "alerts.db")); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.store.Close() })
	log := slog.New(slog.NewTextHandler(s.logs, nil))
	links := Links{BaseURL: "https://alerts.example", Secret: bytes.Repeat([]byte{7}, 32)}
	server := NewServer(s.store, 412346, links, s.email, s.chat, log)
	s.svc = NewService(s.store, NewReader(s.client), keeper.NewIndex(s.index.URL, http.DefaultClient),
		map[string]Channel{Email: s.email, Telegram: s.chat}, links, server, log)
	s.time = 1_790_775_000
	return s
}

func (s *service) hold(id [32]byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, map[string]any{"contract": "tapehouse.MarginAccounts", "event": "Deposit",
		"args": map[string]any{"account": holderAccount.Hex(), "position": hexutil.Encode(id[:])}})
}

func (s *service) subscribe(channel, destination string, threshold int) string {
	s.t.Helper()
	ctx := context.Background()
	id, err := s.store.Subscribe(ctx, holderAccount, channel, destination, threshold, s.svc.Links.digestOf)
	if err != nil {
		s.t.Fatal(err)
	}
	if _, err := s.store.Confirm(ctx, id, ""); err != nil {
		s.t.Fatal(err)
	}
	return id
}

func (s *service) tick() {
	s.t.Helper()
	if err := s.svc.Tick(context.Background()); err != nil {
		s.t.Fatal(err)
	}
}

// stocks puts ten NVDA in the position with debt a fraction of the debt at which it falls short in an open market.
func (s *service) stocks(percent int64) {
	debt := boundary(s.m, false) * percent / 100
	s.holdings = map[string]*big.Int{"collateral NVDA": units(10, 18), "debt": units(debt, 6)}
}

func TestAThresholdAlertGoesOutOnceByEmailAndByTelegram(t *testing.T) {
	s := newService(t)
	s.hold([32]byte{1})
	s.stocks(98)
	s.subscribe(Email, "ada@example.org", DefaultThreshold)
	s.subscribe(Telegram, "4242", DefaultThreshold)
	s.tick()
	mail, chat := s.email.take(), s.chat.take()
	if len(mail) != 1 || mail[0].destination != "ada@example.org" || !strings.Contains(mail[0].message.Subject, "close to its requirement") ||
		!strings.Contains(mail[0].message.Body, "/v1/unsubscribe/") {
		t.Fatalf("email %v", mail)
	}
	if len(chat) != 1 || chat[0].destination != "4242" || !strings.Contains(chat[0].message.Body, "/stop") || !strings.Contains(chat[0].message.Body, "under your threshold of 1.250") {
		t.Fatalf("telegram %v", chat)
	}
	for range 3 {
		s.tick()
	}
	if len(s.email.take())+len(s.chat.take()) != 0 {
		t.Fatal("an alert again at a ratio that stayed under the threshold")
	}
	for _, secret := range []string{"ada@example.org", "4242"} {
		if strings.Contains(s.logs.String(), secret) {
			t.Fatalf("a log line names a destination:\n%s", s.logs)
		}
	}
	s.stocks(30)
	s.tick()
	s.stocks(98)
	s.tick()
	if len(s.email.take()) != 1 {
		t.Fatal("no alert after the ratio rose above the threshold and fell again")
	}
}

func TestNothingIsReadForAnUnconfirmedSubscription(t *testing.T) {
	s := newService(t)
	s.hold([32]byte{1})
	s.stocks(98)
	if _, err := s.store.Subscribe(context.Background(), holderAccount, Email, "ada@example.org", DefaultThreshold, s.svc.Links.digestOf); err != nil {
		t.Fatal(err)
	}
	s.tick()
	if len(s.email.take()) != 0 || len(s.calls) != 0 {
		t.Fatalf("an unconfirmed subscription cost %d calls or got an alert", len(s.calls))
	}
}

func TestOneWeekendAlertWithTheSearchedPricesIsSentForEachClosure(t *testing.T) {
	s := newService(t)
	s.hold([32]byte{1})
	s.stocks(0)
	s.holdings["debt"] = units(boundary(s.m, true)*90/100, 6)
	s.subscribe(Email, "ada@example.org", DefaultThreshold)
	s.session = sessionState{state: 2, nyse: 1, nyseNext: 3, boundaryMs: s.time*1000 + 3*3_600_000}
	s.tick()
	sent := s.email.take()
	if len(sent) != 1 {
		t.Fatalf("%d alerts, want the weekend one: %v", len(sent), sent)
	}
	m, err := s.svc.Reader.Market(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	p := s.position(holderAccount)
	weekend, reopening := m.WeekendPrice(p, "NVDA"), m.ReopeningPrice(p, "NVDA")
	if weekend.Status != Priced || weekend.Centre >= m.Assets[0].Mid {
		t.Fatalf("the fixture must have a weekend price under the centre: %+v", weekend)
	}
	body := sent[0].message.Body
	for _, want := range []string{"Weekend liquidation price: " + price(weekend.Centre), "Reopening price: " + price(reopening.Centre),
		"Open-market liquidation price: 77.00 USD", "reopening auction", "requirement across the closure"} {
		if !strings.Contains(body, want) {
			t.Fatalf("the alert lacks %q:\n%s", want, body)
		}
	}
	s.tick()
	s.tick()
	if len(s.email.take()) != 0 {
		t.Fatal("a second weekend alert for one closure")
	}
	s.session.boundaryMs += 7 * 24 * 3_600_000
	s.time += 7 * 24 * 3600
	s.session.boundaryMs = s.time*1000 + 2*3_600_000
	s.tick()
	if got := s.email.take(); len(got) != 1 || !strings.Contains(got[0].message.Subject, "about to close") {
		t.Fatalf("the next closure: %v", got)
	}
}

func TestAServiceSendsOneLiquidationAlertForEachAuction(t *testing.T) {
	s := newService(t)
	s.hold([32]byte{1})
	s.stocks(0)
	s.holdings["debt"] = units(boundary(s.m, false)*110/100, 6)
	s.subscribe(Email, "ada@example.org", DefaultThreshold)
	s.tick()
	if got := s.email.take(); len(got) != 1 || !strings.Contains(got[0].message.Subject, "falls short") {
		t.Fatalf("short now: %v", got)
	}
	s.auction = s.time - 30
	s.tick()
	s.tick()
	if got := s.email.take(); len(got) != 0 {
		t.Fatalf("the auction of a position already warned of: %v", got)
	}
	s.auction = s.time - 10
	s.tick()
	if got := s.email.take(); len(got) != 1 || !strings.Contains(got[0].message.Body, "An auction of it started") {
		t.Fatalf("a new auction: %v", got)
	}
}

func TestAHaltedAssetGivesOneUnjudgedNoticeAndNoPrice(t *testing.T) {
	s := newService(t)
	s.hold([32]byte{1})
	s.stocks(50)
	s.subscribe(Email, "ada@example.org", DefaultThreshold)
	s.m.Assets[0].State = uint8(band.Halted)
	s.m.Assets[0].Low, s.m.Assets[0].High, s.m.Assets[0].Mid = 0, n(0), 0
	s.session = sessionState{state: 2, nyse: 1, nyseNext: 3, boundaryMs: s.time*1000 + 3*3_600_000}
	s.tick()
	got := s.email.take()
	if len(got) != 1 || !strings.Contains(got[0].message.Subject, "cannot be judged") || strings.Contains(got[0].message.Body, "liquidation price") {
		t.Fatalf("%v", got)
	}
	s.tick()
	if len(s.email.take()) != 0 {
		t.Fatal("a second notice for one halt")
	}
}

func TestAFailedSendIsTriedAgainAtTheNextTickAndDroppedAfterThree(t *testing.T) {
	s := newService(t)
	s.hold([32]byte{1})
	s.stocks(98)
	s.subscribe(Email, "ada@example.org", DefaultThreshold)
	s.email.fail = errors.New("email send failed")
	s.tick()
	s.email.fail = nil
	s.tick()
	if got := s.email.take(); len(got) != 1 {
		t.Fatalf("an alert that failed once must go out at the next tick: %v", got)
	}
	other := newService(t)
	other.hold([32]byte{1})
	other.stocks(98)
	id := other.subscribe(Email, "ada@example.org", DefaultThreshold)
	other.email.fail = errors.New("email send failed")
	for range Tries {
		other.tick()
	}
	other.email.fail = nil
	other.tick()
	if got := other.email.take(); len(got) != 0 {
		t.Fatalf("an alert dropped after %d tries was sent: %v", Tries, got)
	}
	if !strings.Contains(other.logs.String(), "alert dropped") || !strings.Contains(other.logs.String(), id) || strings.Contains(other.logs.String(), "ada@example.org") {
		t.Fatalf("the log:\n%s", other.logs)
	}
}

func TestAPositionWhoseRequirementTheEngineDoesNotConfirmIsNotAlertedOn(t *testing.T) {
	s := newService(t)
	s.hold([32]byte{1})
	s.stocks(98)
	s.subscribe(Email, "ada@example.org", DefaultThreshold)
	good := s.requirement
	s.requirement = func(q, p []*big.Int) (*big.Int, uint8) {
		r, g := good(q, p)
		return r.Add(r, n(1)), g
	}
	s.tick()
	if len(s.email.take()) != 0 || !strings.Contains(s.logs.String(), "not evaluated") {
		t.Fatalf("log:\n%s", s.logs)
	}
}

func TestAShortIsAlertedOnByItsHealth(t *testing.T) {
	s := newService(t)
	s.mu.Lock()
	s.events = append(s.events, map[string]any{"contract": "tapehouse.ShortPositions", "event": "Sell",
		"args": map[string]any{"account": holderAccount.Hex(), "symbol": hexutil.Encode(common.RightPadBytes([]byte("NVDA"), 32))}})
	s.mu.Unlock()
	s.subscribe(Telegram, "4242", DefaultThreshold)
	s.short = &shortState{equity: units(900, 18), requirement: units(1000, 18)}
	s.tick()
	got := s.chat.take()
	if len(got) != 1 || !strings.Contains(got[0].message.Subject, "your short of NVDA falls short") {
		t.Fatalf("%v", got)
	}
	s.tick()
	if len(s.chat.take()) != 0 {
		t.Fatal("a second alert for the same shortfall")
	}
	s.short = nil
	s.tick()
}

type updates struct {
	batches [][]Update
	done    context.CancelFunc
}

func (u *updates) Updates(ctx context.Context, _ int64, _ int) ([]Update, error) {
	if len(u.batches) == 0 {
		u.done()
		<-ctx.Done()
		return nil, ctx.Err()
	}
	batch := u.batches[0]
	u.batches = u.batches[1:]
	return batch, nil
}

func TestTheBotsMessagesReachTheOptIn(t *testing.T) {
	s := newService(t)
	id, err := s.store.Subscribe(context.Background(), holderAccount, Telegram, "", DefaultThreshold, s.svc.Links.digestOf)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.store.AddPending(context.Background(), id, Telegram, Digest("a-code"), time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	feed := &updates{done: cancel, batches: [][]Update{{{ID: 5, Chat: "42", Text: "/start a-code"}, {ID: 6}}}}
	if err := s.svc.Poll(ctx, feed); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	subs, _ := s.store.Confirmed(context.Background())
	if len(subs) != 1 || subs[0].Destination != "42" {
		t.Fatalf("%+v", subs)
	}
}
