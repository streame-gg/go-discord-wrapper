package common

import "time"

type Entitlement struct {
	ID            Snowflake       `json:"id"`
	SkuID         Snowflake       `json:"sku_id"`
	ApplicationID Snowflake       `json:"application_id"`
	UserID        *Snowflake      `json:"user_id"`
	Type          EntitlementType `json:"type"`
	Deleted       bool            `json:"deleted"`
	StartsAt      *time.Time      `json:"starts_at,omitempty"`
	EndsAt        *time.Time      `json:"ends_at,omitempty"`
	GuildID       *Snowflake      `json:"guild_id,omitempty"`
	Consumed      bool            `json:"consumed,omitempty"`
}

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
