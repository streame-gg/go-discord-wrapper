package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Param / response types ────────────────────────────────────────────────────

// https://docs.discord.com/developers/resources/channel#start-thread-from-message
type CreateThreadFromMessageParams struct {
	Name                string              `json:"name"`
	AutoArchiveDuration discord.Option[int] `json:"auto_archive_duration,omitempty"`
	RateLimitPerUser    discord.Option[int] `json:"rate_limit_per_user,omitempty"`

	AuditLogReason *string `json:"-"`
}

// https://docs.discord.com/developers/resources/channel#start-thread-without-message
type CreateThreadParams struct {
	Name                string                              `json:"name"`
	AutoArchiveDuration discord.Option[int]                 `json:"auto_archive_duration,omitempty"`
	Type                discord.Option[discord.ChannelType] `json:"type,omitempty"`
	Invitable           discord.Option[bool]                `json:"invitable,omitempty"`
	RateLimitPerUser    discord.Option[int]                 `json:"rate_limit_per_user,omitempty"`

	AuditLogReason *string `json:"-"`
}

// CreateForumThreadParams is used to start a thread in a forum or media channel.
// https://docs.discord.com/developers/resources/channel#start-thread-in-forum-or-media-channel
type CreateForumThreadParams struct {
	Name                string                              `json:"name"`
	AutoArchiveDuration discord.Option[int]                 `json:"auto_archive_duration,omitempty"`
	RateLimitPerUser    discord.Option[int]                 `json:"rate_limit_per_user,omitempty"`
	Message             CreateMessageParams                 `json:"message"`
	AppliedTags         discord.Option[[]discord.Snowflake] `json:"applied_tags,omitempty"`

	AuditLogReason *string `json:"-"`
}

// https://docs.discord.com/developers/resources/channel#list-thread-members
type ListThreadMembersParams struct {
	WithMember *bool
	After      *discord.Snowflake
	Limit      *int
}

