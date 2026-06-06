package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/events/gateway-events#invite-delete
type InviteDeleteEvent struct {
	ChannelID discord.Snowflake  `json:"channel_id"`
	Code      string             `json:"code"`
	GuildID   *discord.Snowflake `json:"guild_id,omitempty"`
}

func init() { RegisterEvent(InviteDeleteEvent{}) }

func (i InviteDeleteEvent) DesiredEventType() Event {
	return &InviteDeleteEvent{}
}

func (i InviteDeleteEvent) Event() EventType {
	return EventInviteDelete
}
