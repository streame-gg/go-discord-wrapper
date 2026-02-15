package commands

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type ApplicationCommandOptionString struct {
	Type                     common.ApplicationCommandOptionType      `json:"type"`
	Name                     string                                   `json:"name"`
	NameLocalizations        map[common.Locale]string                 `json:"name_localizations,omitempty"`
	Description              string                                   `json:"description"`
	DescriptionLocalizations map[common.Locale]string                 `json:"description_localizations,omitempty"`
	Required                 *bool                                    `json:"required,omitempty"`
	Choices                  []ApplicationCommandOptionChoice[string] `json:"choices,omitempty"`
	Autocomplete             *bool                                    `json:"autocomplete,omitempty"`
}

func (o *ApplicationCommandOptionString) ApplicationCommandOptionType() common.ApplicationCommandOptionType {
	return common.ApplicationCommandOptionTypeString
}

func (o *ApplicationCommandOptionString) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionString
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
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
