package components

import (
	"go-discord-wrapper/types/common"
)

type AnyContainerComponent interface {
	MarshalJSON() ([]byte, error)
	GetType() common.ComponentType
	IsAnyContainerComponent() bool
}

type AnyChildComponent interface {
	IsAnyLabelComponent() bool
}

type AnyComponentInteractionResponse interface {
	IsInteractionResponseDataComponent() bool
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

type AnySectionComponent interface {
	IsAnySectionComponent() bool
}

type AnySectionAccessory interface {
	IsAnySectionAccessory() bool
}

type SelectDefaultValue struct {
	ID   common.Snowflake       `json:"id"`
	Type SelectDefaultValueType `json:"type"`
}

type SelectDefaultValueType string

const (
	SelectDefaultValueTypeUser    SelectDefaultValueType = "user"
	SelectDefaultValueTypeRole    SelectDefaultValueType = "role"
	SelectDefaultValueTypeChannel SelectDefaultValueType = "channel"
)

type ApplicationCommandInteractionOptionType int

const (
	ApplicationCommandInteractionOptionTypeSubCommand      ApplicationCommandInteractionOptionType = 1
	ApplicationCommandInteractionOptionTypeSubCommandGroup ApplicationCommandInteractionOptionType = 2
	ApplicationCommandInteractionOptionTypeString          ApplicationCommandInteractionOptionType = 3
	ApplicationCommandInteractionOptionTypeInteger         ApplicationCommandInteractionOptionType = 4
	ApplicationCommandInteractionOptionTypeBoolean         ApplicationCommandInteractionOptionType = 5
	ApplicationCommandInteractionOptionTypeUser            ApplicationCommandInteractionOptionType = 6
	ApplicationCommandInteractionOptionTypeChannel         ApplicationCommandInteractionOptionType = 7
	ApplicationCommandInteractionOptionTypeRole            ApplicationCommandInteractionOptionType = 8
	ApplicationCommandInteractionOptionTypeMentionable     ApplicationCommandInteractionOptionType = 9
	ApplicationCommandInteractionOptionTypeNumber          ApplicationCommandInteractionOptionType = 10
	ApplicationCommandInteractionOptionTypeAttachment      ApplicationCommandInteractionOptionType = 11
)

type UnfurledMediaItem struct {
	URL          string            `json:"url"`
	ProxyURL     string            `json:"proxy_url,omitempty"`
	Height       int               `json:"height,omitempty"`
	Width        int               `json:"width,omitempty"`
	ContentType  string            `json:"content_type,omitempty"`
	AttachmentID *common.Snowflake `json:"attachment_id,omitempty"`
}
