package discord

import (
	"time"
)

// https://docs.discord.com/developers/resources/guild#integration-object
type IntegrationType string

const (
	IntegrationTypeTwitch            IntegrationType = "twitch"
	IntegrationTypeYouTube           IntegrationType = "youtube"
	IntegrationTypeDiscord           IntegrationType = "discord"
	IntegrationTypeGuildSubscription IntegrationType = "guild_subscription"
)

// https://docs.discord.com/developers/resources/guild#integration-object-integration-expire-behaviors
type IntegrationExpireBehavior int

const (
	IntegrationExpireBehaviorRemoveRole IntegrationExpireBehavior = 0
	IntegrationExpireBehaviorKick       IntegrationExpireBehavior = 1
)

// https://docs.discord.com/developers/resources/guild#integration-account-object
type IntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// https://docs.discord.com/developers/resources/guild#integration-application-object
type IntegrationApplication struct {
	ID          Snowflake `json:"id"`
	Name        string    `json:"name"`
	IconHash    *string   `json:"icon,omitempty"`
	Description string    `json:"description"`
	Bot         *User     `json:"bot,omitempty"`
}

// https://docs.discord.com/developers/resources/guild#integration-object
type Integration struct {
	hClient EntityClient
	GuildID Snowflake `json:"-"`

	ID                Snowflake                  `json:"id"`
	Name              string                     `json:"name"`
	Type              IntegrationType            `json:"type"`
	Enabled           bool                       `json:"enabled"`
	Syncing           bool                       `json:"syncing,omitempty"`
	RoleID            *Snowflake                 `json:"role_id,omitempty"`
	EnableEmoticons   bool                       `json:"enable_emoticons,omitempty"`
	ExpireBehavior    *IntegrationExpireBehavior `json:"expire_behavior,omitempty"`
	ExpireGracePeriod int                        `json:"expire_grace_period,omitempty"`
	User              *User                      `json:"user,omitempty"`
	Account           IntegrationAccount         `json:"account"`
	SyncedAt          *time.Time                 `json:"synced_at,omitempty"`
	SubscriberCount   int                        `json:"subscriber_count,omitempty"`
	Revoked           bool                       `json:"revoked,omitempty"`
	Application       *IntegrationApplication    `json:"application,omitempty"`
	Scopes            []Scope                    `json:"scopes,omitempty"`
}
