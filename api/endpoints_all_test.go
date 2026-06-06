package api

import (
	"context"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/commands"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

// Additional resource IDs used by the table-driven endpoint tests below.
const (
	testUserID        = discord.Snowflake(110000000000000001)
	testMsgID         = discord.Snowflake(120000000000000002)
	testEventID       = discord.Snowflake(130000000000000003)
	testRuleID        = discord.Snowflake(140000000000000004)
	testStickerID     = discord.Snowflake(150000000000000005)
	testSoundID       = discord.Snowflake(160000000000000006)
	testIntegrationID = discord.Snowflake(170000000000000007)
	testEntID         = discord.Snowflake(180000000000000008)
	testInteractionID = discord.Snowflake(190000000000000009)
	testWebhookID     = discord.Snowflake(210000000000000010)
	testOverwriteID   = discord.Snowflake(220000000000000011)
	testSkuID         = discord.Snowflake(230000000000000012)
	testSubID         = discord.Snowflake(240000000000000013)
	testPackID        = discord.Snowflake(250000000000000014)
	testReqID         = discord.Snowflake(260000000000000015)
)

// epCase describes one endpoint round-trip: invoke call, then assert the mock
// server received the expected verb and path. status/body default to 200/"{}".
type epCase struct {
	name       string
	method     string
	path       string // path suffix (without the version prefix); wrapped by apiPath
	status     int
	body       string
	wantReason string // when set, assert the X-Audit-Log-Reason header was forwarded
	call       func() error
}

// runEndpoints drives each case as a subtest: it arms the canned response, runs
// the client call, and asserts no error plus the captured verb and path. This
// verifies URL construction, verb selection and response decoding across the
// whole surface in one place.
func (s *endpointSuite) runEndpoints(cases []epCase) {
	for _, c := range cases {
		s.Run(c.name, func() {
			status := c.status
			if status == 0 {
				status = http.StatusOK
			}
			body := c.body
			if body == "" {
				body = "{}"
			}
			s.respondWith(status, body)

			err := c.call()
			s.Require().NoError(err)
			s.Equal(c.method, s.last.Method, "HTTP method")
			s.Equal(apiPath(c.path), s.last.Path, "request path")
			if c.wantReason != "" {
				s.Equal(c.wantReason, s.last.Reason, "X-Audit-Log-Reason header")
			}
		})
	}
}

// drop adapts a (T, error) call to the func() error shape used by epCase.
func drop[T any](_ T, err error) error { return err }

func (s *endpointSuite) TestGuildEndpoints() {
	ctx := context.Background()
	g := testGuildID.String()
	s.runEndpoints([]epCase{
		{name: "GetGuild", method: "GET", path: "/guilds/" + g, call: func() error {
			return drop(s.client.GetGuild(ctx, testGuildID, false))
		}},
		{name: "GetGuildPreview", method: "GET", path: "/guilds/" + g + "/preview", call: func() error {
			return drop(s.client.GetGuildPreview(ctx, testGuildID))
		}},
		{name: "ModifyGuild", method: "PATCH", path: "/guilds/" + g, call: func() error {
			return drop(s.client.ModifyGuild(ctx, testGuildID, ModifyGuildParams{}, nil))
		}},
		{name: "DeleteGuild", method: "DELETE", path: "/guilds/" + g, call: func() error {
			return s.client.DeleteGuild(ctx, testGuildID)
		}},
		{name: "ListGuildChannels", method: "GET", path: "/guilds/" + g + "/channels", body: "[]", call: func() error {
			return drop(s.client.ListGuildChannels(ctx, testGuildID))
		}},
		{name: "CreateGuildChannel", method: "POST", path: "/guilds/" + g + "/channels", status: http.StatusCreated, call: func() error {
			return drop(s.client.CreateGuildChannel(ctx, testGuildID, CreateGuildChannelParams{}, nil))
		}},
		{name: "ModifyGuildChannelPositions", method: "PATCH", path: "/guilds/" + g + "/channels", call: func() error {
			return s.client.ModifyGuildChannelPositions(ctx, testGuildID, nil, nil)
		}},
		{name: "ListGuildRoles", method: "GET", path: "/guilds/" + g + "/roles", body: "[]", call: func() error {
			return drop(s.client.ListGuildRoles(ctx, testGuildID))
		}},
		{name: "GetGuildRoleMemberCounts", method: "GET", path: "/guilds/" + g + "/roles/member-counts", call: func() error {
			return drop(s.client.GetGuildRoleMemberCounts(ctx, testGuildID))
		}},
		{name: "ModifyGuildRolePositions", method: "PATCH", path: "/guilds/" + g + "/roles", body: "[]", call: func() error {
			return drop(s.client.ModifyGuildRolePositions(ctx, testGuildID, nil, nil))
		}},
		{name: "ModifyGuildRole", method: "PATCH", path: "/guilds/" + g + "/roles/" + testRoleID.String(), call: func() error {
			return drop(s.client.ModifyGuildRole(ctx, testGuildID, testRoleID, ModifyGuildRoleParams{}, nil))
		}},
		{name: "ListGuildBans", method: "GET", path: "/guilds/" + g + "/bans", body: "[]", call: func() error {
			return drop(s.client.ListGuildBans(ctx, testGuildID, GetGuildBansParams{}))
		}},
		{name: "GetGuildBan", method: "GET", path: "/guilds/" + g + "/bans/" + testUserID.String(), call: func() error {
			return drop(s.client.GetGuildBan(ctx, testGuildID, testUserID))
		}},
		{name: "CreateGuildBan", method: "PUT", path: "/guilds/" + g + "/bans/" + testUserID.String(), call: func() error {
			return s.client.CreateGuildBan(ctx, testGuildID, testUserID, CreateGuildBanParams{}, nil)
		}},
		{name: "RemoveGuildBan", method: "DELETE", path: "/guilds/" + g + "/bans/" + testUserID.String(), call: func() error {
			return s.client.RemoveGuildBan(ctx, testGuildID, testUserID, nil)
		}},
		{name: "GetGuildPruneCount", method: "GET", path: "/guilds/" + g + "/prune", call: func() error {
			return drop(s.client.GetGuildPruneCount(ctx, testGuildID, GetGuildPruneCountParams{}))
		}},
		{name: "BeginGuildPrune", method: "POST", path: "/guilds/" + g + "/prune", call: func() error {
			return drop(s.client.BeginGuildPrune(ctx, testGuildID, BeginGuildPruneParams{}, nil))
		}},
		{name: "ListGuildInvites", method: "GET", path: "/guilds/" + g + "/invites", body: "[]", call: func() error {
			return drop(s.client.ListGuildInvites(ctx, testGuildID))
		}},
		{name: "GetGuildVanityURL", method: "GET", path: "/guilds/" + g + "/vanity-url", call: func() error {
			return drop(s.client.GetGuildVanityURL(ctx, testGuildID))
		}},
		{name: "GetGuildAuditLog", method: "GET", path: "/guilds/" + g + "/audit-logs", call: func() error {
			return drop(s.client.GetGuildAuditLog(ctx, testGuildID, GetGuildAuditLogParams{}))
		}},
		{name: "GetGuildWidgetSettings", method: "GET", path: "/guilds/" + g + "/widget", call: func() error {
			return drop(s.client.GetGuildWidgetSettings(ctx, testGuildID))
		}},
		{name: "ModifyGuildWidgetSettings", method: "PATCH", path: "/guilds/" + g + "/widget", call: func() error {
			return drop(s.client.ModifyGuildWidgetSettings(ctx, testGuildID, ModifyGuildWidgetParams{}))
		}},
		{name: "GetGuildWidgetImage", method: "GET", path: "/guilds/" + g + "/widget.png", call: func() error {
			return drop(s.client.GetGuildWidgetImage(ctx, testGuildID, ""))
		}},
		{name: "GetGuildWidget", method: "GET", path: "/guilds/" + g + "/widget.json", call: func() error {
			return drop(s.client.GetGuildWidget(ctx, testGuildID))
		}},
		{name: "GetGuildWelcomeScreen", method: "GET", path: "/guilds/" + g + "/welcome-screen", call: func() error {
			return drop(s.client.GetGuildWelcomeScreen(ctx, testGuildID))
		}},
		{name: "ModifyGuildWelcomeScreen", method: "PATCH", path: "/guilds/" + g + "/welcome-screen", call: func() error {
			return drop(s.client.ModifyGuildWelcomeScreen(ctx, testGuildID, ModifyGuildWelcomeScreenParams{}, nil))
		}},
		{name: "GetGuildOnboarding", method: "GET", path: "/guilds/" + g + "/onboarding", call: func() error {
			return drop(s.client.GetGuildOnboarding(ctx, testGuildID))
		}},
		{name: "ModifyGuildOnboarding", method: "PUT", path: "/guilds/" + g + "/onboarding", call: func() error {
			return drop(s.client.ModifyGuildOnboarding(ctx, testGuildID, ModifyGuildOnboardingParams{}, nil))
		}},
		{name: "CreateGuild", method: "POST", path: "/guilds", status: http.StatusCreated, call: func() error {
			return drop(s.client.CreateGuild(ctx, CreateGuildParams{}))
		}},
		{name: "BulkBanGuildMembers", method: "POST", path: "/guilds/" + g + "/bulk-ban", call: func() error {
			return drop(s.client.BulkBanGuildMembers(ctx, testGuildID, BulkBanParams{}, nil))
		}},
		{name: "GetGuildJoinRequests", method: "GET", path: "/guilds/" + g + "/requests", call: func() error {
			return drop(s.client.GetGuildJoinRequests(ctx, testGuildID, GetGuildJoinRequestsParams{}))
		}},
		{name: "ActionGuildJoinRequest", method: "PATCH", path: "/guilds/" + g + "/requests/" + testReqID.String(), call: func() error {
			return drop(s.client.ActionGuildJoinRequest(ctx, testGuildID, testReqID, ActionGuildJoinRequestParams{}))
		}},
		{name: "GetGuildNewMemberWelcome", method: "GET", path: "/guilds/" + g + "/new-member-welcome", call: func() error {
			return drop(s.client.GetGuildNewMemberWelcome(ctx, testGuildID))
		}},
		{name: "ListGuildIntegrations", method: "GET", path: "/guilds/" + g + "/integrations", body: "[]", call: func() error {
			return drop(s.client.ListGuildIntegrations(ctx, testGuildID))
		}},
		{name: "DeleteGuildIntegration", method: "DELETE", path: "/guilds/" + g + "/integrations/" + testIntegrationID.String(), call: func() error {
			return s.client.DeleteGuildIntegration(ctx, testGuildID, testIntegrationID, nil)
		}},
		{name: "ModifyGuildIncidentActions", method: "PUT", path: "/guilds/" + g + "/incident-actions", body: "{}", call: func() error {
			return drop(s.client.ModifyGuildIncidentActions(ctx, testGuildID, nil))
		}},
	})
}

func (s *endpointSuite) TestMemberEndpoints() {
	ctx := context.Background()
	g := testGuildID.String()
	u := testUserID.String()
	s.runEndpoints([]epCase{
		{name: "GetGuildMember", method: "GET", path: "/guilds/" + g + "/members/" + u, call: func() error {
			return drop(s.client.GetGuildMember(ctx, testGuildID, testUserID))
		}},
		{name: "ListGuildMembers", method: "GET", path: "/guilds/" + g + "/members", body: "[]", call: func() error {
			return drop(s.client.ListGuildMembers(ctx, testGuildID, GetGuildMembersParams{}))
		}},
		{name: "SearchGuildMembers", method: "GET", path: "/guilds/" + g + "/members/search", body: "[]", call: func() error {
			return drop(s.client.SearchGuildMembers(ctx, testGuildID, SearchGuildMembersParams{}))
		}},
		{name: "ModifyGuildMember", method: "PATCH", path: "/guilds/" + g + "/members/" + u, call: func() error {
			return drop(s.client.ModifyGuildMember(ctx, testGuildID, testUserID, ModifyGuildMemberParams{}, nil))
		}},
		{name: "ModifyCurrentMember", method: "PATCH", path: "/guilds/" + g + "/members/@me", call: func() error {
			return drop(s.client.ModifyCurrentMember(ctx, testGuildID, ModifyCurrentMemberParams{}))
		}},
		{name: "AddGuildMemberRole", method: "PUT", path: "/guilds/" + g + "/members/" + u + "/roles/" + testRoleID.String(), call: func() error {
			return s.client.AddGuildMemberRole(ctx, testGuildID, testUserID, testRoleID, nil)
		}},
		{name: "RemoveGuildMemberRole", method: "DELETE", path: "/guilds/" + g + "/members/" + u + "/roles/" + testRoleID.String(), call: func() error {
			return s.client.RemoveGuildMemberRole(ctx, testGuildID, testUserID, testRoleID, nil)
		}},
		{name: "KickGuildMember", method: "DELETE", path: "/guilds/" + g + "/members/" + u, call: func() error {
			return s.client.KickGuildMember(ctx, testGuildID, testUserID, nil)
		}},
	})
}

func (s *endpointSuite) TestChannelEndpoints() {
	ctx := context.Background()
	ch := testChanID.String()
	s.runEndpoints([]epCase{
		{name: "ModifyChannel", method: "PATCH", path: "/channels/" + ch, call: func() error {
			return drop(s.client.ModifyChannel(ctx, testChanID, ModifyChannelParams{}, nil))
		}},
		{name: "ListChannelInvites", method: "GET", path: "/channels/" + ch + "/invites", body: "[]", call: func() error {
			return drop(s.client.ListChannelInvites(ctx, testChanID))
		}},
		{name: "CreateChannelInvite", method: "POST", path: "/channels/" + ch + "/invites", call: func() error {
			return drop(s.client.CreateChannelInvite(ctx, testChanID, CreateChannelInviteParams{}, nil))
		}},
		{name: "EditChannelPermissions", method: "PUT", path: "/channels/" + ch + "/permissions/" + testOverwriteID.String(), call: func() error {
			return s.client.EditChannelPermissions(ctx, testChanID, testOverwriteID, EditChannelPermissionsParams{}, nil)
		}},
		{name: "DeleteChannelPermission", method: "DELETE", path: "/channels/" + ch + "/permissions/" + testOverwriteID.String(), call: func() error {
			return s.client.DeleteChannelPermission(ctx, testChanID, testOverwriteID, nil)
		}},
		{name: "SetVoiceChannelStatus", method: "PUT", path: "/channels/" + ch + "/voice-status", call: func() error {
			return s.client.SetVoiceChannelStatus(ctx, testChanID, nil)
		}},
		{name: "FollowAnnouncementChannel", method: "POST", path: "/channels/" + ch + "/followers", call: func() error {
			return drop(s.client.FollowAnnouncementChannel(ctx, testChanID, testGuildID, nil))
		}},
		{name: "AddGroupDMRecipient", method: "PUT", path: "/channels/" + ch + "/recipients/" + testUserID.String(), call: func() error {
			return s.client.AddGroupDMRecipient(ctx, testChanID, testUserID, AddGroupDMRecipientParams{})
		}},
		{name: "RemoveGroupDMRecipient", method: "DELETE", path: "/channels/" + ch + "/recipients/" + testUserID.String(), call: func() error {
			return s.client.RemoveGroupDMRecipient(ctx, testChanID, testUserID)
		}},
	})
}

func (s *endpointSuite) TestCommandEndpoints() {
	ctx := context.Background()
	app := testAppID.String()
	g := testGuildID.String()
	cmd := &commands.ApplicationCommand{Name: "c", Description: "d"}
	s.runEndpoints([]epCase{
		{name: "BulkRegisterCommands", method: "PUT", path: "/applications/" + app + "/commands", body: "[]", call: func() error {
			return drop(s.client.BulkRegisterCommands(ctx, testAppID, nil))
		}},
		{name: "GetGlobalApplicationCommand", method: "GET", path: "/applications/" + app + "/commands/" + testCmdID.String(), call: func() error {
			return drop(s.client.GetGlobalApplicationCommand(ctx, testAppID, testCmdID))
		}},
		{name: "EditGlobalApplicationCommand", method: "PATCH", path: "/applications/" + app + "/commands/" + testCmdID.String(), call: func() error {
			return drop(s.client.EditGlobalApplicationCommand(ctx, testAppID, testCmdID, cmd))
		}},
		{name: "ListGuildApplicationCommands", method: "GET", path: "/applications/" + app + "/guilds/" + g + "/commands", body: "[]", call: func() error {
			return drop(s.client.ListGuildApplicationCommands(ctx, testAppID, testGuildID, false))
		}},
		{name: "CreateGuildApplicationCommand", method: "POST", path: "/applications/" + app + "/guilds/" + g + "/commands", status: http.StatusCreated, call: func() error {
			return drop(s.client.CreateGuildApplicationCommand(ctx, testAppID, testGuildID, cmd))
		}},
		{name: "GetGuildApplicationCommand", method: "GET", path: "/applications/" + app + "/guilds/" + g + "/commands/" + testCmdID.String(), call: func() error {
			return drop(s.client.GetGuildApplicationCommand(ctx, testAppID, testGuildID, testCmdID))
		}},
		{name: "EditGuildApplicationCommand", method: "PATCH", path: "/applications/" + app + "/guilds/" + g + "/commands/" + testCmdID.String(), call: func() error {
			return drop(s.client.EditGuildApplicationCommand(ctx, testAppID, testGuildID, testCmdID, cmd))
		}},
		{name: "DeleteGuildApplicationCommand", method: "DELETE", path: "/applications/" + app + "/guilds/" + g + "/commands/" + testCmdID.String(), call: func() error {
			return s.client.DeleteGuildApplicationCommand(ctx, testAppID, testGuildID, testCmdID)
		}},
		{name: "BulkOverwriteGuildApplicationCommands", method: "PUT", path: "/applications/" + app + "/guilds/" + g + "/commands", body: "[]", call: func() error {
			return drop(s.client.BulkOverwriteGuildApplicationCommands(ctx, testAppID, testGuildID, nil))
		}},
		{name: "ListGuildApplicationCommandPermissions", method: "GET", path: "/applications/" + app + "/guilds/" + g + "/commands/permissions", body: "[]", call: func() error {
			return drop(s.client.ListGuildApplicationCommandPermissions(ctx, testAppID, testGuildID))
		}},
		{name: "GetApplicationCommandPermissions", method: "GET", path: "/applications/" + app + "/guilds/" + g + "/commands/" + testCmdID.String() + "/permissions", call: func() error {
			return drop(s.client.GetApplicationCommandPermissions(ctx, testAppID, testGuildID, testCmdID))
		}},
		{name: "EditApplicationCommandPermissions", method: "PUT", path: "/applications/" + app + "/guilds/" + g + "/commands/" + testCmdID.String() + "/permissions", call: func() error {
			return drop(s.client.EditApplicationCommandPermissions(ctx, testAppID, testGuildID, testCmdID, nil))
		}},
	})
}

func (s *endpointSuite) TestEmojiEndpoints() {
	ctx := context.Background()
	g := testGuildID.String()
	app := testAppID.String()
	e := testEmojiID.String()
	s.runEndpoints([]epCase{
		{name: "GetGuildEmoji", method: "GET", path: "/guilds/" + g + "/emojis/" + e, call: func() error {
			return drop(s.client.GetGuildEmoji(ctx, testGuildID, testEmojiID))
		}},
		{name: "CreateGuildEmoji", method: "POST", path: "/guilds/" + g + "/emojis", status: http.StatusCreated, call: func() error {
			return drop(s.client.CreateGuildEmoji(ctx, testGuildID, CreateGuildEmojiParams{}, nil))
		}},
		{name: "ModifyGuildEmoji", method: "PATCH", path: "/guilds/" + g + "/emojis/" + e, call: func() error {
			return drop(s.client.ModifyGuildEmoji(ctx, testGuildID, testEmojiID, ModifyGuildEmojiParams{}, nil))
		}},
		{name: "ListApplicationEmojis", method: "GET", path: "/applications/" + app + "/emojis", call: func() error {
			return drop(s.client.ListApplicationEmojis(ctx, testAppID))
		}},
		{name: "GetApplicationEmoji", method: "GET", path: "/applications/" + app + "/emojis/" + e, call: func() error {
			return drop(s.client.GetApplicationEmoji(ctx, testAppID, testEmojiID))
		}},
		{name: "CreateApplicationEmoji", method: "POST", path: "/applications/" + app + "/emojis", status: http.StatusCreated, call: func() error {
			return drop(s.client.CreateApplicationEmoji(ctx, testAppID, CreateEmojiParams{}))
		}},
		{name: "ModifyApplicationEmoji", method: "PATCH", path: "/applications/" + app + "/emojis/" + e, call: func() error {
			return drop(s.client.ModifyApplicationEmoji(ctx, testAppID, testEmojiID, ModifyEmojiParams{}))
		}},
		{name: "DeleteApplicationEmoji", method: "DELETE", path: "/applications/" + app + "/emojis/" + e, call: func() error {
			return s.client.DeleteApplicationEmoji(ctx, testAppID, testEmojiID)
		}},
	})
}

func (s *endpointSuite) TestScheduledEventEndpoints() {
	ctx := context.Background()
	g := testGuildID.String()
	ev := testEventID.String()
	s.runEndpoints([]epCase{
		{name: "ListGuildScheduledEvents", method: "GET", path: "/guilds/" + g + "/scheduled-events", body: "[]", call: func() error {
			return drop(s.client.ListGuildScheduledEvents(ctx, testGuildID, false))
		}},
		{name: "CreateGuildScheduledEvent", method: "POST", path: "/guilds/" + g + "/scheduled-events", call: func() error {
			return drop(s.client.CreateGuildScheduledEvent(ctx, testGuildID, CreateGuildScheduledEventParams{}))
		}},
		{name: "GetGuildScheduledEvent", method: "GET", path: "/guilds/" + g + "/scheduled-events/" + ev, call: func() error {
			return drop(s.client.GetGuildScheduledEvent(ctx, testGuildID, testEventID, false))
		}},
		{name: "ModifyGuildScheduledEvent", method: "PATCH", path: "/guilds/" + g + "/scheduled-events/" + ev, call: func() error {
			return drop(s.client.ModifyGuildScheduledEvent(ctx, testGuildID, testEventID, ModifyGuildScheduledEventParams{}))
		}},
		{name: "DeleteGuildScheduledEvent", method: "DELETE", path: "/guilds/" + g + "/scheduled-events/" + ev, call: func() error {
			return s.client.DeleteGuildScheduledEvent(ctx, testGuildID, testEventID)
		}},
		{name: "ListGuildScheduledEventUsers", method: "GET", path: "/guilds/" + g + "/scheduled-events/" + ev + "/users", body: "[]", call: func() error {
			return drop(s.client.ListGuildScheduledEventUsers(ctx, testGuildID, testEventID, GetGuildScheduledEventUsersParams{}))
		}},
	})
}

func (s *endpointSuite) TestThreadEndpoints() {
	ctx := context.Background()
	ch := testChanID.String()
	s.runEndpoints([]epCase{
		{name: "CreateThreadFromMessage", method: "POST", path: "/channels/" + ch + "/messages/" + testMsgID.String() + "/threads", status: http.StatusCreated, call: func() error {
			return drop(s.client.CreateThreadFromMessage(ctx, testChanID, testMsgID, CreateThreadFromMessageParams{}, nil))
		}},
		{name: "CreateThread", method: "POST", path: "/channels/" + ch + "/threads", status: http.StatusCreated, call: func() error {
			return drop(s.client.CreateThread(ctx, testChanID, CreateThreadParams{}, nil))
		}},
		{name: "CreateForumThread", method: "POST", path: "/channels/" + ch + "/threads", status: http.StatusCreated, call: func() error {
			return drop(s.client.CreateForumThread(ctx, testChanID, CreateForumThreadParams{}, nil))
		}},
		{name: "JoinThread", method: "PUT", path: "/channels/" + ch + "/thread-members/@me", call: func() error {
			return s.client.JoinThread(ctx, testChanID)
		}},
		{name: "LeaveThread", method: "DELETE", path: "/channels/" + ch + "/thread-members/@me", call: func() error {
			return s.client.LeaveThread(ctx, testChanID)
		}},
		{name: "AddThreadMember", method: "PUT", path: "/channels/" + ch + "/thread-members/" + testUserID.String(), call: func() error {
			return s.client.AddThreadMember(ctx, testChanID, testUserID)
		}},
		{name: "RemoveThreadMember", method: "DELETE", path: "/channels/" + ch + "/thread-members/" + testUserID.String(), call: func() error {
			return s.client.RemoveThreadMember(ctx, testChanID, testUserID)
		}},
		{name: "GetThreadMember", method: "GET", path: "/channels/" + ch + "/thread-members/" + testUserID.String(), call: func() error {
			return drop(s.client.GetThreadMember(ctx, testChanID, testUserID, false))
		}},
		{name: "ListThreadMembers", method: "GET", path: "/channels/" + ch + "/thread-members", body: "[]", call: func() error {
			return drop(s.client.ListThreadMembers(ctx, testChanID, ListThreadMembersParams{}))
		}},
		{name: "ListPublicArchivedThreads", method: "GET", path: "/channels/" + ch + "/threads/archived/public", call: func() error {
			return drop(s.client.ListPublicArchivedThreads(ctx, testChanID, ListArchivedThreadsParams{}))
		}},
		{name: "ListPrivateArchivedThreads", method: "GET", path: "/channels/" + ch + "/threads/archived/private", call: func() error {
			return drop(s.client.ListPrivateArchivedThreads(ctx, testChanID, ListArchivedThreadsParams{}))
		}},
		{name: "ListJoinedPrivateArchivedThreads", method: "GET", path: "/channels/" + ch + "/users/@me/threads/archived/private", call: func() error {
			return drop(s.client.ListJoinedPrivateArchivedThreads(ctx, testChanID, ListArchivedThreadsParams{}))
		}},
		{name: "SearchThreads", method: "GET", path: "/channels/" + ch + "/threads/search", call: func() error {
			return drop(s.client.SearchThreads(ctx, testChanID, SearchThreadsParams{}))
		}},
		{name: "ListActiveGuildThreads", method: "GET", path: "/guilds/" + testGuildID.String() + "/threads/active", call: func() error {
			return drop(s.client.ListActiveGuildThreads(ctx, testGuildID))
		}},
	})
}

func (s *endpointSuite) TestMessageEndpoints() {
	ctx := context.Background()
	ch := testChanID.String()
	msg := testMsgID.String()
	s.runEndpoints([]epCase{
		{name: "ListMessages", method: "GET", path: "/channels/" + ch + "/messages", body: "[]", call: func() error {
			return drop(s.client.ListMessages(ctx, testChanID, GetMessagesParams{}))
		}},
		{name: "GetMessage", method: "GET", path: "/channels/" + ch + "/messages/" + msg, call: func() error {
			return drop(s.client.GetMessage(ctx, testChanID, testMsgID))
		}},
		{name: "EditMessage", method: "PATCH", path: "/channels/" + ch + "/messages/" + msg, call: func() error {
			return drop(s.client.EditMessage(ctx, testChanID, testMsgID, EditMessageParams{}))
		}},
		{name: "DeleteMessage", method: "DELETE", path: "/channels/" + ch + "/messages/" + msg, call: func() error {
			return s.client.DeleteMessage(ctx, testChanID, testMsgID, nil)
		}},
		{name: "BulkDeleteMessages", method: "POST", path: "/channels/" + ch + "/messages/bulk-delete", call: func() error {
			return s.client.BulkDeleteMessages(ctx, testChanID, []discord.Snowflake{testMsgID, testUserID}, nil)
		}},
		{name: "CrosspostMessage", method: "POST", path: "/channels/" + ch + "/messages/" + msg + "/crosspost", call: func() error {
			return drop(s.client.CrosspostMessage(ctx, testChanID, testMsgID))
		}},
		{name: "GetChannelPins", method: "GET", path: "/channels/" + ch + "/messages/pins", body: `{"items":[],"has_more":false}`, call: func() error {
			return drop(s.client.GetChannelPins(ctx, testChanID, GetChannelPinsParams{}))
		}},
		{name: "ListPinnedMessages", method: "GET", path: "/channels/" + ch + "/messages/pins", body: `{"items":[],"has_more":false}`, call: func() error {
			return drop(s.client.ListPinnedMessages(ctx, testChanID))
		}},
		{name: "PinMessage", method: "PUT", path: "/channels/" + ch + "/messages/pins/" + msg, call: func() error {
			return s.client.PinMessage(ctx, testChanID, testMsgID, nil)
		}},
		{name: "UnpinMessage", method: "DELETE", path: "/channels/" + ch + "/messages/pins/" + msg, call: func() error {
			return s.client.UnpinMessage(ctx, testChanID, testMsgID, nil)
		}},
		{name: "AddReaction", method: "PUT", path: "/channels/" + ch + "/messages/" + msg + "/reactions/fire/@me", call: func() error {
			return s.client.AddReaction(ctx, testChanID, testMsgID, "fire")
		}},
		{name: "DeleteOwnReaction", method: "DELETE", path: "/channels/" + ch + "/messages/" + msg + "/reactions/fire/@me", call: func() error {
			return s.client.DeleteOwnReaction(ctx, testChanID, testMsgID, "fire")
		}},
		{name: "DeleteUserReaction", method: "DELETE", path: "/channels/" + ch + "/messages/" + msg + "/reactions/fire/" + testUserID.String(), call: func() error {
			return s.client.DeleteUserReaction(ctx, testChanID, testMsgID, "fire", testUserID)
		}},
		{name: "ListReactions", method: "GET", path: "/channels/" + ch + "/messages/" + msg + "/reactions/fire", body: "[]", call: func() error {
			return drop(s.client.ListReactions(ctx, testChanID, testMsgID, "fire", GetReactionsParams{}))
		}},
		{name: "DeleteAllReactions", method: "DELETE", path: "/channels/" + ch + "/messages/" + msg + "/reactions", call: func() error {
			return s.client.DeleteAllReactions(ctx, testChanID, testMsgID)
		}},
		{name: "DeleteAllReactionsForEmoji", method: "DELETE", path: "/channels/" + ch + "/messages/" + msg + "/reactions/fire", call: func() error {
			return s.client.DeleteAllReactionsForEmoji(ctx, testChanID, testMsgID, "fire")
		}},
		{name: "SearchGuildMessages", method: "GET", path: "/guilds/" + testGuildID.String() + "/messages/search", call: func() error {
			return drop(s.client.SearchGuildMessages(ctx, testGuildID, SearchGuildMessagesParams{}))
		}},
		{name: "CreatePoll", method: "POST", path: "/channels/" + ch + "/messages", call: func() error {
			return drop(s.client.CreatePoll(ctx, testChanID, discord.PollRequest{}))
		}},
		{name: "EndPoll", method: "POST", path: "/channels/" + ch + "/polls/" + msg + "/expire", call: func() error {
			return drop(s.client.EndPoll(ctx, testChanID, testMsgID))
		}},
	})
}

// TestAuditReasonForwarding exercises the opts != nil branch of every endpoint
// that accepts an audit-log reason, asserting the reason is forwarded as the
// X-Audit-Log-Reason header. This complements the nil-opts paths covered by the
// resource tables above. A single-token reason avoids percent-encoding (that is
// covered separately in restclient_auditlog_test.go).
func (s *endpointSuite) TestAuditReasonForwarding() {
	ctx := context.Background()
	const r = "audit-me"
	g := testGuildID.String()
	ch := testChanID.String()
	created := http.StatusCreated
	s.runEndpoints([]epCase{
		{name: "ModifyGuild", method: "PATCH", path: "/guilds/" + g, wantReason: r, call: func() error {
			return drop(s.client.ModifyGuild(ctx, testGuildID, ModifyGuildParams{}, &ModifyGuildOptions{Reason: r}))
		}},
		{name: "CreateGuildChannel", method: "POST", path: "/guilds/" + g + "/channels", status: created, wantReason: r, call: func() error {
			return drop(s.client.CreateGuildChannel(ctx, testGuildID, CreateGuildChannelParams{}, &CreateGuildChannelOptions{Reason: r}))
		}},
		{name: "ModifyGuildChannelPositions", method: "PATCH", path: "/guilds/" + g + "/channels", wantReason: r, call: func() error {
			return s.client.ModifyGuildChannelPositions(ctx, testGuildID, nil, &ModifyGuildChannelPositionsOptions{Reason: r})
		}},
		{name: "ModifyGuildRolePositions", method: "PATCH", path: "/guilds/" + g + "/roles", body: "[]", wantReason: r, call: func() error {
			return drop(s.client.ModifyGuildRolePositions(ctx, testGuildID, nil, &ModifyGuildRolePositionsOptions{Reason: r}))
		}},
		{name: "ModifyGuildRole", method: "PATCH", path: "/guilds/" + g + "/roles/" + testRoleID.String(), wantReason: r, call: func() error {
			return drop(s.client.ModifyGuildRole(ctx, testGuildID, testRoleID, ModifyGuildRoleParams{}, &ModifyGuildRoleOptions{Reason: r}))
		}},
		{name: "CreateGuildBan", method: "PUT", path: "/guilds/" + g + "/bans/" + testUserID.String(), wantReason: r, call: func() error {
			return s.client.CreateGuildBan(ctx, testGuildID, testUserID, CreateGuildBanParams{}, &CreateGuildBanOptions{Reason: r})
		}},
		{name: "RemoveGuildBan", method: "DELETE", path: "/guilds/" + g + "/bans/" + testUserID.String(), wantReason: r, call: func() error {
			return s.client.RemoveGuildBan(ctx, testGuildID, testUserID, &RemoveGuildBanOptions{Reason: r})
		}},
		{name: "BeginGuildPrune", method: "POST", path: "/guilds/" + g + "/prune", wantReason: r, call: func() error {
			return drop(s.client.BeginGuildPrune(ctx, testGuildID, BeginGuildPruneParams{}, &BeginGuildPruneOptions{Reason: r}))
		}},
		{name: "ModifyGuildWelcomeScreen", method: "PATCH", path: "/guilds/" + g + "/welcome-screen", wantReason: r, call: func() error {
			return drop(s.client.ModifyGuildWelcomeScreen(ctx, testGuildID, ModifyGuildWelcomeScreenParams{}, &ModifyGuildWelcomeScreenOptions{Reason: r}))
		}},
		{name: "ModifyGuildOnboarding", method: "PUT", path: "/guilds/" + g + "/onboarding", wantReason: r, call: func() error {
			return drop(s.client.ModifyGuildOnboarding(ctx, testGuildID, ModifyGuildOnboardingParams{}, &ModifyGuildOnboardingOptions{Reason: r}))
		}},
		{name: "BulkBanGuildMembers", method: "POST", path: "/guilds/" + g + "/bulk-ban", wantReason: r, call: func() error {
			return drop(s.client.BulkBanGuildMembers(ctx, testGuildID, BulkBanParams{}, &BulkBanOptions{Reason: r}))
		}},
		{name: "DeleteGuildIntegration", method: "DELETE", path: "/guilds/" + g + "/integrations/" + testIntegrationID.String(), wantReason: r, call: func() error {
			return s.client.DeleteGuildIntegration(ctx, testGuildID, testIntegrationID, &DeleteGuildIntegrationOptions{Reason: r})
		}},
		{name: "ModifyChannel", method: "PATCH", path: "/channels/" + ch, wantReason: r, call: func() error {
			return drop(s.client.ModifyChannel(ctx, testChanID, ModifyChannelParams{}, &ModifyChannelOptions{Reason: r}))
		}},
		{name: "CreateChannelInvite", method: "POST", path: "/channels/" + ch + "/invites", wantReason: r, call: func() error {
			return drop(s.client.CreateChannelInvite(ctx, testChanID, CreateChannelInviteParams{}, &CreateChannelInviteOptions{Reason: r}))
		}},
		{name: "EditChannelPermissions", method: "PUT", path: "/channels/" + ch + "/permissions/" + testOverwriteID.String(), wantReason: r, call: func() error {
			return s.client.EditChannelPermissions(ctx, testChanID, testOverwriteID, EditChannelPermissionsParams{}, &EditChannelPermissionsOptions{Reason: r})
		}},
		{name: "DeleteChannelPermission", method: "DELETE", path: "/channels/" + ch + "/permissions/" + testOverwriteID.String(), wantReason: r, call: func() error {
			return s.client.DeleteChannelPermission(ctx, testChanID, testOverwriteID, &DeleteChannelPermissionOptions{Reason: r})
		}},
		{name: "ModifyGuildMember", method: "PATCH", path: "/guilds/" + g + "/members/" + testUserID.String(), wantReason: r, call: func() error {
			return drop(s.client.ModifyGuildMember(ctx, testGuildID, testUserID, ModifyGuildMemberParams{}, &ModifyGuildMemberOptions{Reason: r}))
		}},
		{name: "AddGuildMemberRole", method: "PUT", path: "/guilds/" + g + "/members/" + testUserID.String() + "/roles/" + testRoleID.String(), wantReason: r, call: func() error {
			return s.client.AddGuildMemberRole(ctx, testGuildID, testUserID, testRoleID, &AddGuildMemberRoleOptions{Reason: r})
		}},
		{name: "RemoveGuildMemberRole", method: "DELETE", path: "/guilds/" + g + "/members/" + testUserID.String() + "/roles/" + testRoleID.String(), wantReason: r, call: func() error {
			return s.client.RemoveGuildMemberRole(ctx, testGuildID, testUserID, testRoleID, &RemoveGuildMemberRoleOptions{Reason: r})
		}},
		{name: "KickGuildMember", method: "DELETE", path: "/guilds/" + g + "/members/" + testUserID.String(), wantReason: r, call: func() error {
			return s.client.KickGuildMember(ctx, testGuildID, testUserID, &KickGuildMemberOptions{Reason: r})
		}},
		{name: "CreateGuildEmoji", method: "POST", path: "/guilds/" + g + "/emojis", status: created, wantReason: r, call: func() error {
			return drop(s.client.CreateGuildEmoji(ctx, testGuildID, CreateGuildEmojiParams{}, &CreateGuildEmojiOptions{Reason: r}))
		}},
		{name: "ModifyGuildEmoji", method: "PATCH", path: "/guilds/" + g + "/emojis/" + testEmojiID.String(), wantReason: r, call: func() error {
			return drop(s.client.ModifyGuildEmoji(ctx, testGuildID, testEmojiID, ModifyGuildEmojiParams{}, &ModifyGuildEmojiOptions{Reason: r}))
		}},
		{name: "CreateAutoModerationRule", method: "POST", path: "/guilds/" + g + "/auto-moderation/rules", wantReason: r, call: func() error {
			return drop(s.client.CreateAutoModerationRule(ctx, testGuildID, CreateAutoModerationRuleParams{}, &CreateAutoModerationRuleOptions{Reason: r}))
		}},
		{name: "ModifyAutoModerationRule", method: "PATCH", path: "/guilds/" + g + "/auto-moderation/rules/" + testRuleID.String(), wantReason: r, call: func() error {
			return drop(s.client.ModifyAutoModerationRule(ctx, testGuildID, testRuleID, ModifyAutoModerationRuleParams{}, &ModifyAutoModerationRuleOptions{Reason: r}))
		}},
		{name: "DeleteAutoModerationRule", method: "DELETE", path: "/guilds/" + g + "/auto-moderation/rules/" + testRuleID.String(), wantReason: r, call: func() error {
			return s.client.DeleteAutoModerationRule(ctx, testGuildID, testRuleID, &DeleteAutoModerationRuleOptions{Reason: r})
		}},
		{name: "CreateStageInstance", method: "POST", path: "/stage-instances", status: created, wantReason: r, call: func() error {
			return drop(s.client.CreateStageInstance(ctx, CreateStageInstanceParams{ChannelID: testChanID}, &CreateStageInstanceOptions{Reason: r}))
		}},
		{name: "ModifyStageInstance", method: "PATCH", path: "/stage-instances/" + ch, wantReason: r, call: func() error {
			return drop(s.client.ModifyStageInstance(ctx, testChanID, ModifyStageInstanceParams{}, &ModifyStageInstanceOptions{Reason: r}))
		}},
		{name: "DeleteStageInstance", method: "DELETE", path: "/stage-instances/" + ch, wantReason: r, call: func() error {
			return s.client.DeleteStageInstance(ctx, testChanID, &DeleteStageInstanceOptions{Reason: r})
		}},
		{name: "CreateGuildSoundboardSound", method: "POST", path: "/guilds/" + g + "/soundboard-sounds", wantReason: r, call: func() error {
			return drop(s.client.CreateGuildSoundboardSound(ctx, testGuildID, CreateGuildSoundboardSoundParams{}, &CreateGuildSoundboardSoundOptions{Reason: r}))
		}},
		{name: "ModifyGuildSoundboardSound", method: "PATCH", path: "/guilds/" + g + "/soundboard-sounds/" + testSoundID.String(), wantReason: r, call: func() error {
			return drop(s.client.ModifyGuildSoundboardSound(ctx, testGuildID, testSoundID, ModifyGuildSoundboardSoundParams{}, &ModifyGuildSoundboardSoundOptions{Reason: r}))
		}},
		{name: "DeleteGuildSoundboardSound", method: "DELETE", path: "/guilds/" + g + "/soundboard-sounds/" + testSoundID.String(), wantReason: r, call: func() error {
			return s.client.DeleteGuildSoundboardSound(ctx, testGuildID, testSoundID, &DeleteGuildSoundboardSoundOptions{Reason: r})
		}},
		{name: "ModifyGuildSticker", method: "PATCH", path: "/guilds/" + g + "/stickers/" + testStickerID.String(), wantReason: r, call: func() error {
			return drop(s.client.ModifyGuildSticker(ctx, testGuildID, testStickerID, ModifyGuildStickerParams{}, &ModifyGuildStickerOptions{Reason: r}))
		}},
		{name: "DeleteGuildSticker", method: "DELETE", path: "/guilds/" + g + "/stickers/" + testStickerID.String(), wantReason: r, call: func() error {
			return s.client.DeleteGuildSticker(ctx, testGuildID, testStickerID, &DeleteGuildStickerOptions{Reason: r})
		}},
		{name: "CreateGuildSticker", method: "POST", path: "/guilds/" + g + "/stickers", status: created, wantReason: r, call: func() error {
			return drop(s.client.CreateGuildSticker(ctx, testGuildID, CreateGuildStickerParams{
				Name: "s", Description: "d", Tags: "t", File: []byte("png"), ContentType: "image/png",
			}, &CreateGuildStickerOptions{Reason: r}))
		}},
		{name: "CreateWebhook", method: "POST", path: "/channels/" + ch + "/webhooks", wantReason: r, call: func() error {
			return drop(s.client.CreateWebhook(ctx, testChanID, CreateWebhookParams{}, &CreateWebhookOptions{Reason: r}))
		}},
		{name: "ModifyWebhook", method: "PATCH", path: "/webhooks/" + testWebhookID.String(), wantReason: r, call: func() error {
			return drop(s.client.ModifyWebhook(ctx, testWebhookID, ModifyWebhookParams{}, &ModifyWebhookOptions{Reason: r}))
		}},
		{name: "DeleteWebhook", method: "DELETE", path: "/webhooks/" + testWebhookID.String(), wantReason: r, call: func() error {
			return s.client.DeleteWebhook(ctx, testWebhookID, &DeleteWebhookOptions{Reason: r})
		}},
		{name: "CreateThreadFromMessage", method: "POST", path: "/channels/" + ch + "/messages/" + testMsgID.String() + "/threads", status: created, wantReason: r, call: func() error {
			return drop(s.client.CreateThreadFromMessage(ctx, testChanID, testMsgID, CreateThreadFromMessageParams{}, &CreateThreadFromMessageOptions{Reason: r}))
		}},
		{name: "CreateThread", method: "POST", path: "/channels/" + ch + "/threads", status: created, wantReason: r, call: func() error {
			return drop(s.client.CreateThread(ctx, testChanID, CreateThreadParams{}, &CreateThreadOptions{Reason: r}))
		}},
		{name: "CreateForumThread", method: "POST", path: "/channels/" + ch + "/threads", status: created, wantReason: r, call: func() error {
			return drop(s.client.CreateForumThread(ctx, testChanID, CreateForumThreadParams{}, &CreateForumThreadOptions{Reason: r}))
		}},
		{name: "DeleteMessage", method: "DELETE", path: "/channels/" + ch + "/messages/" + testMsgID.String(), wantReason: r, call: func() error {
			return s.client.DeleteMessage(ctx, testChanID, testMsgID, &DeleteMessageOptions{Reason: r})
		}},
		{name: "PinMessage", method: "PUT", path: "/channels/" + ch + "/messages/pins/" + testMsgID.String(), wantReason: r, call: func() error {
			return s.client.PinMessage(ctx, testChanID, testMsgID, &PinMessageOptions{Reason: r})
		}},
		{name: "UnpinMessage", method: "DELETE", path: "/channels/" + ch + "/messages/pins/" + testMsgID.String(), wantReason: r, call: func() error {
			return s.client.UnpinMessage(ctx, testChanID, testMsgID, &UnpinMessageOptions{Reason: r})
		}},
		{name: "BulkDeleteMessages", method: "POST", path: "/channels/" + ch + "/messages/bulk-delete", wantReason: r, call: func() error {
			return s.client.BulkDeleteMessages(ctx, testChanID, []discord.Snowflake{testMsgID, testUserID}, &BulkDeleteMessagesOptions{Reason: r})
		}},
	})
}

// TestGetPollAnswerVoters_PathValidation verifies the small integer poll answer
// ID (1, 2, 3 …) is accepted by the path guard: such all-digit segments are
// injection-safe, and Snowflake parameters are validated at the call site, so
// the request reaches the server.
func (s *endpointSuite) TestGetPollAnswerVoters_PathValidation() {
	s.respondWith(http.StatusOK, "{}")

	_, err := s.client.GetPollAnswerVoters(context.Background(), testChanID, testMsgID, 1, GetPollAnswerVotersParams{})

	s.Require().NoError(err, "small answer ID must be accepted, not rejected as a non-Snowflake")
	s.Equal(int32(1), s.last.CallCount, "request must reach the server")
	s.Contains(s.last.Path, "/answers/1", "answer id must appear in the request path")
}

// TestValidateAPIPath_RejectsOversizedDigitRun keeps a guard on abnormally long
// digit runs (no real Discord ID/index is longer than 20 digits).
func (s *endpointSuite) TestValidateAPIPath_RejectsOversizedDigitRun() {
	err := validateAPIPath("/channels/123456789012345678901234567890/messages")
	s.Require().Error(err)
	s.Contains(err.Error(), "too long")
}

func (s *endpointSuite) TestInteractionEndpoints() {
	ctx := context.Background()
	wh := testWebhookID.String()
	app := testAppID.String()
	s.runEndpoints([]epCase{
		{name: "CreateInteractionResponse", method: "POST", path: "/interactions/" + testInteractionID.String() + "/tok/callback", call: func() error {
			return drop(s.client.CreateInteractionResponse(ctx, testInteractionID, "tok", responses.InteractionResponse{}, true))
		}},
		{name: "GetOriginalInteractionResponse", method: "GET", path: "/webhooks/" + wh + "/tok/messages/@original", call: func() error {
			return drop(s.client.GetOriginalInteractionResponse(ctx, testWebhookID, "tok"))
		}},
		{name: "EditOriginalInteractionResponse", method: "PATCH", path: "/webhooks/" + wh + "/tok/messages/@original", call: func() error {
			return drop(s.client.EditOriginalInteractionResponse(ctx, testWebhookID, "tok", EditMessageParams{}))
		}},
		{name: "DeleteOriginalInteractionResponse", method: "DELETE", path: "/webhooks/" + wh + "/tok/messages/@original", call: func() error {
			return s.client.DeleteOriginalInteractionResponse(ctx, testWebhookID, "tok")
		}},
		{name: "CreateFollowupMessage", method: "POST", path: "/webhooks/" + app + "/tok", call: func() error {
			return drop(s.client.CreateFollowupMessage(ctx, testAppID, "tok", CreateMessageParams{}))
		}},
		{name: "GetFollowupMessage", method: "GET", path: "/webhooks/" + app + "/tok/messages/" + testMsgID.String(), call: func() error {
			return drop(s.client.GetFollowupMessage(ctx, testAppID, "tok", testMsgID))
		}},
		{name: "EditFollowupMessage", method: "PATCH", path: "/webhooks/" + app + "/tok/messages/" + testMsgID.String(), call: func() error {
			return drop(s.client.EditFollowupMessage(ctx, testAppID, "tok", testMsgID, EditMessageParams{}))
		}},
		{name: "DeleteFollowupMessage", method: "DELETE", path: "/webhooks/" + app + "/tok/messages/" + testMsgID.String(), call: func() error {
			return s.client.DeleteFollowupMessage(ctx, testAppID, "tok", testMsgID)
		}},
	})
}

func (s *endpointSuite) TestMiscEndpoints() {
	ctx := context.Background()
	g := testGuildID.String()
	app := testAppID.String()
	s.runEndpoints([]epCase{
		// templates
		{name: "GetTemplate", method: "GET", path: "/guilds/templates/code1", call: func() error {
			return drop(s.client.GetTemplate(ctx, "code1"))
		}},
		{name: "CreateGuildFromTemplate", method: "POST", path: "/guilds/templates/code1", status: http.StatusCreated, call: func() error {
			return drop(s.client.CreateGuildFromTemplate(ctx, "code1", CreateGuildFromTemplateParams{}))
		}},
		{name: "ListGuildTemplates", method: "GET", path: "/guilds/" + g + "/templates", body: "[]", call: func() error {
			return drop(s.client.ListGuildTemplates(ctx, testGuildID))
		}},
		{name: "CreateGuildTemplate", method: "POST", path: "/guilds/" + g + "/templates", status: http.StatusCreated, call: func() error {
			return drop(s.client.CreateGuildTemplate(ctx, testGuildID, CreateGuildTemplateParams{}))
		}},
		{name: "SyncGuildTemplate", method: "PUT", path: "/guilds/" + g + "/templates/code1", call: func() error {
			return drop(s.client.SyncGuildTemplate(ctx, testGuildID, "code1"))
		}},
		{name: "ModifyGuildTemplate", method: "PATCH", path: "/guilds/" + g + "/templates/code1", call: func() error {
			return drop(s.client.ModifyGuildTemplate(ctx, testGuildID, "code1", ModifyGuildTemplateParams{}))
		}},
		{name: "DeleteGuildTemplate", method: "DELETE", path: "/guilds/" + g + "/templates/code1", call: func() error {
			return drop(s.client.DeleteGuildTemplate(ctx, testGuildID, "code1"))
		}},
		// voice
		{name: "ListVoiceRegions", method: "GET", path: "/voice/regions", body: "[]", call: func() error {
			return drop(s.client.ListVoiceRegions(ctx))
		}},
		{name: "ListGuildVoiceRegions", method: "GET", path: "/guilds/" + g + "/regions", body: "[]", call: func() error {
			return drop(s.client.ListGuildVoiceRegions(ctx, testGuildID))
		}},
		{name: "ModifyCurrentUserVoiceState", method: "PATCH", path: "/guilds/" + g + "/voice-states/@me", call: func() error {
			return s.client.ModifyCurrentUserVoiceState(ctx, testGuildID, ModifyCurrentUserVoiceStateParams{})
		}},
		{name: "ModifyUserVoiceState", method: "PATCH", path: "/guilds/" + g + "/voice-states/" + testUserID.String(), call: func() error {
			return s.client.ModifyUserVoiceState(ctx, testGuildID, testUserID, ModifyUserVoiceStateParams{})
		}},
		{name: "GetCurrentUserVoiceState", method: "GET", path: "/guilds/" + g + "/voice-states/@me", call: func() error {
			return drop(s.client.GetCurrentUserVoiceState(ctx, testGuildID))
		}},
		{name: "GetUserVoiceState", method: "GET", path: "/guilds/" + g + "/voice-states/" + testUserID.String(), call: func() error {
			return drop(s.client.GetUserVoiceState(ctx, testGuildID, testUserID))
		}},
		// users
		{name: "GetCurrentUser", method: "GET", path: "/users/@me", call: func() error {
			return drop(s.client.GetCurrentUser(ctx, nil))
		}},
		{name: "GetUser", method: "GET", path: "/users/" + testUserID.String(), call: func() error {
			return drop(s.client.GetUser(ctx, testUserID))
		}},
		{name: "ModifyCurrentUser", method: "PATCH", path: "/users/@me", call: func() error {
			return drop(s.client.ModifyCurrentUser(ctx, ModifyCurrentUserParams{}))
		}},
		{name: "ListCurrentUserGuilds", method: "GET", path: "/users/@me/guilds", body: "[]", call: func() error {
			return drop(s.client.ListCurrentUserGuilds(ctx, GetCurrentUserGuildsParams{}))
		}},
		{name: "GetCurrentUserGuildMember", method: "GET", path: "/users/@me/guilds/" + g + "/member", call: func() error {
			return drop(s.client.GetCurrentUserGuildMember(ctx, testGuildID, nil))
		}},
		{name: "LeaveGuild", method: "DELETE", path: "/users/@me/guilds/" + g, call: func() error {
			return s.client.LeaveGuild(ctx, testGuildID)
		}},
		{name: "CreateDM", method: "POST", path: "/users/@me/channels", call: func() error {
			return drop(s.client.CreateDM(ctx, testUserID))
		}},
		{name: "CreateGroupDM", method: "POST", path: "/users/@me/channels", call: func() error {
			return drop(s.client.CreateGroupDM(ctx, CreateGroupDMParams{}))
		}},
		// invites
		{name: "GetInvite", method: "GET", path: "/invites/inv1", call: func() error {
			return drop(s.client.GetInvite(ctx, "inv1", GetInviteParams{}))
		}},
		{name: "DeleteInvite", method: "DELETE", path: "/invites/inv1", call: func() error {
			return drop(s.client.DeleteInvite(ctx, "inv1", nil))
		}},
		// oauth2
		{name: "GetCurrentAuthorizationInformation", method: "GET", path: "/oauth2/@me", call: func() error {
			return drop(s.client.GetCurrentAuthorizationInformation(ctx, "utok"))
		}},
		{name: "DeleteCurrentUserApplicationRoleConnection", method: "DELETE", path: "/users/@me/applications/" + app + "/role-connection", call: func() error {
			return s.client.DeleteCurrentUserApplicationRoleConnection(ctx, testAppID, "utok")
		}},
		{name: "ListCurrentUserConnections", method: "GET", path: "/users/@me/connections", body: "[]", call: func() error {
			return drop(s.client.ListCurrentUserConnections(ctx, "utok"))
		}},
		{name: "GetCurrentUserApplicationRoleConnection", method: "GET", path: "/users/@me/applications/" + app + "/role-connection", call: func() error {
			return drop(s.client.GetCurrentUserApplicationRoleConnection(ctx, testAppID, "utok"))
		}},
		{name: "UpdateCurrentUserApplicationRoleConnection", method: "PUT", path: "/users/@me/applications/" + app + "/role-connection", call: func() error {
			return drop(s.client.UpdateCurrentUserApplicationRoleConnection(ctx, testAppID, "utok", UpdateRoleConnectionParams{}))
		}},
		{name: "AddGuildMember", method: "PUT", path: "/guilds/" + g + "/members/" + testUserID.String(), status: http.StatusCreated, call: func() error {
			return drop(s.client.AddGuildMember(ctx, testGuildID, testUserID, AddGuildMemberParams{}))
		}},
		// entitlements
		{name: "ListEntitlements", method: "GET", path: "/applications/" + app + "/entitlements", body: "[]", call: func() error {
			return drop(s.client.ListEntitlements(ctx, testAppID, ListEntitlementsParams{}))
		}},
		{name: "GetEntitlement", method: "GET", path: "/applications/" + app + "/entitlements/" + testEntID.String(), call: func() error {
			return drop(s.client.GetEntitlement(ctx, testAppID, testEntID))
		}},
		{name: "CreateTestEntitlement", method: "POST", path: "/applications/" + app + "/entitlements", call: func() error {
			return drop(s.client.CreateTestEntitlement(ctx, testAppID, CreateTestEntitlementParams{}))
		}},
		{name: "ConsumeEntitlement", method: "POST", path: "/applications/" + app + "/entitlements/" + testEntID.String() + "/consume", call: func() error {
			return s.client.ConsumeEntitlement(ctx, testAppID, testEntID)
		}},
		{name: "ListCurrentUserApplicationEntitlements", method: "GET", path: "/users/@me/applications/" + app + "/entitlements", body: "[]", call: func() error {
			return drop(s.client.ListCurrentUserApplicationEntitlements(ctx, testAppID, GetCurrentUserApplicationEntitlementsParams{}, "utok"))
		}},
		{name: "DeleteTestEntitlement", method: "DELETE", path: "/applications/" + app + "/entitlements/" + testEntID.String(), call: func() error {
			return s.client.DeleteTestEntitlement(ctx, testAppID, testEntID)
		}},
		// automoderation
		{name: "ListAutoModerationRules", method: "GET", path: "/guilds/" + g + "/auto-moderation/rules", body: "[]", call: func() error {
			return drop(s.client.ListAutoModerationRules(ctx, testGuildID))
		}},
		{name: "GetAutoModerationRule", method: "GET", path: "/guilds/" + g + "/auto-moderation/rules/" + testRuleID.String(), call: func() error {
			return drop(s.client.GetAutoModerationRule(ctx, testGuildID, testRuleID))
		}},
		{name: "CreateAutoModerationRule", method: "POST", path: "/guilds/" + g + "/auto-moderation/rules", call: func() error {
			return drop(s.client.CreateAutoModerationRule(ctx, testGuildID, CreateAutoModerationRuleParams{}, nil))
		}},
		{name: "ModifyAutoModerationRule", method: "PATCH", path: "/guilds/" + g + "/auto-moderation/rules/" + testRuleID.String(), call: func() error {
			return drop(s.client.ModifyAutoModerationRule(ctx, testGuildID, testRuleID, ModifyAutoModerationRuleParams{}, nil))
		}},
		{name: "DeleteAutoModerationRule", method: "DELETE", path: "/guilds/" + g + "/auto-moderation/rules/" + testRuleID.String(), call: func() error {
			return s.client.DeleteAutoModerationRule(ctx, testGuildID, testRuleID, nil)
		}},
		// stage instances
		{name: "CreateStageInstance", method: "POST", path: "/stage-instances", status: http.StatusCreated, call: func() error {
			return drop(s.client.CreateStageInstance(ctx, CreateStageInstanceParams{ChannelID: testChanID}, nil))
		}},
		{name: "GetStageInstance", method: "GET", path: "/stage-instances/" + testChanID.String(), call: func() error {
			return drop(s.client.GetStageInstance(ctx, testChanID))
		}},
		{name: "ModifyStageInstance", method: "PATCH", path: "/stage-instances/" + testChanID.String(), call: func() error {
			return drop(s.client.ModifyStageInstance(ctx, testChanID, ModifyStageInstanceParams{}, nil))
		}},
		{name: "DeleteStageInstance", method: "DELETE", path: "/stage-instances/" + testChanID.String(), call: func() error {
			return s.client.DeleteStageInstance(ctx, testChanID, nil)
		}},
		// subscriptions
		{name: "ListSKUSubscriptions", method: "GET", path: "/skus/" + testSkuID.String() + "/subscriptions", body: "[]", call: func() error {
			return drop(s.client.ListSKUSubscriptions(ctx, testSkuID, ListSKUSubscriptionsParams{}))
		}},
		{name: "GetSKUSubscription", method: "GET", path: "/skus/" + testSkuID.String() + "/subscriptions/" + testSubID.String(), call: func() error {
			return drop(s.client.GetSKUSubscription(ctx, testSkuID, testSubID))
		}},
		// stickers
		{name: "GetSticker", method: "GET", path: "/stickers/" + testStickerID.String(), call: func() error {
			return drop(s.client.GetSticker(ctx, testStickerID))
		}},
		{name: "ListStickerPacks", method: "GET", path: "/sticker-packs", call: func() error {
			return drop(s.client.ListStickerPacks(ctx))
		}},
		{name: "ListGuildStickers", method: "GET", path: "/guilds/" + g + "/stickers", body: "[]", call: func() error {
			return drop(s.client.ListGuildStickers(ctx, testGuildID))
		}},
		{name: "GetGuildSticker", method: "GET", path: "/guilds/" + g + "/stickers/" + testStickerID.String(), call: func() error {
			return drop(s.client.GetGuildSticker(ctx, testGuildID, testStickerID))
		}},
		{name: "ModifyGuildSticker", method: "PATCH", path: "/guilds/" + g + "/stickers/" + testStickerID.String(), call: func() error {
			return drop(s.client.ModifyGuildSticker(ctx, testGuildID, testStickerID, ModifyGuildStickerParams{}, nil))
		}},
		{name: "DeleteGuildSticker", method: "DELETE", path: "/guilds/" + g + "/stickers/" + testStickerID.String(), call: func() error {
			return s.client.DeleteGuildSticker(ctx, testGuildID, testStickerID, nil)
		}},
		{name: "GetStickerPack", method: "GET", path: "/sticker-packs/" + testPackID.String(), call: func() error {
			return drop(s.client.GetStickerPack(ctx, testPackID))
		}},
		{name: "CreateGuildSticker", method: "POST", path: "/guilds/" + g + "/stickers", status: http.StatusCreated, call: func() error {
			return drop(s.client.CreateGuildSticker(ctx, testGuildID, CreateGuildStickerParams{
				Name: "s", Description: "d", Tags: "t", File: []byte("png"), ContentType: "image/png",
			}, nil))
		}},
		// soundboard
		{name: "ListDefaultSoundboardSounds", method: "GET", path: "/soundboard-default-sounds", body: "[]", call: func() error {
			return drop(s.client.ListDefaultSoundboardSounds(ctx))
		}},
		{name: "ListGuildSoundboardSounds", method: "GET", path: "/guilds/" + g + "/soundboard-sounds", call: func() error {
			return drop(s.client.ListGuildSoundboardSounds(ctx, testGuildID))
		}},
		{name: "GetGuildSoundboardSound", method: "GET", path: "/guilds/" + g + "/soundboard-sounds/" + testSoundID.String(), call: func() error {
			return drop(s.client.GetGuildSoundboardSound(ctx, testGuildID, testSoundID))
		}},
		{name: "CreateGuildSoundboardSound", method: "POST", path: "/guilds/" + g + "/soundboard-sounds", call: func() error {
			return drop(s.client.CreateGuildSoundboardSound(ctx, testGuildID, CreateGuildSoundboardSoundParams{}, nil))
		}},
		{name: "ModifyGuildSoundboardSound", method: "PATCH", path: "/guilds/" + g + "/soundboard-sounds/" + testSoundID.String(), call: func() error {
			return drop(s.client.ModifyGuildSoundboardSound(ctx, testGuildID, testSoundID, ModifyGuildSoundboardSoundParams{}, nil))
		}},
		{name: "DeleteGuildSoundboardSound", method: "DELETE", path: "/guilds/" + g + "/soundboard-sounds/" + testSoundID.String(), call: func() error {
			return s.client.DeleteGuildSoundboardSound(ctx, testGuildID, testSoundID, nil)
		}},
		{name: "SendSoundboardSound", method: "POST", path: "/channels/" + testChanID.String() + "/send-soundboard-sound", call: func() error {
			return s.client.SendSoundboardSound(ctx, testChanID, SendSoundboardSoundParams{})
		}},
		// application
		{name: "GetCurrentApplication", method: "GET", path: "/applications/@me", call: func() error {
			return drop(s.client.GetCurrentApplication(ctx))
		}},
		{name: "ListApplicationRoleConnectionsMetadata", method: "GET", path: "/applications/" + app + "/role-connections/metadata", body: "[]", call: func() error {
			return drop(s.client.ListApplicationRoleConnectionsMetadata(ctx, testAppID))
		}},
		{name: "UpdateApplicationRoleConnectionsMetadata", method: "PUT", path: "/applications/" + app + "/role-connections/metadata", body: "[]", call: func() error {
			return drop(s.client.UpdateApplicationRoleConnectionsMetadata(ctx, testAppID, nil))
		}},
		{name: "ModifyCurrentApplication", method: "PATCH", path: "/applications/@me", call: func() error {
			return drop(s.client.ModifyCurrentApplication(ctx, ModifyCurrentApplicationParams{}))
		}},
		// activity & skus
		{name: "GetActivityInstance", method: "GET", path: "/applications/" + app + "/activity-instances/inst1", call: func() error {
			return drop(s.client.GetActivityInstance(ctx, testAppID, "inst1"))
		}},
		{name: "ListSKUs", method: "GET", path: "/applications/" + app + "/skus", body: "[]", call: func() error {
			return drop(s.client.ListSKUs(ctx, testAppID))
		}},
	})
}
