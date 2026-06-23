package connection

import (
	"testing"
	"time"

	"github.com/streame-gg/go-discord-wrapper/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// ── MESSAGE_CREATE test harness ────────────────────────────────────────────────

// msgCreateHarness wires a real Client (fake token) to an on-demand mock gateway
// with a registered OnMessageCreate listener. Tests dispatch a packet with
// send() and read the decoded event the client produced from recv — letting you
// customise each case around the same dispatch -> decode -> handler path.
//
// pushGateway only emits packets when send() is called, so the listener is always
// registered before any dispatch arrives (no timing race to work around).
type msgCreateHarness struct {
	client *Client
	send   func(data map[string]interface{}) // dispatches one MESSAGE_CREATE packet
	recv   <-chan *events.MessageCreateEvent // events the handler received
	stop   func()
}

// newMsgCreateHarness registers the client with a fake token and the given
// intents/options, waits for READY, then attaches the OnMessageCreate listener.
func newMsgCreateHarness(t *testing.T, intents discord.Intent, opts ...options.Option) *msgCreateHarness {
	t.Helper()

	wsURL, push, closeServer := pushGateway(t)

	client, err := NewClient("Bot fake-token", intents, opts...)
	require.NoError(t, err)
	require.NoError(t, client.connectWebsocket(wsURL, false, nil, nil))
	go func() { _ = client.listenWebsocket() }()
	select {
	case <-client.Ready():
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for READY")
	}
	stop := func() { client.Shutdown() }

	recv := make(chan *events.MessageCreateEvent, 16)
	client.OnMessageCreate(func(_ *Client, ev *events.MessageCreateEvent) {
		recv <- ev
	})

	return &msgCreateHarness{
		client: client,
		send: func(data map[string]interface{}) {
			push(dispatchPacket("MESSAGE_CREATE", data))
		},
		recv: recv,
		stop: func() { stop(); closeServer() },
	}
}

// await blocks for the next received event, failing the test on timeout.
func (h *msgCreateHarness) await(t *testing.T, timeout time.Duration) *events.MessageCreateEvent {
	t.Helper()
	select {
	case ev := <-h.recv:
		return ev
	case <-time.After(timeout):
		t.Fatal("timeout waiting for MESSAGE_CREATE")
		return nil
	}
}

// expectNone asserts no event arrives within the window — used to verify the
// handler does NOT fire for a given dispatch.
func (h *msgCreateHarness) expectNone(t *testing.T, within time.Duration) {
	t.Helper()
	select {
	case ev := <-h.recv:
		t.Fatalf("unexpected MESSAGE_CREATE received: %+v", ev)
	case <-time.After(within):
	}
}

// ── Behavioural tests ──────────────────────────────────────────────────────────

// TestMessageCreateRoundTrip asserts a full payload survives the gateway ->
// decode -> dispatch path intact, including the outer guild_id/member fields that
// MessageCreateEvent.UnmarshalJSON splices on top of the embedded Message.
func (cs *ConnectionSuite) TestMessageCreateRoundTrip() {
	t := cs.T()

	h := newMsgCreateHarness(t, discord.IntentGuilds)
	defer h.stop()

	h.send(map[string]interface{}{
		"id":         "917800000000000001",
		"channel_id": "917800000000000002",
		"guild_id":   "917800000000000003",
		"content":    "hello from the gateway",
		"timestamp":  "2024-01-01T00:00:00Z",
		"author": map[string]interface{}{
			"id":            "917800000000000004",
			"username":      "alice",
			"discriminator": "0001",
		},
		"member": map[string]interface{}{
			"nick":  "Ally",
			"roles": []interface{}{},
		},
	})

	ev := h.await(t, 5*time.Second)
	require.NotNil(t, ev)

	// Embedded Message fields.
	assert.Equal(t, "917800000000000001", ev.ID.String())
	assert.Equal(t, "917800000000000002", ev.ChannelID.String())
	assert.Equal(t, "hello from the gateway", ev.Content)
	require.NotNil(t, ev.Author)
	assert.Equal(t, "917800000000000004", ev.Author.ID.String())
	assert.Equal(t, "alice", ev.Author.Username)

	// Outer fields MESSAGE_CREATE appends on top of the Message object.
	require.NotNil(t, ev.GuildID)
	assert.Equal(t, "917800000000000003", ev.GuildID.String())
	require.NotNil(t, ev.Member)
	require.NotNil(t, ev.Member.Nick)
	assert.Equal(t, "Ally", *ev.Member.Nick)
}

// TestMessageCreateDMHasNoGuildID verifies a DM message (no guild_id/member)
// decodes with those outer fields left nil.
func (cs *ConnectionSuite) TestMessageCreateDMHasNoGuildID() {
	t := cs.T()

	h := newMsgCreateHarness(t, discord.IntentGuilds)
	defer h.stop()

	h.send(map[string]interface{}{
		"id":         "917800000000000010",
		"channel_id": "917800000000000011",
		"content":    "dm message",
		"timestamp":  "2024-01-01T00:00:00Z",
		"author": map[string]interface{}{
			"id":            "917800000000000012",
			"username":      "bob",
			"discriminator": "0002",
		},
	})

	ev := h.await(t, 5*time.Second)
	require.NotNil(t, ev)
	assert.Equal(t, "dm message", ev.Content)
	assert.Nil(t, ev.GuildID, "DM message must have no guild_id")
	assert.Nil(t, ev.Member, "DM message must have no member")
}

// TestMessageCreateMemberNil verifies, by listening on the real client, that the
// event's Member field tracks the presence of the outer "member" payload:
// nil when MESSAGE_CREATE omits it (e.g. a DM), non-nil when Discord includes it.
//
// Note: there is no Message.Member — Member lives on MessageCreateEvent and is
// populated by its UnmarshalJSON, so the dispatch path is what exercises it.
func (cs *ConnectionSuite) TestMessageCreateMemberNil() {
	t := cs.T()

	h := newMsgCreateHarness(t, discord.IntentGuilds)
	defer h.stop()

	base := func() map[string]interface{} {
		return map[string]interface{}{
			"id":         "919900000000000001",
			"channel_id": "919900000000000002",
			"content":    "hi",
			"timestamp":  "2024-01-01T00:00:00Z",
			"author": map[string]interface{}{
				"id":            "919900000000000003",
				"username":      "dave",
				"discriminator": "0004",
			},
		}
	}

	// No "member" in the payload -> Member must be nil.
	h.send(base())
	ev := h.await(t, 5*time.Second)
	assert.Nil(t, ev.Member, "Member must be nil when the payload omits it")

	// "member" present -> Member must be populated.
	withMember := base()
	withMember["guild_id"] = "919900000000000005"
	withMember["member"] = map[string]interface{}{
		"nick":  "Davey",
		"roles": []interface{}{},
	}
	h.send(withMember)
	ev = h.await(t, 5*time.Second)
	require.NotNil(t, ev.Member, "Member must be set when the payload includes it")
	require.NotNil(t, ev.Member.Nick)
	assert.Equal(t, "Davey", *ev.Member.Nick)
}

// TestMessageCreateCapturesMember dispatches a fully-populated outer "member"
// and captures it on the received event, asserting the logic that runs between
// send and receive reconstructs each field.
//
// Two stages matter here:
//   - MessageCreateEvent.UnmarshalJSON decodes the outer member from the payload.
//   - hydrateEvent (gateway.go) then backfills the json:"-" linkage fields that
//     are NOT in the wire payload: Member.GuildID from the event's guild_id, and
//     Member.UserID (via Member.Hydrate) from member.user.id.
//
// This is the cache-independent hydration path — the harness has no cache, yet
// the linkage fields are still populated.
func (cs *ConnectionSuite) TestMessageCreateCapturesMember() {
	t := cs.T()

	h := newMsgCreateHarness(t, discord.IntentGuilds)
	defer h.stop()

	h.send(map[string]interface{}{
		"id":         "920000000000000001",
		"channel_id": "920000000000000002",
		"guild_id":   "920000000000000003",
		"member": map[string]interface{}{
			"nick":      "Dave",
			"roles":     []interface{}{},
			"joined_at": "2024-01-01T00:00:00Z",
			"flags":     10,
		},
		"content":   "with a rich member",
		"timestamp": "2024-01-01T00:00:00Z",
		"author": map[string]interface{}{
			"id":            "920000000000000004",
			"username":      "DGamet",
			"discriminator": "0001",
			"bot":           false,
		},
		"type":   0,
		"pinned": false,
	})

	ev := h.await(t, 5*time.Second)
	require.NotNil(t, ev.Member, "Member must be captured from the dispatched payload")

	cs.Require().True(ev.Member.IsHydrated(), "Member must be hydrated")
	cs.Require().Equal("<@920000000000000004>", ev.Member.Mention(), "Member must mention the user ID")
	cs.Require().NotNil(ev.Member.User, "Member must have a user")
	cs.Require().Equal("920000000000000004", ev.Member.User.ID.String(), "Member.User.ID must match the author ID")
	cs.Require().Equal(ev.Author, *ev.Member.User, "Author and user should match")

	ev.Author.Bot = util.PointerOf(true)

	// ev.Member.User should be a copy of ev.Author and not directly point at the memory address of ev.Author
	cs.Require().True(*ev.Author.Bot)
	cs.Require().False(*ev.Member.User.Bot)
}

// sendN dispatches len(contents) MESSAGE_CREATE packets, one per content string.
func (h *msgCreateHarness) sendN(contents []string) {
	for _, body := range contents {
		h.send(map[string]interface{}{
			"id":         "918800000000000000",
			"channel_id": "918800000000000001",
			"content":    body,
			"timestamp":  "2024-01-01T00:00:00Z",
			"author": map[string]interface{}{
				"id":            "918800000000000002",
				"username":      "carol",
				"discriminator": "0003",
			},
		})
	}
}

// TestMessageCreateDeliversAll dispatches several packets through a DEFAULT
// client and asserts every one is delivered. It deliberately does NOT assert
// order: the default client runs a worker pool (MaxConcurrentEvents defaults to
// 64), so handlers for distinct events run concurrently and may interleave.
func (cs *ConnectionSuite) TestMessageCreateDeliversAll() {
	t := cs.T()

	h := newMsgCreateHarness(t, discord.IntentGuilds)
	defer h.stop()

	contents := []string{"first", "second", "third"}
	h.sendN(contents)

	got := make(map[string]bool)
	for range contents {
		got[h.await(t, 5*time.Second).Content] = true
	}
	for _, want := range contents {
		assert.True(t, got[want], "expected to receive %q", want)
	}
}

// TestMessageCreateSerialPreservesOrder shows that pinning the worker pool to a
// single worker (WithMaxConcurrentEvents(1)) makes dispatch serial, so events
// are delivered in send order.
func (cs *ConnectionSuite) TestMessageCreateSerialPreservesOrder() {
	t := cs.T()

	h := newMsgCreateHarness(t, discord.IntentGuilds, options.WithMaxConcurrentEvents(1))
	defer h.stop()

	contents := []string{"first", "second", "third"}
	h.sendN(contents)

	for _, want := range contents {
		assert.Equal(t, want, h.await(t, 5*time.Second).Content)
	}
}
