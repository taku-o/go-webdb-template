package integration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/test/testutil"
)

// TestMasterGroupConnection tests connection to master database group
func TestMasterGroupConnection(t *testing.T) {
	// Setup test database with GroupManager
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	// Test: Get master connection
	t.Run("GetMasterConnection", func(t *testing.T) {
		conn, err := groupManager.GetMasterConnection()
		require.NoError(t, err)
		assert.NotNil(t, conn)
		assert.NotNil(t, conn.DB)
	})

	// Test: Ping master connection
	t.Run("PingMasterConnection", func(t *testing.T) {
		conn, err := groupManager.GetMasterConnection()
		require.NoError(t, err)

		sqlDB, err := conn.DB.DB()
		require.NoError(t, err)

		err = sqlDB.Ping()
		assert.NoError(t, err)
	})
}

// TestShardingGroupConnection tests connection to sharding database group
func TestShardingGroupConnection(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	// Test: Get sharding connection for each table range
	testCases := []struct {
		name        string
		tableNumber int
		expectDBID  int
	}{
		{"Table 0 in DB1", 0, 1},
		{"Table 7 in DB1", 7, 1},
		{"Table 8 in DB2", 8, 2},
		{"Table 15 in DB2", 15, 2},
		{"Table 16 in DB3", 16, 3},
		{"Table 23 in DB3", 23, 3},
		{"Table 24 in DB4", 24, 4},
		{"Table 31 in DB4", 31, 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conn, err := groupManager.GetShardingConnection(tc.tableNumber)
			require.NoError(t, err)
			assert.NotNil(t, conn)
			assert.Equal(t, tc.expectDBID, conn.ShardID)
		})
	}
}

// TestTableSelectionLogic tests the table selection logic (ID % 32)
func TestTableSelectionLogic(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	tableSelector := db.NewTableSelector(32, 8)

	testCases := []struct {
		id              int64
		expectedTable   int
		expectedDBID    int
		expectedSuffix  string
	}{
		{0, 0, 1, "000"},
		{1, 1, 1, "001"},
		{7, 7, 1, "007"},
		{8, 8, 2, "008"},
		{31, 31, 4, "031"},
		{32, 0, 1, "000"},   // wraps around
		{33, 1, 1, "001"},   // wraps around
		{100, 4, 1, "004"},  // 100 % 32 = 4
		{1000, 8, 2, "008"}, // 1000 % 32 = 8
	}

	for _, tc := range testCases {
		t.Run("ID_"+fmt.Sprintf("%d", tc.id), func(t *testing.T) {
			tableNumber := tableSelector.GetTableNumber(tc.id)
			assert.Equal(t, tc.expectedTable, tableNumber)

			dbID := tableSelector.GetDBID(tableNumber)
			assert.Equal(t, tc.expectedDBID, dbID)

			tableName := tableSelector.GetTableName("dm_users", tc.id)
			expectedName := "dm_users_" + tc.expectedSuffix
			assert.Equal(t, expectedName, tableName)
		})
	}
}

// TestCrossTableQueryUsers tests cross-table query for users across all shards
func TestCrossTableQueryUsers(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	ctx := context.Background()

	// Create users with specific IDs that will be distributed across different tables
	testUsers := []struct {
		id    int64
		name  string
		email string
	}{
		{1, "User in Table 1", "user1@test.com"},
		{8, "User in Table 8", "user8@test.com"},
		{16, "User in Table 16", "user16@test.com"},
		{24, "User in Table 24", "user24@test.com"},
	}

	tableSelector := db.NewTableSelector(32, 8)

	// Insert users directly to specific tables
	for _, u := range testUsers {
		tableNumber := tableSelector.GetTableNumber(u.id)
		conn, err := groupManager.GetShardingConnection(tableNumber)
		require.NoError(t, err)

		tableName := tableSelector.GetTableName("dm_users", u.id)
		err = conn.DB.Table(tableName).Create(&model.DmUser{
			ID:    u.id,
			Name:  u.name,
			Email: u.email,
		}).Error
		require.NoError(t, err)
	}

	// Verify users can be retrieved from their respective tables
	for _, u := range testUsers {
		tableNumber := tableSelector.GetTableNumber(u.id)
		conn, err := groupManager.GetShardingConnection(tableNumber)
		require.NoError(t, err)

		tableName := tableSelector.GetTableName("dm_users", u.id)
		var retrieved model.DmUser
		err = conn.DB.Table(tableName).Where("id = ?", u.id).First(&retrieved).Error
		require.NoError(t, err, "Failed to retrieve user %d from %s", u.id, tableName)
		assert.Equal(t, u.name, retrieved.Name)
		assert.Equal(t, u.email, retrieved.Email)
	}

	_ = ctx // context for future use
}

