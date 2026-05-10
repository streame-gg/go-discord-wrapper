package cache

import (
	"sync"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// TestBug10SetAllIsAtomicNoEmptyWindowForReaders verifies that concurrent reads
// during SetAll never observe a zero-entry state. Without Bug 10's fix, the lock
// is released between the delete loop and the insert loop, leaving a window
// where GetByGuild returns an empty slice.
func TestBug10SetAllIsAtomicNoEmptyWindowForReaders(t *testing.T) {
	es := &memEmojiStore{
		s: newGenericStore[common.Snowflake, emojiValue](storeCfg{}),
	}

	guildID := common.Snowflake("123456789012345")
	// Seed two emojis.
	es.s.set("e1", emojiValue{GuildID: guildID, Emoji: &common.Emoji{ID: "e1"}})
	es.s.set("e2", emojiValue{GuildID: guildID, Emoji: &common.Emoji{ID: "e2"}})

	sawEmpty := make(chan struct{}, 1)
	stop := make(chan struct{})

	// Reader goroutine: continuously polls GetByGuild.
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				if len(es.GetByGuild(guildID)) == 0 {
					select {
					case sawEmpty <- struct{}{}:
					default:
					}
				}
			}
		}
	}()

	// Writer: repeatedly calls SetAll with the same two emojis.
	const iters = 500
	var wg sync.WaitGroup
	wg.Add(iters)
	for i := 0; i < iters; i++ {
		go func() {
			defer wg.Done()
			es.SetAll(guildID, []*common.Emoji{
				{ID: "e1"},
				{ID: "e2"},
			})
		}()
	}
	wg.Wait()
	close(stop)

	select {
	case <-sawEmpty:
		t.Error("reader saw empty emoji list during SetAll — delete→insert not atomic (Bug 10)")
	default:
	}
}

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
