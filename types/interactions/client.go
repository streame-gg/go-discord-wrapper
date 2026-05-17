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

// ── Hydrated convenience methods ──────────────────────────────────────────────
//
// These methods use the client and context injected by the gateway dispatcher.
// They are zero-overhead wrappers: if the Interaction was not dispatched via
// the gateway, they return errNotHydrated immediately.
//
// Use WithContext to override the captured context (e.g. in background goroutines).

// Reply sends an immediate message response to the interaction.
//
//	_, err := i.Reply(interactions.ReplyOptions{Content: "Hello!"})
func (i *Interaction) Reply(opts ReplyOptions) (*responses.InteractionCallbackResponse, error) {
	c, ctx, err := i.ensureHydrated()
	if err != nil {
		return nil, err
	}
	return c.Reply(ctx, i, opts.toResponseData(), opts.WithResponse, opts.Files...)
}

// DeferReply acknowledges the interaction without an immediate visible response.
//
//	err := i.DeferReply(interactions.DeferOptions{Ephemeral: true})
func (i *Interaction) DeferReply(opts DeferOptions) error {
	c, ctx, err := i.ensureHydrated()
	if err != nil {
		return err
	}
	return c.DeferReply(ctx, i, opts.Ephemeral)
}

// DeferUpdateMessage acknowledges a component interaction without editing the
// original message. Call EditReply afterwards to push the actual update.
func (i *Interaction) DeferUpdateMessage() error {
	c, ctx, err := i.ensureHydrated()
	if err != nil {
		return err
	}
	return c.DeferUpdateMessage(ctx, i)
}

// UpdateMessage edits the message that triggered a component interaction.
func (i *Interaction) UpdateMessage(opts UpdateMessageOptions) error {
	c, ctx, err := i.ensureHydrated()
	if err != nil {
		return err
	}
	return c.UpdateMessage(ctx, i, opts.toResponseData())
}

// Modal responds to the interaction with a modal dialog.
func (i *Interaction) Modal(opts ModalOptions) error {
	c, ctx, err := i.ensureHydrated()
	if err != nil {
		return err
	}
	return c.ReplyWithModal(ctx, i, opts.Modal)
}

// Autocomplete sends autocomplete choices back to Discord.
func (i *Interaction) Autocomplete(opts AutocompleteOptions) error {
	c, ctx, err := i.ensureHydrated()
	if err != nil {
		return err
	}
	return c.ReplyAutocomplete(ctx, i, opts.Choices)
}

// LaunchActivity responds by launching the app's associated Activity.
func (i *Interaction) LaunchActivity() error {
	c, ctx, err := i.ensureHydrated()
	if err != nil {
		return err
	}
	return c.LaunchActivity(ctx, i)
}

// GetOriginalResponse fetches the original response message for this interaction.
func (i *Interaction) GetOriginalResponse() (*discord.Message, error) {
	c, ctx, err := i.ensureHydrated()
	if err != nil {
		return nil, err
	}
	return c.GetOriginalResponse(ctx, i)
}

// EditReply edits the original interaction response.
func (i *Interaction) EditReply(params api.EditMessageParams) (*discord.Message, error) {
	c, ctx, err := i.ensureHydrated()
	if err != nil {
		return nil, err
	}
	return c.EditReply(ctx, i, params)
}

// DeleteReply deletes the original interaction response.
func (i *Interaction) DeleteReply() error {
	c, ctx, err := i.ensureHydrated()
	if err != nil {
		return err
	}
	return c.DeleteReply(ctx, i)
}

// FollowUp sends a follow-up message (usable up to 15 minutes after the initial response).
func (i *Interaction) FollowUp(opts FollowUpOptions) (*discord.Message, error) {
	c, ctx, err := i.ensureHydrated()
	if err != nil {
		return nil, err
	}
	return c.CreateFollowup(ctx, i, opts.toCreateParams())
}

// GetFollowup fetches a follow-up message by ID.
func (i *Interaction) GetFollowup(messageID discord.Snowflake) (*discord.Message, error) {
	c, ctx, err := i.ensureHydrated()
	if err != nil {
		return nil, err
	}
	return c.GetFollowup(ctx, i, messageID)
}

// EditFollowup edits a follow-up message by ID.
func (i *Interaction) EditFollowup(messageID discord.Snowflake, params api.EditMessageParams) (*discord.Message, error) {
	c, ctx, err := i.ensureHydrated()
	if err != nil {
		return nil, err
	}
	return c.EditFollowup(ctx, i, messageID, params)
}

// DeleteFollowup deletes a follow-up message by ID.
func (i *Interaction) DeleteFollowup(messageID discord.Snowflake) error {
	c, ctx, err := i.ensureHydrated()
	if err != nil {
		return err
	}
	return c.DeleteFollowup(ctx, i, messageID)
}

// DeferAndFollowup acknowledges the interaction immediately (Discord shows a
// loading state) and returns a sender function that delivers the actual
// follow-up when called. Use this pattern for handlers that need to do async
// work before responding.
//
// The returned send function accepts its own context so that the caller can
// control the follow-up deadline independently of the defer call.
//
//	send, err := i.DeferAndFollowup(interactions.DeferOptions{})
//	if err != nil { return }
//	go func() {
//	    // do slow work…
//	    send(context.Background(), interactions.FollowUpOptions{Content: "Done!"})
//	}()
func (i *Interaction) DeferAndFollowup(opts DeferOptions) (func(context.Context, FollowUpOptions) (*discord.Message, error), error) {
	c, ctx, err := i.ensureHydrated()
	if err != nil {
		return nil, err
	}
	if err := c.DeferReply(ctx, i, opts.Ephemeral); err != nil {
		return nil, err
	}
	return func(sendCtx context.Context, followOpts FollowUpOptions) (*discord.Message, error) {
		return c.CreateFollowup(sendCtx, i, followOpts.toCreateParams())
	}, nil
}
