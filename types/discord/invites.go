package discord

import "time"

type InviteTargetUserType int

const (
	InviteTargetUserTypeStream              InviteTargetUserType = 1
	InviteTargetUserTypeEmbeddedApplication InviteTargetUserType = 2
)

type InviteFlags int

const (
	IsGuestInvite InviteFlags = 1 << 0
)

type Invite struct {
	hClient                  EntityClient
	Code                     string                `json:"code"`
	Guild                    *Guild                `json:"guild,omitempty"`
	Channel                  *Channel              `json:"channel,omitempty"`
	Inviter                  *User                 `json:"inviter,omitempty"`
	TargetType               *InviteTargetUserType `json:"target_type,omitempty"`
	TargetUser               *User                 `json:"target_user,omitempty"`
	ApproximatePresenceCount *int                  `json:"approximate_presence_count,omitempty"`
	ApproximateMemberCount   *int                  `json:"approximate_member_count,omitempty"`
	ExpiresAt                *time.Time            `json:"expires_at,omitempty"`
	Flags                    *InviteFlags          `json:"flags,omitempty"`
	Uses                     *int                  `json:"uses,omitempty"`
	MaxUses                  *int                  `json:"max_uses,omitempty"`
	MaxAge                   *int                  `json:"max_age,omitempty"`
	Temporary                *bool                 `json:"temporary,omitempty"`
	CreatedAt                *time.Time            `json:"created_at,omitempty"`
}
