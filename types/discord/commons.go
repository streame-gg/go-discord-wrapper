package discord

// https://docs.discord.com/developers/events/gateway-events#activity-object
type Activity struct {
	Type    int     `json:"type"`
	PartyID *string `json:"party_id,omitempty"`
}

// https://docs.discord.com/developers/resources/message#role-subscription-data-object
type RoleSubscriptionData struct {
	RoleSubscriptionListingID Snowflake `json:"role_subscription_listing_id"`
	TierName                  string    `json:"tier_name"`
	TotalMonthsSubscribed     int       `json:"total_months_subscribed"`
	IsRenewal                 bool      `json:"is_renewal"`
}
