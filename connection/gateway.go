package connection

import (
	"encoding/json"
	"errors"
	"github.com/DatGamet/go-discord-wrapper/types/common"
	"github.com/DatGamet/go-discord-wrapper/types/events"
	"github.com/DatGamet/go-discord-wrapper/util"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type EventHandler func(*Client, events.Event)

type ClientSharding struct {
	TotalShards int
	ShardID     int
}

type Client struct {
	token *string

	APIVersion *common.APIVersion

	Logger *zerolog.Logger

	Intents *common.Intent

	Websocket *Websocket

	Events map[events.EventType][]EventHandler

	mu sync.RWMutex

	UnavailableGuilds map[common.Snowflake]struct{}

	User *common.User

	Sharding *ClientSharding
}

type ClientOption func(*Client)

func WithSharding(sharding *ClientSharding) ClientOption {
	return func(c *Client) {
		c.Sharding = sharding
	}
}

func WithAPIVersion(version common.APIVersion) ClientOption {
	return func(c *Client) {
		c.APIVersion = util.PointerOf(version)
	}
}

func WithLogger(logger *zerolog.Logger) ClientOption {
	return func(c *Client) {
		c.Logger = logger
	}
}

func NewClient(token string, intents common.Intent, options ...ClientOption) *Client {
	c := &Client{
		token:             &token,
		APIVersion:        util.PointerOf(common.APIVersion10),
		Logger:            util.NewLogger(),
		Intents:           &intents,
		UnavailableGuilds: make(map[common.Snowflake]struct{}),
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func (d *Client) initializeGatewayConnection() (*common.BotRegisterResponse, error) {
	do, err := http.DefaultClient.Do(&http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "https",
			Host:   "discord.com",
			Path:   common.APIBaseString(*d.APIVersion) + "gateway/bot",
		},
		Header: http.Header{
			"Authorization": []string{"Bot " + *d.token},
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

	var resp common.BotRegisterResponse
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
		for {
			if err := d.listenWebsocket(); err != nil {
				d.Logger.Err(err).Msg("Error listening to websocket")

				if d.Websocket == nil {
					d.Logger.Debug().Msg("Websocket is nil, stopping listener")
					return
				}

				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					d.Logger.Debug().Msg("Gateway connection closed normally")
					return
				}

				if websocket.IsCloseError(err, 4000, 4001, 4002, 4003, 4005, 4007, 4008, 4009) {
					d.Logger.Debug().Msg("Gateway connection closed by Discord, trying to reconnect")
					if err := d.reconnect(true); err != nil {
						d.Logger.Err(err).Msg("Failed to reconnect")
						return
					}

					continue
				}

				if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) ||
					websocket.IsUnexpectedCloseError(err) {
					d.Logger.Warn().Msg("Abnormal websocket closure, attempting to reconnect")
					if err := d.reconnect(false); err != nil {
						d.Logger.Err(err).Msg("Failed to reconnect after abnormal closure")
						if err := d.reconnect(true); err != nil {
							d.Logger.Err(err).Msg("Failed to reconnect with fresh connection")
							return
						}
					}
					continue
				}

				d.Logger.Warn().Msg("Unexpected error, attempting to reconnect")
				if err := d.reconnect(false); err != nil {
					d.Logger.Warn().Msg("Resume failed, attempting fresh reconnect")
					if err := d.reconnect(true); err != nil {
						d.Logger.Err(err).Msg("Fresh reconnect failed, stopping listener")
						return
					}
				}
				continue
			}

			if d.Websocket != nil {
				d.Logger.Debug().Msg("Restarting websocket listener")
				continue
			}

			return
		}
	}()

	<-d.Websocket.Ready

	d.Logger.Info().Msg("Successfully connected to the Discord gateway")

	return nil
}

func (d *Client) OnEvent(
	eventName events.EventType,
	handler EventHandler,
) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.Events == nil {
		d.Events = make(map[events.EventType][]EventHandler)
	}

	d.Events[eventName] = append(d.Events[eventName], handler)
}

func (d *Client) OnGuildCreate(
	handler func(*Client, *events.GuildCreateEvent),
) {
	d.OnEvent(events.EventGuildCreate, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.GuildCreateEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to GuildCreateEvent: %T", event)
		}
	})
}

func (d *Client) OnMessageCreate(
	handler func(*Client, *events.MessageCreateEvent),
) {
	d.OnEvent(events.EventMessageCreate, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.MessageCreateEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to MessageCreateEvent: %T", event)
		}
	})
}

func (d *Client) OnInteractionCreate(
	handler func(*Client, *events.InteractionCreateEvent),
) {
	d.OnEvent(events.EventInteractionCreate, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.InteractionCreateEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to InteractionCreateEvent: %T", event)
		}
	})
}

func (d *Client) OnReady(
	handler func(*Client, *events.ReadyEvent),
) {
	d.OnEvent(events.EventReady, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.ReadyEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to ReadyEvent: %T", event)
		}
	})
}

func (d *Client) OnGuildDelete(
	handler func(*Client, *events.GuildDeleteEvent),
) {
	d.OnEvent(events.EventGuildDelete, func(
		session *Client,
		event events.Event,
	) {
		if e, ok := event.(*events.GuildDeleteEvent); ok {
			handler(session, e)
		} else {
			d.Logger.Warn().Msgf("Failed to cast event to GuildDeleteEvent: %T", event)
		}
	})
}

func (d *Client) dispatch(event events.Event) {
	handlers := d.Events[event.Event()]
	for _, h := range handlers {
		h(d, event)
	}
}

func (d *Client) internalEventHandler(msg json.RawMessage, event events.EventType) bool {
	switch event {
	case events.EventReady:
		{
			var readyEvent events.ReadyEvent
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
	case events.EventGuildCreate:
		{
			var guildCreateEvent events.GuildCreateEvent
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
	case events.EventGuildDelete:
		{
			var guildDeleteEvent events.GuildDeleteEvent
			if err := json.Unmarshal(msg, &guildDeleteEvent); err != nil {
				d.Logger.Err(err).Msg("Failed to unmarshal GUILD_DELETE event")
				return false
			}

			if guildDeleteEvent.Unavailable != nil && *guildDeleteEvent.Unavailable {
				if !d.IsGuildUnavailable(guildDeleteEvent.ID) {
					return false
				}

				d.addUnavailableGuild(guildDeleteEvent.ID)

				d.Logger.Debug().Msgf("Guild %s became unavailable", guildDeleteEvent.ID)

				return false
			}
		}
	default:
		return true
	}

	return true
}

func (d *Client) addUnavailableGuild(id common.Snowflake) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.UnavailableGuilds[id] = struct{}{}
}

func (d *Client) deleteUnavailableGuild(id common.Snowflake) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.UnavailableGuilds, id)
}

func (d *Client) IsGuildUnavailable(id common.Snowflake) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, exists := d.UnavailableGuilds[id]
	return exists
}

func (d *Client) Shutdown() {
	_ = d.Websocket.Connection.Close()
}
