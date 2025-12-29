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

// =============================================================================
// シャーディング数8対応のテスト
// =============================================================================

// TestNewShardingManager_8Sharding tests ShardingManager with 8 sharding entries
func TestNewShardingManager_8Sharding(t *testing.T) {
	tmpDir := t.TempDir()

	// 8つのシャーディングエントリ、4つのデータベース
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
							TableRange: [2]int{0, 3},
						},
						{
							ID:         2,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding1.db"), // 同じDSN（接続共有）
							WriterDSN:  filepath.Join(tmpDir, "sharding1.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding1.db")},
							TableRange: [2]int{4, 7},
						},
						{
							ID:         3,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding2.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding2.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding2.db")},
							TableRange: [2]int{8, 11},
						},
						{
							ID:         4,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding2.db"), // 同じDSN（接続共有）
							WriterDSN:  filepath.Join(tmpDir, "sharding2.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding2.db")},
							TableRange: [2]int{12, 15},
						},
						{
							ID:         5,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding3.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding3.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding3.db")},
							TableRange: [2]int{16, 19},
						},
						{
							ID:         6,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding3.db"), // 同じDSN（接続共有）
							WriterDSN:  filepath.Join(tmpDir, "sharding3.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding3.db")},
							TableRange: [2]int{20, 23},
						},
						{
							ID:         7,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding4.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding4.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding4.db")},
							TableRange: [2]int{24, 27},
						},
						{
							ID:         8,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding4.db"), // 同じDSN（接続共有）
							WriterDSN:  filepath.Join(tmpDir, "sharding4.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding4.db")},
							TableRange: [2]int{28, 31},
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

