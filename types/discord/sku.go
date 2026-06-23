package discord

// https://docs.discord.com/developers/resources/sku#sku-object
type SKU struct {
	ID            Snowflake `json:"id"`
	Type          SKUType   `json:"type"`
	ApplicationID Snowflake `json:"application_id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Flags         SKUFlag   `json:"flags"`
}

// https://docs.discord.com/developers/resources/sku#sku-object-sku-flags
type SKUFlag uint64

const (
	SKUFlagAvailable         SKUFlag = 1 << 2
	SKUFlagGuildSubscription SKUFlag = 1 << 7
	SKUFlagUserSubscription  SKUFlag = 1 << 8
)

// https://docs.discord.com/developers/resources/sku#sku-object-sku-types
type SKUType uint8

const (
	SKUTypeDurable           SKUType = 2
	SKUTypeConsumable        SKUType = 3
	SKUTypeSubscription      SKUType = 5
	SKUTypeSubscriptionGroup SKUType = 6
)
