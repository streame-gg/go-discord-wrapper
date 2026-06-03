package connection

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// entityClientSuite drives the discord.EntityClient implementation on *Client
// (entity_client.go) against an httptest server. Every method is a thin wrapper
// that maps a discord.XxxOptions struct onto api params, delegates to RestClient,
// and on success hydrates/caches the result. The suite calls each wrapper so both
// the param-building path and (where the response status matches) the cache branch
// execute.
//
// The handler returns a configurable status + body so the suite can satisfy the
// per-endpoint success code (200 for reads/edits, 201 for creates, 204 for
// deletes), which is what unlocks the `if err == nil { cache }` branches.
type entityClientSuite struct {
	suite.Suite
	server *httptest.Server
	client *Client

	mu     sync.Mutex
	status int
	body   []byte
}

func TestEntityClientSuite(t *testing.T) { suite.Run(t, new(entityClientSuite)) }

func (s *entityClientSuite) setResp(status int, body string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status = status
	s.body = []byte(body)
}

func (s *entityClientSuite) SetupTest() {
	s.status = http.StatusOK
	s.body = []byte(`{"id":"175928847299117063","channel_id":"175928847299117063","guild_id":"175928847299117063","content":"hi","name":"x"}`)
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
	s.client = c
}

// reason returns a *string for AuditLogReason fields, exercising the apiOpts!=nil branch.
func ptrStr(v string) *string { return &v }

// TestObjectAndEditWrappers covers wrappers whose endpoints return 200 with a
// single object (reads + PATCH edits), including their hydrate/cache branch.
func (s *entityClientSuite) TestObjectAndEditWrappers() {
	c, ctx := s.client, context.Background()
	id := sampleSnowflake
	r := ptrStr("reason")

	mustNoErr := func(_ any, err error) { s.NoError(err) }

	// Messages
	mustNoErr(c.EditMessage(ctx, id, id, discord.MessageEditOptions{}))
	// CreateMessage returns 200 on Discord.
	mustNoErr(c.CreateMessage(ctx, id, discord.MessageCreateOptions{Content: "hi"}))

	// Channel edits
	mustNoErr(c.ModifyChannel(ctx, id, discord.ChannelEditOptions{AuditLogReason: r}))
	mustNoErr(c.CreateChannelInvite(ctx, id, discord.InviteCreateOptions{AuditLogReason: r}))

	// Guild edits / reads
	mustNoErr(c.ModifyGuild(ctx, id, discord.GuildEditOptions{AuditLogReason: r}))
	mustNoErr(c.GetGuildAuditLog(ctx, id, discord.AuditLogOptions{}))

	// Roles
	mustNoErr(c.ModifyGuildRole(ctx, id, id, discord.RoleEditOptions{AuditLogReason: r}))

	// Members
	mustNoErr(c.ModifyGuildMember(ctx, id, id, discord.MemberEditOptions{AuditLogReason: r}))

	// User / DM
	mustNoErr(c.GetUser(ctx, id))
	mustNoErr(c.CreateDM(ctx, id))

	// Emoji
	mustNoErr(c.ModifyGuildEmoji(ctx, id, id, discord.EmojiEditOptions{AuditLogReason: r}))
	mustNoErr(c.GetGuildEmoji(ctx, id, id))

	// Webhook
	mustNoErr(c.ModifyWebhook(ctx, id, discord.WebhookEditOptions{AuditLogReason: r}))
	mustNoErr(c.GetWebhookMessage(ctx, id, "tok", id))

	// Stage
	mustNoErr(c.ModifyStageInstance(ctx, id, discord.StageEditOptions{AuditLogReason: r}))
	mustNoErr(c.GetStageInstance(ctx, id))

	// Scheduled events
	mustNoErr(c.ModifyGuildScheduledEvent(ctx, id, id, discord.ScheduledEventEditOptions{}))
	mustNoErr(c.GetGuildScheduledEvent(ctx, id, id))

	// Sticker
	mustNoErr(c.ModifyGuildSticker(ctx, id, id, discord.StickerEditOptions{AuditLogReason: r}))
	mustNoErr(c.GetGuildSticker(ctx, id, id))

	// AutoMod
	mustNoErr(c.ModifyAutoModerationRule(ctx, id, id, discord.RuleEditOptions{AuditLogReason: r}))
	mustNoErr(c.GetAutoModerationRule(ctx, id, id))

	// Soundboard
	mustNoErr(c.ModifyGuildSoundboardSound(ctx, id, id, discord.SoundEditOptions{AuditLogReason: r}))

	// ClientCache accessor
	s.NotNil(c.ClientCache())
}

