package types

import (
	"encoding/json"
	"fmt"
)

const (
	ApplicationCommandNameRegex = "^[-_'\\p{L}\\p{N}\\p{sc=Deva}\\p{sc=Thai}]{1,32}$"
)

type CommandHandlerType int

const (
	CommandHandlerTypeAppHandler            CommandHandlerType = 1
	CommandHandlerTypeDiscordLaunchActivity CommandHandlerType = 2
)

type ApplicationCommand struct {
	ID                       Snowflake                               `json:"id"`
	Type                     ApplicationCommandType                  `json:"type"`
	ApplicationID            Snowflake                               `json:"application_id"`
	GuildID                  *Snowflake                              `json:"guild_id,omitempty"`
	Name                     string                                  `json:"name"`
	NameLocalizations        map[Locale]string                       `json:"name_localizations,omitempty"`
	Description              string                                  `json:"description"`
	DescriptionLocalizations map[Locale]string                       `json:"description_localizations,omitempty"`
	DefaultMemberPermissions *string                                 `json:"default_member_permissions,omitempty"`
	NSFW                     *bool                                   `json:"nsfw,omitempty"`
	IntegrationTypes         []InteractionApplicationIntegrationType `json:"integration_types,omitempty"`
	Contexts                 []InteractionContextType                `json:"contexts,omitempty"`
	Version                  Snowflake                               `json:"version"`
	Handler                  CommandHandlerType                      `json:"handler_type,omitempty"`
	Options                  *[]AnyApplicationCommandOption          `json:"options,omitempty"`
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

type ApplicationCommandOptionChoice[T string | int] struct {
	Name              string            `json:"name"`
	NameLocalizations map[Locale]string `json:"name_localizations,omitempty"`
	Value             T                 `json:"value"`
}

type AnyApplicationCommandOption interface {
	ApplicationCommandOptionType() ApplicationCommandOptionType
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

type ApplicationCommandOptionString struct {
	Type                     ApplicationCommandOptionType             `json:"type"`
	Name                     string                                   `json:"name"`
	NameLocalizations        map[Locale]string                        `json:"name_localizations,omitempty"`
	Description              string                                   `json:"description"`
	DescriptionLocalizations map[Locale]string                        `json:"description_localizations,omitempty"`
	Required                 *bool                                    `json:"required,omitempty"`
	Choices                  []ApplicationCommandOptionChoice[string] `json:"choices,omitempty"`
	Autocomplete             *bool                                    `json:"autocomplete,omitempty"`
}

func (o *ApplicationCommandOptionString) ApplicationCommandOptionType() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeString
}

func (o *ApplicationCommandOptionString) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionString
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

func (o *ApplicationCommandOptionString) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionString
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

type ApplicationCommandOptionInteger struct {
	Type                     ApplicationCommandOptionType             `json:"type"`
	Name                     string                                   `json:"name"`
	NameLocalizations        map[Locale]string                        `json:"name_localizations,omitempty"`
	Description              string                                   `json:"description"`
	DescriptionLocalizations map[Locale]string                        `json:"description_localizations,omitempty"`
	Required                 *bool                                    `json:"required,omitempty"`
	Choices                  []ApplicationCommandOptionChoice[string] `json:"choices,omitempty"`
	MinValue                 *int64                                   `json:"min_value,omitempty"`
	MaxValue                 *int64                                   `json:"max_value,omitempty"`
	Autocomplete             *bool                                    `json:"autocomplete,omitempty"`
}

func (o *ApplicationCommandOptionInteger) ApplicationCommandOptionType() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeInteger
}

func (o *ApplicationCommandOptionInteger) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionInteger
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

func (o *ApplicationCommandOptionInteger) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionInteger
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

type ApplicationCommandOptionNumber struct {
	Type                     ApplicationCommandOptionType             `json:"type"`
	Name                     string                                   `json:"name"`
	NameLocalizations        map[Locale]string                        `json:"name_localizations,omitempty"`
	Description              string                                   `json:"description"`
	DescriptionLocalizations map[Locale]string                        `json:"description_localizations,omitempty"`
	Required                 *bool                                    `json:"required,omitempty"`
	Choices                  []ApplicationCommandOptionChoice[string] `json:"choices,omitempty"`
	MinValue                 *int64                                   `json:"min_value,omitempty"`
	MaxValue                 *int64                                   `json:"max_value,omitempty"`
	Autocomplete             *bool                                    `json:"autocomplete,omitempty"`
}

func (o *ApplicationCommandOptionNumber) ApplicationCommandOptionType() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeNumber
}

func (o *ApplicationCommandOptionNumber) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionNumber
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

func (o *ApplicationCommandOptionNumber) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionNumber
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

type ApplicationCommandOptionBoolean struct {
	Type                     ApplicationCommandOptionType `json:"type"`
	Name                     string                       `json:"name"`
	NameLocalizations        map[Locale]string            `json:"name_localizations,omitempty"`
	Description              string                       `json:"description"`
	DescriptionLocalizations map[Locale]string            `json:"description_localizations,omitempty"`
	Required                 *bool                        `json:"required,omitempty"`
}

func (o *ApplicationCommandOptionBoolean) ApplicationCommandOptionType() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeBoolean
}

func (o *ApplicationCommandOptionBoolean) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionBoolean
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

func (o *ApplicationCommandOptionBoolean) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionBoolean
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

