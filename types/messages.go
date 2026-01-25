package types

import (
	"time"
)

type Snowflake string

func (s Snowflake) ToString() string {
	return string(s)
}

type Activity struct {
	Type    int     `json:"type"`
	PartyID *string `json:"party_id,omitempty"`
}

type ApplicationEventWebhookStatus int

const (
	ApplicationEventWebhookStatusDisabled   ApplicationEventWebhookStatus = 0
	ApplicationEventWebhookStatusEnabled    ApplicationEventWebhookStatus = 1
	ApplicationEventWebhookStatusDisabledBy ApplicationEventWebhookStatus = 2
)

type ApplicationTeamMemberMembershipState int

const (
	ApplicationTeamMemberMembershipStateInvited  ApplicationTeamMemberMembershipState = 1
	ApplicationTeamMemberMembershipStateAccepted ApplicationTeamMemberMembershipState = 2
)

type ApplicationTeamMember struct {
	MembershipState ApplicationTeamMemberMembershipState `json:"membership_state"`
	TeamID          *string                              `json:"team_id,omitempty"`
	User            *User                                `json:"user,omitempty"`
	Role            string                               `json:"role"`
}

type ApplicationTeam struct {
	IconHash    *string                 `json:"icon,omitempty"`
	ID          Snowflake               `json:"id"`
	Members     []ApplicationTeamMember `json:"members"`
	Name        string                  `json:"name"`
	OwnerUserID Snowflake               `json:"owner_user_id"`
}

type ApplicationInstallParams struct {
	Scopes      []string `json:"scopes"`
	Permissions string   `json:"permissions"`
}

type Application struct {
	ID                                Snowflake                     `json:"id"`
	Name                              string                        `json:"name"`
	IconHash                          *string                       `json:"icon,omitempty"`
	Description                       string                        `json:"description,omitempty"`
	RpcOrigins                        []string                      `json:"rpc_origins,omitempty"`
	BotPublic                         bool                          `json:"bot_public,omitempty"`
	BotRequireCodeGrant               bool                          `json:"bot_require_code_grant,omitempty"`
	TermsOfServiceURL                 *string                       `json:"terms_of_service_url,omitempty"`
	PrivacyPolicyURL                  *string                       `json:"privacy_policy_url,omitempty"`
	Owner                             *User                         `json:"owner,omitempty"`
	VerifyKey                         *string                       `json:"verify_key,omitempty"`
	Team                              *ApplicationTeam              `json:"team,omitempty"`
	GuildID                           *Snowflake                    `json:"guild_id,omitempty"`
	Guild                             Guild                         `json:"guild,omitempty"`
	PrimarySKUID                      *string                       `json:"primary_sku_id,omitempty"`
	Slug                              *string                       `json:"slug,omitempty"`
	CoverImage                        *string                       `json:"cover_image,omitempty"`
	Flags                             *int                          `json:"flags,omitempty"`
	ApproximateGuildCount             *int                          `json:"approximate_guild_count,omitempty"`
	ApproximateUserInstallCount       *int                          `json:"approximate_user_install_count,omitempty"`
	ApproximateUserAuthorizationCount *int                          `json:"approximate_user_authorization_count,omitempty"`
	RedirectURIs                      *[]string                     `json:"redirect_uris,omitempty"`
	InteractionEndpointURL            *string                       `json:"interaction_endpoint_url,omitempty"`
	RoleConnectionsVerificationURL    *string                       `json:"role_connections_verification_url,omitempty"`
	EventWebhooksURL                  *string                       `json:"event_webhooks_url,omitempty"`
	EventWebhookStatus                ApplicationEventWebhookStatus `json:"event_webhook_status,omitempty"`
	EventWebhooksTypes                *[]string                     `json:"event_webhooks_types,omitempty"`
	Tags                              *[]string                     `json:"tags,omitempty"`
	InstallParams                     *ApplicationInstallParams     `json:"install_params,omitempty"`
	IntegrationTypesConfig            *interface{}                  `json:"integration_types_config,omitempty"`
	CustomInstallURL                  *string                       `json:"custom_install_url,omitempty"`
}

type Call struct {
	Participants   []Snowflake `json:"participants"`
	EndedTimestamp *time.Time  `json:"ended_timestamp,omitempty"`
}

type MessageChannelMention struct {
	ID      string `json:"id"`
	GuildID string `json:"guild_id"`
	Type    int    `json:"type"`
	Name    string `json:"name"`
}

