package components

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type FileUploadComponent struct {
	Type      common.ComponentType `json:"type"`
	ID        *int                 `json:"id,omitempty"`
	CustomID  string               `json:"custom_id"`
	Required  *bool                `json:"required,omitempty"`
	MinValues *int                 `json:"min_values,omitempty"`
	MaxValues *int                 `json:"max_values,omitempty"`
}

func (f *FileUploadComponent) MarshalJSON() ([]byte, error) {
	f.Type = common.ComponentTypeFileUpload
	type Alias FileUploadComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(f),
	})
}

func (f *FileUploadComponent) UnmarshalJSON(data []byte) error {
	type Alias FileUploadComponent
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*f = FileUploadComponent(*raw.Alias)
	return nil
}

func (f *FileUploadComponent) GetType() common.ComponentType {
	return common.ComponentTypeFileUpload
}

func (f *FileUploadComponent) IsAnyLabelComponent() {

}

type FileUploadComponentInteractionResponse struct {
	Type     common.ComponentType `json:"type"`
	ID       *int                 `json:"id,omitempty"`
	CustomID string               `json:"custom_id"`
	Values   []common.Snowflake   `json:"values"`
}

func (f *FileUploadComponentInteractionResponse) IsInteractionResponseDataComponent() {

}

func (f *FileUploadComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	f.Type = common.ComponentTypeFileUpload

	type Alias FileUploadComponentInteractionResponse

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(f),
	})
}

func (f *FileUploadComponentInteractionResponse) UnmarshalJSON(data []byte) error {
	type Alias FileUploadComponentInteractionResponse
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*f = FileUploadComponentInteractionResponse(*raw.Alias)
	return nil
}
