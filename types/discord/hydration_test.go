package discord_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Stub EntityClient ─────────────────────────────────────────────────────────

// stubClient implements discord.EntityClient with zero-value returns.
// Override specific fields on a spy wrapper when you need recorded calls.
type stubClient struct {
	editMessageFn   func(ctx context.Context, channelID, messageID discord.Snowflake, opts discord.MessageEditOptions) (*discord.Message, error)
	deleteMessageFn func(ctx context.Context, channelID, messageID discord.Snowflake, reason *string) error
	createMessageFn func(ctx context.Context, channelID discord.Snowflake, opts discord.MessageCreateOptions) (*discord.Message, error)
	modifyGuildFn   func(ctx context.Context, guildID discord.Snowflake, opts discord.GuildEditOptions) (*discord.Guild, error)
	leaveGuildFn    func(ctx context.Context, guildID discord.Snowflake) error
	modifyRoleFn    func(ctx context.Context, guildID, roleID discord.Snowflake, opts discord.RoleEditOptions) (*discord.Role, error)
	deleteRoleFn    func(ctx context.Context, guildID, roleID discord.Snowflake, reason *string) error
	modifyMemberFn  func(ctx context.Context, guildID, userID discord.Snowflake, opts discord.MemberEditOptions) (*discord.GuildMember, error)
	kickMemberFn    func(ctx context.Context, guildID, userID discord.Snowflake, reason *string) error
	modifyChannelFn func(ctx context.Context, channelID discord.Snowflake, opts discord.ChannelEditOptions) (*discord.Channel, error)
	deleteChannelFn func(ctx context.Context, channelID discord.Snowflake, reason *string) (*discord.Channel, error)
}

