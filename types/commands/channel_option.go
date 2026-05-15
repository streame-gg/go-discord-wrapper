package commands

import (
	"encoding/json"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type ApplicationCommandOptionChannel struct {
	Type                     discord.ApplicationCommandOptionType `json:"type"`
	Name                     string                               `json:"name"`
	NameLocalizations        map[discord.Locale]string            `json:"name_localizations,omitempty"`
	Description              string                               `json:"description"`
	DescriptionLocalizations map[discord.Locale]string            `json:"description_localizations,omitempty"`
	Required                 *bool                                `json:"required,omitempty"`
	ChannelTypes             []discord.ChannelType                `json:"channel_types,omitempty"`
}

func (o *ApplicationCommandOptionChannel) ApplicationCommandOptionType() discord.ApplicationCommandOptionType {
	return discord.ApplicationCommandOptionTypeChannel
}

func (o *ApplicationCommandOptionChannel) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionChannel
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

func (o *ApplicationCommandOptionChannel) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionChannel
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
