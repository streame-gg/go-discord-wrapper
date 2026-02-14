package common

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
