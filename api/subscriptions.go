package api

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// ── Param types ───────────────────────────────────────────────────────────────

type ListSKUSubscriptionsParams struct {
	Before *common.Snowflake
	After  *common.Snowflake
	Limit  *int
	UserID *common.Snowflake
}

func (p ListSKUSubscriptionsParams) toQuery() string {
	q := url.Values{}
	if p.Before != nil {
		q.Set("before", p.Before.String())
	}
	if p.After != nil {
		q.Set("after", p.After.String())
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if p.UserID != nil {
		q.Set("user_id", p.UserID.String())
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// ── Subscription endpoints ────────────────────────────────────────────────────

func (c *RestClient) ListSKUSubscriptions(ctx context.Context, skuID common.Snowflake, params ListSKUSubscriptionsParams) (*[]*common.Subscription, error) {
	if err := skuID.Validate(); err != nil {
		return nil, err
	}

	path := "/skus/" + skuID.String() + "/subscriptions" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[[]*common.Subscription](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

func (c *RestClient) GetSKUSubscription(ctx context.Context, skuID, subscriptionID common.Snowflake) (*common.Subscription, error) {
	if err := skuID.Validate(); err != nil {
		return nil, err
	}

	if err := subscriptionID.Validate(); err != nil {
		return nil, err
	}

	path := "/skus/" + skuID.String() + "/subscriptions/" + subscriptionID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[common.Subscription](c, req, map[int]bool{
		http.StatusOK: true,
	})
}
