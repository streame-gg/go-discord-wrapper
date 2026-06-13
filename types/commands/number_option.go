package commands

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionNumber struct {
	Type                     discord.ApplicationCommandOptionType      `json:"type"`
	Name                     string                                    `json:"name"`
	NameLocalizations        map[discord.Locale]string                 `json:"name_localizations,omitempty"`
	Description              string                                    `json:"description"`
	DescriptionLocalizations map[discord.Locale]string                 `json:"description_localizations,omitempty"`
	Required                 bool                                      `json:"required"`
	Choices                  []ApplicationCommandOptionChoice[float64] `json:"choices,omitempty"`
	MinValue                 *float64                                  `json:"min_value,omitempty"`
	MaxValue                 *float64                                  `json:"max_value,omitempty"`
	Autocomplete             bool                                      `json:"autocomplete"`
}

func (o *ApplicationCommandOptionNumber) ApplicationCommandOptionType() discord.ApplicationCommandOptionType {
	return discord.ApplicationCommandOptionTypeNumber
}

func (o *ApplicationCommandOptionNumber) MarshalJSON() ([]byte, error) {
	type Alias ApplicationCommandOptionNumber
	return json.Marshal(struct {
		Alias
		Type discord.ApplicationCommandOptionType `json:"type"`
	}{
		Alias: Alias(*o),
		Type:  o.ApplicationCommandOptionType(),
	})
}

func (o *ApplicationCommandOptionNumber) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionNumber
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
