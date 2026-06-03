package connection

import (
	"sync/atomic"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// TestBug122_EntitlementUpdateEvent_OldEntitlementPopulated verifies that
// OldEntitlement is populated from the in-memory cache when ENTITLEMENT_UPDATE
// arrives after a preceding ENTITLEMENT_CREATE (Bug 122).
func (cs *ConnectionSuite) TestBug122_EntitlementUpdateEvent_OldEntitlementPopulated() {
	t := cs.T()
	entitlementID := "175928847299117063"

	createPacket := dispatchPacket("ENTITLEMENT_CREATE", map[string]interface{}{
		"id":             entitlementID,
		"sku_id":         "100000000000000001",
		"application_id": "200000000000000001",
		"type":           1,
		"deleted":        false,
	})

	updatePacket := dispatchPacket("ENTITLEMENT_UPDATE", map[string]interface{}{
		"id":             entitlementID,
		"sku_id":         "100000000000000001",
		"application_id": "200000000000000001",
		"type":           1,
		"deleted":        false,
		"consumed":       true,
	})

	wsURL, closeServer := mockGateway(t, createPacket, updatePacket)
	defer closeServer()

	client, stop := launchClient(t, wsURL, discord.IntentGuilds)
	defer stop()

	var gotOld atomic.Pointer[discord.Entitlement]
	done := make(chan struct{})
	once := make(chan struct{}, 1)
	client.OnEntitlementUpdate(func(c *Client, ev *events.EntitlementUpdateEvent) {
		select {
		case once <- struct{}{}:
			gotOld.Store(ev.OldEntitlement)
			close(done)
		default:
		}
	})

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for ENTITLEMENT_UPDATE")
	}

	old := gotOld.Load()
	require.NotNil(t, old, "OldEntitlement must be non-nil after ENTITLEMENT_CREATE was seen (Bug 122)")
	assert.Equal(t, discord.Snowflake(175928847299117063), old.ID)
}

// TestBug122_EntitlementUpdateEvent_OldEntitlementNilWhenMissed verifies that
// OldEntitlement is nil when no ENTITLEMENT_CREATE preceded the update (cache miss).
func (cs *ConnectionSuite) TestBug122_EntitlementUpdateEvent_OldEntitlementNilWhenMissed() {
	t := cs.T()
	updatePacket := dispatchPacket("ENTITLEMENT_UPDATE", map[string]interface{}{
		"id":             "175928847299117064",
		"sku_id":         "100000000000000002",
		"application_id": "200000000000000001",
		"type":           1,
		"deleted":        false,
	})

	wsURL, closeServer := mockGateway(t, updatePacket)
	defer closeServer()

	client, stop := launchClient(t, wsURL, discord.IntentGuilds)
	defer stop()

	var gotOld atomic.Pointer[discord.Entitlement]
	done := make(chan struct{})
	once := make(chan struct{}, 1)
	client.OnEntitlementUpdate(func(c *Client, ev *events.EntitlementUpdateEvent) {
		select {
		case once <- struct{}{}:
			gotOld.Store(ev.OldEntitlement)
			close(done)
		default:
		}
	})

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for ENTITLEMENT_UPDATE (miss case)")
	}

	assert.Nil(t, gotOld.Load(), "OldEntitlement must be nil when no preceding CREATE was cached (Bug 122)")
}
