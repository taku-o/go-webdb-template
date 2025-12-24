# GORM導入設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、GORM導入の詳細設計を定義する。既存のレイヤードアーキテクチャを維持しながら、GORM、dbresolver、shardingプラグインを統合する。

### 1.2 設計の範囲
- データベース接続層のGORM移行設計
- Writer/Reader分離のアーキテクチャ設計
- シャーディング統合の設計
- モデル定義のGORMタグ設計
- リポジトリ層のGORM API移行設計
- エラーハンドリング設計
- テスト戦略

## 2. アーキテクチャ設計

### 2.1 全体アーキテクチャ

既存のレイヤードアーキテクチャを維持し、DB層のみをGORMベースに変更する。

```
┌─────────────────────────────────────────────────────────────┐
│                         Client Layer                         │
│                      (Next.js 14 + React)                    │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP/REST
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      Server Layer (Go)                       │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │                    API Layer                            │ │
│  │  • HTTP Handlers (user_handler, post_handler)          │ │
│  └──────────────────────┬─────────────────────────────────┘ │
│                         │                                     │
│  ┌──────────────────────▼─────────────────────────────────┐ │
│  │                 Service Layer                           │ │
│  │  • Business logic                                       │ │
│  └──────────────────────┬─────────────────────────────────┘ │
│                         │                                     │
│  ┌──────────────────────▼─────────────────────────────────┐ │
│  │               Repository Layer                          │ │
│  │  • GORM API (Create, First, Find, Updates, Delete)      │ │
│  └──────────────────────┬─────────────────────────────────┘ │
│                         │                                     │
│  ┌──────────────────────▼─────────────────────────────────┐ │
│  │                  DB Layer (GORM)                        │ │
│  │  • GORM Manager (複数シャードの*gorm.DB管理)            │ │
│  │  • dbresolver Plugin (Writer/Reader分離)               │ │
│  │  • sharding Plugin (シャーディング)                     │ │
│  └──────────────────────┬─────────────────────────────────┘ │
└─────────────────────────┼───────────────────────────────────┘
                          │
         ┌────────────────┴────────────────┐
         ▼                                  ▼
    ┌─────────┐                        ┌─────────┐
    │ Shard 1 │                        │ Shard 2 │
    │ Writer  │                        │ Writer  │
    │ Reader  │                        │ Reader  │
    └─────────┘                        └─────────┘
```

### 2.2 DB層のアーキテクチャ変更

#### 2.2.1 変更前（database/sql）
```
Connection
  - *sql.DB
  - ShardID
  - Driver

Manager
  - map[int]*Connection
  - ShardingStrategy
```

#### 2.2.2 変更後（GORM）
```
GORMConnection
  - *gorm.DB (Writer/Reader分離済み)
  - ShardID
  - Driver

GORMManager
  - map[int]*GORMConnection
  - ShardingStrategy
```

### 2.3 Writer/Reader分離アーキテクチャ

各シャードごとにWriter/Reader接続を分離する。

```
Shard 1
  ├─ Writer: *gorm.DB (dbresolver設定済み)
  └─ Reader: *gorm.DB (dbresolver設定済み、複数可)

Shard 2
  ├─ Writer: *gorm.DB (dbresolver設定済み)
  └─ Reader: *gorm.DB (dbresolver設定済み、複数可)
```

**dbresolverプラグインの動作**:
- 読み取り操作（SELECT）: 自動的にReader接続を使用
- 書き込み操作（INSERT, UPDATE, DELETE）: 自動的にWriter接続を使用
- トランザクション: 常にWriter接続を使用

### 2.4 シャーディングアーキテクチャ

GORM Shardingプラグインを使用し、既存のHash-based sharding戦略を統合する。

```
Query with user_id (shard key)
  ↓
Hash-based Sharding Strategy
  ↓
Route to appropriate shard
  ↓
Execute on shard's GORM instance
```

**クロスシャードクエリ**:
- シャードキーを含まないクエリは全シャードに実行
- 結果をアプリケーションレベルでマージ

## 3. データモデル設計

### 3.1 Userモデル

