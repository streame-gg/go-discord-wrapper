package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestThreadCreate() {
	s.T().Log("Testing Thread Create Unmarshal Logic")

	sub := testutil.InitSub[ThreadCreateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewChannel()
	payload["newly_created"] = testutil.RandomBool()

	sub.RunCases([]testutil.UnmarshalTestCase[ThreadCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got ThreadCreateEvent) {
				s.EqualValues(payload["newly_created"], *got.NewlyCreated)
				s.compareChannel(payload, got.Thread)
				s.Nil(got.Guild)
			},
		},
	})
}

func (s *eventSuite) TestThreadUpdate() {
	s.T().Log("Testing Thread Update Unmarshal Logic")

	sub := testutil.InitSub[ThreadUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewChannel()

	sub.RunCases([]testutil.UnmarshalTestCase[ThreadUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got ThreadUpdateEvent) {
				s.compareChannel(payload, got.NewThread)

				s.Nil(got.Guild)
				s.Nil(got.OldThread)
			},
		},
	})
}

func (s *eventSuite) TestThreadDelete() {
	s.T().Log("Testing Thread Delete Unmarshal Logic")

	sub := testutil.InitSub[ThreadDeleteEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"id":        discord.RandomSnowflake(),
		"guild_id":  discord.RandomSnowflake(),
		"parent_id": discord.RandomSnowflake(),
		"type": testutil.RandomItem(
			discord.ChannelTypeAnnouncementThread,
			discord.ChannelTypePublicThread,
			discord.ChannelTypePrivateThread,
		),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[ThreadDeleteEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got ThreadDeleteEvent) {
				s.EqualValues(payload["id"], got.ID)
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.EqualValues(payload["parent_id"], *got.ParentID)
				s.EqualValues(payload["type"], got.Type)

				s.Nil(got.Guild)
			},
		},
	})
}

func (s *eventSuite) TestThreadListSync() {
	s.T().Log("Testing Thread Delete Unmarshal Logic")

	sub := testutil.InitSub[ThreadListSyncEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"guild_id":    discord.RandomSnowflake(),
		"channel_ids": testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 30)),
		"threads": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 32), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, testdata.NewChannel())
		}),
		"members": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 32), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, testdata.NewThreadMember())
		}),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[ThreadListSyncEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got ThreadListSyncEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.EqualValues(payload["channel_ids"], got.ChannelIDs)

				threads := payload["threads"].([]map[string]interface{})
				s.Len(got.Threads, len(threads))
				for i, thread := range got.Threads {
					s.compareChannel(threads[i], thread)
				}

				members := payload["members"].([]map[string]interface{})
				s.Len(got.Members, len(members))
				for i, member := range got.Members {
					s.compareThreadMember(members[i], member)
				}

				s.Nil(got.Guild)
			},
		},
	})
}

func (s *eventSuite) TestThreadMemberUpdate() {
	s.T().Log("Testing Thread Member Update Unmarshal Logic")

	sub := testutil.InitSub[ThreadMemberUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewThreadMember()
	payload["guild_id"] = discord.RandomSnowflake()

	sub.RunCases([]testutil.UnmarshalTestCase[ThreadMemberUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got ThreadMemberUpdateEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.compareThreadMember(payload, got.NewMember)
				s.Nil(got.Guild)
			},
		},
	})
}

func (s *eventSuite) TestThreadMembersUpdate() {
	s.T().Log("Testing Thread Members Update Unmarshal Logic")

	sub := testutil.InitSub[ThreadMembersUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"id":           discord.RandomSnowflake(),
		"guild_id":     discord.RandomSnowflake(),
		"member_count": testutil.RandomIntInRange(1, 50),
		"added_members": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 50), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, testdata.NewThreadMember())
		}),
		"removed_member_ids": testutil.RandomSnowflakeArray(testutil.RandomIntInRange(1, 50)),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[ThreadMembersUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got ThreadMembersUpdateEvent) {
				s.EqualValues(payload["id"], got.ID)
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.EqualValues(payload["member_count"], got.MemberCount)
				s.EqualValues(payload["removed_member_ids"], got.RemovedMemberIDs)
				s.Nil(got.Guild)

				addedMembers := payload["added_members"].([]map[string]interface{})
				s.Len(got.AddedMembers, len(addedMembers))
				for i, member := range got.AddedMembers {
					s.compareThreadMember(addedMembers[i], member)
				}
			},
		},
	})
}
