package connection

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// restWrappersSuite drives the thin Client REST wrappers against an httptest
// server, exercising their success path (delegate → hydrate → cache). The
// server returns a generic JSON body (object or array, switchable per test) that
// decodes cleanly into the wrapper return types — unknown fields are ignored and
// missing fields take zero values, so a single fixture satisfies most endpoints.
//
// Interaction-based wrappers (which need a populated *interactions.Interaction
// and application id) and gateway-backed wrappers (RequestGuildMembers,
// UpdatePresence, FetchAllGuildMembers) are exercised elsewhere / require a live
// socket and are intentionally out of scope here.
type restWrappersSuite struct {
	suite.Suite
	server *httptest.Server
	client *Client

	mu   sync.Mutex
	body []byte
}

func TestRestWrappersSuite(t *testing.T) { suite.Run(t, new(restWrappersSuite)) }

const sampleSnowflake = discord.Snowflake(175928847299117063)

func (s *restWrappersSuite) setBody(b string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.body = []byte(b)
}

func (s *restWrappersSuite) SetupTest() {
	s.body = []byte(`{"id":"175928847299117063","channel_id":"175928847299117063","content":"hi","name":"x"}`)
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		b := s.body
		s.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}))
	s.T().Cleanup(s.server.Close)

	c, err := NewClient("Bot fake-token", discord.IntentGuilds, options.WithBaseURL(s.server.URL))
	s.Require().NoError(err)
	c.Cache = cache.NewMemoryCache(cache.Options{Messages: cache.MessageOptions{MaxPerChannel: 100}})
	s.client = c
}

// TestObjectWrappers calls every single-object-returning wrapper (excluding the
// interaction-based ones) against an object response.
func (s *restWrappersSuite) TestObjectWrappers() {
	c, ctx := s.client, context.Background()
	id := sampleSnowflake

	mustNoErr := func(_ any, err error) { s.NoError(err) }
	// ign exercises the wrapper's delegation path for endpoints whose success
	// code (201/204), application-id, or file-validation requirements aren't met
	// by this generic fixture; the error is expected and ignored.
	ign := func(_ any, _ error) {}

	mustNoErr(c.GetMessage(ctx, id, id))
	mustNoErr(c.SendMessage(ctx, id, api.CreateMessageParams{}))
	mustNoErr(c.EditMessageRaw(ctx, id, id, api.EditMessageParams{}))
	mustNoErr(c.CrosspostMessage(ctx, id, id))
	mustNoErr(c.EndPoll(ctx, id, id))
	mustNoErr(c.GetChannel(ctx, id))
	mustNoErr(c.ModifyChannelRaw(ctx, id, api.ModifyChannelParams{}))
	mustNoErr(c.DeleteChannel(ctx, id, nil))
	mustNoErr(c.CreateChannelInviteRaw(ctx, id, api.CreateChannelInviteParams{}))
	mustNoErr(c.GetGuild(ctx, id))
	mustNoErr(c.GetGuildPreview(ctx, id))
	mustNoErr(c.ModifyGuildRaw(ctx, id, api.ModifyGuildParams{}))
	ign(c.CreateGuildChannelRaw(ctx, id, api.CreateGuildChannelParams{}))
	mustNoErr(c.GetGuildRole(ctx, id, id))
	mustNoErr(c.CreateGuildRoleRaw(ctx, id, api.CreateGuildRoleParams{}))
	mustNoErr(c.ModifyGuildRoleRaw(ctx, id, id, api.ModifyGuildRoleParams{}))
	mustNoErr(c.GetGuildBan(ctx, id, id))
	mustNoErr(c.GetGuildPruneCount(ctx, id, api.GetGuildPruneCountParams{}))
	mustNoErr(c.BeginGuildPrune(ctx, id, api.BeginGuildPruneParams{}))
	mustNoErr(c.GetGuildVanityURL(ctx, id))
	mustNoErr(c.GetGuildAuditLogRaw(ctx, id, api.GetGuildAuditLogParams{}))
	mustNoErr(c.GetGuildMember(ctx, id, id))
	mustNoErr(c.ModifyGuildMemberRaw(ctx, id, id, api.ModifyGuildMemberParams{}))
	mustNoErr(c.ModifyCurrentMember(ctx, id, api.ModifyCurrentMemberParams{}))
	mustNoErr(c.GetGuildWidgetSettings(ctx, id))
	mustNoErr(c.ModifyGuildWidgetSettings(ctx, id, api.ModifyGuildWidgetParams{}))
	mustNoErr(c.GetGuildWidget(ctx, id))
	mustNoErr(c.GetGuildWelcomeScreen(ctx, id))
	mustNoErr(c.ModifyGuildWelcomeScreen(ctx, id, api.ModifyGuildWelcomeScreenParams{}))
	mustNoErr(c.GetGuildOnboarding(ctx, id))
	mustNoErr(c.ModifyGuildOnboarding(ctx, id, api.ModifyGuildOnboardingParams{}))
	mustNoErr(c.GetGuildSoundboardSound(ctx, id, id))
	mustNoErr(c.CreateGuildSoundboardSound(ctx, id, api.CreateGuildSoundboardSoundParams{}))
	mustNoErr(c.ModifyGuildSoundboardSoundRaw(ctx, id, id, api.ModifyGuildSoundboardSoundParams{}))
	mustNoErr(c.GetCurrentApplication(ctx))
	mustNoErr(c.ModifyCurrentApplication(ctx, api.ModifyCurrentApplicationParams{}))
	ign(c.GetApplicationCommandPermissions(ctx, id, id))
	mustNoErr(c.GetApplicationEmoji(ctx, id, id))
	ign(c.CreateApplicationEmoji(ctx, id, api.CreateEmojiParams{}))
	mustNoErr(c.ModifyApplicationEmoji(ctx, id, id, api.ModifyEmojiParams{}))
	mustNoErr(c.GetEntitlement(ctx, id, id))
	mustNoErr(c.CreateTestEntitlement(ctx, id, api.CreateTestEntitlementParams{}))
	mustNoErr(c.GetSKUSubscription(ctx, id, id))
	mustNoErr(c.GetActivityInstance(ctx, id, "inst"))
	mustNoErr(c.FollowAnnouncementChannel(ctx, id, id))
	ign(c.CreateGuild(ctx, api.CreateGuildParams{}))
	ign(c.AddGuildMember(ctx, id, id, api.AddGuildMemberParams{}))
	mustNoErr(c.GetCurrentUserVoiceState(ctx, id))
	mustNoErr(c.GetUserVoiceState(ctx, id, id))
	mustNoErr(c.GetStickerPack(ctx, id))
	ign(c.CreateGuildStickerRaw(ctx, id, api.CreateGuildStickerParams{}))
	mustNoErr(c.GetInvite(ctx, "code", api.GetInviteParams{}))
	mustNoErr(c.DeleteInvite(ctx, "code", nil))
}

