package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/topliceanu/cache"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	var (
		// size of the input set
		m = 1000000
		// cardinality of the input set
		n = 10000
		// cache size
		k = 1000
		// input set generated randomly
		values = generate(n, m)
		// all the caches under test
		cacheTypes = []string{ cache.LRU, cache.LFU, cache.MRU, cache.SLRU, cache.LFRU, cache.ARC }
		hitRate float64
		missRate float64
	)
	fmt.Printf("Cache type    Hit rate    Miss rate \n")
	for _, cacheType := range cacheTypes {
		c := cache.Factory(cacheType, k)
		misses := 0
		for _, value := range values {
			_, isCacheMiss := c.Read(value)
			if isCacheMiss {
				misses += 1
				c.Write(value, value)
			}
		}
		hitRate = float64(m - misses) / float64(m) * 100
		missRate = float64(misses) / float64(m) * 100
		fmt.Printf("%10s    %2.3f      %2.3f\n", cacheType, hitRate, missRate)
	}
}

func generate(cardinality, length int) []int {
	out := make([]int, length)
	for i := 0; i < length; i ++ {
		out[i] = rand.Intn(cardinality + 1)
	}
	return out
}