func (s *stubClient) EditMessage(ctx context.Context, chID, msgID discord.Snowflake, opts discord.MessageEditOptions) (*discord.Message, error) {
	if s.editMessageFn != nil {
		return s.editMessageFn(ctx, chID, msgID, opts)
	}
	return nil, nil
}
func (s *stubClient) DeleteMessage(ctx context.Context, chID, msgID discord.Snowflake, reason *string) error {
	if s.deleteMessageFn != nil {
		return s.deleteMessageFn(ctx, chID, msgID, reason)
	}
	return nil
}
func (s *stubClient) CreateMessage(ctx context.Context, chID discord.Snowflake, opts discord.MessageCreateOptions) (*discord.Message, error) {
	if s.createMessageFn != nil {
		return s.createMessageFn(ctx, chID, opts)
	}
	return nil, nil
}
func (s *stubClient) PinMessage(ctx context.Context, chID, msgID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubClient) UnpinMessage(ctx context.Context, chID, msgID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubClient) CrosspostMessage(ctx context.Context, chID, msgID discord.Snowflake) (*discord.Message, error) {
	return nil, nil
}
func (s *stubClient) AddReaction(ctx context.Context, chID, msgID discord.Snowflake, emoji string) error {
	return nil
}
func (s *stubClient) ModifyChannel(ctx context.Context, chID discord.Snowflake, opts discord.ChannelEditOptions) (*discord.Channel, error) {
	if s.modifyChannelFn != nil {
		return s.modifyChannelFn(ctx, chID, opts)
	}
	return nil, nil
}
func (s *stubClient) DeleteChannel(ctx context.Context, chID discord.Snowflake, reason *string) (*discord.Channel, error) {
	if s.deleteChannelFn != nil {
		return s.deleteChannelFn(ctx, chID, reason)
	}
	return nil, nil
}
func (s *stubClient) BulkDeleteMessages(ctx context.Context, chID discord.Snowflake, ids []discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubClient) ListChannelMessages(ctx context.Context, chID discord.Snowflake, opts discord.FetchMessagesOptions) ([]*discord.Message, error) {
	return nil, nil
}
func (s *stubClient) TriggerTypingIndicator(ctx context.Context, chID discord.Snowflake) error {
	return nil
}
func (s *stubClient) SetVoiceChannelStatus(ctx context.Context, chID discord.Snowflake, status *string) error {
	return nil
}
func (s *stubClient) CreateChannelInvite(ctx context.Context, chID discord.Snowflake, opts discord.InviteCreateOptions) (*discord.Invite, error) {
	return nil, nil
}
func (s *stubClient) ModifyGuild(ctx context.Context, guildID discord.Snowflake, opts discord.GuildEditOptions) (*discord.Guild, error) {
	if s.modifyGuildFn != nil {
		return s.modifyGuildFn(ctx, guildID, opts)
	}
	return nil, nil
}
func (s *stubClient) LeaveGuild(ctx context.Context, guildID discord.Snowflake) error {
	if s.leaveGuildFn != nil {
		return s.leaveGuildFn(ctx, guildID)
	}
	return nil
}
func (s *stubClient) ListGuildMembers(ctx context.Context, guildID discord.Snowflake, opts discord.FetchMembersOptions) ([]*discord.GuildMember, error) {
	return nil, nil
}
func (s *stubClient) CreateGuildRole(ctx context.Context, guildID discord.Snowflake, opts discord.RoleCreateOptions) (*discord.Role, error) {
	return nil, nil
}
func (s *stubClient) CreateGuildChannel(ctx context.Context, guildID discord.Snowflake, opts discord.ChannelCreateOptions) (*discord.Channel, error) {
	return nil, nil
}
func (s *stubClient) CreateGuildEmoji(ctx context.Context, guildID discord.Snowflake, opts discord.EmojiCreateOptions) (*discord.Emoji, error) {
	return nil, nil
}
func (s *stubClient) CreateGuildSticker(ctx context.Context, guildID discord.Snowflake, opts discord.StickerCreateOptions) (*discord.Sticker, error) {
	return nil, nil
}
func (s *stubClient) GetGuildAuditLog(ctx context.Context, guildID discord.Snowflake, opts discord.AuditLogOptions) (*discord.AuditLog, error) {
	return nil, nil
}
func (s *stubClient) CreateGuildScheduledEvent(ctx context.Context, guildID discord.Snowflake, opts discord.ScheduledEventCreateOptions) (*discord.GuildScheduledEvent, error) {
	return nil, nil
}
func (s *stubClient) ModifyGuildMember(ctx context.Context, guildID, userID discord.Snowflake, opts discord.MemberEditOptions) (*discord.GuildMember, error) {
	if s.modifyMemberFn != nil {
		return s.modifyMemberFn(ctx, guildID, userID, opts)
	}
	return nil, nil
}
func (s *stubClient) KickGuildMember(ctx context.Context, guildID, userID discord.Snowflake, reason *string) error {
	if s.kickMemberFn != nil {
		return s.kickMemberFn(ctx, guildID, userID, reason)
	}
	return nil
}
func (s *stubClient) CreateGuildBan(ctx context.Context, guildID, userID discord.Snowflake, opts discord.BanOptions) error {
	return nil
}
func (s *stubClient) AddGuildMemberRole(ctx context.Context, guildID, userID, roleID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubClient) RemoveGuildMemberRole(ctx context.Context, guildID, userID, roleID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubClient) ModifyGuildRole(ctx context.Context, guildID, roleID discord.Snowflake, opts discord.RoleEditOptions) (*discord.Role, error) {
	if s.modifyRoleFn != nil {
		return s.modifyRoleFn(ctx, guildID, roleID, opts)
	}
	return nil, nil
}
func (s *stubClient) DeleteGuildRole(ctx context.Context, guildID, roleID discord.Snowflake, reason *string) error {
	if s.deleteRoleFn != nil {
		return s.deleteRoleFn(ctx, guildID, roleID, reason)
	}
	return nil
}
func (s *stubClient) ModifyGuildRolePositions(ctx context.Context, guildID discord.Snowflake, opts discord.RolePositionOptions) ([]*discord.Role, error) {
	return nil, nil
}
func (s *stubClient) GetUser(ctx context.Context, userID discord.Snowflake) (*discord.User, error) {
	return nil, nil
}
func (s *stubClient) CreateDM(ctx context.Context, recipientID discord.Snowflake) (*discord.Channel, error) {
	return nil, nil
}
func (s *stubClient) ModifyGuildEmoji(ctx context.Context, guildID, emojiID discord.Snowflake, opts discord.EmojiEditOptions) (*discord.Emoji, error) {
	return nil, nil
}
func (s *stubClient) DeleteGuildEmoji(ctx context.Context, guildID, emojiID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubClient) ModifyWebhook(ctx context.Context, webhookID discord.Snowflake, opts discord.WebhookEditOptions) (*discord.Webhook, error) {
	return nil, nil
}
func (s *stubClient) DeleteWebhook(ctx context.Context, webhookID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubClient) ExecuteWebhook(ctx context.Context, webhookID discord.Snowflake, token string, opts discord.WebhookExecuteOptions) (*discord.Message, error) {
	return nil, nil
}
func (s *stubClient) GetWebhookMessage(ctx context.Context, webhookID discord.Snowflake, token string, messageID discord.Snowflake) (*discord.Message, error) {
	return nil, nil
}
func (s *stubClient) DeleteInvite(ctx context.Context, code string, reason *string) (*discord.Invite, error) {
	return nil, nil
}
func (s *stubClient) ModifyStageInstance(ctx context.Context, chID discord.Snowflake, opts discord.StageEditOptions) (*discord.StageInstance, error) {
	return nil, nil
}
func (s *stubClient) DeleteStageInstance(ctx context.Context, chID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubClient) ModifyGuildScheduledEvent(ctx context.Context, guildID, eventID discord.Snowflake, opts discord.ScheduledEventEditOptions) (*discord.GuildScheduledEvent, error) {
	return nil, nil
}
func (s *stubClient) DeleteGuildScheduledEvent(ctx context.Context, guildID, eventID discord.Snowflake) error {
	return nil
}
func (s *stubClient) ListGuildScheduledEventUsers(ctx context.Context, guildID, eventID discord.Snowflake, opts discord.FetchUsersOptions) ([]*discord.GuildScheduledEventUser, error) {
	return nil, nil
}
func (s *stubClient) ModifyGuildSticker(ctx context.Context, guildID, stickerID discord.Snowflake, opts discord.StickerEditOptions) (*discord.Sticker, error) {
	return nil, nil
}
func (s *stubClient) DeleteGuildSticker(ctx context.Context, guildID, stickerID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubClient) ModifyAutoModerationRule(ctx context.Context, guildID, ruleID discord.Snowflake, opts discord.RuleEditOptions) (*discord.AutoModerationRule, error) {
	return nil, nil
}
func (s *stubClient) DeleteAutoModerationRule(ctx context.Context, guildID, ruleID discord.Snowflake) error {
	return nil
}
func (s *stubClient) ModifyGuildSoundboardSound(ctx context.Context, guildID, soundID discord.Snowflake, opts discord.SoundEditOptions) (*discord.SoundboardSound, error) {
	return nil, nil
}
func (s *stubClient) DeleteGuildSoundboardSound(ctx context.Context, guildID, soundID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubClient) DeleteGuildIntegration(ctx context.Context, guildID, integrationID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubClient) ListGuildIntegrations(ctx context.Context, guildID discord.Snowflake) ([]*discord.Integration, error) {
	return nil, nil
}
func (s *stubClient) GetGuildMember(ctx context.Context, guildID, userID discord.Snowflake) (*discord.GuildMember, error) {
	return nil, nil
}
func (s *stubClient) GetGuildRole(ctx context.Context, guildID, roleID discord.Snowflake) (*discord.Role, error) {
	return nil, nil
}
func (s *stubClient) ListGuildRoles(ctx context.Context, guildID discord.Snowflake) ([]*discord.Role, error) {
	return nil, nil
}
func (s *stubClient) GetChannel(ctx context.Context, channelID discord.Snowflake) (*discord.Channel, error) {
	return nil, nil
}
func (s *stubClient) ListGuildChannels(ctx context.Context, guildID discord.Snowflake) ([]*discord.Channel, error) {
	return nil, nil
}
func (s *stubClient) GetGuildEmoji(ctx context.Context, guildID, emojiID discord.Snowflake) (*discord.Emoji, error) {
	return nil, nil
}
func (s *stubClient) ListGuildEmojis(ctx context.Context, guildID discord.Snowflake) ([]*discord.Emoji, error) {
	return nil, nil
}
func (s *stubClient) GetGuildSticker(ctx context.Context, guildID, stickerID discord.Snowflake) (*discord.Sticker, error) {
	return nil, nil
}
func (s *stubClient) ListGuildStickers(ctx context.Context, guildID discord.Snowflake) ([]*discord.Sticker, error) {
	return nil, nil
}
func (s *stubClient) GetGuildBan(ctx context.Context, guildID, userID discord.Snowflake) (*discord.Ban, error) {
	return nil, nil
}
func (s *stubClient) ListGuildBans(ctx context.Context, guildID discord.Snowflake, opts discord.FetchBansOptions) ([]*discord.Ban, error) {
	return nil, nil
}
func (s *stubClient) RemoveGuildBan(ctx context.Context, guildID, userID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubClient) GetGuildScheduledEvent(ctx context.Context, guildID, eventID discord.Snowflake) (*discord.GuildScheduledEvent, error) {
	return nil, nil
}
func (s *stubClient) ListGuildScheduledEvents(ctx context.Context, guildID discord.Snowflake) ([]*discord.GuildScheduledEvent, error) {
	return nil, nil
}
func (s *stubClient) GetStageInstance(ctx context.Context, channelID discord.Snowflake) (*discord.StageInstance, error) {
	return nil, nil
}
func (s *stubClient) CreateStageInstance(ctx context.Context, opts discord.StageCreateOptions) (*discord.StageInstance, error) {
	return nil, nil
}
func (s *stubClient) GetGuildSoundboardSound(ctx context.Context, guildID, soundID discord.Snowflake) (*discord.SoundboardSound, error) {
	return nil, nil
}
func (s *stubClient) ListGuildSoundboardSounds(ctx context.Context, guildID discord.Snowflake) ([]*discord.SoundboardSound, error) {
	return nil, nil
}
func (s *stubClient) ListGuildInvites(ctx context.Context, guildID discord.Snowflake) ([]*discord.Invite, error) {
	return nil, nil
}
func (s *stubClient) ListGuildWebhooks(ctx context.Context, guildID discord.Snowflake) ([]*discord.Webhook, error) {
	return nil, nil
}
func (s *stubClient) GetAutoModerationRule(ctx context.Context, guildID, ruleID discord.Snowflake) (*discord.AutoModerationRule, error) {
	return nil, nil
}
func (s *stubClient) ListAutoModerationRules(ctx context.Context, guildID discord.Snowflake) ([]*discord.AutoModerationRule, error) {
	return nil, nil
}
func (s *stubClient) CreateAutoModerationRule(ctx context.Context, guildID discord.Snowflake, opts discord.RuleCreateOptions) (*discord.AutoModerationRule, error) {
	return nil, nil
}
func (s *stubClient) GetMessage(ctx context.Context, channelID, messageID discord.Snowflake) (*discord.Message, error) {
	return nil, nil
}
func (s *stubClient) GetGuild(ctx context.Context, guildID discord.Snowflake) (*discord.Guild, error) {
	return nil, nil
}
func (s *stubClient) ClientCache() discord.Cache {
	return nil
}
func (s *stubClient) ThreadsForParent(parentID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Channel] {
	return collection.New[discord.Snowflake, *discord.Channel]()
}

