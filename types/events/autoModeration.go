package events

import "github.com/streame-gg/go-discord-wrapper/types/common"

type AutoModerationRuleCreateEvent struct {
	common.AutoModerationRule
}

type AutoModerationRuleUpdateEvent struct {
	common.AutoModerationRule
}

type AutoModerationRuleDeleteEvent struct {
	common.AutoModerationRule
}

type AutoModerationActionExecutionEvent struct {
	GuildID              common.Snowflake                 `json:"guild_id"`
	Action               common.AutoModerationAction      `json:"action"`
	RuleID               common.Snowflake                 `json:"rule_id"`
	RuleTriggerType      common.AutoModerationTriggerType `json:"rule_trigger_type"`
	UserID               common.Snowflake                 `json:"user_id"`
	ChannelID            *common.Snowflake                `json:"channel_id,omitempty"`
	MessageID            *common.Snowflake                `json:"message_id,omitempty"`
	AlertSystemMessageID *common.Snowflake                `json:"alert_system_message_id,omitempty"`
	Content              string                           `json:"content"`
	MatchedKeyword       *string                          `json:"matched_keyword"`
	MatchedContent       *string                          `json:"matched_content"`
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
