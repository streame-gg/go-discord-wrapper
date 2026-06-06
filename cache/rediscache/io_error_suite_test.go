package rediscache_test

import (
	"context"
	"errors"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/cache/rediscache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// These tests close the last reachable gap in the read paths: the
// "index lookup succeeds but the bulk value fetch (MGET) fails" branch that
// every All / GetByGuild / AllInGuild / Channel method guards with
// `if err != nil { return coll }`. The happy-path suites never reach it and the
// after-Close suite trips the *first* command (SMembers/Scan), not MGET, so a
// per-command fault is required.
//
// Intentionally NOT covered (provably unreachable via the public API): the
// json.Marshal-error `continue` in the emoji and sticker SetAll loops.
// discord.Emoji and discord.Sticker (and the discord.User they embed) contain
// no field encoding/json can reject, so json.Marshal cannot fail for them.
// Covering these would require a test-only marshal seam in production code,
// which we choose not to add — the flat/guild Set helpers' marshal guards are
// already covered via Channel/Member/Soundboard/Message in TestMarshalErrorBranches.

// mgetFaultHook is a go-redis hook that returns an error for MGET commands once
// armed, while letting every other command (including the SMembers/Scan/ZRange
// index lookups and the pipelined writes used to seed data) succeed.
type mgetFaultHook struct{ armed atomic.Bool }

var errInjectedMGet = errors.New("injected mget failure")

func (h *mgetFaultHook) DialHook(next redis.DialHook) redis.DialHook { return next }

func (h *mgetFaultHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		if h.armed.Load() && cmd.Name() == "mget" {
			return errInjectedMGet
		}
		return next(ctx, cmd)
	}
}

func (h *mgetFaultHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		if h.armed.Load() {
			for _, cmd := range cmds {
				if cmd.Name() == "mget" {
					return errInjectedMGet
				}
			}
		}
		return next(ctx, cmds)
	}
}

// faultCache builds a sync-write cache whose client carries an (initially
// disarmed) MGET fault hook, returning the cache and the hook so the test can
// seed data first and then arm the fault before reading.
func (s *RedisCacheTestSuite) faultCache() (*rediscache.RedisCache, *mgetFaultHook) {
	s.T().Helper()
	hook := &mgetFaultHook{}
	client := redis.NewClient(&redis.Options{Addr: s.redisAddr})
	client.AddHook(hook)
	s.T().Cleanup(func() { _ = client.Close() })
	c := rediscache.NewRedisCache(client, cache.Options{
		Messages: cache.MessageOptions{MaxPerChannel: 100},
	}).WithKeyPrefix("mgetfault:" + s.T().Name())
	c.EnableSyncWrites()
	s.T().Cleanup(func() { _ = c.Close() })
	return c, hook
}

// TestMGetErrorBranches seeds one live entity per store, then makes MGET fail so
// the index lookup still returns ids but the value fetch errors out, hitting the
// `if err != nil { return coll }` guard in every read aggregation method.
func (s *RedisCacheTestSuite) TestMGetErrorBranches() {
	c, hook := s.faultCache()

	// Seed one live entity in every store (writes do not use a standalone MGET).
	c.Guilds().Set(guild("1"))
	c.Channels().Set(channel("1"))
	c.Users().Set(user("1"))
	c.Roles().Set(gA, &discord.Role{ID: 1, Name: "r1"})
	c.Members().Set(gA, member("1"))
	c.VoiceStates().Set(gA, &discord.VoiceState{UserID: 1})
	c.Soundboard().Set(gA, &discord.SoundboardSound{SoundID: 1, Name: "s1"})
	c.ScheduledEvents().Set(&discord.GuildScheduledEvent{ID: 1, GuildID: gA, Name: "e1"})
	c.StageInstances().Set(&discord.StageInstance{ID: 1, GuildID: gA, Topic: "t1"})
	c.Emojis().Set(gA, &discord.Emoji{ID: 1, Name: "e1"})
	c.Stickers().Set(gA, sticker("1"))
	c.Presences().Set(&discord.Presence{GuildID: gA, User: discord.PartialPresenceUser{ID: 1}})
	c.Bans().Set(gA, &discord.Ban{User: discord.User{ID: 1}})
	c.AutoModRules().Set(gA, &discord.AutoModerationRule{ID: 1, GuildID: gA, Name: "r1"})
	c.Invites().SetWithGuild(gA, &discord.Invite{Code: "c1"})
	c.Messages().Add(message("1", "9"))

	// Sanity: everything is readable before the fault is armed.
	s.Require().Equal(1, c.Guilds().All().Len())

	hook.armed.Store(true)

	// Index lookups succeed, MGET fails → each method returns an empty collection
	// via its error guard rather than panicking.
	s.Equal(0, c.Guilds().All().Len())
	s.Equal(0, c.Channels().All().Len())
	s.Equal(0, c.Users().All().Len())
	s.Equal(0, c.Roles().All().Len())
	s.Equal(0, c.Roles().GetByGuild(gA).Len())
	s.Equal(0, c.Members().AllInGuild(gA).Len())
	s.Equal(0, c.VoiceStates().AllInGuild(gA).Len())
	s.Equal(0, c.Soundboard().GetByGuild(gA).Len())
	s.Equal(0, c.ScheduledEvents().GetByGuild(gA).Len())
	s.Equal(0, c.StageInstances().GetByGuild(gA).Len())
	s.Equal(0, c.Emojis().GetByGuild(gA).Len())
	s.Equal(0, c.Stickers().GetByGuild(gA).Len())
	s.Equal(0, c.Presences().GetByGuild(gA).Len())
	s.Equal(0, c.Bans().AllInGuild(gA).Len())
	s.Equal(0, c.AutoModRules().GetByGuild(gA).Len())
	s.Equal(0, c.Invites().GetByGuild(gA).Len())
	s.Equal(0, c.Messages().Channel(mustSnowflake("9")).Len())
}

