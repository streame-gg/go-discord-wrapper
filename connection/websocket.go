package connection

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/common"
	"github.com/streame-gg/go-discord-wrapper/types/events"

	"github.com/gorilla/websocket"
)

type Websocket struct {
	Connection *websocket.Conn

	HeartbeatInterval time.Duration

	LastHeartBeat *time.Time

	SessionID *string

	LastEventNum atomic.Pointer[int]

	ReconnectURL *string

	Closed chan struct{}

	Ready chan struct{}

	writeMu              sync.Mutex
	heartbeatMu          sync.Mutex
	lastHeartbeatSentAt  time.Time
	awaitingHeartbeatAck bool
	missedHeartbeatAcks  int
	closeOnce            sync.Once
	readyOnce            sync.Once
}

func NewWebsocket(bot *Client, host string, isReconnect bool, lastEventNum *int, sessionID *string) (*Websocket, error) {
	if bot.token == nil {
		return nil, errors.New("bot token is nil")
	}

	dialer := *websocket.DefaultDialer
	dialer.HandshakeTimeout = 30 * time.Second

	c, _, err := dialer.Dial(host+"?v="+(*bot.APIVersion).ToString()+"&encoding=json", nil)
	if err != nil {
		return nil, err
	}

	_ = c.SetWriteDeadline(time.Time{})
	c.SetReadLimit(4 * 1024 * 1024) // 4 MiB — protects against rogue/compromised endpoint flooding RAM

	c.SetPongHandler(func(string) error {
		bot.Logger.Debug("Received pong from Discord")
		return nil
	})

	_, message, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}

	var payload common.Payload
	if err := json.Unmarshal(message, &payload); err != nil {
		return nil, err
	}

	var hello common.HelloPayloadData
	if err := json.Unmarshal(payload.D, &hello); err != nil {
		return nil, err
	}

	ws := &Websocket{
		Connection:        c,
		HeartbeatInterval: time.Millisecond * time.Duration(hello.HeartbeatInterval),
		Closed:            make(chan struct{}),
		Ready:             make(chan struct{}),
		SessionID:         sessionID,
	}

	bot.Logger.Info("Connected to Discord gateway", slog.Float64("heartbeatIntervalMs", hello.HeartbeatInterval))

	if isReconnect && lastEventNum != nil && sessionID != nil {
		if err := ws.writeJSON(map[string]interface{}{
			"op": 6,
			"d": map[string]interface{}{
				"token":      *bot.token,
				"session_id": *sessionID,
				"seq":        *lastEventNum,
			},
		}); err != nil {
			return nil, err
		}
	} else {
		data := map[string]interface{}{
			"op": 2,
			"d": map[string]interface{}{
				"token":   *bot.token,
				"intents": *bot.Intents,
				"properties": map[string]string{
					"os":      runtime.GOOS,
					"browser": "https://github.com/streame-gg/go-discord-wrapper@alpha",
					"device":  "https://github.com/streame-gg/go-discord-wrapper@alpha",
				},
			},
		}

		if bot.Sharding != nil {
			data["d"].(map[string]interface{})["shard"] = []int{bot.Sharding.ShardID, bot.Sharding.TotalShards}
		}

		if err := ws.writeJSON(data); err != nil {
			return nil, err
		}
	}

	// Heartbeat goroutine starts only after a successful Identify/Resume so that
	// a write failure does not leave an orphaned goroutine waiting on ws.Closed.
	go func() {
		timeoutThreshold := ws.HeartbeatInterval * 3
		if timeoutThreshold < 10*time.Second {
			timeoutThreshold = 10 * time.Second
		}
		const maxMissedHeartbeatAcks = 3

		// Returns false if the goroutine should stop.
		checkTimeout := func() bool {
			ws.heartbeatMu.Lock()
			if !ws.awaitingHeartbeatAck || time.Since(ws.lastHeartbeatSentAt) <= timeoutThreshold {
				ws.heartbeatMu.Unlock()
				return true
			}
			ws.missedHeartbeatAcks++
			missedAcks := ws.missedHeartbeatAcks
			ws.heartbeatMu.Unlock()

			bot.Logger.Warn("Heartbeat ACK timeout", slog.Int("missed", missedAcks), slog.Int("max", maxMissedHeartbeatAcks))
			if missedAcks >= maxMissedHeartbeatAcks {
				bot.Logger.Warn("Heartbeat ACK timeout threshold reached, closing connection")
				ws.close()
				return false
			}
			return true
		}

		// Returns false if the goroutine should stop.
		sendHeartbeat := func() bool {
			var heartbeatData json.RawMessage
			if seq := ws.LastEventNum.Load(); seq != nil {
				data, _ := json.Marshal(*seq)
				heartbeatData = data
			} else {
				heartbeatData = json.RawMessage("null")
			}

			if err := ws.writeJSONDeadline(common.Payload{Op: 1, D: heartbeatData}, 10*time.Second); err != nil {
				bot.Logger.Error("Failed to send heartbeat", slog.Any("err", err))
				if websocket.IsUnexpectedCloseError(err) {
					bot.Logger.Warn("Heartbeat failed due to closed connection, stopping heartbeat loop")
					return false
				}
				return true
			}

			ws.heartbeatMu.Lock()
			now := time.Now()
			ws.lastHeartbeatSentAt = now
			ws.awaitingHeartbeatAck = true
			if ws.LastHeartBeat == nil {
				ws.LastHeartBeat = &now
			}
			ws.heartbeatMu.Unlock()
			bot.Logger.Debug("Heartbeat sent")
			return true
		}

		// Discord Gateway spec §4.2: first heartbeat must fire after heartbeat_interval * jitter (0–1).
		jitter := time.Duration(rand.Float64() * float64(ws.HeartbeatInterval))
		select {
		case <-time.After(jitter):
		case <-ws.Closed:
			bot.Logger.Debug("Heartbeat stopped: websocket closed")
			return
		}

		if !sendHeartbeat() {
			return
		}

		ticker := time.NewTicker(ws.HeartbeatInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if !checkTimeout() {
					return
				}
				if !sendHeartbeat() {
					return
				}
			case <-ws.Closed:
				bot.Logger.Debug("Heartbeat stopped: websocket closed")
				return
			}
		}
	}()

	return ws, nil
}

