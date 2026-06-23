package events

import (
	"log/slog"

	"github.com/streame-gg/go-discord-wrapper/connection"
	devents "github.com/streame-gg/go-discord-wrapper/types/events"
)

func init() {
	On(devents.EventMessageCreate, func(c *connection.Client, e *devents.MessageCreateEvent) {
		// Ignore our own and other bots' messages.
		if e.Author.Bot != nil && *e.Author.Bot {
			return
		}
		c.Logger.Info("message",
			slog.String("author", e.Author.Username),
			slog.String("content", e.Content),
		)
	})
}
