package discord

import (
	"time"
)

type IntegrationType string

const (
	IntegrationTypeTwitch            IntegrationType = "twitch"
	IntegrationTypeYouTube           IntegrationType = "youtube"
	IntegrationTypeDiscord           IntegrationType = "discord"
	IntegrationTypeGuildSubscription IntegrationType = "guild_subscription"
)

type IntegrationExpireBehavior int

const (
	IntegrationExpireBehaviorRemoveRole IntegrationExpireBehavior = 0
	IntegrationExpireBehaviorKick       IntegrationExpireBehavior = 1
)

type IntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type IntegrationApplication struct {
	ID          Snowflake `json:"id"`
	Name        string    `json:"name"`
	IconHash    *string   `json:"icon,omitempty"`
	Description string    `json:"description"`
	Bot         *User     `json:"bot,omitempty"`
}

type Integration struct {
	ID                Snowflake                  `json:"id"`
	Name              string                     `json:"name"`
	Type              IntegrationType            `json:"type"`
	Enabled           bool                       `json:"enabled"`
	Syncing           *bool                      `json:"syncing,omitempty"`
	RoleID            *Snowflake                 `json:"role_id,omitempty"`
	EnableEmoticons   *bool                      `json:"enable_emoticons,omitempty"`
	ExpireBehavior    *IntegrationExpireBehavior `json:"expire_behavior,omitempty"`
	ExpireGracePeriod *int                       `json:"expire_grace_period,omitempty"`
	User              *User                      `json:"user,omitempty"`
	Account           IntegrationAccount         `json:"account"`
	SyncedAt          *time.Time                 `json:"synced_at,omitempty"`
	SubscriberCount   *int                       `json:"subscriber_count,omitempty"`
	Revoked           *bool                      `json:"revoked,omitempty"`
	Application       *IntegrationApplication    `json:"application,omitempty"`
	Scopes            *[]string                  `json:"scopes,omitempty"`
}
