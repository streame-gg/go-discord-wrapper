package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DatGamet/go-discord-wrapper/connection"
	"github.com/DatGamet/go-discord-wrapper/types/commands"
	"github.com/DatGamet/go-discord-wrapper/types/common"
	"github.com/DatGamet/go-discord-wrapper/types/components"
	"github.com/DatGamet/go-discord-wrapper/types/events"
	"github.com/DatGamet/go-discord-wrapper/types/interactions/responses"
	"github.com/DatGamet/go-discord-wrapper/util"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	bot := connection.NewClient(
		os.Getenv("TOKEN"),
		common.AllIntentsExceptDirectMessage,
		connection.WithSharding(&connection.ClientSharding{
			TotalShards: 1,
			ShardID:     0,
		}),
	)

	bot.OnMessageCreate(func(session *connection.Client, event *events.MessageCreateEvent) {
		session.Logger.Info().Msgf("Received message: %s", event.Content)
	})

	bot.OnReady(func(session *connection.Client, event *events.ReadyEvent) {
		bot.Logger.Info().Msgf("Logged in as %s#%s", event.User.Username, event.User.Discriminator)
	})

	bot.OnInteractionCreate(func(session *connection.Client, event *events.InteractionCreateEvent) {
		if event.GetFullCommand() == "info channel" {
			bot.Logger.Debug().Msgf("Received info channel command from %s", event.Member.User.DisplayName())

			if err := event.ReplyWithModal(&components.Modal{
				Title:    "Modal",
				CustomID: "modal",
				Components: &[]components.LabelComponent{
					{
						Label:       "Input 1",
						Description: "lololol",
						Component: &components.FileUploadComponent{
							CustomID: "input_1",
							Required: util.PointerOf(false),
						},
					}, {
						Label:       "Input 2",
						Description: "adadadadadad",
						Component: &components.TextInputComponent{
							CustomID: "input_2",
							Style:    components.TextInputStyleParagraph,
							Required: util.PointerOf(false),
						},
					}, {
						Label:       "Input 3",
						Description: "qdwdwqwqddqwdqw",
						Component: &components.CheckboxGroupComponent{
							CustomID: "checkbox",
							Options: &[]components.CheckboxGroupComponentOption{
								{
									Value: "One",
									Label: "one",
								}, {
									Value: "Two",
									Label: "two",
								}, {
									Value: "Three",
									Label: "three",
								},
							},
							MinValues: util.PointerOf(1),
							MaxValues: util.PointerOf(3),
							Required:  util.PointerOf(true),
						},
					},
					{
						Label:       "Input 4",
						Description: "dwqdqwwdqdwqdwqdqwqdw",
						Component: &components.CheckboxComponent{
							CustomID: "lololololol",
						},
					},
					{
						Label:       "Input 5",
						Description: "12973123",
						Component: &components.RadioGroupComponent{
							CustomID: "radiogroup",
							Options: &[]components.RadioGroupComponentOption{
								{
									Value: "one",
									Label: "Eins",
								}, {
									Value: "two",
									Label: "Zwei",
								}, {
									Value: "three",
									Label: "Drei",
								},
							},
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

			_, err := event.Reply(&responses.InteractionResponseDataDefault{
				Flags: common.MessageFlagEphemeral | common.MessageFlagIsComponentsV2,
				Components: &[]common.AnyComponent{
					&components.TextDisplayComponent{
						Content: "Hello, " + event.Member.User.DisplayName() + "!",
					},

					&components.SeparatorComponent{
						SeparatorComponentSpacing: components.SeparatorComponentSpacingSmall,
					},

					&components.MediaGalleryComponent{
						Items: &[]components.MediaGalleryItem{
							{
								Media: &components.UnfurledMediaItem{
									URL: "https://i.imgur.com/AfFp7pu.png",
								},
							},
							{
								Media: &components.UnfurledMediaItem{
									URL: "https://i.imgur.com/AfFp7pu.png",
								},
							},
						},
					},

					&components.Container{
						Components: &[]components.AnyContainerComponent{
							&components.TextDisplayComponent{
								Content: "## Hey!",
							},

							&components.Section{
								Components: &[]components.AnySectionComponent{
									&components.TextDisplayComponent{
										Content: "You used the command **" + event.GetFullCommand() + "**",
									},
								},
								Accessory: &components.ButtonComponent{
									Style:    components.ButtonStylePrimary,
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
				_, err := event.Reply(&responses.InteractionResponseDataDefault{
					Content: "You clicked the button!",
					Flags:   common.MessageFlagEphemeral,
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

	res, err := bot.BulkRegisterCommands([]*commands.ApplicationCommand{
		{
			Name:        "info",
			Description: "Get information",
			Type:        common.ApplicationCommandTypeChatInput,
			Options: &[]commands.AnyApplicationCommandOption{
				&commands.ApplicationCommandOptionSubCommand{
					Name:        "channel",
					Description: "Get information about the channel",
					Options: &[]commands.AnyApplicationCommandOption{
						&commands.ApplicationCommandOptionChannel{
							Name:        "channel",
							Description: "The channel to get information about",
							Required:    util.PointerOf(true),
						},
					},
				},
			},
		}, {
			Name:        "welcome",
			Description: "Do more stuff!!!!",
			Type:        common.ApplicationCommandTypeChatInput,
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
