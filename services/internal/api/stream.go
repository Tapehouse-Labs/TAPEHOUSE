// SPDX-License-Identifier: MIT OR Apache-2.0

package api

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/tapehouse/tapehouse/services/internal/indexer"
	"github.com/tapehouse/tapehouse/services/internal/store"
)

// StreamBuffer is how many messages a subscriber may fall behind before the stream drops it.
const StreamBuffer = 256

// Message is what the stream sends: "head" after every step of the indexer, "event" for each event the subscription
// selects, and "rewind" when a reorganisation removed every event above Block, which the client must drop.
type Message struct {
	Type      string      `json:"type"`
	Head      *store.Head `json:"head,omitempty"`
	Finalized uint64      `json:"finalized"`
	Event     *Event      `json:"event,omitempty"`
	Block     *uint64     `json:"block,omitempty"`
}

// Hub fans the indexer's updates out to the stream's subscribers.
type Hub struct {
	mu          sync.Mutex
	subscribers map[*subscriber]struct{}
	served      func([]store.Event) []Event
}

type subscriber struct {
	contract, event string
	messages        chan []byte
	dropped         chan struct{}
}

// NewHub returns a hub with no subscribers.
func NewHub() *Hub {
	return &Hub{subscribers: map[*subscriber]struct{}{}}
}

// Publish sends u to every subscriber, and drops a subscriber that has fallen StreamBuffer messages behind.
func (h *Hub) Publish(u indexer.Update) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.subscribers) == 0 || h.served == nil {
		return
	}
	var messages []Message
	if u.Rewound != nil {
		messages = append(messages, Message{Type: "rewind", Block: u.Rewound, Finalized: u.Finalized})
	}
	for _, event := range h.served(u.Events) {
		messages = append(messages, Message{Type: "event", Event: &event, Finalized: u.Finalized})
	}
	head := u.Head
	messages = append(messages, Message{Type: "head", Head: &head, Finalized: u.Finalized})
	encoded := make([][]byte, len(messages))
	for i, message := range messages {
		encoded[i], _ = json.Marshal(message)
	}
	for sub := range h.subscribers {
		for i, message := range messages {
			if message.Event != nil && !sub.selects(message.Event) {
				continue
			}
			select {
			case sub.messages <- encoded[i]:
			default:
				delete(h.subscribers, sub)
				close(sub.dropped)
			}
			if _, ok := h.subscribers[sub]; !ok {
				break
			}
		}
	}
}

func (s *subscriber) selects(event *Event) bool {
	return (s.contract == "" || s.contract == event.Contract) && (s.event == "" || s.event == event.Event.Event)
}

// Subscribers is how many streams are open.
func (h *Hub) Subscribers() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.subscribers)
}

func (h *Hub) subscribe(contract, event string) *subscriber {
	sub := &subscriber{contract, event, make(chan []byte, StreamBuffer), make(chan struct{})}
	h.mu.Lock()
	h.subscribers[sub] = struct{}{}
	h.mu.Unlock()
	return sub
}

func (h *Hub) unsubscribe(sub *subscriber) {
	h.mu.Lock()
	delete(h.subscribers, sub)
	h.mu.Unlock()
}

// stream upgrades to a WebSocket that sends the indexer's updates as they happen, optionally only the events of one
// contract, or of one event of it. A client too slow to keep up is disconnected; it resumes from /v1/events.
func (s *Server) stream(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	contract, event := query.Get("contract"), query.Get("event")
	for key := range query {
		if !slices.Contains([]string{"contract", "event"}, key) {
			s.reply(w, nil, invalid("unknown parameter %s", key))
			return
		}
	}
	if contract != "" {
		if _, ok := s.Catalog.Contract(contract); !ok {
			s.reply(w, nil, missing("the registry names no contract %s", contract))
			return
		}
	}
	release, ok := s.Tiers.OpenStream(w, r)
	if !ok {
		return
	}
	defer release()
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: []string{"*"}})
	if err != nil {
		return
	}
	defer func() { _ = conn.CloseNow() }()
	ctx := conn.CloseRead(r.Context())
	sub := s.Hub.subscribe(contract, event)
	defer s.Hub.unsubscribe(sub)
	for {
		select {
		case <-ctx.Done():
			return
		case <-sub.dropped:
			_ = conn.Close(websocket.StatusPolicyViolation, "too slow: resume from /v1/events")
			return
		case message := <-sub.messages:
			write, cancel := context.WithTimeout(ctx, 10*time.Second)
			err := conn.Write(write, websocket.MessageText, message)
			cancel()
			if err != nil {
				return
			}
		}
	}
}
