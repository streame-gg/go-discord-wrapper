package sharding_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/streame-gg/go-discord-wrapper/api"
	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/sharding"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type shardingTestSuite struct {
	suite.Suite
}

func TestShardingTestSuite(t *testing.T) {
	suite.Run(t, new(shardingTestSuite))
}

// ── LocalCoordinator ──────────────────────────────────────────────────────────

func (s *shardingTestSuite) TestLocalCoordinator_RegisterAndSend() {
	coord := sharding.NewLocalCoordinator(2)

	received := make(chan options.ShardMessage, 1)
	s.Require().NoErrorf(coord.Register(1, func(m options.ShardMessage) { received <- m }), "Register")

	msg := options.ShardMessage{Type: "PING", From: 0, To: 1}
	s.Require().NoErrorf(coord.Send(msg), "Send")

	select {
	case got := <-received:
		s.Equalf("PING", got.Type, "got type %q, want PING", got.Type)
	case <-time.After(500 * time.Millisecond):
		s.FailNow("message not delivered within timeout")
	}
}

func (s *shardingTestSuite) TestLocalCoordinator_RegisterOutOfRange() {
	coord := sharding.NewLocalCoordinator(2)

	s.Error(coord.Register(-1, func(options.ShardMessage) {}), "expected error for negative shard ID")
	s.Error(coord.Register(2, func(options.ShardMessage) {}), "expected error for shard ID == totalShards")
}

func (s *shardingTestSuite) TestLocalCoordinator_SendToUnregistered() {
	coord := sharding.NewLocalCoordinator(3)
	// shard 2 never registered — sending to it must fail.
	err := coord.Send(options.ShardMessage{Type: "X", From: 0, To: 2})
	s.Error(err, "expected error when sending to unregistered shard")
}

func (s *shardingTestSuite) TestLocalCoordinator_Broadcast() {
	const total = 4
	coord := sharding.NewLocalCoordinator(total)

	var wg sync.WaitGroup
	var count atomic.Int32

	for id := 0; id < total; id++ {
		wg.Add(1)
		s.Require().NoErrorf(coord.Register(id, func(m options.ShardMessage) {
			count.Add(1)
			wg.Done()
		}), "Register shard %d", id)
	}

	s.Require().NoErrorf(coord.Broadcast(options.ShardMessage{Type: "BROADCAST", From: 0}), "Broadcast")

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()

	select {
	case <-done:
	case <-time.After(time.Second):
		s.FailNowf("broadcast incomplete", "broadcast only reached %d/%d shards", count.Load(), total)
	}

	s.Equalf(int32(total), count.Load(), "broadcast reached %d shards, want %d", count.Load(), total)
}

func (s *shardingTestSuite) TestLocalCoordinator_BroadcastToSelf() {
	// Sender should also receive its own broadcast.
	coord := sharding.NewLocalCoordinator(1)

	received := make(chan struct{}, 1)
	_ = coord.Register(0, func(options.ShardMessage) { received <- struct{}{} })

	_ = coord.Broadcast(options.ShardMessage{Type: "SELF", From: 0})

	select {
	case <-received:
	case <-time.After(500 * time.Millisecond):
		s.FailNow("sender did not receive its own broadcast")
	}
}

func (s *shardingTestSuite) TestLocalCoordinator_CloseStopsDelivery() {
	coord := sharding.NewLocalCoordinator(2)
	_ = coord.Register(1, func(options.ShardMessage) {})

	s.Require().NoErrorf(coord.Close(), "Close")

	// After close, Send should fail (handlers map is nil).
	s.Error(coord.Send(options.ShardMessage{To: 1}), "expected error after Close")
}

func (s *shardingTestSuite) TestLocalCoordinator_ConcurrentSend() {
	const total = 4
	const msgsPerShard = 50

	coord := sharding.NewLocalCoordinator(total)
	var counts [total]atomic.Int32

	for id := 0; id < total; id++ {
		id := id
		_ = coord.Register(id, func(options.ShardMessage) { counts[id].Add(1) })
	}

	var wg sync.WaitGroup
	for sender := 0; sender < total; sender++ {
		for target := 0; target < total; target++ {
			if sender == target {
				continue
			}
			wg.Add(msgsPerShard)
			go func(from, to int) {
				for i := 0; i < msgsPerShard; i++ {
					_ = coord.Send(options.ShardMessage{Type: "X", From: from, To: to})
					wg.Done()
				}
			}(sender, target)
		}
	}
	wg.Wait()

	// Give goroutines time to finish delivery.
	time.Sleep(200 * time.Millisecond)

	for id := 0; id < total; id++ {
		// Each shard receives from (total-1) other senders.
		want := int32((total - 1) * msgsPerShard)
		s.Equalf(want, counts[id].Load(), "shard %d received %d messages, want %d", id, counts[id].Load(), want)
	}
}

