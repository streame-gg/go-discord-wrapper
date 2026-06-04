package selectmenus

import (
	"strings"

	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
)

func init() { Register(color{}) }

// color handles the string select menu sent by /demo. It reads the chosen
// option(s) with SelectedValues and echoes them back ephemerally.
type color struct{}

func (color) CustomID() string { return "demo_color" }

func (color) Handle(c *connection.Client, ev *events.InteractionCreateEvent) {
	picked := SelectedValues(&ev.Interaction)
	if _, err := ev.Reply(interactions.ReplyOptions{
		Content: "🎨 You picked: " + strings.Join(picked, ", "),
		Flags:   discord.MessageFlagEphemeral,
	}); err != nil {
		c.Logger.Error("color select reply failed", "err", err)
	}
}
