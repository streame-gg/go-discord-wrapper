package components

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type ActionRow struct {
	Type       common.ComponentType  `json:"type"`
	ID         *int                  `json:"id"`
	Components []common.AnyComponent `json:"components"`
}

func (a *ActionRow) GetType() common.ComponentType {
	return common.ComponentTypeActionRow
}

func (a *ActionRow) MarshalJSON() ([]byte, error) {
	a.Type = common.ComponentTypeActionRow
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
			Type common.ComponentType `json:"type"`
		}

		if err := json.Unmarshal(c, &probe); err != nil {
			return err
		}

		switch probe.Type {
		case common.ComponentTypeButton:
			var b *ButtonComponent
			if err := json.Unmarshal(c, &b); err != nil {
				return err
			}
			a.Components = append(a.Components, b)
		}
	}

	return nil
}
