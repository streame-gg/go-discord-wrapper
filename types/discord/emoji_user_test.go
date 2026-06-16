package discord

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// J0-#26: Emoji.Users must decode as a slice, not a single pointer.
func (su *emojiUserSuite) TestEmojiUsersDecodeSlice() {
	t := su.T()
	raw := `{"user":{"id":"123","username":"a"}}`
	var e Emoji
	require.NoError(t, json.Unmarshal([]byte(raw), &e))
	assert.Equal(t, "a", e.User.Username)
}

func (su *emojiUserSuite) TestEmojiUsersRoundtrip() {
	t := su.T()
	e := Emoji{User: &User{Username: "a"}}
	b, err := json.Marshal(e)
	require.NoError(t, err)
	var e2 Emoji
	require.NoError(t, json.Unmarshal(b, &e2))
	su.Equal(e, e2)
}

// J0-#27: User.AvatarHash must be *string so null and absent are distinguishable.
func (su *emojiUserSuite) TestUserAvatarHashNull() {
	t := su.T()
	raw := `{"id":"1","username":"u","discriminator":"0","avatar":null}`
	var u User
	require.NoError(t, json.Unmarshal([]byte(raw), &u))
	assert.Nil(t, u.AvatarHash, "null avatar must decode to nil pointer")
}

func (su *emojiUserSuite) TestUserAvatarHashPresent() {
	t := su.T()
	raw := `{"id":"1","username":"u","discriminator":"0","avatar":"abc123"}`
	var u User
	require.NoError(t, json.Unmarshal([]byte(raw), &u))
	require.NotNil(t, u.AvatarHash)
	assert.Equal(t, "abc123", *u.AvatarHash)
}

type emojiUserSuite struct{ suite.Suite }

func TestEmojiUserSuite(t *testing.T) { suite.Run(t, new(emojiUserSuite)) }
