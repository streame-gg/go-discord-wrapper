package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
)

func (s *eventSuite) TestEntitlementCreate() {
	s.T().Log("Testing Entitlement Create Unmarshal Logic")

	sub := testutil.InitSub[EntitlementCreateEvent](s)

	sub.RunCommonEdgeCases()

	entitlement := testdata.NewEntitlement()

	sub.RunCases([]testutil.UnmarshalTestCase[EntitlementCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(entitlement),
			Validate: func(got EntitlementCreateEvent) {
				s.EqualValues(entitlement["id"], got.Entitlement.ID)
				s.EqualValues(entitlement["sku_id"], got.Entitlement.SkuID)
				s.EqualValues(entitlement["application_id"], got.Entitlement.ApplicationID)
				s.EqualValues(entitlement["user_id"], *got.Entitlement.UserID)
				s.EqualValues(entitlement["type"], got.Entitlement.Type)
				s.EqualValues(entitlement["deleted"], got.Entitlement.Deleted)
				s.EqualValues(entitlement["starts_at"], *got.Entitlement.StartsAt)
				s.EqualValues(entitlement["ends_at"], *got.Entitlement.EndsAt)
				s.EqualValues(entitlement["guild_id"], *got.Entitlement.GuildID)
				s.EqualValues(entitlement["consumed"], *got.Entitlement.Consumed)
			},
		},
	})
}

func (s *eventSuite) TestEntitlementUpdate() {
	s.T().Log("Testing Entitlement Update Unmarshal Logic")

	sub := testutil.InitSub[EntitlementUpdateEvent](s)

	sub.RunCommonEdgeCases()

	entitlement := testdata.NewEntitlement()

	sub.RunCases([]testutil.UnmarshalTestCase[EntitlementUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(entitlement),
			Validate: func(got EntitlementUpdateEvent) {
				s.EqualValues(entitlement["id"], got.NewEntitlement.ID)
				s.EqualValues(entitlement["sku_id"], got.NewEntitlement.SkuID)
				s.EqualValues(entitlement["application_id"], got.NewEntitlement.ApplicationID)
				s.EqualValues(entitlement["user_id"], *got.NewEntitlement.UserID)
				s.EqualValues(entitlement["type"], got.NewEntitlement.Type)
				s.EqualValues(entitlement["deleted"], got.NewEntitlement.Deleted)
				s.EqualValues(entitlement["starts_at"], *got.NewEntitlement.StartsAt)
				s.EqualValues(entitlement["ends_at"], *got.NewEntitlement.EndsAt)
				s.EqualValues(entitlement["guild_id"], *got.NewEntitlement.GuildID)
				s.EqualValues(entitlement["consumed"], *got.NewEntitlement.Consumed)
			},
		},
	})
}

func (s *eventSuite) TestEntitlementDelete() {
	s.T().Log("Testing Entitlement Delete Unmarshal Logic")

	sub := testutil.InitSub[EntitlementDeleteEvent](s)

	sub.RunCommonEdgeCases()

	entitlement := testdata.NewEntitlement()

	sub.RunCases([]testutil.UnmarshalTestCase[EntitlementDeleteEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(entitlement),
			Validate: func(got EntitlementDeleteEvent) {
				s.EqualValues(entitlement["id"], got.Entitlement.ID)
				s.EqualValues(entitlement["sku_id"], got.Entitlement.SkuID)
				s.EqualValues(entitlement["application_id"], got.Entitlement.ApplicationID)
				s.EqualValues(entitlement["user_id"], *got.Entitlement.UserID)
				s.EqualValues(entitlement["type"], got.Entitlement.Type)
				s.EqualValues(entitlement["deleted"], got.Entitlement.Deleted)
				s.EqualValues(entitlement["starts_at"], *got.Entitlement.StartsAt)
				s.EqualValues(entitlement["ends_at"], *got.Entitlement.EndsAt)
				s.EqualValues(entitlement["guild_id"], *got.Entitlement.GuildID)
				s.EqualValues(entitlement["consumed"], *got.Entitlement.Consumed)
			},
		},
	})
}
