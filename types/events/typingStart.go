package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &TypingStartEvent{}

func init() { RegisterEvent(&TypingStartEvent{}) }

// https://docs.discord.com/developers/events/gateway-events#typing-start
type TypingStartEvent struct {
	ChannelID discord.Snowflake    `json:"channel_id"`
	GuildID   *discord.Snowflake   `json:"guild_id,omitempty"`
	UserID    discord.Snowflake    `json:"user_id"`
	Timestamp int64                `json:"timestamp"`
	Member    *discord.GuildMember `json:"member,omitempty"`

	Guild *discord.Guild `json:"-"`
}

func (e *TypingStartEvent) DesiredEventType() Event { return &TypingStartEvent{} }
func (e *TypingStartEvent) Event() EventType        { return EventTypingStart }
