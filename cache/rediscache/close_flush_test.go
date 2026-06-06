package rediscache_test

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/cache/rediscache"
)

// TestClose_FlushesQueuedWrites is a regression test: Close must drain the async
// write queue while the context is still live. Previously Close cancelled the
// context first, so the drained write executed with an already-cancelled
// callCtx and was silently lost.
func (s *RedisCacheTestSuite) TestClose_FlushesQueuedWrites() {
	client := redis.NewClient(&redis.Options{Addr: s.redisAddr})
	s.T().Cleanup(func() { _ = client.Close() })

	prefix := "closeflush:" + s.T().Name()
	// Async write worker enabled (no EnableSyncWrites): the Set is queued.
	c := rediscache.NewRedisCache(client, cache.Options{}).WithKeyPrefix(prefix)

	c.Guilds().Set(guild("1"))

	s.Require().NoError(c.Close(), "Close drains then cancels")

	// Verify via an independent client that the queued write was flushed.
	raw := redis.NewClient(&redis.Options{Addr: s.redisAddr})
	s.T().Cleanup(func() { _ = raw.Close() })
	val, err := raw.Get(context.Background(), prefix+":guild:1").Result()
	s.Require().NoError(err, "queued write must be flushed during Close, not dropped")
	s.Contains(val, "guild-1")
}
