# データベース遅延接続・自動再接続機能設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、データベースへの遅延接続とダウン後の自動再接続を実現するための詳細設計を定義する。`database/sql`の標準機能を活用し、既存のアーキテクチャに統合する。

### 1.2 設計の範囲
- 起動時のDB接続確認処理の削除
- 接続作成時のPing削除
- 接続プール設定の確認と追加
- PostgreSQL環境の構築（Docker Compose設定）
- 遅延接続と自動再接続の動作確認
- DB接続エラー時のリトライ機能の実装
- エラーハンドリング設計
- テスト戦略

### 1.3 設計方針
- **`database/sql`の標準機能の活用**: 遅延接続と自動再接続は`database/sql`の標準機能を活用
- **既存パターンの遵守**: 既存のDocker Compose設定や起動スクリプトのパターンに従う
- **後方互換性の保持**: 既存のDB接続機能を壊さない
- **設定ファイルからの読み込み**: 接続プール設定は設定ファイルから読み込む
- **リトライ機能の実装**: `avast/retry-go`を使用して接続エラー時にリトライ

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
go-webdb-template/
├── docker-compose.mailpit.yml
├── docker-compose.metabase.yml
├── docker-compose.cloudbeaver.yml
├── scripts/
│   ├── start-mailpit.sh
│   ├── metabase-start.sh
│   └── cloudbeaver-start.sh
├── server/
│   ├── internal/
│   │   ├── db/
│   │   │   ├── connection.go        # Ping()呼び出しあり
│   │   │   ├── manager.go
│   │   │   └── group_manager.go
│   │   └── config/
│   │       └── config.go
│   └── cmd/
│       └── server/
│           └── main.go              # PingAll()呼び出しあり
└── ...
```

#### 2.1.2 変更後の構造
```
go-webdb-template/
├── docker-compose.postgres.yml     # 新規: PostgreSQL用Docker Compose設定
├── docker-compose.mailpit.yml
├── docker-compose.metabase.yml
├── docker-compose.cloudbeaver.yml
├── scripts/
│   ├── start-postgres.sh            # 新規: PostgreSQL起動スクリプト
│   ├── start-mailpit.sh
│   ├── metabase-start.sh
│   └── cloudbeaver-start.sh
├── server/
│   ├── internal/
│   │   ├── db/
│   │   │   ├── connection.go        # 変更: Ping()呼び出し削除、リトライ機能追加
│   │   │   ├── manager.go           # 変更: リトライ機能追加
│   │   │   └── group_manager.go     # 変更: リトライ機能追加
│   │   └── config/
│   │       └── config.go            # 確認: 接続プール設定の確認
│   └── cmd/
│       └── server/
│           └── main.go              # 変更: PingAll()呼び出し削除
└── ...
```

### 2.2 接続確立の実行フロー

#### 2.2.1 変更前のフロー（起動時に接続確認）
```
┌─────────────────────────────────────────────────────────────┐
│              1. アプリケーション起動                           │
│              server/cmd/server/main.go                        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              2. GroupManagerの初期化                          │
│              db.NewGroupManager()                            │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              3. 各接続の作成                                  │
│              NewGORMConnection()                            │
│              - sql.Open()で接続オブジェクト作成                │
│              - 接続プール設定                                │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. 起動時の接続確認（削除対象）                    │
│              groupManager.PingAll()                         │
│              - すべての接続に対してPing()を実行                │
│              - 失敗時はサーバーを終了                        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              5. サーバー起動                                  │
│              - DB接続が利用できない場合は起動失敗              │
└─────────────────────────────────────────────────────────────┘
```

#### 2.2.2 変更後のフロー（遅延接続）
```
┌─────────────────────────────────────────────────────────────┐
│              1. アプリケーション起動                           │
│              server/cmd/server/main.go                        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              2. GroupManagerの初期化                          │
│              db.NewGroupManager()                            │
│              - リトライ機能付きで接続作成                      │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              3. 各接続の作成（遅延接続）                       │
│              NewGORMConnection()                            │
│              - sql.Open()で接続オブジェクト作成                │
│              - 接続プール設定                                │
│              - Ping()は実行しない（遅延接続）                   │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. サーバー起動（DB接続不要）                      │
│              - DB接続が利用できない場合でも起動成功            │
│              - 警告ログを出力                                │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              5. 最初のクエリ実行時                            │
│              - この時点で接続が確立される（遅延接続）          │
│              - 接続エラー時はリトライ（最大3回）               │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 既存アーキテクチャとの統合

