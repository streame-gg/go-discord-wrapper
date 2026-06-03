package connection

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// eventRegistrationSuite exercises the thin On* handler-registration wrappers on
// Client and the read-only cacheAdapter delegators, which are otherwise only
// reached indirectly.
type eventRegistrationSuite struct {
	suite.Suite
	client *Client
}

func TestEventRegistrationSuite(t *testing.T) { suite.Run(t, new(eventRegistrationSuite)) }

func (s *eventRegistrationSuite) SetupTest() {
	c, err := NewClient("Bot fake-token", discord.IntentGuilds)
	s.Require().NoError(err)
	s.client = c
}

// TestRegisterAllHandlers calls every typed On* registration helper with a
// no-op handler, covering the full set of one-line wrappers in clientEvents.go.
func (s *eventRegistrationSuite) TestRegisterAllHandlers() {
	c := s.client
	c.OnReady(func(*Client, *events.ReadyEvent) {})
	c.OnResumed(func(*Client, *events.ResumedEvent) {})
	c.OnGuildCreate(func(*Client, *events.GuildCreateEvent) {})
	c.OnGuildUpdate(func(*Client, *events.GuildUpdateEvent) {})
	c.OnGuildDelete(func(*Client, *events.GuildDeleteEvent) {})
	c.OnGuildBanAdd(func(*Client, *events.GuildBanAddEvent) {})
	c.OnGuildBanRemove(func(*Client, *events.GuildBanRemoveEvent) {})
	c.OnGuildEmojisUpdate(func(*Client, *events.GuildEmojisUpdateEvent) {})
	c.OnGuildStickersUpdate(func(*Client, *events.GuildStickersUpdateEvent) {})
	c.OnGuildMemberAdd(func(*Client, *events.GuildMemberAddEvent) {})
	c.OnGuildMemberRemove(func(*Client, *events.GuildMemberRemoveEvent) {})
	c.OnGuildMemberUpdate(func(*Client, *events.GuildMemberUpdateEvent) {})
	c.OnGuildMembersChunk(func(*Client, *events.GuildMembersChunkEvent) {})
	c.OnGuildRoleCreate(func(*Client, *events.GuildRoleCreateEvent) {})
	c.OnGuildRoleUpdate(func(*Client, *events.GuildRoleUpdateEvent) {})
	c.OnGuildRoleDelete(func(*Client, *events.GuildRoleDeleteEvent) {})
	c.OnGuildAuditLogEntryCreate(func(*Client, *events.GuildAuditLogEntryCreateEvent) {})
	c.OnGuildScheduledEventCreate(func(*Client, *events.GuildScheduledEventCreateEvent) {})
	c.OnGuildScheduledEventUpdate(func(*Client, *events.GuildScheduledEventUpdateEvent) {})
	c.OnGuildScheduledEventDelete(func(*Client, *events.GuildScheduledEventDeleteEvent) {})
	c.OnGuildScheduledEventUserAdd(func(*Client, *events.GuildScheduledEventUserAddEvent) {})
	c.OnGuildScheduledEventUserRemove(func(*Client, *events.GuildScheduledEventUserRemoveEvent) {})
	c.OnChannelCreate(func(*Client, *events.ChannelCreateEvent) {})
	c.OnChannelUpdate(func(*Client, *events.ChannelUpdateEvent) {})
	c.OnChannelDelete(func(*Client, *events.ChannelDeleteEvent) {})
	c.OnChannelPinsUpdate(func(*Client, *events.ChannelPinsUpdateEvent) {})
	c.OnThreadCreate(func(*Client, *events.ThreadCreateEvent) {})
	c.OnThreadUpdate(func(*Client, *events.ThreadUpdateEvent) {})
	c.OnThreadDelete(func(*Client, *events.ThreadDeleteEvent) {})
	c.OnThreadListSync(func(*Client, *events.ThreadListSyncEvent) {})
	c.OnThreadMemberUpdate(func(*Client, *events.ThreadMemberUpdateEvent) {})
	c.OnThreadMembersUpdate(func(*Client, *events.ThreadMembersUpdateEvent) {})
	c.OnMessageCreate(func(*Client, *events.MessageCreateEvent) {})
	c.OnMessageUpdate(func(*Client, *events.MessageUpdateEvent) {})
	c.OnMessageDelete(func(*Client, *events.MessageDeleteEvent) {})
	c.OnMessageDeleteBulk(func(*Client, *events.MessageDeleteBulkEvent) {})
	c.OnMessageReactionAdd(func(*Client, *events.MessageReactionAddEvent) {})
	c.OnMessageReactionRemove(func(*Client, *events.MessageReactionRemoveEvent) {})
	c.OnMessageReactionRemoveAll(func(*Client, *events.MessageReactionRemoveAllEvent) {})
	c.OnMessageReactionRemoveEmoji(func(*Client, *events.MessageReactionRemoveEmojiEvent) {})
	c.OnMessagePollVoteAdd(func(*Client, *events.MessagePollVoteAddEvent) {})
	c.OnMessagePollVoteRemove(func(*Client, *events.MessagePollVoteRemoveEvent) {})
	c.OnInteractionCreate(func(*Client, *events.InteractionCreateEvent) {})
	c.OnPresenceUpdate(func(*Client, *events.PresenceUpdateEvent) {})
	c.OnTypingStart(func(*Client, *events.TypingStartEvent) {})
	c.OnUserUpdate(func(*Client, *events.UserUpdateEvent) {})
	c.OnVoiceStateUpdate(func(*Client, *events.VoiceStateUpdateEvent) {})
	c.OnVoiceServerUpdate(func(*Client, *events.VoiceServerUpdateEvent) {})
	c.OnVoiceChannelEffectSend(func(*Client, *events.VoiceChannelEffectSendEvent) {})
	c.OnGuildEmojiAdd(func(*Client, *events.GuildEmojiAddEvent) {})
	c.OnGuildEmojiRemove(func(*Client, *events.GuildEmojiRemoveEvent) {})
	c.OnGuildEmojiUpdate(func(*Client, *events.GuildEmojiUpdateEvent) {})
	c.OnGuildStickerAdd(func(*Client, *events.GuildStickerAddEvent) {})
	c.OnGuildStickerRemove(func(*Client, *events.GuildStickerRemoveEvent) {})
	c.OnGuildStickerUpdate(func(*Client, *events.GuildStickerUpdateEvent) {})
	c.OnStageInstanceCreate(func(*Client, *events.StageInstanceCreateEvent) {})
	c.OnStageInstanceUpdate(func(*Client, *events.StageInstanceUpdateEvent) {})
	c.OnStageInstanceDelete(func(*Client, *events.StageInstanceDeleteEvent) {})
	c.OnInviteCreate(func(*Client, *events.InviteCreateEvent) {})
	c.OnInviteDelete(func(*Client, *events.InviteDeleteEvent) {})
	c.OnIntegrationCreate(func(*Client, *events.IntegrationCreateEvent) {})
	c.OnIntegrationUpdate(func(*Client, *events.IntegrationUpdateEvent) {})
	c.OnIntegrationDelete(func(*Client, *events.IntegrationDeleteEvent) {})
	c.OnWebhooksUpdate(func(*Client, *events.WebhooksUpdateEvent) {})
	c.OnAutoModerationRuleCreate(func(*Client, *events.AutoModerationRuleCreateEvent) {})
	c.OnAutoModerationRuleUpdate(func(*Client, *events.AutoModerationRuleUpdateEvent) {})
	c.OnAutoModerationRuleDelete(func(*Client, *events.AutoModerationRuleDeleteEvent) {})
	c.OnAutoModerationActionExecution(func(*Client, *events.AutoModerationActionExecutionEvent) {})
	c.OnEntitlementCreate(func(*Client, *events.EntitlementCreateEvent) {})
	c.OnEntitlementUpdate(func(*Client, *events.EntitlementUpdateEvent) {})
	c.OnEntitlementDelete(func(*Client, *events.EntitlementDeleteEvent) {})
	c.OnApplicationCommandPermissionsUpdate(func(*Client, *events.ApplicationCommandPermissionsUpdateEvent) {})
	c.OnGuildSoundboardSoundCreate(func(*Client, *events.GuildSoundboardSoundCreateEvent) {})
	c.OnGuildSoundboardSoundUpdate(func(*Client, *events.GuildSoundboardSoundUpdateEvent) {})
	c.OnGuildSoundboardSoundDelete(func(*Client, *events.GuildSoundboardSoundDeleteEvent) {})
	c.OnGuildSoundboardSoundsUpdate(func(*Client, *events.GuildSoundboardSoundsUpdateEvent) {})
	c.OnSubscriptionCreate(func(*Client, *events.SubscriptionCreateEvent) {})
	c.OnSubscriptionUpdate(func(*Client, *events.SubscriptionUpdateEvent) {})
	c.OnSubscriptionDelete(func(*Client, *events.SubscriptionDeleteEvent) {})
	c.OnVoiceChannelStatusUpdate(func(*Client, *events.VoiceChannelStatusUpdateEvent) {})
	c.OnVoiceChannelStartTimeUpdate(func(*Client, *events.VoiceChannelStartTimeUpdateEvent) {})
	c.OnSoundboardSounds(func(*Client, *events.SoundboardSoundsEvent) {})
}

// TestCacheAdapterDelegators verifies every read-only cacheAdapter method
// returns the corresponding store from the wrapped cache.
func (s *eventRegistrationSuite) TestCacheAdapterDelegators() {
	a := cacheAdapter{inner: cache.NewMemoryCache(cache.Options{})}
	s.NotNil(a.Guilds())
	s.NotNil(a.Channels())
	s.NotNil(a.Users())
	s.NotNil(a.Members())
	s.NotNil(a.Messages())
	s.NotNil(a.Roles())
	s.NotNil(a.VoiceStates())
	s.NotNil(a.Soundboard())
	s.NotNil(a.ScheduledEvents())
	s.NotNil(a.StageInstances())
	s.NotNil(a.Emojis())
	s.NotNil(a.Stickers())
	s.NotNil(a.Bans())
	s.NotNil(a.AutoModRules())
	s.NotNil(a.Invites())
	var _ discord.Cache = a
}
