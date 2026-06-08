package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/options"
)

func (s *endpointSuite) TestWebhookRestEndpoints() {
	ctx := context.Background()
	wh := testWebhookID.String()
	waitTrue := true
	s.runEndpoints([]epCase{
		{name: "CreateWebhook", method: "POST", path: "/channels/" + testChanID.String() + "/webhooks", call: func() error {
			return drop(s.client.CreateWebhook(ctx, testChanID, CreateWebhookParams{}))
		}},
		{name: "ListChannelWebhooks", method: "GET", path: "/channels/" + testChanID.String() + "/webhooks", body: "[]", call: func() error {
			return drop(s.client.ListChannelWebhooks(ctx, testChanID))
		}},
		{name: "ListGuildWebhooks", method: "GET", path: "/guilds/" + testGuildID.String() + "/webhooks", body: "[]", call: func() error {
			return drop(s.client.ListGuildWebhooks(ctx, testGuildID))
		}},
		{name: "GetWebhook", method: "GET", path: "/webhooks/" + wh, call: func() error {
			return drop(s.client.GetWebhook(ctx, testWebhookID))
		}},
		{name: "GetWebhookWithToken", method: "GET", path: "/webhooks/" + wh + "/tok", call: func() error {
			return drop(s.client.GetWebhookWithToken(ctx, testWebhookID, "tok"))
		}},
		{name: "ModifyWebhook", method: "PATCH", path: "/webhooks/" + wh, call: func() error {
			return drop(s.client.ModifyWebhook(ctx, testWebhookID, ModifyWebhookParams{}))
		}},
		{name: "ModifyWebhookWithToken", method: "PATCH", path: "/webhooks/" + wh + "/tok", call: func() error {
			return drop(s.client.ModifyWebhookWithToken(ctx, testWebhookID, "tok", ModifyWebhookParams{}))
		}},
		{name: "DeleteWebhook", method: "DELETE", path: "/webhooks/" + wh, call: func() error {
			return s.client.DeleteWebhook(ctx, testWebhookID, nil)
		}},
		{name: "DeleteWebhookWithToken", method: "DELETE", path: "/webhooks/" + wh + "/tok", call: func() error {
			return s.client.DeleteWebhookWithToken(ctx, testWebhookID, "tok")
		}},
		{name: "ExecuteWebhook", method: "POST", path: "/webhooks/" + wh + "/tok", call: func() error {
			return drop(s.client.ExecuteWebhook(ctx, testWebhookID, "tok", ExecuteWebhookParams{Content: "hi"}, ExecuteWebhookQueryParams{Wait: &waitTrue}))
		}},
		{name: "GetWebhookMessage", method: "GET", path: "/webhooks/" + wh + "/tok/messages/" + testMsgID.String(), call: func() error {
			return drop(s.client.GetWebhookMessage(ctx, testWebhookID, "tok", testMsgID))
		}},
		{name: "EditWebhookMessage", method: "PATCH", path: "/webhooks/" + wh + "/tok/messages/" + testMsgID.String(), call: func() error {
			return drop(s.client.EditWebhookMessage(ctx, testWebhookID, "tok", testMsgID, EditMessageParams{}))
		}},
		{name: "DeleteWebhookMessage", method: "DELETE", path: "/webhooks/" + wh + "/tok/messages/" + testMsgID.String(), call: func() error {
			return s.client.DeleteWebhookMessage(ctx, testWebhookID, "tok", testMsgID)
		}},
		// Slack/GitHub compatibility webhooks return 204 when wait=false.
		{name: "ExecuteSlackWebhook", method: "POST", path: "/webhooks/" + wh + "/tok/slack", status: http.StatusNoContent, call: func() error {
			return s.client.ExecuteSlackWebhook(ctx, testWebhookID, "tok", false, json.RawMessage(`{}`))
		}},
		{name: "ExecuteGitHubWebhook", method: "POST", path: "/webhooks/" + wh + "/tok/github", status: http.StatusNoContent, call: func() error {
			return s.client.ExecuteGitHubWebhook(ctx, testWebhookID, "tok", false, json.RawMessage(`{}`))
		}},
	})
}

// webhookClientSuite covers the high-level WebhookClient wrapper. It uses its
// own mock server (the WebhookClient owns an pkg RestClient) and verifies
// the convenience methods that the existing webhook_client_test.go leaves
// uncovered (Edit, Delete, SendSlack, SendGitHub, FetchMessage, EditMessage,
// the Reset* default helpers, and NewWebhookClientFromURL).
type webhookClientSuite struct {
	suite.Suite

	server     *httptest.Server
	wc         *WebhookClient
	lastMethod string
	lastPath   string
	status     int
	body       string
}

