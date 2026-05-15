package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Param types ───────────────────────────────────────────────────────────────

type CreateAutoModerationRuleParams struct {
	Name            string                                 `json:"name"`
	EventType       discord.AutoModerationEventType        `json:"event_type"`
	TriggerType     discord.AutoModerationTriggerType      `json:"trigger_type"`
	TriggerMetadata *discord.AutoModerationTriggerMetadata `json:"trigger_metadata,omitempty"`
	Actions         []discord.AutoModerationAction         `json:"actions"`
	Enabled         *bool                                  `json:"enabled,omitempty"`
	ExemptRoles     []discord.Snowflake                    `json:"exempt_roles,omitempty"`
	ExemptChannels  []discord.Snowflake                    `json:"exempt_channels,omitempty"`
}

type CreateAutoModerationRuleOptions struct {
	Reason string
}

type ModifyAutoModerationRuleParams struct {
	Name            *string                                `json:"name,omitempty"`
	EventType       *discord.AutoModerationEventType       `json:"event_type,omitempty"`
	TriggerMetadata *discord.AutoModerationTriggerMetadata `json:"trigger_metadata,omitempty"`
	Actions         []discord.AutoModerationAction         `json:"actions,omitempty"`
	Enabled         *bool                                  `json:"enabled,omitempty"`
	ExemptRoles     []discord.Snowflake                    `json:"exempt_roles,omitempty"`
	ExemptChannels  []discord.Snowflake                    `json:"exempt_channels,omitempty"`
}

type ModifyAutoModerationRuleOptions struct {
	Reason string
}

type DeleteAutoModerationRuleOptions struct {
	Reason string
}

// ── Auto moderation endpoints ─────────────────────────────────────────────────

// ListAutoModerationRules returns all auto moderation rules for a guild.
func (c *RestClient) ListAutoModerationRules(ctx context.Context, guildID discord.Snowflake) (*[]*discord.AutoModerationRule, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodGet, "/guilds/"+guildID.String()+"/auto-moderation/rules", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[[]*discord.AutoModerationRule](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetAutoModerationRule returns a single auto moderation rule.
func (c *RestClient) GetAutoModerationRule(ctx context.Context, guildID, ruleID discord.Snowflake) (*discord.AutoModerationRule, error) {
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

	return doRequest[discord.AutoModerationRule](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// CreateAutoModerationRule creates a new auto moderation rule in a guild.
func (c *RestClient) CreateAutoModerationRule(ctx context.Context, guildID discord.Snowflake, params CreateAutoModerationRuleParams, opts *CreateAutoModerationRuleOptions) (*discord.AutoModerationRule, error) {
	if err := guildID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/auto-moderation/rules", bytes.NewReader(body), c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		return doRequest[discord.AutoModerationRule](c, req, map[int]bool{
			http.StatusOK: true,
		})
	}

	req, err := c.generateRequest(ctx, http.MethodPost, "/guilds/"+guildID.String()+"/auto-moderation/rules", bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.AutoModerationRule](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyAutoModerationRule updates an existing auto moderation rule.
func (c *RestClient) ModifyAutoModerationRule(ctx context.Context, guildID, ruleID discord.Snowflake, params ModifyAutoModerationRuleParams, opts *ModifyAutoModerationRuleOptions) (*discord.AutoModerationRule, error) {
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

	if opts == nil {
		req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		return doRequest[discord.AutoModerationRule](c, req, map[int]bool{
			http.StatusOK: true,
		})
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, path, bytes.NewReader(body), c.WithBotAuthorization(), WithAuditLogReason(opts.Reason))
	if err != nil {
		return nil, err
	}

	return doRequest[discord.AutoModerationRule](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// DeleteAutoModerationRule deletes an auto moderation rule.
func (c *RestClient) DeleteAutoModerationRule(ctx context.Context, guildID, ruleID discord.Snowflake, opts *DeleteAutoModerationRuleOptions) error {
	if err := guildID.Validate(); err != nil {
		return err
	}

	if err := ruleID.Validate(); err != nil {
		return err
	}

	path := "/guilds/" + guildID.String() + "/auto-moderation/rules/" + ruleID.String()

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
