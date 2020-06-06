package cache

func newLRU(size int) *lru {
	return &lru{
		size: size,
		head: nil,
		last: nil,
		hash: make(map[int]*lruNode),
	}
}

// lru implements Cache and iCache interfaces
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

// Cache interface

func (c *lru) Read(key int) (value int, isCacheMiss bool) {
	node := c.read(key)
	if node == nil {
		return 0, true
	}
	return node.value, false
}

func (c *lru) Write(key, value int) {
	_, _ = c.write(key, value)
}

// iCache interface

func (c *lru) read(key int) *lruNode {
	if node, exists := c.hash[key]; exists {
		c.promote(key)
		return node
	}
	return nil
}

func (c *lru) write(key, value int) (node, evicted *lruNode) {
	if node, exists := c.hash[key]; exists {
		node.value = value
		c.promote(key)
		return node, nil
	}
	node = c.insert(key, value)
	for c.isOverflowing() {
		evicted = c.remove(c.last.key)
	}
	return node, evicted
}

func (c *lru) remove(key int) *lruNode {
	node, found := c.hash[key]
	if !found {
		return nil
	}
	if node == c.head && node == c.last { // only one element in the cache
		c.head = nil
		c.last = nil
		delete(c.hash, node.key)
		return node
	}
	if node == c.head { // it's the head
		node.next.previous = nil
		c.head = node.next
		node.next = nil
		delete(c.hash, node.key)
		return node
	}
	if node == c.last { // it's the last node
		c.last = node.previous
		node.previous.next = nil
		node.previous = nil
		delete(c.hash, node.key)
		return node
	}
	// it's not the last nor the first.
	node.previous.next = node.next
	node.next.previous = node.previous
	node.next = nil
	node.previous = nil
	delete(c.hash, node.key)
	return node
}

// Helpers

// insert assumes node does not yet exist in the hash table.
func (c *lru) insert(key, value int) *lruNode {
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
	return newNode
}

func (c *lru) promote(key int) {
	node, exists := c.hash[key]
	if !exists {
		return
	}
	if node.previous == nil { // it's the head
		return
	}
	if node.next == nil { // it's the last node
		// detach from previous
		node.previous.next = nil
		c.last = node.previous
		node.previous = nil
		// move to the front of the head
		c.head.previous = node
		node.next = c.head
		c.head = node
		return
	}
	// it's not last nor head.
	// detach from previous and next
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

type lruState struct {
	size int
	list [][]int
	hash map[int]int
}

func (c *lru) printable() lruState {
	output := lruState{
		size: c.size,
		list: [][]int{},
		hash: make(map[int]int),
	}
	for node := c.head; node != nil; node = node.next {
		output.list = append(output.list, []int{node.key, node.value})
	}
	for key, node := range c.hash {
		output.hash[key] = node.value
	}
	return output
}
