package db_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/example/go-webdb-template/internal/config"
	"github.com/example/go-webdb-template/internal/db"
	"github.com/example/go-webdb-template/internal/model"
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

// TestGORMManagerSharding tests that GORMManager correctly routes to shards
func TestGORMManagerSharding(t *testing.T) {
	tmpDir := t.TempDir()

	// Create 3 shards
	shards := []config.ShardConfig{
		{
			ID:         1,
			Driver:     "sqlite3",
			DSN:        filepath.Join(tmpDir, "shard1.db"),
			WriterDSN:  filepath.Join(tmpDir, "shard1.db"),
			ReaderDSNs: []string{filepath.Join(tmpDir, "shard1.db")},
		},
		{
			ID:         2,
			Driver:     "sqlite3",
			DSN:        filepath.Join(tmpDir, "shard2.db"),
			WriterDSN:  filepath.Join(tmpDir, "shard2.db"),
			ReaderDSNs: []string{filepath.Join(tmpDir, "shard2.db")},
		},
		{
			ID:         3,
			Driver:     "sqlite3",
			DSN:        filepath.Join(tmpDir, "shard3.db"),
			WriterDSN:  filepath.Join(tmpDir, "shard3.db"),
			ReaderDSNs: []string{filepath.Join(tmpDir, "shard3.db")},
		},
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{Shards: shards},
	}

	manager, err := db.NewGORMManager(cfg)
	require.NoError(t, err)
	defer manager.CloseAll()

	// Initialize schema on all shards
	for i := 1; i <= 3; i++ {
		database, err := manager.GetGORM(i)
		require.NoError(t, err)
		err = database.AutoMigrate(&model.User{})
		require.NoError(t, err)
	}

	ctx := context.Background()

	// Create users with different keys to distribute across shards
	userKeys := []int64{100, 200, 300, 400, 500}
	for _, key := range userKeys {
		database, err := manager.GetGORMByKey(key)
		require.NoError(t, err)

		user := &model.User{
			ID:    key,
			Name:  fmt.Sprintf("User %d", key),
			Email: fmt.Sprintf("user%d@example.com", key),
		}
		err = database.WithContext(ctx).Create(user).Error
		require.NoError(t, err)
	}

	// Verify that users are distributed across different shards
	shardCounts := make(map[int]int)
	for _, key := range userKeys {
		database, err := manager.GetGORMByKey(key)
		require.NoError(t, err)

		var user model.User
		err = database.WithContext(ctx).First(&user, key).Error
		require.NoError(t, err)

		// Determine which shard this key belongs to
		shardID := manager.GetShardIDByKey(key)
		shardCounts[shardID]++
	}

	// Verify that at least 2 different shards are used (probabilistic)
	assert.GreaterOrEqual(t, len(shardCounts), 2, "Users should be distributed across multiple shards")
}

// TestGORMManagerGetAllConnections tests retrieving all GORM connections
func TestGORMManagerGetAllConnections(t *testing.T) {
	tmpDir := t.TempDir()

	shards := []config.ShardConfig{
		{
			ID:         1,
			Driver:     "sqlite3",
			DSN:        filepath.Join(tmpDir, "shard1.db"),
			WriterDSN:  filepath.Join(tmpDir, "shard1.db"),
			ReaderDSNs: []string{filepath.Join(tmpDir, "shard1.db")},
		},
		{
			ID:         2,
			Driver:     "sqlite3",
			DSN:        filepath.Join(tmpDir, "shard2.db"),
			WriterDSN:  filepath.Join(tmpDir, "shard2.db"),
			ReaderDSNs: []string{filepath.Join(tmpDir, "shard2.db")},
		},
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{Shards: shards},
	}

	manager, err := db.NewGORMManager(cfg)
	require.NoError(t, err)
	defer manager.CloseAll()

	// Get all connections
	connections := manager.GetAllGORMConnections()
	assert.Len(t, connections, 2)

	// Verify shard IDs
	shardIDs := make(map[int]bool)
	for _, conn := range connections {
		shardIDs[conn.ShardID] = true
	}
	assert.True(t, shardIDs[1])
	assert.True(t, shardIDs[2])
}

