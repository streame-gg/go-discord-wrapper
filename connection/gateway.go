package connection

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"sync"
	"sync/atomic"

	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/common"
	"github.com/streame-gg/go-discord-wrapper/types/events"
	"github.com/streame-gg/go-discord-wrapper/util"

	"github.com/gorilla/websocket"
)

type EventHandler func(*Client, events.Event)

type Client struct {
	token *string

	APIVersion *common.APIVersion

	Logger *slog.Logger

	Intents *common.Intent

	Websocket *Websocket

	events map[events.EventType][]EventHandler

	discordEventEmitter *util.EventEmitter[events.EventType, EventHandler]

	mu sync.RWMutex

	reconnectMu  sync.Mutex
	reconnecting bool

	// shutdown is set to true by Shutdown() and prevents reconnect attempts.
	shutdown atomic.Bool

	// maxReconnectRetries mirrors options.Config.MaxReconnectRetries.
	maxReconnectRetries int

	// disableCacheAutoPopulation mirrors options.Config.DisableCacheAutoPopulation.
	disableCacheAutoPopulation bool

	UnavailableGuilds map[common.Snowflake]struct{}

	User *common.User

	Sharding *options.Sharding

	RestClient *api.RestClient

	// Coordinator handles inter-shard message routing.
	// Nil when sharding is not configured or no coordinator was provided.
	Coordinator options.ShardCoordinator

	// Cache stores Discord entities populated automatically from gateway events.
	// Nil when no cache was configured via options.WithCache.
	Cache cache.Cache

	// guildMemberCounts tracks the member_count for every available guild on
	// this shard, updated from GUILD_CREATE / GUILD_DELETE gateway events.
	guildMemberCounts map[common.Snowflake]int
	guildMu           sync.RWMutex

	// voiceStates caches the latest voice state per (guildID:userID) key so
	// VoiceStateUpdateEvent can carry OldState alongside the new state.
	voiceStates   map[string]*common.VoiceState
	voiceStatesMu sync.RWMutex

	// Client lifecycle event handlers.
	onConnect      []func(*Client)
	onDisconnect   []func(*Client, error)
	onReconnect    []func(*Client)
	onPacketError  []func(*Client, error)
	clientEventsMu sync.RWMutex

	// shardHandlers are persistent handlers registered via OnShardMessage.
	shardHandlers   []func(options.ShardMessage)
	shardHandlersMu sync.RWMutex

	// pendingShardReqs maps "responseType:corrID" to a buffered channel that
	// RequestAll uses to collect responses without holding a permanent handler.
	pendingShardReqs   map[string]chan options.ShardMessage
	pendingShardReqsMu sync.Mutex
}

// NewClient creates a new Discord gateway client.
// Configure it with options from the options package:
//
//	connection.NewClient(token, intents,
//	    options.WithSharding(4, 0),
//	    options.WithCoordinator(coord),
//	    options.WithLogger(&myLogger),
//	)
func NewClient(token string, intents common.Intent, opts ...options.Option) *Client {
	cfg := options.Build(options.Config{APIVersion: common.APIVersion10}, opts)
	if err := cfg.Validate(); err != nil {
		panic("go-discord-wrapper: " + err.Error())
	}

	c := &Client{
		token:                      &token,
		APIVersion:                 util.PointerOf(cfg.APIVersion),
		Intents:                    &intents,
		discordEventEmitter:        util.NewEventEmitter[events.EventType, EventHandler](),
		UnavailableGuilds:          make(map[common.Snowflake]struct{}),
		guildMemberCounts:          make(map[common.Snowflake]int),
		voiceStates:                make(map[string]*common.VoiceState),
		RestClient:                 api.NewRestClient(token, opts...),
		Sharding:                   cfg.Sharding,
		Coordinator:                cfg.Coordinator,
		Cache:                      cfg.Cache,
		maxReconnectRetries:        cfg.MaxReconnectRetries,
		disableCacheAutoPopulation: cfg.DisableCacheAutoPopulation,
	}

	switch {
	case cfg.Logger != nil:
		c.Logger = cfg.Logger
	case cfg.LogLevel != nil:
		c.Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: *cfg.LogLevel}))
	default:
		c.Logger = slog.Default()
	}

	// Register this shard with the coordinator so messages can arrive.
	if c.Coordinator != nil && c.Sharding != nil {
		if err := c.Coordinator.Register(c.Sharding.ShardID, c.dispatchShardMessage); err != nil {
			c.Logger.Error("Failed to register shard with coordinator", slog.Int("shardId", c.Sharding.ShardID), slog.Any("err", err))
		}
	}

	return c
}