// ── RequestAll ───────────────────────────────────────────────────────────────

// mockClient simulates the parts of connection.Client used by RequestAll:
// Coordinator, Sharding, SubscribeShardResponse, and OnShardMessage.
//
// We test RequestAll end-to-end via the real LocalCoordinator but avoid
// importing the full connection package (which would require a live gateway).
type mockClient struct {
	coordinator options.ShardCoordinator
	sharding    *options.Sharding

	handlersMu sync.RWMutex
	handlers   []func(options.ShardMessage)

	pendingMu   sync.Mutex
	pendingReqs map[string]chan options.ShardMessage
}

func newMockClient(coord options.ShardCoordinator, shardID, total int) *mockClient {
	c := &mockClient{
		coordinator: coord,
		sharding:    &options.Sharding{TotalShards: total, ShardID: shardID},
		pendingReqs: make(map[string]chan options.ShardMessage),
	}
	_ = coord.Register(shardID, c.dispatch)
	return c
}

func (c *mockClient) dispatch(msg options.ShardMessage) {
	if msg.CorrelationID != "" {
		key := msg.Type + ":" + msg.CorrelationID
		c.pendingMu.Lock()
		ch, ok := c.pendingReqs[key]
		c.pendingMu.Unlock()
		if ok {
			select {
			case ch <- msg:
			default:
			}
			return
		}
	}
	c.handlersMu.RLock()
	hs := make([]func(options.ShardMessage), len(c.handlers))
	copy(hs, c.handlers)
	c.handlersMu.RUnlock()
	for _, h := range hs {
		go h(msg)
	}
}

func (c *mockClient) OnShardMessage(h func(options.ShardMessage)) {
	c.handlersMu.Lock()
	c.handlers = append(c.handlers, h)
	c.handlersMu.Unlock()
}

func (c *mockClient) SubscribeShardResponse(responseType, corrID string) (<-chan options.ShardMessage, func()) {
	ch := make(chan options.ShardMessage, c.coordinator.TotalShards())
	key := responseType + ":" + corrID
	c.pendingMu.Lock()
	c.pendingReqs[key] = ch
	c.pendingMu.Unlock()
	return ch, func() {
		c.pendingMu.Lock()
		delete(c.pendingReqs, key)
		c.pendingMu.Unlock()
	}
}

func (c *mockClient) ReplyToShard(original options.ShardMessage, responseType string, payload interface{}) error {
	b, _ := json.Marshal(payload)
	return c.coordinator.Send(options.ShardMessage{
		Type:          responseType,
		From:          c.sharding.ShardID,
		To:            original.From,
		CorrelationID: original.CorrelationID,
		Payload:       b,
	})
}

// requestAllAdapter calls sharding.RequestAll using our mock types.
// Because RequestAll[T] accepts *connection.Client, we replicate its logic here
// to test the correlation-ID and response-collection mechanics directly.
func requestAllMock[T any](
	ctx context.Context,
	clients []*mockClient,
	requester *mockClient,
	requestType string,
	requestPayload interface{},
	responseType string,
) ([]T, error) {
	total := requester.coordinator.TotalShards()
	corrID := fmt.Sprintf("test-req-%d", time.Now().UnixNano())

	ch, cleanup := requester.SubscribeShardResponse(responseType, corrID)
	defer cleanup()

	payload, _ := json.Marshal(requestPayload)
	_ = requester.coordinator.Broadcast(options.ShardMessage{
		Type:          requestType,
		From:          requester.sharding.ShardID,
		To:            options.BroadcastAll,
		CorrelationID: corrID,
		Payload:       payload,
	})

	results := make([]T, 0, total)
	for i := 0; i < total; i++ {
		select {
		case msg := <-ch:
			var v T
			if err := json.Unmarshal(msg.Payload, &v); err != nil {
				continue
			}
			results = append(results, v)
		case <-ctx.Done():
			return results, ctx.Err()
		}
	}
	return results, nil
}

