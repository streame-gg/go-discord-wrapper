package connection

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// Bug 5: cacheMember must respect cacheStoreEnabled(CategoryMembers).
// When CategoryMembers is disabled, a GUILD_MEMBER_ADD event must not
// write to the member store.
func TestBug5_CacheMemberRespectsStoreEnabled(t *testing.T) {
	mc := cache.NewMemoryCache(cache.Options{})
	c, err := NewClient("test-token", discord.IntentGuilds,
		options.WithCache(mc),
		options.WithCacheStores(cache.CategoryGuilds), // Members excluded
	)
	require.NoError(t, err)

	guildID := discord.Snowflake(111222333444555666)
	userID := discord.Snowflake(999888777666555444)

	payload := map[string]any{
		"guild_id": strconv.FormatUint(uint64(guildID), 10),
		"user": map[string]any{
			"id":            strconv.FormatUint(uint64(userID), 10),
			"username":      "testuser",
			"discriminator": "0000",
		},
		"roles":     []string{},
		"joined_at": "2024-01-01T00:00:00Z",
		"deaf":      false,
		"mute":      false,
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(json.RawMessage(raw), events.EventGuildMemberAdd, nil)

	assert.Equal(t, 0, c.Cache.Members().Size(),
		"Members store must remain empty when CategoryMembers is disabled (Bug 5)")
}

// Bug 5: cacheMember must still cache the member when CategoryMembers is enabled.
func TestBug5_CacheMemberWorksWhenEnabled(t *testing.T) {
	mc := cache.NewMemoryCache(cache.Options{})
	c, err := NewClient("test-token", discord.IntentGuilds,
		options.WithCache(mc),
		options.WithCacheStores(cache.CategoryMembers|cache.CategoryUsers),
	)
	require.NoError(t, err)

	guildID := discord.Snowflake(111222333444555666)
	userID := discord.Snowflake(999888777666555444)

	payload := map[string]any{
		"guild_id": strconv.FormatUint(uint64(guildID), 10),
		"user": map[string]any{
			"id":            strconv.FormatUint(uint64(userID), 10),
			"username":      "testuser",
			"discriminator": "0000",
		},
		"roles":     []string{},
		"joined_at": "2024-01-01T00:00:00Z",
		"deaf":      false,
		"mute":      false,
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(json.RawMessage(raw), events.EventGuildMemberAdd, nil)

	assert.Equal(t, 1, c.Cache.Members().Size(),
		"Member must be cached when CategoryMembers is enabled (Bug 5)")
}
