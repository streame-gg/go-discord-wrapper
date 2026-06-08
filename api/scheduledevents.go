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

// https://docs.discord.com/developers/resources/guild-scheduled-event#create-guild-scheduled-event
type CreateGuildScheduledEventParams struct {
	ChannelID          *discord.Snowflake                         `json:"channel_id,omitempty"`
	EntityMetadata     *discord.GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Name               string                                     `json:"name"`
	PrivacyLevel       discord.GuildScheduledEventPrivacyLevel    `json:"privacy_level"`
	ScheduledStartTime time.Time                                  `json:"scheduled_start_time"`
	ScheduledEndTime   *time.Time                                 `json:"scheduled_end_time,omitempty"`
	Description        *string                                    `json:"description,omitempty"`
	EntityType         discord.GuildScheduledEventEntityType      `json:"entity_type"`
	// Image is a base64-encoded image data URI for the event cover image.
	Image *string `json:"image,omitempty"`
}

// https://docs.discord.com/developers/resources/guild-scheduled-event#modify-guild-scheduled-event
type ModifyGuildScheduledEventParams struct {
	ChannelID          *discord.Snowflake                         `json:"channel_id,omitempty"`
	EntityMetadata     *discord.GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Name               *string                                    `json:"name,omitempty"`
	PrivacyLevel       *discord.GuildScheduledEventPrivacyLevel   `json:"privacy_level,omitempty"`
	ScheduledStartTime *time.Time                                 `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   *time.Time                                 `json:"scheduled_end_time,omitempty"`
	Description        *string                                    `json:"description,omitempty"`
	EntityType         *discord.GuildScheduledEventEntityType     `json:"entity_type,omitempty"`
	Status             *discord.GuildScheduledEventStatus         `json:"status,omitempty"`
	Image              *string                                    `json:"image,omitempty"`
}

// https://docs.discord.com/developers/resources/guild-scheduled-event#get-guild-scheduled-event-users
type GetGuildScheduledEventUsersParams struct {
	Limit      *int
	WithMember *bool
	Before     *discord.Snowflake
	After      *discord.Snowflake
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

// ── Scheduled event endpoints ─────────────────────────────────────────────────

// ListGuildScheduledEvents returns all scheduled events for a guild.
// Set withUserCount to true to include subscriber counts.
// https://docs.discord.com/developers/resources/guild-scheduled-event#list-scheduled-events-for-guild
func (c *RestClient) ListGuildScheduledEvents(ctx context.Context, guildID discord.Snowflake, withUserCount bool) ([]*discord.GuildScheduledEvent, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/scheduled-events"
	if withUserCount {
		path += "?with_user_count=true"
	}

	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.GuildScheduledEvent](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// CreateGuildScheduledEvent creates a new scheduled event in a guild.
// https://docs.discord.com/developers/resources/guild-scheduled-event#create-guild-scheduled-event
func (c *RestClient) CreateGuildScheduledEvent(ctx context.Context, guildID discord.Snowflake, params CreateGuildScheduledEventParams) (*discord.GuildScheduledEvent, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/scheduled-events", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildScheduledEvent](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetGuildScheduledEvent returns a single scheduled event.
// Set withUserCount to true to include subscriber count.
// https://docs.discord.com/developers/resources/guild-scheduled-event#get-guild-scheduled-event
func (c *RestClient) GetGuildScheduledEvent(ctx context.Context, guildID, eventID discord.Snowflake, withUserCount bool) (*discord.GuildScheduledEvent, error) {
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

	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildScheduledEvent](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyGuildScheduledEvent updates a scheduled event.
// https://docs.discord.com/developers/resources/guild-scheduled-event#modify-guild-scheduled-event
func (c *RestClient) ModifyGuildScheduledEvent(ctx context.Context, guildID, eventID discord.Snowflake, params ModifyGuildScheduledEventParams) (*discord.GuildScheduledEvent, error) {
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
	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildScheduledEvent](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteGuildScheduledEvent deletes a scheduled event.
// https://docs.discord.com/developers/resources/guild-scheduled-event#delete-guild-scheduled-event
func (c *RestClient) DeleteGuildScheduledEvent(ctx context.Context, guildID, eventID discord.Snowflake) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	if err := eventID.Validate(); err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/scheduled-events/" + eventID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// ListGuildScheduledEventUsers returns users subscribed to a scheduled event.
// https://docs.discord.com/developers/resources/guild-scheduled-event#get-guild-scheduled-event-users
func (c *RestClient) ListGuildScheduledEventUsers(ctx context.Context, guildID, eventID discord.Snowflake, params GetGuildScheduledEventUsersParams) ([]*discord.GuildScheduledEventUser, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := eventID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/scheduled-events/" + eventID.String() + "/users" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.GuildScheduledEventUser](c, req, map[int]bool{
		http.StatusOK: true,
	})
}
