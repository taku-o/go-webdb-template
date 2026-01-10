# AdminサーバーのMySQL対応の実装タスク一覧

## 概要
Adminサーバー（GoAdmin管理画面）をPostgreSQLに加えて、MySQLでも動作するように修正するためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: データベース接続設定の修正

#### タスク 1.1: getDatabaseConfig()メソッドの修正
**目的**: データベースタイプに応じて適切なDSN形式を生成するように`getDatabaseConfig()`メソッドを修正する。

**作業内容**:
- `server/internal/admin/config.go`の`getDatabaseConfig()`メソッドを修正
- `c.appConfig.Database.Groups.Master[0].Driver`を参照してデータベースタイプを判定
- PostgreSQL用とMySQL用のDSN形式をそれぞれ生成
- ドライバー名もデータベースタイプに応じて設定
- エラーハンドリングを追加（ドライバー未指定、未対応ドライバー）

**実装内容**:
- 修正対象: `server/internal/admin/config.go`の`getDatabaseConfig()`メソッド（53-83行目）
- 修正前:
  ```go
  // PostgreSQL用のDSN形式を構築
  // host=... port=... user=... password=... dbname=... sslmode=disable
  dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
      masterDB.Host,
      masterDB.Port,
      masterDB.User,
      masterDB.Password,
      masterDB.Name,
  )

  return goadminConfig.DatabaseList{
      "default": {
          Driver: "postgresql",
          Dsn:    dsn,
      },
  }
  ```
- 修正後:
  ```go
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
  ```

**受け入れ基準**:
- [ ] `getDatabaseConfig()`メソッドで`masterDB.Driver`を参照してデータベースタイプを判定している
- [ ] PostgreSQLの場合、PostgreSQL用のDSN形式を生成している（`host=... port=... user=... password=... dbname=... sslmode=disable`）
- [ ] MySQLの場合、MySQL用のDSN形式を生成している（`user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=true&loc=Local`）
- [ ] ドライバー名がデータベースタイプに応じて設定されている（`postgresql`または`mysql`）
- [ ] ドライバーが指定されていない場合にエラーを返している
- [ ] 未対応のドライバーが指定された場合にエラーを返している
- [ ] 既存のPostgreSQLコードに影響がない（後方互換性を維持）

- _Requirements: 3.1.1, 3.1.2, 6.1_
- _Design: 3.1.1_

---

### Phase 2: MySQLドライバーのインポート追加

#### タスク 2.1: main.goにMySQLドライバーのインポートを追加
**目的**: GoAdminでMySQLを使用できるようにするため、MySQLドライバーをインポートする。

**作業内容**:
- `server/cmd/admin/main.go`にMySQLドライバーのインポートを追加
- 既存のPostgreSQLドライバーのインポートは維持

**実装内容**:
- 修正対象: `server/cmd/admin/main.go`のインポートセクション（13行目付近）
- 修正前:
  ```go
  _ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"
  _ "github.com/GoAdminGroup/themes/adminlte"
  ```
- 修正後:
  ```go
  _ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"
  _ "github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql"
  _ "github.com/GoAdminGroup/themes/adminlte"
  ```

**受け入れ基準**:
- [ ] `server/cmd/admin/main.go`にMySQLドライバーのインポートが追加されている
- [ ] 既存のPostgreSQLドライバーのインポートが維持されている
- [ ] ブランクインポート（`_`）を使用している

- _Requirements: 3.2.1, 6.2_
- _Design: 3.2.1_

---

### Phase 3: テストコードの追加

#### タスク 3.1: getDatabaseConfig()メソッドの単体テストを追加
**目的**: `getDatabaseConfig()`メソッドの動作を検証するため、単体テストを追加する。

**作業内容**:
- `server/internal/admin/config_test.go`にテストケースを追加
- PostgreSQL用のDSN生成テスト
- MySQL用のDSN生成テスト
- エラーハンドリングのテスト（ドライバー未指定、未対応ドライバー）

