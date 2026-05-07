package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/streame-gg/go-discord-wrapper/cache"
	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/commands"
	"github.com/streame-gg/go-discord-wrapper/types/common"
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/events"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"

	"github.com/joho/godotenv"
)

func main() {
	slog.SetDefault(slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	))

	_ = godotenv.Load()

	c := cache.NewMemoryCache(cache.Options{
		TTL:           30 * time.Minute,
		SweepInterval: 5 * time.Minute,
		Limits: cache.Limits{
			MaxMessages: 10_000,
			MaxUsers:    5_000,
		},
		Messages: cache.MessageOptions{
			MaxPerChannel: 200,
		},
	})

	bot, clientErr := connection.NewClient(
		os.Getenv("TOKEN"),
		common.AllIntentsExceptDirectMessage,
		options.WithSharding(1, 0),
		options.WithCache(c),
	)
	if clientErr != nil {
		slog.Error("Failed to create client", slog.Any("err", clientErr))
		os.Exit(1)
	}

	bot.OnEvent(events.EventChannelCreate, func(client *connection.Client, event *events.ChannelCreateEvent) {
		client.Logger.Info("Channel created", slog.String("name", event.Name), slog.String("id", event.ID.String()))
	})

	bot.OnEvent(events.EventInteractionCreate, func(client *connection.Client, event *events.InteractionCreateEvent) {
		i := &event.Interaction

		if event.GetFullCommand() == "info channel" {
			bot.Logger.Debug("Received info channel command", slog.String("executingUser", event.Member.User.DisplayName()))

			if err := client.ReplyWithModal(context.Background(), i, &components.Modal{
				Title:    "Modal",
				CustomID: "modal",
				Components: &[]components.LabelComponent{
					{
						Label:       "Input 1",
						Description: "lololol",
						Component: &components.FileUploadComponent{
							CustomID: "input_1",
							Required: options.Ptr(false),
						},
					}, {
						Label:       "Input 2",
						Description: "adadadadadad",
						Component: &components.TextInputComponent{
							CustomID: "input_2",
							Style:    components.TextInputStyleParagraph,
							Required: options.Ptr(false),
						},
					}, {
						Label:       "Input 3",
						Description: "qdwdwqwqddqwdqw",
						Component: &components.CheckboxGroupComponent{
							CustomID: "checkbox",
							Options: &[]components.CheckboxGroupComponentOption{
								{Value: "One", Label: "one"},
								{Value: "Two", Label: "two"},
								{Value: "Three", Label: "three"},
							},
							MinValues: options.Ptr(1),
							MaxValues: options.Ptr(3),
							Required:  options.Ptr(true),
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
								{Value: "one", Label: "Eins"},
								{Value: "two", Label: "Zwei"},
								{Value: "three", Label: "Drei"},
							},
						},
					},
				},
			}); err != nil {
				bot.Logger.Error("Failed to create modal interaction response", slog.Any("err", err))
				return
			}
			return
		}

		if event.IsCommand() {
			bot.Logger.Debug("Received interaction command", slog.String("fullCommand", event.GetFullCommand()), slog.String("user", event.Member.User.DisplayName()))

			_, err := client.Reply(context.Background(), i, &responses.InteractionResponseDataDefault{
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
							{Media: &components.UnfurledMediaItem{URL: "https://i.imgur.com/AfFp7pu.png"}},
							{Media: &components.UnfurledMediaItem{URL: "https://i.imgur.com/AfFp7pu.png"}},
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
			}, false)

			if err != nil {
				bot.Logger.Error("Failed to create interaction response", slog.Any("err", err))
			}
		}

		if event.IsButton() {
			bot.Logger.Info("Received button interaction", slog.String("customID", event.GetCustomID()), slog.String("user", event.Member.User.DisplayName()))

			if event.GetCustomID() == "button_click_me" {
				_, err := client.Reply(context.Background(), i, &responses.InteractionResponseDataDefault{
					Content: "You clicked the button!",
					Flags:   common.MessageFlagEphemeral,
				}, false)

				if err != nil {
					bot.Logger.Error("Failed to create button interaction response", slog.Any("err", err))
				}
			}
		}

		if event.IsAnySelectMenu() {
			bot.Logger.Debug("Received any select menu interaction", slog.String("customId", event.GetCustomID()), slog.String("user", event.Member.User.DisplayName()))
		}

		if event.IsAutocomplete() {
			bot.Logger.Debug("Received autocomplete interaction", slog.String("customId", event.GetCustomID()), slog.String("user", event.Member.User.DisplayName()))
		}

		if event.IsModalSubmit() {
			bot.Logger.Debug("Received modal submit interaction", slog.String("customId", event.GetCustomID()), slog.String("user", event.Member.User.DisplayName()))

			if event.ChannelID != nil {
				// Prefer the cache; only hit the REST API on a miss.
				if ch, ok := bot.Cache.Channels().Get(*event.ChannelID); ok {
					bot.Logger.Debug("Channel resolved from cache", slog.String("channelId", ch.ID.String()), slog.String("name", ch.Name))
				} else {
					ch, err := client.GetChannel(context.Background(), *event.ChannelID)
					if err != nil {
						bot.Logger.Error("Failed to get channel", slog.Any("err", err))
					} else {
						bot.Logger.Debug("Channel fetched from REST", slog.String("channelId", ch.ID.String()))
					}
				}
			}
		}
	})

	if err := bot.Login(); err != nil {
		panic(err)
	}

	_, err := bot.BulkRegisterCommands(context.Background(), []*commands.ApplicationCommand{
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
							Required:    options.Ptr(true),
						},
					},
				},
			},
		},
		{
			Name:        "welcome",
			Description: "Do more stuff!!!!",
			Type:        common.ApplicationCommandTypeChatInput,
		},
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	bot.Logger.Info("Shutting down...")
	if err := bot.Shutdown(); err != nil {
		slog.Error("Shutdown failed", slog.Any("err", err))
	}
}
