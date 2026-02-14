package common

import "time"

type EmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type EmbedImage struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   *int   `json:"height,omitempty"`
	Width    *int   `json:"width,omitempty"`
}

type EmbedThumbnail struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   *int   `json:"height,omitempty"`
	Width    *int   `json:"width,omitempty"`
}

type EmbedVideo struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   *int   `json:"height,omitempty"`
	Width    *int   `json:"width,omitempty"`
}

type EmbedProvider struct {
	Name *string `json:"name,omitempty"`
	URL  *string `json:"url,omitempty"`
}

type EmbedAuthor struct {
	Name         string  `json:"name"`
	URL          *string `json:"url"`
	IconURL      *string `json:"icon_url,omitempty"`
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"`
}

type EmbedFields struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline *bool  `json:"inline,omitempty"`
}

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

type Embed struct {
	Title       *string         `json:"title,omitempty"`
	Type        *EmbedType      `json:"type,omitempty"`
	Description *string         `json:"description,omitempty"`
	URL         *string         `json:"url,omitempty"`
	Timestamp   *time.Time      `json:"timestamp,omitempty"`
	Color       *int            `json:"color,omitempty"`
	Footer      *EmbedFooter    `json:"footer,omitempty"`
	Image       *EmbedImage     `json:"image,omitempty"`
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *EmbedVideo     `json:"video,omitempty"`
	Provider    *EmbedProvider  `json:"provider,omitempty"`
	Author      *EmbedAuthor    `json:"author,omitempty"`
	Fields      []EmbedFields   `json:"fields,omitempty"`
}
