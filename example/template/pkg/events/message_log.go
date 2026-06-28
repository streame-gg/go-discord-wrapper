package events

import (
	"log/slog"

	"github.com/streame-gg/go-discord-wrapper/connection"
	devents "github.com/streame-gg/go-discord-wrapper/types/events"
)

func init() {
	On(devents.EventMessageCreate, func(c *connection.Client, e *devents.MessageCreateEvent) {
		// Ignore our own and other bots' messages.
		if e.Message.Author.Bot != nil && *e.Message.Author.Bot {
			return
		}
		c.Logger.Info("message",
			slog.String("author", e.Message.Author.Username),
			slog.String("content", e.Message.Content),
		)
	})
}
