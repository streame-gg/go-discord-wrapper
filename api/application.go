package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// ModifyCurrentApplicationParams holds the fields that can be updated on the current application.
type ModifyCurrentApplicationParams struct {
	Description                    *string                               `json:"description,omitempty"`
	Icon                           *string                               `json:"icon,omitempty"`
	CoverImage                     *string                               `json:"cover_image,omitempty"`
	InteractionEndpointURL         *string                               `json:"interactions_endpoint_url,omitempty"`
	Tags                           []string                              `json:"tags,omitempty"`
	InstallParams                  *common.ApplicationInstallParams      `json:"install_params,omitempty"`
	CustomInstallURL               *string                               `json:"custom_install_url,omitempty"`
	RoleConnectionsVerificationURL *string                               `json:"role_connections_verification_url,omitempty"`
	EventWebhooksURL               *string                               `json:"event_webhooks_url,omitempty"`
	EventWebhooksStatus            *common.ApplicationEventWebhookStatus `json:"event_webhooks_status,omitempty"`
	EventWebhooksTypes             []string                              `json:"event_webhooks_types,omitempty"`
}

// GetCurrentApplication returns the application object for the bot's application.
func (c *RestClient) GetCurrentApplication(ctx context.Context) (*common.Application, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/applications/@me", nil)
	if err != nil {
		return nil, err
	}

	var app common.Application
	if _, err := c.do(req, http.StatusOK, &app); err != nil {
		return nil, err
	}

	return &app, nil
}

// ModifyCurrentApplication updates the current application. Returns the updated application.
func (c *RestClient) ModifyCurrentApplication(ctx context.Context, params ModifyCurrentApplicationParams) (*common.Application, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/applications/@me", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var app common.Application
	if _, err := c.do(req, http.StatusOK, &app); err != nil {
		return nil, err
	}

	return &app, nil
}
