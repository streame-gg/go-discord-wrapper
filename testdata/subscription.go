package testdata

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func NewSubscription() map[string]interface{} {
	return map[string]interface{}{
		"id":                   discord.RandomSnowflake(),
		"user_id":              discord.RandomSnowflake(),
		"sku_ids":              testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 5)),
		"entitlement_ids":      testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 5)),
		"renewal_sku_ids":      testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 5)),
		"current_period_start": testutil.RandomTime(),
		"current_period_end":   testutil.RandomTime(),
		"status": testutil.RandomItem(
			discord.SubscriptionStatusActive,
			discord.SubscriptionStatusInactive,
			discord.SubscriptionStatusEnding,
		),
		"canceled_at": testutil.RandomTime(),
		"country":     testutil.RandomString(2),
	}
}
