package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type threadManager struct {
	channelID discord.Snowflake
	client    discord.EntityClient
}

func NewThreadManager(channelID discord.Snowflake, c discord.EntityClient) discord.ThreadManager {
	return &threadManager{channelID: channelID, client: c}
}

func (m *threadManager) Cache() *collection.Collection[discord.Snowflake, *discord.Channel] {
	if c := m.client.ClientCache(); c != nil {
		return c.Channels().All().Filter(func(ch *discord.Channel) bool {
			return ch.ParentID != nil && *ch.ParentID == m.channelID && isThread(ch)
		})
	}
	return collection.New[discord.Snowflake, *discord.Channel]()
}

func (m *threadManager) Get(threadID discord.Snowflake) (*discord.Channel, bool) {
	if c := m.client.ClientCache(); c != nil {
		ch, ok := c.Channels().Get(threadID)
		if ok && isThread(ch) && ch.ParentID != nil && *ch.ParentID == m.channelID {
			return ch, true
		}
		return nil, false
	}
	return nil, false
}

func (m *threadManager) Fetch(ctx context.Context, threadID discord.Snowflake) (*discord.Channel, error) {
	return m.client.GetChannel(ctx, threadID)
}

func (m *threadManager) Resolve(input any) *discord.Channel {
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

func (m *threadManager) Size() int {
	return m.Cache().Len()
}

func isThread(ch *discord.Channel) bool {
	return ch.Type == discord.ChannelTypeAnnouncementThread ||
		ch.Type == discord.ChannelTypePublicThread ||
		ch.Type == discord.ChannelTypePrivateThread
}