func (s *stubClient) ChannelsForGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Channel] {
	return collection.New[discord.Snowflake, *discord.Channel]()
}

var _ discord.EntityClient = (*stubClient)(nil)

// ── Infrastructure tests ──────────────────────────────────────────────────────

func (su *hydrationSuite) TestErrEntityNotHydrated_MessageEdit() {
	t := su.T()
	var msg discord.Message
	_, err := msg.Edit(context.Background(), discord.MessageEditOptions{})
	if !errors.Is(err, discord.ErrEntityNotHydrated) {
		t.Fatalf("expected ErrEntityNotHydrated, got %v", err)
	}
}

func (su *hydrationSuite) TestErrEntityNotHydrated_MessageDelete() {
	t := su.T()
	var msg discord.Message
	err := msg.Delete(context.Background(), nil)
	if !errors.Is(err, discord.ErrEntityNotHydrated) {
		t.Fatalf("expected ErrEntityNotHydrated, got %v", err)
	}
}

func (su *hydrationSuite) TestErrEntityNotHydrated_GuildLeave() {
	t := su.T()
	var g discord.Guild
	err := g.Leave(context.Background())
	if !errors.Is(err, discord.ErrEntityNotHydrated) {
		t.Fatalf("expected ErrEntityNotHydrated, got %v", err)
	}
}

