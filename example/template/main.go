// Command template is a batteries-included starter bot you can copy as the
// skeleton for a real project.
//
// Layout:
//
//	main.go                      entry point — config, start, signal handling
//	pkg/bot/                wiring: client + cache + interaction router
//	pkg/config/             in-code settings (owner ID, dev guild)
//	pkg/commands/           one file per slash command, self-registering
//	pkg/components/         shared Handler/Registry, split by interaction kind:
//	pkg/components/buttons/      button handlers
//	pkg/components/selectmenus/  select-menu handlers
//	pkg/components/modals/       modal-submit handlers
//	pkg/events/             one file per gateway handler, self-registering
//
// Adding a command, component, or event is just dropping a new file in the
// matching folder: each file registers itself from an init(), so there is no
// central list to keep in sync. See pkg/commands/ping.go for the smallest
// example.
//
// The owner-only /dev command can reload commands, events, or components at
// runtime without a restart — see pkg/commands/dev.go and the README.
//
// Run it:
//
//	DISCORD_TOKEN=your-token go run ./example/template
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/streame-gg/go-discord-wrapper/example/template/pkg/bot"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		slog.Error("DISCORD_TOKEN is not set")
		os.Exit(1)
	}

	b, err := bot.New(token)
	if err != nil {
		slog.Error("failed to build bot", "err", err)
		os.Exit(1)
	}

	if err := b.Run(context.Background()); err != nil {
		slog.Error("failed to start bot", "err", err)
		os.Exit(1)
	}

	// Block until interrupted, then shut down cleanly.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	slog.Info("shutting down…")
	if err := b.Client.Shutdown(); err != nil {
		slog.Error("shutdown error", "err", err)
	}
}
