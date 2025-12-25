# シャーディング規則修正設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、masterグループとshardingグループの2つのデータベースグループを導入し、テーブル単位での32分割シャーディングを実現するための詳細設計を定義する。既存のアーキテクチャを拡張し、より柔軟なデータ分散戦略を実現する。

### 1.2 設計の範囲
- データベースグループ管理の設計
- 設定ファイル構造の拡張設計
- グループ別接続管理の設計
- テーブル選択ロジックの設計
- テンプレートベースのマイグレーション管理システムの設計
- Repository層の変更設計
- データモデル設計（newsテーブル）
- GoAdmin管理画面のnewsデータ参照ページ設計
- エラーハンドリング設計
- テスト戦略

### 1.3 設計方針
- **既存アーキテクチャの拡張**: 既存のレイヤードアーキテクチャとGORM接続管理を維持しつつ拡張
- **設定ファイルベース**: データベースグループの定義は設定ファイルで管理
- **テンプレートベースのマイグレーション**: 1つのテンプレートから32個のテーブル定義を生成
- **後方互換性の維持**: 既存のAPIエンドポイントの動作は維持
- **パフォーマンス重視**: テーブル選択はO(1)で実行、クロステーブルクエリは並列実行

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
config/
├── develop/
│   └── database.yaml      # shards配列のみ
├── staging/
│   └── database.yaml      # shards配列のみ
└── production/
    └── database.yaml.example  # shards配列のみ

db/
└── migrations/
    ├── shard1/
    │   └── 001_init.sql
    ├── shard2/
    │   └── 001_init.sql
    ├── shard3/
    │   └── 001_init.sql
    └── shard4/
        └── 001_init.sql

server/
└── internal/
    ├── config/
    │   └── config.go      # DatabaseConfig構造体
    ├── db/
    │   ├── manager.go      # Manager/GORMManager
    │   ├── connection.go   # Connection/GORMConnection
    │   └── sharding.go     # HashBasedSharding
    └── repository/
        ├── user_repository.go
        └── post_repository.go
```

#### 2.1.2 変更後の構造
```
config/
├── develop/
│   └── database.yaml      # groups構造に変更
├── staging/
│   └── database.yaml      # groups構造に変更
└── production/
    └── database.yaml.example  # groups構造に変更

db/
└── migrations/
    ├── master/
    │   └── 001_init.sql          # newsテーブル
    └── sharding/
        ├── templates/
        │   ├── users.sql.template    # usersテーブルのテンプレート
        │   └── posts.sql.template     # postsテーブルのテンプレート
        └── generated/              # 生成されたマイグレーション（オプション）

server/
└── internal/
    ├── config/
    │   └── config.go      # DatabaseConfig構造体を拡張
    ├── db/
    │   ├── manager.go      # GroupManager追加
    │   ├── connection.go   # 変更不要
    │   ├── sharding.go     # TableSelector追加
    │   └── group_manager.go  # 新規: グループ別接続管理
    ├── model/
    │   └── news.go         # 新規: newsモデル
    ├── admin/
    │   └── tables.go       # GetNewsTable関数を追加
    └── repository/
        ├── user_repository.go      # 動的テーブル名を使用
        └── post_repository.go      # 動的テーブル名を使用
```

### 2.2 データベースグループ管理の実行フロー

```
┌─────────────────────────────────────────────────────────────┐
│              1. アプリケーション起動                           │
│              server/cmd/server/main.go                        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              2. 設定ファイル読み込み                           │
│              config.Load()                                  │
│              - config/{env}/database.yaml を読み込み         │
│              - groups構造を解析                              │
│              - masterグループ: 1つのデータベース               │
│              - shardingグループ: 4つのデータベース            │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              3. GroupManager初期化                           │
│              NewGroupManager(cfg)                            │
│              - MasterManagerを作成                           │
│              - ShardingManagerを作成                         │
│              - 各グループの接続を確立                         │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. Repository層での使用                          │
│              - masterグループ: GetMasterConnection()          │
│              - shardingグループ: GetShardingConnection()      │
│              - テーブル名の動的生成: getTableName()           │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 テーブル選択の実行フロー

