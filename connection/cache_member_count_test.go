package connection

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// TestReconnect_CleansStaleGuildMemberCount verifies that a READY event during
// reconnect removes guilds from guildMemberCounts when the bot is no longer a
// member, so GuildCount() does not return stale data (P1-16).
func TestReconnect_CleansStaleGuildMemberCount(t *testing.T) {
	c := newClientWithCache(t)

	guild1 := discord.Snowflake(111111111111111)
	guild2 := discord.Snowflake(222222222222222) // will be absent from the reconnect READY

	// Seed both guilds into the cache and the member-count map.
	c.Cache.Guilds().Set(&discord.Guild{ID: guild1})
	c.Cache.Guilds().Set(&discord.Guild{ID: guild2})
	c.guildMu.Lock()
	c.guildMemberCounts[guild1] = 5
	c.guildMemberCounts[guild2] = 3
	c.guildMu.Unlock()

	require.Equal(t, 2, c.GuildCount(), "setup: expected 2 guilds before reconnect")

	// The READY handler writes to d.Websocket — init a minimal stub.
	c.Websocket = &Websocket{
		Closed: make(chan struct{}),
		Ready:  make(chan struct{}),
	}

	// Simulate a reconnect READY that only lists guild1 (bot was kicked from guild2 while offline).
	readyPayload := map[string]any{
		"user":               map[string]any{"id": "1"},
		"session_id":         "sess-reconnect",
		"resume_gateway_url": "wss://fake.discord.gg",
		"guilds":             []map[string]any{{"id": strconv.FormatUint(uint64(guild1), 10)}},
	}
	raw, err := json.Marshal(readyPayload)
	require.NoError(t, err)
	c.internalEventHandler(json.RawMessage(raw), events.EventReady, nil)

	assert.Equal(t, 1, c.GuildCount(),
		"GuildCount must be 1 after reconnect; guild2 was dropped from READY (P1-16)")
}

// TestReconnect_CleansStaleGuildMemberCount_NoCache verifies stale reconnect
// state is pruned from guildMemberCounts even when cache is disabled.
func TestReconnect_CleansStaleGuildMemberCount_NoCache(t *testing.T) {
	c, err := NewClient("test-token", discord.IntentGuilds)
	require.NoError(t, err)

	guild1 := discord.Snowflake(111111111111111)
	guild2 := discord.Snowflake(222222222222222) // will be absent from the reconnect READY

	c.guildMu.Lock()
	c.guildMemberCounts[guild1] = 5
	c.guildMemberCounts[guild2] = 3
	c.guildMu.Unlock()

	require.Equal(t, 2, c.GuildCount(), "setup: expected 2 guilds before reconnect")

	// The READY handler writes to d.Websocket — init a minimal stub.
	c.Websocket = &Websocket{
		Closed: make(chan struct{}),
		Ready:  make(chan struct{}),
	}

	// Simulate a reconnect READY that only lists guild1 (bot was kicked from guild2 while offline).
	readyPayload := map[string]any{
		"user":               map[string]any{"id": "1"},
		"session_id":         "sess-reconnect",
		"resume_gateway_url": "wss://fake.discord.gg",
		"guilds":             []map[string]any{{"id": strconv.FormatUint(uint64(guild1), 10)}},
	}
	raw, err := json.Marshal(readyPayload)
	require.NoError(t, err)
	c.internalEventHandler(json.RawMessage(raw), events.EventReady, nil)

	assert.Equal(t, 1, c.GuildCount(),
		"GuildCount must be 1 after reconnect with cache disabled; guild2 was dropped from READY")
}

// TestGuildDelete_CleansGuildMemberCount verifies that a GUILD_DELETE kick
// removes the entry from guildMemberCounts so GuildCount() is accurate.
func TestGuildDelete_CleansGuildMemberCount(t *testing.T) {
	c := newClientWithCache(t)

	guild1 := discord.Snowflake(333333333333333)
	guild2 := discord.Snowflake(444444444444444)

	c.guildMu.Lock()
	c.guildMemberCounts[guild1] = 10
	c.guildMemberCounts[guild2] = 7
	c.guildMu.Unlock()

	require.Equal(t, 2, c.GuildCount())

	// Kick from guild2 via GUILD_DELETE (unavailable omitted → kick path).
	dispatchGuildDelete(t, c, guild2)

	assert.Equal(t, 1, c.GuildCount(),
		"GuildCount must drop to 1 after GUILD_DELETE kick (P1-16)")
}
