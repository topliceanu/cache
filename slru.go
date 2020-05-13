package cache

type slru struct {
	protected *lru
	probation *lru
}

func newSLRU(size int) *slru {
	first, second := (size + 1) / 2, size / 2
	return &slru{
		protected: newLRU(first),
		probation: newLRU(second),
	}
}

func (c *slru) Read(key int) (int, bool) {
	// Search the protected section.
	value, cacheMiss := c.protected.Read(key)
	if !cacheMiss {
		return value, false
	}
	// Search the probation section.
	value, cacheMiss = c.probation.Read(key)
	if cacheMiss {
		return 0, true
	}
	// Promote from probation to protected taking care of any
	// evicted overflow from protected.
	_ = c.probation.remove(key)
	_, evicted := c.protected.write(key, value)
	if evicted != (*lruNode)(nil) {
		c.probation.Write(evicted.(*lruNode).key, evicted.(*lruNode).value)
	}
	return value, false
}

func (c *slru) Write(key, value int) {
	// Key is in protected so we update the value.
	node := c.protected.read(key)
	if node != nil {
		node.(*lruNode).value = value
		return
	}
	// Key is not in probation so we write the new page in probation.
	node = c.probation.read(key)
	if node == nil {
		c.probation.Write(key, value)
		return
	}
	// Key is in probabation. We move the key to protected, update the
	// value and handle any overflow from protected.
	_ = c.probation.remove(key)
	_, evicted := c.protected.write(key, value)
	if evicted != (*lruNode)(nil) {
		c.probation.Write(evicted.(*lruNode).key, evicted.(*lruNode).value)
	}
}
