package commands

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionSubCommandGroup struct {
	Type                     discord.ApplicationCommandOptionType `json:"type"`
	Name                     string                               `json:"name"`
	NameLocalizations        map[discord.Locale]string            `json:"name_localizations,omitempty"`
	Description              string                               `json:"description"`
	DescriptionLocalizations map[discord.Locale]string            `json:"description_localizations,omitempty"`
	Options                  []AnyApplicationCommandOption        `json:"options,omitempty"`
}

func (o *ApplicationCommandOptionSubCommandGroup) ApplicationCommandOptionType() discord.ApplicationCommandOptionType {
	return discord.ApplicationCommandOptionTypeSubCommandGroup
}

func (o *ApplicationCommandOptionSubCommandGroup) MarshalJSON() ([]byte, error) {
	type Alias ApplicationCommandOptionSubCommandGroup
	return json.Marshal(struct {
		Alias
		Type discord.ApplicationCommandOptionType `json:"type"`
	}{
		Alias: Alias(*o),
		Type:  o.ApplicationCommandOptionType(),
	})
}

func (o *ApplicationCommandOptionSubCommandGroup) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionSubCommandGroup
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
		o.Options = opts
	}

	return nil
}
