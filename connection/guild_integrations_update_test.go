package connection

import (
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// TestGuildIntegrationsUpdateDispatch verifies the coarse GUILD_INTEGRATIONS_UPDATE
// gateway event is delivered end-to-end to OnGuildIntegrationsUpdate handlers with
// the guild ID parsed. The event has no cache side-effect; it must still dispatch.
func (cs *ConnectionSuite) TestGuildIntegrationsUpdateDispatch() {
	t := cs.T()
	const guildID = discord.Snowflake(100000000000000001)

	packet := dispatchPacket("GUILD_INTEGRATIONS_UPDATE", map[string]interface{}{
		"guild_id": guildID.String(),
	})

	wsURL, closeServer := mockGateway(t, packet)
	defer closeServer()

	client, stop := launchClient(t, wsURL, discord.IntentGuildIntegrations)
	defer stop()

	got := make(chan discord.Snowflake, 1)
	client.OnGuildIntegrationsUpdate(func(_ *Client, ev *events.GuildIntegrationsUpdateEvent) {
		got <- ev.GuildID
	})

	select {
	case id := <-got:
		assert.Equal(t, guildID, id)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for GUILD_INTEGRATIONS_UPDATE handler")
	}
}
