# シャーディング数8対応設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、シャーディング数を4から8に増やすための詳細設計を定義する。設定ファイル上では8つのシャーディングエントリを定義するが、実際のデータベース接続は4つのまま維持し、同じDSNを持つ複数のエントリが接続を共有する仕組みを実装する。

### 1.2 設計の範囲
- 設定ファイルに8つのシャーディングエントリを追加
- 接続管理ロジックの変更（table_rangeベースの接続選択）
- 接続共有の実装（同じDSNを持つ複数のエントリが接続を共有）
- 既存のデータベースファイルはそのまま使用
- 32分割のテーブル構造は維持

### 1.3 設計方針
- **接続の共有**: 同じDSNを持つ複数のシャーディングエントリは、同じ接続オブジェクトを共有する
- **table_rangeベースの選択**: テーブル番号から接続を選択する際は、初期化時に構築したマッピングテーブルを使用してO(1)で接続を取得する
- **既存コードへの影響を最小化**: 既存のAPIインターフェースは維持し、内部実装のみを変更する
- **パフォーマンス**: 接続選択はO(1)で実行される。初期化時にテーブル番号からエントリIDへのマッピングテーブルを構築し、接続選択時はマップルックアップでO(1)を実現する

## 2. アーキテクチャ設計

### 2.1 システム構成図

```
┌─────────────────────────────────────────────────────────┐
│              ShardingManager                            │
│                                                         │
│  ┌──────────────────────────────────────────────────┐  │
│  │  connections: map[int]*GORMConnection            │  │
│  │  - Key: シャーディングエントリID (1-8)          │  │
│  │  - Value: 接続オブジェクト（共有可能）          │  │
│  └──────────────────────────────────────────────────┘  │
│                                                         │
│  ┌──────────────────────────────────────────────────┐  │
│  │  tableRange: map[int][2]int                     │  │
│  │  - Key: シャーディングエントリID (1-8)          │  │
│  │  - Value: [min, max] テーブル番号範囲           │  │
│  └──────────────────────────────────────────────────┘  │
│                                                         │
│  ┌──────────────────────────────────────────────────┐  │
│  │  connectionPool: map[string]*GORMConnection     │  │
│  │  - Key: DSN文字列                                │  │
│  │  - Value: 接続オブジェクト（共有用）            │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
                    │
                    │ GetConnectionByTableNumber(tableNumber)
                    │
                    ▼
        ┌───────────────────────────┐
        │ table_rangeを確認         │
        │ 該当するエントリIDを取得  │
        └───────────┬───────────────┘
                    │
                    ▼
        ┌───────────────────────────┐
        │ connections[dbID]を返す   │
        └───────────────────────────┘
```

### 2.2 接続共有の仕組み

#### 2.2.1 接続プールの管理
同じDSNを持つ複数のシャーディングエントリは、同じ接続オブジェクトを共有する：

```
設定ファイル:
  - id: 1, dsn: ./data/sharding_db_1.db, table_range: [0, 3]
  - id: 2, dsn: ./data/sharding_db_1.db, table_range: [4, 7]

接続プール:
  - "./data/sharding_db_1.db" → GORMConnection (共有)

connections map:
  - 1 → GORMConnection (共有接続への参照)
  - 2 → GORMConnection (共有接続への参照)
```

#### 2.2.2 接続の初期化フロー
1. 設定ファイルから8つのシャーディングエントリを読み込む
2. 各エントリのDSNを確認
3. 同じDSNを持つエントリが既に存在する場合は、既存の接続を再利用
4. 新しいDSNの場合は、新しい接続を作成してプールに追加
5. 各エントリIDと接続オブジェクトのマッピングを`connections`に保存

### 2.3 接続選択ロジック

#### 2.3.1 テーブル番号から接続を取得
```go
func (sm *ShardingManager) GetConnectionByTableNumber(tableNumber int) (*GORMConnection, error) {
    // 1. テーブル番号の範囲チェック
    if tableNumber < 0 || tableNumber >= 32 {
        return nil, fmt.Errorf("invalid table number: %d", tableNumber)
    }
    
    // 2. O(1)ルックアップ: テーブル番号からエントリIDを取得
    dbID, exists := sm.tableNumberToDBID[tableNumber]
    if !exists {
        return nil, fmt.Errorf("no connection found for table number %d", tableNumber)
    }
    
    // 3. エントリIDから接続を取得
    conn, exists := sm.connections[dbID]
    if !exists {
        return nil, fmt.Errorf("connection for sharding DB %d not found", dbID)
    }
    
    return conn, nil
}
```

