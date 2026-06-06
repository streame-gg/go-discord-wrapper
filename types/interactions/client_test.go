package interactions_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

// mockClient records the last call made to it and returns canned values, so the
// hydrated convenience methods on Interaction can be exercised without a live
// gateway connection.
type mockClient struct {
	lastMethod string
	lastCtx    context.Context
	lastData   *responses.InteractionResponseDataDefault
	lastParams interface{}
	lastID     discord.Snowflake
	ephemeral  bool
	err        error
}

func (m *mockClient) Reply(ctx context.Context, _ *interactions.Interaction, data *responses.InteractionResponseDataDefault, _ bool, _ ...discord.MessageFile) (*responses.InteractionCallbackResponse, error) {
	m.lastMethod, m.lastCtx, m.lastData = "Reply", ctx, data
	return &responses.InteractionCallbackResponse{}, m.err
}

func (m *mockClient) DeferReply(ctx context.Context, _ *interactions.Interaction, ephemeral bool) error {
	m.lastMethod, m.lastCtx, m.ephemeral = "DeferReply", ctx, ephemeral
	return m.err
}

func (m *mockClient) DeferUpdateMessage(ctx context.Context, _ *interactions.Interaction) error {
	m.lastMethod, m.lastCtx = "DeferUpdateMessage", ctx
	return m.err
}

func (m *mockClient) UpdateMessage(ctx context.Context, _ *interactions.Interaction, data *responses.InteractionResponseDataDefault, _ ...discord.MessageFile) error {
	m.lastMethod, m.lastCtx, m.lastData = "UpdateMessage", ctx, data
	return m.err
}

func (m *mockClient) ReplyWithModal(ctx context.Context, _ *interactions.Interaction, _ *components.Modal) error {
	m.lastMethod, m.lastCtx = "ReplyWithModal", ctx
	return m.err
}

func (m *mockClient) ReplyAutocomplete(ctx context.Context, _ *interactions.Interaction, _ []responses.AutocompleteChoice) error {
	m.lastMethod, m.lastCtx = "ReplyAutocomplete", ctx
	return m.err
}

func (m *mockClient) LaunchActivity(ctx context.Context, _ *interactions.Interaction) error {
	m.lastMethod, m.lastCtx = "LaunchActivity", ctx
	return m.err
}

func (m *mockClient) GetOriginalResponse(ctx context.Context, _ *interactions.Interaction) (*discord.Message, error) {
	m.lastMethod, m.lastCtx = "GetOriginalResponse", ctx
	return &discord.Message{ID: 1}, m.err
}

func (m *mockClient) EditReply(ctx context.Context, _ *interactions.Interaction, params api.EditMessageParams) (*discord.Message, error) {
	m.lastMethod, m.lastCtx, m.lastParams = "EditReply", ctx, params
	return &discord.Message{ID: 2}, m.err
}

func (m *mockClient) DeleteReply(ctx context.Context, _ *interactions.Interaction) error {
	m.lastMethod, m.lastCtx = "DeleteReply", ctx
	return m.err
}

func (m *mockClient) CreateFollowup(ctx context.Context, _ *interactions.Interaction, params api.CreateMessageParams) (*discord.Message, error) {
	m.lastMethod, m.lastCtx, m.lastParams = "CreateFollowup", ctx, params
	return &discord.Message{ID: 3}, m.err
}

func (m *mockClient) GetFollowup(ctx context.Context, _ *interactions.Interaction, messageID discord.Snowflake) (*discord.Message, error) {
	m.lastMethod, m.lastCtx, m.lastID = "GetFollowup", ctx, messageID
	return &discord.Message{ID: messageID}, m.err
}

func (m *mockClient) EditFollowup(ctx context.Context, _ *interactions.Interaction, messageID discord.Snowflake, params api.EditMessageParams) (*discord.Message, error) {
	m.lastMethod, m.lastCtx, m.lastID, m.lastParams = "EditFollowup", ctx, messageID, params
	return &discord.Message{ID: messageID}, m.err
}

func (m *mockClient) DeleteFollowup(ctx context.Context, _ *interactions.Interaction, messageID discord.Snowflake) error {
	m.lastMethod, m.lastCtx, m.lastID = "DeleteFollowup", ctx, messageID
	return m.err
}

// ── suite ─────────────────────────────────────────────────────────────────────

type clientSuite struct {
	suite.Suite
}

func TestClientSuite(t *testing.T) { suite.Run(t, new(clientSuite)) }

// hydrated returns an Interaction wired up to a fresh mockClient.
func (s *clientSuite) hydrated() (*interactions.Interaction, *mockClient) {
	m := &mockClient{}
	i := &interactions.Interaction{ID: 42}
	i.Hydrate(m, context.Background())
	return i, m
}

