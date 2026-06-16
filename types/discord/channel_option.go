package discord

import (
	"encoding/json"
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionChannel struct {
	Type                     ApplicationCommandOptionType `json:"type"`
	Name                     string                       `json:"name"`
	NameLocalizations        map[Locale]string            `json:"name_localizations,omitempty"`
	Description              string                       `json:"description"`
	DescriptionLocalizations map[Locale]string            `json:"description_localizations,omitempty"`
	Required                 bool                         `json:"required"`
	ChannelTypes             []ChannelType                `json:"channel_types,omitempty"`
}

func (o *ApplicationCommandOptionChannel) ApplicationCommandOptionType() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeChannel
}

func (o *ApplicationCommandOptionChannel) MarshalJSON() ([]byte, error) {
	type Alias ApplicationCommandOptionChannel
	return json.Marshal(struct {
		Alias
		Type ApplicationCommandOptionType `json:"type"`
	}{
		Alias: Alias(*o),
		Type:  o.ApplicationCommandOptionType(),
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

func NewChannelOptionBuilder() *ApplicationCommandOptionChannel {
	return &ApplicationCommandOptionChannel{
		Type: ApplicationCommandOptionTypeChannel,
	}
}

func (o *ApplicationCommandOptionChannel) SetName(name string) *ApplicationCommandOptionChannel {
	o.Name = name
	return o
}

func (o *ApplicationCommandOptionChannel) SetNameLocalizations(localizations map[Locale]string) *ApplicationCommandOptionChannel {
	o.NameLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionChannel) SetDescription(description string) *ApplicationCommandOptionChannel {
	o.Description = description
	return o
}

func (o *ApplicationCommandOptionChannel) SetDescriptionLocalizations(localizations map[Locale]string) *ApplicationCommandOptionChannel {
	o.DescriptionLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionChannel) SetRequired(required bool) *ApplicationCommandOptionChannel {
	o.Required = required
	return o
}

func (o *ApplicationCommandOptionChannel) SetChannelTypes(channelTypes []ChannelType) *ApplicationCommandOptionChannel {
	o.ChannelTypes = channelTypes
	return o
}

func (o *ApplicationCommandOptionChannel) Build() ApplicationCommandOptionChannel {
	return *o
}
