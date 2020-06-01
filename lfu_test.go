package cache

import (
	"testing"
)

func TestLFU(t *testing.T) {
	t.Run("cache state for base cases is correct", func(t *testing.T) {
		c := newLFU(2)
		if len(c.heap) != 0 || len(c.hash) != 0 || c.size != 2 {
			t.Fatalf("empty cache state is incorrect: %#v", c)
		}
		c.Write(1, 10)
		if len(c.heap) != 1 || len(c.hash) != 1 || c.heap[0].key != 1 ||
			c.heap[0].value != 10 || c.heap[0].numRequests != 1 || c.heap[0].index != 0 {
			t.Fatalf("cache with one value has incorrect state: %#v", c)
		}
		c.Read(1)
		if len(c.heap) != 1 || len(c.hash) != 1 || c.heap[0].key != 1 ||
			c.heap[0].value != 10 || c.heap[0].numRequests != 2 || c.heap[0].index != 0 {
			t.Fatalf("cache with one value has incorrect state: %#v", c.heap[0])
		}
	})
	t.Run("cache state after eviction is correct", func(t *testing.T) {
		c := newLFU(2)
		c.Write(1, 10)
		if len(c.heap) != 1 || len(c.hash) != 1 || c.heap[0].key != 1 ||
			c.heap[0].value != 10 || c.heap[0].numRequests != 1 || c.heap[0].index != 0 {
			t.Fatalf("cache with one value has incorrect state: %#v", c.heap[0])
		}
		c.Write(2, 20)
		if len(c.heap) != 2 || len(c.hash) != 2 || c.heap[1].key != 2 ||
			c.heap[1].value != 20 || c.heap[1].numRequests != 1 || c.heap[1].index != 1 {
			t.Fatalf("cache with two values has incorrect state: %#v", c.heap)
		}
		actualValue, isCacheMiss := c.Read(1) // {1, 2}
		if actualValue != 10 || isCacheMiss == true || len(c.heap) != 2 ||
			len(c.hash) != 2 || c.heap[0].key != 1 || c.heap[0].value != 10 ||
			c.heap[0].numRequests != 2 || c.heap[0].index != 0 {
			t.Fatalf("cache after a read has incorrect state: %#v", c.heap[0])
		}
		c.Write(3, 30) // {1, 3}
		if len(c.heap) != 2 || len(c.hash) != 2 || c.heap[0].key != 1 || c.heap[0].value != 10 ||
			c.heap[1].key != 3 || c.heap[1].value != 30 || c.heap[1].numRequests != 1 ||
			c.heap[1].index != 1 {
				t.Fatalf("cache after another write has incorrect state: #0:%#v #1:%#v", c.heap[0], c.heap[1])
		}
		actualValue, isCacheMiss = c.Read(2)
		if isCacheMiss != true || actualValue != 0 {
			t.Fatalf("cache does not return a miss for a key that should have been evicted: %v", c.heap)
		}
	})
	t.Run("check cache state after multiple reads and writes", func(t *testing.T) {
		c := newLFU(2)
		c.Write(2, 20)
		c.Write(1, 10)
		node := c.remove(2)
		if node.key != 2 || node.value != 20 || len(c.heap) != 1 ||
		c.heap[0].key != 1 || c.heap[0].value != 10 ||
		c.heap[0].numRequests != 1 || c.heap[0].index != 0 ||
		c.hash[1].key != 1 || c.hash[1].value != 10 {
			t.Fatalf("cache is in an inconsistent state %#v", c.heap)
		}
	})
}