```go
type User struct {
    ID        int64     `gorm:"primaryKey" json:"id,string"`
    Name      string    `gorm:"type:varchar(100);not null" json:"name"`
    Email     string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_users_email" json:"email"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
    return "users"
}
```

**GORMタグの説明**:
- `primaryKey`: 主キー指定
- `type:varchar(100)`: カラム型指定
- `not null`: NOT NULL制約
- `uniqueIndex:idx_users_email`: ユニークインデックス
- `autoCreateTime`: 作成時の自動タイムスタンプ
- `autoUpdateTime`: 更新時の自動タイムスタンプ

### 3.2 Postモデル

```go
type Post struct {
    ID        int64     `gorm:"primaryKey" json:"id,string"`
    UserID    int64     `gorm:"type:bigint;not null;index:idx_posts_user_id" json:"user_id,string"`
    Title     string    `gorm:"type:varchar(200);not null" json:"title"`
    Content   string    `gorm:"type:text;not null" json:"content"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Post) TableName() string {
    return "posts"
}
```

**GORMタグの説明**:
- `index:idx_posts_user_id`: インデックス指定（シャードキー用）
- 外部キー制約はGORMのリレーション機能を使用せず、アプリケーションレベルで管理

### 3.3 UserPostモデル（JOIN結果用）

```go
type UserPost struct {
    PostID      int64     `gorm:"column:post_id" json:"post_id,string"`
    PostTitle   string    `gorm:"column:post_title" json:"post_title"`
    PostContent string    `gorm:"column:post_content" json:"post_content"`
    UserID      int64     `gorm:"column:user_id" json:"user_id,string"`
    UserName    string    `gorm:"column:user_name" json:"user_name"`
    UserEmail   string    `gorm:"column:user_email" json:"user_email"`
    CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}

func (UserPost) TableName() string {
    return "" // テーブル名なし（JOIN結果用）
}
```

## 4. コンポーネント設計

### 4.1 GORMConnection構造体

```go
package db

import (
    "gorm.io/gorm"
    "github.com/example/go-webdb-template/internal/config"
)

// GORMConnection は単一のシャードのGORM接続を管理
type GORMConnection struct {
    DB       *gorm.DB  // dbresolver設定済みのGORMインスタンス
    ShardID  int
    Driver   string
    config   *config.ShardConfig
}

