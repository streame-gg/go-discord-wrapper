package events

import "github.com/streame-gg/go-discord-wrapper/types/discord"

type UserUpdateEvent struct {
	discord.User
}

func init() { RegisterEvent(UserUpdateEvent{}) }

func (e UserUpdateEvent) DesiredEventType() Event { return &UserUpdateEvent{} }
func (e UserUpdateEvent) Event() EventType        { return EventUserUpdate }
