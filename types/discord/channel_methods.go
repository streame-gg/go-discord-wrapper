package discord

import "context"

// ── Options ───────────────────────────────────────────────────────────────────

// ChannelEditOptions configures Channel.Edit.
type ChannelEditOptions struct {
	Name                          *string
	Type                          *ChannelType
	Position                      *int
	Topic                         *string
	NSFW                          *bool
	RateLimitPerUser              *int
	Bitrate                       *int
	UserLimit                     *int
	PermissionOverwrites          []ChannelPermissionOverwrite
	ParentID                      *Snowflake
	RTCRegion                     *string
	VideoQualityMode              *VideoQualityMode
	DefaultAutoArchiveDuration    *int
	Flags                         *ChannelFlags
	AvailableTags                 []ChannelTag
	DefaultReactionEmoji          *DefaultReactionEmoji
	DefaultThreadRateLimitPerUser *int
	DefaultSortOrder              *DefaultSortOrder
	DefaultForumLayout            *ChannelForumLayoutType
	AuditLogReason                *string
}

// ChannelCreateOptions configures Guild.CreateChannel.
type ChannelCreateOptions struct {
	Name                          string
	Type                          *ChannelType
	Topic                         *string
	Bitrate                       *int
	UserLimit                     *int
	RateLimitPerUser              *int
	Position                      *int
	PermissionOverwrites          []ChannelPermissionOverwrite
	ParentID                      *Snowflake
	NSFW                          *bool
	RTCRegion                     *string
	VideoQualityMode              *VideoQualityMode
	DefaultAutoArchiveDuration    *int
	DefaultReactionEmoji          *DefaultReactionEmoji
	AvailableTags                 []ChannelTag
	DefaultSortOrder              *DefaultSortOrder
	DefaultForumLayout            *ChannelForumLayoutType
	DefaultThreadRateLimitPerUser *int
	AuditLogReason                *string
}

// FetchMessagesOptions configures Channel.FetchMessages.
type FetchMessagesOptions struct {
	Around *Snowflake
	Before *Snowflake
	After  *Snowflake
	Limit  *int
}

// InviteCreateOptions configures Channel.CreateInvite.
type InviteCreateOptions struct {
	MaxAge              *int
	MaxUses             *int
	Temporary           *bool
	Unique              *bool
	TargetType          *int
	TargetUserID        *Snowflake
	TargetApplicationID *Snowflake
	AuditLogReason      *string
}

// ── Hydration ─────────────────────────────────────────────────────────────────

// Hydrate injects the gateway client so that the convenience methods can be
// called without passing ctx and client explicitly.
func (ch *Channel) Hydrate(c EntityClient) {
	ch.hClient = c
}

// WithClient returns a shallow copy of the Channel with the client replaced.
func (ch *Channel) WithClient(c EntityClient) *Channel {
	cp := *ch
	cp.hClient = c
	return &cp
}

// IsHydrated reports whether the channel has an associated client.
func (ch *Channel) IsHydrated() bool { return ch.hClient != nil }

// ── Entity methods ────────────────────────────────────────────────────────────

// Send sends a new message to this channel.
func (ch *Channel) Send(ctx context.Context, opts MessageCreateOptions) (*Message, error) {
	c, err := ensureClient(ch.hClient)
	if err != nil {
		return nil, err
	}
	return c.CreateMessage(ctx, ch.ID, opts)
}

// Edit modifies this channel's settings. Requires MANAGE_CHANNELS.
func (ch *Channel) Edit(ctx context.Context, opts ChannelEditOptions) (*Channel, error) {
	c, err := ensureClient(ch.hClient)
	if err != nil {
		return nil, err
	}
	return c.ModifyChannel(ctx, ch.ID, opts)
}

// Delete deletes this channel. Requires MANAGE_CHANNELS or MANAGE_THREADS.
// Pass a non-nil reason for the audit log.
func (ch *Channel) Delete(ctx context.Context, reason *string) (*Channel, error) {
	c, err := ensureClient(ch.hClient)
	if err != nil {
		return nil, err
	}
	return c.DeleteChannel(ctx, ch.ID, reason)
}

// BulkDelete deletes multiple messages at once. Requires MANAGE_MESSAGES.
// Messages must be ≤ 14 days old.
func (ch *Channel) BulkDelete(ctx context.Context, messageIDs []Snowflake, reason *string) error {
	c, err := ensureClient(ch.hClient)
	if err != nil {
		return err
	}
	return c.BulkDeleteMessages(ctx, ch.ID, messageIDs, reason)
}

// FetchMessages retrieves up to 100 messages from this channel.
func (ch *Channel) FetchMessages(ctx context.Context, opts FetchMessagesOptions) ([]*Message, error) {
	c, err := ensureClient(ch.hClient)
	if err != nil {
		return nil, err
	}
	return c.GetChannelMessages(ctx, ch.ID, opts)
}

// TriggerTyping posts a typing indicator for ~10 seconds.
func (ch *Channel) TriggerTyping(ctx context.Context) error {
	c, err := ensureClient(ch.hClient)
	if err != nil {
		return err
	}
	return c.TriggerTypingIndicator(ctx, ch.ID)
}

// SetVoiceStatus sets the voice channel status string. Pass nil to clear it.
func (ch *Channel) SetVoiceStatus(ctx context.Context, status *string) error {
	c, err := ensureClient(ch.hClient)
	if err != nil {
		return err
	}
	return c.SetVoiceChannelStatus(ctx, ch.ID, status)
}

// CreateInvite creates a new invite for this channel. Requires CREATE_INSTANT_INVITE.
func (ch *Channel) CreateInvite(ctx context.Context, opts InviteCreateOptions) (*Invite, error) {
	c, err := ensureClient(ch.hClient)
	if err != nil {
		return nil, err
	}
	return c.CreateChannelInvite(ctx, ch.ID, opts)
}
