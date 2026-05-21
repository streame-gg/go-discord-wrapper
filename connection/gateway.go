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
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
	"github.com/streame-gg/go-discord-wrapper/util"

	"github.com/gorilla/websocket"
)

type EventHandler func(*Client, events.Event)

type Client struct {
	token *string

	APIVersion *discord.APIVersion

	Logger *slog.Logger

	Intents *discord.Intent

	Websocket *Websocket

	// httpClient is the shared HTTP client used for REST calls inside the
	// gateway package (e.g. GET /gateway/bot). Reusing it avoids allocating a
	// new transport and connection pool on every call (Bug 40).
	httpClient *http.Client

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

	// shutdownCh is closed by Shutdown() to signal concurrent goroutines (e.g. eventCh
	// senders) that shutdown is in progress.  Initialized in NewClient.
	shutdownCh   chan struct{}
	shutdownOnce sync.Once

	// readyCh is closed exactly once when the first READY event is received.
	// Unlike Websocket.Ready, it lives on the Client and survives reconnects,
	// so Login() can safely wait on it without racing against d.Websocket being
	// replaced by a concurrent reconnect.
	readyCh   chan struct{}
	readyOnce sync.Once

	// eventWg tracks unlimited-mode goroutines (one per event) so Shutdown can
	// drain them before waiting on handler goroutines tracked by dispatchWg.
	eventWg sync.WaitGroup

	// maxReconnectRetries mirrors options.Config.MaxReconnectRetries.
	maxReconnectRetries int

	// cacheAutoPopulate mirrors options.Config.CacheStores, zero disables auto-population.
	cacheAutoPopulate cache.OverflowCategory

	UnavailableGuilds map[discord.Snowflake]struct{}

	User *discord.User

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
	channelsByGuild map[discord.Snowflake]map[discord.Snowflake]struct{}
	// guildByChannel is the reverse mapping: channel ID → guild ID.
	guildByChannel map[discord.Snowflake]discord.Snowflake

	// threadIndexMu protects threadsByParent.
	threadIndexMu sync.RWMutex
	// threadsByParent maps parent channel ID → set of thread IDs.
	// Used to evict stale threads on THREAD_LIST_SYNC (Bug 43).
	threadsByParent map[discord.Snowflake]map[discord.Snowflake]struct{}

	// guildMemberCounts tracks the member_count for every available guild on
	// this shard, updated from GUILD_CREATE / GUILD_DELETE gateway events.
	guildMemberCounts map[discord.Snowflake]int
	guildMu           sync.RWMutex

	// voiceStates caches the latest voice state per (guildID:userID) key so
	// VoiceStateUpdateEvent can carry OldState alongside the new state.
	voiceStates   map[string]*discord.VoiceState
	voiceStatesMu sync.RWMutex

	// Client lifecycle event handlers.
	onConnect      []func(*Client)
	onDisconnect   []func(*Client, error)
	onReconnect    []func(*Client)
	onPacketError  []func(*Client, error)
	clientEventsMu sync.RWMutex

	// shardHandlers are persistent handlers registered via OnShardMessage.
	shardHandlers    []func(options.ShardMessage)
	shardHandlersMu  sync.RWMutex
	shardDispatchSem chan struct{} // limits concurrent shard-message handler goroutines (Bug 35)

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
func NewClient(token string, intents discord.Intent, opts ...options.Option) (*Client, error) {
	cfg := options.Build(options.Config{
		APIVersion:  discord.APIVersion10,
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
		UnavailableGuilds:   make(map[discord.Snowflake]struct{}),
		guildMemberCounts:   make(map[discord.Snowflake]int),
		voiceStates:         make(map[string]*discord.VoiceState),
		channelsByGuild:     make(map[discord.Snowflake]map[discord.Snowflake]struct{}),
		guildByChannel:      make(map[discord.Snowflake]discord.Snowflake),
		threadsByParent:     make(map[discord.Snowflake]map[discord.Snowflake]struct{}),
		httpClient:          &http.Client{Timeout: 10 * time.Second},
		RestClient:          rc,
		Sharding:            cfg.Sharding,
		Coordinator:         cfg.Coordinator,
		Cache:               cfg.Cache,
		maxReconnectRetries: cfg.MaxReconnectRetries,
		cacheAutoPopulate:   cacheStores,
		shardDispatchSem:    make(chan struct{}, 16),
		shutdownCh:          make(chan struct{}),
		readyCh:             make(chan struct{}),
	}

	if cfg.MaxConcurrentEvents == 0 {
		cfg.MaxConcurrentEvents = 64
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

// Ready returns a channel that is closed once the first READY event is
// received from the Discord gateway.  Unlike Websocket.Ready, this channel
// lives on the Client and is never replaced, so it is safe to use across
// reconnects.
func (d *Client) Ready() <-chan struct{} {
	return d.readyCh
}

// dispatchJob carries a single gateway event through the worker pool.
// Cache updates (internalEventHandler) run synchronously in listenWebsocket
// before the job is enqueued, so rawPayload is no longer needed here.
type dispatchJob struct {
	event events.Event
}

// processEvent dispatches one gateway event to user handlers.
// Cache updates have already been applied synchronously in listenWebsocket.
func (d *Client) processEvent(job dispatchJob) {
	defer func() {
		if r := recover(); r != nil {
			d.Logger.Error("panic in event dispatch", slog.Any("recover", r))
		}
	}()
	d.dispatch(job.event)
}

// runEventWorker processes events from d.eventCh until shutdownCh is closed.
// It drains any buffered events before returning so in-flight work completes.
// eventCh is never closed; workers exit via shutdownCh to avoid close-vs-send
// data races (Bug 48).
func (d *Client) runEventWorker() {
	defer d.workerWg.Done()
	for {
		select {
		case job := <-d.eventCh:
			d.processEvent(job)
		case <-d.shutdownCh:
			for {
				select {
				case job := <-d.eventCh:
					d.processEvent(job)
				default:
					return
				}
			}
		}
	}
}

// enqueueOrDispatch submits event for handler dispatch.
// In pool mode it enqueues to eventCh; in unlimited mode it spawns a goroutine.
// Returns true if shutdown is in progress (caller should return nil from listenWebsocket).
func (d *Client) enqueueOrDispatch(event events.Event) bool {
	job := dispatchJob{event: event}
	if d.eventCh != nil {
		select {
		case d.eventCh <- job:
		case <-d.shutdownCh:
			return true
		default:
			d.Logger.Warn("Event queue full, dropping event",
				slog.String("type", string(event.Event())),
				slog.Int("queueCap", cap(d.eventCh)),
			)
		}
		return false
	}
	select {
	case <-d.shutdownCh:
		return true
	default:
	}
	d.eventWg.Add(1)
	go func() {
		defer d.eventWg.Done()
		select {
		case <-d.shutdownCh:
			return
		default:
		}
		d.processEvent(job)
	}()
	return false
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
		select {
		case d.shardDispatchSem <- struct{}{}:
			// Acquired a slot: run handler in its own goroutine.
			go func() {
				defer func() { <-d.shardDispatchSem }()
				h(msg)
			}()
		default:
			// All 16 slots taken: execute synchronously to apply back-pressure
			// rather than spawning an unbounded number of goroutines (Bug 35).
			h(msg)
		}
	}
}

// ── Gateway connection ───────────────────────────────────────────────────────

func (d *Client) initializeGatewayConnection(ctx context.Context) (*discord.BotRegisterResponse, error) {
	const userAgent = "DiscordBot (https://github.com/streame-gg/go-discord-wrapper, alpha)"

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://discord.com"+discord.APIBaseString(*d.APIVersion)+"gateway/bot", nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bot "+*d.token)
		req.Header.Set("User-Agent", userAgent)

		resp, err := d.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			_ = resp.Body.Close()
			retryAfter := 1 * time.Second
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if secs, parseErr := time.ParseDuration(ra + "s"); parseErr == nil {
					retryAfter = secs
				}
			}
			d.Logger.Warn("Gateway /gateway/bot rate-limited", slog.Duration("retryAfter", retryAfter))
			select {
			case <-time.After(retryAfter):
				continue
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			return nil, errors.New("failed to register bot gateway connection, status code: " + resp.Status)
		}

		var result discord.BotRegisterResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}
		return &result, nil
	}
}

