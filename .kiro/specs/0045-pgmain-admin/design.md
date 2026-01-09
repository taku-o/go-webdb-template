# GoAdminサーバーPostgreSQL対応設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、GoAdminサーバーでPostgreSQLを利用するように修正するための詳細設計を定義する。既存のSQLite設定をPostgreSQL設定に切り替え、設定ファイルからPostgreSQL接続情報を読み込んでGoAdminの設定に反映する。

### 1.2 設計の範囲
- `server/internal/admin/config.go`の`getDatabaseConfig()`関数をPostgreSQL対応に修正（SQLite設定を削除）
- GoAdminのデータベース接続設定をPostgreSQLドライバーに変更
- 設定ファイル（`config/{env}/database.yaml`）からPostgreSQL接続情報を読み込む
- SQLite用ライブラリのインポート確認（既にPostgreSQLドライバーはインポート済み）
- ドキュメントの更新

### 1.3 設計方針
- **既存システムとの統合**: 既存の設定ファイル構造（`config/{env}/database.yaml`）を活用し、`config.Load()`で読み込んだPostgreSQL接続情報を使用
- **設定ファイルベース**: 環境変数ではなく設定ファイル（`config/{env}/database.yaml`）から接続情報を読み込む
- **環境別対応**: `APP_ENV`環境変数で環境を指定（develop/staging/production）
- **masterデータベースのみ使用**: GoAdminはmasterデータベース（`webdb_master`）のみを使用
- **既存設定ファイルの維持**: `config/{env}/database.yaml`の構造は維持（Issue #87でPostgreSQL設定が有効になっている）

## 2. アーキテクチャ設計

### 2.1 既存アーキテクチャの分析

#### 2.1.1 現在の構成
- **設定ファイル**: `config/develop/database.yaml`でPostgreSQL設定が有効になっている（Issue #87で対応済み）
- **データベース接続**: `server/cmd/admin/main.go`で既にPostgreSQLドライバー（`github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres`）がインポートされている
- **GoAdmin設定**: `server/internal/admin/config.go`の`getDatabaseConfig()`関数でSQLite設定（`Driver: "sqlite"`、`File: dsn`）がハードコードされている
- **設定読み込み**: `server/cmd/admin/main.go`で`config.Load()`で設定を読み込み、`appdb.NewGroupManager(cfg)`でデータベース接続を初期化しているが、GoAdminの設定では使用されていない
- **データベース構成**:
  - マスターデータベース: 1台（`webdb_master`）
  - シャーディングデータベース: 4台（`webdb_sharding_1` ～ `webdb_sharding_4`）
  - **論理シャーディング数: 8**（物理DB 4台 × 2論理シャード）
  - GoAdminはmasterデータベースのみを使用

#### 2.1.2 既存パターンの維持
- **設定ファイル構造**: `config/{env}/database.yaml`の構造は維持（Issue #87でPostgreSQL設定が有効になっている）
- **GoAdminの初期化**: `server/cmd/admin/main.go`の初期化処理は変更不要
- **PostgreSQLドライバー**: 既にインポートされているため、追加のインポートは不要

### 2.2 システム構成図

