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

type EventHandler func(*Client, types.Event)

type ClientSharding struct {
	TotalShards int
	ShardID     int
}

type Client struct {
	Token *string

	APIVersion *types.APIVersion

	Logger *zerolog.Logger

	Intents *types.Intent

	Websocket *Websocket

	Events map[types.EventType][]EventHandler

	mu sync.RWMutex

	UnavailableGuilds map[types.Snowflake]struct{}

	User *types.User

	Sharding *ClientSharding
}

func NewClient(token string, intents types.Intent, sharding *ClientSharding) *Client {
	return &Client{
		Token:             &token,
		APIVersion:        functions.PointerTo(types.APIVersion10),
		Logger:            util.NewLogger(),
		Intents:           &intents,
		UnavailableGuilds: make(map[types.Snowflake]struct{}),
		Sharding:          sharding,
	}
}

func (d *Client) initializeGatewayConnection() (*types.BotRegisterResponse, error) {
	do, err := http.DefaultClient.Do(&http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "https",
			Host:   "discord.com",
			Path:   types.APIBaseString(*d.APIVersion) + "gateway/bot",
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

	var resp types.BotRegisterResponse
	if err := json.NewDecoder(do.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (d *Client) Login() error {
	gatewayResp, err := d.initializeGatewayConnection()
	if err != nil {
		return err
	}

	d.Logger.Debug().Msgf("Connecting to gateway websocket at %s with %d shards", gatewayResp.Url, gatewayResp.Shards)

	if err := d.connectWebsocket(gatewayResp.Url, false, nil); err != nil {
		return err
	}

	go func() {
		if err := d.listenWebsocket(); err != nil {
			d.Logger.Err(err).Msg("Error listening to websocket")
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, 4000, 4001, 4002, 4003, 4005, 4007, 4008, 4009) {
				d.Logger.Debug().Msg(" gateway connection closed by , trying to reconnect")
				if err := d.reconnect(true); err != nil {
					d.Logger.Err(err).Msg("Failed to reconnect")
				}
			}

			d.Logger.Debug().Msg(" gateway connection closed by , no reconnecting attempt will be made")

			return
		}
	}()

	<-d.Websocket.Ready

	d.Logger.Info().Msg("Successfully connected to  gateway")

	return nil
}

func (d *Client) OnEvent(
	eventName types.EventType,
	handler EventHandler,
) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.Events == nil {
		d.Events = make(map[types.EventType][]EventHandler)
	}

	d.Events[eventName] = append(d.Events[eventName], handler)
}

func (d *Client) OnGuildCreate(
	handler func(*Client, *types.GuildCreateEvent),
) {
	d.OnEvent(types.EventGuildCreate, func(
		session *Client,
		event types.Event,
	) {
		if e, ok := event.(*types.GuildCreateEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to GuildCreateEvent: %T", event)
		}
	})
}

func (d *Client) OnMessageCreate(
	handler func(*Client, *types.MessageCreateEvent),
) {
	d.OnEvent(types.EventMessageCreate, func(
		session *Client,
		event types.Event,
	) {
		if e, ok := event.(*types.MessageCreateEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to MessageCreateEvent: %T", event)
		}
	})
}

func (d *Client) OnInteractionCreate(
	handler func(*Client, *types.InteractionCreateEvent),
) {
	d.OnEvent(types.EventInteractionCreate, func(
		session *Client,
		event types.Event,
	) {
		if e, ok := event.(*types.InteractionCreateEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to InteractionCreateEvent: %T", event)
		}
	})
}

func (d *Client) dispatch(event types.Event) {
	handlers := d.Events[event.Event()]
	for _, h := range handlers {
		h(d, event)
	}
}

func (d *Client) internalEventHandler(msg json.RawMessage, event types.EventType) bool {
	switch event {
	case types.EventReady:
		{
			var readyEvent types.ReadyEvent
			if err := json.Unmarshal(msg, &readyEvent); err != nil {
				d.Logger.Err(err).Msg("Failed to unmarshal READY event")
			}

			d.Websocket.SessionID = &readyEvent.SessionID
			d.Websocket.ReconnectURL = &readyEvent.ResumeGatewayURL
			d.User = &readyEvent.User

			if readyEvent.Shard != nil {
				d.Logger.Debug().Msgf("Connected to shard %d of %d", readyEvent.Shard[0]+1, readyEvent.Shard[1])
			}

			for _, guild := range readyEvent.Guilds {
				if !guild.Guild.IsAvailable() {
					d.addUnavailableGuild(guild.Guild.GetID())
				}
			}

			close(d.Websocket.Ready)

			return true
		}
	case types.EventGuildCreate:
		{
			var guildCreateEvent types.GuildCreateEvent
			if err := json.Unmarshal(msg, &guildCreateEvent); err != nil {
				d.Logger.Err(err).Msg("Failed to unmarshal GUILD_CREATE event")
				return false
			}

			if guildCreateEvent.Guild.IsAvailable() && d.IsGuildUnavailable(guildCreateEvent.Guild.GetID()) {
				// is telling us that a guild is available again by firing a GUILD_CREATE event with available = true

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

func (d *Client) addUnavailableGuild(id types.Snowflake) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.UnavailableGuilds[id] = struct{}{}
}

func (d *Client) deleteUnavailableGuild(id types.Snowflake) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.UnavailableGuilds, id)
}

func (d *Client) IsGuildUnavailable(id types.Snowflake) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, exists := d.UnavailableGuilds[id]
	return exists
}

func (d *Client) Shutdown() {
	_ = d.Websocket.Connection.Close()
}
