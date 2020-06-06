package cache

import (
	"testing"
)

func TestARC(t *testing.T) {
	t.Run("it uses the full capacity of the cache", func(t *testing.T) {
		c := newARC(8)
		c.Write(1, 10)
		c.Write(2, 20)
		c.Write(3, 30)
		c.Write(4, 40)
		value, isCacheMiss := c.Read(1)
		if isCacheMiss != false {
			t.Fatalf("expected read of key 1 to be a hit because it was stored in b1 but got value=%d, isCacheMiss=%t", value, isCacheMiss)
		}
		value, isCacheMiss = c.Read(1)
		if isCacheMiss == true || value != 10 {
			t.Fatalf("expected to read value for key 1 but got value=%d, isCacheMiss=%t", value, isCacheMiss)
		}
		value, isCacheMiss = c.Read(4)
		if isCacheMiss == true || value != 40 {
			t.Fatalf("expected to read value for key 4 but got value=%d, isCacheMiss=%t", value, isCacheMiss)
		}
		c.Write(5, 50)
		c.Write(6, 60)
		value, isCacheMiss = c.Read(5)
		if isCacheMiss == true || value != 50 {
			t.Fatalf("expected to read value for key 5 but got value=%d, isCacheMiss=%t", value, isCacheMiss)
		}
	})
}
