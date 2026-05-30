package stores

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func newTestInvite(code string) *discord.Invite {
	return &discord.Invite{Code: code}
}

func newTestInviteWithGuild(code string, guildID discord.Snowflake) *discord.Invite {
	return &discord.Invite{Code: code, Guild: &discord.Guild{ID: guildID}}
}

// ── MemInviteStore ────────────────────────────────────────────────────────────

func TestMemInviteStore_SetGet(t *testing.T) {
	s := NewInviteStore(Defaults())
	defer s.Close()

	inv := newTestInvite("abc123")
	s.Set(inv)

	got, ok := s.Get("abc123")
	require.True(t, ok)
	assert.Equal(t, "abc123", got.Code)
}

func TestMemInviteStore_GetMiss(t *testing.T) {
	s := NewInviteStore(Defaults())
	defer s.Close()

	_, ok := s.Get("nope")
	assert.False(t, ok)
}

func TestMemInviteStore_SetNilIsNoop(t *testing.T) {
	s := NewInviteStore(Defaults())
	defer s.Close()

	s.Set(nil)
	assert.Equal(t, 0, s.Size())
}

func TestMemInviteStore_SetEmptyCodeIsNoop(t *testing.T) {
	s := NewInviteStore(Defaults())
	defer s.Close()

	s.Set(&discord.Invite{Code: ""})
	assert.Equal(t, 0, s.Size())
}

func TestMemInviteStore_Delete(t *testing.T) {
	s := NewInviteStore(Defaults())
	defer s.Close()

	s.Set(newTestInvite("abc123"))
	s.Delete("abc123")

	_, ok := s.Get("abc123")
	assert.False(t, ok)
	assert.Equal(t, 0, s.Size())
}

func TestMemInviteStore_SetWithGuild_GetByGuild(t *testing.T) {
	s := NewInviteStore(Defaults())
	defer s.Close()

	guildID := discord.Snowflake(1000)
	s.SetWithGuild(guildID, newTestInvite("code1"))
	s.SetWithGuild(guildID, newTestInvite("code2"))

	all := s.GetByGuild(guildID)
	assert.Equal(t, 2, all.Len())

	_, ok1 := all.Get("code1")
	_, ok2 := all.Get("code2")
	assert.True(t, ok1)
	assert.True(t, ok2)
}

func TestMemInviteStore_Set_WithGuildField_GetByGuild(t *testing.T) {
	s := NewInviteStore(Defaults())
	defer s.Close()

	guildID := discord.Snowflake(1000)
	s.Set(newTestInviteWithGuild("xyzabc", guildID))

	all := s.GetByGuild(guildID)
	assert.Equal(t, 1, all.Len())
}

func TestMemInviteStore_GetByGuild_OtherGuildIsolated(t *testing.T) {
	s := NewInviteStore(Defaults())
	defer s.Close()

	guildA := discord.Snowflake(111)
	guildB := discord.Snowflake(222)

	s.SetWithGuild(guildA, newTestInvite("codeA"))
	s.SetWithGuild(guildB, newTestInvite("codeB"))

	allA := s.GetByGuild(guildA)
	assert.Equal(t, 1, allA.Len())
	_, okB := allA.Get("codeB")
	assert.False(t, okB, "guild B invite must not appear in guild A GetByGuild")
}

func TestMemInviteStore_DeleteGuild(t *testing.T) {
	s := NewInviteStore(Defaults())
	defer s.Close()

	guildA := discord.Snowflake(100)
	guildB := discord.Snowflake(200)

	s.SetWithGuild(guildA, newTestInvite("codeA1"))
	s.SetWithGuild(guildA, newTestInvite("codeA2"))
	s.SetWithGuild(guildB, newTestInvite("codeB1"))

	s.DeleteGuild(guildA)

	_, okA1 := s.Get("codeA1")
	_, okA2 := s.Get("codeA2")
	_, okB1 := s.Get("codeB1")
	assert.False(t, okA1, "codeA1 must be gone after DeleteGuild(guildA)")
	assert.False(t, okA2, "codeA2 must be gone after DeleteGuild(guildA)")
	assert.True(t, okB1, "codeB1 in guild B must survive DeleteGuild(guildA)")
	assert.Equal(t, 1, s.Size())
}

func TestMemInviteStore_Delete_RemovesFromGuildIndex(t *testing.T) {
	s := NewInviteStore(Defaults())
	defer s.Close()

	guildID := discord.Snowflake(1000)
	s.SetWithGuild(guildID, newTestInvite("code1"))
	s.SetWithGuild(guildID, newTestInvite("code2"))

	s.Delete("code1")

	all := s.GetByGuild(guildID)
	assert.Equal(t, 1, all.Len(), "deleted invite must not appear in GetByGuild")
	_, ok := all.Get("code1")
	assert.False(t, ok)
}
