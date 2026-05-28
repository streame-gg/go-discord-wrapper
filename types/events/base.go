package events

// eventFactories maps each EventType to a factory that returns a new zero-value
// instance of the corresponding event struct. Populated via RegisterEvent.
var eventFactories = map[EventType]func() Event{}

// RegisterEvent registers an event type so the gateway knows how to unmarshal
// it. Call this from an init() function in your event file:
//
//	func init() { events.RegisterEvent(MyEvent{}) }
func RegisterEvent(e Event) {
	eventFactories[e.Event()] = e.DesiredEventType
}

// GetEventFactory returns the factory function for the given event type and
// whether it was found. Use this instead of accessing the internal map directly.
func GetEventFactory(t EventType) (func() Event, bool) {
	f, ok := eventFactories[t]
	return f, ok
}

// Event values are shared across concurrent handlers and must not be mutated.
type Event interface {
	Event() EventType
	DesiredEventType() Event
}

type EventType string

const (
	EventReady                         EventType = "READY"
	EventGuildCreate                   EventType = "GUILD_CREATE"
	EventGuildDelete                   EventType = "GUILD_DELETE"
	EventGuildUpdate                   EventType = "GUILD_UPDATE"
	EventGuildBanAdd                   EventType = "GUILD_BAN_ADD"
	EventGuildBanRemove                EventType = "GUILD_BAN_REMOVE"
	EventGuildEmojisUpdate             EventType = "GUILD_EMOJIS_UPDATE"
	EventGuildStickersUpdate           EventType = "GUILD_STICKERS_UPDATE"
	EventGuildMemberAdd                EventType = "GUILD_MEMBER_ADD"
	EventGuildMemberRemove             EventType = "GUILD_MEMBER_REMOVE"
	EventGuildMemberUpdate             EventType = "GUILD_MEMBER_UPDATE"
	EventGuildRoleCreate               EventType = "GUILD_ROLE_CREATE"
	EventGuildRoleUpdate               EventType = "GUILD_ROLE_UPDATE"
	EventGuildRoleDelete               EventType = "GUILD_ROLE_DELETE"
	EventGuildScheduledEventCreate     EventType = "GUILD_SCHEDULED_EVENT_CREATE"
	EventGuildScheduledEventUpdate     EventType = "GUILD_SCHEDULED_EVENT_UPDATE"
	EventGuildScheduledEventDelete     EventType = "GUILD_SCHEDULED_EVENT_DELETE"
	EventGuildScheduledEventUserAdd    EventType = "GUILD_SCHEDULED_EVENT_USER_ADD"
	EventGuildScheduledEventUserRemove EventType = "GUILD_SCHEDULED_EVENT_USER_REMOVE"
	EventGuildAuditLogEntryCreate      EventType = "GUILD_AUDIT_LOG_ENTRY_CREATE"
	EventChannelCreate                 EventType = "CHANNEL_CREATE"
	EventChannelUpdate                 EventType = "CHANNEL_UPDATE"
	EventChannelDelete                 EventType = "CHANNEL_DELETE"
	EventChannelPinsUpdate             EventType = "CHANNEL_PINS_UPDATE"
	EventThreadCreate                  EventType = "THREAD_CREATE"
	EventThreadUpdate                  EventType = "THREAD_UPDATE"
	EventThreadDelete                  EventType = "THREAD_DELETE"
	EventThreadListSync                EventType = "THREAD_LIST_SYNC"
	EventMessageCreate                 EventType = "MESSAGE_CREATE"
	EventMessageUpdate                 EventType = "MESSAGE_UPDATE"
	EventMessageDelete                 EventType = "MESSAGE_DELETE"
	EventMessageDeleteBulk             EventType = "MESSAGE_DELETE_BULK"
	EventMessageReactionAdd            EventType = "MESSAGE_REACTION_ADD"
	EventMessageReactionRemove         EventType = "MESSAGE_REACTION_REMOVE"
	EventMessageReactionRemoveAll      EventType = "MESSAGE_REACTION_REMOVE_ALL"
	EventMessageReactionRemoveEmoji    EventType = "MESSAGE_REACTION_REMOVE_EMOJI"
	EventMessagePollVoteAdd            EventType = "MESSAGE_POLL_VOTE_ADD"
	EventMessagePollVoteRemove         EventType = "MESSAGE_POLL_VOTE_REMOVE"
	EventInteractionCreate             EventType = "INTERACTION_CREATE"
	EventInviteCreate                  EventType = "INVITE_CREATE"
	EventInviteDelete                  EventType = "INVITE_DELETE"
	EventPresenceUpdate                EventType = "PRESENCE_UPDATE"
	EventTypingStart                   EventType = "TYPING_START"
	EventUserUpdate                    EventType = "USER_UPDATE"
	EventVoiceStateUpdate              EventType = "VOICE_STATE_UPDATE"
	EventStageInstanceCreate           EventType = "STAGE_INSTANCE_CREATE"
	EventStageInstanceUpdate           EventType = "STAGE_INSTANCE_UPDATE"
	EventStageInstanceDelete           EventType = "STAGE_INSTANCE_DELETE"
	EventIntegrationCreate             EventType = "INTEGRATION_CREATE"
	EventIntegrationUpdate             EventType = "INTEGRATION_UPDATE"
	EventIntegrationDelete             EventType = "INTEGRATION_DELETE"
	EventWebhooksUpdate                EventType = "WEBHOOKS_UPDATE"
	EventAutoModerationRuleCreate      EventType = "AUTO_MODERATION_RULE_CREATE"
	EventAutoModerationRuleUpdate      EventType = "AUTO_MODERATION_RULE_UPDATE"
	EventAutoModerationRuleDelete      EventType = "AUTO_MODERATION_RULE_DELETE"
	EventAutoModerationActionExecution EventType = "AUTO_MODERATION_ACTION_EXECUTION"
	EventEntitlementCreate             EventType = "ENTITLEMENT_CREATE"
	EventEntitlementUpdate             EventType = "ENTITLEMENT_UPDATE"
	EventEntitlementDelete             EventType = "ENTITLEMENT_DELETE"

	EventResumed EventType = "RESUMED"

	EventThreadMemberUpdate  EventType = "THREAD_MEMBER_UPDATE"
	EventThreadMembersUpdate EventType = "THREAD_MEMBERS_UPDATE"

	EventGuildMembersChunk EventType = "GUILD_MEMBERS_CHUNK"

	EventApplicationCommandPermissionsUpdate EventType = "APPLICATION_COMMAND_PERMISSIONS_UPDATE"

	EventGuildSoundboardSoundCreate  EventType = "GUILD_SOUNDBOARD_SOUND_CREATE"
	EventGuildSoundboardSoundUpdate  EventType = "GUILD_SOUNDBOARD_SOUND_UPDATE"
	EventGuildSoundboardSoundDelete  EventType = "GUILD_SOUNDBOARD_SOUND_DELETE"
	EventGuildSoundboardSoundsUpdate EventType = "GUILD_SOUNDBOARD_SOUNDS_UPDATE"

	EventVoiceServerUpdate      EventType = "VOICE_SERVER_UPDATE"
	EventVoiceChannelEffectSend EventType = "VOICE_CHANNEL_EFFECT_SEND"

	EventSubscriptionCreate EventType = "SUBSCRIPTION_CREATE"
	EventSubscriptionUpdate EventType = "SUBSCRIPTION_UPDATE"
	EventSubscriptionDelete EventType = "SUBSCRIPTION_DELETE"

	EventVoiceChannelStatusUpdate    EventType = "VOICE_CHANNEL_STATUS_UPDATE"
	EventVoiceChannelStartTimeUpdate EventType = "VOICE_CHANNEL_START_TIME_UPDATE"
	EventSoundboardSounds            EventType = "SOUNDBOARD_SOUNDS"

	// Wrapper-synthesized events — derived by go-discord-wrapper, never sent by Discord directly.
	EventWrapperVoiceMemberJoin   EventType = "WRAPPER_VOICE_MEMBER_JOIN"
	EventWrapperVoiceMemberLeave  EventType = "WRAPPER_VOICE_MEMBER_LEAVE"
	EventWrapperVoiceMemberMove   EventType = "WRAPPER_VOICE_MEMBER_MOVE"
	EventWrapperVoiceMemberUpdate EventType = "WRAPPER_VOICE_MEMBER_UPDATE"

	EventWrapperGuildMemberRoleAdd    EventType = "WRAPPER_GUILD_MEMBER_ROLE_ADD"
	EventWrapperGuildMemberRoleRemove EventType = "WRAPPER_GUILD_MEMBER_ROLE_REMOVE"
	EventWrapperGuildMemberNickChange EventType = "WRAPPER_GUILD_MEMBER_NICK_CHANGE"
	EventWrapperGuildMemberTimeout    EventType = "WRAPPER_GUILD_MEMBER_TIMEOUT"
	EventWrapperGuildMemberBoostStart EventType = "WRAPPER_GUILD_MEMBER_BOOST_START"
	EventWrapperGuildMemberBoostEnd   EventType = "WRAPPER_GUILD_MEMBER_BOOST_END"

	EventWrapperUserOnline         EventType = "WRAPPER_USER_ONLINE"
	EventWrapperUserOffline        EventType = "WRAPPER_USER_OFFLINE"
	EventWrapperUserActivityChange EventType = "WRAPPER_USER_ACTIVITY_CHANGE"
	EventWrapperUserUpdate         EventType = "WRAPPER_USER_UPDATE"

	EventWrapperGuildEmojiAdd    EventType = "WRAPPER_GUILD_EMOJI_ADD"
	EventWrapperGuildEmojiRemove EventType = "WRAPPER_GUILD_EMOJI_REMOVE"
	EventWrapperGuildEmojiUpdate EventType = "WRAPPER_GUILD_EMOJI_UPDATE"

	EventWrapperGuildStickerAdd    EventType = "WRAPPER_GUILD_STICKER_ADD"
	EventWrapperGuildStickerRemove EventType = "WRAPPER_GUILD_STICKER_REMOVE"
	EventWrapperGuildStickerUpdate EventType = "WRAPPER_GUILD_STICKER_UPDATE"

	EventWrapperGuildRolePermissionsChange EventType = "WRAPPER_GUILD_ROLE_PERMISSIONS_CHANGE"
)
