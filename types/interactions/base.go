package interactions

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type Interaction struct {
	ID                           discord.Snowflake                                             `json:"id"`
	ApplicationID                discord.Snowflake                                             `json:"application_id"`
	Type                         discord.InteractionType                                       `json:"type"`
	Data                         discord.InteractionData                                       `json:"data,omitempty"`
	GuildID                      *discord.Snowflake                                            `json:"guild_id,omitempty"`
	ChannelID                    *discord.Snowflake                                            `json:"channel_id,omitempty"`
	Guild                        *discord.Guild                                                `json:"guild,omitempty"`
	Channel                      *discord.Channel                                              `json:"channel,omitempty"`
	Member                       *discord.GuildMember                                          `json:"member,omitempty"`
	User                         *discord.User                                                 `json:"user,omitempty"`
	Token                        string                                                        `json:"token"`
	Version                      int                                                           `json:"version"`
	Message                      *discord.Message                                              `json:"message,omitempty"`
	AppPermissions               string                                                        `json:"app_permissions,omitempty"`
	Locale                       *discord.Locale                                               `json:"locale,omitempty"`
	GuildLocale                  string                                                        `json:"guild_locale,omitempty"`
	Entitlements                 []discord.Entitlement                                         `json:"entitlements,omitempty"`
	AuthorizingIntegrationOwners map[discord.InteractionApplicationIntegrationType]interface{} `json:"authorizing_integration_owners,omitempty"`
	Context                      discord.InteractionContextType                                `json:"context,omitempty"`
	AttachmentSizeLimit          int                                                           `json:"attachment_size_limit,omitempty"`
}