func (su *hydrationSuite) TestErrEntityNotHydrated_RoleDelete() {
	t := su.T()
	var r discord.Role
	err := r.Delete(context.Background(), nil)
	if !errors.Is(err, discord.ErrEntityNotHydrated) {
		t.Fatalf("expected ErrEntityNotHydrated, got %v", err)
	}
}

func (su *hydrationSuite) TestErrEntityNotHydrated_ChannelEdit() {
	t := su.T()
	var ch discord.Channel
	_, err := ch.Edit(context.Background(), discord.ChannelEditOptions{})
	if !errors.Is(err, discord.ErrEntityNotHydrated) {
		t.Fatalf("expected ErrEntityNotHydrated, got %v", err)
	}
}

func (su *hydrationSuite) TestErrEntityNotHydrated_MemberKick() {
	t := su.T()
	m := &discord.GuildMember{}
	err := m.Kick(context.Background(), nil)
	if !errors.Is(err, discord.ErrEntityNotHydrated) {
		t.Fatalf("expected ErrEntityNotHydrated, got %v", err)
	}
}

// ── IsHydrated ────────────────────────────────────────────────────────────────

func (su *hydrationSuite) TestIsHydrated_FalseBeforeHydrate() {
	t := su.T()
	var msg discord.Message
	if msg.IsHydrated() {
		t.Fatal("expected IsHydrated to be false before Hydrate call")
	}
}

