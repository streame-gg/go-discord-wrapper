package managers

import (
	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type voiceStateManager struct {
	guildID discord.Snowflake
	client  discord.EntityClient
}

func NewVoiceStateManager(guildID discord.Snowflake, c discord.EntityClient) discord.VoiceStateManager {
	return &voiceStateManager{guildID: guildID, client: c}
}

func (m *voiceStateManager) Cache() *collection.Collection[discord.Snowflake, *discord.VoiceState] {
	if c := m.client.ClientCache(); c != nil {
		return c.VoiceStates().AllInGuild(m.guildID)
	}
	return collection.New[discord.Snowflake, *discord.VoiceState]()
}

func (m *voiceStateManager) Get(userID discord.Snowflake) (*discord.VoiceState, bool) {
	if c := m.client.ClientCache(); c != nil {
		return c.VoiceStates().Get(m.guildID, userID)
	}
	return nil, false
}

func (m *voiceStateManager) Resolve(input any) *discord.VoiceState {
	switch v := input.(type) {
	case *discord.VoiceState:
		return v
	case discord.Snowflake:
		vs, _ := m.Get(v)
		return vs
	case string:
		vs, _ := m.Get(discord.Snowflake(v))
		return vs
	}
	return nil
}

func (m *voiceStateManager) Size() int {
	if c := m.client.ClientCache(); c != nil {
		return c.VoiceStates().AllInGuild(m.guildID).Len()
	}
	return 0
}
