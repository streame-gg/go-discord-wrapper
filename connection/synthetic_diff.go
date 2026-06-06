package connection

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// deriveSyntheticEvents returns the synthetic events to dispatch after ev.
// Called synchronously on the websocket reader goroutine after internalEventHandler
// has already applied cache updates and populated Old fields.
func (d *Client) deriveSyntheticEvents(ev events.Event) []events.Event {
	switch e := ev.(type) {
	case *events.GuildEmojisUpdateEvent:
		return deriveGuildEmojiSyntheticEvents(e)
	case *events.GuildStickersUpdateEvent:
		return deriveGuildStickerSyntheticEvents(e)
	}
	return nil
}

// deriveGuildEmojiSyntheticEvents diffs the old and new emoji sets from a
// GUILD_EMOJIS_UPDATE and fires Add/Remove/Update events per changed emoji.
// Returns nil when OldEmojis is nil (cache cold).
func deriveGuildEmojiSyntheticEvents(ev *events.GuildEmojisUpdateEvent) []events.Event {
	if ev.OldEmojis == nil {
		return nil
	}

	oldByID := make(map[discord.Snowflake]*discord.Emoji, len(ev.OldEmojis))
	for _, e := range ev.OldEmojis {
		oldByID[e.ID] = e
	}
	newByID := make(map[discord.Snowflake]*discord.Emoji, len(ev.NewEmojis))
	for _, e := range ev.NewEmojis {
		newByID[e.ID] = e
	}

	var result []events.Event

	for id, newEmoji := range newByID {
		if old, exists := oldByID[id]; !exists {
			result = append(result, &events.GuildEmojiAddEvent{
				GuildID: ev.GuildID,
				Emoji:   newEmoji,
			})
		} else if emojiChanged(old, newEmoji) {
			result = append(result, &events.GuildEmojiUpdateEvent{
				GuildID:  ev.GuildID,
				OldEmoji: old,
				NewEmoji: newEmoji,
			})
		}
	}
	for id, oldEmoji := range oldByID {
		if _, exists := newByID[id]; !exists {
			result = append(result, &events.GuildEmojiRemoveEvent{
				GuildID: ev.GuildID,
				Emoji:   oldEmoji,
			})
		}
	}

	return result
}

// deriveGuildStickerSyntheticEvents diffs the old and new sticker sets from a
// GUILD_STICKERS_UPDATE and fires Add/Remove/Update events per changed sticker.
// Returns nil when OldStickers is nil (cache cold).
func deriveGuildStickerSyntheticEvents(ev *events.GuildStickersUpdateEvent) []events.Event {
	if ev.OldStickers == nil {
		return nil
	}

	oldByID := make(map[discord.Snowflake]*discord.Sticker, len(ev.OldStickers))
	for _, s := range ev.OldStickers {
		oldByID[s.ID] = s
	}
	newByID := make(map[discord.Snowflake]*discord.Sticker, len(ev.NewStickers))
	for _, s := range ev.NewStickers {
		newByID[s.ID] = s
	}

	var result []events.Event

	for id, newSticker := range newByID {
		if old, exists := oldByID[id]; !exists {
			result = append(result, &events.GuildStickerAddEvent{
				GuildID: ev.GuildID,
				Sticker: newSticker,
			})
		} else if stickerChanged(old, newSticker) {
			result = append(result, &events.GuildStickerUpdateEvent{
				GuildID:    ev.GuildID,
				OldSticker: old,
				NewSticker: newSticker,
			})
		}
	}
	for id, oldSticker := range oldByID {
		if _, exists := newByID[id]; !exists {
			result = append(result, &events.GuildStickerRemoveEvent{
				GuildID: ev.GuildID,
				Sticker: oldSticker,
			})
		}
	}

	return result
}

func emojiChanged(old, new_ *discord.Emoji) bool {
	return old.Name != new_.Name ||
		!boolPtrEqual(old.Available, new_.Available) ||
		!boolPtrEqual(old.RequireColons, new_.RequireColons) ||
		!snowflakeSliceSetEqual(old.Roles, new_.Roles)
}

func snowflakeSliceSetEqual(a, b []discord.Snowflake) bool {
	if len(a) != len(b) {
		return false
	}
	counts := make(map[discord.Snowflake]int, len(a))
	for _, s := range a {
		counts[s]++
	}
	for _, s := range b {
		if counts[s] <= 0 {
			return false
		}
		counts[s]--
	}
	return true
}

func stickerChanged(old, new_ *discord.Sticker) bool {
	return old.Name != new_.Name ||
		old.Tags != new_.Tags ||
		!stringPtrEqual(old.Description, new_.Description) ||
		!boolPtrEqual(old.Available, new_.Available)
}

func boolPtrEqual(a, b *bool) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func stringPtrEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
