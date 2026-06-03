package rediscache_test

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// These tests exercise the guild-scoped and composite-key stores that the
// original container-based suite left untested. They run against the same
// miniredis instance as the rest of RedisCacheTestSuite.

const (
	gA = discord.Snowflake(1001)
	gB = discord.Snowflake(2002)
)

// ── accessors ─────────────────────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestAccessors_NotNil() {
	c := s.newCache(cache.Options{})
	s.NotNil(c.Guilds())
	s.NotNil(c.Channels())
	s.NotNil(c.Users())
	s.NotNil(c.Members())
	s.NotNil(c.Roles())
	s.NotNil(c.Messages())
	s.NotNil(c.VoiceStates())
	s.NotNil(c.Soundboard())
	s.NotNil(c.ScheduledEvents())
	s.NotNil(c.StageInstances())
	s.NotNil(c.Emojis())
	s.NotNil(c.Stickers())
	s.NotNil(c.Presences())
}

// ── Role store ────────────────────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestRoleStore() {
	c := s.newCache(cache.Options{})
	r := c.Roles()

	r.Set(gA, nil) // no-op
	r.Set(gA, &discord.Role{ID: 1, Name: "r1"})
	r.Set(gA, &discord.Role{ID: 2, Name: "r2"})
	r.Set(gB, &discord.Role{ID: 3, Name: "r3"})

	got, ok := r.Get(1)
	s.Require().True(ok)
	s.Equal("r1", got.Name)

	_, ok = r.Get(999)
	s.False(ok)

	s.Equal(2, r.GetByGuild(gA).Len())
	s.Equal(3, r.All().Len())
	s.Equal(3, r.Size())

	r.Delete(1)
	s.Equal(1, r.GetByGuild(gA).Len())

	r.DeleteGuild(gA)
	s.Equal(0, r.GetByGuild(gA).Len())
	s.Equal(1, r.Size())
}

// ── VoiceState store ──────────────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestVoiceStateStore() {
	c := s.newCache(cache.Options{})
	v := c.VoiceStates()

	v.Set(gA, nil) // no-op
	v.Set(gA, &discord.VoiceState{UserID: 1})
	v.Set(gA, &discord.VoiceState{UserID: 2})
	v.Set(gB, &discord.VoiceState{UserID: 3})

	got, ok := v.Get(gA, 1)
	s.Require().True(ok)
	s.Equal(discord.Snowflake(1), got.UserID)

	_, ok = v.Get(gA, 999)
	s.False(ok)

	s.Equal(2, v.AllInGuild(gA).Len())
	s.Equal(3, v.Size())

	v.Delete(gA, 1)
	s.Equal(1, v.AllInGuild(gA).Len())

	v.DeleteGuild(gA)
	s.Equal(0, v.AllInGuild(gA).Len())
	s.Equal(1, v.Size())
}

// ── Soundboard store ──────────────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestSoundboardStore() {
	c := s.newCache(cache.Options{})
	sb := c.Soundboard()

	sb.Set(gA, nil) // no-op
	sb.Set(gA, &discord.SoundboardSound{SoundID: 1, Name: "s1"})
	sb.Set(gB, &discord.SoundboardSound{SoundID: 2, Name: "s2"})

	got, ok := sb.Get(1)
	s.Require().True(ok)
	s.Equal("s1", got.Name)

	_, ok = sb.Get(999)
	s.False(ok)

	s.Equal(1, sb.GetByGuild(gA).Len())
	s.Equal(2, sb.Size())

	// SetAll replaces guild A's set, skipping nils.
	sb.SetAll(gA, []*discord.SoundboardSound{{SoundID: 10}, nil, {SoundID: 11}})
	s.Equal(2, sb.GetByGuild(gA).Len())

	sb.Delete(10)
	s.Equal(1, sb.GetByGuild(gA).Len())

	sb.DeleteGuild(gA)
	s.Equal(0, sb.GetByGuild(gA).Len())
	s.Equal(1, sb.Size())

	// SetAll with only nils → delete-only path.
	sb.SetAll(gB, []*discord.SoundboardSound{nil})
	s.Equal(0, sb.GetByGuild(gB).Len())
}

// ── ScheduledEvent store ──────────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestScheduledEventStore() {
	c := s.newCache(cache.Options{})
	se := c.ScheduledEvents()

	se.Set(nil) // no-op
	se.Set(&discord.GuildScheduledEvent{ID: 1, GuildID: gA, Name: "e1"})
	se.Set(&discord.GuildScheduledEvent{ID: 2, GuildID: gA})
	se.Set(&discord.GuildScheduledEvent{ID: 3, GuildID: gB})

	got, ok := se.Get(1)
	s.Require().True(ok)
	s.Equal("e1", got.Name)

	_, ok = se.Get(999)
	s.False(ok)

	s.Equal(2, se.GetByGuild(gA).Len())
	s.Equal(3, se.Size())

	se.Delete(1)
	s.Equal(1, se.GetByGuild(gA).Len())

	se.DeleteGuild(gA)
	s.Equal(0, se.GetByGuild(gA).Len())
	s.Equal(1, se.Size())
}

