package discord

import (
	"encoding/json"
)

// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-structure
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
	type Alias ApplicationCommandOptionSubCommandGroup
	return json.Marshal(struct {
		Alias
		Type ApplicationCommandOptionType `json:"type"`
	}{
		Alias: Alias(*o),
		Type:  o.ApplicationCommandOptionType(),
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

func NewSubCommandGroupOptionBuilder() *ApplicationCommandOptionSubCommandGroup {
	return &ApplicationCommandOptionSubCommandGroup{
		Type: ApplicationCommandOptionTypeSubCommandGroup,
	}
}

func (o *ApplicationCommandOptionSubCommandGroup) SetName(name string) *ApplicationCommandOptionSubCommandGroup {
	o.Name = name
	return o
}

func (o *ApplicationCommandOptionSubCommandGroup) SetNameLocalizations(localizations map[Locale]string) *ApplicationCommandOptionSubCommandGroup {
	o.NameLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionSubCommandGroup) SetDescription(description string) *ApplicationCommandOptionSubCommandGroup {
	o.Description = description
	return o
}

func (o *ApplicationCommandOptionSubCommandGroup) SetDescriptionLocalizations(localizations map[Locale]string) *ApplicationCommandOptionSubCommandGroup {
	o.DescriptionLocalizations = localizations
	return o
}

func (o *ApplicationCommandOptionSubCommandGroup) AddOptions(opt ...AnyApplicationCommandOption) *ApplicationCommandOptionSubCommandGroup {
	o.Options = append(o.Options, opt...)
	return o
}

func (o *ApplicationCommandOptionSubCommandGroup) AddOption(opt AnyApplicationCommandOption) *ApplicationCommandOptionSubCommandGroup {
	o.Options = append(o.Options, opt)
	return o
}

func (o *ApplicationCommandOptionSubCommandGroup) SetOptions(opt ...AnyApplicationCommandOption) *ApplicationCommandOptionSubCommandGroup {
	o.Options = opt
	return o
}

func (o *ApplicationCommandOptionSubCommandGroup) Build() ApplicationCommandOptionSubCommandGroup {
	return *o
}
