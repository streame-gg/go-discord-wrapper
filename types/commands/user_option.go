package commands

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionUser struct {
	Type                     discord.ApplicationCommandOptionType `json:"type"`
	Name                     string                               `json:"name"`
	NameLocalizations        map[discord.Locale]string            `json:"name_localizations,omitempty"`
	Description              string                               `json:"description"`
	DescriptionLocalizations map[discord.Locale]string            `json:"description_localizations,omitempty"`
	Required                 bool                                 `json:"required"`
}

func (o *ApplicationCommandOptionUser) ApplicationCommandOptionType() discord.ApplicationCommandOptionType {
	return discord.ApplicationCommandOptionTypeUser
}

func (o *ApplicationCommandOptionUser) MarshalJSON() ([]byte, error) {
	type Alias ApplicationCommandOptionUser
	return json.Marshal(struct {
		Alias
		Type discord.ApplicationCommandOptionType `json:"type"`
	}{
		Alias: Alias(*o),
		Type:  o.ApplicationCommandOptionType(),
	})
}

func (o *ApplicationCommandOptionUser) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionUser
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

func NewUserOptionBuilder() *ApplicationCommandOptionUser {
	return &ApplicationCommandOptionUser{
		Type: discord.ApplicationCommandOptionTypeUser,
	}
}

func (o *ApplicationCommandOptionUser) SetName(name string) *ApplicationCommandOptionUser {
	o.Name = name
	return o
}

func (o *ApplicationCommandOptionUser) SetNameLocalizations(localizations map[discord.Locale]string) *ApplicationCommandOptionUser {
	o.NameLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionUser) SetDescription(description string) *ApplicationCommandOptionUser {
	o.Description = description
	return o
}

func (o *ApplicationCommandOptionUser) SetDescriptionLocalizations(localizations map[discord.Locale]string) *ApplicationCommandOptionUser {
	o.DescriptionLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionUser) SetRequired(required bool) *ApplicationCommandOptionUser {
	o.Required = required
	return o
}

func (o *ApplicationCommandOptionUser) Build() ApplicationCommandOptionUser {
	return *o
}
