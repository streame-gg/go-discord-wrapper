// commands demonstrates the full slash-command lifecycle:
//
//  1. Registering commands (with typed options and subcommands) after Login.
//  2. Routing an interaction by its full command path (GetFullCommand).
//  3. Reading option values the user supplied.
//  4. Replying immediately, ephemerally, or after a defer.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/commands"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

func main() {
	bot, err := connection.NewClient(os.Getenv("DISCORD_TOKEN"), discord.IntentGuilds)
	if err != nil {
		slog.Error("failed to create client", "err", err)
		os.Exit(1)
	}

	// ── Route interactions ─────────────────────────────────────────────────────
	//
	// GetFullCommand returns the command name plus any subcommand-group and
	// subcommand, space-separated — e.g. "echo", "config set". Branch on it.
	bot.OnInteractionCreate(func(c *connection.Client, ev *events.InteractionCreateEvent) {
		i := &ev.Interaction
		if !ev.IsCommand() {
			return
		}

		switch ev.GetFullCommand() {
		case "echo":
			// Read the required "text" string option the user typed.
			text := stringOption(i, "text")
			_, _ = i.Reply(interactions.ReplyOptions{
				Content: "🔊 " + text,
			})

		case "secret":
			// Flags.Ephemeral (or DeferOptions.Ephemeral) makes the reply
			// visible only to the invoking user.
			_, _ = i.Reply(interactions.ReplyOptions{
				Content: "Only you can see this.",
				Flags:   discord.MessageFlagEphemeral,
			})

		case "config set":
			key := stringOption(i, "key")
			_, _ = i.Reply(interactions.ReplyOptions{
				Content: "Updated `" + key + "`.",
				Flags:   discord.MessageFlagEphemeral,
			})
		}
	})

	if err := bot.Login(context.Background()); err != nil {
		slog.Error("login failed", "err", err)
		os.Exit(1)
	}

	// ── Register commands ──────────────────────────────────────────────────────
	//
	// BulkRegisterCommands overwrites the application's global command set. Global
	// commands can take up to an hour to propagate; for instant iteration during
	// development, register guild-scoped commands instead.
	if _, err := bot.BulkRegisterCommands(context.Background(), []*commands.ApplicationCommand{
		{
			Name:        "echo",
			Description: "Repeat back what you say",
			Type:        discord.ApplicationCommandTypeChatInput,
			Options: &[]commands.AnyApplicationCommandOption{
				&commands.ApplicationCommandOptionString{
					Type:        discord.ApplicationCommandOptionTypeString,
					Name:        "text",
					Description: "What to echo",
					Required:    options.Ptr(true),
				},
			},
		},
		{
			Name:        "secret",
			Description: "Send an ephemeral reply",
			Type:        discord.ApplicationCommandTypeChatInput,
		},
		{
			Name:        "config",
			Description: "Manage settings",
			Type:        discord.ApplicationCommandTypeChatInput,
			Options: &[]commands.AnyApplicationCommandOption{
				&commands.ApplicationCommandOptionSubCommand{
					Type:        discord.ApplicationCommandOptionTypeSubCommand,
					Name:        "set",
					Description: "Set a configuration value",
					Options: &[]commands.AnyApplicationCommandOption{
						&commands.ApplicationCommandOptionString{
							Type:        discord.ApplicationCommandOptionTypeString,
							Name:        "key",
							Description: "Which setting to change",
							Required:    options.Ptr(true),
						},
					},
				},
			},
		},
	}); err != nil {
		bot.Logger.Error("failed to register commands", "err", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	_ = bot.Shutdown()
}

// stringOption returns the value of a top-level string option by name, walking
// into the first subcommand/subcommand-group if present. Returns "" if absent.
func stringOption(i *interactions.Interaction, name string) string {
	data, ok := i.Data.(*responses.InteractionDataApplicationCommand)
	if !ok || data.Options == nil {
		return ""
	}
	return findString(*data.Options, name)
}

func findString(opts []responses.ApplicationCommandInteractionDataOption[interface{}], name string) string {
	for _, opt := range opts {
		switch opt.Type {
		case discord.ApplicationCommandOptionTypeSubCommand,
			discord.ApplicationCommandOptionTypeSubCommandGroup:
			if v := findString(opt.Options, name); v != "" {
				return v
			}
		default:
			if opt.Name == name && opt.Value != nil {
				if s, ok := (*opt.Value).(string); ok {
					return s
				}
			}
		}
	}
	return ""
}
