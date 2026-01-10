# テスト用データベース導入の実装タスク一覧

## 概要
テスト実行時に使用する専用のデータベースを導入するためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 設定ファイルの作成

#### タスク 1.1: `config/test/database.yaml`の作成
**目的**: テスト環境用のデータベース設定ファイルを作成する。

**作業内容**:
- `config/test/`ディレクトリを作成（存在しない場合）
- `config/develop/database.yaml`をベースに`config/test/database.yaml`を作成
- すべてのデータベース名に`_test`サフィックスを付与
- その他の設定（ホスト、ポート、ユーザー、パスワード等）は開発環境と同じ設定を使用

**実装内容**:
- マスターデータベース名: `webdb_master` → `webdb_master_test`
- シャーディングデータベース名: `webdb_sharding_1` → `webdb_sharding_1_test`（同様に2-4も）
- 8つのシャーディングエントリすべてに対応（id: 1-8）

**ファイル構成**:
```yaml
# PostgreSQL設定（テスト環境用）
# テスト用データベースは事前に作成し、マイグレーションを実行しておく必要がある
database:
  groups:
    master:
      - id: 1
        driver: postgres
        host: localhost
        port: 5432
        user: webdb
        password: webdb
        name: webdb_master_test
        max_connections: 25
        max_idle_connections: 5
        connection_max_lifetime: 1h

    sharding:
      databases:
        # 論理シャード 1: テーブル _000-003 → postgres-sharding-1
        - id: 1
          driver: postgres
          host: localhost
          port: 5433
          user: webdb
          password: webdb
          name: webdb_sharding_1_test
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [0, 3]
        # ... 同様に2-8も定義（webdb_sharding_1_test, webdb_sharding_2_test, webdb_sharding_3_test, webdb_sharding_4_test）

      tables:
        - name: dm_users
          suffix_count: 32
        - name: dm_posts
          suffix_count: 32
```

**受け入れ基準**:
- `config/test/database.yaml`が作成されている
- データベース名に`_test`サフィックスが付与されている
- マスターデータベース名が`webdb_master_test`である
- シャーディングデータベース名が`webdb_sharding_1_test` ~ `webdb_sharding_4_test`である（8つのエントリすべてに対応）
- その他の設定（ホスト、ポート、ユーザー等）が開発環境と同じである
- コメントが適切に追加されている

- _Requirements: 3.1, 6.1_
- _Design: 3.1.1_

---

#### タスク 1.2: `config/test/atlas.hcl`の作成
**目的**: テスト環境用のAtlas設定ファイルを作成する。

**作業内容**:
- `config/develop/atlas.hcl`をベースに`config/test/atlas.hcl`を作成
- 各環境（master, sharding_1 ~ sharding_4）のURLにテスト用データベース名を使用
- マイグレーションディレクトリは開発環境と同じ設定を使用

**実装内容**:
- マスターデータベースURL: `postgres://webdb:webdb@localhost:5432/webdb_master_test?sslmode=disable`
- シャーディングデータベースURL: `postgres://webdb:webdb@localhost:5433/webdb_sharding_1_test?sslmode=disable`（同様に2-4も）

**ファイル構成**:
```hcl
// テスト環境用Atlas設定ファイル (PostgreSQL)

// マスターデータベース用環境
env "master" {
  src = "file://db/schema/master.hcl"
  url = "postgres://webdb:webdb@localhost:5432/webdb_master_test?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://db/migrations/master"
  }
}

// シャーディングDB 1
env "sharding_1" {
  src = "file://db/schema/sharding_1"
  url = "postgres://webdb:webdb@localhost:5433/webdb_sharding_1_test?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://db/migrations/sharding_1"
  }
}

// シャーディングDB 2
env "sharding_2" {
  src = "file://db/schema/sharding_2"
  url = "postgres://webdb:webdb@localhost:5434/webdb_sharding_2_test?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://db/migrations/sharding_2"
  }
}

// シャーディングDB 3
env "sharding_3" {
  src = "file://db/schema/sharding_3"
  url = "postgres://webdb:webdb@localhost:5435/webdb_sharding_3_test?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://db/migrations/sharding_3"
  }
}

// シャーディングDB 4
env "sharding_4" {
  src = "file://db/schema/sharding_4"
  url = "postgres://webdb:webdb@localhost:5436/webdb_sharding_4_test?sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://db/migrations/sharding_4"
  }
}
```

