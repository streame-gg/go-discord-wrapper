package events

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

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

func init() { RegisterEvent(InviteCreateEvent{}) }

func (i InviteCreateEvent) DesiredEventType() Event {
	return &InviteCreateEvent{}
}

func (i InviteCreateEvent) Event() EventType {
	return EventInviteCreate
}
