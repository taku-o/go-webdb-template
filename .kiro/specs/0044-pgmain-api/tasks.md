# APIサーバーPostgreSQL対応実装タスク一覧

## 概要
APIサーバーでPostgreSQLを利用するように修正を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: 設定ファイルの修正

#### タスク 1.1: config/develop/database.yamlの修正
**目的**: SQLite設定を削除し、PostgreSQL設定を有効化。論理シャーディング数8を定義。

**作業内容**:
- `config/develop/database.yaml`を開く
- SQLite設定（`database.groups.master`と`database.groups.sharding`のSQLite設定）を完全に削除
- PostgreSQL設定のコメントアウトを解除（または新規追加）
- masterグループのPostgreSQL設定を追加:
  - `id: 1`
  - `driver: postgres`
  - `host: localhost`
  - `port: 5432`
  - `user: webdb`
  - `password: webdb`
  - `name: webdb_master`
  - `max_connections: 25`
  - `max_idle_connections: 5`
  - `connection_max_lifetime: 1h`
- shardingグループに8つのデータベース設定（id: 1-8）を追加:
  - id: 1, table_range: [0, 3], host: localhost, port: 5433, name: webdb_sharding_1
  - id: 2, table_range: [4, 7], host: localhost, port: 5433, name: webdb_sharding_1
  - id: 3, table_range: [8, 11], host: localhost, port: 5434, name: webdb_sharding_2
  - id: 4, table_range: [12, 15], host: localhost, port: 5434, name: webdb_sharding_2
  - id: 5, table_range: [16, 19], host: localhost, port: 5435, name: webdb_sharding_3
  - id: 6, table_range: [20, 23], host: localhost, port: 5435, name: webdb_sharding_3
  - id: 7, table_range: [24, 27], host: localhost, port: 5436, name: webdb_sharding_4
  - id: 8, table_range: [28, 31], host: localhost, port: 5436, name: webdb_sharding_4
- 各shardingデータベース設定に以下を追加:
  - `driver: postgres`
  - `user: webdb`
  - `password: webdb`
  - `max_connections: 25`
  - `max_idle_connections: 5`
  - `connection_max_lifetime: 1h`
  - `sslmode: disable`
- tables定義を追加:
  - `name: dm_users`, `suffix_count: 32`
  - `name: dm_posts`, `suffix_count: 32`

**受け入れ基準**:
- SQLite設定が完全に削除されている（コメントアウトではない）
- PostgreSQL設定が有効になっている
- masterグループの設定が正しく定義されている
- shardingグループに8つのデータベース設定（id: 1-8）が定義されている
- 各論理シャードのtable_rangeが正しく設定されている（[0,3], [4,7], [8,11], [12,15], [16,19], [20,23], [24,27], [28,31]）
- 各論理シャードが正しい物理DB（host/port/name）を参照している

- _Requirements: 3.1.1_

---

#### タスク 1.2: config/staging/database.yamlの確認・修正
**目的**: PostgreSQL設定を確認し、論理シャーディング数8を確認。SQLite設定があれば削除。

**作業内容**:
- `config/staging/database.yaml`を開く
- PostgreSQL設定が正しく定義されているか確認
- シャーディング構成が正しいか確認（物理DB 4台、論理シャーディング8）
- shardingグループに8つのデータベース設定（id: 1-8）が定義されているか確認
- 各論理シャードのtable_rangeが正しく設定されているか確認
- SQLite設定が存在する場合は削除
- 不備があれば修正

**受け入れ基準**:
- PostgreSQL設定が正しく定義されている
- シャーディング構成が正しい（物理DB 4台、論理シャーディング8）
- shardingグループに8つのデータベース設定（id: 1-8）が定義されている
- 各論理シャードのtable_rangeが正しく設定されている
- SQLite設定が削除されている（存在する場合）

- _Requirements: 3.1.2_

---

#### タスク 1.3: config/production/database.yamlの確認・修正
**目的**: PostgreSQL設定を確認し、論理シャーディング数8を確認。SQLite設定があれば削除（存在する場合）。

