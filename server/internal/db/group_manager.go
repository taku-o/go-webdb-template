package db

import (
	"fmt"
	"log"
	"sync"

	"github.com/example/go-webdb-template/internal/config"
)

// =============================================================================
// タスク4.1, 4.2, 4.3: GroupManager, MasterManager, ShardingManagerの実装
// =============================================================================

// GroupManager はmaster/shardingグループの接続を統合管理
type GroupManager struct {
	masterManager   *MasterManager
	shardingManager *ShardingManager
	mu              sync.RWMutex
}

// NewGroupManager は新しいGroupManagerを作成
func NewGroupManager(cfg *config.Config) (*GroupManager, error) {
	// MasterManagerの作成
	masterManager, err := NewMasterManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create master manager: %w", err)
	}

	// ShardingManagerの作成
	shardingManager, err := NewShardingManager(cfg)
	if err != nil {
		masterManager.CloseAll()
		return nil, fmt.Errorf("failed to create sharding manager: %w", err)
	}

	return &GroupManager{
		masterManager:   masterManager,
		shardingManager: shardingManager,
	}, nil
}

// GetMasterConnection はmasterグループの接続を取得
func (gm *GroupManager) GetMasterConnection() (*GORMConnection, error) {
	return gm.masterManager.GetConnection()
}

// GetShardingConnection はテーブル番号からshardingグループの接続を取得
func (gm *GroupManager) GetShardingConnection(tableNumber int) (*GORMConnection, error) {
	return gm.shardingManager.GetConnectionByTableNumber(tableNumber)
}

// GetShardingConnectionByID はIDからshardingグループの接続を取得
func (gm *GroupManager) GetShardingConnectionByID(id int64, tableName string) (*GORMConnection, error) {
	tableNumber := int(id % 32)
	return gm.GetShardingConnection(tableNumber)
}

// GetAllShardingConnections はすべてのsharding接続を取得（クロステーブルクエリ用）
func (gm *GroupManager) GetAllShardingConnections() []*GORMConnection {
	return gm.shardingManager.GetAllConnections()
}

// CloseAll はすべての接続をクローズ
func (gm *GroupManager) CloseAll() error {
	var lastErr error

	if err := gm.masterManager.CloseAll(); err != nil {
		lastErr = err
	}

	if err := gm.shardingManager.CloseAll(); err != nil {
		lastErr = err
	}

	return lastErr
}

// PingAll はすべての接続を確認
func (gm *GroupManager) PingAll() error {
	if err := gm.masterManager.PingAll(); err != nil {
		return fmt.Errorf("master group ping failed: %w", err)
	}

	if err := gm.shardingManager.PingAll(); err != nil {
		return fmt.Errorf("sharding group ping failed: %w", err)
	}

	return nil
}

// =============================================================================
// MasterManager
// =============================================================================

// MasterManager はmasterグループの接続を管理
type MasterManager struct {
	connection *GORMConnection
	mu         sync.RWMutex
}

// NewMasterManager は新しいMasterManagerを作成
func NewMasterManager(cfg *config.Config) (*MasterManager, error) {
	if len(cfg.Database.Groups.Master) == 0 {
		return nil, fmt.Errorf("master group configuration not found")
	}

	masterCfg := cfg.Database.Groups.Master[0]

	// SQL Loggerの作成
	sqlLogger, err := NewSQLLogger(
		1, // masterはID=1
		"master",
		masterCfg.Driver,
		cfg.Logging.SQLLogOutputDir,
		cfg.Logging.SQLLogEnabled,
	)
	if err != nil {
		log.Printf("Warning: Failed to create SQL logger for master: %v", err)
	}

	conn, err := NewGORMConnection(&masterCfg, sqlLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create master connection: %w", err)
	}

	return &MasterManager{
		connection: conn,
	}, nil
}

// GetConnection はmasterグループの接続を取得
func (mm *MasterManager) GetConnection() (*GORMConnection, error) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	if mm.connection == nil {
		return nil, fmt.Errorf("master connection not initialized")
	}

	return mm.connection, nil
}

