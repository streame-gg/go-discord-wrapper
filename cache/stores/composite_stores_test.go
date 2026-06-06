package stores

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// compositeStoresSuite exercises the composite-key stores (member, voice state,
// presence) end to end, covering CRUD, guild isolation and no-op guards.
type compositeStoresSuite struct{ suite.Suite }

func TestCompositeStoresSuite(t *testing.T) {
	suite.Run(t, new(compositeStoresSuite))
}

// ── MemMemberStore ────────────────────────────────────────────────────────────

func (su *compositeStoresSuite) TestMemberStore() {
	t := su.T()
	s := NewMemberStore(Defaults())
	defer s.Close()

	assert.True(t, s.IsEnabled())

	// no-ops
	s.Set(gA, nil)
	s.Set(gA, &discord.GuildMember{User: nil})
	assert.Equal(t, 0, s.Size())

	u1, u2, u3 := discord.Snowflake(1), discord.Snowflake(2), discord.Snowflake(3)
	s.Set(gA, newTestMember(u1))
	s.Set(gA, newTestMember(u2))
	s.Set(gB, newTestMember(u3))

	got, ok := s.Get(gA, u1)
	require.True(t, ok)
	assert.Equal(t, u1, got.User.ID)

	assert.True(t, s.Has(gA, u1))
	assert.False(t, s.Has(gA, discord.Snowflake(999)))
	assert.Equal(t, 3, s.Size())
	assert.Equal(t, 2, s.AllInGuild(gA).Len())
	assert.Equal(t, 1, s.AllInGuild(gB).Len())

	s.Delete(gA, u1)
	assert.False(t, s.Has(gA, u1))
	assert.Equal(t, 1, s.AllInGuild(gA).Len())

	s.DeleteGuild(gA)
	assert.Equal(t, 0, s.AllInGuild(gA).Len())
	assert.Equal(t, 1, s.Size())

	s.SweepExpired(0)
	assert.EqualValues(t, 0, s.Bytes())
	s.TrimTo(0)
	assert.Equal(t, 0, s.Size())
}

// ── MemVoiceStateStore ────────────────────────────────────────────────────────

func (su *compositeStoresSuite) TestVoiceStateStore() {
	t := su.T()
	s := NewVoiceStateStore(Defaults())
	defer s.Close()

	assert.True(t, s.IsEnabled())
	s.Set(gA, nil) // no-op
	assert.Equal(t, 0, s.Size())

	u1, u2 := discord.Snowflake(1), discord.Snowflake(2)
	s.Set(gA, newTestVoiceState(u1))
	s.Set(gA, newTestVoiceState(u2))

	got, ok := s.Get(gA, u1)
	require.True(t, ok)
	assert.Equal(t, u1, got.UserID)

	assert.True(t, s.Has(gA, u1))
	assert.False(t, s.Has(gA, discord.Snowflake(999)))
	assert.Equal(t, 2, s.AllInGuild(gA).Len())

	s.Delete(gA, u1)
	assert.Equal(t, 1, s.AllInGuild(gA).Len())

	s.DeleteGuild(gA)
	assert.Equal(t, 0, s.Size())

	s.SweepExpired(0)
	assert.EqualValues(t, 0, s.Bytes())
	s.TrimTo(0)
}

// ── MemPresenceStore ──────────────────────────────────────────────────────────

func (su *compositeStoresSuite) TestPresenceStore() {
	t := su.T()
	// PresenceStore is disabled by default; enable it explicitly here.
	s := NewPresenceStore(StoreOptions{})
	defer s.Close()

	assert.True(t, s.IsEnabled())

	// no-ops
	s.Set(nil)
	s.Set(&discord.Presence{}) // empty User.ID
	assert.Equal(t, 0, s.Size())

	u1, u2 := discord.Snowflake(1), discord.Snowflake(2)
	s.Set(newTestPresence(gA, u1))
	s.Set(newTestPresence(gA, u2))

	got, ok := s.Get(gA, u1)
	require.True(t, ok)
	assert.Equal(t, u1, got.User.ID)

	assert.Equal(t, 2, s.GetByGuild(gA).Len())

	s.Delete(gA, u1)
	assert.Equal(t, 1, s.GetByGuild(gA).Len())

	s.DeleteGuild(gA)
	assert.Equal(t, 0, s.Size())

	s.SweepExpired(0)
	assert.EqualValues(t, 0, s.Bytes())
	s.TrimTo(0)
}

func (su *compositeStoresSuite) TestPresenceStore_DisabledByDefault() {
	t := su.T()
	s := NewPresenceStore(PresenceDefaults())
	defer s.Close()

	assert.False(t, s.IsEnabled())
	s.Set(newTestPresence(gA, discord.Snowflake(1)))
	assert.Equal(t, 0, s.Size())
}
