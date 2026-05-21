package api

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func mustSnowflake(s string) discord.Snowflake {
	n, err := strconv.ParseUint(s, 10, 64)
	if err == nil {
		return discord.Snowflake(n)
	}
	h := fnv.New64a()
	h.Write([]byte(s))
	return discord.Snowflake(h.Sum64())
}

func newPaginationClient(ts *httptest.Server) *RestClient {
	rc, err := NewRestClient("test-token",
		options.WithBaseURL(ts.URL),
		options.WithRateLimiting(options.RateLimiterOptions{Disabled: true}),
		options.WithRetry(options.RetryOptions{
			MaxRetries:          0,
			RetryOnRateLimit:    false,
			RetryOnServerErrors: false,
		}),
	)
	if err != nil {
		panic(err)
	}
	return rc
}

func makeMembersPage(n int, startID int) []*discord.GuildMember {
	members := make([]*discord.GuildMember, n)
	for i := 0; i < n; i++ {
		members[i] = &discord.GuildMember{
			User: &discord.User{ID: discord.Snowflake(jsonIntStr(startID + i))},
		}
	}
	return members
}

// jsonIntStr converts an int to a Snowflake.
func jsonIntStr(n int) discord.Snowflake {
	return discord.Snowflake(uint64(n))
}

func makeMembersPageWithIDs(ids []string) []*discord.GuildMember {
	members := make([]*discord.GuildMember, len(ids))
	for i, id := range ids {
		members[i] = &discord.GuildMember{
			User: &discord.User{ID: mustSnowflake(id)},
		}
	}
	return members
}

func makeBanPage(userIDs []string) []*discord.Ban {
	bans := make([]*discord.Ban, len(userIDs))
	for i, id := range userIDs {
		bans[i] = &discord.Ban{
			User: discord.User{ID: mustSnowflake(id)},
		}
	}
	return bans
}

func makeMessagePage(ids []string, channelID string) []*discord.Message {
	msgs := make([]*discord.Message, len(ids))
	for i, id := range ids {
		msgs[i] = &discord.Message{
			ID:        mustSnowflake(id),
			ChannelID: mustSnowflake(channelID),
		}
	}
	return msgs
}

