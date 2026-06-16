package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// capturedRequest records what the mock server received so a test can assert on
// the method, path, query, headers and body the RestClient produced.
type capturedRequest struct {
	Method    string
	Path      string
	RawQuery  string
	Auth      string
	Reason    string
	Body      []byte
	CallCount int32
}

// endpointSuite spins up a single httptest.Server shared across the suite. Each
// test configures the canned response (status + body) the server returns and
// then inspects suite.last to verify the request the client sent.
//
// This exercises the real generateRequest → rate-limit/retry → doRequest
// pipeline against a mock Discord backend, so it covers URL construction, verb
// selection, auth/audit headers and response decoding without a live API.
type endpointSuite struct {
	suite.Suite

	server *httptest.Server
	client *RestClient

	respStatus int
	respBody   string
	last       capturedRequest
	calls      int32
}

func TestEndpointSuite(t *testing.T) {
	suite.Run(t, new(endpointSuite))
}

func (s *endpointSuite) SetupTest() {
	s.respStatus = http.StatusOK
	s.respBody = "{}"
	s.last = capturedRequest{}
	atomic.StoreInt32(&s.calls, 0)

	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		s.last = capturedRequest{
			Method:    r.Method,
			Path:      r.URL.Path,
			RawQuery:  r.URL.RawQuery,
			Auth:      r.Header.Get("Authorization"),
			Reason:    r.Header.Get("X-Audit-Log-Reason"),
			Body:      body,
			CallCount: atomic.AddInt32(&s.calls, 1),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(s.respStatus)
		_, _ = io.WriteString(w, s.respBody)
	}))

	client, err := NewRestClient("test-token",
		options.WithBaseURL(s.server.URL),
		options.WithRateLimiting(options.RateLimiterOptions{Disabled: true}),
		options.WithRetry(options.RetryOptions{MaxRetries: 0}),
	)
	s.Require().NoError(err)
	s.client = client
}

func (s *endpointSuite) TearDownTest() {
	s.server.Close()
}

// respondWith sets the canned response for the next request.
func (s *endpointSuite) respondWith(status int, body string) {
	s.respStatus = status
	s.respBody = body
}

const (
	testAppID   = discord.Snowflake(900000000000000001)
	testGuildID = discord.Snowflake(800000000000000002)
	testRoleID  = discord.Snowflake(700000000000000003)
	testEmojiID = discord.Snowflake(600000000000000004)
	testChanID  = discord.Snowflake(500000000000000005)
	testCmdID   = discord.Snowflake(400000000000000006)
)

// apiPath prefixes an endpoint path with the API version segment the client
// prepends to every request (e.g. "/v10"), so assertions stay correct if the
// default version is bumped.
func apiPath(suffix string) string {
	return "/v" + discord.APIVersion10.ToString() + suffix
}

// ── Gateway ─────────────────────────────────────────────────────────────────

func (s *endpointSuite) TestGetGateway() {
	s.respondWith(http.StatusOK, `{"url":"wss://gateway.discord.gg"}`)

	resp, err := s.client.GetGateway(context.Background())

	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Equal("wss://gateway.discord.gg", resp.URL)
	s.Equal(http.MethodGet, s.last.Method)
	s.Equal(apiPath("/gateway"), s.last.Path)
	s.Equal("Bot test-token", s.last.Auth)
}

func (s *endpointSuite) TestGetBotGateway() {
	s.respondWith(http.StatusOK, `{"url":"wss://gateway.discord.gg","shards":4,"session_start_limit":{"max_concurrency":2,"total":1000,"remaining":999}}`)

	resp, err := s.client.GetBotGateway(context.Background())

	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Equal(4, resp.Shards)
	s.Equal(2, resp.SessionStartLimit.MaxConcurrency)
	s.Equal(999, resp.SessionStartLimit.Remaining)
	s.Equal(apiPath("/gateway/bot"), s.last.Path)
}

// ── Application commands ──────────────────────────────────────────────────────

func (s *endpointSuite) TestRegisterCommand() {
	// Discord answers command creation with 201 Created; the helper must accept it.
	s.respondWith(http.StatusCreated, `{"id":"400000000000000006","name":"ping","description":"pong"}`)

	cmd := &discord.ApplicationCommand{Name: "ping", Description: "pong"}
	got, err := s.client.RegisterCommand(context.Background(), testAppID, cmd)

	s.Require().NoError(err)
	s.Require().NotNil(got)
	s.Equal(testCmdID, got.ID)
	s.Equal("ping", got.Name)

	s.Equal(http.MethodPost, s.last.Method)
	s.Equal(apiPath("/applications/"+testAppID.String()+"/commands"), s.last.Path)

	// The marshalled command must be sent in the request body.
	var sent map[string]any
	s.Require().NoError(json.Unmarshal(s.last.Body, &sent))
	s.Equal("ping", sent["name"])
	s.Equal("pong", sent["description"])
}

