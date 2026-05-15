package discord

import (
	"time"
)

type Subscription struct {
	ID                 Snowflake   `json:"id"`
	UserID             Snowflake   `json:"user_id"`
	SKUIDs             []Snowflake `json:"sku_ids"`
	EntitlementIDs     []Snowflake `json:"entitlement_ids"`
	CurrentPeriodStart time.Time   `json:"current_period_start"`
	CurrentPeriodEnd   time.Time   `json:"current_period_end"`
	Status             int         `json:"status"`
	CanceledAt         *time.Time  `json:"canceled_at,omitempty"`
	Country            *string     `json:"country,omitempty"`
}
