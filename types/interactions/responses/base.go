package responses

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type AnyInteractionResponseData interface {
	IsInteractionResponseData() bool
	MarshalJSON() ([]byte, error)
}

type InteractionResponseDataDefault struct {
	TTS             bool                    `json:"tts,omitempty"`
	Content         string                  `json:"content,omitempty"`
	Embeds          *[]common.Embed         `json:"embeds,omitempty"`
	AllowedMentions *common.AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           common.MessageFlag      `json:"flags,omitempty"`
	Components      *[]common.AnyComponent  `json:"components,omitempty"`
	//TODO partial
	Attachment   *[]common.Attachment `json:"attachment,omitempty"`
	Poll         *common.PollRequest  `json:"poll,omitempty"`
	WithResponse bool                 `json:"with_response,omitempty"`
}

func (d *InteractionResponseDataDefault) IsInteractionResponseData() bool {
	return true
}

func (d *InteractionResponseDataDefault) MarshalJSON() ([]byte, error) {
	type Alias InteractionResponseDataDefault
	return json.Marshal((*Alias)(d))
}

type InteractionResponse struct {
	Type common.InteractionCallbackType `json:"type"`
	Data AnyInteractionResponseData     `json:"data,omitempty"`
}

type InteractionCallback struct {
	ID                       common.Snowflake       `json:"id"`
	Type                     common.InteractionType `json:"type"`
	ActivityInstanceID       *common.Snowflake      `json:"activity_instance_id,omitempty"`
	ResponseMessageID        *common.Snowflake      `json:"response_message_id,omitempty"`
	ResponseMessageLoading   *bool                  `json:"response_message_loading,omitempty"`
	ResponseMessageEphemeral *bool                  `json:"response_message_ephemeral,omitempty"`
}

type InteractionCallbackActivityInstance struct {
	ID string `json:"id"`
}

type InteractionCallbackResource struct {
	Type             common.InteractionCallbackType       `json:"type"`
	ActivityInstance *InteractionCallbackActivityInstance `json:"activity_instance,omitempty"`
	Message          *common.Message                      `json:"message,omitempty"`
}

type InteractionCallbackResponse struct {
	Interaction InteractionCallback          `json:"interaction"`
	Resource    *InteractionCallbackResource `json:"resource"`
}
