package stores

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// MemStageInstanceStore is the in-memory implementation of cache.StageInstanceStore.
type MemStageInstanceStore struct {
	base  *BaseStore[discord.Snowflake, *discord.StageInstance]
	index *guildIndex
}

// NewStageInstanceStore creates a MemStageInstanceStore with the given options.
func NewStageInstanceStore(opts StoreOptions) *MemStageInstanceStore {
	return &MemStageInstanceStore{
		base:  NewBaseStore[discord.Snowflake, *discord.StageInstance](opts),
		index: newGuildIndex(),
	}
}

func (s *MemStageInstanceStore) Set(instance *discord.StageInstance) {
	if instance == nil {
		return
	}
	s.index.add(instance.GuildID, instance.ID)
	s.base.Set(instance.ID, instance)
}

func (s *MemStageInstanceStore) Get(instanceID discord.Snowflake) (*discord.StageInstance, bool) {
	return s.base.Get(instanceID)
}

// GetByGuild returns a snapshot Collection of all stage instances for guildID.
func (s *MemStageInstanceStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.StageInstance] {
	ids := s.index.forGuild(guildID)
	out := collection.NewWithCapacity[discord.Snowflake, *discord.StageInstance](len(ids))
	for _, id := range ids {
		if si, ok := s.base.Get(id); ok {
			out.Set(id, si)
		}
	}
	return out
}

func (s *MemStageInstanceStore) Delete(instanceID discord.Snowflake) {
	s.index.removeEntity(instanceID)
	s.base.Delete(instanceID)
}

// DeleteGuild removes all stage instances for guildID. Call on GUILD_DELETE.
func (s *MemStageInstanceStore) DeleteGuild(guildID discord.Snowflake) {
	ids := s.index.removeGuild(guildID)
	for _, id := range ids {
		s.base.Delete(id)
	}
}

func (s *MemStageInstanceStore) Size() int { return s.base.Len() }

func (s *MemStageInstanceStore) SweepExpired(unusedWindow time.Duration) {
	s.base.SweepExpired(unusedWindow)
}

func (s *MemStageInstanceStore) Bytes() int64 { return s.base.Bytes() }

func (s *MemStageInstanceStore) TrimTo(n int) { s.base.TrimTo(n) }

func (s *MemStageInstanceStore) IsEnabled() bool { return s.base.IsEnabled() }

func (s *MemStageInstanceStore) Close() { s.base.Close() }
