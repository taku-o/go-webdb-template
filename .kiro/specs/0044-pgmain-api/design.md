# APIサーバーPostgreSQL対応設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、APIサーバーでPostgreSQLを利用するように修正するための詳細設計を定義する。既存のSQLite設定をPostgreSQL設定に切り替え、SQLite用ライブラリと処理分岐を削除し、開発環境、staging環境、production環境でPostgreSQLを利用できるようにする。

### 1.2 設計の範囲
- `config/develop/database.yaml`のPostgreSQL設定を有効化（SQLite設定を削除、論理シャーディング8つを定義）
- `config/staging/database.yaml`のPostgreSQL設定を確認・修正
- `config/production/database.yaml`のPostgreSQL設定を確認・修正（存在する場合）
- `config/production/database.yaml.example`のPostgreSQL設定を確認・修正
- SQLite用ライブラリのインポート削除
- ソースコード中のSQLite用処理分岐の削除
- テストコードのPostgreSQL対応
- ドキュメントの更新

### 1.3 設計方針
- **既存システムとの統合**: 既存のデータベース接続コード（`server/internal/db/connection.go`）は既にPostgreSQLドライバーをサポートしているため、SQLite用の処理のみを削除
- **設定ファイルベース**: 環境変数ではなく設定ファイル（`config/{env}/database.yaml`）から接続情報を読み込む
- **環境別対応**: `APP_ENV`環境変数で環境を指定（develop/staging/production）
- **論理シャーディング数8**: 物理DB 4台、各物理DBに2つの論理シャードを割り当て、合計8つの論理シャーディング設定を定義
- **既存設定ファイルの維持**: `config/{env}/database.yaml`の構造は維持（SQLite設定の削除とPostgreSQL設定の追加のみ）

## 2. アーキテクチャ設計

### 2.1 既存アーキテクチャの分析

#### 2.1.1 現在の構成
- **設定ファイル**: `config/develop/database.yaml`でSQLite設定が有効、PostgreSQL設定がコメントアウトされている（または未定義）
- **データベース接続**: `server/internal/db/connection.go`で既にPostgreSQLドライバー（`gorm.io/driver/postgres`）がサポートされている
- **APIサーバー**: `server/cmd/server/main.go`で`config.Load()`で設定を読み込み、`db.NewGroupManager(cfg)`でデータベース接続を初期化
- **SQLite用ライブラリ**: 
  - `server/internal/db/connection.go`: `_ "github.com/mattn/go-sqlite3"`、`"gorm.io/driver/sqlite"`がインポートされている
  - `server/cmd/admin/main.go`: `_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"`がインポートされている
  - `server/go.mod`: `gorm.io/driver/sqlite v1.5.6`が依存関係に含まれている
- **SQLite用処理分岐**: `server/internal/db/connection.go`内に複数の`case "sqlite3":`分岐が存在

#### 2.1.2 既存パターンの維持
- **設定ファイル構造**: `config/{env}/database.yaml`の構造は維持
- **データベース接続コード**: PostgreSQLドライバーは既にサポートされているため、SQLite用の処理のみを削除
- **APIサーバーの初期化**: `server/cmd/server/main.go`の初期化処理は変更不要

### 2.2 システム構成図

