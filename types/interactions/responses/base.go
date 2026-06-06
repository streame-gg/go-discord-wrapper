package responses

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-response-object
type AnyInteractionResponseData interface {
	IsInteractionResponseData() bool
	MarshalJSON() ([]byte, error)
}

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-response-object-messages
type InteractionResponseDataDefault struct {
	TTS             bool                     `json:"tts,omitempty"`
	Content         string                   `json:"content,omitempty"`
	Embeds          *[]discord.Embed         `json:"embeds,omitempty"`
	AllowedMentions *discord.AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           discord.MessageFlag      `json:"flags,omitempty"`
	Components      *[]discord.AnyComponent  `json:"components,omitempty"`
	Attachments     *[]discord.Attachment    `json:"attachments,omitempty"`
	Poll            *discord.PollRequest     `json:"poll,omitempty"`

	Files []discord.MessageFile `json:"-"`
}

func (d *InteractionResponseDataDefault) IsInteractionResponseData() bool {
	return true
}

func (d *InteractionResponseDataDefault) MarshalJSON() ([]byte, error) {
	type Alias InteractionResponseDataDefault
	return json.Marshal((*Alias)(d))
}

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-response-object
type InteractionResponse struct {
	Type discord.InteractionCallbackType `json:"type"`
	Data AnyInteractionResponseData      `json:"data,omitempty"`
}

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-callback-object-interaction-callback-structure
type InteractionCallback struct {
	ID                       discord.Snowflake       `json:"id"`
	Type                     discord.InteractionType `json:"type"`
	ActivityInstanceID       *discord.Snowflake      `json:"activity_instance_id,omitempty"`
	ResponseMessageID        *discord.Snowflake      `json:"response_message_id,omitempty"`
	ResponseMessageLoading   *bool                   `json:"response_message_loading,omitempty"`
	ResponseMessageEphemeral *bool                   `json:"response_message_ephemeral,omitempty"`
}

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-callback-object-interaction-callback-activity-instance-resource
type InteractionCallbackActivityInstance struct {
	ID string `json:"id"`
}

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-callback-object-interaction-callback-resource-object
type InteractionCallbackResource struct {
	Type             discord.InteractionCallbackType      `json:"type"`
	ActivityInstance *InteractionCallbackActivityInstance `json:"activity_instance,omitempty"`
	Message          *discord.Message                     `json:"message,omitempty"`
}

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-callback-object-interaction-callback-response-object
type InteractionCallbackResponse struct {
	Interaction InteractionCallback          `json:"interaction"`
	Resource    *InteractionCallbackResource `json:"resource,omitempty"`
}
