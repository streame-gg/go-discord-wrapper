package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-object-resolved-data-structure
type ResolvedData struct {
	Users       map[Snowflake]User        `json:"users,omitempty"`
	Members     map[Snowflake]GuildMember `json:"members,omitempty"`
	Messages    map[Snowflake]Message     `json:"messages,omitempty"`
	Channels    map[Snowflake]Channel     `json:"channels,omitempty"`
	Roles       map[Snowflake]Role        `json:"roles,omitempty"`
	Attachments map[Snowflake]Attachment  `json:"attachments,omitempty"`
}

// MessagePin is one entry returned by Get Channel Pins: a pinned message and
// when it was pinned. See
// https://docs.discord.com/developers/resources/message#message-pin-object.
type MessagePin struct {
	// PinnedAt is when the message was pinned.
	PinnedAt time.Time `json:"pinned_at"`
	// Message is the pinned message.
	Message Message `json:"message"`
}

// ChannelPins is the paginated response from Get Channel Pins. Items are ordered
// newest-pin-first; HasMore reports whether older pins remain (fetch them with a
// before cursor set to the last item's PinnedAt).
//
// https://docs.discord.com/developers/resources/message#get-channel-pins
type ChannelPins struct {
	Items   []MessagePin `json:"items"`
	HasMore bool         `json:"has_more"`
}

// https://docs.discord.com/developers/resources/message#message-object-message-activity-types
type MessageActivityType uint64

const (
	MessageActivityTypeJoin        MessageActivityType = 1
	MessageActivityTypeSpectate    MessageActivityType = 2
	MessageActivityTypeListen      MessageActivityType = 3
	MessageActivityTypeJoinRequest MessageActivityType = 4
)

// https://docs.discord.com/developers/resources/message#message-object-message-activity-structure
type MessageActivity struct {
	Type    MessageActivityType `json:"type"`
	PartyID string              `json:"party_id,omitempty"`
}

// https://docs.discord.com/developers/resources/message#message-object
type Message struct {
	hClient EntityClient

	Activity      *MessageActivity `json:"activity,omitempty"`
	Application   *Application     `json:"application,omitempty"`
	ApplicationID *Snowflake       `json:"application_id,omitempty"`
	Attachments   []Attachment     `json:"attachments,omitempty"`
	Author        User             `json:"author"`
	Call          *Call            `json:"call,omitempty"`
	ChannelID     Snowflake        `json:"channel_id"`
	// Components has custom marshal und unmarshal due to some issues with the default MarshalJSON and UnmarshalJSON.
	Components          []AnyComponent              `json:"-"`
	Content             string                      `json:"content"`
	EditedTimestamp     *time.Time                  `json:"edited_timestamp"`
	Embeds              []Embed                     `json:"embeds,omitempty"`
	Flags               MessageFlag                 `json:"flags,omitempty"`
	ID                  Snowflake                   `json:"id"`
	InteractionMetadata *MessageInteractionMetadata `json:"interaction_metadata,omitempty"`
	MentionEveryone     bool                        `json:"mention_everyone"`
	MentionChannels     []MessageChannelMention     `json:"mention_channels,omitempty"`
	MentionRoles        []Snowflake                 `json:"mention_roles"`
	Mentions            []User                      `json:"mentions"`
	MessageReference    *MessageMessageReference    `json:"message_reference,omitempty"`
	MessageSnapshots    []MessageMessageSnapshot    `json:"message_snapshots,omitempty"`
	// Nonce is either an integer or a string
	Nonce                interface{}           `json:"nonce,omitempty"`
	Pinned               bool                  `json:"pinned"`
	Poll                 *Poll                 `json:"poll,omitempty"`
	Position             *int                  `json:"position,omitempty"`
	Reactions            []Reaction            `json:"reactions,omitempty"`
	Resolved             *ResolvedData         `json:"resolved,omitempty"`
	ReferencedMessage    *Message              `json:"referenced_message,omitempty"`
	RoleSubscriptionData *RoleSubscriptionData `json:"role_subscription_data,omitempty"`
	StickerItems         []MessageStickerItem  `json:"sticker_items,omitempty"`
	Thread               *Channel              `json:"thread,omitempty"`
	Timestamp            time.Time             `json:"timestamp"`
	TTS                  bool                  `json:"tts"`
	Type                 MessageType           `json:"type"`
	WebhookID            *Snowflake            `json:"webhook_id,omitempty"`
	SharedClientTheme    *SharedClientTheme    `json:"shared_client_theme,omitempty"`
}

