package stores

import (
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func newTestBan(guildID, userID discord.Snowflake) *discord.Ban {
	reason := "spam"
	return &discord.Ban{Reason: &reason, User: discord.User{ID: userID}}
}

// ── MemBanStore ───────────────────────────────────────────────────────────────

func (su *storesBanSuite) TestMemBanStore_SetGet() {
	t := su.T()
	s := NewBanStore(Defaults())
	defer s.Close()

	guildID := discord.Snowflake(1000)
	userID := discord.Snowflake(2000)
	ban := newTestBan(guildID, userID)

	s.Set(guildID, ban)

	got, ok := s.Get(guildID, userID)
	require.True(t, ok)
	assert.Equal(t, ban.Reason, got.Reason)
	assert.Equal(t, userID, got.User.ID)
}

func (su *storesBanSuite) TestMemBanStore_GetMiss() {
	t := su.T()
	s := NewBanStore(Defaults())
	defer s.Close()

	_, ok := s.Get(discord.Snowflake(1), discord.Snowflake(2))
	assert.False(t, ok)
}

func (su *storesBanSuite) TestMemBanStore_SetNilIsNoop() {
	t := su.T()
	s := NewBanStore(Defaults())
	defer s.Close()

	s.Set(discord.Snowflake(1), nil)
	assert.Equal(t, 0, s.Size())
}

func (su *storesBanSuite) TestMemBanStore_Delete() {
	t := su.T()
	s := NewBanStore(Defaults())
	defer s.Close()

	guildID := discord.Snowflake(1000)
	userID := discord.Snowflake(2000)
	s.Set(guildID, newTestBan(guildID, userID))
	s.Delete(guildID, userID)

	_, ok := s.Get(guildID, userID)
	assert.False(t, ok)
	assert.Equal(t, 0, s.Size())
}

func (su *storesBanSuite) TestMemBanStore_DeleteGuild() {
	t := su.T()
	s := NewBanStore(Defaults())
	defer s.Close()

	guildID := discord.Snowflake(1000)
	user1 := discord.Snowflake(2001)
	user2 := discord.Snowflake(2002)
	other := discord.Snowflake(9999)
	otherUser := discord.Snowflake(8888)

	s.Set(guildID, newTestBan(guildID, user1))
	s.Set(guildID, newTestBan(guildID, user2))
	s.Set(other, newTestBan(other, otherUser))

	s.DeleteGuild(guildID)

	_, ok1 := s.Get(guildID, user1)
	_, ok2 := s.Get(guildID, user2)
	_, okOther := s.Get(other, otherUser)
	assert.False(t, ok1, "ban for user1 should be gone after DeleteGuild")
	assert.False(t, ok2, "ban for user2 should be gone after DeleteGuild")
	assert.True(t, okOther, "ban in other guild must survive DeleteGuild")
	assert.Equal(t, 1, s.Size())
}

func (su *storesBanSuite) TestMemBanStore_AllInGuild() {
	t := su.T()
	s := NewBanStore(Defaults())
	defer s.Close()

	guildID := discord.Snowflake(1000)
	user1 := discord.Snowflake(2001)
	user2 := discord.Snowflake(2002)

	s.Set(guildID, newTestBan(guildID, user1))
	s.Set(guildID, newTestBan(guildID, user2))

	all := s.AllInGuild(guildID)
	assert.Equal(t, 2, all.Len())

	_, ok1 := all.Get(user1)
	_, ok2 := all.Get(user2)
	assert.True(t, ok1)
	assert.True(t, ok2)
}

func (su *storesBanSuite) TestMemBanStore_AllInGuild_OtherGuildIsolated() {
	t := su.T()
	s := NewBanStore(Defaults())
	defer s.Close()

	guildA := discord.Snowflake(111)
	guildB := discord.Snowflake(222)
	userA := discord.Snowflake(11)
	userB := discord.Snowflake(22)

	s.Set(guildA, newTestBan(guildA, userA))
	s.Set(guildB, newTestBan(guildB, userB))

	allA := s.AllInGuild(guildA)
	assert.Equal(t, 1, allA.Len(), "guild A should only see its own bans")
	_, okB := allA.Get(userB)
	assert.False(t, okB, "guild B user must not appear in guild A AllInGuild")
}

type storesBanSuite struct{ suite.Suite }

func TestStoresBanSuite(t *testing.T) { suite.Run(t, new(storesBanSuite)) }
