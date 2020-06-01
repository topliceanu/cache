package cache

import (
	"testing"
)

func TestLFRU(t *testing.T) {
	t.Run("cache state should be correct for a cache with two elements in each side", func(t *testing.T) {
		c := newLFRU(4)
		c.Write(1, 10)
		if len(c.privileged.hash) != 0 || len(c.unprivileged.hash) != 1 ||
		c.unprivileged.hash[1].key != 1 || c.unprivileged.hash[1].value != 10 ||
		c.unprivileged.hash[1].numRequests != 1 || c.unprivileged.hash[1].index != 0 {
			t.Fatalf("incorrect cache state after first write: %#v", c.unprivileged)
		}
		c.Write(2, 20)
		if len(c.privileged.hash) != 0 || len(c.unprivileged.hash) != 2 ||
		c.unprivileged.hash[1].key != 1 || c.unprivileged.hash[1].value != 10 ||
		c.unprivileged.hash[1].numRequests != 1 || c.unprivileged.hash[1].index != 0 ||
		c.unprivileged.hash[2].key != 2 || c.unprivileged.hash[2].value != 20 ||
		c.unprivileged.hash[2].numRequests != 1 || c.unprivileged.hash[2].index != 1 {
			t.Fatalf("incorrect cache state after the second write: %#v", c.unprivileged)
		}
		value, cacheMiss := c.Read(2)
		if value != 20 || cacheMiss != false ||
		len(c.privileged.hash) != 1 || len(c.unprivileged.hash) != 1 ||
		c.privileged.hash[2].key != 2 || c.privileged.hash[2].value != 20 ||
		c.privileged.hash[2].next != nil || c.privileged.hash[2].previous != nil ||
		c.unprivileged.hash[1].key != 1 || c.unprivileged.hash[1].value != 10 ||
		c.unprivileged.hash[1].numRequests != 1 || c.unprivileged.hash[1].index != 0 {
			t.Fatalf("incorrect cache state after the a read: %#v, %#v", c.privileged, c.unprivileged)
		}
	})
	t.Run("cache state after multiple evictions and promotions is correct", func(t *testing.T) {
		c := newLFRU(4)
		c.Write(1, 10)                // (_, 1); (_, _)
		c.Write(2, 20)                // (2, 1); (_, _)
		c.Write(3, 30)                // (3, 1); (_, _)
		_, _ = c.Read(1)              // (_, 3); (_, 1)
		_, _ = c.Read(3)              // (_, _); (3, 1)
		c.Write(1, 10)                // (_, 1); (2, 3)
		c.Write(4, 40)                // (4, 1); (2, 3)
		value, cacheMiss := c.Read(1) // (4, 2); (1, 3)
		if value != 10 || cacheMiss != false {
			t.Fatalf("failed to get the correct value after a read: %#v", c.unprivileged)
		}
	})
}
