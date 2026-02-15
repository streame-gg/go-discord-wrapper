package commands

import (
	"encoding/json"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

const (
	ApplicationCommandNameRegex = "^[-_'\\p{L}\\p{N}\\p{sc=Deva}\\p{sc=Thai}]{1,32}$"
)

type CommandHandlerType int

const (
	CommandHandlerTypeAppHandler            CommandHandlerType = 1
	CommandHandlerTypeDiscordLaunchActivity CommandHandlerType = 2
)

type ApplicationCommandOptionChoice[T string | int] struct {
	Name              string                   `json:"name"`
	NameLocalizations map[common.Locale]string `json:"name_localizations,omitempty"`
	Value             T                        `json:"value"`
}

type AnyApplicationCommandOption interface {
	ApplicationCommandOptionType() common.ApplicationCommandOptionType
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
