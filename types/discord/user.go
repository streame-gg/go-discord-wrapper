package discord

type AvatarDecorationData struct {
	Asset     string `json:"asset"`
	ExpiresAt int64  `json:"expires_at"`
	SkuID     string `json:"sku_id"`
}

type Clan struct {
	Badge           string `json:"badge"`
	IdentityEnabled bool   `json:"identity_enabled"`
	IdentityGuildID string `json:"identity_guild_id"`
	Tag             string `json:"tag"`
}

type User struct {
	hClient              EntityClient
	AccentColor          *int                  `json:"accent_color,omitempty"`
	AvatarHash           string                `json:"avatar"`
	AvatarDecorationData *AvatarDecorationData `json:"avatar_decoration,omitempty"`
	Bot                  bool                  `json:"bot,omitempty"`
	Discriminator        string                `json:"discriminator"`
	Flags                int                   `json:"flags"`
	GlobalName           *string               `json:"global_name,omitempty"`
	ID                   Snowflake             `json:"id"`
	Locale               *string               `json:"locale,omitempty"`
	MFAEnabled           bool                  `json:"mfa_enabled"`
	PrimaryGuild         *Clan                 `json:"primary_guild,omitempty"`
	PublicFlags          int                   `json:"public_flags"`
	System               bool                  `json:"system,omitempty"`
	Username             string                `json:"username"`
}

func (u *User) DisplayName() string {
	if u.GlobalName != nil && *u.GlobalName != "" {
		return *u.GlobalName
	}

	return u.Username
}

// UserConnection represents an external account linked to a Discord user.
type UserConnection struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Revoked      *bool          `json:"revoked,omitempty"`
	Integrations []*Integration `json:"integrations,omitempty"`
	Verified     bool           `json:"verified"`
	FriendSync   bool           `json:"friend_sync"`
	ShowActivity bool           `json:"show_activity"`
	TwoWayLink   bool           `json:"two_way_link"`
	Visibility   int            `json:"visibility"`
}

// OAuth2Authorization represents the current OAuth2 authorization info for a bearer token.
type OAuth2Authorization struct {
	Application Application `json:"application"`
	Expires     string      `json:"expires"`
	Scopes      []string    `json:"scopes"`
	User        *User       `json:"user,omitempty"`
}

// ApplicationRoleConnection represents the user's role connection for an application.
type ApplicationRoleConnection struct {
	PlatformName     *string           `json:"platform_name,omitempty"`
	PlatformUsername *string           `json:"platform_username,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}