```
┌─────────────────────────────────────────────────────────────┐
│          GoAdminサーバー (server/cmd/admin/main.go)           │
│                                                              │
│  ┌────────────────────────────────────────────────────┐   │
│  │ config.Load()                                       │   │
│  │   ↓                                                 │   │
│  │ config/{env}/database.yaml                         │   │
│  │   - master: 1台 (webdb_master)                      │   │
│  │   - sharding: 8つの論理シャード (物理DB 4台)        │   │
│  └────────────────────────────────────────────────────┘   │
│                          │                                  │
│                          ▼                                  │
│  ┌────────────────────────────────────────────────────┐   │
│  │ admin.NewConfig(cfg)                                │   │
│  │   ↓                                                 │   │
│  │ server/internal/admin/config.go                     │   │
│  │   - getDatabaseConfig(): PostgreSQL設定に修正       │   │
│  │   - Driver: "postgres"                              │   │
│  │   - Dsn: PostgreSQL接続文字列                       │   │
│  └────────────────────────────────────────────────────┘   │
│                          │                                  │
│                          ▼                                  │
│  ┌────────────────────────────────────────────────────┐   │
│  │ GoAdmin Engine                                       │   │
│  │   - PostgreSQLドライバー: インポート済み            │   │
│  │   - 接続: masterデータベースのみ                    │   │
│  └────────────────────────────────────────────────────┘   │
│                          │                                  │
│                          ▼                                  │
│  ┌────────────────────────────────────────────────────┐   │
│  │ PostgreSQL接続                                      │   │
│  │   - master: postgres-master:5432                    │   │
│  │   - DB: webdb_master                                │   │
│  └────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│              Docker Compose (docker-compose.postgres.yml)    │
│                                                              │
│  ┌──────────────────┐                                      │
│  │ postgres-master   │                                      │
│  │ (ポート: 5432)    │                                      │
│  │ DB: webdb_master  │                                      │
│  └──────────────────┘                                      │
│                                                              │
│  ┌──────────────────┐  ┌──────────────────┐              │
│  │ postgres-sharding │  │ postgres-sharding │              │
│  │ -1 (ポート: 5433) │  │ -2 (ポート: 5434) │              │
│  │ DB: webdb_sharding_1│ │ DB: webdb_sharding_2│             │
│  └──────────────────┘  └──────────────────┘              │
│                                                              │
│  ┌──────────────────┐  ┌──────────────────┐              │
│  │ postgres-sharding │  │ postgres-sharding │              │
│  │ -3 (ポート: 5435) │  │ -4 (ポート: 5436) │              │
│  │ DB: webdb_sharding_3│ │ DB: webdb_sharding_4│             │
│  └──────────────────┘  └──────────────────┘              │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 ディレクトリ構造

#### 2.3.1 変更前の構造
```
.
├── config/
│   ├── develop/
│   │   └── database.yaml          # PostgreSQL設定が有効（Issue #87で対応済み）
│   ├── staging/
│   │   └── database.yaml          # PostgreSQL設定（Issue #87で対応済み）
│   └── production/
│       ├── database.yaml          # PostgreSQL設定（存在する場合、Issue #87で対応済み）
│       └── database.yaml.example  # PostgreSQL設定（Issue #87で対応済み）
├── server/
│   ├── internal/
│   │   └── admin/
│   │       └── config.go          # SQLite設定（Driver: "sqlite", File: dsn）がハードコード
│   ├── cmd/
│   │   └── admin/
│   │       └── main.go            # PostgreSQLドライバーがインポート済み
│   └── go.mod                     # 変更不要
└── README.md                      # GoAdminサーバーのPostgreSQL利用に関する記述が未整備
```

#### 2.3.2 変更後の構造
```
.
├── config/
│   ├── develop/
│   │   └── database.yaml          # 変更なし（Issue #87でPostgreSQL設定が有効）
│   ├── staging/
│   │   └── database.yaml          # 変更なし（Issue #87でPostgreSQL設定が有効）
│   └── production/
│       ├── database.yaml          # 変更なし（Issue #87でPostgreSQL設定が有効）
│       └── database.yaml.example  # 変更なし（Issue #87でPostgreSQL設定が有効）
├── server/
│   ├── internal/
│   │   └── admin/
│   │       └── config.go          # PostgreSQL設定（Driver: "postgres", Dsn: dsn）に修正
│   ├── cmd/
│   │   └── admin/
│   │       └── main.go            # 変更なし（PostgreSQLドライバーがインポート済み）
│   └── go.mod                     # 変更なし
└── README.md                      # GoAdminサーバーのPostgreSQL利用に関する記述を追加
```

## 3. 詳細設計

### 3.1 GoAdmin設定の修正

#### 3.1.1 server/internal/admin/config.go

| フィールド | 詳細 |
|-----------|------|
| Intent | `getDatabaseConfig()`関数をPostgreSQL対応に修正 |
| Requirements | 3.1.1, 3.1.2 |

**変更内容**:

1. **SQLite設定の削除**:
   - `Driver: "sqlite"`を削除
   - `File: dsn`を削除

2. **PostgreSQL設定の追加**:
   - `Driver: "postgres"`を設定
   - `Dsn: dsn`を設定（PostgreSQL用のDSN形式）

3. **設定ファイルからの接続情報取得**:
   - `c.appConfig.Database.Groups.Master[0]`から接続情報を取得
   - 接続情報が存在しない場合のエラーハンドリングを追加

4. **PostgreSQL用DSN形式の構築**:
   - `postgres://user:password@host:port/dbname?sslmode=disable`形式でDSNを構築
   - 開発環境では`sslmode=disable`、本番環境では適切なSSL設定を推奨

**修正後のコード例**:
```go
// getDatabaseConfig はGoAdmin用のデータベース設定を返す
func (c *Config) getDatabaseConfig() goadminConfig.DatabaseList {
	// masterグループのデータベースをGoAdmin用データベースとして使用
	if len(c.appConfig.Database.Groups.Master) == 0 {
		panic("no database configuration found: master group is required")
	}

	masterDB := c.appConfig.Database.Groups.Master[0]
	
	// 接続情報の検証
	if masterDB.Host == "" || masterDB.Port == 0 || masterDB.User == "" || masterDB.Name == "" {
		panic("incomplete database configuration: host, port, user, and name are required")
	}

	// PostgreSQL用のDSN形式を構築
	// postgres://user:password@host:port/dbname?sslmode=disable
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		masterDB.Host,
		masterDB.Port,
		masterDB.User,
		masterDB.Password,
		masterDB.Name,
	)

	// 開発環境以外ではsslmodeを適切に設定（将来の拡張項目）
	env := os.Getenv("APP_ENV")
	if env == "production" {
		// 本番環境では適切なSSL設定を推奨
		// 現時点ではsslmode=disableのまま（設定ファイルから読み込む場合は追加実装が必要）
	}

	return goadminConfig.DatabaseList{
		"default": {
			Driver: "postgres",
			Dsn:    dsn,
		},
	}
}
```

