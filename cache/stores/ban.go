package stores

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// BanKey is a composite cache key for guild bans (same shape as MemberKey).
type BanKey = MemberKey

// MemBanStore is the in-memory implementation of cache.BanStore.
type MemBanStore struct {
	base  *BaseStore[BanKey, *discord.Ban]
	index *compositeGuildIndex
}

func NewBanStore(opts StoreOptions) *MemBanStore {
	return &MemBanStore{
		base:  NewBaseStore[BanKey, *discord.Ban](opts),
		index: newCompositeGuildIndex(),
	}
}

func (s *MemBanStore) Set(guildID discord.Snowflake, ban *discord.Ban) {
	if ban == nil {
		return
	}
	s.base.Set(BanKey{GuildID: guildID, UserID: ban.User.ID}, ban)
	s.index.add(guildID, ban.User.ID)
}

func (s *MemBanStore) Get(guildID, userID discord.Snowflake) (*discord.Ban, bool) {
	return s.base.Get(BanKey{GuildID: guildID, UserID: userID})
}

func (s *MemBanStore) Delete(guildID, userID discord.Snowflake) {
	s.base.Delete(BanKey{GuildID: guildID, UserID: userID})
	s.index.remove(guildID, userID)
}

func (s *MemBanStore) Has(guildID, userID discord.Snowflake) bool {
	return s.base.Has(BanKey{GuildID: guildID, UserID: userID})
}

// DeleteGuild removes all bans for guildID. Call on GUILD_DELETE.
func (s *MemBanStore) DeleteGuild(guildID discord.Snowflake) {
	userIDs := s.index.removeGuild(guildID)
	for _, uid := range userIDs {
		s.base.Delete(BanKey{GuildID: guildID, UserID: uid})
	}
}

// AllInGuild returns a snapshot Collection of all bans for guildID, keyed by user ID.
func (s *MemBanStore) AllInGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Ban] {
	userIDs := s.index.forGuild(guildID)
	out := collection.NewWithCapacity[discord.Snowflake, *discord.Ban](len(userIDs))
	for _, uid := range userIDs {
		if b, ok := s.base.Get(BanKey{GuildID: guildID, UserID: uid}); ok {
			out.Set(uid, b)
		}
	}
	return out
}

func (s *MemBanStore) Size() int { return s.base.Len() }

func (s *MemBanStore) SweepExpired(unusedWindow time.Duration) {
	s.base.SweepExpired(unusedWindow)
}

func (s *MemBanStore) Bytes() int64 { return s.base.Bytes() }

func (s *MemBanStore) TrimTo(n int) { s.base.TrimTo(n) }

func (s *MemBanStore) IsEnabled() bool { return s.base.IsEnabled() }

func (s *MemBanStore) Close() { s.base.Close() }