func (d *Client) Login(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	gatewayResp, err := d.initializeGatewayConnection(ctx)
	if err != nil {
		return err
	}

	d.Logger.Debug("Connecting to gateway websocket", slog.String("url", gatewayResp.URL), slog.Int("shards", gatewayResp.Shards))

	if err := d.connectWebsocket(gatewayResp.URL, false, nil, nil); err != nil {
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

				// Non-recoverable close codes — reconnecting will not help.
				if websocket.IsCloseError(err, 4004, 4010, 4011, 4012, 4013, 4014) {
					d.Logger.Error("Gateway connection closed with a non-recoverable code; cannot reconnect", slog.Any("err", err))
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

	select {
	case <-d.readyCh:
	case <-ctx.Done():
		_ = d.Shutdown()
		return ctx.Err()
	case <-d.shutdownCh:
		return errors.New("gateway closed before READY")
	}

	d.Logger.Info("Successfully connected to the Discord gateway")
	d.emitConnect()

	return nil
}

func (d *Client) dispatch(event events.Event) {
	handlers := d.discordEventEmitter.Handlers(event.Event())

	if d.eventCh != nil {
		// Pool mode: handlers run serially, inline in the worker goroutine.
		// MaxConcurrentEvents is the number of workers, so it is the true upper
		// bound on concurrent handler execution.
		for _, h := range handlers {
			h := h
			func() {
				defer func() {
					if r := recover(); r != nil {
						d.Logger.Error("panic in event handler", slog.Any("recover", r))
					}
				}()
				h(d, event)
			}()
		}
		return
	}

	// Unlimited mode: each handler gets its own goroutine.
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
				return false
			}

			d.wsMu.Lock()
			d.Websocket.SessionID = &readyEvent.SessionID
			d.Websocket.ReconnectURL = &readyEvent.ResumeGatewayURL
			d.wsMu.Unlock()
			d.User = &readyEvent.User

			if len(readyEvent.Shard) >= 2 {
				d.Logger.Debug("Connected to shard", slog.Int("shard", readyEvent.Shard[0]+1), slog.Int("total", readyEvent.Shard[1]))
			}

			// Build the set of guilds the bot currently belongs to.
			currentGuildIDs := make(map[discord.Snowflake]struct{}, len(readyEvent.Guilds))
			for _, g := range readyEvent.Guilds {
				currentGuildIDs[g.Guild.GetID()] = struct{}{}
			}

			// Prune stale member-count entries even when cache is disabled.
			d.guildMu.Lock()
			for guildID := range d.guildMemberCounts {
				if _, present := currentGuildIDs[guildID]; !present {
					delete(d.guildMemberCounts, guildID)
				}
			}
			d.guildMu.Unlock()

			// Remove cached guilds that are no longer present (bot was kicked while offline).
			if d.cacheEnabled() {
				for _, cachedGuild := range d.Cache.Guilds().All().Values() {
					if _, present := currentGuildIDs[cachedGuild.ID]; !present {
						d.removeGuildFromCache(cachedGuild.ID)
					}
				}
			}

			// Reset unavailable-guild tracking and rebuild from the READY payload.
			d.mu.Lock()
			d.UnavailableGuilds = make(map[discord.Snowflake]struct{})
			d.mu.Unlock()
			for _, guild := range readyEvent.Guilds {
				if !guild.Guild.IsAvailable() {
					d.addUnavailableGuild(guild.Guild.GetID())
				}
			}

			d.Websocket.readyOnce.Do(func() { close(d.Websocket.Ready) })
			d.readyOnce.Do(func() { close(d.readyCh) })

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
				var guild discord.Guild
				var gatewayGuild discord.GatewayGuild
				hasGateway := false
				switch g := guildCreateEvent.Guild.(type) {
				case discord.GatewayGuild:
					guild = g.Guild
					gatewayGuild = g
					hasGateway = true
					if memberCount == 0 && g.MemberCount > 0 {
						memberCount = g.MemberCount
					}
				case discord.Guild:
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
						key := guildID.String() + ":" + vs.UserID.String()
						vscopy := vs
						if vscopy.GuildID == nil {
							vscopy.GuildID = &guildID
						}
						d.voiceStates[key] = &vscopy
					}
					d.voiceStatesMu.Unlock()
				}

				if d.cacheEnabled() {
					if !guild.ID.IsEmpty() {
						guild.Hydrate(d)
						d.setGuildManagers(&guild)
					}

					if !guild.ID.IsEmpty() && d.cacheStoreEnabled(cache.CategoryGuilds) {
						d.Cache.Guilds().Set(&guild)
					}

					if !guild.ID.IsEmpty() && d.cacheStoreEnabled(cache.CategoryRoles) {
						// Delete stale roles before re-adding so removed roles don't persist.
						d.Cache.Roles().DeleteGuild(guild.ID)
						for i := range guild.RawRoles {
							role := guild.RawRoles[i]
							role.GuildID = guild.ID
							role.Hydrate(d)
							d.Cache.Roles().Set(guild.ID, &role)
						}
					}

					if !guild.ID.IsEmpty() && d.cacheStoreEnabled(cache.CategoryEmojis) {
						emojis := make([]*discord.Emoji, 0, len(guild.RawEmojis))
						for i := range guild.RawEmojis {
							emoji := guild.RawEmojis[i]
							if !emoji.ID.IsEmpty() {
								emoji.GuildID = guild.ID
								emoji.Hydrate(d)
								emojis = append(emojis, &emoji)
							}
						}
						d.Cache.Emojis().SetAll(guild.ID, emojis)
					}
					if !guild.ID.IsEmpty() && d.cacheStoreEnabled(cache.CategoryStickers) && len(guild.RawStickers) > 0 {
						stickers := make([]*discord.Sticker, 0, len(guild.RawStickers))
						for i := range guild.RawStickers {
							sticker := guild.RawStickers[i]
							if !sticker.ID.IsEmpty() {
								if sticker.GuildID == nil {
									sticker.GuildID = &guild.ID
								}
								sticker.Hydrate(d)
								stickers = append(stickers, &sticker)
							}
						}
						d.Cache.Stickers().SetAll(guild.ID, stickers)
					}

					if hasGateway {
						guildID := guildCreateEvent.Guild.GetID()
						if d.cacheStoreEnabled(cache.CategoryChannels) {
							// Drain stale channels before re-adding so deleted channels don't persist.
							for _, oldID := range d.drainGuildChannelIDs(guildID) {
								d.Cache.Channels().Delete(oldID)
								d.Cache.Messages().DeleteChannel(oldID)
							}
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
								d.trackThread(&ch)
							}
						}

						if d.cacheStoreEnabled(cache.CategoryMembers) {
							// Delete stale members before re-adding.
							d.Cache.Members().DeleteGuild(guildID)
							for i := range gatewayGuild.Members {
								m := gatewayGuild.Members[i]
								m.GuildID = guildID
								if m.User != nil {
									m.UserID = m.User.ID
								}
								m.Hydrate(d)
								d.Cache.Members().Set(guildID, &m)
							}
						}

						if d.cacheStoreEnabled(cache.CategoryUsers) {
							for i := range gatewayGuild.Members {
								if gatewayGuild.Members[i].User == nil {
									continue
								}
								u := *gatewayGuild.Members[i].User
								u.Hydrate(d)
								d.Cache.Users().Set(&u)
							}
						}

						if d.cacheStoreEnabled(cache.CategoryVoiceStates) {
							// Delete stale voice states before re-adding.
							d.Cache.VoiceStates().DeleteGuild(guildID)
							for i := range gatewayGuild.VoiceStates {
								vs := gatewayGuild.VoiceStates[i]
								if vs.GuildID == nil {
									vs.GuildID = &guildID
								}
								d.Cache.VoiceStates().Set(guildID, &vs)
							}
						}

						if d.cacheStoreEnabled(cache.CategoryPresences) {
							// Delete stale presences before re-adding.
							d.Cache.Presences().DeleteGuild(guildID)
							for i := range gatewayGuild.Presences {
								p := gatewayGuild.Presences[i]
								presence := discord.Presence{
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
							sounds := make([]*discord.SoundboardSound, 0, len(gatewayGuild.SoundboardSounds))
							for i := range gatewayGuild.SoundboardSounds {
								s := gatewayGuild.SoundboardSounds[i]
								s.Hydrate(d)
								sounds = append(sounds, &s)
							}
							d.Cache.Soundboard().SetAll(guildID, sounds)
						}

						if d.cacheStoreEnabled(cache.CategoryScheduledEvents) {
							// Delete stale scheduled events before re-adding.
							d.Cache.ScheduledEvents().DeleteGuild(guildID)
							for i := range gatewayGuild.GuildScheduledEvents {
								ev := gatewayGuild.GuildScheduledEvents[i]
								ev.Hydrate(d)
								d.Cache.ScheduledEvents().Set(&ev)
							}
						}
						if d.cacheStoreEnabled(cache.CategoryStageInstances) {
							// Delete stale stage instances before re-adding.
							d.Cache.StageInstances().DeleteGuild(guildID)
							for i := range gatewayGuild.StageInstances {
								instance := gatewayGuild.StageInstances[i]
								instance.Hydrate(d)
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
			d.removeGuildFromCache(guildDeleteEvent.ID)

			prefix := guildDeleteEvent.ID.String() + ":"
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
				key := vse.GuildID.String() + ":" + vse.UserID.String()

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
				u := ev.NewUser
				u.Hydrate(d)
				d.Cache.Users().Set(&u)
			}
		}
	case events.EventPresenceUpdate:
		{
			ev, ok := event.(*events.PresenceUpdateEvent)
			if !ok {
				return false
			}
			if d.cacheStoreEnabled(cache.CategoryPresences) {
				if old, exists := d.Cache.Presences().Get(ev.NewPresence.GuildID, ev.NewPresence.User.ID); exists {
					ev.OldPresence = old
				}
				presence := ev.NewPresence
				d.Cache.Presences().Set(&presence)
			}
		}
	case events.EventGuildUpdate:
		{
			ev, ok := event.(*events.GuildUpdateEvent)
			if !ok {
				return false
			}
			g := ev.NewGuild
			g.Hydrate(d)
			d.setGuildManagers(&g)
			ev.NewGuild = g
			if d.cacheStoreEnabled(cache.CategoryGuilds) {
				if old, exists := d.Cache.Guilds().Get(g.ID); exists {
					ev.OldGuild = old
				}
				d.Cache.Guilds().Set(&g)
			}
			if d.cacheStoreEnabled(cache.CategoryRoles) {
				// Delete stale roles before re-adding so removed roles don't persist.
				d.Cache.Roles().DeleteGuild(g.ID)
				for i := range g.RawRoles {
					role := g.RawRoles[i]
					role.GuildID = g.ID
					role.Hydrate(d)
					d.Cache.Roles().Set(g.ID, &role)
				}
			}
		}
	case events.EventGuildEmojisUpdate:
		{
			ev, ok := event.(*events.GuildEmojisUpdateEvent)
			if !ok {
				return false
			}
			if d.cacheStoreEnabled(cache.CategoryEmojis) {
				if col := d.Cache.Emojis().GetByGuild(ev.GuildID); col != nil && col.Len() > 0 {
					ev.OldEmojis = col.Values()
				}
				for _, emoji := range ev.NewEmojis {
					emoji.GuildID = ev.GuildID
					emoji.Hydrate(d)
				}
				d.Cache.Emojis().SetAll(ev.GuildID, ev.NewEmojis)
			}
		}
	case events.EventGuildStickersUpdate:
		{
			ev, ok := event.(*events.GuildStickersUpdateEvent)
			if !ok {
				return false
			}
			if d.cacheStoreEnabled(cache.CategoryStickers) {
				if col := d.Cache.Stickers().GetByGuild(ev.GuildID); col != nil && col.Len() > 0 {
					ev.OldStickers = col.Values()
				}
				for _, sticker := range ev.NewStickers {
					if sticker.GuildID == nil {
						sticker.GuildID = &ev.GuildID
					}
					sticker.Hydrate(d)
				}
				d.Cache.Stickers().SetAll(ev.GuildID, ev.NewStickers)
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
			ev, ok := event.(*events.GuildMemberUpdateEvent)
			if !ok {
				return false
			}
			if d.cacheStoreEnabled(cache.CategoryMembers) {
				if existing, exists := d.Cache.Members().Get(ev.GuildID, ev.NewMember.UserID); exists {
					ev.OldMember = existing
				}
				m := ev.NewMember
				m.GuildID = ev.GuildID
				if m.User != nil {
					m.UserID = m.User.ID
				}
				m.Hydrate(d)
				d.Cache.Members().Set(ev.GuildID, &m)
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
				role.GuildID = ev.GuildID
				role.Hydrate(d)
				d.Cache.Roles().Set(ev.GuildID, &role)
			}
		}
	case events.EventGuildRoleUpdate:
		{
			ev, ok := event.(*events.GuildRoleUpdateEvent)
			if !ok {
				return false
			}
			if d.cacheStoreEnabled(cache.CategoryRoles) {
				if old, exists := d.Cache.Roles().Get(ev.NewRole.ID); exists {
					ev.OldRole = old
				}
				role := ev.NewRole
				role.GuildID = ev.GuildID
				role.Hydrate(d)
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

				// Purge the deleted role from every cached member's Roles slice.
				// Discord may send GUILD_MEMBER_UPDATEs for affected members, but
				// delivery is not guaranteed — removing here keeps the cache
				// consistent immediately so permission checks cannot see stale roles.
				if d.cacheStoreEnabled(cache.CategoryMembers) {
					roleStr := ev.RoleID.String()
					for _, m := range d.Cache.Members().AllInGuild(ev.GuildID).Values() {
						found := false
						for _, r := range m.Roles {
							if r == roleStr {
								found = true
								break
							}
						}
						if !found {
							continue
						}
						// Copy and filter — never mutate the cached pointer directly.
						updated := *m
						filtered := make([]string, 0, len(m.Roles)-1)
						for _, r := range m.Roles {
							if r != roleStr {
								filtered = append(filtered, r)
							}
						}
						updated.Roles = filtered
						d.Cache.Members().Set(ev.GuildID, &updated)
					}
				}
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
				scheduled.Hydrate(d)
				d.Cache.ScheduledEvents().Set(&scheduled)
			}
		}
	case events.EventGuildScheduledEventUpdate:
		{
			ev, ok := event.(*events.GuildScheduledEventUpdateEvent)
			if !ok {
				return false
			}
			if d.cacheStoreEnabled(cache.CategoryScheduledEvents) {
				if old, exists := d.Cache.ScheduledEvents().Get(ev.NewEvent.ID); exists {
					ev.OldEvent = old
				}
				scheduled := ev.NewEvent
				scheduled.Hydrate(d)
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
				instance.Hydrate(d)
				d.Cache.StageInstances().Set(&instance)
			}
		}
	case events.EventStageInstanceUpdate:
		{
			ev, ok := event.(*events.StageInstanceUpdateEvent)
			if !ok {
				return false
			}
			if d.cacheStoreEnabled(cache.CategoryStageInstances) {
				if old, exists := d.Cache.StageInstances().Get(ev.NewInstance.ID); exists {
					ev.OldInstance = old
				}
				instance := ev.NewInstance
				instance.Hydrate(d)
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
				sound.Hydrate(d)
				if sound.GuildID != nil {
					d.Cache.Soundboard().Set(*sound.GuildID, &sound)
				}
			}
		}
	case events.EventGuildSoundboardSoundUpdate:
		{
			ev, ok := event.(*events.GuildSoundboardSoundUpdateEvent)
			if !ok {
				return false
			}
			if d.cacheStoreEnabled(cache.CategorySoundboard) {
				if old, exists := d.Cache.Soundboard().Get(ev.NewSound.SoundID); exists {
					ev.OldSound = old
				}
				sound := ev.NewSound
				sound.Hydrate(d)
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
				for _, sound := range ev.NewSoundboardSounds {
					sound.Hydrate(d)
				}
				d.Cache.Soundboard().SetAll(ev.GuildID, ev.NewSoundboardSounds)
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
				sounds := make([]*discord.SoundboardSound, 0, len(ev.SoundboardSounds))
				for i := range ev.SoundboardSounds {
					sound := ev.SoundboardSounds[i]
					sound.Hydrate(d)
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
			ev, ok := event.(*events.ChannelUpdateEvent)
			if !ok {
				return false
			}
			if d.cacheStoreEnabled(cache.CategoryChannels) {
				if old, exists := d.Cache.Channels().Get(ev.NewChannel.ID); exists {
					ev.OldChannel = old
				}
				ch := ev.NewChannel
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
				d.trackThread(&ch)
			}
		}
	case events.EventThreadUpdate:
		{
			if d.cacheStoreEnabled(cache.CategoryChannels) {
				var ch discord.Channel
				if ev, ok := event.(*events.ThreadUpdateEvent); ok {
					if old, exists := d.Cache.Channels().Get(ev.NewThread.ID); exists {
						ev.OldThread = old
					}
					ch = ev.NewThread
				} else {
					var fallback events.ThreadUpdateEvent
					if err := json.Unmarshal(msg, &fallback); err != nil {
						return false
					}
					ch = fallback.NewThread
				}
				d.cacheChannel(&ch)
				d.trackThread(&ch)
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
				if ev.ParentID != nil {
					d.untrackThread(ev.ID, *ev.ParentID)
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

				// Build a set of incoming thread IDs for O(1) lookup.
				incoming := make(map[discord.Snowflake]struct{}, len(ev.Threads))
				for i := range ev.Threads {
					incoming[ev.Threads[i].ID] = struct{}{}
				}

				// Determine which parent channels are in scope for this sync.
				// If ChannelIDs is absent, all threads in the guild are synced.
				var parents []discord.Snowflake
				if len(ev.ChannelIDs) > 0 {
					parents = ev.ChannelIDs
				} else {
					d.channelIndexMu.RLock()
					for parentID := range d.channelsByGuild[ev.GuildID] {
						parents = append(parents, parentID)
					}
					d.channelIndexMu.RUnlock()
				}

				// Evict threads that were cached for a synced parent but are
				// absent from Discord's authoritative list.
				for _, parentID := range parents {
					for _, threadID := range d.drainParentThreadIDs(parentID) {
						if _, ok := incoming[threadID]; !ok {
							d.removeChannelFromCache(threadID)
						}
					}
				}

				// Cache and index the authoritative thread list.
				for i := range ev.Threads {
					ch := ev.Threads[i]
					d.cacheChannel(&ch)
					d.trackThread(&ch)
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
					msg.Hydrate(d)
					d.Cache.Messages().Add(&msg)
				}
				if ev.Author != nil && d.cacheStoreEnabled(cache.CategoryUsers) {
					d.cacheUser(ev.Author)
				}
			}
		}
	case events.EventMessageUpdate:
		{
			ev, ok := event.(*events.MessageUpdateEvent)
			if !ok {
				return false
			}
			if d.cacheStoreEnabled(cache.CategoryMessages) {
				if old, exists := d.Cache.Messages().Get(ev.NewMessage.ChannelID, ev.NewMessage.ID); exists {
					ev.OldMessage = old
				}
				msg := ev.NewMessage
				msg.Hydrate(d)
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
					m := ev.Members[i]
					m.GuildID = ev.GuildID
					if m.User != nil {
						m.UserID = m.User.ID
					}
					m.Hydrate(d)
					d.Cache.Members().Set(ev.GuildID, &m)
					if d.cacheStoreEnabled(cache.CategoryUsers) && m.User != nil {
						d.cacheUser(m.User)
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
					// Copy before mutating — never modify the cached pointer in-place.
					updated := *ch
					updated.Status = ev.Status
					d.Cache.Channels().Set(&updated)
				}
			}
		}
	case events.EventInteractionCreate:
		if ev, ok := event.(*events.InteractionCreateEvent); ok {
			ev.Hydrate(d, context.Background())
			// Resolve Guild from cache so handlers get the fully hydrated object
			// with sub-managers. The payload only carries a partial Guild stub.
			if ev.GuildID != nil {
				var resolvedGuild *discord.Guild
				if d.cacheEnabled() {
					resolvedGuild, _ = d.Cache.Guilds().Get(*ev.GuildID)
				}
				if resolvedGuild == nil {
					// Cache miss or cache disabled. Use the partial stub if Discord
					// sent one; otherwise synthesize a minimal stub with just the ID.
					// Either way, inject managers so accessor methods never panic.
					if ev.Guild != nil {
						resolvedGuild = ev.Guild
					} else {
						resolvedGuild = &discord.Guild{ID: *ev.GuildID}
					}
					resolvedGuild.Hydrate(d)
					d.setGuildManagers(resolvedGuild)
				}
				ev.Guild = resolvedGuild
			}
			// Hydrate the invoking member and the channel if present.
			if ev.Member != nil {
				if ev.GuildID != nil {
					ev.Member.GuildID = *ev.GuildID
					if ev.Member.User != nil {
						ev.Member.UserID = ev.Member.User.ID
					}
				}
				ev.Member.Hydrate(d)
			}
			if ev.User != nil {
				ev.User.Hydrate(d)
			}
			if ev.Channel != nil {
				ev.Channel.Hydrate(d)
				d.setChannelManagers(ev.Channel)
			}
		}
		return true

	default:
		return true
	}

	return true
}

// ensure interactions package is used (import guard)
var _ interactions.Client = (*Client)(nil)

func (d *Client) addUnavailableGuild(id discord.Snowflake) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.UnavailableGuilds[id] = struct{}{}
}

func (d *Client) deleteUnavailableGuild(id discord.Snowflake) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.UnavailableGuilds, id)
}

func (d *Client) IsGuildUnavailable(id discord.Snowflake) bool {
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
// Shutdown is idempotent: concurrent or repeated calls are safe.
func (d *Client) Shutdown() error {
	var errs []error
	d.shutdownOnce.Do(func() {
		d.shutdown.Store(true)
		close(d.shutdownCh) // unblocks any eventCh senders still in listenWebsocket

		d.wsMu.RLock()
		ws := d.Websocket
		d.wsMu.RUnlock()
		if ws != nil {
			if err := ws.close(); err != nil {
				errs = append(errs, err)
			}
		}

		// Worker-pool mode: shutdownCh (already closed above) signals workers to
		// drain their queue and exit.  Wait for them before dispatchWg.Wait()
		// because workers call dispatch() which adds to dispatchWg.
		if d.eventCh != nil {
			d.workerWg.Wait()
		}

		// Unlimited mode: wait for per-event goroutines to finish before waiting
		// on handler goroutines they may have spawned via dispatch().
		d.eventWg.Wait()

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
	})
	return errors.Join(errs...)
}

// ── Event entity hydration ────────────────────────────────────────────────────

// hydrateEvent injects the client reference into entities embedded in event
// payloads so that users can call convenience methods (e.g. ev.Edit) on them
// directly without an explicit client argument.
func (d *Client) hydrateEvent(event events.Event) {
	switch ev := event.(type) {
	case *events.MessageCreateEvent:
		ev.Hydrate(d)
		if ev.Member != nil {
			if ev.GuildID != nil {
				ev.Member.GuildID = *ev.GuildID
			}
			ev.Member.Hydrate(d)
		}
	case *events.MessageUpdateEvent:
		ev.NewMessage.Hydrate(d)
	case *events.ChannelCreateEvent:
		ev.Hydrate(d)
		d.setChannelManagers(&ev.Channel)
	case *events.ChannelUpdateEvent:
		ev.NewChannel.Hydrate(d)
		d.setChannelManagers(&ev.NewChannel)
	case *events.ChannelDeleteEvent:
		ev.Hydrate(d)
		d.setChannelManagers(&ev.Channel)
	case *events.ThreadCreateEvent:
		ev.Hydrate(d)
		d.setChannelManagers(&ev.Channel)
	case *events.ThreadUpdateEvent:
		ev.NewThread.Hydrate(d)
		d.setChannelManagers(&ev.NewThread)
	case *events.GuildCreateEvent:
		switch g := ev.Guild.(type) {
		case discord.GatewayGuild:
			g.Hydrate(d)
			d.setGuildManagers(&g.Guild)
			ev.Guild = g
		case discord.Guild:
			g.Hydrate(d)
			d.setGuildManagers(&g)
			ev.Guild = g
		}
	case *events.GuildUpdateEvent:
		ev.NewGuild.Hydrate(d)
		d.setGuildManagers(&ev.NewGuild)
	case *events.GuildMemberAddEvent:
		ev.GuildMember.GuildID = ev.GuildID
		if ev.User != nil {
			ev.UserID = ev.User.ID
		}
		ev.Hydrate(d)
	case *events.GuildRoleCreateEvent:
		ev.Role.GuildID = ev.GuildID
		ev.Role.Hydrate(d)
	case *events.GuildRoleUpdateEvent:
		ev.NewRole.GuildID = ev.GuildID
		ev.NewRole.Hydrate(d)
	case *events.GuildScheduledEventCreateEvent:
		ev.Hydrate(d)
	case *events.GuildScheduledEventUpdateEvent:
		ev.NewEvent.Hydrate(d)
	case *events.GuildScheduledEventDeleteEvent:
		ev.Hydrate(d)
	case *events.AutoModerationRuleCreateEvent:
		ev.Hydrate(d)
	case *events.AutoModerationRuleUpdateEvent:
		ev.NewRule.Hydrate(d)
	case *events.AutoModerationRuleDeleteEvent:
		ev.Hydrate(d)
	case *events.StageInstanceCreateEvent:
		ev.Hydrate(d)
	case *events.StageInstanceUpdateEvent:
		ev.NewInstance.Hydrate(d)
	case *events.StageInstanceDeleteEvent:
		ev.Hydrate(d)
	case *events.UserUpdateEvent:
		ev.NewUser.Hydrate(d)
	case *events.GuildSoundboardSoundCreateEvent:
		ev.Hydrate(d)
	case *events.GuildSoundboardSoundUpdateEvent:
		ev.NewSound.Hydrate(d)
	case *events.GuildMemberUpdateEvent:
		if ev.NewMember.User != nil {
			ev.NewMember.User.Hydrate(d)
		}
		ev.NewMember.Hydrate(d)
	case *events.GuildMemberRemoveEvent:
		ev.User.Hydrate(d)
	case *events.ThreadListSyncEvent:
		for i := range ev.Threads {
			ev.Threads[i].Hydrate(d)
			d.setChannelManagers(&ev.Threads[i])
		}
	}
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
