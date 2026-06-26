package stores

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// MemAutoModerationRuleStore is the in-memory implementation of cache.AutoModerationRuleStore.
type MemAutoModerationRuleStore struct {
	base  *BaseStore[discord.Snowflake, *discord.AutoModerationRule]
	index *guildIndex
}

func NewAutoModerationRuleStore(opts StoreOptions) *MemAutoModerationRuleStore {
	s := &MemAutoModerationRuleStore{
		base:  NewBaseStore[discord.Snowflake, *discord.AutoModerationRule](opts),
		index: newGuildIndex(),
	}
	// Keep the guild index in sync when TTL sweeps or LRU/TrimTo evictions remove
	// entries straight from the base store.
	s.base.onEvict = func(ruleID discord.Snowflake) { s.index.removeEntity(ruleID) }
	return s
}

func (s *MemAutoModerationRuleStore) Set(guildID discord.Snowflake, rule *discord.AutoModerationRule) {
	if rule == nil {
		return
	}
	s.base.Set(rule.ID, rule)
	s.index.add(guildID, rule.ID)
}

func (s *MemAutoModerationRuleStore) Get(ruleID discord.Snowflake) (*discord.AutoModerationRule, bool) {
	return s.base.Get(ruleID)
}

// GetByGuild returns a snapshot Collection of all automod rules for guildID, keyed by rule ID.
func (s *MemAutoModerationRuleStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, discord.AutoModerationRule] {
	ruleIDs := s.index.forGuild(guildID)
	out := collection.NewWithCapacity[discord.Snowflake, discord.AutoModerationRule](len(ruleIDs))
	for _, id := range ruleIDs {
		if r, ok := s.base.Get(id); ok {
			out.Set(id, *r)
		}
	}
	return out
}

func (s *MemAutoModerationRuleStore) Delete(ruleID discord.Snowflake) {
	s.index.removeEntity(ruleID)
	s.base.Delete(ruleID)
}

// DeleteGuild removes all automod rules for guildID. Call on GUILD_DELETE.
func (s *MemAutoModerationRuleStore) DeleteGuild(guildID discord.Snowflake) {
	ruleIDs := s.index.removeGuild(guildID)
	for _, id := range ruleIDs {
		s.base.Delete(id)
	}
}

func (s *MemAutoModerationRuleStore) Size() int { return s.base.Len() }

func (s *MemAutoModerationRuleStore) SweepExpired(unusedWindow time.Duration) {
	s.base.SweepExpired(unusedWindow)
}

func (s *MemAutoModerationRuleStore) Bytes() int64 { return s.base.Bytes() }

func (s *MemAutoModerationRuleStore) TrimTo(n int) { s.base.TrimTo(n) }

func (s *MemAutoModerationRuleStore) IsEnabled() bool { return s.base.IsEnabled() }

func (s *MemAutoModerationRuleStore) Close() { s.base.Close() }
