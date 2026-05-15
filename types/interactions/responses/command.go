package responses

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type InteractionDataApplicationCommand struct {
	ID          discord.Snowflake                                       `json:"id"`
	CommandName string                                                  `json:"name"`
	Type        discord.ApplicationCommandType                          `json:"type"`
	GuildID     *discord.Snowflake                                      `json:"guild_id,omitempty"`
	TargetID    *discord.Snowflake                                      `json:"target_id,omitempty"`
	Resolved    *discord.ResolvedData                                   `json:"resolved,omitempty"`
	Options     *[]ApplicationCommandInteractionDataOption[interface{}] `json:"options,omitempty"`
}

func (d *InteractionDataApplicationCommand) GetType() discord.InteractionDataType {
	return discord.InteractionDataTypeApplicationCommand
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
		}
		d.Options = &options
	}

	return nil
}
