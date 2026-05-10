package cache

import (
	"sync"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// TestBug11TotalBytesNeverNegative verifies that the subtract(-old) and add(+new)
// operations on totalBytes happen inside the same lock, so no concurrent reader
// can observe a transient undercount between the two ops.
func TestBug11TotalBytesNeverNegative(t *testing.T) {
	s := newGenericStore[string, *common.Guild](storeCfg{trackBytes: true})

	// Seed one entry so every subsequent Set is an update (subtract+add path).
	s.set("1", &common.Guild{ID: "1", Name: "initial"})

	negative := make(chan int64, 1)
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				if v := s.totalBytes.Load(); v < 0 {
					select {
					case negative <- v:
					default:
					}
				}
			}
		}
	}()

	const goroutines = 500
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			s.set("1", &common.Guild{ID: "1", Name: "updated"})
		}()
	}
	wg.Wait()
	close(stop)

	select {
	case v := <-negative:
		t.Errorf("totalBytes went negative (%d) during concurrent updates — non-atomic byte accounting (Bug 11)", v)
	default:
	}
}