// NewGORMConnection は新しいGORM接続を作成
func NewGORMConnection(cfg *config.ShardConfig) (*GORMConnection, error) {
    // 1. Writer接続を作成
    writerDB, err := createGORMConnection(cfg, true)
    if err != nil {
        return nil, fmt.Errorf("failed to create writer connection: %w", err)
    }

    // 2. Reader接続を作成（複数可）
    readerDBs := make([]*gorm.DB, 0)
    for _, readerDSN := range cfg.ReaderDSNs {
        readerDB, err := createGORMConnectionFromDSN(readerDSN, cfg.Driver)
        if err != nil {
            return nil, fmt.Errorf("failed to create reader connection: %w", err)
        }
        readerDBs = append(readerDBs, readerDB)
    }

    // 3. dbresolverプラグインを設定
    if len(readerDBs) > 0 {
        err = writerDB.Use(dbresolver.Register(dbresolver.Config{
            Sources:  []gorm.Dialector{writerDB.Dialector},
            Replicas: readerDBs,
            Policy:   dbresolver.RandomPolicy(), // または RoundRobinPolicy
        }))
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
```

### 4.2 GORMManager構造体

```go
package db

import (
    "gorm.io/gorm"
    "sync"
    "github.com/example/go-webdb-template/internal/config"
)

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

    // 各シャードへの接続を確立
    for _, shardCfg := range cfg.Database.Shards {
        conn, err := NewGORMConnection(&shardCfg)
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
```

### 4.3 設定構造の拡張

```go
package config

// ShardConfig は各シャードの設定（拡張版）
type ShardConfig struct {
    ID                    int           `mapstructure:"id"`
    Driver                string        `mapstructure:"driver"`
    Host                  string        `mapstructure:"host"`
    Port                  int           `mapstructure:"port"`
    Name                  string        `mapstructure:"name"`
    User                  string        `mapstructure:"user"`
    Password              string        `mapstructure:"password"`
    DSN                   string        `mapstructure:"dsn"` // SQLite用のDSN
    MaxConnections        int           `mapstructure:"max_connections"`
    MaxIdleConnections    int           `mapstructure:"max_idle_connections"`
    ConnectionMaxLifetime time.Duration `mapstructure:"connection_max_lifetime"`
    
    // Writer/Reader分離用の設定（新規追加）
    WriterDSN             string        `mapstructure:"writer_dsn"` // Writer接続用DSN
    ReaderDSNs            []string      `mapstructure:"reader_dsns"` // Reader接続用DSNリスト
    ReaderPolicy          string        `mapstructure:"reader_policy"` // "random" or "round_robin"
}

// GetWriterDSN はWriter接続用DSNを取得
func (s *ShardConfig) GetWriterDSN() string {
    if s.WriterDSN != "" {
        return s.WriterDSN
    }
    // 後方互換性: 既存のDSNをWriterとして使用
    return s.GetDSN()
}

// GetReaderDSNs はReader接続用DSNリストを取得
func (s *ShardConfig) GetReaderDSNs() []string {
    if len(s.ReaderDSNs) > 0 {
        return s.ReaderDSNs
    }
    // 後方互換性: Writerと同じDSNをReaderとして使用（開発環境用）
    return []string{s.GetWriterDSN()}
}
```

### 4.4 リポジトリ層の変更

#### 4.4.1 UserRepository

```go
package repository

import (
    "context"
    "gorm.io/gorm"
    "github.com/example/go-webdb-template/internal/model"
)

type UserRepository struct {
    dbManager *db.GORMManager
}

func NewUserRepository(dbManager *db.GORMManager) *UserRepository {
    return &UserRepository{
        dbManager: dbManager,
    }
}

// Create はユーザーを作成
func (r *UserRepository) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
    user := &model.User{
        Name:  req.Name,
        Email: req.Email,
    }

    // ID生成（タイムスタンプベース、既存ロジック維持）
    user.ID = time.Now().UnixNano()

    // シャードキーに基づいてGORMインスタンスを取得
    db, err := r.dbManager.GetGORMByKey(user.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to get database: %w", err)
    }

    // GORM APIで作成
    if err := db.WithContext(ctx).Create(user).Error; err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    return user, nil
}

// GetByID はIDでユーザーを取得
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
    db, err := r.dbManager.GetGORMByKey(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get database: %w", err)
    }

    var user model.User
    if err := db.WithContext(ctx).First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("user not found: %d", id)
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    return &user, nil
}

// List はすべてのユーザーを取得（クロスシャードクエリ）
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
    connections := r.dbManager.GetAllGORMConnections()
    users := make([]*model.User, 0)

    // 各シャードから並列にデータを取得
    for _, conn := range connections {
        var shardUsers []*model.User
        if err := conn.DB.WithContext(ctx).
            Order("id").
            Limit(limit).
            Offset(offset).
            Find(&shardUsers).Error; err != nil {
            return nil, fmt.Errorf("failed to query shard %d: %w", conn.ShardID, err)
        }
        users = append(users, shardUsers...)
    }

    return users, nil
}

// Update はユーザーを更新
func (r *UserRepository) Update(ctx context.Context, id int64, req *model.UpdateUserRequest) (*model.User, error) {
    db, err := r.dbManager.GetGORMByKey(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get database: %w", err)
    }

    updates := make(map[string]interface{})
    if req.Name != "" {
        updates["name"] = req.Name
    }
    if req.Email != "" {
        updates["email"] = req.Email
    }
    updates["updated_at"] = time.Now()

    result := db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(updates)
    if result.Error != nil {
        return nil, fmt.Errorf("failed to update user: %w", result.Error)
    }
    if result.RowsAffected == 0 {
        return nil, fmt.Errorf("user not found: %d", id)
    }

    return r.GetByID(ctx, id)
}