```
┌─────────────────────────────────────────────────────────────┐
│                    APIサーバー (server/cmd/server/main.go)   │
│                                                              │
│  ┌────────────────────────────────────────────────────┐   │
│  │ config.Load()                                       │   │
│  │   ↓                                                 │   │
│  │ config/{env}/database.yaml                         │   │
│  │   - master: 1台                                     │   │
│  │   - sharding: 8つの論理シャード (物理DB 4台)        │   │
│  └────────────────────────────────────────────────────┘   │
│                          │                                  │
│                          ▼                                  │
│  ┌────────────────────────────────────────────────────┐   │
│  │ db.NewGroupManager(cfg)                            │   │
│  │   ↓                                                 │   │
│  │ server/internal/db/connection.go                   │   │
│  │   - PostgreSQLドライバー: サポート済み              │   │
│  │   - SQLiteドライバー: 削除対象                      │   │
│  └────────────────────────────────────────────────────┘   │
│                          │                                  │
│                          ▼                                  │
│  ┌────────────────────────────────────────────────────┐   │
│  │ PostgreSQL接続                                      │   │
│  │   - master: postgres-master:5432                    │   │
│  │   - sharding_1: postgres-sharding-1:5433            │   │
│  │   - sharding_2: postgres-sharding-2:5434            │   │
│  │   - sharding_3: postgres-sharding-3:5435            │   │
│  │   - sharding_4: postgres-sharding-4:5436            │   │
│  └────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│              Docker Compose (docker-compose.postgres.yml)    │
│                                                              │
│  ┌──────────────────┐  ┌──────────────────┐              │
│  │ postgres-master   │  │ postgres-sharding │              │
│  │ (ポート: 5432)    │  │ -1 (ポート: 5433) │              │
│  │ DB: webdb_master  │  │ DB: webdb_sharding_1│             │
│  └──────────────────┘  └────────┬─────────┘              │
│                                  │                          │
│                                  ├── postgres-sharding-2    │
│                                  │   (ポート: 5434)        │
│                                  │   DB: webdb_sharding_2   │
│                                  │                          │
│                                  ├── postgres-sharding-3    │
│                                  │   (ポート: 5435)        │
│                                  │   DB: webdb_sharding_3  │
│                                  │                          │
│                                  └── postgres-sharding-4    │
│                                      (ポート: 5436)        │
│                                      DB: webdb_sharding_4   │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 ディレクトリ構造

#### 2.3.1 変更前の構造
```
.
├── config/
│   ├── develop/
│   │   └── database.yaml          # SQLite設定が有効、PostgreSQL設定がコメントアウト
│   ├── staging/
│   │   └── database.yaml          # PostgreSQL設定（確認・修正が必要）
│   └── production/
│       ├── database.yaml          # PostgreSQL設定（存在する場合、確認・修正が必要）
│       └── database.yaml.example  # PostgreSQL設定（確認・修正が必要）
├── server/
│   ├── internal/
│   │   └── db/
│   │       └── connection.go      # SQLite用ライブラリと処理分岐が含まれている
│   ├── cmd/
│   │   ├── server/
│   │   │   └── main.go            # 変更不要
│   │   └── admin/
│   │       └── main.go            # SQLite用ライブラリがインポートされている
│   ├── test/
│   │   └── testutil/
│   │       └── db.go              # SQLite設定が含まれている
│   └── go.mod                     # gorm.io/driver/sqliteが依存関係に含まれている
└── README.md                      # SQLiteに関する記述
```

#### 2.3.2 変更後の構造
```
.
├── config/
│   ├── develop/
│   │   └── database.yaml          # PostgreSQL設定が有効、SQLite設定を削除、論理シャーディング8つ
│   ├── staging/
│   │   └── database.yaml          # PostgreSQL設定（確認・修正済み、論理シャーディング8つ）
│   └── production/
│       ├── database.yaml          # PostgreSQL設定（存在する場合、確認・修正済み、論理シャーディング8つ）
│       └── database.yaml.example  # PostgreSQL設定（確認・修正済み、論理シャーディング8つ）
├── server/
│   ├── internal/
│   │   └── db/
│   │       └── connection.go      # SQLite用ライブラリと処理分岐を削除
│   ├── cmd/
│   │   ├── server/
│   │   │   └── main.go            # 変更なし
│   │   └── admin/
│   │       └── main.go            # SQLite用ライブラリのインポートを削除
│   ├── test/
│   │   └── testutil/
│   │       └── db.go              # PostgreSQL設定に変更
│   └── go.mod                     # gorm.io/driver/sqliteを削除
└── README.md                      # PostgreSQLに関する記述に更新
```

## 3. コンポーネント設計

### 3.1 設定ファイルの修正

#### 3.1.1 config/develop/database.yaml

| フィールド | 詳細 |
|-----------|------|
| Intent | PostgreSQL設定を有効化し、SQLite設定を削除。論理シャーディング数8を定義 |
| Requirements | 3.1.1 |

**変更内容**:
1. **SQLite設定の削除**: 既存のSQLite設定（`database.groups.master`と`database.groups.sharding`）を完全に削除
2. **PostgreSQL設定の有効化**: コメントアウトされているPostgreSQL設定を有効化（または新規追加）
3. **論理シャーディング数8の定義**: shardingグループに8つのデータベース設定（id: 1-8）を定義

**設定構造**:
```yaml
database:
  groups:
    master:
      - id: 1
        driver: postgres
        host: localhost
        port: 5432
        user: webdb
        password: webdb
        name: webdb_master
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
          name: webdb_sharding_1
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [0, 3]
        # 論理シャード 2: テーブル _004-007 → postgres-sharding-1
        - id: 2
          driver: postgres
          host: localhost
          port: 5433
          user: webdb
          password: webdb
          name: webdb_sharding_1
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [4, 7]
        # 論理シャード 3: テーブル _008-011 → postgres-sharding-2
        - id: 3
          driver: postgres
          host: localhost
          port: 5434
          user: webdb
          password: webdb
          name: webdb_sharding_2
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [8, 11]
        # 論理シャード 4: テーブル _012-015 → postgres-sharding-2
        - id: 4
          driver: postgres
          host: localhost
          port: 5434
          user: webdb
          password: webdb
          name: webdb_sharding_2
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [12, 15]
        # 論理シャード 5: テーブル _016-019 → postgres-sharding-3
        - id: 5
          driver: postgres
          host: localhost
          port: 5435
          user: webdb
          password: webdb
          name: webdb_sharding_3
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [16, 19]
        # 論理シャード 6: テーブル _020-023 → postgres-sharding-3
        - id: 6
          driver: postgres
          host: localhost
          port: 5435
          user: webdb
          password: webdb
          name: webdb_sharding_3
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [20, 23]
        # 論理シャード 7: テーブル _024-027 → postgres-sharding-4
        - id: 7
          driver: postgres
          host: localhost
          port: 5436
          user: webdb
          password: webdb
          name: webdb_sharding_4
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [24, 27]
        # 論理シャード 8: テーブル _028-031 → postgres-sharding-4
        - id: 8
          driver: postgres
          host: localhost
          port: 5436
          user: webdb
          password: webdb
          name: webdb_sharding_4
          max_connections: 25
          max_idle_connections: 5
          connection_max_lifetime: 1h
          table_range: [28, 31]

      tables:
        - name: dm_users
          suffix_count: 32
        - name: dm_posts
          suffix_count: 32
