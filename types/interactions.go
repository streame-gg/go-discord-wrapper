package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Entitlement struct {
	ID            Snowflake       `json:"id"`
	SkuID         Snowflake       `json:"sku_id"`
	ApplicationID Snowflake       `json:"application_id"`
	UserID        *Snowflake      `json:"user_id"`
	Type          EntitlementType `json:"type"`
	Deleted       bool            `json:"deleted"`
	StartsAt      *time.Time      `json:"starts_at,omitempty"`
	EndsAt        *time.Time      `json:"ends_at,omitempty"`
	GuildID       *Snowflake      `json:"guild_id,omitempty"`
	Consumed      bool            `json:"consumed,omitempty"`
}

type EntitlementType int

const (
	EntitlementTypePurchase                EntitlementType = 1
	EntitlementTypePremiumSubscription     EntitlementType = 2
	EntitlementTypeDeveloperGift           EntitlementType = 3
	EntitlementTypeTestModePurchase        EntitlementType = 4
	EntitlementTypeFreePurchase            EntitlementType = 5
	EntitlementTypeUserGift                EntitlementType = 6
	EntitlementTypePremiumPurchase         EntitlementType = 7
	EntitlementTypeApplicationSubscription EntitlementType = 8
)

type InteractionType int

const (
	InteractionTypePing                           InteractionType = 1
	InteractionTypeApplicationCommand             InteractionType = 2
	InteractionTypeMessageComponent               InteractionType = 3
	InteractionTypeApplicationCommandAutocomplete InteractionType = 4
	InteractionTypeModalSubmit                    InteractionType = 5
)

type InteractionApplicationIntegrationType int

const (
	InteractionApplicationIntegrationTypeGuildInstall InteractionApplicationIntegrationType = 0
	InteractionApplicationIntegrationTypeUserInstall  InteractionApplicationIntegrationType = 1
)

type MessageMessageReference struct {
	Type            *MessageMessageReferenceType `json:"type,omitempty"`
	MessageID       *Snowflake                   `json:"message_id,omitempty"`
	ChannelID       *Snowflake                   `json:"channel_id,omitempty"`
	GuildID         *Snowflake                   `json:"guild_id,omitempty"`
	FailIfNotExists *bool                        `json:"fail_if_not_exists,omitempty"`
}

type MessageInteractionMetadataApplicationCommand struct {
	ID                           Snowflake                                             `json:"id"`
	Type                         InteractionType                                       `json:"type"`
	User                         User                                                  `json:"user,omitempty"`
	AuthorizingIntegrationOwners map[InteractionApplicationIntegrationType]interface{} `json:"authorizing_integration_owners,omitempty"`
	OriginalResponseMessageID    *Snowflake                                            `json:"original_response_message_id,omitempty"`
	TargetUser                   *User                                                 `json:"target_user,omitempty"`
	TargetMessageID              *Snowflake                                            `json:"target_message_id,omitempty"`
}

type MessageInteractionMetadataMessageComponent struct {
	ID                           Snowflake                                             `json:"id"`
	Type                         InteractionType                                       `json:"type"`
	User                         User                                                  `json:"user,omitempty"`
	AuthorizingIntegrationOwners map[InteractionApplicationIntegrationType]interface{} `json:"authorizing_integration_owners,omitempty"`
	OriginalResponseMessageID    *Snowflake                                            `json:"original_response_message_id,omitempty"`
	InteractedMessageID          *Snowflake                                            `json:"interacted_message_id,omitempty"`
}

type AnyMessageInteractionMetadata interface{}

type MessageInteractionMetadata struct {
	Value AnyMessageInteractionMetadata
}

type AnyMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata interface{}

type MessageInteractionMetadataModalSubmitTriggering struct {
	AnyMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata
}

func (d *MessageInteractionMetadataModalSubmitTriggering) UnmarshalJSON(data []byte) error {
	var a MessageInteractionMetadataApplicationCommand
	if err := json.Unmarshal(data, &a); err == nil && a.ID != "" {
		d.AnyMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata = &a
		return nil
	}

	var b MessageInteractionMetadataMessageComponent
	if err := json.Unmarshal(data, &b); err == nil && b.ID != "" {
		d.AnyMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata = &b
		return nil
	}

	return nil
}

