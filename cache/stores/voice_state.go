package stores

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// MemVoiceStateStore is the in-memory implementation of cache.VoiceStateStore.
type MemVoiceStateStore struct {
	base  *BaseStore[VoiceStateKey, *discord.VoiceState]
	index *compositeGuildIndex
}

// NewVoiceStateStore creates a MemVoiceStateStore with the given options.
func NewVoiceStateStore(opts StoreOptions) *MemVoiceStateStore {
	return &MemVoiceStateStore{
		base:  NewBaseStore[VoiceStateKey, *discord.VoiceState](opts),
		index: newCompositeGuildIndex(),
	}
}

// Set stores the voice state. guildID is passed explicitly because
// discord.VoiceState.GuildID is a pointer and may be nil on some gateway events.
func (s *MemVoiceStateStore) Set(guildID discord.Snowflake, state *discord.VoiceState) {
	if state == nil {
		return
	}
	s.base.Set(VoiceStateKey{GuildID: guildID, UserID: state.UserID}, state)
	s.index.add(guildID, state.UserID)
}

func (s *MemVoiceStateStore) Get(guildID, userID discord.Snowflake) (*discord.VoiceState, bool) {
	return s.base.Get(VoiceStateKey{GuildID: guildID, UserID: userID})
}

func (s *MemVoiceStateStore) Delete(guildID, userID discord.Snowflake) {
	s.base.Delete(VoiceStateKey{GuildID: guildID, UserID: userID})
	s.index.remove(guildID, userID)
}

func (s *MemVoiceStateStore) Has(guildID, userID discord.Snowflake) bool {
	return s.base.Has(VoiceStateKey{GuildID: guildID, UserID: userID})
}

// DeleteGuild removes all voice states for guildID. Call on GUILD_DELETE.
func (s *MemVoiceStateStore) DeleteGuild(guildID discord.Snowflake) {
	userIDs := s.index.removeGuild(guildID)
	for _, uid := range userIDs {
		s.base.Delete(VoiceStateKey{GuildID: guildID, UserID: uid})
	}
}

// AllInGuild returns a snapshot Collection of all voice states for guildID, keyed by user ID.
func (s *MemVoiceStateStore) AllInGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.VoiceState] {
	userIDs := s.index.forGuild(guildID)
	out := collection.NewWithCapacity[discord.Snowflake, *discord.VoiceState](len(userIDs))
	for _, uid := range userIDs {
		if vs, ok := s.base.Get(VoiceStateKey{GuildID: guildID, UserID: uid}); ok {
			out.Set(uid, vs)
		}
	}
	return out
}

func (s *MemVoiceStateStore) Size() int { return s.base.Len() }

func (s *MemVoiceStateStore) SweepExpired(unusedWindow time.Duration) {
	s.base.SweepExpired(unusedWindow)
}

func (s *MemVoiceStateStore) Bytes() int64 { return s.base.Bytes() }

func (s *MemVoiceStateStore) TrimTo(n int) { s.base.TrimTo(n) }

func (s *MemVoiceStateStore) IsEnabled() bool { return s.base.IsEnabled() }

func (s *MemVoiceStateStore) Close() { s.base.Close() }
