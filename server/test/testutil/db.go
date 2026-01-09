package testutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/taku-o/go-webdb-template/internal/auth"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/db"
)

// TestSecretKey はテスト用の秘密鍵
const TestSecretKey = "test-secret-key-for-jwt-signing"

// TestEnv はテスト用の環境
const TestEnv = "develop"

// PostgreSQL接続設定（テスト用）
const (
	TestDBHost     = "localhost"
	TestDBUser     = "webdb"
	TestDBPassword = "webdb"
)

// GetTestAPIToken はテスト用のAPIトークンを生成
func GetTestAPIToken() (string, error) {
	return auth.GeneratePublicAPIKey(TestSecretKey, "v2", TestEnv, time.Now().Unix())
}

// GetTestConfig returns a test configuration
func GetTestConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Port: 8080,
		},
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"*"},
		},
		API: config.APIConfig{
			CurrentVersion:     "v2",
			SecretKey:          TestSecretKey,
			InvalidVersions:    []string{"v1"},
			Auth0IssuerBaseURL: "https://dev-oaa5vtzmld4dsxtd.jp.auth0.com",
		},
	}
}

// SetupTestGroupManager creates a GroupManager with PostgreSQL databases for testing
// dbCount: number of sharding databases (typically 4)
// tablesPerDB: number of tables per database (typically 8, total 32 tables)
func SetupTestGroupManager(t *testing.T, dbCount int, tablesPerDB int) *db.GroupManager {
	// Create master database config (PostgreSQL)
	masterDB := config.ShardConfig{
		ID:           1,
		Driver:       "postgres",
		Host:         TestDBHost,
		Port:         5432,
		User:         TestDBUser,
		Password:     TestDBPassword,
		Name:         "webdb_master",
		ReaderPolicy: "random",
	}

	// Create sharding databases config
	totalTables := dbCount * tablesPerDB
	shardingDBs := make([]config.ShardConfig, dbCount)
	for i := 0; i < dbCount; i++ {
		startTable := i * tablesPerDB
		endTable := startTable + tablesPerDB - 1
		if endTable >= totalTables {
			endTable = totalTables - 1
		}
		shardingDBs[i] = config.ShardConfig{
			ID:           i + 1,
			Driver:       "postgres",
			Host:         TestDBHost,
			Port:         5433 + i, // 5433, 5434, 5435, 5436
			User:         TestDBUser,
			Password:     TestDBPassword,
			Name:         fmt.Sprintf("webdb_sharding_%d", i+1),
			ReaderPolicy: "random",
			TableRange:   [2]int{startTable, endTable},
		}
	}

	// Create table configs
	tables := []config.ShardingTableConfig{
		{Name: "dm_users", SuffixCount: totalTables},
		{Name: "dm_posts", SuffixCount: totalTables},
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{masterDB},
				Sharding: config.ShardingGroupConfig{
					Databases: shardingDBs,
					Tables:    tables,
				},
			},
		},
	}

	manager, err := db.NewGroupManager(cfg)
	require.NoError(t, err)

	// Initialize master database schema (news table)
	masterConn, err := manager.GetMasterConnection()
	require.NoError(t, err)
	InitMasterSchema(t, masterConn.DB)

	// Initialize sharding database schemas (users_XXX and posts_XXX tables)
	connections := manager.GetAllShardingConnections()
	for _, conn := range connections {
		// Calculate table range for this database
		startTable := (conn.ShardID - 1) * tablesPerDB
		endTable := startTable + tablesPerDB - 1
		if endTable >= totalTables {
			endTable = totalTables - 1
		}
		InitShardingSchema(t, conn.DB, startTable, endTable)
	}

	return manager
}

// InitMasterSchema initializes the master database schema (dm_news table)
func InitMasterSchema(t *testing.T, database *gorm.DB) {
	schema := `
		CREATE TABLE IF NOT EXISTS dm_news (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			author_id INTEGER,
			published_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`
	err := database.Exec(schema).Error
	require.NoError(t, err)
}

// InitShardingSchema initializes the sharding database schema
// Creates dm_users_XXX and dm_posts_XXX tables for the given table range
func InitShardingSchema(t *testing.T, database *gorm.DB, startTable, endTable int) {
	for i := startTable; i <= endTable; i++ {
		suffix := fmt.Sprintf("%03d", i)

		usersSchema := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS dm_users_%s (
				id TEXT PRIMARY KEY,
				name TEXT NOT NULL,
				email TEXT NOT NULL UNIQUE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
		`, suffix)
		err := database.Exec(usersSchema).Error
		require.NoError(t, err)

		postsSchema := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS dm_posts_%s (
				id TEXT PRIMARY KEY,
				user_id TEXT NOT NULL,
				title TEXT NOT NULL,
				content TEXT NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
		`, suffix)
		err = database.Exec(postsSchema).Error
		require.NoError(t, err)
	}
}

