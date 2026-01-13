package types

import (
	"time"
)

type DiscordSnowflake string

type DiscordActivity struct {
	Type    int     `json:"type"`
	PartyID *string `json:"party_id,omitempty"`
}

type DiscordApplicationEventWebhookStatus int

const (
	DiscordApplicationEventWebhookStatusDisabled          DiscordApplicationEventWebhookStatus = 0
	DiscordApplicationEventWebhookStatusEnabled           DiscordApplicationEventWebhookStatus = 1
	DiscordApplicationEventWebhookStatusDisabledByDiscord DiscordApplicationEventWebhookStatus = 2
)

type DiscordApplicationTeamMemberMembershipState int

const (
	DiscordApplicationTeamMemberMembershipStateInvited  DiscordApplicationTeamMemberMembershipState = 1
	DiscordApplicationTeamMemberMembershipStateAccepted DiscordApplicationTeamMemberMembershipState = 2
)

type DiscordApplicationTeamMember struct {
	MembershipState DiscordApplicationTeamMemberMembershipState `json:"membership_state"`
	TeamID          *string                                     `json:"team_id,omitempty"`
	User            *DiscordUser                                `json:"user,omitempty"`
	Role            string                                      `json:"role"`
}

type DiscordApplicationTeam struct {
	IconHash    *string                        `json:"icon,omitempty"`
	ID          DiscordSnowflake               `json:"id"`
	Members     []DiscordApplicationTeamMember `json:"members"`
	Name        string                         `json:"name"`
	OwnerUserID DiscordSnowflake               `json:"owner_user_id"`
}

type DiscordApplicationInstallParams struct {
	Scopes      []string `json:"scopes"`
	Permissions string   `json:"permissions"`
}

type DiscordApplication struct {
	ID                                DiscordSnowflake                     `json:"id"`
	Name                              string                               `json:"name"`
	IconHash                          *string                              `json:"icon,omitempty"`
	Description                       string                               `json:"description,omitempty"`
	RpcOrigins                        []string                             `json:"rpc_origins,omitempty"`
	BotPublic                         bool                                 `json:"bot_public,omitempty"`
	BotRequireCodeGrant               bool                                 `json:"bot_require_code_grant,omitempty"`
	TermsOfServiceURL                 *string                              `json:"terms_of_service_url,omitempty"`
	PrivacyPolicyURL                  *string                              `json:"privacy_policy_url,omitempty"`
	Owner                             *DiscordUser                         `json:"owner,omitempty"`
	VerifyKey                         *string                              `json:"verify_key,omitempty"`
	Team                              *DiscordApplicationTeam              `json:"team,omitempty"`
	GuildID                           *DiscordSnowflake                    `json:"guild_id,omitempty"`
	Guild                             Guild                                `json:"guild,omitempty"`
	PrimarySKUID                      *string                              `json:"primary_sku_id,omitempty"`
	Slug                              *string                              `json:"slug,omitempty"`
	CoverImage                        *string                              `json:"cover_image,omitempty"`
	Flags                             *int                                 `json:"flags,omitempty"`
	ApproximateGuildCount             *int                                 `json:"approximate_guild_count,omitempty"`
	ApproximateUserInstallCount       *int                                 `json:"approximate_user_install_count,omitempty"`
	ApproximateUserAuthorizationCount *int                                 `json:"approximate_user_authorization_count,omitempty"`
	RedirectURIs                      *[]string                            `json:"redirect_uris,omitempty"`
	InteractionEndpointURL            *string                              `json:"interaction_endpoint_url,omitempty"`
	RoleConnectionsVerificationURL    *string                              `json:"role_connections_verification_url,omitempty"`
	EventWebhooksURL                  *string                              `json:"event_webhooks_url,omitempty"`
	EventWebhookStatus                DiscordApplicationEventWebhookStatus `json:"event_webhook_status,omitempty"`
	EventWebhooksTypes                *[]string                            `json:"event_webhooks_types,omitempty"`
	Tags                              *[]string                            `json:"tags,omitempty"`
	InstallParams                     *DiscordApplicationInstallParams     `json:"install_params,omitempty"`
	IntegrationTypesConfig            *interface{}                         `json:"integration_types_config,omitempty"`
	CustomInstallURL                  *string                              `json:"custom_install_url,omitempty"`
}

