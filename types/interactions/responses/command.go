package responses

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type InteractionDataApplicationCommand struct {
	ID          common.Snowflake                                        `json:"id"`
	CommandName string                                                  `json:"name"`
	Type        common.ApplicationCommandType                           `json:"type"`
	GuildID     *common.Snowflake                                       `json:"guild_id,omitempty"`
	TargetID    *common.Snowflake                                       `json:"target_id,omitempty"`
	Resolved    *common.ResolvedData                                    `json:"resolved,omitempty"`
	Options     *[]ApplicationCommandInteractionDataOption[interface{}] `json:"options,omitempty"`
}

func (d *InteractionDataApplicationCommand) GetType() common.InteractionDataType {
	return common.InteractionDataTypeApplicationCommand
}

func (d *InteractionDataApplicationCommand) UnmarshalJSON(data []byte) error {
	type Alias InteractionDataApplicationCommand
	raw := &struct {
		*Alias
		Options []json.RawMessage `json:"options,omitempty"`
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	d.ID = raw.ID
	d.CommandName = raw.CommandName
	d.Type = raw.Type
	d.GuildID = raw.GuildID
	d.TargetID = raw.TargetID
	d.Resolved = raw.Resolved

	if raw.Options != nil {
		var options []ApplicationCommandInteractionDataOption[interface{}]
		for _, optionData := range raw.Options {
			var option ApplicationCommandInteractionDataOption[interface{}]
			if err := json.Unmarshal(optionData, &option); err != nil {
				return err
			}
			options = append(options, option)

			if option.Options != nil {
				var option ApplicationCommandInteractionDataOption[interface{}]
				if err := json.Unmarshal(optionData, &option); err != nil {
					return err
				}

			}
		}
		d.Options = &options
	}

	return nil
}
