package connection

import (
	"encoding/json"
	"go-discord-wrapper/types"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func NewWebsocket(bot *DiscordClient, host string, isReconnect bool) (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial(host+"?v=10&encoding=json", nil)
	if err != nil {
		log.Fatal("dial:", err)
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

	go func() {
		ticker := time.NewTicker(time.Millisecond * time.Duration(hello.HeartbeatInterval))
		defer ticker.Stop()
		for range ticker.C {
			if bot.LastHeartbeat != nil && !bot.LastHeartbeat.IsZero() &&
				time.Since(*bot.LastHeartbeat) > time.Duration(hello.HeartbeatInterval)*time.Millisecond*2 {

				bot.Logger.Warn().Msg("Heartbeat ACK timeout, reconnecting")
				_ = c.Close()
				_ = bot.reconnect(true)
				return
			}

			if err := c.WriteJSON(types.Payload{
				Op: 1,
				D:  nil,
			}); err != nil {
				bot.Logger.Err(err).Msg("Failed to send heartbeat")
				return
			}

			bot.Logger.Debug().Msg("Heartbeat sent")
		}
	}()

	if isReconnect {
		if err := c.WriteJSON(map[string]interface{}{
			"op": 6,
			"d": map[string]interface{}{
				"token":      *bot.Token,
				"session_id": bot.SessionID,
				"seq":        *bot.LastEventNum,
			},
		}); err != nil {
			return nil, err
		}
	} else {
		if err := c.WriteJSON(map[string]interface{}{
			"op": 2,
			"d": map[string]interface{}{
				"token":   *bot.Token,
				"intents": *bot.Intents,
				"properties": map[string]string{
					"$os":      "windows",
					"$browser": "dat_bot_go",
					"$device":  "dat_bot_go",
				},
			},
		}); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (d *DiscordClient) connectWebsocket(url string, isReconnect bool) error {
	ws, err := NewWebsocket(d, url, isReconnect)
	if err != nil {
		return err
	}

	d.Websocket = ws
	return nil
}

func (d *DiscordClient) reconnect(freshConnect bool) error {
	d.Logger.Warn().Msg("Reconnecting to Discord gateway")

	if d.Websocket != nil {
		_ = d.Websocket.Close()
		d.Websocket = nil
	}

	if err := d.connectWebsocket("wss://gateway.discord.gg", !freshConnect); err != nil {
		return err
	}

	d.Logger.Info().Msg("Reconnected to Discord gateway")
	return nil
}

func (d *DiscordClient) listenWebsocket() error {
	for {
		_, message, err := d.Websocket.ReadMessage()
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
			d.Logger.Debug().Msg("Reconnecting to gateway; requested by Discord")
			return d.reconnect(false)
		}

		if payload.Op == 9 {
			var invalidSession types.DiscordInvalidSessionPayload
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
			d.LastHeartbeat = &now
			d.Logger.Debug().Msg("Heartbeat Ack Received")
		}

		if payload.S != nil {
			d.LastEventNum = payload.S
		}

		if payload.T != "" {
			factory, ok := types.EventFactories[payload.T]
			if !ok {
				d.Logger.Warn().Msgf("No factory found for event type %s", payload.T)
				continue
			}

			event := factory()

			if err := json.Unmarshal(payload.D, event); err != nil {
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
