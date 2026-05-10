package sharding

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/options"
)

// corrSeq is a process-wide monotonic counter used to generate unique correlation IDs.
var corrSeq atomic.Int64

// newCorrID returns a unique correlation ID string for the given shard.
func newCorrID(shardID int) string {
	return fmt.Sprintf("shard%d-req%d", shardID, corrSeq.Add(1))
}

// RequestAll broadcasts a request to all shards and collects typed responses.
//
// # How it works
//
//  1. A unique CorrelationID is generated and embedded in the broadcast message.
//  2. Every shard (including the caller) receives the broadcast via OnShardMessage.
//  3. Each shard handler computes its local answer and calls client.ReplyToShard,
//     which echoes the CorrelationID back to the originating shard.
//  4. RequestAll collects one response per shard and returns when all have
//     replied or when ctx / timeout expires.
//
// Partial results are returned on timeout — check len(results) vs TotalShards().
//
// # Example — total guild count across all shards
//
//	// Register on every shard (including the one that calls RequestAll):
//	client.OnShardMessage(func(msg options.ShardMessage) {
//	    if msg.Type == "GUILD_COUNT_REQ" {
//	        count := myGuildCount() // local to this shard
//	        if err := client.ReplyToShard(msg, "GUILD_COUNT_RESP", count); err != nil {
//	            log.Println("reply failed:", err)
//	        }
//	    }
//	})
//
//	// On the shard that needs the aggregate:
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	counts, err := sharding.RequestAll[int](ctx, client, "GUILD_COUNT_REQ", nil, "GUILD_COUNT_RESP")
//	if err != nil { ... }
//	total := 0
//	for _, n := range counts { total += n }
func RequestAll[T any](
	ctx context.Context,
	client *connection.Client,
	requestType string,
	requestPayload interface{},
	responseType string,
) ([]T, error) {
	if client.Coordinator == nil {
		return nil, fmt.Errorf("sharding: client has no coordinator; use options.WithCoordinator")
	}
	if client.Sharding == nil {
		return nil, fmt.Errorf("sharding: client has no sharding config; use options.WithSharding")
	}

	total := client.Coordinator.TotalShards()
	corrID := newCorrID(client.Sharding.ShardID)

	// Subscribe for responses before broadcasting so we don't miss any fast replies.
	ch, cleanup := client.SubscribeShardResponse(responseType, corrID)
	defer cleanup()

	// Broadcast the request.
	payload, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("sharding: marshal request payload: %w", err)
	}
	if err := client.Coordinator.Broadcast(options.ShardMessage{
		Type:          requestType,
		From:          client.Sharding.ShardID,
		To:            options.BroadcastAll,
		CorrelationID: corrID,
		Payload:       payload,
	}); err != nil {
		return nil, fmt.Errorf("sharding: broadcast failed: %w", err)
	}

	// Collect one response per shard. Malformed responses do not consume a
	// slot so the loop keeps waiting until total valid results arrive or the
	// context expires.
	results := make([]T, 0, total)
	for len(results) < total {
		select {
		case msg := <-ch:
			var v T
			if err := json.Unmarshal(msg.Payload, &v); err != nil {
				// Skip malformed response; do NOT advance the counter so we
				// still wait for a valid reply (or context cancellation).
				continue
			}
			results = append(results, v)
		case <-ctx.Done():
			// Return whatever arrived so far; caller can detect partial results.
			return results, ctx.Err()
		}
	}

	return results, nil
}
