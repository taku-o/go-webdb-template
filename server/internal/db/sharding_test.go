package db_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/model"
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
		err = database.AutoMigrate(&model.DmUser{})
		require.NoError(t, err)
	}

	ctx := context.Background()

	// Create users with different keys to distribute across shards
	userKeys := []int64{100, 200, 300, 400, 500}
	for _, key := range userKeys {
		database, err := manager.GetGORMByKey(key)
		require.NoError(t, err)

		user := &model.DmUser{
			ID:    fmt.Sprintf("550e8400e29b41d4a7164466554400%02x", key%256),
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

		var user model.DmUser
		userID := fmt.Sprintf("550e8400e29b41d4a7164466554400%02x", key%256)
		err = database.WithContext(ctx).First(&user, "id = ?", userID).Error
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
		err = database.AutoMigrate(&model.DmUser{})
		require.NoError(t, err)
	}

	ctx := context.Background()

	// Create users on specific shards
	db1, err := manager.GetGORM(1)
	require.NoError(t, err)

	user1 := &model.DmUser{
		ID:    "550e8400e29b41d4a716446655440001",
		Name:  "User on Shard 1",
		Email: "user1@example.com",
	}
	err = db1.WithContext(ctx).Create(user1).Error
	require.NoError(t, err)

	db2, err := manager.GetGORM(2)
	require.NoError(t, err)

	user2 := &model.DmUser{
		ID:    "550e8400e29b41d4a716446655440002",
		Name:  "User on Shard 2",
		Email: "user2@example.com",
	}
	err = db2.WithContext(ctx).Create(user2).Error
	require.NoError(t, err)

	// Query all shards
	allUsers := make([]*model.DmUser, 0)
	connections := manager.GetAllGORMConnections()

	for _, conn := range connections {
		var users []*model.DmUser
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

// =============================================================================
// タスク3.1, 3.2, 3.3: TableSelectorのテスト
// =============================================================================

// TestNewTableSelector tests TableSelector creation
func TestNewTableSelector(t *testing.T) {
	tests := []struct {
		name            string
		tableCount      int
		tablesPerDB     int
		wantTableCount  int
		wantTablesPerDB int
	}{
		{
			name:            "default values when zero",
			tableCount:      0,
			tablesPerDB:     0,
			wantTableCount:  32,
			wantTablesPerDB: 8,
		},
		{
			name:            "custom values",
			tableCount:      64,
			tablesPerDB:     16,
			wantTableCount:  64,
			wantTablesPerDB: 16,
		},
		{
			name:            "negative values use defaults",
			tableCount:      -1,
			tablesPerDB:     -1,
			wantTableCount:  32,
			wantTablesPerDB: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector := db.NewTableSelector(tt.tableCount, tt.tablesPerDB)

			assert.Equal(t, tt.wantTableCount, selector.GetTableCount())
		})
	}
}

// TestTableSelector_GetTableNumber tests GetTableNumber method
func TestTableSelector_GetTableNumber(t *testing.T) {
	selector := db.NewTableSelector(32, 8)

	tests := []struct {
		id              int64
		wantTableNumber int
	}{
		{id: 0, wantTableNumber: 0},
		{id: 1, wantTableNumber: 1},
		{id: 31, wantTableNumber: 31},
		{id: 32, wantTableNumber: 0},
		{id: 33, wantTableNumber: 1},
		{id: 100, wantTableNumber: 4},   // 100 % 32 = 4
		{id: 1000, wantTableNumber: 8},  // 1000 % 32 = 8
		{id: 10000, wantTableNumber: 16}, // 10000 % 32 = 16
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("id=%d", tt.id), func(t *testing.T) {
			tableNumber := selector.GetTableNumber(tt.id)
			assert.Equal(t, tt.wantTableNumber, tableNumber)
		})
	}
}

// TestTableSelector_GetTableName tests GetTableName method
func TestTableSelector_GetTableName(t *testing.T) {
	selector := db.NewTableSelector(32, 8)

	tests := []struct {
		baseName      string
		id            int64
		wantTableName string
	}{
		{baseName: "dm_users", id: 0, wantTableName: "dm_users_000"},
		{baseName: "dm_users", id: 1, wantTableName: "dm_users_001"},
		{baseName: "dm_users", id: 7, wantTableName: "dm_users_007"},
		{baseName: "dm_users", id: 31, wantTableName: "dm_users_031"},
		{baseName: "dm_users", id: 32, wantTableName: "dm_users_000"},
		{baseName: "dm_posts", id: 15, wantTableName: "dm_posts_015"},
		{baseName: "dm_posts", id: 100, wantTableName: "dm_posts_004"}, // 100 % 32 = 4
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_id=%d", tt.baseName, tt.id), func(t *testing.T) {
			tableName := selector.GetTableName(tt.baseName, tt.id)
			assert.Equal(t, tt.wantTableName, tableName)
		})
	}
}