// TestShardingManager_ConnectionSharing tests that entries with the same DSN share the same connection
func TestShardingManager_ConnectionSharing(t *testing.T) {
	tmpDir := t.TempDir()

	// 2つのエントリが同じDSNを共有
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Sharding: config.ShardingGroupConfig{
					Databases: []config.ShardConfig{
						{
							ID:         1,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "shared.db"),
							WriterDSN:  filepath.Join(tmpDir, "shared.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "shared.db")},
							TableRange: [2]int{0, 3},
						},
						{
							ID:         2,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "shared.db"), // 同じDSN
							WriterDSN:  filepath.Join(tmpDir, "shared.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "shared.db")},
							TableRange: [2]int{4, 7},
						},
						{
							ID:         3,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "other.db"), // 異なるDSN
							WriterDSN:  filepath.Join(tmpDir, "other.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "other.db")},
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

	// GetAllConnectionsはユニークな接続のみを返す（2つの接続）
	connections := manager.GetAllConnections()
	assert.Len(t, connections, 2, "接続共有により、ユニークな接続は2つのみ")
}

// TestShardingManager_GetConnectionByTableNumber_8Sharding tests connection by table number with 8 sharding entries
func TestShardingManager_GetConnectionByTableNumber_8Sharding(t *testing.T) {
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
							TableRange: [2]int{0, 3},
						},
						{
							ID:         2,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding1.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding1.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding1.db")},
							TableRange: [2]int{4, 7},
						},
						{
							ID:         3,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding2.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding2.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding2.db")},
							TableRange: [2]int{8, 11},
						},
						{
							ID:         4,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding2.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding2.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding2.db")},
							TableRange: [2]int{12, 15},
						},
						{
							ID:         5,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding3.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding3.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding3.db")},
							TableRange: [2]int{16, 19},
						},
						{
							ID:         6,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding3.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding3.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding3.db")},
							TableRange: [2]int{20, 23},
						},
						{
							ID:         7,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding4.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding4.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding4.db")},
							TableRange: [2]int{24, 27},
						},
						{
							ID:         8,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding4.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding4.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding4.db")},
							TableRange: [2]int{28, 31},
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

	// 接続共有により、同じDSNを持つ複数のエントリは同じ接続を共有する
	// ShardIDは最初に登録されたエントリのIDになるため、
	// テーブル番号0-7はShardID=1、テーブル番号8-15はShardID=3、
	// テーブル番号16-23はShardID=5、テーブル番号24-31はShardID=7となる
	tests := []struct {
		tableNumber  int
		wantShardID  int // 接続共有により、最初に登録されたエントリのShardIDを期待
		expectError  bool
		errorMessage string
	}{
		// sharding_db_1.db（エントリID 1, 2が共有）: テーブル番号 0-7
		{tableNumber: 0, wantShardID: 1, expectError: false},
		{tableNumber: 3, wantShardID: 1, expectError: false},
		{tableNumber: 4, wantShardID: 1, expectError: false}, // 接続共有: ShardID=1
		{tableNumber: 7, wantShardID: 1, expectError: false}, // 接続共有: ShardID=1
		// sharding_db_2.db（エントリID 3, 4が共有）: テーブル番号 8-15
		{tableNumber: 8, wantShardID: 3, expectError: false},
		{tableNumber: 11, wantShardID: 3, expectError: false},
		{tableNumber: 12, wantShardID: 3, expectError: false}, // 接続共有: ShardID=3
		{tableNumber: 15, wantShardID: 3, expectError: false}, // 接続共有: ShardID=3
		// sharding_db_3.db（エントリID 5, 6が共有）: テーブル番号 16-23
		{tableNumber: 16, wantShardID: 5, expectError: false},
		{tableNumber: 19, wantShardID: 5, expectError: false},
		{tableNumber: 20, wantShardID: 5, expectError: false}, // 接続共有: ShardID=5
		{tableNumber: 23, wantShardID: 5, expectError: false}, // 接続共有: ShardID=5
		// sharding_db_4.db（エントリID 7, 8が共有）: テーブル番号 24-31
		{tableNumber: 24, wantShardID: 7, expectError: false},
		{tableNumber: 27, wantShardID: 7, expectError: false},
		{tableNumber: 28, wantShardID: 7, expectError: false}, // 接続共有: ShardID=7
		{tableNumber: 31, wantShardID: 7, expectError: false}, // 接続共有: ShardID=7
		// エラーケース
		{tableNumber: -1, expectError: true, errorMessage: "invalid table number"},
		{tableNumber: 32, expectError: true, errorMessage: "invalid table number"},
	}

	for _, tt := range tests {
		t.Run("tableNumber="+string(rune('0'+tt.tableNumber)), func(t *testing.T) {
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
				assert.Equal(t, tt.wantShardID, conn.ShardID)
			}
		})
	}
}

// TestShardingManager_GetAllConnections_8Sharding tests GetAllConnections returns unique connections with 8 sharding entries
func TestShardingManager_GetAllConnections_8Sharding(t *testing.T) {
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
							TableRange: [2]int{0, 3},
						},
						{
							ID:         2,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding1.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding1.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding1.db")},
							TableRange: [2]int{4, 7},
						},
						{
							ID:         3,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding2.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding2.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding2.db")},
							TableRange: [2]int{8, 11},
						},
						{
							ID:         4,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding2.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding2.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding2.db")},
							TableRange: [2]int{12, 15},
						},
						{
							ID:         5,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding3.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding3.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding3.db")},
							TableRange: [2]int{16, 19},
						},
						{
							ID:         6,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding3.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding3.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding3.db")},
							TableRange: [2]int{20, 23},
						},
						{
							ID:         7,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding4.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding4.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding4.db")},
							TableRange: [2]int{24, 27},
						},
						{
							ID:         8,
							Driver:     "sqlite3",
							DSN:        filepath.Join(tmpDir, "sharding4.db"),
							WriterDSN:  filepath.Join(tmpDir, "sharding4.db"),
							ReaderDSNs: []string{filepath.Join(tmpDir, "sharding4.db")},
							TableRange: [2]int{28, 31},
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

	// 8つのシャーディングエントリがあるが、実際のDB接続は4つのみ
	connections := manager.GetAllConnections()
	assert.Len(t, connections, 4, "接続共有により、ユニークな接続は4つのみ")
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

// =============================================================================
// 異なるDB構成のテスト（8台DB、1台DB）
// =============================================================================

// TestShardingManager_8Databases tests 8 sharding entries with 8 separate databases
// Each entry has its own database (no connection sharing)
func TestShardingManager_8Databases(t *testing.T) {
	tmpDir := t.TempDir()

	// 8つの別々のデータベースファイルを使用
	shardingDBs := make([]config.ShardConfig, 8)
	for i := 0; i < 8; i++ {
		dbPath := filepath.Join(tmpDir, "sharding_db_"+string(rune('1'+i))+".db")
		startTable := i * 4
		endTable := startTable + 3
		shardingDBs[i] = config.ShardConfig{
			ID:           i + 1,
			Driver:       "sqlite3",
			DSN:          dbPath,
			WriterDSN:    dbPath,
			ReaderDSNs:   []string{dbPath},
			ReaderPolicy: "random",
			TableRange:   [2]int{startTable, endTable},
		}
	}

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
					Databases: shardingDBs,
					Tables: []config.ShardingTableConfig{
						{Name: "users", SuffixCount: 32},
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

	// 8つのユニークな接続が作成されることを確認
	connections := manager.GetAllShardingConnections()
	assert.Len(t, connections, 8, "Should have 8 unique connections for 8 separate databases")

	// 各テーブル番号が正しい接続を返すことを確認
	testCases := []struct {
		tableNumber int
		expectDBID  int
	}{
		{0, 1}, {3, 1},   // Entry 1: tables 0-3
		{4, 2}, {7, 2},   // Entry 2: tables 4-7
		{8, 3}, {11, 3},  // Entry 3: tables 8-11
		{12, 4}, {15, 4}, // Entry 4: tables 12-15
		{16, 5}, {19, 5}, // Entry 5: tables 16-19
		{20, 6}, {23, 6}, // Entry 6: tables 20-23
		{24, 7}, {27, 7}, // Entry 7: tables 24-27
		{28, 8}, {31, 8}, // Entry 8: tables 28-31
	}

	for _, tc := range testCases {
		conn, err := manager.GetShardingConnection(tc.tableNumber)
		require.NoError(t, err)
		assert.Equal(t, tc.expectDBID, conn.ShardID, "Table %d should return connection with ShardID %d", tc.tableNumber, tc.expectDBID)
	}
}

// TestShardingManager_1Database tests 8 sharding entries with 1 shared database
// All entries share the same database connection
func TestShardingManager_1Database(t *testing.T) {
	tmpDir := t.TempDir()

	// 1つのデータベースファイルを全エントリで共有
	sharedDBPath := filepath.Join(tmpDir, "sharding_shared.db")
	shardingDBs := make([]config.ShardConfig, 8)
	for i := 0; i < 8; i++ {
		startTable := i * 4
		endTable := startTable + 3
		shardingDBs[i] = config.ShardConfig{
			ID:           i + 1,
			Driver:       "sqlite3",
			DSN:          sharedDBPath, // 全エントリで同じDSN
			WriterDSN:    sharedDBPath,
			ReaderDSNs:   []string{sharedDBPath},
			ReaderPolicy: "random",
			TableRange:   [2]int{startTable, endTable},
		}
	}

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
					Databases: shardingDBs,
					Tables: []config.ShardingTableConfig{
						{Name: "users", SuffixCount: 32},
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

	// 1つのユニークな接続のみが作成されることを確認
	connections := manager.GetAllShardingConnections()
	assert.Len(t, connections, 1, "Should have only 1 unique connection for 1 shared database")

	// 全テーブル番号が同じ接続（最初のエントリのShardID=1）を返すことを確認
	for tableNumber := 0; tableNumber < 32; tableNumber++ {
		conn, err := manager.GetShardingConnection(tableNumber)
		require.NoError(t, err)
		assert.Equal(t, 1, conn.ShardID, "Table %d should return connection with ShardID 1 (first entry)", tableNumber)
	}

	// 全てのテーブル番号が同じ接続オブジェクトを返すことを確認
	conn0, _ := manager.GetShardingConnection(0)
	conn15, _ := manager.GetShardingConnection(15)
	conn31, _ := manager.GetShardingConnection(31)
	assert.Same(t, conn0, conn15, "Tables 0 and 15 should share the same connection")
	assert.Same(t, conn0, conn31, "Tables 0 and 31 should share the same connection")
}