// https://docs.discord.com/developers/resources/message#role-subscription-data-object
type RoleSubscriptionData struct {
	RoleSubscriptionListingID Snowflake `json:"role_subscription_listing_id"`
	TierName                  string    `json:"tier_name"`
	TotalMonthsSubscribed     int       `json:"total_months_subscribed"`
	IsRenewal                 bool      `json:"is_renewal"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	aux := &struct {
		Components []json.RawMessage `json:"components"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Unmarshal components as RawComponent
	for _, comp := range aux.Components {
		var raw RawComponent
		if err := json.Unmarshal(comp, &raw); err != nil {
			return fmt.Errorf("failed to unmarshal component: %w", err)
		}
		m.Components = append(m.Components, &raw)
	}

	return nil
}

func (m *Message) MarshalJSON() ([]byte, error) {
	type Alias Message
	return json.Marshal(&struct {
		Components []AnyComponent `json:"components,omitempty"`
		*Alias
	}{
		Components: m.Components,
		Alias:      (*Alias)(m),
	})
}

// https://docs.discord.com/developers/resources/message#message-object-message-flags
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

// https://docs.discord.com/developers/resources/message#message-object-message-types
type MessageType uint64

const (
	MessageTypeDefault                                 MessageType = 0
	MessageTypeRecipientAdd                            MessageType = 1
	MessageTypeRecipientRemove                         MessageType = 2
	MessageTypeCall                                    MessageType = 3
	MessageTypeChannelNameChange                       MessageType = 4
	MessageTypeChannelIconChange                       MessageType = 5
	MessageTypeChannelPinnedMessage                    MessageType = 6
	MessageTypeGuildUserJoin                           MessageType = 7
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

// https://docs.discord.com/developers/resources/message#message-call-object
type Call struct {
	Participants   []Snowflake `json:"participants"`
	EndedTimestamp *time.Time  `json:"ended_timestamp,omitempty"`
}

// https://docs.discord.com/developers/resources/message#channel-mention-object
type MessageChannelMention struct {
	ID      Snowflake   `json:"id"`
	GuildID Snowflake   `json:"guild_id"`
	Type    ChannelType `json:"type"`
	Name    string      `json:"name"`
}

// https://docs.discord.com/developers/resources/sticker#sticker-item-object
type MessageStickerItem struct {
	ID         Snowflake         `json:"id"`
	Name       string            `json:"name"`
	FormatType StickerFormatType `json:"format_type"`
}

// https://docs.discord.com/developers/resources/message#message-snapshot-object
type MessageMessageSnapshot struct {
	Message Message `json:"message"`
}

// https://docs.discord.com/developers/resources/message#message-reference-types
type MessageReferenceType int

const (
	MessageReferenceTypeDefault MessageReferenceType = 0
	MessageReferenceTypeForward MessageReferenceType = 1
)

// https://docs.discord.com/developers/resources/message#message-reference-structure
type MessageMessageReference struct {
	Type            *MessageReferenceType `json:"type,omitempty"`
	MessageID       *Snowflake            `json:"message_id,omitempty"`
	ChannelID       *Snowflake            `json:"channel_id,omitempty"`
	GuildID         *Snowflake            `json:"guild_id,omitempty"`
	FailIfNotExists *bool                 `json:"fail_if_not_exists,omitempty"`
}

// https://docs.discord.com/developers/resources/message#message-interaction-metadata-object-application-command-interaction-metadata-structure
type MessageInteractionMetadataApplicationCommand struct {
	ID                           Snowflake                                                 `json:"id"`
	Type                         InteractionType                                           `json:"type"`
	User                         User                                                      `json:"user"`
	AuthorizingIntegrationOwners map[ApplicationIntegrationType]ApplicationIntegrationType `json:"authorizing_integration_owners"`
	OriginalResponseMessageID    *Snowflake                                                `json:"original_response_message_id,omitempty"`
	TargetUser                   *User                                                     `json:"target_user,omitempty"`
	TargetMessageID              *Snowflake                                                `json:"target_message_id,omitempty"`
}

// https://docs.discord.com/developers/resources/message#message-interaction-metadata-object-message-component-interaction-metadata-structure
type MessageInteractionMetadataMessageComponent struct {
	ID                           Snowflake                                  `json:"id"`
	Type                         InteractionType                            `json:"type"`
	User                         User                                       `json:"user"`
	AuthorizingIntegrationOwners map[ApplicationIntegrationType]interface{} `json:"authorizing_integration_owners"`
	OriginalResponseMessageID    *Snowflake                                 `json:"original_response_message_id,omitempty"`
	InteractedMessageID          Snowflake                                  `json:"interacted_message_id"`
}

// https://docs.discord.com/developers/resources/message#message-interaction-metadata-object
type AnyMessageInteractionMetadata interface{}

// https://docs.discord.com/developers/resources/message#message-interaction-metadata-object
type MessageInteractionMetadata struct {
	Value AnyMessageInteractionMetadata
}

// https://docs.discord.com/developers/resources/message#message-interaction-metadata-object
type AnyMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata interface{}

// https://docs.discord.com/developers/resources/message#message-interaction-metadata-object
type MessageInteractionMetadataModalSubmitTriggering struct {
	AnyMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata
}

func (d *MessageInteractionMetadataModalSubmitTriggering) UnmarshalJSON(data []byte) error {
	var a MessageInteractionMetadataApplicationCommand
	if err := json.Unmarshal(data, &a); err == nil && !a.ID.IsEmpty() {
		d.AnyMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata = &a
		return nil
	}

	var b MessageInteractionMetadataMessageComponent
	if err := json.Unmarshal(data, &b); err == nil && !b.ID.IsEmpty() {
		d.AnyMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata = &b
		return nil
	}

	return nil
}

func (d *MessageInteractionMetadata) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		d.Value = nil
		return nil
	}

	var a MessageInteractionMetadataApplicationCommand
	if err := json.Unmarshal(data, &a); err == nil && !a.ID.IsEmpty() {
		d.Value = &a
		return nil
	}

	var b MessageInteractionMetadataMessageComponent
	if err := json.Unmarshal(data, &b); err == nil && !b.ID.IsEmpty() {
		d.Value = &b
		return nil
	}

	var c MessageInteractionMetadataModalSubmit
	if err := json.Unmarshal(data, &c); err == nil && !c.ID.IsEmpty() {
		d.Value = &c
		return nil
	}

	return fmt.Errorf("unknown MessageInteractionMetadata: %s", string(data))
}

