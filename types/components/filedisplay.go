package components

import (
	"encoding/json"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type FileComponent struct {
	Type    discord.ComponentType `json:"type"`
	ID      *int                  `json:"id,omitempty"`
	Spoiler bool                  `json:"spoiler,omitempty"`
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

	*f = FileComponent(*raw.Alias)
	return nil
}

func (f *FileComponent) MarshalJSON() ([]byte, error) {
	f.Type = discord.ComponentTypeFileDisplay
	type Alias FileComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(f),
	})
}

func (f *FileComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeFileDisplay
}

func (f *FileComponent) IsAnyContainerComponent() {

}
