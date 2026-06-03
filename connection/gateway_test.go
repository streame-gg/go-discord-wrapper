package connection

import (
	"io"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// TestClientImplementsCloser is a compile-time check that *Client satisfies io.Closer.
var _ io.Closer = (*Client)(nil)

func (cs *ConnectionSuite) TestNewClientValidConfig() {
	t := cs.T()
	c, err := NewClient("test-token", discord.IntentGuilds)
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func (cs *ConnectionSuite) TestNewClientErrorOnBadSharding() {
	t := cs.T()
	// ShardID (5) >= TotalShards (4) — invalid configuration, must return error.
	_, err := NewClient("test-token", discord.IntentGuilds, options.WithSharding(4, 5))
	assert.Error(t, err)
}

// TestBug64_BotPrefixStripped verifies that NewClient strips a leading "Bot "
// prefix (case-insensitive) so the gateway never sends "Bot xyz" as the token.
func (cs *ConnectionSuite) TestBug64_BotPrefixStripped() {
	t := cs.T()
	cases := []struct {
		input string
		want  string
	}{
		{"Bot mytoken", "mytoken"},
		{"bot mytoken", "mytoken"},
		{"BOT mytoken", "mytoken"},
		{"mytoken", "mytoken"},
	}
	for _, tc := range cases {
		c, err := NewClient(tc.input, discord.IntentGuilds)
		require.NoError(t, err)
		assert.Equal(t, tc.want, *c.token, "input %q", tc.input)
	}
}
