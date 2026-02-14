package common

type ApplicationEventWebhookStatus int

const (
	ApplicationEventWebhookStatusDisabled   ApplicationEventWebhookStatus = 0
	ApplicationEventWebhookStatusEnabled    ApplicationEventWebhookStatus = 1
	ApplicationEventWebhookStatusDisabledBy ApplicationEventWebhookStatus = 2
)

type ApplicationTeamMemberMembershipState int

const (
	ApplicationTeamMemberMembershipStateInvited  ApplicationTeamMemberMembershipState = 1
	ApplicationTeamMemberMembershipStateAccepted ApplicationTeamMemberMembershipState = 2
)

type ApplicationTeamMember struct {
	MembershipState ApplicationTeamMemberMembershipState `json:"membership_state"`
	TeamID          *string                              `json:"team_id,omitempty"`
	User            *User                                `json:"user,omitempty"`
	Role            string                               `json:"role"`
}

type ApplicationTeam struct {
	IconHash    *string                 `json:"icon,omitempty"`
	ID          Snowflake               `json:"id"`
	Members     []ApplicationTeamMember `json:"members"`
	Name        string                  `json:"name"`
	OwnerUserID Snowflake               `json:"owner_user_id"`
}

type ApplicationInstallParams struct {
	Scopes      []string `json:"scopes"`
	Permissions string   `json:"permissions"`
}

type Application struct {
	ID                                Snowflake                     `json:"id"`
	Name                              string                        `json:"name"`
	IconHash                          *string                       `json:"icon,omitempty"`
	Description                       string                        `json:"description,omitempty"`
	RpcOrigins                        []string                      `json:"rpc_origins,omitempty"`
	BotPublic                         bool                          `json:"bot_public,omitempty"`
	BotRequireCodeGrant               bool                          `json:"bot_require_code_grant,omitempty"`
	TermsOfServiceURL                 *string                       `json:"terms_of_service_url,omitempty"`
	PrivacyPolicyURL                  *string                       `json:"privacy_policy_url,omitempty"`
	Owner                             *User                         `json:"owner,omitempty"`
	VerifyKey                         *string                       `json:"verify_key,omitempty"`
	Team                              *ApplicationTeam              `json:"team,omitempty"`
	GuildID                           *Snowflake                    `json:"guild_id,omitempty"`
	Guild                             Guild                         `json:"guild,omitempty"`
	PrimarySKUID                      *string                       `json:"primary_sku_id,omitempty"`
	Slug                              *string                       `json:"slug,omitempty"`
	CoverImage                        *string                       `json:"cover_image,omitempty"`
	Flags                             *int                          `json:"flags,omitempty"`
	ApproximateGuildCount             *int                          `json:"approximate_guild_count,omitempty"`
	ApproximateUserInstallCount       *int                          `json:"approximate_user_install_count,omitempty"`
	ApproximateUserAuthorizationCount *int                          `json:"approximate_user_authorization_count,omitempty"`
	RedirectURIs                      *[]string                     `json:"redirect_uris,omitempty"`
	InteractionEndpointURL            *string                       `json:"interaction_endpoint_url,omitempty"`
	RoleConnectionsVerificationURL    *string                       `json:"role_connections_verification_url,omitempty"`
	EventWebhooksURL                  *string                       `json:"event_webhooks_url,omitempty"`
	EventWebhookStatus                ApplicationEventWebhookStatus `json:"event_webhook_status,omitempty"`
	EventWebhooksTypes                *[]string                     `json:"event_webhooks_types,omitempty"`
	Tags                              *[]string                     `json:"tags,omitempty"`
	InstallParams                     *ApplicationInstallParams     `json:"install_params,omitempty"`
	IntegrationTypesConfig            *interface{}                  `json:"integration_types_config,omitempty"`
	CustomInstallURL                  *string                       `json:"custom_install_url,omitempty"`
}
