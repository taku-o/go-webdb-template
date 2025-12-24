package testutil

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"

	"github.com/example/go-webdb-template/internal/config"
	"github.com/example/go-webdb-template/internal/db"
	"github.com/example/go-webdb-template/internal/model"
)

// SetupTestShards creates temporary file-based multi-shard databases for testing
func SetupTestShards(t *testing.T, shardCount int) *db.Manager {
	// Create temporary directory for test databases
	tmpDir := t.TempDir()

	shards := make([]config.ShardConfig, shardCount)
	for i := 0; i < shardCount; i++ {
		dbPath := filepath.Join(tmpDir, fmt.Sprintf("test_shard_%d.db", i+1))
		shards[i] = config.ShardConfig{
			ID:     i + 1,
			Driver: "sqlite3",
			DSN:    dbPath,
		}
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{Shards: shards},
	}

	manager, err := db.NewManager(cfg)
	require.NoError(t, err)

	// Initialize schema on all shards
	// Important: Use the same connection that repositories will use
	for i := 1; i <= shardCount; i++ {
		database, err := manager.GetDB(i)
		require.NoError(t, err)
		InitSchema(t, database)
	}

	return manager
}

// InitSchema initializes the database schema for testing
func InitSchema(t *testing.T, database *sql.DB) {
	schema := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`
	_, err := database.Exec(schema)
	require.NoError(t, err)
}

// CleanupTestDB closes all database connections
func CleanupTestDB(manager *db.Manager) {
	if manager != nil {
		manager.CloseAll()
	}
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
	}
}

// SetupTestGORMShards creates temporary file-based multi-shard GORM databases for testing
func SetupTestGORMShards(t *testing.T, shardCount int) *db.GORMManager {
	// Create temporary directory for test databases
	tmpDir := t.TempDir()

	shards := make([]config.ShardConfig, shardCount)
	for i := 0; i < shardCount; i++ {
		dbPath := filepath.Join(tmpDir, fmt.Sprintf("test_gorm_shard_%d.db", i+1))
		shards[i] = config.ShardConfig{
			ID:          i + 1,
			Driver:      "sqlite3",
			DSN:         dbPath,
			WriterDSN:   dbPath,
			ReaderDSNs:  []string{dbPath},
			ReaderPolicy: "random",
		}
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{Shards: shards},
	}

	manager, err := db.NewGORMManager(cfg)
	require.NoError(t, err)

	// Initialize schema on all shards using GORM AutoMigrate
	for i := 1; i <= shardCount; i++ {
		database, err := manager.GetGORM(i)
		require.NoError(t, err)
		InitGORMSchema(t, database)
	}

	return manager
}

// InitGORMSchema initializes the database schema for testing using GORM AutoMigrate
func InitGORMSchema(t *testing.T, database *gorm.DB) {
	err := database.AutoMigrate(&model.User{}, &model.Post{})
	require.NoError(t, err)
}

// CleanupTestGORMDB closes all GORM database connections
func CleanupTestGORMDB(manager *db.GORMManager) {
	if manager != nil {
		manager.CloseAll()
	}
}
