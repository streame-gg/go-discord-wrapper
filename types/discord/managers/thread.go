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
	return m.client.ThreadsForParent(m.channelID)
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

func (m *threadManager) Resolve(input any) (*discord.Channel, error) {
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

func (m *threadManager) Size() int {
	return m.Cache().Len()
}

func isThread(ch *discord.Channel) bool {
	return ch.IsThread()
}
