package api

import (
	"context"
	"net/http"
	"net/url"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

func (c *RestClient) GetActivityInstance(ctx context.Context, appID common.Snowflake, instanceID string) (*common.ActivityInstance, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/activity-instances/" + url.PathEscape(instanceID)
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result common.ActivityInstance
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
