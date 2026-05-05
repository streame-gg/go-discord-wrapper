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

type GuildIntegrationsUpdateEvent struct {
	GuildID common.Snowflake `json:"guild_id"`
}

func init() {
	RegisterEvent(GuildEmojisUpdateEvent{})
	RegisterEvent(GuildStickersUpdateEvent{})
	RegisterEvent(GuildIntegrationsUpdateEvent{})
}

func (e GuildEmojisUpdateEvent) DesiredEventType() Event { return &GuildEmojisUpdateEvent{} }
func (e GuildEmojisUpdateEvent) Event() EventType        { return EventGuildEmojisUpdate }

func (e GuildStickersUpdateEvent) DesiredEventType() Event { return &GuildStickersUpdateEvent{} }
func (e GuildStickersUpdateEvent) Event() EventType        { return EventGuildStickersUpdate }

func (e GuildIntegrationsUpdateEvent) DesiredEventType() Event {
	return &GuildIntegrationsUpdateEvent{}
}
func (e GuildIntegrationsUpdateEvent) Event() EventType { return EventGuildIntegrationsUpdate }
