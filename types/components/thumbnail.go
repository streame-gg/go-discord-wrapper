package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#thumbnail
type ThumbnailComponent struct {
	Type        discord.ComponentType `json:"type"`
	ID          *int                  `json:"id,omitempty"`
	Description string                `json:"description,omitempty"`
	Spoiler     bool                  `json:"spoiler"`
	Media       UnfurledMediaItem     `json:"media"`
}

func (t *ThumbnailComponent) UnmarshalJSON(data []byte) error {
	type Alias ThumbnailComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*t = ThumbnailComponent(*raw.Alias)
	return nil
}

func (t *ThumbnailComponent) MarshalJSON() ([]byte, error) {
	type Alias ThumbnailComponent
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*t),
		Type:  discord.ComponentTypeThumbnail,
	})
}

func (t *ThumbnailComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeThumbnail
}

func (t *ThumbnailComponent) IsAnySectionAccessory() {}
