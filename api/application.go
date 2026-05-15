package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ModifyCurrentApplicationParams holds the fields that can be updated on the current application.
type ModifyCurrentApplicationParams struct {
	Description                    *string                                `json:"description,omitempty"`
	Icon                           *string                                `json:"icon,omitempty"`
	CoverImage                     *string                                `json:"cover_image,omitempty"`
	InteractionEndpointURL         *string                                `json:"interactions_endpoint_url,omitempty"`
	Tags                           []string                               `json:"tags,omitempty"`
	InstallParams                  *discord.ApplicationInstallParams      `json:"install_params,omitempty"`
	CustomInstallURL               *string                                `json:"custom_install_url,omitempty"`
	RoleConnectionsVerificationURL *string                                `json:"role_connections_verification_url,omitempty"`
	EventWebhooksURL               *string                                `json:"event_webhooks_url,omitempty"`
	EventWebhooksStatus            *discord.ApplicationEventWebhookStatus `json:"event_webhooks_status,omitempty"`
	EventWebhooksTypes             []string                               `json:"event_webhooks_types,omitempty"`
}

// GetCurrentApplication returns the application object for the bot's application.
func (c *RestClient) GetCurrentApplication(ctx context.Context) (*discord.Application, error) {
	req, err := c.generateRequest(ctx, http.MethodGet, "/applications/@me", nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Application](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// ModifyCurrentApplication updates the current application. Returns the updated application.
func (c *RestClient) ModifyCurrentApplication(ctx context.Context, params ModifyCurrentApplicationParams) (*discord.Application, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.generateRequest(ctx, http.MethodPatch, "/applications/@me", bytes.NewReader(body), c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Application](c, req, map[int]bool{
		http.StatusOK: true,
	})
}
