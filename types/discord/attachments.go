package discord

import "time"

// https://docs.discord.com/developers/resources/message#attachment-object
type Attachment struct {
	ID                 Snowflake      `json:"id"`
	Filename           string         `json:"filename"`
	Title              string         `json:"title,omitempty"`
	Description        string         `json:"description,omitempty"`
	ContentType        string         `json:"content_type,omitempty"`
	Size               int            `json:"size"`
	URL                string         `json:"url"`
	ProxyURL           string         `json:"proxy_url"`
	Height             *int           `json:"height,omitempty"`
	Width              *int           `json:"width,omitempty"`
	Placeholder        string         `json:"placeholder,omitempty"`
	PlaceholderVersion int            `json:"placeholder_version,omitempty"`
	Ephemeral          bool           `json:"ephemeral,omitempty"`
	DurationSecs       float64        `json:"duration_secs,omitempty"`
	Waveform           string         `json:"waveform,omitempty"`
	Flags              AttachmentFlag `json:"flags,omitempty"`
	ClipParticipants   []User         `json:"clip_participants,omitempty"`
	ClipCreatedAt      time.Time      `json:"clip_created_at,omitempty"`
	Application        *Application   `json:"application,omitempty"`
}

// https://docs.discord.com/developers/resources/message#attachment-object-attachment-flags
type AttachmentFlag uint64

const (
	AttachmentFlagIsClip      AttachmentFlag = 1 << 0
	AttachmentFlagIsThumbnail AttachmentFlag = 1 << 1
	AttachmentFlagIsRemix     AttachmentFlag = 1 << 2
	AttachmentFlagIsSpoiler   AttachmentFlag = 1 << 3
	AttachmentFlagIsAnimated  AttachmentFlag = 1 << 5
)
