// slash_with_defer demonstrates three patterns:
//
//  1. Middleware — a simple logging middleware applied to all handlers.
//  2. DeferAndFollowup — defer the interaction, do slow work, then reply.
//  3. Typed error handling — branch on api.ErrMissingPermissions.
package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	bot, clientErr := connection.NewClient("EXAMPLE_TOKEN", discord.IntentGuilds|discord.IntentGuildMessages)
	if clientErr != nil {
		slog.Error("Failed to create client", slog.Any("err", clientErr))
		os.Exit(1)
	}

	// ── 1. Logging middleware ─────────────────────────────────────────────────
	//
	// Use registers middleware that wraps every handler registered after this call.
	// Multiple middleware are applied outermost-first (first registered = outermost).
	bot.Use(func(next connection.EventHandler) connection.EventHandler {
		return func(c *connection.Client, ev events.Event) {
			start := time.Now()
			c.Logger.Debug("event received", slog.String("type", string(ev.Event())))
			next(c, ev)
			c.Logger.Debug("event handled", slog.String("type", string(ev.Event())), slog.Duration("elapsed", time.Since(start)))
		}
	})

	// ── 2. Slash command with DeferAndFollowup ────────────────────────────────
	//
	// /slow — simulates a command that needs to do async work before replying.
	// The gateway hydrates i with the client and context automatically, so
	// ctx and client are no longer passed explicitly.
	if err := bot.OnEvent(events.EventInteractionCreate, func(c *connection.Client, ev *events.InteractionCreateEvent) {
		i := &ev.Interaction

		if ev.GetFullCommand() != "slow" {
			return
		}

		// DeferAndFollowup acknowledges the interaction right away and returns
		// a sender function for the eventual reply.
		send, err := i.DeferAndFollowup(interactions.DeferOptions{})
		if err != nil {
			c.Logger.Error("defer failed", slog.Any("err", err))
			return
		}

		go func() {
			// Simulate slow work (database query, external API call, …).
			time.Sleep(2 * time.Second)

			if _, err := send(context.Background(), interactions.FollowUpOptions{
				Content: "Done! Your slow result is here.",
			}); err != nil {
				c.Logger.Error("followup failed", slog.Any("err", err))
			}
		}()
	}); err != nil {
		slog.Error("failed to register slow handler", "err", err)
		os.Exit(1)
	}

	// ── 3. Typed error handling ───────────────────────────────────────────────
	//
	// /ban — shows how to distinguish a permission error from other failures.
	if err := bot.OnEvent(events.EventInteractionCreate, func(c *connection.Client, ev *events.InteractionCreateEvent) {
		i := &ev.Interaction

		if ev.GetFullCommand() != "ban" || ev.GuildID == nil || ev.Member == nil {
			return
		}

		// Immediately acknowledge — ban can take a moment.
		send, err := i.DeferAndFollowup(interactions.DeferOptions{Ephemeral: true})
		if err != nil {
			return
		}

		// Pretend we extracted the target user ID from the command options.
		targetID := discord.Snowflake(123456789012345678)

		go func() {
			err := c.CreateGuildBanRaw(context.Background(), *ev.GuildID, targetID, api.CreateGuildBanParams{})
			if err != nil {
				var msg string
				switch {
				case errors.Is(err, api.ErrMissingPermissions):
					msg = "I don't have permission to ban that member."
				case errors.Is(err, api.ErrUnknownMember):
					msg = "That member isn't in this server."
				default:
					msg = "Something went wrong: " + err.Error()
				}
				_, _ = send(context.Background(), interactions.FollowUpOptions{Content: msg})
				return
			}
			_, _ = send(context.Background(), interactions.FollowUpOptions{Content: "Member banned."})
		}()
	}); err != nil {
		slog.Error("failed to register ban handler", "err", err)
		os.Exit(1)
	}

	if err := bot.Login(context.Background()); err != nil {
		panic(err)
	}

	_, err := bot.BulkRegisterCommands(context.Background(), []*discord.ApplicationCommand{
		{Name: "slow", Description: "Simulate a slow operation with DeferAndFollowup", Type: discord.ApplicationCommandTypeChatInput},
		{Name: "ban", Description: "Ban a member (typed-error example)", Type: discord.ApplicationCommandTypeChatInput},
	})
	if err != nil {
		bot.Logger.Error("failed to register commands", slog.Any("err", err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	bot.Logger.Info("Shutting down…")
	if err := bot.Shutdown(); err != nil {
		slog.Error("Shutdown failed", slog.Any("err", err))
	}
}
