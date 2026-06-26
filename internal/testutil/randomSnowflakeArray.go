package testutil

import "github.com/streame-gg/go-discord-wrapper/types/discord"

func RandomSnowflakeArray(arraySize int) []discord.Snowflake {
	arr := make([]discord.Snowflake, 0, arraySize)
	for i := 0; i < arraySize; i++ {
		arr = append(arr, discord.RandomSnowflake())
	}
	return arr
}
