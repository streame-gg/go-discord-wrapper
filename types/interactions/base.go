package interactions

import (
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type Interaction struct {
	ID                           common.Snowflake                                             `json:"id"`
	ApplicationID                common.Snowflake                                             `json:"application_id"`
	Type                         common.InteractionType                                       `json:"type"`
	Data                         common.InteractionData                                       `json:"data,omitempty"`
	GuildID                      *common.Snowflake                                            `json:"guild_id,omitempty"`
	ChannelID                    *common.Snowflake                                            `json:"channel_id,omitempty"`
	Guild                        *common.Guild                                                `json:"guild,omitempty"`
	Channel                      *common.Channel                                              `json:"channel,omitempty"`
	Member                       *common.GuildMember                                          `json:"member,omitempty"`
	User                         *common.User                                                 `json:"user,omitempty"`
	Token                        string                                                       `json:"token"`
	Version                      int                                                          `json:"version"`
	Message                      *common.Message                                              `json:"message,omitempty"`
	AppPermissions               string                                                       `json:"app_permissions,omitempty"`
	Locale                       *common.Locale                                               `json:"locale,omitempty"`
	GuildLocale                  string                                                       `json:"guild_locale,omitempty"`
	Entitlements                 []common.Entitlement                                         `json:"entitlements,omitempty"`
	AuthorizingIntegrationOwners map[common.InteractionApplicationIntegrationType]interface{} `json:"authorizing_integration_owners,omitempty"`
	Context                      common.InteractionContextType                                `json:"context,omitempty"`
	AttachmentSizeLimit          int                                                          `json:"attachment_size_limit,omitempty"`
}
