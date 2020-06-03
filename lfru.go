package cache

// lfru implements Cache
type lfru struct {
	privileged   *lru
	unprivileged *lfu
}

func newLFRU(size int) *lfru {
	first, second := (size+1)/2, size/2
	return &lfru{
		privileged:   newLRU(first),
		unprivileged: newLFU(second),
	}
}

func (c *lfru) Read(key int) (int, bool) {
	// check privileged, if there, read and promote.
	value, cacheMiss := c.privileged.Read(key)
	if !cacheMiss {
		return value, false
	}
	// check unprivileged, if not there, cache miss.
	value, cacheMiss = c.unprivileged.Read(key)
	if cacheMiss {
		return 0, true
	}
	// otherwise  delete from unpriviledged, insert to privileged, handle overflow.
	_ = c.unprivileged.remove(key)
	_, evicted := c.privileged.write(key, value)
	if evicted != nil {
		c.unprivileged.Write(evicted.key, evicted.value)
	}
	return value, false
}

func (c *lfru) Write(key, value int) {
	// check privileged, if there, update and promote
	pnode := c.privileged.read(key)
	if pnode != nil {
		pnode.value = value
		return
	}
	// check unprivileged, if not there, insert and evict potential overflow.
	unode := c.unprivileged.read(key)
	if unode == nil {
		c.unprivileged.Write(key, value)
		return
	}
	// otherwise delete from unprivileged, insert into privileged,
	// then move potential node evicted from privileged into unprivileged.
	_ = c.unprivileged.remove(key)
	_, evicted := c.privileged.write(key, value)
	if evicted != nil {
		c.unprivileged.Write(evicted.key, evicted.value)
	}
}
