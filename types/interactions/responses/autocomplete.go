package responses

import (
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type InteractionDataAutocomplete struct {
	ID          common.Snowflake                                        `json:"id"`
	CommandName string                                                  `json:"name"`
	Type        common.ApplicationCommandType                           `json:"type"`
	GuildID     *common.Snowflake                                       `json:"guild_id,omitempty"`
	TargetID    *common.Snowflake                                       `json:"target_id,omitempty"`
	Resolved    *common.ResolvedData                                    `json:"resolved,omitempty"`
	Options     *[]ApplicationCommandInteractionDataOption[interface{}] `json:"options,omitempty"`
}

func (d *InteractionDataAutocomplete) GetType() common.InteractionDataType {
	return common.InteractionDataTypeApplicationCommandAutocomplete
}
