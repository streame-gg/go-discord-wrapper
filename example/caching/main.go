// caching demonstrates how to attach a cache to the client and read entities
// back out of it without hitting the REST API.
//
// It covers:
//
//  1. Configuring an in-process memory cache (TTL, sweep, limits, overflow).
//  2. Reading entities from the cache in an event handler.
//  3. Selecting which stores are auto-populated with options.WithCacheStores.
//  4. Working with the Collection values that cache reads return.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

func main() {
	// ── 1. Build a memory cache ────────────────────────────────────────────────
	//
	// Every field is optional; the zero value is a usable unbounded cache.
	mc := cache.NewMemoryCache(cache.Options{
		// Entries expire 30 minutes after they are written…
		TTL:           30 * time.Minute,
		SweepInterval: 5 * time.Minute,

		// …unless they have not been read recently. EvictUnused also drops
		// entries that have been idle longer than UnusedWindow, even if their
		// TTL has not elapsed — useful for keeping the working set hot.
		EvictBehavior: cache.EvictUnused,
		UnusedWindow:  10 * time.Minute,

		// Global ceilings. When a per-category cap is exceeded the oldest
		// entries in that store are evicted to make room.
		Limits: cache.Limits{
			MaxUsers:    5_000,
			MaxMembers:  20_000,
			MaxMessages: 10_000,
		},

		// Message history is a bounded ring buffer per channel.
		Messages: cache.MessageOptions{
			MaxPerChannel: 200,
		},
	})

	// ── 2. Attach the cache to the client ──────────────────────────────────────
	//
	// WithCacheStores narrows auto-population to just the stores you care about.
	// Here we skip presences and messages to keep memory down. Omit
	// WithCacheStores entirely to populate every store (the default).
	bot, err := connection.NewClient(
		os.Getenv("DISCORD_TOKEN"),
		discord.IntentGuilds|discord.IntentGuildMessages|discord.IntentGuildMembers,
		options.WithCache(mc),
		options.WithCacheStores(
			cache.CategoryGuilds|
				cache.CategoryChannels|
				cache.CategoryUsers|
				cache.CategoryMembers|
				cache.CategoryRoles,
		),
	)
	if err != nil {
		slog.Error("failed to create client", "err", err)
		os.Exit(1)
	}

	// ── 3. Read from the cache in a handler ────────────────────────────────────
	bot.OnGuildCreate(func(c *connection.Client, e *events.GuildCreateEvent) {
		guildID := e.Guild.GetID()

		// Channels and roles arrive inside GUILD_CREATE and are cached
		// immediately. Read them back without a REST round-trip.
		channels := c.ChannelsForGuild(guildID)
		roles := c.Cache.Roles().GetByGuild(guildID)

		// Members in the cache are only the subset Discord sent on connect
		// (see docs/CACHE.md). Filter the snapshot Collection in-place.
		members := c.Cache.Members().AllInGuild(guildID)
		bots := members.Filter(func(m *discord.GuildMember) bool {
			return m.User != nil && m.User.Bot != nil && *m.User.Bot
		})

		c.Logger.Info("guild ready",
			slog.String("guild", guildID.String()),
			slog.Int("channels", channels.Len()),
			slog.Int("roles", roles.Len()),
			slog.Int("cached_members", members.Len()),
			slog.Int("cached_bots", bots.Len()),
		)
	})

	// A single-entity lookup returns (value, ok) — no Collection needed.
	bot.OnMessageCreate(func(c *connection.Client, e *events.MessageCreateEvent) {
		if ch, ok := c.Cache.Channels().Get(e.Message.ChannelID); ok {
			c.Logger.Info("message", "channel", ch.Name, "author", e.Message.Author.Username)
		}
	})

	if err := bot.Login(context.Background()); err != nil {
		slog.Error("login failed", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	_ = bot.Shutdown()
}