func TestFetchAllGuildMembersOnePage(t *testing.T) {
	var reqCount int32

	// Build 5 members
	page := make([]*discord.GuildMember, 5)
	for i := range page {
		page[i] = &discord.GuildMember{
			User: &discord.User{ID: mustSnowflake(string(rune('a' + i)))},
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
	members, err := client.FetchAllGuildMembers(context.Background(), discord.Snowflake(1234567890123456789))

	require.NoError(t, err)
	assert.Len(t, members, 5)
	assert.Equal(t, int32(1), atomic.LoadInt32(&reqCount), "expected single request for one page")
}

func TestFetchAllGuildMembersMultiplePages(t *testing.T) {
	// Build page 1: exactly 1000 members with IDs "1" through "1000"
	page1 := make([]*discord.GuildMember, 1000)
	for i := range page1 {
		page1[i] = &discord.GuildMember{
			User: &discord.User{ID: discord.Snowflake(uint64(i + 1))},
		}
	}
	lastPage1ID := page1[999].User.ID

	// Build page 2: 5 members
	page2 := make([]*discord.GuildMember, 5)
	for i := range page2 {
		page2[i] = &discord.GuildMember{
			User: &discord.User{ID: discord.Snowflake(uint64(1001 + i))},
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
	members, err := client.FetchAllGuildMembers(context.Background(), discord.Snowflake(1234567890123456789))

	require.NoError(t, err)
	assert.Len(t, members, 1005)
	assert.Equal(t, int32(2), atomic.LoadInt32(&reqCount))
	assert.Equal(t, lastPage1ID.String(), secondRequestAfter, "second request should advance cursor past last member of page 1")
}

func intToSnowflake(n int) string {
	b, _ := json.Marshal(n)
	return string(b)
}

func TestFetchAllMessagesMultiplePages(t *testing.T) {
	// Page 1: 100 messages with IDs "100" down to "1" (newest-first)
	page1 := make([]*discord.Message, 100)
	for i := range page1 {
		page1[i] = &discord.Message{
			ID:        discord.Snowflake(uint64(100 - i)),
			ChannelID: discord.Snowflake(1234567890123456789),
			Author:    &discord.User{},
		}
	}
	oldestPage1ID := page1[99].ID // Snowflake(1)

	// Page 2: 10 messages
	page2 := make([]*discord.Message, 10)
	for i := range page2 {
		page2[i] = &discord.Message{
			ID:        mustSnowflake(intToSnowflake(-(i + 1))),
			ChannelID: discord.Snowflake(1234567890123456789),
			Author:    &discord.User{},
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
	msgs, err := client.FetchAllMessages(context.Background(), discord.Snowflake(1234567890123456789))

	require.NoError(t, err)
	assert.Len(t, msgs, 110)
	assert.Equal(t, int32(2), atomic.LoadInt32(&reqCount))
	assert.Equal(t, oldestPage1ID.String(), secondRequestBefore, "second request should use before=last_id_of_page1")
}

func TestFetchAllGuildBansOnePage(t *testing.T) {
	bans := make([]*discord.Ban, 7)
	for i := range bans {
		bans[i] = &discord.Ban{
			User: discord.User{
				ID:            discord.Snowflake(uint64(i + 1)),
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
	result, err := client.FetchAllGuildBans(context.Background(), discord.Snowflake(1234567890123456789))

	require.NoError(t, err)
	assert.Len(t, result, 7)
	assert.Equal(t, int32(1), atomic.LoadInt32(&reqCount), "expected single page request")

	// verify IDs round-trip
	first := result[0]
	assert.Equal(t, discord.Snowflake(1), first.User.ID)

	// suppress unused import
	_ = time.Second
}

// ── FetchAllAuditLogEntries ───────────────────────────────────────────────────

func makeAuditLogPage(ids []string) *discord.AuditLog {
	entries := make([]discord.AuditLogEntry, len(ids))
	for i, id := range ids {
		entries[i] = discord.AuditLogEntry{ID: mustSnowflake(id)}
	}
	return &discord.AuditLog{AuditLogEntries: entries}
}

func TestFetchAllAuditLogEntriesOnePage(t *testing.T) {
	page := makeAuditLogPage([]string{"10", "9", "8", "7", "6"})

	var reqCount int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&reqCount, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(page)
	}))
	defer ts.Close()

	client := newPaginationClient(ts)
	entries, err := client.FetchAllAuditLogEntries(context.Background(), discord.Snowflake(1234567890123456789), GetGuildAuditLogParams{})

	require.NoError(t, err)
	assert.Len(t, entries, 5)
	assert.Equal(t, int32(1), atomic.LoadInt32(&reqCount))
}

func TestFetchAllAuditLogEntriesMultiplePages(t *testing.T) {
	// Build 100 entries (page 1) + 3 entries (page 2)
	ids1 := make([]string, 100)
	for i := range ids1 {
		ids1[i] = intToSnowflake(1000 - i)
	}
	page1 := makeAuditLogPage(ids1)
	oldestID1 := ids1[99]

	ids2 := []string{"901", "902", "903"}
	page2 := makeAuditLogPage(ids2)

	var reqCount int32
	var secondBefore string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&reqCount, 1)
		parsed, _ := url.ParseQuery(r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if n == 1 {
			_ = json.NewEncoder(w).Encode(page1)
		} else {
			secondBefore = parsed.Get("before")
			_ = json.NewEncoder(w).Encode(page2)
		}
	}))
	defer ts.Close()

	client := newPaginationClient(ts)
	entries, err := client.FetchAllAuditLogEntries(context.Background(), discord.Snowflake(1234567890123456789), GetGuildAuditLogParams{})

	require.NoError(t, err)
	assert.Len(t, entries, 103)
	assert.Equal(t, int32(2), atomic.LoadInt32(&reqCount))
	assert.Equal(t, oldestID1, secondBefore, "second request should use before=oldest entry on page 1")
}

// ── FetchAllEntitlements ──────────────────────────────────────────────────────

func makeEntitlementPage(ids []string) []*discord.Entitlement {
	page := make([]*discord.Entitlement, len(ids))
	for i, id := range ids {
		page[i] = &discord.Entitlement{ID: mustSnowflake(id)}
	}
	return page
}

func TestFetchAllEntitlementsOnePage(t *testing.T) {
	page := makeEntitlementPage([]string{"1", "2", "3"})

	var reqCount int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&reqCount, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(page)
	}))
	defer ts.Close()

	client := newPaginationClient(ts)
	result, err := client.FetchAllEntitlements(context.Background(), discord.Snowflake(1234567890123456789), ListEntitlementsParams{})

	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, int32(1), atomic.LoadInt32(&reqCount))
}