// ── Shard messaging ─────────────────────────────────────────────────────────

// OnShardMessage registers a persistent handler for messages received from
// other shards. Multiple handlers may be registered; all are called concurrently.
func (d *Client) OnShardMessage(handler func(options.ShardMessage)) {
	d.shardHandlersMu.Lock()
	d.shardHandlers = append(d.shardHandlers, handler)
	d.shardHandlersMu.Unlock()
}

// SendToShard sends a typed message to a single shard.
// payload is JSON-marshalled automatically.
func (d *Client) SendToShard(shardID int, msgType string, payload interface{}) error {
	if d.Coordinator == nil {
		return errors.New("no shard coordinator configured; use options.WithCoordinator")
	}
	if d.Sharding == nil {
		return errors.New("sharding is not configured; use options.WithSharding")
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal shard payload: %w", err)
	}
	return d.Coordinator.Send(options.ShardMessage{
		Type:    msgType,
		From:    d.Sharding.ShardID,
		To:      shardID,
		Payload: b,
	})
}

// BroadcastToShards sends a typed message to every shard including self.
// payload is JSON-marshalled automatically.
func (d *Client) BroadcastToShards(msgType string, payload interface{}) error {
	if d.Coordinator == nil {
		return errors.New("no shard coordinator configured; use options.WithCoordinator")
	}
	if d.Sharding == nil {
		return errors.New("sharding is not configured; use options.WithSharding")
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal shard payload: %w", err)
	}
	return d.Coordinator.Broadcast(options.ShardMessage{
		Type:    msgType,
		From:    d.Sharding.ShardID,
		To:      options.BroadcastAll,
		Payload: b,
	})
}

// ReplyToShard sends a response back to the originator of msg, automatically
// echoing the CorrelationID so RequestAll can match it.
// payload is JSON-marshalled automatically.
func (d *Client) ReplyToShard(original options.ShardMessage, responseType string, payload interface{}) error {
	if d.Coordinator == nil {
		return errors.New("no shard coordinator configured; use options.WithCoordinator")
	}
	if d.Sharding == nil {
		return errors.New("sharding is not configured; use options.WithSharding")
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal shard reply payload: %w", err)
	}
	return d.Coordinator.Send(options.ShardMessage{
		Type:          responseType,
		From:          d.Sharding.ShardID,
		To:            original.From,
		CorrelationID: original.CorrelationID,
		Payload:       b,
	})
}

// SubscribeShardResponse registers a temporary response channel for the given
// responseType and corrID combination. The returned cleanup function must be
// called (e.g. via defer) to release resources.
//
// This is a low-level building block used by sharding.RequestAll. Prefer that
// helper unless you need custom aggregation logic.
func (d *Client) SubscribeShardResponse(responseType, corrID string) (<-chan options.ShardMessage, func()) {
	ch := make(chan options.ShardMessage, d.Coordinator.TotalShards())
	key := responseType + ":" + corrID

	d.pendingShardReqsMu.Lock()
	if d.pendingShardReqs == nil {
		d.pendingShardReqs = make(map[string]chan options.ShardMessage)
	}
	d.pendingShardReqs[key] = ch
	d.pendingShardReqsMu.Unlock()

	cleanup := func() {
		d.pendingShardReqsMu.Lock()
		delete(d.pendingShardReqs, key)
		d.pendingShardReqsMu.Unlock()
	}
	return ch, cleanup
}

// ── Per-shard statistics ─────────────────────────────────────────────────────

// GuildCount returns the number of guilds this shard is currently connected to.
// The value is updated automatically from GUILD_CREATE and GUILD_DELETE events.
// It may be lower than the eventual total immediately after Login() returns, as
// Discord streams GUILD_CREATE events asynchronously after the READY payload.
func (d *Client) GuildCount() int {
	d.guildMu.RLock()
	defer d.guildMu.RUnlock()
	return len(d.guildMemberCounts)
}

