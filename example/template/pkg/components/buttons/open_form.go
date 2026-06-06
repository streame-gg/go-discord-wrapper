package buttons

import (
	"github.com/streame-gg/go-discord-wrapper/builder"
	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/events"
	"github.com/streame-gg/go-discord-wrapper/types/interactions"
)

func init() { Register(openForm{}) }

// openForm is the button sent by /demo. Clicking it responds with a modal whose
// custom ID ("demo_feedback") is handled in the modals package — showing how a
// button hands off to a modal.
type openForm struct{}

func (openForm) CustomID() string { return "demo_open_form" }

func (openForm) Handle(c *connection.Client, ev *events.InteractionCreateEvent) {
	modal := builder.NewModal().
		SetCustomID("demo_feedback").
		SetTitle("Feedback").
		AddComponents(
			builder.NewLabel().
				SetLabel("Your message").
				SetComponent(
					builder.NewTextInput().
						SetCustomID("message").
						SetStyle(components.TextInputStyleParagraph).
						SetPlaceholder("Tell us what you think…").
						Build(),
				).
				Build(),
		).
		Build()

	if err := ev.Modal(interactions.ModalOptions{Modal: modal}); err != nil {
		c.Logger.Error("open_form modal failed", "err", err)
	}
}