func (d *MessageInteractionMetadata) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		d.Value = nil
		return nil
	}

	var a MessageInteractionMetadataApplicationCommand
	if err := json.Unmarshal(data, &a); err == nil && a.ID != "" {
		d.Value = &a
		return nil
	}

	var b MessageInteractionMetadataMessageComponent
	if err := json.Unmarshal(data, &b); err == nil && b.ID != "" {
		d.Value = &b
		return nil
	}

	var c MessageInteractionMetadataModalSubmit
	if err := json.Unmarshal(data, &c); err == nil && c.ID != "" {
		d.Value = &c
		return nil
	}

	return fmt.Errorf("unknown MessageInteractionMetadata: %s", string(data))
}

type MessageInteractionMetadataModalSubmit struct {
	ID                            Snowflake                                             `json:"id"`
	Type                          InteractionType                                       `json:"type"`
	User                          User                                                  `json:"user,omitempty"`
	AuthorizingIntegrationOwners  map[InteractionApplicationIntegrationType]interface{} `json:"authorizing_integration_owners,omitempty"`
	OriginalResponseMessageID     *Snowflake                                            `json:"original_response_message_id,omitempty"`
	TriggeringInteractionMetadata *MessageInteractionMetadataModalSubmitTriggering      `json:"triggering_interaction_metadata,omitempty"`
}

type InteractionContextType int

const (
	InteractionContextTypeGuild          InteractionContextType = 0
	InteractionContextTypeBotDM          InteractionContextType = 1
	InteractionContextTypePrivateChannel InteractionContextType = 2
)

type Interaction struct {
	ID                           Snowflake                                             `json:"id"`
	ApplicationID                Snowflake                                             `json:"application_id"`
	Type                         InteractionType                                       `json:"type"`
	Data                         InteractionData                                       `json:"data,omitempty"`
	GuildID                      *Snowflake                                            `json:"guild_id,omitempty"`
	ChannelID                    *Snowflake                                            `json:"channel_id,omitempty"`
	Guild                        *Guild                                                `json:"guild,omitempty"`
	Channel                      *Channel                                              `json:"channel,omitempty"`
	Member                       *GuildMember                                          `json:"member,omitempty"`
	User                         *User                                                 `json:"user,omitempty"`
	Token                        string                                                `json:"token"`
	Version                      int                                                   `json:"version"`
	Message                      *Message                                              `json:"message,omitempty"`
	AppPermissions               string                                                `json:"app_permissions,omitempty"`
	Locale                       *Locale                                               `json:"locale,omitempty"`
	GuildLocale                  string                                                `json:"guild_locale,omitempty"`
	Entitlements                 []Entitlement                                         `json:"entitlements,omitempty"`
	AuthorizingIntegrationOwners map[InteractionApplicationIntegrationType]interface{} `json:"authorizing_integration_owners,omitempty"`
	Context                      InteractionContextType                                `json:"context,omitempty"`
	AttachmentSizeLimit          int                                                   `json:"attachment_size_limit,omitempty"`
}

