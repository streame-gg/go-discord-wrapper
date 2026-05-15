package events

import (
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

type AutoModerationRuleCreateEvent struct {
	discord.AutoModerationRule
}

type AutoModerationRuleUpdateEvent struct {
	discord.AutoModerationRule
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
	MatchedKeyword       *string                           `json:"matched_keyword"`
	MatchedContent       *string                           `json:"matched_content"`
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
