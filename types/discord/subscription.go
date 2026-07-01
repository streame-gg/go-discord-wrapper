package discord

import (
	"time"
)

// https://docs.discord.com/developers/resources/subscription#subscription-object
type Subscription struct {
	ID                 Snowflake          `json:"id"`
	UserID             Snowflake          `json:"user_id"`
	SKUIDs             []Snowflake        `json:"sku_ids"`
	EntitlementIDs     []Snowflake        `json:"entitlement_ids"`
	RenewalSKUIDs      []Snowflake        `json:"renewal_sku_ids,omitempty"`
	CurrentPeriodStart time.Time          `json:"current_period_start"`
	CurrentPeriodEnd   time.Time          `json:"current_period_end"`
	Status             SubscriptionStatus `json:"status"`
	CanceledAt         *time.Time         `json:"canceled_at"`
	Country            string             `json:"country,omitempty"`

	User *User `json:"-"`
}

// https://docs.discord.com/developers/resources/subscription#subscription-statuses
type SubscriptionStatus uint8

const (
	SubscriptionStatusActive   SubscriptionStatus = 0
	SubscriptionStatusInactive SubscriptionStatus = 1
	SubscriptionStatusEnding   SubscriptionStatus = 2
)
