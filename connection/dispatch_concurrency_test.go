package connection

import (
	"encoding/json"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// Bug 9: With MaxConcurrentEvents=2 (2 workers) and a handler that sleeps
// 200ms, dispatching 10 events serially through the pool must take at least
// 1 second wall-clock time (10 events / 2 workers × 200 ms = 1 s).
//
// Before the fix, handlers ran in goroutines inside each worker so all 10
// could complete in ≈200 ms regardless of the worker count.
func TestBug9_MaxConcurrentEventsLimitsHandlers(t *testing.T) {
	// Build 10 ready-made MESSAGE_CREATE events.
	rawMsg, _ := json.Marshal(map[string]interface{}{
		"id": "1", "channel_id": "2",
		"author":    map[string]interface{}{"id": "3", "username": "u", "discriminator": "0"},
		"content":   "hello",
		"timestamp": "2024-01-01T00:00:00Z",
	})

	const (
		numEvents  = 10
		workers    = 2
		sleepTime  = 100 * time.Millisecond
		minElapsed = time.Duration(numEvents) / workers * sleepTime
	)

	wsURL, closeServer := mockGateway(t)
	defer closeServer()

	c, err := NewClient("Bot fake-token", discord.IntentGuilds,
		options.WithMaxConcurrentEvents(workers),
	)
	assert.NoError(t, err)

	err = c.connectWebsocket(wsURL, false, nil, nil)
	assert.NoError(t, err)
	go func() { _ = c.listenWebsocket() }()

	// Wait for READY before registering handlers.
	select {
	case <-c.Websocket.Ready:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for READY")
	}
	defer c.Shutdown()

	var completed atomic.Int32

	factory := events.EventFactories[events.EventMessageCreate]

	// Register a slow handler.
	c.OnMessageCreate(func(_ *Client, _ *events.MessageCreateEvent) {
		time.Sleep(sleepTime)
		completed.Add(1)
	})

	start := time.Now()

	// Dispatch all 10 events directly to the pool channel.
	for i := 0; i < numEvents; i++ {
		ev := factory()
		if err := json.Unmarshal(rawMsg, ev); err != nil {
			t.Fatal(err)
		}
		job := dispatchJob{event: ev}
		select {
		case c.eventCh <- job:
		case <-time.After(5 * time.Second):
			t.Fatal("event channel blocked unexpectedly")
		}
	}

	// Wait until all handlers have run.
	deadline := time.After(30 * time.Second)
	for completed.Load() < numEvents {
		select {
		case <-deadline:
			t.Fatalf("only %d/%d handlers completed within deadline (Bug 9)", completed.Load(), numEvents)
		case <-time.After(20 * time.Millisecond):
		}
	}

	elapsed := time.Since(start)
	assert.GreaterOrEqual(t, elapsed, minElapsed,
		"with %d workers and %v sleep, %d events must take >= %v (Bug 9); got %v",
		workers, sleepTime, numEvents, minElapsed, elapsed)
}