// Delete はユーザーを削除
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
    db, err := r.dbManager.GetGORMByKey(id)
    if err != nil {
        return fmt.Errorf("failed to get database: %w", err)
    }

    result := db.WithContext(ctx).Delete(&model.User{}, id)
    if result.Error != nil {
        return fmt.Errorf("failed to delete user: %w", result.Error)
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("user not found: %d", id)
    }

    return nil
}
```

#### 4.4.2 PostRepository

同様にGORM APIを使用して実装。主要な変更点：
- `Create`: `db.Create()`を使用
- `GetByID`: `db.First()`を使用
- `ListByUserID`: `db.Where("user_id = ?", userID).Find()`を使用
- `List`: クロスシャードクエリで`db.Find()`を使用
- `GetUserPosts`: クロスシャードJOIN（アプリケーションレベル）
- `Update`: `db.Model().Updates()`を使用
- `Delete`: `db.Delete()`を使用

## 5. インターフェース設計

### 5.1 既存インターフェースの維持

Repository層のインターフェースは変更しない。既存のメソッドシグネチャを維持し、内部実装のみをGORMに変更する。

```go
// 既存のインターフェース（変更なし）
type UserRepository interface {
    Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
    GetByID(ctx context.Context, id int64) (*model.User, error)
    List(ctx context.Context, limit, offset int) ([]*model.User, error)
    Update(ctx context.Context, id int64, req *model.UpdateUserRequest) (*model.User, error)
    Delete(ctx context.Context, id int64) error
}
```

### 5.2 DB層のインターフェース変更

既存の`Manager`インターフェースを`GORMManager`に置き換える。後方互換性のため、既存のメソッド名を維持する。

```go
// 既存メソッド（*sql.DBを返していた）を*gorm.DBを返すように変更
type GORMManager interface {
    GetGORM(shardID int) (*gorm.DB, error)
    GetGORMByKey(key int64) (*gorm.DB, error)
    GetAllGORMConnections() []*GORMConnection
    CloseAll() error
    PingAll() error
}
```

## 6. エラーハンドリング設計

### 6.1 GORMエラーの変換

GORMのエラーを既存のエラーハンドリングパターンに変換する。

```go
package repository

import (
    "errors"
    "gorm.io/gorm"
)

// convertGORMError はGORMエラーを既存のエラー形式に変換
func convertGORMError(err error, entityType string, id int64) error {
    if err == nil {
        return nil
    }

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return fmt.Errorf("%s not found: %d", entityType, id)
    }

    // その他のGORMエラーはそのまま返す（必要に応じてラップ）
    return fmt.Errorf("database error: %w", err)
}
```

### 6.2 エラーハンドリングの一貫性

既存のエラーハンドリングパターンを維持：
- Repository層: エラーを返す（ラップ可）
- Service層: エラーをコンテキスト付きでラップ
- API層: エラーをHTTPステータスコードに変換

## 7. 設定ファイル設計

### 7.1 develop.yaml（開発環境）

```yaml
database:
  shards:
    - id: 1
      driver: sqlite3
      dsn: ./data/shard1.db
      writer_dsn: ./data/shard1.db
      reader_dsns:
        - ./data/shard1.db  # 開発環境では同一DB
      reader_policy: random
      max_connections: 10
      max_idle_connections: 5
      connection_max_lifetime: 300s
    - id: 2
      driver: sqlite3
      dsn: ./data/shard2.db
      writer_dsn: ./data/shard2.db
      reader_dsns:
        - ./data/shard2.db
      reader_policy: random
      max_connections: 10
      max_idle_connections: 5
      connection_max_lifetime: 300s
```

### 7.2 production.yaml.example（本番環境）

```yaml
database:
  shards:
    - id: 1
      driver: postgres
      host: db-shard1-writer.example.com
      port: 5432
      name: app_shard1
      user: app_user
      password: ${DB_SHARD1_PASSWORD}
      writer_dsn: host=db-shard1-writer.example.com port=5432 user=app_user password=${DB_SHARD1_PASSWORD} dbname=app_shard1 sslmode=require
      reader_dsns:
        - host=db-shard1-reader1.example.com port=5432 user=app_user password=${DB_SHARD1_PASSWORD} dbname=app_shard1 sslmode=require
        - host=db-shard1-reader2.example.com port=5432 user=app_user password=${DB_SHARD1_PASSWORD} dbname=app_shard1 sslmode=require
      reader_policy: round_robin
      max_connections: 25
      max_idle_connections: 5
      connection_max_lifetime: 300s
