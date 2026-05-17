package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

// CreateInteractionResponse sends a response to an interaction. Must be called within 3 seconds.
// Set withResponse=true to receive the created message back; returns nil otherwise.
// Pass optional files to send attachments; when present the request is encoded as multipart/form-data.
func (c *RestClient) CreateInteractionResponse(ctx context.Context, interactionID discord.Snowflake, token string, response responses.InteractionResponse, withResponse bool, files ...discord.MessageFile) (*responses.InteractionCallbackResponse, error) {
	if err := interactionID.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	path := "/interactions/" + interactionID.String() + "/" + url.PathEscape(token) + "/callback"
	if withResponse {
		path += "?with_response=true"
	}

	var req *http.Request
	if len(files) > 0 {
		buf, ct, err := buildMultipartMessage(body, files)
		if err != nil {
			return nil, err
		}
		req, err = c.generateRequest(ctx, http.MethodPost, path, buf, c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", ct)
	} else {
		req, err = c.generateRequest(ctx, http.MethodPost, path, bytes.NewReader(body), c.WithBotAuthorization())
		if err != nil {
			return nil, err
		}
	}

	if !withResponse {
		return nil, doRequestWithoutResponse(c, req)
	}

	return doRequest[responses.InteractionCallbackResponse](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// GetOriginalInteractionResponse fetches the initial response message for an interaction.
func (c *RestClient) GetOriginalInteractionResponse(ctx context.Context, webhookID discord.Snowflake, token string) (*discord.Message, error) {
	if err := webhookID.Validate(); err != nil {
		return nil, err
	}

	path := "/webhooks/" + webhookID.String() + "/" + url.PathEscape(token) + "/messages/@original"
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Message](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// EditOriginalInteractionResponse edits the initial response message for an interaction.
// When params.Files is non-empty the request is sent as multipart/form-data.
func (c *RestClient) EditOriginalInteractionResponse(ctx context.Context, webhookID discord.Snowflake, token string, params EditMessageParams) (*discord.Message, error) {
	if err := webhookID.Validate(); err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/webhooks/" + webhookID.String() + "/" + url.PathEscape(token) + "/messages/@original"

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

// DeleteOriginalInteractionResponse deletes the initial response message for an interaction.
func (c *RestClient) DeleteOriginalInteractionResponse(ctx context.Context, webhookID discord.Snowflake, token string) error {
	if err := webhookID.Validate(); err != nil {
		return err
	}

	path := "/webhooks/" + webhookID.String() + "/" + url.PathEscape(token) + "/messages/@original"
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}

// CreateFollowupMessage sends a follow-up message to an interaction (usable up to 15 minutes after the initial response).
// When params.Files is non-empty the request is sent as multipart/form-data.
func (c *RestClient) CreateFollowupMessage(ctx context.Context, appID discord.Snowflake, token string, params CreateMessageParams) (*discord.Message, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/webhooks/" + appID.String() + "/" + url.PathEscape(token)

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

// GetFollowupMessage fetches a follow-up message sent for an interaction.
func (c *RestClient) GetFollowupMessage(ctx context.Context, appID discord.Snowflake, token string, messageID discord.Snowflake) (*discord.Message, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	if err := messageID.Validate(); err != nil {
		return nil, err
	}

	path := "/webhooks/" + appID.String() + "/" + url.PathEscape(token) + "/messages/" + messageID.String()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Message](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

// EditFollowupMessage edits a follow-up message sent for an interaction.
// When params.Files is non-empty the request is sent as multipart/form-data.
func (c *RestClient) EditFollowupMessage(ctx context.Context, appID discord.Snowflake, token string, messageID discord.Snowflake, params EditMessageParams) (*discord.Message, error) {
	if err := appID.Validate(); err != nil {
		return nil, err
	}

	if err := messageID.Validate(); err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/webhooks/" + appID.String() + "/" + url.PathEscape(token) + "/messages/" + messageID.String()

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

// DeleteFollowupMessage deletes a follow-up message sent for an interaction.
func (c *RestClient) DeleteFollowupMessage(ctx context.Context, appID discord.Snowflake, token string, messageID discord.Snowflake) error {
	if err := appID.Validate(); err != nil {
		return err
	}

	if err := messageID.Validate(); err != nil {
		return err
	}

	path := "/webhooks/" + appID.String() + "/" + url.PathEscape(token) + "/messages/" + messageID.String()
	req, err := c.generateRequest(ctx, http.MethodDelete, path, nil, c.WithBotAuthorization())
	if err != nil {
		return err
	}

	return doRequestWithoutResponse(c, req)
}