func (d *Client) connectWebsocket(url string, isReconnect bool, lastEventNum *int, sessionID *string) error {
	ws, err := NewWebsocket(d, url, isReconnect, lastEventNum, sessionID)
	if err != nil {
		return err
	}

	d.wsMu.Lock()
	d.Websocket = ws
	d.wsMu.Unlock()
	return nil
}

func (d *Client) reconnect(freshConnect bool) error {
	d.reconnectMu.Lock()
	if d.reconnecting {
		d.reconnectMu.Unlock()
		d.Logger.Debug("Reconnect already running, skipping duplicate attempt")
		return nil
	}
	d.reconnecting = true
	d.reconnectMu.Unlock()

	defer func() {
		d.reconnectMu.Lock()
		d.reconnecting = false
		d.reconnectMu.Unlock()
	}()

	d.Logger.Warn("Reconnecting to Discord gateway")

	// On a fresh reconnect Discord re-sends GUILD_CREATE for every guild and
	// rebuilds all voice states from scratch. Wipe the local map now so stale
	// entries from users who left channels during the disconnection don't
	// produce wrong OldState values on the next VOICE_STATE_UPDATE.
	if freshConnect {
		d.voiceStatesMu.Lock()
		d.voiceStates = make(map[string]*common.VoiceState)
		d.voiceStatesMu.Unlock()
	}

	var lastEventNum *int
	var sessionID *string
	var reconnectURL string

	d.wsMu.Lock()
	if d.Websocket != nil {
		lastEventNum = d.Websocket.LastEventNum.Load()
		sessionID = d.Websocket.SessionID

		if !freshConnect && d.Websocket.ReconnectURL != nil {
			reconnectURL = *d.Websocket.ReconnectURL
		}

		d.Websocket.close()
		d.Websocket = nil
	}
	d.wsMu.Unlock()

	if reconnectURL == "" {
		reconnectURL = "wss://gateway.discord.gg"
	}

	maxRetries := d.maxReconnectRetries
	if maxRetries == 0 {
		maxRetries = 3
	}

	for i := 0; maxRetries < 0 || i < maxRetries; i++ {
		if d.shutdown.Load() {
			return nil
		}

		if i > 0 {
			const maxBackoff = 30 * time.Second
			exp := time.Duration(1<<uint(i-1)) * time.Second // 1s, 2s, 4s, 8s, …
			if exp > maxBackoff {
				exp = maxBackoff
			}
			backoff := exp + time.Duration(rand.Int63n(int64(exp)+1)) // jitter: [exp, 2*exp]
			if maxRetries < 0 {
				d.Logger.Debug("Waiting before retry", slog.Duration("backoff", backoff), slog.Int("attempt", i+1))
			} else {
				d.Logger.Debug("Waiting before retry", slog.Duration("backoff", backoff), slog.Int("attempt", i+1), slog.Int("max", maxRetries))
			}
			time.Sleep(backoff)
		}

		if err := d.connectWebsocket(reconnectURL, !freshConnect, lastEventNum, sessionID); err != nil {
			if maxRetries < 0 {
				d.Logger.Warn("Reconnect attempt failed, retrying", slog.Int("attempt", i+1), slog.Any("err", err))
			} else {
				d.Logger.Warn("Reconnect attempt failed", slog.Int("attempt", i+1), slog.Int("max", maxRetries), slog.Any("err", err))
				if i == maxRetries-1 {
					return err
				}
			}
			continue
		}

		if !freshConnect && sessionID != nil {
			d.wsMu.RLock()
			if d.Websocket != nil {
				d.Websocket.SessionID = sessionID
			}
			d.wsMu.RUnlock()
		}

		d.Logger.Info("Successfully reconnected to gateway")
		// For a session resume, emitReconnect is called by the RESUMED event
		// handler to avoid firing the callback twice (once here and once there).
		if freshConnect {
			d.emitReconnect()
		}
		return nil
	}

	return fmt.Errorf("failed to reconnect after %d attempts", maxRetries)
}

