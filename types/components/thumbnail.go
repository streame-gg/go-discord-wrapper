package components

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type ThumbnailComponent struct {
	Type        common.ComponentType `json:"type"`
	ID          *int                 `json:"id,omitempty"`
	Description string               `json:"description,omitempty"`
	Spoiler     bool                 `json:"spoiler,omitempty"`
	Media       *UnfurledMediaItem   `json:"media,omitempty"`
}

func (t *ThumbnailComponent) UnmarshalJSON(data []byte) error {
	type Alias ThumbnailComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*t = ThumbnailComponent(*raw.Alias)
	return nil
}

func (t *ThumbnailComponent) MarshalJSON() ([]byte, error) {
	t.Type = common.ComponentTypeThumbnail
	type Alias ThumbnailComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

func (t *ThumbnailComponent) GetType() common.ComponentType {
	return common.ComponentTypeThumbnail
}

func (t *ThumbnailComponent) IsAnySectionAccessory() {}
