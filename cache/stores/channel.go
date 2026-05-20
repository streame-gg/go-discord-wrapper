package stores

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// MemChannelStore is the in-memory implementation of cache.ChannelStore.
type MemChannelStore struct {
	base *BaseStore[discord.Snowflake, *discord.Channel]
}

// NewChannelStore creates a MemChannelStore with the given options.
func NewChannelStore(opts StoreOptions) *MemChannelStore {
	return &MemChannelStore{base: NewBaseStore[discord.Snowflake, *discord.Channel](opts)}
}

func (s *MemChannelStore) Set(ch *discord.Channel) {
	if ch == nil {
		return
	}
	s.base.Set(ch.ID, ch)
}

func (s *MemChannelStore) Get(id discord.Snowflake) (*discord.Channel, bool) {
	return s.base.Get(id)
}

func (s *MemChannelStore) Delete(id discord.Snowflake) { s.base.Delete(id) }

func (s *MemChannelStore) Has(id discord.Snowflake) bool { return s.base.Has(id) }

func (s *MemChannelStore) Size() int { return s.base.Len() }

func (s *MemChannelStore) SweepExpired(unusedWindow time.Duration) {
	s.base.SweepExpired(unusedWindow)
}

func (s *MemChannelStore) All() *collection.Collection[discord.Snowflake, *discord.Channel] {
	return s.base.All()
}

func (s *MemChannelStore) Clear() { s.base.Clear() }

func (s *MemChannelStore) Bytes() int64 { return s.base.Bytes() }

func (s *MemChannelStore) TrimTo(n int) { s.base.TrimTo(n) }

func (s *MemChannelStore) IsEnabled() bool { return s.base.IsEnabled() }

func (s *MemChannelStore) Close() { s.base.Close() }
