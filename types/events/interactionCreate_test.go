package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func newInteractionEnvelope(itype discord.InteractionType, data interface{}) map[string]interface{} {
	payload := map[string]interface{}{
		"id":             discord.RandomSnowflake(),
		"application_id": discord.RandomSnowflake(),
		"type":           itype,
		"token":          testutil.RandomString(32),
		"version":        1,
	}
	if data != nil {
		payload["data"] = data
	}
	return payload
}

func (s *eventSuite) TestInteractionCreate() {
	s.T().Log("Testing Interaction Create Unmarshal Logic")

	sub := testutil.InitSub[InteractionCreateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewInteraction()

	cases := []testutil.UnmarshalTestCase[InteractionCreateEvent]{
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

				if got.Type != discord.InteractionTypePing {
					s.compareInteractionData(payload["data"].(map[string]interface{}), *got.Data)
				}
				s.compareMessage(payload["message"].(map[string]interface{}), *got.Message)
				s.compareGuild(payload["guild"].(map[string]interface{}), *got.Guild)
				s.compareChannel(payload["channel"].(map[string]interface{}), *got.Channel)
				s.compareMember(payload["member"].(map[string]interface{}), *got.Member)
				s.compareUser(payload["user"].(map[string]interface{}), *got.User)
			},
		},
		{
			Name:  "ping without data field: Data stays nil",
			Input: sub.MustMarshal(newInteractionEnvelope(discord.InteractionTypePing, nil)),
			Validate: func(got InteractionCreateEvent) {
				s.Equal(discord.InteractionTypePing, got.Type)
				s.Nil(got.Data)
			},
		},
		{
			Name:  "ping with data field: Data allocated but left empty",
			Input: sub.MustMarshal(newInteractionEnvelope(discord.InteractionTypePing, map[string]interface{}{})),
			Validate: func(got InteractionCreateEvent) {
				s.Require().NotNil(got.Data)
				s.Nil(*got.Data)
			},
		},
		func() testutil.UnmarshalTestCase[InteractionCreateEvent] {
			data := testdata.NewInteractionDataApplicationCommand()
			return testutil.UnmarshalTestCase[InteractionCreateEvent]{
				Name:  "autocomplete: Data decodes to InteractionDataAutocomplete",
				Input: sub.MustMarshal(newInteractionEnvelope(discord.InteractionTypeApplicationCommandAutocomplete, data)),
				Validate: func(got InteractionCreateEvent) {
					s.Require().NotNil(got.Data)
					s.compareInteractionData(data, *got.Data)
				},
			}
		}(),
		func() testutil.UnmarshalTestCase[InteractionCreateEvent] {
			data := testdata.NewModalSubmitData()
			return testutil.UnmarshalTestCase[InteractionCreateEvent]{
				Name:  "modal submit: Data decodes to InteractionDataModalSubmit",
				Input: sub.MustMarshal(newInteractionEnvelope(discord.InteractionTypeModalSubmit, data)),
				Validate: func(got InteractionCreateEvent) {
					s.Require().NotNil(got.Data)
					s.compareInteractionData(data, *got.Data)
				},
			}
		}(),
		func() testutil.UnmarshalTestCase[InteractionCreateEvent] {
			data := testdata.NewInteractionDataMessageComponent()
			return testutil.UnmarshalTestCase[InteractionCreateEvent]{
				Name:  "message component: Data decodes to InteractionDataMessageComponent",
				Input: sub.MustMarshal(newInteractionEnvelope(discord.InteractionTypeMessageComponent, data)),
				Validate: func(got InteractionCreateEvent) {
					s.Require().NotNil(got.Data)
					s.compareInteractionData(data, *got.Data)
				},
			}
		}(),
		func() testutil.UnmarshalTestCase[InteractionCreateEvent] {
			data := testdata.NewInteractionDataApplicationCommand()
			data["type"] = discord.ApplicationCommandTypeChatInput
			return testutil.UnmarshalTestCase[InteractionCreateEvent]{
				Name:  "application command: Data decodes to InteractionDataApplicationCommand",
				Input: sub.MustMarshal(newInteractionEnvelope(discord.InteractionTypeApplicationCommand, data)),
				Validate: func(got InteractionCreateEvent) {
					s.Require().NotNil(got.Data)
					s.compareInteractionData(data, *got.Data)
				},
			}
		}(),
		{
			Name: "malformed data field: probe unmarshal error",
			Input: sub.MustMarshal(newInteractionEnvelope(
				discord.InteractionTypeApplicationCommand,
				[]interface{}{},
			)),
			WantErr: true,
		},
	}

	modalSubmitChildCases := []struct {
		name  string
		child map[string]interface{}
	}{
		{"string select", testdata.NewStringSelectMenuData()},
		{"user select", testdata.NewUserSelectMenuData()},
		{"role select", testdata.NewRoleSelectMenuData()},
		{"mentionable select", testdata.NewMentionableSelectMenuData()},
		{"channel select", testdata.NewChannelSelectMenuData()},
		{"text input", testdata.NewTextInputData()},
		{"file upload", testdata.NewFileUploadData()},
		{"radio group", testdata.NewRadioGroupData()},
		{"checkbox group", testdata.NewCheckboxGroupData()},
		{"checkbox", testdata.NewCheckboxData()},
		{"text display", testdata.NewTextDisplayData()},
		{"nested label", testdata.NewLabelData()},
	}

	for _, tc := range modalSubmitChildCases {
		data := testdata.NewModalSubmitDataWithComponents(testdata.NewLabelDataWithComponent(tc.child))
		cases = append(cases, testutil.UnmarshalTestCase[InteractionCreateEvent]{
			Name:  "modal submit component dispatch: " + tc.name,
			Input: sub.MustMarshal(newInteractionEnvelope(discord.InteractionTypeModalSubmit, data)),
			Validate: func(got InteractionCreateEvent) {
				s.Require().NotNil(got.Data)
				s.compareInteractionData(data, *got.Data)
			},
		})
	}

	cases = append(cases,
		testutil.UnmarshalTestCase[InteractionCreateEvent]{
			Name: "modal submit component dispatch: unrecognized child type errors",
			Input: sub.MustMarshal(newInteractionEnvelope(
				discord.InteractionTypeModalSubmit,
				testdata.NewModalSubmitDataWithComponents(testdata.NewLabelDataWithComponent(
					map[string]interface{}{"type": 9999, "id": 1},
				)),
			)),
			WantErr: true,
		},
	)

	sub.RunCases(cases)
}
