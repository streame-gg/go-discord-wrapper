# Slash commands & interactions

This library handles the full slash-command lifecycle: registering commands
with Discord, receiving the resulting interactions over the gateway, and
replying to them. A complete runnable program lives in
[`example/commands`](../example/commands/main.go).

## Registering commands

`bot.BulkRegisterCommands` replaces your application's command set in one call.

```go
_, err := bot.BulkRegisterCommands(ctx, []*commands.ApplicationCommand{
    {
        Name:        "ping",
        Description: "Replies with Pong!",
        Type:        discord.ApplicationCommandTypeChatInput,
    },
})
```

> **Global vs. guild commands.** `BulkRegisterCommands` registers *global*
> commands, which can take up to an hour to propagate. While iterating during
> development, register guild-scoped commands (via the REST client) so changes
> appear instantly.

Register **after** `Login`: the client needs the application ID from the
gateway READY payload to build the request.

### Command options

Options are typed structs implementing `commands.AnyApplicationCommandOption`.
Set each option's `Type` to the matching `discord.ApplicationCommandOptionType`.
Pointer fields like `Required` take a pointer — use `options.Ptr` for literals:

```go
&commands.ApplicationCommand{
    Name:        "echo",
    Description: "Repeat back what you say",
    Type:        discord.ApplicationCommandTypeChatInput,
    Options: &[]commands.AnyApplicationCommandOption{
        &commands.ApplicationCommandOptionString{
            Type:        discord.ApplicationCommandOptionTypeString,
            Name:        "text",
            Description: "What to echo",
            Required:    options.Ptr(true),
        },
    },
}
```

There is one option struct per Discord option type:
`ApplicationCommandOptionString`, `...Integer`, `...Number`, `...Boolean`,
`...User`, `...Channel`, `...Role`, `...Mentionable`, `...Attachment`,
plus `...SubCommand` and `...SubCommandGroup` for nesting.

## Receiving interactions

Slash invocations arrive as `EventInteractionCreate`. Register a handler and
branch on the command path:

```go
bot.OnEvent(events.EventInteractionCreate, func(c *connection.Client, ev *events.InteractionCreateEvent) {
    i := &ev.Interaction
    if !ev.IsCommand() {
        return
    }
    switch ev.GetFullCommand() {
    case "ping":
        _, _ = i.Reply(interactions.ReplyOptions{Content: "Pong!"})
    }
})
```

`GetFullCommand` returns the command name joined with any subcommand group and
subcommand, space-separated — so a `/config set` invocation yields
`"config set"`. Related helpers:

| Helper | Returns |
|--------|---------|
| `ev.IsCommand()` | whether this interaction is a chat-input command |
| `i.GetFullCommand()` | `"name [group] [subcommand]"` |
| `i.GetSubCommand()` | the invoked subcommand name, if any |
| `i.GetSubCommandGroup()` | the invoked subcommand-group name, if any |
| `i.GetCustomID()` | the custom ID of a component or modal-submit interaction |

### Reading option values

The interaction has typed getters for each option kind. They return the zero
value when the option is absent and descend into the selected subcommand /
subcommand-group, so they work for nested options too:

```go
text  := i.GetStringOption("text")
count := i.GetIntOption("count")
ratio := i.GetFloatOption("ratio")
flag  := i.GetBoolOption("flag")
```

When you need to distinguish "absent" from "zero", use the generic form, which
reports whether the option was present and held that type:

```go
if v, ok := interactions.OptionValue[int](i, "count"); ok {
    // v is the supplied count
}
```

For `User`, `Channel`, `Role`, and `Mentionable` options the value is the
entity's ID (read it with `GetStringOption`); the full object is available in the
data's `Resolved` (`Resolved.Users[id]`, etc.). To get at `Resolved` — or any
other field of the concrete data union — assert it with `interactions.As`:

```go
if data, ok := interactions.As[*responses.InteractionDataApplicationCommand](i.Data); ok {
    user := data.Resolved.Users[discord.Snowflake(...)]
}
```

## Replying

You have one acknowledgement to make within **3 seconds**. Pick the method that
matches how long your work takes:

| Method | Use when |
|--------|----------|
| `i.Reply(ReplyOptions{…})` | you can answer immediately |
| `i.DeferReply(DeferOptions{…})` then `i.FollowUp(…)` | the work takes longer than 3s |
| `i.DeferAndFollowup(DeferOptions{…})` | defer + get a ready-to-use sender in one call |
| `i.UpdateMessage(…)` | editing the message a component lives on |
| `i.Modal(…)` | opening a modal form |
| `i.Autocomplete(…)` | responding to an autocomplete interaction |

### Ephemeral replies

Set `Flags: discord.MessageFlagEphemeral` on `ReplyOptions`, or `Ephemeral: true`
on `DeferOptions` / `FollowUpOptions`, to make the response visible only to the
user who invoked the command.

### Deferring slow work

`DeferAndFollowup` acknowledges the interaction right away and hands back a
sender for the eventual reply — convenient when the real work runs in a
goroutine:

```go
send, err := i.DeferAndFollowup(interactions.DeferOptions{Ephemeral: true})
if err != nil { return }

go func() {
    result := doSlowWork()
    _, _ = send(context.Background(), interactions.FollowUpOptions{Content: result})
}()
```

See [`example/slash_with_defer`](../example/slash_with_defer/main.go) for defer,
middleware, and typed error handling together.

## Building rich responses

Embeds, buttons, select menus, modals, and other components are assembled with
the [`builder`](https://pkg.go.dev/github.com/streame-gg/go-discord-wrapper/builder)
package and passed via `ReplyOptions.Embeds` / `ReplyOptions.Components`.

## See also

- [docs/EVENTS.md](EVENTS.md) — the event-handler and middleware model.
- [docs/CONFIGURATION.md](CONFIGURATION.md) — retry and rate-limiting for the
  REST calls that command registration makes.
