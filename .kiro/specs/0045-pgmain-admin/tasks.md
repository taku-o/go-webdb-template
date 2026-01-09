# GoAdminサーバーPostgreSQL対応実装タスク一覧

## 概要
GoAdminサーバーでPostgreSQLを利用するように修正を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: GoAdmin設定の修正

#### タスク 1.1: server/internal/admin/config.goの修正
**目的**: `getDatabaseConfig()`関数をPostgreSQL対応に修正。SQLite設定を削除し、PostgreSQL設定を追加。

**作業内容**:
- `server/internal/admin/config.go`を開く
- `getDatabaseConfig()`関数を修正:
  - SQLite設定（`Driver: "sqlite"`, `File: dsn`）を削除
  - PostgreSQL設定（`Driver: "postgres"`, `Dsn: dsn`）を追加
  - `c.appConfig.Database.Groups.Master[0]`から接続情報を取得
  - masterグループのデータベース設定が存在しない場合のエラーハンドリングを追加
  - 接続情報が不完全な場合のエラーハンドリングを追加
  - PostgreSQL用のDSN形式を構築（`host=... port=... user=... password=... dbname=... sslmode=disable`）
- 必要なインポートを追加:
  - `fmt`パッケージ（DSN構築のため）
  - `os`パッケージ（環境変数の読み込みのため）

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
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		masterDB.Host,
		masterDB.Port,
		masterDB.User,
		masterDB.Password,
		masterDB.Name,
	)

	return goadminConfig.DatabaseList{
		"default": {
			Driver: "postgres",
			Dsn:    dsn,
		},
	}
}
```

**受け入れ基準**:
- SQLite設定（`Driver: "sqlite"`, `File: dsn`）が削除されている
- PostgreSQL設定（`Driver: "postgres"`, `Dsn: dsn`）が追加されている
- `c.appConfig.Database.Groups.Master[0]`から接続情報を取得している
- masterグループのデータベース設定が存在しない場合のエラーハンドリングが実装されている
- 接続情報が不完全な場合のエラーハンドリングが実装されている
- PostgreSQL用のDSN形式が正しく構築されている
- `fmt`パッケージと`os`パッケージがインポートされている

- _Requirements: 3.1.1, 3.1.2_
- _Design: 3.1.1_

---

### Phase 2: SQLite用ライブラリのインポート確認

#### タスク 2.1: server/cmd/admin/main.goのインポート確認
**目的**: SQLite用ライブラリのインポートが存在しないことを確認。PostgreSQLドライバーがインポートされていることを確認。

**作業内容**:
- `server/cmd/admin/main.go`を開く
- インポートセクションを確認:
  - `_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"`のインポートが存在しないことを確認
  - `_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"`のインポートが存在することを確認
- 必要に応じて`grep`コマンドで確認:
  - `grep -n "sqlite" server/cmd/admin/main.go`でSQLite関連のインポートを検索
  - `grep -n "postgres" server/cmd/admin/main.go`でPostgreSQLドライバーのインポートを検索

**受け入れ基準**:
- SQLite用ライブラリ（`_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"`）のインポートが存在しない
- PostgreSQLドライバー（`_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"`）のインポートが存在する

- _Requirements: 3.2.1_
- _Design: 3.1.2_

---

### Phase 3: ドキュメントの更新

#### タスク 3.1: README.mdの更新
**目的**: GoAdminサーバーのPostgreSQL利用に関する記述を追加。

**作業内容**:
- `README.md`を開く
- GoAdminサーバーのPostgreSQL利用に関するセクションを追加:
  - GoAdminサーバーがPostgreSQLを利用することを明記
  - 設定ファイル（`config/{env}/database.yaml`）から接続情報を読み込むことを明記
  - masterデータベースのみを使用することを明記
- 起動手順を記載:
  - PostgreSQLコンテナの起動手順（`./scripts/start-postgres.sh start`）
  - マイグレーション適用手順（`./scripts/migrate.sh`）
  - GoAdminサーバー起動手順（`cd server && go run cmd/admin/main.go`）
- 設定ファイルの変更方法を記載:
  - `config/{env}/database.yaml`の設定方法
  - PostgreSQL接続情報の設定方法

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

**受け入れ基準**:
- GoAdminサーバーのPostgreSQL利用に関する記述が追加されている
- 起動手順が記載されている
- 設定ファイルの変更方法が記載されている

