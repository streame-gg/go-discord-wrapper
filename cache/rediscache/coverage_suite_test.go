package rediscache_test

import (
	"context"
	"math"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/cache/rediscache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// These tests fill the remaining reachable branches that the happy-path suite
// does not reach: option normalisation, the async write worker, callCtx with a
// write timeout, backend errors (cancelled context after Close), and the
// corrupt-document / stale-index pruning paths in every store.

// asyncCache builds a cache that keeps the async write worker enabled (unlike
// newCache, which forces synchronous writes).
func (s *RedisCacheTestSuite) asyncCache(opts cache.Options) *rediscache.RedisCache {
	s.T().Helper()
	client := redis.NewClient(&redis.Options{Addr: s.redisAddr})
	s.T().Cleanup(func() { _ = client.Close() })
	c := rediscache.NewRedisCache(client, opts).WithKeyPrefix("async:" + s.T().Name())
	s.T().Cleanup(func() { _ = c.Close() })
	return c
}

// rawClient returns a plain redis client and the key prefix used by newCache, so
// tests can inject documents directly into the keyspace.
func (s *RedisCacheTestSuite) rawClient() (*redis.Client, string) {
	client := redis.NewClient(&redis.Options{Addr: s.redisAddr})
	s.T().Cleanup(func() { _ = client.Close() })
	return client, "test:" + s.T().Name()
}

// ── NewRedisCache option normalisation ────────────────────────────────────────

func (s *RedisCacheTestSuite) TestNewRedisCache_OptionNormalisation() {
	// MaxPerChannel < 0 normalises to the default of 100, WriteQueueSize > 0 is
	// honoured as-is.
	c := s.newCache(cache.Options{
		Messages:       cache.MessageOptions{MaxPerChannel: -1},
		WriteQueueSize: 16,
	})
	c.Messages().Add(message("m1", "c1"))
	s.Equal(1, c.Messages().Size(), "caching stays enabled when MaxPerChannel<0 → defaults to 100")
}

// ── callCtx with WriteTimeout ─────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestWriteTimeout_CallCtxBranch() {
	c := s.newCache(cache.Options{WriteTimeout: 5 * time.Second})
	// Sync write goes through callCtx, which derives a deadline context.
	c.Guilds().Set(guild("1"))
	_, ok := c.Guilds().Get(1)
	s.True(ok)
}

// ── Async write worker + enqueueWrite paths ───────────────────────────────────

func (s *RedisCacheTestSuite) TestAsyncWriteWorker() {
	c := s.asyncCache(cache.Options{})

	c.Guilds().Set(guild("1")) // dispatched to the async worker

	s.Require().Eventually(func() bool {
		_, ok := c.Guilds().Get(1)
		return ok
	}, 3*time.Second, 10*time.Millisecond, "async write should land")
}

func (s *RedisCacheTestSuite) TestEnqueueWrite_AfterCloseIsDropped() {
	c := s.asyncCache(cache.Options{})
	s.Require().NoError(c.Close())
	// writerClosed branch: the write is silently dropped, no panic.
	c.Guilds().Set(guild("1"))
}

// ── Backend-error branches (cancelled context after Close) ────────────────────

func (s *RedisCacheTestSuite) TestBackendErrorBranches_AfterClose() {
	c := s.newCache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}})
	c.Close() // cancels the internal context → every command errors

	// A sync write after Close still runs fn directly, exercising the
	// TxPipelined error return in setJSONAndIndex.
	c.Guilds().Set(guild("1"))

	// Reads: every backend command returns context.Canceled → empty results.
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
	s.Equal(0, c.Messages().Channel(gA).Len())

	// Get on a dead backend → miss (Get error branch).
	_, ok := c.Guilds().Get(1)
	s.False(ok)
	_, ok = c.Channels().Get(1)
	s.False(ok)
	_, ok = c.Users().Get(1)
	s.False(ok)
	_, ok = c.Roles().Get(1)
	s.False(ok)
	_, ok = c.Members().Get(gA, 1)
	s.False(ok)
	_, ok = c.VoiceStates().Get(gA, 1)
	s.False(ok)
	_, ok = c.Soundboard().Get(1)
	s.False(ok)
	_, ok = c.ScheduledEvents().Get(1)
	s.False(ok)
	_, ok = c.StageInstances().Get(1)
	s.False(ok)
	_, ok = c.Emojis().Get(1)
	s.False(ok)
	_, ok = c.Stickers().Get(1)
	s.False(ok)
	_, ok = c.Presences().Get(gA, 1)
	s.False(ok)

	// Size on a dead backend returns 0 for every store (Scan/SCard error break).
	s.Equal(0, c.Guilds().Size())
	s.Equal(0, c.Channels().Size())
	s.Equal(0, c.Users().Size())
	s.Equal(0, c.Members().Size())
	s.Equal(0, c.Roles().Size())
	s.Equal(0, c.VoiceStates().Size())
	s.Equal(0, c.Soundboard().Size())
	s.Equal(0, c.ScheduledEvents().Size())
	s.Equal(0, c.StageInstances().Size())
	s.Equal(0, c.Emojis().Size())
	s.Equal(0, c.Stickers().Size())
	s.Equal(0, c.Presences().Size())
	s.Equal(0, c.Messages().Size())
}

