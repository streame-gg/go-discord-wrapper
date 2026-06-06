package events

import (
	"encoding/json"
)

// TestAllRegisteredEvents_RoundTrip exercises every registered event's
// MarshalJSON / UnmarshalJSON / DesiredEventType in one sweep. For each factory
// it marshals the zero-value event, unmarshals the result back into a fresh
// instance, and re-checks the DesiredEventType contract. This pins the
// (de)serialisation surface that the gateway relies on for all event types at
// once.
func (s *registrySuite) TestAllRegisteredEvents_RoundTrip() {
	s.Require().NotEmpty(eventFactories)

	for key, factory := range eventFactories {
		evt := factory()
		s.Require().NotNilf(evt, "factory %q returned nil", key)

		// DesiredEventType must return a non-nil event of the same EventType.
		dt := evt.DesiredEventType()
		s.Require().NotNilf(dt, "%q DesiredEventType() is nil", key)
		s.Equalf(key, dt.Event(), "%q DesiredEventType().Event() mismatch", key)

		// Marshal the zero-value event (covers MarshalJSON where defined).
		b, err := json.Marshal(evt)
		s.Require().NoErrorf(err, "%q marshal", key)

		// Unmarshal back into a fresh instance (covers UnmarshalJSON).
		fresh := factory()
		s.Require().NoErrorf(json.Unmarshal(b, fresh), "%q unmarshal of %s", key, b)
	}
}

// syntheticEvents are wrapper-synthesised: they are never sent by Discord and
// thus not registered via init(), so the registry sweep above doesn't reach
// them. Cover their Event()/DesiredEventType() contract explicitly.
func (s *registrySuite) TestSyntheticExpressionEvents() {
	cases := []struct {
		evt Event
		typ EventType
	}{
		{GuildEmojiAddEvent{}, EventWrapperGuildEmojiAdd},
		{GuildEmojiRemoveEvent{}, EventWrapperGuildEmojiRemove},
		{GuildEmojiUpdateEvent{}, EventWrapperGuildEmojiUpdate},
		{GuildStickerAddEvent{}, EventWrapperGuildStickerAdd},
		{GuildStickerRemoveEvent{}, EventWrapperGuildStickerRemove},
		{GuildStickerUpdateEvent{}, EventWrapperGuildStickerUpdate},
	}

	for _, c := range cases {
		s.Equal(c.typ, c.evt.Event())
		dt := c.evt.DesiredEventType()
		s.Require().NotNil(dt)
		s.Equal(c.typ, dt.Event())
	}
}
