package responses

import (
	"encoding/json"
	"fmt"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type ApplicationCommandInteractionDataOption[T string | int | bool | interface{}] struct {
	Name    string                                                 `json:"name"`
	Type    discord.ApplicationCommandOptionType                   `json:"type"`
	Value   *T                                                     `json:"value,omitempty"`
	Options []ApplicationCommandInteractionDataOption[interface{}] `json:"options,omitempty"`
	Focused *bool                                                  `json:"focused,omitempty"`
}

func (t *ApplicationCommandInteractionDataOption[T]) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandInteractionDataOption[T]
	// Use json.Number for the raw value field so integer options are not
	// decoded via float64 (which loses precision for values > 2^53).
	raw := &struct {
		*Alias
		Value json.Number `json:"value,omitempty"`
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

	if raw.Value != "" {
		switch t.Type {
		case discord.ApplicationCommandOptionTypeString:
			v := any(raw.Value.String()).(T)
			t.Value = &v
		case discord.ApplicationCommandOptionTypeInteger:
			n, err := raw.Value.Int64()
			if err != nil {
				return fmt.Errorf("option %q: integer value %q: %w", t.Name, raw.Value, err)
			}
			v := any(int(n)).(T)
			t.Value = &v
		case discord.ApplicationCommandOptionTypeNumber:
			f, err := raw.Value.Float64()
			if err != nil {
				return fmt.Errorf("option %q: number value %q: %w", t.Name, raw.Value, err)
			}
			v := any(f).(T)
			t.Value = &v
		case discord.ApplicationCommandOptionTypeBoolean:
			switch raw.Value.String() {
			case "true":
				v := any(true).(T)
				t.Value = &v
			case "false":
				v := any(false).(T)
				t.Value = &v
			default:
				return fmt.Errorf("option %q: unexpected boolean value %q", t.Name, raw.Value)
			}
		default:
			v, ok := any(raw.Value.String()).(T)
			if !ok {
				return fmt.Errorf("unexpected value type for option type %v", t.Type)
			}
			t.Value = &v
		}
	}

	return nil
}
