package commands

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionInteger struct {
	Type                     discord.ApplicationCommandOptionType    `json:"type"`
	Name                     string                                  `json:"name"`
	NameLocalizations        map[discord.Locale]string               `json:"name_localizations,omitempty"`
	Description              string                                  `json:"description"`
	DescriptionLocalizations map[discord.Locale]string               `json:"description_localizations,omitempty"`
	Required                 bool                                    `json:"required"`
	Choices                  []ApplicationCommandOptionChoice[int64] `json:"choices,omitempty"`
	MinValue                 *int64                                  `json:"min_value,omitempty"`
	MaxValue                 *int64                                  `json:"max_value,omitempty"`
	Autocomplete             bool                                    `json:"autocomplete"`
}

func (o *ApplicationCommandOptionInteger) ApplicationCommandOptionType() discord.ApplicationCommandOptionType {
	return discord.ApplicationCommandOptionTypeInteger
}

func (o *ApplicationCommandOptionInteger) MarshalJSON() ([]byte, error) {
	type Alias ApplicationCommandOptionInteger
	return json.Marshal(struct {
		Alias
		Type discord.ApplicationCommandOptionType `json:"type"`
	}{
		Alias: Alias(*o),
		Type:  o.ApplicationCommandOptionType(),
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

func NewIntegerOptionBuilder() *ApplicationCommandOptionInteger {
	return &ApplicationCommandOptionInteger{
		Type: discord.ApplicationCommandOptionTypeInteger,
	}
}

func (o *ApplicationCommandOptionInteger) SetName(name string) *ApplicationCommandOptionInteger {
	o.Name = name
	return o
}

func (o *ApplicationCommandOptionInteger) SetNameLocalizations(localizations map[discord.Locale]string) *ApplicationCommandOptionInteger {
	o.NameLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionInteger) SetDescription(description string) *ApplicationCommandOptionInteger {
	o.Description = description
	return o
}

func (o *ApplicationCommandOptionInteger) SetDescriptionLocalizations(localizations map[discord.Locale]string) *ApplicationCommandOptionInteger {
	o.DescriptionLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionInteger) SetRequired(required bool) *ApplicationCommandOptionInteger {
	o.Required = required
	return o
}

func (o *ApplicationCommandOptionInteger) SetMinValue(minValue int64) *ApplicationCommandOptionInteger {
	o.MinValue = &minValue
	return o
}

func (o *ApplicationCommandOptionInteger) SetMaxValue(maxValue int64) *ApplicationCommandOptionInteger {
	o.MaxValue = &maxValue
	return o
}

func (o *ApplicationCommandOptionInteger) SetAutocomplete(autocomplete bool) *ApplicationCommandOptionInteger {
	o.Autocomplete = autocomplete
	return o
}

func (o *ApplicationCommandOptionInteger) SetChoices(choices []ApplicationCommandOptionChoice[int64]) *ApplicationCommandOptionInteger {
	o.Choices = choices
	return o
}

func (o *ApplicationCommandOptionInteger) Build() ApplicationCommandOptionInteger {
	return *o
}
