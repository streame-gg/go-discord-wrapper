package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &UserUpdateEvent{}

func init() { RegisterEvent(&UserUpdateEvent{}) }

// https://docs.discord.com/developers/events/gateway-events#user-update
type UserUpdateEvent struct {
	NewUser discord.User  `json:"-"`
	OldUser *discord.User `json:"-"`
}

func (e *UserUpdateEvent) DesiredEventType() Event { return &UserUpdateEvent{} }
func (e *UserUpdateEvent) Event() EventType        { return EventUserUpdate }
func (e *UserUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewUser)
}
