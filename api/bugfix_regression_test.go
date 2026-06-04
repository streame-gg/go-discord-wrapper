package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

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
