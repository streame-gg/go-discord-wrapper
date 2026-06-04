package rediscache_test

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// These tests cover the Ban, AutoModerationRule and Invite stores added so the
// Redis backend satisfies the full cache.Cache interface. The happy-path tests
// mirror the equivalent mongocache suite for behavioural parity; the stale-index
// tests exercise the corrupt-document / reverse-map prune branches the way
// TestCorruptAndStaleBranches does for the original stores.

// ── Ban store (composite key: guild → user) ───────────────────────────────────

func (s *RedisCacheTestSuite) TestBanStore() {
	c := s.newCache(cache.Options{})
	b := c.Bans()

	b.Set(gA, nil) // no-op
	b.Set(gA, &discord.Ban{User: discord.User{ID: 1}})
	b.Set(gA, &discord.Ban{User: discord.User{ID: 2}})
	b.Set(gB, &discord.Ban{User: discord.User{ID: 3}})

	got, ok := b.Get(gA, 1)
	s.Require().True(ok)
	s.Equal(discord.Snowflake(1), got.User.ID)

	_, ok = b.Get(gA, 999)
	s.False(ok)

	s.Equal(2, b.AllInGuild(gA).Len())
	s.Equal(3, b.Size())

	b.Delete(gA, 1)
	s.Equal(1, b.AllInGuild(gA).Len())

	b.DeleteGuild(gA)
	s.Equal(0, b.AllInGuild(gA).Len())
	s.Equal(1, b.Size())
}

// ── AutoModeration rule store (guild-indexed by rule ID) ──────────────────────

func (s *RedisCacheTestSuite) TestAutoModRuleStore() {
	c := s.newCache(cache.Options{})
	a := c.AutoModRules()

	a.Set(gA, nil) // no-op
	a.Set(gA, &discord.AutoModerationRule{ID: 1, GuildID: gA, Name: "r1"})
	a.Set(gA, &discord.AutoModerationRule{ID: 2, GuildID: gA})
	a.Set(gB, &discord.AutoModerationRule{ID: 3, GuildID: gB})

	got, ok := a.Get(1)
	s.Require().True(ok)
	s.Equal("r1", got.Name)

	_, ok = a.Get(999)
	s.False(ok)

	s.Equal(2, a.GetByGuild(gA).Len())
	s.Equal(3, a.Size())

	a.Delete(1)
	s.Equal(1, a.GetByGuild(gA).Len())

	a.DeleteGuild(gA)
	s.Equal(0, a.GetByGuild(gA).Len())
	s.Equal(1, a.Size())
}

// ── Invite store (guild-indexed by string code) ──────────────────────────────

func (s *RedisCacheTestSuite) TestInviteStore() {
	c := s.newCache(cache.Options{})
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
	s.Require().True(ok)
	s.Equal("abc", got.Code)

	got, ok = i.Get("guildless")
	s.Require().True(ok)
	s.Equal("guildless", got.Code)

	_, ok = i.Get("missing")
	s.False(ok)

	s.Equal(2, i.GetByGuild(gA).Len())
	s.Equal(1, i.GetByGuild(gB).Len())
	s.Equal(3, i.Size())

	i.Delete("abc")
	s.Equal(1, i.GetByGuild(gA).Len())

	i.DeleteGuild(gA)
	s.Equal(0, i.GetByGuild(gA).Len())
	_, ok = i.Get("def")
	s.False(ok, "DeleteGuild must remove the guild's invites")
}

// ── Stale-index / corrupt-document prune branches ─────────────────────────────

