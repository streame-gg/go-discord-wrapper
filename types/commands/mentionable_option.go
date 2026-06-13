package commands

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionMentionable struct {
	Type                     discord.ApplicationCommandOptionType `json:"type"`
	Name                     string                               `json:"name"`
	NameLocalizations        map[discord.Locale]string            `json:"name_localizations,omitempty"`
	Description              string                               `json:"description"`
	DescriptionLocalizations map[discord.Locale]string            `json:"description_localizations,omitempty"`
	Required                 bool                                 `json:"required"`
}

func (o *ApplicationCommandOptionMentionable) ApplicationCommandOptionType() discord.ApplicationCommandOptionType {
	return discord.ApplicationCommandOptionTypeMentionable
}

func (o *ApplicationCommandOptionMentionable) MarshalJSON() ([]byte, error) {
	type Alias ApplicationCommandOptionMentionable
	return json.Marshal(struct {
		Alias
		Type discord.ApplicationCommandOptionType `json:"type"`
	}{
		Alias: Alias(*o),
		Type:  o.ApplicationCommandOptionType(),
	})
}

func (o *ApplicationCommandOptionMentionable) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionMentionable
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

func NewMentionableOptionBuilder() *ApplicationCommandOptionMentionable {
	return &ApplicationCommandOptionMentionable{
		Type: discord.ApplicationCommandOptionTypeMentionable,
	}
}

func (o *ApplicationCommandOptionMentionable) SetName(name string) *ApplicationCommandOptionMentionable {
	o.Name = name
	return o
}

func (o *ApplicationCommandOptionMentionable) SetNameLocalizations(localizations map[discord.Locale]string) *ApplicationCommandOptionMentionable {
	o.NameLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionMentionable) SetDescription(description string) *ApplicationCommandOptionMentionable {
	o.Description = description
	return o
}

func (o *ApplicationCommandOptionMentionable) SetDescriptionLocalizations(localizations map[discord.Locale]string) *ApplicationCommandOptionMentionable {
	o.DescriptionLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionMentionable) SetRequired(required bool) *ApplicationCommandOptionMentionable {
	o.Required = required
	return o
}

func (o *ApplicationCommandOptionMentionable) Build() ApplicationCommandOptionMentionable {
	return *o
}
