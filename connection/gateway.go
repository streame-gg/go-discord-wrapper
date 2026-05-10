// Package connection provides the Discord gateway (WebSocket) client, including
// automatic reconnection, session resume, event dispatch, and middleware support.
package connection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

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

	discordEventEmitter *util.EventEmitter[events.EventType, EventHandler]

	mu sync.RWMutex

	// wsMu protects concurrent reads and writes of the Websocket pointer.
	wsMu sync.RWMutex

	// dispatchWg tracks in-flight user event handler goroutines for graceful shutdown.
	dispatchWg sync.WaitGroup

	// eventCh is non-nil when MaxConcurrentEvents > 0 (worker-pool mode).
	// Events are enqueued here and drained by a fixed set of worker goroutines.
	// When nil, each event gets its own goroutine (unlimited mode — the default).
	eventCh chan dispatchJob

	// workerWg tracks the fixed event-worker goroutines started by NewClient.
	workerWg sync.WaitGroup

	reconnectMu  sync.Mutex
	reconnecting bool

	// shutdown is set to true by Shutdown() and prevents reconnect attempts.
	shutdown atomic.Bool

	// maxReconnectRetries mirrors options.Config.MaxReconnectRetries.
	maxReconnectRetries int

	// cacheAutoPopulate mirrors options.Config.CacheStores, zero disables auto-population.
	cacheAutoPopulate cache.OverflowCategory

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

	// channelIndexMu protects channelsByGuild and guildByChannel.
	// Both maps must always be updated together to keep the bidirectional
	// index consistent; acquire channelIndexMu (write lock) before any write.
	channelIndexMu sync.RWMutex
	// channelsByGuild maps each guild ID to the set of channel IDs that
	// belong to it. Used to efficiently evict all guild channels on GUILD_DELETE.
	channelsByGuild map[common.Snowflake]map[common.Snowflake]struct{}
	// guildByChannel is the reverse mapping: channel ID → guild ID.
	guildByChannel map[common.Snowflake]common.Snowflake

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

	// middleware is the ordered chain applied to every event handler registered
	// after a Use call. Middleware registered first wraps outermost.
	middleware   []Middleware
	middlewareMu sync.RWMutex
}

// NewClient creates a new Discord gateway client.
// Configure it with options from the options package:
//
//	bot, err := connection.NewClient(token, intents,
//	    options.WithSharding(4, 0),
//	    options.WithCoordinator(coord),
//	    options.WithLogger(&myLogger),
//	)
func NewClient(token string, intents common.Intent, opts ...options.Option) (*Client, error) {
	cfg := options.Build(options.Config{
		APIVersion:  common.APIVersion10,
		CacheStores: cache.CategoryAll,
	}, opts)
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("go-discord-wrapper: %w", err)
	}

	rc, err := api.NewRestClient(token, opts...)
	if err != nil {
		return nil, err
	}

	cacheStores := cfg.CacheStores
	if cfg.DisableCacheAutoPopulation {
		cacheStores = 0
	}

	c := &Client{
		token:               &token,
		APIVersion:          util.PointerOf(cfg.APIVersion),
		Intents:             &intents,
		discordEventEmitter: util.NewEventEmitter[events.EventType, EventHandler](),
		UnavailableGuilds:   make(map[common.Snowflake]struct{}),
		guildMemberCounts:   make(map[common.Snowflake]int),
		voiceStates:         make(map[string]*common.VoiceState),
		channelsByGuild:     make(map[common.Snowflake]map[common.Snowflake]struct{}),
		guildByChannel:      make(map[common.Snowflake]common.Snowflake),
		RestClient:          rc,
		Sharding:            cfg.Sharding,
		Coordinator:         cfg.Coordinator,
		Cache:               cfg.Cache,
		maxReconnectRetries: cfg.MaxReconnectRetries,
		cacheAutoPopulate:   cacheStores,
	}

	if cfg.MaxConcurrentEvents > 0 {
		n := cfg.MaxConcurrentEvents
		queueDepth := n * 4
		if queueDepth < 256 {
			queueDepth = 256
		}
		c.eventCh = make(chan dispatchJob, queueDepth)
		for i := 0; i < n; i++ {
			c.workerWg.Add(1)
			go c.runEventWorker()
		}
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

	return c, nil
}

