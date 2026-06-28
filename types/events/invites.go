package events

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func init() {
	RegisterEvent(InviteCreateEvent{})
	RegisterEvent(InviteDeleteEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#invite-create
type InviteCreateEvent struct {
	ChannelID         discord.Snowflake             `json:"channel_id"`
	Code              string                        `json:"code"`
	CreatedAt         time.Time                     `json:"created_at"`
	GuildID           *discord.Snowflake            `json:"guild_id,omitempty"`
	Inviter           *discord.User                 `json:"inviter,omitempty"`
	MaxAge            *int                          `json:"max_age,omitempty"`
	MaxUses           *int                          `json:"max_uses,omitempty"`
	TargetUserType    *discord.InviteTargetUserType `json:"target_user_type,omitempty"`
	TargetUser        *discord.User                 `json:"target_user,omitempty"`
	TargetApplication *discord.Application          `json:"target_application,omitempty"`
	Temporary         bool                          `json:"temporary"`
	Uses              int                           `json:"uses,omitempty"`
	ExpiresAt         *time.Time                    `json:"expires_at,omitempty"`
	RoleIDs           *[]discord.Snowflake          `json:"role_ids,omitempty"`
}

// https://docs.discord.com/developers/events/gateway-events#invite-delete
type InviteDeleteEvent struct {
	ChannelID discord.Snowflake  `json:"channel_id"`
	Code      string             `json:"code"`
	GuildID   *discord.Snowflake `json:"guild_id,omitempty"`

	Invite *discord.Invite `json:"invite,omitempty"`
}

func (i InviteCreateEvent) DesiredEventType() Event {
	return &InviteCreateEvent{}
}
func (i InviteCreateEvent) Event() EventType {
	return EventInviteCreate
}

func (i InviteDeleteEvent) DesiredEventType() Event {
	return &InviteDeleteEvent{}
}
func (i InviteDeleteEvent) Event() EventType {
	return EventInviteDelete
}
