package main

import (
	"context"
	"fmt"
	"go-discord-wrapper/connection"
	"go-discord-wrapper/types"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	bot := connection.NewDiscordClient(os.Getenv("TOKEN"), 33281)
	if err := bot.Login(); err != nil {
		panic(err)
	}

	bot.On("MESSAGE_CREATE", func(event types.Payload) {
		msg, err := connection.UnwrapEvent[types.DiscordMessageCreateEvent](event)
		if err != nil {
			bot.Logger.Err(err).Msg("Failed to unwrap MESSAGE_CREATE event")
		}
		fmt.Println(fmt.Sprintf("%s said: %s", msg.Author.DisplayName(), msg.Content))
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	bot.Logger.Info().Msg("Shutting down bot")
	if err := bot.Shutdown(); err != nil {
		bot.Logger.Err(err).Msg("Error shutting down bot")
	}
}
