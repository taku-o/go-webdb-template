package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/taku-o/go-webdb-template/internal/api/handler"
	"github.com/taku-o/go-webdb-template/internal/auth"
	"github.com/taku-o/go-webdb-template/internal/config"
	"github.com/taku-o/go-webdb-template/internal/db"
	"github.com/taku-o/go-webdb-template/internal/service"
	"github.com/taku-o/go-webdb-template/internal/usecase"
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

// LoadTestConfig はテスト環境の設定を読み込む
func LoadTestConfig() (*config.Config, error) {
	// 既存のAPP_ENVを保存
	oldEnv := os.Getenv("APP_ENV")
	defer func() {
		if oldEnv != "" {
			os.Setenv("APP_ENV", oldEnv)
		} else {
			os.Unsetenv("APP_ENV")
		}
	}()

	// テスト環境を設定
	os.Setenv("APP_ENV", "test")

	// 設定を読み込む
	return config.Load()
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
// dbCount: number of sharding databases (typically 4) - パラメータは互換性のため維持、設定ファイルの値を使用
// tablesPerDB: number of tables per database (typically 8, total 32 tables) - パラメータは互換性のため維持、設定ファイルの値を使用
func SetupTestGroupManager(t *testing.T, dbCount int, tablesPerDB int) *db.GroupManager {
	// ロックを取得
	fileLock, err := AcquireTestLock(t)
	if err != nil {
		t.Fatalf("Failed to acquire test lock: %v", err)
	}
	defer func() {
		if err := fileLock.Unlock(); err != nil {
			t.Logf("Warning: failed to unlock test lock: %v", err)
		}
	}()

	// 設定ファイルから読み込む
	cfg, err := LoadTestConfig()
	require.NoError(t, err)

	// 設定からGroupManagerを作成
	manager, err := db.NewGroupManager(cfg)
	require.NoError(t, err)

	// データベースをクリア
	ClearTestDatabase(t, manager)

	// Initialize master database schema (news table)
	masterConn, err := manager.GetMasterConnection()
	require.NoError(t, err)
	InitMasterSchema(t, masterConn.DB)

	// Initialize sharding database schemas (users_XXX and posts_XXX tables)
	// 設定ファイルのテーブル範囲を使用
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
	// ロックを取得
	fileLock, err := AcquireTestLock(t)
	if err != nil {
		t.Fatalf("Failed to acquire test lock: %v", err)
	}
	defer func() {
		if err := fileLock.Unlock(); err != nil {
			t.Logf("Warning: failed to unlock test lock: %v", err)
		}
	}()

	// 設定ファイルから読み込む
	cfg, err := LoadTestConfig()
	require.NoError(t, err)

	// 設定からGroupManagerを作成
	manager, err := db.NewGroupManager(cfg)
	require.NoError(t, err)

	// データベースをクリア
	ClearTestDatabase(t, manager)

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

// clearDatabaseTables は指定されたデータベースの全テーブルのデータをクリアする
func clearDatabaseTables(t *testing.T, database *gorm.DB) {
	// テーブル一覧を取得
	var tables []string
	err := database.Raw(`
		SELECT tablename
		FROM pg_tables
		WHERE schemaname = 'public'
	`).Scan(&tables).Error
	require.NoError(t, err)

	// 各テーブルをTRUNCATE
	for _, tableName := range tables {
		err := database.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tableName)).Error
		require.NoError(t, err)
	}
}

// ClearTestDatabase はテスト用データベースの全テーブルのデータをクリアする
func ClearTestDatabase(t *testing.T, manager *db.GroupManager) {
	// マスターデータベースのクリア
	masterConn, err := manager.GetMasterConnection()
	require.NoError(t, err)
	clearDatabaseTables(t, masterConn.DB)

	// シャーディングデータベースのクリア
	connections := manager.GetAllShardingConnections()
	for _, conn := range connections {
		clearDatabaseTables(t, conn.DB)
	}
}

// MockDateService はテスト用のDateServiceモック
type MockDateService struct{}

// GetToday は今日の日付を返す
func (m *MockDateService) GetToday(ctx context.Context) (string, error) {
	return time.Now().Format("2006-01-02"), nil
}

// CreateTodayHandler はテスト用のTodayHandlerを作成するヘルパー関数
func CreateTodayHandler() *handler.TodayHandler {
	mockDateService := &MockDateService{}
	todayUsecase := usecase.NewTodayUsecase(mockDateService)
	return handler.NewTodayHandler(todayUsecase)
}

// CreateDmUserHandler はテスト用のDmUserHandlerを作成するヘルパー関数
func CreateDmUserHandler(dmUserService *service.DmUserService) *handler.DmUserHandler {
	dmUserUsecase := usecase.NewDmUserUsecase(dmUserService)
	return handler.NewDmUserHandler(dmUserUsecase)
}

// CreateDmPostHandler はテスト用のDmPostHandlerを作成するヘルパー関数
func CreateDmPostHandler(dmPostService *service.DmPostService) *handler.DmPostHandler {
	dmPostUsecase := usecase.NewDmPostUsecase(dmPostService)
	return handler.NewDmPostHandler(dmPostUsecase)
}

// CreateEmailHandler はテスト用のEmailHandlerを作成するヘルパー関数
func CreateEmailHandler(emailService usecase.EmailServiceInterface, templateService usecase.TemplateServiceInterface) *handler.EmailHandler {
	emailUsecase := usecase.NewEmailUsecase(emailService, templateService)
	return handler.NewEmailHandler(emailUsecase)
}

// CreateDmJobqueueHandler はテスト用のDmJobqueueHandlerを作成するヘルパー関数
// jobQueueClient が nil の場合、usecase層も nil クライアントで初期化される
func CreateDmJobqueueHandler(jobQueueClient usecase.JobQueueClientInterface) *handler.DmJobqueueHandler {
	dmJobqueueUsecase := usecase.NewDmJobqueueUsecase(jobQueueClient)
	return handler.NewDmJobqueueHandler(dmJobqueueUsecase)
}
