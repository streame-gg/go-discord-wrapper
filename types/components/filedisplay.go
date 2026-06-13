package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#file
type FileComponent struct {
	Type    discord.ComponentType `json:"type"`
	ID      *int                  `json:"id,omitempty"`
	Spoiler bool                  `json:"spoiler"`
	Name    string                `json:"name,omitempty"`
	Size    int                   `json:"size,omitempty"`
	File    *UnfurledMediaItem    `json:"file,omitempty"`
}

func (f *FileComponent) UnmarshalJSON(data []byte) error {
	type Alias FileComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*f = FileComponent(*raw.Alias)
	return nil
}

func (f *FileComponent) MarshalJSON() ([]byte, error) {
	type Alias FileComponent
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*f),
		Type:  discord.ComponentTypeFileDisplay,
	})
}

func (f *FileComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeFileDisplay
}

func (f *FileComponent) IsAnyContainerComponent() {

}
