package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type banManager struct {
	guildID discord.Snowflake
	client  discord.EntityClient
}

func NewBanManager(guildID discord.Snowflake, c discord.EntityClient) discord.BanManager {
	return &banManager{guildID: guildID, client: c}
}

func (m *banManager) Fetch(ctx context.Context, userID discord.Snowflake) (*discord.Ban, error) {
	return m.client.GetGuildBan(ctx, m.guildID, userID)
}

func (m *banManager) FetchAll(ctx context.Context, opts discord.FetchBansOptions) ([]*discord.Ban, error) {
	return m.client.GetGuildBans(ctx, m.guildID, opts)
}

func (m *banManager) Create(ctx context.Context, userID discord.Snowflake, opts discord.BanOptions) error {
	return m.client.CreateGuildBan(ctx, m.guildID, userID, opts)
}

func (m *banManager) Remove(ctx context.Context, userID discord.Snowflake, reason *string) error {
	return m.client.RemoveGuildBan(ctx, m.guildID, userID, reason)
}
