package managers_test

import (
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/discord/managers"
)

func sf(v uint64) *discord.Snowflake {
	s := discord.Snowflake(v)
	return &s
}

func (su *guildChannelSuite) TestGuildChannelManager_Get_HitMatchingGuild() {
	t := su.T()
	const guildID discord.Snowflake = 100
	const channelID discord.Snowflake = 200
	ch := &discord.Channel{ID: channelID, GuildID: sf(uint64(guildID))}

	cache := newStubCache()
	cache.channels.channels[channelID] = ch
	m := managers.NewGuildChannelManager(guildID, &stubEntityClient{cache: cache})

	got, ok := m.Get(channelID)
	if !ok || got != ch {
		t.Fatalf("expected hit, got %v, %v", got, ok)
	}
}

func (su *guildChannelSuite) TestGuildChannelManager_Get_GuildMismatch() {
	t := su.T()
	const guildID discord.Snowflake = 100
	const channelID discord.Snowflake = 200
	// Channel belongs to a different guild.
	ch := &discord.Channel{ID: channelID, GuildID: sf(999)}

	cache := newStubCache()
	cache.channels.channels[channelID] = ch
	m := managers.NewGuildChannelManager(guildID, &stubEntityClient{cache: cache})

	if got, ok := m.Get(channelID); ok || got != nil {
		t.Fatalf("expected miss for mismatched guild, got %v, %v", got, ok)
	}
}

func (su *guildChannelSuite) TestGuildChannelManager_Get_NilGuildID() {
	t := su.T()
	const guildID discord.Snowflake = 100
	const channelID discord.Snowflake = 200
	ch := &discord.Channel{ID: channelID, GuildID: nil} // e.g. a DM channel

	cache := newStubCache()
	cache.channels.channels[channelID] = ch
	m := managers.NewGuildChannelManager(guildID, &stubEntityClient{cache: cache})

	if got, ok := m.Get(channelID); ok || got != nil {
		t.Fatalf("expected miss for nil GuildID, got %v, %v", got, ok)
	}
}

func (su *guildChannelSuite) TestGuildChannelManager_Get_NotInCache() {
	t := su.T()
	cache := newStubCache()
	m := managers.NewGuildChannelManager(100, &stubEntityClient{cache: cache})
	if got, ok := m.Get(200); ok || got != nil {
		t.Fatalf("expected miss for absent channel, got %v, %v", got, ok)
	}
}

func (su *guildChannelSuite) TestGuildChannelManager_Get_NoCache() {
	t := su.T()
	m := managers.NewGuildChannelManager(100, &stubEntityClient{cache: nil})
	if got, ok := m.Get(200); ok || got != nil {
		t.Fatalf("expected miss with nil cache, got %v, %v", got, ok)
	}
}

func (su *guildChannelSuite) TestGuildChannelManager_Resolve() {
	t := su.T()
	const guildID discord.Snowflake = 100
	const channelID discord.Snowflake = 200000000000000000 // valid 18-digit snowflake
	ch := &discord.Channel{ID: channelID, GuildID: sf(uint64(guildID))}

	cache := newStubCache()
	cache.channels.channels[channelID] = ch
	m := managers.NewGuildChannelManager(guildID, &stubEntityClient{cache: cache})

	// *Channel returns directly.
	if got, err := m.Resolve(ch); err != nil || got != ch {
		t.Fatalf("Resolve(*Channel) = %v, %v", got, err)
	}
	// Snowflake resolves via cache.
	if got, err := m.Resolve(channelID); err != nil || got != ch {
		t.Fatalf("Resolve(Snowflake) = %v, %v", got, err)
	}
	// Valid numeric string resolves via cache.
	if got, err := m.Resolve("200000000000000000"); err != nil || got != ch {
		t.Fatalf("Resolve(string) = %v, %v", got, err)
	}
	// Invalid string returns a parse error.
	if _, err := m.Resolve("not-a-snowflake"); err == nil {
		t.Fatal("Resolve(invalid string) should error")
	}
	// Unsupported type returns ErrNotConvertable.
	if _, err := m.Resolve(123); err != discord.ErrNotConvertable {
		t.Fatalf("Resolve(int) = %v, want ErrNotConvertable", err)
	}
}

type guildChannelSuite struct{ suite.Suite }

func TestGuildChannelSuite(t *testing.T) { suite.Run(t, new(guildChannelSuite)) }
