package discord

import (
	"time"
)

// https://docs.discord.com/developers/resources/entitlement#entitlement-object
type Entitlement struct {
	ID            Snowflake       `json:"id"`
	SkuID         Snowflake       `json:"sku_id"`
	ApplicationID Snowflake       `json:"application_id"`
	UserID        *Snowflake      `json:"user_id,omitempty"`
	Type          EntitlementType `json:"type"`
	Deleted       bool            `json:"deleted"`
	StartsAt      *time.Time      `json:"starts_at"`
	EndsAt        *time.Time      `json:"ends_at"`
	GuildID       *Snowflake      `json:"guild_id,omitempty"`
	Consumed      *bool           `json:"consumed,omitempty"`
}

// https://docs.discord.com/developers/resources/entitlement#entitlement-object-entitlement-types
type EntitlementType int

const (
	EntitlementTypePurchase                EntitlementType = 1
	EntitlementTypePremiumSubscription     EntitlementType = 2
	EntitlementTypeDeveloperGift           EntitlementType = 3
	EntitlementTypeTestModePurchase        EntitlementType = 4
	EntitlementTypeFreePurchase            EntitlementType = 5
	EntitlementTypeUserGift                EntitlementType = 6
	EntitlementTypePremiumPurchase         EntitlementType = 7
	EntitlementTypeApplicationSubscription EntitlementType = 8
)