// TestTableSelector_GetDBID tests GetDBID method
func TestTableSelector_GetDBID(t *testing.T) {
	selector := db.NewTableSelector(32, 8)

	tests := []struct {
		tableNumber int
		wantDBID    int
	}{
		// DB1: テーブル番号 0-7
		{tableNumber: 0, wantDBID: 1},
		{tableNumber: 7, wantDBID: 1},
		// DB2: テーブル番号 8-15
		{tableNumber: 8, wantDBID: 2},
		{tableNumber: 15, wantDBID: 2},
		// DB3: テーブル番号 16-23
		{tableNumber: 16, wantDBID: 3},
		{tableNumber: 23, wantDBID: 3},
		// DB4: テーブル番号 24-31
		{tableNumber: 24, wantDBID: 4},
		{tableNumber: 31, wantDBID: 4},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("tableNumber=%d", tt.tableNumber), func(t *testing.T) {
			dbID := selector.GetDBID(tt.tableNumber)
			assert.Equal(t, tt.wantDBID, dbID)
		})
	}
}

// TestTableSelector_GetTableCount tests GetTableCount method
func TestTableSelector_GetTableCount(t *testing.T) {
	selector := db.NewTableSelector(32, 8)
	assert.Equal(t, 32, selector.GetTableCount())

	selector64 := db.NewTableSelector(64, 16)
	assert.Equal(t, 64, selector64.GetTableCount())
}

// TestGetShardingTableName tests the utility function
func TestGetShardingTableName(t *testing.T) {
	tests := []struct {
		baseName      string
		id            int64
		wantTableName string
	}{
		{baseName: "dm_users", id: 0, wantTableName: "dm_users_000"},
		{baseName: "dm_users", id: 31, wantTableName: "dm_users_031"},
		{baseName: "dm_users", id: 32, wantTableName: "dm_users_000"},
		{baseName: "dm_posts", id: 100, wantTableName: "dm_posts_004"}, // 100 % 32 = 4
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_id=%d", tt.baseName, tt.id), func(t *testing.T) {
			tableName := db.GetShardingTableName(tt.baseName, tt.id)
			assert.Equal(t, tt.wantTableName, tableName)
		})
	}
}

// TestGetShardingTableNumber tests the utility function
func TestGetShardingTableNumber(t *testing.T) {
	tests := []struct {
		id              int64
		wantTableNumber int
	}{
		{id: 0, wantTableNumber: 0},
		{id: 1, wantTableNumber: 1},
		{id: 31, wantTableNumber: 31},
		{id: 32, wantTableNumber: 0},
		{id: 100, wantTableNumber: 4},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("id=%d", tt.id), func(t *testing.T) {
			tableNumber := db.GetShardingTableNumber(tt.id)
			assert.Equal(t, tt.wantTableNumber, tableNumber)
		})
	}
}

// TestGetShardingDBID tests the utility function
func TestGetShardingDBID(t *testing.T) {
	tests := []struct {
		tableNumber int
		wantDBID    int
	}{
		{tableNumber: 0, wantDBID: 1},
		{tableNumber: 7, wantDBID: 1},
		{tableNumber: 8, wantDBID: 2},
		{tableNumber: 15, wantDBID: 2},
		{tableNumber: 16, wantDBID: 3},
		{tableNumber: 23, wantDBID: 3},
		{tableNumber: 24, wantDBID: 4},
		{tableNumber: 31, wantDBID: 4},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("tableNumber=%d", tt.tableNumber), func(t *testing.T) {
			dbID := db.GetShardingDBID(tt.tableNumber)
			assert.Equal(t, tt.wantDBID, dbID)
		})
	}
}

// TestValidateTableName tests the table name validation function
func TestValidateTableName(t *testing.T) {
	allowedBaseNames := []string{"dm_users", "dm_posts"}

	tests := []struct {
		tableName string
		want      bool
	}{
		// Valid names
		{tableName: "dm_users_000", want: true},
		{tableName: "dm_users_031", want: true},
		{tableName: "dm_posts_000", want: true},
		{tableName: "dm_posts_031", want: true},
		{tableName: "dm_users_015", want: true},
		// Invalid names
		{tableName: "dm_users_032", want: false},  // Out of range
		{tableName: "dm_users_100", want: false},  // Out of range
		{tableName: "dm_news", want: false},       // Not in allowed list
		{tableName: "dm_users", want: false},      // No suffix
		{tableName: "dm_users_00", want: false},   // Wrong suffix format
		{tableName: "dm_users_0000", want: false}, // Wrong suffix format
		{tableName: "other_000", want: false},  // Not in allowed list
		{tableName: "'; DROP TABLE dm_users; --", want: false}, // SQL injection attempt
	}

	for _, tt := range tests {
		t.Run(tt.tableName, func(t *testing.T) {
			valid := db.ValidateTableName(tt.tableName, allowedBaseNames)
			assert.Equal(t, tt.want, valid)
		})
	}
}


