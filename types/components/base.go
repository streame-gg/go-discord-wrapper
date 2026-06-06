package components

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#what-is-a-component
type AnyContainerComponent interface {
	MarshalJSON() ([]byte, error)
	GetType() discord.ComponentType
	IsAnyContainerComponent()
}

// https://docs.discord.com/developers/components/reference#what-is-a-component
type AnyChildComponent interface {
	IsAnyLabelComponent()
}

// https://docs.discord.com/developers/components/reference#what-is-a-component
type AnyComponentInteractionResponse interface {
	IsInteractionResponseDataComponent()
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

// https://docs.discord.com/developers/components/reference#section
type AnySectionComponent interface {
	IsAnySectionComponent()
}

// https://docs.discord.com/developers/components/reference#section
type AnySectionAccessory interface {
	IsAnySectionAccessory()
}

// https://docs.discord.com/developers/components/reference#user-select-select-default-value-structure
type SelectDefaultValue struct {
	ID   discord.Snowflake      `json:"id"`
	Type SelectDefaultValueType `json:"type"`
}

// https://docs.discord.com/developers/components/reference#user-select-select-default-value-structure
type SelectDefaultValueType string

const (
	SelectDefaultValueTypeUser    SelectDefaultValueType = "user"
	SelectDefaultValueTypeRole    SelectDefaultValueType = "role"
	SelectDefaultValueTypeChannel SelectDefaultValueType = "channel"
)

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-object-application-command-interaction-data-option-structure
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

// https://docs.discord.com/developers/components/reference#unfurled-media-item-structure
type UnfurledMediaItem struct {
	URL          string             `json:"url"`
	ProxyURL     string             `json:"proxy_url,omitempty"`
	Height       int                `json:"height,omitempty"`
	Width        int                `json:"width,omitempty"`
	ContentType  string             `json:"content_type,omitempty"`
	AttachmentID *discord.Snowflake `json:"attachment_id,omitempty"`
}
