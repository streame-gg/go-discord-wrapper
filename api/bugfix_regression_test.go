package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// TestMaxResponseBodySize_UnlimitedReadsFullBody is a regression test for the
// WithMaxResponseBodySize(-1) bug: a negative limit was passed straight to
// io.LimitReader, which returns EOF immediately, so the whole response body was
// dropped and decoding failed. -1 must mean "no limit".
func (su *restSuite) TestMaxResponseBodySize_UnlimitedReadsFullBody() {
	// A body comfortably larger than any small default would allow.
	big := strings.Repeat("a", 2<<20) // 2 MiB
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"1","name":"`+big+`"}`)
	}))
	defer ts.Close()

	rc, err := NewRestClient("test-token",
		options.WithBaseURL(ts.URL),
		options.WithRateLimiting(options.RateLimiterOptions{Disabled: true}),
		options.WithMaxResponseBodySize(-1), // unlimited
	)
	su.Require().NoError(err)

	ch, err := rc.GetChannel(context.Background(), discord.Snowflake(175928847299117063))
	su.Require().NoError(err, "unlimited body size must not truncate the response")
	su.Require().NotNil(ch)
	su.Len(ch.Name, len(big), "the full (large) body should have been read and decoded")
}

// TestSetVoiceChannelStatus_NilClearsWithNull is a regression test for the
// nullable status bug: a nil *string status must marshal to JSON null (which
// clears the channel status), not "".
func (su *restSuite) TestSetVoiceChannelStatus_NilClearsWithNull() {
	var gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	rc, err := NewRestClient("test-token",
		options.WithBaseURL(ts.URL),
		options.WithRateLimiting(options.RateLimiterOptions{Disabled: true}),
	)
	su.Require().NoError(err)

	su.Require().NoError(rc.SetVoiceChannelStatus(context.Background(), discord.Snowflake(175928847299117063), nil))
	su.JSONEq(`{"status":null}`, gotBody, "nil status must clear via JSON null, not empty string")

	status := "live"
	su.Require().NoError(rc.SetVoiceChannelStatus(context.Background(), discord.Snowflake(175928847299117063), &status))
	su.JSONEq(`{"status":"live"}`, gotBody, "non-nil status must send the string value")
}

// TestModifyGuildIncidentActions_NilClearsWithNull verifies the PUT body sends
// nil pause timestamps as JSON null (re-enabling invites/DMs), and targets the
// correct method/path. The endpoint replaces the guild's incident actions, so
// omitting fields would be wrong here.
func (su *restSuite) TestModifyGuildIncidentActions_NilClearsWithNull() {
	var gotBody, gotMethod, gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody, gotMethod, gotPath = string(b), r.Method, r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}))
	defer ts.Close()

	rc, err := NewRestClient("test-token",
		options.WithBaseURL(ts.URL),
		options.WithRateLimiting(options.RateLimiterOptions{Disabled: true}),
	)
	su.Require().NoError(err)

	const guildID = discord.Snowflake(175928847299117063)
	data, err := rc.ModifyGuildIncidentActions(context.Background(), guildID, nil)
	su.Require().NoError(err)
	su.Require().NotNil(data)
	su.Equal(http.MethodPut, gotMethod)
	su.True(strings.HasSuffix(gotPath, "/guilds/"+guildID.String()+"/incident-actions"), "unexpected path %q", gotPath)
	su.JSONEq(`{"invites_disabled_until":null,"dms_disabled_until":null}`, gotBody)
}

// TestModifyGuildIncidentActions_SendsTimestamps verifies non-nil pause times
// serialise as RFC3339 timestamps while the unset field stays null.
func (su *restSuite) TestModifyGuildIncidentActions_SendsTimestamps() {
	var gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}))
	defer ts.Close()

	rc, err := NewRestClient("test-token",
		options.WithBaseURL(ts.URL),
		options.WithRateLimiting(options.RateLimiterOptions{Disabled: true}),
	)
	su.Require().NoError(err)

	until := time.Date(2026, 6, 5, 12, 0, 0, 0, time.UTC)
	_, err = rc.ModifyGuildIncidentActions(context.Background(), discord.Snowflake(175928847299117063),
		&ModifyGuildIncidentActionsOptions{InvitesDisabledUntil: &until})
	su.Require().NoError(err)
	su.JSONEq(`{"invites_disabled_until":"2026-06-05T12:00:00Z","dms_disabled_until":null}`, gotBody)
}

// TestGetChannelPins_UsesNewEndpoint verifies GetChannelPins hits the new
// /messages/pins path, forwards limit, and parses the {items, has_more} envelope
// including each pin's pinned_at and message.
func (su *restSuite) TestGetChannelPins_UsesNewEndpoint() {
	var gotPath, gotQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotQuery = r.URL.Path, r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"items":[{"pinned_at":"2026-06-05T12:00:00Z","message":{"id":"175928847299117042","channel_id":"175928847299117063","content":"hi"}}],"has_more":false}`)
	}))
	defer ts.Close()

	rc, err := NewRestClient("test-token",
		options.WithBaseURL(ts.URL),
		options.WithRateLimiting(options.RateLimiterOptions{Disabled: true}),
	)
	su.Require().NoError(err)

	const chanID = discord.Snowflake(175928847299117063)
	limit := 10
	pins, err := rc.GetChannelPins(context.Background(), chanID, GetChannelPinsParams{Limit: &limit})
	su.Require().NoError(err)
	su.Require().Len(pins.Items, 1)
	su.False(pins.HasMore)
	su.Require().NotNil(pins.Items[0].Message)
	su.Equal(discord.Snowflake(175928847299117042), pins.Items[0].Message.ID)
	su.Equal(2026, pins.Items[0].PinnedAt.Year())
	su.True(strings.HasSuffix(gotPath, "/channels/"+chanID.String()+"/messages/pins"), "unexpected path %q", gotPath)
	su.Contains(gotQuery, "limit=10")
}

// TestListPinnedMessages_PaginatesAllPages verifies ListPinnedMessages walks
// every page of the new paginated endpoint, using the previous page's last
// pinned_at as the before cursor, and aggregates the messages in order.
func (su *restSuite) TestListPinnedMessages_PaginatesAllPages() {
	var befores []string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		before := r.URL.Query().Get("before")
		befores = append(befores, before)
		w.Header().Set("Content-Type", "application/json")
		if before == "" {
			_, _ = io.WriteString(w, `{"items":[{"pinned_at":"2026-06-05T12:00:00Z","message":{"id":"175928847299117001","channel_id":"175928847299117063","content":"a"}}],"has_more":true}`)
			return
		}
		_, _ = io.WriteString(w, `{"items":[{"pinned_at":"2026-06-01T09:00:00Z","message":{"id":"175928847299117002","channel_id":"175928847299117063","content":"b"}}],"has_more":false}`)
	}))
	defer ts.Close()

	rc, err := NewRestClient("test-token",
		options.WithBaseURL(ts.URL),
		options.WithRateLimiting(options.RateLimiterOptions{Disabled: true}),
	)
	su.Require().NoError(err)

	msgs, err := rc.ListPinnedMessages(context.Background(), discord.Snowflake(175928847299117063))
	su.Require().NoError(err)
	su.Require().Len(msgs, 2)
	su.Equal(discord.Snowflake(175928847299117001), msgs[0].ID)
	su.Equal(discord.Snowflake(175928847299117002), msgs[1].ID)

	// Two requests: first with no cursor, second with before = page 1's last pinned_at.
	su.Require().Len(befores, 2)
	su.Empty(befores[0])
	su.Equal("2026-06-05T12:00:00Z", befores[1])
}
