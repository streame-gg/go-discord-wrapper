package testdata

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func NewGuildEmojisUpdateEventPayload() map[string]interface{} {
	return map[string]interface{}{
		"guild_id": discord.RandomSnowflake(),
		"emojis": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 50), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewEmoji())
		}),
	}
}

func NewGuildStickersUpdateEventPayload() map[string]interface{} {
	return map[string]interface{}{
		"guild_id": discord.RandomSnowflake(),
		"stickers": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 50), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewSticker())
		}),
	}
}
