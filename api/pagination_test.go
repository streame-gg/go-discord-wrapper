package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/common"
)

func newPaginationClient(ts *httptest.Server) *RestClient {
	return NewRestClient("test-token",
		options.WithBaseURL(ts.URL),
		options.WithRateLimiting(options.RateLimiterOptions{Disabled: true}),
		options.WithRetry(options.RetryOptions{
			MaxRetries:          0,
			RetryOnRateLimit:    false,
			RetryOnServerErrors: false,
		}),
	)
}

func makeMembersPage(n int, startID int) []*common.GuildMember {
	members := make([]*common.GuildMember, n)
	for i := 0; i < n; i++ {
		id := common.Snowflake(string(rune('0' + startID + i))) // simple unique IDs
		// Use a proper string-based snowflake
		idStr := common.Snowflake(jsonIntStr(startID + i))
		members[i] = &common.GuildMember{
			User: &common.User{ID: idStr},
		}
		_ = id
	}
	return members
}

// jsonIntStr formats an int as a string Snowflake.
func jsonIntStr(n int) common.Snowflake {
	b, _ := json.Marshal(n)
	s := string(b)
	return common.Snowflake(s)
}

func makeMembersPageWithIDs(ids []string) []*common.GuildMember {
	members := make([]*common.GuildMember, len(ids))
	for i, id := range ids {
		members[i] = &common.GuildMember{
			User: &common.User{ID: common.Snowflake(id)},
		}
	}
	return members
}

func makeBanPage(userIDs []string) []*Ban {
	bans := make([]*Ban, len(userIDs))
	for i, id := range userIDs {
		bans[i] = &Ban{
			User: common.User{ID: common.Snowflake(id)},
		}
	}
	return bans
}

func makeMessagePage(ids []string, channelID string) []*common.Message {
	msgs := make([]*common.Message, len(ids))
	for i, id := range ids {
		msgs[i] = &common.Message{
			ID:        common.Snowflake(id),
			ChannelID: common.Snowflake(channelID),
		}
	}
	return msgs
}

func TestFetchAllGuildMembersOnePage(t *testing.T) {
	var reqCount int32

	// Build 5 members
	page := make([]*common.GuildMember, 5)
	for i := range page {
		page[i] = &common.GuildMember{
			User: &common.User{ID: common.Snowflake(string(rune('a' + i)))},
		}
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&reqCount, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(page)
	}))
	defer ts.Close()

	client := newPaginationClient(ts)
	members, err := client.FetchAllGuildMembers(context.Background(), "guild1")

	require.NoError(t, err)
	assert.Len(t, members, 5)
	assert.Equal(t, int32(1), atomic.LoadInt32(&reqCount), "expected single request for one page")
}

func TestFetchAllGuildMembersMultiplePages(t *testing.T) {
	// Build page 1: exactly 1000 members with IDs "1" through "1000"
	page1 := make([]*common.GuildMember, 1000)
	for i := range page1 {
		page1[i] = &common.GuildMember{
			User: &common.User{ID: common.Snowflake(intToSnowflake(i + 1))},
		}
	}
	lastPage1ID := page1[999].User.ID

	// Build page 2: 5 members
	page2 := make([]*common.GuildMember, 5)
	for i := range page2 {
		page2[i] = &common.GuildMember{
			User: &common.User{ID: common.Snowflake(intToSnowflake(1001 + i))},
		}
	}

	var reqCount int32
	var secondRequestAfter string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&reqCount, 1)
		parsed, _ := url.ParseQuery(r.URL.RawQuery)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if n == 1 {
			_ = json.NewEncoder(w).Encode(page1)
		} else {
			secondRequestAfter = parsed.Get("after")
			_ = json.NewEncoder(w).Encode(page2)
		}
	}))
	defer ts.Close()

	client := newPaginationClient(ts)
	members, err := client.FetchAllGuildMembers(context.Background(), "guild1")

	require.NoError(t, err)
	assert.Len(t, members, 1005)
	assert.Equal(t, int32(2), atomic.LoadInt32(&reqCount))
	assert.Equal(t, string(lastPage1ID), secondRequestAfter, "second request should advance cursor past last member of page 1")
}

func intToSnowflake(n int) string {
	b, _ := json.Marshal(n)
	return string(b)
}

func TestFetchAllMessagesMultiplePages(t *testing.T) {
	// Page 1: 100 messages with IDs "100" down to "1" (newest-first)
	page1 := make([]*common.Message, 100)
	for i := range page1 {
		page1[i] = &common.Message{
			ID:        common.Snowflake(intToSnowflake(100 - i)),
			ChannelID: "chan1",
			Author:    &common.User{},
		}
	}
	oldestPage1ID := page1[99].ID // "1"

	// Page 2: 10 messages
	page2 := make([]*common.Message, 10)
	for i := range page2 {
		page2[i] = &common.Message{
			ID:        common.Snowflake(intToSnowflake(-(i + 1))),
			ChannelID: "chan1",
			Author:    &common.User{},
		}
	}

	var reqCount int32
	var secondRequestBefore string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&reqCount, 1)
		parsed, _ := url.ParseQuery(r.URL.RawQuery)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if n == 1 {
			_ = json.NewEncoder(w).Encode(page1)
		} else {
			secondRequestBefore = parsed.Get("before")
			_ = json.NewEncoder(w).Encode(page2)
		}
	}))
	defer ts.Close()

	client := newPaginationClient(ts)
	msgs, err := client.FetchAllMessages(context.Background(), "chan1")

	require.NoError(t, err)
	assert.Len(t, msgs, 110)
	assert.Equal(t, int32(2), atomic.LoadInt32(&reqCount))
	assert.Equal(t, string(oldestPage1ID), secondRequestBefore, "second request should use before=last_id_of_page1")
}

func TestFetchAllGuildBansOnePage(t *testing.T) {
	bans := make([]*Ban, 7)
	for i := range bans {
		bans[i] = &Ban{
			User: common.User{
				ID:            common.Snowflake(intToSnowflake(i + 1)),
				Discriminator: "0",
			},
		}
	}

	var reqCount int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&reqCount, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(bans)
	}))
	defer ts.Close()

	client := newPaginationClient(ts)
	result, err := client.FetchAllGuildBans(context.Background(), "guild1")

	require.NoError(t, err)
	assert.Len(t, result, 7)
	assert.Equal(t, int32(1), atomic.LoadInt32(&reqCount), "expected single page request")

	// verify IDs round-trip
	assert.Equal(t, common.Snowflake("1"), result[0].User.ID)

	// suppress unused import
	_ = time.Second
}
