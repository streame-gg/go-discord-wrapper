package events

import (
	"log/slog"

	"github.com/streame-gg/go-discord-wrapper/connection"
	devents "github.com/streame-gg/go-discord-wrapper/types/events"
)

func init() {
	On(devents.EventReady, func(c *connection.Client, e *devents.ReadyEvent) {
		c.Logger.Info("bot ready",
			slog.String("user", e.User.Username),
			slog.Int("guilds", len(e.Guilds)),
		)
	})
}
