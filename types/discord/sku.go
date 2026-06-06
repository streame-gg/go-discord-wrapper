package discord

// https://docs.discord.com/developers/resources/sku#sku-object
type SKU struct {
	ID            Snowflake `json:"id"`
	Type          int       `json:"type"`
	ApplicationID Snowflake `json:"application_id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Flags         int       `json:"flags"`
}
