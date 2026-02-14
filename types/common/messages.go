package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type ResolvedData struct {
	Users       map[Snowflake]*User        `json:"users"`
	Members     map[Snowflake]*GuildMember `json:"members,omitempty"`
	Messages    map[Snowflake]*Message     `json:"messages,omitempty"`
	Channels    map[Snowflake]*Channel     `json:"channels,omitempty"`
	Roles       map[Snowflake]*Role        `json:"roles,omitempty"`
	Attachments map[Snowflake]*Attachment  `json:"attachments,omitempty"`
}

type Message struct {
	Activity             *Activity                   `json:"activity,omitempty"`
	Application          *Application                `json:"application,omitempty"`
	ApplicationID        *string                     `json:"application_id,omitempty"`
	Attachments          []Attachment                `json:"attachments,omitempty"`
	Author               *User                       `json:"author"`
	Call                 *Call                       `json:"call,omitempty"`
	ChannelID            Snowflake                   `json:"channel_id"`
	ChannelType          ChannelType                 `json:"channel_type"`
	Components           []AnyComponent              `json:"components"`
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

type MessageMessageSnapshot struct {
	Message PartialMessage `json:"message,omitempty"`
}

type MessageMessageReferenceType int

const (
	MessageMessageReferenceTypeDefault MessageMessageReferenceType = 0
	MessageMessageReferenceTypeForward MessageMessageReferenceType = 1
)

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

type MessageMessageReference struct {
	Type            *MessageMessageReferenceType `json:"type,omitempty"`
	MessageID       *Snowflake                   `json:"message_id,omitempty"`
	ChannelID       *Snowflake                   `json:"channel_id,omitempty"`
	GuildID         *Snowflake                   `json:"guild_id,omitempty"`
	FailIfNotExists *bool                        `json:"fail_if_not_exists,omitempty"`
}

type MessageInteractionMetadataApplicationCommand struct {
	ID                           Snowflake                                             `json:"id"`
	Type                         InteractionType                                       `json:"type"`
	User                         User                                                  `json:"user,omitempty"`
	AuthorizingIntegrationOwners map[InteractionApplicationIntegrationType]interface{} `json:"authorizing_integration_owners,omitempty"`
	OriginalResponseMessageID    *Snowflake                                            `json:"original_response_message_id,omitempty"`
	TargetUser                   *User                                                 `json:"target_user,omitempty"`
	TargetMessageID              *Snowflake                                            `json:"target_message_id,omitempty"`
}

type MessageInteractionMetadataMessageComponent struct {
	ID                           Snowflake                                             `json:"id"`
	Type                         InteractionType                                       `json:"type"`
	User                         User                                                  `json:"user,omitempty"`
	AuthorizingIntegrationOwners map[InteractionApplicationIntegrationType]interface{} `json:"authorizing_integration_owners,omitempty"`
	OriginalResponseMessageID    *Snowflake                                            `json:"original_response_message_id,omitempty"`
	InteractedMessageID          *Snowflake                                            `json:"interacted_message_id,omitempty"`
}

type AnyMessageInteractionMetadata interface{}

type MessageInteractionMetadata struct {
	Value AnyMessageInteractionMetadata
}

type AnyMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata interface{}

type MessageInteractionMetadataModalSubmitTriggering struct {
	AnyMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata
}

func (d *MessageInteractionMetadataModalSubmitTriggering) UnmarshalJSON(data []byte) error {
	var a MessageInteractionMetadataApplicationCommand
	if err := json.Unmarshal(data, &a); err == nil && a.ID != "" {
		d.AnyMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata = &a
		return nil
	}

	var b MessageInteractionMetadataMessageComponent
	if err := json.Unmarshal(data, &b); err == nil && b.ID != "" {
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
	if err := json.Unmarshal(data, &a); err == nil && a.ID != "" {
		d.Value = &a
		return nil
	}

	var b MessageInteractionMetadataMessageComponent
	if err := json.Unmarshal(data, &b); err == nil && b.ID != "" {
		d.Value = &b
		return nil
	}

	var c MessageInteractionMetadataModalSubmit
	if err := json.Unmarshal(data, &c); err == nil && c.ID != "" {
		d.Value = &c
		return nil
	}

	return fmt.Errorf("unknown MessageInteractionMetadata: %s", string(data))
}

type MessageInteractionMetadataModalSubmit struct {
	ID                            Snowflake                                             `json:"id"`
	Type                          InteractionType                                       `json:"type"`
	User                          User                                                  `json:"user,omitempty"`
	AuthorizingIntegrationOwners  map[InteractionApplicationIntegrationType]interface{} `json:"authorizing_integration_owners,omitempty"`
	OriginalResponseMessageID     *Snowflake                                            `json:"original_response_message_id,omitempty"`
	TriggeringInteractionMetadata *MessageInteractionMetadataModalSubmitTriggering      `json:"triggering_interaction_metadata,omitempty"`
}
