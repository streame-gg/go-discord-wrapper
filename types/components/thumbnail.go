package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#thumbnail
type Thumbnail struct {
	Type        discord.ComponentType `json:"type"`
	ID          *int                  `json:"id,omitempty"`
	Description *string               `json:"description"`
	Spoiler     bool                  `json:"spoiler"`
	Media       UnfurledMediaItem     `json:"media"`
}

func (t *Thumbnail) UnmarshalJSON(data []byte) error {
	type Alias Thumbnail
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*t = Thumbnail(*raw.Alias)
	return nil
}

func (t *Thumbnail) MarshalJSON() ([]byte, error) {
	type Alias Thumbnail
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*t),
		Type:  discord.ComponentTypeThumbnail,
	})
}

func (t *Thumbnail) GetType() discord.ComponentType {
	return discord.ComponentTypeThumbnail
}

func (t *Thumbnail) IsAnySectionAccessory() {}
