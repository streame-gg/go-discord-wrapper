package testdata

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func NewApplicationCommandPermissions() map[string]interface{} {
	return map[string]interface{}{
		"id": discord.RandomSnowflake(),
		"type": testutil.RandomItem(
			discord.ApplicationCommandPermissionTypeRole,
			discord.ApplicationCommandPermissionTypeUser,
			discord.ApplicationCommandPermissionTypeChannel,
		),
		"permission": testutil.RandomBool(),
	}
}

func NewGuildApplicationCommandPermissions() map[string]interface{} {
	return map[string]interface{}{
		"id":             discord.RandomSnowflake(),
		"application_id": discord.RandomSnowflake(),
		"guild_id":       discord.RandomSnowflake(),
		"permissions": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 3), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, NewApplicationCommandPermissions())
		}),
	}
}
