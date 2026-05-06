# Changelog

All notable changes to this project will be documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
This project uses [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `context.Context` as first parameter on all 208 REST methods — enables per-request timeouts and cancellation.
- Typed API errors: `api.APIError` struct with `HTTPStatus`, `Code`, `Message`, and `Errors` fields. All REST methods now return `*api.APIError` on failure; use `errors.Is(err, api.ErrNotFound)` or `errors.As(err, &apiErr)`.
- Sentinel errors: `api.ErrUnauthorized`, `api.ErrForbidden`, `api.ErrNotFound`, `api.ErrRateLimited`.
- `options.Config.Validate()` — validates sharding range, retry counts, and interval signs. Called automatically by `NewClient` and `NewRestClient`; panics on invalid config.
- Pagination helpers: `FetchAllGuildMembers`, `FetchAllMessages`, `FetchAllGuildBans` — automatically page through cursor-based list endpoints.
- `connection.Client.Close()` — alias for `Shutdown()`, satisfies `io.Closer`.
- CI: `go vet`, `gofmt` check, `golangci-lint`, race detector, coverage report, `govulncheck`.
- `.golangci.yml` linting configuration.

### Fixed
- `GatewayError.Error()` returned a Unicode character for the error code (`string(rune(code))`); now uses `strconv.Itoa`.
- `context.Context` cancellation now correctly aborts in-flight retries and rate-limit sleeps, not just the initial HTTP call.
- `Client.Events` map (now unexported as `events`) was a public field reachable without a mutex; direct external writes could cause a concurrent-map panic.
- Gateway event handlers now run in separate goroutines, consistent with lifecycle event handlers (`OnConnect`, `OnDisconnect`). Slow handlers no longer stall the event loop.

## [0.1.0] — Upcoming

Initial stable release target.

[Unreleased]: https://github.com/streame-gg/go-discord-wrapper/compare/HEAD...HEAD
