package stores

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func newTestAutoModRule(guildID, ruleID discord.Snowflake) *discord.AutoModerationRule {
	return &discord.AutoModerationRule{
		ID:      ruleID,
		GuildID: guildID,
		Name:    "test-rule",
	}
}

// ── MemAutoModerationRuleStore ────────────────────────────────────────────────

func TestMemAutoModStore_SetGet(t *testing.T) {
	s := NewAutoModerationRuleStore(Defaults())
	defer s.Close()

	guildID := discord.Snowflake(1000)
	ruleID := discord.Snowflake(5000)
	rule := newTestAutoModRule(guildID, ruleID)

	s.Set(guildID, rule)

	got, ok := s.Get(ruleID)
	require.True(t, ok)
	assert.Equal(t, "test-rule", got.Name)
}

func TestMemAutoModStore_GetMiss(t *testing.T) {
	s := NewAutoModerationRuleStore(Defaults())
	defer s.Close()

	_, ok := s.Get(discord.Snowflake(999))
	assert.False(t, ok)
}

func TestMemAutoModStore_SetNilIsNoop(t *testing.T) {
	s := NewAutoModerationRuleStore(Defaults())
	defer s.Close()

	s.Set(discord.Snowflake(1), nil)
	assert.Equal(t, 0, s.Size())
}

func TestMemAutoModStore_Delete(t *testing.T) {
	s := NewAutoModerationRuleStore(Defaults())
	defer s.Close()

	guildID := discord.Snowflake(1000)
	ruleID := discord.Snowflake(5000)
	s.Set(guildID, newTestAutoModRule(guildID, ruleID))
	s.Delete(ruleID)

	_, ok := s.Get(ruleID)
	assert.False(t, ok)
	assert.Equal(t, 0, s.Size())
}

func TestMemAutoModStore_DeleteGuild(t *testing.T) {
	s := NewAutoModerationRuleStore(Defaults())
	defer s.Close()

	guildA := discord.Snowflake(100)
	guildB := discord.Snowflake(200)
	ruleA1 := discord.Snowflake(1)
	ruleA2 := discord.Snowflake(2)
	ruleB1 := discord.Snowflake(3)

	s.Set(guildA, newTestAutoModRule(guildA, ruleA1))
	s.Set(guildA, newTestAutoModRule(guildA, ruleA2))
	s.Set(guildB, newTestAutoModRule(guildB, ruleB1))

	s.DeleteGuild(guildA)

	_, ok1 := s.Get(ruleA1)
	_, ok2 := s.Get(ruleA2)
	_, okB := s.Get(ruleB1)
	assert.False(t, ok1, "ruleA1 must be gone after DeleteGuild(guildA)")
	assert.False(t, ok2, "ruleA2 must be gone after DeleteGuild(guildA)")
	assert.True(t, okB, "ruleB1 in guild B must survive DeleteGuild(guildA)")
	assert.Equal(t, 1, s.Size())
}

func TestMemAutoModStore_GetByGuild(t *testing.T) {
	s := NewAutoModerationRuleStore(Defaults())
	defer s.Close()

	guildID := discord.Snowflake(1000)
	r1 := discord.Snowflake(1)
	r2 := discord.Snowflake(2)

	s.Set(guildID, newTestAutoModRule(guildID, r1))
	s.Set(guildID, newTestAutoModRule(guildID, r2))

	all := s.GetByGuild(guildID)
	assert.Equal(t, 2, all.Len())

	_, ok1 := all.Get(r1)
	_, ok2 := all.Get(r2)
	assert.True(t, ok1)
	assert.True(t, ok2)
}

func TestMemAutoModStore_GetByGuild_OtherGuildIsolated(t *testing.T) {
	s := NewAutoModerationRuleStore(Defaults())
	defer s.Close()

	guildA := discord.Snowflake(111)
	guildB := discord.Snowflake(222)
	ruleA := discord.Snowflake(1)
	ruleB := discord.Snowflake(2)

	s.Set(guildA, newTestAutoModRule(guildA, ruleA))
	s.Set(guildB, newTestAutoModRule(guildB, ruleB))

	allA := s.GetByGuild(guildA)
	assert.Equal(t, 1, allA.Len())
	_, okB := allA.Get(ruleB)
	assert.False(t, okB, "guild B rule must not appear in guild A GetByGuild")
}
