// SPDX-License-Identifier: MIT OR Apache-2.0

package alerts

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tapehouse/tapehouse/services/internal/keeper"
	"github.com/tapehouse/tapehouse/services/internal/margin"
)

// Tries is how many times an alert is tried, once at each tick, before it is dropped.
const Tries = 3

// Service evaluates every subscribed account's positions and sends the alerts that are due.
type Service struct {
	Store    *Store
	Reader   *Reader
	Index    *keeper.Index
	Channels map[string]Channel
	Links    Links
	Server   *Server
	Log      *slog.Logger
	Now      func() time.Time
	tries    map[string]int
}

// NewService returns a service over its store and the chain, with the channels it may send by.
func NewService(store *Store, reader *Reader, index *keeper.Index, channels map[string]Channel, links Links, server *Server,
	log *slog.Logger) *Service {
	return &Service{Store: store, Reader: reader, Index: index, Channels: channels, Links: links, Server: server, Log: log,
		Now: time.Now, tries: map[string]int{}}
}

// Run ticks now and then every interval until ctx ends.
func (s *Service) Run(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		if err := s.Tick(ctx); err != nil && ctx.Err() == nil {
			s.Log.Error("tick", "error", err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

// Updater is a Telegram bot's incoming messages.
type Updater interface {
	Updates(ctx context.Context, offset int64, waitS int) ([]Update, error)
}

// Poll hands the messages sent to the bot to the opt-in until ctx ends.
func (s *Service) Poll(ctx context.Context, bot Updater) error {
	var offset int64
	for ctx.Err() == nil {
		updates, err := bot.Updates(ctx, offset, 25)
		if err != nil {
			s.Log.Error("telegram updates", "error", err)
			select {
			case <-ctx.Done():
			case <-time.After(5 * time.Second):
			}
			continue
		}
		for _, u := range updates {
			offset = max(offset, u.ID+1)
			if u.Text != "" {
				s.Server.Chat(ctx, u.Chat, u.Text)
			}
		}
	}
	return ctx.Err()
}

// Tick deletes what expired, reads the chain once and evaluates every confirmed subscription's positions at that block.
func (s *Service) Tick(ctx context.Context) error {
	if err := s.Store.Sweep(ctx, s.Now()); err != nil {
		return err
	}
	subs, err := s.Store.Confirmed(ctx)
	if err != nil || len(subs) == 0 {
		return err
	}
	m, err := s.Reader.Market(ctx, nil)
	if err != nil {
		return fmt.Errorf("the market: %w", err)
	}
	for _, sub := range subs {
		if err := s.subscription(ctx, m, sub); err != nil {
			s.Log.Error("subscription", "subscription", sub.ID, "error", err)
		}
	}
	return nil
}

// positions are the positions an account has held through the index: its margin positions and its shorts.
func (s *Service) positions(ctx context.Context, account common.Address) (ids [][32]byte, shorts []string, err error) {
	args := map[string]string{"account": account.Hex()}
	for _, event := range []string{"Deposit", "Borrow"} {
		events, err := s.Index.Events(ctx, "tapehouse.MarginAccounts", event, args)
		if err != nil {
			return nil, nil, err
		}
		for _, e := range events {
			raw, err := hexutil.Decode(e.Arg("position"))
			if err != nil || len(raw) != 32 {
				return nil, nil, fmt.Errorf("the index names a position %q", e.Arg("position"))
			}
			if id := [32]byte(raw); !slices.Contains(ids, id) {
				ids = append(ids, id)
			}
		}
	}
	events, err := s.Index.Events(ctx, "tapehouse.ShortPositions", "Sell", args)
	if err != nil {
		return nil, nil, err
	}
	for _, e := range events {
		asset, err := symbolName(e.Arg("symbol"))
		if err != nil {
			return nil, nil, err
		}
		if !slices.Contains(shorts, asset) {
			shorts = append(shorts, asset)
		}
	}
	return ids, shorts, nil
}

func (s *Service) subscription(ctx context.Context, m *Market, sub Subscription) error {
	channel, ok := s.Channels[sub.Channel]
	if !ok {
		return errors.New("its channel is closed")
	}
	ids, shorts, err := s.positions(ctx, sub.Account)
	if err != nil {
		return err
	}
	state := make(map[string]PositionState, len(sub.State))
	for k, v := range sub.State {
		state[k] = v
	}
	var failures []error
	send := func(key string, old PositionState, o Observation, texts func(Kind) Message) {
		due, next := Decide(sub.Threshold, old, o)
		for _, kind := range due {
			message := texts(kind)
			message.Body += s.footer(sub)
			if err := channel.Send(ctx, sub.Destination, message); err != nil {
				tries := fmt.Sprintf("%s/%s/%d", sub.ID, key, kind)
				s.tries[tries]++
				if s.tries[tries] < Tries {
					revert(&next, old, kind)
					s.Log.Warn("alert not sent", "subscription", sub.ID, "kind", kind, "attempt", s.tries[tries], "error", err)
				} else {
					s.Log.Error("alert dropped", "subscription", sub.ID, "kind", kind, "error", err)
					delete(s.tries, tries)
				}
				continue
			}
			delete(s.tries, fmt.Sprintf("%s/%s/%d", sub.ID, key, kind))
		}
		state[key] = next
	}
	for _, id := range ids {
		if err := s.position(ctx, m, sub, id, state, send); err != nil {
			failures = append(failures, err)
		}
	}
	for _, asset := range shorts {
		if err := s.short(ctx, m, sub, asset, state, send); err != nil {
			failures = append(failures, err)
		}
	}
	if err := s.Store.SaveState(ctx, sub.ID, state); err != nil {
		failures = append(failures, err)
	}
	return errors.Join(failures...)
}

type sender func(key string, old PositionState, o Observation, texts func(Kind) Message)

func (s *Service) position(ctx context.Context, m *Market, sub Subscription, id [32]byte, state map[string]PositionState, send sender) error {
	key := hexutil.Encode(id[:])
	j, err := s.Reader.Judge(ctx, m, sub.Account, id)
	if err != nil {
		return fmt.Errorf("position %s: %w", key, err)
	}
	if !j.Evaluated {
		s.Log.Warn("not evaluated", "subscription", sub.ID, "position", key)
		return nil
	}
	v := j.Valuation
	o := Observation{Judged: v.Judged, Held: v.Held, Closing: m.Regime.Kind == margin.Closing, Closure: m.BoundaryMs}
	if v.Judged {
		o.Equity, o.Requirement, o.Short = v.Equity, v.Requirement, v.Short()
	}
	if o.Short {
		if o.Auction, err = s.Reader.Auction(ctx, m, sub.Account, id); err != nil {
			return fmt.Errorf("position %s: %w", key, err)
		}
	}
	send(key, state[key], o, func(kind Kind) Message {
		switch kind {
		case Threshold:
			return thresholdText(key, sub.Threshold, v.Equity, v.Requirement)
		case Liquidation:
			return liquidationText(key, v.Equity, v.Requirement, o.Auction)
		case Weekend:
			return s.weekend(ctx, m, key, j)
		}
		return unjudgedText(key)
	})
	return nil
}

func (s *Service) weekend(ctx context.Context, m *Market, key string, j *Judgement) Message {
	p := j.Position
	closed := m.Value(p, m.Current(), closedRegime(), false)
	var assets []assetPrices
	for _, stock := range p.Stocks {
		a := m.Assets[m.index(stock.Asset)]
		prices := assetPrices{Asset: stock.Asset, Centre: a.Mid, Weekend: m.WeekendPrice(p, stock.Asset),
			Reopening: m.ReopeningPrice(p, stock.Asset)}
		if open, err := s.Reader.LiquidationPrice(ctx, m, p, stock.Asset); err == nil {
			prices.Open = open
		}
		assets = append(assets, prices)
	}
	return weekendText(key, m.BoundaryMs, j.Valuation.Equity, closed.Requirement, assets)
}

func (s *Service) short(ctx context.Context, m *Market, sub Subscription, asset string, state map[string]PositionState, send sender) error {
	key := "short:" + asset
	equity, requirement, open, err := s.Reader.Short(ctx, m, sub.Account, asset)
	if err != nil {
		return fmt.Errorf("short %s: %w", asset, err)
	}
	if !open {
		delete(state, key)
		return nil
	}
	o := Observation{Judged: true, Equity: equity, Requirement: requirement, Short: equity.Cmp(requirement) < 0}
	send(key, state[key], o, func(kind Kind) Message {
		if kind == Liquidation {
			return liquidationText(key, equity, requirement, 0)
		}
		return thresholdText(key, sub.Threshold, equity, requirement)
	})
	return nil
}

func (s *Service) footer(sub Subscription) string {
	if sub.Channel == Email {
		return "\n--\nUnsubscribe: " + s.Links.Unsubscribe(sub.ID) + "\n"
	}
	return "\nSend /stop to stop these alerts.\n"
}