// SetupTestGroupManager8Sharding creates a GroupManager with 8 sharding entries and 4 databases
// This simulates the production 8-sharding configuration with connection sharing
// - 8 sharding entries, each handling 4 tables
// - 4 actual database files (entries with same DB share connections)
// - Entries 1,2 -> postgres-sharding-1 (port 5433, tables 0-7)
// - Entries 3,4 -> postgres-sharding-2 (port 5434, tables 8-15)
// - Entries 5,6 -> postgres-sharding-3 (port 5435, tables 16-23)
// - Entries 7,8 -> postgres-sharding-4 (port 5436, tables 24-31)
func SetupTestGroupManager8Sharding(t *testing.T) *db.GroupManager {
	// Create master database config (PostgreSQL)
	masterDB := config.ShardConfig{
		ID:           1,
		Driver:       "postgres",
		Host:         TestDBHost,
		Port:         5432,
		User:         TestDBUser,
		Password:     TestDBPassword,
		Name:         "webdb_master",
		ReaderPolicy: "random",
	}

	// Create 8 sharding entries with connection sharing (4 actual databases)
	shardingDBs := []config.ShardConfig{
		// Entries 1,2 -> postgres-sharding-1 (port 5433)
		{ID: 1, Driver: "postgres", Host: TestDBHost, Port: 5433, User: TestDBUser, Password: TestDBPassword, Name: "webdb_sharding_1", ReaderPolicy: "random", TableRange: [2]int{0, 3}},
		{ID: 2, Driver: "postgres", Host: TestDBHost, Port: 5433, User: TestDBUser, Password: TestDBPassword, Name: "webdb_sharding_1", ReaderPolicy: "random", TableRange: [2]int{4, 7}},
		// Entries 3,4 -> postgres-sharding-2 (port 5434)
		{ID: 3, Driver: "postgres", Host: TestDBHost, Port: 5434, User: TestDBUser, Password: TestDBPassword, Name: "webdb_sharding_2", ReaderPolicy: "random", TableRange: [2]int{8, 11}},
		{ID: 4, Driver: "postgres", Host: TestDBHost, Port: 5434, User: TestDBUser, Password: TestDBPassword, Name: "webdb_sharding_2", ReaderPolicy: "random", TableRange: [2]int{12, 15}},
		// Entries 5,6 -> postgres-sharding-3 (port 5435)
		{ID: 5, Driver: "postgres", Host: TestDBHost, Port: 5435, User: TestDBUser, Password: TestDBPassword, Name: "webdb_sharding_3", ReaderPolicy: "random", TableRange: [2]int{16, 19}},
		{ID: 6, Driver: "postgres", Host: TestDBHost, Port: 5435, User: TestDBUser, Password: TestDBPassword, Name: "webdb_sharding_3", ReaderPolicy: "random", TableRange: [2]int{20, 23}},
		// Entries 7,8 -> postgres-sharding-4 (port 5436)
		{ID: 7, Driver: "postgres", Host: TestDBHost, Port: 5436, User: TestDBUser, Password: TestDBPassword, Name: "webdb_sharding_4", ReaderPolicy: "random", TableRange: [2]int{24, 27}},
		{ID: 8, Driver: "postgres", Host: TestDBHost, Port: 5436, User: TestDBUser, Password: TestDBPassword, Name: "webdb_sharding_4", ReaderPolicy: "random", TableRange: [2]int{28, 31}},
	}

	// Create table configs
	tables := []config.ShardingTableConfig{
		{Name: "dm_users", SuffixCount: 32},
		{Name: "dm_posts", SuffixCount: 32},
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{masterDB},
				Sharding: config.ShardingGroupConfig{
					Databases: shardingDBs,
					Tables:    tables,
				},
			},
		},
	}

	manager, err := db.NewGroupManager(cfg)
	require.NoError(t, err)

	// Initialize master database schema (news table)
	masterConn, err := manager.GetMasterConnection()
	require.NoError(t, err)
	InitMasterSchema(t, masterConn.DB)

	// Initialize sharding database schemas (users_XXX and posts_XXX tables)
	// With connection sharing, we need to create tables based on table_range for each entry
	// DB1 (ShardID=1): tables 0-7
	// DB2 (ShardID=3): tables 8-15
	// DB3 (ShardID=5): tables 16-23
	// DB4 (ShardID=7): tables 24-31
	tableRanges := map[int][2]int{
		1: {0, 7},   // Entries 1,2 -> tables 0-7
		3: {8, 15},  // Entries 3,4 -> tables 8-15
		5: {16, 23}, // Entries 5,6 -> tables 16-23
		7: {24, 31}, // Entries 7,8 -> tables 24-31
	}

	connections := manager.GetAllShardingConnections()
	for _, conn := range connections {
		tableRange, ok := tableRanges[conn.ShardID]
		if ok {
			InitShardingSchema(t, conn.DB, tableRange[0], tableRange[1])
		}
	}

	return manager
}

// CleanupTestGroupManager closes all GroupManager database connections
func CleanupTestGroupManager(manager *db.GroupManager) {
	if manager != nil {
		manager.CloseAll()
	}
}