// https://docs.discord.com/developers/resources/message#message-interaction-metadata-object-modal-submit-interaction-metadata-structure
type MessageInteractionMetadataModalSubmit struct {
	ID                            Snowflake                                       `json:"id"`
	Type                          InteractionType                                 `json:"type"`
	User                          User                                            `json:"user"`
	AuthorizingIntegrationOwners  map[ApplicationIntegrationType]interface{}      `json:"authorizing_integration_owners"`
	OriginalResponseMessageID     *Snowflake                                      `json:"original_response_message_id,omitempty"`
	TriggeringInteractionMetadata MessageInteractionMetadataModalSubmitTriggering `json:"triggering_interaction_metadata"`
}

// SharedClientTheme https://docs.discord.com/developers/resources/message#shared-client-theme-object
type SharedClientTheme struct {
	Colors        []string  `json:"colors"`
	GradientAngle int       `json:"gradient_angle"`
	BaseMix       int       `json:"base_mix"`
	BaseTheme     BaseTheme `json:"base_theme,omitempty"`
}

// BaseTheme https://docs.discord.com/developers/resources/message#base-theme-types
// BaseThemeUnset is equal to BaseThemeDark.
type BaseTheme int

const (
	BaseThemeUnset BaseTheme = iota
	BaseThemeDark
	BaseThemeLight
	BaseThemeDarker
	BaseThemeMidnight
)

// GuildSearchResponse is returned by the guild message search endpoint.
// Messages is a list of context windows; each window is a list of adjacent messages.
//
// https://docs.discord.com/developers/resources/message#search-guild-messages
type GuildSearchResponse struct {
	// Messages is a double array because it used to provide surrounding context. The surrounding context is no longer returned.
	Messages                 [][]Message   `json:"messages"`
	DoingDeepHistoricalIndex bool          `json:"doing_deep_historical_index"`
	TotalResults             int           `json:"total_results"`
	Threads                  []Channel     `json:"threads,omitempty"`
	Members                  []GuildMember `json:"members,omitempty"`
	DocumentsIndexed         *int          `json:"documents_indexed,omitempty"`
}

// ThreadSearchResponse is returned by the thread search endpoint.
//
// https://docs.discord.com/developers/resources/channel#list-public-archived-threads
type ThreadSearchResponse struct {
	Threads []Channel      `json:"threads"`
	Members []ThreadMember `json:"members"`
	HasMore bool           `json:"has_more"`
}
