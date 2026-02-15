package connection

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
	"github.com/DatGamet/go-discord-wrapper/types/events"
	"time"

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
}

func NewWebsocket(bot *Client, host string, isReconnect bool, lastEventNum *int) (*Websocket, error) {
	dialer := *websocket.DefaultDialer
	dialer.HandshakeTimeout = 30 * time.Second

	c, _, err := dialer.Dial(host+"?v=10&encoding=json", nil)
	if err != nil {
		return nil, err
	}

	_ = c.SetReadDeadline(time.Time{})
	_ = c.SetWriteDeadline(time.Time{})

	c.SetPongHandler(func(string) error {
		bot.Logger.Debug().Msg("Received pong from Discord")
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
	}

	bot.Logger.Info().Msgf("Connected to Discord gateway, heartbeat interval: %f ms", hello.HeartbeatInterval)

	go func() {
		ticker := time.NewTicker(time.Millisecond * time.Duration(hello.HeartbeatInterval))
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				{
					if ws.LastHeartBeat != nil && !ws.LastHeartBeat.IsZero() &&
						time.Since(*ws.LastHeartBeat) > time.Duration(hello.HeartbeatInterval)*time.Millisecond*2 {

						bot.Logger.Warn().Msg("Heartbeat ACK timeout, reconnecting")
						ws.close()
						_ = bot.reconnect(true)
						return
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

					// Set write deadline to prevent hanging
					if err := c.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
						bot.Logger.Err(err).Msg("Failed to set write deadline")
						return
					}

					if err := c.WriteJSON(heartbeatPayload); err != nil {
						bot.Logger.Err(err).Msg("Failed to send heartbeat")

						if websocket.IsUnexpectedCloseError(err) {
							bot.Logger.Warn().Msg("Heartbeat failed due to closed connection, stopping heartbeat loop")
							return
						}
						continue
					}

					_ = c.SetWriteDeadline(time.Time{})

					bot.Logger.Debug().Msg("Heartbeat sent")
				}
			case <-ws.Closed:
				bot.Logger.Debug().Msg("Heartbeat stopped: websocket closed")
				return
			}
		}
	}()

	if isReconnect && lastEventNum != nil {
		if err := c.WriteJSON(map[string]interface{}{
			"op": 6,
			"d": map[string]interface{}{
				"token":      *bot.token,
				"session_id": ws.SessionID,
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
					"$browser": "https://github.com/DatGamet/go-discord-wrapper@alpha",
					"$device":  "https://github.com/DatGamet/go-discord-wrapper@alpha",
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

func (d *Client) connectWebsocket(url string, isReconnect bool, lastEventNum *int) error {
	ws, err := NewWebsocket(d, url, isReconnect, lastEventNum)
	if err != nil {
		return err
	}

	d.Websocket = ws
	return nil
}

func (d *Client) reconnect(freshConnect bool) error {
	d.Logger.Warn().Msg("Reconnecting to Discord gateway")

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
			d.Logger.Debug().Msgf("Waiting %v before retry %d/%d", backoff, i+1, maxRetries)
			time.Sleep(backoff)
		}

		if err := d.connectWebsocket(reconnectURL, !freshConnect, lastEventNum); err != nil {
			d.Logger.Warn().Msgf("Reconnect attempt %d/%d failed: %v", i+1, maxRetries, err)
			if i == maxRetries-1 {
				return err
			}
			continue
		}

		if !freshConnect && sessionID != nil {
			d.Websocket.SessionID = sessionID
		}

		d.Logger.Info().Msg("Successfully reconnected to gateway")
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

		d.Logger.Debug().Msgf("Received payload: %s %d", payload.T, payload.Op)

		if payload.Op == 6 {
			d.Logger.Debug().Msg("Resuming session")
		}

		if payload.Op == 7 {
			d.Logger.Debug().Msg("Reconnecting to gateway; requested by Discord")
			if err := d.reconnect(false); err != nil {
				d.Logger.Err(err).Msg("Failed to reconnect")
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
				d.Logger.Debug().Msg("Invalid session, attempting to resume")
				if err := d.reconnect(false); err != nil {
					d.Logger.Err(err).Msg("Failed to resume session")
					return err
				}
			} else {
				d.Logger.Debug().Msg("Invalid session, re-identifying")
				if err := d.reconnect(true); err != nil {
					d.Logger.Err(err).Msg("Failed to re-identify")
					return err
				}
			}

			return nil
		}

		if payload.Op == 11 {
			now := time.Now()
			d.Websocket.LastHeartBeat = &now
			d.Logger.Debug().Msg("Heartbeat Ack Received")
		}

		if payload.S != nil {
			d.Websocket.LastEventNum = payload.S
		}

		if payload.T != "" {
			factory, ok := events.EventFactories[events.EventType(payload.T)]
			if !ok {
				d.Logger.Warn().Msgf("No factory found for event type %s", payload.T)
				continue
			}

			event := factory()

			if err := json.Unmarshal(payload.D, event); err != nil {
				var anyVal any
				_ = json.Unmarshal(payload.D, &anyVal)
				d.Logger.Debug().Msgf("Failed Event payload: %+v", anyVal)

				d.Logger.Err(err).Msgf("Failed to unmarshal event %s", payload.T)
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
