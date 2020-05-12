package cache

func newLRU(size int) *lru {
	return &lru{
		size: size,
		head: nil,
		last: nil,
		hash: make(map[int]*lruNode),
	}
}

// lru implements Cache
// LRU evicts the least-recently used key.
type lru struct {
	size int
	head *lruNode
	last *lruNode
	hash map[int]*lruNode
}

type lruNode struct {
	key      int
	value    int
	next     *lruNode
	previous *lruNode
}

func (c *lru) Read(key int) (int, bool) {
	if _, exists := c.hash[key]; exists {
		c.promote(key)
		return c.hash[key].value, true
	}
	return 0, false
}

func (c *lru) Write(key, value int) {
	if _, exists := c.hash[key]; exists {
		c.hash[key].value = value
		c.promote(key)
		return
	}
	c.insert(key, value)
	for c.isOverflowing() {
		c.evictLast()
	}
}

// insert assumes node does not yet exist in the hash table.
func (c *lru) insert(key, value int) {
	newNode := &lruNode{
		key:      key,
		value:    value,
		previous: nil,
		next:     nil,
	}
	c.hash[key] = newNode
	if c.head == nil {
		c.head = newNode
		c.last = newNode
	} else {
		newNode.next = c.head
		c.head.previous = newNode
		c.head = newNode
	}
}

func (c *lru) promote(key int) {
	node := c.hash[key]
	if node.previous == nil { // it's the head
		return
	}
	if node.next == nil { // it's the last node
		// detach from previous
		node.previous.next = nil
		// move to the front of the head
		c.head.previous = node
		node.next = c.head
		node.previous = nil
		c.head = node
		c.last = node
	}
	// detach from the head
	node.previous.next = node.next
	node.next.previous = node.previous
	// move to the front of the head
	node.previous = nil
	node.next = c.head
	c.head.previous = node
	c.head = node

}

func (c *lru) isOverflowing() bool {
	return len(c.hash) > c.size
}

func (c *lru) evictLast() {
	delete(c.hash, c.last.key)
	c.last.previous.next = nil
	c.last = c.last.previous
}
