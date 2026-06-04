package connection

import (
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
	"github.com/streame-gg/go-discord-wrapper/types/interactions/responses"
)

// TestE2E_InteractionCreate_StringOptionReachesHandler is an end-to-end test of
// the full receive path: the mock gateway pushes an INTERACTION_CREATE for a
// slash command with a string option, and the test asserts the registered
// handler receives it with the option value decoded correctly. This guards the
// option-decoder fix (string/bool options previously failed to unmarshal) at the
// gateway/dispatch level, not just in the responses package.
func (cs *ConnectionSuite) TestE2E_InteractionCreate_StringOptionReachesHandler() {
	t := cs.T()

	interaction := map[string]interface{}{
		"id":             "175928847299117063",
		"application_id": "175928847299117064",
		"type":           2, // APPLICATION_COMMAND
		"token":          "tok",
		"version":        1,
		"data": map[string]interface{}{
			"id":   "175928847299117065",
			"name": "echo",
			"type": 1,
			"options": []interface{}{
				map[string]interface{}{"type": 3, "name": "text", "value": "hello"},
			},
		},
	}
	pkt := dispatchPacket("INTERACTION_CREATE", interaction)

	wsURL, closeServer := mockGateway(t, pkt)
	defer closeServer()

	client, stop := launchClient(t, wsURL, discord.IntentGuilds)
	defer stop()

	type result struct {
		full  string
		value string
		okCmd bool
	}
	got := make(chan result, 1)

	client.OnInteractionCreate(func(_ *Client, ev *events.InteractionCreateEvent) {
		r := result{full: ev.GetFullCommand()}
		if cmd, ok := ev.Data.(*responses.InteractionDataApplicationCommand); ok && cmd.Options != nil {
			r.okCmd = true
			if opts := *cmd.Options; len(opts) == 1 && opts[0].Value != nil {
				if sv, isStr := (*opts[0].Value).(string); isStr {
					r.value = sv
				}
			}
		}
		got <- r
	})

	select {
	case r := <-got:
		cs.True(r.okCmd, "interaction data should decode as an application command")
		cs.Equal("echo", r.full, "command name should round-trip through the gateway")
		cs.Equal("hello", r.value, "string option value must decode end-to-end")
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for INTERACTION_CREATE handler")
	}
}
