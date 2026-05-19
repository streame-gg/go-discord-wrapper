package connection

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// Middleware wraps an EventHandler to add cross-cutting behaviour such as
// logging, permission checks, or metrics collection. The next handler must be
// called for the chain to continue; not calling it short-circuits processing.
//
//	client.Use(func(next connection.EventHandler) connection.EventHandler {
//	    return func(c *connection.Client, ev events.Event) {
//	        log.Println("before", ev.Event())
//	        next(c, ev)
//	        log.Println("after", ev.Event())
//	    }
//	})
type Middleware func(next EventHandler) EventHandler

// Use registers one or more middleware that wrap every event handler registered
// after this call. Middleware are applied in registration order — the first
// registered middleware is the outermost wrapper.
func (d *Client) Use(mw ...Middleware) {
	d.middlewareMu.Lock()
	d.middleware = append(d.middleware, mw...)
	d.middlewareMu.Unlock()
}

// applyMiddleware wraps handler with the current middleware chain (snapshot).
func (d *Client) applyMiddleware(handler EventHandler) EventHandler {
	d.middlewareMu.RLock()
	chain := make([]Middleware, len(d.middleware))
	copy(chain, d.middleware)
	d.middlewareMu.RUnlock()

	for i := len(chain) - 1; i >= 0; i-- {
		handler = chain[i](handler)
	}
	return handler
}

// ── Client lifecycle events ───────────────────────────────────────────────────

// OnConnect registers a handler fired once after the initial gateway connection
// is established (i.e. after Login() returns successfully).
func (d *Client) OnConnect(handler func(*Client)) {
	d.clientEventsMu.Lock()
	d.onConnect = append(d.onConnect, handler)
	d.clientEventsMu.Unlock()
}

// OnDisconnect registers a handler fired whenever the gateway connection drops.
// The error is the underlying cause; it may be nil on a clean close.
func (d *Client) OnDisconnect(handler func(*Client, error)) {
	d.clientEventsMu.Lock()
	d.onDisconnect = append(d.onDisconnect, handler)
	d.clientEventsMu.Unlock()
}

// OnReconnect registers a handler fired after each successful gateway reconnect.
func (d *Client) OnReconnect(handler func(*Client)) {
	d.clientEventsMu.Lock()
	d.onReconnect = append(d.onReconnect, handler)
	d.clientEventsMu.Unlock()
}

// OnPacketError registers a handler fired when an incoming gateway packet
// cannot be unmarshalled. The error describes the parsing failure.
func (d *Client) OnPacketError(handler func(*Client, error)) {
	d.clientEventsMu.Lock()
	d.onPacketError = append(d.onPacketError, handler)
	d.clientEventsMu.Unlock()
}

// ── Typed Discord event helpers ───────────────────────────────────────────────

