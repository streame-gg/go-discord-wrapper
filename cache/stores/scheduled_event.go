package stores

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// MemScheduledEventStore is the in-memory implementation of cache.ScheduledEventStore.
type MemScheduledEventStore struct {
	base  *BaseStore[discord.Snowflake, *discord.GuildScheduledEvent]
	index *guildIndex
}

// NewScheduledEventStore creates a MemScheduledEventStore with the given options.
func NewScheduledEventStore(opts StoreOptions) *MemScheduledEventStore {
	return &MemScheduledEventStore{
		base:  NewBaseStore[discord.Snowflake, *discord.GuildScheduledEvent](opts),
		index: newGuildIndex(),
	}
}

func (s *MemScheduledEventStore) Set(event *discord.GuildScheduledEvent) {
	if event == nil {
		return
	}
	s.index.add(event.GuildID, event.ID)
	s.base.Set(event.ID, event)
}

func (s *MemScheduledEventStore) Get(eventID discord.Snowflake) (*discord.GuildScheduledEvent, bool) {
	return s.base.Get(eventID)
}

// GetByGuild returns a snapshot Collection of all scheduled events for guildID.
func (s *MemScheduledEventStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.GuildScheduledEvent] {
	ids := s.index.forGuild(guildID)
	out := collection.NewWithCapacity[discord.Snowflake, *discord.GuildScheduledEvent](len(ids))
	for _, id := range ids {
		if ev, ok := s.base.Get(id); ok {
			out.Set(id, ev)
		}
	}
	return out
}

func (s *MemScheduledEventStore) Delete(eventID discord.Snowflake) {
	s.index.removeEntity(eventID)
	s.base.Delete(eventID)
}

// DeleteGuild removes all scheduled events for guildID. Call on GUILD_DELETE.
func (s *MemScheduledEventStore) DeleteGuild(guildID discord.Snowflake) {
	ids := s.index.removeGuild(guildID)
	for _, id := range ids {
		s.base.Delete(id)
	}
}

func (s *MemScheduledEventStore) Size() int { return s.base.Len() }

func (s *MemScheduledEventStore) SweepExpired(unusedWindow time.Duration) {
	s.base.SweepExpired(unusedWindow)
}

func (s *MemScheduledEventStore) Bytes() int64 { return s.base.Bytes() }

func (s *MemScheduledEventStore) TrimTo(n int) { s.base.TrimTo(n) }

func (s *MemScheduledEventStore) IsEnabled() bool { return s.base.IsEnabled() }

func (s *MemScheduledEventStore) Close() { s.base.Close() }