```

**実装上の注意事項**:
- SQLite設定は完全に削除（コメントアウトではない）
- 論理シャーディング数は必ず8つ（id: 1-8）を定義
- 各論理シャードのtable_rangeが正しく設定されていることを確認
- 物理DBと論理シャードの対応関係を確認（同じ物理DBを参照する論理シャードは同じhost/port/nameを使用）

#### 3.1.2 config/staging/database.yaml

| フィールド | 詳細 |
|-----------|------|
| Intent | PostgreSQL設定を確認・修正し、論理シャーディング数8を確認 |
| Requirements | 3.1.2 |

**確認内容**:
1. PostgreSQL設定が正しく定義されているか確認
2. シャーディング構成が正しいか確認（物理DB 4台、論理シャーディング8）
3. shardingグループに8つのデータベース設定（id: 1-8）が定義されているか確認
4. 各論理シャードのtable_rangeが正しく設定されているか確認
5. SQLite設定が存在する場合は削除

**修正内容**:
- 不備があれば修正
- SQLite設定があれば削除
- 論理シャーディング数が8つになっていることを確認

#### 3.1.3 config/production/database.yaml

| フィールド | 詳細 |
|-----------|------|
| Intent | PostgreSQL設定を確認・修正し、論理シャーディング数8を確認（存在する場合） |
| Requirements | 3.1.3 |

**確認内容**:
1. ファイルが存在するか確認
2. PostgreSQL設定が正しく定義されているか確認
3. シャーディング構成が正しいか確認（物理DB 4台、論理シャーディング8）
4. shardingグループに8つのデータベース設定（id: 1-8）が定義されているか確認
5. 各論理シャードのtable_rangeが正しく設定されているか確認
6. SQLite設定が存在する場合は削除

**修正内容**:
- 不備があれば修正
- SQLite設定があれば削除
- 論理シャーディング数が8つになっていることを確認

#### 3.1.4 config/production/database.yaml.example

| フィールド | 詳細 |
|-----------|------|
| Intent | PostgreSQL設定を確認・修正し、論理シャーディング数8を確認 |
| Requirements | 3.1.4 |

**確認内容**:
1. PostgreSQL設定が正しく定義されているか確認
2. シャーディング構成が正しいか確認（物理DB 4台、論理シャーディング8）
3. shardingグループに8つのデータベース設定（id: 1-8）が定義されているか確認
4. 各論理シャードのtable_rangeが正しく設定されているか確認
5. 接続情報が適切か確認
6. SQLite設定が存在する場合は削除

**修正内容**:
- 不備があれば修正
- SQLite設定があれば削除
- 接続情報をIssue #86で定義された構成に合わせる
- 論理シャーディング数が8つになっていることを確認

**注意事項**:
- `.example`ファイルは本番環境の設定例として使用されるため、適切な設定例を記載する必要がある

### 3.2 SQLite用ライブラリと処理分岐の削除

#### 3.2.1 server/internal/db/connection.go

| フィールド | 詳細 |
|-----------|------|
| Intent | SQLite用ライブラリのインポートと処理分岐を削除 |
| Requirements | 3.2.1, 3.2.2 |

**削除対象**:
1. **インポート削除**:
   - `_ "github.com/mattn/go-sqlite3"`を削除
   - `"gorm.io/driver/sqlite"`を削除

2. **処理分岐削除**:
   - `NewConnection`関数内の`driver = "sqlite3"`デフォルト値設定を削除
   - `createGORMConnection`関数内の`case "sqlite3":`分岐を削除
   - `createGORMConnectionFromDSN`関数内の`case "sqlite3":`分岐を削除
   - `NewGORMConnection`関数内の`case "sqlite3":`分岐を削除（Reader接続作成部分）

**修正後のコード例**:
```go
// インポート部分
import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/taku-o/go-webdb-template/internal/config"
	// SQLite用ライブラリのインポートを削除
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	// "gorm.io/driver/sqlite" を削除
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// NewConnection関数
func NewConnection(cfg *config.ShardConfig) (*Connection, error) {
	// ... 既存のコード ...
	
	driver := cfg.Driver
	if driver == "" {
		// driver = "sqlite3" のデフォルト値設定を削除
		return nil, fmt.Errorf("driver is required")
	}
	
	// ... 既存のコード ...
}

