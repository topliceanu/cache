package cache

import (
	"reflect"
	"testing"
)

func TestLRU(t *testing.T) {
	t.Run("empty cache has correct state", func(t *testing.T) {
		c := newLRU(2)
		if c.size != 2 {
			t.Errorf("expected size of new cache to be %d but is %d", 2, c.size)
		}
		if c.head != nil {
			t.Errorf("expected head of cache to be nil but is %#v", c.head)
		}
		if c.last != nil {
			t.Errorf("expected last of cache to be nil but is %#v", c.last)
		}
		if len(c.hash) != 0 {
			t.Errorf("expected size of hashtable to be 0 but is %#v", len(c.hash))
		}
	})
	t.Run("cache with one value has correct state", func(t *testing.T) {
		c := newLRU(2)
		c.Write(10, 100)
		if c.head == nil || c.head.key != 10 && c.head.value != 100 {
			t.Errorf("expected head of cache to be (10, 100) but is %#v", c.head)
		}
		if c.last == nil || c.last != c.head || c.last.key != 10 || c.last.value != 100 {
			t.Errorf("expected last of cache to be the same as the head but is %#v", c.last)
		}
		if len(c.hash) != 1 || c.hash[10].value != 100 {
			t.Errorf("expected hashtable to have one key but is %#v", c.hash)
		}
	})
	t.Run("cache with two values has correct state", func(t *testing.T) {
		c := newLRU(2)
		c.Write(10, 100)
		c.Write(20, 200)
		if c.head == nil || c.head.key != 20 && c.head.value != 200 {
			t.Errorf("expected head of cache to be (20, 200) but is %#v", c.head)
		}
		if c.last == nil || c.last == c.head || c.last.key != 10 || c.last.value != 100 {
			t.Errorf("expected last of cache to be different than the head but is %#v", c.last)
		}
		if c.head.next != c.last || c.last.previous != c.head {
			t.Error("expected head and last to be linked by references")
		}
		if len(c.hash) != 2 || c.hash[10].value != 100 || c.hash[20].value != 200 {
			t.Errorf("expected hashtable to have two keys but is %#v", c.hash)
		}
	})
	t.Run("cache correctly evicts least recently used value", func(t *testing.T) {
		c := newLRU(2)
		c.Write(10, 100)
		c.Write(20, 200)
		c.Write(30, 300)
		if c.head == nil || c.head.key != 30 && c.head.value != 300 {
			t.Errorf("expected head of cache to be (30, 300) but is %#v", c.head)
		}
		if c.last == nil || c.last == c.head || c.last.key != 20 || c.last.value != 200 {
			t.Errorf("expected last of cache to be different than the head but is %#v", c.last)
		}
		if c.head.next != c.last || c.last.previous != c.head {
			t.Error("expected head and last to be linked by references")
		}
		if len(c.hash) != 2 || c.hash[10] != nil || c.hash[20].value != 200 || c.hash[30].value != 300 {
			t.Errorf("expected hashtable to have two keys but is %#v", c.hash)
		}
	})
	t.Run(".remove() correctly delete a key in the middle of the linked list", func(t *testing.T) {
		c := newLRU(3)
		c.Write(10, 100)
		c.Write(20, 200)
		c.Write(30, 300)
		node := c.remove(20)
		if node.key != 20 || node.value != 200 {
			t.Fatalf("expected the correct reeturn value from remove but got %#v", node)
		}
		state := c.printable()
		if !reflect.DeepEqual(state.list, [][]int{ {30, 300}, {10, 100} }) {
			t.Fatalf("unexpected cache linked list state: %#v", state.list)
		}
		if !reflect.DeepEqual(state.hash, map[int]int{ 30: 300, 10: 100 }) {
			t.Fatalf("unexpected cache map state: %#v", state.hash)
		}
	})
}
