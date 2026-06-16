package discord

import (
	"encoding/json"
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionNumber struct {
	Type                     ApplicationCommandOptionType              `json:"type"`
	Name                     string                                    `json:"name"`
	NameLocalizations        map[Locale]string                         `json:"name_localizations,omitempty"`
	Description              string                                    `json:"description"`
	DescriptionLocalizations map[Locale]string                         `json:"description_localizations,omitempty"`
	Required                 bool                                      `json:"required"`
	Choices                  []ApplicationCommandOptionChoice[float64] `json:"choices,omitempty"`
	MinValue                 *float64                                  `json:"min_value,omitempty"`
	MaxValue                 *float64                                  `json:"max_value,omitempty"`
	Autocomplete             bool                                      `json:"autocomplete"`
}

func (o *ApplicationCommandOptionNumber) ApplicationCommandOptionType() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeNumber
}

func (o *ApplicationCommandOptionNumber) MarshalJSON() ([]byte, error) {
	type Alias ApplicationCommandOptionNumber
	return json.Marshal(struct {
		Alias
		Type ApplicationCommandOptionType `json:"type"`
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

func NewNumberOptionBuilder() *ApplicationCommandOptionNumber {
	return &ApplicationCommandOptionNumber{
		Type: ApplicationCommandOptionTypeNumber,
	}
}

func (o *ApplicationCommandOptionNumber) SetName(name string) *ApplicationCommandOptionNumber {
	o.Name = name
	return o
}

func (o *ApplicationCommandOptionNumber) SetNameLocalizations(localizations map[Locale]string) *ApplicationCommandOptionNumber {
	o.NameLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionNumber) SetDescription(description string) *ApplicationCommandOptionNumber {
	o.Description = description
	return o
}

func (o *ApplicationCommandOptionNumber) SetDescriptionLocalizations(localizations map[Locale]string) *ApplicationCommandOptionNumber {
	o.DescriptionLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionNumber) SetRequired(required bool) *ApplicationCommandOptionNumber {
	o.Required = required
	return o
}

func (o *ApplicationCommandOptionNumber) SetMinValue(minValue float64) *ApplicationCommandOptionNumber {
	o.MinValue = &minValue
	return o
}

func (o *ApplicationCommandOptionNumber) SetMaxValue(maxValue float64) *ApplicationCommandOptionNumber {
	o.MaxValue = &maxValue
	return o
}

func (o *ApplicationCommandOptionNumber) SetAutocomplete(autocomplete bool) *ApplicationCommandOptionNumber {
	o.Autocomplete = autocomplete
	return o
}

func (o *ApplicationCommandOptionNumber) SetChoices(choices []ApplicationCommandOptionChoice[float64]) *ApplicationCommandOptionNumber {
	o.Choices = choices
	return o
}

func (o *ApplicationCommandOptionNumber) Build() ApplicationCommandOptionNumber {
	return *o
}
