package interactions

import (
	"context"
	"errors"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// errNotHydrated is returned by Interaction methods when the interaction was
// not dispatched through the gateway (i.e. client/ctx were never set).
var errNotHydrated = errors.New("interactions: Interaction was not dispatched via the gateway — use the client and ctx arguments directly")

// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-object
type Interaction struct {
	ID                           discord.Snowflake                                        `json:"id"`
	ApplicationID                discord.Snowflake                                        `json:"application_id"`
	Type                         discord.InteractionType                                  `json:"type"`
	Data                         *discord.InteractionData                                 `json:"data,omitempty"`
	GuildID                      *discord.Snowflake                                       `json:"guild_id,omitempty"`
	ChannelID                    *discord.Snowflake                                       `json:"channel_id,omitempty"`
	Guild                        *discord.Guild                                           `json:"guild,omitempty"`
	Channel                      *discord.Channel                                         `json:"channel,omitempty"`
	Member                       *discord.GuildMember                                     `json:"member,omitempty"`
	User                         *discord.User                                            `json:"user,omitempty"`
	Token                        string                                                   `json:"token"`
	Version                      int                                                      `json:"version"`
	Message                      *discord.Message                                         `json:"message,omitempty"`
	AppPermissions               discord.Permission                                       `json:"app_permissions,omitempty"`
	Locale                       discord.Locale                                           `json:"locale,omitempty"`
	GuildLocale                  discord.Locale                                           `json:"guild_locale,omitempty"`
	Entitlements                 []discord.Entitlement                                    `json:"entitlements,omitempty"`
	AuthorizingIntegrationOwners map[discord.ApplicationIntegrationType]discord.Snowflake `json:"authorizing_integration_owners,omitempty"`
	Context                      discord.InteractionContextType                           `json:"context,omitempty"`
	AttachmentSizeLimit          int                                                      `json:"attachment_size_limit,omitempty"`

	// hClient and hCtx are injected by the gateway dispatcher before calling
	// user handlers. They are nil/zero when the Interaction is constructed
	// outside the gateway (e.g. in tests). Use WithContext to override hCtx.
	hClient Client
	hCtx    context.Context
}

// Hydrate injects the gateway client and dispatch context into the Interaction
// so that the convenience methods (Reply, DeferReply, etc.) can be called
// without passing ctx and client explicitly. Called by the gateway dispatcher.
func (i *Interaction) Hydrate(c Client, ctx context.Context) {
	i.hClient = c
	i.hCtx = ctx
}

// WithContext returns a shallow copy of the Interaction with the context
// replaced by ctx. Use this inside goroutines or when you need a different
// deadline for follow-up calls than the one captured at dispatch time.
//
//	go func() {
//	    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
//	    defer cancel()
//	    i.WithContext(ctx).Reply(interactions.ReplyOptions{Content: "Late reply"})
//	}()
func (i *Interaction) WithContext(ctx context.Context) *Interaction {
	c := *i
	c.hCtx = ctx
	return &c
}

// ensureHydrated returns the client and ctx, or errNotHydrated if the
// Interaction was not dispatched via the gateway.
func (i *Interaction) ensureHydrated() (Client, context.Context, error) {
	if i.hClient == nil {
		return nil, nil, errNotHydrated
	}
	ctx := i.hCtx
	if ctx == nil {
		ctx = context.Background()
	}
	return i.hClient, ctx, nil
}
