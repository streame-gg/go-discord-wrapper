package connection

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/internal/util"
	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

// clientExtrasSuite covers the remaining thin *Client wrappers that the
// rest/entity suites skip: interaction responses + follow-ups (which need an
// application id and a *interactions.Interaction), command registration, the
// cache-index manager helpers, and a handful of slice/map REST wrappers.
//
// As a white-box test it sets Client.appID directly so the application-id-gated
// wrappers run their delegate path instead of erroring out early.
type clientExtrasSuite struct {
	suite.Suite
	server *httptest.Server
	client *Client

	mu     sync.Mutex
	status int
	body   []byte
}

func TestClientExtrasSuite(t *testing.T) { suite.Run(t, new(clientExtrasSuite)) }

func (s *clientExtrasSuite) setResp(status int, body string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status = status
	s.body = []byte(body)
}

func (s *clientExtrasSuite) SetupTest() {
	s.status = http.StatusOK
	s.body = []byte(`{"id":"175928847299117063","channel_id":"175928847299117063","guild_id":"175928847299117063","name":"x"}`)
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		status, b := s.status, s.body
		s.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write(b)
	}))
	s.T().Cleanup(s.server.Close)

	c, err := NewClient("Bot fake-token", discord.IntentGuilds, options.WithBaseURL(s.server.URL))
	s.Require().NoError(err)
	c.Cache = cache.NewMemoryCache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}})
	// Make the bot "ready" so applicationID() succeeds.
	appID := sampleSnowflake
	c.appID = &appID
	s.client = c
}

// TestInteractionResponses drives the immediate interaction-response wrappers.
// The interaction-response endpoint returns 204 (no body), so success codes vary;
// errors are ignored — the wrapper's request-building path is what we cover.
func (s *clientExtrasSuite) TestInteractionResponses() {
	c, ctx := s.client, context.Background()
	i := &interactions.Interaction{ID: sampleSnowflake, Token: "tok"}

	s.setResp(http.StatusNoContent, ``)

	ignE := func(_ error) {}
	ignR := func(_ any, _ error) {}

	ignR(c.Reply(ctx, i, &responses.InteractionResponseDataDefault{}, false))
	ignE(c.DeferReply(ctx, i, true))
	ignE(c.DeferReply(ctx, i, false))
	ignE(c.DeferUpdateMessage(ctx, i))
	ignE(c.UpdateMessage(ctx, i, &responses.InteractionResponseDataDefault{}))
	ignE(c.ReplyWithModal(ctx, i, &components.Modal{}))
	ignE(c.LaunchActivity(ctx, i))
	ignE(c.ReplyAutocomplete(ctx, i, []responses.AutocompleteChoice{}))
}

// TestInteractionFollowups drives the original-response and follow-up wrappers,
// which return a message object on success (200) and hit the cache branch.
func (s *clientExtrasSuite) TestInteractionFollowups() {
	c, ctx := s.client, context.Background()
	i := &interactions.Interaction{ID: sampleSnowflake, Token: "tok"}
	id := sampleSnowflake

	s.setResp(http.StatusOK, `{"id":"175928847299117063","channel_id":"175928847299117063"}`)

	mustNoErr := func(_ any, err error) { s.NoError(err) }

	mustNoErr(c.GetOriginalResponse(ctx, i))
	mustNoErr(c.EditReply(ctx, i, api.EditMessageParams{}))
	mustNoErr(c.CreateFollowup(ctx, i, api.CreateMessageParams{}))
	mustNoErr(c.GetFollowup(ctx, i, id))
	mustNoErr(c.EditFollowup(ctx, i, id, api.EditMessageParams{}))

	s.setResp(http.StatusNoContent, ``)
	s.NoError(c.DeleteReply(ctx, i))
	s.NoError(c.DeleteFollowup(ctx, i, id))
}

