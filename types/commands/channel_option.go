package commands

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type ApplicationCommandOptionChannel struct {
	Type                     common.ApplicationCommandOptionType `json:"type"`
	Name                     string                              `json:"name"`
	NameLocalizations        map[common.Locale]string            `json:"name_localizations,omitempty"`
	Description              string                              `json:"description"`
	DescriptionLocalizations map[common.Locale]string            `json:"description_localizations,omitempty"`
	Required                 *bool                               `json:"required,omitempty"`
	ChannelTypes             []common.ChannelType                `json:"channel_types,omitempty"`
}

func (o *ApplicationCommandOptionChannel) ApplicationCommandOptionType() common.ApplicationCommandOptionType {
	return common.ApplicationCommandOptionTypeChannel
}

func (o *ApplicationCommandOptionChannel) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionChannel
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
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
