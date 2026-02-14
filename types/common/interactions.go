package common

type InteractionType int

const (
	InteractionTypePing                           InteractionType = 1
	InteractionTypeApplicationCommand             InteractionType = 2
	InteractionTypeMessageComponent               InteractionType = 3
	InteractionTypeApplicationCommandAutocomplete InteractionType = 4
	InteractionTypeModalSubmit                    InteractionType = 5
)

type InteractionApplicationIntegrationType int

const (
	InteractionApplicationIntegrationTypeGuildInstall InteractionApplicationIntegrationType = 0
	InteractionApplicationIntegrationTypeUserInstall  InteractionApplicationIntegrationType = 1
)

type InteractionContextType int

const (
	InteractionContextTypeGuild          InteractionContextType = 0
	InteractionContextTypeBotDM          InteractionContextType = 1
	InteractionContextTypePrivateChannel InteractionContextType = 2
)

type InteractionCallbackType int

const (
	InteractionCallbackTypePong                                 InteractionCallbackType = 1
	InteractionCallbackTypeChannelMessageWithSource             InteractionCallbackType = 4
	InteractionCallbackTypeDeferredChannelMessageWithSource     InteractionCallbackType = 5
	InteractionCallbackTypeDeferredUpdateMessage                InteractionCallbackType = 6
	InteractionCallbackTypeUpdateMessage                        InteractionCallbackType = 7
	InteractionCallbackTypeApplicationCommandAutocompleteResult InteractionCallbackType = 8
	InteractionCallbackTypeModal                                InteractionCallbackType = 9
	InteractionCallbackTypePremiumRequired                      InteractionCallbackType = 10
	InteractionCallbackTypeLaunchActivity                       InteractionCallbackType = 12
)

type InteractionDataType int

const (
	InteractionDataTypePing                           InteractionDataType = 1
	InteractionDataTypeApplicationCommand             InteractionDataType = 2
	InteractionDataTypeMessageComponent               InteractionDataType = 3
	InteractionDataTypeApplicationCommandAutocomplete InteractionDataType = 4
	InteractionDataTypeModalSubmit                    InteractionDataType = 5
)

type InteractionData interface {
	GetType() InteractionDataType
}

type ApplicationCommandType int

const (
	ApplicationCommandTypeChatInput       ApplicationCommandType = 1
	ApplicationCommandTypeUser            ApplicationCommandType = 2
	ApplicationCommandTypeMessage         ApplicationCommandType = 3
	ApplicationCommandTypePrimaryEndpoint ApplicationCommandType = 4
)

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