func (d *Client) listenWebsocket() error {
	for {
		d.wsMu.RLock()
		ws := d.Websocket
		d.wsMu.RUnlock()

		if ws == nil {
			return nil
		}

		// Refresh read deadline before blocking — if no message arrives within
		// two heartbeat intervals the connection is considered a zombie.
		if ws.HeartbeatInterval > 0 {
			_ = ws.Connection.SetReadDeadline(time.Now().Add(ws.HeartbeatInterval * 2))
		}

		_, message, err := ws.Connection.ReadMessage()
		if err != nil {
			return err
		}

		var payload common.Payload
		if err := json.Unmarshal(message, &payload); err != nil {
			return err
		}

		d.Logger.Debug("Received payload", slog.String("type", payload.T), slog.Int("op", int(payload.Op)))

		if payload.Op == 6 {
			d.Logger.Debug("Resuming session")
		}

		if payload.Op == 7 {
			d.Logger.Debug("Reconnecting to gateway; requested by Discord")
			if err := d.reconnect(false); err != nil {
				d.Logger.Error("Failed to reconnect", slog.Any("err", err))
				return err
			}

			return nil
		}

		if payload.Op == 9 {
			var canResume bool
			if err := json.Unmarshal(payload.D, &canResume); err != nil {
				return err
			}

			if canResume {
				d.Logger.Debug("Invalid session, attempting to resume")
				if err := d.reconnect(false); err != nil {
					d.Logger.Error("Failed to resume session", slog.Any("err", err))
					return err
				}
			} else {
				d.Logger.Debug("Invalid session, re-identifying")
				if err := d.reconnect(true); err != nil {
					d.Logger.Error("Failed to re-identify", slog.Any("err", err))
					return err
				}
			}

			return nil
		}

		if payload.Op == 11 {
			now := time.Now()
			ws.heartbeatMu.Lock()
			ws.awaitingHeartbeatAck = false
			ws.missedHeartbeatAcks = 0
			ws.LastHeartBeat = &now
			ws.heartbeatMu.Unlock()
			d.Logger.Debug("Heartbeat ACK received")
		}

		if payload.S != nil {
			ws.LastEventNum.Store(payload.S)
		}

		if payload.T != "" {
			factory, ok := events.EventFactories[events.EventType(payload.T)]
			if !ok {
				d.Logger.Warn("No factory found for event type", slog.String("type", string(payload.T)))
				continue
			}

			event := factory()

			if err := json.Unmarshal(payload.D, event); err != nil {
				var anyVal any
				_ = json.Unmarshal(payload.D, &anyVal)
				d.Logger.Debug("Failed event payload", slog.Any("payload", anyVal))
				d.Logger.Error("Failed to unmarshal event", slog.String("type", string(payload.T)), slog.Any("err", err))
				d.emitPacketError(err)
				continue
			}

			// Apply cache updates synchronously on the websocket reader goroutine.
			// This preserves the serial ordering Discord guarantees on the wire:
			// concurrent goroutines / workers would otherwise race and let
			// ROLE_UPDATE land after MESSAGE_CREATE that depends on the new role.
			canDispatch := func() (ok bool) {
				defer func() {
					if r := recover(); r != nil {
						d.Logger.Error("panic in internal event handler", slog.Any("recover", r))
					}
				}()
				return d.internalEventHandler(payload.D, event.Event(), event)
			}()

			if canDispatch {
				job := dispatchJob{event: event}
				if d.eventCh != nil {
					// Worker-pool mode: enqueue or drop if queue is full.
					// shutdownCh case prevents sending on a closed channel after Shutdown.
					select {
					case d.eventCh <- job:
					case <-d.shutdownCh:
						return nil
					default:
						d.Logger.Warn("Event queue full, dropping event",
							slog.String("type", string(payload.T)),
							slog.Int("queueCap", cap(d.eventCh)),
						)
					}
				} else {
					// Unlimited mode (default): spawn one goroutine per event.
					// Track via eventWg so Shutdown can drain before returning.
					d.eventWg.Add(1)
					go func() {
						defer d.eventWg.Done()
						d.processEvent(job)
					}()
				}
			}
		}
	}
}

func (ws *Websocket) writeJSON(v any) error {
	ws.writeMu.Lock()
	defer ws.writeMu.Unlock()
	return ws.Connection.WriteJSON(v)
}

func (ws *Websocket) writeJSONDeadline(v any, d time.Duration) error {
	ws.writeMu.Lock()
	defer ws.writeMu.Unlock()
	if err := ws.Connection.SetWriteDeadline(time.Now().Add(d)); err != nil {
		return err
	}
	err := ws.Connection.WriteJSON(v)
	_ = ws.Connection.SetWriteDeadline(time.Time{})
	return err
}

func (d *Websocket) close() error {
	if d == nil {
		return nil
	}
	var err error
	d.closeOnce.Do(func() {
		err = d.Connection.Close()
		close(d.Closed)
	})
	return err
}
