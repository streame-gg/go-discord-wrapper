package sharding_test

import (
	"testing"

	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/sharding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Bug 10: Register after Close must return an error, not panic with a nil-map write.
func TestBug10_RegisterAfterCloseMustNotPanic(t *testing.T) {
	c := sharding.NewLocalCoordinator(2)
	require.NoError(t, c.Close())

	var err error
	assert.NotPanics(t, func() {
		err = c.Register(0, func(options.ShardMessage) {})
	})
	assert.Error(t, err, "Register after Close must return an error (Bug 10)")
}

// Bug 10: Send after Close must return an error, not panic.
func TestBug10_SendAfterCloseMustNotPanic(t *testing.T) {
	c := sharding.NewLocalCoordinator(2)
	require.NoError(t, c.Register(0, func(options.ShardMessage) {}))
	require.NoError(t, c.Close())

	var err error
	assert.NotPanics(t, func() {
		err = c.Send(options.ShardMessage{To: 0, Type: "test"})
	})
	assert.Error(t, err, "Send after Close must return an error (Bug 10)")
}

// Bug 10: Broadcast after Close must return an error, not panic.
func TestBug10_BroadcastAfterCloseMustNotPanic(t *testing.T) {
	c := sharding.NewLocalCoordinator(2)
	require.NoError(t, c.Close())

	var err error
	assert.NotPanics(t, func() {
		err = c.Broadcast(options.ShardMessage{Type: "test"})
	})
	assert.Error(t, err, "Broadcast after Close must return an error (Bug 10)")
}
