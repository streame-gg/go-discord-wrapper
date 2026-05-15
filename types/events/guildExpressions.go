package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type GuildEmojisUpdateEvent struct {
	GuildID discord.Snowflake `json:"guild_id"`
	Emojis  []discord.Emoji   `json:"emojis"`
}

type GuildStickersUpdateEvent struct {
	GuildID  discord.Snowflake `json:"guild_id"`
	Stickers []discord.Sticker `json:"stickers"`
}

func init() {
	RegisterEvent(GuildEmojisUpdateEvent{})
	RegisterEvent(GuildStickersUpdateEvent{})
}

func (e GuildEmojisUpdateEvent) DesiredEventType() Event { return &GuildEmojisUpdateEvent{} }
func (e GuildEmojisUpdateEvent) Event() EventType        { return EventGuildEmojisUpdate }

func (e GuildStickersUpdateEvent) DesiredEventType() Event { return &GuildStickersUpdateEvent{} }
func (e GuildStickersUpdateEvent) Event() EventType        { return EventGuildStickersUpdate }
