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

	bot.OnInteractionCreate(func(session *connection.DiscordClient, event *types.DiscordInteractionCreateEvent) {
		if event.IsCommand() {
			bot.Logger.Debug().Msgf("Received interaction command %s from %s", event.GetFullCommand(), event.Member.User.DisplayName())

			_, err := event.CreateInteractionResponse(&types.DiscordInteractionResponse{
				Data: &types.DiscordInteractionResponseData{
					Content: fmt.Sprintf("You invoked the command: %s", event.GetFullCommand()),
					Flags:   types.DiscordMessageFlagSuppressEmbeds,
				},
				Type: types.DiscordInteractionCallbackTypeChannelMessageWithSource,
			})

			if err != nil {
				bot.Logger.Error().Msgf("Failed to create interaction response: %v", err)
			}
		}

		if event.IsButton() {
			bot.Logger.Debug().Msgf("Received button interaction with custom ID %s from %s", event.GetCustomID(), event.Member.User.DisplayName())
		}

		if event.IsAnySelectMenu() {
			bot.Logger.Debug().Msgf("Received select menu interaction with custom ID %s from %s", event.GetCustomID(), event.Member.User.DisplayName())
		}

		if event.IsAutocomplete() {
			bot.Logger.Debug().Msgf("Received autocomplete interaction for command %s from %s", event.GetFullCommand(), event.Member.User.DisplayName())
		}

		if event.IsModalSubmit() {
			bot.Logger.Debug().Msgf("Received modal submit interaction with custom ID %s from %s", event.GetCustomID(), event.Member.User.DisplayName())
		}
	})

	if err := bot.Login(); err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	bot.Logger.Info().Msg("Shutting down bot")
	bot.Shutdown()
}
