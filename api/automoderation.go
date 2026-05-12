package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// ── Param types ───────────────────────────────────────────────────────────────

type CreateAutoModerationRuleParams struct {
	Name            string                                `json:"name"`
	EventType       common.AutoModerationEventType        `json:"event_type"`
	TriggerType     common.AutoModerationTriggerType      `json:"trigger_type"`
	TriggerMetadata *common.AutoModerationTriggerMetadata `json:"trigger_metadata,omitempty"`
	Actions         []common.AutoModerationAction         `json:"actions"`
	Enabled         *bool                                 `json:"enabled,omitempty"`
	ExemptRoles     []common.Snowflake                    `json:"exempt_roles,omitempty"`
	ExemptChannels  []common.Snowflake                    `json:"exempt_channels,omitempty"`
}

type ModifyAutoModerationRuleParams struct {
	Name            *string                               `json:"name,omitempty"`
	EventType       *common.AutoModerationEventType       `json:"event_type,omitempty"`
	TriggerMetadata *common.AutoModerationTriggerMetadata `json:"trigger_metadata,omitempty"`
	Actions         []common.AutoModerationAction         `json:"actions,omitempty"`
	Enabled         *bool                                 `json:"enabled,omitempty"`
	ExemptRoles     []common.Snowflake                    `json:"exempt_roles,omitempty"`
	ExemptChannels  []common.Snowflake                    `json:"exempt_channels,omitempty"`
}

// ── Auto moderation endpoints ─────────────────────────────────────────────────

// ListAutoModerationRules returns all auto moderation rules for a guild.
func (c *RestClient) ListAutoModerationRules(ctx context.Context, guildID common.Snowflake) ([]*common.AutoModerationRule, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/auto-moderation/rules", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var rules []*common.AutoModerationRule
	if _, err := c.do(req, []int{http.StatusOK}, &rules); err != nil {
		return nil, err
	}

	return rules, nil
}

// GetAutoModerationRule returns a single auto moderation rule.
func (c *RestClient) GetAutoModerationRule(ctx context.Context, guildID, ruleID common.Snowflake) (*common.AutoModerationRule, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := ruleID.Validate(); err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/auto-moderation/rules/" + ruleID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var rule common.AutoModerationRule
	if _, err := c.do(req, []int{http.StatusOK}, &rule); err != nil {
		return nil, err
	}

	return &rule, nil
}

// CreateAutoModerationRule creates a new auto moderation rule in a guild.
func (c *RestClient) CreateAutoModerationRule(ctx context.Context, guildID common.Snowflake, params CreateAutoModerationRuleParams) (*common.AutoModerationRule, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/auto-moderation/rules", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var rule common.AutoModerationRule
	if _, err := c.do(req, []int{http.StatusOK}, &rule); err != nil {
		return nil, err
	}

	return &rule, nil
}

// ModifyAutoModerationRule updates an existing auto moderation rule.
func (c *RestClient) ModifyAutoModerationRule(ctx context.Context, guildID, ruleID common.Snowflake, params ModifyAutoModerationRuleParams) (*common.AutoModerationRule, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	if err := ruleID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/guilds/" + guildID.String() + "/auto-moderation/rules/" + ruleID.String()
	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	var rule common.AutoModerationRule
	if _, err := c.do(req, []int{http.StatusOK}, &rule); err != nil {
		return nil, err
	}

	return &rule, nil
}

// DeleteAutoModerationRule deletes an auto moderation rule.
func (c *RestClient) DeleteAutoModerationRule(ctx context.Context, guildID, ruleID common.Snowflake) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	if err := ruleID.Validate(); err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/auto-moderation/rules/" + ruleID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	_, err = c.do(req, []int{http.StatusNoContent}, nil)
	return err
}
