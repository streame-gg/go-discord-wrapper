package connection

import (
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSubscribeShardResponse_NilCoordinator verifies that SubscribeShardResponse
// does not panic when no Coordinator is configured (Bug 101).
func (cs *ConnectionSuite) TestSubscribeShardResponse_NilCoordinator() {
	t := cs.T()
	c, err := NewClient("test-token", discord.IntentGuilds)
	require.NoError(t, err)
	require.Nil(t, c.Coordinator, "precondition: no coordinator")

	assert.NotPanics(t, func() {
		ch, cleanup := c.SubscribeShardResponse("TYPE", "corr-1")
		defer cleanup()
		assert.NotNil(t, ch)
	})
}

// TestSubscribeShardResponse_DuplicateCorrID verifies that a second call with the
// same responseType+corrID does not overwrite the first subscription (Bug 100).
func (cs *ConnectionSuite) TestSubscribeShardResponse_DuplicateCorrID() {
	t := cs.T()
	c, err := NewClient("test-token", discord.IntentGuilds)
	require.NoError(t, err)

	ch1, cleanup1 := c.SubscribeShardResponse("TYPE", "corr-dup")
	defer cleanup1()

	ch2, cleanup2 := c.SubscribeShardResponse("TYPE", "corr-dup")
	defer cleanup2()

	// The first channel must still be registered (no overwrite).
	msg := options.ShardMessage{Type: "TYPE", CorrelationID: "corr-dup", From: 1}
	c.pendingShardReqsMu.Lock()
	registered := c.pendingShardReqs["TYPE:corr-dup"]
	c.pendingShardReqsMu.Unlock()

	assert.Equal(t, ch1, (<-chan options.ShardMessage)(registered),
		"first channel must still be registered after duplicate subscribe")
	assert.NotEqual(t, ch1, ch2,
		"second call must return a distinct (non-registered) channel")

	// Send a message directly into the registered channel; only ch1 receives it.
	registered <- msg
	select {
	case got := <-ch1:
		assert.Equal(t, msg, got)
	default:
		t.Fatal("ch1 should have received the message")
	}
	select {
	case <-ch2:
		t.Fatal("ch2 must not receive any message — it was not registered")
	default:
		// correct
	}
}
