package interactions

import (
	"encoding/json"
	"fmt"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

func (i *Interaction) GetSubCommand() string {
	if i.Data == nil {
		return ""
	}

	cmdData, ok := (*i.Data).(*responses.InteractionDataApplicationCommand)
	if !ok {
		return ""
	}

	if cmdData.Options == nil || len(cmdData.Options) == 0 {
		return ""
	}

	for _, option := range cmdData.Options {
		if option.Type == discord.ApplicationCommandOptionTypeSubCommand {
			return option.Name
		}

		if option.Type == discord.ApplicationCommandOptionTypeSubCommandGroup {
			if option.Options != nil {
				for _, subOption := range option.Options {
					if subOption.Type == discord.ApplicationCommandOptionTypeSubCommand {
						return subOption.Name
					}
				}
			}
		}
	}

	return ""
}

func (i *Interaction) GetSubCommandGroup() string {
	if i.Data == nil {
		return ""
	}

	cmdData, ok := (*i.Data).(*responses.InteractionDataApplicationCommand)

	if !ok {
		return ""
	}

	if cmdData.Options == nil || len(cmdData.Options) == 0 {
		return ""
	}

	for _, option := range cmdData.Options {
		if option.Type == discord.ApplicationCommandOptionTypeSubCommandGroup {
			return option.Name
		}
	}

	return ""
}

func (i *Interaction) GetFullCommand() (fullCommand string) {
	if i.Data == nil {
		return ""
	}

	cmdData, ok := (*i.Data).(*responses.InteractionDataApplicationCommand)
	if !ok {
		return ""
	}

	fullCommand += cmdData.Name

	if sub := i.GetSubCommandGroup(); sub != "" {
		fullCommand += " " + sub
	}

	if sub := i.GetSubCommand(); sub != "" {
		fullCommand += " " + sub
	}

	return fullCommand
}

// InvokingUser returns the user who triggered the interaction, whether it
// arrived from a guild (carried on Member.User) or a DM (carried on User). It
// returns nil when neither is populated.
func (i *Interaction) InvokingUser() *discord.User {
	if i == nil {
		return nil
	}
	if i.Member != nil && i.Member.User != nil {
		return i.Member.User
	}
	return i.User
}

// UserID returns the ID of the user who triggered the interaction, or a zero
// Snowflake when it cannot be determined.
func (i *Interaction) UserID() discord.Snowflake {
	if u := i.InvokingUser(); u != nil {
		return u.ID
	}
	return 0
}

// InGuild reports whether the interaction was invoked inside a guild (as opposed
// to a DM or group DM).
func (i *Interaction) InGuild() bool {
	return i != nil && i.GuildID != nil && !i.GuildID.IsEmpty()
}

func (i *Interaction) GetCustomID() string {
	if i.Data == nil {
		return ""
	}

	if comp, ok := (*i.Data).(*responses.InteractionDataMessageComponent); ok {
		return comp.CustomID
	}

	if modal, ok := (*i.Data).(*responses.InteractionDataModalSubmit); ok {
		return modal.CustomID
	}

	return ""
}

// As asserts an interaction's data union to a concrete data type, reporting
// whether the assertion held:
//
//	cmd, ok := interactions.As[*responses.InteractionDataApplicationCommand](ev.Data)
//
// It is a thin, generic wrapper over a type assertion that reads at the call
// site and avoids spelling out the concrete type twice.
func As[T discord.InteractionData](data discord.InteractionData) (T, bool) {
	v, ok := data.(T)
	return v, ok
}

// GetOption returns the command option named name, descending into the selected
// subcommand / subcommand-group so options nested under a subcommand are found
// too. The bool reports whether the option was present.
func (i *Interaction) GetOption(name string) (responses.ApplicationCommandInteractionDataOption[interface{}], bool) {
	data, ok := (*i.Data).(*responses.InteractionDataApplicationCommand)
	if !ok || data.Options == nil {
		return responses.ApplicationCommandInteractionDataOption[interface{}]{}, false
	}
	return findOption(data.Options, name)
}

// OptionValue returns the value of the command option named name as T, reporting
// whether such an option was present and held a value of type T. Discord decodes
// option values to concrete Go types: string options to string, integer options
// to int, number options to float64, boolean options to bool, and snowflake
// (user/channel/role/mentionable) options to their ID as a string.
func OptionValue[T any](i *Interaction, name string) (T, bool) {
	var zero T
	opt, ok := i.GetOption(name)
	if !ok || opt.Value == nil {
		return zero, false
	}
	v, ok := (*opt.Value).(T)
	return v, ok
}

// GetStringOption returns the string option named name, or "" if absent. Also
// use it for user/channel/role/mentionable options, whose values are snowflake
// ID strings.
func (i *Interaction) GetStringOption(name string) string {
	v, _ := OptionValue[string](i, name)
	return v
}

// GetIntOption returns the integer option named name, or 0 if absent.
func (i *Interaction) GetIntOption(name string) int {
	v, _ := OptionValue[int](i, name)
	return v
}

// GetFloatOption returns the number option named name, or 0 if absent.
func (i *Interaction) GetFloatOption(name string) float64 {
	v, _ := OptionValue[float64](i, name)
	return v
}

// GetBoolOption returns the boolean option named name, or false if absent.
func (i *Interaction) GetBoolOption(name string) bool {
	v, _ := OptionValue[bool](i, name)
	return v
}

// findOption returns the option named name, descending into nested subcommand
// options. The bool reports whether it was found.
func findOption(
	opts []responses.ApplicationCommandInteractionDataOption[interface{}],
	name string,
) (responses.ApplicationCommandInteractionDataOption[interface{}], bool) {
	for _, opt := range opts {
		if opt.Name == name {
			return opt, true
		}
		if len(opt.Options) > 0 {
			if found, ok := findOption(opt.Options, name); ok {
				return found, true
			}
		}
	}
	return responses.ApplicationCommandInteractionDataOption[interface{}]{}, false
}

func (i *Interaction) UnmarshalJSON(data []byte) error {
	type Alias Interaction
	aux := &struct {
		Data json.RawMessage `json:"data"`
		*Alias
	}{
		Alias: (*Alias)(i),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if aux.Data == nil {
		return nil
	}

	var typeProbe struct {
		Type          discord.ApplicationCommandType `json:"type"`
		ComponentType discord.ComponentType          `json:"component_type"`
	}

	if err := json.Unmarshal(aux.Data, &typeProbe); err != nil {
		return err
	}
	if i.Data == nil {
		i.Data = new(discord.InteractionData)
	}

	// Autocomplete and regular commands both carry type=ChatInput in the data
	// payload, so check the interaction type first to distinguish them.
	switch i.Type {
	case discord.InteractionTypePing:
		return nil
	case discord.InteractionTypeApplicationCommand:
		var auto responses.InteractionDataApplicationCommand
		if err := json.Unmarshal(aux.Data, &auto); err != nil {
			return err
		}
		*i.Data = &auto
		return nil

	case discord.InteractionTypeApplicationCommandAutocomplete:
		var auto responses.InteractionDataAutocomplete
		if err := json.Unmarshal(aux.Data, &auto); err != nil {
			return err
		}
		*i.Data = &auto
		return nil
	case discord.InteractionTypeModalSubmit:
		var modal responses.InteractionDataModalSubmit
		if err := json.Unmarshal(aux.Data, &modal); err != nil {
			return err
		}
		*i.Data = &modal
		return nil
	}

	switch typeProbe.ComponentType {
	case discord.ComponentTypeButton,
		discord.ComponentTypeStringSelect,
		discord.ComponentTypeUserSelect,
		discord.ComponentTypeRoleSelect,
		discord.ComponentTypeMentionableSelect,
		discord.ComponentTypeChannelSelect,
		discord.ComponentTypeTextInput,
		discord.ComponentTypeCheckbox,
		discord.ComponentTypeCheckboxGroup,
		discord.ComponentTypeRadioGroup,
		discord.ComponentTypeFileUpload:
		var comp responses.InteractionDataMessageComponent
		if err := json.Unmarshal(aux.Data, &comp); err != nil {
			return err
		}
		*i.Data = &comp
		return nil
	}

	switch typeProbe.Type {
	case discord.ApplicationCommandTypeChatInput, discord.ApplicationCommandTypeUser, discord.ApplicationCommandTypeMessage:
		var cmd responses.InteractionDataApplicationCommand
		if err := json.Unmarshal(aux.Data, &cmd); err != nil {
			return err
		}
		*i.Data = &cmd
		return nil
	case discord.ApplicationCommandTypePrimaryEndpoint:
		return nil
	}

	return fmt.Errorf("unknown interaction data type: interaction_type=%d data_type=%d", i.Type, typeProbe.Type)
}
