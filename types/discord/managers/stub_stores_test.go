package managers_test

import (
	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// Additional stub stores for managers backed by the less common caches. Each is
// keyed guildID -> entityID -> entity so GetByGuild can return a per-guild view
// while Get searches across guilds (mirroring the real store semantics).

type stubVoiceStateStore struct {
	states map[discord.Snowflake]map[discord.Snowflake]*discord.VoiceState
}

func (s *stubVoiceStateStore) Get(guildID, userID discord.Snowflake) (*discord.VoiceState, bool) {
	g, ok := s.states[guildID]
	if !ok {
		return nil, false
	}
	v, ok := g[userID]
	return v, ok
}
func (s *stubVoiceStateStore) AllInGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.VoiceState] {
	c := collection.New[discord.Snowflake, *discord.VoiceState]()
	for uid, v := range s.states[guildID] {
		c.Set(uid, v)
	}
	return c
}
func (s *stubVoiceStateStore) Size() int { return len(s.states) }

type stubSoundboardStore struct {
	sounds map[discord.Snowflake]map[discord.Snowflake]*discord.SoundboardSound
}

func (s *stubSoundboardStore) Get(soundID discord.Snowflake) (*discord.SoundboardSound, bool) {
	for _, g := range s.sounds {
		if snd, ok := g[soundID]; ok {
			return snd, true
		}
	}
	return nil, false
}
func (s *stubSoundboardStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.SoundboardSound] {
	c := collection.New[discord.Snowflake, *discord.SoundboardSound]()
	for id, snd := range s.sounds[guildID] {
		c.Set(id, snd)
	}
	return c
}
func (s *stubSoundboardStore) Size() int { return len(s.sounds) }

type stubScheduledEventStore struct {
	events map[discord.Snowflake]map[discord.Snowflake]*discord.GuildScheduledEvent
}

func (s *stubScheduledEventStore) Get(eventID discord.Snowflake) (*discord.GuildScheduledEvent, bool) {
	for _, g := range s.events {
		if e, ok := g[eventID]; ok {
			return e, true
		}
	}
	return nil, false
}
func (s *stubScheduledEventStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.GuildScheduledEvent] {
	c := collection.New[discord.Snowflake, *discord.GuildScheduledEvent]()
	for id, e := range s.events[guildID] {
		c.Set(id, e)
	}
	return c
}
func (s *stubScheduledEventStore) Size() int { return len(s.events) }

type stubStageInstanceStore struct {
	instances map[discord.Snowflake]map[discord.Snowflake]*discord.StageInstance
}

func (s *stubStageInstanceStore) Get(instanceID discord.Snowflake) (*discord.StageInstance, bool) {
	for _, g := range s.instances {
		if si, ok := g[instanceID]; ok {
			return si, true
		}
	}
	return nil, false
}
func (s *stubStageInstanceStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.StageInstance] {
	c := collection.New[discord.Snowflake, *discord.StageInstance]()
	for id, si := range s.instances[guildID] {
		c.Set(id, si)
	}
	return c
}
func (s *stubStageInstanceStore) Size() int { return len(s.instances) }

type stubEmojiStore struct {
	emojis map[discord.Snowflake]map[discord.Snowflake]discord.Emoji
}

func (s *stubEmojiStore) Get(emojiID discord.Snowflake) (*discord.Emoji, bool) {
	for _, g := range s.emojis {
		if e, ok := g[emojiID]; ok {
			return &e, true
		}
	}
	return nil, false
}
func (s *stubEmojiStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, discord.Emoji] {
	c := collection.New[discord.Snowflake, discord.Emoji]()
	for id, e := range s.emojis[guildID] {
		c.Set(id, e)
	}
	return c
}
func (s *stubEmojiStore) Size() int { return len(s.emojis) }

type stubStickerStore struct {
	stickers map[discord.Snowflake]map[discord.Snowflake]*discord.Sticker
}

func (s *stubStickerStore) Get(stickerID discord.Snowflake) (*discord.Sticker, bool) {
	for _, g := range s.stickers {
		if st, ok := g[stickerID]; ok {
			return st, true
		}
	}
	return nil, false
}
func (s *stubStickerStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Sticker] {
	c := collection.New[discord.Snowflake, *discord.Sticker]()
	for id, st := range s.stickers[guildID] {
		c.Set(id, st)
	}
	return c
}
func (s *stubStickerStore) Size() int { return len(s.stickers) }