- _Requirements: 3.3.1_
- _Design: 3.2.1_

---

#### タスク 3.2: その他のドキュメントの更新
**目的**: GoAdminサーバーのSQLiteに関する記述をPostgreSQLに変更。

**作業内容**:
- `docs/`配下の関連ドキュメントを確認:
  - `docs/Architecture.md`
  - `docs/Initial-Setup.md`
  - その他の関連ドキュメント
- GoAdminサーバーのSQLiteに関する記述をPostgreSQLに変更
- GoAdminサーバーのPostgreSQL利用に関する記述を追加（必要に応じて）

**受け入れ基準**:
- GoAdminサーバーのSQLiteに関する記述がPostgreSQLに変更されている
- GoAdminサーバーのPostgreSQL利用に関する記述が追加されている（必要に応じて）

- _Requirements: 3.3.2_
- _Design: 3.2.2_

---

### Phase 4: 動作確認

#### タスク 4.1: GoAdminサーバーの起動確認
**目的**: GoAdminサーバーがPostgreSQLに正常に接続できることを確認。

**作業内容**:
- PostgreSQLコンテナが起動していることを確認:
  - `./scripts/start-postgres.sh status`で確認
  - または`docker ps`で確認
- マイグレーションが適用されていることを確認:
  - `./scripts/migrate.sh`を実行（必要に応じて）
- GoAdminサーバーを起動:
  - `cd server && go run cmd/admin/main.go`
- 起動ログを確認:
  - PostgreSQLへの接続エラーが発生していないことを確認
  - 正常に起動していることを確認
- 管理画面にアクセス:
  - `http://localhost:{port}/admin`にアクセス（portは設定ファイルで指定されたポート）
  - ログイン画面が表示されることを確認

**受け入れ基準**:
- GoAdminサーバーが正常に起動する
- PostgreSQLへの接続エラーが発生しない
- 管理画面にアクセスできる

- _Requirements: 6.3_

---

#### タスク 4.2: GoAdmin管理画面の動作確認
**目的**: GoAdmin管理画面でデータベースのテーブルが正常に表示できることを確認。

**作業内容**:
- GoAdmin管理画面にログイン:
  - 設定ファイルで指定された認証情報でログイン
- データベースのテーブル一覧を確認:
  - テーブル一覧が正常に表示されることを確認
  - masterデータベース（`webdb_master`）のテーブルが表示されることを確認
- テーブルのデータを確認:
  - 任意のテーブルを選択してデータが正常に表示されることを確認

**受け入れ基準**:
- GoAdmin管理画面にログインできる
- データベースのテーブル一覧が正常に表示される
- テーブルのデータが正常に表示される

- _Requirements: 6.3_

---

### Phase 5: テスト

#### タスク 5.1: 単体テストの作成
**目的**: `server/internal/admin/config.go`の`getDatabaseConfig()`関数の単体テストを作成。

**作業内容**:
- `server/internal/admin/config_test.go`を開く（存在しない場合は新規作成）
- テストケースを追加:
  1. **正常系**: masterグループのデータベース設定が存在する場合、PostgreSQL設定が正しく構築されることを確認
  2. **異常系**: masterグループのデータベース設定が存在しない場合、エラーが発生することを確認
  3. **異常系**: 接続情報が不完全な場合、エラーが発生することを確認
  4. **正常系**: PostgreSQL用のDSN形式が正しく構築されることを確認
- テストを実行:
  - `cd server && go test ./internal/admin/... -v`

**受け入れ基準**:
- すべてのテストケースが実装されている
- すべてのテストが成功する

- _Requirements: 6.1_
- _Design: 6.1.1_

---

## タスク実行順序

1. **Phase 1**: GoAdmin設定の修正（タスク 1.1）
2. **Phase 2**: SQLite用ライブラリのインポート確認（タスク 2.1）
3. **Phase 3**: ドキュメントの更新（タスク 3.1, 3.2）
4. **Phase 4**: 動作確認（タスク 4.1, 4.2）
5. **Phase 5**: テスト（タスク 5.1）

## 注意事項

- 各タスクの実行前に、関連する要件定義書と設計書を確認すること
- タスク実行後は、受け入れ基準を確認すること
- エラーが発生した場合は、エラーメッセージを確認し、適切に対処すること
- 動作確認は、PostgreSQLコンテナが起動している状態で実施すること
