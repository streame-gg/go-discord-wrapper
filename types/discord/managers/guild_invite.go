package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type guildInviteManager struct {
	guildID discord.Snowflake
	client  discord.EntityClient
}

func NewGuildInviteManager(guildID discord.Snowflake, c discord.EntityClient) discord.GuildInviteManager {
	return &guildInviteManager{guildID: guildID, client: c}
}

func (m *guildInviteManager) FetchAll(ctx context.Context) ([]*discord.Invite, error) {
	return m.client.GetGuildInvites(ctx, m.guildID)
}