#### 2.3.2 パフォーマンス考慮
- **実装**: O(1)でマップルックアップ（初期化時にマッピングテーブルを構築）
- **初期化時**: `buildTableNumberMap()`で32エントリのマッピングテーブルを構築（O(32) = O(1)）
- **接続選択時**: `tableNumberToDBID[tableNumber]`でO(1)ルックアップ
- **メモリ使用量**: 32エントリのマップ（約256バイト程度）

## 3. データ構造設計

### 3.1 ShardingManagerのデータ構造

#### 3.1.1 変更前
```go
type ShardingManager struct {
    connections map[int]*GORMConnection // DB ID -> Connection
    tableRange  map[int][2]int       // DB ID -> [min, max]
    mu          sync.RWMutex
}
```

#### 3.1.2 変更後
```go
type ShardingManager struct {
    connections       map[int]*GORMConnection // シャーディングエントリID -> Connection
    tableRange        map[int][2]int          // シャーディングエントリID -> [min, max]
    connectionPool    map[string]*GORMConnection // DSN -> Connection (共有用)
    tableNumberToDBID map[int]int             // テーブル番号 -> エントリID (O(1)ルックアップ用)
    mu                sync.RWMutex
}
```

### 3.2 接続プールの管理

#### 3.2.1 connectionPoolの役割
- **Key**: DSN文字列（例: `./data/sharding_db_1.db`）
- **Value**: 接続オブジェクト（共有用）
- **目的**: 同じDSNを持つ複数のエントリが接続を共有する

#### 3.2.2 接続の共有ロジック
```go
func (sm *ShardingManager) getOrCreateConnection(dsn string, cfg *config.ShardConfig, sqlLogger *SQLLogger) (*GORMConnection, error) {
    // 1. 既存の接続を確認
    if conn, exists := sm.connectionPool[dsn]; exists {
        return conn, nil
    }
    
    // 2. 新しい接続を作成
    conn, err := NewGORMConnection(cfg, sqlLogger)
    if err != nil {
        return nil, err
    }
    
    // 3. プールに追加
    sm.connectionPool[dsn] = conn
    
    return conn, nil
}
```

### 3.3 設定ファイルの構造

#### 3.3.1 設定ファイル例
```yaml
database:
  groups:
    sharding:
      databases:
        - id: 1
          driver: sqlite3
          dsn: ./data/sharding_db_1.db
          table_range: [0, 3]
        - id: 2
          driver: sqlite3
          dsn: ./data/sharding_db_1.db  # 同じDSN
          table_range: [4, 7]
        # ... 残り6つのエントリ
```

## 4. 実装詳細設計

### 4.1 NewShardingManagerの変更

#### 4.1.1 変更前
```go
func NewShardingManager(cfg *config.Config) (*ShardingManager, error) {
    manager := &ShardingManager{
        connections: make(map[int]*GORMConnection),
        tableRange:  make(map[int][2]int),
    }
    
    for _, dbCfg := range shardingCfg.Databases {
        conn, err := NewGORMConnection(&dbCfg, sqlLogger)
        if err != nil {
            return nil, err
        }
        
        manager.connections[dbCfg.ID] = conn
        manager.tableRange[dbCfg.ID] = dbCfg.TableRange
    }
    
    return manager, nil
}
```

#### 4.1.2 変更後
```go
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
```

### 4.2 getOrCreateConnectionメソッドの実装

```go
// getOrCreateConnection は接続を取得または作成（共有対応）
func (sm *ShardingManager) getOrCreateConnection(dsn string, cfg *config.ShardConfig, sqlLogger *SQLLogger) (*GORMConnection, error) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
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
```

### 4.3 GetConnectionByTableNumberの変更

#### 4.3.1 変更前
```go
func (sm *ShardingManager) GetConnectionByTableNumber(tableNumber int) (*GORMConnection, error) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    
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
```

#### 4.3.2 変更後（O(1)最適化版）
```go
func (sm *ShardingManager) GetConnectionByTableNumber(tableNumber int) (*GORMConnection, error) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    
    // テーブル番号が範囲内か確認
    if tableNumber < 0 || tableNumber >= 32 {
        return nil, fmt.Errorf("invalid table number: %d (must be 0-31)", tableNumber)
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
```

### 4.4 buildTableNumberMapメソッドの実装

```go
// buildTableNumberMap はテーブル番号からエントリIDへのマッピングテーブルを構築
func (sm *ShardingManager) buildTableNumberMap() {
    sm.tableNumberToDBID = make(map[int]int)
    for dbID, tableRange := range sm.tableRange {
        for i := tableRange[0]; i <= tableRange[1]; i++ {
            sm.tableNumberToDBID[i] = dbID
        }
    }
}
```

