package connection

import (
	"sync/atomic"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// TestBug66_LifecycleEmittersTrackedByDispatchWg verifies that goroutines
// spawned by emitConnect/Disconnect/Reconnect/PacketError are counted in
// dispatchWg, so Shutdown waits for them to finish.
func (cs *ConnectionSuite) TestBug66_LifecycleEmittersTrackedByDispatchWg() {
	t := cs.T()
	c, err := NewClient("fake-token", discord.IntentGuilds)
	require.NoError(t, err)

	blocker := make(chan struct{})
	var count atomic.Int32

	c.OnConnect(func(_ *Client) { <-blocker; count.Add(1) })
	c.OnDisconnect(func(_ *Client, _ error) { <-blocker; count.Add(1) })
	c.OnReconnect(func(_ *Client) { <-blocker; count.Add(1) })
	c.OnPacketError(func(_ *Client, _ error) { <-blocker; count.Add(1) })

	c.emitConnect()
	c.emitDisconnect(nil)
	c.emitReconnect()
	c.emitPacketError(nil)

	// Handlers are blocked; dispatchWg must have 4 outstanding entries.
	// Release them and confirm Wait() resolves only after all four complete.
	close(blocker)

	done := make(chan struct{})
	go func() { c.dispatchWg.Wait(); close(done) }()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("dispatchWg.Wait() timed out — lifecycle emitters not tracked by dispatchWg (Bug 66)")
	}

	assert.Equal(t, int32(4), count.Load(), "all four lifecycle handlers must have run")
}

// TestBug66_DispatchWgWaitsBeforeHandlersFinish is the inverse: confirm that
// dispatchWg.Wait() does NOT return while handlers are still blocked, proving
// the Add/Done pairing is not trivially vacuous.
func (cs *ConnectionSuite) TestBug66_DispatchWgBlocksUntilHandlersDone() {
	t := cs.T()
	c, err := NewClient("fake-token", discord.IntentGuilds)
	require.NoError(t, err)

	blocker := make(chan struct{})

	c.OnConnect(func(_ *Client) { <-blocker })

	c.emitConnect()

	// Wait() must not return while the handler is still blocked.
	waitDone := make(chan struct{})
	go func() { c.dispatchWg.Wait(); close(waitDone) }()

	select {
	case <-waitDone:
		t.Fatal("dispatchWg.Wait() returned before handler finished")
	case <-time.After(50 * time.Millisecond):
		// expected: still waiting
	}

	close(blocker)

	select {
	case <-waitDone:
	case <-time.After(5 * time.Second):
		t.Fatal("dispatchWg.Wait() never returned after handler finished")
	}
}

// TestBug67_LifecycleEmitterPanicDoesNotCrash verifies that a panic inside any
// lifecycle handler is caught by the per-goroutine recover, leaving the process
// alive and dispatchWg balanced.
func (cs *ConnectionSuite) TestBug67_LifecycleEmitterPanicDoesNotCrash() {
	t := cs.T()
	c, err := NewClient("fake-token", discord.IntentGuilds)
	require.NoError(t, err)

	c.OnConnect(func(_ *Client) { panic("connect panic") })
	c.OnDisconnect(func(_ *Client, _ error) { panic("disconnect panic") })
	c.OnReconnect(func(_ *Client) { panic("reconnect panic") })
	c.OnPacketError(func(_ *Client, _ error) { panic("packet-error panic") })

	// None of the calls below should crash the test process.
	c.emitConnect()
	c.emitDisconnect(nil)
	c.emitReconnect()
	c.emitPacketError(nil)

	done := make(chan struct{})
	go func() { c.dispatchWg.Wait(); close(done) }()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("dispatchWg.Wait() timed out after panicking handlers (Bug 67)")
	}
}
