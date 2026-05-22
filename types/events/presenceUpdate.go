package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type PresenceUpdateEvent struct {
	NewPresence discord.Presence  `json:"-"`
	OldPresence *discord.Presence `json:"-"`
}

func (e *PresenceUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewPresence)
}

func (e PresenceUpdateEvent) MarshalJSON() ([]byte, error) {
	type wire struct {
		discord.Presence
		OldPresence *discord.Presence `json:"old_presence,omitempty"`
	}
	return json.Marshal(wire{e.NewPresence, e.OldPresence})
}

func init() { RegisterEvent(PresenceUpdateEvent{}) }

func (e PresenceUpdateEvent) DesiredEventType() Event { return &PresenceUpdateEvent{} }
func (e PresenceUpdateEvent) Event() EventType        { return EventPresenceUpdate }
