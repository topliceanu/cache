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

// The Cache interface

func (c *lfu) Read(key int) (value int, isCacheMiss bool) {
	node := c.read(key)
	if node == nil {
		return 0, true
	}
	return node.value, false
}

func (c *lfu) Write(key, value int) {
	_, _ = c.write(key, value)
}

// The iCache interface

func (c *lfu) read(key int) *lfuNode {
	node, present := c.hash[key]
	if !present {
		return nil
	}
	c.increment(node)
	return node
}

func (c *lfu) write(key, value int) (node, evicted *lfuNode) {
	if node, present := c.hash[key]; present {
		node.value = value
		c.increment(node)
		return node, nil
	}
	node = &lfuNode{
		key:         key,
		value:       value,
		numRequests: 1,
	}
	if len(c.heap) >= c.size {
		evicted = c.remove(c.heap[len(c.heap) - 1].key)
	}
	c.hash[key] = node
	c.heap = heapPush(c.heap, node)
	return node, evicted
}

// increment will bump the numRequests property and promote the node in the heap.
// increment assumes the node is still in the cache.
func (c *lfu) increment(node *lfuNode) {
	node.numRequests++
	heapBubbleUp(c.heap, node.index)
}

// remove moves the node corresponding to key to the slot in the heap then
// bubbles down the interchanged value and resizes the heap.
func (c *lfu) remove(key int) *lfuNode {
	node, exists := c.hash[key]
	if !exists {
		return nil
	}

	index, lastIndex := node.index, len(c.heap) - 1
	c.heap[index], c.heap[lastIndex] = c.heap[lastIndex], c.heap[index]
	c.heap[index].index = index
	node.index = lastIndex

	heapBubbleDown(c.heap, index)

	delete(c.hash, node.key)
	c.heap = c.heap[:lastIndex]
	return node
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

func heapBubbleDown(heap []*lfuNode, parentIndex int) {
	leftIndex, rightIndex := parentIndex * 2 + 1, parentIndex * 2 + 2
	maxIndex := getMaxIndex(heap, parentIndex, leftIndex, rightIndex)
	if maxIndex == parentIndex {
		return
	}
	heap[maxIndex], heap[parentIndex] = heap[parentIndex], heap[maxIndex]
	heap[parentIndex].index = parentIndex
	heap[maxIndex].index = maxIndex

	heapBubbleDown(heap, maxIndex)
}

func getMaxIndex(heap []*lfuNode, parent, left, right int) int {
	maxIndex := parent
	for _, i := range([]int{left, right}) {
		if i >= len(heap) {
			continue
		}
		if (heap[maxIndex].value < heap[i].value) {
			maxIndex = i
		}
	}
	return maxIndex
}
