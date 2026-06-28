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
				s.Equal(applicationCommandPermissions, e.NewPermissions)
			},
		},
	})
}