### 4.5 CloseAllメソッドの変更

#### 4.5.1 変更前
```go
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
```

#### 4.5.2 変更後
```go
func (sm *ShardingManager) CloseAll() error {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    var lastErr error
    
    // connectionPoolの接続をクローズ（重複を避けるため）
    closedDSNs := make(map[string]bool)
    for dsn, conn := range sm.connectionPool {
        if !closedDSNs[dsn] {
            if err := conn.Close(); err != nil {
                lastErr = fmt.Errorf("failed to close connection for DSN %s: %w", dsn, err)
            }
            closedDSNs[dsn] = true
        }
    }
    
    // マップをクリア
    sm.connections = make(map[int]*GORMConnection)
    sm.tableRange = make(map[int][2]int)
    sm.connectionPool = make(map[string]*GORMConnection)
    sm.tableNumberToDBID = make(map[int]int)
    
    return lastErr
}
```

### 4.6 GetAllConnectionsメソッドの変更

#### 4.6.1 変更前
```go
func (sm *ShardingManager) GetAllConnections() []*GORMConnection {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    
    conns := make([]*GORMConnection, 0, len(sm.connections))
    for _, conn := range sm.connections {
        conns = append(conns, conn)
    }
    
    return conns
}
```

#### 4.5.2 変更後
```go
func (sm *ShardingManager) GetAllConnections() []*GORMConnection {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    
    // connectionPoolからユニークな接続のみを返す（重複を避ける）
    conns := make([]*GORMConnection, 0, len(sm.connectionPool))
    seen := make(map[*GORMConnection]bool)
    
    for _, conn := range sm.connectionPool {
        if !seen[conn] {
            conns = append(conns, conn)
            seen[conn] = true
        }
    }
    
    return conns
}
```

## 5. エラーハンドリング

### 5.1 接続作成時のエラー

#### 5.1.1 エラーケース
- DSNが無効な場合
- データベース接続に失敗した場合
- 接続プールの作成に失敗した場合

#### 5.1.2 エラーハンドリング
```go
conn, err := manager.getOrCreateConnection(dsn, &dbCfg, sqlLogger)
if err != nil {
    // 既に作成された接続をクローズ
    manager.CloseAll()
    return nil, fmt.Errorf("failed to create connection for sharding DB %d: %w", dbCfg.ID, err)
}
```

### 5.2 接続選択時のエラー

#### 5.2.1 エラーケース
- テーブル番号が範囲外の場合
- 該当する接続が見つからない場合

#### 5.2.2 エラーハンドリング
```go
if tableNumber < 0 || tableNumber >= 32 {
    return nil, fmt.Errorf("invalid table number: %d (must be 0-31)", tableNumber)
}

// table_rangeを確認
for dbID, tableRange := range sm.tableRange {
    if tableNumber >= tableRange[0] && tableNumber <= tableRange[1] {
        conn, exists := sm.connections[dbID]
        if !exists {
            return nil, fmt.Errorf("connection for sharding DB %d not found", dbID)
        }
        return conn, nil
    }
}

return nil, fmt.Errorf("no connection found for table number %d", tableNumber)
```

### 5.3 接続クローズ時のエラー

#### 5.3.1 エラーケース
- 接続のクローズに失敗した場合

#### 5.3.2 エラーハンドリング
```go
var lastErr error
for dsn, conn := range sm.connectionPool {
    if !closedDSNs[dsn] {
        if err := conn.Close(); err != nil {
            lastErr = fmt.Errorf("failed to close connection for DSN %s: %w", dsn, err)
        }
        closedDSNs[dsn] = true
    }
}
return lastErr
```

## 6. テスト戦略

### 6.1 単体テスト

#### 6.1.1 NewShardingManagerのテスト
- 8つのシャーディングエントリが正しく初期化されること
- 同じDSNを持つ複数のエントリが接続を共有すること
- 異なるDSNを持つエントリが異なる接続を使用すること

#### 6.1.2 GetConnectionByTableNumberのテスト
- テーブル番号0-31が正しい接続を返すこと
- 範囲外のテーブル番号でエラーが返ること
- 各テーブル番号が正しいエントリにマッピングされること
- O(1)で接続選択が実行されること（パフォーマンステスト）

#### 6.1.3 buildTableNumberMapのテスト
- マッピングテーブルが正しく構築されること
- すべてのテーブル番号（0-31）がマッピングされること
- 重複するマッピングがないこと

#### 6.1.4 GetAllConnectionsのテスト
- ユニークな接続のみが返されること（重複がないこと）
- 4つの接続が返されること（実際のデータベース数）

### 6.2 統合テスト

