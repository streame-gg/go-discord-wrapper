package api

import (
	"context"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (c *RestClient) ListSKUs(ctx context.Context, appID discord.Snowflake) (*[]*discord.SKU, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/applications/"+appID.String()+"/skus", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[[]*discord.SKU](c, req, map[int]bool{
		http.StatusOK: true,
	})
}
