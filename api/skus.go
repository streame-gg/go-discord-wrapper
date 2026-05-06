package api

import (
	"context"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

func (c *RestClient) ListSKUs(ctx context.Context, appID common.Snowflake) ([]*common.SKU, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/applications/"+appID.String()+"/skus", nil)
	if err != nil {
		return nil, err
	}

	var result []*common.SKU
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return result, nil
}
