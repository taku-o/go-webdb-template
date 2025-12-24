package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/go-webdb-template/internal/config"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
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

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 接続プールの設定
	db.SetMaxOpenConns(cfg.MaxConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)

	// 接続確認
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
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
func createGORMConnection(cfg *config.ShardConfig, isWriter bool) (*gorm.DB, error) {
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

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 接続プール設定
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(cfg.MaxConnections)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)

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
func NewGORMConnection(cfg *config.ShardConfig) (*GORMConnection, error) {
	// 1. Writer接続を作成
	writerDB, err := createGORMConnection(cfg, true)
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