// MemberCount returns the sum of member counts across all guilds on this shard.
//
// Important: this counts member slots, not unique users. A user who belongs to
// three guilds on this shard is counted three times. Use sharding.GetTotalUsers
// for a cross-shard aggregate with the same semantics.
func (d *Client) MemberCount() int {
	d.guildMu.RLock()
	defer d.guildMu.RUnlock()
	total := 0
	for _, n := range d.guildMemberCounts {
		total += n
	}
	return total
}

// dispatchShardMessage is the handler registered with the coordinator.
// It first handles built-in stat requests, then routes correlated responses to
// pending RequestAll channels, and finally fans out to OnShardMessage handlers.
func (d *Client) dispatchShardMessage(msg options.ShardMessage) {
	// Handle built-in cross-shard stat requests automatically so that
	// sharding.GetTotalServers and sharding.GetTotalUsers work on every shard
	// without any manual OnShardMessage registration.
	switch msg.Type {
	case options.ShardMsgServerCountReq:
		if err := d.ReplyToShard(msg, options.ShardMsgServerCountResp, d.GuildCount()); err != nil {
			d.Logger.Error("Failed to reply to server count request", slog.Any("err", err))
		}
		return
	case options.ShardMsgMemberCountReq:
		if err := d.ReplyToShard(msg, options.ShardMsgMemberCountResp, d.MemberCount()); err != nil {
			d.Logger.Error("Failed to reply to member count request", slog.Any("err", err))
		}
		return
	}

	// Route to a pending request-response channel if this is a correlated response.
	if msg.CorrelationID != "" {
		key := msg.Type + ":" + msg.CorrelationID
		d.pendingShardReqsMu.Lock()
		ch, ok := d.pendingShardReqs[key]
		d.pendingShardReqsMu.Unlock()
		if ok {
			// Non-blocking send: if the buffer is full the response is dropped.
			// This can happen when more shards respond than TotalShards() reports.
			select {
			case ch <- msg:
			default:
			}
			return
		}
	}

	// Fan out to all registered persistent handlers.
	d.shardHandlersMu.RLock()
	handlers := make([]func(options.ShardMessage), len(d.shardHandlers))
	copy(handlers, d.shardHandlers)
	d.shardHandlersMu.RUnlock()

	for _, h := range handlers {
		h := h
		go h(msg)
	}
}

// ── Gateway connection ───────────────────────────────────────────────────────

func (d *Client) initializeGatewayConnection() (*common.BotRegisterResponse, error) {
	do, err := http.DefaultClient.Do(&http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "https",
			Host:   "discord.com",
			Path:   common.APIBaseString(*d.APIVersion) + "gateway/bot",
		},
		Header: http.Header{
			"Authorization": []string{"Bot " + *d.token},
		},
	})

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = do.Body.Close()
	}()

	if do.StatusCode != http.StatusOK {
		return nil, errors.New("failed to register bot gateway connection, status code: " + do.Status)
	}

	var resp common.BotRegisterResponse
	if err := json.NewDecoder(do.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (d *Client) Login() error {
	gatewayResp, err := d.initializeGatewayConnection()
	if err != nil {
		return err
	}

	d.Logger.Debug("Connecting to gateway websocket", slog.String("url", gatewayResp.Url), slog.Int("shards", gatewayResp.Shards))

	if err := d.connectWebsocket(gatewayResp.Url, false, nil, nil); err != nil {
		return err
	}

	go func() {
		for {
			if err := d.listenWebsocket(); err != nil {
				// Shutdown() was called — do not reconnect.
				if d.shutdown.Load() {
					d.Logger.Info("Shutting down gracefully")
					return
				}

				d.Logger.Error("Error listening to websocket", slog.Any("err", err))
				d.emitDisconnect(err)

				if d.Websocket == nil {
					d.Logger.Debug("Websocket is nil, stopping listener")
					return
				}

				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					d.Logger.Debug("Gateway connection closed normally")
					return
				}

				// 4014: privileged intent not enabled — reconnecting will not help.
				if websocket.IsCloseError(err, 4014) {
					d.Logger.Error("A privileged intent is not enabled in the Discord developer portal; cannot reconnect")
					return
				}

				if websocket.IsCloseError(err, 4000, 4001, 4002, 4003, 4005, 4007, 4008, 4009) {
					d.Logger.Debug("Gateway connection closed by Discord, trying to reconnect")
					if err := d.reconnect(true); err != nil {
						d.Logger.Error("Failed to reconnect", slog.Any("err", err))
						return
					}

					continue
				}

				if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) ||
					websocket.IsUnexpectedCloseError(err) {
					d.Logger.Warn("Abnormal websocket closure, attempting to reconnect")
					if err := d.reconnect(false); err != nil {
						d.Logger.Error("Failed to reconnect after abnormal closure", slog.Any("err", err))
						if err := d.reconnect(true); err != nil {
							d.Logger.Error("Failed to reconnect with fresh connection", slog.Any("err", err))
							return
						}
					}
					continue
				}

				d.Logger.Warn("Unexpected error, attempting to reconnect")
				if err := d.reconnect(false); err != nil {
					d.Logger.Warn("Resume failed, attempting fresh reconnect")
					if err := d.reconnect(true); err != nil {
						d.Logger.Error("Fresh reconnect failed, stopping listener", slog.Any("err", err))
						return
					}
				}
				continue
			}

			if d.Websocket != nil {
				d.Logger.Debug("Restarting websocket listener")
				continue
			}

			return
		}
	}()

	<-d.Websocket.Ready

	d.Logger.Info("Successfully connected to the Discord gateway")
	d.emitConnect()

	return nil
}