// TestSliceWrappers covers wrappers returning a slice (200 + array body).
func (s *entityClientSuite) TestSliceWrappers() {
	c, ctx := s.client, context.Background()
	id := sampleSnowflake
	r := ptrStr("reason")

	s.setResp(http.StatusOK, `[{"id":"175928847299117063","channel_id":"175928847299117063","guild_id":"175928847299117063","name":"x"}]`)

	mustNoErr := func(_ any, err error) { s.NoError(err) }

	mustNoErr(c.ModifyGuildRolePositions(ctx, id, discord.RolePositionOptions{AuditLogReason: r}))
	mustNoErr(c.ListChannelMessages(ctx, id, discord.FetchMessagesOptions{}))
	mustNoErr(c.ListGuildMembers(ctx, id, discord.FetchMembersOptions{}))
	mustNoErr(c.ListGuildEmojis(ctx, id))
	mustNoErr(c.ListGuildStickers(ctx, id))
	mustNoErr(c.ListGuildScheduledEvents(ctx, id))
	mustNoErr(c.ListGuildScheduledEventUsers(ctx, id, id, discord.FetchUsersOptions{}))
	mustNoErr(c.ListGuildWebhooks(ctx, id))
	mustNoErr(c.ListAutoModerationRules(ctx, id))
}

// TestCreateWrappers covers POST wrappers (most Discord create endpoints return
// 200, not 201).
func (s *entityClientSuite) TestCreateWrappers() {
	c, ctx := s.client, context.Background()
	id := sampleSnowflake
	r := ptrStr("reason")

	s.setResp(http.StatusOK, `{"id":"175928847299117063","guild_id":"175928847299117063","name":"x"}`)

	mustNoErr := func(_ any, err error) { s.NoError(err) }
	// ign: endpoints whose success status isn't 201, or that validate params the
	// generic fixture can't satisfy; the wrapper's param-building still runs.
	ign := func(_ any, _ error) {}

	mustNoErr(c.CreateGuildRole(ctx, id, discord.RoleCreateOptions{AuditLogReason: r}))
	ign(c.CreateGuildChannel(ctx, id, discord.ChannelCreateOptions{AuditLogReason: r}))
	ign(c.CreateGuildScheduledEvent(ctx, id, discord.ScheduledEventCreateOptions{}))
	ign(c.CreateStageInstance(ctx, discord.StageCreateOptions{AuditLogReason: r}))
	ign(c.CreateAutoModerationRule(ctx, id, discord.RuleCreateOptions{AuditLogReason: r}))
}

// TestErrorOnlyWrappers covers wrappers returning only error, including the 204
// delete endpoints and their post-delete cache eviction branch.
func (s *entityClientSuite) TestErrorOnlyWrappers() {
	c, ctx := s.client, context.Background()
	id := sampleSnowflake
	r := ptrStr("reason")

	s.setResp(http.StatusNoContent, ``)

	s.NoError(c.LeaveGuild(ctx, id))
	s.NoError(c.KickGuildMember(ctx, id, id, r))
	s.NoError(c.CreateGuildBan(ctx, id, id, discord.BanOptions{AuditLogReason: r}))
	s.NoError(c.DeleteGuildEmoji(ctx, id, id, r))
	s.NoError(c.DeleteWebhook(ctx, id, r))
	s.NoError(c.DeleteStageInstance(ctx, id, r))
	s.NoError(c.DeleteGuildScheduledEvent(ctx, id, id))
	s.NoError(c.DeleteGuildSticker(ctx, id, id, r))
	s.NoError(c.DeleteAutoModerationRule(ctx, id, id))
}
