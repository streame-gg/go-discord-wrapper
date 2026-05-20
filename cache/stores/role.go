package stores

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// MemRoleStore is the in-memory implementation of cache.RoleStore.
type MemRoleStore struct {
	base  *BaseStore[discord.Snowflake, *discord.Role]
	index *guildIndex
}

// NewRoleStore creates a MemRoleStore with the given options.
func NewRoleStore(opts StoreOptions) *MemRoleStore {
	return &MemRoleStore{
		base:  NewBaseStore[discord.Snowflake, *discord.Role](opts),
		index: newGuildIndex(),
	}
}

func (s *MemRoleStore) Set(guildID discord.Snowflake, role *discord.Role) {
	if role == nil {
		return
	}
	s.index.add(guildID, role.ID)
	s.base.Set(role.ID, role)
}

func (s *MemRoleStore) Get(roleID discord.Snowflake) (*discord.Role, bool) {
	return s.base.Get(roleID)
}

// GetByGuild returns a snapshot Collection of all roles for guildID, keyed by role ID.
func (s *MemRoleStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Role] {
	ids := s.index.forGuild(guildID)
	out := collection.NewWithCapacity[discord.Snowflake, *discord.Role](len(ids))
	for _, id := range ids {
		if r, ok := s.base.Get(id); ok {
			out.Set(id, r)
		}
	}
	return out
}

func (s *MemRoleStore) Delete(roleID discord.Snowflake) {
	s.index.removeEntity(roleID)
	s.base.Delete(roleID)
}

// DeleteGuild removes all roles for guildID. Call on GUILD_DELETE.
func (s *MemRoleStore) DeleteGuild(guildID discord.Snowflake) {
	ids := s.index.removeGuild(guildID)
	for _, id := range ids {
		s.base.Delete(id)
	}
}

// SetAll atomically replaces all roles for guildID.
// Uses BaseStore.ReplaceKeys so no reader sees a half-replaced state.
// Fixes the non-atomic DeleteGuild+loop-Set pattern that was in the old implementation.
func (s *MemRoleStore) SetAll(guildID discord.Snowflake, roles []*discord.Role) {
	newIDs := make([]discord.Snowflake, 0, len(roles))
	add := make(map[discord.Snowflake]*discord.Role, len(roles))
	for _, r := range roles {
		if r == nil {
			continue
		}
		newIDs = append(newIDs, r.ID)
		add[r.ID] = r
	}
	toDelete := s.index.setForGuild(guildID, newIDs)
	s.base.ReplaceKeys(toDelete, add)
}

// All returns a snapshot Collection of all roles across all guilds, keyed by role ID.
func (s *MemRoleStore) All() *collection.Collection[discord.Snowflake, *discord.Role] {
	return s.base.All()
}

func (s *MemRoleStore) Size() int { return s.base.Len() }

func (s *MemRoleStore) SweepExpired(unusedWindow time.Duration) {
	s.base.SweepExpired(unusedWindow)
}

func (s *MemRoleStore) Bytes() int64 { return s.base.Bytes() }

func (s *MemRoleStore) TrimTo(n int) { s.base.TrimTo(n) }

func (s *MemRoleStore) IsEnabled() bool { return s.base.IsEnabled() }

func (s *MemRoleStore) Close() { s.base.Close() }
