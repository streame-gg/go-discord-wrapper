package events

import (
	"testing"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
	"github.com/stretchr/testify/suite"
)

// registrySuite exercises the event-dispatch contract: the factory registry
// populated by each event's init(), and the InteractionCreateEvent
// classification helpers the gateway exposes to handlers.
type registrySuite struct {
	suite.Suite
}

func TestRegistrySuite(t *testing.T) {
	suite.Run(t, new(registrySuite))
}

// TestRegistryIsSelfConsistent asserts the central invariant of the dispatch
// table: for every registered EventType, the factory produces an Event whose
// Event() reports that same type. A copy-paste bug where DesiredEventType
// returns the wrong struct (and thus the wrong Event()) would surface here as a
// mismatched key, because the gateway keys handlers by EventType.
func (s *registrySuite) TestRegistryIsSelfConsistent() {
	s.Require().NotEmpty(eventFactories, "no events were registered via init()")

	for key, factory := range eventFactories {
		s.Require().NotNilf(factory, "factory for %q is nil", key)
		evt := factory()
		s.Require().NotNilf(evt, "factory for %q produced a nil event", key)
		s.Equalf(key, evt.Event(), "factory registered under %q returns Event()=%q", key, evt.Event())
	}
}

// TestGetEventFactory verifies the public lookup wrapper resolves a known type
// and reports absence for an unknown one.
func (s *registrySuite) TestGetEventFactory() {
	factory, ok := GetEventFactory(EventInteractionCreate)
	s.Require().True(ok, "INTERACTION_CREATE must be registered")
	s.Require().NotNil(factory)
	s.Equal(EventInteractionCreate, factory().Event())

	_, ok = GetEventFactory(EventType("NOT_A_REAL_EVENT"))
	s.False(ok, "unknown event type must report not found")
}

// TestRegisterEvent confirms RegisterEvent installs a factory keyed by the
// event's own Event() value. We register into a temporary type and restore the
// table afterwards so the global registry is left untouched.
func (s *registrySuite) TestRegisterEvent() {
	const key = EventType("TEST_ONLY_EVENT")
	s.Require().NotContains(eventFactories, key, "precondition: key must be unused")
	defer delete(eventFactories, key)

	RegisterEvent(testEvent{})

	factory, ok := GetEventFactory(key)
	s.Require().True(ok)
	s.Equal(key, factory().Event())
}

// testEvent is a throwaway Event used only by TestRegisterEvent.
type testEvent struct{}

func (testEvent) Event() EventType        { return EventType("TEST_ONLY_EVENT") }
func (testEvent) DesiredEventType() Event { return testEvent{} }

// ── InteractionCreateEvent classification helpers ───────────────────────────────

func (s *registrySuite) TestInteractionIsCommand() {
	s.True(newInteraction(discord.InteractionTypeApplicationCommand, nil).IsCommand())
	s.False(newInteraction(discord.InteractionTypeMessageComponent, nil).IsCommand())
}

func (s *registrySuite) TestInteractionIsAutocomplete() {
	s.True(newInteraction(discord.InteractionTypeApplicationCommandAutocomplete, nil).IsAutocomplete())
	s.False(newInteraction(discord.InteractionTypeApplicationCommand, nil).IsAutocomplete())
}

func (s *registrySuite) TestInteractionIsModalSubmit() {
	s.True(newInteraction(discord.InteractionTypeModalSubmit, nil).IsModalSubmit())
	s.False(newInteraction(discord.InteractionTypePing, nil).IsModalSubmit())
}

func (s *registrySuite) TestInteractionIsButton() {
	button := &responses.InteractionDataMessageComponent{ComponentType: discord.ComponentTypeButton}
	s.True(newInteraction(discord.InteractionTypeMessageComponent, button).IsButton(),
		"button component must be reported as a button")

	selectMenu := &responses.InteractionDataMessageComponent{ComponentType: discord.ComponentTypeStringSelect}
	s.False(newInteraction(discord.InteractionTypeMessageComponent, selectMenu).IsButton(),
		"select menu must not be reported as a button")

	// Wrong interaction type short-circuits before inspecting Data.
	s.False(newInteraction(discord.InteractionTypeApplicationCommand, button).IsButton())

	// Message component interaction with mismatched Data type is not a button.
	s.False(newInteraction(discord.InteractionTypeMessageComponent, nil).IsButton())
}

func (s *registrySuite) TestInteractionIsAnySelectMenu() {
	selectMenu := &responses.InteractionDataMessageComponent{ComponentType: discord.ComponentTypeStringSelect}
	s.True(newInteraction(discord.InteractionTypeMessageComponent, selectMenu).IsAnySelectMenu())

	userSelect := &responses.InteractionDataMessageComponent{ComponentType: discord.ComponentTypeUserSelect}
	s.True(newInteraction(discord.InteractionTypeMessageComponent, userSelect).IsAnySelectMenu())

	button := &responses.InteractionDataMessageComponent{ComponentType: discord.ComponentTypeButton}
	s.False(newInteraction(discord.InteractionTypeMessageComponent, button).IsAnySelectMenu())

	// Non-component interaction type is never a select menu.
	s.False(newInteraction(discord.InteractionTypeApplicationCommand, selectMenu).IsAnySelectMenu())
}

// newInteraction builds an InteractionCreateEvent fixture with the given
// interaction type and (optional) component data.
func newInteraction(t discord.InteractionType, data discord.InteractionData) InteractionCreateEvent {
	return InteractionCreateEvent{
		Interaction: interactions.Interaction{Type: t, Data: data},
	}
}