func (d *Client) OnReady(h func(*Client, *events.ReadyEvent)) {
	_ = d.OnEvent(events.EventReady, h)
}
func (d *Client) OnResumed(h func(*Client, *events.ResumedEvent)) {
	_ = d.OnEvent(events.EventResumed, h)
}
func (d *Client) OnGuildCreate(h func(*Client, *events.GuildCreateEvent)) {
	_ = d.OnEvent(events.EventGuildCreate, h)
}
func (d *Client) OnGuildUpdate(h func(*Client, *events.GuildUpdateEvent)) {
	_ = d.OnEvent(events.EventGuildUpdate, h)
}
func (d *Client) OnGuildDelete(h func(*Client, *events.GuildDeleteEvent)) {
	_ = d.OnEvent(events.EventGuildDelete, h)
}
func (d *Client) OnGuildBanAdd(h func(*Client, *events.GuildBanAddEvent)) {
	_ = d.OnEvent(events.EventGuildBanAdd, h)
}
func (d *Client) OnGuildBanRemove(h func(*Client, *events.GuildBanRemoveEvent)) {
	_ = d.OnEvent(events.EventGuildBanRemove, h)
}
func (d *Client) OnGuildEmojisUpdate(h func(*Client, *events.GuildEmojisUpdateEvent)) {
	_ = d.OnEvent(events.EventGuildEmojisUpdate, h)
}
func (d *Client) OnGuildStickersUpdate(h func(*Client, *events.GuildStickersUpdateEvent)) {
	_ = d.OnEvent(events.EventGuildStickersUpdate, h)
}
func (d *Client) OnGuildMemberAdd(h func(*Client, *events.GuildMemberAddEvent)) {
	_ = d.OnEvent(events.EventGuildMemberAdd, h)
}
func (d *Client) OnGuildMemberRemove(h func(*Client, *events.GuildMemberRemoveEvent)) {
	_ = d.OnEvent(events.EventGuildMemberRemove, h)
}
func (d *Client) OnGuildMemberUpdate(h func(*Client, *events.GuildMemberUpdateEvent)) {
	_ = d.OnEvent(events.EventGuildMemberUpdate, h)
}
func (d *Client) OnGuildMembersChunk(h func(*Client, *events.GuildMembersChunkEvent)) {
	_ = d.OnEvent(events.EventGuildMembersChunk, h)
}
func (d *Client) OnGuildRoleCreate(h func(*Client, *events.GuildRoleCreateEvent)) {
	_ = d.OnEvent(events.EventGuildRoleCreate, h)
}
func (d *Client) OnGuildRoleUpdate(h func(*Client, *events.GuildRoleUpdateEvent)) {
	_ = d.OnEvent(events.EventGuildRoleUpdate, h)
}
func (d *Client) OnGuildRoleDelete(h func(*Client, *events.GuildRoleDeleteEvent)) {
	_ = d.OnEvent(events.EventGuildRoleDelete, h)
}
func (d *Client) OnGuildAuditLogEntryCreate(h func(*Client, *events.GuildAuditLogEntryCreateEvent)) {
	_ = d.OnEvent(events.EventGuildAuditLogEntryCreate, h)
}
func (d *Client) OnGuildScheduledEventCreate(h func(*Client, *events.GuildScheduledEventCreateEvent)) {
	_ = d.OnEvent(events.EventGuildScheduledEventCreate, h)
}
func (d *Client) OnGuildScheduledEventUpdate(h func(*Client, *events.GuildScheduledEventUpdateEvent)) {
	_ = d.OnEvent(events.EventGuildScheduledEventUpdate, h)
}
func (d *Client) OnGuildScheduledEventDelete(h func(*Client, *events.GuildScheduledEventDeleteEvent)) {
	_ = d.OnEvent(events.EventGuildScheduledEventDelete, h)
}
func (d *Client) OnGuildScheduledEventUserAdd(h func(*Client, *events.GuildScheduledEventUserAddEvent)) {
	_ = d.OnEvent(events.EventGuildScheduledEventUserAdd, h)
}
func (d *Client) OnGuildScheduledEventUserRemove(h func(*Client, *events.GuildScheduledEventUserRemoveEvent)) {
	_ = d.OnEvent(events.EventGuildScheduledEventUserRemove, h)
}
func (d *Client) OnChannelCreate(h func(*Client, *events.ChannelCreateEvent)) {
	_ = d.OnEvent(events.EventChannelCreate, h)
}
func (d *Client) OnChannelUpdate(h func(*Client, *events.ChannelUpdateEvent)) {
	_ = d.OnEvent(events.EventChannelUpdate, h)
}
func (d *Client) OnChannelDelete(h func(*Client, *events.ChannelDeleteEvent)) {
	_ = d.OnEvent(events.EventChannelDelete, h)
}
func (d *Client) OnChannelPinsUpdate(h func(*Client, *events.ChannelPinsUpdateEvent)) {
	_ = d.OnEvent(events.EventChannelPinsUpdate, h)
}
func (d *Client) OnThreadCreate(h func(*Client, *events.ThreadCreateEvent)) {
	_ = d.OnEvent(events.EventThreadCreate, h)
}
func (d *Client) OnThreadUpdate(h func(*Client, *events.ThreadUpdateEvent)) {
	_ = d.OnEvent(events.EventThreadUpdate, h)
}
func (d *Client) OnThreadDelete(h func(*Client, *events.ThreadDeleteEvent)) {
	_ = d.OnEvent(events.EventThreadDelete, h)
}
func (d *Client) OnThreadListSync(h func(*Client, *events.ThreadListSyncEvent)) {
	_ = d.OnEvent(events.EventThreadListSync, h)
}
func (d *Client) OnThreadMemberUpdate(h func(*Client, *events.ThreadMemberUpdateEvent)) {
	_ = d.OnEvent(events.EventThreadMemberUpdate, h)
}
func (d *Client) OnThreadMembersUpdate(h func(*Client, *events.ThreadMembersUpdateEvent)) {
	_ = d.OnEvent(events.EventThreadMembersUpdate, h)
}
func (d *Client) OnMessageCreate(h func(*Client, *events.MessageCreateEvent)) {
	_ = d.OnEvent(events.EventMessageCreate, h)
}
func (d *Client) OnMessageUpdate(h func(*Client, *events.MessageUpdateEvent)) {
	_ = d.OnEvent(events.EventMessageUpdate, h)
}
func (d *Client) OnMessageDelete(h func(*Client, *events.MessageDeleteEvent)) {
	_ = d.OnEvent(events.EventMessageDelete, h)
}
func (d *Client) OnMessageDeleteBulk(h func(*Client, *events.MessageDeleteBulkEvent)) {
	_ = d.OnEvent(events.EventMessageDeleteBulk, h)
}
func (d *Client) OnMessageReactionAdd(h func(*Client, *events.MessageReactionAddEvent)) {
	_ = d.OnEvent(events.EventMessageReactionAdd, h)
}
func (d *Client) OnMessageReactionRemove(h func(*Client, *events.MessageReactionRemoveEvent)) {
	_ = d.OnEvent(events.EventMessageReactionRemove, h)
}
func (d *Client) OnMessageReactionRemoveAll(h func(*Client, *events.MessageReactionRemoveAllEvent)) {
	_ = d.OnEvent(events.EventMessageReactionRemoveAll, h)
}
func (d *Client) OnMessageReactionRemoveEmoji(h func(*Client, *events.MessageReactionRemoveEmojiEvent)) {
	_ = d.OnEvent(events.EventMessageReactionRemoveEmoji, h)
}
func (d *Client) OnMessagePollVoteAdd(h func(*Client, *events.MessagePollVoteAddEvent)) {
	_ = d.OnEvent(events.EventMessagePollVoteAdd, h)
}
func (d *Client) OnMessagePollVoteRemove(h func(*Client, *events.MessagePollVoteRemoveEvent)) {
	_ = d.OnEvent(events.EventMessagePollVoteRemove, h)
}
func (d *Client) OnInteractionCreate(h func(*Client, *events.InteractionCreateEvent)) {
	_ = d.OnEvent(events.EventInteractionCreate, h)
}
func (d *Client) OnPresenceUpdate(h func(*Client, *events.PresenceUpdateEvent)) {
	_ = d.OnEvent(events.EventPresenceUpdate, h)
}
func (d *Client) OnTypingStart(h func(*Client, *events.TypingStartEvent)) {
	_ = d.OnEvent(events.EventTypingStart, h)
}
func (d *Client) OnUserUpdate(h func(*Client, *events.UserUpdateEvent)) {
	_ = d.OnEvent(events.EventUserUpdate, h)
}
func (d *Client) OnVoiceStateUpdate(h func(*Client, *events.VoiceStateUpdateEvent)) {
	_ = d.OnEvent(events.EventVoiceStateUpdate, h)
}
func (d *Client) OnVoiceServerUpdate(h func(*Client, *events.VoiceServerUpdateEvent)) {
	_ = d.OnEvent(events.EventVoiceServerUpdate, h)
}
func (d *Client) OnVoiceChannelEffectSend(h func(*Client, *events.VoiceChannelEffectSendEvent)) {
	_ = d.OnEvent(events.EventVoiceChannelEffectSend, h)
}
func (d *Client) OnUserOnline(h func(*Client, *events.UserOnlineEvent)) {
	_ = d.OnEvent(events.EventWrapperUserOnline, h)
}
func (d *Client) OnUserOffline(h func(*Client, *events.UserOfflineEvent)) {
	_ = d.OnEvent(events.EventWrapperUserOffline, h)
}
func (d *Client) OnUserActivityChange(h func(*Client, *events.UserActivityChangeEvent)) {
	_ = d.OnEvent(events.EventWrapperUserActivityChange, h)
}
func (d *Client) OnUserProfileUpdate(h func(*Client, *events.UserProfileUpdateEvent)) {
	_ = d.OnEvent(events.EventWrapperUserUpdate, h)
}
func (d *Client) OnGuildMemberRoleAdd(h func(*Client, *events.GuildMemberRoleAddEvent)) {
	_ = d.OnEvent(events.EventWrapperGuildMemberRoleAdd, h)
}
func (d *Client) OnGuildMemberRoleRemove(h func(*Client, *events.GuildMemberRoleRemoveEvent)) {
	_ = d.OnEvent(events.EventWrapperGuildMemberRoleRemove, h)
}
func (d *Client) OnGuildMemberNickChange(h func(*Client, *events.GuildMemberNickChangeEvent)) {
	_ = d.OnEvent(events.EventWrapperGuildMemberNickChange, h)
}
func (d *Client) OnGuildMemberTimeout(h func(*Client, *events.GuildMemberTimeoutEvent)) {
	_ = d.OnEvent(events.EventWrapperGuildMemberTimeout, h)
}
func (d *Client) OnGuildMemberBoostStart(h func(*Client, *events.GuildMemberBoostStartEvent)) {
	_ = d.OnEvent(events.EventWrapperGuildMemberBoostStart, h)
}
func (d *Client) OnGuildMemberBoostEnd(h func(*Client, *events.GuildMemberBoostEndEvent)) {
	_ = d.OnEvent(events.EventWrapperGuildMemberBoostEnd, h)
}
func (d *Client) OnVoiceMemberJoin(h func(*Client, *events.VoiceMemberJoinEvent)) {
	_ = d.OnEvent(events.EventWrapperVoiceMemberJoin, h)
}
func (d *Client) OnVoiceMemberLeave(h func(*Client, *events.VoiceMemberLeaveEvent)) {
	_ = d.OnEvent(events.EventWrapperVoiceMemberLeave, h)
}
func (d *Client) OnVoiceMemberMove(h func(*Client, *events.VoiceMemberMoveEvent)) {
	_ = d.OnEvent(events.EventWrapperVoiceMemberMove, h)
}
func (d *Client) OnVoiceMemberUpdate(h func(*Client, *events.VoiceMemberUpdateEvent)) {
	_ = d.OnEvent(events.EventWrapperVoiceMemberUpdate, h)
}
func (d *Client) OnStageInstanceCreate(h func(*Client, *events.StageInstanceCreateEvent)) {
	_ = d.OnEvent(events.EventStageInstanceCreate, h)
}
func (d *Client) OnStageInstanceUpdate(h func(*Client, *events.StageInstanceUpdateEvent)) {
	_ = d.OnEvent(events.EventStageInstanceUpdate, h)
}
func (d *Client) OnStageInstanceDelete(h func(*Client, *events.StageInstanceDeleteEvent)) {
	_ = d.OnEvent(events.EventStageInstanceDelete, h)
}
func (d *Client) OnInviteCreate(h func(*Client, *events.InviteCreateEvent)) {
	_ = d.OnEvent(events.EventInviteCreate, h)
}
func (d *Client) OnInviteDelete(h func(*Client, *events.InviteDeleteEvent)) {
	_ = d.OnEvent(events.EventInviteDelete, h)
}
func (d *Client) OnIntegrationCreate(h func(*Client, *events.IntegrationCreateEvent)) {
	_ = d.OnEvent(events.EventIntegrationCreate, h)
}
func (d *Client) OnIntegrationUpdate(h func(*Client, *events.IntegrationUpdateEvent)) {
	_ = d.OnEvent(events.EventIntegrationUpdate, h)
}
func (d *Client) OnIntegrationDelete(h func(*Client, *events.IntegrationDeleteEvent)) {
	_ = d.OnEvent(events.EventIntegrationDelete, h)
}
func (d *Client) OnWebhooksUpdate(h func(*Client, *events.WebhooksUpdateEvent)) {
	_ = d.OnEvent(events.EventWebhooksUpdate, h)
}
func (d *Client) OnAutoModerationRuleCreate(h func(*Client, *events.AutoModerationRuleCreateEvent)) {
	_ = d.OnEvent(events.EventAutoModerationRuleCreate, h)
}
func (d *Client) OnAutoModerationRuleUpdate(h func(*Client, *events.AutoModerationRuleUpdateEvent)) {
	_ = d.OnEvent(events.EventAutoModerationRuleUpdate, h)
}
func (d *Client) OnAutoModerationRuleDelete(h func(*Client, *events.AutoModerationRuleDeleteEvent)) {
	_ = d.OnEvent(events.EventAutoModerationRuleDelete, h)
}
func (d *Client) OnAutoModerationActionExecution(h func(*Client, *events.AutoModerationActionExecutionEvent)) {
	_ = d.OnEvent(events.EventAutoModerationActionExecution, h)
}
func (d *Client) OnEntitlementCreate(h func(*Client, *events.EntitlementCreateEvent)) {
	_ = d.OnEvent(events.EventEntitlementCreate, h)
}
func (d *Client) OnEntitlementUpdate(h func(*Client, *events.EntitlementUpdateEvent)) {
	_ = d.OnEvent(events.EventEntitlementUpdate, h)
}
func (d *Client) OnEntitlementDelete(h func(*Client, *events.EntitlementDeleteEvent)) {
	_ = d.OnEvent(events.EventEntitlementDelete, h)
}
func (d *Client) OnApplicationCommandPermissionsUpdate(h func(*Client, *events.ApplicationCommandPermissionsUpdateEvent)) {
	_ = d.OnEvent(events.EventApplicationCommandPermissionsUpdate, h)
}
func (d *Client) OnGuildSoundboardSoundCreate(h func(*Client, *events.GuildSoundboardSoundCreateEvent)) {
	_ = d.OnEvent(events.EventGuildSoundboardSoundCreate, h)
}
func (d *Client) OnGuildSoundboardSoundUpdate(h func(*Client, *events.GuildSoundboardSoundUpdateEvent)) {
	_ = d.OnEvent(events.EventGuildSoundboardSoundUpdate, h)
}
func (d *Client) OnGuildSoundboardSoundDelete(h func(*Client, *events.GuildSoundboardSoundDeleteEvent)) {
	_ = d.OnEvent(events.EventGuildSoundboardSoundDelete, h)
}
func (d *Client) OnGuildSoundboardSoundsUpdate(h func(*Client, *events.GuildSoundboardSoundsUpdateEvent)) {
	_ = d.OnEvent(events.EventGuildSoundboardSoundsUpdate, h)
}
func (d *Client) OnSubscriptionCreate(h func(*Client, *events.SubscriptionCreateEvent)) {
	_ = d.OnEvent(events.EventSubscriptionCreate, h)
}
func (d *Client) OnSubscriptionUpdate(h func(*Client, *events.SubscriptionUpdateEvent)) {
	_ = d.OnEvent(events.EventSubscriptionUpdate, h)
}
func (d *Client) OnSubscriptionDelete(h func(*Client, *events.SubscriptionDeleteEvent)) {
	_ = d.OnEvent(events.EventSubscriptionDelete, h)
}
func (d *Client) OnVoiceChannelStatusUpdate(h func(*Client, *events.VoiceChannelStatusUpdateEvent)) {
	_ = d.OnEvent(events.EventVoiceChannelStatusUpdate, h)
}
func (d *Client) OnVoiceChannelStartTimeUpdate(h func(*Client, *events.VoiceChannelStartTimeUpdateEvent)) {
	_ = d.OnEvent(events.EventVoiceChannelStartTimeUpdate, h)
}
func (d *Client) OnSoundboardSounds(h func(*Client, *events.SoundboardSoundsEvent)) {
	_ = d.OnEvent(events.EventSoundboardSounds, h)
}

