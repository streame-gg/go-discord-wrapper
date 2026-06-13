package commands

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionString struct {
	Type                     discord.ApplicationCommandOptionType     `json:"type"`
	Name                     string                                   `json:"name"`
	NameLocalizations        map[discord.Locale]string                `json:"name_localizations,omitempty"`
	Description              string                                   `json:"description"`
	DescriptionLocalizations map[discord.Locale]string                `json:"description_localizations,omitempty"`
	Required                 bool                                     `json:"required"`
	Choices                  []ApplicationCommandOptionChoice[string] `json:"choices,omitempty"`
	Autocomplete             bool                                     `json:"autocomplete"`
	MinLength                *int64                                   `json:"min_length,omitempty"`
	MaxLength                *int64                                   `json:"max_length,omitempty"`
}

func (o *ApplicationCommandOptionString) ApplicationCommandOptionType() discord.ApplicationCommandOptionType {
	return discord.ApplicationCommandOptionTypeString
}

func (o *ApplicationCommandOptionString) MarshalJSON() ([]byte, error) {
	type Alias ApplicationCommandOptionString
	return json.Marshal(struct {
		Alias
		Type discord.ApplicationCommandOptionType `json:"type"`
	}{
		Alias: Alias(*o),
		Type:  o.ApplicationCommandOptionType(),
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

func NewStringOptionBuilder() *ApplicationCommandOptionString {
	return &ApplicationCommandOptionString{
		Type: discord.ApplicationCommandOptionTypeString,
	}
}

func (o *ApplicationCommandOptionString) SetName(name string) *ApplicationCommandOptionString {
	o.Name = name
	return o
}

func (o *ApplicationCommandOptionString) SetNameLocalizations(localizations map[discord.Locale]string) *ApplicationCommandOptionString {
	o.NameLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionString) SetDescription(description string) *ApplicationCommandOptionString {
	o.Description = description
	return o
}

func (o *ApplicationCommandOptionString) SetDescriptionLocalizations(localizations map[discord.Locale]string) *ApplicationCommandOptionString {
	o.DescriptionLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionString) SetRequired(required bool) *ApplicationCommandOptionString {
	o.Required = required
	return o
}

func (o *ApplicationCommandOptionString) SetMinLength(minValue int64) *ApplicationCommandOptionString {
	o.MinLength = &minValue
	return o
}

func (o *ApplicationCommandOptionString) SetMaxLength(maxValue int64) *ApplicationCommandOptionString {
	o.MaxLength = &maxValue
	return o
}

func (o *ApplicationCommandOptionString) SetAutocomplete(autocomplete bool) *ApplicationCommandOptionString {
	o.Autocomplete = autocomplete
	return o
}

func (o *ApplicationCommandOptionString) SetChoices(choices []ApplicationCommandOptionChoice[string]) *ApplicationCommandOptionString {
	o.Choices = choices
	return o
}

func (o *ApplicationCommandOptionString) Build() ApplicationCommandOptionString {
	return *o
}
