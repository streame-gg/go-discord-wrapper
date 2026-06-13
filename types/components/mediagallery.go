package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#media-gallery
type MediaGalleryComponent struct {
	Type  discord.ComponentType `json:"type"`
	ID    *int                  `json:"id,omitempty"`
	Items []MediaGalleryItem    `json:"items"`
}

func (m *MediaGalleryComponent) UnmarshalJSON(data []byte) error {
	type Alias MediaGalleryComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*m = MediaGalleryComponent(*raw.Alias)
	return nil
}

func (m *MediaGalleryComponent) MarshalJSON() ([]byte, error) {
	type Alias MediaGalleryComponent
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*m),
		Type:  discord.ComponentTypeMediaGallery,
	})
}

func (m *MediaGalleryComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeMediaGallery
}

func (m *MediaGalleryComponent) IsAnyContainerComponent() {

}

// https://docs.discord.com/developers/components/reference#media-gallery-media-gallery-item-structure
type MediaGalleryItem struct {
	Media       UnfurledMediaItem `json:"media"`
	Description *string           `json:"description,omitempty"`
	Spoiler     bool              `json:"spoiler"`
}