```
┌─────────────────────────────────────────────────────────────┐
│              GroupManager (internal/db)                     │
│              - MasterManager                                │
│              - ShardingManager                              │
│              - リトライ機能付き接続作成                       │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              GORMConnection (internal/db)                   │
│              - 遅延接続対応                                  │
│              - 接続プール設定済み                            │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              database/sql (標準ライブラリ)                    │
│              - 遅延接続（sql.Open()は接続を確立しない）        │
│              - 自動再接続（接続プール設定により実現）          │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              PostgreSQL / SQLite (データベース)              │
│              - PostgreSQL: 動作確認用（Docker Compose）      │
│              - SQLite: 既存の開発環境                        │
└─────────────────────────────────────────────────────────────┘
```

## 3. コンポーネント設計

### 3.1 起動時のDB接続確認処理の削除

#### 3.1.1 server/cmd/server/main.go の修正

**変更前**:
```go
// すべてのデータベースへの接続確認
if err := groupManager.PingAll(); err != nil {
    log.Fatalf("Failed to ping databases: %v", err)
}
log.Println("Successfully connected to all database groups")
```

**変更後**:
```go
// 起動時のDB接続確認は削除（遅延接続のため）
// 最初のクエリ実行時に接続が確立される
log.Println("Database connections will be established on first query execution (lazy connection)")
```

**設計ポイント**:
- `PingAll()`呼び出しを削除
- 警告ログを出力して、遅延接続であることを明示
- DB接続が利用できない場合でもサーバーを起動できる

### 3.2 接続作成時のPing削除

#### 3.2.1 server/internal/db/connection.go の修正

**変更前** (`NewConnection`関数):
```go
// 接続プールの設定
db.SetMaxOpenConns(cfg.MaxConnections)
db.SetMaxIdleConns(cfg.MaxIdleConnections)
db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)

// 接続確認
if err := db.Ping(); err != nil {
    db.Close()
    return nil, fmt.Errorf("failed to ping database: %w", err)
}
```

**変更後**:
```go
// 接続プールの設定
db.SetMaxOpenConns(cfg.MaxConnections)
db.SetMaxIdleConns(cfg.MaxIdleConnections)
db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)

// 接続確認は削除（遅延接続のため）
// sql.Open()は接続を確立せず、接続オブジェクトを作成するのみ
// 実際のクエリ実行時に接続が確立される
```

**設計ポイント**:
- `db.Ping()`呼び出しを削除
- `sql.Open()`は接続を確立しないため、エラーが発生しない
- 実際のクエリ実行時に接続が確立される

**注意**: `NewGORMConnection`関数では既にPingは実行されていないため、変更不要。

### 3.3 接続プール設定の確認と追加

#### 3.3.1 接続プール設定の確認

**現在の実装状況**:
- `NewConnection`関数: 接続プール設定が実装済み（50-52行目）
- `createGORMConnection`関数: 接続プール設定が実装済み（119-121行目）

**確認項目**:
1. 設定値が適切か確認
2. 設定ファイルから読み込まれているか確認
3. デフォルト値が適切か確認

#### 3.3.2 接続プール設定のデフォルト値

**設定ファイル** (`config/{env}/database.yaml`):
```yaml
database:
  groups:
    master:
      - id: 1
        driver: "sqlite3"
        dsn: "file:master.db?mode=rwc"
        max_connections: 25          # 最大同時接続数
        max_idle_connections: 5      # アイドル状態のコネクション数
        connection_max_lifetime: 1h   # 接続の最大有効期間
```

