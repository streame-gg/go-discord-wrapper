package mongocache_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/cache/mongocache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// mongoStoresSuite exercises the guild-scoped and composite-key stores end to
// end against a real MongoDB container. Each test gets an isolated database.
type mongoStoresSuite struct {
	suite.Suite
	db *mongo.Database // database backing the most recent cache() call
}

func TestMongoStoresSuite(t *testing.T) { suite.Run(t, new(mongoStoresSuite)) }

// cache builds a sync-write cache backed by a fresh database named after the
// running test, dropped on cleanup. The backing database is exposed via s.db so
// tests can inject corrupt documents directly.
func (s *mongoStoresSuite) cache(opts cache.Options) *mongocache.MongoDBCache {
	t := s.T()
	if mongoClient == nil {
		t.Fatal("MongoDB client not initialised — TestMain should have failed first")
	}
	db := mongoClient.Database(uniqueDBName(t.Name()))
	s.db = db
	t.Cleanup(func() { _ = db.Drop(context.Background()) })
	c := mongocache.NewMongoDBCache(db, opts)
	c.EnableSyncWrites()
	t.Cleanup(func() { c.Close() })
	return c
}

// asyncCache is like cache but leaves the async write worker enabled.
func (s *mongoStoresSuite) asyncCache(opts cache.Options) *mongocache.MongoDBCache {
	t := s.T()
	if mongoClient == nil {
		t.Fatal("MongoDB client not initialised — TestMain should have failed first")
	}
	db := mongoClient.Database(uniqueDBName(t.Name()))
	s.db = db
	t.Cleanup(func() { _ = db.Drop(context.Background()) })
	c := mongocache.NewMongoDBCache(db, opts)
	t.Cleanup(func() { c.Close() })
	return c
}

const (
	gA = discord.Snowflake(1001)
	gB = discord.Snowflake(2002)
)

// ── accessors ─────────────────────────────────────────────────────────────────

func (s *mongoStoresSuite) TestAccessors_NotNil() {
	c := s.cache(cache.Options{})
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
	s.NotNil(c.Bans())
	s.NotNil(c.AutoModRules())
	s.NotNil(c.Invites())
}

// ── Role store ────────────────────────────────────────────────────────────────

func (s *mongoStoresSuite) TestRoleStore() {
	t := s.T()
	c := s.cache(cache.Options{})
	r := c.Roles()

	r.Set(gA, nil) // no-op
	r.Set(gA, &discord.Role{ID: 1, Name: "r1"})
	r.Set(gA, &discord.Role{ID: 2, Name: "r2"})
	r.Set(gB, &discord.Role{ID: 3, Name: "r3"})

	got, ok := r.Get(1)
	require.True(t, ok)
	assert.Equal(t, "r1", got.Name)

	_, ok = r.Get(999)
	assert.False(t, ok)

	assert.Equal(t, 2, r.GetByGuild(gA).Len())
	assert.Equal(t, 3, r.All().Len())
	assert.Equal(t, 3, r.Size())

	r.Delete(1)
	assert.Equal(t, 1, r.GetByGuild(gA).Len())

	r.DeleteGuild(gA)
	assert.Equal(t, 0, r.GetByGuild(gA).Len())
	assert.Equal(t, 1, r.Size())
}

// ── VoiceState store ──────────────────────────────────────────────────────────

func (s *mongoStoresSuite) TestVoiceStateStore() {
	t := s.T()
	c := s.cache(cache.Options{})
	v := c.VoiceStates()

	v.Set(gA, nil) // no-op
	v.Set(gA, &discord.VoiceState{UserID: 1})
	v.Set(gA, &discord.VoiceState{UserID: 2})
	v.Set(gB, &discord.VoiceState{UserID: 3})

	got, ok := v.Get(gA, 1)
	require.True(t, ok)
	assert.Equal(t, discord.Snowflake(1), got.UserID)

	_, ok = v.Get(gA, 999)
	assert.False(t, ok)

	assert.Equal(t, 2, v.AllInGuild(gA).Len())
	assert.Equal(t, 3, v.Size())

	v.Delete(gA, 1)
	assert.Equal(t, 1, v.AllInGuild(gA).Len())

	v.DeleteGuild(gA)
	assert.Equal(t, 0, v.AllInGuild(gA).Len())
	assert.Equal(t, 1, v.Size())
}

// ── Soundboard store ──────────────────────────────────────────────────────────