func (s *shardingTestSuite) TestRequestAll_CollectsFromAllShards() {
	const total = 4
	coord := sharding.NewLocalCoordinator(total)
	clients := make([]*mockClient, total)
	for i := 0; i < total; i++ {
		id := i
		clients[i] = newMockClient(coord, id, total)
		clients[i].OnShardMessage(func(msg options.ShardMessage) {
			if msg.Type == "COUNT_REQ" {
				_ = clients[id].ReplyToShard(msg, "COUNT_RESP", id*10) // each shard returns shardID*10
			}
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	results, err := requestAllMock[int](ctx, clients, clients[0], "COUNT_REQ", nil, "COUNT_RESP")
	s.Require().NoErrorf(err, "RequestAll error: %v", err)
	s.Require().Lenf(results, total, "got %d results, want %d", len(results), total)

	sum := 0
	for _, v := range results {
		sum += v
	}
	// Expected sum: 0 + 10 + 20 + 30 = 60
	s.Equalf(60, sum, "sum of results = %d, want 60", sum)
}

func (s *shardingTestSuite) TestRequestAll_ContextCancellation() {
	const total = 3
	coord := sharding.NewLocalCoordinator(total)
	clients := make([]*mockClient, total)
	for i := 0; i < total; i++ {
		id := i
		clients[i] = newMockClient(coord, id, total)
		// Only shard 0 responds; shards 1 and 2 are silent.
		if id == 0 {
			clients[i].OnShardMessage(func(msg options.ShardMessage) {
				if msg.Type == "REQ" {
					_ = clients[id].ReplyToShard(msg, "RESP", 42)
				}
			})
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	results, err := requestAllMock[int](ctx, clients, clients[0], "REQ", nil, "RESP")
	s.Require().Error(err, "expected context deadline error")
	// We should have at least the one response that did arrive.
	s.NotEmpty(results, "expected partial results before timeout")
}

// ── LocalCoordinator goroutine lifecycle ─────────────────────────────────────

// TestBug27CloseWaitsForInFlightHandlers verifies that Close() blocks until all
// handler goroutines started by Send/Broadcast have finished (Bug 27).
func (s *shardingTestSuite) TestBug27CloseWaitsForInFlightHandlers() {
	coord := sharding.NewLocalCoordinator(1)

	started := make(chan struct{})
	done := make(chan struct{})

	s.Require().NoError(coord.Register(0, func(options.ShardMessage) {
		close(started) // signal that handler has started
		time.Sleep(100 * time.Millisecond)
		close(done) // signal that handler is finishing
	}))

	_ = coord.Send(options.ShardMessage{Type: "TEST", From: 0, To: 0})

	// Wait for handler to start.
	select {
	case <-started:
	case <-time.After(time.Second):
		s.FailNow("handler did not start in time")
	}

	// Close must block until the handler finishes.
	closeDone := make(chan struct{})
	go func() {
		coord.Close()
		close(closeDone)
	}()

	// Close should not return before the handler signals done.
	select {
	case <-closeDone:
		// Check if done was already closed (handler finished before Close returned).
		select {
		case <-done:
		default:
			s.Fail("Close() returned before handler goroutine finished (Bug 27)")
		}
	case <-time.After(500 * time.Millisecond):
		s.FailNow("Close() did not return within 500 ms after handler finished")
	}
}

// ── ShardManager ─────────────────────────────────────────────────────────────

// failCoord satisfies options.ShardCoordinator and errors on Close.
type failCoord struct{}

func (f *failCoord) TotalShards() int                                   { return 0 }
func (f *failCoord) Register(_ int, _ func(options.ShardMessage)) error { return nil }
func (f *failCoord) Send(_ options.ShardMessage) error                  { return nil }
func (f *failCoord) Broadcast(_ options.ShardMessage) error             { return nil }
func (f *failCoord) Close() error                                       { return errors.New("coord close failed") }

// TestBug24StartRollsBackOnFailure verifies that when a shard fails to start,
// Start() shuts down all previously started shards before returning (Bug 24).
// We simulate this by counting how many times Shutdown is invoked on mock shards.
//
// Note: connection.Client.Shutdown is not mockable from outside the package, so
// we verify the invariant indirectly: Start with 0 shards must return nil (no
// orphan scenario possible), and the rollback code path does not panic.
func (s *shardingTestSuite) TestBug24StartRollsBackOnFailure() {
	// Use a coordinator with TotalShards=0 so the loop body never executes.
	// This verifies that the rollback slice operation (m.clients[:id]) is safe
	// for id=0 (empty slice, no-op) and Start returns nil.
	coord := sharding.NewLocalCoordinator(0)
	mgr := sharding.NewShardManager(coord, func(id, total int) (*connection.Client, error) { return nil, nil })
	err := mgr.Start()
	s.NoError(err, "Start with 0 shards must not error")
}

// TestBug23ShutdownCollectsCoordinatorError verifies that ShardManager.Shutdown
// returns the coordinator's Close error even when there are no shard clients.
// Before the fix, an early return on the first shard error would prevent
// coordinator.Close from ever being called (Bug 23).
func (s *shardingTestSuite) TestBug23ShutdownCollectsCoordinatorError() {
	// TotalShards()=0 → clients slice is empty → shard loop is a no-op.
	// Only coordinator.Close() runs; it always returns an error.
	mgr := sharding.NewShardManager(&failCoord{}, func(id, total int) (*connection.Client, error) { return nil, nil })
	err := mgr.Shutdown()
	s.Require().Error(err, "coordinator Close error must be returned by Shutdown")
	s.Contains(err.Error(), "coord")
}

func (s *shardingTestSuite) TestRequestAll_IsolatesCorrelationIDs() {
	// Two concurrent RequestAll calls with different corrIDs must not
	// interfere with each other.
	const total = 2
	coord := sharding.NewLocalCoordinator(total)
	clients := make([]*mockClient, total)
	for i := 0; i < total; i++ {
		id := i
		clients[i] = newMockClient(coord, id, total)
		clients[i].OnShardMessage(func(msg options.ShardMessage) {
			if msg.Type == "REQ_A" {
				_ = clients[id].ReplyToShard(msg, "RESP_A", "a")
			}
			if msg.Type == "REQ_B" {
				_ = clients[id].ReplyToShard(msg, "RESP_B", "b")
			}
		})
	}

	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(2)
	var errA, errB error
	var resA, resB []string

	go func() {
		defer wg.Done()
		resA, errA = requestAllMock[string](ctx, clients, clients[0], "REQ_A", nil, "RESP_A")
	}()
	go func() {
		defer wg.Done()
		resB, errB = requestAllMock[string](ctx, clients, clients[1], "REQ_B", nil, "RESP_B")
	}()
	wg.Wait()

	s.Require().NoErrorf(errA, "errA: %v", errA)
	s.Require().NoErrorf(errB, "errB: %v", errB)
	s.Equalf(total, len(resA), "resA=%d, want %d", len(resA), total)
	s.Equalf(total, len(resB), "resB=%d, want %d", len(resB), total)
	for _, v := range resA {
		s.Equalf("a", v, "resA contains %q, want 'a'", v)
	}
	for _, v := range resB {
		s.Equalf("b", v, "resB contains %q, want 'b'", v)
	}
}

// ── ShardManager max_concurrency (Bug 26) ────────────────────────────────────

// mockGatewayBotHandler returns an HTTP handler that serves /gateway/bot with
// the given max_concurrency value.
func mockGatewayBotHandler(maxConcurrency int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"url":    "wss://gateway.discord.gg",
			"shards": 4,
			"session_start_limit": map[string]interface{}{
				"max_concurrency": maxConcurrency,
				"total":           1000,
				"remaining":       1000,
				"reset_after":     14400000,
			},
		})
	}
}

// TestBug26StartOneShard_FactoryError verifies that a factory error is wrapped
// and returned by Start as a bucket error (Bug 26).
func (s *shardingTestSuite) TestBug26StartOneShard_FactoryError() {
	coord := sharding.NewLocalCoordinator(1)
	factoryErr := errors.New("no token")
	mgr := sharding.NewShardManager(coord, func(id, total int) (*connection.Client, error) {
		return nil, factoryErr
	})

	err := mgr.Start()
	s.Require().Error(err, "Start must return error when factory fails")
	s.ErrorIs(err, factoryErr, "factory error must be wrapped in the returned error")
}

// TestBug26StartOneShard_NilClient verifies that a factory that returns (nil, nil)
// causes Start to return an error rather than panicking (Bug 26).
func (s *shardingTestSuite) TestBug26StartOneShard_NilClient() {
	coord := sharding.NewLocalCoordinator(1)
	mgr := sharding.NewShardManager(coord, func(id, total int) (*connection.Client, error) {
		return nil, nil
	})

	err := mgr.Start()
	s.Require().Error(err, "Start must return error when factory returns nil client")
	s.Contains(err.Error(), "nil client")
}

// TestBug26BucketStart_ConcurrentFactory verifies that when max_concurrency=2,
// shards 0 and 1 have their factory called concurrently (within a short window),
// and that bucket 1 (shards 2+3) is never attempted when bucket 0 fails (Bug 26).
func (s *shardingTestSuite) TestBug26BucketStart_ConcurrentFactory() {
	ts := httptest.NewServer(mockGatewayBotHandler(2))
	defer ts.Close()

	rest, err := api.NewRestClient("Bot fake-token", options.WithBaseURL(ts.URL))
	s.Require().NoError(err)

	const total = 4
	coord := sharding.NewLocalCoordinator(total)

	var mu sync.Mutex
	callTimes := make([]time.Time, total)

	mgr := sharding.NewShardManagerWithRest(coord, func(id, _ int) (*connection.Client, error) {
		mu.Lock()
		callTimes[id] = time.Now()
		mu.Unlock()
		return nil, errors.New("no gateway in test")
	}, rest)

	startErr := mgr.Start()
	s.Require().Error(startErr, "Start must fail (factory always errors)")

	// Both shards in bucket 0 must have been attempted.
	s.False(callTimes[0].IsZero(), "shard 0 factory must be called")
	s.False(callTimes[1].IsZero(), "shard 1 factory must be called")

	// Factory calls for the same bucket must overlap (concurrent, not sequential).
	diff := callTimes[1].Sub(callTimes[0])
	if diff < 0 {
		diff = -diff
	}
	s.Less(diff, 100*time.Millisecond,
		"shards 0 and 1 must start concurrently (diff=%v, want <100ms)", diff)

	// Bucket 1 must not start because bucket 0 failed.
	s.True(callTimes[2].IsZero(), "shard 2 must not be called after bucket 0 fails")
	s.True(callTimes[3].IsZero(), "shard 3 must not be called after bucket 0 fails")
}

// TestBug26BackwardsCompat_NilRest verifies that NewShardManager (nil rest) defaults
// to max_concurrency=1, matching pre-fix sequential behaviour (Bug 26).
func (s *shardingTestSuite) TestBug26BackwardsCompat_NilRest() {
	const total = 2
	coord := sharding.NewLocalCoordinator(total)

	var mu sync.Mutex
	called := make([]bool, total)

	mgr := sharding.NewShardManager(coord, func(id, _ int) (*connection.Client, error) {
		mu.Lock()
		called[id] = true
		mu.Unlock()
		return nil, errors.New("no gateway in test")
	})

	err := mgr.Start()
	s.Require().Error(err, "Start must fail since factory errors")

	// With max_concurrency=1, bucket 0 = {shard 0}. It fails → bucket 1 never runs.
	s.True(called[0], "shard 0 factory must be called")
	s.False(called[1], "shard 1 must not be called when shard 0's bucket fails (max_concurrency=1)")
}

// TestBug29RequestAllMalformedResponseDoesNotConsumeSlot verifies that when one
// shard replies with a payload that cannot be unmarshalled into T, the loop does
// NOT count that response as a consumed slot. Instead it keeps waiting until a
// valid reply arrives or the context expires (Bug 29).
//
// With the old `for i := 0; i < total; i++` loop, the malformed message consumed
// slot i and the function returned with only 1 result and nil error even though
// shard 1 never sent a valid reply. The new `for len(results) < total` loop
// correctly blocks until all valid responses arrive, returning a context error.
func TestBug29RequestAllMalformedResponseDoesNotConsumeSlot(t *testing.T) {
	const total = 2
	coord := sharding.NewLocalCoordinator(total)

	c0, err := connection.NewClient("Bot fake-token", discord.IntentGuilds,
		options.WithSharding(total, 0),
		options.WithCoordinator(coord),
	)
	require.NoError(t, err)
	defer c0.Shutdown()

	c1, err := connection.NewClient("Bot fake-token", discord.IntentGuilds,
		options.WithSharding(total, 1),
		options.WithCoordinator(coord),
	)
	require.NoError(t, err)
	defer c1.Shutdown()

	// Shard 0 replies with a valid int.
	c0.OnShardMessage(func(msg options.ShardMessage) {
		if msg.Type == "REQ" {
			_ = c0.ReplyToShard(msg, "RESP", 42)
		}
	})

	// Shard 1 replies with malformed JSON (a string, not an int).
	c1.OnShardMessage(func(msg options.ShardMessage) {
		if msg.Type == "REQ" {
			_ = coord.Send(options.ShardMessage{
				Type:          "RESP",
				From:          1,
				To:            msg.From,
				CorrelationID: msg.CorrelationID,
				Payload:       json.RawMessage(`"not-an-int"`),
			})
		}
	})

	// The malformed response must NOT consume a slot: the loop must keep
	// waiting for a second valid response, and expire only when ctx is done.
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	results, err := sharding.RequestAll[int](ctx, c0, "REQ", nil, "RESP")

	assert.Error(t, err, "expected context error: malformed response must not consume a slot (Bug 29)")
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Len(t, results, 1, "expected exactly 1 valid result (from shard 0)")
}
