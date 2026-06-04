package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/options"
	dcmd "github.com/streame-gg/go-discord-wrapper/types/commands"
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

func (dev) Definition() *dcmd.ApplicationCommand {
	return &dcmd.ApplicationCommand{
		Name:        "dev",
		Description: "Owner-only developer tools",
		Type:        discord.ApplicationCommandTypeChatInput,
		Options: &[]dcmd.AnyApplicationCommandOption{
			&dcmd.ApplicationCommandOptionSubCommand{
				Type:        discord.ApplicationCommandOptionTypeSubCommand,
				Name:        "reload",
				Description: "Reload commands, events, or components",
				Options: &[]dcmd.AnyApplicationCommandOption{
					&dcmd.ApplicationCommandOptionString{
						Type:        discord.ApplicationCommandOptionTypeString,
						Name:        "choice",
						Description: "What to reload",
						Required:    options.Ptr(true),
						Choices: []dcmd.ApplicationCommandOptionChoice[string]{
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
	if !config.Current.IsOwner(invokingUser(&ev.Interaction)) {
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
// user. It returns an error only when re-registering commands with Discord
// fails; the in-memory reloads cannot fail.
func reload(c *connection.Client, i *interactions.Interaction, choice string) (string, error) {
	switch choice {
	case "events":
		tevents.Reload()
		return "🔁 Reloaded event handlers.", nil
	case "components":
		reloadComponents()
		return "🔁 Reloaded button, select-menu, and modal handlers.", nil
	case "commands":
		n, err := syncCommands(c, i.ApplicationID)
		if err != nil {
			return "", fmt.Errorf("command sync failed: %w", err)
		}
		return fmt.Sprintf("🔁 Re-registered %d commands.", n), nil
	case "all":
		tevents.Reload()
		reloadComponents()
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
// and modals.
func reloadComponents() {
	buttons.Reload()
	selectmenus.Reload()
	modals.Reload()
}

// syncCommands re-registers the command set with Discord, scoped to the
// configured dev guild when set. It shares commands.Sync with startup so the
// two never disagree about where commands live.
func syncCommands(c *connection.Client, appID discord.Snowflake) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return Sync(ctx, c, appID, config.Current.DevGuild)
}

// invokingUser returns the ID of the user who ran the interaction, whether it
// arrived from a guild (Member) or a DM (User).
func invokingUser(i *interactions.Interaction) discord.Snowflake {
	if i.Member != nil && i.Member.User != nil {
		return i.Member.User.ID
	}
	if i.User != nil {
		return i.User.ID
	}
	return 0
}
