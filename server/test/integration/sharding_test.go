package integration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/example/go-webdb-template/internal/db"
	"github.com/example/go-webdb-template/internal/model"
	"github.com/example/go-webdb-template/test/testutil"
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

			tableName := tableSelector.GetTableName("users", tc.id)
			expectedName := "users_" + tc.expectedSuffix
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

		tableName := tableSelector.GetTableName("users", u.id)
		err = conn.DB.Table(tableName).Create(&model.User{
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

		tableName := tableSelector.GetTableName("users", u.id)
		var retrieved model.User
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
		news := &model.News{
			Title:   "Test News Title",
			Content: "Test news content",
		}

		err := conn.DB.Create(news).Error
		require.NoError(t, err)
		assert.NotZero(t, news.ID)

		// Test Read
		t.Run("Read News by ID", func(t *testing.T) {
			var retrieved model.News
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

			var retrieved model.News
			err = conn.DB.First(&retrieved, news.ID).Error
			require.NoError(t, err)
			assert.Equal(t, "Updated Title", retrieved.Title)
		})

		// Test Delete
		t.Run("Delete News", func(t *testing.T) {
			err := conn.DB.Delete(news).Error
			require.NoError(t, err)

			var retrieved model.News
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
			conn, err := groupManager.GetShardingConnectionByID(tc.userID, "users")
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
