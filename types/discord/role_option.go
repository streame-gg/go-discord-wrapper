package discord

import (
	"encoding/json"
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionRole struct {
	Type                     ApplicationCommandOptionType `json:"type"`
	Name                     string                       `json:"name"`
	NameLocalizations        map[Locale]string            `json:"name_localizations,omitempty"`
	Description              string                       `json:"description"`
	DescriptionLocalizations map[Locale]string            `json:"description_localizations,omitempty"`
	Required                 bool                         `json:"required"`
}

func (o *ApplicationCommandOptionRole) ApplicationCommandOptionType() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeRole
}

func (o *ApplicationCommandOptionRole) MarshalJSON() ([]byte, error) {
	type Alias ApplicationCommandOptionRole
	return json.Marshal(struct {
		Alias
		Type ApplicationCommandOptionType `json:"type"`
	}{
		Alias: Alias(*o),
		Type:  o.ApplicationCommandOptionType(),
	})
}

func (o *ApplicationCommandOptionRole) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionRole
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

func NewRoleOptionBuilder() *ApplicationCommandOptionRole {
	return &ApplicationCommandOptionRole{
		Type: ApplicationCommandOptionTypeRole,
	}
}

func (o *ApplicationCommandOptionRole) SetName(name string) *ApplicationCommandOptionRole {
	o.Name = name
	return o
}

func (o *ApplicationCommandOptionRole) SetNameLocalizations(localizations map[Locale]string) *ApplicationCommandOptionRole {
	o.NameLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionRole) SetDescription(description string) *ApplicationCommandOptionRole {
	o.Description = description
	return o
}

func (o *ApplicationCommandOptionRole) SetDescriptionLocalizations(localizations map[Locale]string) *ApplicationCommandOptionRole {
	o.DescriptionLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionRole) SetRequired(required bool) *ApplicationCommandOptionRole {
	o.Required = required
	return o
}

func (o *ApplicationCommandOptionRole) Build() ApplicationCommandOptionRole {
	return *o
}
