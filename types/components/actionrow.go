package components

import (
	"encoding/json"
	"fmt"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#action-row
type ActionRow struct {
	Type       discord.ComponentType  `json:"type"`
	ID         *int                   `json:"id,omitempty"`
	Components []discord.AnyComponent `json:"components"`
}

func (a *ActionRow) GetType() discord.ComponentType {
	return discord.ComponentTypeActionRow
}

func (a *ActionRow) MarshalJSON() ([]byte, error) {
	type Alias ActionRow
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*a),
		Type:  discord.ComponentTypeActionRow,
	})
}

func (a *ActionRow) IsAnyContainerComponent() {

}

func (a *ActionRow) UnmarshalJSON(data []byte) error {
	type Alias ActionRow

	var raw struct {
		Alias
		Components []json.RawMessage `json:"components"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*a = ActionRow(raw.Alias)

	for _, c := range raw.Components {
		var probe struct {
			Type discord.ComponentType `json:"type"`
		}

		if err := json.Unmarshal(c, &probe); err != nil {
			return err
		}

		switch probe.Type {
		case discord.ComponentTypeButton:
			var b *ButtonComponent
			if err := json.Unmarshal(c, &b); err != nil {
				return err
			}
			a.Components = append(a.Components, b)
		case discord.ComponentTypeStringSelect:
			var s *StringSelectMenuComponent
			if err := json.Unmarshal(c, &s); err != nil {
				return err
			}
			a.Components = append(a.Components, s)
		case discord.ComponentTypeUserSelect:
			var u *UserSelectMenuComponent
			if err := json.Unmarshal(c, &u); err != nil {
				return err
			}
			a.Components = append(a.Components, u)
		case discord.ComponentTypeRoleSelect:
			var r *RoleSelectMenuComponent
			if err := json.Unmarshal(c, &r); err != nil {
				return err
			}
			a.Components = append(a.Components, r)
		case discord.ComponentTypeMentionableSelect:
			var m *MentionableSelectMenuComponent
			if err := json.Unmarshal(c, &m); err != nil {
				return err
			}
			a.Components = append(a.Components, m)
		case discord.ComponentTypeChannelSelect:
			var ch *ChannelSelectMenuComponent
			if err := json.Unmarshal(c, &ch); err != nil {
				return err
			}
			a.Components = append(a.Components, ch)
		case discord.ComponentTypeTextInput:
			var t *TextInputComponent
			if err := json.Unmarshal(c, &t); err != nil {
				return err
			}
			a.Components = append(a.Components, t)
		default:
			return fmt.Errorf("unknown action row component type: %d", probe.Type)
		}
	}

	return nil
}
