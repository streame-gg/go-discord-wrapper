package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/common"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

func (c *RestClient) CreateInteractionResponse(interactionID common.Snowflake, token string, response responses.InteractionResponse, withResponse bool) (*responses.InteractionCallbackResponse, error) {
	body, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	path := "/interactions/" + interactionID.String() + "/" + token + "/callback"
	if withResponse {
		path += "?with_response=true"
	}

	req, err := c.generateRequest(http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	if !withResponse {
		_, err = c.do(req, http.StatusNoContent, nil)
		return nil, err
	}

	var result responses.InteractionCallbackResponse
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *RestClient) GetOriginalInteractionResponse(appID common.Snowflake, token string) (*common.Message, error) {
	path := "/webhooks/" + appID.String() + "/" + token + "/messages/@original"
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

func (c *RestClient) EditOriginalInteractionResponse(appID common.Snowflake, token string, params EditMessageParams) (*common.Message, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/webhooks/" + appID.String() + "/" + token + "/messages/@original"
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

func (c *RestClient) DeleteOriginalInteractionResponse(appID common.Snowflake, token string) error {
	path := "/webhooks/" + appID.String() + "/" + token + "/messages/@original"
	req, err := c.generateRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}

func (c *RestClient) CreateFollowupMessage(appID common.Snowflake, token string, params CreateMessageParams) (*common.Message, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/webhooks/" + appID.String() + "/" + token
	req, err := c.generateRequest(http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var msg common.Message
	if _, err := c.do(req, http.StatusOK, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

func (c *RestClient) GetFollowupMessage(appID common.Snowflake, token string, messageID common.Snowflake) (*common.Message, error) {
	path := "/webhooks/" + appID.String() + "/" + token + "/messages/" + messageID.String()
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

func (c *RestClient) EditFollowupMessage(appID common.Snowflake, token string, messageID common.Snowflake, params EditMessageParams) (*common.Message, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	path := "/webhooks/" + appID.String() + "/" + token + "/messages/" + messageID.String()
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

func (c *RestClient) DeleteFollowupMessage(appID common.Snowflake, token string, messageID common.Snowflake) error {
	path := "/webhooks/" + appID.String() + "/" + token + "/messages/" + messageID.String()
	req, err := c.generateRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, http.StatusNoContent, nil)
	return err
}