```
┌─────────────────────────────────────────────────────────────┐
│              1. Repository層でのリクエスト                    │
│              userRepo.GetByID(ctx, userID)                    │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              2. テーブル番号の計算                             │
│              tableNumber := userID % 32                       │
│              - 0-31の範囲                                   │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              3. テーブル名の生成                               │
│              tableName := fmt.Sprintf("users_%03d",          │
│                                       tableNumber)            │
│              - users_000, users_001, ..., users_031           │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. データベースの選択                             │
│              dbID := (tableNumber / 8) + 1                   │
│              - 1-4の範囲（各DBに8テーブル）                    │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              5. 接続の取得                                   │
│              conn := groupManager.GetShardingConnection(     │
│                          tableNumber)                        │
│              - sharding_db_1: _000-007                       │
│              - sharding_db_2: _008-015                      │
│              - sharding_db_3: _016-023                      │
│              - sharding_db_4: _024-031                       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              6. クエリの実行                                 │
│              query := fmt.Sprintf("SELECT * FROM %s         │
│                                    WHERE id = ?",            │
│                                    tableName)                │
│              - 動的なテーブル名を使用                         │
└─────────────────────────────────────────────────────────────┘
```

### 2.4 既存アーキテクチャとの統合

#### 2.4.1 設定構造体への統合
- `server/internal/config/config.go`の`DatabaseConfig`構造体を拡張
- `groups`フィールドを追加して、master/shardingグループを定義
- 既存の`shards`配列は後方互換性のために残す（非推奨）

#### 2.4.2 接続管理への統合
- `server/internal/db/manager.go`に`GroupManager`を追加
- 既存の`GORMManager`は後方互換性のために残す（非推奨）
- 新しいコードでは`GroupManager`を使用

#### 2.4.3 Repository層への統合
- `server/internal/repository/`の各Repositoryを拡張
- テーブル名の動的生成ロジックを追加
- 既存のAPIインターフェースは維持

## 3. データモデル設計

### 3.1 masterグループのテーブル

#### 3.1.1 newsテーブル
```sql
CREATE TABLE IF NOT EXISTS news (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author_id INTEGER,
    published_at DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_news_published_at ON news(published_at);
CREATE INDEX IF NOT EXISTS idx_news_author_id ON news(author_id);
```

**フィールド説明**:
- `id`: 主キー（AUTOINCREMENT）
- `title`: ニュースタイトル
- `content`: ニュース本文
- `author_id`: 作成者ID（オプション）
- `published_at`: 公開日時（オプション）
- `created_at`: 作成日時
- `updated_at`: 更新日時

### 3.2 shardingグループのテーブル

#### 3.2.1 usersテーブル（32分割）
各テーブル（users_000からusers_031）は同じスキーマを持つ：

```sql
-- テンプレート: users.sql.template
CREATE TABLE IF NOT EXISTS {TABLE_NAME} (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_{TABLE_NAME}_email ON {TABLE_NAME}(email);
```

**展開例**:
- `{TABLE_NAME}` → `users_000`, `users_001`, ..., `users_031`

#### 3.2.2 postsテーブル（32分割）
各テーブル（posts_000からposts_031）は同じスキーマを持つ：

```sql
-- テンプレート: posts.sql.template
CREATE TABLE IF NOT EXISTS {TABLE_NAME} (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users_{TABLE_SUFFIX}(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_{TABLE_NAME}_user_id ON {TABLE_NAME}(user_id);
CREATE INDEX IF NOT EXISTS idx_{TABLE_NAME}_created_at ON {TABLE_NAME}(created_at);
```

**展開例**:
- `{TABLE_NAME}` → `posts_000`, `posts_001`, ..., `posts_031`
- `{TABLE_SUFFIX}` → `000`, `001`, ..., `031`

**注意**: 外部キー制約は、同一データベース内のテーブル間でのみ有効。異なるデータベース間の外部キー制約はサポートされない。

### 3.3 テーブル分散マッピング

