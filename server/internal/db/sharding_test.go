package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/example/go-db-prj-sample/internal/db"
)

func TestHashBasedSharding_GetShardID(t *testing.T) {
	tests := []struct {
		name       string
		shardCount int
		key        int64
		wantMin    int
		wantMax    int
	}{
		{
			name:       "2 shards",
			shardCount: 2,
			key:        1,
			wantMin:    1,
			wantMax:    2,
		},
		{
			name:       "4 shards",
			shardCount: 4,
			key:        1,
			wantMin:    1,
			wantMax:    4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := db.NewHashBasedSharding(tt.shardCount)
			shardID := strategy.GetShardID(tt.key)

			assert.GreaterOrEqual(t, shardID, tt.wantMin)
			assert.LessOrEqual(t, shardID, tt.wantMax)
		})
	}
}

func TestHashBasedSharding_Consistency(t *testing.T) {
	// Same key should always return same shard
	strategy := db.NewHashBasedSharding(2)

	key := int64(12345)
	shard1 := strategy.GetShardID(key)
	shard2 := strategy.GetShardID(key)
	shard3 := strategy.GetShardID(key)

	assert.Equal(t, shard1, shard2)
	assert.Equal(t, shard2, shard3)
}

func TestHashBasedSharding_Distribution(t *testing.T) {
	// Test that keys are distributed across shards
	strategy := db.NewHashBasedSharding(2)

	distribution := make(map[int]int)
	for i := int64(1); i <= 100; i++ {
		shardID := strategy.GetShardID(i)
		distribution[shardID]++
	}

	// Both shards should have at least some keys
	assert.Greater(t, distribution[1], 0, "Shard 1 should have some keys")
	assert.Greater(t, distribution[2], 0, "Shard 2 should have some keys")

	// Distribution should be somewhat balanced (not perfect, but reasonable)
	// Allow up to 70-30 split for 100 keys
	assert.Greater(t, distribution[1], 20, "Shard 1 should have at least 20% of keys")
	assert.Greater(t, distribution[2], 20, "Shard 2 should have at least 20% of keys")
}