**受け入れ基準**:
- `config/test/atlas.hcl`が作成されている
- 各環境のURLにテスト用データベース名が使用されている
- マイグレーションディレクトリが正しく設定されている
- すべての環境（master, sharding_1 ~ sharding_4）が定義されている

- _Requirements: 3.2, 6.2_
- _Design: 3.1.2_

---

#### タスク 1.3: テスト用マイグレーションスクリプトの作成
**目的**: テスト用データベースに対してマイグレーションを実行するスクリプトを作成する。

**作業内容**:
- `scripts/migrate-test.sh`を新規作成
- 既存の`scripts/migrate.sh`をベースに、テスト環境用のスクリプトを作成
- `APP_ENV=test`を設定してAtlasマイグレーションを実行
- `config/test/atlas.hcl`を使用してマイグレーションを適用

**実装内容**:
- `APP_ENV=test`を設定
- `config/test/atlas.hcl`を使用してAtlasマイグレーションを実行
- 既存の`migrate.sh`と同様の使い勝手を提供（`[master|sharding|all]`オプション対応）

**ファイル構成**:
```bash
#!/bin/bash
# PostgreSQL用マイグレーションスクリプト（テスト環境用）
# 使用方法: ./scripts/migrate-test.sh [master|sharding|all]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# テスト環境を設定
APP_ENV="test"
ATLAS_CONFIG="$PROJECT_ROOT/config/$APP_ENV/atlas.hcl"

# デフォルトのPostgreSQL接続情報（テスト環境用）
MASTER_HOST="localhost"
MASTER_PORT="5432"
MASTER_USER="webdb"
MASTER_PASSWORD="webdb"
MASTER_DB="webdb_master_test"

SHARDING_1_HOST="localhost"
SHARDING_1_PORT="5433"
SHARDING_1_USER="webdb"
SHARDING_1_PASSWORD="webdb"
SHARDING_1_DB="webdb_sharding_1_test"

SHARDING_2_HOST="localhost"
SHARDING_2_PORT="5434"
SHARDING_2_USER="webdb"
SHARDING_2_PASSWORD="webdb"
SHARDING_2_DB="webdb_sharding_2_test"

SHARDING_3_HOST="localhost"
SHARDING_3_PORT="5435"
SHARDING_3_USER="webdb"
SHARDING_3_PASSWORD="webdb"
SHARDING_3_DB="webdb_sharding_3_test"

SHARDING_4_HOST="localhost"
SHARDING_4_PORT="5436"
SHARDING_4_USER="webdb"
SHARDING_4_PASSWORD="webdb"
SHARDING_4_DB="webdb_sharding_4_test"

# PostgreSQL URL形式を構築
build_postgres_url() {
    local host=$1
    local port=$2
    local user=$3
    local password=$4
    local dbname=$5
    echo "postgres://${user}:${password}@${host}:${port}/${dbname}?sslmode=disable"
}

# 使用方法を表示
usage() {
    echo "Usage: $0 [master|sharding|all]"
    echo ""
    echo "Commands:"
    echo "  master    Apply migrations to master database only"
    echo "  sharding  Apply migrations to sharding databases only"
    echo "  all       Apply migrations to all databases (default)"
    echo ""
    echo "This script uses APP_ENV=test and config/test/atlas.hcl"
    exit 1
}

# マスターグループのマイグレーション
migrate_master() {
    echo "Migrating master database (test environment)..."
    local url=$(build_postgres_url "$MASTER_HOST" "$MASTER_PORT" "$MASTER_USER" "$MASTER_PASSWORD" "$MASTER_DB")

    # Atlasマイグレーション適用
    echo "  Applying Atlas migrations..."
    atlas migrate apply \
        --dir "file://$PROJECT_ROOT/db/migrations/master" \
        --url "$url"

    # Viewマイグレーション適用（生SQL）
    echo "  Applying View migrations..."
    local view_dir="$PROJECT_ROOT/db/migrations/view_master"
    if [ -d "$view_dir" ]; then
        for sql_file in $(ls "$view_dir"/*.sql 2>/dev/null | sort); do
            echo "    Applying $(basename "$sql_file")..."
            docker exec -i postgres-master psql -U "$MASTER_USER" -d "$MASTER_DB" < "$sql_file"
        done
    fi

    echo "Master database migration applied."
}

# シャーディンググループのマイグレーション
migrate_sharding() {
    echo "Migrating sharding databases (test environment)..."

    # Sharding 1
    echo "  Migrating sharding_1..."
    local url1=$(build_postgres_url "$SHARDING_1_HOST" "$SHARDING_1_PORT" "$SHARDING_1_USER" "$SHARDING_1_PASSWORD" "$SHARDING_1_DB")
    atlas migrate apply \
        --dir "file://$PROJECT_ROOT/db/migrations/sharding_1" \
        --url "$url1"

    # Sharding 2
    echo "  Migrating sharding_2..."
    local url2=$(build_postgres_url "$SHARDING_2_HOST" "$SHARDING_2_PORT" "$SHARDING_2_USER" "$SHARDING_2_PASSWORD" "$SHARDING_2_DB")
    atlas migrate apply \
        --dir "file://$PROJECT_ROOT/db/migrations/sharding_2" \
        --url "$url2"

    # Sharding 3
    echo "  Migrating sharding_3..."
    local url3=$(build_postgres_url "$SHARDING_3_HOST" "$SHARDING_3_PORT" "$SHARDING_3_USER" "$SHARDING_3_PASSWORD" "$SHARDING_3_DB")
    atlas migrate apply \
        --dir "file://$PROJECT_ROOT/db/migrations/sharding_3" \
        --url "$url3"

    # Sharding 4
    echo "  Migrating sharding_4..."
    local url4=$(build_postgres_url "$SHARDING_4_HOST" "$SHARDING_4_PORT" "$SHARDING_4_USER" "$SHARDING_4_PASSWORD" "$SHARDING_4_DB")
    atlas migrate apply \
        --dir "file://$PROJECT_ROOT/db/migrations/sharding_4" \
        --url "$url4"

    echo "Sharding databases migration applied."
}

# メイン処理
case "${1:-all}" in
    master)
        migrate_master
        ;;
    sharding)
        migrate_sharding
        ;;
    all)
        migrate_master
        migrate_sharding
        ;;
    -h|--help)
        usage
        ;;
    *)
        usage
        ;;
esac

echo "All migrations applied successfully!"
```