func (s *mongoStoresSuite) TestSoundboardStore() {
	t := s.T()
	c := s.cache(cache.Options{})
	sb := c.Soundboard()

	sb.Set(gA, nil) // no-op
	sb.Set(gA, &discord.SoundboardSound{SoundID: 1, Name: "s1"})
	sb.Set(gB, &discord.SoundboardSound{SoundID: 2, Name: "s2"})

	got, ok := sb.Get(1)
	require.True(t, ok)
	assert.Equal(t, "s1", got.Name)

	_, ok = sb.Get(999)
	assert.False(t, ok)

	assert.Equal(t, 1, sb.GetByGuild(gA).Len())
	assert.Equal(t, 2, sb.Size())

	// SetAll replaces guild A's set, skipping nils.
	sb.SetAll(gA, []*discord.SoundboardSound{{SoundID: 10}, nil, {SoundID: 11}})
	assert.Equal(t, 2, sb.GetByGuild(gA).Len())

	sb.Delete(10)
	assert.Equal(t, 1, sb.GetByGuild(gA).Len())

	sb.DeleteGuild(gA)
	assert.Equal(t, 0, sb.GetByGuild(gA).Len())
	assert.Equal(t, 1, sb.Size())

	// SetAll with only nils → empty docs slice (delete-only path).
	sb.SetAll(gB, []*discord.SoundboardSound{nil})
	assert.Equal(t, 0, sb.GetByGuild(gB).Len())
}

// ── ScheduledEvent store ──────────────────────────────────────────────────────

func (s *mongoStoresSuite) TestScheduledEventStore() {
	t := s.T()
	c := s.cache(cache.Options{})
	se := c.ScheduledEvents()

	se.Set(nil) // no-op
	se.Set(&discord.GuildScheduledEvent{ID: 1, GuildID: gA, Name: "e1"})
	se.Set(&discord.GuildScheduledEvent{ID: 2, GuildID: gA})
	se.Set(&discord.GuildScheduledEvent{ID: 3, GuildID: gB})

	got, ok := se.Get(1)
	require.True(t, ok)
	assert.Equal(t, "e1", got.Name)

	_, ok = se.Get(999)
	assert.False(t, ok)

	assert.Equal(t, 2, se.GetByGuild(gA).Len())
	assert.Equal(t, 3, se.Size())

	se.Delete(1)
	assert.Equal(t, 1, se.GetByGuild(gA).Len())

	se.DeleteGuild(gA)
	assert.Equal(t, 0, se.GetByGuild(gA).Len())
	assert.Equal(t, 1, se.Size())
}

// ── StageInstance store ───────────────────────────────────────────────────────

func (s *mongoStoresSuite) TestStageInstanceStore() {
	t := s.T()
	c := s.cache(cache.Options{})
	si := c.StageInstances()

	si.Set(nil) // no-op
	si.Set(&discord.StageInstance{ID: 1, GuildID: gA, Topic: "t1"})
	si.Set(&discord.StageInstance{ID: 2, GuildID: gA})
	si.Set(&discord.StageInstance{ID: 3, GuildID: gB})

	got, ok := si.Get(1)
	require.True(t, ok)
	assert.Equal(t, "t1", got.Topic)

	_, ok = si.Get(999)
	assert.False(t, ok)

	assert.Equal(t, 2, si.GetByGuild(gA).Len())
	assert.Equal(t, 3, si.Size())

	si.Delete(1)
	assert.Equal(t, 1, si.GetByGuild(gA).Len())

	si.DeleteGuild(gA)
	assert.Equal(t, 0, si.GetByGuild(gA).Len())
	assert.Equal(t, 1, si.Size())
}

// ── Emoji store ───────────────────────────────────────────────────────────────

func (s *mongoStoresSuite) TestEmojiStore() {
	t := s.T()
	c := s.cache(cache.Options{})
	e := c.Emojis()

	e.Set(gA, nil) // no-op
	e.Set(gA, &discord.Emoji{ID: 1, Name: "e1"})
	e.Set(gB, &discord.Emoji{ID: 2})

	got, ok := e.Get(1)
	require.True(t, ok)
	assert.Equal(t, "e1", got.Name)

	_, ok = e.Get(999)
	assert.False(t, ok)

	assert.Equal(t, 1, e.GetByGuild(gA).Len())
	assert.Equal(t, 2, e.Size())

	// SetAll skips nil and empty-ID entries.
	e.SetAll(gA, []*discord.Emoji{{ID: 10}, nil, {ID: 0}, {ID: 11}})
	assert.Equal(t, 2, e.GetByGuild(gA).Len())

	e.Delete(10)
	assert.Equal(t, 1, e.GetByGuild(gA).Len())

	e.DeleteGuild(gA)
	assert.Equal(t, 0, e.GetByGuild(gA).Len())
	assert.Equal(t, 1, e.Size())
}

// ── Presence store ────────────────────────────────────────────────────────────

func (s *mongoStoresSuite) TestPresenceStore() {
	t := s.T()
	c := s.cache(cache.Options{})
	p := c.Presences()

	p.Set(nil)                 // no-op
	p.Set(&discord.Presence{}) // empty User.ID → no-op
	p.Set(&discord.Presence{GuildID: gA, User: discord.PartialPresenceUser{ID: 1}})
	p.Set(&discord.Presence{GuildID: gA, User: discord.PartialPresenceUser{ID: 2}})
	p.Set(&discord.Presence{GuildID: gB, User: discord.PartialPresenceUser{ID: 3}})

	got, ok := p.Get(gA, 1)
	require.True(t, ok)
	assert.Equal(t, discord.Snowflake(1), got.User.ID)

	_, ok = p.Get(gA, 999)
	assert.False(t, ok)

	assert.Equal(t, 2, p.GetByGuild(gA).Len())
	assert.Equal(t, 3, p.Size())

	p.Delete(gA, 1)
	assert.Equal(t, 1, p.GetByGuild(gA).Len())

	p.DeleteGuild(gA)
	assert.Equal(t, 0, p.GetByGuild(gA).Len())
	assert.Equal(t, 1, p.Size())
}

