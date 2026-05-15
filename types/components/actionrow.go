package components

import (
	"encoding/json"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type ActionRow struct {
	Type       discord.ComponentType  `json:"type"`
	ID         *int                   `json:"id"`
	Components []discord.AnyComponent `json:"components"`
}

func (a *ActionRow) GetType() discord.ComponentType {
	return discord.ComponentTypeActionRow
}

func (a *ActionRow) MarshalJSON() ([]byte, error) {
	a.Type = discord.ComponentTypeActionRow
	type Alias ActionRow
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(a),
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
		}
	}

	return nil
}