// TestMasterGroupNewsTable tests CRUD operations on the news table in master group
func TestMasterGroupNewsTable(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	conn, err := groupManager.GetMasterConnection()
	require.NoError(t, err)

	// Test Create
	t.Run("Create News", func(t *testing.T) {
		news := &model.DmNews{
			Title:   "Test News Title",
			Content: "Test news content",
		}

		err := conn.DB.Create(news).Error
		require.NoError(t, err)
		assert.NotZero(t, news.ID)

		// Test Read
		t.Run("Read News by ID", func(t *testing.T) {
			var retrieved model.DmNews
			err := conn.DB.First(&retrieved, news.ID).Error
			require.NoError(t, err)
			assert.Equal(t, news.Title, retrieved.Title)
			assert.Equal(t, news.Content, retrieved.Content)
		})

		// Test Update
		t.Run("Update News", func(t *testing.T) {
			news.Title = "Updated Title"
			err := conn.DB.Save(news).Error
			require.NoError(t, err)

			var retrieved model.DmNews
			err = conn.DB.First(&retrieved, news.ID).Error
			require.NoError(t, err)
			assert.Equal(t, "Updated Title", retrieved.Title)
		})

		// Test Delete
		t.Run("Delete News", func(t *testing.T) {
			err := conn.DB.Delete(news).Error
			require.NoError(t, err)

			var retrieved model.DmNews
			err = conn.DB.First(&retrieved, news.ID).Error
			assert.Error(t, err) // Should not find deleted news
		})
	})
}

// TestShardingConnectionByID tests getting connection using entity ID
func TestShardingConnectionByID(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	testCases := []struct {
		userID      int64
		expectedDB  int
	}{
		{1, 1},    // 1 % 32 = 1, DB1
		{8, 2},    // 8 % 32 = 8, DB2
		{16, 3},   // 16 % 32 = 16, DB3
		{24, 4},   // 24 % 32 = 24, DB4
		{32, 1},   // 32 % 32 = 0, DB1
		{100, 1},  // 100 % 32 = 4, DB1
	}

	for _, tc := range testCases {
		t.Run("UserID_"+fmt.Sprintf("%d", tc.userID), func(t *testing.T) {
			conn, err := groupManager.GetShardingConnectionByID(tc.userID, "dm_users")
			require.NoError(t, err)
			assert.Equal(t, tc.expectedDB, conn.ShardID)
		})
	}
}

// TestGetAllShardingConnections tests retrieval of all sharding connections
func TestGetAllShardingConnections(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	connections := groupManager.GetAllShardingConnections()
	assert.Len(t, connections, 4) // 4 databases

	// Verify all connections are valid
	dbIDs := make(map[int]bool)
	for _, conn := range connections {
		assert.NotNil(t, conn.DB)
		dbIDs[conn.ShardID] = true
	}

	// Verify we have connections for all 4 DBs
	assert.True(t, dbIDs[1])
	assert.True(t, dbIDs[2])
	assert.True(t, dbIDs[3])
	assert.True(t, dbIDs[4])
}

// =============================================================================
// 8シャーディング構成テスト
// =============================================================================

// TestShardingGroupConnection8Sharding tests the 8-sharding configuration with connection sharing
func TestShardingGroupConnection8Sharding(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager8Sharding(t)
	defer testutil.CleanupTestGroupManager(groupManager)

	// Test: Get sharding connection for each table range
	// With 8 sharding entries and connection sharing:
	// - Entries 1,2 share same connection (ShardID=1)
	// - Entries 3,4 share same connection (ShardID=3)
	// - Entries 5,6 share same connection (ShardID=5)
	// - Entries 7,8 share same connection (ShardID=7)
	testCases := []struct {
		name        string
		tableNumber int
		expectDBID  int // ShardID of the shared connection
	}{
		// DB1 (entries 1,2)
		{"Table 0 in Entry1", 0, 1},
		{"Table 3 in Entry1", 3, 1},
		{"Table 4 in Entry2 (shared)", 4, 1},
		{"Table 7 in Entry2 (shared)", 7, 1},
		// DB2 (entries 3,4)
		{"Table 8 in Entry3", 8, 3},
		{"Table 11 in Entry3", 11, 3},
		{"Table 12 in Entry4 (shared)", 12, 3},
		{"Table 15 in Entry4 (shared)", 15, 3},
		// DB3 (entries 5,6)
		{"Table 16 in Entry5", 16, 5},
		{"Table 19 in Entry5", 19, 5},
		{"Table 20 in Entry6 (shared)", 20, 5},
		{"Table 23 in Entry6 (shared)", 23, 5},
		// DB4 (entries 7,8)
		{"Table 24 in Entry7", 24, 7},
		{"Table 27 in Entry7", 27, 7},
		{"Table 28 in Entry8 (shared)", 28, 7},
		{"Table 31 in Entry8 (shared)", 31, 7},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conn, err := groupManager.GetShardingConnection(tc.tableNumber)
			require.NoError(t, err)
			assert.NotNil(t, conn)
			assert.Equal(t, tc.expectDBID, conn.ShardID, "Table %d should return connection with ShardID %d", tc.tableNumber, tc.expectDBID)
		})
	}
}

