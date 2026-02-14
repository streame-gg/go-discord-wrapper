package common

import (
	"encoding/json"
	"time"
)

type AnyGuildWrapper struct {
	Guild AnyGuild
}

func (ag *AnyGuildWrapper) UnmarshalJSON(data []byte) error {
	var probe struct {
		Unavailable *bool `json:"unavailable"`
	}

	if err := json.Unmarshal(data, &probe); err != nil {
		return err
	}

	if probe.Unavailable != nil && *probe.Unavailable {
		var ug UnavailableGuild
		if err := json.Unmarshal(data, &ug); err != nil {
			return err
		}
		ag.Guild = ug
		return nil
	}

	var g Guild
	if err := json.Unmarshal(data, &g); err != nil {
		return err
	}
	ag.Guild = g
	return nil
}

type AnyGuild interface {
	IsAvailable() bool
	GetID() Snowflake
}

type Guild struct {
	ID                          Snowflake                       `json:"id"`
	Name                        string                          `json:"name"`
	IconHash                    *string                         `json:"icon,omitempty"`
	Splash                      *string                         `json:"splash,omitempty"`
	DiscoverySplash             *string                         `json:"discovery_splash,omitempty"`
	Owner                       bool                            `json:"owner,omitempty"`
	OwnerID                     Snowflake                       `json:"owner_id,omitempty"`
	Permissions                 *string                         `json:"permissions,omitempty"`
	Region                      *string                         `json:"region,omitempty"`
	AfkChannelID                *Snowflake                      `json:"afk_channel_id,omitempty"`
	AfkTimeout                  *int                            `json:"afk_timeout,omitempty"`
	WidgetEnabled               *bool                           `json:"widget_enabled,omitempty"`
	WidgetChannelID             *Snowflake                      `json:"widget_channel_id,omitempty"`
	VerificationChannel         *Snowflake                      `json:"verification_channel_id,omitempty"`
	VerificationLevel           GuildVerificationLevel          `json:"verification_level,omitempty"`
	DefaultMessageNotifications DefaultMessageNotificationLevel `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       GuildExplicitContentFilterLevel `json:"explicit_content_filter,omitempty"`
	Roles                       []Role                          `json:"roles,omitempty"`
	Emojis                      []Emoji                         `json:"emojis,omitempty"`
	Features                    []GuildFeatures                 `json:"features,omitempty"`
	MfaLevel                    GuildMFALevel                   `json:"mfa_level,omitempty"`
	ApplicationID               *Snowflake                      `json:"application_id,omitempty"`
	SystemChannelID             *Snowflake                      `json:"system_channel_id,omitempty"`
	SystemChannelFlags          GuildSystemChannelFlags         `json:"system_channel_flags,omitempty"`
	RulesChannelID              *Snowflake                      `json:"rules_channel_id,omitempty"`
	MaxPresences                *int                            `json:"max_presences,omitempty"`
	MaxMembers                  int                             `json:"max_members,omitempty"`
	VanityUrlCode               *string                         `json:"vanity_url_code,omitempty"`
	PremiumTier                 GuildPremiumTier                `json:"premium_tier,omitempty"`
	PremiumSubscriptionCount    *int                            `json:"premium_subscription_count,omitempty"`
	PreferredLocale             Locale                          `json:"preferred_locale,omitempty"`
	PublicUpdatesChannelID      *Snowflake                      `json:"public_updates_channel_id,omitempty"`
	MaxVideoChannelUsers        *int                            `json:"max_video_channel_users,omitempty"`
	MaxStageVideoChannelUsers   *int                            `json:"max_stage_video_channel_users,omitempty"`
	ApproximateMemberCount      *int                            `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount    *int                            `json:"approximate_presence_count,omitempty"`
	WelcomeScreen               *GuildWelcomeScreen             `json:"welcome_screen,omitempty"`
	NSFWLevel                   GuildNSFWLevel                  `json:"nsfw_level,omitempty"`
	Stickers                    *[]Sticker                      `json:"stickers,omitempty"`
	PremiumProgressBarEnabled   bool                            `json:"premium_progress_bar_enabled,omitempty"`
	SafetyAlertsChannelID       *Snowflake                      `json:"safety_alerts_channel_id,omitempty"`
	IncidentsData               *GuildIncidentsData             `json:"incidents_data,omitempty"`
}

func (g Guild) IsAvailable() bool {
	return true
}

func (g Guild) GetID() Snowflake {
	return g.ID
}

type UnavailableGuild struct {
	ID          Snowflake `json:"id"`
	Unavailable *bool     `json:"unavailable"`
}

func (ug UnavailableGuild) IsAvailable() bool {
	return !*ug.Unavailable
}

