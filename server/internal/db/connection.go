package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/taku-o/go-webdb-template/internal/config"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// 接続プール設定のデフォルト値
const (
	DefaultMaxConnections        = 25
	DefaultMaxIdleConnections    = 5
	DefaultConnectionMaxLifetime = 1 * time.Hour
)

// リトライ設定
const (
	MaxRetryAttempts = 3              // 最大3回（初回 + 2回のリトライ）
	RetryDelay       = 1 * time.Second // リトライ間隔
)

// Connection は単一のDB接続を管理
// Deprecated: 新規コードではGORMConnectionを使用してください
type Connection struct {
	DB       *sql.DB
	ShardID  int
	Driver   string
	config   *config.ShardConfig
}

// NewConnection は新しいDB接続を作成
// Deprecated: 新規コードではNewGORMConnectionを使用してください
func NewConnection(cfg *config.ShardConfig) (*Connection, error) {
	dsn := cfg.GetDSN()
	if dsn == "" && cfg.DSN == "" {
		return nil, fmt.Errorf("invalid DSN for shard %d", cfg.ID)
	}

	if dsn == "" {
		dsn = cfg.DSN
	}

	driver := cfg.Driver
	if driver == "" {
		driver = "sqlite3"
	}

	var db *sql.DB
	var openErr error

	// リトライ機能付きで接続作成
	err := retry.Do(
		func() error {
			db, openErr = sql.Open(driver, dsn)
			if openErr != nil {
				return openErr
			}

			// 接続プールの設定（設定値が0以下の場合はデフォルト値を使用）
			maxConnections := cfg.MaxConnections
			if maxConnections <= 0 {
				maxConnections = DefaultMaxConnections
			}
			db.SetMaxOpenConns(maxConnections)

			maxIdleConnections := cfg.MaxIdleConnections
			if maxIdleConnections <= 0 {
				maxIdleConnections = DefaultMaxIdleConnections
			}
			db.SetMaxIdleConns(maxIdleConnections)

			connectionMaxLifetime := cfg.ConnectionMaxLifetime
			if connectionMaxLifetime <= 0 {
				connectionMaxLifetime = DefaultConnectionMaxLifetime
			}
			db.SetConnMaxLifetime(connectionMaxLifetime)

			// 接続確認は削除（遅延接続のため）
			// sql.Open()は接続を確立せず、接続オブジェクトを作成するのみ
			// 実際のクエリ実行時に接続が確立される
			return nil
		},
		retry.Attempts(MaxRetryAttempts),
		retry.Delay(RetryDelay),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Retrying database connection (attempt %d/%d): %v", n+1, MaxRetryAttempts, err)
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create database connection after %d attempts: %w", MaxRetryAttempts, err)
	}

	return &Connection{
		DB:      db,
		ShardID: cfg.ID,
		Driver:  driver,
		config:  cfg,
	}, nil
}

// Close はDB接続をクローズ
func (c *Connection) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

// Ping はDB接続を確認
func (c *Connection) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.DB.PingContext(ctx)
}

// createGORMConnection はGORM接続を作成するヘルパー関数
func createGORMConnection(cfg *config.ShardConfig, isWriter bool, sqlLogger *SQLLogger) (*gorm.DB, error) {
	var dialector gorm.Dialector

	dsn := cfg.GetWriterDSN()
	if !isWriter && len(cfg.ReaderDSNs) > 0 {
		dsn = cfg.ReaderDSNs[0] // 最初のReaderを使用
	}

	switch cfg.Driver {
	case "sqlite3":
		dialector = sqlite.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	case "mysql":
		dialector = mysql.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
	}

	// GORM ConfigにLoggerを設定
	gormConfig := &gorm.Config{}
	if sqlLogger != nil {
		gormConfig.Logger = sqlLogger
	}

	var db *gorm.DB
	var openErr error

	// リトライ機能付きで接続作成
	err := retry.Do(
		func() error {
			db, openErr = gorm.Open(dialector, gormConfig)
			if openErr != nil {
				return openErr
			}

			// 接続プール設定（設定値が0以下の場合はデフォルト値を使用）
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}

			maxConnections := cfg.MaxConnections
			if maxConnections <= 0 {
				maxConnections = DefaultMaxConnections
			}
			sqlDB.SetMaxOpenConns(maxConnections)

			maxIdleConnections := cfg.MaxIdleConnections
			if maxIdleConnections <= 0 {
				maxIdleConnections = DefaultMaxIdleConnections
			}
			sqlDB.SetMaxIdleConns(maxIdleConnections)

			connectionMaxLifetime := cfg.ConnectionMaxLifetime
			if connectionMaxLifetime <= 0 {
				connectionMaxLifetime = DefaultConnectionMaxLifetime
			}
			sqlDB.SetConnMaxLifetime(connectionMaxLifetime)

			return nil
		},
		retry.Attempts(MaxRetryAttempts),
		retry.Delay(RetryDelay),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Retrying GORM connection (attempt %d/%d): %v", n+1, MaxRetryAttempts, err)
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create GORM connection after %d attempts: %w", MaxRetryAttempts, err)
	}

	return db, nil
}

