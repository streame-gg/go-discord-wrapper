package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type guildWebhookManager struct {
	guildID discord.Snowflake
	client  discord.EntityClient
}

func NewGuildWebhookManager(guildID discord.Snowflake, c discord.EntityClient) discord.GuildWebhookManager {
	return &guildWebhookManager{guildID: guildID, client: c}
}

func (m *guildWebhookManager) FetchAll(ctx context.Context) ([]*discord.Webhook, error) {
	return m.client.ListGuildWebhooks(ctx, m.guildID)
}
