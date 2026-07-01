package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (s *eventSuite) TestGuildSoundboardSoundCreate() {
	s.T().Log("Testing Guild Soundboard Sound Create Unmarshal Logic")

	sub := testutil.InitSub[GuildSoundboardSoundCreateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewSoundboardSound()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildSoundboardSoundCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got GuildSoundboardSoundCreateEvent) {
				s.compareSoundboardSound(payload, got.SoundboardSound)
			},
		},
	})
}

func (s *eventSuite) TestGuildSoundboardSoundUpdate() {
	s.T().Log("Testing Guild Soundboard Sound Update Unmarshal Logic")

	sub := testutil.InitSub[GuildSoundboardSoundUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := testdata.NewSoundboardSound()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildSoundboardSoundUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got GuildSoundboardSoundUpdateEvent) {
				s.compareSoundboardSound(payload, got.NewSound)
				s.Nil(got.OldSound)
			},
		},
	})
}

func (s *eventSuite) TestGuildSoundboardSoundDelete() {
	s.T().Log("Testing Guild Soundboard Sound Delete Unmarshal Logic")

	sub := testutil.InitSub[GuildSoundboardSoundDeleteEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"guild_id": discord.RandomSnowflake(),
		"sound_id": discord.RandomSnowflake(),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[GuildSoundboardSoundDeleteEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got GuildSoundboardSoundDeleteEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)
				s.EqualValues(payload["sound_id"], got.SoundID)
			},
		},
	})
}

func (s *eventSuite) TestGuildSoundboardSoundsUpdateEvent() {
	s.T().Log("Testing Guild Soundboard Sounds Update Unmarshal Logic")

	sub := testutil.InitSub[GuildSoundboardSoundsUpdateEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"guild_id": discord.RandomSnowflake(),
		"soundboard_sounds": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, testdata.NewSoundboardSound())
		}),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[GuildSoundboardSoundsUpdateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got GuildSoundboardSoundsUpdateEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)

				sounds := payload["soundboard_sounds"].([]map[string]interface{})
				s.Len(got.NewSoundboardSounds, len(sounds))
				for i, sound := range got.NewSoundboardSounds {
					s.compareSoundboardSound(sounds[i], sound)
				}
			},
		},
	})
}

func (s *eventSuite) TestSoundboardSoundsEvent() {
	s.T().Log("Testing Soundboard Sounds Update Unmarshal Logic")

	sub := testutil.InitSub[SoundboardSoundsEvent](s)

	sub.RunCommonEdgeCases()

	payload := map[string]interface{}{
		"guild_id": discord.RandomSnowflake(),
		"soundboard_sounds": testutil.RandomArrayWithFilledItems(testutil.RandomIntInRange(1, 10), func(arrayToFill *[]map[string]interface{}) {
			*arrayToFill = append(*arrayToFill, testdata.NewSoundboardSound())
		}),
	}

	sub.RunCases([]testutil.UnmarshalTestCase[SoundboardSoundsEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(payload),
			Validate: func(got SoundboardSoundsEvent) {
				s.EqualValues(payload["guild_id"], got.GuildID)

				sounds := payload["soundboard_sounds"].([]map[string]interface{})
				s.Len(got.SoundboardSounds, len(sounds))
				for i, sound := range got.SoundboardSounds {
					s.compareSoundboardSound(sounds[i], sound)
				}
			},
		},
	})
}
