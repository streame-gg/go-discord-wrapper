package types

import "time"

type AvatarDecorationData struct {
	Asset     string `json:"asset"`
	ExpiresAt int64  `json:"expires_at"`
	SkuID     string `json:"sku_id"`
}

type DiscordClan struct {
	Badge           string `json:"badge"`
	IdentityEnabled bool   `json:"identity_enabled"`
	IdentityGuildID string `json:"identity_guild_id"`
	Tag             string `json:"tag"`
}

type DiscordUser struct {
	AccentColor          *int                  `json:"accent_color,omitempty"`
	AvatarHash           string                `json:"avatar"`
	AvatarDecorationData *AvatarDecorationData `json:"avatar_decoration,omitempty"`
	Bot                  bool                  `json:"bot,omitempty"`
	Discriminator        string                `json:"discriminator"`
	Flags                int                   `json:"flags"`
	GlobalName           *string               `json:"global_name,omitempty"`
	ID                   string                `json:"id"`
	Locale               *string               `json:"locale,omitempty"`
	MFAEnabled           bool                  `json:"mfa_enabled"`
	PrimaryGuild         *DiscordClan          `json:"primary_guild,omitempty"`
	PublicFlags          int                   `json:"public_flags"`
	System               bool                  `json:"system,omitempty"`
	Username             string                `json:"username"`
}

func (u *DiscordUser) DisplayName() string {
	if u.GlobalName != nil && *u.GlobalName != "" {
		return *u.GlobalName
	}

	return u.Username
}

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
}
