// SPDX-License-Identifier: MIT OR Apache-2.0

package api

import (
	"testing"

	"github.com/tapehouse/tapehouse/services/internal/indexer"
	"github.com/tapehouse/tapehouse/services/internal/store"
)

func TestASubscriberThatFallsBehindIsDropped(t *testing.T) {
	hub := NewHub()
	hub.served = func([]store.Event) []Event { return nil }
	slow, keeping := hub.subscribe("", ""), hub.subscribe("", "")
	for i := range StreamBuffer + 1 {
		hub.Publish(indexer.Update{Head: store.Head{Number: uint64(i)}})
		if i < StreamBuffer {
			<-keeping.messages
		}
	}
	select {
	case <-slow.dropped:
	default:
		t.Fatal("a subscriber StreamBuffer messages behind was kept")
	}
	select {
	case <-keeping.dropped:
		t.Fatal("a subscriber that kept up was dropped")
	default:
	}
	if hub.Subscribers() != 1 {
		t.Fatalf("%d subscribers", hub.Subscribers())
	}
}
