package discord

// https://docs.discord.com/developers/resources/application#application-object-application-event-webhook-status
type ApplicationEventWebhookStatus int

const (
	ApplicationEventWebhookStatusDisabled          ApplicationEventWebhookStatus = 1
	ApplicationEventWebhookStatusEnabled           ApplicationEventWebhookStatus = 2
	ApplicationEventWebhookStatusDisabledByDiscord ApplicationEventWebhookStatus = 3
)

// https://docs.discord.com/developers/topics/teams#data-models-membership-state-enum
type ApplicationTeamMemberMembershipState int

const (
	ApplicationTeamMemberMembershipStateInvited  ApplicationTeamMemberMembershipState = 1
	ApplicationTeamMemberMembershipStateAccepted ApplicationTeamMemberMembershipState = 2
)

// https://docs.discord.com/developers/topics/teams#data-models-team-member-object
type ApplicationTeamMember struct {
	MembershipState ApplicationTeamMemberMembershipState `json:"membership_state"`
	TeamID          Snowflake                            `json:"team_id"`
	User            User                                 `json:"user"`
	Role            string                               `json:"role"`
}

// https://docs.discord.com/developers/topics/teams#data-models-team-object
type ApplicationTeam struct {
	Icon        *string                 `json:"icon"`
	ID          Snowflake               `json:"id"`
	Members     []ApplicationTeamMember `json:"members"`
	Name        string                  `json:"name"`
	OwnerUserID Snowflake               `json:"owner_user_id"`
}

// https://docs.discord.com/developers/resources/application#install-params-object
type ApplicationInstallParams struct {
	Scopes      []Scope    `json:"scopes"`
	Permissions Permission `json:"permissions"`
}

// https://docs.discord.com/developers/resources/application-role-connection-metadata#application-role-connection-metadata-object-application-role-connection-metadata-type
type ApplicationRoleConnectionsMetadataType int

const (
	ApplicationRoleConnectionsMetadataTypeIntegerLessThanOrEqual     ApplicationRoleConnectionsMetadataType = 1
	ApplicationRoleConnectionsMetadataTypeIntegerGreaterThanOrEqual  ApplicationRoleConnectionsMetadataType = 2
	ApplicationRoleConnectionsMetadataTypeIntegerEqual               ApplicationRoleConnectionsMetadataType = 3
	ApplicationRoleConnectionsMetadataTypeIntegerNotEqual            ApplicationRoleConnectionsMetadataType = 4
	ApplicationRoleConnectionsMetadataTypeDatetimeLessThanOrEqual    ApplicationRoleConnectionsMetadataType = 5
	ApplicationRoleConnectionsMetadataTypeDatetimeGreaterThanOrEqual ApplicationRoleConnectionsMetadataType = 6
	ApplicationRoleConnectionsMetadataTypeBooleanEqual               ApplicationRoleConnectionsMetadataType = 7
	ApplicationRoleConnectionsMetadataTypeBooleanNotEqual            ApplicationRoleConnectionsMetadataType = 8
)

// https://docs.discord.com/developers/resources/application-role-connection-metadata#application-role-connection-metadata-object
type ApplicationRoleConnectionMetadata struct {
	Type                     ApplicationRoleConnectionsMetadataType `json:"type"`
	Key                      string                                 `json:"key"`
	Name                     string                                 `json:"name"`
	NameLocalizations        map[Locale]string                      `json:"name_localizations,omitempty"`
	Description              string                                 `json:"description"`
	DescriptionLocalizations map[Locale]string                      `json:"description_localizations,omitempty"`
}

// https://docs.discord.com/developers/resources/application#application-object-application-integration-types
type ApplicationIntegrationTypeConfiguration struct {
	OAuth2InstallParams *ApplicationInstallParams `json:"oauth2_install_params"`
}

// https://docs.discord.com/developers/resources/application#application-object-application-integration-types
type ApplicationIntegrationType string

const (
	ApplicationIntegrationTypeGuildInstall ApplicationIntegrationType = "0"
	ApplicationIntegrationTypeUserInstall  ApplicationIntegrationType = "1"
)

// https://docs.discord.com/developers/resources/application#application-object
type Application struct {
	ID                                Snowflake                                                              `json:"id"`
	Name                              string                                                                 `json:"name"`
	Icon                              *string                                                                `json:"icon"`
	Description                       string                                                                 `json:"description"`
	RpcOrigins                        []string                                                               `json:"rpc_origins,omitempty"`
	BotPublic                         bool                                                                   `json:"bot_public"`
	BotRequireCodeGrant               bool                                                                   `json:"bot_require_code_grant"`
	Bot                               *User                                                                  `json:"bot,omitempty"`
	TermsOfServiceURL                 string                                                                 `json:"terms_of_service_url,omitempty"`
	PrivacyPolicyURL                  string                                                                 `json:"privacy_policy_url,omitempty"`
	Owner                             *User                                                                  `json:"owner,omitempty"`
	VerifyKey                         string                                                                 `json:"verify_key"`
	Team                              *ApplicationTeam                                                       `json:"team"`
	GuildID                           *Snowflake                                                             `json:"guild_id,omitempty"`
	Guild                             *Guild                                                                 `json:"guild,omitempty"`
	PrimarySKUID                      *Snowflake                                                             `json:"primary_sku_id,omitempty"`
	Slug                              string                                                                 `json:"slug,omitempty"`
	CoverImage                        string                                                                 `json:"cover_image,omitempty"`
	Flags                             *ApplicationFlags                                                      `json:"flags,omitempty"`
	FlagsNew                          *ApplicationFlags                                                      `json:"flags_new,omitempty"`
	ApproximateGuildCount             *int                                                                   `json:"approximate_guild_count,omitempty"`
	ApproximateUserInstallCount       *int                                                                   `json:"approximate_user_install_count,omitempty"`
	ApproximateUserAuthorizationCount *int                                                                   `json:"approximate_user_authorization_count,omitempty"`
	RedirectURIs                      []string                                                               `json:"redirect_uris,omitempty"`
	InteractionsEndpointURL           *string                                                                `json:"interactions_endpoint_url,omitempty"`
	RoleConnectionsVerificationURL    *string                                                                `json:"role_connections_verification_url,omitempty"`
	EventWebhooksURL                  *string                                                                `json:"event_webhooks_url,omitempty"`
	EventWebhooksStatus               *ApplicationEventWebhookStatus                                         `json:"event_webhooks_status,omitempty"`
	EventWebhooksTypes                []string                                                               `json:"event_webhooks_types,omitempty"`
	Tags                              []string                                                               `json:"tags,omitempty"`
	InstallParams                     *ApplicationInstallParams                                              `json:"install_params,omitempty"`
	IntegrationTypesConfig            map[ApplicationIntegrationType]ApplicationIntegrationTypeConfiguration `json:"integration_types_config,omitempty"`
	CustomInstallURL                  string                                                                 `json:"custom_install_url,omitempty"`
}
