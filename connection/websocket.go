package connection

import (
	"encoding/json"
	"log/slog"
	"sync"
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

	LastEventNum *int

	ReconnectURL *string

	Closed chan struct{}

	Ready chan struct{}

	heartbeatMu          sync.Mutex
	lastHeartbeatSentAt  time.Time
	awaitingHeartbeatAck bool
	missedHeartbeatAcks  int
}

func NewWebsocket(bot *Client, host string, isReconnect bool, lastEventNum *int, sessionID *string) (*Websocket, error) {
	dialer := *websocket.DefaultDialer
	dialer.HandshakeTimeout = 30 * time.Second

	c, _, err := dialer.Dial(host+"?v=10&encoding=json", nil)
	if err != nil {
		return nil, err
	}

	_ = c.SetReadDeadline(time.Time{})
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
								bot.Logger.Warn("Heartbeat ACK timeout threshold reached, reconnecting")
								ws.close()
								_ = bot.reconnect(false)
								return
							}
						} else {
							ws.heartbeatMu.Unlock()
						}
					} else {
						ws.heartbeatMu.Unlock()
					}

					var heartbeatData json.RawMessage
					if ws.LastEventNum != nil {
						data, _ := json.Marshal(*ws.LastEventNum)
						heartbeatData = data
					} else {
						heartbeatData = json.RawMessage("null")
					}

					heartbeatPayload := common.Payload{
						Op: 1,
						D:  heartbeatData,
					}

					if err := c.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
						bot.Logger.Error("Failed to set write deadline", slog.Any("err", err))
						return
					}

					if err := c.WriteJSON(heartbeatPayload); err != nil {
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

					_ = c.SetWriteDeadline(time.Time{})

					bot.Logger.Debug("Heartbeat sent")
				}
			case <-ws.Closed:
				bot.Logger.Debug("Heartbeat stopped: websocket closed")
				return
			}
		}
	}()

	if isReconnect && lastEventNum != nil && sessionID != nil {
		if err := c.WriteJSON(map[string]interface{}{
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
					"$os":      "windows",
					"$browser": "https://github.com/streame-gg/go-discord-wrapper@alpha",
					"$device":  "https://github.com/streame-gg/go-discord-wrapper@alpha",
				},
			},
		}

		if bot.Sharding != nil {
			data["d"].(map[string]interface{})["shard"] = []int{bot.Sharding.ShardID, bot.Sharding.TotalShards}
		}

		if err := c.WriteJSON(data); err != nil {
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

	d.Websocket = ws
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

	var lastEventNum *int
	var sessionID *string
	var reconnectURL string

	if d.Websocket != nil {
		lastEventNum = d.Websocket.LastEventNum
		sessionID = d.Websocket.SessionID

		if !freshConnect && d.Websocket.ReconnectURL != nil {
			reconnectURL = *d.Websocket.ReconnectURL
		}

		d.Websocket.close()
		d.Websocket = nil
	}

	if reconnectURL == "" {
		reconnectURL = "wss://gateway.discord.gg"
	}

	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			backoff := time.Duration(i) * time.Second
			d.Logger.Debug("Waiting before retry", slog.Duration("backoff", backoff), slog.Int("attempt", i+1), slog.Int("max", maxRetries))
			time.Sleep(backoff)
		}

		if err := d.connectWebsocket(reconnectURL, !freshConnect, lastEventNum, sessionID); err != nil {
			d.Logger.Warn("Reconnect attempt failed", slog.Int("attempt", i+1), slog.Int("max", maxRetries), slog.Any("err", err))
			if i == maxRetries-1 {
				return err
			}
			continue
		}

		if !freshConnect && sessionID != nil {
			d.Websocket.SessionID = sessionID
		}

		d.Logger.Info("Successfully reconnected to gateway")
		return nil
	}

	return nil
}

func (d *Client) listenWebsocket() error {
	for {
		_, message, err := d.Websocket.Connection.ReadMessage()
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
			d.Websocket.heartbeatMu.Lock()
			d.Websocket.awaitingHeartbeatAck = false
			d.Websocket.missedHeartbeatAcks = 0
			d.Websocket.LastHeartBeat = &now
			d.Websocket.heartbeatMu.Unlock()
			d.Logger.Debug("Heartbeat ACK received")
		}

		if payload.S != nil {
			d.Websocket.LastEventNum = payload.S
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
				continue
			}

			go func() {
				if canContinue := d.internalEventHandler(payload.D, event.Event()); canContinue {
					d.dispatch(event)
				}
			}()
		}
	}
}

func (d *Websocket) close() {
	if d != nil {
		_ = d.Connection.Close()
		select {
		case <-d.Closed:
		default:
			close(d.Closed)
		}
	}
}
