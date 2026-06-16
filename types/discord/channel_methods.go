package discord

import "context"

// ── Options ───────────────────────────────────────────────────────────────────

// ChannelEditOptions configures Channel.Edit.
// https://docs.discord.com/developers/resources/channel#modify-channel
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
// https://docs.discord.com/developers/resources/guild#create-guild-channel
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
// https://docs.discord.com/developers/resources/channel#get-channel-messages
type FetchMessagesOptions struct {
	Around *Snowflake
	Before *Snowflake
	After  *Snowflake
	Limit  *int
}

// InviteCreateOptions configures Channel.CreateInvite.
// https://docs.discord.com/developers/resources/channel#create-channel-invite
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
func (c *Channel) Hydrate(ec EntityClient) {
	c.hClient = ec
	ec.SetChannelManagers(c)
}

// WithClient returns a shallow copy of the Channel with the client replaced.
func (c *Channel) WithClient(ec EntityClient) *Channel {
	cp := *c
	cp.hClient = ec
	return &cp
}

// IsHydrated reports whether the channel has an associated client.
func (c *Channel) IsHydrated() bool { return c.hClient != nil }

// ── Entity methods ────────────────────────────────────────────────────────────

// Send sends a new message to this channel.
func (c *Channel) Send(ctx context.Context, opts MessageCreateOptions) (*Message, error) {
	cl, err := ensureClient(c.hClient)
	if err != nil {
		return nil, err
	}
	return cl.CreateMessage(ctx, c.ID, opts)
}

// Edit modifies this channel's settings. Requires MANAGE_CHANNELS.
func (c *Channel) Edit(ctx context.Context, opts ChannelEditOptions) (*Channel, error) {
	cl, err := ensureClient(c.hClient)
	if err != nil {
		return nil, err
	}
	return cl.ModifyChannel(ctx, c.ID, opts)
}

// Delete deletes this channel. Requires MANAGE_CHANNELS or MANAGE_THREADS.
// Pass a non-nil reason for the audit log.
func (c *Channel) Delete(ctx context.Context, reason *string) (*Channel, error) {
	cl, err := ensureClient(c.hClient)
	if err != nil {
		return nil, err
	}
	return cl.DeleteChannel(ctx, c.ID, reason)
}

// BulkDelete deletes multiple messages at once. Requires MANAGE_MESSAGES.
// Messages must be ≤ 14 days old.
func (c *Channel) BulkDelete(ctx context.Context, messageIDs []Snowflake, reason *string) error {
	cl, err := ensureClient(c.hClient)
	if err != nil {
		return err
	}
	return cl.BulkDeleteMessages(ctx, c.ID, messageIDs, reason)
}

// FetchMessages retrieves up to 100 messages from this channel.
func (c *Channel) FetchMessages(ctx context.Context, opts FetchMessagesOptions) ([]*Message, error) {
	cl, err := ensureClient(c.hClient)
	if err != nil {
		return nil, err
	}
	return cl.ListChannelMessages(ctx, c.ID, opts)
}

// TriggerTyping posts a typing indicator for ~10 seconds.
func (c *Channel) TriggerTyping(ctx context.Context) error {
	cl, err := ensureClient(c.hClient)
	if err != nil {
		return err
	}
	return cl.TriggerTypingIndicator(ctx, c.ID)
}

// SetVoiceStatus sets the voice channel status string. Pass nil to clear it.
func (c *Channel) SetVoiceStatus(ctx context.Context, status *string) error {
	cl, err := ensureClient(c.hClient)
	if err != nil {
		return err
	}
	return cl.SetVoiceChannelStatus(ctx, c.ID, status)
}

// CreateInvite creates a new invite for this channel. Requires CREATE_INSTANT_INVITE.
func (c *Channel) CreateInvite(ctx context.Context, opts InviteCreateOptions) (*Invite, error) {
	cl, err := ensureClient(c.hClient)
	if err != nil {
		return nil, err
	}
	return cl.CreateChannelInvite(ctx, c.ID, opts)
}
