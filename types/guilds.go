package types

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
	GetID() DiscordSnowflake
}

type Guild struct {
	ID                          DiscordSnowflake                       `json:"id"`
	Name                        string                                 `json:"name"`
	IconHash                    *string                                `json:"icon,omitempty"`
	Splash                      *string                                `json:"splash,omitempty"`
	DiscoverySplash             *string                                `json:"discovery_splash,omitempty"`
	Owner                       bool                                   `json:"owner,omitempty"`
	OwnerID                     DiscordSnowflake                       `json:"owner_id,omitempty"`
	Permissions                 *string                                `json:"permissions,omitempty"`
	Region                      *string                                `json:"region,omitempty"`
	AfkChannelID                *DiscordSnowflake                      `json:"afk_channel_id,omitempty"`
	AfkTimeout                  *int                                   `json:"afk_timeout,omitempty"`
	WidgetEnabled               *bool                                  `json:"widget_enabled,omitempty"`
	WidgetChannelID             *DiscordSnowflake                      `json:"widget_channel_id,omitempty"`
	VerificationChannel         *DiscordSnowflake                      `json:"verification_channel_id,omitempty"`
	VerificationLevel           DiscordGuildVerificationLevel          `json:"verification_level,omitempty"`
	DefaultMessageNotifications DiscordDefaultMessageNotificationLevel `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       DiscordGuildExplicitContentFilterLevel `json:"explicit_content_filter,omitempty"`
	Roles                       []DiscordRole                          `json:"roles,omitempty"`
	Emojis                      []DiscordEmoji                         `json:"emojis,omitempty"`
	Features                    []DiscordGuildFeatures                 `json:"features,omitempty"`
	MfaLevel                    DiscordGuildMFALevel                   `json:"mfa_level,omitempty"`
	ApplicationID               *DiscordSnowflake                      `json:"application_id,omitempty"`
	SystemChannelID             *DiscordSnowflake                      `json:"system_channel_id,omitempty"`
	SystemChannelFlags          DiscordGuildSystemChannelFlags         `json:"system_channel_flags,omitempty"`
	RulesChannelID              *DiscordSnowflake                      `json:"rules_channel_id,omitempty"`
	MaxPresences                *int                                   `json:"max_presences,omitempty"`
	MaxMembers                  int                                    `json:"max_members,omitempty"`
	VanityUrlCode               *string                                `json:"vanity_url_code,omitempty"`
	PremiumTier                 DiscordGuildPremiumTier                `json:"premium_tier,omitempty"`
	PremiumSubscriptionCount    *int                                   `json:"premium_subscription_count,omitempty"`
	PreferredLocale             DiscordLocale                          `json:"preferred_locale,omitempty"`
	PublicUpdatesChannelID      *DiscordSnowflake                      `json:"public_updates_channel_id,omitempty"`
	MaxVideoChannelUsers        *int                                   `json:"max_video_channel_users,omitempty"`
	MaxStageVideoChannelUsers   *int                                   `json:"max_stage_video_channel_users,omitempty"`
	ApproximateMemberCount      *int                                   `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount    *int                                   `json:"approximate_presence_count,omitempty"`
	WelcomeScreen               *DiscordGuildWelcomeScreen             `json:"welcome_screen,omitempty"`
	NSFWLevel                   DiscordGuildNSFWLevel                  `json:"nsfw_level,omitempty"`
	Stickers                    *[]DiscordSticker                      `json:"stickers,omitempty"`
	PremiumProgressBarEnabled   bool                                   `json:"premium_progress_bar_enabled,omitempty"`
	SafetyAlertsChannelID       *DiscordSnowflake                      `json:"safety_alerts_channel_id,omitempty"`
	IncidentsData               *DiscordGuildIncidentsData             `json:"incidents_data,omitempty"`
}

func (g Guild) IsAvailable() bool {
	return true
}

func (g Guild) GetID() DiscordSnowflake {
	return g.ID
}

type UnavailableGuild struct {
	ID          DiscordSnowflake `json:"id"`
	Unavailable *bool            `json:"unavailable"`
}

func (ug UnavailableGuild) IsAvailable() bool {
	return !*ug.Unavailable
}

func (ug UnavailableGuild) GetID() DiscordSnowflake {
	return ug.ID
}

type DiscordGuildVerificationLevel int

const (
	DiscordGuildVerificationLevelNone                                  DiscordGuildVerificationLevel = 0
	DiscordGuildVerificationLevelLow                                   DiscordGuildVerificationLevel = 1
	DiscordGuildVerificationLevelMedium                                DiscordGuildVerificationLevel = 2
	DiscordGuildVerificationLevelHigh                                  DiscordGuildVerificationLevel = 3
	DiscordGuildVerificationLevelVeryHighDiscordGuildVerificationLevel                               = 4
)

type DiscordDefaultMessageNotificationLevel int

const (
	DiscordDefaultMessageNotificationLevelAllMessages  DiscordDefaultMessageNotificationLevel = 0
	DiscordDefaultMessageNotificationLevelOnlyMentions DiscordDefaultMessageNotificationLevel = 1
)

type DiscordGuildPremiumTier int

