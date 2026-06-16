package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"

	"github.com/streame-gg/go-discord-wrapper/example/template/pkg/components/buttons"
	"github.com/streame-gg/go-discord-wrapper/example/template/pkg/components/modals"
	"github.com/streame-gg/go-discord-wrapper/example/template/pkg/components/selectmenus"
	"github.com/streame-gg/go-discord-wrapper/example/template/pkg/config"
	tevents "github.com/streame-gg/go-discord-wrapper/example/template/pkg/events"
)

func init() { Register(dev{}) }

// dev is an owner-only command for live-reloading parts of the bot without a
// restart. The interesting reload is "commands", which re-pushes the slash
// command set to Discord; "events" and "components" rebuild their in-memory
// registries (the seam a real bot extends — see each package's Reload doc).
//
// Access is gated on config.Current.OwnerID; with no owner configured the
// command refuses to run for anyone.
type dev struct{}

func (dev) Definition() *discord.ApplicationCommand {
	return &discord.ApplicationCommand{
		Name:        "dev",
		Description: "Owner-only developer tools",
		Type:        discord.ApplicationCommandTypeChatInput,
		Options: []discord.AnyApplicationCommandOption{
			&discord.ApplicationCommandOptionSubCommand{
				Type:        discord.ApplicationCommandOptionTypeSubCommand,
				Name:        "reload",
				Description: "Reload commands, events, or components",
				Options: []discord.AnyApplicationCommandOption{
					&discord.ApplicationCommandOptionString{
						Type:        discord.ApplicationCommandOptionTypeString,
						Name:        "choice",
						Description: "What to reload",
						Required:    true,
						Choices: []discord.ApplicationCommandOptionChoice[string]{
							{Name: "Commands", Value: "commands"},
							{Name: "Events", Value: "events"},
							{Name: "Components", Value: "components"},
							{Name: "All", Value: "all"},
						},
					},
				},
			},
		},
	}
}

func (dev) Handle(c *connection.Client, ev *events.InteractionCreateEvent) {
	if !config.Current.IsOwner(ev.UserID()) {
		_, _ = ev.Reply(interactions.ReplyOptions{
			Content: "⛔ This command is owner-only.",
			Flags:   discord.MessageFlagEphemeral,
		})
		return
	}

	if ev.GetSubCommand() != "reload" {
		return
	}

	choice := ev.GetStringOption("choice")
	msg, err := reload(c, &ev.Interaction, choice)
	if err != nil {
		msg = "⚠️ " + err.Error()
	}
	_, _ = ev.Reply(interactions.ReplyOptions{Content: msg, Flags: discord.MessageFlagEphemeral})
}

// reload performs the requested reload and returns a status line to show the
// user. It returns an error when re-registering commands with Discord fails or a
// component registry has a duplicate custom ID, so /dev reload reports the
// problem instead of crashing the bot.
func reload(c *connection.Client, i *interactions.Interaction, choice string) (string, error) {
	switch choice {
	case "events":
		tevents.Reload()
		return "🔁 Reloaded event handlers.", nil
	case "components":
		if err := reloadComponents(); err != nil {
			return "", fmt.Errorf("component reload failed: %w", err)
		}
		return "🔁 Reloaded button, select-menu, and modal handlers.", nil
	case "commands":
		n, err := syncCommands(c, i.ApplicationID)
		if err != nil {
			return "", fmt.Errorf("command sync failed: %w", err)
		}
		return fmt.Sprintf("🔁 Re-registered %d commands.", n), nil
	case "all":
		tevents.Reload()
		if err := reloadComponents(); err != nil {
			return "", fmt.Errorf("component reload failed: %w", err)
		}
		n, err := syncCommands(c, i.ApplicationID)
		if err != nil {
			return "", fmt.Errorf("reloaded events and components, but command sync failed: %w", err)
		}
		return fmt.Sprintf("🔁 Reloaded events, components, and re-registered %d commands.", n), nil
	default:
		return "", fmt.Errorf("unknown reload choice %q", choice)
	}
}

// reloadComponents rebuilds every component registry: buttons, select menus,
// and modals. It returns the first duplicate-custom-ID error so the caller can
// report it; a failed reload leaves that registry's previous table intact.
func reloadComponents() error {
	for _, reload := range []func() error{buttons.Reload, selectmenus.Reload, modals.Reload} {
		if err := reload(); err != nil {
			return err
		}
	}
	return nil
}

// syncCommands re-registers the command set with Discord, scoped to the
// configured dev guild when set. It shares commands.Sync with startup so the
// two never disagree about where commands live.
func syncCommands(c *connection.Client, appID discord.Snowflake) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return Sync(ctx, c, appID, config.Current.DevGuild)
}