#### 6.2.1 データベース操作のテスト
- 各テーブル番号に対してCRUD操作が正常に動作すること
- クロステーブルクエリが正常に動作すること

#### 6.2.2 接続共有のテスト
- 同じDSNを持つ複数のエントリが同じ接続を使用すること
- 接続のクローズが正しく動作すること

### 6.3 パフォーマンステスト

#### 6.3.1 接続選択のパフォーマンス
- 接続選択がO(1)で実行されること
- マッピングテーブルの構築がO(32)で実行されること
- シャーディング数が増えても接続選択のパフォーマンスが一定であること

## 7. 設定ファイルの更新

### 7.1 開発環境設定ファイル

#### 7.1.1 ファイルパス
- `config/develop/database.yaml`

#### 7.1.2 更新内容
- `sharding.databases`セクションを8つのエントリに拡張
- 各エントリの`table_range`を設定
- 同じDSNを持つ複数のエントリを定義

### 7.2 その他の環境設定ファイル

#### 7.2.1 更新対象
- `config/staging/database.yaml`（存在する場合）
- `config/production/database.yaml.example`（存在する場合）
- `server/internal/config/testdata/develop/database.yaml`

#### 7.2.2 更新内容
- 開発環境と同様に8つのエントリに拡張

## 8. ドキュメント更新

### 8.1 技術ドキュメント

#### 8.1.1 更新対象
- `docs/Sharding.md`

#### 8.1.2 更新内容
- シャーディング数の変更（4→8）を反映
- 接続共有の仕組みを説明
- 設定ファイルの構造を更新

### 8.2 コードコメント

#### 8.2.1 更新対象
- `server/internal/db/group_manager.go`

#### 8.2.2 更新内容
- 接続共有の仕組みをコメントで説明
- `table_range`ベースの接続選択を説明

## 9. 移行計画

### 9.1 実装順序

1. **設定ファイルの更新**
   - 開発環境設定ファイルを8つのエントリに拡張
   - テスト用設定ファイルを更新

2. **ShardingManagerの変更**
   - `connectionPool`フィールドを追加
   - `tableNumberToDBID`フィールドを追加（O(1)最適化用）
   - `getOrCreateConnection`メソッドを実装
   - `buildTableNumberMap`メソッドを実装
   - `NewShardingManager`を変更
   - `GetConnectionByTableNumber`を変更（O(1)最適化）
   - `CloseAll`を変更
   - `GetAllConnections`を変更

3. **テストの実装・更新**
   - 単体テストを実装
   - 統合テストを更新

4. **ドキュメントの更新**
   - `docs/Sharding.md`を更新

### 9.2 後方互換性

#### 9.2.1 既存APIの維持
- `GetShardingConnectionByID`メソッドは変更なし
- `GetAllShardingConnections`メソッドは変更なし（内部実装のみ変更）

#### 9.2.2 既存テストの動作
- 既存のテストは可能な限り動作する
- 大幅な変更が必要な場合はテストコードの更新も許容

## 10. 将来の拡張性

### 10.1 パフォーマンス最適化

#### 10.1.1 接続選択の最適化
本実装では、初期化時にテーブル番号からエントリIDへのマッピングテーブル（`tableNumberToDBID`）を構築し、接続選択をO(1)で実行します。

- **初期化時**: `buildTableNumberMap()`で32エントリのマッピングテーブルを構築（O(32) = O(1)）
- **接続選択時**: `tableNumberToDBID[tableNumber]`でO(1)ルックアップ
- **メモリ使用量**: 32エントリのマップ（約256バイト程度）

この最適化により、シャーディング数が増えても接続選択のパフォーマンスは一定です。

### 10.2 シャーディング数の拡張

#### 10.2.1 16シャーディングへの拡張
将来的に16シャーディングに拡張する場合も、同様の設計で対応可能：
- 設定ファイルに16つのエントリを追加
- 接続共有の仕組みはそのまま使用可能

### 10.3 データベース数の拡張

#### 10.3.1 8データベースへの拡張
将来的に8つのデータベースファイルに分割する場合：
- 設定ファイルで各エントリが異なるDSNを参照するように変更
- 接続共有の仕組みは自動的に適応

## 11. 参考情報

### 11.1 関連ドキュメント
- `.kiro/specs/0024-sharding8/requirements.md`: 要件定義書
- `.kiro/specs/0012-sharding/design.md`: シャーディング規則修正の設計書
- `docs/Sharding.md`: シャーディングの詳細仕様

### 11.2 既存実装
- `server/internal/db/group_manager.go`: 既存の接続管理
- `server/internal/db/connection.go`: 接続作成の実装
- `config/develop/database.yaml`: 既存の設定ファイル
