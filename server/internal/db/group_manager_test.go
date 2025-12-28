package db_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/db"
)

// =============================================================================
// タスク4.1, 4.2, 4.3, 4.4: GroupManager, MasterManager, ShardingManagerのテスト
// =============================================================================

// TestNewMasterManager tests MasterManager creation
func TestNewMasterManager(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{
					{
						ID:         1,
						Driver:     "sqlite3",
						DSN:        filepath.Join(tmpDir, "master.db"),
						WriterDSN:  filepath.Join(tmpDir, "master.db"),
						ReaderDSNs: []string{filepath.Join(tmpDir, "master.db")},
					},
				},
			},
		},
		Logging: config.LoggingConfig{
			SQLLogEnabled:   false,
			SQLLogOutputDir: tmpDir,
		},
	}

	manager, err := db.NewMasterManager(cfg)
	require.NoError(t, err)
	require.NotNil(t, manager)
	defer manager.CloseAll()

	// 接続を取得できることを確認
	conn, err := manager.GetConnection()
	require.NoError(t, err)
	require.NotNil(t, conn)

	// Pingが成功することを確認
	err = manager.PingAll()
	assert.NoError(t, err)
}

// TestNewMasterManager_NoConfig tests MasterManager with no config
func TestNewMasterManager_NoConfig(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{}, // 空の設定
			},
		},
	}

	manager, err := db.NewMasterManager(cfg)
	require.Error(t, err)
	assert.Nil(t, manager)
	assert.Contains(t, err.Error(), "master group configuration not found")
}

// TestNewShardingManager tests ShardingManager creation
func TestNewShardingManager(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Sharding: config.ShardingGroupConfig{
					Databases: []config.ShardConfig{
						{
							ID:         1,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding1.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding1.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding1.db")},
							TableRange: [2]int{0, 7},
						},
						{
							ID:         2,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding2.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding2.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding2.db")},
							TableRange: [2]int{8, 15},
						},
						{
							ID:         3,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding3.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding3.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding3.db")},
							TableRange: [2]int{16, 23},
						},
						{
							ID:         4,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding4.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding4.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding4.db")},
							TableRange: [2]int{24, 31},
						},
					},
					Tables: []config.ShardingTableConfig{
						{Name: "users", SuffixCount: 32},
						{Name: "posts", SuffixCount: 32},
					},
				},
			},
		},
		Logging: config.LoggingConfig{
			SQLLogEnabled:   false,
			SQLLogOutputDir: tmpDir,
		},
	}

	manager, err := db.NewShardingManager(cfg)
	require.NoError(t, err)
	require.NotNil(t, manager)
	defer manager.CloseAll()

	// Pingが成功することを確認
	err = manager.PingAll()
	assert.NoError(t, err)
}

// TestShardingManager_GetConnectionByTableNumber tests connection by table number
func TestShardingManager_GetConnectionByTableNumber(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Sharding: config.ShardingGroupConfig{
					Databases: []config.ShardConfig{
						{
							ID:         1,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding1.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding1.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding1.db")},
							TableRange: [2]int{0, 7},
						},
						{
							ID:         2,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding2.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding2.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding2.db")},
							TableRange: [2]int{8, 15},
						},
						{
							ID:         3,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding3.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding3.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding3.db")},
							TableRange: [2]int{16, 23},
						},
						{
							ID:         4,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding4.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding4.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding4.db")},
							TableRange: [2]int{24, 31},
						},
					},
				},
			},
		},
		Logging: config.LoggingConfig{
			SQLLogEnabled:   false,
			SQLLogOutputDir: tmpDir,
		},
	}

	manager, err := db.NewShardingManager(cfg)
	require.NoError(t, err)
	defer manager.CloseAll()

	tests := []struct {
		tableNumber  int
		wantDBID     int
		expectError  bool
		errorMessage string
	}{
		// DB1: テーブル番号 0-7
		{tableNumber: 0, wantDBID: 1, expectError: false},
		{tableNumber: 7, wantDBID: 1, expectError: false},
		// DB2: テーブル番号 8-15
		{tableNumber: 8, wantDBID: 2, expectError: false},
		{tableNumber: 15, wantDBID: 2, expectError: false},
		// DB3: テーブル番号 16-23
		{tableNumber: 16, wantDBID: 3, expectError: false},
		{tableNumber: 23, wantDBID: 3, expectError: false},
		// DB4: テーブル番号 24-31
		{tableNumber: 24, wantDBID: 4, expectError: false},
		{tableNumber: 31, wantDBID: 4, expectError: false},
		// エラーケース
		{tableNumber: -1, expectError: true, errorMessage: "invalid table number"},
		{tableNumber: 32, expectError: true, errorMessage: "invalid table number"},
	}

	for _, tt := range tests {
		t.Run("tableNumber="+string(rune(tt.tableNumber+'0')), func(t *testing.T) {
			conn, err := manager.GetConnectionByTableNumber(tt.tableNumber)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, conn)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, conn)
				assert.Equal(t, tt.wantDBID, conn.ShardID)
			}
		})
	}
}