func (su *hydrationSuite) TestIsHydrated_TrueAfterHydrate() {
	t := su.T()
	msg := &discord.Message{}
	msg.Hydrate(&stubClient{})
	if !msg.IsHydrated() {
		t.Fatal("expected IsHydrated to be true after Hydrate call")
	}
}

func (su *hydrationSuite) TestWithClient_ReturnsHydratedCopy() {
	t := su.T()
	c := &stubClient{}
	original := &discord.Message{}
	if original.IsHydrated() {
		t.Fatal("original should not be hydrated")
	}

	copy := original.WithClient(c)
	if !copy.IsHydrated() {
		t.Fatal("copy should be hydrated")
	}
	if original.IsHydrated() {
		t.Fatal("original must not be mutated by WithClient")
	}
}

// ── JSON round-trip: hClient field must not be serialized ─────────────────────

func (su *hydrationSuite) TestHydrate_JSONOmitsClient() {
	t := su.T()
	msg := &discord.Message{
		ID:        123,
		ChannelID: 456,
	}
	msg.Hydrate(&stubClient{})

	b, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	if !msg.IsHydrated() {
		t.Fatal("message should be hydrated before marshal")
	}

	var out discord.Message
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if out.IsHydrated() {
		t.Fatal("deserialized message must not be hydrated (hClient is pkg)")
	}
	if out.ID != 123 || out.ChannelID != 456 {
		t.Fatalf("JSON round-trip lost fields: %+v", out)
	}
}

// ── Hydrate propagates to nested entities ─────────────────────────────────────

func (su *hydrationSuite) TestHydrate_PropagatesAuthorHydration() {
	t := su.T()
	c := &stubClient{}
	author := &discord.User{ID: 9}
	msg := &discord.Message{Author: author}
	msg.Hydrate(c)

	if !author.IsHydrated() {
		t.Fatal("author should be hydrated after Message.Hydrate")
	}
}

func (su *hydrationSuite) TestGuildHydrate_PropagatesRolesAndEmojis() {
	t := su.T()
	c := &stubClient{}
	g := &discord.Guild{
		ID:        1,
		RawRoles:  []discord.Role{{ID: 12}, {ID: 123}},
		RawEmojis: []discord.Emoji{{ID: 1234}},
	}
	g.Hydrate(c)

	if !g.IsHydrated() {
		t.Fatal("guild should be hydrated")
	}
	for _, r := range g.RawRoles {
		if !r.IsHydrated() {
			t.Fatalf("role %s should be hydrated after Guild.Hydrate", r.ID.String())
		}
	}
	for _, e := range g.RawEmojis {
		if !e.IsHydrated() {
			t.Fatalf("emoji %s should be hydrated after Guild.Hydrate", e.ID.String())
		}
	}
}

