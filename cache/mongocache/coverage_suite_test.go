package mongocache_test

import (
	"context"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// These tests fill the remaining reachable branches: TTL query filters, JSON
// unmarshal failures (via corrupt documents), backend-error paths (cancelled
// context), the async write worker, callCtx with a write timeout, and the
// less-common message-store paths.

// ── Size for the stores not asserted elsewhere ────────────────────────────────

func (s *mongoStoresSuite) TestSize_ChannelUserSticker() {
	c := s.cache(cache.Options{})
	c.Channels().Set(channel("1"))
	c.Users().Set(user("2"))
	c.Stickers().Set(gA, sticker("s"))

	s.Equal(1, c.Channels().Size())
	s.Equal(1, c.Users().Size())
	s.Equal(1, c.Stickers().Size())
}

// ── NewMongoDBCache: MaxPerChannel < 0 normalises to default ──────────────────

func (s *mongoStoresSuite) TestNewMongoDBCache_NegativeMaxPerChannelDefaults() {
	c := s.cache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: -1}})
	c.Messages().Add(message("m1", "c1"))
	s.Equal(1, c.Messages().Size(), "caching stays enabled when MaxPerChannel<0 → defaults to 100")
}

// ── TTL filter branches across every store ────────────────────────────────────

func (s *mongoStoresSuite) TestTTLFilterBranches_AllStores() {
	c := s.cache(cache.Options{TTL: time.Hour, Messages: cache.MessageOptions{MaxPerChannel: 100, TTL: time.Hour}})

	c.Guilds().Set(guild("1"))
	c.Channels().Set(channel("1"))
	c.Users().Set(user("1"))
	c.Members().Set(gA, member("1"))
	c.Roles().Set(gA, &discord.Role{ID: 1})
	c.VoiceStates().Set(gA, &discord.VoiceState{UserID: 1})
	c.Soundboard().Set(gA, &discord.SoundboardSound{SoundID: 1})
	c.ScheduledEvents().Set(&discord.GuildScheduledEvent{ID: 1, GuildID: gA})
	c.StageInstances().Set(&discord.StageInstance{ID: 1, GuildID: gA})
	c.Emojis().Set(gA, &discord.Emoji{ID: 1})
	c.Stickers().Set(gA, sticker("s"))
	c.Presences().Set(&discord.Presence{GuildID: gA, User: discord.PartialPresenceUser{ID: 1}})
	c.Bans().Set(gA, &discord.Ban{User: discord.User{ID: 1}})
	c.AutoModRules().Set(gA, &discord.AutoModerationRule{ID: 1, GuildID: gA})
	c.Invites().SetWithGuild(gA, &discord.Invite{Code: "c1"})
	c.Messages().Add(message("m1", "c1"))

	// Get with TTL filter (expires_at > now).
	_, ok := c.Guilds().Get(1)
	s.True(ok)
	_, ok = c.Channels().Get(1)
	s.True(ok)
	_, ok = c.Users().Get(1)
	s.True(ok)
	_, ok = c.Members().Get(gA, 1)
	s.True(ok)
	_, ok = c.Roles().Get(1)
	s.True(ok)
	_, ok = c.VoiceStates().Get(gA, 1)
	s.True(ok)
	_, ok = c.Soundboard().Get(1)
	s.True(ok)
	_, ok = c.ScheduledEvents().Get(1)
	s.True(ok)
	_, ok = c.StageInstances().Get(1)
	s.True(ok)
	_, ok = c.Emojis().Get(1)
	s.True(ok)
	_, ok = c.Stickers().Get(mustSnowflake("s"))
	s.True(ok)
	_, ok = c.Presences().Get(gA, 1)
	s.True(ok)
	_, ok = c.Bans().Get(gA, 1)
	s.True(ok)
	_, ok = c.AutoModRules().Get(1)
	s.True(ok)
	_, ok = c.Invites().Get("c1")
	s.True(ok)
	_, ok = c.Messages().Get(mustSnowflake("c1"), mustSnowflake("m1"))
	s.True(ok)

	// All / GetByGuild / AllInGuild / Channel with TTL filter.
	s.Equal(1, c.Guilds().All().Len())
	s.Equal(1, c.Channels().All().Len())
	s.Equal(1, c.Users().All().Len())
	s.Equal(1, c.Roles().All().Len())
	s.Equal(1, c.Members().AllInGuild(gA).Len())
	s.Equal(1, c.Roles().GetByGuild(gA).Len())
	s.Equal(1, c.VoiceStates().AllInGuild(gA).Len())
	s.Equal(1, c.Soundboard().GetByGuild(gA).Len())
	s.Equal(1, c.ScheduledEvents().GetByGuild(gA).Len())
	s.Equal(1, c.StageInstances().GetByGuild(gA).Len())
	s.Equal(1, c.Emojis().GetByGuild(gA).Len())
	s.Equal(1, c.Stickers().GetByGuild(gA).Len())
	s.Equal(1, c.Presences().GetByGuild(gA).Len())
	s.Equal(1, c.Bans().AllInGuild(gA).Len())
	s.Equal(1, c.AutoModRules().GetByGuild(gA).Len())
	s.Equal(1, c.Invites().GetByGuild(gA).Len())
	s.Equal(1, c.Messages().Channel(mustSnowflake("c1")).Len())
}

