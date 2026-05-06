package options

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config with sharding",
			config: Config{
				Sharding: &Sharding{TotalShards: 4, ShardID: 0},
			},
			wantErr: false,
		},
		{
			name: "shard ID equals total shards",
			config: Config{
				Sharding: &Sharding{TotalShards: 4, ShardID: 4},
			},
			wantErr: true,
		},
		{
			name: "shard ID greater than total shards",
			config: Config{
				Sharding: &Sharding{TotalShards: 4, ShardID: 5},
			},
			wantErr: true,
		},
		{
			name: "total shards zero",
			config: Config{
				Sharding: &Sharding{TotalShards: 0, ShardID: 0},
			},
			wantErr: true,
		},
		{
			name: "total shards negative",
			config: Config{
				Sharding: &Sharding{TotalShards: -1, ShardID: 0},
			},
			wantErr: true,
		},
		{
			name: "shard ID negative",
			config: Config{
				Sharding: &Sharding{TotalShards: 4, ShardID: -1},
			},
			wantErr: true,
		},
		{
			name: "negative MaxRetries",
			config: Config{
				Retry: RetryOptions{MaxRetries: -1},
			},
			wantErr: true,
		},
		{
			name: "negative MinRequestInterval",
			config: Config{
				MinRequestInterval: -1 * time.Millisecond,
			},
			wantErr: true,
		},
		{
			name: "negative SafetyMargin",
			config: Config{
				RateLimiter: &RateLimiterOptions{SafetyMargin: -1},
			},
			wantErr: true,
		},
		{
			name: "nil sharding with negative MaxRetries",
			config: Config{
				Sharding: nil,
				Retry:    RetryOptions{MaxRetries: -1},
			},
			wantErr: true,
		},
		{
			name:    "all zero values",
			config:  Config{},
			wantErr: false,
		},
		{
			name: "zero MaxRetries is valid",
			config: Config{
				Retry: RetryOptions{MaxRetries: 0},
			},
			wantErr: false,
		},
		{
			name: "zero MinRequestInterval is valid",
			config: Config{
				MinRequestInterval: 0,
			},
			wantErr: false,
		},
		{
			name: "safety margin of zero is valid",
			config: Config{
				RateLimiter: &RateLimiterOptions{SafetyMargin: 0},
			},
			wantErr: false,
		},
		{
			name: "last valid shard ID",
			config: Config{
				Sharding: &Sharding{TotalShards: 4, ShardID: 3},
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			if tc.wantErr {
				assert.Error(t, err, "expected an error but got nil")
			} else {
				assert.NoError(t, err, "expected no error but got: %v", err)
			}
		})
	}
}
