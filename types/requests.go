package types

type BotRegisterResponse struct {
	Url               string `json:"url"`
	Shards            int    `json:"shards"`
	SessionStartLimit struct {
		Total      int `json:"total"`
		Remaining  int `json:"remaining"`
		ResetAfter int `json:"reset_after"`
	} `json:"session_start_limit"`
}
