package components

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#file-upload
type FileUpload struct {
	Type      discord.ComponentType `json:"type"`
	ID        *int                  `json:"id,omitempty"`
	CustomID  string                `json:"custom_id"`
	Required  bool                  `json:"required"`
	MinValues *int                  `json:"min_values,omitempty"`
	MaxValues *int                  `json:"max_values,omitempty"`
}

func (f *FileUpload) MarshalJSON() ([]byte, error) {
	type Alias FileUpload
	return json.Marshal(struct {
		Alias
		Type discord.ComponentType `json:"type"`
	}{
		Alias: Alias(*f),
		Type:  discord.ComponentTypeFileUpload,
	})
}

func (f *FileUpload) UnmarshalJSON(data []byte) error {
	type Alias FileUpload
	var raw struct {
		*Alias
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Alias == nil {
		return nil
	}
	*f = FileUpload(*raw.Alias)
	return nil
}

func (f *FileUpload) GetType() discord.ComponentType {
	return discord.ComponentTypeFileUpload
}

func (f *FileUpload) IsAnyLabelComponent() {

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
