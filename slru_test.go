package cache

import (
	"testing"
)

func TestSLRU(t *testing.T) {
	t.Run("empty cache state is correct", func(t *testing.T) {
		c := newSLRU(2)
		c.Write(1, 10)
		if len(c.protected.hash) != 0 || len(c.probation.hash) != 1 ||
		c.probation.hash[1].key != 1 || c.probation.hash[1].value != 10 ||
		c.probation.hash[1].next != nil || c.probation.hash[1].previous != nil {
			t.Fatalf("unexpected cache state after first write: %#v", c.probation.hash)
		}
		c.Write(2, 20)
		if len(c.protected.hash) != 0 || len(c.probation.hash) != 1 ||
		c.probation.hash[2].key != 2 || c.probation.hash[2].value != 20 ||
		c.probation.hash[2].next != nil || c.probation.hash[2].previous != nil {
			t.Fatalf("unexpected cache state eviction from probation: %#v", c.probation.hash)
		}
		value, cacheMiss := c.Read(2)
		if value != 20 || cacheMiss == true ||
		len(c.protected.hash) != 1 || len(c.probation.hash) != 0 ||
		c.protected.hash[2].key != 2 || c.protected.hash[2].value != 20 ||
		c.protected.hash[2].next != nil || c.protected.hash[2].previous != nil {
			t.Fatalf("unexpected cache state after a promotion from probation to protected: %#v", c.protected)
		}
		c.Write(1, 10)
		if len(c.protected.hash) != 1 || len(c.probation.hash) != 1 ||
		c.probation.hash[1].key != 1 || c.probation.hash[1].value != 10 ||
		c.probation.hash[1].next != nil || c.probation.hash[1].previous != nil ||
		c.protected.hash[2].key != 2 || c.protected.hash[2].value != 20 ||
		c.protected.hash[2].next != nil || c.protected.hash[2].previous != nil {
			t.Fatalf("unexpected cache state after a promotion from probation to protected: %#v, %#v", c.protected, c.probation)
		}
		value, cacheMiss = c.Read(1)
		if value != 10 || cacheMiss == true ||
		len(c.protected.hash) != 1 || len(c.probation.hash) != 1 ||
		c.probation.hash[2].key != 2 || c.probation.hash[2].value != 20 ||
		c.probation.hash[2].next != nil || c.probation.hash[2].previous != nil ||
		c.protected.hash[1].key != 1 || c.protected.hash[1].value != 10 ||
		c.protected.hash[1].next != nil || c.protected.hash[1].previous != nil {
			t.Fatalf("unexpected cache state after another read and an eviction from protected: %#v, %#v", c.protected, c.probation)
		}
	})
	t.Run("cache state after multiple evictions and promotions is correct", func(t *testing.T) {
		c := newSLRU(4)
		c.Write(1, 10)                // (_, 1); (_, _)
		c.Write(2, 20)                // (1, 2); (_, _)
		c.Write(3, 30)                // (2, 3); (_, _)
		_, _ = c.Read(2)              // (_, 3); (_, 2)
		_, _ = c.Read(3)              // (_, _); (2, 3)
		c.Write(1, 10)                // (_, 1); (2, 3)
		c.Write(4, 40)                // (1, 4); (2, 3)
		value, cacheMiss := c.Read(1) // (4, 2); (3, 1)
		if value != 10 || cacheMiss == true ||
		len(c.protected.hash) != 2 || len(c.probation.hash) != 2 ||
		c.probation.hash[2].key != 2 || c.probation.hash[2].value != 20 ||
		c.probation.hash[4].key != 4 || c.probation.hash[4].value != 40 ||
		c.probation.hash[2].next.key != 4 || c.probation.hash[4].previous.key != 2 ||
		c.protected.hash[1].key != 1 || c.protected.hash[1].value != 10 ||
		c.protected.hash[3].key != 3 || c.protected.hash[3].value != 30 ||
		c.protected.hash[1].next.key != 3 || c.protected.hash[3].previous.key != 1 {
			t.Fatalf("unexpected cache state after a series of read/writes: %#v, %#v", c.protected, c.probation)
		}
	})
}
