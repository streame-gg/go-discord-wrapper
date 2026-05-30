package stores

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// MemMemberStore is the in-memory implementation of cache.MemberStore.
type MemMemberStore struct {
	base  *BaseStore[MemberKey, *discord.GuildMember]
	index *compositeGuildIndex
}

// NewMemberStore creates a MemMemberStore with the given options.
func NewMemberStore(opts StoreOptions) *MemMemberStore {
	s := &MemMemberStore{
		base:  NewBaseStore[MemberKey, *discord.GuildMember](opts),
		index: newCompositeGuildIndex(),
	}
	s.base.onEvict = func(k MemberKey) { s.index.remove(k.GuildID, k.UserID) }
	return s
}

func (s *MemMemberStore) Set(guildID discord.Snowflake, member *discord.GuildMember) {
	if member == nil || member.User == nil {
		return
	}
	s.base.Set(MemberKey{GuildID: guildID, UserID: member.User.ID}, member)
	s.index.add(guildID, member.User.ID)
}

func (s *MemMemberStore) Get(guildID, userID discord.Snowflake) (*discord.GuildMember, bool) {
	return s.base.Get(MemberKey{GuildID: guildID, UserID: userID})
}

func (s *MemMemberStore) Delete(guildID, userID discord.Snowflake) {
	s.base.Delete(MemberKey{GuildID: guildID, UserID: userID})
	s.index.remove(guildID, userID)
}

func (s *MemMemberStore) Has(guildID, userID discord.Snowflake) bool {
	return s.base.Has(MemberKey{GuildID: guildID, UserID: userID})
}

// DeleteGuild removes all members for guildID. Call on GUILD_DELETE.
func (s *MemMemberStore) DeleteGuild(guildID discord.Snowflake) {
	userIDs := s.index.removeGuild(guildID)
	for _, uid := range userIDs {
		s.base.Delete(MemberKey{GuildID: guildID, UserID: uid})
	}
}

// AllInGuild returns a snapshot Collection of all members for guildID, keyed by user ID.
func (s *MemMemberStore) AllInGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.GuildMember] {
	userIDs := s.index.forGuild(guildID)
	out := collection.NewWithCapacity[discord.Snowflake, *discord.GuildMember](len(userIDs))
	for _, uid := range userIDs {
		if m, ok := s.base.Get(MemberKey{GuildID: guildID, UserID: uid}); ok {
			out.Set(uid, m)
		}
	}
	return out
}

func (s *MemMemberStore) Size() int { return s.base.Len() }

func (s *MemMemberStore) SweepExpired(unusedWindow time.Duration) {
	s.base.SweepExpired(unusedWindow)
}

func (s *MemMemberStore) Bytes() int64 { return s.base.Bytes() }

func (s *MemMemberStore) TrimTo(n int) { s.base.TrimTo(n) }

func (s *MemMemberStore) IsEnabled() bool { return s.base.IsEnabled() }

func (s *MemMemberStore) Close() { s.base.Close() }
