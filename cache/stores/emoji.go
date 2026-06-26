package stores

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// MemEmojiStore is the in-memory implementation of cache.EmojiStore.
type MemEmojiStore struct {
	base  *BaseStore[discord.Snowflake, *discord.Emoji]
	index *guildIndex
}

// NewEmojiStore creates a MemEmojiStore with the given options.
func NewEmojiStore(opts StoreOptions) *MemEmojiStore {
	return &MemEmojiStore{
		base:  NewBaseStore[discord.Snowflake, *discord.Emoji](opts),
		index: newGuildIndex(),
	}
}

func (s *MemEmojiStore) Set(guildID discord.Snowflake, emoji *discord.Emoji) {
	if emoji == nil {
		return
	}
	s.index.add(guildID, emoji.ID)
	s.base.Set(emoji.ID, emoji)
}

func (s *MemEmojiStore) Get(emojiID discord.Snowflake) (*discord.Emoji, bool) {
	return s.base.Get(emojiID)
}

// GetByGuild returns a snapshot Collection of all emojis for guildID, keyed by emoji ID.
func (s *MemEmojiStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, discord.Emoji] {
	ids := s.index.forGuild(guildID)
	out := collection.NewWithCapacity[discord.Snowflake, discord.Emoji](len(ids))
	for _, id := range ids {
		if e, ok := s.base.Get(id); ok {
			out.Set(id, *e)
		}
	}
	return out
}

func (s *MemEmojiStore) Delete(emojiID discord.Snowflake) {
	s.index.removeEntity(emojiID)
	s.base.Delete(emojiID)
}

// DeleteGuild removes all emojis for guildID. Call on GUILD_DELETE.
func (s *MemEmojiStore) DeleteGuild(guildID discord.Snowflake) {
	ids := s.index.removeGuild(guildID)
	for _, id := range ids {
		s.base.Delete(id)
	}
}

// SetAll atomically replaces all emojis for guildID.
func (s *MemEmojiStore) SetAll(guildID discord.Snowflake, emojis []*discord.Emoji) {
	newIDs := make([]discord.Snowflake, 0, len(emojis))
	add := make(map[discord.Snowflake]*discord.Emoji, len(emojis))
	for _, e := range emojis {
		if e == nil {
			continue
		}
		newIDs = append(newIDs, e.ID)
		add[e.ID] = e
	}
	toDelete := s.index.setForGuild(guildID, newIDs)
	s.base.ReplaceKeys(toDelete, add)
}

func (s *MemEmojiStore) Size() int { return s.base.Len() }

func (s *MemEmojiStore) SweepExpired(unusedWindow time.Duration) {
	s.base.SweepExpired(unusedWindow)
}

func (s *MemEmojiStore) Bytes() int64 { return s.base.Bytes() }

func (s *MemEmojiStore) TrimTo(n int) { s.base.TrimTo(n) }

func (s *MemEmojiStore) IsEnabled() bool { return s.base.IsEnabled() }

func (s *MemEmojiStore) Close() { s.base.Close() }