func (d *Client) dispatch(event events.Event) {
	handlers := d.discordEventEmitter.Handlers(event.Event())
	if len(handlers) == 0 {
		d.mu.RLock()
		handlers = d.events[event.Event()]
		d.mu.RUnlock()
	}

	// Each handler runs in its own goroutine so a slow handler cannot stall
	// the gateway event loop or block subsequent event processing.
	for _, h := range handlers {
		go h(d, event)
	}
}

func (d *Client) cacheEnabled() bool {
	return d.Cache != nil && !d.disableCacheAutoPopulation
}

func (d *Client) internalEventHandler(msg json.RawMessage, eventType events.EventType, event events.Event) bool {
	switch eventType {
	case events.EventReady:
		{
			var readyEvent events.ReadyEvent
			if err := json.Unmarshal(msg, &readyEvent); err != nil {
				d.Logger.Error("Failed to unmarshal READY event", slog.Any("err", err))
			}

			d.Websocket.SessionID = &readyEvent.SessionID
			d.Websocket.ReconnectURL = &readyEvent.ResumeGatewayURL
			d.User = &readyEvent.User

			if readyEvent.Shard != nil {
				d.Logger.Debug("Connected to shard", slog.Int("shard", readyEvent.Shard[0]+1), slog.Int("total", readyEvent.Shard[1]))
			}

			for _, guild := range readyEvent.Guilds {
				if !guild.Guild.IsAvailable() {
					d.addUnavailableGuild(guild.Guild.GetID())
				}
			}

			close(d.Websocket.Ready)

			return true
		}
	case events.EventGuildCreate:
		{
			var guildCreateEvent events.GuildCreateEvent
			if err := json.Unmarshal(msg, &guildCreateEvent); err != nil {
				d.Logger.Error("Failed to unmarshal GUILD_CREATE event", slog.Any("err", err))
				return false
			}

			if guildCreateEvent.Guild.IsAvailable() {
				// Record (or refresh) the member count for this guild.
				d.guildMu.Lock()
				d.guildMemberCounts[guildCreateEvent.Guild.GetID()] = guildCreateEvent.MemberCount
				d.guildMu.Unlock()

				if d.cacheEnabled() {
					if g, ok := guildCreateEvent.Guild.(common.Guild); ok {
						gcopy := g
						d.Cache.Guilds().Set(&gcopy)
					}
				}
			}

			if guildCreateEvent.Guild.IsAvailable() && d.IsGuildUnavailable(guildCreateEvent.Guild.GetID()) {
				// Discord fires GUILD_CREATE with available=true to signal a guild is back online.
				d.Logger.Debug("Guild is available again", slog.Any("guildId", guildCreateEvent.Guild.GetID()))
				d.deleteUnavailableGuild(guildCreateEvent.Guild.GetID())

				return false
			}
		}
	case events.EventGuildDelete:
		{
			var guildDeleteEvent events.GuildDeleteEvent
			if err := json.Unmarshal(msg, &guildDeleteEvent); err != nil {
				d.Logger.Error("Failed to unmarshal GUILD_DELETE event", slog.Any("err", err))
				return false
			}

			if guildDeleteEvent.Unavailable != nil && *guildDeleteEvent.Unavailable {
				// Temporary outage — guild is still ours, just unreachable right now.
				// Keep the member count in the cache so GuildCount stays accurate.
				if d.IsGuildUnavailable(guildDeleteEvent.ID) {
					return false
				}

				d.addUnavailableGuild(guildDeleteEvent.ID)

				d.Logger.Debug("Guild became unavailable", slog.Any("guildId", guildDeleteEvent.ID))

				return false
			}

			// Bot was kicked or the guild was deleted.
			d.guildMu.Lock()
			delete(d.guildMemberCounts, guildDeleteEvent.ID)
			d.guildMu.Unlock()

			if d.cacheEnabled() {
				d.Cache.Guilds().Delete(guildDeleteEvent.ID)
				d.Cache.Members().DeleteGuild(guildDeleteEvent.ID)
			}
		}
	case events.EventVoiceStateUpdate:
		{
			if vse, ok := event.(*events.VoiceStateUpdateEvent); ok && vse.GuildID != nil {
				key := string(*vse.GuildID) + ":" + string(vse.UserID)
				d.voiceStatesMu.RLock()
				oldState := d.voiceStates[key]
				d.voiceStatesMu.RUnlock()
				vse.OldState = oldState

				d.voiceStatesMu.Lock()
				if vse.ChannelID == nil {
					delete(d.voiceStates, key)
				} else {
					s := vse.VoiceState
					d.voiceStates[key] = &s
				}
				d.voiceStatesMu.Unlock()
			}
		}
	case events.EventUserUpdate:
		{
			if d.cacheEnabled() {
				var ev events.UserUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal USER_UPDATE event", slog.Any("err", err))
					return false
				}
				u := ev.User
				d.Cache.Users().Set(&u)
			}
		}
	case events.EventGuildUpdate:
		{
			if d.cacheEnabled() {
				var ev events.GuildUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal GUILD_UPDATE event", slog.Any("err", err))
					return false
				}
				g := ev.Guild
				d.Cache.Guilds().Set(&g)
			}
		}
	case events.EventGuildMemberAdd:
		{
			var ev events.GuildMemberAddEvent
			if err := json.Unmarshal(msg, &ev); err != nil {
				d.Logger.Error("Failed to unmarshal GUILD_MEMBER_ADD event", slog.Any("err", err))
				return false
			}
			d.guildMu.Lock()
			d.guildMemberCounts[ev.GuildID]++
			d.guildMu.Unlock()
			if d.cacheEnabled() {
				m := ev.GuildMember
				d.Cache.Members().Set(ev.GuildID, &m)
			}
		}
	case events.EventGuildMemberRemove:
		{
			var ev events.GuildMemberRemoveEvent
			if err := json.Unmarshal(msg, &ev); err != nil {
				d.Logger.Error("Failed to unmarshal GUILD_MEMBER_REMOVE event", slog.Any("err", err))
				return false
			}
			d.guildMu.Lock()
			if d.guildMemberCounts[ev.GuildID] > 0 {
				d.guildMemberCounts[ev.GuildID]--
			}
			d.guildMu.Unlock()
			if d.cacheEnabled() {
				d.Cache.Members().Delete(ev.GuildID, ev.User.ID)
			}
		}
	case events.EventGuildMemberUpdate:
		{
			if d.cacheEnabled() {
				var ev events.GuildMemberUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal GUILD_MEMBER_UPDATE event", slog.Any("err", err))
					return false
				}
				m := common.GuildMember{
					User:                       &ev.User,
					Nick:                       ev.Nick,
					AvatarHash:                 ev.AvatarHash,
					PremiumSince:               ev.PremiumSince,
					CommunicationDisabledUntil: ev.CommunicationDisabledUntil,
				}
				if ev.JoinedAt != nil {
					m.JoinedAt = *ev.JoinedAt
				}
				if ev.Deaf != nil {
					m.Deaf = *ev.Deaf
				}
				if ev.Mute != nil {
					m.Mute = *ev.Mute
				}
				if ev.Pending != nil {
					m.Pending = *ev.Pending
				}
				if ev.Flags != nil {
					m.Flags = *ev.Flags
				}
				for _, r := range ev.Roles {
					m.Roles = append(m.Roles, r.String())
				}
				d.Cache.Members().Set(ev.GuildID, &m)
			}
		}
	case events.EventChannelCreate:
		{
			if d.cacheEnabled() {
				var ev events.ChannelCreateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal CHANNEL_CREATE event", slog.Any("err", err))
					return false
				}
				ch := ev.Channel
				d.Cache.Channels().Set(&ch)
			}
		}
	case events.EventChannelUpdate:
		{
			if d.cacheEnabled() {
				var ev events.ChannelUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal CHANNEL_UPDATE event", slog.Any("err", err))
					return false
				}
				ch := ev.Channel
				d.Cache.Channels().Set(&ch)
			}
		}
	case events.EventChannelDelete:
		{
			if d.cacheEnabled() {
				var ev events.ChannelDeleteEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal CHANNEL_DELETE event", slog.Any("err", err))
					return false
				}
				d.Cache.Channels().Delete(ev.ID)
				d.Cache.Messages().DeleteChannel(ev.ID)
			}
		}
	case events.EventMessageCreate:
		{
			if d.cacheEnabled() {
				var ev events.MessageCreateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal MESSAGE_CREATE event", slog.Any("err", err))
					return false
				}
				msg := ev.Message
				d.Cache.Messages().Add(&msg)
				if ev.Author != nil {
					d.Cache.Users().Set(ev.Author)
				}
			}
		}
	case events.EventMessageUpdate:
		{
			if d.cacheEnabled() {
				var ev events.MessageUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal MESSAGE_UPDATE event", slog.Any("err", err))
					return false
				}
				msg := ev.Message
				d.Cache.Messages().Update(&msg)
			}
		}
	case events.EventMessageDelete:
		{
			if d.cacheEnabled() {
				var ev events.MessageDeleteEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal MESSAGE_DELETE event", slog.Any("err", err))
					return false
				}
				d.Cache.Messages().Delete(ev.ChannelID, ev.ID)
			}
		}
	case events.EventMessageDeleteBulk:
		{
			if d.cacheEnabled() {
				var ev events.MessageDeleteBulkEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal MESSAGE_DELETE_BULK event", slog.Any("err", err))
					return false
				}
				d.Cache.Messages().DeleteBulk(ev.ChannelID, ev.IDs)
			}
		}
	default:
		return true
	}

	return true
}

