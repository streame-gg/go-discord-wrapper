package discord

import "context"

// ── Options ───────────────────────────────────────────────────────────────────

// MessageEditOptions configures Message.Edit.
// https://docs.discord.com/developers/resources/message#edit-message
type MessageEditOptions struct {
	Content         *string
	Embeds          *[]Embed
	Flags           *MessageFlag
	AllowedMentions *AllowedMentions
	Components      []AnyComponent
	Attachments     *[]Attachment
	Files           []MessageFile
}

// MessageCreateOptions configures Message.Reply and Channel.Send.
// https://docs.discord.com/developers/resources/message#create-message
type MessageCreateOptions struct {
	Content          string
	TTS              bool
	Embeds           []Embed
	AllowedMentions  *AllowedMentions
	MessageReference *MessageMessageReference
	Components       []AnyComponent
	StickerIDs       []Snowflake
	Flags            MessageFlag
	Files            []MessageFile
}

// ── Hydration ─────────────────────────────────────────────────────────────────

// Hydrate injects the gateway client so that the convenience methods (Edit,
// Delete, Reply, etc.) can be called without passing ctx and client explicitly.
// Called automatically by the gateway dispatcher and cache stores.
func (m *Message) Hydrate(c EntityClient) {
	m.hClient = c
	m.hydrateNested(c)
}

// WithClient returns a shallow copy of the Message with the client replaced by
// c. Use this to transfer a message to a different client (e.g. shard bots).
func (m *Message) WithClient(c EntityClient) *Message {
	cp := *m
	cp.hClient = c
	return &cp
}

// IsHydrated reports whether the message has an associated client.
func (m *Message) IsHydrated() bool { return m.hClient != nil }

// hydrateNested injects the client into sub-entities embedded in this message.
func (m *Message) hydrateNested(c EntityClient) {
	if m.Author != nil {
		m.Author.Hydrate(c)
	}
	for i := range m.Mentions {
		m.Mentions[i].Hydrate(c)
	}
	if m.ReferencedMessage != nil {
		m.ReferencedMessage.Hydrate(c)
	}
}

// ── Entity methods ────────────────────────────────────────────────────────────

// Edit modifies this message. Requires SEND_MESSAGES (own message) or
// MANAGE_MESSAGES (other message).
func (m *Message) Edit(ctx context.Context, opts MessageEditOptions) (*Message, error) {
	c, err := ensureClient(m.hClient)
	if err != nil {
		return nil, err
	}
	return c.EditMessage(ctx, m.ChannelID, m.ID, opts)
}

// Delete removes this message. Pass a non-nil reason to set an audit log entry.
func (m *Message) Delete(ctx context.Context, reason *string) error {
	c, err := ensureClient(m.hClient)
	if err != nil {
		return err
	}
	return c.DeleteMessage(ctx, m.ChannelID, m.ID, reason)
}

// Reply sends a new message in the same channel with a reply reference set to
// this message.
func (m *Message) Reply(ctx context.Context, opts MessageCreateOptions) (*Message, error) {
	c, err := ensureClient(m.hClient)
	if err != nil {
		return nil, err
	}
	ref := MessageMessageReference{MessageID: &m.ID, ChannelID: &m.ChannelID}
	opts.MessageReference = &ref
	return c.CreateMessage(ctx, m.ChannelID, opts)
}

// Pin pins this message in its channel. Requires MANAGE_MESSAGES.
func (m *Message) Pin(ctx context.Context, reason *string) error {
	c, err := ensureClient(m.hClient)
	if err != nil {
		return err
	}
	return c.PinMessage(ctx, m.ChannelID, m.ID, reason)
}

// Unpin removes this message from the channel's pins. Requires MANAGE_MESSAGES.
func (m *Message) Unpin(ctx context.Context, reason *string) error {
	c, err := ensureClient(m.hClient)
	if err != nil {
		return err
	}
	return c.UnpinMessage(ctx, m.ChannelID, m.ID, reason)
}

// CrossPost publishes this message to follower channels. Only valid in
// announcement channels.
func (m *Message) CrossPost(ctx context.Context) (*Message, error) {
	c, err := ensureClient(m.hClient)
	if err != nil {
		return nil, err
	}
	return c.CrosspostMessage(ctx, m.ChannelID, m.ID)
}

// React adds a reaction to this message. emoji should be a Unicode character
// (e.g. "👍") or "name:id" for a custom emoji.
func (m *Message) React(ctx context.Context, emoji string) error {
	c, err := ensureClient(m.hClient)
	if err != nil {
		return err
	}
	return c.AddReaction(ctx, m.ChannelID, m.ID, emoji)
}
