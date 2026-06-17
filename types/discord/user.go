package discord

import "time"

// https://docs.discord.com/developers/resources/user#avatar-decoration-data-object
type AvatarDecorationData struct {
	Asset     string `json:"asset"`
	ExpiresAt int64  `json:"expires_at"`
	SkuID     string `json:"sku_id"`
}

// https://docs.discord.com/developers/resources/user#user-object
type Clan struct {
	Badge           string `json:"badge"`
	IdentityEnabled bool   `json:"identity_enabled"`
	IdentityGuildID string `json:"identity_guild_id"`
	Tag             string `json:"tag"`
}

// https://docs.discord.com/developers/resources/user#collectibles
type Collectible struct {
	Nameplate *Nameplate `json:"nameplate,omitempty"`
}

// https://docs.discord.com/developers/resources/user#nameplate-nameplate-structure
type Nameplate struct {
	SkuID   string           `json:"sku_id"`
	Asset   string           `json:"asset"`
	Label   string           `json:"label"`
	Palette NameplatePalette `json:"palette"`
}

// https://docs.discord.com/developers/resources/user#nameplate-nameplate-structure
type NameplatePalette string

const (
	NameplatePaletteCrimson   NameplatePalette = "crimson"
	NameplatePaletteBerry     NameplatePalette = "berry"
	NameplatePaletteSky       NameplatePalette = "sky"
	NameplatePaletteTeal      NameplatePalette = "teal"
	NameplatePaletteForest    NameplatePalette = "forest"
	NameplatePaletteBubbleGum NameplatePalette = "bubble_gum"
	NameplatePaletteViolet    NameplatePalette = "violet"
	NameplatePaletteCobalt    NameplatePalette = "cobalt"
	NameplatePaletteClover    NameplatePalette = "clover"
	NameplatePaletteLemon     NameplatePalette = "lemon"
	NameplatePaletteWhite     NameplatePalette = "white"
)

// https://docs.discord.com/developers/resources/user#user-object
type User struct {
	hClient              EntityClient
	AccentColor          *int                  `json:"accent_color,omitempty"`
	AvatarHash           *string               `json:"avatar"`
	AvatarDecorationData *AvatarDecorationData `json:"avatar_decoration,omitempty"`
	Bot                  bool                  `json:"bot,omitempty"`
	Discriminator        string                `json:"discriminator"`
	Flags                UserFlags             `json:"flags"`
	GlobalName           *string               `json:"global_name,omitempty"`
	ID                   Snowflake             `json:"id"`
	Locale               *string               `json:"locale,omitempty"`
	MFAEnabled           bool                  `json:"mfa_enabled"`
	PrimaryGuild         *Clan                 `json:"primary_guild,omitempty"`
	PublicFlags          UserFlags             `json:"public_flags"`
	System               bool                  `json:"system,omitempty"`
	Username             string                `json:"username"`
}

func (u *User) DisplayName() string {
	if u.GlobalName != nil && *u.GlobalName != "" {
		return *u.GlobalName
	}

	return u.Username
}

func (u *User) CreatedAt() time.Time {
	return u.ID.Time()
}

// UserConnection represents an external account linked to a Discord user.
//
// https://docs.discord.com/developers/resources/user#connection-object
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
//
// https://docs.discord.com/developers/topics/oauth2#get-current-authorization-information
type OAuth2Authorization struct {
	Application Application `json:"application"`
	Expires     string      `json:"expires"`
	Scopes      []Scope     `json:"scopes"`
	User        *User       `json:"user,omitempty"`
}

// ApplicationRoleConnection represents the user's role connection for an application.
//
// https://docs.discord.com/developers/resources/user#application-role-connection-object
type ApplicationRoleConnection struct {
	PlatformName     *string           `json:"platform_name,omitempty"`
	PlatformUsername *string           `json:"platform_username,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}
