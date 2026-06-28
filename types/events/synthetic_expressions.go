package events

import "github.com/streame-gg/go-discord-wrapper/types/discord"

// GuildEmojiAddEvent fires when an emoji is added to a guild.
// Derived from GUILD_EMOJIS_UPDATE by diffing the old and new emoji sets.
// Wrapper-synthesized; derived from https://docs.discord.com/developers/events/gateway-events#guild-emojis-update
type GuildEmojiAddEvent struct {
	GuildID discord.Snowflake
	Emoji   *discord.Emoji
}

// GuildEmojiRemoveEvent fires when an emoji is removed from a guild.
// Wrapper-synthesized; derived from https://docs.discord.com/developers/events/gateway-events#guild-emojis-update
type GuildEmojiRemoveEvent struct {
	GuildID discord.Snowflake
	Emoji   *discord.Emoji
}

// GuildEmojiUpdateEvent fires when an emoji's name changes.
// Wrapper-synthesized; derived from https://docs.discord.com/developers/events/gateway-events#guild-emojis-update
type GuildEmojiUpdateEvent struct {
	GuildID  discord.Snowflake
	OldEmoji *discord.Emoji
	NewEmoji *discord.Emoji
}

// GuildStickerAddEvent fires when a sticker is added to a guild.
// Derived from GUILD_STICKERS_UPDATE by diffing the old and new sticker sets.
// Wrapper-synthesized; derived from https://docs.discord.com/developers/events/gateway-events#guild-stickers-update
type GuildStickerAddEvent struct {
	GuildID discord.Snowflake
	Sticker *discord.Sticker
}

// GuildStickerRemoveEvent fires when a sticker is removed from a guild.
// Wrapper-synthesized; derived from https://docs.discord.com/developers/events/gateway-events#guild-stickers-update
type GuildStickerRemoveEvent struct {
	GuildID discord.Snowflake
	Sticker *discord.Sticker
}

// GuildStickerUpdateEvent fires when a sticker's name changes.
// Wrapper-synthesized; derived from https://docs.discord.com/developers/events/gateway-events#guild-stickers-update
type GuildStickerUpdateEvent struct {
	GuildID    discord.Snowflake
	OldSticker *discord.Sticker
	NewSticker *discord.Sticker
}

func (e *GuildEmojiAddEvent) Event() EventType        { return EventWrapperGuildEmojiAdd }
func (e *GuildEmojiAddEvent) DesiredEventType() Event { return &GuildEmojiAddEvent{} }

func (e *GuildEmojiUpdateEvent) Event() EventType        { return EventWrapperGuildEmojiUpdate }
func (e *GuildEmojiUpdateEvent) DesiredEventType() Event { return &GuildEmojiUpdateEvent{} }

func (e *GuildEmojiRemoveEvent) Event() EventType        { return EventWrapperGuildEmojiRemove }
func (e *GuildEmojiRemoveEvent) DesiredEventType() Event { return &GuildEmojiRemoveEvent{} }

func (e *GuildStickerAddEvent) Event() EventType        { return EventWrapperGuildStickerAdd }
func (e *GuildStickerAddEvent) DesiredEventType() Event { return &GuildStickerAddEvent{} }

func (e *GuildStickerUpdateEvent) Event() EventType        { return EventWrapperGuildStickerUpdate }
func (e *GuildStickerUpdateEvent) DesiredEventType() Event { return &GuildStickerUpdateEvent{} }

func (e *GuildStickerRemoveEvent) Event() EventType        { return EventWrapperGuildStickerRemove }
func (e *GuildStickerRemoveEvent) DesiredEventType() Event { return &GuildStickerRemoveEvent{} }