// ── Corrupt-document → JSON unmarshal failure branches ────────────────────────

// insertBad writes a document whose json field is invalid into col so that
// Get/All/GetByGuild exercise their json.Unmarshal error branch.
func (s *mongoStoresSuite) insertBad(col, id, guildID string) {
	doc := bson.M{"_id": id, "json": "{ not valid json"}
	if guildID != "" {
		doc["guild_id"] = guildID
	}
	_, err := s.db.Collection(col).InsertOne(context.Background(), doc)
	s.Require().NoError(err)
}

func (s *mongoStoresSuite) TestUnmarshalErrorBranches() {
	c := s.cache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}})

	gid := "1001"
	s.insertBad("guilds", "1", "")
	s.insertBad("channels", "1", "")
	s.insertBad("users", "1", "")
	s.insertBad("members", "1001:1", gid)
	s.insertBad("roles", "1", gid)
	s.insertBad("voice_states", "1001:1", gid)
	s.insertBad("soundboard_sounds", "1", gid)
	s.insertBad("scheduled_events", "1", gid)
	s.insertBad("stage_instances", "1", gid)
	s.insertBad("emojis", "1", gid)
	s.insertBad("stickers", "1", gid)
	s.insertBad("presences", "1001:1", gid)
	s.insertBad("bans", "1001:1", gid)
	s.insertBad("automod_rules", "1", gid)
	s.insertBad("invites", "c1", gid)
	s.insertBad("messages", "1001:1", "")
	// messages need a channel_id field for Channel(); add one.
	_, err := s.db.Collection("messages").UpdateOne(context.Background(),
		bson.M{"_id": "1001:1"}, bson.M{"$set": bson.M{"channel_id": "1001"}})
	s.Require().NoError(err)

	// Get → decode succeeds, json.Unmarshal fails → miss.
	_, ok := c.Guilds().Get(1)
	s.False(ok)
	_, ok = c.Channels().Get(1)
	s.False(ok)
	_, ok = c.Users().Get(1)
	s.False(ok)
	_, ok = c.Members().Get(gA, 1)
	s.False(ok)
	_, ok = c.Roles().Get(1)
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
	_, ok = c.Bans().Get(gA, 1)
	s.False(ok)
	_, ok = c.AutoModRules().Get(1)
	s.False(ok)
	_, ok = c.Invites().Get("c1")
	s.False(ok)
	_, ok = c.Messages().Get(gA, 1)
	s.False(ok)

	// All / GetByGuild / Channel → corrupt rows are skipped (empty result).
	s.Equal(0, c.Guilds().All().Len())
	s.Equal(0, c.Channels().All().Len())
	s.Equal(0, c.Users().All().Len())
	s.Equal(0, c.Roles().All().Len())
	s.Equal(0, c.Members().AllInGuild(gA).Len())
	s.Equal(0, c.Roles().GetByGuild(gA).Len())
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
	s.Equal(0, c.Messages().Channel(gA).Len())
}

// insertBadType writes a document whose json field has the wrong BSON type (an
// int instead of a string) so the *outer* document decode performed by
// cursor.All / FindOne.Decode fails — distinct from insertBad, which stores a
// valid string that only trips the inner json.Unmarshal.
func (s *mongoStoresSuite) insertBadType(col, id, guildID string) {
	doc := bson.M{"_id": id, "json": 12345}
	if guildID != "" {
		doc["guild_id"] = guildID
	}
	_, err := s.db.Collection(col).InsertOne(context.Background(), doc)
	s.Require().NoError(err)
}

