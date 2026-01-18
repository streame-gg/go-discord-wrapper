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

type DiscordEntitlement struct {
	ID            DiscordSnowflake       `json:"id"`
	SkuID         DiscordSnowflake       `json:"sku_id"`
	ApplicationID DiscordSnowflake       `json:"application_id"`
	UserID        *DiscordSnowflake      `json:"user_id"`
	Type          DiscordEntitlementType `json:"type"`
	Deleted       bool                   `json:"deleted"`
	StartsAt      *time.Time             `json:"starts_at,omitempty"`
	EndsAt        *time.Time             `json:"ends_at,omitempty"`
	GuildID       *DiscordSnowflake      `json:"guild_id,omitempty"`
	Consumed      bool                   `json:"consumed,omitempty"`
}

type DiscordEntitlementType int

const (
	DiscordEntitlementTypePurchase                DiscordEntitlementType = 1
	DiscordEntitlementTypePremiumSubscription     DiscordEntitlementType = 2
	DiscordEntitlementTypeDeveloperGift           DiscordEntitlementType = 3
	DiscordEntitlementTypeTestModePurchase        DiscordEntitlementType = 4
	DiscordEntitlementTypeFreePurchase            DiscordEntitlementType = 5
	DiscordEntitlementTypeUserGift                DiscordEntitlementType = 6
	DiscordEntitlementTypePremiumPurchase         DiscordEntitlementType = 7
	DiscordEntitlementTypeApplicationSubscription DiscordEntitlementType = 8
)

type DiscordInteractionType int

const (
	DiscordInteractionTypePing                           DiscordInteractionType = 1
	DiscordInteractionTypeApplicationCommand             DiscordInteractionType = 2
	DiscordInteractionTypeMessageComponent               DiscordInteractionType = 3
	DiscordInteractionTypeApplicationCommandAutocomplete DiscordInteractionType = 4
	DiscordInteractionTypeModalSubmit                    DiscordInteractionType = 5
)

type DiscordInteractionApplicationIntegrationType int

const (
	DiscordInteractionApplicationIntegrationTypeGuildInstall DiscordInteractionApplicationIntegrationType = 0
	DiscordInteractionApplicationIntegrationTypeUserInstall  DiscordInteractionApplicationIntegrationType = 1
)

type DiscordMessageMessageReference struct {
	Type            *DiscordMessageMessageReferenceType `json:"type,omitempty"`
	MessageID       *DiscordSnowflake                   `json:"message_id,omitempty"`
	ChannelID       *DiscordSnowflake                   `json:"channel_id,omitempty"`
	GuildID         *DiscordSnowflake                   `json:"guild_id,omitempty"`
	FailIfNotExists *bool                               `json:"fail_if_not_exists,omitempty"`
}

type DiscordMessageInteractionMetadataApplicationCommand struct {
	ID                           DiscordSnowflake                                             `json:"id"`
	Type                         DiscordInteractionType                                       `json:"type"`
	User                         DiscordUser                                                  `json:"user,omitempty"`
	AuthorizingIntegrationOwners map[DiscordInteractionApplicationIntegrationType]interface{} `json:"authorizing_integration_owners,omitempty"`
	OriginalResponseMessageID    *DiscordSnowflake                                            `json:"original_response_message_id,omitempty"`
	TargetUser                   *DiscordUser                                                 `json:"target_user,omitempty"`
	TargetMessageID              *DiscordSnowflake                                            `json:"target_message_id,omitempty"`
}

type DiscordMessageInteractionMetadataMessageComponent struct {
	ID                           DiscordSnowflake                                             `json:"id"`
	Type                         DiscordInteractionType                                       `json:"type"`
	User                         DiscordUser                                                  `json:"user,omitempty"`
	AuthorizingIntegrationOwners map[DiscordInteractionApplicationIntegrationType]interface{} `json:"authorizing_integration_owners,omitempty"`
	OriginalResponseMessageID    *DiscordSnowflake                                            `json:"original_response_message_id,omitempty"`
	InteractedMessageID          *DiscordSnowflake                                            `json:"interacted_message_id,omitempty"`
}

type AnyDiscordMessageInteractionMetadata interface{}

type DiscordMessageInteractionMetadata struct {
	Value AnyDiscordMessageInteractionMetadata
}

type AnyDiscordMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata interface{}

type DiscordMessageInteractionMetadataModalSubmitTriggering struct {
	AnyDiscordMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata
}