// TestExtraStores_StaleAndCorruptBranches seeds corrupt values and dangling
// index entries directly into the keyspace so the prune paths in Get/AllInGuild/
// GetByGuild are exercised: SRem of a stale composite key, the reverse-map SRem +
// Del on a guild-indexed Get miss, and the v==nil / bad-json skips in the bulk
// MGet loops.
func (s *RedisCacheTestSuite) TestExtraStores_StaleAndCorruptBranches() {
	c := s.newCache(cache.Options{})
	raw, p := s.rawClient()
	ctx := context.Background()
	bad := "{ not valid json"
	g := gA.String()

	// ── Ban: corrupt value at "1", stale index members "2" (Get prune) and "3"
	// (AllInGuild prune). "2"/"3" have no backing value.
	s.Require().NoError(raw.Set(ctx, p+":ban:"+g+":1", bad, 0).Err())
	s.Require().NoError(raw.SAdd(ctx, p+":ban:guild:"+g, "1", "2", "3").Err())

	// AutoMod & Invite: corrupt value + stale members + reverse-map entries so the
	// guild-indexed Get prune path can SRem them.
	s.Require().NoError(raw.Set(ctx, p+":automod:1", bad, 0).Err())
	s.Require().NoError(raw.SAdd(ctx, p+":automod:guild:"+g, "1", "2", "3").Err())
	s.Require().NoError(raw.Set(ctx, p+":automod:map:2", g, 0).Err())
	s.Require().NoError(raw.Set(ctx, p+":automod:map:3", g, 0).Err())

	s.Require().NoError(raw.Set(ctx, p+":invite:c1", bad, 0).Err())
	s.Require().NoError(raw.SAdd(ctx, p+":invite:guild:"+g, "c1", "c2", "c3").Err())
	s.Require().NoError(raw.Set(ctx, p+":invite:map:c2", g, 0).Err())
	s.Require().NoError(raw.Set(ctx, p+":invite:map:c3", g, 0).Err())

	// Get on the corrupt id → json.Unmarshal fails → miss.
	_, ok := c.Bans().Get(gA, 1)
	s.False(ok)
	_, ok = c.AutoModRules().Get(1)
	s.False(ok)
	_, ok = c.Invites().Get("c1")
	s.False(ok)

	// Get on a stale id (in the index, no backing value) → prune branch.
	// Ban uses the composite-key SRem; automod/invite use the reverse-map SRem.
	_, ok = c.Bans().Get(gA, 2)
	s.False(ok)
	_, ok = c.AutoModRules().Get(2)
	s.False(ok)
	_, ok = c.Invites().Get("c2")
	s.False(ok)

	// Bulk reads → corrupt rows skipped, stale rows pruned → empty result.
	s.Equal(0, c.Bans().AllInGuild(gA).Len())
	s.Equal(0, c.AutoModRules().GetByGuild(gA).Len())
	s.Equal(0, c.Invites().GetByGuild(gA).Len())
}

// TestBanStore_AllInGuildEdgeIDs covers AllInGuild's User.ID fallback: a stored
// ban whose JSON unmarshals but carries an empty User.ID is recovered from the
// numeric index member, while a non-numeric index member with an empty-ID ban is
// treated as stale and pruned.
func (s *RedisCacheTestSuite) TestBanStore_AllInGuildEdgeIDs() {
	c := s.newCache(cache.Options{})
	raw, p := s.rawClient()
	ctx := context.Background()
	g := gA.String()
	empty := "{}" // valid JSON, but Ban.User.ID is the zero Snowflake.

	// Numeric index member "5" → empty-ID ban recovered as Snowflake(5).
	s.Require().NoError(raw.Set(ctx, p+":ban:"+g+":5", empty, 0).Err())
	// Non-numeric index member → empty ID, ParseUint fails → stale, pruned.
	s.Require().NoError(raw.Set(ctx, p+":ban:"+g+":notanum", empty, 0).Err())
	s.Require().NoError(raw.SAdd(ctx, p+":ban:guild:"+g, "5", "notanum").Err())

	all := c.Bans().AllInGuild(gA)
	s.Require().Equal(1, all.Len(), "the numeric-id ban is recovered, the non-numeric one is pruned")
	_, ok := all.Get(discord.Snowflake(5))
	s.True(ok, "ban id recovered from the numeric index member")

	// The stale non-numeric member was SRem'd from the index.
	s.Equal(int64(1), raw.SCard(ctx, p+":ban:guild:"+g).Val())
}
