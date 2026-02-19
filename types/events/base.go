package events

var EventFactories = map[EventType]func() Event{
	EventMessageCreate:            MessageCreateEvent{}.DesiredEventType,
	EventReady:                    ReadyEvent{}.DesiredEventType,
	EventGuildCreate:              GuildCreateEvent{}.DesiredEventType,
	EventInteractionCreate:        InteractionCreateEvent{}.DesiredEventType,
	EventGuildDelete:              GuildDeleteEvent{}.DesiredEventType,
	EventInviteCreate:             InviteCreateEvent{}.DesiredEventType,
	EventInviteDelete:             InviteDeleteEvent{}.DesiredEventType,
	EventChannelCreate:            ChannelCreateEvent{}.DesiredEventType,
	EventChannelDelete:            ChannelDeleteEvent{}.DesiredEventType,
	EventMessageDelete:            MessageDeleteEvent{}.DesiredEventType,
	EventMessageDeleteBulk:        MessageDeleteBulkEvent{}.DesiredEventType,
	EventMessageUpdate:            MessageUpdateEvent{}.DesiredEventType,
	EventGuildAuditLogEntryCreate: GuildAuditLogEntryCreateEvent{}.DesiredEventType,
}

type Event interface {
	Event() EventType
	DesiredEventType() Event
}

type EventType string

/*
CHANNEL_UPDATE will not be implemented yet, use EventGuildAuditLogEntryCreate instead
*/

const (
	EventMessageCreate            EventType = "MESSAGE_CREATE"
	EventReady                    EventType = "READY"
	EventGuildCreate              EventType = "GUILD_CREATE"
	EventInteractionCreate        EventType = "INTERACTION_CREATE"
	EventGuildDelete              EventType = "GUILD_DELETE"
	EventMessageDelete            EventType = "MESSAGE_DELETE"
	EventMessageDeleteBulk        EventType = "MESSAGE_DELETE_BULK"
	EventMessageUpdate            EventType = "MESSAGE_UPDATE"
	EventGuildAuditLogEntryCreate EventType = "GUILD_AUDIT_LOG_ENTRY_CREATE"

	EventChannelCreate EventType = "CHANNEL_CREATE"
	EventChannelDelete EventType = "CHANNEL_DELETE"
	/*
		ChannelPinsUpdate EventType = "CHANNEL_PINS_UPDATE"


			EventChannelDelete
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

	*/
	EventInviteCreate EventType = "INVITE_CREATE"
	EventInviteDelete EventType = "INVITE_DELETE"

	/*

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