// ── Corrupt-document + stale-index pruning branches ───────────────────────────

func (s *RedisCacheTestSuite) TestCorruptAndStaleBranches() {
	c := s.newCache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}})
	raw, p := s.rawClient()
	ctx := context.Background()
	bad := "{ not valid json"

	// seed writes a corrupt value at id "1" plus two stale index entries (no
	// backing value): "2" is consumed by the Get prune path and "3" by the
	// GetByGuild/All stale prune path, so both branches are exercised.
	seed := func(valKey, idxKey string) {
		s.Require().NoError(raw.Set(ctx, valKey, bad, 0).Err())
		s.Require().NoError(raw.SAdd(ctx, idxKey, "1", "2", "3").Err())
	}

	g := gA.String()
	seed(p+":guild:1", p+":guild:index")
	seed(p+":channel:1", p+":channel:index")
	seed(p+":user:1", p+":user:index")
	seed(p+":member:"+g+":1", p+":member:guild:"+g)
	seed(p+":voice_state:"+g+":1", p+":voice_state:guild:"+g)
	seed(p+":presence:"+g+":1", p+":presence:guild:"+g)
	seed(p+":role:1", p+":role:guild:"+g)
	seed(p+":soundboard:1", p+":soundboard:guild:"+g)
	seed(p+":scheduled_event:1", p+":scheduled_event:guild:"+g)
	seed(p+":stage_instance:1", p+":stage_instance:guild:"+g)
	seed(p+":emoji:1", p+":emoji:guild:"+g)
	seed(p+":sticker:1", p+":sticker:guild:"+g)

	// Message index is a sorted set keyed by channel "1".
	s.Require().NoError(raw.Set(ctx, p+":msg:1:1", bad, 0).Err())
	s.Require().NoError(raw.ZAdd(ctx, p+":msg:ch:1",
		redis.Z{Score: 1, Member: "1"}, redis.Z{Score: 2, Member: "2"}, redis.Z{Score: 3, Member: "3"}).Err())

	// For the guild-indexed stores the stale ids need a reverse-lookup map entry
	// so Get's prune path can SRem them from the guild index.
	for _, typ := range []string{"role", "soundboard", "scheduled_event", "stage_instance", "emoji", "sticker"} {
		s.Require().NoError(raw.Set(ctx, p+":"+typ+":map:2", g, 0).Err())
		s.Require().NoError(raw.Set(ctx, p+":"+typ+":map:3", g, 0).Err())
	}

	// Get → corrupt value decodes but json.Unmarshal fails → miss.
	_, ok := c.Guilds().Get(1)
	s.False(ok)
	_, ok = c.Channels().Get(1)
	s.False(ok)
	_, ok = c.Users().Get(1)
	s.False(ok)
	_, ok = c.Members().Get(gA, 1)
	s.False(ok)
	_, ok = c.VoiceStates().Get(gA, 1)
	s.False(ok)
	_, ok = c.Presences().Get(gA, 1)
	s.False(ok)
	_, ok = c.Roles().Get(1)
	s.False(ok)
	_, ok = c.Soundboard().Get(1)
	s.False(ok)
	_, ok = c.ScheduledEvents().Get(1)
	s.False(ok)
	_, ok = c.StageInstances().Get(1)
	s.False(ok)
	_, ok = c.Emojis().Get(1)
	s.False(ok)
	_, ok = c.Stickers().Get(1)
	s.False(ok)

	// Get on the stale id "2" (in the index, no backing value) exercises the
	// prune branch of Get. Done before All/GetByGuild below, which would
	// otherwise delete the reverse-lookup map entries first.
	_, ok = c.Guilds().Get(2)
	s.False(ok)
	_, ok = c.Channels().Get(2)
	s.False(ok)
	_, ok = c.Users().Get(2)
	s.False(ok)
	_, ok = c.Members().Get(gA, 2)
	s.False(ok)
	_, ok = c.VoiceStates().Get(gA, 2)
	s.False(ok)
	_, ok = c.Presences().Get(gA, 2)
	s.False(ok)
	_, ok = c.Roles().Get(2)
	s.False(ok)
	_, ok = c.Soundboard().Get(2)
	s.False(ok)
	_, ok = c.ScheduledEvents().Get(2)
	s.False(ok)
	_, ok = c.StageInstances().Get(2)
	s.False(ok)
	_, ok = c.Emojis().Get(2)
	s.False(ok)
	_, ok = c.Stickers().Get(2)
	s.False(ok)

	// All / GetByGuild / AllInGuild / Channel → corrupt rows skipped, stale rows
	// pruned → empty result.
	s.Equal(0, c.Guilds().All().Len())
	s.Equal(0, c.Channels().All().Len())
	s.Equal(0, c.Users().All().Len())
	s.Equal(0, c.Members().AllInGuild(gA).Len())
	s.Equal(0, c.VoiceStates().AllInGuild(gA).Len())
	s.Equal(0, c.Presences().GetByGuild(gA).Len())
	s.Equal(0, c.Roles().GetByGuild(gA).Len())
	s.Equal(0, c.Roles().All().Len())
	s.Equal(0, c.Soundboard().GetByGuild(gA).Len())
	s.Equal(0, c.ScheduledEvents().GetByGuild(gA).Len())
	s.Equal(0, c.StageInstances().GetByGuild(gA).Len())
	s.Equal(0, c.Emojis().GetByGuild(gA).Len())
	s.Equal(0, c.Stickers().GetByGuild(gA).Len())
	s.Equal(0, c.Messages().Channel(mustSnowflake("1")).Len())
}