| データベース | テーブル番号範囲 | テーブル例 |
|------------|----------------|----------|
| sharding_db_1 | _000 〜 _007 | users_000, users_001, ..., users_007 |
| sharding_db_2 | _008 〜 _015 | users_008, users_009, ..., users_015 |
| sharding_db_3 | _016 〜 _023 | users_016, users_017, ..., users_023 |
| sharding_db_4 | _024 〜 _031 | users_024, users_025, ..., users_031 |

## 4. 実装設計

### 4.1 設定構造体の拡張

#### 4.1.1 DatabaseConfig構造体
```go
// server/internal/config/config.go

type DatabaseConfig struct {
    // 後方互換性のため残す（非推奨）
    Shards []ShardConfig `mapstructure:"shards"`
    
    // 新規: データベースグループ
    Groups DatabaseGroupsConfig `mapstructure:"groups"`
}

type DatabaseGroupsConfig struct {
    Master   []ShardConfig        `mapstructure:"master"`
    Sharding ShardingGroupConfig  `mapstructure:"sharding"`
}

type ShardingGroupConfig struct {
    Databases []ShardConfig       `mapstructure:"databases"`
    Tables    []ShardingTableConfig `mapstructure:"tables"`
}

type ShardingTableConfig struct {
    Name        string `mapstructure:"name"`         // テーブル名（例: "users"）
    SuffixCount int    `mapstructure:"suffix_count"` // 分割数（例: 32）
}

type ShardConfig struct {
    ID                    int           `mapstructure:"id"`
    Driver                string        `mapstructure:"driver"`
    Host                  string        `mapstructure:"host"`
    Port                  int           `mapstructure:"port"`
    Name                  string        `mapstructure:"name"`
    User                  string        `mapstructure:"user"`
    Password              string        `mapstructure:"password"`
    DSN                   string        `mapstructure:"dsn"`
    MaxConnections        int           `mapstructure:"max_connections"`
    MaxIdleConnections    int           `mapstructure:"max_idle_connections"`
    ConnectionMaxLifetime time.Duration `mapstructure:"connection_max_lifetime"`
    WriterDSN             string        `mapstructure:"writer_dsn"`
    ReaderDSNs            []string      `mapstructure:"reader_dsns"`
    ReaderPolicy          string        `mapstructure:"reader_policy"`
    
    // 新規: shardingグループ用
    TableRange            [2]int        `mapstructure:"table_range"` // [min, max]
}
```

#### 4.1.2 設定ファイル例
```yaml
# config/develop/database.yaml
database:
  groups:
    master:
      - id: 1
        driver: sqlite3
        dsn: ./data/master.db
        writer_dsn: ./data/master.db
        reader_dsns:
          - ./data/master.db
        reader_policy: random
        max_connections: 10
        max_idle_connections: 5
        connection_max_lifetime: 300s
    
    sharding:
      databases:
        - id: 1
          driver: sqlite3
          dsn: ./data/sharding_db_1.db
          writer_dsn: ./data/sharding_db_1.db
          reader_dsns:
            - ./data/sharding_db_1.db
          reader_policy: random
          max_connections: 10
          max_idle_connections: 5
          connection_max_lifetime: 300s
          table_range: [0, 7]  # _000-007
        - id: 2
          driver: sqlite3
          dsn: ./data/sharding_db_2.db
          writer_dsn: ./data/sharding_db_2.db
          reader_dsns:
            - ./data/sharding_db_2.db
          reader_policy: random
          max_connections: 10
          max_idle_connections: 5
          connection_max_lifetime: 300s
          table_range: [8, 15]  # _008-015
        - id: 3
          driver: sqlite3
          dsn: ./data/sharding_db_3.db
          writer_dsn: ./data/sharding_db_3.db
          reader_dsns:
            - ./data/sharding_db_3.db
          reader_policy: random
          max_connections: 10
          max_idle_connections: 5
          connection_max_lifetime: 300s
          table_range: [16, 23]  # _016-023
        - id: 4
          driver: sqlite3
          dsn: ./data/sharding_db_4.db
          writer_dsn: ./data/sharding_db_4.db
          reader_dsns:
            - ./data/sharding_db_4.db
          reader_policy: random
          max_connections: 10
          max_idle_connections: 5
          connection_max_lifetime: 300s
          table_range: [24, 31]  # _024-031
      
      tables:
        - name: users
          suffix_count: 32
        - name: posts
          suffix_count: 32
```

