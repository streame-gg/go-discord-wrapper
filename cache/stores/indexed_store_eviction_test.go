package stores

// Regression tests for the guild-indexed stores added later (Ban, AutoMod,
// Invite): like the composite-key stores, their secondary guild index must stay
// in sync when the base store drops entries via LRU/TrimTo eviction or via the
// lazy TTL deletion that happens on Get. Without the onEvict hook (and the
// onEvict call in BaseStore.Get) these indexes leaked phantom entries.

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Ban store ─────────────────────────────────────────────────────────────────

func (su *storesEvictionSuite) TestBanStore_LRUEvictionKeepsIndexInSync() {
	guildID := discord.Snowflake(1000)
	s := NewBanStore(StoreOptions{MaxItems: 1, EvictionStrategy: EvictLRU})
	defer s.Close()

	s.Set(guildID, newTestBan(guildID, 1))
	time.Sleep(time.Millisecond)
	s.Set(guildID, newTestBan(guildID, 2)) // evicts user 1

	all := s.AllInGuild(guildID)
	su.Equal(1, all.Len(), "evicted ban must not linger in the guild index")
	_, ok := all.Get(1)
	su.False(ok, "user 1 was evicted; index must not resolve it")
	_, ok = all.Get(2)
	su.True(ok)
}

func (su *storesEvictionSuite) TestBanStore_LazyTTLEvictionKeepsIndexInSync() {
	guildID := discord.Snowflake(1000)
	s := NewBanStore(StoreOptions{TTL: 20 * time.Millisecond})
	defer s.Close()

	s.Set(guildID, newTestBan(guildID, 1))
	time.Sleep(40 * time.Millisecond)

	// Lazy expiry on Get must also clean the index (no background sweeper here).
	_, ok := s.Get(guildID, 1)
	su.False(ok, "entry should have expired")
	su.Equal(0, s.AllInGuild(guildID).Len(), "lazy-expired ban must be pruned from the index")
}

// ── AutoMod rule store ────────────────────────────────────────────────────────

func (su *storesEvictionSuite) TestAutoModStore_LRUEvictionKeepsIndexInSync() {
	guildID := discord.Snowflake(1000)
	s := NewAutoModerationRuleStore(StoreOptions{MaxItems: 1, EvictionStrategy: EvictLRU})
	defer s.Close()

	s.Set(guildID, &discord.AutoModerationRule{ID: 1, GuildID: guildID})
	time.Sleep(time.Millisecond)
	s.Set(guildID, &discord.AutoModerationRule{ID: 2, GuildID: guildID}) // evicts rule 1

	all := s.GetByGuild(guildID)
	su.Equal(1, all.Len(), "evicted rule must not linger in the guild index")
	_, ok := all.Get(1)
	su.False(ok)
}

func (su *storesEvictionSuite) TestAutoModStore_LazyTTLEvictionKeepsIndexInSync() {
	guildID := discord.Snowflake(1000)
	s := NewAutoModerationRuleStore(StoreOptions{TTL: 20 * time.Millisecond})
	defer s.Close()

	s.Set(guildID, &discord.AutoModerationRule{ID: 1, GuildID: guildID})
	time.Sleep(40 * time.Millisecond)

	_, ok := s.Get(1)
	su.False(ok, "rule should have expired")
	su.Equal(0, s.GetByGuild(guildID).Len(), "lazy-expired rule must be pruned from the index")
}

// ── Invite store ──────────────────────────────────────────────────────────────

func (su *storesEvictionSuite) TestInviteStore_LRUEvictionKeepsIndexInSync() {
	guildID := discord.Snowflake(1000)
	s := NewInviteStore(StoreOptions{MaxItems: 1, EvictionStrategy: EvictLRU})
	defer s.Close()

	s.Set(&discord.Invite{Code: "aaa", Guild: &discord.Guild{ID: guildID}})
	time.Sleep(time.Millisecond)
	s.Set(&discord.Invite{Code: "bbb", Guild: &discord.Guild{ID: guildID}}) // evicts "aaa"

	all := s.GetByGuild(guildID)
	su.Equal(1, all.Len(), "evicted invite must not linger in the guild index")
	_, ok := all.Get("aaa")
	su.False(ok)
}

func (su *storesEvictionSuite) TestInviteStore_LazyTTLEvictionKeepsIndexInSync() {
	guildID := discord.Snowflake(1000)
	s := NewInviteStore(StoreOptions{TTL: 20 * time.Millisecond})
	defer s.Close()

	s.Set(&discord.Invite{Code: "aaa", Guild: &discord.Guild{ID: guildID}})
	time.Sleep(40 * time.Millisecond)

	_, ok := s.Get("aaa")
	su.False(ok, "invite should have expired")
	su.Equal(0, s.GetByGuild(guildID).Len(), "lazy-expired invite must be pruned from the index")
}
