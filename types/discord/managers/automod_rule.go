package managers

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type autoModRuleManager struct {
	guildID discord.Snowflake
	client  discord.EntityClient
}

func NewAutoModRuleManager(guildID discord.Snowflake, c discord.EntityClient) discord.AutoModRuleManager {
	return &autoModRuleManager{guildID: guildID, client: c}
}

func (m *autoModRuleManager) Fetch(ctx context.Context, ruleID discord.Snowflake) (*discord.AutoModerationRule, error) {
	return m.client.GetAutoModerationRule(ctx, m.guildID, ruleID)
}

func (m *autoModRuleManager) FetchAll(ctx context.Context) ([]*discord.AutoModerationRule, error) {
	return m.client.ListAutoModerationRules(ctx, m.guildID)
}

func (m *autoModRuleManager) Create(ctx context.Context, opts discord.RuleCreateOptions) (*discord.AutoModerationRule, error) {
	return m.client.CreateAutoModerationRule(ctx, m.guildID, opts)
}
