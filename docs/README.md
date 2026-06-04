# Documentation

Human-readable guides for `go-discord-wrapper`. For the exhaustive API
reference, see
[pkg.go.dev](https://pkg.go.dev/github.com/streame-gg/go-discord-wrapper).

New here? Start with [GETTING_STARTED.md](GETTING_STARTED.md), then copy
[`example/template`](../example/template) as your project skeleton.

## Guides

### Basics
| Guide | What it covers |
|-------|----------------|
| [GETTING_STARTED.md](GETTING_STARTED.md) | From empty folder to a running bot answering a command |
| [CONFIGURATION.md](CONFIGURATION.md) | Intents and every `options.With…` — caching, retries, rate limiting, logging |
| [EVENTS.md](EVENTS.md) | Registering handlers, middleware, concurrency, synthetic events |

### Interactions
| Guide | What it covers |
|-------|----------------|
| [COMMANDS.md](COMMANDS.md) | Registering slash commands, reading options, replying & deferring |
| [COMMAND_MANAGEMENT.md](COMMAND_MANAGEMENT.md) | Create, edit, delete, list, and scope commands; permissions |
| [COMPONENTS.md](COMPONENTS.md) | Buttons, select menus, and Components V2 layouts |
| [MODALS.md](MODALS.md) | Modal forms: build, show, read the submission |

### Content
| Guide | What it covers |
|-------|----------------|
| [MESSAGES.md](MESSAGES.md) | Sending, editing, replying, attachments, allowed mentions |
| [EMBEDS.md](EMBEDS.md) | Building and validating rich embeds |

### Platform
| Guide | What it covers |
|-------|----------------|
| [CACHE.md](CACHE.md) | What is cached, when, and how to work with incomplete member caches |
| [SHARDING.md](SHARDING.md) | `ShardManager`, the local coordinator, and cross-shard messaging |
| [REST.md](REST.md) | Using the REST client standalone and handling typed errors |

## Runnable examples

| Example | Demonstrates |
|---------|--------------|
| [`example/template`](../example/template/) | **Starter project** — auto-loading `commands/` and `events/` folders, cache, router |
| [`example/caching`](../example/caching/main.go) | Configuring a memory cache and reading entities back |
| [`example/commands`](../example/commands/main.go) | Command registration, option parsing, and replies |
| [`example/sharding`](../example/sharding/main.go) | A 4-shard `ShardManager` with a local coordinator |
| [`example/slash_with_defer`](../example/slash_with_defer/main.go) | Middleware, `DeferAndFollowup`, and typed error handling |

Run any example with its token in the environment:

```sh
DISCORD_TOKEN=... go run ./example/template
```
