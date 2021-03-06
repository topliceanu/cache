package cache

import (
	"fmt"
)

// Cache is the main interface implemented by all strageties in this project.
type Cache interface {
	Read(key int) (value int, isCacheMiss bool)
	Write(key, value int)
}

// iCache is an internal interface for cache implementations to expose the data structures used.
// It's helpful for combining different caches into more complex algorithms, like SLRU, LFRU or AR.
// It's only for documentation purposes, cache implementations will return, for convenience,
// their respective node types, not interfaces{}.
type iCache interface {
	// read and promote if cache hit
	read(key int) (node interface{})
	// write new page or update existing one. Either case, promote.
	write(key, value int) (node, evicted interface{})
	// remove a page by key. Nothing happens if key is not found.
	remove(key int) (node interface{})
	// expose cache's internal state. Should only be used in tests!
	state() interface{}
}

// Supress the linter
var _ iCache

const (
	// Cache replacement strategies
	LRU  = "cache-lru"
	LFU  = "cache-lfu"
	MRU  = "cache-mru"
	SLRU = "cache-slru"
	LFRU = "cache-lfru"
	ARC  = "cache-arc"
)

// Factory produces instances of the requested cache replacement strategy.
func Factory(algorithm string, size int) Cache {
	switch algorithm {
	case LRU:
		return newLRU(size)
	case MRU:
		return newMRU(size)
	case LFU:
		return newLFU(size)
	case SLRU:
		return newSLRU(size)
	case LFRU:
		return newLFRU(size)
	case ARC:
		return newARC(size)
	default:
		panic(fmt.Sprintf("unsupported caching algorithm %s", algorithm))
	}
}
