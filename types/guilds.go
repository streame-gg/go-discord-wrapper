package types

import "encoding/json"

type AnyGuildWrapper struct {
	Guild AnyGuild
}

func (ag *AnyGuildWrapper) UnmarshalJSON(data []byte) error {
	var probe struct {
		Unavailable *bool `json:"unavailable"`
	}

	if err := json.Unmarshal(data, &probe); err != nil {
		return err
	}

	if probe.Unavailable != nil && *probe.Unavailable {
		var ug UnavailableGuild
		if err := json.Unmarshal(data, &ug); err != nil {
			return err
		}
		ag.Guild = ug
		return nil
	}

	var g Guild
	if err := json.Unmarshal(data, &g); err != nil {
		return err
	}
	ag.Guild = g
	return nil
}

type AnyGuild interface {
	IsAvailable() bool
	GetID() DiscordSnowflake
}

type Guild struct {
	ID                          DiscordSnowflake                       `json:"id"`
	Name                        string                                 `json:"name"`
	IconHash                    *string                                `json:"icon,omitempty"`
	Splash                      *string                                `json:"splash,omitempty"`
	DiscoverySplash             *string                                `json:"discovery_splash,omitempty"`
	Owner                       bool                                   `json:"owner,omitempty"`
	OwnerID                     DiscordSnowflake                       `json:"owner_id,omitempty"`
	Permissions                 *string                                `json:"permissions,omitempty"`
	Region                      *string                                `json:"region,omitempty"`
	AfkChannelID                *DiscordSnowflake                      `json:"afk_channel_id,omitempty"`
	AfkTimeout                  *int                                   `json:"afk_timeout,omitempty"`
	WidgetEnabled               *bool                                  `json:"widget_enabled,omitempty"`
	WidgetChannelID             *DiscordSnowflake                      `json:"widget_channel_id,omitempty"`
	VerificationChannel         *DiscordSnowflake                      `json:"verification_channel_id,omitempty"`
	VerificationLevel           DiscordGuildVerificationLevel          `json:"verification_level,omitempty"`
	DefaultMessageNotifications DiscordDefaultMessageNotificationLevel `json:"default_message_notifications,omitempty"`
}

func (g Guild) IsAvailable() bool {
	return true
}

func (g Guild) GetID() DiscordSnowflake {
	return g.ID
}

type UnavailableGuild struct {
	ID          DiscordSnowflake `json:"id"`
	Unavailable *bool            `json:"unavailable"`
	Name        string           `json:"name,omitempty"`
}

func (ug UnavailableGuild) IsAvailable() bool {
	return !*ug.Unavailable
}

func (ug UnavailableGuild) GetID() DiscordSnowflake {
	return ug.ID
}

type DiscordGuildVerificationLevel int

const (
	DiscordGuildVerificationLevelNone                                  DiscordGuildVerificationLevel = 0
	DiscordGuildVerificationLevelLow                                   DiscordGuildVerificationLevel = 1
	DiscordGuildVerificationLevelMedium                                DiscordGuildVerificationLevel = 2
	DiscordGuildVerificationLevelHigh                                  DiscordGuildVerificationLevel = 3
	DiscordGuildVerificationLevelVeryHighDiscordGuildVerificationLevel                               = 4
)

type DiscordDefaultMessageNotificationLevel int

const (
	DiscordDefaultMessageNotificationLevelAllMessages  DiscordDefaultMessageNotificationLevel = 0
	DiscordDefaultMessageNotificationLevelOnlyMentions DiscordDefaultMessageNotificationLevel = 1
)
