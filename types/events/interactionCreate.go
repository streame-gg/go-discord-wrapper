package events

import (
	"github.com/DatGamet/go-discord-wrapper/types/common"
	"github.com/DatGamet/go-discord-wrapper/types/interactions"
	"github.com/DatGamet/go-discord-wrapper/types/interactions/responses"
)

type InteractionCreateEvent struct {
	interactions.Interaction
}

func (e InteractionCreateEvent) DesiredEventType() Event {
	return &InteractionCreateEvent{}
}

func (e InteractionCreateEvent) Event() EventType {
	return EventInteractionCreate
}

func (e InteractionCreateEvent) IsCommand() bool {
	return e.Type == common.InteractionTypeApplicationCommand
}

func (e InteractionCreateEvent) IsButton() bool {
	if e.Type != common.InteractionTypeMessageComponent {
		return false
	}

	return e.Data.(*responses.InteractionDataMessageComponent).ComponentType == common.ComponentTypeButton
}

func (e InteractionCreateEvent) IsAnySelectMenu() bool {
	if e.Type != common.InteractionTypeMessageComponent {
		return false
	}

	return e.Data.(*responses.InteractionDataMessageComponent).ComponentType.IsAnySelectMenu()
}

func (e InteractionCreateEvent) IsAutocomplete() bool {
	return e.Type == common.InteractionTypeApplicationCommandAutocomplete
}

func (e InteractionCreateEvent) IsModalSubmit() bool {
	return e.Type == common.InteractionTypeModalSubmit
}