### 4.2 グループ別接続管理の実装

#### 4.2.1 GroupManager構造体
```go
// server/internal/db/group_manager.go

package db

import (
    "database/sql"
    "fmt"
    "sync"
    
    "github.com/example/go-webdb-template/internal/config"
    "gorm.io/gorm"
)

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
```

#### 4.2.2 MasterManager構造体
```go
// server/internal/db/group_manager.go

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
```

#### 4.2.3 ShardingManager構造体
```go
// server/internal/db/group_manager.go

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
```

### 4.3 テーブル選択ロジックの実装

#### 4.3.1 TableSelector構造体
```go
// server/internal/db/sharding.go

// TableSelector はテーブル選択ロジックを提供
type TableSelector struct {
    tableCount int  // 全テーブル数（デフォルト: 32）
    tablesPerDB int // データベースあたりのテーブル数（デフォルト: 8）
}

// NewTableSelector は新しいTableSelectorを作成
func NewTableSelector(tableCount, tablesPerDB int) *TableSelector {
    if tableCount <= 0 {
        tableCount = 32
    }
    if tablesPerDB <= 0 {
        tablesPerDB = 8
    }
    
    return &TableSelector{
        tableCount:  tableCount,
        tablesPerDB: tablesPerDB,
    }
}

// GetTableNumber はIDからテーブル番号を取得
func (ts *TableSelector) GetTableNumber(id int64) int {
    return int(id % int64(ts.tableCount))
}

// GetTableName はベース名とIDからテーブル名を生成
func (ts *TableSelector) GetTableName(baseName string, id int64) string {
    tableNumber := ts.GetTableNumber(id)
    return fmt.Sprintf("%s_%03d", baseName, tableNumber)
}

// GetDBID はテーブル番号からデータベースIDを取得
func (ts *TableSelector) GetDBID(tableNumber int) int {
    return (tableNumber / ts.tablesPerDB) + 1
}

// GetTableCount は全テーブル数を返す
func (ts *TableSelector) GetTableCount() int {
    return ts.tableCount
}
```

#### 4.3.2 テーブル名生成のユーティリティ関数
```go
// server/internal/db/table_selector.go（新規ファイル）

package db

import "fmt"

// GetShardingTableName はshardingグループのテーブル名を生成
func GetShardingTableName(baseName string, id int64) string {
    tableNumber := int(id % 32)
    return fmt.Sprintf("%s_%03d", baseName, tableNumber)
}

// GetShardingTableNumber はIDからテーブル番号を取得
func GetShardingTableNumber(id int64) int {
    return int(id % 32)
}

// GetShardingDBID はテーブル番号からデータベースIDを取得
func GetShardingDBID(tableNumber int) int {
    return (tableNumber / 8) + 1
}

// ValidateTableName はテーブル名が有効か検証（SQLインジェクション対策）
func ValidateTableName(tableName string, allowedBaseNames []string) bool {
    for _, baseName := range allowedBaseNames {
        // users_000, users_001, ..., users_031 の形式をチェック
        for i := 0; i < 32; i++ {
            expectedName := fmt.Sprintf("%s_%03d", baseName, i)
            if tableName == expectedName {
                return true
            }
        }
    }
    return false
}
```

### 4.4 Repository層の変更