// ── json.Marshal error branches ───────────────────────────────────────────────

// TestMarshalErrorBranches forces json.Marshal to fail so the marshal-error
// guards in the shared write helpers (setJSONAndIndex / setJSONAndIndexMap), the
// SetAll loops, and the message Add/Update paths are exercised. encoding/json
// rejects NaN floats and time.Time years outside [0,9999].
func (s *RedisCacheTestSuite) TestMarshalErrorBranches() {
	c := s.newCache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}})
	ff := time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)
	nan := math.NaN()

	// setJSONAndIndex (flat store) marshal-error: Channel.LastPinTimestamp overflow.
	c.Channels().Set(&discord.Channel{ID: 1, LastPinTimestamp: &ff})
	s.Equal(0, c.Channels().Size())

	// setJSONAndIndexMap (guild store) marshal-error: GuildMember.JoinedAt overflow.
	c.Members().Set(gA, &discord.GuildMember{User: user("1"), JoinedAt: ff})
	s.Equal(0, c.Members().Size())

	// Soundboard Set + SetAll marshal-error via NaN volume.
	bad := &discord.SoundboardSound{SoundID: 7, Volume: nan}
	c.Soundboard().Set(gA, bad)
	c.Soundboard().SetAll(gA, []*discord.SoundboardSound{bad})
	s.Equal(0, c.Soundboard().Size())

	// Message Add + Update marshal-error via attachment NaN duration.
	badMsg := &discord.Message{
		ID:          5,
		ChannelID:   9,
		Attachments: []discord.Attachment{{ID: 1, DurationSecs: &nan}},
	}
	c.Messages().Add(badMsg)
	c.Messages().Update(badMsg)
	s.Equal(0, c.Messages().Size())
}

// TestMemberAllInGuild_EdgeIDs covers the member-resolution branches of
// AllInGuild: a non-numeric index id (ParseUint error) whose payload carries a
// valid User (so the user id wins), and a member that resolves to an empty id
// (pruned as stale).
func (s *RedisCacheTestSuite) TestMemberAllInGuild_EdgeIDs() {
	c := s.newCache(cache.Options{})
	raw, p := s.rawClient()
	ctx := context.Background()
	g := gA.String()

	// Non-numeric index id → ParseUint fails; payload has a valid User → its id wins.
	s.Require().NoError(raw.Set(ctx, p+":member:"+g+":abc", `{"user":{"id":"175928847299117063","username":"u"}}`, 0).Err())
	s.Require().NoError(raw.SAdd(ctx, p+":member:guild:"+g, "abc").Err())

	// Index id "0" with no User → resolved id is empty → pruned as stale.
	s.Require().NoError(raw.Set(ctx, p+":member:"+g+":0", `{}`, 0).Err())
	s.Require().NoError(raw.SAdd(ctx, p+":member:guild:"+g, "0").Err())

	coll := c.Members().AllInGuild(gA)
	s.Equal(1, coll.Len())
	_, ok := coll.Get(discord.Snowflake(175928847299117063))
	s.True(ok)
}

// ── Message Update edge paths ─────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestMessageUpdate_NilAndExisting() {
	c := s.newCache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 10}})

	c.Messages().Add(message("m1", "c1"))
	updated := message("m1", "c1")
	updated.Content = "edited"
	c.Messages().Update(updated)

	got, ok := c.Messages().Get(mustSnowflake("c1"), mustSnowflake("m1"))
	s.Require().True(ok)
	s.Equal("edited", got.Content)
}
