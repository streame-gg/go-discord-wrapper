package discord

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type soundboardTestSuite struct {
	suite.Suite
}

func TestSoundboardTestSuite(t *testing.T) {
	suite.Run(t, new(soundboardTestSuite))
}

func (s *soundboardTestSuite) TestUnmarshalSoundboardSound_1() {
	rawSound := "{\n  \"name\": \"quack\",\n  \"sound_id\": \"1\",\n  \"volume\": 1.0,\n  \"emoji_id\": null,\n  \"emoji_name\": \"🦆\",\n  \"available\": true\n}"
	var sound SoundboardSound
	s.Require().NoError(json.Unmarshal([]byte(rawSound), &sound))
	s.Equal("quack", sound.Name)
	s.Equal(Snowflake(1), sound.SoundID)
	s.Equal(1.0, sound.Volume)
	s.Nil(sound.EmojiID)
	s.Equal("🦆", *sound.EmojiName)
	s.True(sound.Available)
	s.Nil(sound.User)
	s.Nil(sound.GuildID)
}

func (s *soundboardTestSuite) TestUnmarshalSoundboardSound_2() {
	rawSound := "{\n  \"name\": \"Yay\",\n  \"sound_id\": \"1106714396018884649\",\n  \"volume\": 1,\n  \"emoji_id\": \"989193655938064464\",\n  \"emoji_name\": null,\n  \"guild_id\": \"613425648685547541\",\n  \"available\": true\n}"
	var sound SoundboardSound
	s.Require().NoError(json.Unmarshal([]byte(rawSound), &sound))
	s.Equal("Yay", sound.Name)
	s.Equal(Snowflake(1106714396018884649), sound.SoundID)
	s.Equal(float64(1), sound.Volume)
	s.Equal(Snowflake(989193655938064464), *sound.EmojiID)
	s.Nil(sound.EmojiName)
	s.True(sound.Available)
	s.Nil(sound.User)
	s.Equal(Snowflake(613425648685547541), *sound.GuildID)
}
