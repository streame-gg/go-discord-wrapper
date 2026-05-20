package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type clientChannelManager struct {
	client discord.EntityClient
}

func NewClientChannelManager(c discord.EntityClient) discord.ChannelManager {
	return &clientChannelManager{client: c}
}

func (m *clientChannelManager) Cache() *collection.Collection[discord.Snowflake, *discord.Channel] {
	if c := m.client.ClientCache(); c != nil {
		return c.Channels().All()
	}
	return collection.New[discord.Snowflake, *discord.Channel]()
}

func (m *clientChannelManager) Get(channelID discord.Snowflake) (*discord.Channel, bool) {
	if c := m.client.ClientCache(); c != nil {
		return c.Channels().Get(channelID)
	}
	return nil, false
}

func (m *clientChannelManager) Fetch(ctx context.Context, channelID discord.Snowflake) (*discord.Channel, error) {
	return m.client.GetChannel(ctx, channelID)
}

func (m *clientChannelManager) Resolve(input any) *discord.Channel {
	switch v := input.(type) {
	case *discord.Channel:
		return v
	case discord.Snowflake:
		ch, _ := m.Get(v)
		return ch
	case string:
		ch, _ := m.Get(discord.Snowflake(v))
		return ch
	}
	return nil
}

func (m *clientChannelManager) Size() int {
	if c := m.client.ClientCache(); c != nil {
		return c.Channels().Size()
	}
	return 0
}