func (d *Client) onRawEvent(
	eventName events.EventType,
	handler EventHandler,
) {
	wrapped := d.applyMiddleware(handler)
	d.discordEventEmitter.On(eventName, wrapped)
}

// OnEvent registers a handler for the given gateway event type.
// handler must have the signature func(*connection.Client, *events.XxxEvent).
// Returns an error if handler has the wrong type; the error should be checked
// by callers — a wrong type means the handler will never be invoked.
// Use the typed helpers (OnMessageCreate, OnInteractionCreate, etc.) for
// compile-time-safe registration that cannot fail.
func (d *Client) OnEvent(eventName events.EventType, handler interface{}) error {
	wrapped, err := toEventHandler(handler, eventName)
	if err != nil {
		return err
	}

	d.onRawEvent(eventName, wrapped)
	return nil
}

func toEventHandler(handler interface{}, eventName events.EventType) (EventHandler, error) {
	if typed, ok := handler.(EventHandler); ok {
		return typed, nil
	}

	handlerValue := reflect.ValueOf(handler)
	if !handlerValue.IsValid() || handlerValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("invalid handler for event %s: expected function", eventName)
	}

	handlerType := handlerValue.Type()
	if handlerType.NumIn() != 2 || handlerType.NumOut() != 0 {
		return nil, fmt.Errorf("invalid handler for event %s: expected func(*connection.Client, *events.X)", eventName)
	}

	clientType := reflect.TypeOf(&Client{})
	if !clientType.AssignableTo(handlerType.In(0)) {
		return nil, fmt.Errorf("invalid handler for event %s: first argument must be *connection.Client", eventName)
	}

	eventType := handlerType.In(1)

	return func(session *Client, event events.Event) {
		eventValue := reflect.ValueOf(event)
		if !eventValue.IsValid() || !eventValue.Type().AssignableTo(eventType) {
			session.Logger.Warn(
				"Failed to cast event to expected type for event",
				slog.Any("event", event), slog.Any("eventName", eventName),
			)
			return
		}

		handlerValue.Call([]reflect.Value{reflect.ValueOf(session), eventValue})
	}, nil
}
