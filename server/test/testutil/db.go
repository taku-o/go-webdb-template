package testutil

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"

	"github.com/taku-o/go-webdb-template/internal/auth"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/db"
)

// TestSecretKey はテスト用の秘密鍵
const TestSecretKey = "test-secret-key-for-jwt-signing"

// TestEnv はテスト用の環境
const TestEnv = "develop"

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

// SetupTestGroupManager creates a GroupManager with temporary file-based databases for testing
// dbCount: number of sharding databases (typically 4)
// tablesPerDB: number of tables per database (typically 8, total 32 tables)
func SetupTestGroupManager(t *testing.T, dbCount int, tablesPerDB int) *db.GroupManager {
	// Create temporary directory for test databases
	tmpDir := t.TempDir()

	// Create master database config
	masterDBPath := filepath.Join(tmpDir, "test_master.db")
	masterDB := config.ShardConfig{
		ID:           1,
		Driver:       "sqlite3",
		DSN:          masterDBPath,
		WriterDSN:    masterDBPath,
		ReaderDSNs:   []string{masterDBPath},
		ReaderPolicy: "random",
	}

	// Create sharding databases config
	totalTables := dbCount * tablesPerDB
	shardingDBs := make([]config.ShardConfig, dbCount)
	for i := 0; i < dbCount; i++ {
		dbPath := filepath.Join(tmpDir, fmt.Sprintf("test_sharding_%d.db", i+1))
		startTable := i * tablesPerDB
		endTable := startTable + tablesPerDB - 1
		if endTable >= totalTables {
			endTable = totalTables - 1
		}
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

	// Create table configs
	tables := []config.ShardingTableConfig{
		{Name: "users", SuffixCount: totalTables},
		{Name: "posts", SuffixCount: totalTables},
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

// InitMasterSchema initializes the master database schema (news table)
func InitMasterSchema(t *testing.T, database *gorm.DB) {
	schema := `
		CREATE TABLE IF NOT EXISTS news (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			author_id INTEGER,
			published_at DATETIME,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`
	err := database.Exec(schema).Error
	require.NoError(t, err)
}

// InitShardingSchema initializes the sharding database schema
// Creates users_XXX and posts_XXX tables for the given table range
func InitShardingSchema(t *testing.T, database *gorm.DB, startTable, endTable int) {
	for i := startTable; i <= endTable; i++ {
		suffix := fmt.Sprintf("%03d", i)

		usersSchema := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS users_%s (
				id INTEGER PRIMARY KEY,
				name TEXT NOT NULL,
				email TEXT NOT NULL UNIQUE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
		`, suffix)
		err := database.Exec(usersSchema).Error
		require.NoError(t, err)

		postsSchema := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS posts_%s (
				id INTEGER PRIMARY KEY,
				user_id INTEGER NOT NULL,
				title TEXT NOT NULL,
				content TEXT NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (user_id) REFERENCES users_%s(id)
			);
		`, suffix, suffix)
		err = database.Exec(postsSchema).Error
		require.NoError(t, err)
	}
}

// CleanupTestGroupManager closes all GroupManager database connections
func CleanupTestGroupManager(manager *db.GroupManager) {
	if manager != nil {
		manager.CloseAll()
	}
}
