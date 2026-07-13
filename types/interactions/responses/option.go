package responses

import (
	"encoding/json"
	"fmt"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-object-application-command-interaction-data-option-structure
type ApplicationCommandInteractionDataOption[T string | discord.Snowflake | int | bool | float64 | interface{}] struct {
	Name    string                                                 `json:"name"`
	Type    discord.ApplicationCommandOptionType                   `json:"type"`
	Value   *T                                                     `json:"value,omitempty"`
	Options []ApplicationCommandInteractionDataOption[interface{}] `json:"options,omitempty"`
	Focused *bool                                                  `json:"focused,omitempty"`
}

func (t *ApplicationCommandInteractionDataOption[T]) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandInteractionDataOption[T]
	raw := &struct {
		*Alias
		Value json.RawMessage `json:"value,omitempty"`
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if raw == nil {
		return nil
	}

	t.Name = raw.Name
	t.Type = raw.Type
	t.Options = raw.Options
	t.Focused = raw.Focused

	if len(raw.Value) == 0 || string(raw.Value) == "null" {
		return nil
	}

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
	case discord.ApplicationCommandOptionTypeUser, discord.ApplicationCommandOptionTypeAttachment, discord.ApplicationCommandOptionTypeChannel, discord.ApplicationCommandOptionTypeRole, discord.ApplicationCommandOptionTypeMentionable:
		var s discord.Snowflake
		if err := json.Unmarshal(raw.Value, &s); err != nil {
			return fmt.Errorf("option %q: value %s: %w", t.Name, raw.Value, err)
		}
		v := any(s).(T)
		t.Value = &v
	case discord.ApplicationCommandOptionTypeSubCommand, discord.ApplicationCommandOptionTypeSubCommandGroup:
	default:
		return fmt.Errorf("option %q: unknown option type", t.Name)
	}

	return nil
}