type MessageStickerItemFormatType int

const (
	MessageStickerItemFormatTypePNG    MessageStickerItemFormatType = 1
	MessageStickerItemFormatTypeAPNG   MessageStickerItemFormatType = 2
	MessageStickerItemFormatTypeLottie MessageStickerItemFormatType = 3
	MessageStickerItemFormatTypeGIF    MessageStickerItemFormatType = 4
)

type MessageStickerItem struct {
	ID         Snowflake                    `json:"id"`
	Name       string                       `json:"name"`
	FormatType MessageStickerItemFormatType `json:"format_type"`
}

type RoleSubscriptionData struct {
	RoleSubscriptionListingID Snowflake `json:"role_subscription_listing_id"`
	TierName                  string    `json:"tier_name"`
	TotalMonthsSubscribed     int       `json:"total_months_subscribed"`
	IsRenewal                 bool      `json:"is_renewal"`
}

type ResolvedData struct {
	Users       map[Snowflake]*User        `json:"users"`
	Members     map[Snowflake]*GuildMember `json:"members,omitempty"`
	Messages    map[Snowflake]*Message     `json:"messages,omitempty"`
	Channels    map[Snowflake]*Channel     `json:"channels,omitempty"`
	Roles       map[Snowflake]*Role        `json:"roles,omitempty"`
	Attachments map[Snowflake]*Attachment  `json:"attachments,omitempty"`
}

type PartialMessage struct {
	Type            MessageType          `json:"type"`
	Content         string               `json:"content"`
	Embeds          []Embed              `json:"embeds,omitempty"`
	Attachments     []Attachment         `json:"attachments,omitempty"`
	Timestamp       *time.Time           `json:"timestamp,omitempty"`
	EditedTimestamp *time.Time           `json:"edited_timestamp,omitempty"`
	Flags           MessageFlag          `json:"flags,omitempty"`
	Mentions        *[]any               `json:"mentions"`
	MentionRoles    []string             `json:"mention_roles"`
	StickerItems    []MessageStickerItem `json:"sticker_items,omitempty"`
	Components      []AnyComponent       `json:"components,omitempty"`
}

type MessageMessageSnapshot struct {
	Message PartialMessage `json:"message,omitempty"`
}

type MessageMessageReferenceType int

const (
	MessageMessageReferenceTypeDefault MessageMessageReferenceType = 0
	MessageMessageReferenceTypeForward MessageMessageReferenceType = 1
)

type Message struct {
	Activity             *Activity                   `json:"activity,omitempty"`
	Application          *Application                `json:"application,omitempty"`
	ApplicationID        *string                     `json:"application_id,omitempty"`
	Attachments          []Attachment                `json:"attachments,omitempty"`
	Author               *User                       `json:"author"`
	Call                 *Call                       `json:"call,omitempty"`
	ChannelID            Snowflake                   `json:"channel_id"`
	ChannelType          ChannelType                 `json:"channel_type"`
	Components           []ActionRow                 `json:"components"`
	Content              string                      `json:"content"`
	EditedTimestamp      *time.Time                  `json:"edited_timestamp,omitempty"`
	Embeds               []Embed                     `json:"embeds,omitempty"`
	Flags                MessageFlag                 `json:"flags"`
	ID                   Snowflake                   `json:"id"`
	InteractionMetadata  *MessageInteractionMetadata `json:"interaction_metadata,omitempty"`
	MentionEveryone      bool                        `json:"mention_everyone"`
	MentionChannels      *[]MessageChannelMention    `json:"mention_channels,omitempty"`
	MentionRoles         []string                    `json:"mention_roles"`
	MessageReference     *MessageMessageReference    `json:"message_reference,omitempty"`
	MessageSnapshots     []MessageMessageSnapshot    `json:"message_snapshots,omitempty"`
	Nonce                interface{}                 `json:"nonce,omitempty"`
	Pinned               bool                        `json:"pinned"`
	Poll                 *Poll                       `json:"poll,omitempty"`
	Position             *int                        `json:"position,omitempty"`
	Reactions            *[]Reaction                 `json:"reactions,omitempty"`
	Resolved             *ResolvedData               `json:"resolved,omitempty"`
	ReferencedMessage    *Message                    `json:"referenced_message,omitempty"`
	RoleSubscriptionData *RoleSubscriptionData       `json:"role_subscription_data,omitempty"`
	StickerItems         []MessageStickerItem        `json:"sticker_items,omitempty"`
	Thread               *Channel                    `json:"thread,omitempty"`
	Timestamp            *time.Time                  `json:"timestamp,omitempty"`
	TTS                  bool                        `json:"tts"`
	Type                 MessageType                 `json:"type"`
	WebhookID            *string                     `json:"webhook_id,omitempty"`
}

