package interactions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DatGamet/go-discord-wrapper/types/common"
	"github.com/DatGamet/go-discord-wrapper/types/components"
	"github.com/DatGamet/go-discord-wrapper/types/interactions/responses"
	"io"
	"net/http"
	"net/url"
)

func (i *Interaction) GetSubCommand() string {
	if i.Data == nil {
		return ""
	}

	cmdData, ok := i.Data.(*responses.InteractionDataApplicationCommand)
	if !ok {
		return ""
	}

	if cmdData.Options == nil || len(*cmdData.Options) == 0 {
		return ""
	}

	for _, option := range *cmdData.Options {
		if option.Type == common.ApplicationCommandOptionTypeSubCommand {
			return option.Name
		}

		if option.Type == common.ApplicationCommandOptionTypeSubCommandGroup {
			if option.Options != nil {
				for _, subOption := range option.Options {
					if subOption.Type == common.ApplicationCommandOptionTypeSubCommandGroup {
						return subOption.Name
					}
				}
			}
		}
	}

	return ""
}

func (i *Interaction) GetSubCommandGroup() string {
	if i.Data == nil {
		return ""
	}

	cmdData, ok := i.Data.(*responses.InteractionDataApplicationCommand)
	if !ok {
		return ""
	}

	if cmdData.Options == nil || len(*cmdData.Options) == 0 {
		return ""
	}

	for _, option := range *cmdData.Options {
		if option.Type == common.ApplicationCommandOptionTypeSubCommandGroup {
			return option.Name
		}
	}

	return ""
}

func (i *Interaction) GetFullCommand() (fullCommand string) {
	if i.Data == nil {
		return ""
	}

	cmdData, ok := i.Data.(*responses.InteractionDataApplicationCommand)
	if !ok {
		return ""
	}

	fullCommand += cmdData.CommandName

	subCommandGroup := i.GetSubCommandGroup()
	if subCommandGroup != "" {
		fullCommand += " " + subCommandGroup
	}

	subCommand := i.GetSubCommand()
	if subCommand != "" {
		fullCommand += " " + subCommand
	}

	return fullCommand
}

func (i *Interaction) GetCustomID() string {
	if i.Data == nil {
		return ""
	}

	componentData, ok := i.Data.(*responses.InteractionDataMessageComponent)
	componentData2, ok2 := i.Data.(*responses.InteractionDataModalSubmit)

	if !ok && !ok2 {
		return ""
	}

	if ok2 {
		return componentData2.CustomID
	}

	return componentData.CustomID
}

func (i *Interaction) UnmarshalJSON(data []byte) error {
	type Alias Interaction
	aux := &struct {
		Data json.RawMessage `json:"data"`
		*Alias
	}{
		Alias: (*Alias)(i),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if aux.Data == nil {
		return nil
	}

	var typeProbe struct {
		Type          common.ApplicationCommandType `json:"type"`
		ComponentType common.ComponentType          `json:"component_type"`
	}

	if err := json.Unmarshal(aux.Data, &typeProbe); err != nil {
		return err
	}

	switch typeProbe.Type {
	case common.ApplicationCommandTypeChatInput, common.ApplicationCommandTypeUser, common.ApplicationCommandTypeMessage:
		var cmd responses.InteractionDataApplicationCommand
		if err := json.Unmarshal(aux.Data, &cmd); err != nil {
			return err
		}
		i.Data = &cmd
		return nil
	}

	switch typeProbe.ComponentType {
	case common.ComponentTypeButton, common.ComponentTypeStringSelectMenu, common.ComponentTypeUserSelectMenu, common.ComponentTypeRoleSelectMenu,
		common.ComponentTypeMentionableMenu, common.ComponentTypeChannelSelect:
		var comp responses.InteractionDataMessageComponent
		if err := json.Unmarshal(aux.Data, &comp); err != nil {
			return err
		}
		i.Data = &comp
		return nil
	}

	switch aux.Type {
	case common.InteractionTypeModalSubmit:
		var modal responses.InteractionDataModalSubmit
		if err := json.Unmarshal(aux.Data, &modal); err != nil {
			return err
		}
		i.Data = &modal
		return nil
	case common.InteractionTypeApplicationCommandAutocomplete:
		var auto responses.InteractionDataAutocomplete
		if err := json.Unmarshal(aux.Data, &auto); err != nil {
			return err
		}
		i.Data = &auto
		return nil
	}

	return fmt.Errorf("unknown interaction data type %d", typeProbe.Type)
}

func (i *Interaction) DeferReply() error {
	//TODO
	return nil
}

func (i *Interaction) EditReply(responseData *responses.AnyInteractionResponseData, clientID string) error {
	bodyBytes, err := json.Marshal(*responseData)
	if err != nil {
		return err
	}

	req, err := http.DefaultClient.Do(&http.Request{
		Method: "PATCH",
		URL: &url.URL{
			Scheme: "https",
			Host:   "discord.com",
			Path:   "/api/v10/webhooks/" + clientID + "/" + i.Token + "/messages/@original",
		},
		Header: http.Header{
			"Authorization": []string{"Bot " + ""},
			"Content-Type":  []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader(bodyBytes)),
	})

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(req.Body)

	if req.StatusCode != http.StatusOK {
		var respErr map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&respErr); err != nil {
			return err
		}

		return fmt.Errorf("expected 204 No Content, got %d: %v", req.StatusCode, respErr)
	}

	return nil
}

