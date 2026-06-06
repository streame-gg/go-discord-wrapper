package api

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (su *apiWebhookClientSuite) TestParseWebhookURL() {
	t := su.T()
	t.Run("standard url", func(t *testing.T) {
		id, token, err := parseWebhookURL("https://discord.com/api/webhooks/123456789012345678/abc-token")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if id != 123456789012345678 {
			t.Errorf("id: got %d", id)
		}
		if token != "abc-token" {
			t.Errorf("token: got %q", token)
		}
	})

	t.Run("versioned url with trailing path ignored", func(t *testing.T) {
		id, token, err := parseWebhookURL("https://discord.com/api/v10/webhooks/123456789012345678/tok/extra")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if id != 123456789012345678 || token != "tok" {
			t.Errorf("got id=%d token=%q", id, token)
		}
	})

	t.Run("missing token", func(t *testing.T) {
		if _, _, err := parseWebhookURL("https://discord.com/api/webhooks/123456789012345678"); err == nil {
			t.Error("expected error for missing token")
		}
	})

	t.Run("not a webhook url", func(t *testing.T) {
		if _, _, err := parseWebhookURL("https://discord.com/api/channels/123/messages"); err == nil {
			t.Error("expected error for non-webhook url")
		}
	})
}

func (su *apiWebhookClientSuite) TestNewWebhookClientValidation() {
	t := su.T()
	if _, err := NewWebhookClient(123456789012345678, "   "); err == nil {
		t.Error("expected error for empty token")
	}
	if _, err := NewWebhookClient(0, "tok"); err == nil {
		t.Error("expected error for invalid id")
	}
	wc, err := NewWebhookClient(123456789012345678, "tok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer wc.Close()
	if wc.ID() != 123456789012345678 {
		t.Errorf("ID(): got %d", wc.ID())
	}
}

// newTestWebhookClient builds a WebhookClient whose requests are routed to srv.
func newTestWebhookClient(t *testing.T, srv *httptest.Server) *WebhookClient {
	t.Helper()
	wc, err := NewWebhookClient(
		123456789012345678, "tok",
		options.WithBaseURL(srv.URL),
		options.WithRateLimiting(options.RateLimiterOptions{Disabled: true}),
	)
	if err != nil {
		t.Fatalf("NewWebhookClient: %v", err)
	}
	t.Cleanup(wc.Close)
	return wc
}

func (su *apiWebhookClientSuite) TestWebhookClientSendAndManage() {
	t := su.T()
	var gotMethod, gotPath, gotWait string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotWait = r.URL.Query().Get("wait")
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"555","content":"hi"}`))
		case r.Method == http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"123456789012345678","type":1}`))
		case r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer srv.Close()

	wc := newTestWebhookClient(t, srv)
	ctx := context.Background()
	wait := true

	t.Run("send with wait returns message", func(t *testing.T) {
		msg, err := wc.Send(ctx,
			ExecuteWebhookParams{Content: "hi"},
			ExecuteWebhookQueryParams{Wait: &wait},
		)
		if err != nil {
			t.Fatalf("Send: %v", err)
		}
		if msg == nil || msg.ID != 555 {
			t.Fatalf("Send returned unexpected message: %+v", msg)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method: got %s", gotMethod)
		}
		if gotPath != "/v10/webhooks/123456789012345678/tok" {
			t.Errorf("path: got %s", gotPath)
		}
		if gotWait != "true" {
			t.Errorf("wait query: got %q", gotWait)
		}
	})

	t.Run("fetch", func(t *testing.T) {
		wh, err := wc.Fetch(ctx)
		if err != nil {
			t.Fatalf("Fetch: %v", err)
		}
		if wh.ID != 123456789012345678 {
			t.Errorf("fetched webhook id: got %d", wh.ID)
		}
		if gotPath != "/v10/webhooks/123456789012345678/tok" {
			t.Errorf("path: got %s", gotPath)
		}
	})

	t.Run("delete message", func(t *testing.T) {
		if err := wc.DeleteMessage(ctx, 987654321098765432); err != nil {
			t.Fatalf("DeleteMessage: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("method: got %s", gotMethod)
		}
		if gotPath != "/v10/webhooks/123456789012345678/tok/messages/987654321098765432" {
			t.Errorf("path: got %s", gotPath)
		}
	})
}

func (su *apiWebhookClientSuite) TestWebhookClientSendDefaults() {
	t := su.T()
	var gotBody struct {
		Content   string  `json:"content"`
		Username  *string `json:"username"`
		AvatarURL *string `json:"avatar_url"`
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody.Content, gotBody.Username, gotBody.AvatarURL = "", nil, nil // reset; omitted fields aren't overwritten by Unmarshal
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	wc := newTestWebhookClient(t, srv)
	ctx := context.Background()

	// Chaining returns the client; defaults applied when params omit the fields.
	wc.SetUsername("Logger").SetAvatarURL("https://example.com/a.png")
	if _, err := wc.Send(ctx, ExecuteWebhookParams{Content: "hi"}, ExecuteWebhookQueryParams{}); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if gotBody.Username == nil || *gotBody.Username != "Logger" {
		t.Errorf("default username not applied: %+v", gotBody.Username)
	}
	if gotBody.AvatarURL == nil || *gotBody.AvatarURL != "https://example.com/a.png" {
		t.Errorf("default avatar not applied: %+v", gotBody.AvatarURL)
	}

	// Per-message values take precedence over defaults.
	override := "Override"
	if _, err := wc.Send(ctx, ExecuteWebhookParams{Content: "hi", Username: &override}, ExecuteWebhookQueryParams{}); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if gotBody.Username == nil || *gotBody.Username != "Override" {
		t.Errorf("per-message username should win: %+v", gotBody.Username)
	}

	// Reset clears defaults.
	wc.ResetDefaults()
	if _, err := wc.Send(ctx, ExecuteWebhookParams{Content: "hi"}, ExecuteWebhookQueryParams{}); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if gotBody.Username != nil || gotBody.AvatarURL != nil {
		t.Errorf("defaults should be cleared: username=%v avatar=%v", gotBody.Username, gotBody.AvatarURL)
	}
}

// Ensure the WebhookClient needs no bot token: construction succeeds with only
// an ID and token, unlike NewRestClient.
func (su *apiWebhookClientSuite) TestWebhookClientNeedsNoBotToken() {
	t := su.T()
	if _, err := NewRestClient(""); err == nil {
		t.Error("NewRestClient should require a token")
	}
	wc, err := NewWebhookClient(discord.Snowflake(123456789012345678), "tok")
	if err != nil {
		t.Fatalf("NewWebhookClient should not require a bot token: %v", err)
	}
	wc.Close()
}

type apiWebhookClientSuite struct{ suite.Suite }

func TestApiWebhookClientSuite(t *testing.T) { suite.Run(t, new(apiWebhookClientSuite)) }
