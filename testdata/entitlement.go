package testdata

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func NewEntitlement() map[string]interface{} {
	return map[string]interface{}{
		"id":             discord.RandomSnowflake(),
		"sku_id":         discord.RandomSnowflake(),
		"application_id": discord.RandomSnowflake(),
		"user_id":        discord.RandomSnowflake(),
		"type": testutil.RandomItem(
			discord.EntitlementTypePurchase,
			discord.EntitlementTypePremiumSubscription,
			discord.EntitlementTypeDeveloperGift,
			discord.EntitlementTypeTestModePurchase,
			discord.EntitlementTypeFreePurchase,
			discord.EntitlementTypeUserGift,
			discord.EntitlementTypePremiumPurchase,
			discord.EntitlementTypeApplicationSubscription,
		),
		"deleted":   testutil.RandomBool(),
		"starts_at": testutil.RandomTime(),
		"ends_at":   testutil.RandomTime(),
		"guild_id":  discord.RandomSnowflake(),
		"consumed":  testutil.RandomBool(),
	}
}
