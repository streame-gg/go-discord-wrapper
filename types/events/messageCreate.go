package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type MessageCreateEvent struct {
	discord.Message
	GuildID  *discord.Snowflake   `json:"guild_id,omitempty"`
	Member   *discord.GuildMember `json:"member,omitempty"`
	Mentions *[]discord.User      `json:"mentions"`
}

func init() { RegisterEvent(MessageCreateEvent{}) }

func (e MessageCreateEvent) DesiredEventType() Event {
	return &MessageCreateEvent{}
}

func (e MessageCreateEvent) Event() EventType {
	return EventMessageCreate
}
