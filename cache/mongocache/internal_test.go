package mongocache

import (
	"strconv"
	"sync"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// TestBug114_ExtractID_UnknownTypePanics verifies that extractID panics for any
// document type it does not recognise.
//
// Root cause: the missing default branch caused extractID to silently return ""
// for unknown types, so all such documents shared _id:"" in MongoDB and
// overwrote each other. The fix adds a panic so programmer errors surface
// immediately rather than corrupting the database.
func TestBug114_ExtractID_UnknownTypePanics(t *testing.T) {
	type unknownDoc struct{ X string }
	defer func() {
		if r := recover(); r == nil {
			t.Error("Bug114: extractID with unknown type should panic, but did not")
		}
	}()
	extractID(unknownDoc{"hello"})
}

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
		go func() {
			defer wg.Done()
			// LoadOrStore directly, mimicking Add().
			c.msgChannelMu.LoadOrStore(strconv.Itoa(i), &sync.Mutex{})
			// Then DeleteChannel removes it.
			store := &mongoMessageStore{c: c}
			store.DeleteChannel(discord.Snowflake(i))
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
