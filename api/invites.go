package api

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/resources/invite#get-invite
type GetInviteParams struct {
	WithCounts            *bool
	GuildScheduledEventID *discord.Snowflake
}

// https://docs.discord.com/developers/resources/invite#delete-invite
type DeleteInviteOptions struct {
	Reason string
}

func (p GetInviteParams) toQuery() string {
	q := url.Values{}
	if p.WithCounts != nil {
		q.Set("with_counts", strconv.FormatBool(*p.WithCounts))
	}
	if p.GuildScheduledEventID != nil {
		q.Set("guild_scheduled_event_id", (*p.GuildScheduledEventID).String())
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// GetInvite returns the invite object for the given invite code.
func (c *RestClient) GetInvite(ctx context.Context, code string, params GetInviteParams) (*discord.Invite, error) {
	path := "/invites/" + url.PathEscape(code) + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Invite](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteInvite deletes an invite by its code. Requires MANAGE_CHANNELS or MANAGE_GUILD.
func (c *RestClient) DeleteInvite(ctx context.Context, code string, opts *DeleteInviteOptions) (*discord.Invite, error) {
	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodDelete, "/invites/"+url.PathEscape(code), nil, c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		return doRequest[discord.Invite](c, req, map[int]bool{
			http.StatusOK: true,
		})
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, "/invites/"+url.PathEscape(code), nil, c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Invite](c, req, map[int]bool{
		http.StatusOK: true,
	})
}