**デフォルト値の設定** (`server/internal/db/connection.go`):
```go
// 接続プール設定のデフォルト値
const (
    DefaultMaxConnections        = 25
    DefaultMaxIdleConnections    = 5
    DefaultConnectionMaxLifetime = 1 * time.Hour
)

// NewConnection は新しいDB接続を作成
func NewConnection(cfg *config.ShardConfig) (*Connection, error) {
    // ... DSNの取得処理 ...

    db, err := sql.Open(driver, dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // 接続プール設定（設定値が0以下の場合はデフォルト値を使用）
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

    // Ping()は削除（遅延接続のため）

    return &Connection{
        DB:      db,
        ShardID: cfg.ID,
        Driver:  driver,
        config:  cfg,
    }, nil
}
```

**設計ポイント**:
- 設定値が0以下の場合はデフォルト値を使用
- デフォルト値は定数で定義
- 接続プール設定により、接続の再利用と自動再接続が実現される

### 3.4 PostgreSQL環境の構築

#### 3.4.1 docker-compose.postgres.yml

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: postgres
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: webdb
      POSTGRES_PASSWORD: webdb
      POSTGRES_DB: webdb
    volumes:
      - ./postgres/data:/var/lib/postgresql/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U webdb"]
      interval: 10s
      timeout: 5s
      retries: 5
```

**設計ポイント**:
- PostgreSQL 15 Alpineイメージを使用（軽量）
- ポート5432を公開
- データ永続化: `./postgres/data`にマウント
- ヘルスチェックを実装
- 開発用途（本番・staging環境では別の方法でPostgreSQLを起動）

#### 3.4.2 scripts/start-postgres.sh

```bash
#!/bin/bash

# PostgreSQL起動スクリプト
# 使用方法: ./scripts/start-postgres.sh {start|stop}

SCRIPT_DIR=$(cd "$(dirname "$0")/.." && pwd)
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.postgres.yml"

case "$1" in
  start)
    echo "Starting PostgreSQL..."
    docker-compose -f "$COMPOSE_FILE" up -d
    echo "PostgreSQL started. Port: 5432"
    echo "Connection: postgresql://webdb:webdb@localhost:5432/webdb"
    ;;
  stop)
    echo "Stopping PostgreSQL..."
    docker-compose -f "$COMPOSE_FILE" down
    echo "PostgreSQL stopped."
    ;;
  *)
    echo "Usage: $0 {start|stop}"
    exit 1
    ;;
esac
```

**設計ポイント**:
- 既存の`start-mailpit.sh`と同じパターン
- `start`/`stop`コマンドをサポート
- 適切なフィードバックを提供
- 開発用途（本番・staging環境では別の方法でPostgreSQLを起動）

#### 3.4.3 PostgreSQL設定ファイルの追加

**設定ファイル** (`config/{env}/database.yaml`):
```yaml
database:
  groups:
    master:
      - id: 1
        driver: "postgres"
        host: "localhost"
        port: 5432
        user: "webdb"
        password: "webdb"
        name: "webdb"
        max_connections: 25
        max_idle_connections: 5
        connection_max_lifetime: 1h
```

**設計ポイント**:
- 既存のSQLite設定と併用できる
- 環境変数で切り替え可能
- PostgreSQL用のDSNが自動生成される

### 3.5 DB接続エラー時のリトライ機能

#### 3.5.1 リトライライブラリの導入

**依存関係の追加** (`go.mod`):
```go
require (
    github.com/avast/retry-go/v4 v4.5.0
)
```

#### 3.5.2 リトライ機能の実装

**server/internal/db/connection.go**:
```go
import (
    "github.com/avast/retry-go/v4"
    "database/sql"
)

// リトライ設定
const (
    MaxRetryAttempts = 3  // 最大3回（初回 + 2回のリトライ）
    RetryDelay       = 1 * time.Second
)

