# Contributing

## Reporting bugs

Open an issue on GitHub. Include:
- Go version and OS
- A minimal reproducible example
- The full error message or unexpected behavior

## Pull requests

1. Fork the repo and create a branch from `master`.
2. Keep changes focused — one concern per PR.
3. Add or update tests for any changed behavior.
4. Run `go test ./...` and `go vet ./...` before opening the PR.
5. Open the PR against `master`.

## Code style

- Standard `gofmt` formatting (enforced by `go fmt ./...`).
- Exported symbols must have a godoc comment.
- No `fmt.Println` or `log.Print*` in library code — use the `slog.Logger` that is already threaded through the client.
- Errors are returned, not swallowed.
- Avoid adding external dependencies. The only current dependencies are `gorilla/websocket`, `go-redis`, and `mongo-driver`, all of which are isolated to their respective optional packages.

## Testing

- Unit tests go in `*_test.go` files alongside the code they test.
- Integration tests that need a running service (Redis, MongoDB) use [testcontainers-go](https://golang.testcontainers.org/) — see `cache/rediscache` and `cache/mongocache` for examples.
- Gateway behaviour is harder to test without a live token; focus on unit-testable logic (serialization, cache operations, rate limiting).

## Adding a new gateway event

1. Create `types/events/<eventName>.go` with the event struct, an `init()` that calls `RegisterEvent`, and the two interface methods.
2. If the event should update the cache, add a `case events.EventXxx` block in `connection/gateway.go` `internalEventHandler`.
3. Add the new `EventType` constant to `types/events/base.go` if it is not already there.
4. Update the supported events table in `README.md`.

## Adding a new REST endpoint

1. Add the method to the appropriate file in `api/` (or create a new file for a new resource group).
2. Follow the existing pattern: define param/response types at the top of the file, then the method(s) below.
3. Include a single-line GoDoc comment on every exported symbol.