// TestErrorWrappers calls every error-only wrapper (excluding interaction- and
// gateway-based ones).
func (s *restWrappersSuite) TestErrorWrappers() {
	c, ctx := s.client, context.Background()
	id := sampleSnowflake
	reason := "cleanup"
	status := "live"

	// Seed a cached message so Pin/Unpin's in-place update path runs.
	c.Cache.Messages().Add(&discord.Message{ID: id, ChannelID: id})

	s.NoError(c.DeleteMessage(ctx, id, id, nil))
	s.NoError(c.BulkDeleteMessages(ctx, id, []discord.Snowflake{id, id + 1}, nil))
	s.NoError(c.PinMessage(ctx, id, id, &reason))
	s.NoError(c.UnpinMessage(ctx, id, id, &reason))
	s.NoError(c.AddReaction(ctx, id, id, "👍"))
	s.NoError(c.DeleteOwnReaction(ctx, id, id, "👍"))
	s.NoError(c.DeleteUserReaction(ctx, id, id, "👍", id))
	s.NoError(c.DeleteAllReactions(ctx, id, id))
	s.NoError(c.DeleteAllReactionsForEmoji(ctx, id, id, "👍"))
	s.NoError(c.EditChannelPermissions(ctx, id, id, api.EditChannelPermissionsParams{}))
	s.NoError(c.DeleteChannelPermission(ctx, id, id))
	s.NoError(c.TriggerTypingIndicator(ctx, id))
	s.NoError(c.DeleteGuild(ctx, id))
	s.NoError(c.ModifyGuildChannelPositions(ctx, id, []api.ModifyGuildChannelPositionsEntry{}))
	s.NoError(c.DeleteGuildRole(ctx, id, id, &reason))
	s.NoError(c.CreateGuildBanRaw(ctx, id, id, api.CreateGuildBanParams{}))
	s.NoError(c.RemoveGuildBan(ctx, id, id, &reason))
	s.NoError(c.AddGuildMemberRole(ctx, id, id, id, &reason))
	s.NoError(c.RemoveGuildMemberRole(ctx, id, id, id, &reason))
	s.NoError(c.RemoveGuildMember(ctx, id, id))
	s.NoError(c.ModifyCurrentUserVoiceState(ctx, id, api.ModifyCurrentUserVoiceStateParams{}))
	s.NoError(c.ModifyUserVoiceState(ctx, id, id, api.ModifyUserVoiceStateParams{}))
	s.NoError(c.DeleteGuildSoundboardSound(ctx, id, id, &reason))
	s.NoError(c.SendSoundboardSound(ctx, id, api.SendSoundboardSoundParams{}))
	s.NoError(c.DeleteApplicationEmoji(ctx, id, id))
	s.NoError(c.ConsumeEntitlement(ctx, id, id))
	s.NoError(c.DeleteTestEntitlement(ctx, id, id))
	s.NoError(c.SetVoiceChannelStatus(ctx, id, &status))
	s.NoError(c.AddGroupDMRecipient(ctx, id, id, api.AddGroupDMRecipientParams{}))
	s.NoError(c.RemoveGroupDMRecipient(ctx, id, id))
	s.NoError(c.DeleteGuildIntegration(ctx, id, id, &reason))
	_ = c.ExecuteSlackWebhook(ctx, id, "tok", false, json.RawMessage(`{}`))
	_ = c.ExecuteGitHubWebhook(ctx, id, "tok", false, json.RawMessage(`{}`))
}

