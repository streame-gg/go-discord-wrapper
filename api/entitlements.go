package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Param / response types ────────────────────────────────────────────────────

// https://docs.discord.com/developers/resources/entitlement#list-entitlements
type ListEntitlementsParams struct {
	UserID       *discord.Snowflake
	SkuIDs       []discord.Snowflake
	Before       *discord.Snowflake
	After        *discord.Snowflake
	Limit        *int
	GuildID      *discord.Snowflake
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

// https://docs.discord.com/developers/resources/entitlement#create-test-entitlement
type CreateTestEntitlementParams struct {
	SkuID     discord.Snowflake `json:"sku_id"`
	OwnerID   discord.Snowflake `json:"owner_id"`
	OwnerType int               `json:"owner_type"`
}

// ── Entitlement endpoints ─────────────────────────────────────────────────────

// ListEntitlements returns entitlements for an application.
// https://docs.discord.com/developers/resources/entitlement#list-entitlements
func (c *RestClient) ListEntitlements(ctx context.Context, appID discord.Snowflake, params ListEntitlementsParams) ([]*discord.Entitlement, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	path := "/applications/" + appID.String() + "/entitlements" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.Entitlement](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// https://docs.discord.com/developers/resources/entitlement#get-entitlement
func (c *RestClient) GetEntitlement(ctx context.Context, appID, entitlementID discord.Snowflake) (*discord.Entitlement, error) {
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

	return doRequest[discord.Entitlement](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// https://docs.discord.com/developers/resources/entitlement#create-test-entitlement
func (c *RestClient) CreateTestEntitlement(ctx context.Context, appID discord.Snowflake, params CreateTestEntitlementParams) (*discord.Entitlement, error) {
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

	return doRequest[discord.Entitlement](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// https://docs.discord.com/developers/resources/entitlement#consume-an-entitlement
func (c *RestClient) ConsumeEntitlement(ctx context.Context, appID, entitlementID discord.Snowflake) error {
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

	return doRequestWithoutResponse(c, req)
}

// https://docs.discord.com/developers/resources/entitlement#list-entitlements
type GetCurrentUserApplicationEntitlementsParams struct {
	SkuIDs          []discord.Snowflake
	ExcludeConsumed *bool
}

func (p GetCurrentUserApplicationEntitlementsParams) toQuery() string {
	q := url.Values{}
	if len(p.SkuIDs) > 0 {
		ids := make([]string, len(p.SkuIDs))
		for i, id := range p.SkuIDs {
			ids[i] = id.String()
		}
		q.Set("sku_ids", strings.Join(ids, ","))
	}
	if p.ExcludeConsumed != nil {
		q.Set("exclude_consumed", strconv.FormatBool(*p.ExcludeConsumed))
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// ListCurrentUserApplicationEntitlements returns entitlements for the current user for the given
// application. Requires an OAuth2 bearer token.
// https://docs.discord.com/developers/resources/entitlement#list-entitlements
func (c *RestClient) ListCurrentUserApplicationEntitlements(ctx context.Context, appID discord.Snowflake, params GetCurrentUserApplicationEntitlementsParams, userToken string) ([]*discord.Entitlement, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	path := "/users/@me/applications/" + appID.String() + "/entitlements" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, WithUserAuthorization(userToken))
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.Entitlement](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// https://docs.discord.com/developers/resources/entitlement#delete-test-entitlement
func (c *RestClient) DeleteTestEntitlement(ctx context.Context, appID, entitlementID discord.Snowflake) error {
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

	return doRequestWithoutResponse(c, req)
}
