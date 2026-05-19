package discord_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Stub EntityClient ─────────────────────────────────────────────────────────

// stubClient implements discord.EntityClient with zero-value returns.
// Override specific fields on a spy wrapper when you need recorded calls.
type stubClient struct {
	editMessageFn func(ctx context.Context, channelID, messageID discord.Snowflake, opts discord.MessageEditOptions) (*discord.Message, error)
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
func (s *stubClient) GetChannelMessages(ctx context.Context, chID discord.Snowflake, opts discord.FetchMessagesOptions) ([]*discord.Message, error) {
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
func (s *stubClient) GetGuildScheduledEventUsers(ctx context.Context, guildID, eventID discord.Snowflake, opts discord.FetchUsersOptions) ([]*discord.GuildScheduledEventUser, error) {
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

var _ discord.EntityClient = (*stubClient)(nil)

// ── Infrastructure tests ──────────────────────────────────────────────────────

func TestErrEntityNotHydrated_MessageEdit(t *testing.T) {
	var msg discord.Message
	_, err := msg.Edit(context.Background(), discord.MessageEditOptions{})
	if !errors.Is(err, discord.ErrEntityNotHydrated) {
		t.Fatalf("expected ErrEntityNotHydrated, got %v", err)
	}
}

func TestErrEntityNotHydrated_MessageDelete(t *testing.T) {
	var msg discord.Message
	err := msg.Delete(context.Background(), nil)
	if !errors.Is(err, discord.ErrEntityNotHydrated) {
		t.Fatalf("expected ErrEntityNotHydrated, got %v", err)
	}
}

func TestErrEntityNotHydrated_GuildLeave(t *testing.T) {
	var g discord.Guild
	err := g.Leave(context.Background())
	if !errors.Is(err, discord.ErrEntityNotHydrated) {
		t.Fatalf("expected ErrEntityNotHydrated, got %v", err)
	}
}

func TestErrEntityNotHydrated_RoleDelete(t *testing.T) {
	var r discord.Role
	err := r.Delete(context.Background(), nil)
	if !errors.Is(err, discord.ErrEntityNotHydrated) {
		t.Fatalf("expected ErrEntityNotHydrated, got %v", err)
	}
}

func TestErrEntityNotHydrated_ChannelEdit(t *testing.T) {
	var ch discord.Channel
	_, err := ch.Edit(context.Background(), discord.ChannelEditOptions{})
	if !errors.Is(err, discord.ErrEntityNotHydrated) {
		t.Fatalf("expected ErrEntityNotHydrated, got %v", err)
	}
}

func TestErrEntityNotHydrated_MemberKick(t *testing.T) {
	m := &discord.GuildMember{}
	err := m.Kick(context.Background(), nil)
	if !errors.Is(err, discord.ErrEntityNotHydrated) {
		t.Fatalf("expected ErrEntityNotHydrated, got %v", err)
	}
}

// ── IsHydrated ────────────────────────────────────────────────────────────────

func TestIsHydrated_FalseBeforeHydrate(t *testing.T) {
	var msg discord.Message
	if msg.IsHydrated() {
		t.Fatal("expected IsHydrated to be false before Hydrate call")
	}
}

func TestIsHydrated_TrueAfterHydrate(t *testing.T) {
	msg := &discord.Message{}
	msg.Hydrate(&stubClient{})
	if !msg.IsHydrated() {
		t.Fatal("expected IsHydrated to be true after Hydrate call")
	}
}

func TestWithClient_ReturnsHydratedCopy(t *testing.T) {
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

func TestHydrate_JSONOmitsClient(t *testing.T) {
	msg := &discord.Message{
		ID:        "123",
		ChannelID: "456",
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
		t.Fatal("deserialized message must not be hydrated (hClient is internal)")
	}
	if out.ID != "123" || out.ChannelID != "456" {
		t.Fatalf("JSON round-trip lost fields: %+v", out)
	}
}

// ── Hydrate propagates to nested entities ─────────────────────────────────────

func TestHydrate_PropagatesAuthorHydration(t *testing.T) {
	c := &stubClient{}
	author := &discord.User{ID: "9"}
	msg := &discord.Message{Author: author}
	msg.Hydrate(c)

	if !author.IsHydrated() {
		t.Fatal("author should be hydrated after Message.Hydrate")
	}
}

func TestGuildHydrate_PropagatesRolesAndEmojis(t *testing.T) {
	c := &stubClient{}
	g := &discord.Guild{
		ID:    "1",
		Roles: []discord.Role{{ID: "r1"}, {ID: "r2"}},
		Emojis: []discord.Emoji{{ID: "e1"}},
	}
	g.Hydrate(c)

	if !g.IsHydrated() {
		t.Fatal("guild should be hydrated")
	}
	for _, r := range g.Roles {
		if !r.IsHydrated() {
			t.Fatalf("role %s should be hydrated after Guild.Hydrate", r.ID)
		}
	}
	for _, e := range g.Emojis {
		if !e.IsHydrated() {
			t.Fatalf("emoji %s should be hydrated after Guild.Hydrate", e.ID)
		}
	}
}

// ── Smoke tests: method routes to correct EntityClient call ───────────────────

func TestMessageEdit_Smoke(t *testing.T) {
	want := &discord.Message{ID: "edited"}
	c := &stubClient{
		editMessageFn: func(_ context.Context, chID, msgID discord.Snowflake, opts discord.MessageEditOptions) (*discord.Message, error) {
			if chID != "ch1" || msgID != "msg1" {
				t.Errorf("unexpected IDs: channelID=%s messageID=%s", chID, msgID)
			}
			return want, nil
		},
	}
	msg := &discord.Message{ID: "msg1", ChannelID: "ch1"}
	msg.Hydrate(c)

	got, err := msg.Edit(context.Background(), discord.MessageEditOptions{})
	if err != nil || got != want {
		t.Fatalf("Edit: err=%v got=%v want=%v", err, got, want)
	}
}

func TestMessageDelete_Smoke(t *testing.T) {
	var called bool
	c := &stubClient{
		deleteMessageFn: func(_ context.Context, chID, msgID discord.Snowflake, _ *string) error {
			called = true
			return nil
		},
	}
	msg := &discord.Message{ID: "msg1", ChannelID: "ch1"}
	msg.Hydrate(c)

	if err := msg.Delete(context.Background(), nil); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if !called {
		t.Fatal("DeleteMessage was not called")
	}
}

func TestGuildLeave_Smoke(t *testing.T) {
	var gotID discord.Snowflake
	c := &stubClient{
		leaveGuildFn: func(_ context.Context, guildID discord.Snowflake) error {
			gotID = guildID
			return nil
		},
	}
	g := &discord.Guild{ID: "g1"}
	g.Hydrate(c)

	if err := g.Leave(context.Background()); err != nil {
		t.Fatalf("Leave: %v", err)
	}
	if gotID != "g1" {
		t.Fatalf("expected guildID g1, got %s", gotID)
	}
}

func TestRoleDelete_Smoke(t *testing.T) {
	var called bool
	c := &stubClient{
		deleteRoleFn: func(_ context.Context, guildID, roleID discord.Snowflake, _ *string) error {
			if guildID != "g1" || roleID != "r1" {
				t.Errorf("unexpected IDs: guildID=%s roleID=%s", guildID, roleID)
			}
			called = true
			return nil
		},
	}
	r := &discord.Role{ID: "r1", GuildID: "g1"}
	r.Hydrate(c)

	if err := r.Delete(context.Background(), nil); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if !called {
		t.Fatal("DeleteGuildRole was not called")
	}
}

func TestMemberKick_Smoke(t *testing.T) {
	var called bool
	c := &stubClient{
		kickMemberFn: func(_ context.Context, guildID, userID discord.Snowflake, _ *string) error {
			if guildID != "g1" || userID != "u1" {
				t.Errorf("unexpected IDs: guildID=%s userID=%s", guildID, userID)
			}
			called = true
			return nil
		},
	}
	m := &discord.GuildMember{GuildID: "g1", UserID: "u1"}
	m.Hydrate(c)

	if err := m.Kick(context.Background(), nil); err != nil {
		t.Fatalf("Kick: %v", err)
	}
	if !called {
		t.Fatal("KickGuildMember was not called")
	}
}

func TestChannelEdit_Smoke(t *testing.T) {
	want := &discord.Channel{ID: "ch1"}
	c := &stubClient{
		modifyChannelFn: func(_ context.Context, chID discord.Snowflake, opts discord.ChannelEditOptions) (*discord.Channel, error) {
			if chID != "ch1" {
				t.Errorf("unexpected channelID: %s", chID)
			}
			return want, nil
		},
	}
	ch := &discord.Channel{ID: "ch1"}
	ch.Hydrate(c)

	got, err := ch.Edit(context.Background(), discord.ChannelEditOptions{})
	if err != nil || got != want {
		t.Fatalf("Edit: err=%v got=%v", err, got)
	}
}
