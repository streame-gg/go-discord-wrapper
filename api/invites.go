package api

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// ── Param types ───────────────────────────────────────────────────────────────

type GetInviteParams struct {
	WithCounts            *bool
	WithExpiration        *bool
	GuildScheduledEventID *common.Snowflake
}

func (p GetInviteParams) toQuery() string {
	q := url.Values{}
	if p.WithCounts != nil {
		q.Set("with_counts", strconv.FormatBool(*p.WithCounts))
	}
	if p.WithExpiration != nil {
		q.Set("with_expiration", strconv.FormatBool(*p.WithExpiration))
	}
	if p.GuildScheduledEventID != nil {
		q.Set("guild_scheduled_event_id", (*p.GuildScheduledEventID).String())
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// ── Invite endpoints ──────────────────────────────────────────────────────────

// GetInvite returns the invite object for the given invite code.
func (c *RestClient) GetInvite(ctx context.Context, code string, params GetInviteParams) (*Invite, error) {
	path := "/invites/" + code + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var invite Invite
	if _, err := c.do(req, http.StatusOK, &invite); err != nil {
		return nil, err
	}

	return &invite, nil
}

// DeleteInvite deletes an invite by its code. Requires MANAGE_CHANNELS or MANAGE_GUILD.
func (c *RestClient) DeleteInvite(ctx context.Context, code string) (*Invite, error) {
	req, err := c.generateRequest(ctx, http.MethodDelete, "/invites/"+code, nil)
	if err != nil {
		return nil, err
	}

	var invite Invite
	if _, err := c.do(req, http.StatusOK, &invite); err != nil {
		return nil, err
	}

	return &invite, nil
}
