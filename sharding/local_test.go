package sharding

import (
	"encoding/json"
	"fmt"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/options"
)

type shardingSuite struct{ suite.Suite }

func TestShardingSuite(t *testing.T) { suite.Run(t, new(shardingSuite)) }

func (s *shardingSuite) TestCreateAndTotalShards() {
	c := NewLocalCoordinator(2)
	s.Equal(2, c.TotalShards())

	c = NewLocalCoordinator(-1)
	s.Equal(0, c.TotalShards())
}

func (s *shardingSuite) TestRegister() {
	c := NewLocalCoordinator(2)
	s.EqualError(c.Register(3, func(options.ShardMessage) {}), "sharding: shard ID 3 out of range [0, 2]")
	s.EqualError(c.Register(-1, func(options.ShardMessage) {}), "sharding: shard ID -1 out of range [0, 2]")

	s.Require().NoError(c.Close())
	s.EqualError(c.Register(2, func(options.ShardMessage) {}), "sharding: coordinator is closed")

	handler := func(options.ShardMessage) {
		fmt.Print("lol!")
	}

	s.Require().NoError(c.Register(2, handler))

	s.Equal(handler, c.handlers[1])
}

func (s *shardingSuite) TestSend() {
	c := NewLocalCoordinator(2)
	s.Require().NoError(c.Close())
	s.EqualError(c.Send(options.ShardMessage{}), "sharding: coordinator is closed")

	c = NewLocalCoordinator(2)
	s.EqualError(c.Send(options.ShardMessage{To: 3}), "sharding: shard 3 is not registered")

	received := make(chan options.ShardMessage, 1)

	c = NewLocalCoordinator(2)
	s.Require().NoError(c.Register(0, func(msg options.ShardMessage) {
		received <- msg
	}))

	want := options.ShardMessage{To: 0, Payload: json.RawMessage("hello")}

	s.NoError(c.Send(want))

	select {
	case got := <-received:
		s.Equal(want, got)
	case <-time.After(time.Second):
		s.Fail("handler was not called in time")
	}
}

func (s *shardingSuite) TestBroadcast() {
	c := NewLocalCoordinator(2)
	s.NoError(c.Broadcast(options.ShardMessage{})) // no handlers registered shouldn't return an error

	c = NewLocalCoordinator(2)

	s.Require().NoError(c.Close())
	s.EqualError(c.Broadcast(options.ShardMessage{}), "sharding: coordinator is closed")

	receivedShardOne := make(chan options.ShardMessage, 1)
	receivedShardTwo := make(chan options.ShardMessage, 1)

	c = NewLocalCoordinator(2)
	s.Require().NoError(c.Register(0, func(msg options.ShardMessage) {
		receivedShardOne <- msg
	}))
	s.Require().NoError(c.Register(1, func(msg options.ShardMessage) {
		receivedShardTwo <- msg
	}))

	want := options.ShardMessage{To: options.BroadcastAll, Payload: json.RawMessage("broadcast all")}

	s.NoError(c.Broadcast(want), "Sent broadcast to 2 shards")

	gotMessages := 0

	for {
		select {
		case got := <-receivedShardOne:
			s.T().Log("Received message for Shard 1")
			gotMessages++
			s.Equal(want, got)

			if gotMessages == 2 {
				return
			}
		case got := <-receivedShardTwo:
			s.T().Log("Received message for Shard 2")
			gotMessages++
			s.Equal(want, got)

			if gotMessages == 2 {
				return
			}
		case <-time.After(time.Second):
			s.Fail("handler was not called in time")
		}
	}
}

func (s *shardingSuite) TestClose() {
	c := NewLocalCoordinator(2)
	s.Require().NoError(c.Close())

	synctest.Test(s.T(), func(t *testing.T) {
		c := NewLocalCoordinator(1)

		proceed := make(chan struct{})
		finished := false

		s.Require().NoError(c.Register(0, func(options.ShardMessage) {
			<-proceed
			finished = true
		}))
		s.Require().NoError(c.Send(options.ShardMessage{To: 0}))

		closeDone := make(chan struct{})
		go func() {
			c.Close()
			close(closeDone)
		}()

		synctest.Wait()

		select {
		case <-closeDone:
			s.Fail("Close() kehrte zurück, bevor der Handler fertig war")
		default:
		}
		s.False(finished)

		close(proceed)
		synctest.Wait()

		select {
		case <-closeDone:
		default:
			s.Fail("Close() kehrte nicht zurück, obwohl der Handler fertig ist")
		}
		s.True(finished)
	})
}