// NewConnection は新しいDB接続を作成（リトライ機能付き）
func NewConnection(cfg *config.ShardConfig) (*Connection, error) {
    // ... DSNの取得処理 ...

    var db *sql.DB
    var err error

    // リトライ機能付きで接続作成
    err = retry.Do(
        func() error {
            db, err = sql.Open(driver, dsn)
            if err != nil {
                return err
            }

            // 接続プール設定
            // ... 設定処理 ...

            // 接続確認（リトライ対象）
            // 注意: 遅延接続のため、実際の接続確認は最初のクエリ実行時に行われる
            // ここでは接続オブジェクトの作成のみを確認
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
```

**設計ポイント**:
- `avast/retry-go/v4`を使用
- 最大3回までリトライ（初回 + 2回のリトライ）
- リトライ間隔: 1秒
- リトライ時はログを出力
- すべてのリトライが失敗した場合はエラーを返す

**注意**: 遅延接続のため、`sql.Open()`自体はエラーを返さない。リトライは実際のクエリ実行時の接続エラーに対して行う。

#### 3.5.3 クエリ実行時のリトライ

**Repository層でのリトライ実装**（必要に応じて）:
```go
// クエリ実行時のリトライ（例）
func (r *Repository) QueryWithRetry(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    var rows *sql.Rows
    var err error

    err = retry.Do(
        func() error {
            rows, err = r.db.QueryContext(ctx, query, args...)
            if err != nil {
                // 接続エラーの場合はリトライ
                if isConnectionError(err) {
                    return err
                }
                // その他のエラーはリトライしない
                return retry.Unrecoverable(err)
            }
            return nil
        },
        retry.Attempts(MaxRetryAttempts),
        retry.Delay(RetryDelay),
    )

    return rows, err
}

// 接続エラーの判定
func isConnectionError(err error) bool {
    if err == nil {
        return false
    }
    // sql.ErrConnDone などの接続エラーを判定
    return errors.Is(err, sql.ErrConnDone) || 
           strings.Contains(err.Error(), "connection") ||
           strings.Contains(err.Error(), "network")
}
```

**設計ポイント**:
- クエリ実行時の接続エラーに対してリトライ
- 接続エラー以外のエラーはリトライしない
- リトライ回数と間隔は接続作成時と同じ

### 3.6 GroupManager のリトライ機能追加

#### 3.6.1 server/internal/db/group_manager.go の修正

**NewGroupManager**関数:
```go
// NewGroupManager は新しいGroupManagerを作成（リトライ機能付き）
func NewGroupManager(cfg *config.Config) (*GroupManager, error) {
    // MasterManagerの作成（リトライ機能付き）
    masterManager, err := NewMasterManager(cfg)
    if err != nil {
        // リトライ機能により、接続エラー時は自動的にリトライされる
        // すべてのリトライが失敗した場合のみエラーを返す
        return nil, fmt.Errorf("failed to create master manager: %w", err)
    }

    // ShardingManagerの作成（リトライ機能付き）
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
```

**設計ポイント**:
- `NewMasterManager`と`NewShardingManager`内でリトライ機能が動作
- 接続エラー時は自動的にリトライされる
- すべてのリトライが失敗した場合のみエラーを返す

## 4. データモデル

### 4.1 接続プール設定

**設定構造体** (`server/internal/config/config.go`):
```go
type ShardConfig struct {
    // ... 既存のフィールド ...
    
    MaxConnections        int           `mapstructure:"max_connections"`
    MaxIdleConnections    int           `mapstructure:"max_idle_connections"`
    ConnectionMaxLifetime time.Duration `mapstructure:"connection_max_lifetime"`
}
```

**デフォルト値**:
- `MaxConnections`: 25
- `MaxIdleConnections`: 5
- `ConnectionMaxLifetime`: 1時間

### 4.2 リトライ設定

**定数定義** (`server/internal/db/connection.go`):
```go
const (
    MaxRetryAttempts = 3
    RetryDelay       = 1 * time.Second
)
```

## 5. エラーハンドリング

### 5.1 起動時のエラーハンドリング

**server/cmd/server/main.go**:
```go
// GroupManagerの初期化
groupManager, err := db.NewGroupManager(cfg)
if err != nil {
    // リトライ機能により、接続エラー時は自動的にリトライされる
    // すべてのリトライが失敗した場合のみエラーを返す
    // ただし、起動時はDB接続が不要なため、警告ログを出力して続行
    log.Printf("WARNING: Failed to create group manager: %v", err)
    log.Printf("WARNING: Database connections will be retried on first query execution")
    // サーバーは起動を続行（遅延接続のため）
}
```

**設計ポイント**:
- 起動時の接続エラーは警告ログを出力して続行
- 最初のクエリ実行時に接続が確立される

### 5.2 クエリ実行時のエラーハンドリング

**Repository層でのエラーハンドリング**:
```go
// クエリ実行時の接続エラーはリトライされる
// すべてのリトライが失敗した場合はエラーを返す
rows, err := r.db.QueryContext(ctx, query, args...)
if err != nil {
    // 接続エラーの場合は、次のクエリ実行時に自動的に再接続される
    // （接続プール設定により実現）
    return nil, fmt.Errorf("query execution failed: %w", err)
}
```

**設計ポイント**:
- 接続エラー時は、次のクエリ実行時に自動的に再接続される
- 接続プール設定により、古い接続は破棄され、新しい接続が作成される

## 6. テスト戦略

### 6.1 単体テスト

#### 6.1.1 接続作成のテスト
- `NewConnection`関数のテスト
- `NewGORMConnection`関数のテスト
- リトライ機能のテスト

#### 6.1.2 接続プール設定のテスト
- デフォルト値のテスト
- 設定ファイルからの読み込みテスト
- 設定値が0以下の場合のデフォルト値適用テスト

### 6.2 統合テスト

#### 6.2.1 遅延接続のテスト
- PostgreSQL環境でAPIサーバーを起動（DB接続なし）
- 最初のクエリ実行時に接続が確立されることを確認

#### 6.2.2 自動再接続のテスト
- PostgreSQLを停止
- クエリ実行時にエラーが発生することを確認
- PostgreSQLを再起動
- 次のクエリ実行時に自動的に再接続されることを確認

#### 6.2.3 リトライ機能のテスト
- DB接続エラー時のリトライ動作を確認
- リトライ回数と間隔を確認

### 6.3 E2Eテスト

#### 6.3.1 動作確認テスト
- PostgreSQL環境での動作確認
- SQLite環境での動作確認（既存機能の回帰テスト）

## 7. 実装順序

1. **起動時のDB接続確認処理の削除**
   - `server/cmd/server/main.go`の修正
   - 警告ログの追加

2. **接続作成時のPing削除**
   - `server/internal/db/connection.go`の修正

3. **接続プール設定の確認と追加**
   - 設定値の確認
   - デフォルト値の追加

4. **PostgreSQL環境の構築**
   - `docker-compose.postgres.yml`の作成
   - `scripts/start-postgres.sh`の作成
   - 設定ファイルの追加

5. **リトライ機能の実装**
   - `avast/retry-go/v4`の導入
   - 接続作成時のリトライ実装
   - クエリ実行時のリトライ実装（必要に応じて）

6. **動作確認**
   - 遅延接続の確認
   - 自動再接続の確認
   - リトライ機能の確認

## 8. 注意事項

### 8.1 遅延接続の注意点
- `sql.Open()`は接続を確立しないため、エラーが発生しない
- 実際のクエリ実行時に接続が確立される
- 接続エラーは最初のクエリ実行時に検出される

### 8.2 自動再接続の注意点
- 接続プール設定により、古い接続は破棄され、新しい接続が作成される
- `SetConnMaxLifetime`により、古い接続は定期的に破棄される
- データベースが復旧した際に、次のクエリ実行時に自動的に再接続される

### 8.3 リトライ機能の注意点
- リトライは接続エラーに対してのみ実行される
- その他のエラー（SQL構文エラーなど）はリトライしない
- リトライ回数と間隔は適切に設定する

### 8.4 PostgreSQL環境の注意点
- PostgreSQL環境は動作確認用（開発用途）
- 本番・staging環境では別の方法でPostgreSQLを起動
- SQLite環境との併用が可能
