package discord

import "time"

// https://docs.discord.com/developers/resources/message#embed-object-embed-footer-structure
type EmbedFooter struct {
	Text         string `json:"text"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// https://docs.discord.com/developers/resources/message#embed-object-embed-image-structure
type EmbedImage struct {
	URL                string          `json:"url"`
	ProxyURL           string          `json:"proxy_url,omitempty"`
	Height             *int            `json:"height,omitempty"`
	Width              *int            `json:"width,omitempty"`
	ContentType        string          `json:"content_type,omitempty"`
	Placeholder        string          `json:"placeholder,omitempty"`
	PlaceholderVersion *int            `json:"placeholder_version,omitempty"`
	Description        string          `json:"description,omitempty"`
	Flags              EmbedMediaFlags `json:"flags,omitempty"`
}

// https://docs.discord.com/developers/resources/message#embed-object-embed-media-flags
type EmbedMediaFlags uint64

const (
	EmbedMediaFlagIsAnimated EmbedMediaFlags = 1 << 5
)

// https://docs.discord.com/developers/resources/message#embed-object-embed-image-structure
type EmbedThumbnail struct {
	URL                string          `json:"url"`
	ProxyURL           string          `json:"proxy_url,omitempty"`
	Height             *int            `json:"height,omitempty"`
	Width              *int            `json:"width,omitempty"`
	ContentType        string          `json:"content_type,omitempty"`
	Placeholder        string          `json:"placeholder,omitempty"`
	PlaceholderVersion *int            `json:"placeholder_version,omitempty"`
	Description        string          `json:"description,omitempty"`
	Flags              EmbedMediaFlags `json:"flags,omitempty"`
}

// https://docs.discord.com/developers/resources/message#embed-object-embed-video-structure
type EmbedVideo struct {
	URL                string          `json:"url"`
	ProxyURL           string          `json:"proxy_url,omitempty"`
	Height             *int            `json:"height,omitempty"`
	Width              *int            `json:"width,omitempty"`
	ContentType        string          `json:"content_type,omitempty"`
	Placeholder        string          `json:"placeholder,omitempty"`
	PlaceholderVersion *int            `json:"placeholder_version,omitempty"`
	Description        string          `json:"description,omitempty"`
	Flags              EmbedMediaFlags `json:"flags,omitempty"`
}

// https://docs.discord.com/developers/resources/message#embed-object-embed-provider-structure
type EmbedProvider struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// https://docs.discord.com/developers/resources/message#embed-object-embed-author-structure
type EmbedAuthor struct {
	Name         string `json:"name"`
	URL          string `json:"url,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// https://docs.discord.com/developers/resources/message#embed-object-embed-field-structure
type EmbedFields struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// https://docs.discord.com/developers/resources/message#embed-object-embed-types
type EmbedType string

const (
	EmbedTypeRich       EmbedType = "rich"
	EmbedTypeImage      EmbedType = "image"
	EmbedTypeVideo      EmbedType = "video"
	EmbedTypeGIFV       EmbedType = "gifv"
	EmbedTypeArticle    EmbedType = "article"
	EmbedTypeLink       EmbedType = "link"
	EmbedTypePollResult EmbedType = "poll_result"
)

// https://docs.discord.com/developers/resources/message#embed-object
type Embed struct {
	Title       string          `json:"title,omitempty"`
	Type        EmbedType       `json:"type,omitempty"`
	Description string          `json:"description,omitempty"`
	URL         string          `json:"url,omitempty"`
	Timestamp   *time.Time      `json:"timestamp,omitempty"`
	Color       *int            `json:"color,omitempty"`
	Footer      *EmbedFooter    `json:"footer,omitempty"`
	Image       *EmbedImage     `json:"image,omitempty"`
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *EmbedVideo     `json:"video,omitempty"`
	Provider    *EmbedProvider  `json:"provider,omitempty"`
	Author      *EmbedAuthor    `json:"author,omitempty"`
	Fields      []EmbedFields   `json:"fields,omitempty"`
	Flags       EmbedFlags      `json:"flags,omitempty"`
}

// https://docs.discord.com/developers/resources/message#embed-object-embed-flags
type EmbedFlags uint64

const (
	EmbedFlagIsContentInventoryEntry EmbedFlags = 1 << 5
)
