package discord

import (
	"context"
	"time"
)

// ── Options ───────────────────────────────────────────────────────────────────

// GuildEditOptions configures Guild.Edit.
// https://docs.discord.com/developers/resources/guild#modify-guild
type GuildEditOptions struct {
	Name                        *string
	VerificationLevel           *GuildVerificationLevel
	DefaultMessageNotifications *DefaultMessageNotificationLevel
	ExplicitContentFilter       *GuildExplicitContentFilterLevel
	AFKChannelID                *Snowflake
	AFKTimeout                  *int
	Icon                        *string
	OwnerID                     *Snowflake
	Splash                      *string
	DiscoverySplash             *string
	Banner                      *string
	SystemChannelID             *Snowflake
	SystemChannelFlags          *GuildSystemChannelFlags
	RulesChannelID              *Snowflake
	PublicUpdatesChannelID      *Snowflake
	PreferredLocale             *Locale
	Features                    []GuildFeatures
	Description                 *string
	PremiumProgressBarEnabled   *bool
	SafetyAlertsChannelID       *Snowflake
	AuditLogReason              *string
}

// FetchMembersOptions configures Guild.FetchMembers.
// https://docs.discord.com/developers/resources/guild#list-guild-members
type FetchMembersOptions struct {
	After *Snowflake
	Limit *int
}

// FetchBansOptions configures BanManager.FetchAll.
// https://docs.discord.com/developers/resources/guild#get-guild-bans
type FetchBansOptions struct {
	Before *Snowflake
	After  *Snowflake
	Limit  *int
}

// AuditLogOptions configures Guild.FetchAuditLog.
// https://docs.discord.com/developers/resources/audit-log#get-guild-audit-log
type AuditLogOptions struct {
	UserID     *Snowflake
	ActionType *AuditLogActionType
	Before     *Snowflake
	After      *Snowflake
	Limit      *int
}

// ScheduledEventCreateOptions configures Guild.CreateScheduledEvent.
// https://docs.discord.com/developers/resources/guild-scheduled-event#create-guild-scheduled-event
type ScheduledEventCreateOptions struct {
	ChannelID          *Snowflake
	EntityMetadata     *GuildScheduledEventEntityMetadata
	Name               string
	PrivacyLevel       GuildScheduledEventPrivacyLevel
	ScheduledStartTime time.Time
	ScheduledEndTime   *time.Time
	Description        *string
	EntityType         GuildScheduledEventEntityType
	Image              *string
}

// ── Hydration ─────────────────────────────────────────────────────────────────

// Hydrate injects the gateway client so that the convenience methods can be
// called without passing ctx and client explicitly.
func (g *Guild) Hydrate(c EntityClient) {
	g.hClient = c
	g.hydrateNested(c)
}

// WithClient returns a shallow copy of the Guild with the client replaced.
func (g *Guild) WithClient(c EntityClient) *Guild {
	cp := *g
	cp.hClient = c
	return &cp
}

// IsHydrated reports whether the guild has an associated client.
func (g *Guild) IsHydrated() bool { return g.hClient != nil }

// hydrateNested injects the client into sub-entities embedded in this guild.
func (g *Guild) hydrateNested(c EntityClient) {
	for i := range g.RawRoles {
		g.RawRoles[i].GuildID = g.ID
		g.RawRoles[i].Hydrate(c)
	}
	for i := range g.RawEmojis {
		g.RawEmojis[i].GuildID = g.ID
		g.RawEmojis[i].Hydrate(c)
	}
	for i := range g.RawStickers {
		g.RawStickers[i].GuildID = &g.ID
		g.RawStickers[i].Hydrate(c)
	}
}

// ── Entity methods ────────────────────────────────────────────────────────────

// Edit modifies this guild's settings. Requires MANAGE_GUILD.
func (g *Guild) Edit(ctx context.Context, opts GuildEditOptions) (*Guild, error) {
	c, err := ensureClient(g.hClient)
	if err != nil {
		return nil, err
	}
	return c.ModifyGuild(ctx, g.ID, opts)
}

// Leave causes the bot to leave this guild.
func (g *Guild) Leave(ctx context.Context) error {
	c, err := ensureClient(g.hClient)
	if err != nil {
		return err
	}
	return c.LeaveGuild(ctx, g.ID)
}

// FetchMembers retrieves a paginated list of guild members.
func (g *Guild) FetchMembers(ctx context.Context, opts FetchMembersOptions) ([]*GuildMember, error) {
	c, err := ensureClient(g.hClient)
	if err != nil {
		return nil, err
	}
	return c.ListGuildMembers(ctx, g.ID, opts)
}

// CreateRole creates a new role in this guild. Requires MANAGE_ROLES.
func (g *Guild) CreateRole(ctx context.Context, opts RoleCreateOptions) (*Role, error) {
	c, err := ensureClient(g.hClient)
	if err != nil {
		return nil, err
	}
	return c.CreateGuildRole(ctx, g.ID, opts)
}

// CreateChannel creates a new channel in this guild. Requires MANAGE_CHANNELS.
func (g *Guild) CreateChannel(ctx context.Context, opts ChannelCreateOptions) (*Channel, error) {
	c, err := ensureClient(g.hClient)
	if err != nil {
		return nil, err
	}
	return c.CreateGuildChannel(ctx, g.ID, opts)
}

// CreateEmoji creates a new emoji in this guild. Requires MANAGE_GUILD_EXPRESSIONS.
func (g *Guild) CreateEmoji(ctx context.Context, opts EmojiCreateOptions) (*Emoji, error) {
	c, err := ensureClient(g.hClient)
	if err != nil {
		return nil, err
	}
	return c.CreateGuildEmoji(ctx, g.ID, opts)
}

// CreateSticker creates a new sticker in this guild. Requires MANAGE_GUILD_EXPRESSIONS.
func (g *Guild) CreateSticker(ctx context.Context, opts StickerCreateOptions) (*Sticker, error) {
	c, err := ensureClient(g.hClient)
	if err != nil {
		return nil, err
	}
	return c.CreateGuildSticker(ctx, g.ID, opts)
}

// FetchAuditLog retrieves the guild's audit log entries. Requires VIEW_AUDIT_LOG.
func (g *Guild) FetchAuditLog(ctx context.Context, opts AuditLogOptions) (*AuditLog, error) {
	c, err := ensureClient(g.hClient)
	if err != nil {
		return nil, err
	}
	return c.GetGuildAuditLog(ctx, g.ID, opts)
}

// CreateScheduledEvent creates a new scheduled event in this guild.
// Requires MANAGE_EVENTS.
func (g *Guild) CreateScheduledEvent(ctx context.Context, opts ScheduledEventCreateOptions) (*GuildScheduledEvent, error) {
	c, err := ensureClient(g.hClient)
	if err != nil {
		return nil, err
	}
	return c.CreateGuildScheduledEvent(ctx, g.ID, opts)
}
