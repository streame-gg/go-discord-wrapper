package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#file-upload
type FileUploadComponent struct {
	Type      discord.ComponentType `json:"type"`
	ID        *int                  `json:"id,omitempty"`
	CustomID  string                `json:"custom_id"`
	Required  bool                  `json:"required"`
	MinValues *int                  `json:"min_values,omitempty"`
	MaxValues *int                  `json:"max_values,omitempty"`
}

func (f *FileUploadComponent) MarshalJSON() ([]byte, error) {
	type Alias FileUploadComponent
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*f),
		Type:  discord.ComponentTypeFileUpload,
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

	if raw.Alias == nil {
		return nil
	}
	*f = FileUploadComponent(*raw.Alias)
	return nil
}

func (f *FileUploadComponent) GetType() discord.ComponentType {
	return discord.ComponentTypeFileUpload
}

func (f *FileUploadComponent) IsAnyLabelComponent() {

}

// https://docs.discord.com/developers/components/reference#file-upload
type FileUploadComponentInteractionResponse struct {
	Type     discord.ComponentType `json:"type"`
	ID       *int                  `json:"id,omitempty"`
	CustomID string                `json:"custom_id"`
	Values   []discord.Snowflake   `json:"values"`
}

func (f *FileUploadComponentInteractionResponse) IsInteractionResponseDataComponent() {

}

func (f *FileUploadComponentInteractionResponse) MarshalJSON() ([]byte, error) {
	type Alias FileUploadComponentInteractionResponse
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*f),
		Type:  discord.ComponentTypeFileUpload,
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

	if raw.Alias == nil {
		return nil
	}
	*f = FileUploadComponentInteractionResponse(*raw.Alias)
	return nil
}
