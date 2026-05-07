package connection

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// TestClientImplementsCloser is a compile-time check that *Client satisfies io.Closer.
var _ io.Closer = (*Client)(nil)

func TestNewClientValidConfig(t *testing.T) {
	c, err := NewClient("test-token", common.IntentGuilds)
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestNewClientErrorOnBadSharding(t *testing.T) {
	// ShardID (5) >= TotalShards (4) — invalid configuration, must return error.
	_, err := NewClient("test-token", common.IntentGuilds, options.WithSharding(4, 5))
	assert.Error(t, err)
}
