package discord

import (
	"time"
)

// GuildMember.User is nil on MESSAGE_UPDATE gateway events.
//
// https://docs.discord.com/developers/resources/guild#guild-member-object
type GuildMember struct {
	hClient                    EntityClient
	GuildID                    Snowflake        `json:"-"`
	UserID                     Snowflake        `json:"-"`
	AvatarHash                 *string          `json:"avatar,omitempty"`
	BannerHash                 *string          `json:"banner,omitempty"`
	CommunicationDisabledUntil *time.Time       `json:"communication_disabled_until,omitempty"`
	Deaf                       *bool            `json:"deaf,omitempty"`
	Flags                      GuildMemberFlags `json:"flags"`
	JoinedAt                   time.Time        `json:"joined_at"`
	Mute                       *bool            `json:"mute,omitempty"`
	Nick                       *string          `json:"nick,omitempty"`
	Pending                    bool             `json:"pending,omitempty"`
	PremiumSince               *time.Time       `json:"premium_since,omitempty"`
	Roles                      []Snowflake      `json:"roles"`
	User                       *User            `json:"user,omitempty"`
}

func (m *GuildMember) DisplayName() string {
	if m == nil {
		return ""
	}

	if m.Nick != nil {
		return *m.Nick
	}

	if m.User != nil {
		return m.User.DisplayName()
	}

	return ""
}

// https://docs.discord.com/developers/resources/channel#thread-member-object
type ThreadMember struct {
	ID            *Snowflake   `json:"id,omitempty"`
	UserID        *Snowflake   `json:"user_id,omitempty"`
	JoinTimestamp time.Time    `json:"join_timestamp,omitempty"`
	Flags         int          `json:"flags,omitempty"`
	Member        *GuildMember `json:"member,omitempty"`
}
