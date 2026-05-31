package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type guildChannelManager struct {
	guildID discord.Snowflake
	client  discord.EntityClient
}

func NewGuildChannelManager(guildID discord.Snowflake, c discord.EntityClient) discord.GuildChannelManager {
	return &guildChannelManager{guildID: guildID, client: c}
}

func (m *guildChannelManager) Cache() *collection.Collection[discord.Snowflake, *discord.Channel] {
	return m.client.ChannelsForGuild(m.guildID)
}

func (m *guildChannelManager) Get(channelID discord.Snowflake) (*discord.Channel, bool) {
	if c := m.client.ClientCache(); c != nil {
		ch, ok := c.Channels().Get(channelID)
		if ok && ch.GuildID != nil && *ch.GuildID == m.guildID {
			return ch, true
		}
		return nil, false
	}
	return nil, false
}

func (m *guildChannelManager) Fetch(ctx context.Context, channelID discord.Snowflake) (*discord.Channel, error) {
	return m.client.GetChannel(ctx, channelID)
}

func (m *guildChannelManager) FetchAll(ctx context.Context) (*collection.Collection[discord.Snowflake, *discord.Channel], error) {
	channels, err := m.client.GetGuildChannels(ctx, m.guildID)
	if err != nil {
		return nil, err
	}
	return collection.FromSlice(channels, func(ch *discord.Channel) discord.Snowflake {
		return ch.ID
	}), nil
}

func (m *guildChannelManager) Create(ctx context.Context, opts discord.ChannelCreateOptions) (*discord.Channel, error) {
	return m.client.CreateGuildChannel(ctx, m.guildID, opts)
}

func (m *guildChannelManager) Resolve(input any) (*discord.Channel, error) {
	switch v := input.(type) {
	case *discord.Channel:
		return v, nil
	case discord.Snowflake:
		ch, _ := m.Get(v)
		return ch, nil
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

func (m *guildChannelManager) Size() int {
	return m.Cache().Len()
}