// =============================================================================
// シャーディング規則のテスト（dm_users.id と dm_posts.user_id の関係）
// =============================================================================

// TestShardingRuleConsistency tests that dm_users.id and dm_posts.user_id
// result in the same table number when they have the same value.
// This is a critical sharding rule: a user and their posts must be in the
// same numbered tables (e.g., dm_users_009 and dm_posts_009).
func TestShardingRuleConsistency(t *testing.T) {
	selector := db.NewTableSelector(32, 8)

	tests := []struct {
		name   string
		userID int64
	}{
		{name: "small_id", userID: 5},
		{name: "boundary_id_31", userID: 31},
		{name: "boundary_id_32", userID: 32},
		{name: "large_id", userID: 12345},
		{name: "very_large_id", userID: 1234567890123456789},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// dm_users uses id as sharding key
			dmUsersTableNumber := selector.GetTableNumber(tt.userID)
			dmUsersTableName := selector.GetTableName("dm_users", tt.userID)

			// dm_posts uses user_id as sharding key
			// When dm_posts.user_id == dm_users.id, they should be in the same table number
			dmPostsTableNumber := selector.GetTableNumber(tt.userID)
			dmPostsTableName := selector.GetTableName("dm_posts", tt.userID)

			// Verify both have the same table number
			assert.Equal(t, dmUsersTableNumber, dmPostsTableNumber,
				"dm_users with id=%d and dm_posts with user_id=%d should have same table number",
				tt.userID, tt.userID)

			// Verify table names have the same suffix number
			expectedSuffix := fmt.Sprintf("_%03d", dmUsersTableNumber)
			assert.Contains(t, dmUsersTableName, expectedSuffix)
			assert.Contains(t, dmPostsTableName, expectedSuffix)

			// Log for clarity
			t.Logf("userID=%d -> dm_users=%s, dm_posts=%s (table number: %d)",
				tt.userID, dmUsersTableName, dmPostsTableName, dmUsersTableNumber)
		})
	}
}

// TestShardingRuleWithMultiplePosts tests that all posts from the same user
// go to the same table as the user.
func TestShardingRuleWithMultiplePosts(t *testing.T) {
	selector := db.NewTableSelector(32, 8)

	// Simulate a user with ID 12345
	userID := int64(12345)
	userTableNumber := selector.GetTableNumber(userID)
	userTableName := selector.GetTableName("dm_users", userID)

	// Simulate multiple posts from this user (each post has a different ID
	// but all share the same user_id)
	postIDs := []int64{100001, 100002, 100003, 100004, 100005}

	for _, postID := range postIDs {
		// Posts are sharded by user_id, not by post ID
		postTableNumber := selector.GetTableNumber(userID)
		postTableName := selector.GetTableName("dm_posts", userID)

		assert.Equal(t, userTableNumber, postTableNumber,
			"Post with id=%d should be in same table as user with id=%d",
			postID, userID)

		assert.Equal(t, fmt.Sprintf("dm_posts_%03d", userTableNumber), postTableName)
	}

	t.Logf("User id=%d is in %s, all posts are in dm_posts_%03d",
		userID, userTableName, userTableNumber)
}

// =============================================================================
// UUIDv7ベースのシャーディングキー計算テスト
// =============================================================================