func (s *endpointSuite) TestListGlobalApplicationCommands_WithLocalizations() {
	s.respondWith(http.StatusOK, `[{"id":"400000000000000006","name":"ping"}]`)

	got, err := s.client.ListGlobalApplicationCommands(context.Background(), testAppID, true)

	s.Require().NoError(err)
	s.Require().Len(got, 1)
	s.Equal("ping", got[0].Name)
	s.Equal(http.MethodGet, s.last.Method)
	s.Equal(apiPath("/applications/"+testAppID.String()+"/commands"), s.last.Path)
	s.Equal("with_localizations=true", s.last.RawQuery)
}

func (s *endpointSuite) TestListGlobalApplicationCommands_NoLocalizations() {
	s.respondWith(http.StatusOK, `[]`)

	got, err := s.client.ListGlobalApplicationCommands(context.Background(), testAppID, false)

	s.Require().NoError(err)
	s.Empty(got)
	s.Empty(s.last.RawQuery, "with_localizations must be omitted when false")
}

func (s *endpointSuite) TestDeleteGlobalApplicationCommand() {
	s.respondWith(http.StatusNoContent, "")

	err := s.client.DeleteGlobalApplicationCommand(context.Background(), testAppID, testCmdID)

	s.Require().NoError(err)
	s.Equal(http.MethodDelete, s.last.Method)
	s.Equal(apiPath("/applications/"+testAppID.String()+"/commands/"+testCmdID.String()), s.last.Path)
}

// ── Guild roles ───────────────────────────────────────────────────────────────

func (s *endpointSuite) TestGetGuildRole() {
	s.respondWith(http.StatusOK, `{"id":"700000000000000003","name":"Admins"}`)

	role, err := s.client.GetGuildRole(context.Background(), testGuildID, testRoleID)

	s.Require().NoError(err)
	s.Require().NotNil(role)
	s.Equal(testRoleID, role.ID)
	s.Equal("Admins", role.Name)
	s.Equal(http.MethodGet, s.last.Method)
	s.Equal(apiPath("/guilds/"+testGuildID.String()+"/roles/"+testRoleID.String()), s.last.Path)
}

func (s *endpointSuite) TestCreateGuildRole_SendsBodyAndReason() {
	s.respondWith(http.StatusOK, `{"id":"700000000000000003","name":"Mods","color":5}`)

	reason := "spring-cleanup"
	params := CreateGuildRoleParams{Name: discord.Some("Mods"), Color: discord.Some(5), AuditLogReason: &reason}
	role, err := s.client.CreateGuildRole(context.Background(), testGuildID, params)

	s.Require().NoError(err)
	s.Require().NotNil(role)
	s.Equal("Mods", role.Name)

	s.Equal(http.MethodPost, s.last.Method)
	s.Equal(apiPath("/guilds/"+testGuildID.String()+"/roles"), s.last.Path)
	s.Equal("spring-cleanup", s.last.Reason, "audit-log reason must be forwarded as a header")

	var sent map[string]any
	s.Require().NoError(json.Unmarshal(s.last.Body, &sent))
	s.Equal("Mods", sent["name"])
	s.EqualValues(5, sent["color"])
}

func (s *endpointSuite) TestCreateGuildRole_NoReasonWhenOptsNil() {
	s.respondWith(http.StatusOK, `{"id":"700000000000000003","name":"NoReason"}`)

	_, err := s.client.CreateGuildRole(context.Background(), testGuildID, CreateGuildRoleParams{})

	s.Require().NoError(err)
	s.Empty(s.last.Reason, "no audit-log header expected when opts is nil")
}

func (s *endpointSuite) TestDeleteGuildRole() {
	s.respondWith(http.StatusNoContent, "")

	err := s.client.DeleteGuildRole(context.Background(), testGuildID, testRoleID, &DeleteGuildRoleOptions{Reason: "obsolete"})

	s.Require().NoError(err)
	s.Equal(http.MethodDelete, s.last.Method)
	s.Equal(apiPath("/guilds/"+testGuildID.String()+"/roles/"+testRoleID.String()), s.last.Path)
	s.Equal("obsolete", s.last.Reason)
}

// ── Channels ──────────────────────────────────────────────────────────────────

