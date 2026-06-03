package rediscache_test

import (
	"github.com/redis/go-redis/v9"

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

// TestBug161_StoreAccessorsReturnSingletons verifies that each Xxx() accessor
// returns the same pointer on every call.
//
// Root cause: every call to Guilds(), Channels(), etc. allocated a new
// &redisXxxStore{c} on the heap, causing measurable GC pressure on high-traffic
// bots. The fix stores pre-allocated singletons and returns them directly.
func (s *RedisCacheTestSuite) TestBug161_StoreAccessorsReturnSingletons() {
	c := nilRedisCache()

	s.True(c.Guilds() == c.Guilds(), "Bug161: Guilds() must return the same pointer")
	s.True(c.Channels() == c.Channels(), "Bug161: Channels() must return the same pointer")
	s.True(c.Users() == c.Users(), "Bug161: Users() must return the same pointer")
	s.True(c.Members() == c.Members(), "Bug161: Members() must return the same pointer")
	s.True(c.Roles() == c.Roles(), "Bug161: Roles() must return the same pointer")
	s.True(c.Messages() == c.Messages(), "Bug161: Messages() must return the same pointer")
	s.True(c.VoiceStates() == c.VoiceStates(), "Bug161: VoiceStates() must return the same pointer")
	s.True(c.Soundboard() == c.Soundboard(), "Bug161: Soundboard() must return the same pointer")
	s.True(c.ScheduledEvents() == c.ScheduledEvents(), "Bug161: ScheduledEvents() must return the same pointer")
	s.True(c.StageInstances() == c.StageInstances(), "Bug161: StageInstances() must return the same pointer")
	s.True(c.Emojis() == c.Emojis(), "Bug161: Emojis() must return the same pointer")
	s.True(c.Stickers() == c.Stickers(), "Bug161: Stickers() must return the same pointer")
	s.True(c.Presences() == c.Presences(), "Bug161: Presences() must return the same pointer")
}

// TestRedisNilGuards verifies that all Set/Add methods on the Redis stores
// return without panicking when called with a nil entity argument (Issue 9).
func (s *RedisCacheTestSuite) TestRedisNilGuards() {
	c := nilRedisCache()
	guildID := discord.Snowflake(1)
	r := s.Require()

	r.NotPanics(func() { c.Guilds().Set(nil) }, "Guilds.Set(nil)")
	r.NotPanics(func() { c.Channels().Set(nil) }, "Channels.Set(nil)")
	r.NotPanics(func() { c.Users().Set(nil) }, "Users.Set(nil)")
	r.NotPanics(func() { c.Members().Set(guildID, nil) }, "Members.Set(guildID, nil)")
	r.NotPanics(func() { c.Roles().Set(guildID, nil) }, "Roles.Set(guildID, nil)")
	r.NotPanics(func() { c.VoiceStates().Set(guildID, nil) }, "VoiceStates.Set(guildID, nil)")
	r.NotPanics(func() { c.Soundboard().Set(guildID, nil) }, "Soundboard.Set(guildID, nil)")
	r.NotPanics(func() { c.ScheduledEvents().Set(nil) }, "ScheduledEvents.Set(nil)")
	r.NotPanics(func() { c.StageInstances().Set(nil) }, "StageInstances.Set(nil)")
	r.NotPanics(func() { c.Emojis().Set(guildID, nil) }, "Emojis.Set(guildID, nil)")
	r.NotPanics(func() { c.Stickers().Set(guildID, nil) }, "Stickers.Set(guildID, nil)")
	r.NotPanics(func() { c.Presences().Set(nil) }, "Presences.Set(nil)")
	r.NotPanics(func() {
		c.Presences().Set(&discord.Presence{}) // User.ID == 0 → no-op
	}, "Presences.Set(zeroUserID)")
	r.NotPanics(func() { c.Messages().Add(nil) }, "Messages.Add(nil)")
}
