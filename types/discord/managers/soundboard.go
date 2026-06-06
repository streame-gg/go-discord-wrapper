package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type soundboardManager struct {
	guildID discord.Snowflake
	client  discord.EntityClient
}

func NewSoundboardManager(guildID discord.Snowflake, c discord.EntityClient) discord.SoundboardManager {
	return &soundboardManager{guildID: guildID, client: c}
}

func (m *soundboardManager) Cache() *collection.Collection[discord.Snowflake, *discord.SoundboardSound] {
	if c := m.client.ClientCache(); c != nil {
		return c.Soundboard().GetByGuild(m.guildID)
	}
	return collection.New[discord.Snowflake, *discord.SoundboardSound]()
}

func (m *soundboardManager) Get(soundID discord.Snowflake) (*discord.SoundboardSound, bool) {
	if c := m.client.ClientCache(); c != nil {
		return c.Soundboard().GetByGuild(m.guildID).Get(soundID)
	}
	return nil, false
}

func (m *soundboardManager) Fetch(ctx context.Context, soundID discord.Snowflake) (*discord.SoundboardSound, error) {
	return m.client.GetGuildSoundboardSound(ctx, m.guildID, soundID)
}

func (m *soundboardManager) FetchAll(ctx context.Context) (*collection.Collection[discord.Snowflake, *discord.SoundboardSound], error) {
	sounds, err := m.client.ListGuildSoundboardSounds(ctx, m.guildID)
	if err != nil {
		return nil, err
	}
	return collection.FromSlice(sounds, func(s *discord.SoundboardSound) discord.Snowflake {
		return s.SoundID
	}), nil
}

func (m *soundboardManager) Resolve(input any) (*discord.SoundboardSound, error) {
	switch v := input.(type) {
	case *discord.SoundboardSound:
		return v, nil
	case discord.Snowflake:
		s, _ := m.Get(v)
		return s, nil
	case string:
		snowflake, err := discord.ParseSnowflake(v)
		if err != nil {
			return nil, err
		}
		mem, _ := m.Get(*snowflake)
		return mem, nil
	}
	return nil, discord.ErrNotConvertable
}

func (m *soundboardManager) Size() int {
	if c := m.client.ClientCache(); c != nil {
		return c.Soundboard().GetByGuild(m.guildID).Len()
	}
	return 0
}