**エラーハンドリング**:
- masterグループのデータベース設定が存在しない場合: `panic("no database configuration found: master group is required")`
- 接続情報が不完全な場合: `panic("incomplete database configuration: host, port, user, and name are required")`
- エラーメッセージが適切か確認

**実装上の注意事項**:
- `fmt`パッケージをインポートする必要がある（DSN構築のため）
- `os`パッケージをインポートする必要がある（環境変数の読み込みのため）
- DSN形式は`host=... port=... user=... password=... dbname=... sslmode=...`形式を使用（GoAdminのPostgreSQLドライバーがサポートする形式）
- パスワードに特殊文字が含まれる場合は適切にエスケープする必要がある（現時点では考慮不要）

#### 3.1.2 SQLite用ライブラリのインポート確認

| フィールド | 詳細 |
|-----------|------|
| Intent | `server/cmd/admin/main.go`にSQLite用ライブラリのインポートが存在しないことを確認 |
| Requirements | 3.2.1 |

**確認内容**:
1. `_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"`のインポートが存在しないことを確認
2. PostgreSQLドライバー（`_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"`）がインポートされていることを確認

**確認方法**:
- `grep`コマンドでSQLite用ライブラリのインポートを検索
- `grep`コマンドでPostgreSQLドライバーのインポートを検索

### 3.2 ドキュメントの更新

#### 3.2.1 README.md

| フィールド | 詳細 |
|-----------|------|
| Intent | GoAdminサーバーのPostgreSQL利用に関する記述を追加 |
| Requirements | 3.3.1 |

**更新内容**:
1. **GoAdminサーバーのPostgreSQL利用に関する記述を追加**:
   - GoAdminサーバーがPostgreSQLを利用することを明記
   - 設定ファイル（`config/{env}/database.yaml`）から接続情報を読み込むことを明記
   - masterデータベースのみを使用することを明記

2. **設定ファイルの変更方法を記載**:
   - `config/{env}/database.yaml`の設定方法を記載
   - PostgreSQL接続情報の設定方法を記載

3. **開発環境でのPostgreSQL起動手順を記載**:
   - PostgreSQLコンテナの起動手順（`./scripts/start-postgres.sh start`）
   - マイグレーション適用手順（`./scripts/migrate.sh`）
   - GoAdminサーバー起動手順

**追加するセクション例**:
```markdown
## GoAdminサーバーのPostgreSQL利用

GoAdminサーバーはPostgreSQLを利用します。設定ファイル（`config/{env}/database.yaml`）から接続情報を読み込み、masterデータベース（`webdb_master`）に接続します。

### 起動手順

1. PostgreSQLコンテナの起動:
   ```bash
   ./scripts/start-postgres.sh start
   ```

2. マイグレーションの適用:
   ```bash
   ./scripts/migrate.sh
   ```

3. GoAdminサーバーの起動:
   ```bash
   cd server
   go run cmd/admin/main.go
   ```

### 設定ファイル

`config/{env}/database.yaml`でPostgreSQL接続情報を設定します。詳細は[設定ファイルの説明](#設定ファイル)を参照してください。
```

#### 3.2.2 その他のドキュメント

| フィールド | 詳細 |
|-----------|------|
| Intent | GoAdminサーバーのSQLiteに関する記述をPostgreSQLに変更 |
| Requirements | 3.3.2 |

**更新内容**:
- `docs/`配下の関連ドキュメントでGoAdminサーバーのSQLiteに関する記述をPostgreSQLに変更
- GoAdminサーバーのPostgreSQL利用に関する記述を追加

## 4. データモデル

### 4.1 設定ファイル構造

設定ファイル（`config/{env}/database.yaml`）の構造は既存のまま維持（Issue #87でPostgreSQL設定が有効になっている）。

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
```

### 4.2 GoAdmin設定構造

GoAdminの設定構造（`goadminConfig.DatabaseList`）:

```go
type DatabaseList map[string]Database