// ── Smoke tests: method routes to correct EntityClient call ───────────────────

func (su *hydrationSuite) TestMessageEdit_Smoke() {
	t := su.T()
	want := &discord.Message{ID: 1}
	c := &stubClient{
		editMessageFn: func(_ context.Context, chID, msgID discord.Snowflake, opts discord.MessageEditOptions) (*discord.Message, error) {
			if chID != 11 || msgID != 12 {
				t.Errorf("unexpected IDs: channelID=%s messageID=%s", chID.String(), msgID.String())
			}
			return want, nil
		},
	}
	msg := &discord.Message{ID: 12, ChannelID: 11}
	msg.Hydrate(c)

	got, err := msg.Edit(context.Background(), discord.MessageEditOptions{})
	if err != nil || got != want {
		t.Fatalf("Edit: err=%v got=%v want=%v", err, got, want)
	}
}

func (su *hydrationSuite) TestMessageDelete_Smoke() {
	t := su.T()
	var called bool
	c := &stubClient{
		deleteMessageFn: func(_ context.Context, chID, msgID discord.Snowflake, _ *string) error {
			called = true
			return nil
		},
	}
	msg := &discord.Message{ID: 11, ChannelID: 12}
	msg.Hydrate(c)

	if err := msg.Delete(context.Background(), nil); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if !called {
		t.Fatal("DeleteMessage was not called")
	}
}

func (su *hydrationSuite) TestGuildLeave_Smoke() {
	t := su.T()
	var gotID discord.Snowflake
	c := &stubClient{
		leaveGuildFn: func(_ context.Context, guildID discord.Snowflake) error {
			gotID = guildID
			return nil
		},
	}
	g := &discord.Guild{ID: 12}
	g.Hydrate(c)

	if err := g.Leave(context.Background()); err != nil {
		t.Fatalf("Leave: %v", err)
	}
	if gotID != 12 {
		t.Fatalf("expected guildID g1, got %s", gotID.String())
	}
}

func (su *hydrationSuite) TestRoleDelete_Smoke() {
	t := su.T()
	var called bool
	c := &stubClient{
		deleteRoleFn: func(_ context.Context, guildID, roleID discord.Snowflake, _ *string) error {
			if guildID != 12 || roleID != 11 {
				t.Errorf("unexpected IDs: guildID=%s roleID=%s", guildID.String(), roleID.String())
			}
			called = true
			return nil
		},
	}
	r := &discord.Role{ID: 11, GuildID: 12}
	r.Hydrate(c)

	if err := r.Delete(context.Background(), nil); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if !called {
		t.Fatal("DeleteGuildRole was not called")
	}
}

func (su *hydrationSuite) TestMemberKick_Smoke() {
	t := su.T()
	var called bool
	c := &stubClient{
		kickMemberFn: func(_ context.Context, guildID, userID discord.Snowflake, _ *string) error {
			if guildID != 11 || userID != 12 {
				t.Errorf("unexpected IDs: guildID=%s userID=%s", guildID.String(), userID.Validate())
			}
			called = true
			return nil
		},
	}
	m := &discord.GuildMember{GuildID: 11, UserID: 12}
	m.Hydrate(c)

	if err := m.Kick(context.Background(), nil); err != nil {
		t.Fatalf("Kick: %v", err)
	}
	if !called {
		t.Fatal("KickGuildMember was not called")
	}
}

func (su *hydrationSuite) TestChannelEdit_Smoke() {
	t := su.T()
	want := &discord.Channel{ID: 13}
	c := &stubClient{
		modifyChannelFn: func(_ context.Context, chID discord.Snowflake, opts discord.ChannelEditOptions) (*discord.Channel, error) {
			if chID != 13 {
				t.Errorf("unexpected channelID: %s", chID.String())
			}
			return want, nil
		},
	}
	ch := &discord.Channel{ID: 13}
	ch.Hydrate(c)

	got, err := ch.Edit(context.Background(), discord.ChannelEditOptions{})
	if err != nil || got != want {
		t.Fatalf("Edit: err=%v got=%v", err, got)
	}
}

type hydrationSuite struct{ suite.Suite }

func TestHydrationSuite(t *testing.T) { suite.Run(t, new(hydrationSuite)) }
