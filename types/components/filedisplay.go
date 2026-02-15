package components

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type FileComponent struct {
	Type    common.ComponentType `json:"type"`
	ID      *int                 `json:"id,omitempty"`
	Spoiler bool                 `json:"spoiler,omitempty"`
	Name    string               `json:"name,omitempty"`
	Size    int                  `json:"size,omitempty"`
	File    *UnfurledMediaItem   `json:"file,omitempty"`
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
	f.Type = common.ComponentTypeFileDisplay
	type Alias FileComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(f),
	})
}

func (f *FileComponent) GetType() common.ComponentType {
	return common.ComponentTypeFileDisplay
}

func (f *FileComponent) IsAnyContainerComponent() {

}