**注意事項**:
- `config/test/atlas.hcl`が存在することを前提とする
- 既存の`migrate.sh`と同様の使い勝手を提供する
- Viewマイグレーションは既存の`migrate.sh`と同様に実装するか、必要に応じて追加

**受け入れ基準**:
- `scripts/migrate-test.sh`が作成されている
- スクリプトが正常に実行できる
- テスト用データベースに対してマイグレーションが正常に適用される
- 既存の`migrate.sh`と同様の使い勝手が提供される

- _Requirements: 6.5_
- _Design: 4.3.2_

---

### Phase 2: 設定読み込み機能の実装

#### タスク 2.1: `LoadTestConfig()`関数の実装
**目的**: テスト環境の設定を読み込む関数を実装する。

**作業内容**:
- `server/test/testutil/db.go`に`LoadTestConfig()`関数を追加
- 既存の`APP_ENV`を保存し、`APP_ENV=test`を設定
- `config.Load()`を呼び出して設定を読み込む
- 関数終了時に元の`APP_ENV`を復元

**実装コード**:
```go
// LoadTestConfig はテスト環境の設定を読み込む
func LoadTestConfig() (*config.Config, error) {
    // 既存のAPP_ENVを保存
    oldEnv := os.Getenv("APP_ENV")
    defer func() {
        if oldEnv != "" {
            os.Setenv("APP_ENV", oldEnv)
        } else {
            os.Unsetenv("APP_ENV")
        }
    }()
    
    // テスト環境を設定
    os.Setenv("APP_ENV", "test")
    
    // 設定を読み込む
    return config.Load()
}
```

