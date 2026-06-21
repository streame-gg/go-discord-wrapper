package api

import (
	"context"
	"fmt"

	"github.com/streame-gg/go-discord-wrapper/internal/util"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
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

		last := page[len(page)-1]
		if last.User == nil {
			return nil, fmt.Errorf("go-discord-wrapper: guild member returned by API has nil User; cannot paginate")
		}
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
		page, err := c.ListMessages(ctx, channelID, params)
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
		page, err := c.ListGuildBans(ctx, guildID, params)
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
		page, err := c.ListGuildScheduledEventUsers(ctx, guildID, eventID, params)
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

// FetchAllReactions fetches every user who reacted to a message with the given
// emoji, paginating forward by user ID (max 100 per request). The filter's Type
// (normal vs burst) is preserved; it's After/Limit are managed internally.
func (c *RestClient) FetchAllReactions(ctx context.Context, channelID, messageID discord.Snowflake, emoji string, filter GetReactionsParams) ([]*discord.User, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}
	if err := messageID.Validate(); err != nil {
		return nil, err
	}

	const pageSize = 100

	filter.Limit = util.PointerOf(pageSize)
	filter.After = nil

	var all []*discord.User

	for {
		page, err := c.ListReactions(ctx, channelID, messageID, emoji, filter)
		if err != nil {
			return nil, err
		}

		all = append(all, page...)

		if len(page) < pageSize {
			return all, nil
		}

		filter.After = &page[len(page)-1].ID
	}
}

// FetchAllPollAnswerVoters fetches every user who voted for a poll answer,
// paginating forward by user ID (max 100 per request).
func (c *RestClient) FetchAllPollAnswerVoters(ctx context.Context, channelID, messageID discord.Snowflake, answerID int) ([]*discord.User, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}
	if err := messageID.Validate(); err != nil {
		return nil, err
	}

	const pageSize = 100

	params := GetPollAnswerVotersParams{Limit: util.PointerOf(pageSize)}

	var all []*discord.User

	for {
		resp, err := c.GetPollAnswerVoters(ctx, channelID, messageID, answerID, params)
		if err != nil {
			return nil, err
		}

		all = append(all, resp.Users...)

		if len(resp.Users) < pageSize {
			return all, nil
		}

		params.After = &resp.Users[len(resp.Users)-1].ID
	}
}

// FetchAllCurrentUserGuilds fetches every guild the current user/bot is in,
// paginating forward by guild ID (max 200 per request). WithCounts and
// UserToken from the filter are preserved; Before/After/Limit are managed
// internally.
func (c *RestClient) FetchAllCurrentUserGuilds(ctx context.Context, filter GetCurrentUserGuildsParams) ([]*discord.CurrentUserGuild, error) {
	const pageSize = 200

	filter.Limit = util.PointerOf(pageSize)
	filter.Before = nil
	filter.After = nil

	var all []*discord.CurrentUserGuild

	for {
		page, err := c.ListCurrentUserGuilds(ctx, filter)
		if err != nil {
			return nil, err
		}

		all = append(all, page...)

		if len(page) < pageSize {
			return all, nil
		}

		filter.After = &page[len(page)-1].ID
	}
}

// FetchAllThreadMembers fetches every member of a thread, paginating forward by
// user ID (max 100 per request). It always requests the full member object,
// which is what makes the endpoint paginate.
func (c *RestClient) FetchAllThreadMembers(ctx context.Context, channelID discord.Snowflake) ([]*discord.ThreadMember, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	const pageSize = 100

	params := ListThreadMembersParams{
		WithMember: util.PointerOf(true),
		Limit:      util.PointerOf(pageSize),
	}

	var all []*discord.ThreadMember

	for {
		page, err := c.ListThreadMembers(ctx, channelID, params)
		if err != nil {
			return nil, err
		}

		all = append(all, page...)

		if len(page) < pageSize {
			return all, nil
		}

		last := page[len(page)-1]
		if last.UserID.IsEmpty() {
			return nil, fmt.Errorf("go-discord-wrapper: thread member returned by API has nil UserID; cannot paginate")
		}
		params.After = &last.UserID
	}
}

