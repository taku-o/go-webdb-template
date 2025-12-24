package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/go-webdb-template/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

// Connection は単一のDB接続を管理
type Connection struct {
	DB       *sql.DB
	ShardID  int
	Driver   string
	config   *config.ShardConfig
}

// NewConnection は新しいDB接続を作成
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
