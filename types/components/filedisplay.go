package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#file
type File struct {
	Type    discord.ComponentType `json:"type"`
	ID      *int                  `json:"id,omitempty"`
	Spoiler bool                  `json:"spoiler"`
	Name    string                `json:"name,omitempty"`
	Size    *int                  `json:"size,omitempty"`
	File    *UnfurledMediaItem    `json:"file,omitempty"`
}

func (f *File) UnmarshalJSON(data []byte) error {
	type Alias File
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*f = File(*raw.Alias)
	return nil
}

func (f *File) MarshalJSON() ([]byte, error) {
	type Alias File
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*f),
		Type:  discord.ComponentTypeFile,
	})
}

func (f *File) GetType() discord.ComponentType {
	return discord.ComponentTypeFile
}

func (f *File) IsAnyContainerComponent() {

}