// TestShardingManager_GetAllConnections tests getting all connections
func TestShardingManager_GetAllConnections(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Sharding: config.ShardingGroupConfig{
					Databases: []config.ShardConfig{
						{
							ID:         1,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding1.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding1.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding1.db")},
							TableRange: [2]int{0, 7},
						},
						{
							ID:         2,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding2.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding2.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding2.db")},
							TableRange: [2]int{8, 15},
						},
					},
				},
			},
		},
		Logging: config.LoggingConfig{
			SQLLogEnabled:   false,
			SQLLogOutputDir: tmpDir,
		},
	}

	manager, err := db.NewShardingManager(cfg)
	require.NoError(t, err)
	defer manager.CloseAll()

	connections := manager.GetAllConnections()
	assert.Len(t, connections, 2)

	// 全接続のShardIDを確認
	shardIDs := make(map[int]bool)
	for _, conn := range connections {
		shardIDs[conn.ShardID] = true
	}
	assert.True(t, shardIDs[1])
	assert.True(t, shardIDs[2])
}

// TestNewGroupManager tests GroupManager creation
func TestNewGroupManager(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{
					{
						ID:         1,
						Driver:     "sqlite3",
						DSN:        filepath.Join(tmpDir, "master.db"),
						WriterDSN:  filepath.Join(tmpDir, "master.db"),
						ReaderDSNs: []string{filepath.Join(tmpDir, "master.db")},
					},
				},
				Sharding: config.ShardingGroupConfig{
					Databases: []config.ShardConfig{
						{
							ID:         1,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding1.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding1.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding1.db")},
							TableRange: [2]int{0, 7},
						},
						{
							ID:         2,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding2.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding2.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding2.db")},
							TableRange: [2]int{8, 15},
						},
						{
							ID:         3,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding3.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding3.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding3.db")},
							TableRange: [2]int{16, 23},
						},
						{
							ID:         4,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding4.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding4.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding4.db")},
							TableRange: [2]int{24, 31},
						},
					},
					Tables: []config.ShardingTableConfig{
						{Name: "users", SuffixCount: 32},
						{Name: "posts", SuffixCount: 32},
					},
				},
			},
		},
		Logging: config.LoggingConfig{
			SQLLogEnabled:   false,
			SQLLogOutputDir: tmpDir,
		},
	}

	manager, err := db.NewGroupManager(cfg)
	require.NoError(t, err)
	require.NotNil(t, manager)
	defer manager.CloseAll()

	// Master接続を取得できることを確認
	masterConn, err := manager.GetMasterConnection()
	require.NoError(t, err)
	require.NotNil(t, masterConn)

	// Sharding接続をテーブル番号で取得できることを確認
	shardingConn, err := manager.GetShardingConnection(0)
	require.NoError(t, err)
	require.NotNil(t, shardingConn)

	// PingAllが成功することを確認
	err = manager.PingAll()
	assert.NoError(t, err)
}

// TestGroupManager_GetShardingConnectionByID tests connection by user ID
func TestGroupManager_GetShardingConnectionByID(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{
					{
						ID:         1,
						Driver:     "sqlite3",
						DSN:        filepath.Join(tmpDir, "master.db"),
						WriterDSN:  filepath.Join(tmpDir, "master.db"),
						ReaderDSNs: []string{filepath.Join(tmpDir, "master.db")},
					},
				},
				Sharding: config.ShardingGroupConfig{
					Databases: []config.ShardConfig{
						{
							ID:         1,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding1.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding1.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding1.db")},
							TableRange: [2]int{0, 7},
						},
						{
							ID:         2,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding2.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding2.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding2.db")},
							TableRange: [2]int{8, 15},
						},
						{
							ID:         3,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding3.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding3.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding3.db")},
							TableRange: [2]int{16, 23},
						},
						{
							ID:         4,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding4.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding4.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding4.db")},
							TableRange: [2]int{24, 31},
						},
					},
				},
			},
		},
		Logging: config.LoggingConfig{
			SQLLogEnabled:   false,
			SQLLogOutputDir: tmpDir,
		},
	}

	manager, err := db.NewGroupManager(cfg)
	require.NoError(t, err)
	defer manager.CloseAll()

	tests := []struct {
		id       int64
		wantDBID int
	}{
		// ID % 32 = tableNumber, tableNumber / 8 + 1 = DB ID
		{id: 0, wantDBID: 1},   // 0 % 32 = 0, 0 / 8 + 1 = 1
		{id: 1, wantDBID: 1},   // 1 % 32 = 1, 1 / 8 + 1 = 1
		{id: 7, wantDBID: 1},   // 7 % 32 = 7, 7 / 8 + 1 = 1
		{id: 8, wantDBID: 2},   // 8 % 32 = 8, 8 / 8 + 1 = 2
		{id: 15, wantDBID: 2},  // 15 % 32 = 15, 15 / 8 + 1 = 2
		{id: 16, wantDBID: 3},  // 16 % 32 = 16, 16 / 8 + 1 = 3
		{id: 24, wantDBID: 4},  // 24 % 32 = 24, 24 / 8 + 1 = 4
		{id: 31, wantDBID: 4},  // 31 % 32 = 31, 31 / 8 + 1 = 4
		{id: 32, wantDBID: 1},  // 32 % 32 = 0, 0 / 8 + 1 = 1
		{id: 100, wantDBID: 1}, // 100 % 32 = 4, 4 / 8 + 1 = 1
	}

	for _, tt := range tests {
		t.Run("id="+string(rune(tt.id)), func(t *testing.T) {
			conn, err := manager.GetShardingConnectionByID(tt.id, "users")
			require.NoError(t, err)
			require.NotNil(t, conn)
			assert.Equal(t, tt.wantDBID, conn.ShardID)
		})
	}
}
