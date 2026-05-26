package connection

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// newClientWithCache creates a Client backed by a fresh MemoryCache configured
// with all stores enabled, suitable for gateway event handler tests.
func newClientWithCache(t *testing.T) *Client {
	t.Helper()
	mc := cache.NewMemoryCache(cache.Options{
		Messages: cache.MessageOptions{MaxPerChannel: 100},
	})
	c, err := NewClient("test-token", discord.IntentGuilds,
		options.WithCache(mc),
	)
	require.NoError(t, err)
	return c
}

// dispatchGuildDelete fires internalEventHandler with a GUILD_DELETE payload
// for guildID, simulating the gateway receiving the event.
func dispatchGuildDelete(t *testing.T, c *Client, guildID discord.Snowflake) {
	t.Helper()
	payload := map[string]any{"id": guildID.String()}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)
	c.internalEventHandler(json.RawMessage(raw), events.EventGuildDelete, nil)
}

// TestGuildDelete_CleansCacheChannels verifies that channels belonging to the
// deleted guild are removed from the channel cache.
func TestGuildDelete_CleansCacheChannels(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(111000111000111)
	chID := discord.Snowflake(222000222000222)

	ch := &discord.Channel{
		ID:      chID,
		GuildID: &guildID,
	}
	c.cacheChannel(ch)

	_, ok := c.Cache.Channels().Get(chID)
	require.True(t, ok, "channel should be in cache before GUILD_DELETE")

	dispatchGuildDelete(t, c, guildID)

	_, ok = c.Cache.Channels().Get(chID)
	assert.False(t, ok, "channel should be evicted from cache after GUILD_DELETE")
}

// TestGuildDelete_CleansChannelIndex verifies that both sides of the
// bidirectional channelsByGuild / guildByChannel index are removed.
func TestGuildDelete_CleansChannelIndex(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(333000333000333)
	ch1 := discord.Snowflake(444000444000444)
	ch2 := discord.Snowflake(555000555000555)

	c.trackChannel(&discord.Channel{ID: ch1, GuildID: &guildID})
	c.trackChannel(&discord.Channel{ID: ch2, GuildID: &guildID})

	c.channelIndexMu.RLock()
	assert.Len(t, c.channelsByGuild[guildID], 2, "expected 2 channels before delete")
	c.channelIndexMu.RUnlock()

	dispatchGuildDelete(t, c, guildID)

	c.channelIndexMu.RLock()
	defer c.channelIndexMu.RUnlock()

	assert.Empty(t, c.channelsByGuild[guildID], "channelsByGuild should be empty after GUILD_DELETE")
	_, ok1 := c.guildByChannel[ch1]
	_, ok2 := c.guildByChannel[ch2]
	assert.False(t, ok1, "guildByChannel entry for ch1 should be gone after GUILD_DELETE")
	assert.False(t, ok2, "guildByChannel entry for ch2 should be gone after GUILD_DELETE")
}

// TestGuildDelete_CleansMessageCache verifies that messages cached for guild
// channels are removed when the guild is deleted.
func TestGuildDelete_CleansMessageCache(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(666000666000666)
	chID := discord.Snowflake(777000777000777)
	msgID := discord.Snowflake(888000888000888)

	ch := &discord.Channel{ID: chID, GuildID: &guildID}
	c.cacheChannel(ch)
	c.cacheMessage(&discord.Message{ID: msgID, ChannelID: chID})

	_, ok := c.Cache.Messages().Get(chID, msgID)
	require.True(t, ok, "message should be in cache before GUILD_DELETE")

	dispatchGuildDelete(t, c, guildID)

	_, ok = c.Cache.Messages().Get(chID, msgID)
	assert.False(t, ok, "message should be evicted from cache after GUILD_DELETE")
}

// TestGuildDelete_CleansGuildAndMembers verifies that the guild itself plus
// its members are removed from the cache.
func TestGuildDelete_CleansGuildAndMembers(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(999000999000999)
	userID := discord.Snowflake(101010101010101)

	c.Cache.Guilds().Set(&discord.Guild{ID: guildID})
	c.Cache.Members().Set(guildID, &discord.GuildMember{
		User: &discord.User{ID: userID},
	})

	_, guildOK := c.Cache.Guilds().Get(guildID)
	require.True(t, guildOK, "guild should be in cache before GUILD_DELETE")

	dispatchGuildDelete(t, c, guildID)

	_, guildOK = c.Cache.Guilds().Get(guildID)
	assert.False(t, guildOK, "guild should be evicted from cache after GUILD_DELETE")

	members := c.Cache.Members().AllInGuild(guildID)
	assert.Equal(t, 0, members.Len(), "guild members should be evicted from cache after GUILD_DELETE")
}

// TestGuildDelete_UnavailableIsNoop verifies that a GUILD_DELETE with
// unavailable=true does NOT evict the cache (temporary outage, not a removal).
func TestGuildDelete_UnavailableIsNoop(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(121212121212121)

	c.Cache.Guilds().Set(&discord.Guild{ID: guildID})

	unavailable := true
	payload := map[string]any{"id": guildID.String(), "unavailable": unavailable}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)
	c.internalEventHandler(json.RawMessage(raw), events.EventGuildDelete, nil)

	_, ok := c.Cache.Guilds().Get(guildID)
	assert.True(t, ok, "guild should remain in cache for a temporary outage (unavailable=true)")
}

// TestGuildDelete_UnavailableDispatches verifies that the first GUILD_DELETE
// with unavailable=true returns true (dispatch to user handlers) so bots can
// react to a guild going offline.  A second identical event should return false
// (deduplicated).
func TestGuildDelete_UnavailableDispatches(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(131313131313131)
	unavailable := true
	payload := map[string]any{"id": guildID.String(), "unavailable": unavailable}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	var ev events.GuildDeleteEvent
	canDispatch := c.internalEventHandler(json.RawMessage(raw), events.EventGuildDelete, &ev)
	assert.True(t, canDispatch, "first GUILD_DELETE unavailable=true must be dispatched to user handlers")

	// Second event for same guild — already tracked as unavailable, deduplicate.
	canDispatch2 := c.internalEventHandler(json.RawMessage(raw), events.EventGuildDelete, &ev)
	assert.False(t, canDispatch2, "duplicate GUILD_DELETE unavailable=true must not be dispatched again")
}
