package cache

import (
	"testing"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func BenchmarkGuildStore_Set(b *testing.B) {
	c := NewMemoryCache(Options{})
	defer c.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Guilds().Set(&discord.Guild{ID: discord.Snowflake(i)})
	}
}

func BenchmarkGuildStore_Get(b *testing.B) {
	c := NewMemoryCache(Options{})
	defer c.Close()
	for i := 0; i < 1000; i++ {
		c.Guilds().Set(&discord.Guild{ID: discord.Snowflake(i)})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Guilds().Get(discord.Snowflake(i % 1000))
	}
}

func BenchmarkMessageStore_Add(b *testing.B) {
	c := NewMemoryCache(Options{Messages: MessageOptions{MaxPerChannel: 100}})
	defer c.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Messages().Add(&discord.Message{
			ID:        discord.Snowflake(i),
			ChannelID: 123123123123123123,
		})
	}
}