// createGORMConnectionFromDSN はDSNからGORM接続を作成
func createGORMConnectionFromDSN(dsn string, driver string) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch driver {
	case "sqlite3":
		dialector = sqlite.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	case "mysql":
		dialector = mysql.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// GORMConnection は単一のシャードのGORM接続を管理
type GORMConnection struct {
	DB      *gorm.DB // dbresolver設定済みのGORMインスタンス
	ShardID int
	Driver  string
	config  *config.ShardConfig
}

// NewGORMConnection は新しいGORM接続を作成
func NewGORMConnection(cfg *config.ShardConfig, sqlLogger *SQLLogger) (*GORMConnection, error) {
	// 1. Writer接続を作成（Logger設定付き）
	writerDB, err := createGORMConnection(cfg, true, sqlLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create writer connection: %w", err)
	}

	// 2. Reader接続を作成（複数可）
	readerDSNs := cfg.GetReaderDSNs()
	var replicaDialectors []gorm.Dialector

	// WriterのDSNと異なるReaderがある場合のみReplicaを設定
	writerDSN := cfg.GetWriterDSN()
	hasDistinctReplicas := false
	for _, readerDSN := range readerDSNs {
		if readerDSN != writerDSN {
			hasDistinctReplicas = true
			var dialector gorm.Dialector
			switch cfg.Driver {
			case "sqlite3":
				dialector = sqlite.Open(readerDSN)
			case "postgres":
				dialector = postgres.Open(readerDSN)
			case "mysql":
				dialector = mysql.Open(readerDSN)
			default:
				return nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
			}
			replicaDialectors = append(replicaDialectors, dialector)
		}
	}

	// 3. dbresolverプラグインを設定（Readerが異なる場合のみ）
	if hasDistinctReplicas && len(replicaDialectors) > 0 {
		resolverConfig := dbresolver.Config{
			Replicas: replicaDialectors,
			// Policyはデフォルト（random）を使用
		}

		err = writerDB.Use(dbresolver.Register(resolverConfig))
		if err != nil {
			return nil, fmt.Errorf("failed to register dbresolver: %w", err)
		}
	}

	return &GORMConnection{
		DB:      writerDB,
		ShardID: cfg.ID,
		Driver:  cfg.Driver,
		config:  cfg,
	}, nil
}

// Close はGORM接続をクローズ
func (c *GORMConnection) Close() error {
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// Ping はGORM接続を確認
func (c *GORMConnection) Ping() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// IsConnectionError は接続エラーかどうかを判定する
// クエリ実行時のリトライ判定に使用
func IsConnectionError(err error) bool {
	if err == nil {
		return false
	}
	// sql.ErrConnDone などの接続エラーを判定
	if err == sql.ErrConnDone {
		return true
	}
	errStr := err.Error()
	// 一般的な接続エラーのパターンを検出
	connectionErrorPatterns := []string{
		"connection",
		"network",
		"refused",
		"timeout",
		"closed",
	}
	for _, pattern := range connectionErrorPatterns {
		if contains(errStr, pattern) {
			return true
		}
	}
	return false
}

// contains は文字列に部分文字列が含まれているかを判定（大文字小文字を区別しない）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsIgnoreCase(s, substr)))
}

// containsIgnoreCase は大文字小文字を区別せずに部分文字列が含まれているかを判定
func containsIgnoreCase(s, substr string) bool {
	sLower := toLower(s)
	substrLower := toLower(substr)
	for i := 0; i <= len(sLower)-len(substrLower); i++ {
		if sLower[i:i+len(substrLower)] == substrLower {
			return true
		}
	}
	return false
}

// toLower は文字列を小文字に変換
func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

// ExecuteWithRetry はクエリ実行時のリトライ機能を提供する
// 接続エラーの場合のみリトライし、その他のエラーはリトライしない
func ExecuteWithRetry(fn func() error) error {
	return retry.Do(
		func() error {
			err := fn()
			if err != nil {
				// 接続エラーの場合はリトライ
				if IsConnectionError(err) {
					log.Printf("Connection error detected, will retry: %v", err)
					return err
				}
				// その他のエラーはリトライしない
				return retry.Unrecoverable(err)
			}
			return nil
		},
		retry.Attempts(MaxRetryAttempts),
		retry.Delay(RetryDelay),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Retrying query execution (attempt %d/%d): %v", n+1, MaxRetryAttempts, err)
		}),
	)
}
