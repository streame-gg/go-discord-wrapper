package types

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type DiscordInteractionType int

const (
	DiscordInteractionTypePing                           DiscordInteractionType = 1
	DiscordInteractionTypeApplicationCommand             DiscordInteractionType = 2
	DiscordInteractionTypeMessageComponent               DiscordInteractionType = 3
	DiscordInteractionTypeApplicationCommandAutocomplete DiscordInteractionType = 4
	DiscordInteractionTypeModalSubmit                    DiscordInteractionType = 5
)

type DiscordInteractionApplicationIntegrationType int

const (
	DiscordInteractionApplicationIntegrationTypeGuildInstall DiscordInteractionApplicationIntegrationType = 0
	DiscordInteractionApplicationIntegrationTypeUserInstall  DiscordInteractionApplicationIntegrationType = 1
)

type DiscordMessageMessageReference struct {
	Type            *DiscordMessageMessageReferenceType `json:"type,omitempty"`
	MessageID       *DiscordSnowflake                   `json:"message_id,omitempty"`
	ChannelID       *DiscordSnowflake                   `json:"channel_id,omitempty"`
	GuildID         *DiscordSnowflake                   `json:"guild_id,omitempty"`
	FailIfNotExists *bool                               `json:"fail_if_not_exists,omitempty"`
}

type DiscordMessageInteractionMetadataApplicationCommand struct {
	ID                           DiscordSnowflake                                             `json:"id"`
	Type                         DiscordInteractionType                                       `json:"type"`
	User                         DiscordUser                                                  `json:"user,omitempty"`
	AuthorizingIntegrationOwners map[DiscordInteractionApplicationIntegrationType]interface{} `json:"authorizing_integration_owners,omitempty"`
	OriginalResponseMessageID    *DiscordSnowflake                                            `json:"original_response_message_id,omitempty"`
	TargetUser                   *DiscordUser                                                 `json:"target_user,omitempty"`
	TargetMessageID              *DiscordSnowflake                                            `json:"target_message_id,omitempty"`
}

type DiscordMessageInteractionMetadataMessageComponent struct {
	ID                           DiscordSnowflake                                             `json:"id"`
	Type                         DiscordInteractionType                                       `json:"type"`
	User                         DiscordUser                                                  `json:"user,omitempty"`
	AuthorizingIntegrationOwners map[DiscordInteractionApplicationIntegrationType]interface{} `json:"authorizing_integration_owners,omitempty"`
	OriginalResponseMessageID    *DiscordSnowflake                                            `json:"original_response_message_id,omitempty"`
	InteractedMessageID          *DiscordSnowflake                                            `json:"interacted_message_id,omitempty"`
}

type AnyDiscordMessageInteractionMetadata interface{}

type DiscordMessageInteractionMetadata struct {
	Value AnyDiscordMessageInteractionMetadata
}

type AnyDiscordMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata interface{}

type DiscordMessageInteractionMetadataModalSubmitTriggering struct {
	AnyDiscordMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata
}

func (d *DiscordMessageInteractionMetadataModalSubmitTriggering) UnmarshalJSON(data []byte) error {
	var a DiscordMessageInteractionMetadataApplicationCommand
	if err := json.Unmarshal(data, &a); err == nil && a.ID != "" {
		d.AnyDiscordMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata = &a
		return nil
	}

	var b DiscordMessageInteractionMetadataMessageComponent
	if err := json.Unmarshal(data, &b); err == nil && b.ID != "" {
		d.AnyDiscordMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata = &b
		return nil
	}

	return nil
}

func (d *DiscordMessageInteractionMetadata) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		d.Value = nil
		return nil
	}

	var a DiscordMessageInteractionMetadataApplicationCommand
	if err := json.Unmarshal(data, &a); err == nil && a.ID != "" {
		d.Value = &a
		return nil
	}

	var b DiscordMessageInteractionMetadataMessageComponent
	if err := json.Unmarshal(data, &b); err == nil && b.ID != "" {
		d.Value = &b
		return nil
	}

	var c DiscordMessageInteractionMetadataModalSubmit
	if err := json.Unmarshal(data, &c); err == nil && c.ID != "" {
		d.Value = &c
		return nil
	}

	return fmt.Errorf("unknown DiscordMessageInteractionMetadata: %s", string(data))
}

type DiscordMessageInteractionMetadataModalSubmit struct {
	ID                            DiscordSnowflake                                             `json:"id"`
	Type                          DiscordInteractionType                                       `json:"type"`
	User                          DiscordUser                                                  `json:"user,omitempty"`
	AuthorizingIntegrationOwners  map[DiscordInteractionApplicationIntegrationType]interface{} `json:"authorizing_integration_owners,omitempty"`
	OriginalResponseMessageID     *DiscordSnowflake                                            `json:"original_response_message_id,omitempty"`
	TriggeringInteractionMetadata *DiscordMessageInteractionMetadataModalSubmitTriggering      `json:"triggering_interaction_metadata,omitempty"`
}
