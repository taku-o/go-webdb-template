# AdminサーバーのMySQL対応の設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、Adminサーバー（GoAdmin管理画面）をPostgreSQLに加えて、MySQLでも動作するように修正するための詳細設計を定義する。PostgreSQLとMySQLの両方に対応することで、Adminサーバーでもデータベース選択の柔軟性を提供する。

### 1.2 設計の範囲
- Adminサーバーのデータベース接続設定のMySQL対応設計
- GoAdminのMySQLドライバーのインポート追加設計
- DSN生成ロジックのMySQL対応設計

### 1.3 設計方針
- **後方互換性**: 既存のPostgreSQL機能に影響を与えない（追加のみ）
- **一貫性**: 既存のコードスタイルと一貫性を保つ
- **設定ファイルの活用**: メインサーバーと同じ設定ファイル（`database.yaml`または`database.mysql.yaml`）を使用
- **エラーハンドリング**: 設定ファイルの不備は明確なエラーメッセージで報告
- **既存コードの最小限の変更**: 既存のPostgreSQLコードは可能な限り変更せず、MySQL対応を追加する

## 2. アーキテクチャ設計

### 2.1 全体構成

```
設定ファイル読み込み（config.Load()）
  ↓
config.yaml から DB_TYPE を取得（既存の実装）
  ↓
DB_TYPE に応じて適切な database.yaml を読み込む（既存の実装）
  ├── DB_TYPE=postgresql → database.yaml (既存)
  └── DB_TYPE=mysql → database.mysql.yaml (既存)
  ↓
GroupManager の初期化（既存の実装）
  ↓
Adminサーバーの起動
  ↓
admin.NewConfig(cfg) で Config を作成
  ↓
GetGoAdminConfig() で GoAdmin設定を取得
  ↓
getDatabaseConfig() でデータベース設定を取得（修正対象）
  ├── masterDB.Driver を参照
  ├── postgres → PostgreSQL用DSN生成
  └── mysql → MySQL用DSN生成
  ↓
GoAdmin Engine の初期化
  ↓
データベース接続（PostgreSQL または MySQL）
  ├── PostgreSQL: github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres
  └── MySQL: github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql
```

### 2.2 データフロー

```
main.go
  ↓
config.Load() → Config構造体
  ↓
admin.NewConfig(cfg) → admin.Config構造体
  ↓
GetGoAdminConfig() → goadminConfig.Config構造体
  ↓
getDatabaseConfig() → goadminConfig.DatabaseList
  ├── masterDB.Driver を判定
  ├── postgres → PostgreSQL用DSN生成
  └── mysql → MySQL用DSN生成
  ↓
GoAdmin Engine に設定を渡す
  ↓
データベース接続確立
```

## 3. 詳細設計

### 3.1 `server/internal/admin/config.go`の修正

#### 3.1.1 `getDatabaseConfig()`メソッドの修正

**現状**:
- PostgreSQL用のDSN形式をハードコード
- `Driver: "postgresql"`をハードコード

**修正内容**:
- `c.appConfig.Database.Groups.Master[0].Driver`を参照してデータベースタイプを判定
- データベースタイプに応じて適切なDSN形式を生成
- ドライバー名もデータベースタイプに応じて設定

**実装詳細**:

```go
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

	// ドライバーの検証
	if masterDB.Driver == "" {
		panic("database driver is not specified: driver field is required in database configuration")
	}

	var dsn string
	var driverName string

	switch masterDB.Driver {
	case "postgres":
		// PostgreSQL用のDSN形式を構築
		// host=... port=... user=... password=... dbname=... sslmode=disable
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			masterDB.Host,
			masterDB.Port,
			masterDB.User,
			masterDB.Password,
			masterDB.Name,
		)
		driverName = "postgresql"
	case "mysql":
		// MySQL用のDSN形式を構築
		// user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=true&loc=Local
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
			masterDB.User,
			masterDB.Password,
			masterDB.Host,
			masterDB.Port,
			masterDB.Name,
		)
		driverName = "mysql"
	default:
		panic(fmt.Sprintf("unsupported database driver: %s (supported drivers: postgres, mysql)", masterDB.Driver))
	}

	return goadminConfig.DatabaseList{
		"default": {
			Driver: driverName,
			Dsn:    dsn,
		},
	}
}
```

**変更点**:
1. `masterDB.Driver`を参照してデータベースタイプを判定
2. `switch`文で`postgres`と`mysql`を分岐
3. PostgreSQL用とMySQL用のDSN形式をそれぞれ生成
4. ドライバー名を`postgresql`または`mysql`に設定
5. ドライバーが指定されていない場合、または未対応のドライバーの場合はエラーを返す

**エラーハンドリング**:
- ドライバーが指定されていない場合: `panic("database driver is not specified: driver field is required in database configuration")`
- 未対応のドライバーの場合: `panic(fmt.Sprintf("unsupported database driver: %s (supported drivers: postgres, mysql)", masterDB.Driver))`

### 3.2 `server/cmd/admin/main.go`の修正