// ── Ban store ─────────────────────────────────────────────────────────────────

func (s *mongoStoresSuite) TestBanStore() {
	t := s.T()
	c := s.cache(cache.Options{})
	b := c.Bans()

	b.Set(gA, nil) // no-op
	b.Set(gA, &discord.Ban{User: discord.User{ID: 1}})
	b.Set(gA, &discord.Ban{User: discord.User{ID: 2}})
	b.Set(gB, &discord.Ban{User: discord.User{ID: 3}})

	got, ok := b.Get(gA, 1)
	require.True(t, ok)
	assert.Equal(t, discord.Snowflake(1), got.User.ID)

	_, ok = b.Get(gA, 999)
	assert.False(t, ok)

	assert.Equal(t, 2, b.AllInGuild(gA).Len())
	assert.Equal(t, 3, b.Size())

	b.Delete(gA, 1)
	assert.Equal(t, 1, b.AllInGuild(gA).Len())

	b.DeleteGuild(gA)
	assert.Equal(t, 0, b.AllInGuild(gA).Len())
	assert.Equal(t, 1, b.Size())
}

// ── AutoModeration rule store ─────────────────────────────────────────────────

func (s *mongoStoresSuite) TestAutoModRuleStore() {
	t := s.T()
	c := s.cache(cache.Options{})
	a := c.AutoModRules()

	a.Set(gA, nil) // no-op
	a.Set(gA, &discord.AutoModerationRule{ID: 1, GuildID: gA, Name: "r1"})
	a.Set(gA, &discord.AutoModerationRule{ID: 2, GuildID: gA})
	a.Set(gB, &discord.AutoModerationRule{ID: 3, GuildID: gB})

	got, ok := a.Get(1)
	require.True(t, ok)
	assert.Equal(t, "r1", got.Name)

	_, ok = a.Get(999)
	assert.False(t, ok)

	assert.Equal(t, 2, a.GetByGuild(gA).Len())
	assert.Equal(t, 3, a.Size())

	a.Delete(1)
	assert.Equal(t, 1, a.GetByGuild(gA).Len())

	a.DeleteGuild(gA)
	assert.Equal(t, 0, a.GetByGuild(gA).Len())
	assert.Equal(t, 1, a.Size())
}

// ── Invite store ──────────────────────────────────────────────────────────────

func (s *mongoStoresSuite) TestInviteStore() {
	t := s.T()
	c := s.cache(cache.Options{})
	i := c.Invites()

	i.Set(nil)                                // no-op
	i.Set(&discord.Invite{})                  // empty code → no-op
	i.Set(&discord.Invite{Code: "guildless"}) // stored, no guild association

	// Guild-associated via the invite's own Guild field.
	i.Set(&discord.Invite{Code: "abc", Guild: &discord.Guild{ID: gA}})
	// Guild-associated via SetWithGuild (gateway events lack Guild).
	i.SetWithGuild(gA, &discord.Invite{Code: "def"})
	i.SetWithGuild(gB, &discord.Invite{Code: "ghi"})
	i.SetWithGuild(gA, nil)               // no-op
	i.SetWithGuild(gA, &discord.Invite{}) // empty code → no-op

	got, ok := i.Get("abc")
	require.True(t, ok)
	assert.Equal(t, "abc", got.Code)

	got, ok = i.Get("guildless")
	require.True(t, ok)
	assert.Equal(t, "guildless", got.Code)

	_, ok = i.Get("missing")
	assert.False(t, ok)

	assert.Equal(t, 2, i.GetByGuild(gA).Len())
	assert.Equal(t, 1, i.GetByGuild(gB).Len())
	assert.Equal(t, 4, i.Size()) // guildless + abc + def + ghi

	i.Delete("abc")
	assert.Equal(t, 1, i.GetByGuild(gA).Len())

	i.DeleteGuild(gA)
	assert.Equal(t, 0, i.GetByGuild(gA).Len())
	_, ok = i.Get("def")
	assert.False(t, ok, "DeleteGuild must remove the guild's invites")
}

// ── Sticker SetAll skip-nil ───────────────────────────────────────────────────

func (s *mongoStoresSuite) TestStickerStore_SetAllSkipsNilAndEmpty() {
	c := s.cache(cache.Options{})
	st := c.Stickers()

	st.Set(gA, nil)                      // no-op
	st.Set(gA, &discord.Sticker{ID: 99}) // establish the collection first
	st.SetAll(gA, []*discord.Sticker{{ID: 1}, nil, {ID: 0}, {ID: 2}})
	s.Equal(2, st.GetByGuild(gA).Len())
}