**受け入れ基準**:
- `LoadTestConfig()`関数が実装されている
- 関数が正常に設定を読み込む
- データベース名に`_test`サフィックスが付与されている
- 元の`APP_ENV`が適切に復元される
- エラーハンドリングが適切に実装されている

- _Requirements: 3.3.1, 6.3_
- _Design: 3.2.1_

---

### Phase 3: データベースクリア機能の実装

#### タスク 3.1: `clearDatabaseTables()`関数の実装
**目的**: 指定されたデータベースの全テーブルのデータをクリアする関数を実装する。

**作業内容**:
- `server/test/testutil/db.go`に`clearDatabaseTables()`関数を追加
- PostgreSQLの`pg_tables`システムカタログからテーブル一覧を取得
- 各テーブルに対して`TRUNCATE TABLE ... RESTART IDENTITY CASCADE`を実行

**実装コード**:
```go
// clearDatabaseTables は指定されたデータベースの全テーブルのデータをクリアする
func clearDatabaseTables(t *testing.T, database *gorm.DB) {
    // テーブル一覧を取得
    var tables []string
    err := database.Raw(`
        SELECT tablename 
        FROM pg_tables 
        WHERE schemaname = 'public'
    `).Scan(&tables).Error
    require.NoError(t, err)
    
    // 各テーブルをTRUNCATE
    for _, table := range tables {
        tableName := tables[table]
        err := database.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tableName)).Error
        require.NoError(t, err)
    }
}
```

**受け入れ基準**:
- `clearDatabaseTables()`関数が実装されている
- 関数が正常にテーブル一覧を取得する
- 関数が正常に各テーブルをTRUNCATEする
- エラーハンドリングが適切に実装されている

- _Requirements: 3.4, 6.4_
- _Design: 3.3.1_

---

#### タスク 3.2: `ClearTestDatabase()`関数の実装
**目的**: テスト用データベースの全テーブルのデータをクリアする関数を実装する。

**作業内容**:
- `server/test/testutil/db.go`に`ClearTestDatabase()`関数を追加
- マスターデータベースの接続を取得し、`clearDatabaseTables()`を呼び出す
- シャーディングデータベースの接続を取得し、各接続に対して`clearDatabaseTables()`を呼び出す

**実装コード**:
```go
// ClearTestDatabase はテスト用データベースの全テーブルのデータをクリアする
func ClearTestDatabase(t *testing.T, manager *db.GroupManager) {
    // マスターデータベースのクリア
    masterConn, err := manager.GetMasterConnection()
    require.NoError(t, err)
    clearDatabaseTables(t, masterConn.DB)
    
    // シャーディングデータベースのクリア
    connections := manager.GetAllShardingConnections()
    for _, conn := range connections {
        clearDatabaseTables(t, conn.DB)
    }
}
```

**受け入れ基準**:
- `ClearTestDatabase()`関数が実装されている
- 関数が正常にマスターデータベースをクリアする
- 関数が正常にシャーディングデータベースをクリアする
- エラーハンドリングが適切に実装されている

- _Requirements: 3.4, 6.4_
- _Design: 3.3.1_

---

### Phase 4: 既存関数の修正

#### タスク 4.1: `SetupTestGroupManager()`関数の修正
**目的**: ハードコードされた設定を削除し、設定ファイルから読み込むように修正する。

**作業内容**:
- `server/test/testutil/db.go`の`SetupTestGroupManager()`関数を修正
- ハードコードされた設定の作成部分を削除
- `LoadTestConfig()`を呼び出して設定を読み込む
- `db.NewGroupManager(cfg)`を使用してGroupManagerを作成
- `ClearTestDatabase()`を呼び出してデータベースをクリア
- 既存のスキーマ初期化処理は維持