**作業内容**:
- `config/production/database.yaml`が存在するか確認
- 存在する場合:
  - ファイルを開く
  - PostgreSQL設定が正しく定義されているか確認
  - シャーディング構成が正しいか確認（物理DB 4台、論理シャーディング8）
  - shardingグループに8つのデータベース設定（id: 1-8）が定義されているか確認
  - 各論理シャードのtable_rangeが正しく設定されているか確認
  - SQLite設定が存在する場合は削除
  - 不備があれば修正

**受け入れ基準**:
- ファイルが存在しない場合はスキップ
- ファイルが存在する場合:
  - PostgreSQL設定が正しく定義されている
  - シャーディング構成が正しい（物理DB 4台、論理シャーディング8）
  - shardingグループに8つのデータベース設定（id: 1-8）が定義されている
  - 各論理シャードのtable_rangeが正しく設定されている
  - SQLite設定が削除されている（存在する場合）

- _Requirements: 3.1.3_

---

#### タスク 1.4: config/production/database.yaml.exampleの確認・修正
**目的**: PostgreSQL設定を確認し、論理シャーディング数8を確認。SQLite設定があれば削除。

**作業内容**:
- `config/production/database.yaml.example`を開く
- PostgreSQL設定が正しく定義されているか確認
- シャーディング構成が正しいか確認（物理DB 4台、論理シャーディング8）
- shardingグループに8つのデータベース設定（id: 1-8）が定義されているか確認
- 各論理シャードのtable_rangeが正しく設定されているか確認
- SQLite設定が存在する場合は削除
- 不備があれば修正
- テーブル名が`dm_users`と`dm_posts`になっているか確認（`users`と`posts`ではない）

**受け入れ基準**:
- PostgreSQL設定が正しく定義されている
- シャーディング構成が正しい（物理DB 4台、論理シャーディング8）
- shardingグループに8つのデータベース設定（id: 1-8）が定義されている
- 各論理シャードのtable_rangeが正しく設定されている
- SQLite設定が削除されている（存在する場合）
- テーブル名が`dm_users`と`dm_posts`になっている

- _Requirements: 3.1.4_

---

### Phase 2: SQLite用ライブラリと処理分岐の削除

#### タスク 2.1: server/internal/db/connection.goのインポート削除
**目的**: SQLite用ライブラリのインポートを削除。

**作業内容**:
- `server/internal/db/connection.go`を開く
- インポートセクションから以下を削除:
  - `_ "github.com/mattn/go-sqlite3"`
  - `"gorm.io/driver/sqlite"`

**受け入れ基準**:
- `_ "github.com/mattn/go-sqlite3"`のインポートが削除されている
- `"gorm.io/driver/sqlite"`のインポートが削除されている
- 未使用のインポートが残っていない

- _Requirements: 3.2.1_

---

#### タスク 2.2: server/internal/db/connection.goの処理分岐削除
**目的**: SQLite用処理分岐を削除。

**作業内容**:
- `server/internal/db/connection.go`を開く
- `NewConnection`関数内の`driver = "sqlite3"`デフォルト値設定を削除
- `createGORMConnection`関数内の`case "sqlite3":`分岐を削除
- `createGORMConnectionFromDSN`関数内の`case "sqlite3":`分岐を削除
- `NewGORMConnection`関数内の`case "sqlite3":`分岐を削除（Reader接続作成部分）

**受け入れ基準**:
- `NewConnection`関数内の`driver = "sqlite3"`デフォルト値設定が削除されている
- `createGORMConnection`関数内の`case "sqlite3":`分岐が削除されている
- `createGORMConnectionFromDSN`関数内の`case "sqlite3":`分岐が削除されている
- `NewGORMConnection`関数内の`case "sqlite3":`分岐が削除されている
- 未サポートドライバー指定時に適切なエラーが返されることを確認

- _Requirements: 3.2.2_

---

#### タスク 2.3: server/cmd/admin/main.goのインポート削除
**目的**: SQLite用ライブラリのインポートを削除。

**作業内容**:
- `server/cmd/admin/main.go`を開く
- インポートセクションから以下を削除:
  - `_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"`

**受け入れ基準**:
- `_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"`のインポートが削除されている
- 未使用のインポートが残っていない

- _Requirements: 3.2.1_

---

#### タスク 2.4: server/go.modの依存関係削除
**目的**: SQLite関連の依存関係を削除。

