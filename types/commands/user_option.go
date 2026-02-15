package commands

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type ApplicationCommandOptionUser struct {
	Type                     common.ApplicationCommandOptionType `json:"type"`
	Name                     string                              `json:"name"`
	NameLocalizations        map[common.Locale]string            `json:"name_localizations,omitempty"`
	Description              string                              `json:"description"`
	DescriptionLocalizations map[common.Locale]string            `json:"description_localizations,omitempty"`
	Required                 *bool                               `json:"required,omitempty"`
}

func (o *ApplicationCommandOptionUser) ApplicationCommandOptionType() common.ApplicationCommandOptionType {
	return common.ApplicationCommandOptionTypeUser
}

func (o *ApplicationCommandOptionUser) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionUser
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
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