**変更前のコード構造**:
```go
func SetupTestGroupManager(t *testing.T, dbCount int, tablesPerDB int) *db.GroupManager {
    // ハードコードされた設定を作成
    masterDB := config.ShardConfig{...}
    shardingDBs := make([]config.ShardConfig, dbCount)
    // ...
    cfg := &config.Config{...}
    
    manager, err := db.NewGroupManager(cfg)
    // ...
}
```

**変更後のコード構造**:
```go
func SetupTestGroupManager(t *testing.T, dbCount int, tablesPerDB int) *db.GroupManager {
    // 設定ファイルから読み込む
    cfg, err := LoadTestConfig()
    require.NoError(t, err)
    
    // 設定からGroupManagerを作成
    manager, err := db.NewGroupManager(cfg)
    require.NoError(t, err)
    
    // データベースをクリア
    ClearTestDatabase(t, manager)
    
    // スキーマを初期化（既存の実装を維持）
    masterConn, err := manager.GetMasterConnection()
    require.NoError(t, err)
    InitMasterSchema(t, masterConn.DB)
    
    // ... 既存のシャーディングスキーマ初期化処理 ...
    
    return manager
}
```

**注意事項**:
- 関数のシグネチャは変更しない（既存のテストコードへの影響を避ける）
- `dbCount`と`tablesPerDB`パラメータは使用しないが、互換性のため維持
- 既存のスキーマ初期化処理（`InitMasterSchema, InitShardingSchema`）は維持

**受け入れ基準**:
- `SetupTestGroupManager()`関数が設定ファイルから読み込むように修正されている
- ハードコードされた設定が削除されている
- データベースクリア機能が自動実行される
- 既存のスキーマ初期化処理が維持されている
- 関数のシグネチャが変更されていない

- _Requirements: 3.3.2, 6.3_
- _Design: 3.2.2_

---

#### タスク 4.2: `SetupTestGroupManager8Sharding()`関数の修正
**目的**: ハードコードされた設定を削除し、設定ファイルから読み込むように修正する。

**作業内容**:
- `server/test/testutil/db.go`の`SetupTestGroupManager8Sharding()`関数を修正
- ハードコードされた設定の作成部分を削除
- `LoadTestConfig()`を呼び出して設定を読み込む
- `db.NewGroupManager(cfg)`を使用してGroupManagerを作成
- `ClearTestDatabase()`を呼び出してデータベースをクリア
- 既存のスキーマ初期化処理は維持

**変更前のコード構造**:
```go
func SetupTestGroupManager8Sharding(t *testing.T) *db.GroupManager {
    // ハードコードされた設定を作成
    masterDB := config.ShardConfig{...}
    shardingDBs := []config.ShardConfig{...}
    // ...
    cfg := &config.Config{...}
    
    manager, err := db.NewGroupManager(cfg)
    // ...
}
```

**変更後のコード構造**:
```go
func SetupTestGroupManager8Sharding(t *testing.T) *db.GroupManager {
    // 設定ファイルから読み込む
    cfg, err := LoadTestConfig()
    require.NoError(t, err)
    
    // 設定からGroupManagerを作成
    manager, err := db.NewGroupManager(cfg)
    require.NoError(t, err)
    
    // データベースをクリア
    ClearTestDatabase(t, manager)
    
    // スキーマを初期化（既存の実装を維持）
    masterConn, err := manager.GetMasterConnection()
    require.NoError(t, err)
    InitMasterSchema(t, masterConn.DB)
    
    // ... 既存のシャーディングスキーマ初期化処理 ...
    
    return manager
}
```

**注意事項**:
- 関数のシグネチャは変更しない（既存のテストコードへの影響を避ける）
- 既存のスキーマ初期化処理（`InitMasterSchema, InitShardingSchema`）は維持