func (d *DiscordMessageInteractionMetadataModalSubmitTriggering) UnmarshalJSON(data []byte) error {
	var a DiscordMessageInteractionMetadataApplicationCommand
	if err := json.Unmarshal(data, &a); err == nil && a.ID != "" {
		d.AnyDiscordMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata = &a
		return nil
	}

	var b DiscordMessageInteractionMetadataMessageComponent
	if err := json.Unmarshal(data, &b); err == nil && b.ID != "" {
		d.AnyDiscordMessageInteractionMetadataModalSubmitTriggeringInteractionMetadata = &b
		return nil
	}

	return nil
}

func (d *DiscordMessageInteractionMetadata) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		d.Value = nil
		return nil
	}

	var a DiscordMessageInteractionMetadataApplicationCommand
	if err := json.Unmarshal(data, &a); err == nil && a.ID != "" {
		d.Value = &a
		return nil
	}

	var b DiscordMessageInteractionMetadataMessageComponent
	if err := json.Unmarshal(data, &b); err == nil && b.ID != "" {
		d.Value = &b
		return nil
	}

	var c DiscordMessageInteractionMetadataModalSubmit
	if err := json.Unmarshal(data, &c); err == nil && c.ID != "" {
		d.Value = &c
		return nil
	}

	return fmt.Errorf("unknown DiscordMessageInteractionMetadata: %s", string(data))
}

type DiscordMessageInteractionMetadataModalSubmit struct {
	ID                            DiscordSnowflake                                             `json:"id"`
	Type                          DiscordInteractionType                                       `json:"type"`
	User                          DiscordUser                                                  `json:"user,omitempty"`
	AuthorizingIntegrationOwners  map[DiscordInteractionApplicationIntegrationType]interface{} `json:"authorizing_integration_owners,omitempty"`
	OriginalResponseMessageID     *DiscordSnowflake                                            `json:"original_response_message_id,omitempty"`
	TriggeringInteractionMetadata *DiscordMessageInteractionMetadataModalSubmitTriggering      `json:"triggering_interaction_metadata,omitempty"`
}

type InteractionContextType int

const (
	InteractionContextTypeGuild          InteractionContextType = 0
	InteractionContextTypeBotDM          InteractionContextType = 1
	InteractionContextTypePrivateChannel InteractionContextType = 2
)

type DiscordInteraction struct {
	ID                           DiscordSnowflake                                             `json:"id"`
	ApplicationID                DiscordSnowflake                                             `json:"application_id"`
	Type                         DiscordInteractionType                                       `json:"type"`
	Data                         DiscordInteractionData                                       `json:"data,omitempty"`
	GuildID                      *DiscordSnowflake                                            `json:"guild_id,omitempty"`
	ChannelID                    *DiscordSnowflake                                            `json:"channel_id,omitempty"`
	Guild                        *Guild                                                       `json:"guild,omitempty"`
	Channel                      *DiscordChannel                                              `json:"channel,omitempty"`
	Member                       *GuildMember                                                 `json:"member,omitempty"`
	User                         *DiscordUser                                                 `json:"user,omitempty"`
	Token                        string                                                       `json:"token"`
	Version                      int                                                          `json:"version"`
	Message                      *DiscordMessage                                              `json:"message,omitempty"`
	AppPermissions               string                                                       `json:"app_permissions,omitempty"`
	Locale                       *DiscordLocale                                               `json:"locale,omitempty"`
	GuildLocale                  string                                                       `json:"guild_locale,omitempty"`
	Entitlements                 []DiscordEntitlement                                         `json:"entitlements,omitempty"`
	AuthorizingIntegrationOwners map[DiscordInteractionApplicationIntegrationType]interface{} `json:"authorizing_integration_owners,omitempty"`
	Context                      InteractionContextType                                       `json:"context,omitempty"`
	AttachmentSizeLimit          int                                                          `json:"attachment_size_limit,omitempty"`
}

