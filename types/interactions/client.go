package interactions

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

// Client is the subset of *connection.Client needed by Interaction methods.
// Pass the client received in your event handler — *connection.Client satisfies this interface.
type Client interface {
	Reply(ctx context.Context, i *Interaction, data *responses.InteractionResponseDataDefault, withResponse bool, files ...api.MessageFile) (*responses.InteractionCallbackResponse, error)
	DeferReply(ctx context.Context, i *Interaction, ephemeral bool) error
	DeferUpdateMessage(ctx context.Context, i *Interaction) error
	UpdateMessage(ctx context.Context, i *Interaction, data *responses.InteractionResponseDataDefault) error
	ReplyWithModal(ctx context.Context, i *Interaction, modal *components.Modal) error
	ReplyAutocomplete(ctx context.Context, i *Interaction, choices []responses.AutocompleteChoice) error
	LaunchActivity(ctx context.Context, i *Interaction) error
	GetOriginalResponse(ctx context.Context, i *Interaction) (*discord.Message, error)
	EditReply(ctx context.Context, i *Interaction, params api.EditMessageParams) (*discord.Message, error)
	DeleteReply(ctx context.Context, i *Interaction) error
	CreateFollowup(ctx context.Context, i *Interaction, params api.CreateMessageParams) (*discord.Message, error)
	GetFollowup(ctx context.Context, i *Interaction, messageID discord.Snowflake) (*discord.Message, error)
	EditFollowup(ctx context.Context, i *Interaction, messageID discord.Snowflake, params api.EditMessageParams) (*discord.Message, error)
	DeleteFollowup(ctx context.Context, i *Interaction, messageID discord.Snowflake) error
}

// Reply sends an immediate message response to the interaction.
// Pass optional files to include attachments; when present the request is sent as multipart/form-data.
func (i *Interaction) Reply(ctx context.Context, client Client, data *responses.InteractionResponseDataDefault, withResponse bool, files ...api.MessageFile) (*responses.InteractionCallbackResponse, error) {
	return client.Reply(ctx, i, data, withResponse, files...)
}

// DeferReply acknowledges the interaction without an immediate visible response.
// Set ephemeral=true to restrict the eventual follow-up to the invoker only.
func (i *Interaction) DeferReply(ctx context.Context, client Client, ephemeral bool) error {
	return client.DeferReply(ctx, i, ephemeral)
}

// DeferUpdateMessage acknowledges a component interaction without editing the original message.
// Call EditReply afterwards to push the actual update.
func (i *Interaction) DeferUpdateMessage(ctx context.Context, client Client) error {
	return client.DeferUpdateMessage(ctx, i)
}

// UpdateMessage edits the message that triggered a component interaction.
func (i *Interaction) UpdateMessage(ctx context.Context, client Client, data *responses.InteractionResponseDataDefault) error {
	return client.UpdateMessage(ctx, i, data)
}

// Modal responds to the interaction with a modal dialog.
func (i *Interaction) Modal(ctx context.Context, client Client, modal *components.Modal) error {
	return client.ReplyWithModal(ctx, i, modal)
}

// Autocomplete sends autocomplete choices back to Discord.
func (i *Interaction) Autocomplete(ctx context.Context, client Client, choices []responses.AutocompleteChoice) error {
	return client.ReplyAutocomplete(ctx, i, choices)
}

// LaunchActivity responds by launching the app's associated Activity.
func (i *Interaction) LaunchActivity(ctx context.Context, client Client) error {
	return client.LaunchActivity(ctx, i)
}

// GetOriginalResponse fetches the original response message for this interaction.
func (i *Interaction) GetOriginalResponse(ctx context.Context, client Client) (*discord.Message, error) {
	return client.GetOriginalResponse(ctx, i)
}

// EditReply edits the original interaction response.
func (i *Interaction) EditReply(ctx context.Context, client Client, params api.EditMessageParams) (*discord.Message, error) {
	return client.EditReply(ctx, i, params)
}

// DeleteReply deletes the original interaction response.
func (i *Interaction) DeleteReply(ctx context.Context, client Client) error {
	return client.DeleteReply(ctx, i)
}

// FollowUp sends a follow-up message (usable up to 15 minutes after the initial response).
func (i *Interaction) FollowUp(ctx context.Context, client Client, params api.CreateMessageParams) (*discord.Message, error) {
	return client.CreateFollowup(ctx, i, params)
}

// GetFollowup fetches a follow-up message by ID.
func (i *Interaction) GetFollowup(ctx context.Context, client Client, messageID discord.Snowflake) (*discord.Message, error) {
	return client.GetFollowup(ctx, i, messageID)
}

// EditFollowup edits a follow-up message by ID.
func (i *Interaction) EditFollowup(ctx context.Context, client Client, messageID discord.Snowflake, params api.EditMessageParams) (*discord.Message, error) {
	return client.EditFollowup(ctx, i, messageID, params)
}

// DeleteFollowup deletes a follow-up message by ID.
func (i *Interaction) DeleteFollowup(ctx context.Context, client Client, messageID discord.Snowflake) error {
	return client.DeleteFollowup(ctx, i, messageID)
}

// DeferAndFollowup acknowledges the interaction immediately (Discord shows a
// loading state) and returns a sender function that delivers the actual
// follow-up when called. Use this pattern for handlers that need to do async
// work before responding.
//
//	send, err := i.DeferAndFollowup(ctx, client, false)
//	if err != nil { return }
//	_ = send // invoke send(api.CreateMessageParams{Content: result}) when ready
func (i *Interaction) DeferAndFollowup(ctx context.Context, client Client, ephemeral bool) (func(api.CreateMessageParams) (*discord.Message, error), error) {
	if err := client.DeferReply(ctx, i, ephemeral); err != nil {
		return nil, err
	}
	return func(params api.CreateMessageParams) (*discord.Message, error) {
		return client.CreateFollowup(ctx, i, params)
	}, nil
}
