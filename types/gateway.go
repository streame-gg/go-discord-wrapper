package types

import (
	"encoding/json"
	"fmt"
)

type DiscordAPIVersion uint8

var (
	DiscordAPIBaseString = func(v DiscordAPIVersion) string {
		return fmt.Sprintf("/api/v%d/", v)
	}

	DiscordAPIVersion10 DiscordAPIVersion = 10
	DiscordAPIVersion9  DiscordAPIVersion = 9

	DiscordAPIGatewayRequest = "gateway/bot"
)

type Payload struct {
	Op PayloadOpCode   `json:"op"`
	D  json.RawMessage `json:"d"`
	T  string          `json:"t,omitempty"`
	S  *int            `json:"s,omitempty"`
}

type HelloPayloadData struct {
	HeartbeatInterval float64 `json:"heartbeat_interval"`
}

type PayloadOpCode int

const (
	PayloadOpCodeDispatch            PayloadOpCode = 0
	PayloadOpCodeHeartbeat           PayloadOpCode = 1
	PayloadOpCodeIdentify            PayloadOpCode = 2
	PayloadOpCodePresenceUpdate      PayloadOpCode = 3
	PayloadOpCodeVoiceStateUpdate    PayloadOpCode = 4
	PayloadOpCodeResume              PayloadOpCode = 6
	PayloadOpCodeReconnect           PayloadOpCode = 7
	PayloadOpCodeRequestGuildMembers PayloadOpCode = 8
	PayloadOpCodeInvalidSession      PayloadOpCode = 9
	PayloadOpCodeHello               PayloadOpCode = 10
	PayloadOpCodeHeartbeatACK        PayloadOpCode = 11
	PayloadRequestSoundboardSounds   PayloadOpCode = 31
)

type DiscordReadyPayload struct {
	User             DiscordUser       `json:"user"`
	SessionID        string            `json:"session_id"`
	ResumeGatewayURL string            `json:"resume_gateway_url"`
	Shard            []int             `json:"shard,omitempty"`
	Guilds           []AnyGuildWrapper `json:"guilds"`
}

type DiscordIntent uint64

const (
	DiscordIntentGuilds                       DiscordIntent = 1 << 0
	DiscordIntentGuildMembers                 DiscordIntent = 1 << 1
	DiscordIntentGuildModeration              DiscordIntent = 1 << 2
	DiscordIntentGuildExpressions             DiscordIntent = 1 << 3
	DiscordIntentGuildIntegrations            DiscordIntent = 1 << 4
	DiscordIntentGuildWebhooks                DiscordIntent = 1 << 5
	DiscordIntentGuildInvites                 DiscordIntent = 1 << 6
	DiscordIntentGuildVoiceStates             DiscordIntent = 1 << 7
	DiscordIntentGuildPresences               DiscordIntent = 1 << 8
	DiscordIntentGuildMessages                DiscordIntent = 1 << 9
	DiscordIntentGuildMessageReactions        DiscordIntent = 1 << 10
	DiscordIntentGuildMessageTyping           DiscordIntent = 1 << 11
	DiscordIntentDirectMessages               DiscordIntent = 1 << 12
	DiscordIntentDirectMessageReactions       DiscordIntent = 1 << 13
	DiscordIntentDirectMessageTyping          DiscordIntent = 1 << 14
	DiscordIntentMessageContent               DiscordIntent = 1 << 15
	DiscordIntentGuildScheduledEvents         DiscordIntent = 1 << 16
	DiscordIntentGuildModerationConfiguration DiscordIntent = 1 << 20
	DiscordIntentGuildModerationExecution     DiscordIntent = 1 << 21
	DiscordIntentMessagePolls                 DiscordIntent = 1 << 24
	DiscordIntentDirectMessagePolls           DiscordIntent = 1 << 25

	AllIntents = DiscordIntentGuilds | DiscordIntentGuildMembers |
		DiscordIntentGuildModeration | DiscordIntentGuildExpressions |
		DiscordIntentGuildIntegrations | DiscordIntentGuildWebhooks |
		DiscordIntentGuildInvites | DiscordIntentGuildVoiceStates |
		DiscordIntentGuildPresences | DiscordIntentGuildMessages |
		DiscordIntentGuildMessageReactions | DiscordIntentGuildMessageTyping |
		DiscordIntentDirectMessages | DiscordIntentDirectMessageReactions |
		DiscordIntentDirectMessageTyping | DiscordIntentMessageContent |
		DiscordIntentGuildScheduledEvents | DiscordIntentGuildModerationConfiguration |
		DiscordIntentGuildModerationExecution | DiscordIntentMessagePolls |
		DiscordIntentDirectMessagePolls

	AllIntentsExceptDirectMessage = DiscordIntentGuilds | DiscordIntentGuildMembers |
		DiscordIntentGuildModeration | DiscordIntentGuildExpressions |
		DiscordIntentGuildIntegrations | DiscordIntentGuildWebhooks |
		DiscordIntentGuildInvites | DiscordIntentGuildVoiceStates |
		DiscordIntentGuildPresences | DiscordIntentGuildMessages |
		DiscordIntentGuildMessageReactions | DiscordIntentGuildMessageTyping |
		DiscordIntentMessageContent | DiscordIntentGuildScheduledEvents |
		DiscordIntentGuildModerationConfiguration |
		DiscordIntentGuildModerationExecution | DiscordIntentMessagePolls
)
