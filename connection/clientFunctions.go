package connection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DatGamet/go-discord-wrapper/types/commands"
	"github.com/DatGamet/go-discord-wrapper/types/common"
	"net/http"
)

func (d *Client) RegisterSingleCommand(command *commands.ApplicationCommand) (*commands.ApplicationCommand, error) {
	body, marshalErr := json.Marshal(command)
	if marshalErr != nil {
		return nil, marshalErr
	}

	req, errCreatingRequest := http.NewRequest(
		"POST",
		"https://discord.com"+common.APIBaseString(*d.APIVersion)+"applications/"+d.User.ID.ToString()+"/commands",
		bytes.NewBuffer(body),
	)
	if errCreatingRequest != nil {
		return nil, errCreatingRequest
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bot "+*d.token)

	do, errDoingRequest := http.DefaultClient.Do(req)
	if errDoingRequest != nil {
		return nil, errDoingRequest
	}

	defer func() {
		_ = do.Body.Close()
	}()

	if do.StatusCode != http.StatusOK && do.StatusCode != http.StatusCreated {
		var body map[string]interface{}
		if err := json.NewDecoder(do.Body).Decode(&body); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("failed to register command, status code: %s, body: %v", do.Status, body)
	}

	var registeredCommand commands.ApplicationCommand
	if err := json.NewDecoder(do.Body).Decode(&registeredCommand); err != nil {
		return nil, err
	}

	return &registeredCommand, nil
}

func (d *Client) BulkRegisterCommands(providedCommands []*commands.ApplicationCommand) (*[]*commands.ApplicationCommand, error) {
	body, marshalErr := json.Marshal(providedCommands)
	if marshalErr != nil {
		return nil, marshalErr
	}

	req, errCreatingRequest := http.NewRequest(
		"PUT",
		"https://discord.com"+common.APIBaseString(*d.APIVersion)+"applications/"+d.User.ID.ToString()+"/commands",
		bytes.NewBuffer(body),
	)
	if errCreatingRequest != nil {
		return nil, errCreatingRequest
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bot "+*d.token)

	do, errDoingRequest := http.DefaultClient.Do(req)
	if errDoingRequest != nil {
		return nil, errDoingRequest
	}

	defer func() {
		_ = do.Body.Close()
	}()

	if do.StatusCode != http.StatusOK && do.StatusCode != http.StatusCreated {
		var body map[string]interface{}
		if err := json.NewDecoder(do.Body).Decode(&body); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("failed to register command, status code: %s, body: %v", do.Status, body)
	}

	if do.StatusCode != http.StatusOK && do.StatusCode != http.StatusCreated {
		var body map[string]interface{}
		if err := json.NewDecoder(do.Body).Decode(&body); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("failed to register command, status code: %s, body: %v", do.Status, body)
	}

	var registeredCommand []*commands.ApplicationCommand
	if err := json.NewDecoder(do.Body).Decode(&registeredCommand); err != nil {
		return nil, err
	}

	return &registeredCommand, nil
}
