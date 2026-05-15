package commands

import (
	"encoding/json"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type ApplicationCommandOptionString struct {
	Type                     discord.ApplicationCommandOptionType     `json:"type"`
	Name                     string                                   `json:"name"`
	NameLocalizations        map[discord.Locale]string                `json:"name_localizations,omitempty"`
	Description              string                                   `json:"description"`
	DescriptionLocalizations map[discord.Locale]string                `json:"description_localizations,omitempty"`
	Required                 *bool                                    `json:"required,omitempty"`
	Choices                  []ApplicationCommandOptionChoice[string] `json:"choices,omitempty"`
	Autocomplete             *bool                                    `json:"autocomplete,omitempty"`
}

func (o *ApplicationCommandOptionString) ApplicationCommandOptionType() discord.ApplicationCommandOptionType {
	return discord.ApplicationCommandOptionTypeString
}

func (o *ApplicationCommandOptionString) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionString
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

func (o *ApplicationCommandOptionString) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionString
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
