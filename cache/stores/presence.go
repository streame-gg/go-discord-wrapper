package stores

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// MemPresenceStore is the in-memory implementation of cache.PresenceStore.
// Disabled by default — enable with cache.WithPresenceStore(cache.StoreOptions{}).
type MemPresenceStore struct {
	base  *BaseStore[PresenceKey, *discord.Presence]
	index *compositeGuildIndex
}

// NewPresenceStore creates a MemPresenceStore with the given options.
func NewPresenceStore(opts StoreOptions) *MemPresenceStore {
	s := &MemPresenceStore{
		base:  NewBaseStore[PresenceKey, *discord.Presence](opts),
		index: newCompositeGuildIndex(),
	}
	s.base.onEvict = func(k PresenceKey) { s.index.remove(k.GuildID, k.UserID) }
	return s
}

// Set stores the presence. GuildID and UserID are extracted from the presence.
func (s *MemPresenceStore) Set(presence *discord.Presence) {
	if presence == nil || presence.User.ID.IsEmpty() {
		return
	}
	key := PresenceKey{GuildID: presence.GuildID, UserID: presence.User.ID}
	s.base.Set(key, presence)
	s.index.add(presence.GuildID, presence.User.ID)
}

func (s *MemPresenceStore) Get(guildID, userID discord.Snowflake) (*discord.Presence, bool) {
	return s.base.Get(PresenceKey{GuildID: guildID, UserID: userID})
}

func (s *MemPresenceStore) Delete(guildID, userID discord.Snowflake) {
	s.base.Delete(PresenceKey{GuildID: guildID, UserID: userID})
	s.index.remove(guildID, userID)
}

// DeleteGuild removes all presences for guildID. Call on GUILD_DELETE.
func (s *MemPresenceStore) DeleteGuild(guildID discord.Snowflake) {
	userIDs := s.index.removeGuild(guildID)
	for _, uid := range userIDs {
		s.base.Delete(PresenceKey{GuildID: guildID, UserID: uid})
	}
}

// GetByGuild returns a snapshot Collection of all presences for guildID, keyed by user ID.
func (s *MemPresenceStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Presence] {
	userIDs := s.index.forGuild(guildID)
	out := collection.NewWithCapacity[discord.Snowflake, *discord.Presence](len(userIDs))
	for _, uid := range userIDs {
		if p, ok := s.base.Get(PresenceKey{GuildID: guildID, UserID: uid}); ok {
			out.Set(uid, p)
		}
	}
	return out
}

func (s *MemPresenceStore) Size() int { return s.base.Len() }

func (s *MemPresenceStore) SweepExpired(unusedWindow time.Duration) {
	s.base.SweepExpired(unusedWindow)
}

func (s *MemPresenceStore) Bytes() int64 { return s.base.Bytes() }

func (s *MemPresenceStore) TrimTo(n int) { s.base.TrimTo(n) }

func (s *MemPresenceStore) IsEnabled() bool { return s.base.IsEnabled() }

func (s *MemPresenceStore) Close() { s.base.Close() }
