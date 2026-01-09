package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/test/testutil"
)

// =============================================================================
// タスク4.1, 4.2, 4.3, 4.4: GroupManager, MasterManager, ShardingManagerのテスト
// =============================================================================

// TestNewMasterManager tests MasterManager creation
func TestNewMasterManager(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{
					{
						ID:           1,
						Driver:       "postgres",
						Host:         testutil.TestDBHost,
						Port:         5432,
						User:         testutil.TestDBUser,
						Password:     testutil.TestDBPassword,
						Name:         "webdb_master",
						ReaderPolicy: "random",
					},
				},
			},
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
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Sharding: config.ShardingGroupConfig{
					Databases: []config.ShardConfig{
						{
							ID:           1,
							Driver:       "postgres",
							Host:         testutil.TestDBHost,
							Port:         5433,
							User:         testutil.TestDBUser,
							Password:     testutil.TestDBPassword,
							Name:         "webdb_sharding_1",
							ReaderPolicy: "random",
							TableRange:   [2]int{0, 7},
						},
						{
							ID:           2,
							Driver:       "postgres",
							Host:         testutil.TestDBHost,
							Port:         5434,
							User:         testutil.TestDBUser,
							Password:     testutil.TestDBPassword,
							Name:         "webdb_sharding_2",
							ReaderPolicy: "random",
							TableRange:   [2]int{8, 15},
						},
						{
							ID:           3,
							Driver:       "postgres",
							Host:         testutil.TestDBHost,
							Port:         5435,
							User:         testutil.TestDBUser,
							Password:     testutil.TestDBPassword,
							Name:         "webdb_sharding_3",
							ReaderPolicy: "random",
							TableRange:   [2]int{16, 23},
						},
						{
							ID:           4,
							Driver:       "postgres",
							Host:         testutil.TestDBHost,
							Port:         5436,
							User:         testutil.TestDBUser,
							Password:     testutil.TestDBPassword,
							Name:         "webdb_sharding_4",
							ReaderPolicy: "random",
							TableRange:   [2]int{24, 31},
						},
					},
					Tables: []config.ShardingTableConfig{
						{Name: "dm_users", SuffixCount: 32},
						{Name: "dm_posts", SuffixCount: 32},
					},
				},
			},
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
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Sharding: config.ShardingGroupConfig{
					Databases: []config.ShardConfig{
						{
							ID:           1,
							Driver:       "postgres",
							Host:         testutil.TestDBHost,
							Port:         5433,
							User:         testutil.TestDBUser,
							Password:     testutil.TestDBPassword,
							Name:         "webdb_sharding_1",
							ReaderPolicy: "random",
							TableRange:   [2]int{0, 7},
						},
						{
							ID:           2,
							Driver:       "postgres",
							Host:         testutil.TestDBHost,
							Port:         5434,
							User:         testutil.TestDBUser,
							Password:     testutil.TestDBPassword,
							Name:         "webdb_sharding_2",
							ReaderPolicy: "random",
							TableRange:   [2]int{8, 15},
						},
						{
							ID:           3,
							Driver:       "postgres",
							Host:         testutil.TestDBHost,
							Port:         5435,
							User:         testutil.TestDBUser,
							Password:     testutil.TestDBPassword,
							Name:         "webdb_sharding_3",
							ReaderPolicy: "random",
							TableRange:   [2]int{16, 23},
						},
						{
							ID:           4,
							Driver:       "postgres",
							Host:         testutil.TestDBHost,
							Port:         5436,
							User:         testutil.TestDBUser,
							Password:     testutil.TestDBPassword,
							Name:         "webdb_sharding_4",
							ReaderPolicy: "random",
							TableRange:   [2]int{24, 31},
						},
					},
				},
			},
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
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Sharding: config.ShardingGroupConfig{
					Databases: []config.ShardConfig{
						{
							ID:           1,
							Driver:       "postgres",
							Host:         testutil.TestDBHost,
							Port:         5433,
							User:         testutil.TestDBUser,
							Password:     testutil.TestDBPassword,
							Name:         "webdb_sharding_1",
							ReaderPolicy: "random",
							TableRange:   [2]int{0, 7},
						},
						{
							ID:           2,
							Driver:       "postgres",
							Host:         testutil.TestDBHost,
							Port:         5434,
							User:         testutil.TestDBUser,
							Password:     testutil.TestDBPassword,
							Name:         "webdb_sharding_2",
							ReaderPolicy: "random",
							TableRange:   [2]int{8, 15},
						},
					},
				},
			},
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
	// 8つのシャーディングエントリ、4つのデータベース
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Sharding: config.ShardingGroupConfig{
					Databases: []config.ShardConfig{
						// Entry 1,2 -> postgres-sharding-1 (port 5433)
						{ID: 1, Driver: "postgres", Host: testutil.TestDBHost, Port: 5433, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_1", ReaderPolicy: "random", TableRange: [2]int{0, 3}},
						{ID: 2, Driver: "postgres", Host: testutil.TestDBHost, Port: 5433, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_1", ReaderPolicy: "random", TableRange: [2]int{4, 7}},
						// Entry 3,4 -> postgres-sharding-2 (port 5434)
						{ID: 3, Driver: "postgres", Host: testutil.TestDBHost, Port: 5434, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_2", ReaderPolicy: "random", TableRange: [2]int{8, 11}},
						{ID: 4, Driver: "postgres", Host: testutil.TestDBHost, Port: 5434, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_2", ReaderPolicy: "random", TableRange: [2]int{12, 15}},
						// Entry 5,6 -> postgres-sharding-3 (port 5435)
						{ID: 5, Driver: "postgres", Host: testutil.TestDBHost, Port: 5435, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_3", ReaderPolicy: "random", TableRange: [2]int{16, 19}},
						{ID: 6, Driver: "postgres", Host: testutil.TestDBHost, Port: 5435, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_3", ReaderPolicy: "random", TableRange: [2]int{20, 23}},
						// Entry 7,8 -> postgres-sharding-4 (port 5436)
						{ID: 7, Driver: "postgres", Host: testutil.TestDBHost, Port: 5436, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_4", ReaderPolicy: "random", TableRange: [2]int{24, 27}},
						{ID: 8, Driver: "postgres", Host: testutil.TestDBHost, Port: 5436, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_4", ReaderPolicy: "random", TableRange: [2]int{28, 31}},
					},
					Tables: []config.ShardingTableConfig{
						{Name: "dm_users", SuffixCount: 32},
						{Name: "dm_posts", SuffixCount: 32},
					},
				},
			},
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

// TestShardingManager_ConnectionSharing tests that entries with the same connection info share the same connection
func TestShardingManager_ConnectionSharing(t *testing.T) {
	// 2つのエントリが同じ接続情報を共有
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Sharding: config.ShardingGroupConfig{
					Databases: []config.ShardConfig{
						{
							ID:           1,
							Driver:       "postgres",
							Host:         testutil.TestDBHost,
							Port:         5433,
							User:         testutil.TestDBUser,
							Password:     testutil.TestDBPassword,
							Name:         "webdb_sharding_1",
							ReaderPolicy: "random",
							TableRange:   [2]int{0, 3},
						},
						{
							ID:           2,
							Driver:       "postgres",
							Host:         testutil.TestDBHost,
							Port:         5433,
							User:         testutil.TestDBUser,
							Password:     testutil.TestDBPassword,
							Name:         "webdb_sharding_1", // 同じ接続情報
							ReaderPolicy: "random",
							TableRange:   [2]int{4, 7},
						},
						{
							ID:           3,
							Driver:       "postgres",
							Host:         testutil.TestDBHost,
							Port:         5434,
							User:         testutil.TestDBUser,
							Password:     testutil.TestDBPassword,
							Name:         "webdb_sharding_2", // 異なる接続情報
							ReaderPolicy: "random",
							TableRange:   [2]int{8, 15},
						},
					},
				},
			},
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
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Sharding: config.ShardingGroupConfig{
					Databases: []config.ShardConfig{
						{ID: 1, Driver: "postgres", Host: testutil.TestDBHost, Port: 5433, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_1", ReaderPolicy: "random", TableRange: [2]int{0, 3}},
						{ID: 2, Driver: "postgres", Host: testutil.TestDBHost, Port: 5433, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_1", ReaderPolicy: "random", TableRange: [2]int{4, 7}},
						{ID: 3, Driver: "postgres", Host: testutil.TestDBHost, Port: 5434, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_2", ReaderPolicy: "random", TableRange: [2]int{8, 11}},
						{ID: 4, Driver: "postgres", Host: testutil.TestDBHost, Port: 5434, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_2", ReaderPolicy: "random", TableRange: [2]int{12, 15}},
						{ID: 5, Driver: "postgres", Host: testutil.TestDBHost, Port: 5435, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_3", ReaderPolicy: "random", TableRange: [2]int{16, 19}},
						{ID: 6, Driver: "postgres", Host: testutil.TestDBHost, Port: 5435, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_3", ReaderPolicy: "random", TableRange: [2]int{20, 23}},
						{ID: 7, Driver: "postgres", Host: testutil.TestDBHost, Port: 5436, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_4", ReaderPolicy: "random", TableRange: [2]int{24, 27}},
						{ID: 8, Driver: "postgres", Host: testutil.TestDBHost, Port: 5436, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_4", ReaderPolicy: "random", TableRange: [2]int{28, 31}},
					},
				},
			},
		},
	}

	manager, err := db.NewShardingManager(cfg)
	require.NoError(t, err)
	defer manager.CloseAll()

	// 接続共有により、同じ接続情報を持つ複数のエントリは同じ接続を共有する
	// ShardIDは最初に登録されたエントリのIDになるため、
	// テーブル番号0-7はShardID=1、テーブル番号8-15はShardID=3、
	// テーブル番号16-23はShardID=5、テーブル番号24-31はShardID=7となる
	tests := []struct {
		tableNumber  int
		wantShardID  int // 接続共有により、最初に登録されたエントリのShardIDを期待
		expectError  bool
		errorMessage string
	}{
		// sharding_db_1（エントリID 1, 2が共有）: テーブル番号 0-7
		{tableNumber: 0, wantShardID: 1, expectError: false},
		{tableNumber: 3, wantShardID: 1, expectError: false},
		{tableNumber: 4, wantShardID: 1, expectError: false}, // 接続共有: ShardID=1
		{tableNumber: 7, wantShardID: 1, expectError: false}, // 接続共有: ShardID=1
		// sharding_db_2（エントリID 3, 4が共有）: テーブル番号 8-15
		{tableNumber: 8, wantShardID: 3, expectError: false},
		{tableNumber: 11, wantShardID: 3, expectError: false},
		{tableNumber: 12, wantShardID: 3, expectError: false}, // 接続共有: ShardID=3
		{tableNumber: 15, wantShardID: 3, expectError: false}, // 接続共有: ShardID=3
		// sharding_db_3（エントリID 5, 6が共有）: テーブル番号 16-23
		{tableNumber: 16, wantShardID: 5, expectError: false},
		{tableNumber: 19, wantShardID: 5, expectError: false},
		{tableNumber: 20, wantShardID: 5, expectError: false}, // 接続共有: ShardID=5
		{tableNumber: 23, wantShardID: 5, expectError: false}, // 接続共有: ShardID=5
		// sharding_db_4（エントリID 7, 8が共有）: テーブル番号 24-31
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
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Sharding: config.ShardingGroupConfig{
					Databases: []config.ShardConfig{
						{ID: 1, Driver: "postgres", Host: testutil.TestDBHost, Port: 5433, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_1", ReaderPolicy: "random", TableRange: [2]int{0, 3}},
						{ID: 2, Driver: "postgres", Host: testutil.TestDBHost, Port: 5433, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_1", ReaderPolicy: "random", TableRange: [2]int{4, 7}},
						{ID: 3, Driver: "postgres", Host: testutil.TestDBHost, Port: 5434, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_2", ReaderPolicy: "random", TableRange: [2]int{8, 11}},
						{ID: 4, Driver: "postgres", Host: testutil.TestDBHost, Port: 5434, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_2", ReaderPolicy: "random", TableRange: [2]int{12, 15}},
						{ID: 5, Driver: "postgres", Host: testutil.TestDBHost, Port: 5435, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_3", ReaderPolicy: "random", TableRange: [2]int{16, 19}},
						{ID: 6, Driver: "postgres", Host: testutil.TestDBHost, Port: 5435, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_3", ReaderPolicy: "random", TableRange: [2]int{20, 23}},
						{ID: 7, Driver: "postgres", Host: testutil.TestDBHost, Port: 5436, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_4", ReaderPolicy: "random", TableRange: [2]int{24, 27}},
						{ID: 8, Driver: "postgres", Host: testutil.TestDBHost, Port: 5436, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_4", ReaderPolicy: "random", TableRange: [2]int{28, 31}},
					},
				},
			},
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
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{
					{
						ID:           1,
						Driver:       "postgres",
						Host:         testutil.TestDBHost,
						Port:         5432,
						User:         testutil.TestDBUser,
						Password:     testutil.TestDBPassword,
						Name:         "webdb_master",
						ReaderPolicy: "random",
					},
				},
				Sharding: config.ShardingGroupConfig{
					Databases: []config.ShardConfig{
						{ID: 1, Driver: "postgres", Host: testutil.TestDBHost, Port: 5433, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_1", ReaderPolicy: "random", TableRange: [2]int{0, 7}},
						{ID: 2, Driver: "postgres", Host: testutil.TestDBHost, Port: 5434, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_2", ReaderPolicy: "random", TableRange: [2]int{8, 15}},
						{ID: 3, Driver: "postgres", Host: testutil.TestDBHost, Port: 5435, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_3", ReaderPolicy: "random", TableRange: [2]int{16, 23}},
						{ID: 4, Driver: "postgres", Host: testutil.TestDBHost, Port: 5436, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_4", ReaderPolicy: "random", TableRange: [2]int{24, 31}},
					},
					Tables: []config.ShardingTableConfig{
						{Name: "dm_users", SuffixCount: 32},
						{Name: "dm_posts", SuffixCount: 32},
					},
				},
			},
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
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{
					{
						ID:           1,
						Driver:       "postgres",
						Host:         testutil.TestDBHost,
						Port:         5432,
						User:         testutil.TestDBUser,
						Password:     testutil.TestDBPassword,
						Name:         "webdb_master",
						ReaderPolicy: "random",
					},
				},
				Sharding: config.ShardingGroupConfig{
					Databases: []config.ShardConfig{
						{ID: 1, Driver: "postgres", Host: testutil.TestDBHost, Port: 5433, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_1", ReaderPolicy: "random", TableRange: [2]int{0, 7}},
						{ID: 2, Driver: "postgres", Host: testutil.TestDBHost, Port: 5434, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_2", ReaderPolicy: "random", TableRange: [2]int{8, 15}},
						{ID: 3, Driver: "postgres", Host: testutil.TestDBHost, Port: 5435, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_3", ReaderPolicy: "random", TableRange: [2]int{16, 23}},
						{ID: 4, Driver: "postgres", Host: testutil.TestDBHost, Port: 5436, User: testutil.TestDBUser, Password: testutil.TestDBPassword, Name: "webdb_sharding_4", ReaderPolicy: "random", TableRange: [2]int{24, 31}},
					},
				},
			},
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
			conn, err := manager.GetShardingConnectionByID(tt.id, "dm_users")
			require.NoError(t, err)
			require.NotNil(t, conn)
			assert.Equal(t, tt.wantDBID, conn.ShardID)
		})
	}
}

// =============================================================================
// 異なるDB構成のテスト（4台DB、1台DB）
// =============================================================================

// TestShardingManager_4Databases tests 8 sharding entries with 4 physical databases
// This matches the actual production environment: 8 logical shards distributed across 4 PostgreSQL containers
func TestShardingManager_4Databases(t *testing.T) {
	// 4つの物理データベースに8つの論理シャードを分散
	// Entry 1-2: port 5433 (webdb_sharding_1)
	// Entry 3-4: port 5434 (webdb_sharding_2)
	// Entry 5-6: port 5435 (webdb_sharding_3)
	// Entry 7-8: port 5436 (webdb_sharding_4)
	ports := []int{5433, 5433, 5434, 5434, 5435, 5435, 5436, 5436}
	dbNames := []string{
		"webdb_sharding_1", "webdb_sharding_1",
		"webdb_sharding_2", "webdb_sharding_2",
		"webdb_sharding_3", "webdb_sharding_3",
		"webdb_sharding_4", "webdb_sharding_4",
	}

	shardingDBs := make([]config.ShardConfig, 8)
	for i := 0; i < 8; i++ {
		startTable := i * 4
		endTable := startTable + 3
		shardingDBs[i] = config.ShardConfig{
			ID:           i + 1,
			Driver:       "postgres",
			Host:         testutil.TestDBHost,
			Port:         ports[i],
			User:         testutil.TestDBUser,
			Password:     testutil.TestDBPassword,
			Name:         dbNames[i],
			ReaderPolicy: "random",
			TableRange:   [2]int{startTable, endTable},
		}
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{
					{
						ID:           1,
						Driver:       "postgres",
						Host:         testutil.TestDBHost,
						Port:         5432,
						User:         testutil.TestDBUser,
						Password:     testutil.TestDBPassword,
						Name:         "webdb_master",
						ReaderPolicy: "random",
					},
				},
				Sharding: config.ShardingGroupConfig{
					Databases: shardingDBs,
					Tables: []config.ShardingTableConfig{
						{Name: "dm_users", SuffixCount: 32},
					},
				},
			},
		},
	}

	manager, err := db.NewGroupManager(cfg)
	require.NoError(t, err)
	require.NotNil(t, manager)
	defer manager.CloseAll()

	// 4つのユニークな接続が作成されることを確認
	connections := manager.GetAllShardingConnections()
	assert.Len(t, connections, 4, "Should have 4 unique connections for 4 physical databases")

	// テーブル番号と期待されるShardIDの対応を確認
	// tables 0-7 → ShardID 1 (port 5433)
	// tables 8-15 → ShardID 3 (port 5434)
	// tables 16-23 → ShardID 5 (port 5435)
	// tables 24-31 → ShardID 7 (port 5436)
	expectedShardIDs := map[int]int{
		0: 1, 4: 1, 7: 1,      // port 5433
		8: 3, 12: 3, 15: 3,    // port 5434
		16: 5, 20: 5, 23: 5,   // port 5435
		24: 7, 28: 7, 31: 7,   // port 5436
	}

	for tableNumber, expectedShardID := range expectedShardIDs {
		conn, err := manager.GetShardingConnection(tableNumber)
		require.NoError(t, err)
		assert.Equal(t, expectedShardID, conn.ShardID, "Table %d should return connection with ShardID %d", tableNumber, expectedShardID)
	}

	// 同じ物理DBを使うテーブルが同じ接続オブジェクトを共有することを確認
	conn0, _ := manager.GetShardingConnection(0)
	conn4, _ := manager.GetShardingConnection(4)
	assert.Same(t, conn0, conn4, "Tables 0 and 4 should share the same connection (port 5433)")

	conn8, _ := manager.GetShardingConnection(8)
	conn12, _ := manager.GetShardingConnection(12)
	assert.Same(t, conn8, conn12, "Tables 8 and 12 should share the same connection (port 5434)")

	// 異なる物理DBのテーブルは異なる接続オブジェクトを使うことを確認
	assert.NotSame(t, conn0, conn8, "Tables 0 and 8 should use different connections")
}

// TestShardingManager_1Database tests 8 sharding entries with 1 shared database
func TestShardingManager_1Database(t *testing.T) {
	// 1つのデータベースを全エントリで共有
	shardingDBs := make([]config.ShardConfig, 8)
	for i := 0; i < 8; i++ {
		startTable := i * 4
		endTable := startTable + 3
		shardingDBs[i] = config.ShardConfig{
			ID:           i + 1,
			Driver:       "postgres",
			Host:         testutil.TestDBHost,
			Port:         5433, // 全エントリで同じポート
			User:         testutil.TestDBUser,
			Password:     testutil.TestDBPassword,
			Name:         "webdb_sharding_1", // 全エントリで同じDB名
			ReaderPolicy: "random",
			TableRange:   [2]int{startTable, endTable},
		}
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{
					{
						ID:           1,
						Driver:       "postgres",
						Host:         testutil.TestDBHost,
						Port:         5432,
						User:         testutil.TestDBUser,
						Password:     testutil.TestDBPassword,
						Name:         "webdb_master",
						ReaderPolicy: "random",
					},
				},
				Sharding: config.ShardingGroupConfig{
					Databases: shardingDBs,
					Tables: []config.ShardingTableConfig{
						{Name: "dm_users", SuffixCount: 32},
					},
				},
			},
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
