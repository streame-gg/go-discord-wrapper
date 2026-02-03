package main

import (
	"context"
	"fmt"
	"go-discord-wrapper/connection"
	"go-discord-wrapper/functions"
	"go-discord-wrapper/types"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	bot := connection.NewClient(
		os.Getenv("TOKEN"),
		types.AllIntentsExceptDirectMessage,
		&connection.ClientSharding{
			TotalShards: 1,
			ShardID:     0,
		},
	)

	bot.OnMessageCreate(func(session *connection.Client, event *types.MessageCreateEvent) {
		session.Logger.Info().Msgf("Received message: %s", event.Content)
	})

	bot.OnReady(func(session *connection.Client, event *types.ReadyEvent) {
		bot.Logger.Info().Msgf("Logged in as %s#%s", event.User.Username, event.User.Discriminator)
	})

	bot.OnInteractionCreate(func(session *connection.Client, event *types.InteractionCreateEvent) {
		if event.GetFullCommand() == "info channel" {
			bot.Logger.Debug().Msgf("Received info channel command from %s", event.Member.User.DisplayName())

			if err := event.ReplyWithModal(&types.Modal{
				Title:    "Modal",
				CustomID: "modal",
				Components: &[]types.LabelComponent{
					{
						Label:       "Input 1",
						Description: "lololol",
						Component: &types.FileUploadComponent{
							CustomID: "input_1",
							Required: functions.PointerTo(false),
						},
					}, {
						Label:       "Input 2",
						Description: "adadadadadad",
						Component: &types.TextInputComponent{
							CustomID: "input_2",
							Style:    types.TextInputStyleParagraph,
							Required: functions.PointerTo(false),
						},
					},
				},
			}); err != nil {
				bot.Logger.Error().Msgf("Failed to create modal interaction response: %v", err)
				return
			}
			return
		}

		if event.IsCommand() {
			bot.Logger.Debug().Msgf("Received interaction command %s from %s", event.GetFullCommand(), event.Member.User.DisplayName())

			_, err := event.Reply(&types.InteractionResponseDataDefault{
				Flags: types.MessageFlagEphemeral | types.MessageFlagIsComponentsV2,
				Components: &[]types.AnyComponent{
					&types.TextDisplayComponent{
						Content: "Hello, " + event.Member.User.DisplayName() + "!",
					},

					&types.SeparatorComponent{
						SeparatorComponentSpacing: types.SeparatorComponentSpacingSmall,
					},

					&types.MediaGalleryComponent{
						Items: &[]types.MediaGalleryItem{
							{
								Media: &types.UnfurledMediaItem{
									URL: "https://i.imgur.com/AfFp7pu.png",
								},
							},
							{
								Media: &types.UnfurledMediaItem{
									URL: "https://i.imgur.com/AfFp7pu.png",
								},
							},
						},
					},

					&types.Container{
						Components: &[]types.AnyContainerComponent{
							&types.TextDisplayComponent{
								Content: "## Hey!",
							},

							&types.Section{
								Components: &[]types.AnySectionComponent{
									&types.TextDisplayComponent{
										Content: "You used the command **" + event.GetFullCommand() + "**",
									},
								},
								Accessory: &types.ButtonComponent{
									Style:    types.ButtonStylePrimary,
									Label:    "Click Me!",
									CustomID: "button_click_me",
								},
							},
						},
					},
				},
			})

			if err != nil {
				bot.Logger.Error().Msgf("Failed to create interaction response: %v", err)
			}
		}

		if event.IsButton() {
			bot.Logger.Debug().Msgf("Received button interaction with custom ID %s from %s", event.GetCustomID(), event.Member.User.DisplayName())

			if event.GetCustomID() == "button_click_me" {
				_, err := event.Reply(&types.InteractionResponseDataDefault{
					Content: "You clicked the button!",
					Flags:   types.MessageFlagEphemeral,
				})

				if err != nil {
					bot.Logger.Error().Msgf("Failed to create button interaction response: %v", err)
				}
			}
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

	res, err := bot.BulkRegisterCommands([]*types.ApplicationCommand{
		{
			Name:        "info",
			Description: "Get information",
			Type:        types.ApplicationCommandTypeChatInput,
			Options: &[]types.AnyApplicationCommandOption{
				&types.ApplicationCommandOptionSubCommand{
					Name:        "channel",
					Description: "Get information about the channel",
					Options: &[]types.AnyApplicationCommandOption{
						&types.ApplicationCommandOptionChannel{
							Name:        "channel",
							Description: "The channel to get information about",
							Required:    functions.PointerTo(true),
						},
					},
				},
			},
		}, {
			Name:        "welcome",
			Description: "Do more stuff!!!!",
			Type:        types.ApplicationCommandTypeChatInput,
		},
	})
	if err != nil {
		fmt.Println(res)
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	bot.Logger.Info().Msg("Shutting down bot")
	bot.Shutdown()
}
