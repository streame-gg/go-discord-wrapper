package stores

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// guildIndexedStoresSuite exercises the ID-keyed, guild-indexed stores
// (emoji, sticker, soundboard, scheduled event, stage instance).
type guildIndexedStoresSuite struct{ suite.Suite }

func TestGuildIndexedStoresSuite(t *testing.T) {
	suite.Run(t, new(guildIndexedStoresSuite))
}

const (
	gA = discord.Snowflake(100)
	gB = discord.Snowflake(200)
)

// ── MemEmojiStore ─────────────────────────────────────────────────────────────

func (su *guildIndexedStoresSuite) TestEmojiStore() {
	t := su.T()
	s := NewEmojiStore(Defaults())
	defer s.Close()

	assert.True(t, s.IsEnabled())
	s.Set(gA, nil) // no-op

	e1, e2, e3 := discord.Snowflake(1), discord.Snowflake(2), discord.Snowflake(3)
	s.Set(gA, &discord.Emoji{ID: e1})
	s.Set(gA, &discord.Emoji{ID: e2})
	s.Set(gB, &discord.Emoji{ID: e3})

	got, ok := s.Get(e1)
	require.True(t, ok)
	assert.Equal(t, e1, got.ID)

	assert.Equal(t, 2, s.GetByGuild(gA).Len())
	assert.Equal(t, 1, s.GetByGuild(gB).Len())
	assert.Equal(t, 3, s.Size())

	s.Delete(e1)
	assert.Equal(t, 1, s.GetByGuild(gA).Len())

	s.DeleteGuild(gA)
	assert.Equal(t, 0, s.GetByGuild(gA).Len())
	assert.Equal(t, 1, s.Size())

	// SetAll atomically replaces guild B's set.
	s.SetAll(gB, []*discord.Emoji{{ID: discord.Snowflake(4)}, nil, {ID: discord.Snowflake(5)}})
	assert.Equal(t, 2, s.GetByGuild(gB).Len())

	s.SweepExpired(0)
	assert.EqualValues(t, 0, s.Bytes())
	s.TrimTo(0)
	assert.Equal(t, 0, s.Size())
}

// ── MemStickerStore ───────────────────────────────────────────────────────────

func (su *guildIndexedStoresSuite) TestStickerStore() {
	t := su.T()
	s := NewStickerStore(Defaults())
	defer s.Close()

	assert.True(t, s.IsEnabled())
	s.Set(gA, nil)

	st1 := discord.Snowflake(11)
	s.Set(gA, &discord.Sticker{ID: st1})
	s.Set(gA, &discord.Sticker{ID: discord.Snowflake(12)})

	got, ok := s.Get(st1)
	require.True(t, ok)
	assert.Equal(t, st1, got.ID)
	assert.Equal(t, 2, s.GetByGuild(gA).Len())

	s.Delete(st1)
	assert.Equal(t, 1, s.GetByGuild(gA).Len())

	s.SetAll(gA, []*discord.Sticker{{ID: discord.Snowflake(13)}, nil})
	assert.Equal(t, 1, s.GetByGuild(gA).Len())

	s.DeleteGuild(gA)
	assert.Equal(t, 0, s.Size())

	s.SweepExpired(0)
	assert.EqualValues(t, 0, s.Bytes())
	s.TrimTo(0)
}

// ── MemSoundboardStore ────────────────────────────────────────────────────────

func (su *guildIndexedStoresSuite) TestSoundboardStore() {
	t := su.T()
	s := NewSoundboardStore(Defaults())
	defer s.Close()

	assert.True(t, s.IsEnabled())
	s.Set(gA, nil)

	snd := discord.Snowflake(21)
	s.Set(gA, &discord.SoundboardSound{SoundID: snd})
	s.Set(gA, &discord.SoundboardSound{SoundID: discord.Snowflake(22)})

	got, ok := s.Get(snd)
	require.True(t, ok)
	assert.Equal(t, snd, got.SoundID)
	assert.Equal(t, 2, s.GetByGuild(gA).Len())

	s.Delete(snd)
	assert.Equal(t, 1, s.GetByGuild(gA).Len())

	s.SetAll(gA, []*discord.SoundboardSound{{SoundID: discord.Snowflake(23)}, nil})
	assert.Equal(t, 1, s.GetByGuild(gA).Len())

	s.DeleteGuild(gA)
	assert.Equal(t, 0, s.Size())

	s.SweepExpired(0)
	assert.EqualValues(t, 0, s.Bytes())
	s.TrimTo(0)
}

// ── MemScheduledEventStore ────────────────────────────────────────────────────

func (su *guildIndexedStoresSuite) TestScheduledEventStore() {
	t := su.T()
	s := NewScheduledEventStore(Defaults())
	defer s.Close()

	assert.True(t, s.IsEnabled())
	s.Set(nil)

	ev := discord.Snowflake(31)
	s.Set(&discord.GuildScheduledEvent{ID: ev, GuildID: gA})
	s.Set(&discord.GuildScheduledEvent{ID: discord.Snowflake(32), GuildID: gA})
	s.Set(&discord.GuildScheduledEvent{ID: discord.Snowflake(33), GuildID: gB})

	got, ok := s.Get(ev)
	require.True(t, ok)
	assert.Equal(t, ev, got.ID)
	assert.Equal(t, 2, s.GetByGuild(gA).Len())

	s.Delete(ev)
	assert.Equal(t, 1, s.GetByGuild(gA).Len())

	s.DeleteGuild(gA)
	assert.Equal(t, 1, s.Size())

	s.SweepExpired(0)
	assert.EqualValues(t, 0, s.Bytes())
	s.TrimTo(0)
}

// ── MemStageInstanceStore ─────────────────────────────────────────────────────

func (su *guildIndexedStoresSuite) TestStageInstanceStore() {
	t := su.T()
	s := NewStageInstanceStore(Defaults())
	defer s.Close()

	assert.True(t, s.IsEnabled())
	s.Set(nil)

	si := discord.Snowflake(41)
	s.Set(&discord.StageInstance{ID: si, GuildID: gA})
	s.Set(&discord.StageInstance{ID: discord.Snowflake(42), GuildID: gA})

	got, ok := s.Get(si)
	require.True(t, ok)
	assert.Equal(t, si, got.ID)
	assert.Equal(t, 2, s.GetByGuild(gA).Len())

	s.Delete(si)
	assert.Equal(t, 1, s.GetByGuild(gA).Len())

	s.DeleteGuild(gA)
	assert.Equal(t, 0, s.Size())

	s.SweepExpired(0)
	assert.EqualValues(t, 0, s.Bytes())
	s.TrimTo(0)
}
