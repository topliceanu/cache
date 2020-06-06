package cache

// arc is an adaptation of the ARC algorithm for the Cache interface.
type arc struct {
	t1, b1, t2, b2 *lru
	p              int
	c              int
}

func newARC(size int) *arc {
	c := size
	t1Size, t2Size := c/2, (c+1)/2
	b1Size, b2Size := c, c
	return &arc{
		c:  c,
		p:  0,
		t1: newLRU(t1Size),
		t2: newLRU(t2Size),
		b1: newLRU(b1Size),
		b2: newLRU(b2Size),
	}
}

func (a *arc) Read(key int) (value int, isCacheMiss bool) {
	// Case I:
	if value, isCacheMiss = a.t2.Read(key); !isCacheMiss {
		return value, false
	}
	if _, isCacheMiss = a.t1.Read(key); !isCacheMiss {
		node := a.t1.remove(key)
		if _, evicted := a.t2.write(node.key, node.value); evicted != nil {
			a.b2.Write(evicted.key, evicted.value)
		}
		return node.value, false
	}
	// Case II:
	if _, isCacheMiss = a.b1.Read(key); !isCacheMiss {
		b1Size, b2Size := len(a.b1.hash), len(a.b2.hash)
		a.p = min(a.c, a.p + max(b2Size / b1Size, 1))
		a.replace(key)
		node := a.b1.remove(key)
		if _, evicted := a.t2.write(node.key, node.value); evicted != nil {
			a.b2.Write(evicted.key, evicted.value)
		}
		return node.value, false
	}
	// Case III:
	if _, isCacheMiss = a.b2.Read(key); !isCacheMiss {
		b1Size, b2Size := len(a.b1.hash), len(a.b2.hash)
		a.p = max(0, a.p-max(b1Size/b2Size, 1))
		a.replace(key)
		node := a.b2.remove(key)
		if _, evicted := a.t2.write(node.key, node.value); evicted != nil {
			a.b2.Write(evicted.key, evicted.value)
		}
		return node.value, false
	}
	// Case IV:
	t1Size := len(a.t1.hash)
	l1Size, l2Size := len(a.t1.hash)+len(a.b1.hash), len(a.t2.hash)+len(a.b2.hash)
	if l1Size == a.c {
		if t1Size < a.c {
			_ = a.b1.remove(a.b1.last.key)
			a.replace(key)
		} else {
			node := a.t1.remove(a.t1.last.key)
			a.b1.Write(node.key, node.value)
		}
	}
	if l1Size < a.c && l1Size+l2Size >= a.c {
		if l1Size+l2Size == 2*a.c {
			a.b2.remove(a.b2.last.key)
		}
		a.replace(key)
	}
	return 0, true
}

func (a *arc) Write(key, value int) {
	// if it exists in t2, update value and promote to the head of the LRU.
	if node := a.t2.read(key); node != nil {
		node.value = value
		return
	}
	// if it exists in t1, remove it from t1, add it to t2, handle eviction from t2 into b2.
	if node := a.t1.read(key); node != nil {
		a.t1.remove(key)
		if _, evicted := a.t2.write(key, value); evicted != nil {
			a.b2.Write(evicted.key, evicted.value)
		}
		return
	}
	// if it doesn't exist in t1 or t2, insert it in t1, handle eviction from t1 into b1.
	if _, evicted := a.t1.write(key, value); evicted != nil {
		a.b1.write(evicted.key, evicted.value)
		return
	}
}

func (a *arc) replace(key int) {
	t1Size := len(a.t1.hash)
	_, cacheMiss := a.b2.Read(key)
	if t1Size >= 1 && ((!cacheMiss && t1Size == a.p) || t1Size > a.p) {
		if a.t1.last == nil {
			return
		}
		if node := a.t1.remove(a.t1.last.key); node != nil {
			a.b1.Write(node.key, node.value)
		}
	} else {
		if a.t2.last == nil {
			return
		}
		if node := a.t2.remove(a.t2.last.key); node != nil {
			a.b2.Write(node.key, node.value)
		}
	}
}

// Helpers

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type arcState struct {
		c int
		p int
		t1 [][]int
		t2 [][]int
		b1 [][]int
		b2 [][]int
}

func (a *arc) printable() arcState {
	return arcState{
		c: a.c,
		p: a.p,
		t1: a.t1.printable().list,
		t2: a.t2.printable().list,
		b1: a.b1.printable().list,
		b2: a.b2.printable().list,
	}
}
