package stores

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// MemUserStore is the in-memory implementation of cache.UserStore.
type MemUserStore struct {
	base *BaseStore[discord.Snowflake, *discord.User]
}

// NewUserStore creates a MemUserStore with the given options.
func NewUserStore(opts StoreOptions) *MemUserStore {
	return &MemUserStore{base: NewBaseStore[discord.Snowflake, *discord.User](opts)}
}

func (s *MemUserStore) Set(u *discord.User) {
	if u == nil {
		return
	}
	s.base.Set(u.ID, u)
}

func (s *MemUserStore) Get(id discord.Snowflake) (*discord.User, bool) {
	return s.base.Get(id)
}

func (s *MemUserStore) Delete(id discord.Snowflake) { s.base.Delete(id) }

func (s *MemUserStore) Has(id discord.Snowflake) bool { return s.base.Has(id) }

func (s *MemUserStore) Size() int { return s.base.Len() }

func (s *MemUserStore) SweepExpired(unusedWindow time.Duration) {
	s.base.SweepExpired(unusedWindow)
}

func (s *MemUserStore) All() *collection.Collection[discord.Snowflake, *discord.User] {
	return s.base.All()
}

func (s *MemUserStore) Clear() { s.base.Clear() }

func (s *MemUserStore) Bytes() int64 { return s.base.Bytes() }

func (s *MemUserStore) TrimTo(n int) { s.base.TrimTo(n) }

func (s *MemUserStore) IsEnabled() bool { return s.base.IsEnabled() }

func (s *MemUserStore) Close() { s.base.Close() }
