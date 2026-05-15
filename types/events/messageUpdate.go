package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type MessageUpdateEvent struct {
	discord.Message
	GuildID  *discord.Snowflake   `json:"guild_id"`
	Member   *discord.GuildMember `json:"member,omitempty"`
	Mentions *[]discord.User      `json:"mentions"`
}

func init() { RegisterEvent(MessageUpdateEvent{}) }

func (e MessageUpdateEvent) DesiredEventType() Event {
	return &MessageUpdateEvent{}
}

func (e MessageUpdateEvent) Event() EventType {
	return EventMessageUpdate
}
