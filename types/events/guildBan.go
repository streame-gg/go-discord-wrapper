package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/events/gateway-events#guild-ban-add
type GuildBanAddEvent struct {
	GuildID discord.Snowflake `json:"guild_id"`
	User    discord.User      `json:"user"`
}

// https://docs.discord.com/developers/events/gateway-events#guild-ban-remove
type GuildBanRemoveEvent struct {
	GuildID discord.Snowflake `json:"guild_id"`
	User    discord.User      `json:"user"`
}

func init() {
	RegisterEvent(GuildBanAddEvent{})
	RegisterEvent(GuildBanRemoveEvent{})
}

func (e GuildBanAddEvent) DesiredEventType() Event { return &GuildBanAddEvent{} }
func (e GuildBanAddEvent) Event() EventType        { return EventGuildBanAdd }

func (e GuildBanRemoveEvent) DesiredEventType() Event { return &GuildBanRemoveEvent{} }
func (e GuildBanRemoveEvent) Event() EventType        { return EventGuildBanRemove }
