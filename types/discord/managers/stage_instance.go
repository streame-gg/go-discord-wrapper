package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type stageInstanceManager struct {
	guildID discord.Snowflake
	client  discord.EntityClient
}

func NewStageInstanceManager(guildID discord.Snowflake, c discord.EntityClient) discord.StageInstanceManager {
	return &stageInstanceManager{guildID: guildID, client: c}
}

func (m *stageInstanceManager) Cache() *collection.Collection[discord.Snowflake, *discord.StageInstance] {
	if c := m.client.ClientCache(); c != nil {
		return c.StageInstances().GetByGuild(m.guildID)
	}
	return collection.New[discord.Snowflake, *discord.StageInstance]()
}

func (m *stageInstanceManager) Get(instanceID discord.Snowflake) (*discord.StageInstance, bool) {
	if c := m.client.ClientCache(); c != nil {
		return c.StageInstances().Get(instanceID)
	}
	return nil, false
}

func (m *stageInstanceManager) Fetch(ctx context.Context, channelID discord.Snowflake) (*discord.StageInstance, error) {
	return m.client.GetStageInstance(ctx, channelID)
}

func (m *stageInstanceManager) Create(ctx context.Context, opts discord.StageCreateOptions) (*discord.StageInstance, error) {
	return m.client.CreateStageInstance(ctx, opts)
}

func (m *stageInstanceManager) Resolve(input any) (*discord.StageInstance, error) {
	switch v := input.(type) {
	case *discord.StageInstance:
		return v, nil
	case discord.Snowflake:
		si, _ := m.Get(v)
		return si, nil
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

func (m *stageInstanceManager) Size() int {
	if c := m.client.ClientCache(); c != nil {
		return c.StageInstances().GetByGuild(m.guildID).Len()
	}
	return 0
}
