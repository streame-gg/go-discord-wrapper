package responses

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// AutocompleteChoice is a single option shown in an autocomplete dropdown.
// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-choice-structure
type AutocompleteChoice struct {
	Name              string                    `json:"name"`
	NameLocalizations map[discord.Locale]string `json:"name_localizations,omitempty"`
	Value             interface{}               `json:"value"` // string, int, or float64
}

// InteractionResponseDataAutocomplete is the response data for an autocomplete callback.
// https://docs.discord.com/developers/interactions/receiving-and-responding#autocomplete
type InteractionResponseDataAutocomplete struct {
	Choices []AutocompleteChoice `json:"choices"`
}

func (d *InteractionResponseDataAutocomplete) IsInteractionResponseData() bool { return true }

func (d *InteractionResponseDataAutocomplete) MarshalJSON() ([]byte, error) {
	type Alias InteractionResponseDataAutocomplete
	return json.Marshal((*Alias)(d))
}

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-object-application-command-data
type InteractionDataAutocomplete struct {
	ID       discord.Snowflake                                      `json:"id"`
	Name     string                                                 `json:"name"`
	Type     discord.ApplicationCommandType                         `json:"type"`
	GuildID  *discord.Snowflake                                     `json:"guild_id,omitempty"`
	TargetID *discord.Snowflake                                     `json:"target_id,omitempty"`
	Resolved *discord.ResolvedData                                  `json:"resolved,omitempty"`
	Options  []ApplicationCommandInteractionDataOption[interface{}] `json:"options,omitempty"`
}

func (d *InteractionDataAutocomplete) GetType() discord.InteractionDataType {
	return discord.InteractionDataTypeApplicationCommandAutocomplete
}