// TestConnectionSharing8Sharding verifies that connections are shared correctly
func TestConnectionSharing8Sharding(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager8Sharding(t)
	defer testutil.CleanupTestGroupManager(groupManager)

	// Get all unique connections
	connections := groupManager.GetAllShardingConnections()
	assert.Len(t, connections, 4, "Should return 4 unique connections (not 8)")

	// Verify shared connections return same object
	// Table 0 (entry 1) and table 4 (entry 2) should share the same connection
	conn1, err := groupManager.GetShardingConnection(0)
	require.NoError(t, err)
	conn2, err := groupManager.GetShardingConnection(4)
	require.NoError(t, err)
	assert.Same(t, conn1, conn2, "Tables 0 and 4 should share same connection")

	// Table 8 (entry 3) and table 12 (entry 4) should share the same connection
	conn3, err := groupManager.GetShardingConnection(8)
	require.NoError(t, err)
	conn4, err := groupManager.GetShardingConnection(12)
	require.NoError(t, err)
	assert.Same(t, conn3, conn4, "Tables 8 and 12 should share same connection")

	// Different databases should have different connections
	assert.NotSame(t, conn1, conn3, "Different databases should have different connections")
}

// TestCrossTableQuery8Sharding tests CRUD operations across all shards in 8-sharding config
func TestCrossTableQuery8Sharding(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager8Sharding(t)
	defer testutil.CleanupTestGroupManager(groupManager)

	tableSelector := db.NewTableSelector(32, 8)

	// Create users in different shards
	testUsers := []struct {
		id    int64
		name  string
		email string
	}{
		{1, "User in Table 1", "user1@test.com"},    // Table 1, Entry 1, DB1
		{5, "User in Table 5", "user5@test.com"},    // Table 5, Entry 2, DB1
		{9, "User in Table 9", "user9@test.com"},    // Table 9, Entry 3, DB2
		{13, "User in Table 13", "user13@test.com"}, // Table 13, Entry 4, DB2
		{17, "User in Table 17", "user17@test.com"}, // Table 17, Entry 5, DB3
		{21, "User in Table 21", "user21@test.com"}, // Table 21, Entry 6, DB3
		{25, "User in Table 25", "user25@test.com"}, // Table 25, Entry 7, DB4
		{29, "User in Table 29", "user29@test.com"}, // Table 29, Entry 8, DB4
	}

	// Insert users
	for _, u := range testUsers {
		tableNumber := tableSelector.GetTableNumber(u.id)
		conn, err := groupManager.GetShardingConnection(tableNumber)
		require.NoError(t, err)

		tableName := tableSelector.GetTableName("dm_users", u.id)
		err = conn.DB.Table(tableName).Create(&model.DmUser{
			ID:    u.id,
			Name:  u.name,
			Email: u.email,
		}).Error
		require.NoError(t, err, "Failed to create user %d in %s", u.id, tableName)
	}

	// Verify users can be retrieved
	for _, u := range testUsers {
		tableNumber := tableSelector.GetTableNumber(u.id)
		conn, err := groupManager.GetShardingConnection(tableNumber)
		require.NoError(t, err)

		tableName := tableSelector.GetTableName("dm_users", u.id)
		var retrieved model.DmUser
		err = conn.DB.Table(tableName).Where("id = ?", u.id).First(&retrieved).Error
		require.NoError(t, err, "Failed to retrieve user %d from %s", u.id, tableName)
		assert.Equal(t, u.name, retrieved.Name)
		assert.Equal(t, u.email, retrieved.Email)
	}
}
