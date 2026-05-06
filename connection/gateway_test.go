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
	assert.NotPanics(t, func() {
		_ = NewClient("test-token", common.IntentGuilds)
	})
}

func TestNewClientPanicsOnBadSharding(t *testing.T) {
	assert.Panics(t, func() {
		// ShardID (5) >= TotalShards (4) — invalid configuration, must panic.
		_ = NewClient("test-token", common.IntentGuilds, options.WithSharding(4, 5))
	})
}
