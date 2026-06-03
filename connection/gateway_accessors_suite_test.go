package connection

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// stubCoordinator is a minimal in-memory ShardCoordinator. The sharding package's
// LocalCoordinator can't be imported here (it imports connection → import cycle),
// so the test supplies its own. Send/Broadcast invoke registered handlers
// synchronously; failSend forces the error path.
type stubCoordinator struct {
	handlers map[int]func(options.ShardMessage)
	sent     []options.ShardMessage
}

func newStubCoordinator() *stubCoordinator {
	return &stubCoordinator{handlers: map[int]func(options.ShardMessage){}}
}

func (s *stubCoordinator) Register(shardID int, h func(options.ShardMessage)) error {
	s.handlers[shardID] = h
	return nil
}

func (s *stubCoordinator) Send(msg options.ShardMessage) error {
	s.sent = append(s.sent, msg)
	if h, ok := s.handlers[msg.To]; ok {
		h(msg)
	}
	return nil
}

func (s *stubCoordinator) Broadcast(msg options.ShardMessage) error {
	s.sent = append(s.sent, msg)
	for _, h := range s.handlers {
		h(msg)
	}
	return nil
}

func (s *stubCoordinator) TotalShards() int { return len(s.handlers) }
func (s *stubCoordinator) Close() error     { return nil }

// gatewayAccessorsSuite covers the lightweight *Client accessors and the
// shard-messaging surface, which don't need a live websocket: they read internal
// state guarded by mutexes or route through a ShardCoordinator. The shard methods
// are exercised on both the unconfigured path (returns an error) and the
// configured path (LocalCoordinator + Sharding).
type gatewayAccessorsSuite struct {
	suite.Suite
	client *Client
}

func TestGatewayAccessorsSuite(t *testing.T) { suite.Run(t, new(gatewayAccessorsSuite)) }

func (s *gatewayAccessorsSuite) SetupTest() {
	c, err := NewClient("Bot fake-token", discord.IntentGuilds)
	s.Require().NoError(err)
	s.client = c
}

func (s *gatewayAccessorsSuite) TestBotUser() {
	s.Nil(s.client.BotUser(), "nil before READY")
	u := &discord.User{ID: sampleSnowflake, Username: "bot"}
	s.client.setBotUser(u)
	s.Equal(u, s.client.BotUser())
}

func (s *gatewayAccessorsSuite) TestGuildAndMemberCounts() {
	c := s.client
	s.Equal(0, c.GuildCount())
	s.Equal(0, c.MemberCount())

	c.guildMu.Lock()
	c.guildMemberCounts = map[discord.Snowflake]int{sampleSnowflake: 5, sampleSnowflake + 1: 7}
	c.guildMu.Unlock()

	s.Equal(2, c.GuildCount())
	s.Equal(12, c.MemberCount())
}

func (s *gatewayAccessorsSuite) TestUnavailableGuilds() {
	c := s.client
	id := sampleSnowflake
	c.mu.Lock()
	c.UnavailableGuilds[id] = struct{}{}
	c.mu.Unlock()

	s.True(c.IsGuildUnavailable(id))
	c.deleteUnavailableGuild(id)
	s.False(c.IsGuildUnavailable(id))
}

func (s *gatewayAccessorsSuite) TestShardMessagingUnconfigured() {
	c := s.client
	msg := options.ShardMessage{From: 1}
	// No coordinator/sharding → each returns an error rather than panicking.
	s.Error(c.SendToShard(0, "t", map[string]int{"a": 1}))
	s.Error(c.BroadcastToShards("t", map[string]int{"a": 1}))
	s.Error(c.ReplyToShard(msg, "t", map[string]int{"a": 1}))
}

func (s *gatewayAccessorsSuite) TestShardMessagingConfigured() {
	c := s.client
	coord := newStubCoordinator()
	c.Coordinator = coord
	c.Sharding = &options.Sharding{TotalShards: 2, ShardID: 0}

	// Register a handler so OnShardMessage's slice is exercised and delivered
	// messages have a destination.
	received := make(chan options.ShardMessage, 4)
	c.OnShardMessage(func(m options.ShardMessage) { received <- m })
	s.Require().NoError(coord.Register(0, func(m options.ShardMessage) { received <- m }))
	s.Require().NoError(coord.Register(1, func(m options.ShardMessage) {}))

	s.NoError(c.SendToShard(1, "ping", map[string]int{"a": 1}))
	s.NoError(c.BroadcastToShards("ping", map[string]int{"a": 1}))
	s.NoError(c.ReplyToShard(options.ShardMessage{From: 1, CorrelationID: "x"}, "pong", 42))

	// dispatchShardMessage: built-in stat requests auto-reply, a correlated
	// response routes to a pending channel, and plain messages fan out to the
	// registered OnShardMessage handlers.
	c.dispatchShardMessage(options.ShardMessage{Type: options.ShardMsgServerCountReq, From: 1})
	c.dispatchShardMessage(options.ShardMessage{Type: options.ShardMsgMemberCountReq, From: 1})

	pendingKey := "corr-resp:abc"
	pending := make(chan options.ShardMessage, 1)
	c.pendingShardReqsMu.Lock()
	if c.pendingShardReqs == nil {
		c.pendingShardReqs = make(map[string]chan options.ShardMessage)
	}
	c.pendingShardReqs[pendingKey] = pending
	c.pendingShardReqsMu.Unlock()
	c.dispatchShardMessage(options.ShardMessage{Type: "corr-resp", CorrelationID: "abc", From: 1})
	s.Len(pending, 1, "correlated response routed to pending channel")

	c.dispatchShardMessage(options.ShardMessage{Type: "plain", From: 1})

	// Close() routes through the coordinator + Shutdown.
	s.NoError(c.Close())
}
