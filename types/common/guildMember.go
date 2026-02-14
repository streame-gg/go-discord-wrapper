package common

import (
	"time"
)

type GuildMember struct {
	AvatarHash                 *string    `json:"avatar,omitempty"`
	BannerHash                 *string    `json:"banner,omitempty"`
	CommunicationDisabledUntil *string    `json:"communication_disabled_until,omitempty"`
	Deaf                       bool       `json:"deaf"`
	Flags                      int        `json:"flags"`
	JoinedAt                   time.Time  `json:"joined_at"`
	Mute                       bool       `json:"mute"`
	Nick                       *string    `json:"nick,omitempty"`
	Pending                    bool       `json:"pending,omitempty"`
	PremiumSince               *time.Time `json:"premium_since,omitempty"`
	Roles                      []string   `json:"roles"`
	User                       *User      `json:"user,omitempty"`
}

func (m GuildMember) DisplayName() string {
	if m.Nick != nil {
		return *m.Nick
	}

	if m.User != nil {
		return m.User.DisplayName()
	}

	return ""
}

type ThreadMember struct {
	ID            *Snowflake   `json:"id,omitempty"`
	UserID        *Snowflake   `json:"user_id,omitempty"`
	JoinTimestamp time.Time    `json:"join_timestamp,omitempty"`
	Flags         int          `json:"flags,omitempty"`
	Member        *GuildMember `json:"member,omitempty"`
}
