package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
)

func (s *eventSuite) TestInteractionCreate() {
	s.T().Log("Testing Interaction Create Unmarshal Logic")

	sub := testutil.InitSub[InteractionCreateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewInteraction()

	sub.RunCases([]testutil.UnmarshalTestCase[InteractionCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got InteractionCreateEvent) {
				s.EqualValues(payload["id"], got.ID)
				s.EqualValues(payload["application_id"], got.ApplicationID)
				s.EqualValues(payload["type"], got.Type)
				s.EqualValues(payload["guild_id"], *got.GuildID)
				s.EqualValues(payload["channel_id"], *got.ChannelID)
				s.EqualValues(payload["token"], got.Token)
				s.EqualValues(payload["version"], got.Version)
				s.EqualValues(payload["app_permissions"], got.AppPermissions)
				s.EqualValues(payload["locale"], got.Locale)
				s.EqualValues(payload["guild_locale"], got.GuildLocale)
				s.EqualValues(payload["authorizing_integration_owners"], got.AuthorizingIntegrationOwners)
				s.EqualValues(payload["context"], got.Context)
				s.EqualValues(payload["attachment_size_limit"], got.AttachmentSizeLimit)

				entilements := payload["entitlements"].([]map[string]interface{})
				s.Len(got.Entitlements, len(entilements))
				for i, ent := range entilements {
					s.compareEntitlement(ent, got.Entitlements[i])
				}

				s.compareInteractionData(payload["data"].(map[string]interface{}), got.Data)
				s.compareMessage(payload["message"].(map[string]interface{}), *got.Message)
				s.compareGuild(payload["guild"].(map[string]interface{}), *got.Guild)
				s.compareChannel(payload["channel"].(map[string]interface{}), *got.Channel)
				s.compareMember(payload["member"].(map[string]interface{}), *got.Member)
				s.compareUser(payload["user"].(map[string]interface{}), *got.User)
			},
		},
	})
}
