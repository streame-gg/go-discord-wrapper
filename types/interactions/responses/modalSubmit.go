package responses

import (
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-object-modal-submit-data-structure
type InteractionDataModalSubmit struct {
	CustomID   string                                `json:"custom_id"`
	Resolved   *discord.ResolvedData                 `json:"resolved,omitempty"`
	Components *[]components.ComponentLabelComponent `json:"components,omitempty"`
}

func (d *InteractionDataModalSubmit) GetType() discord.InteractionDataType {
	return discord.InteractionDataTypeModalSubmit
}
