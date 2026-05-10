package connection

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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
}

func NewWebsocket(bot *Client, host string, isReconnect bool, lastEventNum *int, sessionID *string) (*Websocket, error) {
	if bot.token == nil {
		return nil, errors.New("bot token is nil")
	}

	dialer := *websocket.DefaultDialer
	dialer.HandshakeTimeout = 30 * time.Second

	c, _, err := dialer.Dial(host+"?v=10&encoding=json", nil)
	if err != nil {
		return nil, err
	}

	_ = c.SetWriteDeadline(time.Time{})

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

	go func() {
		ticker := time.NewTicker(ws.HeartbeatInterval)
		defer ticker.Stop()

		timeoutThreshold := ws.HeartbeatInterval * 3
		if timeoutThreshold < 10*time.Second {
			timeoutThreshold = 10 * time.Second
		}

		const maxMissedHeartbeatAcks = 3

		for {
			select {
			case <-ticker.C:
				{
					ws.heartbeatMu.Lock()
					if ws.awaitingHeartbeatAck {
						timeSinceLastSend := time.Since(ws.lastHeartbeatSentAt)
						if timeSinceLastSend > timeoutThreshold {
							ws.missedHeartbeatAcks++
							missedAcks := ws.missedHeartbeatAcks
							ws.heartbeatMu.Unlock()

							bot.Logger.Warn("Heartbeat ACK timeout", slog.Int("missed", missedAcks), slog.Int("max", maxMissedHeartbeatAcks))

							if missedAcks >= maxMissedHeartbeatAcks {
								bot.Logger.Warn("Heartbeat ACK timeout threshold reached, closing connection")
								ws.close()
								return
							}
						} else {
							ws.heartbeatMu.Unlock()
						}
					} else {
						ws.heartbeatMu.Unlock()
					}

					var heartbeatData json.RawMessage
					if seq := ws.LastEventNum.Load(); seq != nil {
						data, _ := json.Marshal(*seq)
						heartbeatData = data
					} else {
						heartbeatData = json.RawMessage("null")
					}

					heartbeatPayload := common.Payload{
						Op: 1,
						D:  heartbeatData,
					}

					if err := ws.writeJSONDeadline(heartbeatPayload, 10*time.Second); err != nil {
						bot.Logger.Error("Failed to send heartbeat", slog.Any("err", err))

						if websocket.IsUnexpectedCloseError(err) {
							bot.Logger.Warn("Heartbeat failed due to closed connection, stopping heartbeat loop")
							return
						}
						continue
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
				}
			case <-ws.Closed:
				bot.Logger.Debug("Heartbeat stopped: websocket closed")
				return
			}
		}
	}()

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
					"$os":      runtime.GOOS,
					"$browser": "https://github.com/streame-gg/go-discord-wrapper@alpha",
					"$device":  "https://github.com/streame-gg/go-discord-wrapper@alpha",
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
			backoff := time.Duration(i) * time.Second
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

			go func() {
				defer func() {
					if r := recover(); r != nil {
						d.Logger.Error("panic in event dispatch", slog.Any("recover", r))
					}
				}()
				if canContinue := d.internalEventHandler(payload.D, event.Event(), event); canContinue {
					d.dispatch(event)
				}
			}()
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
	if d != nil {
		err := d.Connection.Close()
		select {
		case <-d.Closed:
		default:
			close(d.Closed)
		}
		return err
	}
	return nil
}