// TestCursorDecodeErrorBranches drives the cursor.All / cursor decode error
// path: when the stored document cannot be decoded into the typed doc struct,
// All / GetByGuild / AllInGuild / Channel return an empty result.
func (s *mongoStoresSuite) TestCursorDecodeErrorBranches() {
	c := s.cache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}})

	gid := "1001"
	s.insertBadType("guilds", "1", "")
	s.insertBadType("channels", "1", "")
	s.insertBadType("users", "1", "")
	s.insertBadType("members", "1001:1", gid)
	s.insertBadType("roles", "1", gid)
	s.insertBadType("voice_states", "1001:1", gid)
	s.insertBadType("soundboard_sounds", "1", gid)
	s.insertBadType("scheduled_events", "1", gid)
	s.insertBadType("stage_instances", "1", gid)
	s.insertBadType("emojis", "1", gid)
	s.insertBadType("stickers", "1", gid)
	s.insertBadType("presences", "1001:1", gid)
	s.insertBadType("bans", "1001:1", gid)
	s.insertBadType("automod_rules", "1", gid)
	s.insertBadType("invites", "c1", gid)
	_, err := s.db.Collection("messages").InsertOne(context.Background(),
		bson.M{"_id": "1001:1", "channel_id": "1001", "json": 12345})
	s.Require().NoError(err)

	s.Equal(0, c.Guilds().All().Len())
	s.Equal(0, c.Channels().All().Len())
	s.Equal(0, c.Users().All().Len())
	s.Equal(0, c.Roles().All().Len())
	s.Equal(0, c.Members().AllInGuild(gA).Len())
	s.Equal(0, c.Roles().GetByGuild(gA).Len())
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
	s.Equal(0, c.Messages().Channel(mustSnowflake("1001")).Len())
}

// ── Message eviction & TTL update ─────────────────────────────────────────────

// TestMessageStore_Eviction exercises the MaxPerChannel ring enforcement: once
// the per-channel count exceeds the cap, the oldest messages are deleted.
func (s *mongoStoresSuite) TestMessageStore_Eviction() {
	c := s.cache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 2}})
	c.Messages().Add(message("m1", "c1"))
	c.Messages().Add(message("m2", "c1"))
	c.Messages().Add(message("m3", "c1")) // exceeds cap → one oldest evicted
	s.Equal(2, c.Messages().Channel(mustSnowflake("c1")).Len())
	s.Equal(2, c.Messages().Size())
}

// TestMessageStore_UpdateWithTTL covers Update's expires_at branch (TTL > 0).
func (s *mongoStoresSuite) TestMessageStore_UpdateWithTTL() {
	c := s.cache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 10, TTL: time.Hour}})
	c.Messages().Add(message("m1", "c1"))
	c.Messages().Update(&discord.Message{
		ID:        mustSnowflake("m1"),
		ChannelID: mustSnowflake("c1"),
		Content:   "edited",
	})
	got, ok := c.Messages().Get(mustSnowflake("c1"), mustSnowflake("m1"))
	s.Require().True(ok)
	s.Equal("edited", got.Content)
}

// ── json.Marshal error branches ───────────────────────────────────────────────

// TestMarshalErrorBranches forces json.Marshal to fail by embedding a NaN float
// (which encoding/json rejects) so the marshal-error guards in the soundboard
// and message stores are exercised. These are the only cached entities whose
// payloads contain a float field reachable from the public API.
func (s *mongoStoresSuite) TestMarshalErrorBranches() {
	c := s.cache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}})
	nan := math.NaN()

	badSound := &discord.SoundboardSound{SoundID: 7, Volume: nan}
	c.Soundboard().Set(gA, badSound)                                // Set marshal-error
	c.Soundboard().SetAll(gA, []*discord.SoundboardSound{badSound}) // SetAll marshal-error (skip item)
	s.Equal(0, c.Soundboard().Size())

	badMsg := &discord.Message{
		ID:          5,
		ChannelID:   9,
		Attachments: []discord.Attachment{{ID: 1, DurationSecs: nan}},
	}
	c.Messages().Add(badMsg)    // Add marshal-error
	c.Messages().Update(badMsg) // Update marshal-error
	s.Equal(0, c.Messages().Size())
}

