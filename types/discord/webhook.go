package discord

type WebhookType int

const (
	WebhookTypeIncoming        WebhookType = 1
	WebhookTypeChannelFollower WebhookType = 2
	WebhookTypeApplication     WebhookType = 3
)

type Webhook struct {
	ID            Snowflake   `json:"id"`
	Type          WebhookType `json:"type"`
	GuildID       *Snowflake  `json:"guild_id,omitempty"`
	ChannelID     *Snowflake  `json:"channel_id"`
	User          *User       `json:"user,omitempty"`
	Name          *string     `json:"name"`
	Avatar        *string     `json:"avatar"`
	Token         *string     `json:"token,omitempty"`
	ApplicationID *Snowflake  `json:"application_id"`
	SourceGuild   *Guild      `json:"source_guild,omitempty"`
	SourceChannel *Channel    `json:"source_channel,omitempty"`
	URL           *string     `json:"url,omitempty"`
}
