package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

var _ Event = &InteractionCreateEvent{}

func init() { RegisterEvent(&InteractionCreateEvent{}) }

// https://docs.discord.com/developers/events/gateway-events#interaction-create
type InteractionCreateEvent struct {
	interactions.Interaction
}

func (e *InteractionCreateEvent) DesiredEventType() Event {
	return &InteractionCreateEvent{}
}

func (e *InteractionCreateEvent) Event() EventType {
	return EventInteractionCreate
}

func (e *InteractionCreateEvent) IsCommand() bool {
	return e.Type == discord.InteractionTypeApplicationCommand
}

func (e *InteractionCreateEvent) IsButton() bool {
	if e.Type != discord.InteractionTypeMessageComponent {
		return false
	}

	comp, ok := e.Data.(*responses.InteractionDataMessageComponent)
	if !ok {
		return false
	}
	return comp.ComponentType == discord.ComponentTypeButton
}

func (e *InteractionCreateEvent) IsAnySelectMenu() bool {
	if e.Type != discord.InteractionTypeMessageComponent {
		return false
	}

	comp, ok := e.Data.(*responses.InteractionDataMessageComponent)
	if !ok {
		return false
	}
	return comp.ComponentType.IsAnySelectMenu()
}

func (e *InteractionCreateEvent) IsAutocomplete() bool {
	return e.Type == discord.InteractionTypeApplicationCommandAutocomplete
}

func (e *InteractionCreateEvent) IsModalSubmit() bool {
	return e.Type == discord.InteractionTypeModalSubmit
}
