package discord

import (
	"encoding/json"
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionAttachment struct {
	Type                     ApplicationCommandOptionType `json:"type"`
	Name                     string                       `json:"name"`
	NameLocalizations        map[Locale]string            `json:"name_localizations,omitempty"`
	Description              string                       `json:"description"`
	DescriptionLocalizations map[Locale]string            `json:"description_localizations,omitempty"`
	Required                 bool                         `json:"required"`
}

func (o *ApplicationCommandOptionAttachment) ApplicationCommandOptionType() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeAttachment
}

func (o *ApplicationCommandOptionAttachment) MarshalJSON() ([]byte, error) {
	type Alias ApplicationCommandOptionAttachment
	return json.Marshal(struct {
		Alias
		Type ApplicationCommandOptionType `json:"type"`
	}{
		Alias: Alias(*o),
		Type:  o.ApplicationCommandOptionType(),
	})
}

func (o *ApplicationCommandOptionAttachment) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionAttachment
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

func NewAttachmentOptionBuilder() *ApplicationCommandOptionAttachment {
	return &ApplicationCommandOptionAttachment{
		Type: ApplicationCommandOptionTypeAttachment,
	}
}

func (o *ApplicationCommandOptionAttachment) SetName(name string) *ApplicationCommandOptionAttachment {
	o.Name = name
	return o
}

func (o *ApplicationCommandOptionAttachment) SetNameLocalizations(localizations map[Locale]string) *ApplicationCommandOptionAttachment {
	o.NameLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionAttachment) SetDescription(description string) *ApplicationCommandOptionAttachment {
	o.Description = description
	return o
}

func (o *ApplicationCommandOptionAttachment) SetDescriptionLocalizations(localizations map[Locale]string) *ApplicationCommandOptionAttachment {
	o.DescriptionLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionAttachment) SetRequired(required bool) *ApplicationCommandOptionAttachment {
	o.Required = required
	return o
}

func (o *ApplicationCommandOptionAttachment) Build() ApplicationCommandOptionAttachment {
	return *o
}
