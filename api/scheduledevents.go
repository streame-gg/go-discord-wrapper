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

type CreateGuildScheduledEventParams struct {
	ChannelID          *common.Snowflake                         `json:"channel_id,omitempty"`
	EntityMetadata     *common.GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Name               string                                    `json:"name"`
	PrivacyLevel       common.GuildScheduledEventPrivacyLevel    `json:"privacy_level"`
	ScheduledStartTime time.Time                                 `json:"scheduled_start_time"`
	ScheduledEndTime   *time.Time                                `json:"scheduled_end_time,omitempty"`
	Description        *string                                   `json:"description,omitempty"`
	EntityType         common.GuildScheduledEventEntityType      `json:"entity_type"`
	// Image is a base64-encoded image data URI for the event cover image.
	Image *string `json:"image,omitempty"`
}

type ModifyGuildScheduledEventParams struct {
	ChannelID          *common.Snowflake                         `json:"channel_id,omitempty"`
	EntityMetadata     *common.GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Name               *string                                   `json:"name,omitempty"`
	PrivacyLevel       *common.GuildScheduledEventPrivacyLevel   `json:"privacy_level,omitempty"`
	ScheduledStartTime *time.Time                                `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   *time.Time                                `json:"scheduled_end_time,omitempty"`
	Description        *string                                   `json:"description,omitempty"`
	EntityType         *common.GuildScheduledEventEntityType     `json:"entity_type,omitempty"`
	Status             *common.GuildScheduledEventStatus         `json:"status,omitempty"`
	Image              *string                                   `json:"image,omitempty"`
}

type GetGuildScheduledEventUsersParams struct {
	Limit      *int
	WithMember *bool
	Before     *common.Snowflake
	After      *common.Snowflake
}

func (p GetGuildScheduledEventUsersParams) toQuery() string {
	q := url.Values{}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if p.WithMember != nil && *p.WithMember {
		q.Set("with_member", "true")
	}
	if p.Before != nil {
		q.Set("before", p.Before.String())
	}
	if p.After != nil {
		q.Set("after", p.After.String())
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// GuildScheduledEventUser is an entry in the list returned by GetGuildScheduledEventUsers.
type GuildScheduledEventUser struct {
	GuildScheduledEventID common.Snowflake    `json:"guild_scheduled_event_id"`
	User                  common.User         `json:"user"`
	Member                *common.GuildMember `json:"member,omitempty"`
}

// ── Scheduled event endpoints ─────────────────────────────────────────────────

// ListGuildScheduledEvents returns all scheduled events for a guild.
// Set withUserCount to true to include subscriber counts.
func (c *RestClient) ListGuildScheduledEvents(ctx context.Context, guildID common.Snowflake, withUserCount bool) ([]*common.GuildScheduledEvent, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/scheduled-events"
	if withUserCount {
		path += "?with_user_count=true"
	}

	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var events []*common.GuildScheduledEvent
	if _, err := c.do(req, http.StatusOK, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// CreateGuildScheduledEvent creates a new scheduled event in a guild.
func (c *RestClient) CreateGuildScheduledEvent(ctx context.Context, guildID common.Snowflake, params CreateGuildScheduledEventParams) (*common.GuildScheduledEvent, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/scheduled-events", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var event common.GuildScheduledEvent
	if _, err := c.do(req, http.StatusOK, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

// GetGuildScheduledEvent returns a single scheduled event.
// Set withUserCount to true to include subscriber count.
func (c *RestClient) GetGuildScheduledEvent(ctx context.Context, guildID, eventID common.Snowflake, withUserCount bool) (*common.GuildScheduledEvent, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := eventID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/scheduled-events/" + eventID.String()
	if withUserCount {
		path += "?with_user_count=true"
	}

	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var event common.GuildScheduledEvent
	if _, err := c.do(req, http.StatusOK, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

// ModifyGuildScheduledEvent updates a scheduled event.
func (c *RestClient) ModifyGuildScheduledEvent(ctx context.Context, guildID, eventID common.Snowflake, params ModifyGuildScheduledEventParams) (*common.GuildScheduledEvent, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := eventID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/scheduled-events/" + eventID.String()
	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var event common.GuildScheduledEvent
	if _, err := c.do(req, http.StatusOK, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

// DeleteGuildScheduledEvent deletes a scheduled event.
func (c *RestClient) DeleteGuildScheduledEvent(ctx context.Context, guildID, eventID common.Snowflake) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	if err := eventID.Validate(); err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/scheduled-events/" + eventID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// GetGuildScheduledEventUsers returns users subscribed to a scheduled event.
func (c *RestClient) GetGuildScheduledEventUsers(ctx context.Context, guildID, eventID common.Snowflake, params GetGuildScheduledEventUsersParams) ([]*GuildScheduledEventUser, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := eventID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/scheduled-events/" + eventID.String() + "/users" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var users []*GuildScheduledEventUser
	if _, err := c.do(req, http.StatusOK, &users); err != nil {
		return nil, err
	}

	return users, nil
}