func (s *endpointSuite) TestTriggerTypingIndicator() {
	s.respondWith(http.StatusNoContent, "")

	err := s.client.TriggerTypingIndicator(context.Background(), testChanID)

	s.Require().NoError(err)
	s.Equal(http.MethodPost, s.last.Method)
	s.Equal(apiPath("/channels/"+testChanID.String()+"/typing"), s.last.Path)
	s.Empty(s.last.Body, "typing indicator carries no request body")
}

func (s *endpointSuite) TestDeleteChannel() {
	s.respondWith(http.StatusOK, `{"id":"500000000000000005","name":"old-channel","type":0}`)

	ch, err := s.client.DeleteChannel(context.Background(), testChanID, &DeleteChannelOptions{Reason: "cleanup"})

	s.Require().NoError(err)
	s.Require().NotNil(ch)
	s.Equal(testChanID, ch.ID)
	s.Equal(http.MethodDelete, s.last.Method)
	s.Equal(apiPath("/channels/"+testChanID.String()), s.last.Path)
	s.Equal("cleanup", s.last.Reason)
}

// ── Emojis ────────────────────────────────────────────────────────────────────

func (s *endpointSuite) TestListGuildEmojis() {
	s.respondWith(http.StatusOK, `[{"id":"600000000000000004","name":"party"},{"id":"600000000000000005","name":"wave"}]`)

	emojis, err := s.client.ListGuildEmojis(context.Background(), testGuildID)

	s.Require().NoError(err)
	s.Require().Len(emojis, 2)
	s.Equal("party", emojis[0].Name)
	s.Equal("wave", emojis[1].Name)
	s.Equal(http.MethodGet, s.last.Method)
	s.Equal(apiPath("/guilds/"+testGuildID.String()+"/emojis"), s.last.Path)
}

func (s *endpointSuite) TestDeleteGuildEmoji() {
	s.respondWith(http.StatusNoContent, "")

	err := s.client.DeleteGuildEmoji(context.Background(), testGuildID, testEmojiID, &DeleteGuildEmojiOptions{Reason: "spam-emoji"})

	s.Require().NoError(err)
	s.Equal(http.MethodDelete, s.last.Method)
	s.Equal(apiPath("/guilds/"+testGuildID.String()+"/emojis/"+testEmojiID.String()), s.last.Path)
	s.Equal("spam-emoji", s.last.Reason)
}

// ── Cross-cutting behaviour ─────────────────────────────────────────────────────

// TestErrorStatusMapping verifies an endpoint returning a non-2xx status surfaces
// both the sentinel error and a decoded *Error with the HTTP status preserved.
func (s *endpointSuite) TestErrorStatusMapping() {
	s.respondWith(http.StatusForbidden, `{"code":50013,"message":"Missing Permissions"}`)

	role, err := s.client.GetGuildRole(context.Background(), testGuildID, testRoleID)

	s.Require().Error(err)
	s.Nil(role)
	s.ErrorIs(err, ErrForbidden)

	var apiErr *Error
	s.Require().ErrorAs(err, &apiErr)
	s.Equal(http.StatusForbidden, apiErr.HTTPStatus)
	s.Equal(50013, int(apiErr.Code))
	s.Equal("Missing Permissions", apiErr.Message)
}

// TestInvalidSnowflakeShortCircuits verifies that an obviously invalid ID is
// rejected before any HTTP request is made (the mock server must never be hit).
func (s *endpointSuite) TestInvalidSnowflakeShortCircuits() {
	_, err := s.client.GetGuildRole(context.Background(), discord.Snowflake(0), testRoleID)

	s.Require().Error(err, "zero guild ID must be rejected by validation")
	s.Equal(int32(0), atomic.LoadInt32(&s.calls), "no HTTP request should be sent for an invalid snowflake")
}

// TestEveryRequestCarriesBotAuth confirms the bot token is attached to a
// representative spread of verbs (GET / POST / DELETE).
func (s *endpointSuite) TestEveryRequestCarriesBotAuth() {
	s.respondWith(http.StatusOK, `{"url":"wss://x"}`)
	_, _ = s.client.GetGateway(context.Background())
	s.Equal("Bot test-token", s.last.Auth)

	s.respondWith(http.StatusNoContent, "")
	_ = s.client.TriggerTypingIndicator(context.Background(), testChanID)
	s.Equal("Bot test-token", s.last.Auth)

	s.respondWith(http.StatusNoContent, "")
	_ = s.client.DeleteGuildRole(context.Background(), testGuildID, testRoleID, nil)
	s.Equal("Bot test-token", s.last.Auth)
}
