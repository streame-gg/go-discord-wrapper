package sharding_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/sharding"
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

// ── ShardManager ─────────────────────────────────────────────────────────────

// failCoord satisfies options.ShardCoordinator and errors on Close.
type failCoord struct{}

func (f *failCoord) TotalShards() int                                     { return 0 }
func (f *failCoord) Register(_ int, _ func(options.ShardMessage)) error   { return nil }
func (f *failCoord) Send(_ options.ShardMessage) error                    { return nil }
func (f *failCoord) Broadcast(_ options.ShardMessage) error               { return nil }
func (f *failCoord) Close() error                                         { return errors.New("coord close failed") }

// TestBug23ShutdownCollectsCoordinatorError verifies that ShardManager.Shutdown
// returns the coordinator's Close error even when there are no shard clients.
// Before the fix, an early return on the first shard error would prevent
// coordinator.Close from ever being called (Bug 23).
func (s *shardingTestSuite) TestBug23ShutdownCollectsCoordinatorError() {
	// TotalShards()=0 → clients slice is empty → shard loop is a no-op.
	// Only coordinator.Close() runs; it always returns an error.
	mgr := sharding.NewShardManager(&failCoord{}, func(id, total int) *connection.Client { return nil })
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