// TestGORMManagerCrossShardQuery tests querying across multiple shards
func TestGORMManagerCrossShardQuery(t *testing.T) {
	tmpDir := t.TempDir()

	shards := []config.ShardConfig{
		{
			ID:         1,
			Driver:     "sqlite3",
			DSN:        filepath.Join(tmpDir, "shard1.db"),
			WriterDSN:  filepath.Join(tmpDir, "shard1.db"),
			ReaderDSNs: []string{filepath.Join(tmpDir, "shard1.db")},
		},
		{
			ID:         2,
			Driver:     "sqlite3",
			DSN:        filepath.Join(tmpDir, "shard2.db"),
			WriterDSN:  filepath.Join(tmpDir, "shard2.db"),
			ReaderDSNs: []string{filepath.Join(tmpDir, "shard2.db")},
		},
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{Shards: shards},
	}

	manager, err := db.NewGORMManager(cfg)
	require.NoError(t, err)
	defer manager.CloseAll()

	// Initialize schema on all shards
	for i := 1; i <= 2; i++ {
		database, err := manager.GetGORM(i)
		require.NoError(t, err)
		err = database.AutoMigrate(&model.User{})
		require.NoError(t, err)
	}

	ctx := context.Background()

	// Create users on specific shards
	db1, err := manager.GetGORM(1)
	require.NoError(t, err)

	user1 := &model.User{
		ID:    1,
		Name:  "User on Shard 1",
		Email: "user1@example.com",
	}
	err = db1.WithContext(ctx).Create(user1).Error
	require.NoError(t, err)

	db2, err := manager.GetGORM(2)
	require.NoError(t, err)

	user2 := &model.User{
		ID:    2,
		Name:  "User on Shard 2",
		Email: "user2@example.com",
	}
	err = db2.WithContext(ctx).Create(user2).Error
	require.NoError(t, err)

	// Query all shards
	allUsers := make([]*model.User, 0)
	connections := manager.GetAllGORMConnections()

	for _, conn := range connections {
		var users []*model.User
		err := conn.DB.WithContext(ctx).Find(&users).Error
		require.NoError(t, err)
		allUsers = append(allUsers, users...)
	}

	// Verify we got users from both shards
	assert.Len(t, allUsers, 2)
}

// TestGORMManagerShardingConsistency tests that the same key always maps to the same shard
func TestGORMManagerShardingConsistency(t *testing.T) {
	tmpDir := t.TempDir()

	shards := []config.ShardConfig{
		{
			ID:         1,
			Driver:     "sqlite3",
			DSN:        filepath.Join(tmpDir, "shard1.db"),
			WriterDSN:  filepath.Join(tmpDir, "shard1.db"),
			ReaderDSNs: []string{filepath.Join(tmpDir, "shard1.db")},
		},
		{
			ID:         2,
			Driver:     "sqlite3",
			DSN:        filepath.Join(tmpDir, "shard2.db"),
			WriterDSN:  filepath.Join(tmpDir, "shard2.db"),
			ReaderDSNs: []string{filepath.Join(tmpDir, "shard2.db")},
		},
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{Shards: shards},
	}

	manager, err := db.NewGORMManager(cfg)
	require.NoError(t, err)
	defer manager.CloseAll()

	// Test that the same key always maps to the same shard
	testKeys := []int64{100, 200, 300, 400, 500}
	shardMapping := make(map[int64]int)

	// First pass - record shard mappings
	for _, key := range testKeys {
		shardID := manager.GetShardIDByKey(key)
		shardMapping[key] = shardID
	}

	// Second pass - verify consistency
	for i := 0; i < 10; i++ {
		for _, key := range testKeys {
			shardID := manager.GetShardIDByKey(key)
			assert.Equal(t, shardMapping[key], shardID,
				"Key %d should always map to shard %d, got %d", key, shardMapping[key], shardID)
		}
	}
}

// TestGORMManagerPingAll tests pinging all shard connections
func TestGORMManagerPingAll(t *testing.T) {
	tmpDir := t.TempDir()

	shards := []config.ShardConfig{
		{
			ID:         1,
			Driver:     "sqlite3",
			DSN:        filepath.Join(tmpDir, "shard1.db"),
			WriterDSN:  filepath.Join(tmpDir, "shard1.db"),
			ReaderDSNs: []string{filepath.Join(tmpDir, "shard1.db")},
		},
		{
			ID:         2,
			Driver:     "sqlite3",
			DSN:        filepath.Join(tmpDir, "shard2.db"),
			WriterDSN:  filepath.Join(tmpDir, "shard2.db"),
			ReaderDSNs: []string{filepath.Join(tmpDir, "shard2.db")},
		},
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{Shards: shards},
	}

	manager, err := db.NewGORMManager(cfg)
	require.NoError(t, err)
	defer manager.CloseAll()

	// Ping all connections
	err = manager.PingAll()
	assert.NoError(t, err)
}

// TestGORMManagerCloseAll tests closing all shard connections
func TestGORMManagerCloseAll(t *testing.T) {
	tmpDir := t.TempDir()

	shards := []config.ShardConfig{
		{
			ID:         1,
			Driver:     "sqlite3",
			DSN:        filepath.Join(tmpDir, "shard1.db"),
			WriterDSN:  filepath.Join(tmpDir, "shard1.db"),
			ReaderDSNs: []string{filepath.Join(tmpDir, "shard1.db")},
		},
		{
			ID:         2,
			Driver:     "sqlite3",
			DSN:        filepath.Join(tmpDir, "shard2.db"),
			WriterDSN:  filepath.Join(tmpDir, "shard2.db"),
			ReaderDSNs: []string{filepath.Join(tmpDir, "shard2.db")},
		},
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{Shards: shards},
	}

	manager, err := db.NewGORMManager(cfg)
	require.NoError(t, err)

	// Verify connections work
	err = manager.PingAll()
	require.NoError(t, err)

	// Close all connections
	err = manager.CloseAll()
	assert.NoError(t, err)
}