#### 3.2.1 MySQLドライバーのインポート追加

**現状**:
- PostgreSQLドライバーのみインポート

**修正内容**:
- MySQLドライバーのインポートを追加

**実装詳細**:

```go
import (
	// ... 既存のインポート ...

	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql"  // 追加
	_ "github.com/GoAdminGroup/themes/adminlte"

	// ... 既存のインポート ...
)
```

**変更点**:
1. `_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql"`をインポート追加
2. 既存のPostgreSQLドライバーのインポートは維持
3. ブランクインポート（`_`）を使用してドライバーを登録

## 4. DSN形式の詳細

### 4.1 PostgreSQL用DSN形式

**形式**: `host=... port=... user=... password=... dbname=... sslmode=disable`

**例**:
```
host=localhost port=5432 user=webdb password=webdb dbname=webdb_master sslmode=disable
```

**パラメータ**:
- `host`: データベースホスト名
- `port`: データベースポート番号
- `user`: データベースユーザー名
- `password`: データベースパスワード
- `dbname`: データベース名
- `sslmode=disable`: SSL接続を無効化（開発環境用）

### 4.2 MySQL用DSN形式

**形式**: `user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=true&loc=Local`

**例**:
```
webdb:webdb@tcp(localhost:3306)/webdb_master?charset=utf8mb4&parseTime=true&loc=Local
```

**パラメータ**:
- `user`: データベースユーザー名
- `pass`: データベースパスワード
- `tcp(host:port)`: TCP接続（ホスト名とポート番号）
- `dbname`: データベース名
- `charset=utf8mb4`: 文字セット（UTF-8の完全なサポート）
- `parseTime=true`: 時刻型を`time.Time`に自動変換
- `loc=Local`: タイムゾーンをローカル時間に設定

### 4.3 DSN形式の違い

| 項目 | PostgreSQL | MySQL |
|------|-----------|-------|
| 形式 | `key=value`形式（スペース区切り） | `user:pass@tcp(host:port)/dbname?params`形式 |
| ドライバー名 | `postgresql` | `mysql` |
| 文字セット | 不要 | `charset=utf8mb4`必須 |
| 時刻解析 | 不要 | `parseTime=true`推奨 |
| タイムゾーン | 不要 | `loc=Local`推奨 |
| SSL設定 | `sslmode=disable` | 不要（デフォルトで無効） |

## 5. エラーハンドリング設計

### 5.1 エラーケースと対応

#### 5.1.1 ドライバーが指定されていない場合

**エラー**: `panic("database driver is not specified: driver field is required in database configuration")`

**発生箇所**: `getDatabaseConfig()`メソッド内

**対応**: 設定ファイル（`database.yaml`または`database.mysql.yaml`）の`driver`フィールドを確認

#### 5.1.2 未対応のドライバーが指定された場合

**エラー**: `panic(fmt.Sprintf("unsupported database driver: %s (supported drivers: postgres, mysql)", masterDB.Driver))`

**発生箇所**: `getDatabaseConfig()`メソッド内

**対応**: 設定ファイルの`driver`フィールドを`postgres`または`mysql`に修正

#### 5.1.3 接続情報が不完全な場合

**エラー**: `panic("incomplete database configuration: host, port, user, and name are required")`

**発生箇所**: `getDatabaseConfig()`メソッド内（既存のエラーハンドリング）

**対応**: 設定ファイルの接続情報を確認

#### 5.1.4 masterグループが存在しない場合

**エラー**: `panic("no database configuration found: master group is required")`

**発生箇所**: `getDatabaseConfig()`メソッド内（既存のエラーハンドリング）

**対応**: 設定ファイルにmasterグループを追加

### 5.2 エラーメッセージの設計方針

- **明確性**: エラーの原因を明確に示す
- **具体的な指示**: 修正方法を具体的に示す
- **設定ファイルの参照**: どの設定ファイルを確認すべきかを示す

## 6. テスト設計

### 6.1 単体テスト

#### 6.1.1 `getDatabaseConfig()`メソッドのテスト

**テストケース**:
1. PostgreSQL用のDSNが正しく生成されること
2. MySQL用のDSNが正しく生成されること
3. ドライバーが指定されていない場合にエラーが発生すること
4. 未対応のドライバーが指定された場合にエラーが発生すること
5. 接続情報が不完全な場合にエラーが発生すること（既存のテスト）

**テストファイル**: `server/internal/admin/config_test.go`

**実装例**:

```go
func TestConfig_getDatabaseConfig_PostgreSQL(t *testing.T) {
	// PostgreSQL用の設定を作成
	cfg := &Config{
		appConfig: &config.Config{
			Database: config.DatabaseConfig{
				Groups: config.DatabaseGroupsConfig{
					Master: []config.ShardConfig{
						{
							Driver:   "postgres",
							Host:     "localhost",
							Port:     5432,
							User:     "webdb",
							Password: "webdb",
							Name:     "webdb_master",
						},
					},
				},
			},
		},
	}

	dbConfig := cfg.getDatabaseConfig()
	
	assert.Equal(t, "postgresql", dbConfig["default"].Driver)
	assert.Contains(t, dbConfig["default"].Dsn, "host=localhost")
	assert.Contains(t, dbConfig["default"].Dsn, "port=5432")
	assert.Contains(t, dbConfig["default"].Dsn, "user=webdb")
	assert.Contains(t, dbConfig["default"].Dsn, "dbname=webdb_master")
	assert.Contains(t, dbConfig["default"].Dsn, "sslmode=disable")
}

func TestConfig_getDatabaseConfig_MySQL(t *testing.T) {
	// MySQL用の設定を作成
	cfg := &Config{
		appConfig: &config.Config{
			Database: config.DatabaseConfig{
				Groups: config.DatabaseGroupsConfig{
					Master: []config.ShardConfig{
						{
							Driver:   "mysql",
							Host:     "localhost",
							Port:     3306,
							User:     "webdb",
							Password: "webdb",
							Name:     "webdb_master",
						},
					},
				},
			},
		},
	}

	dbConfig := cfg.getDatabaseConfig()
	
	assert.Equal(t, "mysql", dbConfig["default"].Driver)
	assert.Contains(t, dbConfig["default"].Dsn, "webdb:webdb@tcp(localhost:3306)/webdb_master")
	assert.Contains(t, dbConfig["default"].Dsn, "charset=utf8mb4")
	assert.Contains(t, dbConfig["default"].Dsn, "parseTime=true")
	assert.Contains(t, dbConfig["default"].Dsn, "loc=Local")
}

func TestConfig_getDatabaseConfig_NoDriver(t *testing.T) {
	// ドライバーが指定されていない設定を作成
	cfg := &Config{
		appConfig: &config.Config{
			Database: config.DatabaseConfig{
				Groups: config.DatabaseGroupsConfig{
					Master: []config.ShardConfig{
						{
							Driver:   "",  // ドライバー未指定
							Host:     "localhost",
							Port:     5432,
							User:     "webdb",
							Password: "webdb",
							Name:     "webdb_master",
						},
					},
				},
			},
		},
	}

	assert.Panics(t, func() {
		cfg.getDatabaseConfig()
	}, "should panic when driver is not specified")
}

func TestConfig_getDatabaseConfig_UnsupportedDriver(t *testing.T) {
	// 未対応のドライバーが指定された設定を作成
	cfg := &Config{
		appConfig: &config.Config{
			Database: config.DatabaseConfig{
				Groups: config.DatabaseGroupsConfig{
					Master: []config.ShardConfig{
						{
							Driver:   "sqlite",  // 未対応のドライバー
							Host:     "localhost",
							Port:     5432,
							User:     "webdb",
							Password: "webdb",
							Name:     "webdb_master",
						},
					},
				},
			},
		},
	}

	assert.Panics(t, func() {
		cfg.getDatabaseConfig()
	}, "should panic when unsupported driver is specified")
}
```

### 6.2 統合テスト

#### 6.2.1 Adminサーバーの起動テスト

**テストケース**:
1. PostgreSQL設定でAdminサーバーが正常に起動できること
2. MySQL設定でAdminサーバーが正常に起動できること
3. 管理画面が正常に表示されること

**テストファイル**: `server/cmd/admin/main_test.go`（既存のテストを拡張）

## 7. 実装上の注意事項

### 7.1 コードの一貫性

- **命名規則**: 既存のコードスタイルに従う
- **エラーハンドリング**: 既存のエラーハンドリングパターンに従う
- **コメント**: 既存のコメントスタイルに従う

### 7.2 設定ファイルの整合性

- **ドライバーの一致**: `config.yaml`の`DB_TYPE`と`database.yaml`（または`database.mysql.yaml`）の`driver`が一致している必要がある
- **設定ファイルの読み込み**: メインサーバーと同じ設定ファイル読み込みロジックを使用（既存の実装）

### 7.3 後方互換性

- **既存のPostgreSQL設定**: 既存のPostgreSQL設定ファイルは引き続き動作する
- **既存のコード**: 既存のPostgreSQLコードは変更しない（追加のみ）

### 7.4 パフォーマンス

- **接続プール**: GoAdminが管理する接続プールを使用（既存の実装）
- **DSN生成**: 起動時に1回のみ実行されるため、パフォーマンスへの影響は最小限

## 8. 参考情報

### 8.1 既存実装の参考

- `server/internal/config/config.go`: メインサーバーの`GetDSN()`メソッド（PostgreSQL/MySQL対応）
- `server/internal/admin/config.go`: 既存のPostgreSQL用設定
- `server/cmd/admin/main.go`: 既存のエントリーポイント

### 8.2 GoAdminドキュメント

- PostgreSQLドライバー: `github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres`
- MySQLドライバー: `github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql`

### 8.3 関連ドキュメント

- `.kiro/specs/0055-admin-mysql/requirements.md`: 要件定義書
- `.kiro/specs/0054-mysql/design.md`: メインサーバーのMySQL対応設計書
