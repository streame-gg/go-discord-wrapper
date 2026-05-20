package discord

type WebhookType int

const (
	WebhookTypeIncoming        WebhookType = 1
	WebhookTypeChannelFollower WebhookType = 2
	WebhookTypeApplication     WebhookType = 3
)

type Webhook struct {
	hClient       EntityClient
	ID            Snowflake   `json:"id"`
	Type          WebhookType `json:"type"`
	GuildID       *Snowflake  `json:"guild_id,omitempty"`
	ChannelID     *Snowflake  `json:"channel_id,omitempty"`
	User          *User       `json:"user,omitempty"`
	Name          *string     `json:"name,omitempty"`
	Avatar        *string     `json:"avatar,omitempty"`
	Token         *string     `json:"token,omitempty"`
	ApplicationID *Snowflake  `json:"application_id,omitempty"`
	SourceGuild   *Guild      `json:"source_guild,omitempty"`
	SourceChannel *Channel    `json:"source_channel,omitempty"`
	URL           *string     `json:"url,omitempty"`
}
