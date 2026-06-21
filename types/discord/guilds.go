package discord

import (
	"encoding/json"
	"time"
)

// https://docs.discord.com/developers/events/gateway-events#guild-create
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

// https://docs.discord.com/developers/resources/guild#guild-object
type AnyGuild interface {
	IsAvailable() bool
	GetID() Snowflake
}

// https://docs.discord.com/developers/resources/guild#guild-object
type Guild struct {
	hClient EntityClient

	// Sub-managers populated by the connection layer after hydration.
	members          MemberManager
	roles            RoleManager
	channels         GuildChannelManager
	emojis           EmojiManager
	stickers         StickerManager
	bans             BanManager
	scheduledEvents  ScheduledEventManager
	stageInstances   StageInstanceManager
	soundboardSounds SoundboardManager
	invites          GuildInviteManager
	voiceStates      VoiceStateManager
	autoModRules     AutoModRuleManager
	webhooks         GuildWebhookManager
	integrations     IntegrationManager

	ID              Snowflake  `json:"id"`
	Name            string     `json:"name"`
	Icon            *string    `json:"icon"`
	IconHash        *string    `json:"icon_hash"`
	Splash          *string    `json:"splash"`
	DiscoverySplash *string    `json:"discovery_splash"`
	Owner           bool       `json:"owner,omitempty"`
	OwnerID         Snowflake  `json:"owner_id"`
	Permissions     Permission `json:"permissions,omitempty"`
	// Region is deprecated
	Region                      *string                         `json:"region,omitempty"`
	AfkChannelID                *Snowflake                      `json:"afk_channel_id"`
	AfkTimeout                  int                             `json:"afk_timeout"`
	WidgetEnabled               bool                            `json:"widget_enabled,omitempty"`
	WidgetChannelID             *Snowflake                      `json:"widget_channel_id,omitempty"`
	VerificationLevel           GuildVerificationLevel          `json:"verification_level"`
	DefaultMessageNotifications DefaultMessageNotificationLevel `json:"default_message_notifications"`
	ExplicitContentFilter       GuildExplicitContentFilterLevel `json:"explicit_content_filter"`
	RawRoles                    []Role                          `json:"roles"`
	RawEmojis                   []Emoji                         `json:"emojis"`
	Features                    []GuildFeatures                 `json:"features"`
	MfaLevel                    GuildMFALevel                   `json:"mfa_level"`
	ApplicationID               *Snowflake                      `json:"application_id"`
	SystemChannelID             *Snowflake                      `json:"system_channel_id"`
	SystemChannelFlags          GuildSystemChannelFlags         `json:"system_channel_flags"`
	RulesChannelID              *Snowflake                      `json:"rules_channel_id"`
	MaxPresences                *int                            `json:"max_presences"`
	MaxMembers                  int                             `json:"max_members,omitempty"`
	VanityUrlCode               *string                         `json:"vanity_url_code"`
	Description                 *string                         `json:"description"`
	Banner                      *string                         `json:"banner"`
	PremiumTier                 GuildPremiumTier                `json:"premium_tier"`
	PremiumSubscriptionCount    int                             `json:"premium_subscription_count,omitempty"`
	PreferredLocale             Locale                          `json:"preferred_locale,omitempty"`
	PublicUpdatesChannelID      *Snowflake                      `json:"public_updates_channel_id"`
	MaxVideoChannelUsers        int                             `json:"max_video_channel_users,omitempty"`
	MaxStageVideoChannelUsers   int                             `json:"max_stage_video_channel_users,omitempty"`
	ApproximateMemberCount      int                             `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount    int                             `json:"approximate_presence_count,omitempty"`
	WelcomeScreen               *GuildWelcomeScreen             `json:"welcome_screen,omitempty"`
	NSFWLevel                   GuildNSFWLevel                  `json:"nsfw_level"`
	RawStickers                 []Sticker                       `json:"stickers,omitempty"`
	PremiumProgressBarEnabled   bool                            `json:"premium_progress_bar_enabled"`
	SafetyAlertsChannelID       *Snowflake                      `json:"safety_alerts_channel_id"`
	IncidentsData               *GuildIncidentsData             `json:"incidents_data,omitempty"`
}

func (g Guild) IsAvailable() bool {
	return true
}

func (g Guild) GetID() Snowflake {
	return g.ID
}

// https://docs.discord.com/developers/resources/guild#unavailable-guild-object
type UnavailableGuild struct {
	ID          Snowflake `json:"id"`
	Unavailable *bool     `json:"unavailable"`
}

func (ug UnavailableGuild) IsAvailable() bool {
	if ug.Unavailable == nil {
		return false
	}
	return !*ug.Unavailable
}

func (ug UnavailableGuild) GetID() Snowflake {
	return ug.ID
}

// https://docs.discord.com/developers/resources/guild#guild-object-verification-level
type GuildVerificationLevel int

const (
	GuildVerificationLevelNone     GuildVerificationLevel = 0
	GuildVerificationLevelLow      GuildVerificationLevel = 1
	GuildVerificationLevelMedium   GuildVerificationLevel = 2
	GuildVerificationLevelHigh     GuildVerificationLevel = 3
	GuildVerificationLevelVeryHigh GuildVerificationLevel = 4
)

// https://docs.discord.com/developers/resources/guild#guild-object-default-message-notification-level
type DefaultMessageNotificationLevel int

const (
	DefaultMessageNotificationLevelAllMessages  DefaultMessageNotificationLevel = 0
	DefaultMessageNotificationLevelOnlyMentions DefaultMessageNotificationLevel = 1
)

// https://docs.discord.com/developers/resources/guild#guild-object-premium-tier
type GuildPremiumTier int

const (
	GuildPremiumTierNone  GuildPremiumTier = 0
	GuildPremiumTierTier1 GuildPremiumTier = 1
	GuildPremiumTierTier2 GuildPremiumTier = 2
	GuildPremiumTierTier3 GuildPremiumTier = 3
)

// https://docs.discord.com/developers/resources/guild#guild-object-guild-nsfw-level
type GuildNSFWLevel int

const (
	GuildNSFWLevelDefault       GuildNSFWLevel = 0
	GuildNSFWLevelExplicit      GuildNSFWLevel = 1
	GuildNSFWLevelSafe          GuildNSFWLevel = 2
	GuildNSFWLevelAgeRestricted GuildNSFWLevel = 3
)

// https://docs.discord.com/developers/resources/guild#guild-object-explicit-content-filter-level
type GuildExplicitContentFilterLevel int

const (
	GuildExplicitContentFilterLevelDisabled            GuildExplicitContentFilterLevel = 0
	GuildExplicitContentFilterLevelMembersWithoutRoles GuildExplicitContentFilterLevel = 1
	GuildExplicitContentFilterLevelAllMembers          GuildExplicitContentFilterLevel = 2
)

// https://docs.discord.com/developers/resources/guild#guild-object-mfa-level
type GuildMFALevel int

func (g GuildMFALevel) ToInt() int {
	return int(g)
}

const (
	GuildMFALevelNone     GuildMFALevel = 0
	GuildMFALevelElevated GuildMFALevel = 1
)

// https://docs.discord.com/developers/resources/guild#guild-object-system-channel-flags
type GuildSystemChannelFlags int

const (
	GuildSystemChannelFlagsSuppressJoinNotifications       GuildSystemChannelFlags = 1 << 0
	GuildSystemChannelFlagsSuppressPremiumSubscriptions    GuildSystemChannelFlags = 1 << 1
	GuildSystemChannelFlagsSuppressGuildReminderMessages   GuildSystemChannelFlags = 1 << 2
	GuildSystemChannelFlagsSuppressJoinNotificationReplies GuildSystemChannelFlags = 1 << 3
	GuildSystemChannelFlagsPurchaseNotifications           GuildSystemChannelFlags = 1 << 4
	GuildSystemChannelFlagsPurchaseNotificationReplies     GuildSystemChannelFlags = 1 << 5
)

// https://docs.discord.com/developers/resources/guild#guild-object-guild-features
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

// https://docs.discord.com/developers/resources/guild#incidents-data-object
type GuildIncidentsData struct {
	InvitesDisabledUntil *time.Time `json:"invites_disabled_until"`
	DmsDisabledUntil     *time.Time `json:"dms_disabled_until"`
	DmSpanDetectedAt     *time.Time `json:"dm_spam_detected_at,omitempty"`
	RaidDetectedAt       *time.Time `json:"raid_detected_at,omitempty"`
}

// https://docs.discord.com/developers/resources/guild#welcome-screen-object
type GuildWelcomeScreen struct {
	Description     *string                     `json:"description"`
	WelcomeChannels []GuildWelcomeScreenChannel `json:"welcome_channels"`
}

// https://docs.discord.com/developers/resources/guild#welcome-screen-object-welcome-screen-channel-structure
type GuildWelcomeScreenChannel struct {
	ChannelID   Snowflake  `json:"channel_id"`
	Description string     `json:"description"`
	EmojiID     *Snowflake `json:"emoji_id"`
	EmojiName   *string    `json:"emoji_name"`
}

// https://docs.discord.com/developers/topics/permissions#role-object-role-colors-object
type RoleColors struct {
	PrimaryColor   int  `json:"primary_color"`
	SecondaryColor *int `json:"secondary_color"`
	TertiaryColor  *int `json:"tertiary_color"`
}

// https://docs.discord.com/developers/topics/permissions#role-object
type Role struct {
	hClient EntityClient
	GuildID Snowflake `json:"-"`

	ID           Snowflake  `json:"id"`
	Name         string     `json:"name"`
	Colors       RoleColors `json:"colors"`
	Hoist        bool       `json:"hoist"`
	Icon         *string    `json:"icon,omitempty"`
	UnicodeEmoji *string    `json:"unicode_emoji,omitempty"`
	Position     int        `json:"position"`
	Permissions  Permission `json:"permissions"`
	Managed      bool       `json:"managed"`
	Mentionable  bool       `json:"mentionable"`
	Tags         *RoleTags  `json:"tags,omitempty"`
	Flags        *RoleFlags `json:"flags,omitempty"`
}

// https://docs.discord.com/developers/topics/permissions#role-object-role-tags-structure
type RoleTags struct {
	BotID                 Snowflake `json:"bot_id,omitempty"`
	IntegrationID         Snowflake `json:"integration_id,omitempty"`
	PremiumSubscriber     NullFlag  `json:"premium_subscriber,omitempty"`
	SubscriptionListingID Snowflake `json:"subscription_listing_id,omitempty"`
	AvailableForPurchase  NullFlag  `json:"available_for_purchase,omitempty"`
	GuildConnections      NullFlag  `json:"guild_connections,omitempty"`
}

// https://docs.discord.com/developers/topics/permissions#role-object-role-flags
type RoleFlags int

const (
	RoleFlagsInPrompt RoleFlags = 1 << 0
)

// https://docs.discord.com/developers/resources/sticker#sticker-object
type Sticker struct {
	hClient EntityClient

	ID          Snowflake         `json:"id"`
	PackID      Snowflake         `json:"pack_id,omitempty"`
	Name        string            `json:"name"`
	Description *string           `json:"description"`
	Tags        string            `json:"tags"`
	Type        StickerType       `json:"type"`
	FormatType  StickerFormatType `json:"format_type"`
	Available   *bool             `json:"available,omitempty"`
	GuildID     Snowflake         `json:"guild_id,omitempty"`
	SortValue   int               `json:"sort_value,omitempty"`
	User        *User             `json:"user,omitempty"`
}

// https://docs.discord.com/developers/resources/sticker#sticker-object-sticker-types
type StickerType int

const (
	StickerTypeStandard StickerType = 1
	StickerTypeGuild    StickerType = 2
)

// https://docs.discord.com/developers/resources/sticker#sticker-object-sticker-format-types
type StickerFormatType int

const (
	StickerFormatTypePNG    StickerFormatType = 1
	StickerFormatTypeAPNG   StickerFormatType = 2
	StickerFormatTypeLottie StickerFormatType = 3
	StickerFormatTypeGIF    StickerFormatType = 4
)

// GuildWidgetSettings holds the guild widget configuration.
//
// https://docs.discord.com/developers/resources/guild#guild-widget-settings-object
type GuildWidgetSettings struct {
	Enabled   bool       `json:"enabled"`
	ChannelID *Snowflake `json:"channel_id"`
}

// GuildWidget is the public JSON widget for a guild.
//
// https://docs.discord.com/developers/resources/guild#guild-widget-object
type GuildWidget struct {
	ID            Snowflake `json:"id"`
	Name          string    `json:"name"`
	InstantInvite *string   `json:"instant_invite"`
	Channels      []Channel `json:"channels"`
	Members       []User    `json:"members"`
	PresenceCount int       `json:"presence_count"`
}

// GatewayGuild extends [Guild] with fields that are only present in the
// GUILD_CREATE gateway event payload and are never sent in REST responses.
//
// https://docs.discord.com/developers/events/gateway-events#guild-create-guild-create-extra-fields
type GatewayGuild struct {
	Guild

	JoinedAt             time.Time             `json:"joined_at"`
	Large                bool                  `json:"large"`
	Unavailable          *bool                 `json:"unavailable,omitempty"`
	MemberCount          int                   `json:"member_count,omitempty"`
	VoiceStates          []VoiceState          `json:"voice_states,omitempty"`
	Members              []GuildMember         `json:"members,omitempty"`
	Channels             []Channel             `json:"channels,omitempty"`
	Threads              []Channel             `json:"threads,omitempty"`
	Presences            []GatewayPresence     `json:"presences,omitempty"`
	StageInstances       []StageInstance       `json:"stage_instances,omitempty"`
	GuildScheduledEvents []GuildScheduledEvent `json:"guild_scheduled_events,omitempty"`
	SoundboardSounds     []SoundboardSound     `json:"soundboard_sounds,omitempty"`
}

func (g GatewayGuild) IsAvailable() bool { return true }
func (g GatewayGuild) GetID() Snowflake  { return g.ID }

// https://docs.discord.com/developers/resources/guild#ban-object
type Ban struct {
	Reason *string `json:"reason"`
	User   User    `json:"user"`
}

// https://docs.discord.com/developers/resources/guild#get-guild-vanity-url
type GuildVanityURL struct {
	Code *string `json:"code"`
	Uses int     `json:"uses"`
}

// https://docs.discord.com/developers/resources/guild#get-guild-prune-count
type GuildPruneCountResult struct {
	Pruned int `json:"pruned"`
}

// https://docs.discord.com/developers/resources/sticker#sticker-pack-object
type StickerPack struct {
	ID             Snowflake `json:"id"`
	Stickers       []Sticker `json:"stickers"`
	Name           string    `json:"name"`
	SKUId          Snowflake `json:"sku_id"`
	CoverStickerID Snowflake `json:"cover_sticker_id,omitempty"`
	Description    string    `json:"description"`
	BannerAssetID  Snowflake `json:"banner_asset_id,omitempty"`
}

// https://docs.discord.com/developers/resources/user#get-current-user-guilds
type CurrentUserGuild struct {
	ID                       Snowflake       `json:"id"`
	Name                     string          `json:"name"`
	Icon                     *string         `json:"icon"`
	Owner                    bool            `json:"owner"`
	Permissions              Permission      `json:"permissions"`
	Features                 []GuildFeatures `json:"features"`
	ApproximateMemberCount   int             `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount int             `json:"approximate_presence_count,omitempty"`
}

// GatewayPresence is a stripped-down presence record as sent inside the
// GUILD_CREATE payload. It has no guild_id field; the enclosing guild provides
// that context.
//
// https://docs.discord.com/developers/events/gateway-events#presence-update
type GatewayPresence struct {
	User         PartialPresenceUser `json:"user"`
	GuildID      Snowflake           `json:"guild_id"`
	Status       PresenceStatus      `json:"status"`
	Activities   []FullActivity      `json:"activities"`
	ClientStatus ClientStatus        `json:"client_status"`
}

// GatewayGuildWrapper is like [AnyGuildWrapper] but deserialises available
// guilds as [GatewayGuild] so that gateway-only fields (channels, members,
// voice states, etc.) are retained. Used by [GuildCreateEvent].
//
// https://docs.discord.com/developers/events/gateway-events#guild-create
type GatewayGuildWrapper struct {
	Guild AnyGuild
}

func (ag *GatewayGuildWrapper) UnmarshalJSON(data []byte) error {
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
	var g GatewayGuild
	if err := json.Unmarshal(data, &g); err != nil {
		return err
	}
	ag.Guild = g
	return nil
}