func TestFetchAllEntitlementsMultiplePages(t *testing.T) {
	ids1 := make([]string, 100)
	for i := range ids1 {
		ids1[i] = intToSnowflake(i + 1)
	}
	lastID1 := ids1[99]
	page1 := makeEntitlementPage(ids1)
	page2 := makeEntitlementPage([]string{"101", "102"})

	var reqCount int32
	var secondAfter string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&reqCount, 1)
		parsed, _ := url.ParseQuery(r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if n == 1 {
			_ = json.NewEncoder(w).Encode(page1)
		} else {
			secondAfter = parsed.Get("after")
			_ = json.NewEncoder(w).Encode(page2)
		}
	}))
	defer ts.Close()

	client := newPaginationClient(ts)
	result, err := client.FetchAllEntitlements(context.Background(), discord.Snowflake(1234567890123456789), ListEntitlementsParams{})

	require.NoError(t, err)
	assert.Len(t, result, 102)
	assert.Equal(t, int32(2), atomic.LoadInt32(&reqCount))
	assert.Equal(t, lastID1, secondAfter, "second request should advance cursor past last entitlement on page 1")
}

// ── FetchAllScheduledEventUsers ───────────────────────────────────────────────

func makeScheduledEventUsersPage(userIDs []string) []*discord.GuildScheduledEventUser {
	page := make([]*discord.GuildScheduledEventUser, len(userIDs))
	for i, id := range userIDs {
		page[i] = &discord.GuildScheduledEventUser{
			User: discord.User{ID: mustSnowflake(id)},
		}
	}
	return page
}

func TestFetchAllScheduledEventUsersOnePage(t *testing.T) {
	page := makeScheduledEventUsersPage([]string{"u1", "u2", "u3"})

	var reqCount int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&reqCount, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(page)
	}))
	defer ts.Close()

	client := newPaginationClient(ts)
	result, err := client.FetchAllScheduledEventUsers(context.Background(), discord.Snowflake(1234567890123456789), discord.Snowflake(1234567890123456789), false)

	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, int32(1), atomic.LoadInt32(&reqCount))
}

func TestFetchAllScheduledEventUsersMultiplePages(t *testing.T) {
	ids1 := make([]string, 100)
	for i := range ids1 {
		ids1[i] = intToSnowflake(i + 1)
	}
	lastID1 := ids1[99]
	page1 := makeScheduledEventUsersPage(ids1)
	page2 := makeScheduledEventUsersPage([]string{"101", "102"})

	var reqCount int32
	var secondAfter string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&reqCount, 1)
		parsed, _ := url.ParseQuery(r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if n == 1 {
			_ = json.NewEncoder(w).Encode(page1)
		} else {
			secondAfter = parsed.Get("after")
			_ = json.NewEncoder(w).Encode(page2)
		}
	}))
	defer ts.Close()

	client := newPaginationClient(ts)
	result, err := client.FetchAllScheduledEventUsers(context.Background(), discord.Snowflake(1234567890123456789), discord.Snowflake(1234567890123456789), false)

	require.NoError(t, err)
	assert.Len(t, result, 102)
	assert.Equal(t, int32(2), atomic.LoadInt32(&reqCount))
	assert.Equal(t, lastID1, secondAfter)
}