**受け入れ基準**:
- `SetupTestGroupManager8Sharding()`関数が設定ファイルから読み込むように修正されている
- ハードコードされた設定が削除されている
- データベースクリア機能が自動実行される
- 既存のスキーマ初期化処理が維持されている
- 関数のシグネチャが変更されていない

- _Requirements: 3.3.2, 6.3_
- _Design: 3.2.3_

---

### Phase 5: テストの実装

#### タスク 5.1: 設定ファイル読み込みのテスト実装
**目的**: `LoadTestConfig()`関数のテストを実装する。

**作業内容**:
- `server/test/testutil/db_test.go`を新規作成（存在しない場合）
- `TestLoadTestConfig`関数を実装
- 設定が正常に読み込まれることを確認
- データベース名に`_test`サフィックスが付与されていることを確認
- その他の設定が開発環境と同じであることを確認

**テストコード**:
```go
package testutil

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestLoadTestConfig(t *testing.T) {
    cfg, err := LoadTestConfig()
    require.NoError(t, err)
    
    // マスターデータベース名を確認
    assert.Equal(t, "webdb_master_test", cfg.Database.Groups.Master[0].Name)
    
    // シャーディングデータベース名を確認
    assert.Equal(t, "webdb_sharding_1_test", cfg.Database.Groups.Sharding.Databases[0].Name)
    
    // その他の設定を確認（ホスト、ポート等）
    assert.Equal(t, "localhost", cfg.Database.Groups.Master[0].Host)
    assert.Equal(t, 5432, cfg.Database.Groups.Master[0].Port)
    assert.Equal(t, "webdb", cfg.Database.Groups.Master[0].User)
    assert.Equal(t, "webdb", cfg.Database.Groups.Master[0].Password)
}
```

**受け入れ基準**:
- テストが実装されている
- テストが正常に実行できる
- すべてのアサーションが成功する

- _Requirements: 6.6_
- _Design: 4.1_

---

#### タスク 5.2: データベースクリア機能のテスト実装
**目的**: `ClearTestDatabase()`関数のテストを実装する。

**作業内容**:
- `server/test/testutil/db_test.go`に`TestClearTestDatabase`関数を追加
- テストデータを挿入
- `ClearTestDatabase()`を呼び出し
- データがクリアされたことを確認
- テーブル構造が維持されることを確認

**テストコード**:
```go
func TestClearTestDatabase(t *testing.T) {
    manager := SetupTestGroupManager(t, 4, 8)
    defer CleanupTestGroupManager(manager)
    
    // テストデータを挿入
    masterConn, err := manager.GetMasterConnection()
    require.NoError(t, err)
    masterConn.DB.Exec("INSERT INTO dm_news (title, content) VALUES ('test', 'test')")
    
    // データが存在することを確認
    var count int64
    masterConn.DB.Raw("SELECT COUNT(*) FROM dm_news").Scan(&count)
    assert.Greater(t, count, int64(0))
    
    // データベースをクリア
    ClearTestDatabase(t, manager)
    
    // データがクリアされたことを確認
    masterConn.DB.Raw("SELECT COUNT(*) FROM dm_news").Scan(&count)
    assert.Equal(t, int64(0), count)
    
    // テーブルが存在することを確認（構造が維持される）
    var exists bool
    masterConn.DB.Raw(`
        SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_schema = 'public' 
            AND table_name = 'dm_news'
        )
    `).Scan(&exists)
    assert.True(t, exists)
}
```

**受け入れ基準**:
- テストが実装されている
- テストが正常に実行できる
- すべてのアサーションが成功する
- データがクリアされることを確認できる
- テーブル構造が維持されることを確認できる

- _Requirements: 6.4, 6.6_
- _Design: 4.2_

---

### Phase 6: 動作確認

#### タスク 6.1: テスト用データベースの事前準備
**目的**: テスト用データベースを作成し、マイグレーションを実行する。

**作業内容**:
- テスト用データベースを作成
  - `webdb_master_test`
  - `webdb_sharding_1_test`
  - `webdb_sharding_2_test`
  - `webdb_sharding_3_test`
  - `webdb_sharding_4_test`