// ── not-hydrated: every method returns the errNotHydrated error ───────────────

func (s *clientSuite) TestNotHydrated_AllMethodsError() {
	i := &interactions.Interaction{ID: 1} // never hydrated

	_, err := i.Reply(interactions.ReplyOptions{})
	s.Error(err)
	s.Error(i.DeferReply(interactions.DeferOptions{}))
	s.Error(i.DeferUpdateMessage())
	s.Error(i.UpdateMessage(interactions.UpdateMessageOptions{}))
	s.Error(i.Modal(interactions.ModalOptions{}))
	s.Error(i.Autocomplete(interactions.AutocompleteOptions{}))
	s.Error(i.LaunchActivity())
	_, err = i.GetOriginalResponse()
	s.Error(err)
	_, err = i.EditReply(api.EditMessageParams{})
	s.Error(err)
	s.Error(i.DeleteReply())
	_, err = i.FollowUp(interactions.FollowUpOptions{})
	s.Error(err)
	_, err = i.GetFollowup(1)
	s.Error(err)
	_, err = i.EditFollowup(1, api.EditMessageParams{})
	s.Error(err)
	s.Error(i.DeleteFollowup(1))
	_, err = i.DeferAndFollowup(interactions.DeferOptions{})
	s.Error(err)
}

// ── hydrated: each method delegates to the client ─────────────────────────────

func (s *clientSuite) TestReply() {
	i, m := s.hydrated()
	mentions := &discord.AllowedMentions{}
	poll := &discord.PollRequest{}
	_, err := i.Reply(interactions.ReplyOptions{
		Content:         "hi",
		Embeds:          []discord.Embed{{}},
		Components:      []discord.AnyComponent{&discord.RawComponent{}},
		AllowedMentions: mentions,
		Poll:            poll,
		Flags:           discord.MessageFlagEphemeral,
		WithResponse:    true,
	})
	s.Require().NoError(err)
	s.Equal("Reply", m.lastMethod)
	s.Equal("hi", m.lastData.Content)
	s.NotNil(m.lastData.Embeds)
	s.Equal(mentions, m.lastData.AllowedMentions)
	s.Equal(poll, m.lastData.Poll)
}

func (s *clientSuite) TestReply_Empty() {
	i, m := s.hydrated()
	_, err := i.Reply(interactions.ReplyOptions{})
	s.Require().NoError(err)
	s.Equal("", m.lastData.Content)
	s.Nil(m.lastData.Embeds)
}

func (s *clientSuite) TestDeferReply() {
	i, m := s.hydrated()
	s.Require().NoError(i.DeferReply(interactions.DeferOptions{Ephemeral: true}))
	s.Equal("DeferReply", m.lastMethod)
	s.True(m.ephemeral)
}

func (s *clientSuite) TestDeferUpdateMessage() {
	i, m := s.hydrated()
	s.Require().NoError(i.DeferUpdateMessage())
	s.Equal("DeferUpdateMessage", m.lastMethod)
}

func (s *clientSuite) TestUpdateMessage() {
	i, m := s.hydrated()
	mentions := &discord.AllowedMentions{}
	s.Require().NoError(i.UpdateMessage(interactions.UpdateMessageOptions{
		Content:         "edit",
		Embeds:          []discord.Embed{{}},
		Components:      []discord.AnyComponent{&discord.RawComponent{}},
		AllowedMentions: mentions,
	}))
	s.Equal("UpdateMessage", m.lastMethod)
	s.Equal("edit", m.lastData.Content)
	s.NotNil(m.lastData.Components)
	s.Equal(mentions, m.lastData.AllowedMentions)
}

func (s *clientSuite) TestUpdateMessage_Empty() {
	i, m := s.hydrated()
	s.Require().NoError(i.UpdateMessage(interactions.UpdateMessageOptions{}))
	s.Equal("", m.lastData.Content)
	s.Nil(m.lastData.Embeds)
}

func (s *clientSuite) TestModal() {
	i, m := s.hydrated()
	s.Require().NoError(i.Modal(interactions.ModalOptions{Modal: &components.Modal{}}))
	s.Equal("ReplyWithModal", m.lastMethod)
}

func (s *clientSuite) TestAutocomplete() {
	i, m := s.hydrated()
	s.Require().NoError(i.Autocomplete(interactions.AutocompleteOptions{
		Choices: []responses.AutocompleteChoice{},
	}))
	s.Equal("ReplyAutocomplete", m.lastMethod)
}

