package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &GuildEmojisUpdateEvent{}
var _ Event = &GuildStickersUpdateEvent{}

func init() {
	RegisterEvent(&GuildEmojisUpdateEvent{})
	RegisterEvent(&GuildStickersUpdateEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#guild-emojis-update
type GuildEmojisUpdateEvent struct {
	GuildID   discord.Snowflake `json:"-"`
	NewEmojis []discord.Emoji   `json:"-"`
	OldEmojis []discord.Emoji   `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#guild-stickers-update
type GuildStickersUpdateEvent struct {
	GuildID     discord.Snowflake `json:"-"`
	NewStickers []discord.Sticker `json:"-"`
	OldStickers []discord.Sticker `json:"-"`
}

func (e *GuildEmojisUpdateEvent) DesiredEventType() Event { return &GuildEmojisUpdateEvent{} }
func (e *GuildEmojisUpdateEvent) Event() EventType        { return EventGuildEmojisUpdate }
func (e *GuildEmojisUpdateEvent) UnmarshalJSON(data []byte) error {
	var outer struct {
		GuildID discord.Snowflake `json:"guild_id"`
		Emojis  []discord.Emoji   `json:"emojis"`
	}
	if err := json.Unmarshal(data, &outer); err != nil {
		return err
	}

	e.GuildID = outer.GuildID
	e.NewEmojis = outer.Emojis
	return nil
}

func (e *GuildStickersUpdateEvent) DesiredEventType() Event { return &GuildStickersUpdateEvent{} }
func (e *GuildStickersUpdateEvent) Event() EventType        { return EventGuildStickersUpdate }
func (e *GuildStickersUpdateEvent) UnmarshalJSON(data []byte) error {
	var outer struct {
		GuildID  discord.Snowflake `json:"guild_id"`
		Stickers []discord.Sticker `json:"stickers"`
	}
	if err := json.Unmarshal(data, &outer); err != nil {
		return err
	}
	e.GuildID = outer.GuildID
	e.NewStickers = outer.Stickers
	return nil
}