// createGORMConnection関数
func createGORMConnection(cfg *config.ShardConfig, isWriter bool, sqlLogger *SQLLogger) (*gorm.DB, error) {
	// ... 既存のコード ...
	
	switch cfg.Driver {
	// case "sqlite3": を削除
	case "postgres":
		dialector = postgres.Open(dsn)
	case "mysql":
		dialector = mysql.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
	}
	
	// ... 既存のコード ...
}
```

**エラーハンドリング**:
- 未サポートドライバー指定時は`fmt.Errorf("unsupported driver: %s", driver)`でエラーを返す
- エラーメッセージが適切か確認

#### 3.2.2 server/cmd/admin/main.go

| フィールド | 詳細 |
|-----------|------|
| Intent | SQLite用ライブラリのインポートを削除 |
| Requirements | 3.2.1 |

**削除対象**:
- `_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"`のインポートを削除

**修正後のコード例**:
```go
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// SQLite用ライブラリのインポートを削除
	// _ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"
	_ "github.com/GoAdminGroup/themes/adminlte"
	
	// ... 既存のインポート ...
)
```

#### 3.2.3 server/go.mod

| フィールド | 詳細 |
|-----------|------|
| Intent | SQLite関連の依存関係を削除 |
| Requirements | 3.2.1 |

**削除対象**:
- `gorm.io/driver/sqlite v1.5.6`の依存関係を削除

**実装手順**:
1. `server/go.mod`から`gorm.io/driver/sqlite`の行を削除
2. `go mod tidy`を実行して依存関係を整理

### 3.3 テストコードの修正

#### 3.3.1 server/test/testutil/db.go

