package discord

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// J0-#26: Emoji.Users must decode as a slice, not a single pointer.
func TestEmojiUsersDecodeSlice(t *testing.T) {
	raw := `{"users":[{"id":"123","username":"a"},{"id":"456","username":"b"}]}`
	var e Emoji
	require.NoError(t, json.Unmarshal([]byte(raw), &e))
	assert.Len(t, e.Users, 2)
	assert.Equal(t, "a", e.Users[0].Username)
	assert.Equal(t, "b", e.Users[1].Username)
}

func TestEmojiUsersRoundtrip(t *testing.T) {
	e := Emoji{Users: []User{{Username: "foo"}, {Username: "bar"}}}
	b, err := json.Marshal(e)
	require.NoError(t, err)
	var e2 Emoji
	require.NoError(t, json.Unmarshal(b, &e2))
	assert.Len(t, e2.Users, 2)
}

// J0-#27: User.AvatarHash must be *string so null and absent are distinguishable.
func TestUserAvatarHashNull(t *testing.T) {
	raw := `{"id":"1","username":"u","discriminator":"0","avatar":null}`
	var u User
	require.NoError(t, json.Unmarshal([]byte(raw), &u))
	assert.Nil(t, u.AvatarHash, "null avatar must decode to nil pointer")
}

func TestUserAvatarHashPresent(t *testing.T) {
	raw := `{"id":"1","username":"u","discriminator":"0","avatar":"abc123"}`
	var u User
	require.NoError(t, json.Unmarshal([]byte(raw), &u))
	require.NotNil(t, u.AvatarHash)
	assert.Equal(t, "abc123", *u.AvatarHash)
}
