package sharding

import (
	"context"
	"fmt"

	"github.com/streame-gg/go-discord-wrapper/connection"
	"github.com/streame-gg/go-discord-wrapper/options"
)

// GetTotalServers returns the total number of guilds the bot is in across
// all shards. It broadcasts a request to every shard and sums their individual
// guild counts. The per-shard counts are maintained automatically from
// GUILD_CREATE and GUILD_DELETE gateway events — no manual setup required.
//
// ctx controls the deadline. If some shards do not respond in time, the
// returned error will be non-nil and the count will be a partial sum.
//
// Requires that every shard was created with options.WithCoordinator.
func GetTotalServers(ctx context.Context, client *connection.Client) (int, error) {
	counts, err := RequestAll[int](
		ctx, client,
		options.ShardMsgServerCountReq, nil,
		options.ShardMsgServerCountResp,
	)
	if err != nil {
		return sumInts(counts), fmt.Errorf("GetTotalServers: %w", err)
	}
	return sumInts(counts), nil
}

// GetTotalUsers returns the total number of guild memberships across all shards.
//
// Specifically it returns the sum of each guild's member_count field as
// reported by Discord in the GUILD_CREATE gateway event. This means:
//
//   - A user who is a member of five guilds contributes 5 to the total.
//   - Bots are included in member counts.
//
// True unique-user counting would require caching every user ID across every
// guild (only feasible with the GUILD_MEMBERS privileged intent and substantial
// memory). Use this value as an engagement metric, not an exact head count.
//
// ctx controls the deadline. If some shards do not respond in time, the
// returned error will be non-nil and the count will be a partial sum.
//
// Requires that every shard was created with options.WithCoordinator.
func GetTotalUsers(ctx context.Context, client *connection.Client) (int, error) {
	counts, err := RequestAll[int](
		ctx, client,
		options.ShardMsgMemberCountReq, nil,
		options.ShardMsgMemberCountResp,
	)
	if err != nil {
		return sumInts(counts), fmt.Errorf("GetTotalUsers: %w", err)
	}
	return sumInts(counts), nil
}

func sumInts(ns []int) int {
	total := 0
	for _, n := range ns {
		total += n
	}
	return total
}
