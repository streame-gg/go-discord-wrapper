package options

import "encoding/json"

// BroadcastAll is the sentinel value for ShardMessage.To that means
// "deliver to every registered shard".
const BroadcastAll = -1

// Reserved pkg message types used by the built-in cross-shard statistics
// helpers (sharding.GetTotalServers, sharding.GetTotalUsers).
// Do not use these strings in your own OnShardMessage handlers.
const (
	ShardMsgServerCountReq  = "__shard_server_count_req"
	ShardMsgServerCountResp = "__shard_server_count_resp"
	ShardMsgMemberCountReq  = "__shard_member_count_req"
	ShardMsgMemberCountResp = "__shard_member_count_resp"
)

// ShardMessage is the envelope for all inter-shard communication.
type ShardMessage struct {
	// Type is a user-defined string identifying the kind of message.
	// Use a consistent naming convention, e.g. "GUILD_COUNT_REQ" / "GUILD_COUNT_RESP".
	Type string `json:"type"`

	// From is the sender's shard ID.
	From int `json:"from"`

	// To is the target shard ID, or BroadcastAll (-1) for every shard.
	To int `json:"to"`

	// CorrelationID links a response back to its originating request.
	// Set this when using RequestAll; echo it verbatim in the reply.
	CorrelationID string `json:"correlation_id,omitempty"`

	// Payload is the JSON-encoded message body.
	// Marshal your data with encoding/json before setting this field.
	Payload json.RawMessage `json:"payload,omitempty"`
}

// ShardCoordinator abstracts inter-shard message routing.
//
// For single-machine deployments use sharding.NewLocalCoordinator, which routes
// messages between goroutines in the same process.
//
// For multi-machine deployments implement this interface backed by an external
// message broker — Redis Pub/Sub, NATS, RabbitMQ, gRPC streams, etc. The
// library intentionally does not prescribe a specific transport so you can
// choose whatever fits your infrastructure.
//
// Implementing the interface requires only four methods; see LocalCoordinator
// in the sharding package for a reference implementation.
type ShardCoordinator interface {
	// Register connects shardID to the coordinator and installs handler for
	// all messages addressed to that shard. Must be called once per shard
	// before Send or Broadcast.
	Register(shardID int, handler func(ShardMessage)) error

	// Send delivers msg to the single shard identified by msg.To.
	Send(msg ShardMessage) error

	// Broadcast delivers msg to every registered shard, including the sender.
	Broadcast(msg ShardMessage) error

	// TotalShards returns the total number of shards this coordinator manages.
	// The value must be stable for the lifetime of the coordinator: RequestAll
	// reads it once at the start of a call and uses it as the expected response
	// count. Changing it mid-flight causes RequestAll to under- or over-collect.
	TotalShards() int

	// Close releases coordinator resources. Safe to call multiple times.
	Close() error
}
