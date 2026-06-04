package commands

import (
	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/options"
	dcmd "github.com/streame-gg/go-discord-wrapper/types/commands"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
)

func init() { Register(echo{}) }

// echo shows a command with a required string option, read back with the
// GetStringOption helper on the interaction.
type echo struct{}

func (echo) Definition() *dcmd.ApplicationCommand {
	return &dcmd.ApplicationCommand{
		Name:        "echo",
		Description: "Repeat back what you say",
		Type:        discord.ApplicationCommandTypeChatInput,
		Options: &[]dcmd.AnyApplicationCommandOption{
			&dcmd.ApplicationCommandOptionString{
				Type:        discord.ApplicationCommandOptionTypeString,
				Name:        "text",
				Description: "What to echo",
				Required:    options.Ptr(true),
			},
		},
	}
}

func (echo) Handle(c *connection.Client, ev *events.InteractionCreateEvent) {
	text := ev.GetStringOption("text")
	_, _ = ev.Reply(interactions.ReplyOptions{Content: "🔊 " + text})
}
