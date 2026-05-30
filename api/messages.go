package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Request / query-param types ───────────────────────────────────────────────

// GetMessagesParams are query parameters for listing channel messages.
// At most one of Around / Before / After may be set.
type GetMessagesParams struct {
	Around *discord.Snowflake
	Before *discord.Snowflake
	After  *discord.Snowflake
	Limit  *int // 1–100, default 50
}

func (p GetMessagesParams) toQuery() string {
	q := url.Values{}
	if p.Around != nil {
		q.Set("around", p.Around.String())
	}
	if p.Before != nil {
		q.Set("before", p.Before.String())
	}
	if p.After != nil {
		q.Set("after", p.After.String())
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// basenameFilename returns the basename of a filename, stripping any path
// components. Used to ensure local filesystem paths don't leak into Discord
// uploads via the Content-Disposition header.
func basenameFilename(name string) string {
	if i := strings.LastIndexAny(name, `/\`); i >= 0 {
		name = name[i+1:]
	}
	if name == "" {
		return "file"
	}
	return name
}

// buildMultipartMessage encodes payload as multipart/form-data with optional file parts.
// Returns the body buffer and the Content-Type header value (including boundary).
func buildMultipartMessage(payload []byte, files []discord.MessageFile) (*bytes.Buffer, string, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	// Write the JSON payload as the "payload_json" part.
	jsonPart, err := w.CreatePart(textproto.MIMEHeader{
		"Content-Disposition": []string{`form-data; name="payload_json"`},
		"Content-Type":        []string{"application/json"},
	})
	if err != nil {
		return nil, "", err
	}
	if _, err := jsonPart.Write(payload); err != nil {
		return nil, "", err
	}

	for i, f := range files {
		ct := f.ContentType
		if ct == "" {
			ct = "application/octet-stream"
		}
		h := textproto.MIMEHeader{
			"Content-Disposition": []string{
				mime.FormatMediaType("form-data", map[string]string{
					"name":     fmt.Sprintf("files[%d]", i),
					"filename": basenameFilename(f.Name),
				}),
			},
			"Content-Type": []string{ct},
		}
		part, err := w.CreatePart(h)
		if err != nil {
			return nil, "", err
		}
		if _, err := part.Write(f.Data); err != nil {
			return nil, "", err
		}
	}
	if err := w.Close(); err != nil {
		return nil, "", err
	}
	return &buf, w.FormDataContentType(), nil
}

// CreateMessageParams are the parameters for sending a new message.
type CreateMessageParams struct {
	Content          string                           `json:"content,omitempty"`
	TTS              bool                             `json:"tts,omitempty"`
	Embeds           []discord.Embed                  `json:"embeds,omitempty"`
	AllowedMentions  *discord.AllowedMentions         `json:"allowed_mentions,omitempty"`
	MessageReference *discord.MessageMessageReference `json:"message_reference,omitempty"`
	// Components holds message components (action rows, buttons, etc.).
	Components   []discord.AnyComponent `json:"-"`
	StickerIDs   []discord.Snowflake    `json:"sticker_ids,omitempty"`
	Flags        discord.MessageFlag    `json:"flags,omitempty"`
	Nonce        interface{}            `json:"nonce,omitempty"`
	EnforceNonce bool                   `json:"enforce_nonce,omitempty"`
	// Files are binary attachments sent via multipart/form-data.
	// When set, the request is encoded as multipart rather than JSON.
	Files             []discord.MessageFile     `json:"-"`
	Poll              discord.PollRequest       `json:"poll_request,omitempty"`
	SharedClientTheme discord.SharedClientTheme `json:"shared_client_theme,omitempty"`
}

func (p CreateMessageParams) MarshalJSON() ([]byte, error) {
	type Alias CreateMessageParams
	return json.Marshal(&struct {
		Components []discord.AnyComponent `json:"components,omitempty"`
		Alias
	}{
		Components: p.Components,
		Alias:      (Alias)(p),
	})
}

// EditMessageParams are the parameters for editing an existing message.
// Nil pointer fields are omitted from the request body (the field is left unchanged).
// Set Content to a pointer to "" to clear the message content.
type EditMessageParams struct {
	Content         *string                  `json:"content,omitempty"`
	Embeds          *[]discord.Embed         `json:"embeds,omitempty"`
	Flags           *discord.MessageFlag     `json:"flags,omitempty"`
	AllowedMentions *discord.AllowedMentions `json:"allowed_mentions,omitempty"`
	// Components replaces the message's component list.
	// Set to an empty (non-nil) slice to remove all components.
	Components  []discord.AnyComponent `json:"-"`
	Attachments *[]discord.Attachment  `json:"attachments,omitempty"`
	// Files are binary attachments added via multipart/form-data.
	// When set, the request is encoded as multipart rather than JSON.
	Files []discord.MessageFile `json:"-"`
}

func (p EditMessageParams) MarshalJSON() ([]byte, error) {
	type Alias EditMessageParams
	return json.Marshal(&struct {
		Components []discord.AnyComponent `json:"components,omitempty"`
		Alias
	}{
		Components: p.Components,
		Alias:      (Alias)(p),
	})
}

// GetReactionsParams are optional query parameters for listing reaction users.
type GetReactionsParams struct {
	After *discord.Snowflake
	Limit *int // 1–100, default 25
	Type  *int // 0 = normal, 1 = burst (super-reaction)
}

func (p GetReactionsParams) toQuery() string {
	q := url.Values{}
	if p.After != nil {
		q.Set("after", p.After.String())
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if p.Type != nil {
		q.Set("type", strconv.Itoa(*p.Type))
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// encodeEmoji URL-encodes an emoji for use in a reaction path segment.
// Pass a raw Unicode character (e.g. "👍") or a custom emoji as "name:id".
func encodeEmoji(emoji string) string {
	return url.PathEscape(emoji)
}

// ── Message endpoints ─────────────────────────────────────────────────────────

// GetMessages returns up to 100 messages from a channel.
func (c *RestClient) GetMessages(ctx context.Context, channelID discord.Snowflake, params GetMessagesParams) ([]*discord.Message, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/messages" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.Message](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetMessage returns a single message by ID.
func (c *RestClient) GetMessage(ctx context.Context, channelID, messageID discord.Snowflake) (*discord.Message, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	if err := messageID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/messages/" + messageID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Message](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// CreateMessage sends a new message to a channel.
// When params.Files is non-empty the request is sent as multipart/form-data.
func (c *RestClient) CreateMessage(ctx context.Context, channelID discord.Snowflake, params CreateMessageParams) (*discord.Message, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/messages"

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	if len(params.Files) > 0 {
		buf, ct, err := buildMultipartMessage(jsonBody, params.Files)
		if err != nil {
			return nil, err
		}
		req, err = c.generateRequest(ctx, http.MethodPost, path, buf, c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", ct)
	} else {
		req, err = c.generateRequest(ctx, http.MethodPost, path, bytes.NewReader(jsonBody), c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
	}

	return doRequest[discord.Message](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// EditMessage edits a previously sent message. Only fields set in params are changed.
// When params.Files is non-empty the request is sent as multipart/form-data.
func (c *RestClient) EditMessage(ctx context.Context, channelID, messageID discord.Snowflake, params EditMessageParams) (*discord.Message, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	if err := messageID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/messages/" + messageID.String()

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	if len(params.Files) > 0 {
		buf, ct, err := buildMultipartMessage(jsonBody, params.Files)
		if err != nil {
			return nil, err
		}
		req, err = c.generateRequest(ctx, http.MethodPatch, path, buf, c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", ct)
	} else {
		req, err = c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(jsonBody), c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
	}

	return doRequest[discord.Message](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteMessage deletes a message.
func (c *RestClient) DeleteMessage(ctx context.Context, channelID, messageID discord.Snowflake, opts *DeleteMessageOptions) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	if err := messageID.Validate(); err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/messages/" + messageID.String()

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
		if err != nil {
			return err
		}
		return doRequestWithoutResponse(c, req)
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// BulkDeleteMessages deletes 2–100 messages at once.
// Messages older than 14 days cannot be bulk-deleted and will cause a Discord API error.
func (c *RestClient) BulkDeleteMessages(ctx context.Context, channelID discord.Snowflake, messageIDs []discord.Snowflake, opts *BulkDeleteMessagesOptions) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	for _, id := range messageIDs {
		if err := id.Validate(); err != nil {
			return err
		}
	}

	if len(messageIDs) < 2 || len(messageIDs) > 100 {
		return fmt.Errorf("bulk delete requires 2–100 message IDs, got %d", len(messageIDs))
	}

	body, err := json.Marshal(map[string]any{"messages": messageIDs})
	if err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/messages/bulk-delete"

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPost, path, bytes.NewReader(body), c.WithBotAuthorization())
		if err != nil {
			return err
		}
		return doRequestWithoutResponse(c, req)
	}

	req, err := c.generateRequest(ctx, http.MethodPost, path, bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// CrosspostMessage publishes a message in an announcement channel to all following channels.
func (c *RestClient) CrosspostMessage(ctx context.Context, channelID, messageID discord.Snowflake) (*discord.Message, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	if err := messageID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/messages/" + messageID.String() + "/crosspost"
	req, err := c.generateRequest(ctx, http.MethodPost, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Message](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ── Options types ─────────────────────────────────────────────────────────────

type DeleteMessageOptions struct {
	Reason string
}

type BulkDeleteMessagesOptions struct {
	Reason string
}

// ── Pin endpoints ─────────────────────────────────────────────────────────────

// GetPinnedMessages returns all pinned messages in a channel (max 50).
func (c *RestClient) GetPinnedMessages(ctx context.Context, channelID discord.Snowflake) ([]*discord.Message, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/channels/"+channelID.String()+"/pins", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.Message](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// PinMessage pins a message in a channel. Requires MANAGE_MESSAGES.
func (c *RestClient) PinMessage(ctx context.Context, channelID, messageID discord.Snowflake, opts *PinMessageOptions) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	if err := messageID.Validate(); err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/pins/" + messageID.String()

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPut, path, nil, c.WithBotAuthorization())
		if err != nil {
			return err
		}
		return doRequestWithoutResponse(c, req)
	}

	req, err := c.generateRequest(ctx, http.MethodPut, path, nil, c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// UnpinMessage unpins a message from a channel. Requires MANAGE_MESSAGES.
func (c *RestClient) UnpinMessage(ctx context.Context, channelID, messageID discord.Snowflake, opts *UnpinMessageOptions) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	if err := messageID.Validate(); err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/pins/" + messageID.String()

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
		if err != nil {
			return err
		}
		return doRequestWithoutResponse(c, req)
	}

	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// ── Reaction endpoints ────────────────────────────────────────────────────────

// AddReaction adds a reaction to a message.
// emoji is a raw Unicode character (e.g. "👍") or a custom emoji in "name:id" form.
func (c *RestClient) AddReaction(ctx context.Context, channelID, messageID discord.Snowflake, emoji string) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	if err := messageID.Validate(); err != nil {
		return err
	}

	path := fmt.Sprintf("/channels/%s/messages/%s/reactions/%s/@me", channelID.String(), messageID.String(), encodeEmoji(emoji))
	req, err := c.generateRequest(ctx, http.MethodPut, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// DeleteOwnReaction removes the bot's own reaction from a message.
func (c *RestClient) DeleteOwnReaction(ctx context.Context, channelID, messageID discord.Snowflake, emoji string) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	if err := messageID.Validate(); err != nil {
		return err
	}

	path := fmt.Sprintf("/channels/%s/messages/%s/reactions/%s/@me", channelID.String(), messageID.String(), encodeEmoji(emoji))
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// DeleteUserReaction removes another user's reaction from a message. Requires MANAGE_MESSAGES.
func (c *RestClient) DeleteUserReaction(ctx context.Context, channelID, messageID discord.Snowflake, emoji string, userID discord.Snowflake) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	if err := messageID.Validate(); err != nil {
		return err
	}

	if err := userID.Validate(); err != nil {
		return err
	}

	path := fmt.Sprintf("/channels/%s/messages/%s/reactions/%s/%s", channelID.String(), messageID.String(), encodeEmoji(emoji), &userID)
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// GetReactions returns the users who reacted to a message with the given emoji.
func (c *RestClient) GetReactions(ctx context.Context, channelID, messageID discord.Snowflake, emoji string, params GetReactionsParams) ([]*discord.User, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	if err := messageID.Validate(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/channels/%s/messages/%s/reactions/%s%s", channelID.String(), messageID.String(), encodeEmoji(emoji), params.toQuery())
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequestSlice[discord.User](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteAllReactions removes every reaction from a message. Requires MANAGE_MESSAGES.
func (c *RestClient) DeleteAllReactions(ctx context.Context, channelID, messageID discord.Snowflake) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	if err := messageID.Validate(); err != nil {
		return err
	}

	path := fmt.Sprintf("/channels/%s/messages/%s/reactions", channelID.String(), messageID.String())
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// DeleteAllReactionsForEmoji removes all reactions for a specific emoji. Requires MANAGE_MESSAGES.
func (c *RestClient) DeleteAllReactionsForEmoji(ctx context.Context, channelID, messageID discord.Snowflake, emoji string) error {
	if err := channelID.Validate(); err != nil {
		return err
	}

	if err := messageID.Validate(); err != nil {
		return err
	}

	path := fmt.Sprintf("/channels/%s/messages/%s/reactions/%s", channelID.String(), messageID.String(), encodeEmoji(emoji))
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// SearchGuildMessagesParams holds the query parameters for searching guild messages.
type SearchGuildMessagesParams struct {
	Content        *string
	AuthorID       []discord.Snowflake
	ChannelID      []discord.Snowflake
	MentionsUserID []discord.Snowflake
	MentionsRoleID []discord.Snowflake
	Has            []string
	MinID          *discord.Snowflake
	MaxID          *discord.Snowflake
	Pinned         *bool
	IncludeNSFW    *bool
	Limit          *int
	Offset         *int
}

func (p SearchGuildMessagesParams) toQuery() string {
	q := url.Values{}
	if p.Content != nil {
		q.Set("content", *p.Content)
	}
	for _, id := range p.AuthorID {
		q.Add("author_id", id.String())
	}
	for _, id := range p.ChannelID {
		q.Add("channel_id", id.String())
	}
	for _, id := range p.MentionsUserID {
		q.Add("mentions", id.String())
	}
	for _, id := range p.MentionsRoleID {
		q.Add("mentions_role_id", id.String())
	}
	for _, h := range p.Has {
		q.Add("has", h)
	}
	if p.MinID != nil {
		q.Set("min_id", p.MinID.String())
	}
	if p.MaxID != nil {
		q.Set("max_id", p.MaxID.String())
	}
	if p.Pinned != nil {
		q.Set("pinned", strconv.FormatBool(*p.Pinned))
	}
	if p.IncludeNSFW != nil {
		q.Set("include_nsfw", strconv.FormatBool(*p.IncludeNSFW))
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if p.Offset != nil {
		q.Set("offset", strconv.Itoa(*p.Offset))
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// SearchGuildMessages searches for messages in a guild. Returns 202 while the index is still
// being built (result will be empty). Requires READ_MESSAGE_HISTORY.
func (c *RestClient) SearchGuildMessages(ctx context.Context, guildID discord.Snowflake, params SearchGuildMessagesParams) (*discord.GuildSearchResponse, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/messages/search" + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.GuildSearchResponse](c, req, map[int]bool{
		http.StatusOK:       true,
		http.StatusAccepted: true,
	})
}