| フィールド | 詳細 |
|-----------|------|
| Intent | テストユーティリティをPostgreSQL対応に修正 |
| Requirements | 3.3.2 |

**修正内容**:
1. SQLite固有の設定をPostgreSQL設定に変更
2. テストデータベースの初期化方法をPostgreSQL対応に変更
3. テスト実行時のデータベース接続方法をPostgreSQL対応に変更

**実装上の注意事項**:
- テスト実行時にPostgreSQLデータベースを使用
- テストデータの投入方法をPostgreSQL対応に変更
- 論理シャーディング数8を考慮したテスト設定

#### 3.3.2 その他のテストファイル

| フィールド | 詳細 |
|-----------|------|
| Intent | テストファイルのSQLite依存を確認し、必要に応じて修正 |
| Requirements | 3.3.1, 3.3.2 |

**確認内容**:
1. SQLite固有の設定やコードが含まれているか確認
2. テストデータベースの初期化方法を確認
3. テスト実行時のデータベース接続方法を確認

**修正内容**:
- SQLite固有の設定をPostgreSQL設定に変更
- テストデータベースの初期化方法をPostgreSQL対応に変更
- テスト実行時のデータベース接続方法をPostgreSQL対応に変更

### 3.4 ドキュメントの更新

#### 3.4.1 README.md

| フィールド | 詳細 |
|-----------|------|
| Intent | PostgreSQL利用に関する記述を追加 |
| Requirements | 3.4.1 |

**更新内容**:
1. PostgreSQL利用に関する記述を追加
2. 設定ファイルの変更方法を記載
3. 開発環境でのPostgreSQL起動手順を記載

**記載内容**:
- PostgreSQLの起動方法（`./scripts/start-postgres.sh start`）
- マイグレーションの適用方法（`./scripts/migrate.sh`）
- APIサーバーの起動方法
- 設定ファイル（`config/develop/database.yaml`）の設定方法
- 論理シャーディング数8の説明

#### 3.4.2 その他のドキュメント

| フィールド | 詳細 |
|-----------|------|
| Intent | SQLiteに関する記述をPostgreSQLに変更 |
| Requirements | 3.4.2 |

**更新内容**:
1. SQLiteに関する記述をPostgreSQLに変更
2. PostgreSQL利用に関する記述を追加

**対象ファイル**:
- `docs/`配下の関連ドキュメント

## 4. データモデル設計

### 4.1 設定ファイル構造

#### 4.1.1 database.yaml構造

設定ファイルの構造は既存のものを維持し、以下の変更のみを行う：

1. **SQLite設定の削除**: `database.groups.master`と`database.groups.sharding`のSQLite設定を削除
2. **PostgreSQL設定の追加**: PostgreSQL設定を追加（コメントアウト解除または新規追加）
3. **論理シャーディング数8**: shardingグループに8つのデータベース設定（id: 1-8）を定義

### 4.2 データベース接続構造

#### 4.2.1 接続情報

- **masterグループ**: 1台（`postgres-master`、ポート5432、データベース名`webdb_master`）
- **shardingグループ**: 論理シャーディング数8（物理DB 4台、各物理DBに2つの論理シャード）
  - 物理データベース: 4台（`postgres-sharding-1` ～ `postgres-sharding-4`、ポート5433-5436、データベース名`webdb_sharding_1` ～ `webdb_sharding_4`）
  - 論理シャーディング設定: 8つ（id: 1-8）

## 5. エラーハンドリング設計

### 5.1 未サポートドライバー指定時のエラー

#### 5.1.1 エラーハンドリング

**実装**:
- `server/internal/db/connection.go`の各関数で、未サポートドライバー指定時に`fmt.Errorf("unsupported driver: %s", driver)`でエラーを返す

**エラーメッセージ**:
- 明確で分かりやすいエラーメッセージを返す
- サポートされているドライバー（`postgres`, `mysql`）を明示

### 5.2 設定ファイル読み込みエラー

#### 5.2.1 エラーハンドリング

**実装**:
- `config.Load()`で設定ファイルを読み込む際のエラーハンドリングは既存のものを使用
- 設定ファイルが存在しない、または不正な形式の場合は既存のエラーハンドリングを使用

