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
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type EventHandler func(*DiscordClient, types.DiscordEvent)

type DiscordClient struct {
	Token *string

	APIVersion *types.DiscordAPIVersion

	Logger *zerolog.Logger

	Intents *types.DiscordIntent

	Websocket *websocket.Conn

	Events map[types.DiscordEventType][]EventHandler

	mu sync.RWMutex

	LastEventNum *int

	ReconnectURL *string

	SessionID *string

	UnavailableGuilds map[types.DiscordSnowflake]struct{}

	LastHeartbeat *time.Time
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

	defer func() {
		_ = do.Body.Close()
	}()

	if do.StatusCode != http.StatusOK {
		return nil, errors.New("failed to register bot gateway connection, status code: " + do.Status)
	}

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

	if err := d.connectWebsocket(gatewayResp.Url, false); err != nil {
		return err
	}

	go func() {
		if err := d.listenWebsocket(); err != nil {
			d.Logger.Err(err).Msg("Error listening to websocket")
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, 4000, 4001, 4002, 4003, 4005, 4007, 4008, 4009) {
				d.Logger.Info().Msg("Discord gateway connection closed by Discord, trying to reconnect")
				if err := d.reconnect(true); err != nil {
					d.Logger.Err(err).Msg("Failed to reconnect")
				}
			}

			d.Logger.Info().Msg("Discord gateway connection closed by Discord, no reconnecting attempt will be made")

			return
		}
	}()

	return nil
}

func (d *DiscordClient) onEvent(
	eventName types.DiscordEventType,
	handler EventHandler,
) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.Events == nil {
		d.Events = make(map[types.DiscordEventType][]EventHandler)
	}

	d.Events[eventName] = append(d.Events[eventName], handler)
}

func (d *DiscordClient) OnGuildCreate(
	handler func(*DiscordClient, *types.DiscordGuildCreateEvent),
) {
	d.onEvent(types.DiscordEventReady, func(
		session *DiscordClient,
		event types.DiscordEvent,
	) {
		if e, ok := event.(*types.DiscordGuildCreateEvent); ok {
			handler(session, e)
		}
	})
}

func (d *DiscordClient) OnMessageCreate(
	handler func(*DiscordClient, *types.DiscordMessageCreateEvent),
) {
	d.onEvent(types.DiscordEventMessageCreate, func(
		session *DiscordClient,
		event types.DiscordEvent,
	) {
		if e, ok := event.(*types.DiscordMessageCreateEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to MessageCreateEvent: %T", event)
		}
	})
}

func (d *DiscordClient) dispatch(event types.DiscordEvent) {
	handlers := d.Events[event.Event()]
	for _, h := range handlers {
		h(d, event)
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

				d.Logger.Debug().Msgf("Guild %s is available again", guildCreateEvent.Guild.GetID())
				d.deleteUnavailableGuild(guildCreateEvent.Guild.GetID())

				return false
			}
		}
	default:
		return true
	}

	return true
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
