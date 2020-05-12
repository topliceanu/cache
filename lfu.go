package cache

type lfuNode struct {
	key         int
	value       int
	numRequests int
	index       int // index of the node in the heap
}

type lfu struct {
	hash map[int]*lfuNode
	heap []*lfuNode
	size int
}

func newLFU(size int) *lfu {
	return &lfu{
		hash: make(map[int]*lfuNode),
		heap: []*lfuNode{},
		size: size,
	}
}

func (l *lfu) Read(key int) (value int, isCacheMiss bool) {
	node, present := l.hash[key]
	if !present {
		return 0, true
	}
	l.increment(node)
	return node.value, false
}

func (l *lfu) Write(key, value int) {
	if node, present := l.hash[key]; present {
		node.value = value
		l.increment(node)
		return
	}
	node := &lfuNode{
		key:         key,
		value:       value,
		numRequests: 1,
	}
	l.hash[key] = node
	if len(l.heap) >= l.size {
		l.evict()
	}
	l.heap = heapPush(l.heap, node)
}

// increment will bump the numRequests property and promote the node in the heap.
// increment assumes the node is still in the cache.
func (l *lfu) increment(node *lfuNode) {
	node.numRequests++
	heapBubbleUp(l.heap, node.index)
}

// evict removes the least frequently requested value from the cache.
// evict does not check if the cache is over capacity, callers should do that!
func (l *lfu) evict() {
	node := l.heap[len(l.heap)-1]
	delete(l.hash, node.key)
	l.heap = l.heap[:len(l.heap)-1]
}

// Max heap data structure is modeled as a slice of *lfuNode and maintains the
// node with the largest numRequests at the head of the array.

// heapPush inserts a new node in the heap, preserving the heap invariant.
// heapPush maintains the index property of each node
func heapPush(heap []*lfuNode, node *lfuNode) []*lfuNode {
	heap = append(heap, node)
	node.index = len(heap) - 1
	heapBubbleUp(heap, node.index)
	return heap
}

func heapBubbleUp(heap []*lfuNode, index int) {
	parentIndex := (index - 1) / 2
	if parentIndex < 0 {
		return
	}
	if heap[parentIndex].numRequests >= heap[index].numRequests {
		return
	}
	heap[parentIndex], heap[index] = heap[index], heap[parentIndex]
	heap[parentIndex].index = parentIndex
	heap[index].index = index
	heapBubbleUp(heap, parentIndex)
}
