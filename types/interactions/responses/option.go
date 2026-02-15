package responses

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type ApplicationCommandInteractionDataOption[T string | int | bool | interface{}] struct {
	Name    string                                                 `json:"name"`
	Type    common.ApplicationCommandOptionType                    `json:"type"`
	Value   *T                                                     `json:"value"`
	Options []ApplicationCommandInteractionDataOption[interface{}] `json:"options,omitempty"`
	Focused *bool                                                  `json:"focused,omitempty"`
}

func (t *ApplicationCommandInteractionDataOption[T]) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandInteractionDataOption[T]
	raw := &struct {
		*Alias
		Value interface{} `json:"value,omitempty"`
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	t.Name = raw.Name
	t.Type = raw.Type
	t.Options = raw.Options
	t.Focused = raw.Focused

	if raw.Value != nil {
		switch t.Type {
		case common.ApplicationCommandOptionTypeString:
			if strVal, ok := raw.Value.(string); ok {
				var v T
				v = any(strVal).(T)
				t.Value = &v
			}
		case common.ApplicationCommandOptionTypeInteger:
			if intVal, ok := raw.Value.(int); ok {
				var v T
				v = any(intVal).(T)
				t.Value = &v
			}
		case common.ApplicationCommandOptionTypeBoolean:
			if boolVal, ok := raw.Value.(bool); ok {
				var v T
				v = any(boolVal).(T)
				t.Value = &v
			}
		default:
			var v T
			v = raw.Value.(T)
			t.Value = &v
		}
	}

	return nil
}