func (i *Interaction) GetSubCommand() string {
	if i.Data == nil {
		return ""
	}

	cmdData, ok := i.Data.(*InteractionDataApplicationCommand)
	if !ok {
		return ""
	}

	if cmdData.Options == nil || len(*cmdData.Options) == 0 {
		return ""
	}

	for _, option := range *cmdData.Options {
		if option.Type == ApplicationCommandOptionTypeSubCommand {
			return option.Name
		}

		if option.Type == ApplicationCommandOptionTypeSubCommandGroup {
			if option.Options != nil {
				for _, subOption := range option.Options {
					if subOption.Type == ApplicationCommandOptionTypeSubCommand {
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

	cmdData, ok := i.Data.(*InteractionDataApplicationCommand)
	if !ok {
		return ""
	}

	if cmdData.Options == nil || len(*cmdData.Options) == 0 {
		return ""
	}

	for _, option := range *cmdData.Options {
		if option.Type == ApplicationCommandOptionTypeSubCommandGroup {
			return option.Name
		}
	}

	return ""
}

func (i *Interaction) GetFullCommand() (fullCommand string) {
	if i.Data == nil {
		return ""
	}

	cmdData, ok := i.Data.(*InteractionDataApplicationCommand)
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

	componentData, ok := i.Data.(*InteractionDataMessageComponent)
	componentData2, ok2 := i.Data.(*InteractionDataModalSubmit)

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
		Type          ApplicationCommandType `json:"type"`
		ComponentType ComponentType          `json:"component_type"`
	}

	if err := json.Unmarshal(aux.Data, &typeProbe); err != nil {
		return err
	}

	switch typeProbe.Type {
	case ApplicationCommandTypeChatInput, ApplicationCommandTypeUser, ApplicationCommandTypeMessage:
		var cmd InteractionDataApplicationCommand
		if err := json.Unmarshal(aux.Data, &cmd); err != nil {
			return err
		}
		i.Data = &cmd
		return nil
	}

	switch typeProbe.ComponentType {
	case ComponentTypeButton, ComponentTypeStringSelectMenu, ComponentTypeUserSelectMenu, ComponentTypeRoleSelectMenu,
		ComponentTypeMentionableMenu, ComponentTypeChannelSelect:
		var comp InteractionDataMessageComponent
		if err := json.Unmarshal(aux.Data, &comp); err != nil {
			return err
		}
		i.Data = &comp
		return nil
	}

	switch aux.Type {
	case InteractionTypeModalSubmit:
		var modal InteractionDataModalSubmit
		if err := json.Unmarshal(aux.Data, &modal); err != nil {
			return err
		}
		i.Data = &modal
		return nil
	case InteractionTypeApplicationCommandAutocomplete:
		var auto InteractionDataAutocomplete
		if err := json.Unmarshal(aux.Data, &auto); err != nil {
			return err
		}
		i.Data = &auto
		return nil
	}

	return fmt.Errorf("unknown interaction data type %d", typeProbe.Type)
}

type InteractionCallbackType int

const (
	InteractionCallbackTypePong                                 InteractionCallbackType = 1
	InteractionCallbackTypeChannelMessageWithSource             InteractionCallbackType = 4
	InteractionCallbackTypeDeferredChannelMessageWithSource     InteractionCallbackType = 5
	InteractionCallbackTypeDeferredUpdateMessage                InteractionCallbackType = 6
	InteractionCallbackTypeUpdateMessage                        InteractionCallbackType = 7
	InteractionCallbackTypeApplicationCommandAutocompleteResult InteractionCallbackType = 8
	InteractionCallbackTypeModal                                InteractionCallbackType = 9
	InteractionCallbackTypePremiumRequired                      InteractionCallbackType = 10
	InteractionCallbackTypeLaunchActivity                       InteractionCallbackType = 12
)

type AllowedMentionsType string

const (
	AllowedMentionsTypeRoles    AllowedMentionsType = "roles"
	AllowedMentionsTypeUsers    AllowedMentionsType = "users"
	AllowedMentionsTypeEveryone AllowedMentionsType = "everyone"
)

type AllowedMentions struct {
	Parse       *[]AllowedMentionsType `json:"parse,omitempty"`
	Roles       *[]Snowflake           `json:"roles,omitempty"`
	Users       *[]Snowflake           `json:"users,omitempty"`
	RepliedUser *bool                  `json:"replied_user,omitempty"`
}

type AnyInteractionResponseData interface {
	IsInteractionResponseData() bool
	MarshalJSON() ([]byte, error)
}

type InteractionResponseDataDefault struct {
	TTS             bool             `json:"tts,omitempty"`
	Content         string           `json:"content,omitempty"`
	Embeds          *[]Embed         `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           MessageFlag      `json:"flags,omitempty"`
	Components      *[]AnyComponent  `json:"components,omitempty"`
	//TODO partial
	Attachment   *[]Attachment `json:"attachment,omitempty"`
	Poll         *PollRequest  `json:"poll,omitempty"`
	WithResponse bool          `json:"with_response,omitempty"`
}

func (d *InteractionResponseDataDefault) IsInteractionResponseData() bool {
	return true
}

func (d *InteractionResponseDataDefault) MarshalJSON() ([]byte, error) {
	type Alias InteractionResponseDataDefault
	return json.Marshal((*Alias)(d))
}

type InteractionResponse struct {
	Type InteractionCallbackType    `json:"type"`
	Data AnyInteractionResponseData `json:"data,omitempty"`
}

func (i *Interaction) DeferReply() error {
	return nil
}

func (i *Interaction) EditReply(responseData *AnyInteractionResponseData, clientID string) error {
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

func (i *Interaction) ReplyWithModal(modal *Modal) error {
	bodyBytes, err := json.Marshal(InteractionResponse{
		Type: InteractionCallbackTypeModal,
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

type InteractionCallback struct {
	ID                       Snowflake       `json:"id"`
	Type                     InteractionType `json:"type"`
	ActivityInstanceID       *Snowflake      `json:"activity_instance_id,omitempty"`
	ResponseMessageID        *Snowflake      `json:"response_message_id,omitempty"`
	ResponseMessageLoading   *bool           `json:"response_message_loading,omitempty"`
	ResponseMessageEphemeral *bool           `json:"response_message_ephemeral,omitempty"`
}

type InteractionCallbackActivityInstance struct {
	ID string `json:"id"`
}

type InteractionCallbackResource struct {
	Type             InteractionCallbackType              `json:"type"`
	ActivityInstance *InteractionCallbackActivityInstance `json:"activity_instance,omitempty"`
	Message          *Message                             `json:"message,omitempty"`
}

type InteractionCallbackResponse struct {
	Interaction InteractionCallback          `json:"interaction"`
	Resource    *InteractionCallbackResource `json:"resource"`
}

func (i *Interaction) Reply(data *InteractionResponseDataDefault) (*InteractionCallbackResponse, error) {
	bodyBytes, err := json.Marshal(InteractionResponse{
		Type: InteractionCallbackTypeChannelMessageWithSource,
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

	var resp InteractionCallbackResponse
	if err := json.NewDecoder(req.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type InteractionDataType int

const (
	InteractionDataTypePing                           InteractionDataType = 1
	InteractionDataTypeApplicationCommand             InteractionDataType = 2
	InteractionDataTypeMessageComponent               InteractionDataType = 3
	InteractionDataTypeApplicationCommandAutocomplete InteractionDataType = 4
	InteractionDataTypeModalSubmit                    InteractionDataType = 5
)

type InteractionData interface {
	GetType() InteractionDataType
}

type ApplicationCommandType int

const (
	ApplicationCommandTypeChatInput       ApplicationCommandType = 1
	ApplicationCommandTypeUser            ApplicationCommandType = 2
	ApplicationCommandTypeMessage         ApplicationCommandType = 3
	ApplicationCommandTypePrimaryEndpoint ApplicationCommandType = 4
)

type ApplicationCommandOptionType int

const (
	ApplicationCommandOptionTypeSubCommand      ApplicationCommandOptionType = 1
	ApplicationCommandOptionTypeSubCommandGroup ApplicationCommandOptionType = 2
	ApplicationCommandOptionTypeString          ApplicationCommandOptionType = 3
	ApplicationCommandOptionTypeInteger         ApplicationCommandOptionType = 4
	ApplicationCommandOptionTypeBoolean         ApplicationCommandOptionType = 5
	ApplicationCommandOptionTypeUser            ApplicationCommandOptionType = 6
	ApplicationCommandOptionTypeChannel         ApplicationCommandOptionType = 7
	ApplicationCommandOptionTypeRole            ApplicationCommandOptionType = 8
	ApplicationCommandOptionTypeMentionable     ApplicationCommandOptionType = 9
	ApplicationCommandOptionTypeNumber          ApplicationCommandOptionType = 10
	ApplicationCommandOptionTypeAttachment      ApplicationCommandOptionType = 11
)

type ApplicationCommandInteractionDataOption[T string | int | bool | interface{}] struct {
	Name    string                                                 `json:"name"`
	Type    ApplicationCommandOptionType                           `json:"type"`
	Value   *T                                                     `json:"value"`
	Options []ApplicationCommandInteractionDataOption[interface{}] `json:"options,omitempty"`
	Focused *bool                                                  `json:"focused,omitempty"`
}

func (t *ApplicationCommandInteractionDataOption[T]) UnmarshalJSON(data []byte) error {
	type Alias ApplicationCommandInteractionDataOption[T]
	raw := &struct {
		*Alias
		Value interface{} `json:"value,omitempty"`
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	t.Name = raw.Name
	t.Type = raw.Type
	t.Options = raw.Options
	t.Focused = raw.Focused

	if raw.Value != nil {
		switch t.Type {
		case ApplicationCommandOptionTypeString:
			if strVal, ok := raw.Value.(string); ok {
				var v T
				v = any(strVal).(T)
				t.Value = &v
			}
		case ApplicationCommandOptionTypeInteger:
			if intVal, ok := raw.Value.(int); ok {
				var v T
				v = any(intVal).(T)
				t.Value = &v
			}
		case ApplicationCommandOptionTypeBoolean:
			if boolVal, ok := raw.Value.(bool); ok {
				var v T
				v = any(boolVal).(T)
				t.Value = &v
			}
		default:
			var v T
			v = raw.Value.(T)
			t.Value = &v
		}
	}

	return nil
}

type InteractionDataMessageComponent struct {
	CustomID      string         `json:"custom_id"`
	ComponentType ComponentType  `json:"component_type"`
	Values        *[]interface{} `json:"values,omitempty"`
	Resolved      *ResolvedData  `json:"resolved,omitempty"`
}

func (d *InteractionDataMessageComponent) GetType() InteractionDataType {
	return InteractionDataTypeMessageComponent
}

type InteractionDataApplicationCommand struct {
	ID          Snowflake                                               `json:"id"`
	CommandName string                                                  `json:"name"`
	Type        ApplicationCommandType                                  `json:"type"`
	GuildID     *Snowflake                                              `json:"guild_id,omitempty"`
	TargetID    *Snowflake                                              `json:"target_id,omitempty"`
	Resolved    *ResolvedData                                           `json:"resolved,omitempty"`
	Options     *[]ApplicationCommandInteractionDataOption[interface{}] `json:"options,omitempty"`
}

func (d *InteractionDataApplicationCommand) GetType() InteractionDataType {
	return InteractionDataTypeApplicationCommand
}

func (d *InteractionDataApplicationCommand) UnmarshalJSON(data []byte) error {
	type Alias InteractionDataApplicationCommand
	raw := &struct {
		*Alias
		Options []json.RawMessage `json:"options,omitempty"`
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	d.ID = raw.ID
	d.CommandName = raw.CommandName
	d.Type = raw.Type
	d.GuildID = raw.GuildID
	d.TargetID = raw.TargetID
	d.Resolved = raw.Resolved

	if raw.Options != nil {
		var options []ApplicationCommandInteractionDataOption[interface{}]
		for _, optionData := range raw.Options {
			var option ApplicationCommandInteractionDataOption[interface{}]
			if err := json.Unmarshal(optionData, &option); err != nil {
				return err
			}
			options = append(options, option)

			if option.Options != nil {
				var option ApplicationCommandInteractionDataOption[interface{}]
				if err := json.Unmarshal(optionData, &option); err != nil {
					return err
				}

			}
		}
		d.Options = &options
	}

	return nil
}

type InteractionDataAutocomplete struct {
	ID          Snowflake                                               `json:"id"`
	CommandName string                                                  `json:"name"`
	Type        ApplicationCommandType                                  `json:"type"`
	GuildID     *Snowflake                                              `json:"guild_id,omitempty"`
	TargetID    *Snowflake                                              `json:"target_id,omitempty"`
	Resolved    *ResolvedData                                           `json:"resolved,omitempty"`
	Options     *[]ApplicationCommandInteractionDataOption[interface{}] `json:"options,omitempty"`
}

func (d *InteractionDataAutocomplete) GetType() InteractionDataType {
	return InteractionDataTypeApplicationCommand
}

type InteractionDataModalSubmit struct {
	CustomID   string                     `json:"custom_id"`
	Resolved   *ResolvedData              `json:"resolved,omitempty"`
	Components *[]ComponentLabelComponent `json:"components,omitempty"`
}

type ComponentLabelComponent struct {
	Type        ComponentType                    `json:"type"`
	ID          *int                             `json:"id,omitempty"`
	Label       *string                          `json:"label"`
	Description *string                          `json:"description,omitempty"`
	Component   *AnyComponentInteractionResponse `json:"component,omitempty"`
}

func (l *ComponentLabelComponent) UnmarshalJSON(data []byte) error {
	type Alias ComponentLabelComponent
	raw := &struct {
		*Alias
		Component *json.RawMessage `json:"component,omitempty"`
	}{
		Alias: (*Alias)(l),
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Component == nil {
		return nil
	}

	var probe struct {
		Type ComponentType `json:"type"`
	}
	if err := json.Unmarshal(*raw.Component, &probe); err != nil {
		return err
	}

	var c AnyComponentInteractionResponse

	switch probe.Type {
	case ComponentTypeUserSelectMenu:
		c = &UserSelectComponentInteractionResponse{}
	case ComponentTypeRoleSelectMenu:
		c = &RoleComponentInteractionResponse{}
	case ComponentTypeStringSelectMenu:
		c = &StringSelectComponentInteractionResponse{}
	case ComponentTypeChannelSelect:
		c = &ChannelComponentInteractionResponse{}
	case ComponentTypeMentionableMenu:
		c = &MentionableComponentInteractionResponse{}
	case ComponentTypeTextDisplay:
		c = &TextDisplayComponentInteractionResponse{}
	case ComponentTypeTextInput:
		c = &TextInputComponentInteractionResponse{}
	case ComponentTypeFileUpload:
		c = &FileUploadComponentInteractionResponse{}
	case ComponentTypeLabel:
		c = &LabelComponentInteractionResponse{}

	default:
		return fmt.Errorf("unknown interaction component type: %d", probe.Type)
	}

	if err := json.Unmarshal(*raw.Component, c); err != nil {
		return err
	}

	l.Component = &c

	return nil
}

func (d *InteractionDataModalSubmit) GetType() InteractionDataType {
	return InteractionDataTypeModalSubmit
}