func (ug UnavailableGuild) GetID() Snowflake {
	return ug.ID
}

type GuildVerificationLevel int

const (
	GuildVerificationLevelNone                           GuildVerificationLevel = 0
	GuildVerificationLevelLow                            GuildVerificationLevel = 1
	GuildVerificationLevelMedium                         GuildVerificationLevel = 2
	GuildVerificationLevelHigh                           GuildVerificationLevel = 3
	GuildVerificationLevelVeryHighGuildVerificationLevel                        = 4
)

type DefaultMessageNotificationLevel int

const (
	DefaultMessageNotificationLevelAllMessages  DefaultMessageNotificationLevel = 0
	DefaultMessageNotificationLevelOnlyMentions DefaultMessageNotificationLevel = 1
)

type GuildPremiumTier int

const (
	GuildPremiumTierNone  GuildPremiumTier = 0
	GuildPremiumTierTier1 GuildPremiumTier = 1
	GuildPremiumTierTier2 GuildPremiumTier = 2
	GuildPremiumTierTier3 GuildPremiumTier = 3
)

type GuildNSFWLevel int

const (
	GuildNSFWLevelDefault       GuildNSFWLevel = 0
	GuildNSFWLevelExplicit      GuildNSFWLevel = 1
	GuildNSFWLevelSafe          GuildNSFWLevel = 2
	GuildNSFWLevelAgeRestricted GuildNSFWLevel = 3
)

type GuildExplicitContentFilterLevel int

const (
	GuildExplicitContentFilterLevelDisabled            GuildExplicitContentFilterLevel = 0
	GuildExplicitContentFilterLevelMembersWithoutRoles GuildExplicitContentFilterLevel = 1
	GuildExplicitContentFilterLevelAllMembers          GuildExplicitContentFilterLevel = 2
)

type GuildMFALevel int

const (
	GuildMFALevelNone     GuildMFALevel = 0
	GuildMFALevelElevated GuildMFALevel = 1
)

type GuildSystemChannelFlags int

const (
	GuildSystemChannelFlagsSuppressJoinNotifications       GuildSystemChannelFlags = 1 << 0
	GuildSystemChannelFlagsSuppressPremiumSubscriptions    GuildSystemChannelFlags = 1 << 1
	GuildSystemChannelFlagsSuppressGuildReminderMessages   GuildSystemChannelFlags = 1 << 2
	GuildSystemChannelFlagsSuppressJoinNotificationReplies GuildSystemChannelFlags = 1 << 3
	GuildSystemChannelFlagsPurchaseNotifications           GuildSystemChannelFlags = 1 << 4
	GuildSystemChannelFlagsPurchaseNotificationReplies     GuildSystemChannelFlags = 1 << 5
)

type GuildFeatures string

const (
	GuildFeatureAnimatedBanner                        GuildFeatures = "ANIMATED_BANNER"
	GuildFeatureAnimatedIcon                          GuildFeatures = "ANIMATED_ICON"
	GuildFeatureApplicationCommandPermissionsV2       GuildFeatures = "APPLICATION_COMMAND_PERMISSIONS_V2"
	GuildFeatureAutoModeration                        GuildFeatures = "AUTO_MODERATION"
	GuildFeatureBanner                                GuildFeatures = "BANNER"
	GuildFeatureCommunity                             GuildFeatures = "COMMUNITY"
	GuildFeatureCreatorMonetizableProvisional         GuildFeatures = "CREATOR_MONETIZABLE_PROVISIONAL"
	GuildFeatureCreatorStorePage                      GuildFeatures = "CREATOR_STORE_PAGE"
	GuildFeatureDeveloperSupportServer                GuildFeatures = "DEVELOPER_SUPPORT_SERVER"
	GuildFeatureDiscoverable                          GuildFeatures = "DISCOVERABLE"
	GuildFeatureFeaturable                            GuildFeatures = "FEATURABLE"
	GuildFeatureInvitesDisabled                       GuildFeatures = "INVITES_DISABLED"
	GuildFeatureInviteSplash                          GuildFeatures = "INVITE_SPLASH"
	GuildFeatureMemberVerificationGateEnabled         GuildFeatures = "MEMBER_VERIFICATION_GATE_ENABLED"
	GuildFeatureMoreSoundboard                        GuildFeatures = "MORE_SOUNDBOARD"
	GuildFeatureMoreStickers                          GuildFeatures = "MORE_STICKERS"
	GuildFeatureNews                                  GuildFeatures = "NEWS"
	GuildFeaturePartnered                             GuildFeatures = "PARTNERED"
	GuildFeaturePreviewEnabled                        GuildFeatures = "PREVIEW_ENABLED"
	GuildFeatureRaidAlertsDisabled                    GuildFeatures = "RAID_ALERTS_DISABLED"
	GuildFeatureRoleIcons                             GuildFeatures = "ROLE_ICONS"
	GuildFeatureRoleSubscriptionsAvailableForPurchase GuildFeatures = "ROLE_SUBSCRIPTIONS_AVAILABLE_FOR_PURCHASE"
	GuildFeatureRoleSubscriptionsEnabled              GuildFeatures = "ROLE_SUBSCRIPTIONS_ENABLED"
	GuildFeatureSoundboard                            GuildFeatures = "SOUNDBOARD"
	GuildFeatureTicketedEventsEnabled                 GuildFeatures = "TICKETED_EVENTS_ENABLED"
	GuildFeatureVanityURL                             GuildFeatures = "VANITY_URL"
	GuildFeatureVerified                              GuildFeatures = "VERIFIED"
	GuildFeatureVipRegions                            GuildFeatures = "VIP_REGIONS"
	GuildFeatureWelcomeScreenEnabled                  GuildFeatures = "WELCOME_SCREEN_ENABLED"
	GuildFeatureGuestsEnabled                         GuildFeatures = "GUESTS_ENABLED"
	GuildFeatureGuildTags                             GuildFeatures = "GUILD_TAGS"
	GuildFeatureEnhancedRoleColors                    GuildFeatures = "ENHANCED_ROLE_COLORS"
)

