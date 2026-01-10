package connection

import (
	"encoding/json"
	"go-discord-wrapper/types"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// NewWebsocket TODO: websocket: close 1001 (going away): Discord WebSocket requesting client reconnect; Heartbeat no response handling
func NewWebsocket(bot *DiscordClient, host string) (*websocket.Conn, error) {
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
			hb := types.Payload{
				Op: 1,
				D:  nil,
			}
			if err := c.WriteJSON(hb); err != nil {
				bot.Logger.Err(err).Msg("Failed to send heartbeat")
				return
			}
			bot.Logger.Err(err).Msg("Sent heartbeat")
		}
	}()

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

	return c, nil
}

func (d *DiscordClient) connectWebsocket(url string) error {
	ws, err := NewWebsocket(d, url)
	if err != nil {
		return err
	}

	d.Websocket = ws
	return nil
}

func (d *DiscordClient) reconnectWebsocket() error {
	d.Logger.Info().Msg("Reconnecting to Discord WebSocket")

	oldWs := d.Websocket

	if err := d.connectWebsocket(*d.ReconnectURL); err != nil {
		return err
	}

	go func() {
		if err := d.listenWebsocket(); err != nil {
			d.Logger.Err(err).Msg("Error listening to websocket")
		}
	}()

	_ = oldWs.Close()
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

		if payload.Op == 7 {
			if err := d.reconnectWebsocket(); err != nil {
				return err
			}
			continue
		}

		if payload.S != nil {
			d.LastEventNum = payload.S
		}

		d.Logger.Info().Msgf("Received payload: %s %d", payload.T, payload.Op)

		if payload.T != "" {
			d.dispatch(types.DiscordEventType(payload.T), payload)
		}
	}
}
