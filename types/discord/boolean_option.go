package discord

import (
	"encoding/json"
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionBoolean struct {
	Type                     ApplicationCommandOptionType `json:"type"`
	Name                     string                       `json:"name"`
	NameLocalizations        map[Locale]string            `json:"name_localizations,omitempty"`
	Description              string                       `json:"description"`
	DescriptionLocalizations map[Locale]string            `json:"description_localizations,omitempty"`
	Required                 bool                         `json:"required"`
}

func (o *ApplicationCommandOptionBoolean) ApplicationCommandOptionType() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeBoolean
}

func (o *ApplicationCommandOptionBoolean) MarshalJSON() ([]byte, error) {
	type Alias ApplicationCommandOptionBoolean
	return json.Marshal(struct {
		Alias
		Type ApplicationCommandOptionType `json:"type"`
	}{
		Alias: Alias(*o),
		Type:  o.ApplicationCommandOptionType(),
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

func NewBooleanOptionBuilder() *ApplicationCommandOptionBoolean {
	return &ApplicationCommandOptionBoolean{
		Type: ApplicationCommandOptionTypeBoolean,
	}
}

func (o *ApplicationCommandOptionBoolean) SetName(name string) *ApplicationCommandOptionBoolean {
	o.Name = name
	return o
}

func (o *ApplicationCommandOptionBoolean) SetNameLocalizations(localizations map[Locale]string) *ApplicationCommandOptionBoolean {
	o.NameLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionBoolean) SetDescription(description string) *ApplicationCommandOptionBoolean {
	o.Description = description
	return o
}

func (o *ApplicationCommandOptionBoolean) SetDescriptionLocalizations(localizations map[Locale]string) *ApplicationCommandOptionBoolean {
	o.DescriptionLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionBoolean) SetRequired(required bool) *ApplicationCommandOptionBoolean {
	o.Required = required
	return o
}

func (o *ApplicationCommandOptionBoolean) Build() ApplicationCommandOptionBoolean {
	return *o
}
