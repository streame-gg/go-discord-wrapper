package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type UserUpdateEvent struct {
	NewUser discord.User `json:"-"`
}

func (e *UserUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewUser)
}

func (e UserUpdateEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&e.NewUser)
}

func init() { RegisterEvent(UserUpdateEvent{}) }

func (e UserUpdateEvent) DesiredEventType() Event { return &UserUpdateEvent{} }
func (e UserUpdateEvent) Event() EventType        { return EventUserUpdate }