**実装内容**:
- テストファイル: `server/internal/admin/config_test.go`（新規作成または既存ファイルに追加）
- テストケース:
  1. `TestConfig_getDatabaseConfig_PostgreSQL`: PostgreSQL用のDSNが正しく生成されること
  2. `TestConfig_getDatabaseConfig_MySQL`: MySQL用のDSNが正しく生成されること
  3. `TestConfig_getDatabaseConfig_NoDriver`: ドライバーが指定されていない場合にエラーが発生すること
  4. `TestConfig_getDatabaseConfig_UnsupportedDriver`: 未対応のドライバーが指定された場合にエラーが発生すること

**実装例**:

```go
package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taku-o/go-webdb-template/internal/config"
)

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

**受け入れ基準**:
- [ ] `TestConfig_getDatabaseConfig_PostgreSQL`テストが追加されている
- [ ] `TestConfig_getDatabaseConfig_MySQL`テストが追加されている
- [ ] `TestConfig_getDatabaseConfig_NoDriver`テストが追加されている
- [ ] `TestConfig_getDatabaseConfig_UnsupportedDriver`テストが追加されている
- [ ] すべてのテストが正常に実行できる
- [ ] テストカバレッジが適切である

- _Requirements: 6.1_
- _Design: 6.1.1_

---

### Phase 4: 動作確認

#### タスク 4.1: PostgreSQL環境での動作確認
**目的**: 既存のPostgreSQL機能が正常に動作することを確認する。

**作業内容**:
- PostgreSQL設定でAdminサーバーを起動
- 管理画面が正常に表示されることを確認
- 既存の機能が正常に動作することを確認

**受け入れ基準**:
- [ ] PostgreSQL設定でAdminサーバーが正常に起動できる
- [ ] PostgreSQL環境で管理画面が正常に表示される
- [ ] 既存のPostgreSQL機能が正常に動作する

- _Requirements: 6.3_

---

#### タスク 4.2: MySQL環境での動作確認
**目的**: MySQL環境でAdminサーバーが正常に動作することを確認する。

**作業内容**:
- MySQL設定でAdminサーバーを起動
- 管理画面が正常に表示されることを確認
- MySQL環境で正常に動作することを確認

**前提条件**:
- MySQLデータベースが起動していること
- MySQL用の設定ファイル（`database.mysql.yaml`）が正しく設定されていること
- `config.yaml`の`DB_TYPE`が`mysql`に設定されていること

**受け入れ基準**:
- [ ] MySQL設定でAdminサーバーが正常に起動できる
- [ ] MySQL環境で管理画面が正常に表示される
- [ ] MySQL環境で正常に動作する

- _Requirements: 6.3_

---

## タスクの依存関係

```
Phase 1: データベース接続設定の修正
  └─> Phase 2: MySQLドライバーのインポート追加
      └─> Phase 3: テストコードの追加
          └─> Phase 4: 動作確認
```

## 実装順序

1. **Phase 1**: `getDatabaseConfig()`メソッドの修正（データベースタイプ判定とDSN生成）
2. **Phase 2**: MySQLドライバーのインポート追加
3. **Phase 3**: テストコードの追加
4. **Phase 4**: 動作確認（PostgreSQL環境、MySQL環境）

## 注意事項

### 実装時の注意点

1. **既存コードの保護**: 既存のPostgreSQLコードは変更せず、MySQL対応を追加する
2. **エラーハンドリング**: 設定ファイルの不備は明確なエラーメッセージで報告する
3. **テストカバレッジ**: すべての分岐をテストでカバーする
4. **後方互換性**: 既存のPostgreSQL設定ファイルは引き続き動作することを確認する

### 設定ファイルの整合性

- `config.yaml`の`DB_TYPE`と`database.yaml`（または`database.mysql.yaml`）の`driver`が一致している必要がある
- メインサーバーと同じ設定ファイルを使用する想定

### 参考実装

- `server/internal/config/config.go`の`GetDSN()`メソッド: メインサーバーのDSN生成ロジック（PostgreSQL/MySQL対応）
- `server/internal/admin/config.go`: 既存のPostgreSQL用設定
