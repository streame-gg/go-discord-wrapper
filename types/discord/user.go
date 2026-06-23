package discord

import "time"

// https://docs.discord.com/developers/resources/user#avatar-decoration-data-object
type AvatarDecorationData struct {
	Asset string `json:"asset"`
	SkuID string `json:"sku_id"`
}

// https://docs.discord.com/developers/resources/user#user-object-user-primary-guild
type PrimaryGuild struct {
	Badge           *string `json:"badge"`
	IdentityEnabled *bool   `json:"identity_enabled"`
	IdentityGuildID *string `json:"identity_guild_id"`
	Tag             *string `json:"tag"`
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
	hClient EntityClient

	AccentColor          *int                  `json:"accent_color,omitempty"`
	Avatar               *string               `json:"avatar"`
	AvatarDecorationData *AvatarDecorationData `json:"avatar_decoration,omitempty"`
	Bot                  *bool                 `json:"bot,omitempty"`
	Discriminator        string                `json:"discriminator"`
	Flags                UserFlags             `json:"flags"`
	GlobalName           *string               `json:"global_name"`
	ID                   Snowflake             `json:"id"`
	Locale               *string               `json:"locale,omitempty"`
	MFAEnabled           *bool                 `json:"mfa_enabled"`
	PrimaryGuild         *PrimaryGuild         `json:"primary_guild,omitempty"`
	PublicFlags          UserFlags             `json:"public_flags"`
	System               *bool                 `json:"system,omitempty"`
	Username             string                `json:"username"`
	Banner               *string               `json:"banner,omitempty"`
	Verified             *bool                 `json:"verified,omitempty"`
	Email                *string               `json:"email,omitempty"`
	PremiumType          *UserPremiumType      `json:"premium_type,omitempty"`
	Collectibles         *Collectible          `json:"collectibles,omitempty"`
}

// https://docs.discord.com/developers/resources/user#user-object-premium-types
type UserPremiumType uint8

const (
	UserPremiumTypeNone         UserPremiumType = 0
	UserPremiumTypeNitroClassic UserPremiumType = 1
	UserPremiumTypeNitro        UserPremiumType = 2
	UserPremiumTypeNitroBasic   UserPremiumType = 3
)

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
	ID           string                   `json:"id"`
	Name         string                   `json:"name"`
	Type         UserConnectionType       `json:"type"`
	Revoked      *bool                    `json:"revoked,omitempty"`
	Integrations []Integration            `json:"integrations,omitempty"`
	Verified     bool                     `json:"verified"`
	FriendSync   bool                     `json:"friend_sync"`
	ShowActivity bool                     `json:"show_activity"`
	TwoWayLink   bool                     `json:"two_way_link"`
	Visibility   UserConnectionVisibility `json:"visibility"`
}

// https://docs.discord.com/developers/resources/user#connection-object-visibility-types
type UserConnectionVisibility uint8

const (
	UserConnectionVisibilityNone     UserConnectionVisibility = 0
	UserConnectionVisibilityEveryone UserConnectionVisibility = 1
)

// https://docs.discord.com/developers/resources/user#connection-object-services
type UserConnectionType string

const (
	UserConnectionTypeAmazonMusic     UserConnectionType = "amazon-music"
	UserConnectionTypeBattleNet       UserConnectionType = "battlenet"
	UserConnectionTypeBungie          UserConnectionType = "bungie"
	UserConnectionTypeBluesky         UserConnectionType = "bluesky"
	UserConnectionTypeCrunchyRoll     UserConnectionType = "crunchyroll"
	UserConnectionTypeDomain          UserConnectionType = "domain"
	UserConnectionTypeEbay            UserConnectionType = "ebay"
	UserConnectionTypeEpicGames       UserConnectionType = "epicgames"
	UserConnectionTypeFacebook        UserConnectionType = "facebook"
	UserConnectionTypeGithub          UserConnectionType = "github"
	UserConnectionTypeInstagram       UserConnectionType = "instagram"
	UserConnectionTypeLeagueOfLegends UserConnectionType = "leagueoflegends"
	UserConnectionTypeMastodon        UserConnectionType = "mastodon"
	UserConnectionTypePayPal          UserConnectionType = "paypal"
	UserConnectionTypePlaystation     UserConnectionType = "playstation"
	UserConnectionTypeReddit          UserConnectionType = "reddit"
	UserConnectionTypeRiotGames       UserConnectionType = "riotgames"
	UserConnectionTypeRoblox          UserConnectionType = "roblox"
	UserConnectionTypeSpotify         UserConnectionType = "spotify"
	UserConnectionTypeSkype           UserConnectionType = "skype"
	UserConnectionTypeSteam           UserConnectionType = "steam"
	UserConnectionTypeTikTok          UserConnectionType = "tiktok"
	UserConnectionTypeTwitch          UserConnectionType = "twitch"
	UserConnectionTypeTwitter         UserConnectionType = "twitter"
	UserConnectionTypeXbox            UserConnectionType = "xbox"
	UserConnectionTypeYoutube         UserConnectionType = "youtube"
)

// OAuth2Authorization represents the current OAuth2 authorization info for a bearer token.
//
// https://docs.discord.com/developers/topics/oauth2#get-current-authorization-information
type OAuth2Authorization struct {
	Application Application `json:"application"`
	Expires     time.Time   `json:"expires"`
	Scopes      []Scope     `json:"scopes"`
	User        *User       `json:"user,omitempty"`
}

// ApplicationRoleConnection represents the user's role connection for an application.
//
// https://docs.discord.com/developers/resources/user#application-role-connection-object
type ApplicationRoleConnection struct {
	PlatformName *string                                      `json:"platform_name"`
	Metadata     map[string]ApplicationRoleConnectionMetadata `json:"metadata"`
}