type DiscordCall struct {
	Participants   []DiscordSnowflake `json:"participants"`
	EndedTimestamp *time.Time         `json:"ended_timestamp,omitempty"`
}

type DiscordMessageChannelMention struct {
	ID      string `json:"id"`
	GuildID string `json:"guild_id"`
	Type    int    `json:"type"`
	Name    string `json:"name"`
}

type DiscordMessageStickerItemFormatType int

const (
	DiscordMessageStickerItemFormatTypePNG    DiscordMessageStickerItemFormatType = 1
	DiscordMessageStickerItemFormatTypeAPNG   DiscordMessageStickerItemFormatType = 2
	DiscordMessageStickerItemFormatTypeLottie DiscordMessageStickerItemFormatType = 3
	DiscordMessageStickerItemFormatTypeGIF    DiscordMessageStickerItemFormatType = 4
)

type DiscordMessageStickerItem struct {
	ID         DiscordSnowflake                    `json:"id"`
	Name       string                              `json:"name"`
	FormatType DiscordMessageStickerItemFormatType `json:"format_type"`
}

type DiscordRoleSubscriptionData struct {
	RoleSubscriptionListingID DiscordSnowflake `json:"role_subscription_listing_id"`
	TierName                  string           `json:"tier_name"`
	TotalMonthsSubscribed     int              `json:"total_months_subscribed"`
	IsRenewal                 bool             `json:"is_renewal"`
}

type DiscordResolvedData struct {
	Users       map[DiscordSnowflake]*DiscordUser       `json:"users"`
	Members     map[DiscordSnowflake]*GuildMember       `json:"members,omitempty"`
	Messages    map[DiscordSnowflake]*DiscordMessage    `json:"messages,omitempty"`
	Channels    map[DiscordSnowflake]*DiscordChannel    `json:"channels,omitempty"`
	Roles       map[DiscordSnowflake]*DiscordRole       `json:"roles,omitempty"`
	Attachments map[DiscordSnowflake]*DiscordAttachment `json:"attachments,omitempty"`
}

type PartialDiscordMessage struct {
	Type            DiscordMessageType          `json:"type"`
	Content         string                      `json:"content"`
	Embeds          []DiscordEmbed              `json:"embeds,omitempty"`
	Attachments     []DiscordAttachment         `json:"attachments,omitempty"`
	Timestamp       *time.Time                  `json:"timestamp,omitempty"`
	EditedTimestamp *time.Time                  `json:"edited_timestamp,omitempty"`
	Flags           DiscordMessageFlag          `json:"flags,omitempty"`
	Mentions        *[]any                      `json:"mentions"`
	MentionRoles    []string                    `json:"mention_roles"`
	StickerItems    []DiscordMessageStickerItem `json:"sticker_items,omitempty"`
	Components      []AnyComponent              `json:"components,omitempty"`
}

type DiscordMessageMessageSnapshot struct {
	Message PartialDiscordMessage `json:"message,omitempty"`
}

type DiscordMessageMessageReferenceType int

const (
	DiscordMessageMessageReferenceTypeDefault DiscordMessageMessageReferenceType = 0
	DiscordMessageMessageReferenceTypeForward DiscordMessageMessageReferenceType = 1
)

