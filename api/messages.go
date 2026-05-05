package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// ── Request / query-param types ───────────────────────────────────────────────

// GetMessagesParams are query parameters for listing channel messages.
// At most one of Around / Before / After may be set.
type GetMessagesParams struct {
	Around *common.Snowflake
	Before *common.Snowflake
	After  *common.Snowflake
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

// CreateMessageParams are the parameters for sending a new message.
type CreateMessageParams struct {
	Content          string                          `json:"content,omitempty"`
	TTS              bool                            `json:"tts,omitempty"`
	Embeds           []common.Embed                  `json:"embeds,omitempty"`
	AllowedMentions  *common.AllowedMentions         `json:"allowed_mentions,omitempty"`
	MessageReference *common.MessageMessageReference `json:"message_reference,omitempty"`
	// Components holds message components (action rows, buttons, etc.).
	Components   []common.AnyComponent `json:"-"`
	StickerIDs   []common.Snowflake    `json:"sticker_ids,omitempty"`
	Flags        common.MessageFlag    `json:"flags,omitempty"`
	Nonce        interface{}           `json:"nonce,omitempty"`
	EnforceNonce bool                  `json:"enforce_nonce,omitempty"`
}

func (p CreateMessageParams) MarshalJSON() ([]byte, error) {
	type Alias CreateMessageParams
	return json.Marshal(&struct {
		Components []common.AnyComponent `json:"components,omitempty"`
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
	Content         *string                 `json:"content,omitempty"`
	Embeds          *[]common.Embed         `json:"embeds,omitempty"`
	Flags           *common.MessageFlag     `json:"flags,omitempty"`
	AllowedMentions *common.AllowedMentions `json:"allowed_mentions,omitempty"`
	// Components replaces the message's component list.
	// Set to an empty (non-nil) slice to remove all components.
	Components  []common.AnyComponent `json:"-"`
	Attachments *[]common.Attachment  `json:"attachments,omitempty"`
}

func (p EditMessageParams) MarshalJSON() ([]byte, error) {
	type Alias EditMessageParams
	return json.Marshal(&struct {
		Components []common.AnyComponent `json:"components,omitempty"`
		Alias
	}{
		Components: p.Components,
		Alias:      (Alias)(p),
	})
}

// GetReactionsParams are optional query parameters for listing reaction users.
type GetReactionsParams struct {
	After *common.Snowflake
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
func (c *RestClient) GetMessages(channelID common.Snowflake, params GetMessagesParams) ([]*common.Message, error) {
	path := "/channels/" + channelID.String() + "/messages" + params.toQuery()
	req, err := c.generateRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var msgs []*common.Message
	if _, err := c.do(req, http.StatusOK, &msgs); err != nil {
		return nil, err
	}

	return msgs, nil
}

// GetMessage returns a single message by ID.
func (c *RestClient) GetMessage(channelID, messageID common.Snowflake) (*common.Message, error) {
	path := "/channels/" + channelID.String() + "/messages/" + messageID.String()
	req, err := c.generateRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var msg common.Message
	if _, err := c.do(req, http.StatusOK, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// CreateMessage sends a new message to a channel.
func (c *RestClient) CreateMessage(channelID common.Snowflake, params CreateMessageParams) (*common.Message, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(http.MethodPost, "/channels/"+channelID.String()+"/messages", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var msg common.Message
	if _, err := c.do(req, http.StatusOK, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// EditMessage edits a previously sent message. Only fields set in params are changed.
func (c *RestClient) EditMessage(channelID, messageID common.Snowflake, params EditMessageParams) (*common.Message, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/messages/" + messageID.String()
	req, err := c.generateRequest(http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var msg common.Message
	if _, err := c.do(req, http.StatusOK, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// DeleteMessage deletes a message.
func (c *RestClient) DeleteMessage(channelID, messageID common.Snowflake) error {
	path := "/channels/" + channelID.String() + "/messages/" + messageID.String()
	req, err := c.generateRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// BulkDeleteMessages deletes 2–100 messages at once.
// Messages older than 14 days cannot be bulk-deleted and will cause a Discord API error.
func (c *RestClient) BulkDeleteMessages(channelID common.Snowflake, messageIDs []common.Snowflake) error {
	if len(messageIDs) < 2 || len(messageIDs) > 100 {
		return fmt.Errorf("bulk delete requires 2–100 message IDs, got %d", len(messageIDs))
	}

	body, err := json.Marshal(map[string]any{"messages": messageIDs})
	if err != nil {
		return err
	}

	path := "/channels/" + channelID.String() + "/messages/bulk-delete"
	req, err := c.generateRequest(http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// CrosspostMessage publishes a message in an announcement channel to all following channels.
func (c *RestClient) CrosspostMessage(channelID, messageID common.Snowflake) (*common.Message, error) {
	path := "/channels/" + channelID.String() + "/messages/" + messageID.String() + "/crosspost"
	req, err := c.generateRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	var msg common.Message
	if _, err := c.do(req, http.StatusOK, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// ── Pin endpoints ─────────────────────────────────────────────────────────────

// GetPinnedMessages returns all pinned messages in a channel (max 50).
func (c *RestClient) GetPinnedMessages(channelID common.Snowflake) ([]*common.Message, error) {
	req, err := c.generateRequest(http.MethodGet, "/channels/"+channelID.String()+"/pins", nil)
	if err != nil {
		return nil, err
	}

	var msgs []*common.Message
	if _, err := c.do(req, http.StatusOK, &msgs); err != nil {
		return nil, err
	}

	return msgs, nil
}

// PinMessage pins a message in a channel. Requires MANAGE_MESSAGES.
func (c *RestClient) PinMessage(channelID, messageID common.Snowflake) error {
	path := "/channels/" + channelID.String() + "/pins/" + messageID.String()
	req, err := c.generateRequest(http.MethodPut, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// UnpinMessage unpins a message from a channel. Requires MANAGE_MESSAGES.
func (c *RestClient) UnpinMessage(channelID, messageID common.Snowflake) error {
	path := "/channels/" + channelID.String() + "/pins/" + messageID.String()
	req, err := c.generateRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// ── Reaction endpoints ────────────────────────────────────────────────────────

// AddReaction adds a reaction to a message.
// emoji is a raw Unicode character (e.g. "👍") or a custom emoji in "name:id" form.
func (c *RestClient) AddReaction(channelID, messageID common.Snowflake, emoji string) error {
	path := fmt.Sprintf("/channels/%s/messages/%s/reactions/%s/@me", channelID, messageID, encodeEmoji(emoji))
	req, err := c.generateRequest(http.MethodPut, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// DeleteOwnReaction removes the bot's own reaction from a message.
func (c *RestClient) DeleteOwnReaction(channelID, messageID common.Snowflake, emoji string) error {
	path := fmt.Sprintf("/channels/%s/messages/%s/reactions/%s/@me", channelID, messageID, encodeEmoji(emoji))
	req, err := c.generateRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// DeleteUserReaction removes another user's reaction from a message. Requires MANAGE_MESSAGES.
func (c *RestClient) DeleteUserReaction(channelID, messageID common.Snowflake, emoji string, userID common.Snowflake) error {
	path := fmt.Sprintf("/channels/%s/messages/%s/reactions/%s/%s", channelID, messageID, encodeEmoji(emoji), userID)
	req, err := c.generateRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// GetReactions returns the users who reacted to a message with the given emoji.
func (c *RestClient) GetReactions(channelID, messageID common.Snowflake, emoji string, params GetReactionsParams) ([]*common.User, error) {
	path := fmt.Sprintf("/channels/%s/messages/%s/reactions/%s%s", channelID, messageID, encodeEmoji(emoji), params.toQuery())
	req, err := c.generateRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var users []*common.User
	if _, err := c.do(req, http.StatusOK, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// DeleteAllReactions removes every reaction from a message. Requires MANAGE_MESSAGES.
func (c *RestClient) DeleteAllReactions(channelID, messageID common.Snowflake) error {
	path := fmt.Sprintf("/channels/%s/messages/%s/reactions", channelID, messageID)
	req, err := c.generateRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

// DeleteAllReactionsForEmoji removes all reactions for a specific emoji. Requires MANAGE_MESSAGES.
func (c *RestClient) DeleteAllReactionsForEmoji(channelID, messageID common.Snowflake, emoji string) error {
	path := fmt.Sprintf("/channels/%s/messages/%s/reactions/%s", channelID, messageID, encodeEmoji(emoji))
	req, err := c.generateRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}