#### 4.4.1 UserRepositoryの変更
```go
// server/internal/repository/user_repository.go

package repository

import (
    "context"
    "database/sql"
    "fmt"
    "time"
    
    "github.com/example/go-webdb-template/internal/db"
    "github.com/example/go-webdb-template/internal/model"
)

// UserRepository はユーザーのデータアクセスを担当
type UserRepository struct {
    groupManager *db.GroupManager
    tableSelector *db.TableSelector
}

// NewUserRepository は新しいUserRepositoryを作成
func NewUserRepository(groupManager *db.GroupManager) *UserRepository {
    return &UserRepository{
        groupManager:  groupManager,
        tableSelector: db.NewTableSelector(32, 8),
    }
}

// Create はユーザーを作成
func (r *UserRepository) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
    now := time.Now()
    user := &model.User{
        Name:      req.Name,
        Email:     req.Email,
        CreatedAt: now,
        UpdatedAt: now,
    }
    
    // 仮のIDを生成（実際のアプリケーションではID生成戦略を工夫）
    user.ID = now.UnixNano()
    
    // テーブル名の生成
    tableName := r.tableSelector.GetTableName("users", user.ID)
    
    // 接続の取得
    conn, err := r.groupManager.GetShardingConnectionByID(user.ID, "users")
    if err != nil {
        return nil, fmt.Errorf("failed to get sharding connection: %w", err)
    }
    
    // GORMを使用してクエリを実行
    query := fmt.Sprintf("INSERT INTO %s (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", tableName)
    
    sqlDB, err := conn.DB.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get sql.DB: %w", err)
    }
    
    _, err = sqlDB.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    return user, nil
}

// GetByID はIDでユーザーを取得
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
    // テーブル名の生成
    tableName := r.tableSelector.GetTableName("users", id)
    
    // 接続の取得
    conn, err := r.groupManager.GetShardingConnectionByID(id, "users")
    if err != nil {
        return nil, fmt.Errorf("failed to get sharding connection: %w", err)
    }
    
    query := fmt.Sprintf("SELECT id, name, email, created_at, updated_at FROM %s WHERE id = ?", tableName)
    
    sqlDB, err := conn.DB.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get sql.DB: %w", err)
    }
    
    var user model.User
    err = sqlDB.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user not found: %d", id)
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return &user, nil
}

// List はすべてのユーザーを取得（クロステーブルクエリ）
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
    connections := r.groupManager.GetShardingManager().GetAllConnections()
    users := make([]*model.User, 0)
    
    // 各データベースから並列にデータを取得
    for _, conn := range connections {
        // 各テーブル（users_000-031）からデータを取得
        for i := 0; i < 32; i++ {
            tableName := fmt.Sprintf("users_%03d", i)
            
            query := fmt.Sprintf("SELECT id, name, email, created_at, updated_at FROM %s ORDER BY id LIMIT ? OFFSET ?", tableName)
            
            sqlDB, err := conn.DB.DB()
            if err != nil {
                continue
            }
            
            rows, err := sqlDB.QueryContext(ctx, query, limit, offset)
            if err != nil {
                continue
            }
            defer rows.Close()
            
            for rows.Next() {
                var user model.User
                if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
                    continue
                }
                users = append(users, &user)
            }
        }
    }
    
    return users, nil
}
```

### 4.5 テンプレートベースのマイグレーション管理

#### 4.5.1 マイグレーションテンプレート
```sql
-- db/migrations/sharding/templates/users.sql.template
-- このテンプレートは32個のテーブル（users_000からusers_031）を生成する

-- Users テーブル（テーブル名は{TABLE_NAME}に置換される）
CREATE TABLE IF NOT EXISTS {TABLE_NAME} (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_{TABLE_NAME}_email ON {TABLE_NAME}(email);
```

#### 4.5.2 マイグレーション生成ツール
```go
// server/cmd/migrate-sharding/main.go（新規）

package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: migrate-sharding <template_file> <output_dir>")
        os.Exit(1)
    }
    
    templateFile := os.Args[1]
    outputDir := os.Args[2]
    
    // テンプレートファイルを読み込み
    templateContent, err := ioutil.ReadFile(templateFile)
    if err != nil {
        fmt.Printf("Error reading template file: %v\n", err)
        os.Exit(1)
    }
    
    // テーブル名のベースを取得（例: users.sql.template → users）
    baseName := strings.TrimSuffix(filepath.Base(templateFile), ".sql.template")
    
    // 32個のテーブル定義を生成
    for i := 0; i < 32; i++ {
        tableName := fmt.Sprintf("%s_%03d", baseName, i)
        
        // テンプレート内の{TABLE_NAME}を置換
        content := strings.ReplaceAll(string(templateContent), "{TABLE_NAME}", tableName)
        
        // 出力ファイル名を生成
        outputFile := filepath.Join(outputDir, fmt.Sprintf("%s_%03d.sql", baseName, i))
        
        // ファイルに書き込み
        if err := ioutil.WriteFile(outputFile, []byte(content), 0644); err != nil {
            fmt.Printf("Error writing file %s: %v\n", outputFile, err)
            continue
        }
        
        fmt.Printf("Generated: %s\n", outputFile)
    }
    
    fmt.Println("Migration files generated successfully")
}
```

