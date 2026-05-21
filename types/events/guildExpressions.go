package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type GuildEmojisUpdateEvent struct {
	GuildID   discord.Snowflake `json:"-"`
	NewEmojis []*discord.Emoji  `json:"-"`
	OldEmojis []*discord.Emoji  `json:"-"`
}

func (e *GuildEmojisUpdateEvent) UnmarshalJSON(data []byte) error {
	type wire struct {
		GuildID discord.Snowflake `json:"guild_id"`
		Emojis  []*discord.Emoji  `json:"emojis"`
	}
	var w wire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	e.GuildID = w.GuildID
	e.NewEmojis = w.Emojis
	return nil
}

func (e GuildEmojisUpdateEvent) MarshalJSON() ([]byte, error) {
	type wire struct {
		GuildID   discord.Snowflake `json:"guild_id"`
		Emojis    []*discord.Emoji  `json:"emojis"`
		OldEmojis []*discord.Emoji  `json:"old_emojis,omitempty"`
	}
	return json.Marshal(wire{e.GuildID, e.NewEmojis, e.OldEmojis})
}

type GuildStickersUpdateEvent struct {
	GuildID     discord.Snowflake  `json:"-"`
	NewStickers []*discord.Sticker `json:"-"`
	OldStickers []*discord.Sticker `json:"-"`
}

func (e *GuildStickersUpdateEvent) UnmarshalJSON(data []byte) error {
	type wire struct {
		GuildID  discord.Snowflake  `json:"guild_id"`
		Stickers []*discord.Sticker `json:"stickers"`
	}
	var w wire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	e.GuildID = w.GuildID
	e.NewStickers = w.Stickers
	return nil
}

func (e GuildStickersUpdateEvent) MarshalJSON() ([]byte, error) {
	type wire struct {
		GuildID     discord.Snowflake  `json:"guild_id"`
		Stickers    []*discord.Sticker `json:"stickers"`
		OldStickers []*discord.Sticker `json:"old_stickers,omitempty"`
	}
	return json.Marshal(wire{e.GuildID, e.NewStickers, e.OldStickers})
}

func init() {
	RegisterEvent(GuildEmojisUpdateEvent{})
	RegisterEvent(GuildStickersUpdateEvent{})
}

func (e GuildEmojisUpdateEvent) DesiredEventType() Event { return &GuildEmojisUpdateEvent{} }
func (e GuildEmojisUpdateEvent) Event() EventType        { return EventGuildEmojisUpdate }

func (e GuildStickersUpdateEvent) DesiredEventType() Event { return &GuildStickersUpdateEvent{} }
func (e GuildStickersUpdateEvent) Event() EventType        { return EventGuildStickersUpdate }
