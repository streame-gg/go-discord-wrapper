package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// ── Param / response types ────────────────────────────────────────────────────

type ListEntitlementsParams struct {
	UserID       *common.Snowflake
	SkuIDs       []common.Snowflake
	Before       *common.Snowflake
	After        *common.Snowflake
	Limit        *int
	GuildID      *common.Snowflake
	ExcludeEnded *bool
}

func (p ListEntitlementsParams) toQuery() string {
	q := url.Values{}
	if p.UserID != nil {
		q.Set("user_id", p.UserID.String())
	}
	if len(p.SkuIDs) > 0 {
		ids := make([]string, len(p.SkuIDs))
		for i, id := range p.SkuIDs {
			ids[i] = id.String()
		}
		q.Set("sku_ids", strings.Join(ids, ","))
	}
	if p.Before != nil {
		q.Set("before", p.Before.String())
	}
	if p.After != nil {
		q.Set("after", p.After.String())
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if p.GuildID != nil {
		q.Set("guild_id", p.GuildID.String())
	}
	if p.ExcludeEnded != nil {
		q.Set("exclude_ended", strconv.FormatBool(*p.ExcludeEnded))
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

type CreateTestEntitlementParams struct {
	SkuID     common.Snowflake `json:"sku_id"`
	OwnerID   common.Snowflake `json:"owner_id"`
	OwnerType int              `json:"owner_type"`
}

// ── Entitlement endpoints ─────────────────────────────────────────────────────

func (c *RestClient) ListEntitlements(ctx context.Context, appID common.Snowflake, params ListEntitlementsParams) ([]*common.Entitlement, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/entitlements" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var result []*common.Entitlement
	if _, err := c.do(req, []int{http.StatusOK}, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *RestClient) GetEntitlement(ctx context.Context, appID, entitlementID common.Snowflake) (*common.Entitlement, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	if err := entitlementID.Validate(); err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/entitlements/" + entitlementID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var result common.Entitlement
	if _, err := c.do(req, []int{http.StatusOK}, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *RestClient) CreateTestEntitlement(ctx context.Context, appID common.Snowflake, params CreateTestEntitlementParams) (*common.Entitlement, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/applications/"+appID.String()+"/entitlements", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var result common.Entitlement
	if _, err := c.do(req, []int{http.StatusOK}, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *RestClient) ConsumeEntitlement(ctx context.Context, appID, entitlementID common.Snowflake) error {
	if err := appID.Validate(); err != nil {
		return err
	}

	if err := entitlementID.Validate(); err != nil {
		return err
	}

	path := "/applications/" + appID.String() + "/entitlements/" + entitlementID.String() + "/consume"
	req, err := c.generateRequest(ctx, http.MethodPost, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	_, err = c.do(req, []int{http.StatusNoContent}, nil)
	return err
}

func (c *RestClient) DeleteTestEntitlement(ctx context.Context, appID, entitlementID common.Snowflake) error {
	if err := appID.Validate(); err != nil {
		return err
	}

	if err := entitlementID.Validate(); err != nil {
		return err
	}

	path := "/applications/" + appID.String() + "/entitlements/" + entitlementID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	_, err = c.do(req, []int{http.StatusNoContent}, nil)
	return err
}
