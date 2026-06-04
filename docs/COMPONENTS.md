# Message components

Components are the interactive bits attached to a message: buttons, select
menus, and — with **Components V2** — full layout primitives like containers,
sections, and media galleries. They are assembled with the
[`builder`](https://pkg.go.dev/github.com/streame-gg/go-discord-wrapper/builder)
package and attached via the `Components` field of a reply, follow-up, or
message.

Every builder ends in `.Build()` and returns a concrete component value that
satisfies `discord.AnyComponent`, so they all drop into
`[]discord.AnyComponent{…}`.

## Classic components: action rows, buttons, selects

Buttons and select menus must live inside an **action row** (max 5 buttons per
row, or one select menu per row).

```go
import (
    "github.com/streame-gg/go-discord-wrapper/builder"
    "github.com/streame-gg/go-discord-wrapper/types/components"
    "github.com/streame-gg/go-discord-wrapper/types/discord"
    "github.com/streame-gg/go-discord-wrapper/types/interactions"
)

confirm := builder.NewButton().
    SetStyle(components.ButtonStylePrimary).
    SetLabel("Confirm").
    SetCustomID("confirm:order-42"). // your routing key — see below
    Build()

cancel := builder.NewButton().
    SetStyle(components.ButtonStyleDanger).
    SetLabel("Cancel").
    SetCustomID("cancel:order-42").
    Build()

row := builder.NewActionRow().AddComponents(confirm, cancel).Build()

_, _ = ev.Reply(interactions.ReplyOptions{
    Content:    "Place this order?",
    Components: []discord.AnyComponent{row},
})
```

### Button styles

| Style | Use |
|-------|-----|
| `ButtonStylePrimary` / `Secondary` / `Success` / `Danger` | regular buttons with a `CustomID` |
| `ButtonStyleLink` | opens a URL — set `SetURL`, **not** a custom ID |
| (premium) | set `SetSKUID` to surface a store SKU |

### Select menus

There is a builder per select type: `NewStringSelectMenu`, `NewUserSelectMenu`,
`NewRoleSelectMenu`, `NewChannelSelectMenu`, `NewMentionableSelectMenu`.

```go
menu := builder.NewStringSelectMenu().
    SetCustomID("pick:color").
    SetPlaceholder("Choose a colour").
    SetMinValues(1).
    SetMaxValues(1).
    AddOptions(
        builder.NewSelectOption("Red", "red").SetDescription("Warm").Build(),
        builder.NewSelectOption("Blue", "blue").Build(),
    ).
    Build()

row := builder.NewActionRow().AddComponents(menu).Build()
```

User/role/channel/mentionable selects resolve to IDs instead of arbitrary
values; pre-select entries with `AddDefaultValues`.

## Handling component interactions

When a user clicks a button or picks from a select, you receive an
`EventInteractionCreate` whose data is a *message component* interaction.
Route on the **custom ID** you assigned:

```go
bot.OnEvent(events.EventInteractionCreate, func(c *connection.Client, ev *events.InteractionCreateEvent) {
    id := ev.GetCustomID()
    if id == "" {
        return // not a component/modal interaction
    }

    switch {
    case strings.HasPrefix(id, "confirm:"):
        orderID := strings.TrimPrefix(id, "confirm:")
        _ = ev.UpdateMessage(interactions.UpdateMessageOptions{
            Content:    "✅ Order " + orderID + " confirmed.",
            Components: []discord.AnyComponent{}, // remove the buttons
        })

    case strings.HasPrefix(id, "pick:"):
        if data, ok := ev.Data.(*responses.InteractionDataMessageComponent); ok {
            _, _ = ev.Reply(interactions.ReplyOptions{
                Content: "You picked: " + strings.Join(data.Values, ", "),
                Flags:   discord.MessageFlagEphemeral,
            })
        }
    }
})
```

Key points:

- **Encode routing data in the custom ID** (e.g. `confirm:order-42`). The custom
  ID is the only state Discord echoes back, so it's how you know *which* button
  on *which* message was clicked.
- **Selected values** live on `InteractionDataMessageComponent.Values` (strings;
  IDs for entity selects).
- **Acknowledge within 3 seconds.** Use `i.UpdateMessage` to edit the message the
  component is on, `i.Reply` to send a new message, or `i.DeferReply` /
  `c.DeferUpdateMessage` if you need longer.

## Components V2

Components V2 replaces the `content` + `embeds` model with a composable layout
tree: **containers**, **sections**, **text displays**, **separators**,
**media galleries**, **thumbnails**, and **file displays**. To send a V2
message you **must** set the `MessageFlagIsComponentsV2` flag, and you may not
mix it with the classic `content`/`embeds` fields.

```go
container := builder.NewContainer().
    SetAccentColor(0x5865F2).
    AddComponents(
        builder.NewTextDisplay().SetContent("# Welcome\nThanks for joining!").Build(),
        builder.NewSeparator().SetDivider(true).Build(),
        builder.NewSection().
            AddComponents(builder.NewTextDisplay().SetContent("Read the rules, then say hi.").Build()).
            SetAccessory(builder.NewThumbnail().SetURL("https://example.com/banner.png").Build()).
            Build(),
        builder.NewMediaGallery().
            AddItems(builder.NewMediaGalleryItem("https://example.com/1.png").Build()).
            Build(),
    ).
    Build()

_, _ = ev.Reply(interactions.ReplyOptions{
    Components: []discord.AnyComponent{container},
    Flags:      discord.MessageFlagIsComponentsV2,
})
```

V2 building blocks and their builders:

| Builder | Produces | Notes |
|---------|----------|-------|
| `NewContainer` | a bordered, optionally coloured group | top-level wrapper |
| `NewSection` | text with a right-hand accessory | accessory is a thumbnail or button |
| `NewTextDisplay` | Markdown text | replaces `content` |
| `NewSeparator` | spacing / divider line | `SetDivider`, `SetSpacing` |
| `NewMediaGallery` + `NewMediaGalleryItem` | a grid of images | |
| `NewThumbnail` | a section accessory image | |
| `NewFileDisplay` | an uploaded-file block | pair with a `MessageFile` attachment |

Buttons and select menus still work inside V2 layouts — wrap them in an action
row as usual.

## See also

- [docs/MODALS.md](MODALS.md) — modal forms (text inputs, checkboxes, radios).
- [docs/COMMANDS.md](COMMANDS.md) — replying to interactions.
- [docs/EMBEDS.md](EMBEDS.md) — classic embeds (the pre-V2 rich layout).
