package commands

import (
	"github.com/streame-gg/go-discord-wrapper/builder"
	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
)

// Registering from init() is what makes this command auto-load: the file's mere
// presence in the package is enough — no edit to a central list.
func init() { Register(ping{}) }

// ping is the simplest command, plus a button: it replies immediately and
// attaches a "Ping again" button handled in the components package.
type ping struct{}

func (ping) Definition() discord.ApplicationCommand {
	return discord.ApplicationCommand{
		Name:        "ping",
		Description: "Check that the bot is alive",
		Type:        discord.ApplicationCommandTypeChatInput,
	}
}

func (ping) Handle(c *connection.Client, ev *events.InteractionCreateEvent) {
	// The button's CustomID matches components.pingAgain, which handles clicks.
	row := builder.NewActionRow().AddComponents(
		builder.NewButton().
			SetStyle(components.ButtonStyleSecondary).
			SetLabel("Ping again").
			SetCustomID("ping_again").
			Build(),
	).Build()

	if _, err := ev.Reply(interactions.ReplyOptions{
		Content:    "🏓 Pong!",
		Components: []discord.AnyComponent{row},
	}); err != nil {
		c.Logger.Error("ping reply failed", "err", err)
	}
}
