package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type TypingStartEvent struct {
	ChannelID common.Snowflake    `json:"channel_id"`
	GuildID   *common.Snowflake   `json:"guild_id,omitempty"`
	UserID    common.Snowflake    `json:"user_id"`
	Timestamp int64               `json:"timestamp"`
	Member    *common.GuildMember `json:"member,omitempty"`
}

func init() { RegisterEvent(TypingStartEvent{}) }

func (e TypingStartEvent) DesiredEventType() Event { return &TypingStartEvent{} }
func (e TypingStartEvent) Event() EventType        { return EventTypingStart }