func (i *DiscordInteraction) GetSubCommand() string {
	if i.Data == nil {
		return ""
	}

	cmdData, ok := i.Data.(*DiscordInteractionDataApplicationCommand)
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

func (i *DiscordInteraction) GetSubCommandGroup() string {
	if i.Data == nil {
		return ""
	}

	cmdData, ok := i.Data.(*DiscordInteractionDataApplicationCommand)
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

func (i *DiscordInteraction) GetFullCommand() (fullCommand string) {
	if i.Data == nil {
		return ""
	}

	cmdData, ok := i.Data.(*DiscordInteractionDataApplicationCommand)
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

func (i *DiscordInteraction) GetCustomID() string {
	if i.Data == nil {
		return ""
	}

	componentData, ok := i.Data.(*DiscordInteractionDataMessageComponent)
	componentData2, ok2 := i.Data.(*DiscordInteractionDataModalSubmit)

	if !ok && !ok2 {
		return ""
	}

	if ok2 {
		return componentData2.CustomID
	}

	return componentData.CustomID
}

func (i *DiscordInteraction) UnmarshalJSON(data []byte) error {
	type Alias DiscordInteraction
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
		Type          DiscordInteractionDataApplicationCommandType `json:"type"`
		ComponentType DiscordComponentType                         `json:"component_type"`
	}

	if err := json.Unmarshal(aux.Data, &typeProbe); err != nil {
		return err
	}

	switch typeProbe.Type {
	case DiscordInteractionDataApplicationCommandTypeChatInput, DiscordInteractionDataApplicationCommandTypeUser, DiscordInteractionDataApplicationCommandTypeMessage:
		var cmd DiscordInteractionDataApplicationCommand
		if err := json.Unmarshal(aux.Data, &cmd); err != nil {
			return err
		}
		i.Data = &cmd
		return nil
	}

	switch typeProbe.ComponentType {
	case DiscordComponentTypeButton, DiscordComponentTypeStringSelectMenu, DiscordComponentTypeUserSelectMenu, DiscordComponentTypeRoleSelectMenu,
		DiscordComponentTypeMentionableMenu, DiscordComponentTypeChannelSelect:
		var comp DiscordInteractionDataMessageComponent
		if err := json.Unmarshal(aux.Data, &comp); err != nil {
			return err
		}
		i.Data = &comp
		return nil
	}

	switch aux.Type {
	case DiscordInteractionTypeModalSubmit:
		var modal DiscordInteractionDataModalSubmit
		if err := json.Unmarshal(aux.Data, &modal); err != nil {
			return err
		}
		i.Data = &modal
		return nil
	case DiscordInteractionTypeApplicationCommandAutocomplete:
		var auto DiscordInteractionDataAutocomplete
		if err := json.Unmarshal(aux.Data, &auto); err != nil {
			return err
		}
		i.Data = &auto
		return nil
	}

	return fmt.Errorf("unknown interaction data type %d", typeProbe.Type)
}

type DiscordInteractionCallbackType int

const (
	DiscordInteractionCallbackTypePong                                 DiscordInteractionCallbackType = 1
	DiscordInteractionCallbackTypeChannelMessageWithSource             DiscordInteractionCallbackType = 4
	DiscordInteractionCallbackTypeDeferredChannelMessageWithSource     DiscordInteractionCallbackType = 5
	DiscordInteractionCallbackTypeDeferredUpdateMessage                DiscordInteractionCallbackType = 6
	DiscordInteractionCallbackTypeUpdateMessage                        DiscordInteractionCallbackType = 7
	DiscordInteractionCallbackTypeApplicationCommandAutocompleteResult DiscordInteractionCallbackType = 8
	DiscordInteractionCallbackTypeModal                                DiscordInteractionCallbackType = 9
	DiscordInteractionCallbackTypePremiumRequired                      DiscordInteractionCallbackType = 10
	DiscordInteractionCallbackTypeLaunchActivity                       DiscordInteractionCallbackType = 12
)

type DiscordAllowedMentionsType string

const (
	DiscordAllowedMentionsTypeRoles    DiscordAllowedMentionsType = "roles"
	DiscordAllowedMentionsTypeUsers    DiscordAllowedMentionsType = "users"
	DiscordAllowedMentionsTypeEveryone DiscordAllowedMentionsType = "everyone"
)

type DiscordAllowedMentions struct {
	Parse       *[]DiscordAllowedMentionsType `json:"parse,omitempty"`
	Roles       *[]DiscordSnowflake           `json:"roles,omitempty"`
	Users       *[]DiscordSnowflake           `json:"users,omitempty"`
	RepliedUser *bool                         `json:"replied_user,omitempty"`
}

type DiscordInteractionResponseData struct {
	TTS             bool                    `json:"tts,omitempty"`
	Content         string                  `json:"content,omitempty"`
	Embeds          *[]DiscordEmbed         `json:"embeds,omitempty"`
	AllowedMentions *DiscordAllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           DiscordMessageFlag      `json:"flags,omitempty"`
	Components      *Components             `json:"components,omitempty"`
	//TODO partial
	Attachment   *[]DiscordAttachment `json:"attachment,omitempty"`
	Poll         *DiscordPollRequest  `json:"poll,omitempty"`
	WithResponse bool                 `json:"with_response,omitempty"`
}

type DiscordInteractionResponse struct {
	Type DiscordInteractionCallbackType  `json:"type"`
	Data *DiscordInteractionResponseData `json:"data,omitempty"`
}

func (i *DiscordInteraction) CreateInteractionResponse(responseData *DiscordInteractionResponse) (*any, error) {
	bodyBytes, err := json.Marshal(*responseData)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(req.Body)

	if !responseData.Data.WithResponse {
		if req.StatusCode != 204 {
			var respErr map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&respErr); err != nil {
				return nil, err
			}

			return nil, fmt.Errorf("expected 204 No Content, got %d: %v", req.StatusCode, respErr)
		}

		return nil, nil
	}

	if responseData.Data.WithResponse && req.StatusCode != 200 {
		var respErr map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&respErr); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("expected 204 No Content, got %d: %v", req.StatusCode, respErr)
	}

	var resp any
	if err := json.NewDecoder(req.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type DiscordInteractionDataType int

const (
	DiscordInteractionDataTypePing                           DiscordInteractionDataType = 1
	DiscordInteractionDataTypeApplicationCommand             DiscordInteractionDataType = 2
	DiscordInteractionDataTypeMessageComponent               DiscordInteractionDataType = 3
	DiscordInteractionDataTypeApplicationCommandAutocomplete DiscordInteractionDataType = 4
	DiscordInteractionDataTypeModalSubmit                    DiscordInteractionDataType = 5
)

type DiscordInteractionData interface {
	GetType() DiscordInteractionDataType
}

type DiscordInteractionDataApplicationCommandType int

const (
	DiscordInteractionDataApplicationCommandTypeChatInput       DiscordInteractionDataApplicationCommandType = 1
	DiscordInteractionDataApplicationCommandTypeUser            DiscordInteractionDataApplicationCommandType = 2
	DiscordInteractionDataApplicationCommandTypeMessage         DiscordInteractionDataApplicationCommandType = 3
	DiscordInteractionDataApplicationCommandTypePrimaryEndpoint DiscordInteractionDataApplicationCommandType = 4
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

type ApplicationCommandInteractionDataOption struct {
	Name    string                                    `json:"name"`
	Type    ApplicationCommandOptionType              `json:"type"`
	Value   *interface{}                              `json:"value,omitempty"`
	Options []ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	Focused *bool                                     `json:"focused,omitempty"`
}

type DiscordInteractionDataMessageComponent struct {
	CustomID      string               `json:"custom_id"`
	ComponentType DiscordComponentType `json:"component_type"`
	Values        *[]interface{}       `json:"values,omitempty"`
	Resolved      *DiscordResolvedData `json:"resolved,omitempty"`
}

func (d *DiscordInteractionDataMessageComponent) GetType() DiscordInteractionDataType {
	return DiscordInteractionDataTypeMessageComponent
}

type DiscordInteractionDataApplicationCommand struct {
	ID          DiscordSnowflake                             `json:"id"`
	CommandName string                                       `json:"name"`
	Type        DiscordInteractionDataApplicationCommandType `json:"type"`
	GuildID     *DiscordSnowflake                            `json:"guild_id,omitempty"`
	TargetID    *DiscordSnowflake                            `json:"target_id,omitempty"`
	Resolved    *DiscordResolvedData                         `json:"resolved,omitempty"`
	Options     *[]ApplicationCommandInteractionDataOption   `json:"options,omitempty"`
}

func (d *DiscordInteractionDataApplicationCommand) GetType() DiscordInteractionDataType {
	return DiscordInteractionDataTypeApplicationCommand
}

type DiscordInteractionDataAutocomplete struct {
	ID          DiscordSnowflake                             `json:"id"`
	CommandName string                                       `json:"name"`
	Type        DiscordInteractionDataApplicationCommandType `json:"type"`
	GuildID     *DiscordSnowflake                            `json:"guild_id,omitempty"`
	TargetID    *DiscordSnowflake                            `json:"target_id,omitempty"`
	Resolved    *DiscordResolvedData                         `json:"resolved,omitempty"`
	Options     *[]ApplicationCommandInteractionDataOption   `json:"options,omitempty"`
}

func (d *DiscordInteractionDataAutocomplete) GetType() DiscordInteractionDataType {
	return DiscordInteractionDataTypeApplicationCommand
}

type DiscordInteractionDataModalSubmit struct {
	CustomID string               `json:"custom_id"`
	Resolved *DiscordResolvedData `json:"resolved,omitempty"`
}

func (d *DiscordInteractionDataModalSubmit) GetType() DiscordInteractionDataType {
	return DiscordInteractionDataTypeModalSubmit
}
