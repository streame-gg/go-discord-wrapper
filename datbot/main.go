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

	bot := connection.NewDiscordClient(os.Getenv("TOKEN"), types.AllIntentsExceptDirectMessage)

	bot.OnGuildCreate(func(session *connection.DiscordClient, event *types.DiscordGuildCreateEvent) {
		fmt.Println("New guild")
	})

	bot.OnMessageCreate(func(session *connection.DiscordClient, event *types.DiscordMessageCreateEvent) {
		session.Logger.Info().Msgf("Received message: %s", event.Content)
	})

	if err := bot.Login(); err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	bot.Logger.Info().Msg("Shutting down bot")
	if err := bot.Shutdown(); err != nil {
		bot.Logger.Err(err).Msg("Error shutting down bot")
	}
}