## 6. テスト設計

### 6.1 単体テスト

#### 6.1.1 データベース接続テスト

**テスト内容**:
- PostgreSQL接続が正常に確立できることを確認
- masterデータベースへの接続が正常に確立できることを確認
- shardingデータベースへの接続が正常に確立できることを確認（8つの論理シャードすべて）

**テスト環境**:
- PostgreSQLコンテナが起動していること
- マイグレーションが適用されていること

#### 6.1.2 未サポートドライバー指定テスト

**テスト内容**:
- 未サポートドライバー（`sqlite3`など）を指定した場合にエラーが返されることを確認
- エラーメッセージが適切であることを確認

### 6.2 統合テスト

#### 6.2.1 APIサーバー起動テスト

**テスト内容**:
- APIサーバーがPostgreSQLに正常に接続できることを確認
- APIサーバーが正常に起動し、リクエストを処理できることを確認

**テスト環境**:
- PostgreSQLコンテナが起動していること
- マイグレーションが適用されていること

## 7. 実装上の注意事項

### 7.1 設定ファイルの修正

#### 7.1.1 論理シャーディング数8の確認

**注意事項**:
- `config/develop/database.yaml`のshardingグループに必ず8つのデータベース設定（id: 1-8）を定義
- 各論理シャードのtable_rangeが正しく設定されていることを確認
- 物理DBと論理シャードの対応関係を確認（同じ物理DBを参照する論理シャードは同じhost/port/nameを使用）

#### 7.1.2 SQLite設定の削除

**注意事項**:
- SQLite設定は完全に削除（コメントアウトではない）
- 設定ファイルの構造は維持

### 7.2 SQLite用ライブラリと処理分岐の削除

#### 7.2.1 インポート削除

**注意事項**:
- `server/internal/db/connection.go`と`server/cmd/admin/main.go`からSQLite用ライブラリのインポートを削除
- 未使用のインポートが残らないように注意

#### 7.2.2 処理分岐削除

**注意事項**:
- `server/internal/db/connection.go`内のSQLite用処理分岐（`case "sqlite3":`など）を削除
- デフォルト値設定（`driver = "sqlite3"`）を削除
- 未サポートドライバー指定時のエラーハンドリングを確認

#### 7.2.3 依存関係削除

**注意事項**:
- `server/go.mod`から`gorm.io/driver/sqlite`の依存関係を削除
- `go mod tidy`を実行して依存関係を整理

### 7.3 テストコードの修正

#### 7.3.1 テストユーティリティの修正

**注意事項**:
- `server/test/testutil/db.go`をPostgreSQL対応に修正
- テスト実行時にPostgreSQLデータベースを使用
- 論理シャーディング数8を考慮したテスト設定

### 7.4 ドキュメントの更新

#### 7.4.1 README.mdの更新

**注意事項**:
- PostgreSQL利用に関する記述を追加
- 設定ファイルの変更方法を記載
- 開発環境でのPostgreSQL起動手順を記載
- 論理シャーディング数8の説明を追加

## 8. 参考情報

### 8.1 関連Issue
- GitHub Issue #85: 開発環境はPostgreSQLを利用する前提とする
- GitHub Issue #86: PostgreSQLの起動スクリプトと、Atlasマイグレーションスクリプトの修正
- GitHub Issue #87: APIサーバーの修正

### 8.2 既存ドキュメント
- `README.md`: プロジェクト概要とセットアップ手順
- `docs/Architecture.md`: システムアーキテクチャ
- `docs/Initial-Setup.md`: 初期セットアップ手順
- `config/{env}/database.yaml`: 環境別データベース設定

### 8.3 技術スタック
- **PostgreSQL**: 15-alpine（Dockerイメージ）
- **GORM**: 既存のGORMライブラリ
- **PostgreSQLドライバー**: `gorm.io/driver/postgres`（既存）

### 8.4 参考リンク
- PostgreSQL公式ドキュメント: https://www.postgresql.org/docs/
- GORM公式ドキュメント: https://gorm.io/docs/
- GORM PostgreSQLドライバー: https://github.com/go-gorm/postgres
