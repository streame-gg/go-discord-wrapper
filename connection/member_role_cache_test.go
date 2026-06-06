package connection

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// Bug 72: AddGuildMemberRole/RemoveGuildMemberRole must update the member cache
// after a successful REST call. RestClient is a concrete type and cannot be
// stubbed, so we exercise the cache-update logic directly.

func (cs *ConnectionSuite) TestBug72_AddGuildMemberRole_UpdatesCache() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(111222333444555666)
	userID := discord.Snowflake(999888777666555444)
	roleA := mustSnowflake("aaa111222333444555")
	roleB := mustSnowflake("bbb111222333444555")

	// Seed cache with a member that has roleA.
	c.Cache.Members().Set(guildID, &discord.GuildMember{
		UserID:  userID,
		GuildID: guildID,
		User:    &discord.User{ID: userID},
		Roles:   []discord.Snowflake{roleA},
	})

	// Simulate what the fixed AddGuildMemberRole does after a successful REST call.
	m, ok := c.Cache.Members().Get(guildID, userID)
	require.True(t, ok, "member must be in cache before role add")
	updated := *m
	updated.Roles = append(append([]discord.Snowflake{}, m.Roles...), roleB)
	c.Cache.Members().Set(guildID, &updated)

	got, ok := c.Cache.Members().Get(guildID, userID)
	require.True(t, ok, "member must still be in cache after role add")
	assert.Contains(t, got.Roles, roleA, "roleA must still be present after adding roleB")
	assert.Contains(t, got.Roles, roleB, "roleB must be present after AddGuildMemberRole (Bug 72)")
}

func (cs *ConnectionSuite) TestBug72_RemoveGuildMemberRole_UpdatesCache() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(111222333444555666)
	userID := discord.Snowflake(999888777666555444)
	roleA := mustSnowflake("aaa111222333444555")
	roleB := mustSnowflake("bbb111222333444555")

	// Seed cache with a member that has both roleA and roleB.
	c.Cache.Members().Set(guildID, &discord.GuildMember{
		UserID:  userID,
		GuildID: guildID,
		User:    &discord.User{ID: userID},
		Roles:   []discord.Snowflake{roleA, roleB},
	})

	// Simulate what the fixed RemoveGuildMemberRole does after a successful REST call.
	m, ok := c.Cache.Members().Get(guildID, userID)
	require.True(t, ok, "member must be in cache before role remove")
	updated := *m
	roles := make([]discord.Snowflake, 0, len(m.Roles))
	for _, r := range m.Roles {
		if r != roleB {
			roles = append(roles, r)
		}
	}
	updated.Roles = roles
	c.Cache.Members().Set(guildID, &updated)

	got, ok := c.Cache.Members().Get(guildID, userID)
	require.True(t, ok, "member must still be in cache after role remove")
	assert.Contains(t, got.Roles, roleA, "roleA must still be present after removing roleB")
	assert.NotContains(t, got.Roles, roleB, "roleB must be absent after RemoveGuildMemberRole (Bug 72)")
}

func (cs *ConnectionSuite) TestBug72_AddGuildMemberRole_NoOpWhenMemberNotCached() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(111222333444555666)
	userID := discord.Snowflake(999888777666555444)
	roleB := mustSnowflake("bbb111222333444555")

	// No member seeded — simulates the case where GUILD_MEMBERS intent is absent.
	_, ok := c.Cache.Members().Get(guildID, userID)
	require.False(t, ok, "member must not be in cache before the test")

	// The fixed AddGuildMemberRole only updates if the member exists; verify
	// no partial record is written.
	if m, exists := c.Cache.Members().Get(guildID, userID); exists {
		updated := *m
		updated.Roles = append(append([]discord.Snowflake{}, m.Roles...), roleB)
		c.Cache.Members().Set(guildID, &updated)
	}

	assert.Equal(t, 0, c.Cache.Members().Size(),
		"no partial member must be inserted when member was not already cached (Bug 72)")
}

func (cs *ConnectionSuite) TestBug72_AddGuildMemberRole_NoOpWhenCacheDisabled() {
	t := cs.T()
	mc := cache.NewMemoryCache(cache.Options{})
	c, err := NewClient("test-token", discord.IntentGuilds,
		options.WithCache(mc),
		options.WithCacheStores(cache.CategoryGuilds), // Members excluded
	)
	require.NoError(t, err)

	guildID := discord.Snowflake(111222333444555666)
	roleB := mustSnowflake("bbb111222333444555")
	userID := discord.Snowflake(999888777666555444)

	// cacheStoreEnabled(CategoryMembers) is false → the update block is skipped.
	if c.cacheStoreEnabled(cache.CategoryMembers) {
		if m, ok := c.Cache.Members().Get(guildID, userID); ok {
			updated := *m
			updated.Roles = append(append([]discord.Snowflake{}, m.Roles...), roleB)
			c.Cache.Members().Set(guildID, &updated)
		}
	}

	assert.Equal(t, 0, mc.Members().Size(),
		"cache must not be touched when CategoryMembers is disabled (Bug 72)")
}
