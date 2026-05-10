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
	d.OnEvent(events.EventReady, h)
}
func (d *Client) OnResumed(h func(*Client, *events.ResumedEvent)) {
	d.OnEvent(events.EventResumed, h)
}
func (d *Client) OnGuildCreate(h func(*Client, *events.GuildCreateEvent)) {
	d.OnEvent(events.EventGuildCreate, h)
}
func (d *Client) OnGuildUpdate(h func(*Client, *events.GuildUpdateEvent)) {
	d.OnEvent(events.EventGuildUpdate, h)
}
func (d *Client) OnGuildDelete(h func(*Client, *events.GuildDeleteEvent)) {
	d.OnEvent(events.EventGuildDelete, h)
}
func (d *Client) OnGuildBanAdd(h func(*Client, *events.GuildBanAddEvent)) {
	d.OnEvent(events.EventGuildBanAdd, h)
}
func (d *Client) OnGuildBanRemove(h func(*Client, *events.GuildBanRemoveEvent)) {
	d.OnEvent(events.EventGuildBanRemove, h)
}
func (d *Client) OnGuildEmojisUpdate(h func(*Client, *events.GuildEmojisUpdateEvent)) {
	d.OnEvent(events.EventGuildEmojisUpdate, h)
}
func (d *Client) OnGuildStickersUpdate(h func(*Client, *events.GuildStickersUpdateEvent)) {
	d.OnEvent(events.EventGuildStickersUpdate, h)
}
func (d *Client) OnGuildMemberAdd(h func(*Client, *events.GuildMemberAddEvent)) {
	d.OnEvent(events.EventGuildMemberAdd, h)
}
func (d *Client) OnGuildMemberRemove(h func(*Client, *events.GuildMemberRemoveEvent)) {
	d.OnEvent(events.EventGuildMemberRemove, h)
}
func (d *Client) OnGuildMemberUpdate(h func(*Client, *events.GuildMemberUpdateEvent)) {
	d.OnEvent(events.EventGuildMemberUpdate, h)
}
func (d *Client) OnGuildMembersChunk(h func(*Client, *events.GuildMembersChunkEvent)) {
	d.OnEvent(events.EventGuildMembersChunk, h)
}
func (d *Client) OnGuildRoleCreate(h func(*Client, *events.GuildRoleCreateEvent)) {
	d.OnEvent(events.EventGuildRoleCreate, h)
}
func (d *Client) OnGuildRoleUpdate(h func(*Client, *events.GuildRoleUpdateEvent)) {
	d.OnEvent(events.EventGuildRoleUpdate, h)
}
func (d *Client) OnGuildRoleDelete(h func(*Client, *events.GuildRoleDeleteEvent)) {
	d.OnEvent(events.EventGuildRoleDelete, h)
}
func (d *Client) OnGuildAuditLogEntryCreate(h func(*Client, *events.GuildAuditLogEntryCreateEvent)) {
	d.OnEvent(events.EventGuildAuditLogEntryCreate, h)
}
func (d *Client) OnGuildScheduledEventCreate(h func(*Client, *events.GuildScheduledEventCreateEvent)) {
	d.OnEvent(events.EventGuildScheduledEventCreate, h)
}
func (d *Client) OnGuildScheduledEventUpdate(h func(*Client, *events.GuildScheduledEventUpdateEvent)) {
	d.OnEvent(events.EventGuildScheduledEventUpdate, h)
}
func (d *Client) OnGuildScheduledEventDelete(h func(*Client, *events.GuildScheduledEventDeleteEvent)) {
	d.OnEvent(events.EventGuildScheduledEventDelete, h)
}
func (d *Client) OnGuildScheduledEventUserAdd(h func(*Client, *events.GuildScheduledEventUserAddEvent)) {
	d.OnEvent(events.EventGuildScheduledEventUserAdd, h)
}
func (d *Client) OnGuildScheduledEventUserRemove(h func(*Client, *events.GuildScheduledEventUserRemoveEvent)) {
	d.OnEvent(events.EventGuildScheduledEventUserRemove, h)
}
func (d *Client) OnChannelCreate(h func(*Client, *events.ChannelCreateEvent)) {
	d.OnEvent(events.EventChannelCreate, h)
}
func (d *Client) OnChannelUpdate(h func(*Client, *events.ChannelUpdateEvent)) {
	d.OnEvent(events.EventChannelUpdate, h)
}
func (d *Client) OnChannelDelete(h func(*Client, *events.ChannelDeleteEvent)) {
	d.OnEvent(events.EventChannelDelete, h)
}
func (d *Client) OnChannelPinsUpdate(h func(*Client, *events.ChannelPinsUpdateEvent)) {
	d.OnEvent(events.EventChannelPinsUpdate, h)
}
func (d *Client) OnThreadCreate(h func(*Client, *events.ThreadCreateEvent)) {
	d.OnEvent(events.EventThreadCreate, h)
}
func (d *Client) OnThreadUpdate(h func(*Client, *events.ThreadUpdateEvent)) {
	d.OnEvent(events.EventThreadUpdate, h)
}
func (d *Client) OnThreadDelete(h func(*Client, *events.ThreadDeleteEvent)) {
	d.OnEvent(events.EventThreadDelete, h)
}
func (d *Client) OnThreadListSync(h func(*Client, *events.ThreadListSyncEvent)) {
	d.OnEvent(events.EventThreadListSync, h)
}
func (d *Client) OnThreadMemberUpdate(h func(*Client, *events.ThreadMemberUpdateEvent)) {
	d.OnEvent(events.EventThreadMemberUpdate, h)
}
func (d *Client) OnThreadMembersUpdate(h func(*Client, *events.ThreadMembersUpdateEvent)) {
	d.OnEvent(events.EventThreadMembersUpdate, h)
}
func (d *Client) OnMessageCreate(h func(*Client, *events.MessageCreateEvent)) {
	d.OnEvent(events.EventMessageCreate, h)
}
func (d *Client) OnMessageUpdate(h func(*Client, *events.MessageUpdateEvent)) {
	d.OnEvent(events.EventMessageUpdate, h)
}
func (d *Client) OnMessageDelete(h func(*Client, *events.MessageDeleteEvent)) {
	d.OnEvent(events.EventMessageDelete, h)
}
func (d *Client) OnMessageDeleteBulk(h func(*Client, *events.MessageDeleteBulkEvent)) {
	d.OnEvent(events.EventMessageDeleteBulk, h)
}
func (d *Client) OnMessageReactionAdd(h func(*Client, *events.MessageReactionAddEvent)) {
	d.OnEvent(events.EventMessageReactionAdd, h)
}
func (d *Client) OnMessageReactionRemove(h func(*Client, *events.MessageReactionRemoveEvent)) {
	d.OnEvent(events.EventMessageReactionRemove, h)
}
func (d *Client) OnMessageReactionRemoveAll(h func(*Client, *events.MessageReactionRemoveAllEvent)) {
	d.OnEvent(events.EventMessageReactionRemoveAll, h)
}
func (d *Client) OnMessageReactionRemoveEmoji(h func(*Client, *events.MessageReactionRemoveEmojiEvent)) {
	d.OnEvent(events.EventMessageReactionRemoveEmoji, h)
}
func (d *Client) OnMessagePollVoteAdd(h func(*Client, *events.MessagePollVoteAddEvent)) {
	d.OnEvent(events.EventMessagePollVoteAdd, h)
}
func (d *Client) OnMessagePollVoteRemove(h func(*Client, *events.MessagePollVoteRemoveEvent)) {
	d.OnEvent(events.EventMessagePollVoteRemove, h)
}
func (d *Client) OnInteractionCreate(h func(*Client, *events.InteractionCreateEvent)) {
	d.OnEvent(events.EventInteractionCreate, h)
}
func (d *Client) OnPresenceUpdate(h func(*Client, *events.PresenceUpdateEvent)) {
	d.OnEvent(events.EventPresenceUpdate, h)
}
func (d *Client) OnTypingStart(h func(*Client, *events.TypingStartEvent)) {
	d.OnEvent(events.EventTypingStart, h)
}
func (d *Client) OnUserUpdate(h func(*Client, *events.UserUpdateEvent)) {
	d.OnEvent(events.EventUserUpdate, h)
}
func (d *Client) OnVoiceStateUpdate(h func(*Client, *events.VoiceStateUpdateEvent)) {
	d.OnEvent(events.EventVoiceStateUpdate, h)
}
func (d *Client) OnVoiceServerUpdate(h func(*Client, *events.VoiceServerUpdateEvent)) {
	d.OnEvent(events.EventVoiceServerUpdate, h)
}
func (d *Client) OnVoiceChannelEffectSend(h func(*Client, *events.VoiceChannelEffectSendEvent)) {
	d.OnEvent(events.EventVoiceChannelEffectSend, h)
}
func (d *Client) OnStageInstanceCreate(h func(*Client, *events.StageInstanceCreateEvent)) {
	d.OnEvent(events.EventStageInstanceCreate, h)
}
func (d *Client) OnStageInstanceUpdate(h func(*Client, *events.StageInstanceUpdateEvent)) {
	d.OnEvent(events.EventStageInstanceUpdate, h)
}
func (d *Client) OnStageInstanceDelete(h func(*Client, *events.StageInstanceDeleteEvent)) {
	d.OnEvent(events.EventStageInstanceDelete, h)
}
func (d *Client) OnInviteCreate(h func(*Client, *events.InviteCreateEvent)) {
	d.OnEvent(events.EventInviteCreate, h)
}
func (d *Client) OnInviteDelete(h func(*Client, *events.InviteDeleteEvent)) {
	d.OnEvent(events.EventInviteDelete, h)
}
func (d *Client) OnIntegrationCreate(h func(*Client, *events.IntegrationCreateEvent)) {
	d.OnEvent(events.EventIntegrationCreate, h)
}
func (d *Client) OnIntegrationUpdate(h func(*Client, *events.IntegrationUpdateEvent)) {
	d.OnEvent(events.EventIntegrationUpdate, h)
}
func (d *Client) OnIntegrationDelete(h func(*Client, *events.IntegrationDeleteEvent)) {
	d.OnEvent(events.EventIntegrationDelete, h)
}
func (d *Client) OnWebhooksUpdate(h func(*Client, *events.WebhooksUpdateEvent)) {
	d.OnEvent(events.EventWebhooksUpdate, h)
}
func (d *Client) OnAutoModerationRuleCreate(h func(*Client, *events.AutoModerationRuleCreateEvent)) {
	d.OnEvent(events.EventAutoModerationRuleCreate, h)
}
func (d *Client) OnAutoModerationRuleUpdate(h func(*Client, *events.AutoModerationRuleUpdateEvent)) {
	d.OnEvent(events.EventAutoModerationRuleUpdate, h)
}
func (d *Client) OnAutoModerationRuleDelete(h func(*Client, *events.AutoModerationRuleDeleteEvent)) {
	d.OnEvent(events.EventAutoModerationRuleDelete, h)
}
func (d *Client) OnAutoModerationActionExecution(h func(*Client, *events.AutoModerationActionExecutionEvent)) {
	d.OnEvent(events.EventAutoModerationActionExecution, h)
}
func (d *Client) OnEntitlementCreate(h func(*Client, *events.EntitlementCreateEvent)) {
	d.OnEvent(events.EventEntitlementCreate, h)
}
func (d *Client) OnEntitlementUpdate(h func(*Client, *events.EntitlementUpdateEvent)) {
	d.OnEvent(events.EventEntitlementUpdate, h)
}
func (d *Client) OnEntitlementDelete(h func(*Client, *events.EntitlementDeleteEvent)) {
	d.OnEvent(events.EventEntitlementDelete, h)
}
func (d *Client) OnApplicationCommandPermissionsUpdate(h func(*Client, *events.ApplicationCommandPermissionsUpdateEvent)) {
	d.OnEvent(events.EventApplicationCommandPermissionsUpdate, h)
}
func (d *Client) OnGuildSoundboardSoundCreate(h func(*Client, *events.GuildSoundboardSoundCreateEvent)) {
	d.OnEvent(events.EventGuildSoundboardSoundCreate, h)
}
func (d *Client) OnGuildSoundboardSoundUpdate(h func(*Client, *events.GuildSoundboardSoundUpdateEvent)) {
	d.OnEvent(events.EventGuildSoundboardSoundUpdate, h)
}
func (d *Client) OnGuildSoundboardSoundDelete(h func(*Client, *events.GuildSoundboardSoundDeleteEvent)) {
	d.OnEvent(events.EventGuildSoundboardSoundDelete, h)
}
func (d *Client) OnGuildSoundboardSoundsUpdate(h func(*Client, *events.GuildSoundboardSoundsUpdateEvent)) {
	d.OnEvent(events.EventGuildSoundboardSoundsUpdate, h)
}
func (d *Client) OnSubscriptionCreate(h func(*Client, *events.SubscriptionCreateEvent)) {
	d.OnEvent(events.EventSubscriptionCreate, h)
}
func (d *Client) OnSubscriptionUpdate(h func(*Client, *events.SubscriptionUpdateEvent)) {
	d.OnEvent(events.EventSubscriptionUpdate, h)
}
func (d *Client) OnSubscriptionDelete(h func(*Client, *events.SubscriptionDeleteEvent)) {
	d.OnEvent(events.EventSubscriptionDelete, h)
}
func (d *Client) OnVoiceChannelStatusUpdate(h func(*Client, *events.VoiceChannelStatusUpdateEvent)) {
	d.OnEvent(events.EventVoiceChannelStatusUpdate, h)
}
func (d *Client) OnVoiceChannelStartTimeUpdate(h func(*Client, *events.VoiceChannelStartTimeUpdateEvent)) {
	d.OnEvent(events.EventVoiceChannelStartTimeUpdate, h)
}
func (d *Client) OnSoundboardSounds(h func(*Client, *events.SoundboardSoundsEvent)) {
	d.OnEvent(events.EventSoundboardSounds, h)
}

func (d *Client) onRawEvent(
	eventName events.EventType,
	handler EventHandler,
) {
	wrapped := d.applyMiddleware(handler)
	d.discordEventEmitter.On(eventName, wrapped)
}

func (d *Client) OnEvent(eventName events.EventType, handler interface{}) {
	wrapped, err := toEventHandler(handler, eventName)
	if err != nil {
		d.Logger.Warn(err.Error())
		return
	}

	d.onRawEvent(eventName, wrapped)
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
