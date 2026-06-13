package commands

import (
	"encoding/json"
	"fmt"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object
type ApplicationCommand struct {
	ID                       discord.Snowflake                               `json:"id"`
	Type                     discord.ApplicationCommandType                  `json:"type,omitempty"`
	ApplicationID            discord.Snowflake                               `json:"application_id"`
	GuildID                  discord.Snowflake                               `json:"guild_id,omitempty"`
	Name                     string                                          `json:"name"`
	NameLocalizations        map[discord.Locale]string                       `json:"name_localizations,omitempty"`
	Description              string                                          `json:"description"`
	DescriptionLocalizations map[discord.Locale]string                       `json:"description_localizations,omitempty"`
	DefaultMemberPermissions *discord.Permission                             `json:"default_member_permissions"`
	NSFW                     bool                                            `json:"nsfw,omitempty"`
	IntegrationTypes         []discord.InteractionApplicationIntegrationType `json:"integration_types,omitempty"`
	Contexts                 []discord.InteractionContextType                `json:"contexts,omitempty"`
	Version                  discord.Snowflake                               `json:"version"`
	Handler                  CommandHandlerType                              `json:"handler_type,omitempty"`
	Options                  []AnyApplicationCommandOption                   `json:"options,omitempty"`
}

func unmarshalApplicationCommandOption(data []byte) (AnyApplicationCommandOption, error) {
	var meta struct {
		Type discord.ApplicationCommandOptionType `json:"type"`
	}

	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	var opt AnyApplicationCommandOption

	switch meta.Type {
	case discord.ApplicationCommandOptionTypeString:
		opt = &ApplicationCommandOptionString{}
	case discord.ApplicationCommandOptionTypeInteger:
		opt = &ApplicationCommandOptionInteger{}
	case discord.ApplicationCommandOptionTypeNumber:
		opt = &ApplicationCommandOptionNumber{}
	case discord.ApplicationCommandOptionTypeBoolean:
		opt = &ApplicationCommandOptionBoolean{}
	case discord.ApplicationCommandOptionTypeUser:
		opt = &ApplicationCommandOptionUser{}
	case discord.ApplicationCommandOptionTypeChannel:
		opt = &ApplicationCommandOptionChannel{}
	case discord.ApplicationCommandOptionTypeRole:
		opt = &ApplicationCommandOptionRole{}
	case discord.ApplicationCommandOptionTypeMentionable:
		opt = &ApplicationCommandOptionMentionable{}
	case discord.ApplicationCommandOptionTypeAttachment:
		opt = &ApplicationCommandOptionAttachment{}
	case discord.ApplicationCommandOptionTypeSubCommand:
		opt = &ApplicationCommandOptionSubCommand{}
	case discord.ApplicationCommandOptionTypeSubCommandGroup:
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