func (i *Interaction) DeleteReply(clientID string) error {
	req, err := http.DefaultClient.Do(&http.Request{
		Method: "DELETE",
		URL: &url.URL{
			Scheme: "https",
			Host:   "discord.com",
			Path:   "/api/v10/webhooks/" + clientID + "/" + i.Token + "/messages/@original",
		},
		Header: http.Header{
			"Authorization": []string{"Bot " + ""},
			"Content-Type":  []string{"application/json"},
		},
	})

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(req.Body)

	if req.StatusCode != http.StatusNoContent {
		var respErr map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&respErr); err != nil {
			return err
		}

		return fmt.Errorf("expected 204 No Content, got %d: %v", req.StatusCode, respErr)
	}

	return nil
}

func (i *Interaction) ReplyWithModal(modal *components.Modal) error {
	bodyBytes, err := json.Marshal(responses.InteractionResponse{
		Type: common.InteractionCallbackTypeModal,
		Data: modal,
	})
	if err != nil {
		return err
	}

	req, err := http.DefaultClient.Do(&http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "https",
			Host:   "discord.com",
			Path:   "/api/v10/interactions/" + string(i.ID) + "/" + i.Token + "/callback",
		},
		Header: http.Header{
			"Authorization": []string{"Bot " + ""},
			"Content-Type":  []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader(bodyBytes)),
	})

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(req.Body)

	if req.StatusCode != http.StatusNoContent {
		var respErr map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&respErr); err != nil {
			return err
		}

		return fmt.Errorf("expected 204 No Content, got %d: %v", req.StatusCode, respErr)
	}

	return nil
}

func (i *Interaction) Reply(data *responses.InteractionResponseDataDefault) (*responses.InteractionCallbackResponse, error) {
	bodyBytes, err := json.Marshal(responses.InteractionResponse{
		Type: common.InteractionCallbackTypeChannelMessageWithSource,
		Data: data,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.DefaultClient.Do(&http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme:   "https",
			Host:     "discord.com",
			Path:     "/api/v10/interactions/" + string(i.ID) + "/" + i.Token + "/callback",
			RawQuery: "with_response=" + fmt.Sprintf("%t", data.WithResponse),
		},
		Header: http.Header{
			"Authorization": []string{"Bot " + ""},
			"Content-Type":  []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader(bodyBytes)),
	})

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(req.Body)

	if !data.WithResponse {
		if req.StatusCode != 204 {
			var respErr map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&respErr); err != nil {
				return nil, err
			}

			return nil, fmt.Errorf("expected 204 No Content, got %d: %v", req.StatusCode, respErr)
		}

		return nil, nil
	}

	if data.WithResponse && req.StatusCode != 200 {
		var respErr map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&respErr); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("expected 204 No Content, got %d: %v", req.StatusCode, respErr)
	}

	var resp responses.InteractionCallbackResponse
	if err := json.NewDecoder(req.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
