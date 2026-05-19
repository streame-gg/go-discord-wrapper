package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type clientGuildManager struct {
	client discord.EntityClient
}

func NewClientGuildManager(c discord.EntityClient) discord.GuildManager {
	return &clientGuildManager{client: c}
}

func (m *clientGuildManager) Cache() *collection.Collection[discord.Snowflake, *discord.Guild] {
	if c := m.client.ClientCache(); c != nil {
		return c.Guilds().All()
	}
	return collection.New[discord.Snowflake, *discord.Guild]()
}

func (m *clientGuildManager) Get(guildID discord.Snowflake) (*discord.Guild, bool) {
	if c := m.client.ClientCache(); c != nil {
		return c.Guilds().Get(guildID)
	}
	return nil, false
}

func (m *clientGuildManager) Fetch(ctx context.Context, guildID discord.Snowflake) (*discord.Guild, error) {
	return m.client.GetGuild(ctx, guildID)
}

func (m *clientGuildManager) Resolve(input any) *discord.Guild {
	switch v := input.(type) {
	case *discord.Guild:
		return v
	case discord.Snowflake:
		g, _ := m.Get(v)
		return g
	case string:
		g, _ := m.Get(discord.Snowflake(v))
		return g
	}
	return nil
}

func (m *clientGuildManager) Size() int {
	if c := m.client.ClientCache(); c != nil {
		return c.Guilds().Size()
	}
	return 0
}