// CloseAll はすべての接続をクローズ
func (mm *MasterManager) CloseAll() error {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if mm.connection != nil {
		return mm.connection.Close()
	}

	return nil
}

// PingAll はすべての接続を確認
func (mm *MasterManager) PingAll() error {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	if mm.connection != nil {
		return mm.connection.Ping()
	}

	return fmt.Errorf("master connection not initialized")
}

// =============================================================================
// ShardingManager
// =============================================================================

// ShardingManager はshardingグループの接続を管理
type ShardingManager struct {
	connections map[int]*GORMConnection // DB ID -> Connection
	tableRange  map[int][2]int          // DB ID -> [min, max]
	mu          sync.RWMutex
}

// NewShardingManager は新しいShardingManagerを作成
func NewShardingManager(cfg *config.Config) (*ShardingManager, error) {
	shardingCfg := cfg.Database.Groups.Sharding

	manager := &ShardingManager{
		connections: make(map[int]*GORMConnection),
		tableRange:  make(map[int][2]int),
	}

	// 各データベースへの接続を確立
	for _, dbCfg := range shardingCfg.Databases {
		// SQL Loggerの作成
		sqlLogger, err := NewSQLLogger(
			dbCfg.ID,
			"sharding",
			dbCfg.Driver,
			cfg.Logging.SQLLogOutputDir,
			cfg.Logging.SQLLogEnabled,
		)
		if err != nil {
			log.Printf("Warning: Failed to create SQL logger for sharding DB %d: %v", dbCfg.ID, err)
		}

		conn, err := NewGORMConnection(&dbCfg, sqlLogger)
		if err != nil {
			manager.CloseAll()
			return nil, fmt.Errorf("failed to create connection for sharding DB %d: %w", dbCfg.ID, err)
		}

		manager.connections[dbCfg.ID] = conn
		manager.tableRange[dbCfg.ID] = dbCfg.TableRange
	}

	return manager, nil
}

// GetConnectionByTableNumber はテーブル番号から接続を取得
func (sm *ShardingManager) GetConnectionByTableNumber(tableNumber int) (*GORMConnection, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// テーブル番号が範囲内か確認
	if tableNumber < 0 || tableNumber >= 32 {
		return nil, fmt.Errorf("invalid table number: %d (must be 0-31)", tableNumber)
	}

	// テーブル番号からデータベースIDを決定
	dbID := (tableNumber / 8) + 1

	conn, exists := sm.connections[dbID]
	if !exists {
		return nil, fmt.Errorf("connection for sharding DB %d not found", dbID)
	}

	// テーブル番号がデータベースの範囲内か確認
	tableRange, exists := sm.tableRange[dbID]
	if !exists {
		return nil, fmt.Errorf("table range for sharding DB %d not found", dbID)
	}

	if tableNumber < tableRange[0] || tableNumber > tableRange[1] {
		return nil, fmt.Errorf("table number %d is out of range for DB %d (range: %d-%d)",
			tableNumber, dbID, tableRange[0], tableRange[1])
	}

	return conn, nil
}

// GetAllConnections はすべての接続を取得（クロステーブルクエリ用）
func (sm *ShardingManager) GetAllConnections() []*GORMConnection {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	conns := make([]*GORMConnection, 0, len(sm.connections))
	for _, conn := range sm.connections {
		conns = append(conns, conn)
	}

	return conns
}

// CloseAll はすべての接続をクローズ
func (sm *ShardingManager) CloseAll() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var lastErr error
	for dbID, conn := range sm.connections {
		if err := conn.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close sharding DB %d: %w", dbID, err)
		}
	}

	return lastErr
}

// PingAll はすべての接続を確認
func (sm *ShardingManager) PingAll() error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	for dbID, conn := range sm.connections {
		if err := conn.Ping(); err != nil {
			return fmt.Errorf("failed to ping sharding DB %d: %w", dbID, err)
		}
	}

	return nil
}