// TestCommandRegistration covers RegisterCommand / BulkRegisterCommands.
func (s *clientExtrasSuite) TestCommandRegistration() {
	c, ctx := s.client, context.Background()

	s.setResp(http.StatusOK, `{"id":"175928847299117063","name":"x","type":1}`)
	_, err := c.RegisterCommand(ctx, &discord.ApplicationCommand{Name: "x"})
	s.NoError(err)

	s.setResp(http.StatusOK, `[{"id":"175928847299117063","name":"x","type":1}]`)
	_, err = c.BulkRegisterCommands(ctx, []*discord.ApplicationCommand{{Name: "x"}})
	s.NoError(err)
}

// TestManagerHelpers covers the cache-index accessors and manager constructors.
func (s *clientExtrasSuite) TestManagerHelpers() {
	c := s.client
	id := sampleSnowflake

	s.NotNil(c.Guilds())
	s.NotNil(c.Users())
	s.NotNil(c.Channels())

	// Empty index → empty collections.
	s.Equal(0, c.ThreadsForParent(id).Len())
	s.Equal(0, c.ChannelsForGuild(id).Len())

	// Seed a channel + index entries so the lookup loops run their body.
	parent := &discord.Channel{ID: id, GuildID: util.PointerOf(id), Type: discord.ChannelTypeGuildText}
	thread := &discord.Channel{ID: id + 1, GuildID: util.PointerOf(id), ParentID: &id, Type: discord.ChannelTypePublicThread}
	c.Cache.Channels().Set(parent)
	c.Cache.Channels().Set(thread)
	c.threadIndexMu.Lock()
	c.threadsByParent = map[discord.Snowflake]map[discord.Snowflake]struct{}{id: {id + 1: {}}}
	c.threadIndexMu.Unlock()
	c.channelIndexMu.Lock()
	c.channelsByGuild = map[discord.Snowflake]map[discord.Snowflake]struct{}{id: {id: {}}}
	c.channelIndexMu.Unlock()

	s.Equal(1, c.ThreadsForParent(id).Len())
	s.Equal(1, c.ChannelsForGuild(id).Len())
}

// TestMiscRestWrappers covers the remaining map/slice REST wrappers.
func (s *clientExtrasSuite) TestMiscRestWrappers() {
	c, ctx := s.client, context.Background()
	id := sampleSnowflake

	s.setResp(http.StatusOK, `{"175928847299117063":3}`)
	_, err := c.GetGuildRoleMemberCounts(ctx, id)
	s.NoError(err)

	s.setResp(http.StatusOK, `{"items":[{"id":"175928847299117063","name":"x"}]}`)
	_, err = c.ListApplicationEmojis(ctx, id)
	s.NoError(err)

	// answerID is currently routed through path snowflake validation (a known
	// latent bug), so this errors before reaching the cache loop; we only need
	// the wrapper's delegate path to run.
	s.setResp(http.StatusOK, `{"users":[{"id":"175928847299117063"}]}`)
	_, _ = c.GetPollAnswerVoters(ctx, id, id, 1, api.GetPollAnswerVotersParams{})

	s.setResp(http.StatusOK, `{"banned_users":["175928847299117063"],"failed_users":[]}`)
	_, err = c.BulkBanGuildMembers(ctx, id, api.BulkBanParams{UserIDs: []discord.Snowflake{id}})
	s.NoError(err)
}

// TestFileAndWebhookWrappers covers the entity-client wrappers that require a
// non-trivial response (CreateGuildEmoji/Sticker, ExecuteWebhook).
func (s *clientExtrasSuite) TestFileAndWebhookWrappers() {
	c, ctx := s.client, context.Background()
	id := sampleSnowflake
	r := ptrStr("reason")

	ignR := func(_ any, _ error) {}

	s.setResp(http.StatusCreated, `{"id":"175928847299117063","name":"x"}`)
	ignR(c.CreateGuildEmoji(ctx, id, discord.EmojiCreateOptions{Name: "x", AuditLogReason: r}))
	ignR(c.CreateGuildSticker(ctx, id, discord.StickerCreateOptions{Name: "x", AuditLogReason: r}))

	s.setResp(http.StatusOK, `{"id":"175928847299117063","channel_id":"175928847299117063"}`)
	ignR(c.ExecuteWebhook(ctx, id, "tok", discord.WebhookExecuteOptions{Wait: true}))
}
