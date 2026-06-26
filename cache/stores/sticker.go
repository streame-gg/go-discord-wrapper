package stores

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// MemStickerStore is the in-memory implementation of cache.StickerStore.
type MemStickerStore struct {
	base  *BaseStore[discord.Snowflake, *discord.Sticker]
	index *guildIndex
}

// NewStickerStore creates a MemStickerStore with the given options.
func NewStickerStore(opts StoreOptions) *MemStickerStore {
	return &MemStickerStore{
		base:  NewBaseStore[discord.Snowflake, *discord.Sticker](opts),
		index: newGuildIndex(),
	}
}

func (s *MemStickerStore) Set(guildID discord.Snowflake, sticker *discord.Sticker) {
	if sticker == nil {
		return
	}
	s.index.add(guildID, sticker.ID)
	s.base.Set(sticker.ID, sticker)
}

func (s *MemStickerStore) Get(stickerID discord.Snowflake) (*discord.Sticker, bool) {
	return s.base.Get(stickerID)
}

// GetByGuild returns a snapshot Collection of all stickers for guildID, keyed by sticker ID.
func (s *MemStickerStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, discord.Sticker] {
	ids := s.index.forGuild(guildID)
	out := collection.NewWithCapacity[discord.Snowflake, discord.Sticker](len(ids))
	for _, id := range ids {
		if st, ok := s.base.Get(id); ok {
			out.Set(id, *st)
		}
	}
	return out
}

func (s *MemStickerStore) Delete(stickerID discord.Snowflake) {
	s.index.removeEntity(stickerID)
	s.base.Delete(stickerID)
}

// DeleteGuild removes all stickers for guildID. Call on GUILD_DELETE.
func (s *MemStickerStore) DeleteGuild(guildID discord.Snowflake) {
	ids := s.index.removeGuild(guildID)
	for _, id := range ids {
		s.base.Delete(id)
	}
}

// SetAll atomically replaces all stickers for guildID.
func (s *MemStickerStore) SetAll(guildID discord.Snowflake, stickers []*discord.Sticker) {
	newIDs := make([]discord.Snowflake, 0, len(stickers))
	add := make(map[discord.Snowflake]*discord.Sticker, len(stickers))
	for _, st := range stickers {
		if st == nil {
			continue
		}
		newIDs = append(newIDs, st.ID)
		add[st.ID] = st
	}
	toDelete := s.index.setForGuild(guildID, newIDs)
	s.base.ReplaceKeys(toDelete, add)
}

func (s *MemStickerStore) Size() int { return s.base.Len() }

func (s *MemStickerStore) SweepExpired(unusedWindow time.Duration) {
	s.base.SweepExpired(unusedWindow)
}

func (s *MemStickerStore) Bytes() int64 { return s.base.Bytes() }

func (s *MemStickerStore) TrimTo(n int) { s.base.TrimTo(n) }

func (s *MemStickerStore) IsEnabled() bool { return s.base.IsEnabled() }

func (s *MemStickerStore) Close() { s.base.Close() }
