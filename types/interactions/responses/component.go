package responses

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type InteractionDataMessageComponent struct {
	CustomID      string                `json:"custom_id"`
	ComponentType discord.ComponentType `json:"component_type"`
	Values        []string              `json:"values,omitempty"`
	Resolved      *discord.ResolvedData `json:"resolved,omitempty"`
}

func (d *InteractionDataMessageComponent) GetType() discord.InteractionDataType {
	return discord.InteractionDataTypeMessageComponent
}
