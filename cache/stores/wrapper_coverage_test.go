package stores

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// This file fills coverage gaps on the role, invite, ban and automod stores.
// Each test method is attached to the suite already declared in that store's
// primary test file.

// ── MemRoleStore ──────────────────────────────────────────────────────────────

func (su *storesRoleSuite) TestMemRoleStore_FullCoverage() {
	t := su.T()
	s := NewRoleStore(Defaults())
	defer s.Close()

	assert.True(t, s.IsEnabled())
	s.Set(gA, nil) // no-op

	r1, r2, r3 := discord.Snowflake(1), discord.Snowflake(2), discord.Snowflake(3)
	s.Set(gA, &discord.Role{ID: r1})
	s.Set(gA, &discord.Role{ID: r2})
	s.Set(gB, &discord.Role{ID: r3})

	got, ok := s.Get(r1)
	require.True(t, ok)
	assert.Equal(t, r1, got.ID)

	assert.Equal(t, 2, s.GetByGuild(gA).Len())
	assert.Equal(t, 3, s.All().Len())
	assert.Equal(t, 3, s.Size())

	s.Delete(r1)
	assert.Equal(t, 1, s.GetByGuild(gA).Len())

	// SetAll atomically replaces guild A's set, skipping nils.
	s.SetAll(gA, []*discord.Role{{ID: discord.Snowflake(10)}, nil, {ID: discord.Snowflake(11)}})
	assert.Equal(t, 2, s.GetByGuild(gA).Len())

	s.DeleteGuild(gA)
	assert.Equal(t, 0, s.GetByGuild(gA).Len())
	assert.Equal(t, 1, s.Size())

	s.SweepExpired(0)
	assert.EqualValues(t, 0, s.Bytes())
	s.TrimTo(0)
	assert.Equal(t, 0, s.Size())
}

// ── MemInviteStore ────────────────────────────────────────────────────────────

func (su *storesInviteSuite) TestMemInviteStore_RemainingMethods() {
	t := su.T()
	s := NewInviteStore(Defaults())
	defer s.Close()

	assert.True(t, s.IsEnabled())

	s.SetWithGuild(gA, newTestInvite("c1"))
	assert.Equal(t, 1, s.Size())

	s.SweepExpired(0)
	assert.EqualValues(t, 0, s.Bytes())
	s.TrimTo(0)
	assert.Equal(t, 0, s.Size())
}

func (su *storesInviteSuite) TestMemInviteStore_NoopGuards() {
	t := su.T()
	s := NewInviteStore(Defaults())
	defer s.Close()

	// SetWithGuild no-ops
	s.SetWithGuild(gA, nil)
	s.SetWithGuild(gA, &discord.Invite{Code: ""})
	assert.Equal(t, 0, s.Size())

	// GetByGuild on an unknown guild returns empty.
	assert.Equal(t, 0, s.GetByGuild(discord.Snowflake(777)).Len())
	// DeleteGuild on an unknown guild is a no-op.
	s.DeleteGuild(discord.Snowflake(777))
	// Delete of an unknown code is a no-op.
	s.Delete("unknown")
}

// ── inviteGuildIndex edge cases ───────────────────────────────────────────────

func (su *storesInviteSuite) TestInviteGuildIndex_Edges() {
	t := su.T()
	idx := newInviteGuildIndex()

	// remove on empty index → no panic, no effect.
	idx.remove("missing")
	// removeGuild on empty guild → nil.
	assert.Nil(t, idx.removeGuild(gA))
	// forGuild on empty guild → nil.
	assert.Nil(t, idx.forGuild(gA))

	idx.add(gA, "x")
	idx.add(gA, "y")
	assert.Len(t, idx.forGuild(gA), 2)

	idx.remove("x")
	assert.Len(t, idx.forGuild(gA), 1)

	codes := idx.removeGuild(gA)
	assert.Equal(t, []string{"y"}, codes)
	assert.Nil(t, idx.forGuild(gA))
}

// ── MemBanStore ───────────────────────────────────────────────────────────────

func (su *storesBanSuite) TestMemBanStore_RemainingMethods() {
	t := su.T()
	s := NewBanStore(Defaults())
	defer s.Close()

	assert.True(t, s.IsEnabled())

	guildID, userID := discord.Snowflake(1), discord.Snowflake(2)
	s.Set(guildID, newTestBan(guildID, userID))

	assert.True(t, s.Has(guildID, userID))
	assert.False(t, s.Has(guildID, discord.Snowflake(999)))

	s.SweepExpired(0)
	assert.EqualValues(t, 0, s.Bytes())
	s.TrimTo(0)
	assert.Equal(t, 0, s.Size())
}

// ── MemAutoModerationRuleStore ────────────────────────────────────────────────

func (su *storesAutomodSuite) TestMemAutoModStore_RemainingMethods() {
	t := su.T()
	s := NewAutoModerationRuleStore(Defaults())
	defer s.Close()

	assert.True(t, s.IsEnabled())

	guildID, ruleID := discord.Snowflake(1), discord.Snowflake(2)
	s.Set(guildID, newTestAutoModRule(guildID, ruleID))

	s.SweepExpired(0)
	assert.EqualValues(t, 0, s.Bytes())
	s.TrimTo(0)
	assert.Equal(t, 0, s.Size())
}
