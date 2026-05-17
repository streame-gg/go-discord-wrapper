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

	cmdData, ok := i.Data.(*responses.InteractionDataApplicationCommand)
	if !ok {
		return ""
	}

	if cmdData.Options == nil || len(*cmdData.Options) == 0 {
		return ""
	}

	for _, option := range *cmdData.Options {
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

	cmdData, ok := i.Data.(*responses.InteractionDataApplicationCommand)
	if !ok {
		return ""
	}

	if cmdData.Options == nil || len(*cmdData.Options) == 0 {
		return ""
	}

	for _, option := range *cmdData.Options {
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

	cmdData, ok := i.Data.(*responses.InteractionDataApplicationCommand)
	if !ok {
		return ""
	}

	fullCommand += cmdData.CommandName

	if sub := i.GetSubCommandGroup(); sub != "" {
		fullCommand += " " + sub
	}

	if sub := i.GetSubCommand(); sub != "" {
		fullCommand += " " + sub
	}

	return fullCommand
}

func (i *Interaction) GetCustomID() string {
	if i.Data == nil {
		return ""
	}

	if comp, ok := i.Data.(*responses.InteractionDataMessageComponent); ok {
		return comp.CustomID
	}

	if modal, ok := i.Data.(*responses.InteractionDataModalSubmit); ok {
		return modal.CustomID
	}

	return ""
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

	// Autocomplete and regular commands both carry type=ChatInput in the data
	// payload, so check the interaction type first to distinguish them.
	switch i.Type {
	case discord.InteractionTypeApplicationCommandAutocomplete:
		var auto responses.InteractionDataAutocomplete
		if err := json.Unmarshal(aux.Data, &auto); err != nil {
			return err
		}
		i.Data = &auto
		return nil
	case discord.InteractionTypeModalSubmit:
		var modal responses.InteractionDataModalSubmit
		if err := json.Unmarshal(aux.Data, &modal); err != nil {
			return err
		}
		i.Data = &modal
		return nil
	}

	switch typeProbe.ComponentType {
	case discord.ComponentTypeButton,
		discord.ComponentTypeStringSelect,
		discord.ComponentTypeUserSelect,
		discord.ComponentTypeRoleSelect,
		discord.ComponentTypeMentionableSelect,
		discord.ComponentTypeChannelSelect:
		var comp responses.InteractionDataMessageComponent
		if err := json.Unmarshal(aux.Data, &comp); err != nil {
			return err
		}
		i.Data = &comp
		return nil
	}

	switch typeProbe.Type {
	case discord.ApplicationCommandTypeChatInput, discord.ApplicationCommandTypeUser, discord.ApplicationCommandTypeMessage:
		var cmd responses.InteractionDataApplicationCommand
		if err := json.Unmarshal(aux.Data, &cmd); err != nil {
			return err
		}
		i.Data = &cmd
		return nil
	}

	return fmt.Errorf("unknown interaction data type: interaction_type=%d data_type=%d", i.Type, typeProbe.Type)
}
