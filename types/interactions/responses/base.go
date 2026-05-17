package responses

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type AnyInteractionResponseData interface {
	IsInteractionResponseData() bool
	MarshalJSON() ([]byte, error)
}

type InteractionResponseDataDefault struct {
	TTS             bool                     `json:"tts,omitempty"`
	Content         string                   `json:"content,omitempty"`
	Embeds          *[]discord.Embed         `json:"embeds,omitempty"`
	AllowedMentions *discord.AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           discord.MessageFlag      `json:"flags,omitempty"`
	Components      *[]discord.AnyComponent  `json:"components,omitempty"`
	Attachments     *[]discord.Attachment    `json:"attachments,omitempty"`
	Poll            *discord.PollRequest     `json:"poll,omitempty"`
}

func (d *InteractionResponseDataDefault) IsInteractionResponseData() bool {
	return true
}

func (d *InteractionResponseDataDefault) MarshalJSON() ([]byte, error) {
	type Alias InteractionResponseDataDefault
	return json.Marshal((*Alias)(d))
}

type InteractionResponse struct {
	Type discord.InteractionCallbackType `json:"type"`
	Data AnyInteractionResponseData      `json:"data,omitempty"`
}

type InteractionCallback struct {
	ID                       discord.Snowflake       `json:"id"`
	Type                     discord.InteractionType `json:"type"`
	ActivityInstanceID       *discord.Snowflake      `json:"activity_instance_id,omitempty"`
	ResponseMessageID        *discord.Snowflake      `json:"response_message_id,omitempty"`
	ResponseMessageLoading   *bool                   `json:"response_message_loading,omitempty"`
	ResponseMessageEphemeral *bool                   `json:"response_message_ephemeral,omitempty"`
}

type InteractionCallbackActivityInstance struct {
	ID string `json:"id"`
}

type InteractionCallbackResource struct {
	Type             discord.InteractionCallbackType      `json:"type"`
	ActivityInstance *InteractionCallbackActivityInstance `json:"activity_instance,omitempty"`
	Message          *discord.Message                     `json:"message,omitempty"`
}

type InteractionCallbackResponse struct {
	Interaction InteractionCallback          `json:"interaction"`
	Resource    *InteractionCallbackResource `json:"resource,omitempty"`
}
