package commands

import (
	"encoding/json"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type ApplicationCommandOptionBoolean struct {
	Type                     discord.ApplicationCommandOptionType `json:"type"`
	Name                     string                               `json:"name"`
	NameLocalizations        map[discord.Locale]string            `json:"name_localizations,omitempty"`
	Description              string                               `json:"description"`
	DescriptionLocalizations map[discord.Locale]string            `json:"description_localizations,omitempty"`
	Required                 *bool                                `json:"required,omitempty"`
}

func (o *ApplicationCommandOptionBoolean) ApplicationCommandOptionType() discord.ApplicationCommandOptionType {
	return discord.ApplicationCommandOptionTypeBoolean
}

func (o *ApplicationCommandOptionBoolean) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionBoolean
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

func (o *ApplicationCommandOptionBoolean) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionBoolean
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
