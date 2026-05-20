package stores

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// MemGuildStore is the in-memory implementation of cache.GuildStore.
type MemGuildStore struct {
	base *BaseStore[discord.Snowflake, *discord.Guild]
}

// NewGuildStore creates a MemGuildStore with the given options.
func NewGuildStore(opts StoreOptions) *MemGuildStore {
	return &MemGuildStore{base: NewBaseStore[discord.Snowflake, *discord.Guild](opts)}
}

func (s *MemGuildStore) Set(g *discord.Guild) {
	if g == nil {
		return
	}
	s.base.Set(g.ID, g)
}

func (s *MemGuildStore) Get(id discord.Snowflake) (*discord.Guild, bool) {
	return s.base.Get(id)
}

func (s *MemGuildStore) Delete(id discord.Snowflake) { s.base.Delete(id) }

func (s *MemGuildStore) Has(id discord.Snowflake) bool { return s.base.Has(id) }

func (s *MemGuildStore) Size() int { return s.base.Len() }

func (s *MemGuildStore) SweepExpired(unusedWindow time.Duration) {
	s.base.SweepExpired(unusedWindow)
}

func (s *MemGuildStore) All() *collection.Collection[discord.Snowflake, *discord.Guild] {
	return s.base.All()
}

func (s *MemGuildStore) Clear() { s.base.Clear() }

func (s *MemGuildStore) Bytes() int64 { return s.base.Bytes() }

func (s *MemGuildStore) TrimTo(n int) { s.base.TrimTo(n) }

func (s *MemGuildStore) IsEnabled() bool { return s.base.IsEnabled() }

func (s *MemGuildStore) Close() { s.base.Close() }
