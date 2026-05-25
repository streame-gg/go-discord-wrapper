package stores

import (
	"testing"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// J0-#13: MemMessageStore.Channel() must return copies, not shared pointers.
func TestMemMessageStore_ChannelReturnsCopies(t *testing.T) {
	s := NewMessageStore(MessageDefaults())
	defer s.Close()

	chID := discord.Snowflake(111222333444555)
	msgID := discord.Snowflake(555444333222111)

	s.Add(&discord.Message{ID: msgID, ChannelID: chID, Content: "hello"})

	got := s.Channel(chID)
	msg, ok := got.Get(msgID)
	if !ok {
		t.Fatal("message not found in channel")
	}

	// Mutate the returned copy.
	msg.Content = "mutated"

	// Re-fetch: the cache must still hold the original value.
	got2 := s.Channel(chID)
	msg2, ok2 := got2.Get(msgID)
	if !ok2 {
		t.Fatal("message not found on second fetch")
	}
	if msg2.Content != "hello" {
		t.Fatalf("cache was mutated: got %q, want %q", msg2.Content, "hello")
	}
}
