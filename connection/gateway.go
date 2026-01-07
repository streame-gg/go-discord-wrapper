package connection

import (
	"encoding/json"
	"errors"
	"go-discord-wrapper/functions"
	"go-discord-wrapper/types"
	"go-discord-wrapper/util"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type DiscordClientEvent func(event types.Payload)

type DiscordClient struct {
	Token *string

	APIVersion *types.DiscordAPIVersion

	Logger *zerolog.Logger

	Intents *types.DiscordIntent

	Websocket *websocket.Conn

	Events map[string]DiscordClientEvent

	mu sync.RWMutex
}

func NewDiscordClient(token string, intents types.DiscordIntent) *DiscordClient {
	return &DiscordClient{
		Token:      &token,
		APIVersion: functions.PointerTo(types.DiscordAPIVersion10),
		Logger:     util.NewLogger(),
		Intents:    &intents,
	}
}

func (d *DiscordClient) initializeGatewayConnection() (*types.DiscordBotRegisterResponse, error) {
	do, err := http.DefaultClient.Do(&http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "https",
			Host:   "discord.com",
			Path:   types.DiscordAPIBaseString(*d.APIVersion) + types.DiscordAPIGatewayRequest,
		},
		Header: http.Header{
			"Authorization": []string{"Bot " + *d.Token},
		},
	})
	if err != nil {
		return nil, err
	}

	if do.StatusCode != http.StatusOK {
		return nil, errors.New("failed to register bot gateway connection, status code: " + do.Status)
	}

	defer func() {
		_ = do.Body.Close()
	}()

	var resp types.DiscordBotRegisterResponse
	if err := json.NewDecoder(do.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (d *DiscordClient) Login() error {
	gatewayResp, err := d.initializeGatewayConnection()
	if err != nil {
		return err
	}

	d.Logger.Info().Msgf("Connecting to gateway websocket at %s with %d shards", gatewayResp.Url, gatewayResp.Shards)

	if err := d.connectWebsocket(gatewayResp.Url); err != nil {
		return err
	}

	go func() {
		if err := d.listenWebsocket(); err != nil {
			d.Logger.Err(err).Msg("Error listening to websocket")
		}
	}()

	return nil
}

func UnwrapEvent[V any](payload types.Payload) (*V, error) {
	var data V
	err := json.Unmarshal(payload.D, &data)
	if err != nil {
		return nil, err
	}

	return &data, err
}

func (d *DiscordClient) On(event string, handler DiscordClientEvent) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.Events == nil {
		d.Events = make(map[string]DiscordClientEvent)
	}

	d.Events[event] = handler
}

func (d *DiscordClient) dispatch(event string, payload types.Payload) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.Events == nil {
		return
	}

	rawEvent := d.Events[event]
	if rawEvent != nil {
		go rawEvent(payload)
	}
}

func (d *DiscordClient) Shutdown() error {
	if d.Websocket != nil {
		return d.Websocket.Close()
	}
	return nil
}
