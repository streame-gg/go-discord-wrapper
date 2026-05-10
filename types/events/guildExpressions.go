package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type GuildEmojisUpdateEvent struct {
	GuildID common.Snowflake `json:"guild_id"`
	Emojis  []common.Emoji   `json:"emojis"`
}

type GuildStickersUpdateEvent struct {
	GuildID  common.Snowflake `json:"guild_id"`
	Stickers []common.Sticker `json:"stickers"`
}

func init() {
	RegisterEvent(GuildEmojisUpdateEvent{})
	RegisterEvent(GuildStickersUpdateEvent{})
}

func (e GuildEmojisUpdateEvent) DesiredEventType() Event { return &GuildEmojisUpdateEvent{} }
func (e GuildEmojisUpdateEvent) Event() EventType        { return EventGuildEmojisUpdate }

func (e GuildStickersUpdateEvent) DesiredEventType() Event { return &GuildStickersUpdateEvent{} }
func (e GuildStickersUpdateEvent) Event() EventType        { return EventGuildStickersUpdate }
