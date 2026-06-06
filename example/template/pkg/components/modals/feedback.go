package modals

import (
	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
)

func init() { Register(feedback{}) }

// feedback handles the modal opened by the /demo "Open form" button. It reads
// the submitted text input with TextValue and acknowledges it ephemerally. Its
// CustomID must match the modal's custom ID set in components/buttons/open_form.go.
type feedback struct{}

func (feedback) CustomID() string { return "demo_feedback" }

func (feedback) Handle(c *connection.Client, ev *events.InteractionCreateEvent) {
	msg := TextValue(&ev.Interaction, "message")
	if _, err := ev.Reply(interactions.ReplyOptions{
		Content: "📝 Thanks for the feedback:\n> " + msg,
		Flags:   discord.MessageFlagEphemeral,
	}); err != nil {
		c.Logger.Error("feedback modal reply failed", "err", err)
	}
}