- 各データベースに対してマイグレーションを実行（`scripts/migrate-test.sh`を使用）

**実行コマンド**:
```bash
# データベースの作成
createdb -U webdb webdb_master_test
createdb -U webdb webdb_sharding_1_test
createdb -U webdb webdb_sharding_2_test
createdb -U webdb webdb_sharding_3_test
createdb -U webdb webdb_sharding_4_test

# マイグレーションの実行（スクリプトを使用）
./scripts/migrate-test.sh all

# または、個別に実行
./scripts/migrate-test.sh master
./scripts/migrate-test.sh sharding
```

**受け入れ基準**:
- すべてのテスト用データベースが作成されている
- すべてのデータベースに対してマイグレーションが実行されている
- マイグレーションが正常に完了している

- _Requirements: 6.5_
- _Design: 4.3.2_

---

#### タスク 6.2: 既存テストの動作確認
**目的**: 既存のテストが正常に動作することを確認する。

**作業内容**:
- 既存のテストを実行
- テスト用データベースが使用されることを確認
- すべてのテストが正常に実行されることを確認

**実行コマンド**:
```bash
# テストを実行
APP_ENV=test go test ./server/test/...
```

**受け入れ基準**:
- 既存のテストが正常に実行される
- テスト用データベースが使用される
- すべてのテストが成功する

- _Requirements: 6.3, 6.5_
- _Design: 4.3_

---

#### タスク 6.3: 開発環境との分離確認
**目的**: テスト実行後、開発環境のデータベースに影響がないことを確認する。

**作業内容**:
- テストを実行
- 開発環境のデータベース（`webdb_master`, `webdb_sharding_1`等）のデータを確認
- データが変更されていないことを確認

**確認方法**:
```bash
# 開発環境のデータベースのデータを確認
psql -U webdb -d webdb_master -c "SELECT COUNT(*) FROM dm_news;"
psql -U webdb -d webdb_sharding_1 -c "SELECT COUNT(*) FROM dm_users_000;"
```

**受け入れ基準**:
- テスト実行後、開発環境のデータベースに影響がない
- 開発環境のデータが変更されていない

- _Requirements: 6.5_
- _Design: 5.2_

---

#### タスク 6.4: パフォーマンス確認
**目的**: データベースクリア処理のパフォーマンスを確認する。

**作業内容**:
- テストデータを大量に挿入
- `ClearTestDatabase()`の実行時間を測定
- 10秒以内に完了することを確認

**受け入れ基準**:
- データベースクリア処理が10秒以内に完了する
- パフォーマンスに問題がない

- _Requirements: 4.1, 6.6_
- _Design: 6.3_

---

## 実装順序

1. **Phase 1**: 設定ファイルの作成（タスク 1.1, 1.2, 1.3）
2. **Phase 2**: 設定読み込み機能の実装（タスク 2.1）
3. **Phase 3**: データベースクリア機能の実装（タスク 3.1, 3.2）
4. **Phase 4**: 既存関数の修正（タスク 4.1, 4.2）
5. **Phase 5**: テストの実装（タスク 5.1, 5.2）
6. **Phase 6**: 動作確認（タスク 6.1, 6.2, 6.3, 6.4）

## 注意事項

### 実装時の注意点
- 既存のテストコードへの影響を最小限に抑える
- 関数のシグネチャは変更しない
- エラーハンドリングを適切に実装する
- 設定ファイルのコメントを適切に追加する

### テスト実行時の注意点
- テスト実行前にテスト用データベースを作成し、マイグレーションを実行する
- `APP_ENV=test`を設定してテストを実行する
- テスト実行後、開発環境のデータベースに影響がないことを確認する

### トラブルシューティング
- 設定ファイルが読み込めない場合: `config/test/database.yaml`が存在することを確認
- データベース接続エラー: テスト用データベースが作成されていることを確認
- マイグレーションエラー: マイグレーションが実行されていることを確認
