package stores

import (
	"sync"
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// inviteGuildIndex maintains a guildID → invite code mapping for guild-scoped lookups.
type inviteGuildIndex struct {
	mu      sync.RWMutex
	byGuild map[discord.Snowflake]map[string]struct{} // guildID → codes
	toGuild map[string]discord.Snowflake              // code → guildID
}

func newInviteGuildIndex() *inviteGuildIndex {
	return &inviteGuildIndex{
		byGuild: make(map[discord.Snowflake]map[string]struct{}),
		toGuild: make(map[string]discord.Snowflake),
	}
}

func (idx *inviteGuildIndex) add(guildID discord.Snowflake, code string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	if idx.byGuild[guildID] == nil {
		idx.byGuild[guildID] = make(map[string]struct{})
	}
	idx.byGuild[guildID][code] = struct{}{}
	idx.toGuild[code] = guildID
}

func (idx *inviteGuildIndex) remove(code string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	guildID, ok := idx.toGuild[code]
	if !ok {
		return
	}
	delete(idx.toGuild, code)
	if set := idx.byGuild[guildID]; set != nil {
		delete(set, code)
		if len(set) == 0 {
			delete(idx.byGuild, guildID)
		}
	}
}

func (idx *inviteGuildIndex) removeGuild(guildID discord.Snowflake) []string {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	set := idx.byGuild[guildID]
	if len(set) == 0 {
		delete(idx.byGuild, guildID)
		return nil
	}
	codes := make([]string, 0, len(set))
	for code := range set {
		codes = append(codes, code)
		delete(idx.toGuild, code)
	}
	delete(idx.byGuild, guildID)
	return codes
}

func (idx *inviteGuildIndex) forGuild(guildID discord.Snowflake) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	set := idx.byGuild[guildID]
	if len(set) == 0 {
		return nil
	}
	codes := make([]string, 0, len(set))
	for code := range set {
		codes = append(codes, code)
	}
	return codes
}

// MemInviteStore is the in-memory implementation of cache.InviteStore.
type MemInviteStore struct {
	base  *BaseStore[string, *discord.Invite]
	index *inviteGuildIndex
}

func NewInviteStore(opts StoreOptions) *MemInviteStore {
	return &MemInviteStore{
		base:  NewBaseStore[string, *discord.Invite](opts),
		index: newInviteGuildIndex(),
	}
}

func (s *MemInviteStore) Set(invite *discord.Invite) {
	if invite == nil || invite.Code == "" {
		return
	}
	s.base.Set(invite.Code, invite)
	if invite.Guild != nil {
		s.index.add(invite.Guild.ID, invite.Code)
	}
}

func (s *MemInviteStore) SetWithGuild(guildID discord.Snowflake, invite *discord.Invite) {
	if invite == nil || invite.Code == "" {
		return
	}
	s.base.Set(invite.Code, invite)
	s.index.add(guildID, invite.Code)
}

func (s *MemInviteStore) Get(code string) (*discord.Invite, bool) {
	return s.base.Get(code)
}

// GetByGuild returns a snapshot Collection of all invites for guildID, keyed by code.
func (s *MemInviteStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[string, *discord.Invite] {
	codes := s.index.forGuild(guildID)
	out := collection.NewWithCapacity[string, *discord.Invite](len(codes))
	for _, code := range codes {
		if inv, ok := s.base.Get(code); ok {
			out.Set(code, inv)
		}
	}
	return out
}

func (s *MemInviteStore) Delete(code string) {
	s.index.remove(code)
	s.base.Delete(code)
}

// DeleteGuild removes all invites for guildID. Call on GUILD_DELETE.
func (s *MemInviteStore) DeleteGuild(guildID discord.Snowflake) {
	codes := s.index.removeGuild(guildID)
	for _, code := range codes {
		s.base.Delete(code)
	}
}

func (s *MemInviteStore) Size() int { return s.base.Len() }

func (s *MemInviteStore) SweepExpired(unusedWindow time.Duration) {
	s.base.SweepExpired(unusedWindow)
}

func (s *MemInviteStore) Bytes() int64 { return s.base.Bytes() }

func (s *MemInviteStore) TrimTo(n int) { s.base.TrimTo(n) }

func (s *MemInviteStore) IsEnabled() bool { return s.base.IsEnabled() }

func (s *MemInviteStore) Close() { s.base.Close() }
