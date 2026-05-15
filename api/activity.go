package api

import (
	"context"
	"net/http"
	"net/url"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (c *RestClient) GetActivityInstance(ctx context.Context, appID discord.Snowflake, instanceID string) (*discord.ActivityInstance, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/activity-instances/" + url.PathEscape(instanceID)
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.ActivityInstance](c, req, map[int]bool{
		http.StatusOK: true,
	})
}
