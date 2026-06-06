# The REST client

Every gateway `Client` embeds a `*api.RestClient` (as `bot.RestClient`), so all
REST endpoints are available from the bot. You can also use the REST client
**standalone** — no gateway connection — for scripts, web backends handling
HTTP-based interactions, or tests.

```go
import (
    "github.com/streame-gg/go-discord-wrapper/api"
    "github.com/streame-gg/go-discord-wrapper/options"
)

rest, err := api.NewRestClient(token,
    options.WithRetry(options.RetryOptions{MaxRetries: 3}),
)
```

`NewRestClient` takes the same `options.With…` values that apply to REST: retry,
rate limiting, base URL, HTTP client, API version, timeouts. See
[CONFIGURATION.md](CONFIGURATION.md).

## Calling endpoints

Methods are grouped by resource and follow a consistent shape — they take a
`context.Context` first and return typed values:

```go
guild, err := rest.GetGuild(ctx, guildID, false) // false = without member counts
channels, err := rest.ListGuildChannels(ctx, guildID)
msg, err := rest.CreateMessage(ctx, channelID, api.CreateMessageParams{Content: "hi"})
```

> **Hydration note.** When you call these on `bot` (the gateway client) rather
> than `bot.RestClient` directly, returned entities are *hydrated* — their
> `.Members()`, `.Channels()`, etc. managers work and the result is cached. The
> bare `RestClient` returns plain data. Prefer the `bot.*` wrappers when you have
> a gateway client.

## Error handling

Every REST method returns a typed error on a non-2xx response. There are two
ways to inspect it.

### Sentinels with `errors.Is`

For common cases, compare against the package sentinels:

```go
err := bot.CreateGuildBanRaw(ctx, guildID, userID, api.CreateGuildBanParams{})
switch {
case errors.Is(err, api.ErrMissingPermissions):
    // the bot lacks the Ban Members permission
case errors.Is(err, api.ErrUnknownMember):
    // that user isn't in the guild
case errors.Is(err, api.ErrForbidden):
    // 403 for some other reason
case err != nil:
    // anything else
}
```

Status sentinels: `ErrUnauthorized` (401), `ErrForbidden` (403), `ErrNotFound`
(404), `ErrRateLimited` (429). Discord-code sentinels cover the common JSON error
codes: `ErrMissingPermissions`, `ErrMissingAccess`, `ErrUnknownMember`,
`ErrUnknownChannel`, `ErrUnknownMessage`, `ErrUnknownInteraction`,
`ErrInvalidFormBody`, and more.

### Full detail with `errors.As`

To read the HTTP status, Discord error code, message, and field-level validation
errors, unwrap to `*api.Error`:

```go
var apiErr *api.Error
if errors.As(err, &apiErr) {
    log.Printf("http %d, discord code %d: %s",
        apiErr.HTTPStatus, apiErr.Code, apiErr.Message)
    for field, detail := range apiErr.Errors {
        log.Printf("  %s: %v", field, detail)
    }
}
```

## Pagination

List endpoints that page (guild members, messages, audit log, …) take
before/after/limit params. The pagination helpers iterate pages for you — see
the `api` package's paginated `List…` methods and `pagination.go`.

## Rate limits

The client respects Discord's rate limits **proactively** by default: it reads
the bucket headers and waits before a request would exceed a limit, so you
rarely see a 429. Tune or disable it with `options.WithRateLimiting`, and add a
flat floor between requests with `options.WithMinRequestInterval`. Details in
[CONFIGURATION.md](CONFIGURATION.md#rate-limiting).

## Testing against a mock server

`WithBaseURL` points the client at any URL, so you can drive it against
`net/http/httptest` instead of Discord:

```go
ts := httptest.NewServer(myHandler)
defer ts.Close()

rest, _ := api.NewRestClient("test-token", options.WithBaseURL(ts.URL))
```

## See also

- [docs/CONFIGURATION.md](CONFIGURATION.md) — retry & rate-limit options.
- [docs/MESSAGES.md](MESSAGES.md) — the message endpoints in practice.
- [docs/COMMAND_MANAGEMENT.md](COMMAND_MANAGEMENT.md) — command REST endpoints.
