package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/common"
)

type InviteDeleteEvent struct {
	ChannelID common.Snowflake  `json:"channel_id"`
	Code      string            `json:"code"`
	GuildID   *common.Snowflake `json:"guild_id,omitempty"`
}

func (i InviteDeleteEvent) DesiredEventType() Event {
	return &InviteDeleteEvent{}
}

func (i InviteDeleteEvent) Event() EventType {
	return EventInviteDelete
}