func (p ListThreadMembersParams) toQuery() string {
	q := url.Values{}
	if p.WithMember != nil && *p.WithMember {
		q.Set("with_member", "true")
	}
	if p.After != nil {
		q.Set("after", p.After.String())
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// https://docs.discord.com/developers/resources/channel#list-public-archived-threads
type ListArchivedThreadsParams struct {
	// Before paginates the public and private archived-thread listings, which
	// cursor on the threads' archive timestamp.
	Before *time.Time
	// BeforeID paginates ListJoinedPrivateArchivedThreads, which Discord cursors
	// on thread ID (snowflake) rather than archive timestamp. When set it takes
	// precedence over Before. https://docs.discord.com/developers/resources/channel#list-joined-private-archived-threads
	BeforeID *discord.Snowflake
	Limit    *int
}

func (p ListArchivedThreadsParams) toQuery() string {
	q := url.Values{}
	if p.Before != nil {
		q.Set("before", p.Before.Format(time.RFC3339))
	}
	if p.BeforeID != nil {
		q.Set("before", p.BeforeID.String())
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// ArchivedThreadsResponse is returned by the list-archived-threads endpoints.
// https://docs.discord.com/developers/resources/channel#list-public-archived-threads
type ArchivedThreadsResponse struct {
	Threads []*discord.Channel      `json:"threads"`
	Members []*discord.ThreadMember `json:"members"`
	HasMore bool                    `json:"has_more"`
}

// ActiveThreadsResponse is returned by ListActiveGuildThreads.
// https://docs.discord.com/developers/resources/guild#list-active-guild-threads
type ActiveThreadsResponse struct {
	Threads []*discord.Channel      `json:"threads"`
	Members []*discord.ThreadMember `json:"members"`
}

// ── Thread endpoints ──────────────────────────────────────────────────────────

// CreateThreadFromMessage starts a public thread from an existing message.
// https://docs.discord.com/developers/resources/channel#start-thread-from-message
func (c *RestClient) CreateThreadFromMessage(ctx context.Context, channelID, messageID discord.Snowflake, params CreateThreadFromMessageParams) (*discord.Channel, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	if err := messageID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/messages/" + messageID.String() + "/threads"

	req, err := c.generateRequest(ctx, http.MethodPost, path, bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(params.AuditLogReason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Channel](c, req, map[int]bool{
		http.StatusCreated: true,
	})
}

// CreateThread starts a thread that is not connected to an existing message.
// https://docs.discord.com/developers/resources/channel#start-thread-without-message
func (c *RestClient) CreateThread(ctx context.Context, channelID discord.Snowflake, params CreateThreadParams) (*discord.Channel, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/channels/"+channelID.String()+"/threads", bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(params.AuditLogReason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Channel](c, req, map[int]bool{
		http.StatusCreated: true,
	})
}

// CreateForumThread starts a thread in a forum or media channel.
// When params.Message.Files is non-empty the request is sent as multipart/form-data.
// https://docs.discord.com/developers/resources/channel#start-thread-in-forum-or-media-channel
func (c *RestClient) CreateForumThread(ctx context.Context, channelID discord.Snowflake, params CreateForumThreadParams) (*discord.Channel, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/threads"

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reqOpts := []func(req *http.Request){c.WithBotAuthorization(), WithAuditLogReason(params.AuditLogReason)}

	var req *http.Request
	if len(params.Message.Files) > 0 {
		buf, ct, err := buildMultipartMessage(body, params.Message.Files)
		if err != nil {
			return nil, err
		}
		req, err = c.generateRequest(ctx, http.MethodPost, path, buf, reqOpts...)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", ct)
	} else {
		req, err = c.generateRequest(ctx, http.MethodPost, path, bytes.NewReader(body), reqOpts...)
		if err != nil {
			return nil, err
		}
	}

	return doRequest[discord.Channel](c, req, map[int]bool{
		http.StatusCreated: true,
	})
}

// JoinThread adds the current user to a thread.
// https://docs.discord.com/developers/resources/channel#join-thread
func (c *RestClient) JoinThread(ctx context.Context, channelID discord.Snowflake) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	req, err := c.generateRequest(ctx, http.MethodPut, "/channels/"+channelID.String()+"/thread-members/@me", nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// LeaveThread removes the current user from a thread.
// https://docs.discord.com/developers/resources/channel#leave-thread
func (c *RestClient) LeaveThread(ctx context.Context, channelID discord.Snowflake) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, "/channels/"+channelID.String()+"/thread-members/@me", nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// AddThreadMember adds another user to a thread.
// https://docs.discord.com/developers/resources/channel#add-thread-member
func (c *RestClient) AddThreadMember(ctx context.Context, channelID, userID discord.Snowflake) error {
	if err := channelID.Validate(); err != nil {
		return err
	}
	if err := userID.Validate(); err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/thread-members/" + userID.String()
	req, err := c.generateRequest(ctx, http.MethodPut, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// RemoveThreadMember removes a user from a thread.
// https://docs.discord.com/developers/resources/channel#remove-thread-member
func (c *RestClient) RemoveThreadMember(ctx context.Context, channelID, userID discord.Snowflake) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/thread-members/" + userID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// GetThreadMember returns the thread member object for the given user.
// Set withMember to true to include the guild member object.
// https://docs.discord.com/developers/resources/channel#get-thread-member
func (c *RestClient) GetThreadMember(ctx context.Context, channelID, userID discord.Snowflake, withMember bool) (*discord.ThreadMember, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	if err := userID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/thread-members/" + userID.String()
	if withMember {
		path += "?with_member=true"
	}

	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.ThreadMember](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ListThreadMembers returns the members of a thread.
// https://docs.discord.com/developers/resources/channel#list-thread-members
func (c *RestClient) ListThreadMembers(ctx context.Context, channelID discord.Snowflake, params ListThreadMembersParams) ([]*discord.ThreadMember, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/thread-members" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.ThreadMember](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ListPublicArchivedThreads returns archived public threads in a channel, newest first.
// https://docs.discord.com/developers/resources/channel#list-public-archived-threads
func (c *RestClient) ListPublicArchivedThreads(ctx context.Context, channelID discord.Snowflake, params ListArchivedThreadsParams) (*ArchivedThreadsResponse, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/threads/archived/public" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[ArchivedThreadsResponse](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ListPrivateArchivedThreads returns archived private threads in a channel, newest first.
// https://docs.discord.com/developers/resources/channel#list-private-archived-threads
func (c *RestClient) ListPrivateArchivedThreads(ctx context.Context, channelID discord.Snowflake, params ListArchivedThreadsParams) (*ArchivedThreadsResponse, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/threads/archived/private" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[ArchivedThreadsResponse](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ListJoinedPrivateArchivedThreads returns private archived threads the current user has joined.
// https://docs.discord.com/developers/resources/channel#list-joined-private-archived-threads
func (c *RestClient) ListJoinedPrivateArchivedThreads(ctx context.Context, channelID discord.Snowflake, params ListArchivedThreadsParams) (*ArchivedThreadsResponse, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/users/@me/threads/archived/private" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[ArchivedThreadsResponse](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// SearchThreadsParams holds the query parameters for searching threads in a channel.
// https://docs.discord.com/developers/resources/channel
type SearchThreadsParams struct {
	Name     *string
	Archived *bool
	Tag      *discord.Snowflake
	MinID    *discord.Snowflake
	MaxID    *discord.Snowflake
	Limit    *int
	Offset   *int
}

func (p SearchThreadsParams) toQuery() string {
	q := url.Values{}
	if p.Name != nil {
		q.Set("name", *p.Name)
	}
	if p.Archived != nil {
		q.Set("archived", strconv.FormatBool(*p.Archived))
	}
	if p.Tag != nil {
		q.Set("tag", p.Tag.String())
	}
	if p.MinID != nil {
		q.Set("min_id", p.MinID.String())
	}
	if p.MaxID != nil {
		q.Set("max_id", p.MaxID.String())
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if p.Offset != nil {
		q.Set("offset", strconv.Itoa(*p.Offset))
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// SearchThreads searches for threads in a forum or media channel. Returns 202 while the index
// is still being built (result will be empty).
// https://docs.discord.com/developers/resources/channel
func (c *RestClient) SearchThreads(ctx context.Context, channelID discord.Snowflake, params SearchThreadsParams) (*discord.ThreadSearchResponse, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/threads/search" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.ThreadSearchResponse](c, req, map[int]bool{
		http.StatusOK:       true,
		http.StatusAccepted: true,
	})
}

// ListActiveGuildThreads returns all active threads in the guild that the current user can access.
// https://docs.discord.com/developers/resources/guild#list-active-guild-threads
func (c *RestClient) ListActiveGuildThreads(ctx context.Context, guildID discord.Snowflake) (*ActiveThreadsResponse, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/threads/active", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[ActiveThreadsResponse](c, req, map[int]bool{
		http.StatusOK: true,
	})
}
