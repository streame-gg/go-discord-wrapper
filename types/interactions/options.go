package interactions

import (
	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

// ReplyOptions configures an immediate message response to an interaction.
// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-response-object-messages
type ReplyOptions struct {
	Content         string
	Embeds          []discord.Embed
	Components      []discord.AnyComponent
	Files           []discord.MessageFile
	AllowedMentions *discord.AllowedMentions
	Flags           discord.MessageFlag
	// WithResponse requests the created message in the response body.
	WithResponse bool
	Poll         *discord.PollRequest
}

func (o ReplyOptions) toResponseData() *responses.InteractionResponseDataDefault {
	data := &responses.InteractionResponseDataDefault{}
	if o.Flags != 0 {
		data.Flags = discord.Some(o.Flags)
	}
	if o.Content != "" {
		data.Content = discord.Some(o.Content)
	}
	if len(o.Embeds) > 0 {
		data.Embeds = discord.Some(o.Embeds)
	}
	if len(o.Components) > 0 {
		data.Components = discord.Some(o.Components)
	}
	if o.AllowedMentions != nil {
		data.AllowedMentions = discord.Some(*o.AllowedMentions)
	}
	if o.Poll != nil {
		data.Poll = discord.Some(*o.Poll)
	}

	return data
}

// DeferOptions configures an acknowledgement-only response (loading state).
// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-response-object-interaction-callback-type
type DeferOptions struct {
	// Ephemeral makes the eventual follow-up visible only to the invoking user.
	Ephemeral bool
}

// UpdateMessageOptions configures an in-place edit of the component message.
// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-response-object-interaction-callback-type
type UpdateMessageOptions struct {
	Content         string
	Embeds          []discord.Embed
	Components      []discord.AnyComponent
	AllowedMentions *discord.AllowedMentions

	Files []discord.MessageFile
}

func (o UpdateMessageOptions) toResponseData() *responses.InteractionResponseDataDefault {
	data := &responses.InteractionResponseDataDefault{}
	if o.Content != "" {
		data.Content = discord.Some(o.Content)
	}
	if len(o.Embeds) > 0 {
		data.Embeds = discord.Some(o.Embeds)
	}
	if len(o.Components) > 0 {
		data.Components = discord.Some(o.Components)
	}
	if o.AllowedMentions != nil {
		data.AllowedMentions = discord.Some(*o.AllowedMentions)
	}

	return data
}

// FollowUpOptions configures a follow-up message (usable up to 15 minutes
// after the initial response).
// https://docs.discord.com/developers/interactions/receiving-and-responding#create-followup-message
type FollowUpOptions struct {
	Content               string
	Embeds                []discord.Embed
	Components            []discord.AnyComponent
	Files                 []discord.MessageFile
	AllowedMentions       *discord.AllowedMentions
	Ephemeral             bool
	SuppressEmbeds        bool
	SuppressNotifications bool
	Poll                  *discord.PollRequest
}

func (o FollowUpOptions) toCreateParams() api.CreateMessageParams {
	params := api.CreateMessageParams{
		Components: o.Components,
		Files:      o.Files,
	}
	if o.Content != "" {
		params.Content = discord.Some(o.Content)
	}
	if len(o.Embeds) > 0 {
		params.Embeds = discord.Some(o.Embeds)
	}
	if o.AllowedMentions != nil {
		params.AllowedMentions = discord.Some(*o.AllowedMentions)
	}
	if o.Poll != nil {
		params.Poll = discord.Some(*o.Poll)
	}
	var flags discord.MessageFlag
	if o.Ephemeral {
		flags |= discord.MessageFlagEphemeral
	}
	if o.SuppressEmbeds {
		flags |= discord.MessageFlagSuppressEmbeds
	}
	if o.SuppressNotifications {
		flags |= discord.MessageFlagSuppressNotification
	}
	if flags != 0 {
		params.Flags = discord.Some(flags)
	}
	return params
}

// AutocompleteOptions holds the choices for an autocomplete response.
// https://docs.discord.com/developers/interactions/receiving-and-responding#autocomplete
type AutocompleteOptions struct {
	Choices []responses.AutocompleteChoice
}

// ModalOptions holds the modal to present in response to an interaction.
// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-response-object-modal
type ModalOptions struct {
	Modal *components.Modal
}
