package events

var EventFactories = map[EventType]func() Event{
	EventMessageCreate:     MessageCreateEvent{}.DesiredEventType,
	EventReady:             ReadyEvent{}.DesiredEventType,
	EventGuildCreate:       GuildCreateEvent{}.DesiredEventType,
	EventInteractionCreate: InteractionCreateEvent{}.DesiredEventType,
	EventGuildDelete:       GuildDeleteEvent{}.DesiredEventType,
}

type Event interface {
	Event() EventType
	DesiredEventType() Event
}

type EventType string

const (
	EventMessageCreate     EventType = "MESSAGE_CREATE"
	EventReady             EventType = "READY"
	EventGuildCreate       EventType = "GUILD_CREATE"
	EventInteractionCreate EventType = "INTERACTION_CREATE"
	EventGuildDelete       EventType = "GUILD_DELETE"
	/*TODO
	MessageDelete
	MessageUpdate

	GuildAuditLogEntryCreate

	ChannelCreate
	ChannelUpdate
	ChannelDelete
	ChannelPinsUpdate

	RoleCreate
	RoleUpdate
	RoleDelete

	WebhookUpdate

	IntegrationCreate
	IntegrationUpdate
	IntegrationDelete

	AutoModerationRuleCreate
	AutoModerationRuleUpdate
	AutoModerationRuleDelete
	AutoModerationActionExecute

	ThreadCreate
	ThreadUpdate
	ThreadDelete
	ThreadMemberUpdate
	ThreadMembersUpdate

	EntitlementCreate
	EntitlementUpdate
	EntitlementDelete

	GuildBanAdd
	GuildBanRemove
	GuildEmojisUpdate
	GuildStickersUpdate
	GuildIntegrationsUpdate
	GuildMemberAdd
	GuildMemberRemove

	ScheduledEventCreate
	ScheduledEventUpdate
	ScheduledEventDelete
	ScheduledEventUserAdd
	ScheduledEventUserRemove

	SoundboardSoundsCreate
	SoundboardSoundsUpdate
	SoundboardSoundsDelete

	InviteCreate
	InviteDelete

	MessageReactionAdd
	MessageReactionRemove
	MessageReactionRemoveAll
	MessageReactionRemoveEmoji

	PresenceUpdate

	StageInstanceUpdate
	StageInstanceCreate
	StageInstanceDelete

	SubscriptionCreate
	SubscriptionDelete
	SubscriptionUpdate

	TypingStart

	UserUpdate

	VoiceStateUpdate

	MessagePollVoteAdd
	MessagePollVoteRemove
	*/
)
