package discord

// https://docs.discord.com/developers/topics/oauth2#shared-resources-oauth2-scopes
type Scope string

const (
	ScopeActivitiesRead                       Scope = "activities.read"
	ScopeActivitiesWrite                      Scope = "activities.write"
	ScopeApplicationsBuildRead                Scope = "applications.build.read"
	ScopeApplicationsBuildUpload              Scope = "applications.build.upload"
	ScopeApplicationCommands                  Scope = "applications.commands"
	ScopeApplicationCommandsUpdate            Scope = "applications.commands.update"
	ScopeApplicationCommandsPermissionsUpdate Scope = "applications.commands.permissions.update"
	ScopeApplicationEntitlements              Scope = "applications.entitlements"
	ScopeApplicationsStoreUpdate              Scope = "applications.store.update"
	ScopeBot                                  Scope = "bot"
	ScopeConnections                          Scope = "connections"
	ScopeDMChannelsRead                       Scope = "dm_channels.read"
	ScopeEmail                                Scope = "email"
	ScopeGdmJoin                              Scope = "gdm.join"
	ScopeGuilds                               Scope = "guilds"
	ScopeGuildsJoin                           Scope = "guilds.join"
	ScopeGuildMembersRead                     Scope = "guild.members.read"
	ScopeIdentify                             Scope = "identify"
	ScopeIdentifyPremium                      Scope = "identify.premium"
	ScopeMessagesRead                         Scope = "messages.read"
	ScopeRelationshipsRead                    Scope = "relationships.read"
	ScopeRoleConnectionsWrite                 Scope = "role_connections.write"
	ScopeRPC                                  Scope = "rpc"
	ScopeRPCActivitiesWrite                   Scope = "rpc.activities.write"
	ScopeRPCNotificationsRead                 Scope = "rpc.notifications.read"
	ScopeRPCVoiceRead                         Scope = "rpc.voice.read"
	ScopeRPCVoiceWrite                        Scope = "rpc.voice.write"
	ScopeVoice                                Scope = "voice"
	ScopeWebhookIncoming                      Scope = "webhook.incoming"
)

const (
	OAuth2BaseAuthURL    = "https://discord.com/api/oauth2/authorize"
	OAuth2TokenURL       = "https://discord.com/api/oauth2/token"
	OAuth2TokenRevokeURL = "https://discord.com/api/oauth2/token/revoke"
)
