package responses

import (
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ discord.InteractionData = &InteractionDataMessageComponent{}

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-object-message-component-data-structure
type InteractionDataMessageComponent struct {
	CustomID      string                                       `json:"custom_id"`
	ComponentType discord.ComponentType                        `json:"component_type"`
	Values        []components.StringSelectMenuComponentOption `json:"values,omitempty"`
	Resolved      *discord.ResolvedData                        `json:"resolved,omitempty"`
}

func (d *InteractionDataMessageComponent) GetType() discord.InteractionType {
	return discord.InteractionTypeMessageComponent
}