func TestWebhookClientSuite(t *testing.T) {
	suite.Run(t, new(webhookClientSuite))
}

func (s *webhookClientSuite) SetupTest() {
	s.status = http.StatusOK
	s.body = "{}"
	s.lastMethod, s.lastPath = "", ""

	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.lastMethod = r.Method
		s.lastPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(s.status)
		_, _ = io.WriteString(w, s.body)
	}))

	wc, err := NewWebhookClient(testWebhookID, "tok",
		options.WithBaseURL(s.server.URL),
		options.WithRateLimiting(options.RateLimiterOptions{Disabled: true}),
	)
	s.Require().NoError(err)
	s.wc = wc
}

func (s *webhookClientSuite) TearDownTest() {
	s.wc.Close()
	s.server.Close()
}

func (s *webhookClientSuite) TestEdit() {
	s.body = `{"id":"210000000000000010","type":1}`
	wh, err := s.wc.Edit(context.Background(), ModifyWebhookParams{})
	s.Require().NoError(err)
	s.Equal(testWebhookID, wh.ID)
	s.Equal(http.MethodPatch, s.lastMethod)
	s.Equal(apiPath("/webhooks/"+testWebhookID.String()+"/tok"), s.lastPath)
}

func (s *webhookClientSuite) TestDelete() {
	s.status = http.StatusNoContent
	s.body = ""
	err := s.wc.Delete(context.Background())
	s.Require().NoError(err)
	s.Equal(http.MethodDelete, s.lastMethod)
	s.Equal(apiPath("/webhooks/"+testWebhookID.String()+"/tok"), s.lastPath)
}

func (s *webhookClientSuite) TestFetchMessage() {
	s.body = `{"id":"120000000000000002","content":"hi"}`
	msg, err := s.wc.FetchMessage(context.Background(), testMsgID)
	s.Require().NoError(err)
	s.Equal(testMsgID, msg.ID)
	s.Equal(http.MethodGet, s.lastMethod)
	s.Equal(apiPath("/webhooks/"+testWebhookID.String()+"/tok/messages/"+testMsgID.String()), s.lastPath)
}

func (s *webhookClientSuite) TestEditMessage() {
	s.body = `{"id":"120000000000000002","content":"edited"}`
	msg, err := s.wc.EditMessage(context.Background(), testMsgID, EditMessageParams{})
	s.Require().NoError(err)
	s.Equal("edited", msg.Content)
	s.Equal(http.MethodPatch, s.lastMethod)
	s.Equal(apiPath("/webhooks/"+testWebhookID.String()+"/tok/messages/"+testMsgID.String()), s.lastPath)
}

func (s *webhookClientSuite) TestSendSlack() {
	s.status = http.StatusNoContent
	s.body = ""
	err := s.wc.SendSlack(context.Background(), false, json.RawMessage(`{}`))
	s.Require().NoError(err)
	s.Equal(http.MethodPost, s.lastMethod)
	s.Equal(apiPath("/webhooks/"+testWebhookID.String()+"/tok/slack"), s.lastPath)
}

func (s *webhookClientSuite) TestSendGitHub() {
	s.status = http.StatusNoContent
	s.body = ""
	err := s.wc.SendGitHub(context.Background(), false, json.RawMessage(`{}`))
	s.Require().NoError(err)
	s.Equal(http.MethodPost, s.lastMethod)
	s.Equal(apiPath("/webhooks/"+testWebhookID.String()+"/tok/github"), s.lastPath)
}

// TestResetHelpers exercises the fluent default-setters and their resets. These
// are pure state mutations on the client, so we assert the documented chaining
// contract (each returns the receiver) rather than an HTTP effect.
func (s *webhookClientSuite) TestResetHelpers() {
	s.Same(s.wc, s.wc.SetUsername("bot"))
	s.Same(s.wc, s.wc.SetAvatarURL("https://x/y.png"))
	s.Same(s.wc, s.wc.ResetUsername())
	s.Same(s.wc, s.wc.ResetAvatarURL())
	s.Same(s.wc, s.wc.ResetDefaults())
}

func (s *webhookClientSuite) TestNewWebhookClientFromURL() {
	url := "https://discord.com/api/webhooks/" + testWebhookID.String() + "/tok"
	wc, err := NewWebhookClientFromURL(url)
	s.Require().NoError(err)
	defer wc.Close()
	s.Equal(testWebhookID, wc.ID())
}
