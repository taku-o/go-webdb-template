package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/example/go-webdb-template/internal/config"
	"gorm.io/gorm"
)

// Manager は複数のShard接続を管理
// Deprecated: 新規コードではGORMManagerを使用してください
type Manager struct {
	connections map[int]*Connection // ShardID -> Connection
	strategy    ShardingStrategy
	mu          sync.RWMutex
}

// NewManager は新しいDB Managerを作成
// Deprecated: 新規コードではNewGORMManagerを使用してください
func NewManager(cfg *config.Config) (*Manager, error) {
	manager := &Manager{
		connections: make(map[int]*Connection),
		strategy:    NewHashBasedSharding(len(cfg.Database.Shards)),
	}

	// 各Shardへの接続を確立
	for _, shardCfg := range cfg.Database.Shards {
		conn, err := NewConnection(&shardCfg)
		if err != nil {
			// すでに作成した接続をクローズ
			manager.CloseAll()
			return nil, fmt.Errorf("failed to create connection for shard %d: %w", shardCfg.ID, err)
		}
		manager.connections[shardCfg.ID] = conn
	}

	return manager, nil
}

// GetConnection はShard IDに基づいて接続を取得
func (m *Manager) GetConnection(shardID int) (*Connection, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conn, exists := m.connections[shardID]
	if !exists {
		return nil, fmt.Errorf("connection for shard %d not found", shardID)
	}
	return conn, nil
}

// GetConnectionByKey はキー（user_idなど）に基づいてShard接続を取得
func (m *Manager) GetConnectionByKey(key int64) (*Connection, error) {
	shardID := m.strategy.GetShardID(key)
	return m.GetConnection(shardID)
}

// GetDB はShard IDに基づいてsql.DBを取得
func (m *Manager) GetDB(shardID int) (*sql.DB, error) {
	conn, err := m.GetConnection(shardID)
	if err != nil {
		return nil, err
	}
	return conn.DB, nil
}

// GetDBByKey はキーに基づいてsql.DBを取得
func (m *Manager) GetDBByKey(key int64) (*sql.DB, error) {
	conn, err := m.GetConnectionByKey(key)
	if err != nil {
		return nil, err
	}
	return conn.DB, nil
}

// GetAllConnections はすべてのShard接続を取得（クロスシャードクエリ用）
func (m *Manager) GetAllConnections() []*Connection {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conns := make([]*Connection, 0, len(m.connections))
	for _, conn := range m.connections {
		conns = append(conns, conn)
	}
	return conns
}

// CloseAll はすべてのShard接続をクローズ
func (m *Manager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastErr error
	for shardID, conn := range m.connections {
		if err := conn.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close shard %d: %w", shardID, err)
		}
	}
	return lastErr
}

// PingAll はすべてのShard接続を確認
func (m *Manager) PingAll() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for shardID, conn := range m.connections {
		if err := conn.Ping(); err != nil {
			return fmt.Errorf("failed to ping shard %d: %w", shardID, err)
		}
	}
	return nil
}

// GORMManager は複数のシャードのGORM接続を管理
type GORMManager struct {
	connections map[int]*GORMConnection // ShardID -> GORMConnection
	strategy    ShardingStrategy
	mu          sync.RWMutex
}

// NewGORMManager は新しいGORM Managerを作成
func NewGORMManager(cfg *config.Config) (*GORMManager, error) {
	manager := &GORMManager{
		connections: make(map[int]*GORMConnection),
		strategy:    NewHashBasedSharding(len(cfg.Database.Shards)),
	}

	// 各シャードへの接続を確立（SQL Logger設定付き）
	for _, shardCfg := range cfg.Database.Shards {
		// SQL Loggerの作成
		sqlLogger, err := NewSQLLogger(
			shardCfg.ID,
			"shard",
			shardCfg.Driver,
			cfg.Logging.SQLLogOutputDir,
			cfg.Logging.SQLLogEnabled,
		)
		if err != nil {
			// Logger作成エラーは警告のみ（SQLクエリ実行には影響しない）
			log.Printf("Warning: Failed to create SQL logger for shard %d: %v", shardCfg.ID, err)
		}

		conn, err := NewGORMConnection(&shardCfg, sqlLogger)
		if err != nil {
			manager.CloseAll()
			return nil, fmt.Errorf("failed to create connection for shard %d: %w", shardCfg.ID, err)
		}
		manager.connections[shardCfg.ID] = conn
	}

	return manager, nil
}

// GetGORM はシャードIDに基づいて*gorm.DBを取得
func (m *GORMManager) GetGORM(shardID int) (*gorm.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conn, exists := m.connections[shardID]
	if !exists {
		return nil, fmt.Errorf("connection for shard %d not found", shardID)
	}
	return conn.DB, nil
}

// GetGORMByKey はキー（user_idなど）に基づいて*gorm.DBを取得
func (m *GORMManager) GetGORMByKey(key int64) (*gorm.DB, error) {
	shardID := m.strategy.GetShardID(key)
	return m.GetGORM(shardID)
}

// GetShardIDByKey はキーに基づいてシャードIDを取得
func (m *GORMManager) GetShardIDByKey(key int64) int {
	return m.strategy.GetShardID(key)
}

// GetAllGORMConnections はすべてのシャードのGORM接続を取得（クロスシャードクエリ用）
func (m *GORMManager) GetAllGORMConnections() []*GORMConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conns := make([]*GORMConnection, 0, len(m.connections))
	for _, conn := range m.connections {
		conns = append(conns, conn)
	}
	return conns
}

// CloseAll はすべてのシャード接続をクローズ
func (m *GORMManager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastErr error
	for shardID, conn := range m.connections {
		if err := conn.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close shard %d: %w", shardID, err)
		}
	}
	return lastErr
}

// PingAll はすべてのシャード接続を確認
func (m *GORMManager) PingAll() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for shardID, conn := range m.connections {
		if err := conn.Ping(); err != nil {
			return fmt.Errorf("failed to ping shard %d: %w", shardID, err)
		}
	}
	return nil
}