// ── StageInstance store ───────────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestStageInstanceStore() {
	c := s.newCache(cache.Options{})
	si := c.StageInstances()

	si.Set(nil) // no-op
	si.Set(&discord.StageInstance{ID: 1, GuildID: gA, Topic: "t1"})
	si.Set(&discord.StageInstance{ID: 2, GuildID: gA})
	si.Set(&discord.StageInstance{ID: 3, GuildID: gB})

	got, ok := si.Get(1)
	s.Require().True(ok)
	s.Equal("t1", got.Topic)

	_, ok = si.Get(999)
	s.False(ok)

	s.Equal(2, si.GetByGuild(gA).Len())
	s.Equal(3, si.Size())

	si.Delete(1)
	s.Equal(1, si.GetByGuild(gA).Len())

	si.DeleteGuild(gA)
	s.Equal(0, si.GetByGuild(gA).Len())
	s.Equal(1, si.Size())
}

// ── Emoji store ───────────────────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestEmojiStore() {
	c := s.newCache(cache.Options{})
	e := c.Emojis()

	e.Set(gA, nil) // no-op
	e.Set(gA, &discord.Emoji{ID: 1, Name: "e1"})
	e.Set(gB, &discord.Emoji{ID: 2})

	got, ok := e.Get(1)
	s.Require().True(ok)
	s.Equal("e1", got.Name)

	_, ok = e.Get(999)
	s.False(ok)

	s.Equal(1, e.GetByGuild(gA).Len())
	s.Equal(2, e.Size())

	// SetAll skips nil and empty-ID entries.
	e.SetAll(gA, []*discord.Emoji{{ID: 10}, nil, {ID: 0}, {ID: 11}})
	s.Equal(2, e.GetByGuild(gA).Len())

	e.Delete(10)
	s.Equal(1, e.GetByGuild(gA).Len())

	e.DeleteGuild(gA)
	s.Equal(0, e.GetByGuild(gA).Len())
	s.Equal(1, e.Size())
}

// ── Presence store ────────────────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestPresenceStore() {
	c := s.newCache(cache.Options{})
	p := c.Presences()

	p.Set(nil)                 // no-op
	p.Set(&discord.Presence{}) // empty User.ID → no-op
	p.Set(&discord.Presence{GuildID: gA, User: discord.PartialPresenceUser{ID: 1}})
	p.Set(&discord.Presence{GuildID: gA, User: discord.PartialPresenceUser{ID: 2}})
	p.Set(&discord.Presence{GuildID: gB, User: discord.PartialPresenceUser{ID: 3}})

	got, ok := p.Get(gA, 1)
	s.Require().True(ok)
	s.Equal(discord.Snowflake(1), got.User.ID)

	_, ok = p.Get(gA, 999)
	s.False(ok)

	s.Equal(2, p.GetByGuild(gA).Len())
	s.Equal(3, p.Size())

	p.Delete(gA, 1)
	s.Equal(1, p.GetByGuild(gA).Len())

	p.DeleteGuild(gA)
	s.Equal(0, p.GetByGuild(gA).Len())
	s.Equal(1, p.Size())
}

// ── Sticker SetAll skip-nil ───────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestStickerStore_SetAllSkipsNilAndEmpty() {
	c := s.newCache(cache.Options{})
	st := c.Stickers()

	st.Set(gA, nil)                      // no-op
	st.Set(gA, &discord.Sticker{ID: 99}) // establish the collection first
	st.SetAll(gA, []*discord.Sticker{{ID: 1}, nil, {ID: 0}, {ID: 2}})
	s.Equal(2, st.GetByGuild(gA).Len())
	s.Equal(2, st.Size())
}

// ── Size for the flat stores ──────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestSize_ChannelUserSticker() {
	c := s.newCache(cache.Options{})
	c.Channels().Set(channel("1"))
	c.Users().Set(user("2"))
	c.Stickers().Set(gA, sticker("3"))

	s.Equal(1, c.Channels().Size())
	s.Equal(1, c.Users().Size())
	s.Equal(1, c.Stickers().Size())
}

// ── Member store size (scans guild indices) ───────────────────────────────────

func (s *RedisCacheTestSuite) TestMemberStore_Size() {
	c := s.newCache(cache.Options{})
	c.Members().Set(gA, member("1"))
	c.Members().Set(gA, member("2"))
	c.Members().Set(gB, member("3"))

	s.Equal(3, c.Members().Size())
}

// ── TTL filter branches across every store ────────────────────────────────────

func (s *RedisCacheTestSuite) TestTTLFilterBranches_AllStores() {
	c := s.newCache(cache.Options{TTL: time.Hour, Messages: cache.MessageOptions{MaxPerChannel: 100, TTL: time.Hour}})

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
	c.Stickers().Set(gA, sticker("1"))
	c.Presences().Set(&discord.Presence{GuildID: gA, User: discord.PartialPresenceUser{ID: 1}})
	c.Messages().Add(message("m1", "c1"))

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
	s.Equal(1, c.Messages().Channel(mustSnowflake("c1")).Len())
}

// ── Disabled message store ────────────────────────────────────────────────────

func (s *RedisCacheTestSuite) TestMessageStore_DisabledAndNoops() {
	c := s.newCache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 0}})

	c.Messages().Add(message("m1", "c1")) // disabled → no-op
	c.Messages().Add(nil)                 // nil → no-op
	s.Equal(0, c.Messages().Channel(mustSnowflake("c1")).Len())
	s.Equal(0, c.Messages().Size())

	// DeleteBulk with empty slice → early return.
	c2 := s.newCache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 10}})
	c2.Messages().Add(message("m1", "c1"))
	c2.Messages().DeleteBulk(mustSnowflake("c1"), nil)
	s.Equal(1, c2.Messages().Size())
}
