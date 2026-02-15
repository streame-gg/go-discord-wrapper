package commands

import (
	"encoding/json"
	"fmt"
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type ApplicationCommand struct {
	ID                       *common.Snowflake                              `json:"id,omitempty"`
	Type                     common.ApplicationCommandType                  `json:"type"`
	ApplicationID            common.Snowflake                               `json:"application_id"`
	GuildID                  *common.Snowflake                              `json:"guild_id,omitempty"`
	Name                     string                                         `json:"name"`
	NameLocalizations        map[common.Locale]string                       `json:"name_localizations,omitempty"`
	Description              string                                         `json:"description"`
	DescriptionLocalizations map[common.Locale]string                       `json:"description_localizations,omitempty"`
	DefaultMemberPermissions *string                                        `json:"default_member_permissions,omitempty"`
	NSFW                     *bool                                          `json:"nsfw,omitempty"`
	IntegrationTypes         []common.InteractionApplicationIntegrationType `json:"integration_types,omitempty"`
	Contexts                 []common.InteractionContextType                `json:"contexts,omitempty"`
	Version                  common.Snowflake                               `json:"version"`
	Handler                  CommandHandlerType                             `json:"handler_type,omitempty"`
	Options                  *[]AnyApplicationCommandOption                 `json:"options,omitempty"`
}

func unmarshalApplicationCommandOption(data []byte) (AnyApplicationCommandOption, error) {
	var meta struct {
		Type common.ApplicationCommandOptionType `json:"type"`
	}

	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	var opt AnyApplicationCommandOption

	switch meta.Type {
	case common.ApplicationCommandOptionTypeString:
		opt = &ApplicationCommandOptionString{}
	case common.ApplicationCommandOptionTypeInteger:
		opt = &ApplicationCommandOptionInteger{}
	case common.ApplicationCommandOptionTypeNumber:
		opt = &ApplicationCommandOptionNumber{}
	case common.ApplicationCommandOptionTypeBoolean:
		opt = &ApplicationCommandOptionBoolean{}
	case common.ApplicationCommandOptionTypeUser:
		opt = &ApplicationCommandOptionUser{}
	case common.ApplicationCommandOptionTypeChannel:
		opt = &ApplicationCommandOptionChannel{}
	case common.ApplicationCommandOptionTypeRole:
		opt = &ApplicationCommandOptionRole{}
	case common.ApplicationCommandOptionTypeMentionable:
		opt = &ApplicationCommandOptionMentionable{}
	case common.ApplicationCommandOptionTypeAttachment:
		opt = &ApplicationCommandOptionAttachment{}
	case common.ApplicationCommandOptionTypeSubCommand:
		opt = &ApplicationCommandOptionSubCommand{}
	case common.ApplicationCommandOptionTypeSubCommandGroup:
		opt = &ApplicationCommandOptionSubCommandGroup{}
	default:
		return nil, fmt.Errorf("unknown ApplicationCommandOptionType: %d", meta.Type)
	}

	if err := json.Unmarshal(data, opt); err != nil {
		return nil, err
	}

	return opt, nil
}

func (a *ApplicationCommand) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommand

	var raw struct {
		*Alias
		Options []json.RawMessage `json:"options,omitempty"`
	}

	raw.Alias = (*Alias)(a)

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Options != nil {
		opts, err := unmarshalOptionSlice(raw.Options)
		if err != nil {
			return err
		}
		a.Options = &opts
	}

	return nil
}

func (a *ApplicationCommand) MarshalJSON() ([]byte, error) {
	type Alias ApplicationCommand
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	})
}
