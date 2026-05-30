package connection

// Tests for Bug 79: GUILD_BAN_ADD/REMOVE gateway events must update the ban cache
// and remove the banned member from the member cache.

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

func TestBug79_GuildBanAdd_CachesBan(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100000000000000001)
	userID := discord.Snowflake(200000000000000002)

	payload := map[string]any{
		"guild_id": guildID.String(),
		"user": map[string]any{
			"id":            userID.String(),
			"username":      "banned_user",
			"discriminator": "0000",
		},
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(json.RawMessage(raw), events.EventGuildBanAdd, nil)

	ban, ok := c.Cache.Bans().Get(guildID, userID)
	require.True(t, ok, "GUILD_BAN_ADD must add ban to BanStore (Bug 79)")
	assert.Equal(t, userID, ban.User.ID)
}

func TestBug79_GuildBanAdd_RemovesMemberFromCache(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100000000000000001)
	userID := discord.Snowflake(200000000000000002)

	// Pre-seed member so we can verify it gets removed.
	c.Cache.Members().Set(guildID, &discord.GuildMember{
		User:    &discord.User{ID: userID},
		UserID:  userID,
		GuildID: guildID,
	})
	_, ok := c.Cache.Members().Get(guildID, userID)
	require.True(t, ok, "member must be in cache before GUILD_BAN_ADD")

	payload := map[string]any{
		"guild_id": guildID.String(),
		"user": map[string]any{
			"id":            userID.String(),
			"username":      "banned_user",
			"discriminator": "0000",
		},
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(json.RawMessage(raw), events.EventGuildBanAdd, nil)

	_, ok = c.Cache.Members().Get(guildID, userID)
	assert.False(t, ok, "banned member must be removed from member cache on GUILD_BAN_ADD (Bug 79)")
}

func TestBug79_GuildBanRemove_DeletesBan(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100000000000000001)
	userID := discord.Snowflake(200000000000000002)

	// Pre-seed a ban so we can verify it gets removed.
	c.Cache.Bans().Set(guildID, &discord.Ban{
		User: discord.User{ID: userID},
	})
	_, ok := c.Cache.Bans().Get(guildID, userID)
	require.True(t, ok, "ban must be in cache before GUILD_BAN_REMOVE")

	payload := map[string]any{
		"guild_id": guildID.String(),
		"user": map[string]any{
			"id":            userID.String(),
			"username":      "unbanned_user",
			"discriminator": "0000",
		},
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(json.RawMessage(raw), events.EventGuildBanRemove, nil)

	_, ok = c.Cache.Bans().Get(guildID, userID)
	assert.False(t, ok, "ban must be removed from BanStore on GUILD_BAN_REMOVE (Bug 79)")
}

func TestBug79_GuildDelete_ClearsBans(t *testing.T) {
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100000000000000001)
	userID := discord.Snowflake(200000000000000002)

	c.Cache.Bans().Set(guildID, &discord.Ban{User: discord.User{ID: userID}})
	require.Equal(t, 1, c.Cache.Bans().AllInGuild(guildID).Len())

	dispatchGuildDelete(t, c, guildID)

	assert.Equal(t, 0, c.Cache.Bans().AllInGuild(guildID).Len(),
		"GUILD_DELETE must evict all bans for the guild")
}
