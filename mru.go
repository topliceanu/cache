package cache

func newMRU(size int) *mru {
	return &mru{
		size: size,
		head: nil,
		last: nil,
		hash: make(map[int]*mruNode),
	}
}

// mru implements Cache
// MRU evicts the most recently used key, ie. the key that was just requested.
type mru struct {
	size int
	head *mruNode
	last *mruNode
	hash map[int]*mruNode
}

type mruNode struct {
	key      int
	value    int
	next     *mruNode
	previous *mruNode
}

func (m *mru) Write(key, value int) {
	if node, exists := m.hash[key]; exists {
		node.value = value
		m.promote(key)
		return
	}
	if m.size == len(m.hash) {
		m.evict()
	}
	node := &mruNode{
		key:   key,
		value: value,
	}
	if len(m.hash) == 0 {
		m.head = node
		m.last = node
		m.hash[key] = node
		return
	}
	node.next = m.head
	m.head.previous = node
	m.head = node
	m.hash[key] = node
}

func (m *mru) Read(key int) (value int, isCacheMiss bool) {
	if node, exists := m.hash[key]; exists {
		// to evict the key we just read, we promote it then evict the head.
		m.promote(key)
		m.evict()
		return node.value, false
	}
	return 0, true
}

// promote makes the node matching the given key, the head of the doubly-linked list.
func (m *mru) promote(key int) {
	node, exists := m.hash[key]
	if !exists {
		return
	}
	if node == m.head {
		return
	}
	if node == m.last {
		m.last = node.previous
		// detach from previous
		node.previous.next = nil
		node.previous = nil
		// move to head
		node.next = m.head
		m.head.previous = node
		m.head = node
		return
	}
	node.previous.next = node.next
	node.next.previous = node.previous
	m.head.previous = node
	node.next = m.head
	m.head = node
}

// evict pops the head of the doubly-linked list.
// evict assumes you check that the list is not empty before you called it.
func (m *mru) evict() {
	if len(m.hash) == 0 {
		return
	}
	if len(m.hash) == 1 {
		node := m.head
		m.head = nil
		m.last = nil
		delete(m.hash, node.key)
		return
	}
	delete(m.hash, m.head.key)
	m.head.next.previous = nil
	m.head = m.head.next
}
