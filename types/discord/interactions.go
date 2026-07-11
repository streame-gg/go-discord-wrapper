package discord

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-object-interaction-type
type InteractionType int

const (
	InteractionTypePing                           InteractionType = 1
	InteractionTypeApplicationCommand             InteractionType = 2
	InteractionTypeMessageComponent               InteractionType = 3
	InteractionTypeApplicationCommandAutocomplete InteractionType = 4
	InteractionTypeModalSubmit                    InteractionType = 5
)

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-object-interaction-context-types
type InteractionContextType int

const (
	InteractionContextTypeGuild          InteractionContextType = 0
	InteractionContextTypeBotDM          InteractionContextType = 1
	InteractionContextTypePrivateChannel InteractionContextType = 2
)

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-response-object-interaction-callback-type
type InteractionCallbackType int

const (
	InteractionCallbackTypePong                                 InteractionCallbackType = 1
	InteractionCallbackTypeChannelMessageWithSource             InteractionCallbackType = 4
	InteractionCallbackTypeDeferredChannelMessageWithSource     InteractionCallbackType = 5
	InteractionCallbackTypeDeferredUpdateMessage                InteractionCallbackType = 6
	InteractionCallbackTypeUpdateMessage                        InteractionCallbackType = 7
	InteractionCallbackTypeApplicationCommandAutocompleteResult InteractionCallbackType = 8
	InteractionCallbackTypeModal                                InteractionCallbackType = 9
	// InteractionCallbackTypePremiumRequired is deprecated
	InteractionCallbackTypePremiumRequired InteractionCallbackType = 10
	InteractionCallbackTypeLaunchActivity  InteractionCallbackType = 12
)

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-object-interaction-data
type InteractionData interface {
	GetType() InteractionType
}

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-types
type ApplicationCommandType int

const (
	ApplicationCommandTypeChatInput       ApplicationCommandType = 1
	ApplicationCommandTypeUser            ApplicationCommandType = 2
	ApplicationCommandTypeMessage         ApplicationCommandType = 3
	ApplicationCommandTypePrimaryEndpoint ApplicationCommandType = 4
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-type
type ApplicationCommandOptionType int

const (
	ApplicationCommandOptionTypeSubCommand      ApplicationCommandOptionType = 1
	ApplicationCommandOptionTypeSubCommandGroup ApplicationCommandOptionType = 2
	ApplicationCommandOptionTypeString          ApplicationCommandOptionType = 3
	ApplicationCommandOptionTypeInteger         ApplicationCommandOptionType = 4
	ApplicationCommandOptionTypeBoolean         ApplicationCommandOptionType = 5
	ApplicationCommandOptionTypeUser            ApplicationCommandOptionType = 6
	ApplicationCommandOptionTypeChannel         ApplicationCommandOptionType = 7
	ApplicationCommandOptionTypeRole            ApplicationCommandOptionType = 8
	ApplicationCommandOptionTypeMentionable     ApplicationCommandOptionType = 9
	ApplicationCommandOptionTypeNumber          ApplicationCommandOptionType = 10
	ApplicationCommandOptionTypeAttachment      ApplicationCommandOptionType = 11
)
