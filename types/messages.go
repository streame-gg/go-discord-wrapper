package types

import "time"

type DiscordSnowflake int

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
	IconHash    *string                         `json:"icon,omitempty"`
	ID          DiscordSnowflake                `json:"id"`
	Members     []*DiscordApplicationTeamMember `json:"members"`
	Name        string                          `json:"name"`
	OwnerUserID DiscordSnowflake                `json:"owner_user_id"`
}

type DiscordApplicationInstallParams struct {
	Scopes      []string `json:"scopes"`
	Permissions string   `json:"permissions"`
}

type DiscordApplication struct {
	ID                  DiscordSnowflake        `json:"id"`
	Name                string                  `json:"name"`
	IconHash            *string                 `json:"icon,omitempty"`
	Description         string                  `json:"description,omitempty"`
	RpcOrigins          []string                `json:"rpc_origins,omitempty"`
	BotPublic           bool                    `json:"bot_public,omitempty"`
	BotRequireCodeGrant bool                    `json:"bot_require_code_grant,omitempty"`
	TermsOfServiceURL   *string                 `json:"terms_of_service_url,omitempty"`
	PrivacyPolicyURL    *string                 `json:"privacy_policy_url,omitempty"`
	Owner               *DiscordUser            `json:"owner,omitempty"`
	VerifyKey           *string                 `json:"verify_key,omitempty"`
	Team                *DiscordApplicationTeam `json:"team,omitempty"`
	GuildID             *DiscordSnowflake       `json:"guild_id,omitempty"`
	//TODO
	Guild                             any                                  `json:"guild,omitempty"`
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

type DiscordMessageResolved struct {
	Users    map[DiscordSnowflake]*DiscordUser    `json:"users"`
	Members  map[DiscordSnowflake]*GuildMember    `json:"members,omitempty"`
	Messages map[DiscordSnowflake]*DiscordMessage `json:"messages,omitempty"`
	//TODO
	Channels map[DiscordSnowflake]*any `json:"channels,omitempty"`
	//TODO
	Roles       map[DiscordSnowflake]*any               `json:"roles,omitempty"`
	Attachments map[DiscordSnowflake]*DiscordAttachment `json:"attachments,omitempty"`
}

type DiscordMessage struct {
	Activity      *DiscordActivity     `json:"activity,omitempty"`
	Application   *DiscordApplication  `json:"application,omitempty"`
	ApplicationID *string              `json:"application_id,omitempty"`
	Attachments   []*DiscordAttachment `json:"attachments,omitempty"`
	Author        *DiscordUser         `json:"author"`
	Call          *DiscordCall         `json:"call,omitempty"`
	ChannelID     string               `json:"channel_id"`
	ChannelType   int                  `json:"channel_type"`
	//TODO
	Components      []*any          `json:"components"`
	Content         string          `json:"content"`
	EditedTimestamp *int64          `json:"edited_timestamp,omitempty"`
	Embeds          []*DiscordEmbed `json:"embeds,omitempty"`
	Flags           int64           `json:"flags"`
	ID              string          `json:"id"`
	//TODO
	InteractionMetadata *any                            `json:"interaction_metadata,omitempty"`
	MentionEveryone     bool                            `json:"mention_everyone"`
	MentionChannels     *[]DiscordMessageChannelMention `json:"mention_channels,omitempty"`
	MentionRoles        []string                        `json:"mention_roles"`
	//TODO
	MessageReference *any `json:"message_reference,omitempty"`
	//TODO
	MessageSnapshots []*any      `json:"message_snapshots,omitempty"`
	Nonce            interface{} `json:"nonce,omitempty"`
	Pinned           bool        `json:"pinned"`
	//TODO
	Poll      *any                `json:"poll,omitempty"`
	Position  *int                `json:"position,omitempty"`
	Reactions *[]*DiscordReaction `json:"reactions,omitempty"`
	//TODO
	Resolved             *any                         `json:"resolved,omitempty"`
	ReferencedMessage    *DiscordMessage              `json:"referenced_message,omitempty"`
	RoleSubscriptionData *DiscordRoleSubscriptionData `json:"role_subscription_data,omitempty"`
	StickerItems         []*DiscordMessageStickerItem `json:"sticker_items,omitempty"`
	//TODO
	Thread    *any       `json:"thread,omitempty"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
	TTS       bool       `json:"tts"`
	Type      int        `json:"type"`
	WebhookID *string    `json:"webhook_id,omitempty"`
}

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
	Fields      []*DiscordEmbedFields  `json:"fields,omitempty"`
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
