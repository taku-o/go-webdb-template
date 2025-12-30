package db

import (
	"fmt"
	"log"
	"sync"

	"github.com/taku-o/go-webdb-template/internal/config"
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
	tableNumber := int(id % DBShardingTableCount)
	return gm.GetShardingConnection(tableNumber)
}

// GetShardingConnectionByUUID はUUIDからshardingグループの接続を取得
func (gm *GroupManager) GetShardingConnectionByUUID(uuid string, tableName string) (*GORMConnection, error) {
	selector := NewTableSelector(DBShardingTableCount, DBShardingTablesPerDB)
	tableNumber, err := selector.GetTableNumberFromUUID(uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get table number from UUID: %w", err)
	}
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
//
// 8シャーディングエントリ構成:
// - 8つの論理的なシャーディングエントリ（各4テーブルを担当）
// - 4つの物理的なデータベース（各8テーブルを格納）
// - 同じDSNを持つエントリは接続を共有
//
// エントリ構成:
//   Entry 1,2 → sharding_db_1.db (tables 0-7)
//   Entry 3,4 → sharding_db_2.db (tables 8-15)
//   Entry 5,6 → sharding_db_3.db (tables 16-23)
//   Entry 7,8 → sharding_db_4.db (tables 24-31)
//
// =============================================================================

// ShardingManager はshardingグループの接続を管理
// 同じDSNを持つ複数のシャーディングエントリは、同じ接続オブジェクトを共有する
type ShardingManager struct {
	connections       map[int]*GORMConnection    // シャーディングエントリID -> Connection
	tableRange        map[int][2]int             // シャーディングエントリID -> [min, max]
	connectionPool    map[string]*GORMConnection // DSN -> Connection (共有用)
	tableNumberToDBID map[int]int                // テーブル番号 -> エントリID (O(1)ルックアップ用)
	mu                sync.RWMutex
}

// NewShardingManager は新しいShardingManagerを作成
// 同じDSNを持つ複数のエントリは、同じ接続オブジェクトを共有する
func NewShardingManager(cfg *config.Config) (*ShardingManager, error) {
	shardingCfg := cfg.Database.Groups.Sharding

	manager := &ShardingManager{
		connections:       make(map[int]*GORMConnection),
		tableRange:        make(map[int][2]int),
		connectionPool:    make(map[string]*GORMConnection),
		tableNumberToDBID: make(map[int]int),
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

		// DSNを取得（接続共有のキーとして使用）
		dsn := dbCfg.GetDSN()
		if dsn == "" {
			dsn = dbCfg.DSN
		}

		// 接続を取得または作成（共有対応）
		conn, err := manager.getOrCreateConnection(dsn, &dbCfg, sqlLogger)
		if err != nil {
			manager.CloseAll()
			return nil, fmt.Errorf("failed to create connection for sharding DB %d: %w", dbCfg.ID, err)
		}

		manager.connections[dbCfg.ID] = conn
		manager.tableRange[dbCfg.ID] = dbCfg.TableRange
	}

	// テーブル番号からエントリIDへのマッピングテーブルを構築（O(1)ルックアップ用）
	manager.buildTableNumberMap()

	return manager, nil
}

// getOrCreateConnection は接続を取得または作成（共有対応）
func (sm *ShardingManager) getOrCreateConnection(dsn string, cfg *config.ShardConfig, sqlLogger *SQLLogger) (*GORMConnection, error) {
	// 既存の接続を確認
	if conn, exists := sm.connectionPool[dsn]; exists {
		return conn, nil
	}

	// 新しい接続を作成
	conn, err := NewGORMConnection(cfg, sqlLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection for DSN %s: %w", dsn, err)
	}

	// プールに追加
	sm.connectionPool[dsn] = conn

	return conn, nil
}

// buildTableNumberMap はテーブル番号からエントリIDへのマッピングテーブルを構築
func (sm *ShardingManager) buildTableNumberMap() {
	sm.tableNumberToDBID = make(map[int]int)
	for dbID, tableRange := range sm.tableRange {
		for i := tableRange[0]; i <= tableRange[1]; i++ {
			sm.tableNumberToDBID[i] = dbID
		}
	}
}

// GetConnectionByTableNumber はテーブル番号から接続を取得（O(1)ルックアップ）
func (sm *ShardingManager) GetConnectionByTableNumber(tableNumber int) (*GORMConnection, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// テーブル番号が範囲内か確認
	if tableNumber < 0 || tableNumber >= DBShardingTableCount {
		return nil, fmt.Errorf("invalid table number: %d (must be 0-%d)", tableNumber, DBShardingTableCount-1)
	}

	// O(1)ルックアップ: テーブル番号からエントリIDを取得
	dbID, exists := sm.tableNumberToDBID[tableNumber]
	if !exists {
		return nil, fmt.Errorf("no connection found for table number %d", tableNumber)
	}

	// エントリIDから接続を取得
	conn, exists := sm.connections[dbID]
	if !exists {
		return nil, fmt.Errorf("connection for sharding DB %d not found", dbID)
	}

	return conn, nil
}

// GetAllConnections はすべてのユニークな接続を取得（クロステーブルクエリ用）
// 接続共有により、同じDSNを持つ複数のエントリは同じ接続を返すため、重複を排除する
func (sm *ShardingManager) GetAllConnections() []*GORMConnection {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// connectionPoolからユニークな接続のみを返す（重複を避ける）
	conns := make([]*GORMConnection, 0, len(sm.connectionPool))
	for _, conn := range sm.connectionPool {
		conns = append(conns, conn)
	}

	return conns
}

// CloseAll はすべての接続をクローズ
// 接続共有により、同じ接続を複数回クローズしないように注意する
func (sm *ShardingManager) CloseAll() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var lastErr error

	// connectionPoolの接続をクローズ（重複を避けるため）
	for dsn, conn := range sm.connectionPool {
		if err := conn.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close connection for DSN %s: %w", dsn, err)
		}
	}

	// マップをクリア
	sm.connections = make(map[int]*GORMConnection)
	sm.tableRange = make(map[int][2]int)
	sm.connectionPool = make(map[string]*GORMConnection)
	sm.tableNumberToDBID = make(map[int]int)

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
