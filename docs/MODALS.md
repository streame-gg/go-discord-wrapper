# Modals

A modal is a pop-up form you show in response to an interaction (a command or a
button click). The user fills it in, and Discord sends you a **modal-submit**
interaction with the values.

The flow is always three steps: **build → show → read the submission.**

## 1. Build the modal

Modals use the Components V2 `Label` wrapper: each input is wrapped in a label
with a title and optional description. Build inputs with `NewTextInput` (and, in
V2, `NewCheckbox` / `NewRadioGroup` / `NewStringSelectMenu`).

```go
import (
    "github.com/streame-gg/go-discord-wrapper/builder"
    "github.com/streame-gg/go-discord-wrapper/types/components"
)

modal := builder.NewModal().
    SetCustomID("feedback").                // routing key for the submission
    SetTitle("Send feedback").
    AddComponents(
        builder.NewLabel().
            SetLabel("Subject").
            SetComponent(
                builder.NewTextInput().
                    SetCustomID("subject").
                    SetStyle(components.TextInputStyleShort).
                    SetRequired(true).
                    SetMaxLength(100).
                    Build(),
            ).
            Build(),
        builder.NewLabel().
            SetLabel("Details").
            SetDescription("What happened? What did you expect?").
            SetComponent(
                builder.NewTextInput().
                    SetCustomID("details").
                    SetStyle(components.TextInputStyleParagraph).
                    Build(),
            ).
            Build(),
    ).
    Build()
```

`TextInputStyleShort` is a single line; `TextInputStyleParagraph` is multi-line.

## 2. Show it

Respond to the triggering interaction with the modal. **This must be the first
response** — you cannot show a modal after deferring or replying.

```go
bot.OnEvent(events.EventInteractionCreate, func(c *connection.Client, ev *events.InteractionCreateEvent) {
    if ev.IsCommand() && ev.GetFullCommand() == "feedback" {
        if err := ev.Modal(interactions.ModalOptions{Modal: modal}); err != nil {
            c.Logger.Error("show modal failed", "err", err)
        }
    }
})
```

## 3. Read the submission

When the user submits, you get another `EventInteractionCreate` — this time a
modal submit, identifiable by its custom ID. The submitted values are nested:
each label component holds the input, whose response carries the text value.

```go
bot.OnEvent(events.EventInteractionCreate, func(c *connection.Client, ev *events.InteractionCreateEvent) {
    data, ok := ev.Data.(*responses.InteractionDataModalSubmit)
    if !ok || data.CustomID != "feedback" {
        return
    }

    values := modalValues(data) // map[customID]value — helper below

    _, _ = ev.Reply(interactions.ReplyOptions{
        Content: "Thanks! Subject: " + values["subject"],
        Flags:   discord.MessageFlagEphemeral,
    })
})

// modalValues flattens a modal submission into customID -> entered text.
func modalValues(data *responses.InteractionDataModalSubmit) map[string]string {
    out := map[string]string{}
    if data.Components == nil {
        return out
    }
    for _, label := range *data.Components {
        if ti, ok := label.Component.(*components.TextInputComponentInteractionResponse); ok {
            out[ti.CustomID] = ti.Value
        }
    }
    return out
}
```

## Tips

- **The modal's custom ID is your router key** — branch on
  `data.CustomID` to know which form came back, the same way you route buttons.
- **You can open a modal from a button**, not just a command: call `ev.Modal(…)`
  in the button's component handler.
- **After a submit you still owe a response** — reply, defer, or update a message
  within 3 seconds, exactly like any other interaction.

## See also

- [docs/COMPONENTS.md](COMPONENTS.md) — buttons and selects that can open modals.
- [docs/COMMANDS.md](COMMANDS.md) — the interaction lifecycle and reply methods.
