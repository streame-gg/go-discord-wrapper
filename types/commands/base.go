package commands

import (
	"encoding/json"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

const (
	ApplicationCommandNameRegex = "^[-_'\\p{L}\\p{N}\\p{sc=Deva}\\p{sc=Thai}]{1,32}$"
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-entry-point-command-handler-types
type CommandHandlerType int

const (
	CommandHandlerTypeAppHandler            CommandHandlerType = 1
	CommandHandlerTypeDiscordLaunchActivity CommandHandlerType = 2
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-choice-structure
type ApplicationCommandOptionChoice[T string | int64 | float64] struct {
	Name              string                    `json:"name"`
	NameLocalizations map[discord.Locale]string `json:"name_localizations,omitempty"`
	Value             T                         `json:"value"`
}

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-structure
type AnyApplicationCommandOption interface {
	ApplicationCommandOptionType() discord.ApplicationCommandOptionType
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

func unmarshalOptionSlice(raw []json.RawMessage) ([]AnyApplicationCommandOption, error) {
	opts := make([]AnyApplicationCommandOption, 0)

	for _, r := range raw {
		opt, err := unmarshalApplicationCommandOption(r)
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
	}

	return opts, nil
}