// TestTableSelector_GetTableNumberFromUUID tests GetTableNumberFromUUID method
func TestTableSelector_GetTableNumberFromUUID(t *testing.T) {
	selector := db.NewTableSelector(32, 8)

	tests := []struct {
		name            string
		uuid            string
		wantTableNumber int
		wantError       bool
	}{
		// 正常系：後ろ2文字が00の場合（テーブル番号: 0）
		{name: "suffix_00", uuid: "550e8400e29b41d4a716446655440000", wantTableNumber: 0, wantError: false},
		// 正常系：後ろ2文字が0fの場合（テーブル番号: 15）
		{name: "suffix_0f", uuid: "550e8400e29b41d4a71644665544000f", wantTableNumber: 15, wantError: false},
		// 正常系：後ろ2文字が1fの場合（テーブル番号: 31）
		{name: "suffix_1f", uuid: "550e8400e29b41d4a71644665544001f", wantTableNumber: 31, wantError: false},
		// 正常系：後ろ2文字が20の場合（32 % 32 = 0）
		{name: "suffix_20", uuid: "550e8400e29b41d4a716446655440020", wantTableNumber: 0, wantError: false},
		// 正常系：後ろ2文字がffの場合（255 % 32 = 31）
		{name: "suffix_ff", uuid: "550e8400e29b41d4a7164466554400ff", wantTableNumber: 31, wantError: false},
		// 正常系：後ろ2文字が21の場合（33 % 32 = 1）
		{name: "suffix_21", uuid: "550e8400e29b41d4a716446655440021", wantTableNumber: 1, wantError: false},
		// 正常系：大文字のUUID
		{name: "uppercase_suffix", uuid: "550E8400E29B41D4A716446655440000", wantTableNumber: 0, wantError: false},
		// エラー系：短すぎるUUID（1文字）
		{name: "too_short", uuid: "0", wantTableNumber: 0, wantError: true},
		// エラー系：空文字列
		{name: "empty", uuid: "", wantTableNumber: 0, wantError: true},
		// エラー系：無効な16進数
		{name: "invalid_hex", uuid: "550e8400e29b41d4a7164466554400gg", wantTableNumber: 0, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tableNumber, err := selector.GetTableNumberFromUUID(tt.uuid)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantTableNumber, tableNumber)
			}
		})
	}
}

// TestTableSelector_GetTableNameFromUUID tests GetTableNameFromUUID method
func TestTableSelector_GetTableNameFromUUID(t *testing.T) {
	selector := db.NewTableSelector(32, 8)

	tests := []struct {
		name          string
		baseName      string
		uuid          string
		wantTableName string
		wantError     bool
	}{
		// 正常系
		{name: "dm_users_000", baseName: "dm_users", uuid: "550e8400e29b41d4a716446655440000", wantTableName: "dm_users_000", wantError: false},
		{name: "dm_users_015", baseName: "dm_users", uuid: "550e8400e29b41d4a71644665544000f", wantTableName: "dm_users_015", wantError: false},
		{name: "dm_users_031", baseName: "dm_users", uuid: "550e8400e29b41d4a71644665544001f", wantTableName: "dm_users_031", wantError: false},
		{name: "dm_posts_000", baseName: "dm_posts", uuid: "550e8400e29b41d4a716446655440020", wantTableName: "dm_posts_000", wantError: false},
		// エラー系
		{name: "invalid_uuid", baseName: "dm_users", uuid: "invalid", wantTableName: "", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tableName, err := selector.GetTableNameFromUUID(tt.baseName, tt.uuid)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantTableName, tableName)
			}
		})
	}
}

// TestShardingRuleConsistencyWithUUID tests that UUID-based sharding produces
// consistent results for the same UUID
func TestShardingRuleConsistencyWithUUID(t *testing.T) {
	selector := db.NewTableSelector(32, 8)

	// 同じUUIDは常に同じテーブル番号を返す
	uuid := "550e8400e29b41d4a716446655440012"

	tableNumber1, err := selector.GetTableNumberFromUUID(uuid)
	require.NoError(t, err)

	tableNumber2, err := selector.GetTableNumberFromUUID(uuid)
	require.NoError(t, err)

	tableNumber3, err := selector.GetTableNumberFromUUID(uuid)
	require.NoError(t, err)

	assert.Equal(t, tableNumber1, tableNumber2)
	assert.Equal(t, tableNumber2, tableNumber3)
}

// TestShardingRuleDistributionWithUUID tests that UUIDs are distributed across tables
func TestShardingRuleDistributionWithUUID(t *testing.T) {
	selector := db.NewTableSelector(32, 8)

	// 異なるUUIDが異なるテーブルに分散されることを確認
	uuids := []string{
		"550e8400e29b41d4a716446655440000",
		"550e8400e29b41d4a716446655440001",
		"550e8400e29b41d4a716446655440010",
		"550e8400e29b41d4a71644665544001f",
		"550e8400e29b41d4a716446655440020",
		"550e8400e29b41d4a7164466554400ff",
	}

	tableNumbers := make(map[int]int)
	for _, uuid := range uuids {
		tableNumber, err := selector.GetTableNumberFromUUID(uuid)
		require.NoError(t, err)
		tableNumbers[tableNumber]++
	}

	// 少なくとも2つの異なるテーブルに分散されることを確認
	assert.GreaterOrEqual(t, len(tableNumbers), 2, "UUIDs should be distributed across multiple tables")
}

