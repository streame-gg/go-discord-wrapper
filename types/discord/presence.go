package discord

import (
	"encoding/json"
)

// https://docs.discord.com/developers/events/gateway-events#activity-object-activity-types
type ActivityType int

const (
	ActivityTypePlaying   ActivityType = 0 // -> Playing {name}
	ActivityTypeStreaming ActivityType = 1 // -> Streaming {details}
	ActivityTypeListening ActivityType = 2 // -> Listening to {name}
	ActivityTypeWatching  ActivityType = 3 // -> Watching {name}
	ActivityTypeCustom    ActivityType = 4 // -> {emoji} {state}
	ActivityTypeCompeting ActivityType = 5 // -> Competing in {name}
)

// ActivityTimestamps holds the start and end times for an activity.
// Discord sends these as Unix milliseconds; they are exposed as time.Time.
//
// https://docs.discord.com/developers/events/gateway-events#activity-object-activity-timestamps
type ActivityTimestamps struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

// https://docs.discord.com/developers/events/gateway-events#activity-object-activity-party
type ActivityParty struct {
	ID   *string `json:"id,omitempty"`
	Size *[2]int `json:"size,omitempty"`
}

// https://docs.discord.com/developers/events/gateway-events#activity-object-activity-assets
type ActivityAssets struct {
	LargeImage       *string `json:"large_image,omitempty"`
	LargeText        *string `json:"large_text,omitempty"`
	LargeURL         *string `json:"large_url,omitempty"`
	SmallImage       *string `json:"small_image,omitempty"`
	SmallText        *string `json:"small_text,omitempty"`
	SmallURL         *string `json:"small_url,omitempty"`
	InviteCoverImage *string `json:"invite_cover_image,omitempty"`
}

// https://docs.discord.com/developers/events/gateway-events#activity-object-activity-secrets
type ActivitySecrets struct {
	Join     *string `json:"join,omitempty"`
	Spectate *string `json:"spectate,omitempty"`
	Match    *string `json:"match,omitempty"`
}

// https://docs.discord.com/developers/events/gateway-events#activity-object-activity-buttons
type ActivityButton struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// https://docs.discord.com/developers/events/gateway-events#activity-object-activity-emoji
type ActivityEmoji struct {
	Name     string     `json:"name"`
	ID       *Snowflake `json:"id,omitempty"`
	Animated *bool      `json:"animated,omitempty"`
}

// https://docs.discord.com/developers/events/gateway-events#activity-object-activity-flags
type ActivityFlags int

const (
	ActivityFlagInstance                 ActivityFlags = 1 << 0
	ActivityFlagJoin                     ActivityFlags = 1 << 1
	ActivityFlagSpectate                 ActivityFlags = 1 << 2
	ActivityFlagJoinRequest              ActivityFlags = 1 << 3
	ActivityFlagSync                     ActivityFlags = 1 << 4
	ActivityFlagPlay                     ActivityFlags = 1 << 5
	ActivityFlagPartyPrivacyFriends      ActivityFlags = 1 << 6
	ActivityFlagPartyPrivacyVoiceChannel ActivityFlags = 1 << 7
	ActivityFlagEmbedded                 ActivityFlags = 1 << 8
)

// FullActivity represents a rich presence activity.
// CreatedAt is the activity creation time; Discord sends it as Unix milliseconds.
//
// https://docs.discord.com/developers/events/gateway-events#activity-object
type FullActivity struct {
	Name              string               `json:"name"`
	Type              ActivityType         `json:"type"`
	URL               *string              `json:"url,omitempty"`
	CreatedAt         int64                `json:"created_at"`
	Timestamps        *ActivityTimestamps  `json:"timestamps,omitempty"`
	ApplicationID     *Snowflake           `json:"application_id,omitempty"`
	StatusDisplayType *StatusDisplayType   `json:"status_display_type,omitempty"`
	Details           *string              `json:"details,omitempty"`
	State             *string              `json:"state,omitempty"`
	Emoji             *ActivityEmoji       `json:"emoji,omitempty"`
	Party             *ActivityParty       `json:"party,omitempty"`
	Assets            *ActivityAssets      `json:"assets,omitempty"`
	Secrets           *ActivitySecrets     `json:"secrets,omitempty"`
	Instance          *bool                `json:"instance,omitempty"`
	Flags             *ActivityFlags       `json:"flags,omitempty"`
	Buttons           []FullActivityButton `json:"buttons,omitempty"`
}

type FullActivityButton struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type StatusDisplayType uint8

const (
	StatusDisplayTypeName    StatusDisplayType = 0
	StatusDisplayTypeState   StatusDisplayType = 1
	StatusDisplayTypeDetails StatusDisplayType = 2
)

func (f *FullActivity) UnmarshalJSON(data []byte) error {
	type Alias FullActivity
	var raw struct {
		*Alias
		CreatedAt int64 `json:"created_at"`
	}
	raw.Alias = (*Alias)(f)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	f.CreatedAt = raw.CreatedAt
	return nil
}

func (f FullActivity) MarshalJSON() ([]byte, error) {
	type Alias FullActivity
	return json.Marshal(struct {
		Alias
		CreatedAt int64 `json:"created_at"`
	}{
		Alias:     Alias(f),
		CreatedAt: f.CreatedAt,
	})
}

// https://docs.discord.com/developers/events/gateway-events#client-status-object
type ClientStatus struct {
	Desktop PresenceStatus `json:"desktop,omitempty"`
	Mobile  PresenceStatus `json:"mobile,omitempty"`
	Web     PresenceStatus `json:"web,omitempty"`
}

// https://docs.discord.com/developers/events/gateway-events#update-presence-status-types
type PresenceStatus string

const (
	PresenceStatusOnline    PresenceStatus = "online"
	PresenceStatusDND       PresenceStatus = "dnd"
	PresenceStatusIdle      PresenceStatus = "idle"
	PresenceStatusOffline   PresenceStatus = "offline"
	PresenceStatusInvisible PresenceStatus = "invisible"
)

// https://docs.discord.com/developers/events/gateway-events#presence-update-presence-update-event-fields
type Presence struct {
	User         User           `json:"user"`
	GuildID      Snowflake      `json:"guild_id"`
	Status       PresenceStatus `json:"status"`
	Activities   []FullActivity `json:"activities"`
	ClientStatus ClientStatus   `json:"client_status"`
}