type DiscordMessage struct {
	Activity      *DiscordActivity    `json:"activity,omitempty"`
	Application   *DiscordApplication `json:"application,omitempty"`
	ApplicationID *string             `json:"application_id,omitempty"`
	Attachments   []DiscordAttachment `json:"attachments,omitempty"`
	Author        *DiscordUser        `json:"author"`
	Call          *DiscordCall        `json:"call,omitempty"`
	ChannelID     DiscordSnowflake    `json:"channel_id"`
	ChannelType   DiscordChannelType  `json:"channel_type"`
	//TODO FIXME this does not wort with AnyComponent currently
	Components           interface{}                        `json:"components"`
	Content              string                             `json:"content"`
	EditedTimestamp      *time.Time                         `json:"edited_timestamp,omitempty"`
	Embeds               []DiscordEmbed                     `json:"embeds,omitempty"`
	Flags                DiscordMessageFlag                 `json:"flags"`
	ID                   DiscordSnowflake                   `json:"id"`
	InteractionMetadata  *DiscordMessageInteractionMetadata `json:"interaction_metadata,omitempty"`
	MentionEveryone      bool                               `json:"mention_everyone"`
	MentionChannels      *[]DiscordMessageChannelMention    `json:"mention_channels,omitempty"`
	MentionRoles         []string                           `json:"mention_roles"`
	MessageReference     *DiscordMessageMessageReference    `json:"message_reference,omitempty"`
	MessageSnapshots     []DiscordMessageMessageSnapshot    `json:"message_snapshots,omitempty"`
	Nonce                interface{}                        `json:"nonce,omitempty"`
	Pinned               bool                               `json:"pinned"`
	Poll                 *DiscordPoll                       `json:"poll,omitempty"`
	Position             *int                               `json:"position,omitempty"`
	Reactions            *[]DiscordReaction                 `json:"reactions,omitempty"`
	Resolved             *DiscordResolvedData               `json:"resolved,omitempty"`
	ReferencedMessage    *DiscordMessage                    `json:"referenced_message,omitempty"`
	RoleSubscriptionData *DiscordRoleSubscriptionData       `json:"role_subscription_data,omitempty"`
	StickerItems         []DiscordMessageStickerItem        `json:"sticker_items,omitempty"`
	Thread               *DiscordChannel                    `json:"thread,omitempty"`
	Timestamp            *time.Time                         `json:"timestamp,omitempty"`
	TTS                  bool                               `json:"tts"`
	Type                 DiscordMessageType                 `json:"type"`
	WebhookID            *string                            `json:"webhook_id,omitempty"`
}

type DiscordMessageFlag uint64

const (
	DiscordMessageFlagCrossposted                      DiscordMessageFlag = 1 << 0
	DiscordMessageFlagIsCrosspost                      DiscordMessageFlag = 1 << 1
	DiscordMessageFlagSuppressEmbeds                   DiscordMessageFlag = 1 << 2
	DiscordMessageFlagSourceMessageDeleted             DiscordMessageFlag = 1 << 3
	DiscordMessageFlagUrgent                           DiscordMessageFlag = 1 << 4
	DiscordMessageFlagHasThread                        DiscordMessageFlag = 1 << 5
	DiscordMessageFlagEphemeral                        DiscordMessageFlag = 1 << 6
	DiscordMessageFlagLoading                          DiscordMessageFlag = 1 << 7
	DiscordMessageFlagFailedToMentionSomeRolesInThread DiscordMessageFlag = 1 << 8
	DiscordMessageFlagSuppressNotification             DiscordMessageFlag = 1 << 12
	DiscordMessageFlagIsVoiceMessage                   DiscordMessageFlag = 1 << 13
	DiscordMessageFlagHasSnapshot                      DiscordMessageFlag = 1 << 14
	DiscordMessageFlagIsComponentsV2                   DiscordMessageFlag = 1 << 15
)

type DiscordMessageType uint64

