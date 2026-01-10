package connection

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-discord-wrapper/functions"
	"go-discord-wrapper/types"
	"go-discord-wrapper/util"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type DiscordClient struct {
	Token *string

	APIVersion *types.DiscordAPIVersion

	Logger *zerolog.Logger

	Intents *types.DiscordIntent

	Websocket *websocket.Conn

	Events map[types.DiscordEventType]func(session *DiscordClient, event types.DiscordEvent)

	mu sync.RWMutex

	LastEventNum *int

	ReconnectURL *string

	SessionID *string

	UnavailableGuilds map[types.DiscordSnowflake]struct{}
}

func NewDiscordClient(token string, intents types.DiscordIntent) *DiscordClient {
	return &DiscordClient{
		Token:             &token,
		APIVersion:        functions.PointerTo(types.DiscordAPIVersion10),
		Logger:            util.NewLogger(),
		Intents:           &intents,
		UnavailableGuilds: make(map[types.DiscordSnowflake]struct{}),
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

func OnEvent[T types.DiscordEvent](d *DiscordClient, event types.DiscordEventType, handler func(*DiscordClient, T)) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.Events == nil {
		d.Events = make(map[types.DiscordEventType]func(session *DiscordClient, event types.DiscordEvent))
	}

	d.Events[event] = func(
		session *DiscordClient,
		ev types.DiscordEvent,
	) {
		typed, ok := ev.(T)
		if !ok {
			session.Logger.Warn().
				Str("expected", fmt.Sprintf("%T", *new(T))).
				Str("got", fmt.Sprintf("%T", ev)).
				Msg("event type mismatch")
			return
		}
		handler(session, typed)
	}
}

func (d *DiscordClient) dispatch(event types.DiscordEventType, payload types.Payload) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.Events == nil {
		return
	}

	rawEvent := d.Events[event]
	if rawEvent != nil {
		discordEvent, err := d.convertToEvent(event, payload.D)
		if err != nil {
			d.Logger.Err(err).Msgf("Failed to convert event %s", event)
			return
		}

		go func() {
			if con := d.internalEventHandler(payload.D, event); con {
				rawEvent(d, discordEvent)
			}
		}()
	}
}

func (d *DiscordClient) internalEventHandler(msg json.RawMessage, event types.DiscordEventType) bool {
	switch event {
	case types.DiscordEventReady:
		{
			var readyEvent types.DiscordReadyPayload
			if err := json.Unmarshal(msg, &readyEvent); err != nil {
				d.Logger.Err(err).Msg("Failed to unmarshal READY event")
			}

			d.SessionID = &readyEvent.SessionID
			d.ReconnectURL = &readyEvent.ResumeGatewayURL

			for _, guild := range readyEvent.Guilds {
				if !guild.Guild.IsAvailable() {
					d.addUnavailableGuild(guild.Guild.GetID())
				}
			}

			return true
		}
	case types.DiscordEventGuildCreate:
		{
			var guildCreateEvent types.DiscordGuildCreateEvent
			if err := json.Unmarshal(msg, &guildCreateEvent); err != nil {
				d.Logger.Err(err).Msg("Failed to unmarshal GUILD_CREATE event")
				return false
			}

			if guildCreateEvent.Guild.IsAvailable() && d.IsGuildUnavailable(guildCreateEvent.Guild.GetID()) {
				//Discord is telling us that a guild is available again by firing a GUILD_CREATE event with available = true

				d.Logger.Info().Msgf("Guild %s is available again", guildCreateEvent.Guild.GetID())
				d.deleteUnavailableGuild(guildCreateEvent.Guild.GetID())

				return false
			}
		}
	default:
		return true
	}

	return true
}

func (d *DiscordClient) convertToEvent(event types.DiscordEventType, data json.RawMessage) (types.DiscordEvent, error) {
	switch event {
	case "MESSAGE_CREATE":
		var msgCreateEvent types.DiscordMessageCreateEvent
		if err := json.Unmarshal(data, &msgCreateEvent); err != nil {
			return nil, err
		}
		return &msgCreateEvent, nil
	case "READY":
		var readyEvent types.DiscordReadyEvent
		if err := json.Unmarshal(data, &readyEvent); err != nil {
			return nil, err
		}
		return &readyEvent, nil
	case "GUILD_CREATE":
		var guildCreateEvent types.DiscordGuildCreateEvent
		if err := json.Unmarshal(data, &guildCreateEvent); err != nil {
			return nil, err
		}
		return &guildCreateEvent, nil
	default:
		return nil, errors.New("unsupported event type: " + string(event))
	}
}

func (d *DiscordClient) addUnavailableGuild(id types.DiscordSnowflake) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.UnavailableGuilds[id] = struct{}{}
}

func (d *DiscordClient) deleteUnavailableGuild(id types.DiscordSnowflake) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.UnavailableGuilds, id)
}

func (d *DiscordClient) IsGuildUnavailable(id types.DiscordSnowflake) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, exists := d.UnavailableGuilds[id]
	return exists
}

func (d *DiscordClient) Shutdown() error {
	if d.Websocket != nil {
		return d.Websocket.Close()
	}
	return nil
}
