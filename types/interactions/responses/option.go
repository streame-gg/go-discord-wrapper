package responses

import (
	"encoding/json"
	"fmt"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-object-application-command-interaction-data-option-structure
type ApplicationCommandInteractionDataOption[T string | int | bool | interface{}] struct {
	Name    string                                                 `json:"name"`
	Type    discord.ApplicationCommandOptionType                   `json:"type"`
	Value   *T                                                     `json:"value,omitempty"`
	Options []ApplicationCommandInteractionDataOption[interface{}] `json:"options,omitempty"`
	Focused *bool                                                  `json:"focused,omitempty"`
}

func (t *ApplicationCommandInteractionDataOption[T]) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandInteractionDataOption[T]
	// Capture the raw value bytes so each option type can be decoded with the
	// right Go type. Discord sends string and boolean options as JSON strings
	// and booleans, which cannot decode into json.Number — so the value is held
	// as RawMessage here and interpreted by t.Type below.
	raw := &struct {
		*Alias
		Value json.RawMessage `json:"value,omitempty"`
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

	if len(raw.Value) == 0 || string(raw.Value) == "null" {
		return nil
	}

	// asNumber decodes the raw value as a json.Number so integer options keep
	// full precision (values > 2^53 would be lossy via float64).
	asNumber := func() (json.Number, error) {
		var n json.Number
		if err := json.Unmarshal(raw.Value, &n); err != nil {
			return "", err
		}
		return n, nil
	}

	switch t.Type {
	case discord.ApplicationCommandOptionTypeString:
		var s string
		if err := json.Unmarshal(raw.Value, &s); err != nil {
			return fmt.Errorf("option %q: string value: %w", t.Name, err)
		}
		v := any(s).(T)
		t.Value = &v
	case discord.ApplicationCommandOptionTypeInteger:
		num, err := asNumber()
		if err != nil {
			return fmt.Errorf("option %q: integer value %s: %w", t.Name, raw.Value, err)
		}
		n, err := num.Int64()
		if err != nil {
			return fmt.Errorf("option %q: integer value %q: %w", t.Name, num, err)
		}
		v := any(int(n)).(T)
		t.Value = &v
	case discord.ApplicationCommandOptionTypeNumber:
		num, err := asNumber()
		if err != nil {
			return fmt.Errorf("option %q: number value %s: %w", t.Name, raw.Value, err)
		}
		f, err := num.Float64()
		if err != nil {
			return fmt.Errorf("option %q: number value %q: %w", t.Name, num, err)
		}
		v := any(f).(T)
		t.Value = &v
	case discord.ApplicationCommandOptionTypeBoolean:
		var b bool
		if err := json.Unmarshal(raw.Value, &b); err != nil {
			return fmt.Errorf("option %q: boolean value %s: %w", t.Name, raw.Value, err)
		}
		v := any(b).(T)
		t.Value = &v
	default:
		// User/Channel/Role/Mentionable values are snowflake IDs sent as JSON
		// strings; decode them as a string.
		var s string
		if err := json.Unmarshal(raw.Value, &s); err != nil {
			return fmt.Errorf("option %q: value %s: %w", t.Name, raw.Value, err)
		}
		v, ok := any(s).(T)
		if !ok {
			return fmt.Errorf("unexpected value type for option type %v", t.Type)
		}
		t.Value = &v
	}

	return nil
}
