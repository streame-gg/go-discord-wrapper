package components

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type MediaGalleryComponent struct {
	Type  common.ComponentType `json:"type"`
	ID    *int                 `json:"id,omitempty"`
	Items *[]MediaGalleryItem  `json:"items"`
}

func (m *MediaGalleryComponent) UnmarshalJSON(data []byte) error {
	type Alias MediaGalleryComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*m = MediaGalleryComponent(*raw.Alias)
	return nil
}

func (m *MediaGalleryComponent) MarshalJSON() ([]byte, error) {
	m.Type = common.ComponentTypeMediaGallery
	type Alias MediaGalleryComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

func (m *MediaGalleryComponent) GetType() common.ComponentType {
	return common.ComponentTypeMediaGallery
}

func (m *MediaGalleryComponent) IsAnyContainerComponent() {

}

type MediaGalleryItem struct {
	Media       *UnfurledMediaItem `json:"media"`
	Description string             `json:"description,omitempty"`
	Spoiler     bool               `json:"spoiler,omitempty"`
}
