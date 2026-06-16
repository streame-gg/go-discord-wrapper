package discord

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type allowedMentionsTestSuite struct {
	suite.Suite
}

func TestNewAllowedMentionsTestSuite(t *testing.T) {
	suite.Run(t, new(allowedMentionsTestSuite))
}

func (s *allowedMentionsTestSuite) TestAllowedMentions_MarshalNilValues() {
	mentions := AllowedMentions{
		Parse:       nil,
		Roles:       nil,
		Users:       nil,
		RepliedUser: false,
	}

	marshal, err := json.Marshal(&mentions)
	s.NoError(err)
	s.Require().Equal("{}", string(marshal))
}

func (s *allowedMentionsTestSuite) TestAllowedMentions_MarshalParseIsNilSlice() {
	mentions := AllowedMentions{
		Parse: &[]AllowedMentionsType{},
		Roles: make([]Snowflake, 0),
	}

	marshal, err := json.Marshal(&mentions)
	s.NoError(err)
	s.Require().Equal("{\"parse\":[]}", string(marshal))
}

func (s *allowedMentionsTestSuite) TestAllowedMentions_MarshalWithValues() {
	// In real API use cases this wouldn't work, but we are just trying the marshal logic
	mentions := AllowedMentions{
		Parse: &[]AllowedMentionsType{
			AllowedMentionsTypeEveryone, AllowedMentionsTypeRoles, AllowedMentionsTypeUsers,
		},
		Users:       []Snowflake{123456789012345678, 987654321098765432},
		Roles:       []Snowflake{111111111111111111, 222222222222222222},
		RepliedUser: true,
	}

	marshal, err := json.Marshal(&mentions)
	s.NoError(err)
	s.Require().Equal("{\"parse\":[\"everyone\",\"roles\",\"users\"],\"roles\":[\"111111111111111111\",\"222222222222222222\"],\"users\":[\"123456789012345678\",\"987654321098765432\"],\"replied_user\":true}", string(marshal))
}

func (s *allowedMentionsTestSuite) TestAllowedMentions_UnmarshalWithNilValues() {
	object := "{}"
	var mentions AllowedMentions
	s.Require().NoError(json.Unmarshal([]byte(object), &mentions))
	s.Nil(mentions.Users)
	s.Nil(mentions.Roles)
	s.Nil(mentions.Parse)
	s.False(mentions.RepliedUser)
}

func (s *allowedMentionsTestSuite) TestAllowedMentions_UnmarshalWithNilSlice() {
	object := "{\"parse\":[], \"roles\":[], \"users\":[]}"
	var mentions AllowedMentions
	s.Require().NoError(json.Unmarshal([]byte(object), &mentions))
	s.Equal(0, len(*mentions.Parse))
	s.Equal(0, len(mentions.Users))
	s.Equal(0, len(mentions.Roles))
	s.False(mentions.RepliedUser)
}

func (s *allowedMentionsTestSuite) TestAllowedMentions_UnmarshalWithValues() {
	object := "{\"parse\":[\"everyone\",\"roles\",\"users\"],\"roles\":[\"111111111111111111\",\"222222222222222222\"],\"users\":[\"123456789012345678\",\"987654321098765432\"],\"replied_user\":true}"
	var mentions AllowedMentions
	s.Require().NoError(json.Unmarshal([]byte(object), &mentions))
	s.Equal([]AllowedMentionsType{
		AllowedMentionsTypeEveryone, AllowedMentionsTypeRoles, AllowedMentionsTypeUsers,
	}, *mentions.Parse)
	s.Equal([]Snowflake{123456789012345678, 987654321098765432}, mentions.Users)
	s.Equal([]Snowflake{111111111111111111, 222222222222222222}, mentions.Roles)
	s.True(mentions.RepliedUser)
}

func (s *allowedMentionsTestSuite) TestAllowedMentions_RepliedUserIsFalse() {
	object := "{\"replied_user\":false}"
	var mentions AllowedMentions
	s.Require().NoError(json.Unmarshal([]byte(object), &mentions))
	s.Equal(false, mentions.RepliedUser)
}

type test struct {
	Test bool `json:"test"`
}