type ApplicationCommandOptionUser struct {
	Type                     ApplicationCommandOptionType `json:"type"`
	Name                     string                       `json:"name"`
	NameLocalizations        map[Locale]string            `json:"name_localizations,omitempty"`
	Description              string                       `json:"description"`
	DescriptionLocalizations map[Locale]string            `json:"description_localizations,omitempty"`
	Required                 *bool                        `json:"required,omitempty"`
}

func (o *ApplicationCommandOptionUser) ApplicationCommandOptionType() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeUser
}

func (o *ApplicationCommandOptionUser) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionUser
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

func (o *ApplicationCommandOptionUser) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionUser
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

type ApplicationCommandOptionChannel struct {
	Type                     ApplicationCommandOptionType `json:"type"`
	Name                     string                       `json:"name"`
	NameLocalizations        map[Locale]string            `json:"name_localizations,omitempty"`
	Description              string                       `json:"description"`
	DescriptionLocalizations map[Locale]string            `json:"description_localizations,omitempty"`
	Required                 *bool                        `json:"required,omitempty"`
	ChannelTypes             []ChannelType                `json:"channel_types,omitempty"`
}

func (o *ApplicationCommandOptionChannel) ApplicationCommandOptionType() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeChannel
}

func (o *ApplicationCommandOptionChannel) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionChannel
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

func (o *ApplicationCommandOptionChannel) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionChannel
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

type ApplicationCommandOptionRole struct {
	Type                     ApplicationCommandOptionType `json:"type"`
	Name                     string                       `json:"name"`
	NameLocalizations        map[Locale]string            `json:"name_localizations,omitempty"`
	Description              string                       `json:"description"`
	DescriptionLocalizations map[Locale]string            `json:"description_localizations,omitempty"`
	Required                 *bool                        `json:"required,omitempty"`
}

func (o *ApplicationCommandOptionRole) ApplicationCommandOptionType() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeRole
}

func (o *ApplicationCommandOptionRole) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionRole
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

func (o *ApplicationCommandOptionRole) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionRole
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

type ApplicationCommandOptionMentionable struct {
	Type                     ApplicationCommandOptionType `json:"type"`
	Name                     string                       `json:"name"`
	NameLocalizations        map[Locale]string            `json:"name_localizations,omitempty"`
	Description              string                       `json:"description"`
	DescriptionLocalizations map[Locale]string            `json:"description_localizations,omitempty"`
	Required                 *bool                        `json:"required,omitempty"`
}

func (o *ApplicationCommandOptionMentionable) ApplicationCommandOptionType() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeMentionable
}

func (o *ApplicationCommandOptionMentionable) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionMentionable
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

func (o *ApplicationCommandOptionMentionable) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionMentionable
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

type ApplicationCommandOptionAttachment struct {
	Type                     ApplicationCommandOptionType `json:"type"`
	Name                     string                       `json:"name"`
	NameLocalizations        map[Locale]string            `json:"name_localizations,omitempty"`
	Description              string                       `json:"description"`
	DescriptionLocalizations map[Locale]string            `json:"description_localizations,omitempty"`
	Required                 *bool                        `json:"required,omitempty"`
}

func (o *ApplicationCommandOptionAttachment) ApplicationCommandOptionType() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeAttachment
}

func (o *ApplicationCommandOptionAttachment) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionAttachment
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

func (o *ApplicationCommandOptionAttachment) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionAttachment
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

type ApplicationCommandOptionSubCommandGroup struct {
	Type                     ApplicationCommandOptionType  `json:"type"`
	Name                     string                        `json:"name"`
	NameLocalizations        map[Locale]string             `json:"name_localizations,omitempty"`
	Description              string                        `json:"description"`
	DescriptionLocalizations map[Locale]string             `json:"description_localizations,omitempty"`
	Options                  []AnyApplicationCommandOption `json:"options,omitempty"`
}

func (o *ApplicationCommandOptionSubCommandGroup) ApplicationCommandOptionType() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeSubCommandGroup
}

func (o *ApplicationCommandOptionSubCommandGroup) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionSubCommandGroup
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

func (o *ApplicationCommandOptionSubCommandGroup) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionSubCommandGroup
	raw := &struct {
		*Alias
		Options []json.RawMessage `json:"options,omitempty"`
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Options != nil {
		opts, err := unmarshalOptionSlice(raw.Options)
		if err != nil {
			return err
		}
		o.Options = opts
	}

	return nil
}

type ApplicationCommandOptionSubCommand struct {
	Type                     ApplicationCommandOptionType   `json:"type"`
	Name                     string                         `json:"name"`
	NameLocalizations        map[Locale]string              `json:"name_localizations,omitempty"`
	Description              string                         `json:"description"`
	DescriptionLocalizations map[Locale]string              `json:"description_localizations,omitempty"`
	Options                  *[]AnyApplicationCommandOption `json:"options,omitempty"`
}

func (o *ApplicationCommandOptionSubCommand) ApplicationCommandOptionType() ApplicationCommandOptionType {
	return ApplicationCommandOptionTypeSubCommand
}

func (o *ApplicationCommandOptionSubCommand) MarshalJSON() ([]byte, error) {
	o.Type = o.ApplicationCommandOptionType()
	type Alias ApplicationCommandOptionSubCommand
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	})
}

func (o *ApplicationCommandOptionSubCommand) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandOptionSubCommand
	raw := &struct {
		*Alias
		Options []json.RawMessage `json:"options,omitempty"`
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Options != nil {
		opts, err := unmarshalOptionSlice(raw.Options)
		if err != nil {
			return err
		}
		o.Options = &opts
	}

	return nil
}
