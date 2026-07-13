package components

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/components/reference#what-is-a-component
type AnyContainerComponent interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
	GetType() discord.ComponentType
	IsAnyContainerComponent()
}

// https://docs.discord.com/developers/components/reference#what-is-a-component
type AnyChildComponent interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
	GetType() discord.ComponentType
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
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
	IsAnySectionComponent()
	GetType() discord.ComponentType
}

// https://docs.discord.com/developers/components/reference#section
type AnySectionAccessory interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
	IsAnySectionAccessory()
	GetType() discord.ComponentType
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

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-type
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
	URL                string                 `json:"url"`
	ProxyURL           string                 `json:"proxy_url,omitempty"`
	Height             *int                   `json:"height,omitempty"`
	Width              *int                   `json:"width,omitempty"`
	ContentType        string                 `json:"content_type,omitempty"`
	AttachmentID       *discord.Snowflake     `json:"attachment_id,omitempty"`
	PlaceholderVersion int                    `json:"placeholder_version,omitempty"`
	Placeholder        string                 `json:"placeholder,omitempty"`
	Flags              UnfurledMediaItemFlags `json:"flags,omitempty"`
}

// https://docs.discord.com/developers/components/reference#unfurled-media-item-unfurled-media-item-flags
type UnfurledMediaItemFlags uint64

const (
	UnfurledMediaItemFlagIsAnimated UnfurledMediaItemFlags = 1 << 0
)
