package api

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/util"
)

// FetchAllGuildMembers fetches every member in a guild across as many pages as
// needed. It uses the after-cursor pattern (max 1000 per request) and stops when
// a page returns fewer members than the page size.
func (c *RestClient) FetchAllGuildMembers(ctx context.Context, guildID discord.Snowflake) ([]*discord.GuildMember, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	const pageSize = 1000

	var all []*discord.GuildMember
	params := GetGuildMembersParams{Limit: util.PointerOf(pageSize)}

	for {
		page, err := c.ListGuildMembers(ctx, guildID, params)
		if err != nil {
			return nil, err
		}

		all = append(all, page...)

		if len(page) < pageSize {
			return all, nil
		}

		// Advance the cursor past the last member on this page.
		last := page[len(page)-1]
		params.After = &last.User.ID
	}
}

// FetchAllMessages fetches all messages in a channel, walking backwards from the
// most recent message (max 100 per request) until the beginning of the channel.
func (c *RestClient) FetchAllMessages(ctx context.Context, channelID discord.Snowflake) ([]*discord.Message, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	const pageSize = 100

	var all []*discord.Message
	params := GetMessagesParams{Limit: util.PointerOf(pageSize)}

	for {
		page, err := c.GetMessages(ctx, channelID, params)
		if err != nil {
			return nil, err
		}

		all = append(all, page...)

		if len(page) < pageSize {
			return all, nil
		}

		// Discord returns messages newest-first; the last entry is the oldest on
		// this page, so use its ID as the next before cursor.
		oldest := page[len(page)-1]
		params.Before = &oldest.ID
	}
}

// FetchAllGuildBans fetches every ban in a guild across as many pages as needed
// (max 1000 per request).
func (c *RestClient) FetchAllGuildBans(ctx context.Context, guildID discord.Snowflake) ([]*discord.Ban, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	const pageSize = 1000

	var all []*discord.Ban
	params := GetGuildBansParams{Limit: util.PointerOf(pageSize)}

	for {
		page, err := c.GetGuildBans(ctx, guildID, params)
		if err != nil {
			return nil, err
		}

		all = append(all, page...)

		if len(page) < pageSize {
			return all, nil
		}

		last := page[len(page)-1]
		params.After = &last.User.ID
	}
}

// FetchAllAuditLogEntries fetches all audit log entries for a guild by paginating
// backwards in time (max 100 per request). The optional filter is forwarded to
// every request so callers can scope by user or action type.
func (c *RestClient) FetchAllAuditLogEntries(ctx context.Context, guildID discord.Snowflake, filter GetGuildAuditLogParams) ([]discord.AuditLogEntry, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	const pageSize = 100

	filter.Limit = util.PointerOf(pageSize)

	var all []discord.AuditLogEntry

	for {
		log, err := c.GetGuildAuditLog(ctx, guildID, filter)
		if err != nil {
			return nil, err
		}

		all = append(all, log.AuditLogEntries...)

		if len(log.AuditLogEntries) < pageSize {
			return all, nil
		}

		// Entries are returned newest-first; walk backwards using the oldest entry's ID.
		oldest := log.AuditLogEntries[len(log.AuditLogEntries)-1]
		filter.Before = &oldest.ID
	}
}

// FetchAllEntitlements fetches every entitlement for an application by paginating
// forward (max 100 per request). The filter is forwarded as-is so callers can
// scope by user, guild, or SKU.
func (c *RestClient) FetchAllEntitlements(ctx context.Context, appID discord.Snowflake, filter ListEntitlementsParams) ([]*discord.Entitlement, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	const pageSize = 100

	filter.Limit = util.PointerOf(pageSize)

	var all []*discord.Entitlement

	for {
		page, err := c.ListEntitlements(ctx, appID, filter)
		if err != nil {
			return nil, err
		}

		all = append(all, page...)

		if len(page) < pageSize {
			return all, nil
		}

		last := page[len(page)-1]
		filter.After = &last.ID
	}
}

// FetchAllScheduledEventUsers fetches every subscriber for a scheduled event by
// paginating forward (max 100 per request). Set withMember=true to include the
// full GuildMember object alongside each user.
func (c *RestClient) FetchAllScheduledEventUsers(ctx context.Context, guildID, eventID discord.Snowflake, withMember bool) ([]*discord.GuildScheduledEventUser, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := eventID.Validate(); err != nil {
		return nil, err
	}

	const pageSize = 100

	var all []*discord.GuildScheduledEventUser
	params := GetGuildScheduledEventUsersParams{
		Limit:      util.PointerOf(pageSize),
		WithMember: util.PointerOf(withMember),
	}

	for {
		page, err := c.GetGuildScheduledEventUsers(ctx, guildID, eventID, params)
		if err != nil {
			return nil, err
		}

		all = append(all, page...)

		if len(page) < pageSize {
			return all, nil
		}

		last := page[len(page)-1]
		params.After = &last.User.ID
	}
}
