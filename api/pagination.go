package api

import (
	"context"

	"github.com/streame-gg/go-discord-wrapper/types/common"
	"github.com/streame-gg/go-discord-wrapper/util"
)

// FetchAllGuildMembers fetches every member in a guild across as many pages as
// needed. It uses the after-cursor pattern (max 1000 per request) and stops when
// a page returns fewer members than the page size.
func (c *RestClient) FetchAllGuildMembers(ctx context.Context, guildID common.Snowflake) ([]*common.GuildMember, error) {
	const pageSize = 1000

	var all []*common.GuildMember
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
func (c *RestClient) FetchAllMessages(ctx context.Context, channelID common.Snowflake) ([]*common.Message, error) {
	const pageSize = 100

	var all []*common.Message
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
func (c *RestClient) FetchAllGuildBans(ctx context.Context, guildID common.Snowflake) ([]*Ban, error) {
	const pageSize = 1000

	var all []*Ban
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
