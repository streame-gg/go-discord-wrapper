package cache

import (
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var benchCounter atomic.Int64

func makeStoreWithN(n int) *genericStore[string, *discord.Guild] {
	s := newGenericStore[string, *discord.Guild](storeCfg{
		maxItems: n,
		clearBy:  ClearByAge,
	})
	for i := 0; i < n; i++ {
		s.set(strconv.Itoa(i), &discord.Guild{ID: discord.Snowflake(strconv.Itoa(i))})
	}
	return s
}

func addOne(s *genericStore[string, *discord.Guild]) {
	id := strconv.FormatInt(benchCounter.Add(1), 10)
	s.set(id, &discord.Guild{ID: discord.Snowflake(id)})
}

func BenchmarkEvictN_SmallStore_SingleEvict(b *testing.B) {
	s := makeStoreWithN(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		addOne(s)
	}
}

func BenchmarkEvictN_LargeStore_SingleEvict(b *testing.B) {
	s := makeStoreWithN(50000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		addOne(s)
	}
}

func BenchmarkEvictN_LargeStore_BulkEvict(b *testing.B) {
	s := makeStoreWithN(50000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.evictN(100, ClearByAge)
	}
}

func BenchmarkEvictN_WithConcurrentReads(b *testing.B) {
	s := makeStoreWithN(50000)
	someKey := "25000"

	stop := make(chan struct{})
	var readsDone atomic.Int64
	for r := 0; r < 8; r++ {
		go func() {
			for {
				select {
				case <-stop:
					return
				default:
					s.get(someKey)
					readsDone.Add(1)
				}
			}
		}()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.evictN(1, ClearByAge)
	}
	b.StopTimer()
	close(stop)
	elapsed := b.Elapsed()
	if elapsed > 0 {
		b.ReportMetric(float64(readsDone.Load())/elapsed.Seconds(), "reads/sec")
	}
}