// TestRoleAll_SkipsBadGuildIndexKey covers the role All() per-key guard
// (`if err != nil || len(roleIDs) == 0 { continue }`): a key that matches the
// `role:guild:*` scan pattern but is not a set makes SMembers error, so that key
// is skipped while the real guild's roles still resolve.
func (s *RedisCacheTestSuite) TestRoleAll_SkipsBadGuildIndexKey() {
	c := s.newCache(cache.Options{})
	c.Roles().Set(gA, &discord.Role{ID: 1, Name: "r1"})

	// A non-set value parked under the role:guild: namespace. Scan picks it up but
	// SMembers returns WRONGTYPE → the continue branch runs.
	raw := redis.NewClient(&redis.Options{Addr: s.redisAddr})
	s.T().Cleanup(func() { _ = raw.Close() })
	prefix := "test:" + s.T().Name()
	s.Require().NoError(raw.Set(context.Background(), prefix+":role:guild:zzz", "not-a-set", 0).Err())

	s.Equal(1, c.Roles().All().Len(), "real guild role still resolves; bad key skipped")
}

// TestRoleAll_SkipsMissingValue covers role All()'s `if v == nil { continue }`
// guard: an id present in the guild's role set but with no backing value key, so
// MGET returns a nil slot alongside the live role.
func (s *RedisCacheTestSuite) TestRoleAll_SkipsMissingValue() {
	c := s.newCache(cache.Options{})
	c.Roles().Set(gA, &discord.Role{ID: 1, Name: "r1"})

	// Add a phantom id to the guild role index with no role:<id> value.
	raw := redis.NewClient(&redis.Options{Addr: s.redisAddr})
	s.T().Cleanup(func() { _ = raw.Close() })
	prefix := "test:" + s.T().Name()
	s.Require().NoError(raw.SAdd(context.Background(), prefix+":role:guild:"+gA.String(), "999").Err())

	all := c.Roles().All()
	s.Equal(1, all.Len(), "phantom role id (nil MGET value) is skipped")
}

// slowPipelineHook blocks pipelined writes (the path the async write worker uses)
// once armed, so the worker stalls and the bounded write channel overflows.
type slowPipelineHook struct {
	armed   atomic.Bool
	release chan struct{}
}

func (h *slowPipelineHook) DialHook(next redis.DialHook) redis.DialHook { return next }

func (h *slowPipelineHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook { return next }

func (h *slowPipelineHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		if h.armed.Load() {
			select {
			case <-h.release:
			case <-time.After(2 * time.Second): // safety valve
			}
		}
		return next(ctx, cmds)
	}
}

// TestEnqueueWrite_QueueFullIsDropped covers the `default:` arm of enqueueWrite:
// with the async worker stalled on a blocked write and a one-slot queue, further
// writes find the channel full and are dropped instead of blocking the caller.
func (s *RedisCacheTestSuite) TestEnqueueWrite_QueueFullIsDropped() {
	hook := &slowPipelineHook{release: make(chan struct{})}
	client := redis.NewClient(&redis.Options{Addr: s.redisAddr})
	client.AddHook(hook)
	s.T().Cleanup(func() { _ = client.Close() })

	c := rediscache.NewRedisCache(client, cache.Options{WriteQueueSize: 1}).
		WithKeyPrefix("queuefull:" + s.T().Name())
	s.T().Cleanup(func() {
		close(hook.release) // let any blocked write finish
		_ = c.Close()
	})

	hook.armed.Store(true)

	// First write is taken by the worker and blocks in the pipeline hook; the
	// next fills the 1-slot channel; the remainder hit the default (drop) arm.
	for i := 0; i < 64; i++ {
		c.Guilds().Set(guild(strconv.Itoa(i)))
	}
	// No assertion beyond "did not block or panic": the drop arm has no
	// observable effect other than the write being discarded.
}

// compile-time assertions that the hooks satisfy redis.Hook.
var (
	_ redis.Hook = (*mgetFaultHook)(nil)
	_ redis.Hook = (*slowPipelineHook)(nil)
)
