package managers_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/discord/managers"
)

// restDelegationSuite drives every manager's full method surface against a
// no-op stubEntityClient (cache == nil). This pins the REST-delegating methods
// (Fetch/FetchAll/Create/Remove) and the nil-cache fallbacks of the cache
// accessors (Cache/Get/Size) and the Snowflake branch of Resolve.
type restDelegationSuite struct {
	suite.Suite
	ctx    context.Context
	client *stubEntityClient
}

func TestRestDelegationSuite(t *testing.T) { suite.Run(t, new(restDelegationSuite)) }

func (s *restDelegationSuite) SetupTest() {
	s.ctx = context.Background()
	s.client = &stubEntityClient{} // cache nil → exercises the no-cache branches
}

const rdID = discord.Snowflake(2)

func (s *restDelegationSuite) TestAutoModRuleManager() {
	m := managers.NewAutoModRuleManager(gID, s.client)
	_, _ = m.Fetch(s.ctx, rdID)
	_, _ = m.FetchAll(s.ctx)
	_, _ = m.Create(s.ctx, discord.RuleCreateOptions{})
}

func (s *restDelegationSuite) TestBanManager() {
	m := managers.NewBanManager(gID, s.client)
	_, _ = m.Fetch(s.ctx, rdID)
	_, _ = m.FetchAll(s.ctx, discord.FetchBansOptions{})
	s.NoError(m.Create(s.ctx, rdID, discord.BanOptions{}))
	s.NoError(m.Remove(s.ctx, rdID, nil))
}

func (s *restDelegationSuite) TestClientGuildManager() {
	m := managers.NewClientGuildManager(s.client)
	s.Zero(m.Cache().Len())
	_, ok := m.Get(gID)
	s.False(ok)
	_, _ = m.Fetch(s.ctx, gID)
	_, _ = m.Resolve(gID)
	s.Zero(m.Size())
}

func (s *restDelegationSuite) TestClientChannelManager() {
	m := managers.NewClientChannelManager(s.client)
	s.Zero(m.Cache().Len())
	_, ok := m.Get(gID)
	s.False(ok)
	_, _ = m.Fetch(s.ctx, gID)
	_, _ = m.Resolve(gID)
	s.Zero(m.Size())
}

func (s *restDelegationSuite) TestClientUserManager() {
	m := managers.NewClientUserManager(s.client)
	s.Zero(m.Cache().Len())
	_, ok := m.Get(gID)
	s.False(ok)
	_, _ = m.Fetch(s.ctx, gID)
	_, _ = m.Resolve(gID)
	s.Zero(m.Size())
}

func (s *restDelegationSuite) TestEmojiManager() {
	m := managers.NewEmojiManager(gID, s.client)
	s.Zero(m.Cache().Len())
	_, _ = m.Get(rdID)
	_, _ = m.Fetch(s.ctx, rdID)
	_, _ = m.FetchAll(s.ctx)
	_, _ = m.Create(s.ctx, discord.EmojiCreateOptions{})
	_, _ = m.Resolve(rdID)
	s.Zero(m.Size())
}

func (s *restDelegationSuite) TestGuildChannelManager() {
	m := managers.NewGuildChannelManager(gID, s.client)
	s.Zero(m.Cache().Len())
	_, _ = m.Get(rdID)
	_, _ = m.Fetch(s.ctx, rdID)
	_, _ = m.FetchAll(s.ctx)
	_, _ = m.Create(s.ctx, discord.ChannelCreateOptions{})
	_, _ = m.Resolve(rdID)
	s.Zero(m.Size())
}

func (s *restDelegationSuite) TestGuildInviteManager() {
	m := managers.NewGuildInviteManager(gID, s.client)
	_, _ = m.FetchAll(s.ctx)
}

func (s *restDelegationSuite) TestGuildWebhookManager() {
	m := managers.NewGuildWebhookManager(gID, s.client)
	_, _ = m.FetchAll(s.ctx)
}

