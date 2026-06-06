package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/events/gateway-events#typing-start
type TypingStartEvent struct {
	ChannelID discord.Snowflake    `json:"channel_id"`
	GuildID   *discord.Snowflake   `json:"guild_id,omitempty"`
	UserID    discord.Snowflake    `json:"user_id"`
	Timestamp int64                `json:"timestamp"`
	Member    *discord.GuildMember `json:"member,omitempty"`
}

func init() { RegisterEvent(TypingStartEvent{}) }

func (e TypingStartEvent) DesiredEventType() Event { return &TypingStartEvent{} }
func (e TypingStartEvent) Event() EventType        { return EventTypingStart }