// FetchAllSKUSubscriptions fetches every subscription for an SKU, paginating
// forward by subscription ID (max 100 per request). The UserID filter is
// preserved; Before/After/Limit are managed internally.
func (c *RestClient) FetchAllSKUSubscriptions(ctx context.Context, skuID discord.Snowflake, filter ListSKUSubscriptionsParams) ([]*discord.Subscription, error) {
	if err := skuID.Validate(); err != nil {
		return nil, err
	}

	const pageSize = 100

	filter.Limit = util.PointerOf(pageSize)
	filter.Before = nil
	filter.After = nil

	var all []*discord.Subscription

	for {
		page, err := c.ListSKUSubscriptions(ctx, skuID, filter)
		if err != nil {
			return nil, err
		}

		all = append(all, page...)

		if len(page) < pageSize {
			return all, nil
		}

		filter.After = &page[len(page)-1].ID
	}
}

// FetchAllPublicArchivedThreads fetches every public archived thread in a
// channel, paginating backwards by archive timestamp until has_more is false.
func (c *RestClient) FetchAllPublicArchivedThreads(ctx context.Context, channelID discord.Snowflake) ([]*discord.Channel, error) {
	return c.fetchAllArchivedThreads(ctx, channelID, c.ListPublicArchivedThreads)
}

// FetchAllPrivateArchivedThreads fetches every private archived thread in a
// channel, paginating backwards by archive timestamp until has_more is false.
func (c *RestClient) FetchAllPrivateArchivedThreads(ctx context.Context, channelID discord.Snowflake) ([]*discord.Channel, error) {
	return c.fetchAllArchivedThreads(ctx, channelID, c.ListPrivateArchivedThreads)
}

// FetchAllJoinedPrivateArchivedThreads fetches every private archived thread the
// current user has joined in a channel. Unlike the public/private listings, this
// endpoint paginates backwards by thread ID (snowflake), so it cursors on
// ListArchivedThreadsParams.BeforeID rather than the archive timestamp.
func (c *RestClient) FetchAllJoinedPrivateArchivedThreads(ctx context.Context, channelID discord.Snowflake) ([]*discord.Channel, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	const pageSize = 100

	params := ListArchivedThreadsParams{Limit: util.PointerOf(pageSize)}

	var all []*discord.Channel

	for {
		page, err := c.ListJoinedPrivateArchivedThreads(ctx, channelID, params)
		if err != nil {
			return nil, err
		}

		all = append(all, page.Threads...)

		if !page.HasMore || len(page.Threads) == 0 {
			return all, nil
		}

		last := page.Threads[len(page.Threads)-1]
		params.BeforeID = &last.ID
	}
}

// fetchAllArchivedThreads walks the archive-timestamp cursor shared by the
// public and private archived-thread listings. Joined private archived threads
// paginate by thread ID instead; see FetchAllJoinedPrivateArchivedThreads.
func (c *RestClient) fetchAllArchivedThreads(ctx context.Context, channelID discord.Snowflake, list func(context.Context, discord.Snowflake, ListArchivedThreadsParams) (*ArchivedThreadsResponse, error)) ([]*discord.Channel, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	const pageSize = 100

	params := ListArchivedThreadsParams{Limit: util.PointerOf(pageSize)}

	var all []*discord.Channel

	for {
		page, err := list(ctx, channelID, params)
		if err != nil {
			return nil, err
		}

		all = append(all, page.Threads...)

		if !page.HasMore || len(page.Threads) == 0 {
			return all, nil
		}

		last := page.Threads[len(page.Threads)-1]
		if last.ThreadMetadata == nil {
			return nil, fmt.Errorf("go-discord-wrapper: archived thread returned by API has nil ThreadMetadata; cannot paginate")
		}
		ts := last.ThreadMetadata.ArchiveTimestamp
		params.Before = &ts
	}
}
