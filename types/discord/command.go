package discord

import (
	"encoding/json"
	"fmt"
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object
type ApplicationCommand struct {
	// ID uses omitempty due to command registrations.
	ID                       Snowflake                               `json:"id,omitempty"`
	Type                     ApplicationCommandType                  `json:"type,omitempty"`
	ApplicationID            *Snowflake                              `json:"application_id"`
	GuildID                  *Snowflake                              `json:"guild_id,omitempty"`
	Name                     string                                  `json:"name"`
	NameLocalizations        map[Locale]string                       `json:"name_localizations,omitempty"`
	Description              string                                  `json:"description"`
	DescriptionLocalizations map[Locale]string                       `json:"description_localizations,omitempty"`
	DefaultMemberPermissions *Permission                             `json:"default_member_permissions"`
	NSFW                     bool                                    `json:"nsfw,omitempty"`
	IntegrationTypes         []InteractionApplicationIntegrationType `json:"integration_types,omitempty"`
	Contexts                 []InteractionContextType                `json:"contexts,omitempty"`
	Version                  Snowflake                               `json:"version"`
	Handler                  CommandHandlerType                      `json:"handler_type,omitempty"`
	Options                  []AnyApplicationCommandOption           `json:"options,omitempty"`
}

func unmarshalApplicationCommandOption(data []byte) (AnyApplicationCommandOption, error) {
	var meta struct {
		Type ApplicationCommandOptionType `json:"type"`
	}

	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	var opt AnyApplicationCommandOption

	switch meta.Type {
	case ApplicationCommandOptionTypeString:
		opt = &ApplicationCommandOptionString{}
	case ApplicationCommandOptionTypeInteger:
		opt = &ApplicationCommandOptionInteger{}
	case ApplicationCommandOptionTypeNumber:
		opt = &ApplicationCommandOptionNumber{}
	case ApplicationCommandOptionTypeBoolean:
		opt = &ApplicationCommandOptionBoolean{}
	case ApplicationCommandOptionTypeUser:
		opt = &ApplicationCommandOptionUser{}
	case ApplicationCommandOptionTypeChannel:
		opt = &ApplicationCommandOptionChannel{}
	case ApplicationCommandOptionTypeRole:
		opt = &ApplicationCommandOptionRole{}
	case ApplicationCommandOptionTypeMentionable:
		opt = &ApplicationCommandOptionMentionable{}
	case ApplicationCommandOptionTypeAttachment:
		opt = &ApplicationCommandOptionAttachment{}
	case ApplicationCommandOptionTypeSubCommand:
		opt = &ApplicationCommandOptionSubCommand{}
	case ApplicationCommandOptionTypeSubCommandGroup:
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
		a.Options = opts
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