// TestSliceWrappers calls every slice-returning wrapper against an array
// response (skipping gateway-backed FetchAllGuildMembers and BulkRegisterCommands).
func (s *restWrappersSuite) TestSliceWrappers() {
	c, ctx := s.client, context.Background()
	id := sampleSnowflake

	mustNoErr := func(_ any, err error) { s.NoError(err) }
	// ign: endpoints returning a wrapper object rather than a bare array.
	ign := func(_ any, _ error) {}

	s.setBody(`[{"id":"175928847299117063","channel_id":"175928847299117063","content":"hi","name":"x"}]`)

	mustNoErr(c.ListChannelMessagesRaw(ctx, id, api.GetMessagesParams{}))
	// ListPinnedMessages now walks the paginated Get Channel Pins object.
	s.setBody(`{"items":[],"has_more":false}`)
	mustNoErr(c.ListPinnedMessages(ctx, id))
	s.setBody(`[{"id":"175928847299117063","channel_id":"175928847299117063","content":"hi","name":"x"}]`)
	mustNoErr(c.ListReactions(ctx, id, id, "👍", api.GetReactionsParams{}))
	mustNoErr(c.ListChannelInvites(ctx, id))
	mustNoErr(c.ListGuildChannels(ctx, id))
	mustNoErr(c.ListGuildRoles(ctx, id))
	mustNoErr(c.ModifyGuildRolePositionsRaw(ctx, id, []api.ModifyGuildRolePositionsEntry{}))
	mustNoErr(c.ListGuildBans(ctx, id, discord.FetchBansOptions{}))
	mustNoErr(c.ListGuildInvites(ctx, id))
	mustNoErr(c.ListGuildMembersRaw(ctx, id, api.GetGuildMembersParams{}))
	mustNoErr(c.SearchGuildMembers(ctx, id, api.SearchGuildMembersParams{}))
	mustNoErr(c.ListVoiceRegions(ctx))
	mustNoErr(c.ListGuildVoiceRegions(ctx, id))
	mustNoErr(c.ListDefaultSoundboardSounds(ctx))
	ign(c.ListGuildSoundboardSounds(ctx, id))
	ign(c.ListGuildApplicationCommandPermissions(ctx, id))
	mustNoErr(c.ListEntitlements(ctx, id, api.ListEntitlementsParams{}))
	mustNoErr(c.ListSKUs(ctx, id))
	mustNoErr(c.ListSKUSubscriptions(ctx, id, api.ListSKUSubscriptionsParams{}))
	mustNoErr(c.ListGuildIntegrations(ctx, id))

	// []byte response — raw image bytes; the array body is returned verbatim.
	mustNoErr(c.GetGuildWidgetImage(ctx, id, "shield"))
}
