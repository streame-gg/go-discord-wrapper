package events

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

type InviteCreateEvent struct {
	ChannelID         common.Snowflake             `json:"channel_id"`
	Code              string                       `json:"code"`
	CreatedAt         time.Time                    `json:"created_at"`
	GuildID           *common.Snowflake            `json:"guild_id"`
	Inviter           *common.User                 `json:"inviter,omitempty"`
	MaxAge            *int                         `json:"max_age,omitempty"`
	MaxUses           *int                         `json:"max_uses,omitempty"`
	TargetUserType    *common.InviteTargetUserType `json:"target_user_type,omitempty"`
	TargetUser        *common.User                 `json:"target_user_,omitempty"`
	TargetApplication *common.Application          `json:"target_application,omitempty"`
	Temporary         bool                         `json:"temporary"`
	Uses              int                          `json:"uses,omitempty"`
	ExpiresAt         *time.Time                   `json:"expires_at,omitempty"`
	RoleIDs           *[]common.Snowflake          `json:"role_ids,omitempty"`
}

func (i InviteCreateEvent) DesiredEventType() Event {
	return &InviteCreateEvent{}
}

func (i InviteCreateEvent) Event() EventType {
	return EventInviteCreate
}
