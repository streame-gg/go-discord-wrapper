package testdata

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func NewApplicationCommandPermissions() discord.ApplicationCommandPermission {
	return discord.ApplicationCommandPermission{
		ID: discord.RandomSnowflake(),
		Type: testutil.RandomItem(
			discord.ApplicationCommandPermissionTypeRole,
			discord.ApplicationCommandPermissionTypeUser,
			discord.ApplicationCommandPermissionTypeChannel,
		),
		Permission: testutil.RandomBool(),
	}
}

func NewGuildApplicationCommandPermissions() discord.GuildApplicationCommandPermissions {
	return discord.GuildApplicationCommandPermissions{
		ID:            discord.RandomSnowflake(),
		ApplicationID: discord.RandomSnowflake(),
		GuildID:       discord.RandomSnowflake(),
		Permissions: testutil.RandomArrayWithFilledItems(testutil.RandomNumberInRange(1, 100), func(arrayToFill *[]discord.ApplicationCommandPermission) {
			*arrayToFill = append(*arrayToFill, NewApplicationCommandPermissions())
		}),
	}
}
