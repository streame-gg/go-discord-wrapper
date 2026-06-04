# Sending & editing messages

Outside of interactions, you send messages to a channel with the REST methods on
the client. The client wraps the REST client, hydrates the returned message
(so its managers work), and caches it for you.

## Sending

```go
import "github.com/streame-gg/go-discord-wrapper/api"

msg, err := bot.SendMessage(ctx, channelID, api.CreateMessageParams{
    Content: "Hello, channel!",
})
```

`CreateMessageParams` carries everything a message can hold:

| Field | Purpose |
|-------|---------|
| `Content` | the text body |
| `Embeds` | up to 10 embeds (see [EMBEDS.md](EMBEDS.md)) |
| `Components` | buttons, selects, V2 layouts (see [COMPONENTS.md](COMPONENTS.md)) |
| `Files` | binary attachments (multipart) |
| `AllowedMentions` | which mentions actually ping |
| `MessageReference` | reply-to another message |
| `Flags` | e.g. `MessageFlagSuppressEmbeds`, `MessageFlagIsComponentsV2` |
| `Poll` | attach a poll |
| `StickerIDs` | send stickers |

### Replying to a message

If you already have a `*discord.Message`, its `Reply` helper sets the reference
for you:

```go
_, err := msg.Reply(ctx, discord.MessageCreateOptions{
    Content: "Replying right here.",
})
```

### Attachments

Attach files via `MessageFile`. Use `Data` for in-memory bytes, or `Reader` for
large files streamed without buffering:

```go
api.CreateMessageParams{
    Content: "Here's the log.",
    Files: []discord.MessageFile{
        {Name: "out.log", ContentType: "text/plain", Data: logBytes},
    },
}
```

When `Files` is set the request is sent as `multipart/form-data` automatically.

### Controlling mentions

By default a message pings everyone it mentions. Restrict that with
`AllowedMentions` — e.g. allow role/user mentions but never `@everyone`:

```go
api.CreateMessageParams{
    Content: "Heads up <@&123>!",
    AllowedMentions: &discord.AllowedMentions{
        // Parse is a pointer-to-slice so the empty case is distinguishable:
        //   nil           → Discord's defaults (mention everything)
        //   &[]{}         → suppress all auto-mentions
        //   &[]{"roles"}  → allow only role mentions
        Parse: &[]discord.AllowedMentionsType{discord.AllowedMentionsTypeRoles},
    },
}
```

## Editing & deleting

```go
content := "Edited content."
edited, err := bot.EditMessageRaw(ctx, channelID, messageID, api.EditMessageParams{
    Content: &content, // pointer fields: nil = leave unchanged
})
```

`EditMessageParams` uses pointers so you can change one field without clobbering
the rest. To **remove** all components, set `Components` to a non-nil empty slice
(`[]discord.AnyComponent{}`); leaving it nil keeps the existing ones.

## Reading messages

```go
msg, err := bot.GetMessage(ctx, channelID, messageID)

// Paginated history (also populates the cache):
msgs, err := bot.ListChannelMessagesRaw(ctx, channelID, api.GetMessagesParams{Limit: 50})
```

If a cache is attached, recently-seen messages are also available without a REST
call via `bot.Cache.Messages().Channel(channelID)` — see [CACHE.md](CACHE.md).

## Interaction responses vs. channel messages

Replying to a slash command or component uses the **interaction** methods
(`i.Reply`, `i.FollowUp`, `i.DeferAndFollowup`), not `SendMessage` — they go
through the interaction token and have their own 3-second / 15-minute timing
rules. See [COMMANDS.md](COMMANDS.md). `SendMessage` is for posting to a channel
on your own initiative.

## See also

- [docs/EMBEDS.md](EMBEDS.md) and [docs/COMPONENTS.md](COMPONENTS.md) — rich content.
- [docs/REST.md](REST.md) — error handling and the REST client.
