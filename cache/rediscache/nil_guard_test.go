package rediscache_test

import (
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/cache/rediscache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// nilRedisCache returns a RedisCache backed by a disconnected client.
// The nil-guard tests call Set/Add with nil arguments which return before
// touching the client, so no connection is needed.
func nilRedisCache() *rediscache.RedisCache {
	client := redis.NewClient(&redis.Options{Addr: "localhost:0"})
	return rediscache.NewRedisCache(client, cache.Options{})
}

// TestRedisNilGuards verifies that all Set/Add methods on the Redis stores
// return without panicking when called with a nil entity argument (Issue 9).
func TestRedisNilGuards(t *testing.T) {
	c := nilRedisCache()
	guildID := discord.Snowflake(1)

	require.NotPanics(t, func() { c.Guilds().Set(nil) }, "Guilds.Set(nil)")
	require.NotPanics(t, func() { c.Channels().Set(nil) }, "Channels.Set(nil)")
	require.NotPanics(t, func() { c.Users().Set(nil) }, "Users.Set(nil)")
	require.NotPanics(t, func() { c.Members().Set(guildID, nil) }, "Members.Set(guildID, nil)")
	require.NotPanics(t, func() { c.Roles().Set(guildID, nil) }, "Roles.Set(guildID, nil)")
	require.NotPanics(t, func() { c.VoiceStates().Set(guildID, nil) }, "VoiceStates.Set(guildID, nil)")
	require.NotPanics(t, func() { c.Soundboard().Set(guildID, nil) }, "Soundboard.Set(guildID, nil)")
	require.NotPanics(t, func() { c.ScheduledEvents().Set(nil) }, "ScheduledEvents.Set(nil)")
	require.NotPanics(t, func() { c.StageInstances().Set(nil) }, "StageInstances.Set(nil)")
	require.NotPanics(t, func() { c.Emojis().Set(guildID, nil) }, "Emojis.Set(guildID, nil)")
	require.NotPanics(t, func() { c.Stickers().Set(guildID, nil) }, "Stickers.Set(guildID, nil)")
	require.NotPanics(t, func() { c.Presences().Set(nil) }, "Presences.Set(nil)")
	require.NotPanics(t, func() {
		c.Presences().Set(&discord.Presence{}) // User.ID == 0 → no-op
	}, "Presences.Set(zeroUserID)")
	require.NotPanics(t, func() { c.Messages().Add(nil) }, "Messages.Add(nil)")
}