**作業内容**:
- `server/go.mod`を開く
- `gorm.io/driver/sqlite v1.5.6`の行を削除
- `cd server && go mod tidy`を実行して依存関係を整理

**受け入れ基準**:
- `gorm.io/driver/sqlite`の依存関係が削除されている
- `go mod tidy`が正常に実行される
- 未使用の依存関係が残っていない

- _Requirements: 3.2.1_

---

#### タスク 2.5: エラーハンドリングの確認
**目的**: 未サポートドライバー指定時のエラーハンドリングが適切か確認。

**作業内容**:
- `server/internal/db/connection.go`を確認
- 未サポートドライバー（`sqlite3`など）を指定した場合のエラーハンドリングを確認
- エラーメッセージが適切か確認

**受け入れ基準**:
- 未サポートドライバー指定時に適切なエラーが返される
- エラーメッセージが明確で分かりやすい

- _Requirements: 3.2.3_

---

### Phase 3: テストコードの修正

#### タスク 3.1: テストコードのSQLite依存確認
**目的**: テストコードのSQLite依存を確認。

**作業内容**:
- `server/test/`配下のテストファイルを確認
- SQLite固有の設定やコードが含まれているか確認
- テストデータベースの初期化方法を確認
- テスト実行時のデータベース接続方法を確認
- 確認結果を記録

**受け入れ基準**:
- すべてのテストファイルのSQLite依存が確認されている
- SQLite固有の設定やコードがリストアップされている
- テストデータベースの初期化方法が確認されている
- テスト実行時のデータベース接続方法が確認されている

- _Requirements: 3.3.1_

---

#### タスク 3.2: server/test/testutil/db.goのPostgreSQL対応
**目的**: テストユーティリティをPostgreSQL対応に修正。

**作業内容**:
- `server/test/testutil/db.go`を開く
- SQLite固有の設定をPostgreSQL設定に変更
- テストデータベースの初期化方法をPostgreSQL対応に変更
- テスト実行時のデータベース接続方法をPostgreSQL対応に変更
- 論理シャーディング数8を考慮したテスト設定

**受け入れ基準**:
- SQLite固有の設定がPostgreSQL設定に変更されている
- テストデータベースの初期化方法がPostgreSQL対応になっている
- テスト実行時のデータベース接続方法がPostgreSQL対応になっている
- 論理シャーディング数8を考慮したテスト設定になっている

- _Requirements: 3.3.2_

---

#### タスク 3.3: その他のテストファイルのPostgreSQL対応
**目的**: その他のテストファイルのSQLite依存を修正。

**作業内容**:
- タスク3.1で確認したSQLite依存があるテストファイルを修正
- SQLite固有の設定をPostgreSQL設定に変更
- テストデータベースの初期化方法をPostgreSQL対応に変更
- テスト実行時のデータベース接続方法をPostgreSQL対応に変更

**受け入れ基準**:
- すべてのSQLite依存があるテストファイルが修正されている
- SQLite固有の設定がPostgreSQL設定に変更されている
- テストデータベースの初期化方法がPostgreSQL対応になっている
- テスト実行時のデータベース接続方法がPostgreSQL対応になっている

- _Requirements: 3.3.2_

---

### Phase 4: ドキュメントの更新

#### タスク 4.1: README.mdの更新
**目的**: PostgreSQL利用に関する記述を追加。

**作業内容**:
- `README.md`を開く
- PostgreSQL利用に関する記述を追加
- 設定ファイルの変更方法を記載
- 開発環境でのPostgreSQL起動手順を記載:
  - PostgreSQLコンテナの起動方法（`./scripts/start-postgres.sh start`）
  - マイグレーションの適用方法（`./scripts/migrate.sh`）
  - APIサーバーの起動方法
- 論理シャーディング数8の説明を追加
- SQLiteに関する記述があればPostgreSQLに変更

**受け入れ基準**:
- PostgreSQL利用に関する記述が追加されている
- 設定ファイルの変更方法が記載されている
- 開発環境でのPostgreSQL起動手順が記載されている
- 論理シャーディング数8の説明が追加されている
- SQLiteに関する記述がPostgreSQLに変更されている（存在する場合）

- _Requirements: 3.4.1_

---

