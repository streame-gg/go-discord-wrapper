package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#media-gallery
type MediaGallery struct {
	Type  discord.ComponentType `json:"type"`
	ID    *int                  `json:"id,omitempty"`
	Items []MediaGalleryItem    `json:"items"`
}

func (m *MediaGallery) UnmarshalJSON(data []byte) error {
	type Alias MediaGallery
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*m = MediaGallery(*raw.Alias)
	return nil
}

func (m *MediaGallery) MarshalJSON() ([]byte, error) {
	type Alias MediaGallery
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*m),
		Type:  discord.ComponentTypeMediaGallery,
	})
}

func (m *MediaGallery) GetType() discord.ComponentType {
	return discord.ComponentTypeMediaGallery
}

func (m *MediaGallery) IsAnyContainerComponent() {

}

// https://docs.discord.com/developers/components/reference#media-gallery-media-gallery-item-structure
type MediaGalleryItem struct {
	Media       UnfurledMediaItem `json:"media"`
	Description *string           `json:"description,omitempty"`
	Spoiler     bool              `json:"spoiler"`
}