const (
	DiscordMessageTypeDefault                                 DiscordMessageType = 0
	DiscordMessageTypeRecipientAdd                            DiscordMessageType = 1
	DiscordMessageTypeRecipientRemove                         DiscordMessageType = 2
	DiscordMessageTypeCall                                    DiscordMessageType = 3
	DiscordMessageTypeChannelNameChange                       DiscordMessageType = 4
	DiscordMessageTypeChannelIconChange                       DiscordMessageType = 5
	DiscordMessageTypeChannelPinnedMessage                    DiscordMessageType = 6
	DiscordMessageTypeGuildMemberJoin                         DiscordMessageType = 7
	DiscordMessageTypeGuildBoost                              DiscordMessageType = 8
	DiscordMessageTypeGuildBoostTier1                         DiscordMessageType = 9
	DiscordMessageTypeGuildBoostTier2                         DiscordMessageType = 10
	DiscordMessageTypeGuildBoostTier3                         DiscordMessageType = 11
	DiscordMessageTypeChannelFollowAdd                        DiscordMessageType = 12
	DiscordMessageTypeGuildDiscoveryDisqualified              DiscordMessageType = 14
	DiscordMessageTypeGuildDiscoveryRequalified               DiscordMessageType = 15
	DiscordMessageTypeGuildDiscoveryGracePeriodInitialWarning DiscordMessageType = 16
	DiscordMessageTypeGuildDiscoveryGracePeriodFinalWarning   DiscordMessageType = 17
	DiscordMessageTypeThreadCreated                           DiscordMessageType = 18
	DiscordMessageTypeReply                                   DiscordMessageType = 19
	DiscordMessageTypeChatInputCommand                        DiscordMessageType = 20
	DiscordMessageTypeThreadStarterMessage                    DiscordMessageType = 21
	DiscordMessageTypeGuildInviteReminder                     DiscordMessageType = 22
	DiscordMessageTypeContextMenuCommand                      DiscordMessageType = 23
	DiscordMessageTypeAutoModerationAction                    DiscordMessageType = 24
	DiscordMessageTypeRoleSubscriptionPurchase                DiscordMessageType = 25
	DiscordMessageTypeInteractionPremiumUpsell                DiscordMessageType = 26
	DiscordMessageTypeStageStart                              DiscordMessageType = 27
	DiscordMessageTypeStageEnd                                DiscordMessageType = 28
	DiscordMessageTypeStageSpeaker                            DiscordMessageType = 29
	DiscordMessageTypeStageTopic                              DiscordMessageType = 31
	DiscordMessageTypeGuildApplicationPremiumSubscription     DiscordMessageType = 32
	DiscordMessageTypeGuildIncidentAlertModeEnabled           DiscordMessageType = 36
	DiscordMessageTypeGuildIncidentAlertModeDisabled          DiscordMessageType = 37
	DiscordMessageTypeReportRaid                              DiscordMessageType = 38
	DiscordMessageTypeReportFalseAlarm                        DiscordMessageType = 39
	DiscordMessageTypePurchaseNotification                    DiscordMessageType = 44
	DiscordMessageTypePollResult                              DiscordMessageType = 46
)

type DiscordAttachment struct {
	ID           DiscordSnowflake `json:"id"`
	Filename     string           `json:"filename"`
	Title        *string          `json:"title,omitempty"`
	Description  *string          `json:"description,omitempty"`
	ContentType  *string          `json:"content_type,omitempty"`
	Size         int              `json:"size"`
	URL          string           `json:"url"`
	ProxyURL     string           `json:"proxy_url,omitempty"`
	Height       *int             `json:"height,omitempty"`
	Width        *int             `json:"width,omitempty"`
	Ephemeral    *bool            `json:"ephemeral,omitempty"`
	DurationSecs *float64         `json:"duration_secs,omitempty"`
	Waveform     *string          `json:"waveform,omitempty"`
	Flags        *int             `json:"flags,omitempty"`
}

type DiscordEmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type DiscordEmbedImage struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   *int   `json:"height,omitempty"`
	Width    *int   `json:"width,omitempty"`
}

type DiscordEmbedThumbnail struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   *int   `json:"height,omitempty"`
	Width    *int   `json:"width,omitempty"`
}

type DiscordEmbedVideo struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   *int   `json:"height,omitempty"`
	Width    *int   `json:"width,omitempty"`
}

type DiscordEmbedProvider struct {
	Name *string `json:"name,omitempty"`
	URL  *string `json:"url,omitempty"`
}

