// Package bot wires the client, cache, command registry, event handlers, and
// component handlers together. It is the one place that knows about all the
// moving parts; the commands, events, and components packages stay independent
// and self-registering.
package bot

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	devents "github.com/streame-gg/go-discord-wrapper/types/events"

	"github.com/streame-gg/go-discord-wrapper/example/template/pkg/commands"
	"github.com/streame-gg/go-discord-wrapper/example/template/pkg/components/buttons"
	"github.com/streame-gg/go-discord-wrapper/example/template/pkg/components/modals"
	"github.com/streame-gg/go-discord-wrapper/example/template/pkg/components/selectmenus"
	"github.com/streame-gg/go-discord-wrapper/example/template/pkg/config"
	"github.com/streame-gg/go-discord-wrapper/example/template/pkg/events"
)

// Bot is the running application.
type Bot struct {
	Client *connection.Client

	// regOnce guards initial command registration so it runs on the first READY
	// only, not again on every reconnect.
	regOnce sync.Once
}

// New builds a fully wired bot: cache attached, every event handler registered,
// the component lookup table built, and slash-command and component
// interactions routed to their handlers. It does not connect — call Run.
func New(token string) (*Bot, error) {
	// Bounded cache defaults in the spirit of discord.js: a global TTL with a
	// background sweeper, a per-channel message cap, and entry ceilings so a
	// long-running bot in large guilds doesn't grow without bound. Tune these to
	// your workload, or drop the Limits for an unbounded cache. See docs/CACHE.md.
	mc := cache.NewMemoryCache(cache.Options{
		TTL:           30 * time.Minute,
		SweepInterval: 5 * time.Minute,
		Messages:      cache.MessageOptions{MaxPerChannel: 200},
		Limits: cache.Limits{
			MaxChannels: 10_000,
			MaxUsers:    50_000,
			MaxMembers:  50_000,
			MaxRoles:    10_000,
			MaxMessages: 10_000,
		},
	})

	client, err := connection.NewClient(token,
		discord.IntentGuilds|
			discord.IntentGuildMessages|
			discord.IntentGuildMembers|
			discord.IntentMessageContent,
		options.WithCache(mc),
	)
	if err != nil {
		return nil, err
	}

	b := &Bot{Client: client}

	// Auto-load every event handler from the events package.
	events.Attach(client)

	// Build the component lookup tables (commands build theirs from init()).
	buttons.Reload()
	selectmenus.Reload()
	modals.Reload()

	// Register the command set once the application ID is known. READY carries
	// it, and /dev reload re-runs the same commands.Sync, so startup and reload
	// always agree on where commands live (dev guild vs. global).
	client.OnReady(func(c *connection.Client, e *devents.ReadyEvent) {
		b.regOnce.Do(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if _, err := commands.Sync(ctx, c, e.Application.ID, config.Current.DevGuild); err != nil {
				c.Logger.Error("command registration failed", "err", err)
			}
		})
	})

	// Route every interaction to the matching handler.
	if err := client.OnEvent(devents.EventInteractionCreate, b.route); err != nil {
		return nil, err
	}

	return b, nil
}

// route dispatches an interaction to its handler: slash commands by top-level
// name (so "dev reload" routes to the "dev" command), and message components
// (buttons, select menus, modals) by custom ID.
func (b *Bot) route(c *connection.Client, ev *devents.InteractionCreateEvent) {
	switch {
	case ev.IsCommand():
		fields := strings.Fields(ev.GetFullCommand())
		if len(fields) == 0 {
			return
		}
		if cmd, ok := commands.Lookup(fields[0]); ok {
			cmd.Handle(c, ev)
		}
	case ev.IsButton():
		if h, ok := buttons.Lookup(ev.GetCustomID()); ok {
			h.Handle(c, ev)
		}
	case ev.IsAnySelectMenu():
		if h, ok := selectmenus.Lookup(ev.GetCustomID()); ok {
			h.Handle(c, ev)
		}
	case ev.IsModalSubmit():
		if h, ok := modals.Lookup(ev.GetCustomID()); ok {
			h.Handle(c, ev)
		}
	}
}

// Run connects to the gateway. Command registration happens on READY (see New),
// so this just logs in; block on a signal in main to stay running.
func (b *Bot) Run(ctx context.Context) error {
	return b.Client.Login(ctx)
}
