package commands

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type ApplicationCommandOptionNumber struct {
	Type                     common.ApplicationCommandOptionType      `json:"type"`
	Name                     string                                   `json:"name"`
	NameLocalizations        map[common.Locale]string                 `json:"name_localizations,omitempty"`
	Description              string                                   `json:"description"`
	DescriptionLocalizations map[common.Locale]string                 `json:"description_localizations,omitempty"`
	Required                 *bool                                    `json:"required,omitempty"`
	Choices                  []ApplicationCommandOptionChoice[string] `json:"choices,omitempty"`
	MinValue                 *int64                                   `json:"min_value,omitempty"`
	MaxValue                 *int64                                   `json:"max_value,omitempty"`
	Autocomplete             *bool                                    `json:"autocomplete,omitempty"`
}

func (o *ApplicationCommandOptionNumber) ApplicationCommandOptionType() common.ApplicationCommandOptionType {
	return common.ApplicationCommandOptionTypeNumber
}

func (o *ApplicationCommandOptionNumber) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionNumber
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
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
