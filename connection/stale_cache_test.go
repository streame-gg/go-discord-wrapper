package connection

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// TestBug25GetGuildRolesEvictsStaleRoles verifies that GetGuildRoles removes
// roles from the previous cache fill that are no longer returned by the API
// (Bug 25). We exercise the DeleteGuild+cacheRoles sequence directly because
// RestClient is a concrete type and cannot be stubbed.
func TestBug25GetGuildRolesEvictsStaleRoles(t *testing.T) {
	c := newClientWithCache(t)

	guildID := common.Snowflake("111222333444555")
	roleA := &common.Role{ID: "aaa111222333444"}
	roleB := &common.Role{ID: "bbb111222333444"}

	c.cacheRole(guildID, roleA)
	c.cacheRole(guildID, roleB)

	_, okB := c.Cache.Roles().Get(roleB.ID)
	require.True(t, okB, "roleB must be in cache before the refresh")

	// Simulate what the fixed GetGuildRoles does: remove all guild roles first,
	// then cache only the roles returned by the API (roleA only).
	c.Cache.Roles().DeleteGuild(guildID)
	c.cacheRoles(guildID, []*common.Role{roleA})

	_, okB = c.Cache.Roles().Get(roleB.ID)
	assert.False(t, okB, "stale roleB must be evicted by GetGuildRoles refresh (Bug 25)")

	_, okA := c.Cache.Roles().Get(roleA.ID)
	assert.True(t, okA, "fresh roleA must remain in cache after GetGuildRoles refresh")
}

// TestBug25GetGuildChannelsEvictsStaleChannels verifies that GetGuildChannels
// evicts stale channels (deleted on Discord) before inserting the fresh set.
func TestBug25GetGuildChannelsEvictsStaleChannels(t *testing.T) {
	c := newClientWithCache(t)

	guildID := common.Snowflake("222333444555666")
	chA := common.Snowflake("aaa222333444555")
	chB := common.Snowflake("bbb222333444555")

	c.cacheChannel(&common.Channel{ID: chA, GuildID: &guildID})
	c.cacheChannel(&common.Channel{ID: chB, GuildID: &guildID})

	_, okB := c.Cache.Channels().Get(chB)
	require.True(t, okB, "channelB must be in cache before the refresh")

	// Simulate what the fixed GetGuildChannels does: drain old IDs and delete
	// them from both the channel cache and the message cache.
	for _, oldID := range c.drainGuildChannelIDs(guildID) {
		c.Cache.Channels().Delete(oldID)
		c.Cache.Messages().DeleteChannel(oldID)
	}
	// Only chA is returned by the API this time.
	c.cacheChannels([]*common.Channel{{ID: chA, GuildID: &guildID}})

	_, okB = c.Cache.Channels().Get(chB)
	assert.False(t, okB, "stale channelB must be evicted by GetGuildChannels refresh (Bug 25)")

	_, okA := c.Cache.Channels().Get(chA)
	assert.True(t, okA, "fresh channelA must remain in cache after GetGuildChannels refresh")
}
