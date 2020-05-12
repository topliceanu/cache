package cache

import (
	"testing"
)

func TestMRU(t *testing.T) {
	t.Run("base-case caches have correct states", func(t *testing.T) {
		// That is, empty cache and cache with one node.
		c := newMRU(2)
		if len(c.hash) != 0 || c.head != nil || c.last != nil || c.size != 2 {
			t.Fatalf("incorrect empty cache state: %#v", c)
		}
		c.Write(1, 2)
		if len(c.hash) != 1 || c.head != c.last || c.head.key != 1 || c.head.value != 2 || c.head.previous != nil || c.head.next != nil {
			t.Fatalf("incorrect state for cache with only one value in: %#v", c.head)
		}
		actualValue, isCacheMiss := c.Read(1)
		if actualValue != 2 || isCacheMiss == true {
			t.Fatalf("unexpected return from a Read(): value=%d isCacheMis=%t", actualValue, isCacheMiss)
		}
	})
	t.Run("cache correctly evicts most recently requested key", func(t *testing.T) {
		c := newMRU(2)
		c.Write(1, 10)
		c.Write(2, 20)
		if len(c.hash) != 2 || c.head.key != 2 || c.last.key != 1 {
			t.Fatalf("unexpected cache state before any evictions: head=%#v, last=%#v", c.head, c.last)
		}
		c.Write(3, 30)
		if len(c.hash) != 2 || c.head.key != 3 || c.last.key != 1 {
			t.Fatalf("unexpected cache state after an eviction: head=%#v, last=%#v", c.head, c.last)
		}
	})
}
