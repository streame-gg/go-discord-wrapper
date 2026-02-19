package common

type AutoModerationTriggerType int

const (
	AutoModerationTriggerTypeKeyword       AutoModerationTriggerType = 1
	AutoModerationTriggerTypeSpam          AutoModerationTriggerType = 3
	AutoModerationTriggerTypeKeywordPreset AutoModerationTriggerType = 4
	AutoModerationTriggerTypeMentionSpam   AutoModerationTriggerType = 5
	AutoModerationTriggerTypeMemberProfile AutoModerationTriggerType = 6
)
