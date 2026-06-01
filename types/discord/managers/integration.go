package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type integrationManager struct {
	guildID discord.Snowflake
	client  discord.EntityClient
}

func NewIntegrationManager(guildID discord.Snowflake, c discord.EntityClient) discord.IntegrationManager {
	return &integrationManager{guildID: guildID, client: c}
}

func (m *integrationManager) FetchAll(ctx context.Context) ([]*discord.Integration, error) {
	return m.client.ListGuildIntegrations(ctx, m.guildID)
}

func (m *integrationManager) Remove(ctx context.Context, integrationID discord.Snowflake, reason *string) error {
	return m.client.DeleteGuildIntegration(ctx, m.guildID, integrationID, reason)
}