const (
	DiscordGuildPremiumTierNone  DiscordGuildPremiumTier = 0
	DiscordGuildPremiumTierTier1 DiscordGuildPremiumTier = 1
	DiscordGuildPremiumTierTier2 DiscordGuildPremiumTier = 2
	DiscordGuildPremiumTierTier3 DiscordGuildPremiumTier = 3
)

type DiscordGuildNSFWLevel int

const (
	DiscordGuildNSFWLevelDefault       DiscordGuildNSFWLevel = 0
	DiscordGuildNSFWLevelExplicit      DiscordGuildNSFWLevel = 1
	DiscordGuildNSFWLevelSafe          DiscordGuildNSFWLevel = 2
	DiscordGuildNSFWLevelAgeRestricted DiscordGuildNSFWLevel = 3
)

type DiscordGuildExplicitContentFilterLevel int

const (
	DiscordGuildExplicitContentFilterLevelDisabled            DiscordGuildExplicitContentFilterLevel = 0
	DiscordGuildExplicitContentFilterLevelMembersWithoutRoles DiscordGuildExplicitContentFilterLevel = 1
	DiscordGuildExplicitContentFilterLevelAllMembers          DiscordGuildExplicitContentFilterLevel = 2
)

type DiscordGuildMFALevel int

const (
	DiscordGuildMFALevelNone     DiscordGuildMFALevel = 0
	DiscordGuildMFALevelElevated DiscordGuildMFALevel = 1
)

type DiscordGuildSystemChannelFlags int

const (
	DiscordGuildSystemChannelFlagsSuppressJoinNotifications       DiscordGuildSystemChannelFlags = 1 << 0
	DiscordGuildSystemChannelFlagsSuppressPremiumSubscriptions    DiscordGuildSystemChannelFlags = 1 << 1
	DiscordGuildSystemChannelFlagsSuppressGuildReminderMessages   DiscordGuildSystemChannelFlags = 1 << 2
	DiscordGuildSystemChannelFlagsSuppressJoinNotificationReplies DiscordGuildSystemChannelFlags = 1 << 3
	DiscordGuildSystemChannelFlagsPurchaseNotifications           DiscordGuildSystemChannelFlags = 1 << 4
	DiscordGuildSystemChannelFlagsPurchaseNotificationReplies     DiscordGuildSystemChannelFlags = 1 << 5
)

type DiscordGuildFeatures string

const (
	DiscordGuildFeatureAnimatedBanner                        DiscordGuildFeatures = "ANIMATED_BANNER"
	DiscordGuildFeatureAnimatedIcon                          DiscordGuildFeatures = "ANIMATED_ICON"
	DiscordGuildFeatureApplicationCommandPermissionsV2       DiscordGuildFeatures = "APPLICATION_COMMAND_PERMISSIONS_V2"
	DiscordGuildFeatureAutoModeration                        DiscordGuildFeatures = "AUTO_MODERATION"
	DiscordGuildFeatureBanner                                DiscordGuildFeatures = "BANNER"
	DiscordGuildFeatureCommunity                             DiscordGuildFeatures = "COMMUNITY"
	DiscordGuildFeatureCreatorMonetizableProvisional         DiscordGuildFeatures = "CREATOR_MONETIZABLE_PROVISIONAL"
	DiscordGuildFeatureCreatorStorePage                      DiscordGuildFeatures = "CREATOR_STORE_PAGE"
	DiscordGuildFeatureDeveloperSupportServer                DiscordGuildFeatures = "DEVELOPER_SUPPORT_SERVER"
	DiscordGuildFeatureDiscoverable                          DiscordGuildFeatures = "DISCOVERABLE"
	DiscordGuildFeatureFeaturable                            DiscordGuildFeatures = "FEATURABLE"
	DiscordGuildFeatureInvitesDisabled                       DiscordGuildFeatures = "INVITES_DISABLED"
	DiscordGuildFeatureInviteSplash                          DiscordGuildFeatures = "INVITE_SPLASH"
	DiscordGuildFeatureMemberVerificationGateEnabled         DiscordGuildFeatures = "MEMBER_VERIFICATION_GATE_ENABLED"
	DiscordGuildFeatureMoreSoundboard                        DiscordGuildFeatures = "MORE_SOUNDBOARD"
	DiscordGuildFeatureMoreStickers                          DiscordGuildFeatures = "MORE_STICKERS"
	DiscordGuildFeatureNews                                  DiscordGuildFeatures = "NEWS"
	DiscordGuildFeaturePartnered                             DiscordGuildFeatures = "PARTNERED"
	DiscordGuildFeaturePreviewEnabled                        DiscordGuildFeatures = "PREVIEW_ENABLED"
	DiscordGuildFeatureRaidAlertsDisabled                    DiscordGuildFeatures = "RAID_ALERTS_DISABLED"
	DiscordGuildFeatureRoleIcons                             DiscordGuildFeatures = "ROLE_ICONS"
	DiscordGuildFeatureRoleSubscriptionsAvailableForPurchase DiscordGuildFeatures = "ROLE_SUBSCRIPTIONS_AVAILABLE_FOR_PURCHASE"
	DiscordGuildFeatureRoleSubscriptionsEnabled              DiscordGuildFeatures = "ROLE_SUBSCRIPTIONS_ENABLED"
	DiscordGuildFeatureSoundboard                            DiscordGuildFeatures = "SOUNDBOARD"
	DiscordGuildFeatureTicketedEventsEnabled                 DiscordGuildFeatures = "TICKETED_EVENTS_ENABLED"
	DiscordGuildFeatureVanityURL                             DiscordGuildFeatures = "VANITY_URL"
	DiscordGuildFeatureVerified                              DiscordGuildFeatures = "VERIFIED"
	DiscordGuildFeatureVipRegions                            DiscordGuildFeatures = "VIP_REGIONS"
	DiscordGuildFeatureWelcomeScreenEnabled                  DiscordGuildFeatures = "WELCOME_SCREEN_ENABLED"
	DiscordGuildFeatureGuestsEnabled                         DiscordGuildFeatures = "GUESTS_ENABLED"
	DiscordGuildFeatureGuildTags                             DiscordGuildFeatures = "GUILD_TAGS"
	DiscordGuildFeatureEnhancedRoleColors                    DiscordGuildFeatures = "ENHANCED_ROLE_COLORS"
)

