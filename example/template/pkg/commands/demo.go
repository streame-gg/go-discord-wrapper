package commands

import (
	"github.com/streame-gg/go-discord-wrapper/builder"
	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
)

func init() { Register(demo{}) }

// demo shows the three component kinds in one reply: a select menu (handled in
// components/selectmenus) and a button that opens a modal (handled in
// components/buttons + components/modals). The custom IDs here must match the
// handlers' CustomID values.
type demo struct{}

func (demo) Definition() discord.ApplicationCommand {
	return discord.ApplicationCommand{
		Name:        "demo",
		Description: "Show a select menu and a modal button",
		Type:        discord.ApplicationCommandTypeChatInput,
	}
}

func (demo) Handle(c *connection.Client, ev *events.InteractionCreateEvent) {
	menu := builder.NewActionRow().AddComponents(
		builder.NewStringSelectMenu().
			SetCustomID("demo_color").
			SetPlaceholder("Pick a colour").
			AddOptions(
				builder.NewSelectOption("Red", "red").Build(),
				builder.NewSelectOption("Green", "green").Build(),
				builder.NewSelectOption("Blue", "blue").Build(),
			).
			Build(),
	).Build()

	button := builder.NewActionRow().AddComponents(
		builder.NewButton().
			SetStyle(components.ButtonStylePrimary).
			SetLabel("Open form").
			SetCustomID("demo_open_form").
			Build(),
	).Build()

	if _, err := ev.Reply(interactions.ReplyOptions{
		Content:    "Try the components below:",
		Components: []discord.AnyComponent{menu, button},
		Flags:      discord.MessageFlagEphemeral,
	}); err != nil {
		c.Logger.Error("demo reply failed", "err", err)
	}
}
