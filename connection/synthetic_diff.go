package connection

import (
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// deriveSyntheticEvents returns the synthetic events to dispatch after ev.
// Called synchronously on the websocket reader goroutine after internalEventHandler
// has already applied cache updates and populated Old fields.
func (d *Client) deriveSyntheticEvents(ev events.Event) []events.Event {
	switch e := ev.(type) {
	case *events.VoiceStateUpdateEvent:
		return deriveVoiceSyntheticEvents(e)
	}
	return nil
}

// deriveVoiceSyntheticEvents derives 0 or 1 voice synthetic events from a
// VoiceStateUpdateEvent. The four transitions are mutually exclusive.
func deriveVoiceSyntheticEvents(ev *events.VoiceStateUpdateEvent) []events.Event {
	if ev.GuildID == nil {
		return nil
	}
	guildID := *ev.GuildID

	switch {
	case ev.OldState == nil && ev.ChannelID != nil:
		stateCopy := ev.VoiceState
		return []events.Event{&events.VoiceMemberJoinEvent{
			GuildID:   guildID,
			UserID:    ev.UserID,
			ChannelID: *ev.ChannelID,
			State:     &stateCopy,
		}}

	case ev.OldState != nil && ev.ChannelID == nil:
		if ev.OldState.ChannelID == nil {
			return nil
		}
		return []events.Event{&events.VoiceMemberLeaveEvent{
			GuildID:      guildID,
			UserID:       ev.UserID,
			OldChannelID: *ev.OldState.ChannelID,
			OldState:     ev.OldState,
		}}

	case ev.OldState != nil && ev.ChannelID != nil:
		if ev.OldState.ChannelID == nil {
			return nil
		}
		stateCopy := ev.VoiceState
		if *ev.OldState.ChannelID != *ev.ChannelID {
			return []events.Event{&events.VoiceMemberMoveEvent{
				GuildID:      guildID,
				UserID:       ev.UserID,
				OldChannelID: *ev.OldState.ChannelID,
				NewChannelID: *ev.ChannelID,
				OldState:     ev.OldState,
				NewState:     &stateCopy,
			}}
		}
		return []events.Event{&events.VoiceMemberUpdateEvent{
			GuildID:   guildID,
			UserID:    ev.UserID,
			ChannelID: *ev.ChannelID,
			OldState:  ev.OldState,
			NewState:  &stateCopy,
		}}
	}
	return nil
}
