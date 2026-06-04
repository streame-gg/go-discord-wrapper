# Managing commands (create, edit, delete)

[COMMANDS.md](COMMANDS.md) covers writing a command and handling its
interaction. This page is about the **lifecycle** of the command registration
itself: creating, listing, editing, deleting, and scoping commands — plus
per-command permissions.

All of these are REST calls on the client (which embeds the REST client). They
need your **application ID**, which the client reads from the gateway READY
payload — so call them **after `Login`**. The convenience wrappers
(`bot.RegisterCommand`, `bot.BulkRegisterCommands`) fill in the application ID
for you; the finer-grained operations live on `bot.RestClient` and take the
application ID explicitly.

## Global vs. guild commands

| | Global | Guild |
|--|--------|-------|
| Scope | every server the bot is in + DMs | one specific guild |
| Propagation | up to ~1 hour | instant |
| Best for | production | development & per-server commands |

> **Develop against a guild.** Register to a single test guild while iterating so
> changes show up immediately, then switch to global for release.

## Creating & overwriting

`BulkRegisterCommands` is the workhorse: it **replaces the entire global command
set** in one request. Anything not in the list is deleted, which makes it
idempotent — call it on every startup with your full set.

```go
_, err := bot.BulkRegisterCommands(ctx, []*commands.ApplicationCommand{
    {Name: "ping", Description: "Pong", Type: discord.ApplicationCommandTypeChatInput},
})
```

To add one command without touching the others, use `RegisterCommand`:

```go
created, err := bot.RegisterCommand(ctx, &commands.ApplicationCommand{
    Name: "new", Description: "A new command", Type: discord.ApplicationCommandTypeChatInput,
})
```

### Guild-scoped (instant) registration

The guild variants live on `bot.RestClient` and take the application and guild
IDs. Get the application ID from any interaction (`ev.ApplicationID`) or from
the bot user.

```go
appID := ev.ApplicationID // or your application's known ID

// Replace this guild's command set (instant — great for dev).
_, err := bot.RestClient.BulkOverwriteGuildApplicationCommands(ctx, appID, guildID,
    []*commands.ApplicationCommand{
        {Name: "ping", Description: "Pong", Type: discord.ApplicationCommandTypeChatInput},
    })

// Or add a single guild command.
_, err = bot.RestClient.CreateGuildApplicationCommand(ctx, appID, guildID, cmd)
```

## Listing, editing, deleting

Each operation comes in a **global** and a **guild** flavour:

| Action | Global | Guild |
|--------|--------|-------|
| List | `ListGlobalApplicationCommands` | `ListGuildApplicationCommands` |
| Get one | `GetGlobalApplicationCommand` | `GetGuildApplicationCommand` |
| Edit | `EditGlobalApplicationCommand` | `EditGuildApplicationCommand` |
| Delete | `DeleteGlobalApplicationCommand` | `DeleteGuildApplicationCommand` |

```go
// Find a command by name, then edit its description.
cmds, _ := bot.RestClient.ListGlobalApplicationCommands(ctx, appID, false)
for _, cmd := range cmds {
    if cmd.Name == "ping" {
        cmd.Description = "Replies with Pong! (updated)"
        _, err := bot.RestClient.EditGlobalApplicationCommand(ctx, appID, cmd.ID, cmd)
        // …
    }
}

// Delete a command outright.
err := bot.RestClient.DeleteGlobalApplicationCommand(ctx, appID, cmdID)
```

`Edit…` takes a full `*commands.ApplicationCommand` and PATCHes it, so the usual
pattern is *fetch, mutate, send back*.

## Per-command permissions

You can restrict who may use a command in a given guild — by role, user, or
channel — with `EditApplicationCommandPermissions`:

```go
_, err := bot.RestClient.EditApplicationCommandPermissions(ctx, appID, guildID, cmdID,
    []discord.ApplicationCommandPermission{
        {ID: adminRoleID, Type: discord.ApplicationCommandPermissionTypeRole, Permission: true},
    })
```

Read them back with `ListGuildApplicationCommandPermissions` (all commands) or
`GetApplicationCommandPermissions` (one command). Editing permissions accepts a
bot token or an OAuth2 bearer token with the
`applications.commands.permissions.update` scope.

For coarse gating without an API call, set `DefaultMemberPermissions` on the
command definition itself (a `*discord.Permission` bitmask) so only members with
those permissions see it.

## See also

- [docs/COMMANDS.md](COMMANDS.md) — defining options and handling interactions.
- [docs/REST.md](REST.md) — using the REST client and reading typed errors.
- [`example/template`](../example/template) — registers its command set on every
  startup with `BulkRegisterCommands`.