// dispatchJob carries a single gateway event through the worker pool.
type dispatchJob struct {
	rawPayload json.RawMessage
	event      events.Event
}

// processEvent runs internalEventHandler and dispatch for one gateway event.
// It is called both by worker goroutines (pool mode) and directly spawned goroutines (unlimited mode).
func (d *Client) processEvent(job dispatchJob) {
	defer func() {
		if r := recover(); r != nil {
			d.Logger.Error("panic in event dispatch", slog.Any("recover", r))
		}
	}()
	if canContinue := d.internalEventHandler(job.rawPayload, job.event.Event(), job.event); canContinue {
		d.dispatch(job.event)
	}
}

// runEventWorker drains d.eventCh until it is closed by Shutdown.
func (d *Client) runEventWorker() {
	defer d.workerWg.Done()
	for job := range d.eventCh {
		d.processEvent(job)
	}
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
	httpClient := &http.Client{Timeout: 10 * time.Second}
	do, err := httpClient.Do(&http.Request{
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

func (d *Client) Login(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	gatewayResp, err := d.initializeGatewayConnection()
	if err != nil {
		return err
	}

	d.Logger.Debug("Connecting to gateway websocket", slog.String("url", gatewayResp.Url), slog.Int("shards", gatewayResp.Shards))

	if err := d.connectWebsocket(gatewayResp.Url, false, nil, nil); err != nil {
		return err
	}

	if ctx.Done() != nil {
		go func() {
			<-ctx.Done()
			_ = d.Shutdown()
		}()
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

				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					d.Logger.Debug("Gateway connection closed normally, attempting resume")
					if err := d.reconnect(false); err != nil {
						d.Logger.Warn("Resume failed after normal close, attempting fresh reconnect")
						if err := d.reconnect(true); err != nil {
							d.Logger.Error("Failed to reconnect after normal close", slog.Any("err", err))
							return
						}
					}
					continue
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

			d.wsMu.RLock()
			wsAlive := d.Websocket != nil
			d.wsMu.RUnlock()
			if wsAlive {
				d.Logger.Debug("Restarting websocket listener")
				continue
			}

			return
		}
	}()

	d.wsMu.RLock()
	ready := d.Websocket.Ready
	d.wsMu.RUnlock()
	<-ready

	d.Logger.Info("Successfully connected to the Discord gateway")
	d.emitConnect()

	return nil
}

func (d *Client) dispatch(event events.Event) {
	handlers := d.discordEventEmitter.Handlers(event.Event())

	// Each handler runs in its own goroutine so a slow handler cannot stall
	// the gateway event loop or block subsequent event processing.
	// dispatchWg tracks in-flight handlers so Shutdown can drain them.
	for _, h := range handlers {
		h := h
		d.dispatchWg.Add(1)
		go func() {
			defer d.dispatchWg.Done()
			defer func() {
				if r := recover(); r != nil {
					d.Logger.Error("panic in event handler", slog.Any("recover", r))
				}
			}()
			h(d, event)
		}()
	}
}

func (d *Client) cacheEnabled() bool {
	return d.Cache != nil && d.cacheAutoPopulate != 0
}

func (d *Client) cacheStoreEnabled(category cache.OverflowCategory) bool {
	return d.Cache != nil && d.cacheAutoPopulate&category != 0
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

			if len(readyEvent.Shard) >= 2 {
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
	case events.EventResumed:
		{
			d.Logger.Info("Gateway session resumed")
			d.emitReconnect()
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
				memberCount := guildCreateEvent.MemberCount
				var guild common.Guild
				var gatewayGuild common.GatewayGuild
				hasGateway := false
				switch g := guildCreateEvent.Guild.(type) {
				case common.GatewayGuild:
					guild = g.Guild
					gatewayGuild = g
					hasGateway = true
					if memberCount == 0 && g.MemberCount > 0 {
						memberCount = g.MemberCount
					}
				case common.Guild:
					guild = g
				}

				// Record (or refresh) the member count for this guild.
				d.guildMu.Lock()
				d.guildMemberCounts[guildCreateEvent.Guild.GetID()] = memberCount
				d.guildMu.Unlock()

				if hasGateway {
					guildID := guildCreateEvent.Guild.GetID()
					d.voiceStatesMu.Lock()
					for i := range gatewayGuild.VoiceStates {
						vs := gatewayGuild.VoiceStates[i]
						key := string(guildID) + ":" + string(vs.UserID)
						vscopy := vs
						if vscopy.GuildID == nil {
							vscopy.GuildID = &guildID
						}
						d.voiceStates[key] = &vscopy
					}
					d.voiceStatesMu.Unlock()
				}

				if d.cacheEnabled() {
					if guild.ID != "" && d.cacheStoreEnabled(cache.CategoryGuilds) {
						gcopy := guild
						d.Cache.Guilds().Set(&gcopy)
					}

					if guild.ID != "" && d.cacheStoreEnabled(cache.CategoryRoles) {
						for i := range guild.Roles {
							role := guild.Roles[i]
							d.Cache.Roles().Set(guild.ID, &role)
						}
					}

					if guild.ID != "" && d.cacheStoreEnabled(cache.CategoryEmojis) {
						for i := range guild.Emojis {
							emoji := guild.Emojis[i]
							if emoji.ID != "" {
								d.Cache.Emojis().Set(guild.ID, &emoji)
							}
						}
					}
					if guild.ID != "" && d.cacheStoreEnabled(cache.CategoryStickers) && guild.Stickers != nil {
						for i := range *guild.Stickers {
							sticker := (*guild.Stickers)[i]
							if sticker.ID != "" {
								d.Cache.Stickers().Set(guild.ID, &sticker)
							}
						}
					}

					if hasGateway {
						guildID := guildCreateEvent.Guild.GetID()
						if d.cacheStoreEnabled(cache.CategoryChannels) {
							gid := guildCreateEvent.Guild.GetID()
							for i := range gatewayGuild.Channels {
								ch := gatewayGuild.Channels[i]
								if ch.GuildID == nil {
									ch.GuildID = &gid
								}
								d.cacheChannel(&ch)
							}
							for i := range gatewayGuild.Threads {
								ch := gatewayGuild.Threads[i]
								if ch.GuildID == nil {
									ch.GuildID = &gid
								}
								d.cacheChannel(&ch)
							}
						}

						if d.cacheStoreEnabled(cache.CategoryMembers) {
							for i := range gatewayGuild.Members {
								m := gatewayGuild.Members[i]
								d.Cache.Members().Set(guildID, &m)
							}
						}

						if d.cacheStoreEnabled(cache.CategoryUsers) {
							for i := range gatewayGuild.Members {
								if gatewayGuild.Members[i].User == nil {
									continue
								}
								u := *gatewayGuild.Members[i].User
								d.Cache.Users().Set(&u)
							}
						}

						if d.cacheStoreEnabled(cache.CategoryVoiceStates) {
							for i := range gatewayGuild.VoiceStates {
								vs := gatewayGuild.VoiceStates[i]
								if vs.GuildID == nil {
									vs.GuildID = &guildID
								}
								d.Cache.VoiceStates().Set(guildID, &vs)
							}
						}

						if d.cacheStoreEnabled(cache.CategoryPresences) {
							for i := range gatewayGuild.Presences {
								p := gatewayGuild.Presences[i]
								presence := common.Presence{
									User:         p.User,
									GuildID:      guildID,
									Status:       p.Status,
									Activities:   p.Activities,
									ClientStatus: p.ClientStatus,
								}
								d.Cache.Presences().Set(&presence)
							}
						}

						if d.cacheStoreEnabled(cache.CategorySoundboard) {
							for i := range gatewayGuild.SoundboardSounds {
								sound := gatewayGuild.SoundboardSounds[i]
								d.Cache.Soundboard().Set(guildID, &sound)
							}
						}

						if d.cacheStoreEnabled(cache.CategoryScheduledEvents) {
							for i := range gatewayGuild.GuildScheduledEvents {
								ev := gatewayGuild.GuildScheduledEvents[i]
								d.Cache.ScheduledEvents().Set(&ev)
							}
						}
						if d.cacheStoreEnabled(cache.CategoryStageInstances) {
							for i := range gatewayGuild.StageInstances {
								instance := gatewayGuild.StageInstances[i]
								d.Cache.StageInstances().Set(&instance)
							}
						}
					}
				}
			}

			if guildCreateEvent.Guild.IsAvailable() && d.IsGuildUnavailable(guildCreateEvent.Guild.GetID()) {
				// Discord fires GUILD_CREATE with available=true to signal a guild is back online.
				// Dispatch to user handlers so they can react to the guild becoming reachable again.
				d.Logger.Debug("Guild is available again", slog.Any("guildId", guildCreateEvent.Guild.GetID()))
				d.deleteUnavailableGuild(guildCreateEvent.Guild.GetID())

				return true
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

			d.removeGuildFromCache(guildDeleteEvent.ID)

			prefix := string(guildDeleteEvent.ID) + ":"
			d.voiceStatesMu.Lock()
			for key := range d.voiceStates {
				if strings.HasPrefix(key, prefix) {
					delete(d.voiceStates, key)
				}
			}
			d.voiceStatesMu.Unlock()
		}
	case events.EventVoiceStateUpdate:
		{
			if vse, ok := event.(*events.VoiceStateUpdateEvent); ok && vse.GuildID != nil {
				key := string(*vse.GuildID) + ":" + string(vse.UserID)

				d.voiceStatesMu.Lock()
				vse.OldState = d.voiceStates[key]
				if vse.ChannelID == nil {
					delete(d.voiceStates, key)
				} else {
					s := vse.VoiceState
					d.voiceStates[key] = &s
				}
				d.voiceStatesMu.Unlock()

				if d.cacheStoreEnabled(cache.CategoryVoiceStates) {
					guildID := *vse.GuildID
					if vse.ChannelID == nil {
						d.Cache.VoiceStates().Delete(guildID, vse.UserID)
					} else {
						state := vse.VoiceState
						d.Cache.VoiceStates().Set(guildID, &state)
					}
				}
			}
		}
	case events.EventUserUpdate:
		{
			if d.cacheStoreEnabled(cache.CategoryUsers) {
				var ev events.UserUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal USER_UPDATE event", slog.Any("err", err))
					return false
				}
				u := ev.User
				d.Cache.Users().Set(&u)
			}
		}
	case events.EventPresenceUpdate:
		{
			if d.cacheStoreEnabled(cache.CategoryPresences) {
				var ev events.PresenceUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal PRESENCE_UPDATE event", slog.Any("err", err))
					return false
				}
				presence := common.Presence{
					User:         ev.User,
					GuildID:      ev.GuildID,
					Status:       ev.Status,
					Activities:   ev.Activities,
					ClientStatus: ev.ClientStatus,
				}
				d.Cache.Presences().Set(&presence)
			}
		}
	case events.EventGuildUpdate:
		{
			if d.cacheStoreEnabled(cache.CategoryGuilds) || d.cacheStoreEnabled(cache.CategoryRoles) {
				var ev events.GuildUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal GUILD_UPDATE event", slog.Any("err", err))
					return false
				}
				g := ev.Guild
				if d.cacheStoreEnabled(cache.CategoryGuilds) {
					d.Cache.Guilds().Set(&g)
				}
				if d.cacheStoreEnabled(cache.CategoryRoles) {
					for i := range g.Roles {
						role := g.Roles[i]
						d.Cache.Roles().Set(g.ID, &role)
					}
				}
			}
		}
	case events.EventGuildEmojisUpdate:
		{
			if d.cacheStoreEnabled(cache.CategoryEmojis) {
				var ev events.GuildEmojisUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal GUILD_EMOJIS_UPDATE event", slog.Any("err", err))
					return false
				}
				emojis := make([]*common.Emoji, 0, len(ev.Emojis))
				for i := range ev.Emojis {
					emoji := ev.Emojis[i]
					emojis = append(emojis, &emoji)
				}
				d.Cache.Emojis().SetAll(ev.GuildID, emojis)
			}
		}
	case events.EventGuildStickersUpdate:
		{
			if d.cacheStoreEnabled(cache.CategoryStickers) {
				var ev events.GuildStickersUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal GUILD_STICKERS_UPDATE event", slog.Any("err", err))
					return false
				}
				stickers := make([]*common.Sticker, 0, len(ev.Stickers))
				for i := range ev.Stickers {
					sticker := ev.Stickers[i]
					stickers = append(stickers, &sticker)
				}
				d.Cache.Stickers().SetAll(ev.GuildID, stickers)
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
			d.cacheMember(ev.GuildID, &ev.GuildMember)
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
			if d.cacheStoreEnabled(cache.CategoryMembers) {
				d.Cache.Members().Delete(ev.GuildID, ev.User.ID)
			}
		}
	case events.EventGuildMemberUpdate:
		{
			if d.cacheStoreEnabled(cache.CategoryMembers) {
				var ev events.GuildMemberUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal GUILD_MEMBER_UPDATE event", slog.Any("err", err))
					return false
				}
				m, ok := d.Cache.Members().Get(ev.GuildID, ev.User.ID)
				if !ok {
					m = &common.GuildMember{}
				}
				m.User = &ev.User
				m.Nick = ev.Nick
				m.AvatarHash = ev.AvatarHash
				m.PremiumSince = ev.PremiumSince
				m.CommunicationDisabledUntil = ev.CommunicationDisabledUntil
				m.Roles = make([]string, 0, len(ev.Roles))
				for _, r := range ev.Roles {
					m.Roles = append(m.Roles, r.String())
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
				d.Cache.Members().Set(ev.GuildID, m)
			}
		}
	case events.EventGuildRoleCreate:
		{
			if d.cacheStoreEnabled(cache.CategoryRoles) {
				var ev events.GuildRoleCreateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal GUILD_ROLE_CREATE event", slog.Any("err", err))
					return false
				}
				role := ev.Role
				d.Cache.Roles().Set(ev.GuildID, &role)
			}
		}
	case events.EventGuildRoleUpdate:
		{
			if d.cacheStoreEnabled(cache.CategoryRoles) {
				var ev events.GuildRoleUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal GUILD_ROLE_UPDATE event", slog.Any("err", err))
					return false
				}
				role := ev.Role
				d.Cache.Roles().Set(ev.GuildID, &role)
			}
		}
	case events.EventGuildRoleDelete:
		{
			if d.cacheStoreEnabled(cache.CategoryRoles) {
				var ev events.GuildRoleDeleteEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal GUILD_ROLE_DELETE event", slog.Any("err", err))
					return false
				}
				d.removeRoleFromCache(ev.RoleID)
			}
		}
	case events.EventGuildScheduledEventCreate:
		{
			if d.cacheStoreEnabled(cache.CategoryScheduledEvents) {
				var ev events.GuildScheduledEventCreateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal GUILD_SCHEDULED_EVENT_CREATE event", slog.Any("err", err))
					return false
				}
				scheduled := ev.GuildScheduledEvent
				d.Cache.ScheduledEvents().Set(&scheduled)
			}
		}
	case events.EventGuildScheduledEventUpdate:
		{
			if d.cacheStoreEnabled(cache.CategoryScheduledEvents) {
				var ev events.GuildScheduledEventUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal GUILD_SCHEDULED_EVENT_UPDATE event", slog.Any("err", err))
					return false
				}
				scheduled := ev.GuildScheduledEvent
				d.Cache.ScheduledEvents().Set(&scheduled)
			}
		}
	case events.EventGuildScheduledEventDelete:
		{
			if d.cacheStoreEnabled(cache.CategoryScheduledEvents) {
				var ev events.GuildScheduledEventDeleteEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal GUILD_SCHEDULED_EVENT_DELETE event", slog.Any("err", err))
					return false
				}
				d.Cache.ScheduledEvents().Delete(ev.ID)
			}
		}
	case events.EventStageInstanceCreate:
		{
			if d.cacheStoreEnabled(cache.CategoryStageInstances) {
				var ev events.StageInstanceCreateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal STAGE_INSTANCE_CREATE event", slog.Any("err", err))
					return false
				}
				instance := ev.StageInstance
				d.Cache.StageInstances().Set(&instance)
			}
		}
	case events.EventStageInstanceUpdate:
		{
			if d.cacheStoreEnabled(cache.CategoryStageInstances) {
				var ev events.StageInstanceUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal STAGE_INSTANCE_UPDATE event", slog.Any("err", err))
					return false
				}
				instance := ev.StageInstance
				d.Cache.StageInstances().Set(&instance)
			}
		}
	case events.EventStageInstanceDelete:
		{
			if d.cacheStoreEnabled(cache.CategoryStageInstances) {
				var ev events.StageInstanceDeleteEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal STAGE_INSTANCE_DELETE event", slog.Any("err", err))
					return false
				}
				d.Cache.StageInstances().Delete(ev.ID)
			}
		}
	case events.EventGuildSoundboardSoundCreate:
		{
			if d.cacheStoreEnabled(cache.CategorySoundboard) {
				var ev events.GuildSoundboardSoundCreateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal GUILD_SOUNDBOARD_SOUND_CREATE event", slog.Any("err", err))
					return false
				}
				sound := ev.SoundboardSound
				if sound.GuildID != nil {
					d.Cache.Soundboard().Set(*sound.GuildID, &sound)
				}
			}
		}
	case events.EventGuildSoundboardSoundUpdate:
		{
			if d.cacheStoreEnabled(cache.CategorySoundboard) {
				var ev events.GuildSoundboardSoundUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal GUILD_SOUNDBOARD_SOUND_UPDATE event", slog.Any("err", err))
					return false
				}
				sound := ev.SoundboardSound
				if sound.GuildID != nil {
					d.Cache.Soundboard().Set(*sound.GuildID, &sound)
				}
			}
		}
	case events.EventGuildSoundboardSoundDelete:
		{
			if d.cacheStoreEnabled(cache.CategorySoundboard) {
				var ev events.GuildSoundboardSoundDeleteEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal GUILD_SOUNDBOARD_SOUND_DELETE event", slog.Any("err", err))
					return false
				}
				d.Cache.Soundboard().Delete(ev.SoundID)
			}
		}
	case events.EventGuildSoundboardSoundsUpdate:
		{
			if d.cacheStoreEnabled(cache.CategorySoundboard) {
				var ev events.GuildSoundboardSoundsUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal GUILD_SOUNDBOARD_SOUNDS_UPDATE event", slog.Any("err", err))
					return false
				}
				sounds := make([]*common.SoundboardSound, 0, len(ev.SoundboardSounds))
				for i := range ev.SoundboardSounds {
					sound := ev.SoundboardSounds[i]
					sounds = append(sounds, &sound)
				}
				d.Cache.Soundboard().SetAll(ev.GuildID, sounds)
			}
		}
	case events.EventSoundboardSounds:
		{
			if d.cacheStoreEnabled(cache.CategorySoundboard) {
				var ev events.SoundboardSoundsEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal SOUNDBOARD_SOUNDS event", slog.Any("err", err))
					return false
				}
				sounds := make([]*common.SoundboardSound, 0, len(ev.SoundboardSounds))
				for i := range ev.SoundboardSounds {
					sound := ev.SoundboardSounds[i]
					sounds = append(sounds, &sound)
				}
				d.Cache.Soundboard().SetAll(ev.GuildID, sounds)
			}
		}
	case events.EventChannelCreate:
		{
			if d.cacheStoreEnabled(cache.CategoryChannels) {
				var ev events.ChannelCreateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal CHANNEL_CREATE event", slog.Any("err", err))
					return false
				}
				ch := ev.Channel
				d.cacheChannel(&ch)
			}
		}
	case events.EventChannelUpdate:
		{
			if d.cacheStoreEnabled(cache.CategoryChannels) {
				var ev events.ChannelUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal CHANNEL_UPDATE event", slog.Any("err", err))
					return false
				}
				ch := ev.Channel
				d.cacheChannel(&ch)
			}
		}
	case events.EventChannelDelete:
		{
			if d.cacheStoreEnabled(cache.CategoryChannels) || d.cacheStoreEnabled(cache.CategoryMessages) {
				var ev events.ChannelDeleteEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal CHANNEL_DELETE event", slog.Any("err", err))
					return false
				}
				d.removeChannelFromCache(ev.ID)
			}
		}
	case events.EventThreadCreate:
		{
			if d.cacheStoreEnabled(cache.CategoryChannels) {
				var ev events.ThreadCreateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal THREAD_CREATE event", slog.Any("err", err))
					return false
				}
				ch := ev.Channel
				d.cacheChannel(&ch)
			}
		}
	case events.EventThreadUpdate:
		{
			if d.cacheStoreEnabled(cache.CategoryChannels) {
				var ev events.ThreadUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal THREAD_UPDATE event", slog.Any("err", err))
					return false
				}
				ch := ev.Channel
				d.cacheChannel(&ch)
			}
		}
	case events.EventThreadDelete:
		{
			if d.cacheStoreEnabled(cache.CategoryChannels) || d.cacheStoreEnabled(cache.CategoryMessages) {
				var ev events.ThreadDeleteEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal THREAD_DELETE event", slog.Any("err", err))
					return false
				}
				d.removeChannelFromCache(ev.ID)
			}
		}
	case events.EventThreadListSync:
		{
			if d.cacheStoreEnabled(cache.CategoryChannels) {
				var ev events.ThreadListSyncEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal THREAD_LIST_SYNC event", slog.Any("err", err))
					return false
				}
				for i := range ev.Threads {
					ch := ev.Threads[i]
					d.cacheChannel(&ch)
				}
			}
		}
	case events.EventMessageCreate:
		{
			if d.cacheStoreEnabled(cache.CategoryMessages) || d.cacheStoreEnabled(cache.CategoryUsers) {
				var ev events.MessageCreateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal MESSAGE_CREATE event", slog.Any("err", err))
					return false
				}
				msg := ev.Message
				if d.cacheStoreEnabled(cache.CategoryMessages) {
					d.Cache.Messages().Add(&msg)
				}
				if ev.Author != nil && d.cacheStoreEnabled(cache.CategoryUsers) {
					d.Cache.Users().Set(ev.Author)
				}
			}
		}
	case events.EventMessageUpdate:
		{
			if d.cacheStoreEnabled(cache.CategoryMessages) {
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
			if d.cacheStoreEnabled(cache.CategoryMessages) {
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
			if d.cacheStoreEnabled(cache.CategoryMessages) {
				var ev events.MessageDeleteBulkEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal MESSAGE_DELETE_BULK event", slog.Any("err", err))
					return false
				}
				d.Cache.Messages().DeleteBulk(ev.ChannelID, ev.IDs)
			}
		}
	case events.EventGuildMembersChunk:
		{
			if d.cacheStoreEnabled(cache.CategoryMembers) {
				var ev events.GuildMembersChunkEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal GUILD_MEMBERS_CHUNK event", slog.Any("err", err))
					return false
				}
				for i := range ev.Members {
					d.Cache.Members().Set(ev.GuildID, &ev.Members[i])
					if d.cacheStoreEnabled(cache.CategoryUsers) && ev.Members[i].User != nil {
						u := *ev.Members[i].User
						d.Cache.Users().Set(&u)
					}
				}
				d.Logger.Debug("Cached guild members chunk",
					slog.Any("guildId", ev.GuildID),
					slog.Int("chunk", ev.ChunkIndex+1),
					slog.Int("total", ev.ChunkCount),
					slog.Int("members", len(ev.Members)),
				)
			}
		}
	case events.EventVoiceChannelStatusUpdate:
		{
			if d.cacheStoreEnabled(cache.CategoryChannels) {
				var ev events.VoiceChannelStatusUpdateEvent
				if err := json.Unmarshal(msg, &ev); err != nil {
					d.Logger.Error("Failed to unmarshal VOICE_CHANNEL_STATUS_UPDATE event", slog.Any("err", err))
					return false
				}
				if ch, ok := d.Cache.Channels().Get(ev.ChannelID); ok {
					ch.Status = ev.Status
					d.Cache.Channels().Set(ch)
				}
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
// It waits for all in-flight event handlers to finish before returning.
func (d *Client) Shutdown() error {
	d.shutdown.Store(true)
	var errs []error
	d.wsMu.RLock()
	ws := d.Websocket
	d.wsMu.RUnlock()
	if ws != nil {
		if err := ws.close(); err != nil {
			errs = append(errs, err)
		}
	}

	// Worker-pool mode: close the event channel so workers drain it and exit,
	// then wait for them. This must happen before dispatchWg.Wait() because
	// workers call dispatch() which adds entries to dispatchWg.
	if d.eventCh != nil {
		close(d.eventCh)
		d.workerWg.Wait()
	}

	// Wait for all user event handler goroutines (spawned by dispatch()).
	d.dispatchWg.Wait()

	if d.Cache != nil {
		if err := d.Cache.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if d.RestClient != nil {
		d.RestClient.Close()
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
