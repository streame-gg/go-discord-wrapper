package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// ── Param / response types ────────────────────────────────────────────────────

type CreateThreadFromMessageParams struct {
	Name                string `json:"name"`
	AutoArchiveDuration *int   `json:"auto_archive_duration,omitempty"`
	RateLimitPerUser    *int   `json:"rate_limit_per_user,omitempty"`
}

type CreateThreadParams struct {
	Name                string              `json:"name"`
	AutoArchiveDuration *int                `json:"auto_archive_duration,omitempty"`
	Type                *common.ChannelType `json:"type,omitempty"`
	Invitable           *bool               `json:"invitable,omitempty"`
	RateLimitPerUser    *int                `json:"rate_limit_per_user,omitempty"`
}

// CreateForumThreadParams is used to start a thread in a forum or media channel.
type CreateForumThreadParams struct {
	Name                string              `json:"name"`
	AutoArchiveDuration *int                `json:"auto_archive_duration,omitempty"`
	RateLimitPerUser    *int                `json:"rate_limit_per_user,omitempty"`
	Message             CreateMessageParams `json:"message"`
	AppliedTags         []common.Snowflake  `json:"applied_tags,omitempty"`
}

type ListThreadMembersParams struct {
	WithMember *bool
	After      *common.Snowflake
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

type ListArchivedThreadsParams struct {
	Before *time.Time
	Limit  *int
}

func (p ListArchivedThreadsParams) toQuery() string {
	q := url.Values{}
	if p.Before != nil {
		q.Set("before", p.Before.Format(time.RFC3339))
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
type ArchivedThreadsResponse struct {
	Threads []*common.Channel      `json:"threads"`
	Members []*common.ThreadMember `json:"members"`
	HasMore bool                   `json:"has_more"`
}

// ActiveThreadsResponse is returned by ListActiveGuildThreads.
type ActiveThreadsResponse struct {
	Threads []*common.Channel      `json:"threads"`
	Members []*common.ThreadMember `json:"members"`
}

// ── Thread endpoints ──────────────────────────────────────────────────────────

// CreateThreadFromMessage starts a public thread from an existing message.
func (c *RestClient) CreateThreadFromMessage(ctx context.Context, channelID, messageID common.Snowflake, params CreateThreadFromMessageParams) (*common.Channel, error) {
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
	req, err := c.generateRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var thread common.Channel
	if _, err := c.do(req, http.StatusCreated, &thread); err != nil {
		return nil, err
	}

	return &thread, nil
}

// CreateThread starts a thread that is not connected to an existing message.
func (c *RestClient) CreateThread(ctx context.Context, channelID common.Snowflake, params CreateThreadParams) (*common.Channel, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/channels/"+channelID.String()+"/threads", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var thread common.Channel
	if _, err := c.do(req, http.StatusCreated, &thread); err != nil {
		return nil, err
	}

	return &thread, nil
}

// CreateForumThread starts a thread in a forum or media channel.
func (c *RestClient) CreateForumThread(ctx context.Context, channelID common.Snowflake, params CreateForumThreadParams) (*common.Channel, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/channels/"+channelID.String()+"/threads", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var thread common.Channel
	if _, err := c.do(req, http.StatusCreated, &thread); err != nil {
		return nil, err
	}

	return &thread, nil
}

// JoinThread adds the current user to a thread.
func (c *RestClient) JoinThread(ctx context.Context, channelID common.Snowflake) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	req, err := c.generateRequest(ctx, http.MethodPut, "/channels/"+channelID.String()+"/thread-members/@me", nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// LeaveThread removes the current user from a thread.
func (c *RestClient) LeaveThread(ctx context.Context, channelID common.Snowflake) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, "/channels/"+channelID.String()+"/thread-members/@me", nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// AddThreadMember adds another user to a thread.
func (c *RestClient) AddThreadMember(ctx context.Context, channelID, userID common.Snowflake) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/thread-members/" + userID.String()
	req, err := c.generateRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// RemoveThreadMember removes a user from a thread.
func (c *RestClient) RemoveThreadMember(ctx context.Context, channelID, userID common.Snowflake) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/thread-members/" + userID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// GetThreadMember returns the thread member object for the given user.
// Set withMember to true to include the guild member object.
func (c *RestClient) GetThreadMember(ctx context.Context, channelID, userID common.Snowflake, withMember bool) (*common.ThreadMember, error) {
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

	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var member common.ThreadMember
	if _, err := c.do(req, http.StatusOK, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

// ListThreadMembers returns the members of a thread.
func (c *RestClient) ListThreadMembers(ctx context.Context, channelID common.Snowflake, params ListThreadMembersParams) ([]*common.ThreadMember, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/thread-members" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var members []*common.ThreadMember
	if _, err := c.do(req, http.StatusOK, &members); err != nil {
		return nil, err
	}

	return members, nil
}

// ListPublicArchivedThreads returns archived public threads in a channel, newest first.
func (c *RestClient) ListPublicArchivedThreads(ctx context.Context, channelID common.Snowflake, params ListArchivedThreadsParams) (*ArchivedThreadsResponse, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/threads/archived/public" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result ArchivedThreadsResponse
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListPrivateArchivedThreads returns archived private threads in a channel, newest first.
func (c *RestClient) ListPrivateArchivedThreads(ctx context.Context, channelID common.Snowflake, params ListArchivedThreadsParams) (*ArchivedThreadsResponse, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/threads/archived/private" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result ArchivedThreadsResponse
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListJoinedPrivateArchivedThreads returns private archived threads the current user has joined.
func (c *RestClient) ListJoinedPrivateArchivedThreads(ctx context.Context, channelID common.Snowflake, params ListArchivedThreadsParams) (*ArchivedThreadsResponse, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/users/@me/threads/archived/private" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result ArchivedThreadsResponse
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListActiveGuildThreads returns all active threads in the guild that the current user can access.
func (c *RestClient) ListActiveGuildThreads(ctx context.Context, guildID common.Snowflake) (*ActiveThreadsResponse, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/threads/active", nil)
	if err != nil {
		return nil, err
	}

	var result ActiveThreadsResponse
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