type DiscordEmbedAuthor struct {
	Name         string  `json:"name"`
	URL          *string `json:"url"`
	IconURL      *string `json:"icon_url,omitempty"`
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"`
}

type DiscordEmbedFields struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline *bool  `json:"inline,omitempty"`
}

type DiscordEmbedType string

const (
	DiscordEmbedTypeRich       DiscordEmbedType = "rich"
	DiscordEmbedTypeImage      DiscordEmbedType = "image"
	DiscordEmbedTypeVideo      DiscordEmbedType = "video"
	DiscordEmbedTypeGIFV       DiscordEmbedType = "gifv"
	DiscordEmbedTypeArticle    DiscordEmbedType = "article"
	DiscordEmbedTypeLink       DiscordEmbedType = "link"
	DiscordEmbedTypePollResult DiscordEmbedType = "poll_result"
)

type DiscordEmbed struct {
	Title       *string                `json:"title,omitempty"`
	Type        *DiscordEmbedType      `json:"type,omitempty"`
	Description *string                `json:"description,omitempty"`
	URL         *string                `json:"url,omitempty"`
	Timestamp   *time.Time             `json:"timestamp,omitempty"`
	Color       *int                   `json:"color,omitempty"`
	Footer      *DiscordEmbedFooter    `json:"footer,omitempty"`
	Image       *DiscordEmbedImage     `json:"image,omitempty"`
	Thumbnail   *DiscordEmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *DiscordEmbedVideo     `json:"video,omitempty"`
	Provider    *DiscordEmbedProvider  `json:"provider,omitempty"`
	Author      *DiscordEmbedAuthor    `json:"author,omitempty"`
	Fields      []DiscordEmbedFields   `json:"fields,omitempty"`
}

type DiscordReactionCountDetails struct {
	Burst  int `json:"burst,omitempty"`
	Normal int `json:"normal,omitempty"`
}

type DiscordReaction struct {
	Count        int                         `json:"count"`
	CountDetails DiscordReactionCountDetails `json:"count_details,omitempty"`
	Me           bool                        `json:"me"`
	MeBurst      bool                        `json:"me_burst"`
	Emoji        DiscordEmoji                `json:"emoji"`
	BurstColors  []interface{}               `json:"burst_colors,omitempty"`
}

type DiscordEmoji struct {
	ID            *DiscordSnowflake `json:"id,omitempty"`
	Name          *string           `json:"name,omitempty"`
	Roles         []string          `json:"roles,omitempty"`
	Users         *DiscordUser      `json:"users,omitempty"`
	RequireColons *bool             `json:"require_colons,omitempty"`
	Managed       *bool             `json:"managed,omitempty"`
	Animated      *bool             `json:"animated,omitempty"`
	Available     *bool             `json:"available,omitempty"`
}

type DiscordPollLayoutType int

const (
	DiscordPollLayoutTypeDefault DiscordPollLayoutType = 0
)

type DiscordPollQuestion struct {
	Text  *string       `json:"text"`
	Emoji *DiscordEmoji `json:"emoji,omitempty"`
}

type DiscordPollAnswer struct {
	AnswerID  int                 `json:"answer_id"`
	PollMedia DiscordPollQuestion `json:"poll_media"`
}

type DiscordPollResultsAnswerCounts struct {
	ID      int  `json:"id"`
	Count   int  `json:"count"`
	MeVoted bool `json:"me_voted"`
}

type DiscordPollResults struct {
	IsFinalized  bool                             `json:"is_finalized"`
	AnswerCounts []DiscordPollResultsAnswerCounts `json:"answer_counts"`
}

type DiscordPoll struct {
	Question         *DiscordPollQuestion  `json:"question"`
	Answers          []DiscordPollAnswer   `json:"answers"`
	Expiry           *time.Time            `json:"expiry,omitempty"`
	AllowMultiselect bool                  `json:"allow_multiselect,omitempty"`
	LayoutType       DiscordPollLayoutType `json:"layout_type,omitempty"`
	Results          *DiscordPollResults   `json:"results,omitempty"`
}
