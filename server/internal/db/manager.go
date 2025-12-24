package db

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/example/go-webdb-template/internal/config"
)

// Manager は複数のShard接続を管理
type Manager struct {
	connections map[int]*Connection // ShardID -> Connection
	strategy    ShardingStrategy
	mu          sync.RWMutex
}

// NewManager は新しいDB Managerを作成
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
