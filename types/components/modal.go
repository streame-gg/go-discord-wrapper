package components

import "encoding/json"

type Modal struct {
	Title      string            `json:"title"`
	CustomID   string            `json:"custom_id"`
	Components *[]LabelComponent `json:"components"`
}

func (m Modal) IsInteractionResponseData() bool {
	return true
}

func (m Modal) MarshalJSON() ([]byte, error) {
	type Alias Modal
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&m),
	})
}
