// events/application_command_permissions_update_test.go

package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
)

func (s *eventSuite) TestApplicationCommandPermissionsUpdate() {
	s.T().Log("Testing Application Command Permissions Update")

	sub := testutil.InitSub[ApplicationCommandPermissionsUpdateEvent](s)

	sub.RunCommonEdgeCases()

	applicationCommandPermissions := testdata.NewGuildApplicationCommandPermissions()

	sub.RunCases([]testutil.UnmarshalTestCase[ApplicationCommandPermissionsUpdateEvent]{
		{
			Name:  "valid full payload with permissions",
			Input: sub.MustMarshal(applicationCommandPermissions),
			Validate: func(e ApplicationCommandPermissionsUpdateEvent) {
				s.EqualValues(applicationCommandPermissions["id"], e.NewPermissions.ID)
				s.EqualValues(applicationCommandPermissions["application_id"], e.NewPermissions.ApplicationID)
				s.EqualValues(applicationCommandPermissions["guild_id"], e.NewPermissions.GuildID)

				permissions := applicationCommandPermissions["permissions"].([]map[string]interface{})
				s.Len(e.NewPermissions.Permissions, len(permissions))

				for i, perm := range permissions {
					s.EqualValues(perm["id"], e.NewPermissions.Permissions[i].ID)
					s.EqualValues(perm["type"], e.NewPermissions.Permissions[i].Type)
					s.EqualValues(perm["permission"], e.NewPermissions.Permissions[i].Permission)
				}
			},
		},
	})
}
