package benchmark

import (
	"testing"

	"github.com/topliceanu/cache"
)

func BenchmarkCache(b *testing.B) {
	for _, cacheType := range []string{
		cache.LRU,
		cache.LFU,
		cache.MRU,
		cache.SLRU,
		cache.LFRU,
		cache.ARC,
	} {
		b.Run(cacheType, func(b *testing.B) {
			c := cache.Factory(cacheType, 1000)

			b.Run("Write()", func(b *testing.B) {
				for i := 0; i < b.N; i ++ {
					c.Write(i, i)
				}
			})

			b.Run("Read()", func(b *testing.B) {
				for i := 0; i <= b.N; i ++ {
					_, _ = c.Read(i)
				}
			})
		})
	}
}

/*
import (
	"testing"

	"github.com/topliceanu/cache"
)

func TestCache(t *testing.T) {
	testCases = []*testCases{
		{
			cacheSize: 2,
			sequence: []*testOp{
				{
					typ: OpRead,
					key: 1,
					cacheMiss: true,
				},
				{
					typ: OpWrite,
					key: 1,
					value: 100,
				},
				{
					typ: OpRead,
					key: 1,
					value: 100,
				},
				{
					typ: OpWrite,
					key: 2,
					value: 200,
				},
				{
					typ: OpWrite,
					key: 3,
					value: 300,
				},
				{
					typ: OpRead,
					key: 1,
					cacheMiss: true,
				},
			},
		},
	}

	for strategy := range []string{ cache.LRU, cache.LFU, cache.MRU, cache.SLRU, cache.AR } {
		for tcIdx, tc := range testCases {
			c := cache.Factory(strategy, tc.cacheSize)
			for stepIdx, step := range tc.sequence {
				if step.typ == OpRead {
					value, cacheMiss = c.Read(step.key)
					if cacheMiss != step.cacheMiss {
						t.Fatalf("%s cache, test case #%d, step #%d, incorrect cache miss: expected=%b actual=%b",
							strategy, tcIdx, stepIdx, tc.cacheMiss, cacheMiss)
					}
					if value != step.value {
						t.Fatalf("%s cache, test case #%d, step #%d, incorrect read value: expected=%d actual=%d",
							strategy, tcIdx, stepIdx, tc.value, value)
					}
				} else if step.typ == OpWrite {
					c.Write(step.key)
				}
			}
		}
	}
}

type testCase struct {
	cacheSize int
	sequence []testOp
}

type opType int

const (
	opRead  opType = iota
	opWrite opType
)

type testOp struct {
	typ opType
	key int
	value int
	cacheMiss bool
}
*/
