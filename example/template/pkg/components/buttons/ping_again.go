package buttons

import (
	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/types/events"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
)

func init() { Register(pingAgain{}) }

// pingAgain handles the button attached to the /ping reply. It edits the
// original message in place, demonstrating a full button round-trip: a command
// sends a button, and the handler whose CustomID matches responds. The custom
// ID here must match the one /ping sets on the button.
type pingAgain struct{}

func (pingAgain) CustomID() string { return "ping_again" }

func (pingAgain) Handle(c *connection.Client, ev *events.InteractionCreateEvent) {
	if err := ev.UpdateMessage(interactions.UpdateMessageOptions{
		Content: "🏓 Pong! (again)",
	}); err != nil {
		c.Logger.Error("ping_again update failed", "err", err)
	}
}
