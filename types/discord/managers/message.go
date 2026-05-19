package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type messageManager struct {
	channelID discord.Snowflake
	client    discord.EntityClient
}

func NewMessageManager(channelID discord.Snowflake, c discord.EntityClient) discord.MessageManager {
	return &messageManager{channelID: channelID, client: c}
}

func (m *messageManager) Cache() *collection.Collection[discord.Snowflake, *discord.Message] {
	if c := m.client.ClientCache(); c != nil {
		return c.Messages().Channel(m.channelID)
	}
	return collection.New[discord.Snowflake, *discord.Message]()
}

func (m *messageManager) Get(messageID discord.Snowflake) (*discord.Message, bool) {
	if c := m.client.ClientCache(); c != nil {
		return c.Messages().Get(m.channelID, messageID)
	}
	return nil, false
}

func (m *messageManager) Fetch(ctx context.Context, messageID discord.Snowflake) (*discord.Message, error) {
	return m.client.GetMessage(ctx, m.channelID, messageID)
}

func (m *messageManager) FetchAll(ctx context.Context, opts discord.FetchMessagesOptions) (*collection.Collection[discord.Snowflake, *discord.Message], error) {
	messages, err := m.client.GetChannelMessages(ctx, m.channelID, opts)
	if err != nil {
		return nil, err
	}
	return collection.FromSlice(messages, func(msg *discord.Message) discord.Snowflake {
		return msg.ID
	}), nil
}

func (m *messageManager) Create(ctx context.Context, opts discord.MessageCreateOptions) (*discord.Message, error) {
	return m.client.CreateMessage(ctx, m.channelID, opts)
}

func (m *messageManager) Resolve(input any) *discord.Message {
	switch v := input.(type) {
	case *discord.Message:
		return v
	case discord.Snowflake:
		msg, _ := m.Get(v)
		return msg
	case string:
		msg, _ := m.Get(discord.Snowflake(v))
		return msg
	}
	return nil
}

func (m *messageManager) Size() int {
	if c := m.client.ClientCache(); c != nil {
		return c.Messages().Channel(m.channelID).Len()
	}
	return 0
}