#### タスク 4.2: docs/配下の関連ドキュメントの更新
**目的**: SQLiteに関する記述をPostgreSQLに変更。

**作業内容**:
- `docs/`配下の関連ドキュメントを確認
- SQLiteに関する記述をPostgreSQLに変更
- PostgreSQL利用に関する記述を追加（必要に応じて）

**受け入れ基準**:
- すべての関連ドキュメントのSQLiteに関する記述がPostgreSQLに変更されている
- PostgreSQL利用に関する記述が追加されている（必要に応じて）

- _Requirements: 3.4.2_

---

### Phase 5: 動作確認

#### タスク 5.1: PostgreSQL接続の確認
**目的**: APIサーバーがPostgreSQLに正常に接続できることを確認。

**作業内容**:
- PostgreSQLコンテナが起動していることを確認
- マイグレーションが適用されていることを確認
- APIサーバーを起動
- APIサーバーがPostgreSQLに正常に接続できることを確認
- masterデータベースへの接続を確認
- shardingデータベースへの接続を確認（8つの論理シャードすべて）

**受け入れ基準**:
- PostgreSQLコンテナが起動している
- マイグレーションが適用されている
- APIサーバーがPostgreSQLに正常に接続できる
- masterデータベースへの接続が正常に確立できる
- shardingデータベースへの接続が正常に確立できる（8つの論理シャードすべて）

- _Requirements: 6.4_

---

#### タスク 5.2: APIサーバーの動作確認
**目的**: APIサーバーが正常に起動し、リクエストを処理できることを確認。

**作業内容**:
- APIサーバーを起動
- APIサーバーが正常に起動することを確認
- リクエストを送信して正常に処理されることを確認
- エラーログがないことを確認

**受け入れ基準**:
- APIサーバーが正常に起動する
- リクエストが正常に処理される
- エラーログがない

- _Requirements: 6.4_

---

#### タスク 5.3: テストの実行
**目的**: テストがPostgreSQL環境で正常に実行できることを確認。

**作業内容**:
- PostgreSQLコンテナが起動していることを確認
- マイグレーションが適用されていることを確認
- テストを実行
- すべてのテストが正常に実行されることを確認

**受け入れ基準**:
- PostgreSQLコンテナが起動している
- マイグレーションが適用されている
- すべてのテストが正常に実行される
- テストエラーがない

- _Requirements: 6.3_

---

## タスクの依存関係

### Phase 1 → Phase 2
- Phase 1（設定ファイルの修正）が完了してからPhase 2（SQLite用ライブラリと処理分岐の削除）を開始

### Phase 2 → Phase 3
- Phase 2（SQLite用ライブラリと処理分岐の削除）が完了してからPhase 3（テストコードの修正）を開始

### Phase 3 → Phase 4
- Phase 3（テストコードの修正）が完了してからPhase 4（ドキュメントの更新）を開始

### Phase 4 → Phase 5
- Phase 4（ドキュメントの更新）が完了してからPhase 5（動作確認）を開始

## 注意事項

### 実装時の注意事項
1. **設定ファイルの修正**: SQLite設定は完全に削除（コメントアウトではない）
2. **論理シャーディング数8**: 必ず8つのデータベース設定（id: 1-8）を定義
3. **物理DBと論理シャードの対応**: 同じ物理DBを参照する論理シャードは同じhost/port/nameを使用
4. **依存関係の削除**: `go mod tidy`を実行して依存関係を整理
5. **エラーハンドリング**: 未サポートドライバー指定時のエラーハンドリングを確認

### 動作確認時の注意事項
1. **PostgreSQLコンテナの起動**: マイグレーション適用前にPostgreSQLコンテナを起動
2. **マイグレーションの適用**: APIサーバー起動前にマイグレーションを適用
3. **接続確認**: すべてのデータベース（master 1台、sharding 8つの論理シャード）への接続を確認

## 参考情報

### 関連Issue
- GitHub Issue #85: 開発環境はPostgreSQLを利用する前提とする
- GitHub Issue #86: PostgreSQLの起動スクリプトと、Atlasマイグレーションスクリプトの修正
- GitHub Issue #87: APIサーバーの修正

### 関連ドキュメント
- 要件定義書: `requirements.md`
- 設計書: `design.md`