func (d *Client) addUnavailableGuild(id common.Snowflake) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.UnavailableGuilds[id] = struct{}{}
}

func (d *Client) deleteUnavailableGuild(id common.Snowflake) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.UnavailableGuilds, id)
}

func (d *Client) IsGuildUnavailable(id common.Snowflake) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, exists := d.UnavailableGuilds[id]
	return exists
}

// Close gracefully closes the gateway connection and the entity cache.
// It implements io.Closer.
func (d *Client) Close() error {
	return d.Shutdown()
}

// Shutdown gracefully closes the gateway connection and the entity cache.
// It returns a non-nil error if either close operation fails.
func (d *Client) Shutdown() error {
	d.shutdown.Store(true)
	var errs []error
	if d.Websocket != nil {
		if err := d.Websocket.Connection.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if d.Cache != nil {
		if err := d.Cache.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// ── Client lifecycle event emitters ─────────────────────────────────────────

func (d *Client) emitConnect() {
	d.clientEventsMu.RLock()
	handlers := make([]func(*Client), len(d.onConnect))
	copy(handlers, d.onConnect)
	d.clientEventsMu.RUnlock()
	for _, h := range handlers {
		go h(d)
	}
}

func (d *Client) emitDisconnect(err error) {
	d.clientEventsMu.RLock()
	handlers := make([]func(*Client, error), len(d.onDisconnect))
	copy(handlers, d.onDisconnect)
	d.clientEventsMu.RUnlock()
	for _, h := range handlers {
		go h(d, err)
	}
}

func (d *Client) emitReconnect() {
	d.clientEventsMu.RLock()
	handlers := make([]func(*Client), len(d.onReconnect))
	copy(handlers, d.onReconnect)
	d.clientEventsMu.RUnlock()
	for _, h := range handlers {
		go h(d)
	}
}

func (d *Client) emitPacketError(err error) {
	d.clientEventsMu.RLock()
	handlers := make([]func(*Client, error), len(d.onPacketError))
	copy(handlers, d.onPacketError)
	d.clientEventsMu.RUnlock()
	for _, h := range handlers {
		go h(d, err)
	}
}
