package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
)

func (s *eventSuite) TestGuildCreate() {
	s.T().Log("Testing Guild Create Unmarshal Logic")

	sub := testutil.InitSub[GuildCreateEvent](s)

	sub.RunCommonEdgeCases()

	guild := testdata.NewAvailableGuildWithGuildCreateValues()
	unavailableGuild := testdata.NewUnavailableGuild()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(guild),
			Validate: func(got GuildCreateEvent) {
				s.Require().NotNil(got.Guild)
				s.False(got.Unavailable)
				s.compareGuild(guild, *got.Guild)

				// compare values only for guild create
				s.EqualValues(guild["joined_at"], got.JoinedAt)
				s.EqualValues(guild["large"], got.Large)
				s.EqualValues(guild["member_count"], got.MemberCount)

				stageInstances := guild["stage_instances"].([]map[string]interface{})
				s.Require().Len(got.StageInstances, len(stageInstances))
				for i, instance := range stageInstances {
					s.compareStageInstance(instance, got.StageInstances[i])
				}

				sounds := guild["soundboard_sounds"].([]map[string]interface{})
				s.Require().Len(got.SoundboardSounds, len(sounds))
				for i, sound := range sounds {
					s.compareSoundboardSound(sound, got.SoundboardSounds[i])
				}

				events := guild["guild_scheduled_events"].([]map[string]interface{})
				s.Require().Len(got.GuildScheduledEvents, len(events))
				for i, event := range events {
					s.compareGuildScheduledEvent(event, got.GuildScheduledEvents[i])
				}

				voiceStates := guild["voice_states"].([]map[string]interface{})
				s.Require().Len(got.VoiceStates, len(voiceStates))
				for i, voiceState := range voiceStates {
					s.compareVoiceState(voiceState, got.VoiceStates[i])
				}

				members := guild["members"].([]map[string]interface{})
				s.Require().Len(got.Members, len(members))
				for i, member := range members {
					s.compareMember(member, got.Members[i])
				}

				channels := guild["channels"].([]map[string]interface{})
				s.Require().Len(got.Channels, len(channels))
				for i, channel := range channels {
					s.compareChannel(channel, got.Channels[i])
				}

				threads := guild["threads"].([]map[string]interface{})
				s.Require().Len(got.Threads, len(threads))
				for i, thread := range threads {
					s.compareChannel(thread, got.Threads[i])
				}

				presences := guild["presences"].([]map[string]interface{})
				s.Require().Len(got.Presences, len(presences))
				for i, presence := range presences {
					s.comparePresence(presence, got.Presences[i])
				}
			},
		},
		{
			Name:  "unavailable guild",
			Input: sub.MustMarshal(unavailableGuild),
			Validate: func(got GuildCreateEvent) {
				s.True(got.Unavailable)
				s.EqualValues(unavailableGuild["id"], got.ID)
				s.Nil(got.Guild)
			},
		},
	})
}

func (s *eventSuite) TestGuildUpdate() {
	s.T().Log("Testing Guild Update Unmarshal Logic")

	sub := testutil.InitSub[GuildUpdateEvent](s)

	sub.RunCommonEdgeCases()

	guild := testdata.NewAvailableGuild()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(guild),
			Validate: func(got GuildUpdateEvent) {
				s.compareGuild(guild, got.NewGuild)
			},
		},
	})
}

func (s *eventSuite) TestGuildDelete() {
	s.T().Log("Testing Guild Delete Unmarshal Logic")

	sub := testutil.InitSub[GuildDeleteEvent](s)

	sub.RunCommonEdgeCases()

	unavailableGuild := testdata.NewUnavailableGuild()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildDeleteEvent]{
		{
			Name:  "unavailable guild",
			Input: sub.MustMarshal(unavailableGuild),
			Validate: func(got GuildDeleteEvent) {
				s.True(got.Unavailable)
				s.EqualValues(unavailableGuild["id"], got.ID)
				s.Nil(got.Guild)
			},
		},
	})
}
