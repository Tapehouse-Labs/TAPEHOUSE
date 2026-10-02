// SPDX-License-Identifier: MIT OR Apache-2.0

package api

import (
	"container/list"
	"sync"
)

// lru keeps the latest answers to reads at a block, which never change, up to a number of entries.
type lru struct {
	mu      sync.Mutex
	size    int
	order   *list.List
	entries map[string]*list.Element
}

type entry struct {
	key   string
	value any
}

func newLRU(size int) *lru {
	return &lru{size: size, order: list.New(), entries: map[string]*list.Element{}}
}

// get returns key's value, computing and keeping it where it is not kept. A failed computation is not kept.
func (c *lru) get(key string, compute func() (any, error)) (any, error) {
	c.mu.Lock()
	if element, ok := c.entries[key]; ok {
		c.order.MoveToFront(element)
		c.mu.Unlock()
		return element.Value.(*entry).value, nil
	}
	c.mu.Unlock()
	value, err := compute()
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.entries[key]; !ok {
		c.entries[key] = c.order.PushFront(&entry{key, value})
		if c.order.Len() > c.size {
			oldest := c.order.Back()
			c.order.Remove(oldest)
			delete(c.entries, oldest.Value.(*entry).key)
		}
	}
	return value, nil
}
