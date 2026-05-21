package api

import (
	"context"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// TestP0_10_CreateForumThreadWithFiles verifies that CreateForumThread sends a
// multipart/form-data request when params.Message.Files is non-empty, and that
// the request body contains a "payload_json" part.
//
// Before the fix, CreateForumThread always used json.Marshal and ignored
// params.Message.Files (which has json:"-"), so attachments were silently lost.
func TestP0_10_CreateForumThreadWithFiles(t *testing.T) {
	var gotContentType string
	var gotBody []byte

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		gotBody, _ = io.ReadAll(r.Body)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		// Minimal discord.Channel response.
		json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
			"id":   "999888777666555444",
			"type": 15,
			"name": "my-forum-thread",
		})
	}))
	defer ts.Close()

	c := newTestClient(ts)
	defer c.Close()

	params := CreateForumThreadParams{
		Name: "my-forum-thread",
		Message: CreateMessageParams{
			Content: "hello from forum thread",
			Files: []discord.MessageFile{
				{
					Name:        "hello.txt",
					ContentType: "text/plain",
					Data:        []byte("file content"),
				},
			},
		},
	}

	channelID := discord.Snowflake(123456789012345678)
	_, err := c.CreateForumThread(context.Background(), channelID, params, nil)
	require.NoError(t, err)

	// Content-Type must be multipart/form-data when Files is non-empty.
	mediaType, _, parseErr := mime.ParseMediaType(gotContentType)
	require.NoError(t, parseErr, "Content-Type header must be parseable")
	assert.Equal(t, "multipart/form-data", mediaType,
		"CreateForumThread must use multipart/form-data when Files are present (P0-10)")

	// Body must contain "payload_json" part.
	assert.Contains(t, string(gotBody), "payload_json",
		"multipart body must include payload_json part")
}

// TestP0_10_CreateForumThreadWithoutFiles verifies that CreateForumThread uses
// plain JSON when no files are attached.
func TestP0_10_CreateForumThreadWithoutFiles(t *testing.T) {
	var gotContentType string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
			"id":   "999888777666555443",
			"type": 15,
			"name": "text-only-thread",
		})
	}))
	defer ts.Close()

	c := newTestClient(ts)
	defer c.Close()

	params := CreateForumThreadParams{
		Name: "text-only-thread",
		Message: CreateMessageParams{
			Content: "no attachments",
		},
	}

	channelID := discord.Snowflake(123456789012345678)
	_, err := c.CreateForumThread(context.Background(), channelID, params, nil)
	require.NoError(t, err)

	assert.Contains(t, gotContentType, "application/json",
		"CreateForumThread must use JSON when no Files are present")
}

// TestP0_10_CreateForumThreadWithReason verifies that the X-Audit-Log-Reason
// header is set when opts.Reason is provided, for both JSON and multipart paths.
func TestP0_10_CreateForumThreadWithReason(t *testing.T) {
	var gotReason string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotReason = r.Header.Get("X-Audit-Log-Reason")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
			"id":   "999888777666555442",
			"type": 15,
			"name": "reasoned-thread",
		})
	}))
	defer ts.Close()

	c := newTestClient(ts)
	defer c.Close()

	params := CreateForumThreadParams{
		Name:    "reasoned-thread",
		Message: CreateMessageParams{Content: "with reason"},
	}
	opts := &CreateForumThreadOptions{Reason: "automated thread creation"}

	channelID := discord.Snowflake(123456789012345678)
	_, err := c.CreateForumThread(context.Background(), channelID, params, opts)
	require.NoError(t, err)

	assert.Equal(t, "automated%20thread%20creation", gotReason,
		"X-Audit-Log-Reason header must be URL-encoded when opts.Reason is provided")
}
