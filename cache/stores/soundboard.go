package stores

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// MemSoundboardStore is the in-memory implementation of cache.SoundboardStore.
type MemSoundboardStore struct {
	base  *BaseStore[discord.Snowflake, *discord.SoundboardSound]
	index *guildIndex
}

// NewSoundboardStore creates a MemSoundboardStore with the given options.
func NewSoundboardStore(opts StoreOptions) *MemSoundboardStore {
	return &MemSoundboardStore{
		base:  NewBaseStore[discord.Snowflake, *discord.SoundboardSound](opts),
		index: newGuildIndex(),
	}
}

// Set stores a sound. guildID is passed explicitly because
// discord.SoundboardSound.GuildID is a pointer and may be nil.
func (s *MemSoundboardStore) Set(guildID discord.Snowflake, sound *discord.SoundboardSound) {
	if sound == nil {
		return
	}
	s.index.add(guildID, sound.SoundID)
	s.base.Set(sound.SoundID, sound)
}

func (s *MemSoundboardStore) Get(soundID discord.Snowflake) (*discord.SoundboardSound, bool) {
	return s.base.Get(soundID)
}

// GetByGuild returns a snapshot Collection of all sounds for guildID, keyed by sound ID.
func (s *MemSoundboardStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.SoundboardSound] {
	ids := s.index.forGuild(guildID)
	out := collection.NewWithCapacity[discord.Snowflake, *discord.SoundboardSound](len(ids))
	for _, id := range ids {
		if snd, ok := s.base.Get(id); ok {
			out.Set(id, snd)
		}
	}
	return out
}

func (s *MemSoundboardStore) Delete(soundID discord.Snowflake) {
	s.index.removeEntity(soundID)
	s.base.Delete(soundID)
}

// DeleteGuild removes all sounds for guildID. Call on GUILD_DELETE.
func (s *MemSoundboardStore) DeleteGuild(guildID discord.Snowflake) {
	ids := s.index.removeGuild(guildID)
	for _, id := range ids {
		s.base.Delete(id)
	}
}

// SetAll atomically replaces all sounds for guildID.
func (s *MemSoundboardStore) SetAll(guildID discord.Snowflake, sounds []discord.SoundboardSound) {
	newIDs := make([]discord.Snowflake, 0, len(sounds))
	add := make(map[discord.Snowflake]*discord.SoundboardSound, len(sounds))
	for _, snd := range sounds {
		newIDs = append(newIDs, snd.SoundID)
		add[snd.SoundID] = &snd
	}
	toDelete := s.index.setForGuild(guildID, newIDs)
	s.base.ReplaceKeys(toDelete, add)
}

func (s *MemSoundboardStore) Size() int { return s.base.Len() }

func (s *MemSoundboardStore) SweepExpired(unusedWindow time.Duration) {
	s.base.SweepExpired(unusedWindow)
}

func (s *MemSoundboardStore) Bytes() int64 { return s.base.Bytes() }

func (s *MemSoundboardStore) TrimTo(n int) { s.base.TrimTo(n) }

func (s *MemSoundboardStore) IsEnabled() bool { return s.base.IsEnabled() }

func (s *MemSoundboardStore) Close() { s.base.Close() }
