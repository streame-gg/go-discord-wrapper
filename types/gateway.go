package types

import (
	"encoding/json"
)

type APIVersion uint8

func (a APIVersion) ToString() string {
	switch a {
	case APIVersion10:
		return "10"
	case APIVersion9:
		return "9"
	default:
		return "unknown"
	}
}

var (
	APIBaseString = func(v APIVersion) string {
		return "/api/v" + v.ToString() + "/"
	}

	APIVersion10 APIVersion = 10
	APIVersion9  APIVersion = 9
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

type InvalidSessionPayload struct {
	D bool `json:"d"`
}

type Intent uint64

const (
	IntentGuilds                       Intent = 1 << 0
	IntentGuildMembers                 Intent = 1 << 1
	IntentGuildModeration              Intent = 1 << 2
	IntentGuildExpressions             Intent = 1 << 3
	IntentGuildIntegrations            Intent = 1 << 4
	IntentGuildWebhooks                Intent = 1 << 5
	IntentGuildInvites                 Intent = 1 << 6
	IntentGuildVoiceStates             Intent = 1 << 7
	IntentGuildPresences               Intent = 1 << 8
	IntentGuildMessages                Intent = 1 << 9
	IntentGuildMessageReactions        Intent = 1 << 10
	IntentGuildMessageTyping           Intent = 1 << 11
	IntentDirectMessages               Intent = 1 << 12
	IntentDirectMessageReactions       Intent = 1 << 13
	IntentDirectMessageTyping          Intent = 1 << 14
	IntentMessageContent               Intent = 1 << 15
	IntentGuildScheduledEvents         Intent = 1 << 16
	IntentGuildModerationConfiguration Intent = 1 << 20
	IntentGuildModerationExecution     Intent = 1 << 21
	IntentMessagePolls                 Intent = 1 << 24
	IntentDirectMessagePolls           Intent = 1 << 25

	AllIntents = IntentGuilds | IntentGuildMembers |
		IntentGuildModeration | IntentGuildExpressions |
		IntentGuildIntegrations | IntentGuildWebhooks |
		IntentGuildInvites | IntentGuildVoiceStates |
		IntentGuildPresences | IntentGuildMessages |
		IntentGuildMessageReactions | IntentGuildMessageTyping |
		IntentDirectMessages | IntentDirectMessageReactions |
		IntentDirectMessageTyping | IntentMessageContent |
		IntentGuildScheduledEvents | IntentGuildModerationConfiguration |
		IntentGuildModerationExecution | IntentMessagePolls |
		IntentDirectMessagePolls

	AllIntentsExceptDirectMessage = IntentGuilds | IntentGuildMembers |
		IntentGuildModeration | IntentGuildExpressions |
		IntentGuildIntegrations | IntentGuildWebhooks |
		IntentGuildInvites | IntentGuildVoiceStates |
		IntentGuildPresences | IntentGuildMessages |
		IntentGuildMessageReactions | IntentGuildMessageTyping |
		IntentMessageContent | IntentGuildScheduledEvents |
		IntentGuildModerationConfiguration |
		IntentGuildModerationExecution | IntentMessagePolls
)
