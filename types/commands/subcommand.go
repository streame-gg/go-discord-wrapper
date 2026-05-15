package commands

import (
	"encoding/json"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type ApplicationCommandOptionSubCommand struct {
	Type                     discord.ApplicationCommandOptionType `json:"type"`
	Name                     string                               `json:"name"`
	NameLocalizations        map[discord.Locale]string            `json:"name_localizations,omitempty"`
	Description              string                               `json:"description"`
	DescriptionLocalizations map[discord.Locale]string            `json:"description_localizations,omitempty"`
	Options                  *[]AnyApplicationCommandOption       `json:"options,omitempty"`
}

func (o *ApplicationCommandOptionSubCommand) ApplicationCommandOptionType() discord.ApplicationCommandOptionType {
	return discord.ApplicationCommandOptionTypeSubCommand
}

func (o *ApplicationCommandOptionSubCommand) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionSubCommand
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

func (o *ApplicationCommandOptionSubCommand) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionSubCommand
	raw := &struct {
		*Alias
		Options []json.RawMessage `json:"options,omitempty"`
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Options != nil {
		opts, err := unmarshalOptionSlice(raw.Options)
		if err != nil {
			return err
		}
		o.Options = &opts
	}

	return nil
}
