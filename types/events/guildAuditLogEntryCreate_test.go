package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestGuildAuditLogEntryCreate() {
	s.T().Log("Testing Guild Audit Log Entry Create Unmarshal Logic")

	sub := testutil.InitSub[GuildAuditLogEntryCreateEvent](s)

	sub.RunCommonEdgeCases()

	auditLogEntry := testdata.NewAuditLogEntry()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildAuditLogEntryCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(auditLogEntry),
			Validate: func(got GuildAuditLogEntryCreateEvent) {
				s.EqualValues(auditLogEntry["target_id"], *got.AuditLogEntry.TargetID)
				s.EqualValues(auditLogEntry["id"], got.AuditLogEntry.ID)
				s.EqualValues(auditLogEntry["user_id"], *got.AuditLogEntry.UserID)
				s.EqualValues(auditLogEntry["reason"], got.AuditLogEntry.Reason)
				s.EqualValues(auditLogEntry["action_type"], got.AuditLogEntry.ActionType)

				changes := auditLogEntry["changes"].([]discord.AuditLogEntryChange)
				for i, change := range changes {
					s.EqualValues(change.Key, got.AuditLogEntry.Changes[i].Key)
					s.EqualValues(change.NewValue, got.AuditLogEntry.Changes[i].NewValue)
					s.EqualValues(change.OldValue, got.AuditLogEntry.Changes[i].OldValue)
				}

				options := auditLogEntry["options"].(map[string]interface{})
				s.EqualValues(options["application_id"], *got.AuditLogEntry.Options.ApplicationID)
				s.EqualValues(options["auto_moderation_rule_name"], got.AuditLogEntry.Options.AutoModerationRuleName)
				s.EqualValues(options["auto_moderation_rule_trigger_type"], got.AuditLogEntry.Options.AutoModerationRuleTriggerType)
				s.EqualValues(options["channel_id"], *got.AuditLogEntry.Options.ChannelID)
				s.EqualValues(options["count"], got.AuditLogEntry.Options.Count)
				s.EqualValues(options["delete_member_days"], got.AuditLogEntry.Options.DeleteMemberDays)
				s.EqualValues(options["id"], *got.AuditLogEntry.Options.ID)
				s.EqualValues(options["members_removed"], got.AuditLogEntry.Options.MembersRemoved)
				s.EqualValues(options["message_id"], *got.AuditLogEntry.Options.MessageID)
				s.EqualValues(options["role_name"], got.AuditLogEntry.Options.RoleName)
				s.EqualValues(options["type"], *got.AuditLogEntry.Options.Type)
				s.EqualValues(options["integration_type"], *got.AuditLogEntry.Options.IntegrationType)
				s.EqualValues(options["status"], got.AuditLogEntry.Options.Status)
			},
		},
	})
}
