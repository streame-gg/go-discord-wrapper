package commands

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionAttachment struct {
	Type                     discord.ApplicationCommandOptionType `json:"type"`
	Name                     string                               `json:"name"`
	NameLocalizations        map[discord.Locale]string            `json:"name_localizations,omitempty"`
	Description              string                               `json:"description"`
	DescriptionLocalizations map[discord.Locale]string            `json:"description_localizations,omitempty"`
	Required                 bool                                 `json:"required"`
}

func (o *ApplicationCommandOptionAttachment) ApplicationCommandOptionType() discord.ApplicationCommandOptionType {
	return discord.ApplicationCommandOptionTypeAttachment
}

func (o *ApplicationCommandOptionAttachment) MarshalJSON() ([]byte, error) {
	type Alias ApplicationCommandOptionAttachment
	return json.Marshal(struct {
		Alias
		Type discord.ApplicationCommandOptionType `json:"type"`
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
