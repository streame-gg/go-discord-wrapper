package commands

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type ApplicationCommandOptionSubCommandGroup struct {
	Type                     common.ApplicationCommandOptionType `json:"type"`
	Name                     string                              `json:"name"`
	NameLocalizations        map[common.Locale]string            `json:"name_localizations,omitempty"`
	Description              string                              `json:"description"`
	DescriptionLocalizations map[common.Locale]string            `json:"description_localizations,omitempty"`
	Options                  []AnyApplicationCommandOption       `json:"options,omitempty"`
}

func (o *ApplicationCommandOptionSubCommandGroup) ApplicationCommandOptionType() common.ApplicationCommandOptionType {
	return common.ApplicationCommandOptionTypeSubCommandGroup
}

func (o *ApplicationCommandOptionSubCommandGroup) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionSubCommandGroup
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

func (o *ApplicationCommandOptionSubCommandGroup) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionSubCommandGroup
	raw := &struct {
		*Alias
		Options []json.RawMessage `json:"options,omitempty"`
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Options != nil {
		opts, err := unmarshalOptionSlice(raw.Options)
		if err != nil {
			return err
		}
		o.Options = opts
	}

	return nil
}