type GuildIncidentsData struct {
	InvitesDisabledUntil *time.Time `json:"invites_disabled_until,omitempty"`
	DmsDisabledUntil     *time.Time `json:"dms_disabled_until,omitempty"`
	DmSpanDetectedAt     *time.Time `json:"dm_spam_detected_at,omitempty"`
	RaidDetectedAt       *time.Time `json:"raid_detected_at,omitempty"`
}

type GuildWelcomeScreen struct {
	Description     *string                     `json:"description,omitempty"`
	WelcomeChannels []GuildWelcomeScreenChannel `json:"welcome_channels,omitempty"`
}

type GuildWelcomeScreenChannel struct {
	ChannelID   Snowflake  `json:"channel_id"`
	Description string     `json:"description"`
	EmojiID     *Snowflake `json:"emoji_id,omitempty"`
	EmojiName   *string    `json:"emoji_name,omitempty"`
}

type RoleColors struct {
	PrimaryColor   *int `json:"primary_color,omitempty"`
	SecondaryColor *int `json:"secondary_color,omitempty"`
	TertiaryColor  *int `json:"tertiary_color,omitempty"`
}

type Role struct {
	ID           Snowflake  `json:"id"`
	Name         string     `json:"name"`
	Colors       RoleColors `json:"colors,omitempty"`
	Hoist        bool       `json:"hoist"`
	IconHash     *string    `json:"icon,omitempty"`
	UnicodeEmoji *string    `json:"unicode_emoji,omitempty"`
	Position     int        `json:"position,omitempty"`
	Permissions  string     `json:"permissions"`
	Managed      bool       `json:"managed"`
	Mentionable  bool       `json:"mentionable"`
	Tags         *RoleTags  `json:"tags,omitempty"`
	Flags        *RoleFlags `json:"flags,omitempty"`
}

type RoleTags struct {
	BotID                 *Snowflake   `json:"bot_id,omitempty"`
	IntegrationID         *Snowflake   `json:"integration_id,omitempty"`
	PremiumSubscriber     *interface{} `json:"premium_subscriber,omitempty"`
	SubscriptionListingID *Snowflake   `json:"subscription_listing_id,omitempty"`
	AvailableForPurchase  *interface{} `json:"available_for_purchase,omitempty"`
	GuildConnections      *interface{} `json:"guild_connections,omitempty"`
}

type RoleFlags int

const (
	RoleFlagsInPrompt RoleFlags = 1 << 0
)

type Sticker struct {
	ID          Snowflake         `json:"id"`
	PackID      *Snowflake        `json:"pack_id,omitempty"`
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	Tags        string            `json:"tags"`
	Type        StickerType       `json:"type"`
	FormatType  StickerFormatType `json:"format_type"`
	Available   *bool             `json:"available,omitempty"`
	GuildID     *Snowflake        `json:"guild_id,omitempty"`
	SortValue   *int              `json:"sort_value,omitempty"`
	User        *User             `json:"user,omitempty"`
}

type StickerType int

const (
	StickerTypeStandard StickerType = 1
	StickerTypeGuild    StickerType = 2
)

type StickerFormatType int

const (
	StickerFormatTypePNG    StickerFormatType = 1
	StickerFormatTypeAPNG   StickerFormatType = 2
	StickerFormatTypeLottie StickerFormatType = 3
	StickerFormatTypeGIF    StickerFormatType = 4
)