type MessageFlag uint64

const (
	MessageFlagCrossposted                      MessageFlag = 1 << 0
	MessageFlagIsCrosspost                      MessageFlag = 1 << 1
	MessageFlagSuppressEmbeds                   MessageFlag = 1 << 2
	MessageFlagSourceMessageDeleted             MessageFlag = 1 << 3
	MessageFlagUrgent                           MessageFlag = 1 << 4
	MessageFlagHasThread                        MessageFlag = 1 << 5
	MessageFlagEphemeral                        MessageFlag = 1 << 6
	MessageFlagLoading                          MessageFlag = 1 << 7
	MessageFlagFailedToMentionSomeRolesInThread MessageFlag = 1 << 8
	MessageFlagSuppressNotification             MessageFlag = 1 << 12
	MessageFlagIsVoiceMessage                   MessageFlag = 1 << 13
	MessageFlagHasSnapshot                      MessageFlag = 1 << 14
	MessageFlagIsComponentsV2                   MessageFlag = 1 << 15
)

type MessageType uint64

const (
	MessageTypeDefault                                 MessageType = 0
	MessageTypeRecipientAdd                            MessageType = 1
	MessageTypeRecipientRemove                         MessageType = 2
	MessageTypeCall                                    MessageType = 3
	MessageTypeChannelNameChange                       MessageType = 4
	MessageTypeChannelIconChange                       MessageType = 5
	MessageTypeChannelPinnedMessage                    MessageType = 6
	MessageTypeGuildMemberJoin                         MessageType = 7
	MessageTypeGuildBoost                              MessageType = 8
	MessageTypeGuildBoostTier1                         MessageType = 9
	MessageTypeGuildBoostTier2                         MessageType = 10
	MessageTypeGuildBoostTier3                         MessageType = 11
	MessageTypeChannelFollowAdd                        MessageType = 12
	MessageTypeGuildDiscoveryDisqualified              MessageType = 14
	MessageTypeGuildDiscoveryRequalified               MessageType = 15
	MessageTypeGuildDiscoveryGracePeriodInitialWarning MessageType = 16
	MessageTypeGuildDiscoveryGracePeriodFinalWarning   MessageType = 17
	MessageTypeThreadCreated                           MessageType = 18
	MessageTypeReply                                   MessageType = 19
	MessageTypeChatInputCommand                        MessageType = 20
	MessageTypeThreadStarterMessage                    MessageType = 21
	MessageTypeGuildInviteReminder                     MessageType = 22
	MessageTypeContextMenuCommand                      MessageType = 23
	MessageTypeAutoModerationAction                    MessageType = 24
	MessageTypeRoleSubscriptionPurchase                MessageType = 25
	MessageTypeInteractionPremiumUpsell                MessageType = 26
	MessageTypeStageStart                              MessageType = 27
	MessageTypeStageEnd                                MessageType = 28
	MessageTypeStageSpeaker                            MessageType = 29
	MessageTypeStageTopic                              MessageType = 31
	MessageTypeGuildApplicationPremiumSubscription     MessageType = 32
	MessageTypeGuildIncidentAlertModeEnabled           MessageType = 36
	MessageTypeGuildIncidentAlertModeDisabled          MessageType = 37
	MessageTypeReportRaid                              MessageType = 38
	MessageTypeReportFalseAlarm                        MessageType = 39
	MessageTypePurchaseNotification                    MessageType = 44
	MessageTypePollResult                              MessageType = 46
)

type Attachment struct {
	ID           Snowflake `json:"id"`
	Filename     string    `json:"filename"`
	Title        *string   `json:"title,omitempty"`
	Description  *string   `json:"description,omitempty"`
	ContentType  *string   `json:"content_type,omitempty"`
	Size         int       `json:"size"`
	URL          string    `json:"url"`
	ProxyURL     string    `json:"proxy_url,omitempty"`
	Height       *int      `json:"height,omitempty"`
	Width        *int      `json:"width,omitempty"`
	Ephemeral    *bool     `json:"ephemeral,omitempty"`
	DurationSecs *float64  `json:"duration_secs,omitempty"`
	Waveform     *string   `json:"waveform,omitempty"`
	Flags        *int      `json:"flags,omitempty"`
}

type EmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type EmbedImage struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   *int   `json:"height,omitempty"`
	Width    *int   `json:"width,omitempty"`
}

type EmbedThumbnail struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   *int   `json:"height,omitempty"`
	Width    *int   `json:"width,omitempty"`
}

type EmbedVideo struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   *int   `json:"height,omitempty"`
	Width    *int   `json:"width,omitempty"`
}

type EmbedProvider struct {
	Name *string `json:"name,omitempty"`
	URL  *string `json:"url,omitempty"`
}

type EmbedAuthor struct {
	Name         string  `json:"name"`
	URL          *string `json:"url"`
	IconURL      *string `json:"icon_url,omitempty"`
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"`
}

type EmbedFields struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline *bool  `json:"inline,omitempty"`
}

type EmbedType string

const (
	EmbedTypeRich       EmbedType = "rich"
	EmbedTypeImage      EmbedType = "image"
	EmbedTypeVideo      EmbedType = "video"
	EmbedTypeGIFV       EmbedType = "gifv"
	EmbedTypeArticle    EmbedType = "article"
	EmbedTypeLink       EmbedType = "link"
	EmbedTypePollResult EmbedType = "poll_result"
)

type Embed struct {
	Title       *string         `json:"title,omitempty"`
	Type        *EmbedType      `json:"type,omitempty"`
	Description *string         `json:"description,omitempty"`
	URL         *string         `json:"url,omitempty"`
	Timestamp   *time.Time      `json:"timestamp,omitempty"`
	Color       *int            `json:"color,omitempty"`
	Footer      *EmbedFooter    `json:"footer,omitempty"`
	Image       *EmbedImage     `json:"image,omitempty"`
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *EmbedVideo     `json:"video,omitempty"`
	Provider    *EmbedProvider  `json:"provider,omitempty"`
	Author      *EmbedAuthor    `json:"author,omitempty"`
	Fields      []EmbedFields   `json:"fields,omitempty"`
}

type ReactionCountDetails struct {
	Burst  int `json:"burst,omitempty"`
	Normal int `json:"normal,omitempty"`
}

type Reaction struct {
	Count        int                  `json:"count"`
	CountDetails ReactionCountDetails `json:"count_details,omitempty"`
	Me           bool                 `json:"me"`
	MeBurst      bool                 `json:"me_burst"`
	Emoji        Emoji                `json:"emoji"`
	BurstColors  []interface{}        `json:"burst_colors,omitempty"`
}

type Emoji struct {
	ID            Snowflake `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	Roles         []string  `json:"roles,omitempty"`
	Users         *User     `json:"users,omitempty"`
	RequireColons *bool     `json:"require_colons,omitempty"`
	Managed       *bool     `json:"managed,omitempty"`
	Animated      bool      `json:"animated,omitempty"`
	Available     *bool     `json:"available,omitempty"`
}

type PollLayoutType int

const (
	PollLayoutTypeDefault PollLayoutType = 0
)

type PollQuestion struct {
	Text  *string `json:"text"`
	Emoji *Emoji  `json:"emoji,omitempty"`
}

type PollAnswer struct {
	AnswerID  int          `json:"answer_id"`
	PollMedia PollQuestion `json:"poll_media"`
}

type PollResultsAnswerCounts struct {
	ID      int  `json:"id"`
	Count   int  `json:"count"`
	MeVoted bool `json:"me_voted"`
}

type PollResults struct {
	IsFinalized  bool                      `json:"is_finalized"`
	AnswerCounts []PollResultsAnswerCounts `json:"answer_counts"`
}

type Poll struct {
	Question         *PollQuestion  `json:"question"`
	Answers          []PollAnswer   `json:"answers"`
	Expiry           *time.Time     `json:"expiry,omitempty"`
	AllowMultiselect bool           `json:"allow_multiselect,omitempty"`
	LayoutType       PollLayoutType `json:"layout_type,omitempty"`
	Results          *PollResults   `json:"results,omitempty"`
}

type PollRequest struct {
	Question         *PollQuestion  `json:"question"`
	Answers          []PollAnswer   `json:"answers"`
	Duration         *time.Time     `json:"duration,omitempty"`
	AllowMultiselect bool           `json:"allow_multiselect,omitempty"`
	LayoutType       PollLayoutType `json:"layout_type,omitempty"`
}