func (s *clientSuite) TestLaunchActivity() {
	i, m := s.hydrated()
	s.Require().NoError(i.LaunchActivity())
	s.Equal("LaunchActivity", m.lastMethod)
}

func (s *clientSuite) TestGetOriginalResponse() {
	i, m := s.hydrated()
	msg, err := i.GetOriginalResponse()
	s.Require().NoError(err)
	s.Equal(discord.Snowflake(1), msg.ID)
	s.Equal("GetOriginalResponse", m.lastMethod)
}

func (s *clientSuite) TestEditReply() {
	i, m := s.hydrated()
	_, err := i.EditReply(api.EditMessageParams{})
	s.Require().NoError(err)
	s.Equal("EditReply", m.lastMethod)
}

func (s *clientSuite) TestDeleteReply() {
	i, m := s.hydrated()
	s.Require().NoError(i.DeleteReply())
	s.Equal("DeleteReply", m.lastMethod)
}

func (s *clientSuite) TestFollowUp() {
	i, m := s.hydrated()
	poll := &discord.PollRequest{}
	_, err := i.FollowUp(interactions.FollowUpOptions{
		Content:               "follow",
		Poll:                  poll,
		Ephemeral:             true,
		SuppressEmbeds:        true,
		SuppressNotifications: true,
	})
	s.Require().NoError(err)
	s.Equal("CreateFollowup", m.lastMethod)
	params := m.lastParams.(api.CreateMessageParams)
	s.Equal("follow", params.Content)
	s.NotZero(params.Flags & discord.MessageFlagEphemeral)
	s.NotZero(params.Flags & discord.MessageFlagSuppressEmbeds)
	s.NotZero(params.Flags & discord.MessageFlagSuppressNotification)
}

func (s *clientSuite) TestFollowUp_NoFlags() {
	i, m := s.hydrated()
	_, err := i.FollowUp(interactions.FollowUpOptions{Content: "plain"})
	s.Require().NoError(err)
	params := m.lastParams.(api.CreateMessageParams)
	s.Zero(params.Flags)
}

func (s *clientSuite) TestGetFollowup() {
	i, m := s.hydrated()
	msg, err := i.GetFollowup(99)
	s.Require().NoError(err)
	s.Equal(discord.Snowflake(99), msg.ID)
	s.Equal("GetFollowup", m.lastMethod)
}

func (s *clientSuite) TestEditFollowup() {
	i, m := s.hydrated()
	_, err := i.EditFollowup(7, api.EditMessageParams{})
	s.Require().NoError(err)
	s.Equal("EditFollowup", m.lastMethod)
	s.Equal(discord.Snowflake(7), m.lastID)
}

func (s *clientSuite) TestDeleteFollowup() {
	i, m := s.hydrated()
	s.Require().NoError(i.DeleteFollowup(8))
	s.Equal("DeleteFollowup", m.lastMethod)
	s.Equal(discord.Snowflake(8), m.lastID)
}

func (s *clientSuite) TestDeferAndFollowup() {
	i, m := s.hydrated()
	send, err := i.DeferAndFollowup(interactions.DeferOptions{Ephemeral: true})
	s.Require().NoError(err)
	s.Equal("DeferReply", m.lastMethod)

	msg, err := send(context.Background(), interactions.FollowUpOptions{Content: "done"})
	s.Require().NoError(err)
	s.Equal(discord.Snowflake(3), msg.ID)
	s.Equal("CreateFollowup", m.lastMethod)
}

func (s *clientSuite) TestDeferAndFollowup_DeferError() {
	m := &mockClient{err: errors.New("boom")}
	i := &interactions.Interaction{ID: 1}
	i.Hydrate(m, context.Background())

	_, err := i.DeferAndFollowup(interactions.DeferOptions{})
	s.Error(err)
}

// ── WithContext overrides the captured dispatch context ───────────────────────

func (s *clientSuite) TestWithContext() {
	i, m := s.hydrated()
	type ctxKey string
	ctx := context.WithValue(context.Background(), ctxKey("k"), "v")

	_, err := i.WithContext(ctx).Reply(interactions.ReplyOptions{Content: "x"})
	s.Require().NoError(err)
	s.Equal("v", m.lastCtx.Value(ctxKey("k")))
}

// ── ensureHydrated falls back to context.Background when hCtx is nil ───────────

func (s *clientSuite) TestHydrate_NilContextFallsBack() {
	m := &mockClient{}
	i := &interactions.Interaction{ID: 1}
	i.Hydrate(m, nil) // nil ctx → ensureHydrated substitutes context.Background()

	s.Require().NoError(i.DeleteReply())
	s.NotNil(m.lastCtx)
}
