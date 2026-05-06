package interactions

import (
	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/types/common"
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

// Client is the subset of *connection.Client needed by Interaction methods.
// Pass the client received in your event handler — *connection.Client satisfies this interface.
type Client interface {
	Reply(i *Interaction, data *responses.InteractionResponseDataDefault, withResponse bool) (*responses.InteractionCallbackResponse, error)
	DeferReply(i *Interaction, ephemeral bool) error
	DeferUpdateMessage(i *Interaction) error
	UpdateMessage(i *Interaction, data *responses.InteractionResponseDataDefault) error
	ReplyWithModal(i *Interaction, modal *components.Modal) error
	ReplyAutocomplete(i *Interaction, choices []responses.AutocompleteChoice) error
	GetOriginalResponse(i *Interaction) (*common.Message, error)
	EditReply(i *Interaction, params api.EditMessageParams) (*common.Message, error)
	DeleteReply(i *Interaction) error
	CreateFollowup(i *Interaction, params api.CreateMessageParams) (*common.Message, error)
	GetFollowup(i *Interaction, messageID common.Snowflake) (*common.Message, error)
	EditFollowup(i *Interaction, messageID common.Snowflake, params api.EditMessageParams) (*common.Message, error)
	DeleteFollowup(i *Interaction, messageID common.Snowflake) error
}

// Reply sends an immediate message response to the interaction.
func (i *Interaction) Reply(client Client, data *responses.InteractionResponseDataDefault, withResponse bool) (*responses.InteractionCallbackResponse, error) {
	return client.Reply(i, data, withResponse)
}

// DeferReply acknowledges the interaction without an immediate visible response.
// Set ephemeral=true to restrict the eventual follow-up to the invoker only.
func (i *Interaction) DeferReply(client Client, ephemeral bool) error {
	return client.DeferReply(i, ephemeral)
}

// DeferUpdateMessage acknowledges a component interaction without editing the original message.
// Call EditReply afterwards to push the actual update.
func (i *Interaction) DeferUpdateMessage(client Client) error {
	return client.DeferUpdateMessage(i)
}

// UpdateMessage edits the message that triggered a component interaction.
func (i *Interaction) UpdateMessage(client Client, data *responses.InteractionResponseDataDefault) error {
	return client.UpdateMessage(i, data)
}

// Modal responds to the interaction with a modal dialog.
func (i *Interaction) Modal(client Client, modal *components.Modal) error {
	return client.ReplyWithModal(i, modal)
}

// Autocomplete sends autocomplete choices back to Discord.
func (i *Interaction) Autocomplete(client Client, choices []responses.AutocompleteChoice) error {
	return client.ReplyAutocomplete(i, choices)
}

// GetOriginalResponse fetches the original response message for this interaction.
func (i *Interaction) GetOriginalResponse(client Client) (*common.Message, error) {
	return client.GetOriginalResponse(i)
}

// EditReply edits the original interaction response.
func (i *Interaction) EditReply(client Client, params api.EditMessageParams) (*common.Message, error) {
	return client.EditReply(i, params)
}

// DeleteReply deletes the original interaction response.
func (i *Interaction) DeleteReply(client Client) error {
	return client.DeleteReply(i)
}

// FollowUp sends a follow-up message (usable up to 15 minutes after the initial response).
func (i *Interaction) FollowUp(client Client, params api.CreateMessageParams) (*common.Message, error) {
	return client.CreateFollowup(i, params)
}

// GetFollowup fetches a follow-up message by ID.
func (i *Interaction) GetFollowup(client Client, messageID common.Snowflake) (*common.Message, error) {
	return client.GetFollowup(i, messageID)
}

// EditFollowup edits a follow-up message by ID.
func (i *Interaction) EditFollowup(client Client, messageID common.Snowflake, params api.EditMessageParams) (*common.Message, error) {
	return client.EditFollowup(i, messageID, params)
}

// DeleteFollowup deletes a follow-up message by ID.
func (i *Interaction) DeleteFollowup(client Client, messageID common.Snowflake) error {
	return client.DeleteFollowup(i, messageID)
}
