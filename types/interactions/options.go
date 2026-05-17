package interactions

import (
	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

// ReplyOptions configures an immediate message response to an interaction.
type ReplyOptions struct {
	Content               string
	Embeds                []discord.Embed
	Components            []discord.AnyComponent
	Files                 []discord.MessageFile
	AllowedMentions       *discord.AllowedMentions
	Ephemeral             bool
	SuppressEmbeds        bool
	SuppressNotifications bool
	// WithResponse requests the created message in the response body.
	WithResponse bool
	Poll         *discord.PollRequest
}

func (o ReplyOptions) toResponseData() *responses.InteractionResponseDataDefault {
	data := &responses.InteractionResponseDataDefault{}
	if o.Content != "" {
		data.Content = o.Content
	}
	if len(o.Embeds) > 0 {
		embeds := o.Embeds
		data.Embeds = &embeds
	}
	if len(o.Components) > 0 {
		comps := o.Components
		data.Components = &comps
	}
	if o.AllowedMentions != nil {
		data.AllowedMentions = o.AllowedMentions
	}
	if o.Poll != nil {
		data.Poll = o.Poll
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
	data.Flags = flags
	return data
}

// DeferOptions configures an acknowledgement-only response (loading state).
type DeferOptions struct {
	// Ephemeral makes the eventual follow-up visible only to the invoking user.
	Ephemeral bool
}

// UpdateMessageOptions configures an in-place edit of the component message.
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
		data.Content = o.Content
	}
	if len(o.Embeds) > 0 {
		embeds := o.Embeds
		data.Embeds = &embeds
	}
	if len(o.Components) > 0 {
		comps := o.Components
		data.Components = &comps
	}
	if o.AllowedMentions != nil {
		data.AllowedMentions = o.AllowedMentions
	}

	return data
}

// FollowUpOptions configures a follow-up message (usable up to 15 minutes
// after the initial response).
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
		Content:         o.Content,
		Embeds:          o.Embeds,
		Components:      o.Components,
		Files:           o.Files,
		AllowedMentions: o.AllowedMentions,
	}
	if o.Poll != nil {
		params.Poll = *o.Poll
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
	params.Flags = flags
	return params
}

// AutocompleteOptions holds the choices for an autocomplete response.
type AutocompleteOptions struct {
	Choices []responses.AutocompleteChoice
}

// ModalOptions holds the modal to present in response to an interaction.
type ModalOptions struct {
	Modal *components.Modal
}