// TestMarshalErrorBranches_TimeOverflow forces json.Marshal to fail via a
// time.Time whose year is outside encoding/json's supported range [0,9999].
// This reaches the marshal-error guards of the stores whose entities carry a
// time field (members, scheduled events, voice states, channels, presences).
func (s *mongoStoresSuite) TestMarshalErrorBranches_TimeOverflow() {
	c := s.cache(cache.Options{})
	ff := time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)

	c.Members().Set(gA, &discord.GuildMember{User: user("1"), JoinedAt: &ff})
	s.Equal(0, c.Members().Size())

	c.ScheduledEvents().Set(&discord.GuildScheduledEvent{ID: 1, ScheduledStartTime: ff})
	s.Equal(0, c.ScheduledEvents().Size())

	c.VoiceStates().Set(gA, &discord.VoiceState{UserID: 1, RequestToSpeakTimestamp: &ff})
	s.Equal(0, c.VoiceStates().Size())

	c.Channels().Set(&discord.Channel{ID: 1, LastPinTimestamp: &ff})
	s.Equal(0, c.Channels().Size())

	// Invite carries an *time.Time ExpiresAt, so its store() marshal guard is
	// reachable the same way. Set (no guild) and SetWithGuild both route through store().
	c.Invites().Set(&discord.Invite{Code: "x", ExpiresAt: &ff})
	c.Invites().SetWithGuild(gA, &discord.Invite{Code: "y", ExpiresAt: &ff})
	s.Equal(0, c.Invites().Size())
}

// ── Backend-error branches (cancelled context after Close) ────────────────────

func (s *mongoStoresSuite) TestBackendErrorBranches_AfterClose() {
	c := s.cache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}})
	c.Close() // cancels the pkg context → all queries error

	// Find returns an error → each collection method returns an empty result.
	s.Equal(0, c.Guilds().All().Len())
	s.Equal(0, c.Channels().All().Len())
	s.Equal(0, c.Users().All().Len())
	s.Equal(0, c.Roles().All().Len())
	s.Equal(0, c.Members().AllInGuild(gA).Len())
	s.Equal(0, c.Roles().GetByGuild(gA).Len())
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
	s.Equal(0, c.Messages().Channel(gA).Len())

	// Get on a dead backend → FindOne error → miss.
	_, ok := c.Bans().Get(gA, 1)
	s.False(ok)
	_, ok = c.AutoModRules().Get(1)
	s.False(ok)
	_, ok = c.Invites().Get("c1")
	s.False(ok)

	// Size on a dead backend returns 0 (CountDocuments error ignored).
	s.Equal(0, c.Guilds().Size())
	s.Equal(0, c.Bans().Size())
	s.Equal(0, c.AutoModRules().Size())
	s.Equal(0, c.Invites().Size())
}

// ── Async write worker + enqueueWrite paths ───────────────────────────────────

func (s *mongoStoresSuite) TestAsyncWriteWorker() {
	c := s.asyncCache(cache.Options{})

	c.Guilds().Set(guild("1")) // dispatched to the async worker

	s.Require().Eventually(func() bool {
		_, ok := c.Guilds().Get(1)
		return ok
	}, 3*time.Second, 10*time.Millisecond, "async write should land")
}

func (s *mongoStoresSuite) TestEnqueueWrite_AfterCloseIsDropped() {
	c := s.asyncCache(cache.Options{})
	c.Close()
	// writerClosed branch: the write is silently dropped, no panic.
	c.Guilds().Set(guild("1"))
}

// ── callCtx with WriteTimeout ─────────────────────────────────────────────────

func (s *mongoStoresSuite) TestWriteTimeout_CallCtxBranch() {
	c := s.cache(cache.Options{WriteTimeout: 5 * time.Second})
	// Sync write goes through callCtx, which derives a deadline context.
	c.Guilds().Set(guild("1"))
	_, ok := c.Guilds().Get(1)
	s.True(ok)
}

// ── Message store edge paths ──────────────────────────────────────────────────

func (s *mongoStoresSuite) TestMessageStore_DisabledAndNoops() {
	c := s.cache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 0}})

	c.Messages().Add(message("m1", "c1")) // disabled → no-op
	c.Messages().Add(nil)                 // nil → no-op
	s.Equal(0, c.Messages().Channel(mustSnowflake("c1")).Len())
	s.Equal(0, c.Messages().Size())

	// DeleteBulk with empty slice → early return.
	c2 := s.cache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 10}})
	c2.Messages().Add(message("m1", "c1"))
	c2.Messages().DeleteBulk(mustSnowflake("c1"), nil)
	s.Equal(1, c2.Messages().Size())

	// Update on a non-existent message → no-op (does not insert).
	c2.Messages().Update(message("ghost", "c1"))
	_, ok := c2.Messages().Get(mustSnowflake("c1"), mustSnowflake("ghost"))
	s.False(ok)
}