func (s *restDelegationSuite) TestIntegrationManager() {
	m := managers.NewIntegrationManager(gID, s.client)
	_, _ = m.FetchAll(s.ctx)
	s.NoError(m.Remove(s.ctx, rdID, nil))
}

func (s *restDelegationSuite) TestMemberManager() {
	m := managers.NewMemberManager(gID, s.client)
	s.Zero(m.Cache().Len())
	_, _ = m.Get(rdID)
	_, _ = m.Fetch(s.ctx, rdID)
	_, _ = m.FetchAll(s.ctx, discord.FetchMembersOptions{})
	_, _ = m.Resolve(rdID)
	s.Zero(m.Size())
}

func (s *restDelegationSuite) TestMessageManager() {
	m := managers.NewMessageManager(gID, s.client)
	s.Zero(m.Cache().Len())
	_, _ = m.Get(rdID)
	_, _ = m.Fetch(s.ctx, rdID)
	_, _ = m.FetchAll(s.ctx, discord.FetchMessagesOptions{})
	_, _ = m.Create(s.ctx, discord.MessageCreateOptions{})
	_, _ = m.Resolve(rdID)
	s.Zero(m.Size())
}

func (s *restDelegationSuite) TestRoleManager() {
	m := managers.NewRoleManager(gID, s.client)
	s.Zero(m.Cache().Len())
	_, _ = m.Get(rdID)
	_, _ = m.Fetch(s.ctx, rdID)
	_, _ = m.FetchAll(s.ctx)
	_, _ = m.Create(s.ctx, discord.RoleCreateOptions{})
	_, _ = m.Resolve(rdID)
	s.Zero(m.Size())
}

func (s *restDelegationSuite) TestScheduledEventManager() {
	m := managers.NewScheduledEventManager(gID, s.client)
	s.Zero(m.Cache().Len())
	_, _ = m.Get(rdID)
	_, _ = m.Fetch(s.ctx, rdID)
	_, _ = m.FetchAll(s.ctx)
	_, _ = m.Create(s.ctx, discord.ScheduledEventCreateOptions{})
	_, _ = m.Resolve(rdID)
	s.Zero(m.Size())
}

func (s *restDelegationSuite) TestSoundboardManager() {
	m := managers.NewSoundboardManager(gID, s.client)
	s.Zero(m.Cache().Len())
	_, _ = m.Get(rdID)
	_, _ = m.Fetch(s.ctx, rdID)
	_, _ = m.FetchAll(s.ctx)
	_, _ = m.Resolve(rdID)
	s.Zero(m.Size())
}

func (s *restDelegationSuite) TestStageInstanceManager() {
	m := managers.NewStageInstanceManager(gID, s.client)
	s.Zero(m.Cache().Len())
	_, _ = m.Get(rdID)
	_, _ = m.Fetch(s.ctx, rdID)
	_, _ = m.Create(s.ctx, discord.StageCreateOptions{})
	_, _ = m.Resolve(rdID)
	s.Zero(m.Size())
}

func (s *restDelegationSuite) TestStickerManager() {
	m := managers.NewStickerManager(gID, s.client)
	s.Zero(m.Cache().Len())
	_, _ = m.Get(rdID)
	_, _ = m.Fetch(s.ctx, rdID)
	_, _ = m.FetchAll(s.ctx)
	_, _ = m.Create(s.ctx, discord.StickerCreateOptions{})
	_, _ = m.Resolve(rdID)
	s.Zero(m.Size())
}

func (s *restDelegationSuite) TestThreadManager() {
	m := managers.NewThreadManager(gID, s.client)
	s.Zero(m.Cache().Len())
	_, _ = m.Get(rdID)
	_, _ = m.Fetch(s.ctx, rdID)
	_, _ = m.Resolve(rdID)
	s.Zero(m.Size())
}

func (s *restDelegationSuite) TestVoiceStateManager() {
	m := managers.NewVoiceStateManager(gID, s.client)
	s.Zero(m.Cache().Len())
	_, _ = m.Get(rdID)
	_, _ = m.Resolve(rdID)
	s.Zero(m.Size())
}
