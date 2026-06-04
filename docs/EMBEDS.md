# Embeds

Embeds are the classic rich-content blocks — a titled, coloured card with
fields, images, an author, and a footer. Build them with `builder.NewEmbed()`
and attach via the `Embeds` field of a reply, message, or follow-up.

```go
import "github.com/streame-gg/go-discord-wrapper/builder"

embed := builder.NewEmbed().
    SetTitle("Release v1.0").
    SetDescription("What's new in this version.").
    SetURL("https://example.com/releases/v1").
    SetColor(0x5865F2).
    SetAuthor("changelog-bot", "", "https://example.com/avatar.png").
    AddFields(
        discord.EmbedFields{Name: "Added", Value: "Sharding manager", Inline: options.Ptr(true)},
        discord.EmbedFields{Name: "Fixed", Value: "Cache eviction", Inline: options.Ptr(true)},
    ).
    SetThumbnail("https://example.com/thumb.png").
    SetFooter("go-discord-wrapper", "").
    SetTimestamp(time.Now()).
    Build()

_, _ = ev.Reply(interactions.ReplyOptions{Embeds: []discord.Embed{embed}})
```

## Builder methods

| Method | Sets |
|--------|------|
| `SetTitle` / `SetDescription` / `SetURL` | the headline, body, and title link |
| `SetColor(int)` | the left bar colour (hex, e.g. `0x5865F2`) |
| `SetAuthor(name, url, iconURL)` | the small line above the title |
| `AddFields(...EmbedFields)` | name/value rows; `Inline` packs them side-by-side |
| `SetImage(url)` / `SetThumbnail(url)` | large image / small corner image |
| `SetFooter(text, iconURL)` | the footer line |
| `SetTimestamp(time.Time)` | the timestamp shown in the footer |

`EmbedFields.Inline` is a `*bool` — use `options.Ptr(true)`.

## Validation

Discord enforces limits (title ≤ 256 chars, description ≤ 4096, ≤ 25 fields,
~6000 chars total, etc.). Two ways to catch violations before sending:

```go
// Build then check separately:
embed := builder.NewEmbed().SetTitle(veryLongTitle).Build()
if err := /* builder */ ; err != nil { /* … */ }

// Or build-and-validate in one step:
embed, err := builder.NewEmbed().
    SetTitle(title).
    SetDescription(body).
    BuildChecked()
if err != nil {
    // err describes which limit was exceeded
}
```

`Build()` never fails (it just returns the embed); `BuildChecked()` returns an
error if the embed violates Discord's documented limits, and `Validate()` runs
the same checks on a builder without producing the embed.

## Multiple embeds

A single message may carry up to 10 embeds — pass several in the slice:

```go
interactions.ReplyOptions{Embeds: []discord.Embed{first, second}}
```

## Embeds vs. Components V2

Embeds are the pre-V2 way to lay out rich content. Components V2 (containers,
sections, media galleries) is the newer, more flexible system — but the two are
**mutually exclusive in one message**: a message with the
`MessageFlagIsComponentsV2` flag cannot also carry embeds. Pick one per message.
See [docs/COMPONENTS.md](COMPONENTS.md).

## See also

- [docs/MESSAGES.md](MESSAGES.md) — sending and editing messages.
- [docs/COMPONENTS.md](COMPONENTS.md) — the Components V2 alternative.