type Database struct {
    Driver string
    Dsn    string
    File   string  // SQLite用（削除対象）
    // ... その他のフィールド
}
```

**変更内容**:
- `Driver: "sqlite"` → `Driver: "postgres"`
- `File: dsn` → `Dsn: dsn`（PostgreSQL用のDSN形式）

## 5. エラーハンドリング

### 5.1 エラーケース

| エラーケース | エラーハンドリング | エラーメッセージ |
|------------|------------------|----------------|
| masterグループのデータベース設定が存在しない | `panic` | `"no database configuration found: master group is required"` |
| 接続情報が不完全（host, port, user, nameのいずれかが空） | `panic` | `"incomplete database configuration: host, port, user, and name are required"` |
| PostgreSQL接続エラー | GoAdminのエラーハンドリングに委譲 | GoAdminのエラーメッセージ |

### 5.2 エラーハンドリング方針

- **設定エラー**: `panic`を使用（起動時に検出可能なエラー）
- **接続エラー**: GoAdminのエラーハンドリングに委譲（実行時に検出されるエラー）

## 6. テスト戦略

### 6.1 単体テスト

#### 6.1.1 server/internal/admin/config.goのテスト

**テストケース**:
1. **正常系**: masterグループのデータベース設定が存在する場合、PostgreSQL設定が正しく構築されることを確認
2. **異常系**: masterグループのデータベース設定が存在しない場合、エラーが発生することを確認
3. **異常系**: 接続情報が不完全な場合、エラーが発生することを確認
4. **正常系**: PostgreSQL用のDSN形式が正しく構築されることを確認

**テストファイル**: `server/internal/admin/config_test.go`（新規作成または既存の修正）

### 6.2 統合テスト

#### 6.2.1 GoAdminサーバーの起動テスト

**テストケース**:
1. **正常系**: GoAdminサーバーがPostgreSQLに正常に接続できることを確認
2. **正常系**: GoAdminサーバーがmasterデータベース（`webdb_master`）に正常に接続できることを確認
3. **正常系**: GoAdminサーバーが正常に起動し、管理画面にアクセスできることを確認
4. **正常系**: GoAdmin管理画面でデータベースのテーブルが正常に表示できることを確認

**テスト方法**:
- 手動テスト（ブラウザで管理画面にアクセス）
- 自動テスト（HTTPリクエストで管理画面にアクセス）

## 7. 実装上の注意事項

### 7.1 GoAdmin設定の修正

- **設定ファイル**: `config/{env}/database.yaml`から接続情報を読み込む
- **環境変数**: `APP_ENV`環境変数で環境を指定（develop/staging/production）
- **設定構造**: 既存の`config/{env}/database.yaml`の構造は維持（Issue #87でPostgreSQL設定が有効になっている）
- **PostgreSQL設定**: Issue #86で定義されたPostgreSQL構成に合わせる
- **DSN形式**: PostgreSQL用のDSN形式（`host=... port=... user=... password=... dbname=... sslmode=...`）を使用

### 7.2 エラーハンドリング

- **masterグループの存在確認**: `c.appConfig.Database.Groups.Master`が空でないことを確認
- **接続情報の完全性確認**: 必要な接続情報（host, port, user, name）がすべて存在することを確認
- **エラーメッセージ**: 適切なエラーメッセージを返す

### 7.3 ドキュメント整備

- **起動手順**: PostgreSQLコンテナの起動・マイグレーション適用・GoAdminサーバー起動の手順を記載
- **設定ファイル**: `config/{env}/database.yaml`の設定方法を記載
- **トラブルシューティング**: よくある問題と解決方法を記載

### 7.4 動作確認

- **接続確認**: GoAdminサーバー起動時にPostgreSQLへの接続を確認
- **クエリ実行**: GoAdmin管理画面でデータベースのテーブルが正常に表示できることを確認
- **エラーハンドリング**: 接続エラー時のエラーハンドリングを確認

## 8. 参考情報

### 8.1 関連Issue
- GitHub Issue #85: 開発環境はPostgreSQLを利用する前提とする
- GitHub Issue #86: PostgreSQLの起動スクリプトと、Atlasマイグレーションスクリプトの修正
- GitHub Issue #87: APIサーバーの修正
- GitHub Issue #88: GoAdminサーバーの修正

### 8.2 既存ドキュメント
- `README.md`: プロジェクト概要とセットアップ手順
- `docs/Architecture.md`: システムアーキテクチャ
- `docs/Initial-Setup.md`: 初期セットアップ手順
- `config/{env}/database.yaml`: 環境別データベース設定

### 8.3 技術スタック
- **PostgreSQL**: 15-alpine（Dockerイメージ）
- **GoAdmin**: 既存のGoAdminライブラリ
- **PostgreSQLドライバー**: `github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres`

### 8.4 参考リンク
- PostgreSQL公式ドキュメント: https://www.postgresql.org/docs/
- GoAdmin公式ドキュメント: https://github.com/GoAdminGroup/go-admin
- GoAdmin PostgreSQLドライバー: https://github.com/GoAdminGroup/go-admin/tree/master/modules/db/drivers/postgres
