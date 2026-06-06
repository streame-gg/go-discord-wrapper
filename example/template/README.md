# Template bot

A batteries-included starter you can copy as the skeleton for a real bot. It
shows the structure most projects converge on: a thin `main`, a `bot` package
that wires things together, and **self-registering** command and event folders
so adding behaviour is just dropping in a file.

## Run it

```sh
DISCORD_TOKEN=your-token go run ./example/template
```

The bot connects, registers its commands, logs non-bot messages, and serves
`/ping` (with a button), `/echo`, `/serverinfo`, `/demo` (select menu + modal),
and the owner-only `/dev`.

## In-code config

`pkg/config/config.go` holds settings you edit in source rather than pass via the
environment:

```go
var Current = Config{
	OwnerID:  0, // your Discord user ID — who may run /dev
	DevGuild: 0, // a test guild ID for instant command updates (0 = global)
}
```

Set `DevGuild` to a test guild while developing: guild commands update instantly,
whereas global commands can take up to an hour to propagate.

## Reloading at runtime: `/dev reload`

`/dev reload choice:<commands|events|components|all>` lets the configured owner
re-apply parts of the bot without a restart:

- **commands** — re-registers the slash-command set with Discord (instant when
  `DevGuild` is set). This is the genuinely useful one for a running bot.
- **events** / **components** — rebuild their in-memory handler registries.

> Go is compiled, so reload re-applies the *current* binary — it can't pick up
> source edits without a rebuild. The events/components reloads exist as the seam
> a real bot extends (e.g. handlers gated on config that should re-read it); the
> indirection in `events.Attach`/`Reload` is what lets them swap without ever
> double-registering on the client.

## Layout

```
example/template/
├── main.go                      entry point: config, start, signal handling
└── pkg/
    ├── bot/
    │   └── bot.go               wiring: client + cache + interaction router
    ├── config/
    │   └── config.go            in-code settings (owner ID, dev guild)
    ├── commands/
    │   ├── registry.go          Command interface + registry + Sync
    │   ├── ping.go              ── one command per file ──
    │   ├── echo.go
    │   ├── serverinfo.go
    │   ├── demo.go              sends a select menu + a modal button
    │   └── dev.go               owner-only /dev reload
    ├── components/
    │   ├── registry.go          shared Handler interface + reloadable Registry
    │   ├── buttons/             button handlers (ping_again, open_form)
    │   ├── selectmenus/         select-menu handlers (color) + value helpers
    │   └── modals/              modal-submit handlers (feedback) + value helpers
    └── events/
        ├── registry.go          On() + Attach/Reload dispatchers
        ├── ready.go            ── one handler per file ──
        └── message_log.go
```

## How the auto-loader works

Go has no runtime file scanning, so the idiomatic equivalent is **package-level
self-registration**. Every file in `commands/` belongs to the same package and
registers itself from an `init()`:

```go
func init() { Register(ping{}) }
```

When the `bot` package imports `commands`, Go runs every file's `init()`, so the
registry ends up holding all of them. `commands.Definitions()` then feeds
`BulkRegisterCommands`, and the router in `bot.go` looks commands up by name.

Events and components work the same way: each event file calls `events.On(...)`
from its `init()` and `events.Attach(client)` wires them all up; each component
file calls its kind's `Register(...)` (`buttons.Register`, `selectmenus.Register`,
or `modals.Register`) and is matched by custom ID.

**To add a command:** copy `ping.go`, change the definition and handler. That's
it — no list to update.

**To add an event handler:** copy `message_log.go`, swap the `events.On(...)` call.

**To add a component:** drop a file in the folder for its kind —
`components/buttons/`, `components/selectmenus/`, or `components/modals/` — give
it a `CustomID()`, register it from `init()`, and send a component carrying that
same custom ID (see `commands/ping.go` and `commands/demo.go`). Each kind has its
own registry, so the bot router dispatches by interaction type and a button and a
modal can safely share a custom ID.

## Where to go next

- [docs/GETTING_STARTED.md](../../docs/GETTING_STARTED.md) — first-bot walkthrough
- [docs/COMMANDS.md](../../docs/COMMANDS.md) — options, subcommands, replies
- [docs/COMPONENTS.md](../../docs/COMPONENTS.md) — buttons, selects, Components V2
- [docs/MODALS.md](../../docs/MODALS.md) — modal forms
- [docs/CACHE.md](../../docs/CACHE.md) — what the cache holds