```

## 8. テスト戦略

### 8.1 ユニットテスト

#### 8.1.1 Repository層のテスト

```go
package repository_test

import (
    "testing"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "github.com/example/go-webdb-template/internal/model"
    "github.com/example/go-webdb-template/internal/repository"
)

func TestUserRepository_Create(t *testing.T) {
    // GORMインスタンスを作成（インメモリSQLite）
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    
    // テーブル作成
    db.AutoMigrate(&model.User{})
    
    // Repository作成
    repo := repository.NewUserRepository(mockGORMManager(db))
    
    // テスト実行
    req := &model.CreateUserRequest{
        Name:  "Test User",
        Email: "test@example.com",
    }
    user, err := repo.Create(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
    assert.Equal(t, "Test User", user.Name)
}
```

### 8.2 統合テスト

#### 8.2.1 Writer/Reader分離のテスト

```go
func TestWriterReaderSeparation(t *testing.T) {
    // Writer/Reader分離が正しく動作することを確認
    // 1. 書き込み操作がWriterに送られることを確認
    // 2. 読み取り操作がReaderに送られることを確認
    // 3. トランザクションがWriterに送られることを確認
}
```

#### 8.2.2 シャーディングのテスト

```go
func TestSharding(t *testing.T) {
    // シャーディングが正しく動作することを確認
    // 1. シャードキーに基づくルーティングを確認
    // 2. クロスシャードクエリの動作を確認
    // 3. データが正しいシャードに保存されることを確認
}
```

### 8.3 E2Eテスト

既存のE2Eテストを維持し、GORM実装でも正常に動作することを確認。

## 9. 移行戦略

### 9.1 段階的移行

1. **Phase 1**: GORM本体の導入とモデル定義の更新
2. **Phase 2**: Repository層のGORM APIへの置き換え
3. **Phase 3**: Writer/Reader分離の実装
4. **Phase 4**: シャーディングプラグインの統合

### 9.2 後方互換性の維持

- 既存のAPIエンドポイントの動作を維持
- 既存のテストスイートが全てパスする
- 設定ファイルの後方互換性（既存設定でも動作）

## 10. パフォーマンス考慮事項

### 10.1 接続プール

GORMの接続プール設定を適切に設定：
- `SetMaxOpenConns`: 最大接続数
- `SetMaxIdleConns`: 最大アイドル接続数
- `SetConnMaxLifetime`: 接続の最大生存時間

### 10.2 N+1問題の回避

GORMの`Preload`や`Joins`を使用してN+1問題を回避。

### 10.3 クエリ最適化

- 必要なカラムのみを選択（`Select`）
- インデックスの適切な使用
- バッチ操作の使用（`CreateInBatches`）

## 11. セキュリティ考慮事項

### 11.1 SQLインジェクション

GORMは自動的にパラメータ化クエリを使用するため、SQLインジェクションのリスクは低い。

### 11.2 接続情報の管理

- パスワードは環境変数で管理
- DSNには機密情報を含めない（環境変数参照）

## 12. 参考実装例

### 12.1 GORM接続作成

```go
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
```

## 13. 制約事項と注意点

### 13.1 GORM Shardingプラグインの制約

- シャードキーはクエリに含める必要がある
- クロスシャードクエリはアプリケーションレベルで実装
- 分散トランザクションはサポートされない

### 13.2 dbresolverプラグインの制約

- トランザクションは常にWriter接続を使用
- 読み取り専用クエリのみReader接続を使用
- 書き込みクエリは常にWriter接続を使用

### 13.3 既存機能との統合

- 既存の`HashBasedSharding`戦略を維持
- 既存のクロスシャードJOIN実装を維持
- 既存のエラーハンドリングパターンを維持

