package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

var _ Event = &AutoModerationRuleCreateEvent{}
var _ Event = &AutoModerationRuleUpdateEvent{}
var _ Event = &AutoModerationRuleDeleteEvent{}
var _ Event = &AutoModerationActionExecutionEvent{}

func init() {
	RegisterEvent(&AutoModerationRuleCreateEvent{})
	RegisterEvent(&AutoModerationRuleUpdateEvent{})
	RegisterEvent(&AutoModerationRuleDeleteEvent{})
	RegisterEvent(&AutoModerationActionExecutionEvent{})
}

// https://docs.discord.com/developers/events/gateway-events#auto-moderation-rule-create
type AutoModerationRuleCreateEvent struct {
	Rule discord.AutoModerationRule `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#auto-moderation-rule-update
type AutoModerationRuleUpdateEvent struct {
	NewRule discord.AutoModerationRule `json:"-"`
	// OldRule is nil when the rule was not previously cached.
	OldRule *discord.AutoModerationRule `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#auto-moderation-rule-delete
type AutoModerationRuleDeleteEvent struct {
	Rule discord.AutoModerationRule `json:"-"`
}

// https://docs.discord.com/developers/events/gateway-events#auto-moderation-action-execution
type AutoModerationActionExecutionEvent struct {
	GuildID              discord.Snowflake                 `json:"guild_id"`
	Action               discord.AutoModerationAction      `json:"action"`
	RuleID               discord.Snowflake                 `json:"rule_id"`
	RuleTriggerType      discord.AutoModerationTriggerType `json:"rule_trigger_type"`
	UserID               discord.Snowflake                 `json:"user_id"`
	ChannelID            *discord.Snowflake                `json:"channel_id,omitempty"`
	MessageID            *discord.Snowflake                `json:"message_id,omitempty"`
	AlertSystemMessageID *discord.Snowflake                `json:"alert_system_message_id,omitempty"`
	Content              string                            `json:"content"`
	MatchedKeyword       *string                           `json:"matched_keyword"`
	MatchedContent       *string                           `json:"matched_content"`

	Guild              *discord.Guild              `json:"-"`
	AutomoderationRule *discord.AutoModerationRule `json:"-"`
}

func (e *AutoModerationRuleCreateEvent) DesiredEventType() Event {
	return &AutoModerationRuleCreateEvent{}
}
func (e *AutoModerationRuleCreateEvent) Event() EventType { return EventAutoModerationRuleCreate }
func (e *AutoModerationRuleCreateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.Rule)
}

func (e *AutoModerationRuleUpdateEvent) DesiredEventType() Event {
	return &AutoModerationRuleUpdateEvent{}
}
func (e *AutoModerationRuleUpdateEvent) Event() EventType { return EventAutoModerationRuleUpdate }
func (e *AutoModerationRuleUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewRule)
}

func (e *AutoModerationRuleDeleteEvent) DesiredEventType() Event {
	return &AutoModerationRuleDeleteEvent{}
}
func (e *AutoModerationRuleDeleteEvent) Event() EventType { return EventAutoModerationRuleDelete }
func (e *AutoModerationRuleDeleteEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.Rule)
}

func (e *AutoModerationActionExecutionEvent) DesiredEventType() Event {
	return &AutoModerationActionExecutionEvent{}
}
func (e *AutoModerationActionExecutionEvent) Event() EventType {
	return EventAutoModerationActionExecution
}
