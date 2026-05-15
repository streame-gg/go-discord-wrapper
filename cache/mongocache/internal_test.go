package mongocache

import (
	"fmt"
	"sync"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// TestBug21MsgChannelMuCleanedOnDeleteChannel verifies that per-channel mutexes
// are removed from msgChannelMu when DeleteChannel is called, preventing an
// unbounded memory leak for bots with many short-lived channels.
func TestBug21MsgChannelMuCleanedOnDeleteChannel(t *testing.T) {
	c := &MongoDBCache{
		opts: cache.Options{
			Messages: cache.MessageOptions{MaxPerChannel: 10},
		},
	}
	// We don't need a real MongoDB connection; just exercise the mutex lifecycle.
	// Seed msgChannelMu with 10 000 fake channel entries.
	const n = 10_000
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			cid := discord.Snowflake(fmt.Sprintf("ch%d", i))
			// LoadOrStore directly, mimicking Add().
			c.msgChannelMu.LoadOrStore(string(cid), &sync.Mutex{})
			// Then DeleteChannel removes it.
			store := &mongoMessageStore{c: c}
			store.DeleteChannel(cid)
		}()
	}
	wg.Wait()

	// After all DeleteChannel calls, the map must be empty.
	count := 0
	c.msgChannelMu.Range(func(_, _ any) bool { count++; return true })
	if count != 0 {
		t.Errorf("msgChannelMu has %d entries after DeleteChannel — mutex leak (Bug 21)", count)
	}
}
