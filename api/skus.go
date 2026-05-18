package api

import (
	"context"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ListSKUs returns all SKUs for the given application, keyed by SKU ID.
func (c *RestClient) ListSKUs(ctx context.Context, appID discord.Snowflake) (*collection.Collection[discord.Snowflake, *discord.SKU], error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/applications/"+appID.String()+"/skus", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	skus, err := doRequestSlice[discord.SKU](c, req, map[int]bool{
		http.StatusOK: true,
	})
	if err != nil {
		return nil, err
	}

	coll := collection.NewWithCapacity[discord.Snowflake, *discord.SKU](len(skus))
	for _, s := range skus {
		coll.Set(s.ID, s)
	}
	return coll, nil
}
