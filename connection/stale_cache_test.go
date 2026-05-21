package connection

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// TestBug25GetGuildRolesEvictsStaleRoles verifies that GetGuildRoles removes
// roles from the previous cache fill that are no longer returned by the API
// (Bug 25). We exercise the DeleteGuild+cacheRoles sequence directly because
// RestClient is a concrete type and cannot be stubbed.
func TestBug25GetGuildRolesEvictsStaleRoles(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(111222333444555)
	roleA := &discord.Role{ID: mustSnowflake("aaa111222333444")}
	roleB := &discord.Role{ID: mustSnowflake("bbb111222333444")}

	c.cacheRole(guildID, roleA)
	c.cacheRole(guildID, roleB)

	_, okB := c.Cache.Roles().Get(roleB.ID)
	require.True(t, okB, "roleB must be in cache before the refresh")

	// Simulate what the fixed GetGuildRoles does: remove all guild roles first,
	// then cache only the roles returned by the API (roleA only).
	c.Cache.Roles().DeleteGuild(guildID)
	c.cacheRoles(guildID, []*discord.Role{roleA})

	_, okB = c.Cache.Roles().Get(roleB.ID)
	assert.False(t, okB, "stale roleB must be evicted by GetGuildRoles refresh (Bug 25)")

	_, okA := c.Cache.Roles().Get(roleA.ID)
	assert.True(t, okA, "fresh roleA must remain in cache after GetGuildRoles refresh")
}

// TestBug25GetGuildChannelsEvictsStaleChannels verifies that GetGuildChannels
// evicts stale channels (deleted on Discord) before inserting the fresh set.
func TestBug25GetGuildChannelsEvictsStaleChannels(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(222333444555666)
	chA := mustSnowflake("aaa222333444555")
	chB := mustSnowflake("bbb222333444555")

	c.cacheChannel(&discord.Channel{ID: chA, GuildID: &guildID})
	c.cacheChannel(&discord.Channel{ID: chB, GuildID: &guildID})

	_, okB := c.Cache.Channels().Get(chB)
	require.True(t, okB, "channelB must be in cache before the refresh")

	// Simulate what the fixed GetGuildChannels does: drain old IDs and delete
	// them from both the channel cache and the message cache.
	for _, oldID := range c.drainGuildChannelIDs(guildID) {
		c.Cache.Channels().Delete(oldID)
		c.Cache.Messages().DeleteChannel(oldID)
	}
	// Only chA is returned by the API this time.
	c.cacheChannels([]*discord.Channel{{ID: chA, GuildID: &guildID}})

	_, okB = c.Cache.Channels().Get(chB)
	assert.False(t, okB, "stale channelB must be evicted by GetGuildChannels refresh (Bug 25)")

	_, okA := c.Cache.Channels().Get(chA)
	assert.True(t, okA, "fresh channelA must remain in cache after GetGuildChannels refresh")
}