type DiscordGuildIncidentsData struct {
	InvitesDisabledUntil *time.Time `json:"invites_disabled_until,omitempty"`
	DmsDisabledUntil     *time.Time `json:"dms_disabled_until,omitempty"`
	DmSpanDetectedAt     *time.Time `json:"dm_spam_detected_at,omitempty"`
	RaidDetectedAt       *time.Time `json:"raid_detected_at,omitempty"`
}

type DiscordGuildWelcomeScreen struct {
	Description     *string                            `json:"description,omitempty"`
	WelcomeChannels []DiscordGuildWelcomeScreenChannel `json:"welcome_channels,omitempty"`
}

type DiscordGuildWelcomeScreenChannel struct {
	ChannelID   DiscordSnowflake  `json:"channel_id"`
	Description string            `json:"description"`
	EmojiID     *DiscordSnowflake `json:"emoji_id,omitempty"`
	EmojiName   *string           `json:"emoji_name,omitempty"`
}

type DiscordRoleColors struct {
	PrimaryColor   *int `json:"primary_color,omitempty"`
	SecondaryColor *int `json:"secondary_color,omitempty"`
	TertiaryColor  *int `json:"tertiary_color,omitempty"`
}

type DiscordRole struct {
	ID           DiscordSnowflake  `json:"id"`
	Name         string            `json:"name"`
	Colors       DiscordRoleColors `json:"colors,omitempty"`
	Hoist        bool              `json:"hoist"`
	IconHash     *string           `json:"icon,omitempty"`
	UnicodeEmoji *string           `json:"unicode_emoji,omitempty"`
	Position     int               `json:"position,omitempty"`
	Permissions  string            `json:"permissions"`
	Managed      bool              `json:"managed"`
	Mentionable  bool              `json:"mentionable"`
	Tags         *DiscordRoleTags  `json:"tags,omitempty"`
	Flags        *DiscordRoleFlags `json:"flags,omitempty"`
}

type DiscordRoleTags struct {
	BotID                 *DiscordSnowflake `json:"bot_id,omitempty"`
	IntegrationID         *DiscordSnowflake `json:"integration_id,omitempty"`
	PremiumSubscriber     *interface{}      `json:"premium_subscriber,omitempty"`
	SubscriptionListingID *DiscordSnowflake `json:"subscription_listing_id,omitempty"`
	AvailableForPurchase  *interface{}      `json:"available_for_purchase,omitempty"`
	GuildConnections      *interface{}      `json:"guild_connections,omitempty"`
}

type DiscordRoleFlags int

const (
	DiscordRoleFlagsInPrompt DiscordRoleFlags = 1 << 0
)

type DiscordSticker struct {
	ID          DiscordSnowflake         `json:"id"`
	PackID      *DiscordSnowflake        `json:"pack_id,omitempty"`
	Name        string                   `json:"name"`
	Description *string                  `json:"description,omitempty"`
	Tags        string                   `json:"tags"`
	Type        DiscordStickerType       `json:"type"`
	FormatType  DiscordStickerFormatType `json:"format_type"`
	Available   *bool                    `json:"available,omitempty"`
	GuildID     *DiscordSnowflake        `json:"guild_id,omitempty"`
	SortValue   *int                     `json:"sort_value,omitempty"`
	User        *DiscordUser             `json:"user,omitempty"`
}

type DiscordStickerType int

const (
	DiscordStickerTypeStandard DiscordStickerType = 1
	DiscordStickerTypeGuild    DiscordStickerType = 2
)

type DiscordStickerFormatType int

const (
	DiscordStickerFormatTypePNG    DiscordStickerFormatType = 1
	DiscordStickerFormatTypeAPNG   DiscordStickerFormatType = 2
	DiscordStickerFormatTypeLottie DiscordStickerFormatType = 3
	DiscordStickerFormatTypeGIF    DiscordStickerFormatType = 4
)
