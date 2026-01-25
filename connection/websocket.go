package connection

import (
	"encoding/json"
	"go-discord-wrapper/types"
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
	c, _, err := websocket.DefaultDialer.Dial(host+"?v=10&encoding=json", nil)
	if err != nil {
		return nil, err
	}

	_, message, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}

	var payload types.Payload
	if err := json.Unmarshal(message, &payload); err != nil {
		return nil, err
	}

	var hello types.HelloPayloadData
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

		go func() {
			for range ticker.C {
				if ws.LastHeartBeat != nil && !ws.LastHeartBeat.IsZero() &&
					time.Since(*ws.LastHeartBeat) > time.Duration(hello.HeartbeatInterval)*time.Millisecond*2 {

					bot.Logger.Warn().Msg("Heartbeat ACK timeout, reconnecting")
					ws.close()
					_ = bot.reconnect(true)
					return
				}

				if err := c.WriteJSON(types.Payload{
					Op: 1,
					D:  nil,
				}); err != nil {
					bot.Logger.Err(err).Msg("Failed to send heartbeat")
					ws.close()
					return
				}

				bot.Logger.Debug().Msg("Heartbeat sent")
			}
		}()

		<-ws.Closed
		return
	}()

	if isReconnect && lastEventNum != nil {
		if err := c.WriteJSON(map[string]interface{}{
			"op": 6,
			"d": map[string]interface{}{
				"token":      *bot.Token,
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
				"token":   *bot.Token,
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
	d.Logger.Warn().Msg("Reconnecting to  gateway")

	lastEventNum := d.Websocket.LastEventNum

	if d.Websocket != nil {
		d.Websocket.close()
		d.Websocket = nil
	}

	if err := d.connectWebsocket("wss://gateway.discord.gg", !freshConnect, lastEventNum); err != nil {
		return err
	}

	d.Logger.Debug().Msg("Reconnected to  gateway")

	go func() {
		if err := d.listenWebsocket(); err != nil {
			d.Logger.Err(err).Msg("Error listening to websocket")
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, 4000, 4001, 4002, 4003, 4005, 4007, 4008, 4009) {
				d.Logger.Debug().Msg("gateway connection closed by , trying to reconnect")
				if err := d.reconnect(true); err != nil {
					d.Logger.Err(err).Msg("Failed to reconnect")
				}
			}

			d.Logger.Debug().Msg("gateway connection closed by , no reconnecting attempt will be made")

			return
		}
	}()

	return nil
}

func (d *Client) listenWebsocket() error {
	for {
		_, message, err := d.Websocket.Connection.ReadMessage()
		if err != nil {
			return err
		}

		var payload types.Payload
		if err := json.Unmarshal(message, &payload); err != nil {
			return err
		}

		d.Logger.Debug().Msgf("Received payload: %s %d", payload.T, payload.Op)

		if payload.Op == 6 {
			d.Logger.Debug().Msg("Resuming session")
		}

		if payload.Op == 7 {
			d.Logger.Debug().Msg("Reconnecting to gateway; requested by ")
			return d.reconnect(false)
		}

		if payload.Op == 9 {
			var invalidSession types.InvalidSessionPayload
			if err := json.Unmarshal(payload.D, &invalidSession); err != nil {
				return err
			}

			if invalidSession.D {
				d.Logger.Debug().Msg("Invalid session, re-identifying")
				if err := d.reconnect(true); err != nil {
					return err
				}
			} else {
				d.Logger.Debug().Msg("Invalid session, attempting to resume")
				if err := d.reconnect(false); err != nil {
					return err
				}
			}
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
			factory, ok := types.EventFactories[types.EventType(payload.T)]
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
		close(d.Closed)
	}
}
