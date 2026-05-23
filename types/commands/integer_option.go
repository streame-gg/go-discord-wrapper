package commands

import (
	"encoding/json"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type ApplicationCommandOptionInteger struct {
	Type                     discord.ApplicationCommandOptionType     `json:"type"`
	Name                     string                                   `json:"name"`
	NameLocalizations        map[discord.Locale]string                `json:"name_localizations,omitempty"`
	Description              string                                   `json:"description"`
	DescriptionLocalizations map[discord.Locale]string                `json:"description_localizations,omitempty"`
	Required                 *bool                                    `json:"required,omitempty"`
	Choices                  []ApplicationCommandOptionChoice[int64]  `json:"choices,omitempty"`
	MinValue                 *int64                                   `json:"min_value,omitempty"`
	MaxValue                 *int64                                   `json:"max_value,omitempty"`
	Autocomplete             *bool                                    `json:"autocomplete,omitempty"`
}

func (o *ApplicationCommandOptionInteger) ApplicationCommandOptionType() discord.ApplicationCommandOptionType {
	return discord.ApplicationCommandOptionTypeInteger
}

func (o *ApplicationCommandOptionInteger) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionInteger
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

func (o *ApplicationCommandOptionInteger) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionInteger
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}
