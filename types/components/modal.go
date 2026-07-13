package components

import "encoding/json"

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-response-object-modal
type Modal struct {
	Title      string  `json:"title"`
	CustomID   string  `json:"custom_id"`
	Components []Label `json:"components"`
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