#### 4.5.3 マイグレーション適用スクリプト
```bash
#!/bin/bash
# scripts/apply-sharding-migrations.sh

# テンプレートからマイグレーションファイルを生成
go run server/cmd/migrate-sharding/main.go \
    db/migrations/sharding/templates/users.sql.template \
    db/migrations/sharding/generated/

go run server/cmd/migrate-sharding/main.go \
    db/migrations/sharding/templates/posts.sql.template \
    db/migrations/sharding/generated/

# 各データベースに適切なテーブルを適用
# sharding_db_1: _000-007
for i in {0..7}; do
    sqlite3 server/data/sharding_db_1.db < db/migrations/sharding/generated/users_$(printf "%03d" $i).sql
    sqlite3 server/data/sharding_db_1.db < db/migrations/sharding/generated/posts_$(printf "%03d" $i).sql
done

# sharding_db_2: _008-015
for i in {8..15}; do
    sqlite3 server/data/sharding_db_2.db < db/migrations/sharding/generated/users_$(printf "%03d" $i).sql
    sqlite3 server/data/sharding_db_2.db < db/migrations/sharding/generated/posts_$(printf "%03d" $i).sql
done

# sharding_db_3: _016-023
for i in {16..23}; do
    sqlite3 server/data/sharding_db_3.db < db/migrations/sharding/generated/users_$(printf "%03d" $i).sql
    sqlite3 server/data/sharding_db_3.db < db/migrations/sharding/generated/posts_$(printf "%03d" $i).sql
done

# sharding_db_4: _024-031
for i in {24..31}; do
    sqlite3 server/data/sharding_db_4.db < db/migrations/sharding/generated/users_$(printf "%03d" $i).sql
    sqlite3 server/data/sharding_db_4.db < db/migrations/sharding/generated/posts_$(printf "%03d" $i).sql
done

echo "Sharding migrations applied successfully"
```

## 5. エラーハンドリング設計

### 5.1 接続エラー
- データベース接続失敗時は、エラーメッセージと共に適切なHTTPステータスコードを返す
- 接続プールのエラーはログに記録し、リトライロジックを実装

### 5.2 テーブル選択エラー
- 無効なテーブル番号（0-31の範囲外）の場合はエラーを返す
- テーブル名の検証（SQLインジェクション対策）を実装

### 5.3 マイグレーションエラー
- テンプレートファイルの読み込みエラーは明確なエラーメッセージを返す
- マイグレーション適用時のエラーはログに記録し、処理を中断

## 6. テスト戦略

### 6.1 単体テスト
- `TableSelector`のテーブル選択ロジックのテスト
- `GroupManager`の接続管理のテスト
- テーブル名生成ロジックのテスト

### 6.2 統合テスト
- 実際のデータベースを使用したRepository層のテスト
- クロステーブルクエリのテスト
- マイグレーション適用のテスト

### 6.3 E2Eテスト
- APIエンドポイントを通じたデータ操作のテスト
- 複数のテーブルにまたがる操作のテスト

## 7. パフォーマンス考慮事項

### 7.1 テーブル選択の最適化
- テーブル番号の計算はO(1)で実行
- テーブル名の生成は軽量な文字列操作

### 7.2 クロステーブルクエリの最適化
- 各データベースからのクエリは並列実行
- 結果のマージは効率的なアルゴリズムを使用

### 7.3 接続プールの管理
- グループ別に接続プールを管理
- 接続数の上限を適切に設定

## 8. セキュリティ考慮事項

### 8.1 SQLインジェクション対策
- テーブル名はホワイトリストで検証
- パラメータ化クエリを使用

