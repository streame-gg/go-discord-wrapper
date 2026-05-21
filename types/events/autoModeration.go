package events

import (
	"encoding/json"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type AutoModerationRuleCreateEvent struct {
	discord.AutoModerationRule
}

type AutoModerationRuleUpdateEvent struct {
	NewRule discord.AutoModerationRule `json:"-"`
	// OldRule is nil when the rule was not previously cached (no automod rule store exists).
	OldRule *discord.AutoModerationRule `json:"-"`
}

func (e *AutoModerationRuleUpdateEvent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.NewRule)
}

func (e AutoModerationRuleUpdateEvent) MarshalJSON() ([]byte, error) {
	type wire struct {
		discord.AutoModerationRule
		OldRule *discord.AutoModerationRule `json:"old_rule,omitempty"`
	}
	return json.Marshal(wire{e.NewRule, e.OldRule})
}

type AutoModerationRuleDeleteEvent struct {
	discord.AutoModerationRule
}

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
	MatchedKeyword       *string                           `json:"matched_keyword,omitempty"`
	MatchedContent       *string                           `json:"matched_content,omitempty"`
}

func init() {
	RegisterEvent(AutoModerationRuleCreateEvent{})
	RegisterEvent(AutoModerationRuleUpdateEvent{})
	RegisterEvent(AutoModerationRuleDeleteEvent{})
	RegisterEvent(AutoModerationActionExecutionEvent{})
}

func (e AutoModerationRuleCreateEvent) DesiredEventType() Event {
	return &AutoModerationRuleCreateEvent{}
}
func (e AutoModerationRuleCreateEvent) Event() EventType { return EventAutoModerationRuleCreate }

func (e AutoModerationRuleUpdateEvent) DesiredEventType() Event {
	return &AutoModerationRuleUpdateEvent{}
}
func (e AutoModerationRuleUpdateEvent) Event() EventType { return EventAutoModerationRuleUpdate }

func (e AutoModerationRuleDeleteEvent) DesiredEventType() Event {
	return &AutoModerationRuleDeleteEvent{}
}
func (e AutoModerationRuleDeleteEvent) Event() EventType { return EventAutoModerationRuleDelete }

func (e AutoModerationActionExecutionEvent) DesiredEventType() Event {
	return &AutoModerationActionExecutionEvent{}
}
func (e AutoModerationActionExecutionEvent) Event() EventType {
	return EventAutoModerationActionExecution
}
