package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &PresenceUpdateEvent{}

func init() { RegisterEvent(&PresenceUpdateEvent{}) }

// https://docs.discord.com/developers/events/gateway-events#presence-update
type PresenceUpdateEvent struct {
	NewPresence discord.Presence  `json:"-"`
	OldPresence *discord.Presence `json:"-"`
}

func (e *PresenceUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewPresence)
}
func (e *PresenceUpdateEvent) DesiredEventType() Event { return &PresenceUpdateEvent{} }
func (e *PresenceUpdateEvent) Event() EventType        { return EventPresenceUpdate }