### 8.2 接続情報の保護
- データベースパスワードは環境変数から読み込み
- 設定ファイルには機密情報を含めない

## 9. GoAdmin管理画面のnewsデータ参照ページ設計

### 9.1 GetNewsTable関数の実装
`server/internal/admin/tables.go`に`GetNewsTable`関数を追加：

```go
// GetNewsTable はNewsテーブルのGoAdmin設定を返す
func GetNewsTable(ctx *context.Context) table.Table {
    newsTable := table.NewDefaultTable(ctx, table.Config{
        Driver:     db.DriverSqlite,
        CanAdd:     true,
        Editable:   true,
        Deletable:  true,
        Exportable: true,
        Connection: table.DefaultConnectionName,
        PrimaryKey: table.PrimaryKey{
            Type: db.Int,
            Name: "id",
        },
    })

    // 一覧表示設定
    info := newsTable.GetInfo()
    info.AddField("ID", "id", db.Int).FieldSortable()
    info.AddField("タイトル", "title", db.Varchar).FieldSortable().FieldFilterable()
    info.AddField("内容", "content", db.Text)
    info.AddField("作成者ID", "author_id", db.Int).FieldSortable().FieldFilterable()
    info.AddField("公開日時", "published_at", db.Datetime).FieldSortable().FieldFilterable()
    info.AddField("作成日時", "created_at", db.Datetime).FieldSortable()
    info.AddField("更新日時", "updated_at", db.Datetime).FieldSortable()

    info.SetTable("news").SetTitle("ニュース").SetDescription("ニュース一覧")

    // フォーム設定（新規作成・編集）
    formList := newsTable.GetForm()
    formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowEdit().FieldNotAllowAdd()
    formList.AddField("タイトル", "title", db.Varchar, form.Text).FieldMust()
    formList.AddField("内容", "content", db.Text, form.TextArea).FieldMust()
    formList.AddField("作成者ID", "author_id", db.Int, form.Number)
    formList.AddField("公開日時", "published_at", db.Datetime, form.Datetime)
    formList.AddField("作成日時", "created_at", db.Datetime, form.Datetime).
        FieldHide().
        FieldNowWhenInsert().
        FieldDisableWhenUpdate()
    formList.AddField("更新日時", "updated_at", db.Datetime, form.Datetime).
        FieldHide().
        FieldNowWhenInsert().
        FieldNowWhenUpdate()

    formList.SetTable("news").SetTitle("ニュース").SetDescription("ニュース情報")

    return newsTable
}
```

### 9.2 テーブルジェネレータへの登録
`server/internal/admin/tables.go`の`Generators`マップに`news`を追加：

```go
var Generators = map[string]table.Generator{
    "users": GetUserTable,
    "posts": GetPostTable,
    "news":  GetNewsTable,  // 新規追加
}
```

### 9.3 ホームページへの統計情報追加（オプション）
`server/internal/admin/pages/home.go`の`HomePage`関数にnewsの統計情報を追加：

```go
newsCount := getTableCount(conn, "news")
// HTMLコンテンツにnews統計情報を追加
```

### 9.4 データベース接続の設定
GoAdminのデータベース接続設定で、masterグループのデータベースを使用するように設定：

```go
// server/internal/admin/config.go
func (c *Config) getDatabaseConfig() goadminConfig.DatabaseList {
    // masterグループのデータベースをGoAdmin用データベースとして使用
    dsn := ""
    if len(c.appConfig.Database.Groups.Master) > 0 {
        dsn = c.appConfig.Database.Groups.Master[0].DSN
    }
    
    return goadminConfig.DatabaseList{
        "default": {
            Driver: "sqlite",
            File:   dsn,
        },
    }
}
```

## 10. 将来の拡張性

### 10.1 テーブル数の変更
- テーブル数（32以外）に対応できる設計
- 設定ファイルでテーブル数を指定可能

### 10.2 データベース数の変更
- データベース数（4以外）に対応できる設計
- 設定ファイルでデータベース数を指定可能

### 10.3 新しいテーブルの追加
- テンプレートベースのマイグレーションで容易に追加可能
- 設定ファイルにテーブル定義を追加するだけで対応

